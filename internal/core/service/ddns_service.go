package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/structx/ddns/internal/adapter/port/rpcfx"
	"github.com/structx/ddns/internal/core/domain"
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

	recordbytes, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal record %v", err)
	}

	bucketIDSlice := dd.dht.FindKClosestBuckets(ctx, recordbytes)

BUCKETS:
	for _, bucketID := range bucketIDSlice {
		// if bucket is local node then insert
		if bytes.Equal(bucketID[:], dd.dht.ID[:]) {
			dd.dht.AddOrUpdateNode(ctx, recordbytes, record)
			continue BUCKETS
		}
		// find remote node addresses
		contactAddrSlice := dd.dht.FindClosestNodes(ctx, recordbytes, bucketID)

		for _, addr := range contactAddrSlice {

			// call to remote node

			client, err := rpcfx.NewClient(addr)
			if err != nil {
				return fmt.Errorf("unable to create client %v", err)
			}

			err = client.StoreValue(ctx, []byte{}, recordbytes)
			if err != nil {
				return fmt.Errorf("failed to store value %v", err)
			}
		}
	}

	return nil
}
