// Package contracts defines cross-domain stable port, event, and DTO contracts
// for the FoundationX ecosystem.
//
// contracts follows xlib-standard governance but is NOT the standard source.
// It owns DTO, Event Envelope, Command/Query, Port Interface, and Error Code Registry.
// It does NOT own transport implementations (HTTP/gRPC/Kafka/NATS).
package contracts

// Event represents a domain event in the FoundationX ecosystem.
type Event struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Source  string `json:"source"`
	Version string `json:"version"`
	Data    []byte `json:"data"`
}

// Command represents a domain command.
type Command struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Target string `json:"target"`
	Data   []byte `json:"data"`
}

// Query represents a domain query.
type Query struct {
	ID     string            `json:"id"`
	Type   string            `json:"type"`
	Filter map[string]string `json:"filter,omitempty"`
}

// DTO is a marker interface for data transfer objects.
type DTO interface {
	IsDTO() bool
}

// Port is a marker interface for domain port contracts.
type Port interface {
	IsPort() bool
}

// ErrorCode represents a registered error code in the contract registry.
type ErrorCode struct {
	Code      string `json:"code"`
	Domain    string `json:"domain"`
	Severity  string `json:"severity"`
	Retryable bool   `json:"retryable"`
}
