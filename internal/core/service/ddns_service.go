package service

import (
	"context"

	"githhub.com/structx/ddns/internal/core/domain"
	dpkg "github.com/structx/go-dpkg/domain"
	"github.com/structx/go-dpkg/structs/dht"
)

// DDNS
type DDNS struct {
	dht *dht.Node
}

// interface compliance
var _ domain.DDNS = (*DDNS)(nil)

// NewDDNS
func NewDDNS(ctx context.Context, cfg dpkg.Config) *DDNS {
	return &DDNS{
		dht: dht.NewNodeWithDefault(ctx, "127.0.0.1", 50051),
	}
}

// AddOrUpdateRecord
func (dd *DDNS) AddOrUpdateRecord(ctx context.Context, record domain.Record) error {

	// TODO:
	// implement DHT Put

	return nil
}
