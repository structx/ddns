package domain

import "context"

// DDNS
type DDNS interface {
	AddOrUpdateRecord(ctx context.Context, r Record) error
}
