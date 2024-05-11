package domain

// RecordEnum
type RecordEnum int

const (
	// ARecord
	ARecord RecordEnum = iota
	// CNameRecord
	CNameRecord
)

// Record
type Record interface {
	// GetType
	GetType() RecordEnum
	// GetRoot
	GetRoot() string
	// GetContent
	GetContent() string
	// GetTTL
	GetTTL() int64
}

// A
type A struct {
	root    string
	content string
	ttl     int64
}

// GetType
func (a *A) GetType() RecordEnum {
	return ARecord
}

// GetRoot
func (a *A) GetRoot() string {
	return a.root
}

// GetContent
func (a *A) GetContent() string {
	return a.content
}

// GetTTL
func (a *A) GetTTL() int64 {
	return a.ttl
}

// CName
type CName struct {
	root    string
	content string
	ttl     int64
}

// GetType
func (c *CName) GetType() RecordEnum {
	return CNameRecord
}

// GetRoot
func (c *CName) GetRoot() string {
	return c.root
}

// GetContent
func (c *CName) GetContent() string {
	return c.content
}

// GetTTL
func (c *CName) GetTTL() int64 {
	return c.ttl
}
