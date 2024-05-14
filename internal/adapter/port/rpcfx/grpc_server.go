package rpcfx

import (
	"context"
	"encoding/json"
	"fmt"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/structx/ddns/internal/core/domain"
	pbv1 "github.com/structx/ddns/proto/ddns/v1"
)

// GRPCServer
type GRPCServer struct {
	pbv1.UnimplementedDDNSServiceV1Server

	log     *zap.SugaredLogger
	service domain.DDNS
	srv     *grpc.Server
}

// NewGRPCServer
func NewGRPCServer(logger *zap.Logger, ddns domain.DDNS) *GRPCServer {
	return &GRPCServer{
		log:     logger.Sugar().Named("GrpcServer"),
		service: ddns,
		srv:     nil,
	}
}

// Ping
func (g *GRPCServer) Ping(ctx context.Context, in *pbv1.PingRequest) (*pbv1.PingResponse, error) {

	g.log.Debugw("Ping", "request", in)

	echo := g.service.Echo(ctx)
	return &pbv1.PingResponse{
		Echo: &pbv1.Echo{
			Ip:          echo.IP,
			Port:        echo.Port,
			NodeId:      echo.NodeID,
			CompletedAt: timestamppb.Now(),
		},
	}, nil
}

// Store
func (g *GRPCServer) Store(ctx context.Context, in *pbv1.StoreRequest) (*pbv1.StoreResponse, error) {

	g.log.Debugw("Store", "request", in)

	var record domain.Record
	err := json.Unmarshal(in.GetValue(), &record)
	if err != nil {
		g.log.Errorf("failed to unmarshal record bytes %v", err)
		return nil, status.Error(codes.InvalidArgument, "invalid record")
	}

	echo, err := g.service.AddOrUpdateRecord(ctx, record)
	if err != nil {
		g.log.Errorf("failed to add or update record %v", err)
		return nil, status.Error(codes.Internal, "unable to add or update record")
	}

	return &pbv1.StoreResponse{
		Echo: &pbv1.Echo{
			Ip:          echo.IP,
			Port:        echo.Port,
			NodeId:      echo.NodeID,
			CompletedAt: timestamppb.Now(),
		},
	}, nil
}

// FindNode
func (g *GRPCServer) FindNode(ctx context.Context, in *pbv1.FindNodeRequest) (*pbv1.FindNodeResponse, error) {

	g.log.Debugw("FindNode", "request", in)

	bucketSlice, err := g.service.NodeLookup(ctx, in.GetNodeId())
	if err != nil {
		g.log.Errorf("unable to look up node %x %v", in.GetNodeId(), err)
		return nil, status.Error(codes.Internal, "unable to lookup node")
	}

	kBucketSlice := make([]*pbv1.KBucket, len(bucketSlice))

	for i, levelBucket := range bucketSlice {
		kcontactSlice := make([]*pbv1.KContact, 0, len(levelBucket.Contacts))

		for _, contact := range levelBucket.Contacts {
			kcontactSlice = append(kcontactSlice, &pbv1.KContact{
				NodeId:  contact.NodeID,
				Address: net.JoinHostPort(contact.IP, fmt.Sprintf("%d", contact.Port)),
			})
		}

		kBucketSlice[i].ContactList = kcontactSlice
	}

	echo := g.service.Echo(ctx)
	return &pbv1.FindNodeResponse{
		Echo: &pbv1.Echo{
			Ip:          echo.IP,
			Port:        echo.Port,
			NodeId:      echo.NodeID,
			CompletedAt: timestamppb.Now(),
		},
		BucketList: kBucketSlice,
	}, nil
}

// FindValue
func (g *GRPCServer) FindValue(ctx context.Context, in *pbv1.FindValueRequest) (*pbv1.FindValueResponse, error) {

	g.log.Debugw("FindValue", "request", in)

	return &pbv1.FindValueResponse{
		Echo: &pbv1.Echo{},
	}, nil
}

// Start
func (g *GRPCServer) Start() error {

	g.srv = grpc.NewServer()
	pbv1.RegisterDDNSServiceV1Server(g.srv, g)

	listener, err := net.Listen("tcp", net.JoinHostPort(g.service.GetHost(), fmt.Sprintf("%d", g.service.GetPort())))
	if err != nil {
		return fmt.Errorf("failed to create network listener %v", err)
	}

	go func() {
		if err := g.srv.Serve(listener); err != nil {
			g.log.Fatalf("failed to start gRPC server %v", err)
		}
	}()
	return nil
}

// Shutdown
func (g *GRPCServer) Shutdown() {
	g.srv.GracefulStop()
}
