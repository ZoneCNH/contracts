package contracts_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/ZoneCNH/contracts/pkg/contracts"
)

func TestAlertEvent_JSONRoundTrip(t *testing.T) {
	firedAt := time.Date(2026, 6, 26, 10, 0, 0, 0, time.UTC)
	resolvedAt := time.Date(2026, 6, 26, 10, 5, 0, 0, time.UTC)
	event := contracts.AlertEvent{
		ID:       "alert-uuid-001",
		Source:   "riskx:breach",
		Severity: contracts.SeverityCritical,
		Status:   contracts.AlertStatusFiring,
		Message:  "max drawdown breached: -8.5% < -8.0%",
		Context: map[string]string{
			"symbol":    "BTCUSDT",
			"observed":  "-8.5",
			"threshold": "-8.0",
		},
		FiredAt:    firedAt,
		ResolvedAt: &resolvedAt,
		TraceID:    "trace-abc123",
		DedupKey:   "riskx:breach:BTCUSDT",
	}

	b, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got contracts.AlertEvent
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got.ID != event.ID {
		t.Errorf("ID: got %q want %q", got.ID, event.ID)
	}
	if got.Source != event.Source {
		t.Errorf("Source: got %q want %q", got.Source, event.Source)
	}
	if got.Severity != contracts.SeverityCritical {
		t.Errorf("Severity: got %q want %q", got.Severity, contracts.SeverityCritical)
	}
	if got.Status != contracts.AlertStatusFiring {
		t.Errorf("Status: got %q want %q", got.Status, contracts.AlertStatusFiring)
	}
	if got.Message != event.Message {
		t.Errorf("Message: got %q want %q", got.Message, event.Message)
	}
	if !got.FiredAt.Equal(event.FiredAt) {
		t.Errorf("FiredAt: got %v want %v", got.FiredAt, event.FiredAt)
	}
	if got.ResolvedAt == nil || !got.ResolvedAt.Equal(resolvedAt) {
		t.Errorf("ResolvedAt: got %v want %v", got.ResolvedAt, resolvedAt)
	}
	if got.TraceID != event.TraceID {
		t.Errorf("TraceID: got %q want %q", got.TraceID, event.TraceID)
	}
	if got.DedupKey != event.DedupKey {
		t.Errorf("DedupKey: got %q want %q", got.DedupKey, event.DedupKey)
	}
	if got.Context["symbol"] != "BTCUSDT" {
		t.Errorf("Context[symbol]: got %q want %q", got.Context["symbol"], "BTCUSDT")
	}
}

func TestAlertEvent_ResolvedZeroValue(t *testing.T) {
	// An unresolved alert has a zero ResolvedAt, which MUST be omitted from JSON
	// (omitempty on ResolvedAt).
	event := contracts.AlertEvent{
		ID:       "alert-uuid-002",
		Source:   "system.health",
		Severity: contracts.SeverityInfo,
		Status:   contracts.AlertStatusFiring,
		Message:  "routine health check",
		FiredAt:  time.Now().UTC(),
	}

	b, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		t.Fatalf("unmarshal to map: %v", err)
	}
	if _, present := raw["resolved_at"]; present {
		t.Errorf("resolved_at should be omitted for unresolved alert, got JSON: %s", b)
	}
}

func TestSeverity_Constants(t *testing.T) {
	// Severity constants are stable string values used across modules and in
	// persisted alert records; changing them is a breaking change.
	cases := []struct {
		severity contracts.Severity
		want     string
	}{
		{contracts.SeverityCritical, "critical"},
		{contracts.SeverityWarning, "warning"},
		{contracts.SeverityInfo, "info"},
	}
	for _, tc := range cases {
		if string(tc.severity) != tc.want {
			t.Errorf("severity: got %q want %q", tc.severity, tc.want)
		}
	}
}

func TestAlertStatus_Constants(t *testing.T) {
	cases := []struct {
		status contracts.AlertStatus
		want   string
	}{
		{contracts.AlertStatusFiring, "firing"},
		{contracts.AlertStatusPending, "pending"},
		{contracts.AlertStatusResolved, "resolved"},
		{contracts.AlertStatusSuppressed, "suppressed"},
	}
	for _, tc := range cases {
		if string(tc.status) != tc.want {
			t.Errorf("status: got %q want %q", tc.status, tc.want)
		}
	}
}

func TestAlertRule_JSONRoundTrip(t *testing.T) {
	rule := contracts.AlertRule{
		ID:             "rule-001",
		Name:           "exporter queue saturation",
		Source:         "observex.metrics",
		Severity:       contracts.SeverityWarning,
		Condition:      "metric:foundationx_observex_exporter_queue_size > 1000",
		SuppressWindow: 5 * time.Minute,
		Channels:       []string{"webhook-ops", "pagerduty"},
		Enabled:        true,
	}

	b, err := json.Marshal(rule)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got contracts.AlertRule
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got.ID != rule.ID {
		t.Errorf("ID: got %q want %q", got.ID, rule.ID)
	}
	if got.Severity != rule.Severity {
		t.Errorf("Severity: got %q want %q", got.Severity, rule.Severity)
	}
	if got.Condition != rule.Condition {
		t.Errorf("Condition: got %q want %q", got.Condition, rule.Condition)
	}
	if got.SuppressWindow != rule.SuppressWindow {
		t.Errorf("SuppressWindow: got %v want %v", got.SuppressWindow, rule.SuppressWindow)
	}
	if len(got.Channels) != 2 || got.Channels[0] != "webhook-ops" {
		t.Errorf("Channels: got %v want %v", got.Channels, rule.Channels)
	}
	if !got.Enabled {
		t.Errorf("Enabled: got false want true")
	}
}

func TestAlertSink_Interface(t *testing.T) {
	// 编译期断言：AlertSink 接口可由外部实现满足。
	var _ contracts.AlertSink = (*mockAlertSink)(nil)
}

type mockAlertSink struct{}

func (m *mockAlertSink) IsPort() bool { return true }

func (m *mockAlertSink) SubscribeAlerts(ctx context.Context) (<-chan contracts.AlertEvent, error) {
	return nil, nil
}

func TestAlertRuleStore_Interface(t *testing.T) {
	var _ contracts.AlertRuleStore = (*mockRuleStore)(nil)
}

type mockRuleStore struct{}

func (m *mockRuleStore) IsPort() bool { return true }

func (m *mockRuleStore) Load(ctx context.Context) ([]contracts.AlertRule, error) {
	return nil, nil
}
