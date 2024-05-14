package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	dpkg "github.com/structx/go-dpkg/domain"
	"github.com/structx/go-dpkg/structs/dht"

	"github.com/structx/ddns/internal/adapter/port/rpcfx"
	"github.com/structx/ddns/internal/core/domain"
)

// DDNS
type DDNS struct {
	dht  *dht.Node
	ip   string
	port int64
}

// interface compliance
var _ domain.DDNS = (*DDNS)(nil)

// NewDDNS
func NewDDNS(ctx context.Context, cfg dpkg.Config) *DDNS {
	dcfg := cfg.GetDistributedHashTable()
	return &DDNS{
		dht: dht.NewNodeWithDefault(ctx, dcfg.BindAddr, dcfg.Ports.GRPC),
	}
}

// AddOrUpdateRecord
func (dd *DDNS) AddOrUpdateRecord(ctx context.Context, record domain.Record) (*domain.Echo, error) {

	recordbytes, err := json.Marshal(record)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal record %v", err)
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
				return nil, fmt.Errorf("unable to create client %v", err)
			}

			err = client.StoreValue(ctx, []byte{}, recordbytes, dd.dht.ID[:])
			if err != nil {
				return nil, fmt.Errorf("failed to store value %v", err)
			}
		}
	}

	return &domain.Echo{
		IP:     dd.ip,
		Port:   dd.port,
		NodeID: dd.dht.ID[:],
	}, nil
}

// NodeLookup
func (dd *DDNS) NodeLookup(ctx context.Context, nodeID []byte) ([]*domain.Bucket, error) {

	bucketIDSlice := dd.dht.FindKClosestBuckets(ctx, nodeID)

	for _, bucketID := range bucketIDSlice {

		if bytes.Equal(bucketID[:], dd.dht.ID[:]) {

			// local node is the queried node
			levelBucket := dd.dht.Get(ctx, nodeID)
			bucket := &domain.Bucket{ID: levelBucket.ID[:]}
			bucket.Contacts = make([]*domain.Contact, len(levelBucket.Contacts))

			for i, contact := range levelBucket.Contacts {
				bucket.Contacts[i] = &domain.Contact{
					IP:     contact.IP,
					Port:   int64(contact.Port),
					NodeID: contact.ID[:],
				}
			}
			return []*domain.Bucket{bucket}, nil
		}

		// find remote node addresses
		contactAddrSlice := dd.dht.FindClosestNodes(ctx, nodeID, bucketID)

		for _, addr := range contactAddrSlice {

			// call to remote node
			client, err := rpcfx.NewClient(addr)
			if err != nil {
				return nil, fmt.Errorf("unable to create client %v", err)
			}

			levelBucketSlice, err := client.FindNode(ctx, nodeID, dd.dht.ID[:])
			if err != nil {
				return nil, fmt.Errorf("failed to store value %v", err)
			}

			return transformLevelBucketSlice(levelBucketSlice), nil
		}
	}

	return nil, nil
}

// Health
func (dd *DDNS) Echo(ctx context.Context) *domain.Echo {
	select {
	case <-ctx.Done():
		return nil
	default:
		return &domain.Echo{
			IP:     dd.ip,
			Port:   dd.port,
			NodeID: dd.dht.ID[:],
		}
	}
}

func transformLevelBucketSlice(s []*dpkg.Bucket) []*domain.Bucket {

	bucketslice := make([]*domain.Bucket, len(s))

	for _, levelBucket := range s {

		bucket := &domain.Bucket{ID: levelBucket.ID[:]}
		bucket.Contacts = make([]*domain.Contact, len(levelBucket.Contacts))

		for i, contact := range levelBucket.Contacts {
			bucket.Contacts[i] = &domain.Contact{
				IP:     contact.IP,
				Port:   int64(contact.Port),
				NodeID: contact.ID[:],
			}
		}

	}

	return bucketslice
}
