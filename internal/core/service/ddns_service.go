package service

import (
	"context"

	dpkg "github.com/structx/go-dpkg/domain"
	"github.com/structx/go-dpkg/structs/dht"
)

// DDNS
type DDNS struct {
	dht *dht.Node
}

// NewDDNS
func NewDDNS(ctx context.Context, cfg dpkg.Config) *DDNS {
	return &DDNS{
		dht: dht.NewNodeWithDefault(ctx, "127.0.0.1", 50051),
	}
}
