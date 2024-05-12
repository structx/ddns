package rpcfx

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	pbv1 "github.com/structx/ddns/proto/ddns/v1"
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

// StoreValue
func (c *Client) StoreValue(ctx context.Context, key, value []byte) error {
	c.conn.Connect()

	timeout, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()

	cli := pbv1.NewDDNSServiceV1Client(c.conn)
	_, err := cli.Store(timeout, &pbv1.StoreRequest{
		Sender: &pbv1.Sender{
			NodeId:      []byte{},
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
