package contracts

import (
	"encoding/json"
	"time"
)

// ---- Binance C/S Ingestion Wire Contract (SPEC §8.4) ----

// MarketDataService receives normalized upstream market-data ingestion requests
// from exchange adapters (e.g. module/binance).
// Transport: gRPC unary request/response.
type MarketDataService interface {
	// Ingest accepts one IngestRequest and returns one IngestResult.
	Ingest(in IngestRequest) (IngestResult, error)
}

// IngestRequest is an immutable item submitted by an exchange adapter.
type IngestRequest struct {
	RequestID      string            `json:"request_id"`
	Source         string            `json:"source"`
	ProductLine    string            `json:"product_line"`
	InstrumentKey  json.RawMessage   `json:"instrument_key"`
	EventType      string            `json:"event_type"`
	EventTime      time.Time         `json:"event_time"`
	ReceivedAt     time.Time         `json:"received_at"`
	SchemaVersion  string            `json:"schema_version"`
	Payload        json.RawMessage   `json:"payload"`
	Sequence       int64             `json:"sequence,omitempty"`
	OrderingKey    string            `json:"ordering_key,omitempty"`
	SourceMetadata map[string]string `json:"source_metadata"`
}

// IngestResult is a terminal outcome for exactly one request_id.
// Exactly one of Ack or Reject is non-nil.
type IngestResult struct {
	RequestID string        `json:"request_id"`
	Ack       *IngestAck    `json:"ack,omitempty"`
	Reject    *IngestReject `json:"reject,omitempty"`
}

// IngestAck confirms the receiver accepted one request.
type IngestAck struct {
	RequestID      string `json:"request_id"`
	StreamID       string `json:"stream_id"`
	AcceptedCount  int32  `json:"accepted_count"`
	DuplicateCount int32  `json:"duplicate_count"`
	Durable        bool   `json:"durable"`
}

// IngestReject explains why one request was not accepted.
type IngestReject struct {
	RequestID  string     `json:"request_id"`
	RejectCode RejectCode `json:"reject_code"`
	Reason     string     `json:"reason"`
	Retryable  bool       `json:"retryable"`
}

// RejectCode classifies rejection reasons for adapter retry policy decisions.
type RejectCode string

const (
	// RejectRetryable: temporary unavailability; caller should retry with backoff.
	RejectRetryable RejectCode = "retryable"

	// RejectTerminalValidation: request fails validation; retry won't help.
	RejectTerminalValidation RejectCode = "terminal_validation"

	// RejectTerminalConflict: duplicate request_id with conflicting payload.
	RejectTerminalConflict RejectCode = "terminal_conflict"

	// RejectUnauthorized: caller lacks credentials or permissions.
	RejectUnauthorized RejectCode = "unauthorized"

	// RejectRateLimited: caller exceeded rate limit; should back off.
	RejectRateLimited RejectCode = "rate_limited"

	// RejectServerUnavailable: server cannot accept due to internal state; retryable.
	RejectServerUnavailable RejectCode = "server_unavailable"

	// RejectContractViolation: request violates wire contract.
	RejectContractViolation RejectCode = "contract_violation"

	// RejectQualityRejected: event fails quality gate (stale, future, dirty).
	RejectQualityRejected RejectCode = "quality_rejected"

	// RejectOrderingViolation: sequence gap, reversal, or ordering_key mismatch.
	RejectOrderingViolation RejectCode = "ordering_violation"
	// RejectUnsupportedChannel: caller used an unsupported transport or channel.
	RejectUnsupportedChannel RejectCode = "unsupported_channel"
)

// IsRetryable returns true if the adapter should retry.
func (c RejectCode) IsRetryable() bool {
	switch c {
	case RejectRetryable, RejectServerUnavailable, RejectRateLimited:
		return true
	default:
		return false
	}
}

// IsTerminal returns true if retry will not help.
func (c RejectCode) IsTerminal() bool {
	return !c.IsRetryable()
}

// AllRejectCodes returns all 10 canonical RejectCode values.
func AllRejectCodes() []RejectCode {
	return []RejectCode{
		RejectRetryable,
		RejectTerminalValidation,
		RejectTerminalConflict,
		RejectUnauthorized,
		RejectRateLimited,
		RejectServerUnavailable,
		RejectContractViolation,
		RejectQualityRejected,
		RejectOrderingViolation,
		RejectUnsupportedChannel,
	}
}
