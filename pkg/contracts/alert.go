package contracts

import (
	"context"
	"time"
)

// AlertEvent is the P1 output contract consumed by alertx.
// It represents a fired or resolved alert produced by a rule evaluation,
// emitted by business domains (riskx, strategyx) and the observability layer.
//
// Producer:
//   - alertx: evaluates rules over observex exports (LogEntry/MetricPoint/SpanData)
//     and business-domain state events, producing AlertEvents.
//   - Business domains (riskx, strategyx): may emit AlertEvents directly for
//     policy/risk-class alerts (the business-side input of the dual-subscription
//     model, see module/alertx/ADR-001-foundations.md D2).
//
// Consumer:
//   - alertx: deduplicates, severity-grades, routes to notification channels.
//   - Operations/monitoring: receives routed alert notifications.
//
// Envelope: AlertEvent is carried as the Data payload of a contracts.Event
// (Event.Type = "alert.fired" | "alert.resolved"). Producers MUST wrap it in
// Event; consumers SHOULD decode via Event.Data.
type AlertEvent struct {
	// ID is a unique identifier for this alert instance (typically a UUID v4).
	ID string `json:"id"`
	// Source identifies the producing module and rule (e.g. "riskx:breach",
	// "alertx:metric_threshold:foundationx_observex_exporter_queue_size").
	Source string `json:"source"`
	// Severity is the alert severity grade (critical | warning | info).
	Severity Severity `json:"severity"`
	// Status is the alert lifecycle state (firing | pending | resolved | suppressed).
	Status AlertStatus `json:"status"`
	// Message is the human-readable alert summary.
	Message string `json:"message"`
	// Context holds structured key-value context for the alert (e.g. metric name,
	// observed value, threshold, symbol). Values MUST be redacted by the producer
	// if they contain secrets.
	Context map[string]string `json:"context,omitempty"`
	// FiredAt is the UTC time at which the alert first fired.
	FiredAt time.Time `json:"fired_at"`
	// ResolvedAt is the UTC time at which the alert transitioned to resolved.
	// nil means the alert is still active. Pointer (not time.Time zero value) is
	// used so that omitempty correctly omits unresolved alerts — time.Time is a
	// struct and its zero value is not considered empty by encoding/json.
	ResolvedAt *time.Time `json:"resolved_at,omitempty"`
	// TraceID links the alert to the originating trace context when available,
	// enabling cross-module correlation via observex spans.
	TraceID string `json:"trace_id,omitempty"`
	// DedupKey is the canonical fingerprint used by alertx deduplication.
	// Two AlertEvents with the same DedupKey within a suppression window are
	// collapsed into one notification. Producers SHOULD set a stable key based
	// on (source, subject); alertx computes a fallback if empty.
	DedupKey string `json:"dedup_key,omitempty"`
}

// Severity grades an alert's urgency. It maps to the incident levels I0–I5
// defined in docs/goal/rsi-standard/23-monitoring-incident-response.md:
//
//   - SeverityCritical: I4–I5 (forced-pause conditions, kill-switch violations,
//     systemic failures) — immediate paging.
//   - SeverityWarning:  I2–I3 (degraded operation, threshold breaches) — notify
//     operations, no paging.
//   - SeverityInfo:     I0–I1 (informational, routine state changes) — log only.
type Severity string

const (
	// SeverityCritical requires immediate attention (paging channel).
	SeverityCritical Severity = "critical"
	// SeverityWarning indicates degraded operation; notify but do not page.
	SeverityWarning Severity = "warning"
	// SeverityInfo is informational; log only, no notification channel.
	SeverityInfo Severity = "info"
)

// AlertStatus is the lifecycle state of an alert instance managed by alertx.
type AlertStatus string

const (
	// AlertStatusFiring: the alert condition is currently active and notification
	// has been dispatched (post-dedup).
	AlertStatusFiring AlertStatus = "firing"
	// AlertStatusPending: the alert condition was detected but is within a
	// pending window before notification (e.g. flap suppression).
	AlertStatusPending AlertStatus = "pending"
	// AlertStatusResolved: the alert condition has cleared; a resolve notification
	// may have been sent depending on rule configuration.
	AlertStatusResolved AlertStatus = "resolved"
	// AlertStatusSuppressed: the alert was collapsed by deduplication/suppression
	// and no new notification was sent.
	AlertStatusSuppressed AlertStatus = "suppressed"
)

// AlertRule is the declarative rule contract evaluated by alertx.
// Rules are defined as a YAML DSL (see module/alertx/ADR-001-foundations.md D3)
// and validated at alertx startup; an invalid rule blocks process start.
//
// This contract captures the stable schema of a rule; the DSL parser in alertx
// internal/config maps YAML documents onto AlertRule instances.
type AlertRule struct {
	// ID is the unique rule identifier referenced by alerts and metrics.
	ID string `json:"id"`
	// Name is the human-readable rule name.
	Name string `json:"name"`
	// Source identifies the rule's scope (e.g. "riskx", "observex.metrics",
	// "system.health"). Matches the AlertEvent.Source namespace.
	Source string `json:"source"`
	// Severity is the severity assigned to alerts produced by this rule.
	Severity Severity `json:"severity"`
	// Condition is the rule predicate expression in the alertx rule DSL.
	// The DSL grammar is owned by alertx internal/config; contracts carries the
	// opaque expression string to keep contracts free of parser logic.
	Condition string `json:"condition"`
	// DedupKey is an optional explicit deduplication key template; when empty,
	// alertx derives a key from (ID, subject).
	DedupKey string `json:"dedup_key,omitempty"`
	// SuppressWindow is the dedup/suppression duration. Re-firing alerts with
	// the same DedupKey within this window are suppressed. Zero means no
	// suppression (every match notifies).
	SuppressWindow time.Duration `json:"suppress_window,omitempty"`
	// Channels lists the notification channel IDs this rule routes to.
	// Channel implementations (webhook/email/pagerduty) are owned by alertx
	// internal/channel; contracts carries opaque IDs.
	Channels []string `json:"channels,omitempty"`
	// Enabled controls whether the rule is evaluated. Defaults to true.
	Enabled bool `json:"enabled"`
}

// AlertSink is the P1 port interface exposed by alertx for downstream
// consumers that wish to subscribe to the alert stream (e.g. an incident
// manager, a dashboard, or a replay harness).
//
// Implementations live in alertx; this contract is stable across module
// versions. Producers of AlertEvent (business domains) emit via the
// contracts.Event envelope, not via AlertSink — AlertSink is the outbound
// subscription surface.
type AlertSink interface {
	Port

	// SubscribeAlerts streams AlertEvents to the subscriber.
	// The returned channel is closed when ctx is cancelled or an unrecoverable
	// error occurs (consistent with MarketDataProvider.SubscribeRegimeSnapshots).
	SubscribeAlerts(ctx context.Context) (<-chan AlertEvent, error)
}

// AlertRuleStore is the P1 port interface for loading and hot-reloading rules.
// alertx implements this against its YAML rule files (FR-007 hot-reload);
// alternative implementations may load from a database or remote config.
type AlertRuleStore interface {
	Port

	// Load returns the current set of rules. Implementations MUST return a
	// snapshot that is safe for concurrent evaluation.
	Load(ctx context.Context) ([]AlertRule, error)
}
