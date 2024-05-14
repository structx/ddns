package rpcfx

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	pbv1 "github.com/structx/ddns/proto/ddns/v1"
	"github.com/structx/go-dpkg/domain"
)

// Client
type Client struct {
	conn *grpc.ClientConn
}

// NewClient
func NewClient(addr string) (*Client, error) {

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial %s %v", addr, err)
	}

	return &Client{
		conn: conn,
	}, nil
}

// FindNode
func (c *Client) FindNode(ctx context.Context, nodeID, senderID []byte) ([]*domain.Bucket, error) {
	c.conn.Connect()

	timeout, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()

	cli := pbv1.NewDDNSServiceV1Client(c.conn)
	response, err := cli.FindNode(timeout, &pbv1.FindNodeRequest{
		Sender: &pbv1.Sender{
			NodeId:      senderID,
			RequestedAt: timestamppb.Now(),
		},
		NodeId: nodeID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to send find node request %v", err)
	}

	bucketSlice := make([]*domain.Bucket, len(response.BucketList))
	for i, bucket := range response.GetBucketList() {

		bucketSlice[i] = &domain.Bucket{
			ID: domain.NodeID224(bucket.GetNodeId()),
		}
		contactSlice := make([]*domain.Contact, len(bucket.ContactList))

		for _, contact := range bucket.GetContactList() {
			host, port, _ := net.SplitHostPort(contact.GetAddress())
			num, _ := strconv.Atoi(port)
			contactSlice[i] = &domain.Contact{
				IP:   host,
				Port: num,
				ID:   domain.NodeID224(contact.NodeId),
			}
		}
		bucketSlice[i].Contacts = contactSlice
	}

	return bucketSlice, nil
}

// StoreValue
func (c *Client) StoreValue(ctx context.Context, key, value, senderID []byte) error {
	c.conn.Connect()

	timeout, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()

	cli := pbv1.NewDDNSServiceV1Client(c.conn)
	_, err := cli.Store(timeout, &pbv1.StoreRequest{
		Sender: &pbv1.Sender{
			NodeId:      senderID,
			RequestedAt: timestamppb.Now(),
		},
		Key:   key,
		Value: value,
	})
	if err != nil {
		return fmt.Errorf("failed to call store %v", err)
	}

	return nil
}
