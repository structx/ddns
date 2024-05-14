package domain

import "context"

// Echo
type Echo struct {
	IP     string
	Port   int64
	NodeID []byte
}

// Bucket
type Bucket struct {
	ID       []byte
	Contacts []*Contact
}

// Contact
type Contact struct {
	IP     string
	Port   int64
	NodeID []byte
}

// DDNS
//
//go:generate mockery --name DDNS
type DDNS interface {
	// AddOrUpdateRecord
	AddOrUpdateRecord(ctx context.Context, r Record) (*Echo, error)
	// Echo
	Echo(ctx context.Context) *Echo
	// NodeLookup
	NodeLookup(ctx context.Context, nodeID []byte) ([]*Bucket, error)
	// GetHost
	GetHost() string
	// GetPort
	GetPort() int64
}
