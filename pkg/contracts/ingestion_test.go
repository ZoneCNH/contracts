package contracts

import (
	"encoding/json"
	"testing"
	"time"
)

func TestIngestRequest_JSONRoundTrip(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	req := IngestRequest{
		RequestID:     "req-001",
		Source:        "binance",
		ProductLine:   "spot",
		InstrumentKey: json.RawMessage(`{"venue":"binance","symbol":"BTCUSDT"}`),
		EventType:     "trade",
		EventTime:     now,
		ReceivedAt:    now.Add(10 * time.Millisecond),
		SchemaVersion: "1.0.0",
		Payload:       json.RawMessage(`{"price":"50000","qty":"1.5"}`),
		Sequence:      12345,
		OrderingKey:   "binance:spot:BTCUSDT:trade",
		SourceMetadata: map[string]string{
			"stream_id":         "s-001",
			"connector_version": "0.1.0",
		},
	}

	b, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var req2 IngestRequest
	if err := json.Unmarshal(b, &req2); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if req2.RequestID != "req-001" {
		t.Errorf("RequestID mismatch: %s", req2.RequestID)
	}
	if req2.Source != "binance" {
		t.Errorf("Source mismatch: %s", req2.Source)
	}
	if req2.ProductLine != "spot" {
		t.Errorf("ProductLine mismatch: %s", req2.ProductLine)
	}
	if req2.Sequence != 12345 {
		t.Errorf("Sequence mismatch: %d", req2.Sequence)
	}
	if req2.SourceMetadata["stream_id"] != "s-001" {
		t.Errorf("SourceMetadata missing stream_id")
	}
}

func TestIngestResult_Ack(t *testing.T) {
	result := IngestResult{
		RequestID: "req-001",
		Ack: &IngestAck{
			RequestID:      "req-001",
			StreamID:       "s-001",
			AcceptedCount:  1,
			DuplicateCount: 0,
			Durable:        true,
		},
	}

	if result.Ack == nil {
		t.Error("Ack should be non-nil")
	}
	if result.Reject != nil {
		t.Error("Reject should be nil when Ack is set")
	}
}

func TestIngestResult_Reject(t *testing.T) {
	result := IngestResult{
		RequestID: "req-002",
		Reject: &IngestReject{
			RequestID:  "req-002",
			RejectCode: RejectContractViolation,
			Reason:     "missing product_line",
			Retryable:  false,
		},
	}

	if result.Reject == nil {
		t.Error("Reject should be non-nil")
	}
	if result.Ack != nil {
		t.Error("Ack should be nil when Reject is set")
	}
}

func TestRejectCode_IsRetryable(t *testing.T) {
	tests := []struct {
		code     RejectCode
		retryable bool
	}{
		{RejectRetryable, true},
		{RejectServerUnavailable, true},
		{RejectRateLimited, true},
		{RejectTerminalValidation, false},
		{RejectTerminalConflict, false},
		{RejectUnauthorized, false},
		{RejectContractViolation, false},
		{RejectQualityGate, false},
		{RejectOrderingViolation, false},
	}
	for _, tt := range tests {
		t.Run(string(tt.code), func(t *testing.T) {
			if got := tt.code.IsRetryable(); got != tt.retryable {
				t.Errorf("%s.IsRetryable() = %v, want %v", tt.code, got, tt.retryable)
			}
		})
	}
}

func TestAllRejectCodes(t *testing.T) {
	codes := AllRejectCodes()
	if len(codes) != 9 {
		t.Errorf("expected 9 reject codes, got %d", len(codes))
	}
}
