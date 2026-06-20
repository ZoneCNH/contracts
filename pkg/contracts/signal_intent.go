package contracts

// SignalIntent is the P1 output contract of signal_factory → risk_engine / order_engine.
// It represents a directional trading intent for a single symbol, modulated by
// the regime context from DecisionCard.
//
// Lifecycle:
//   - signal_factory produces SignalIntents via Generate(DecisionCard, symbols).
//   - risk_engine validates and sizes each SignalIntent against position limits.
//   - order_engine converts approved SignalIntents into ExecutionOrder.
type SignalIntent struct {
	// ID is a UUID v4 uniquely identifying this signal intent.
	ID string `json:"id"`
	// GeneratedAt is the creation time in UTC milliseconds.
	GeneratedAt int64 `json:"generated_at"`
	// CardID links this intent back to the DecisionCard that drove it.
	CardID string `json:"card_id"`

	// Symbol is the normalized instrument identifier (e.g. "BTCUSDT").
	Symbol string `json:"symbol"`

	// Action mirrors DecisionCard.Action (A-E) for downstream audit.
	Action Action `json:"action"`
	// Template is the selected strategy template for this regime.
	Template StrategyTemplate `json:"template"`

	// Strength is the signal conviction in [0.0, 1.0].
	// Derived from base strength * conflict gate.
	Strength float64 `json:"strength"`
	// SizePct is the suggested position size as a fraction of allowed maximum.
	// Derived from Strength × DecisionCard.RiskMultiplier.
	SizePct float64 `json:"size_pct"`

	// Conflict is forwarded from DecisionCard; risk_engine applies additional haircut.
	Conflict bool `json:"conflict"`
	// Explain records the signal derivation rationale.
	Explain string `json:"explain,omitempty"`
}

// SignalFactoryProvider is the primary interface of signal_factory.
// It accepts a DecisionCard and a list of symbols to generate per-symbol SignalIntents.
type SignalFactoryProvider interface {
	Generate(card DecisionCard, symbols []string) ([]SignalIntent, error)
}
