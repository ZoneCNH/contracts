package contracts

import "encoding/json"

// DecisionCard is the P0 output contract of regime_engine → signal_factory / risk_engine.
// It is the authoritative cross-domain action directive produced by the M×S joint
// decision matrix (MacroState × MarketState).
//
// Consumers:
//   - signal_factory: uses Action + Profile + Template for signal strength modulation.
//   - risk_engine:    uses RiskTier + ExposureCaps + RiskMultiplier for position limits.
//   - audit/monitor:  uses Conflict + Explain for human review.
type DecisionCard struct {
	// CardID is a UUID v4 uniquely identifying this decision card.
	CardID string `json:"card_id"`
	// GeneratedAt is the card creation time in UTC milliseconds.
	GeneratedAt int64 `json:"generated_at"`

	// Regime inputs (recorded for audit traceability).
	MarketState MarketState `json:"market_state"` // from RegimeSnapshot.RegimeState
	MacroState  MacroState  `json:"macro_state"`  // from RegimeCard.MState

	// Action directive (consumed by signal_factory).

	// Action is the regime-level positioning action from A (aggressive) to E (flat/cash).
	Action Action `json:"action"`
	// Profile refines Action with a sizing modifier.
	Profile ActionProfile `json:"profile"`
	// Template selects the signal_factory strategy template for this regime.
	Template StrategyTemplate `json:"template"`

	// Risk controls (consumed by risk_engine).

	// RiskTier is an integer in [1, 5] where 1 = lowest risk appetite, 5 = highest.
	RiskTier int `json:"risk_tier"`
	// ExposureCaps constrains maximum leverage and position size relative to NAV.
	ExposureCaps ExposureCaps `json:"position_caps"`
	// PositionCaps is a deprecated compatibility alias kept for Go callers that
	// still reference the old field name.
	PositionCaps ExposureCaps `json:"-"`
	// RiskMultiplier scales the base risk budget; range [0.3, 1.0].
	RiskMultiplier float64 `json:"risk_multiplier"`

	// Auditability.

	// Conflict is true when M-State and S-State produce contradictory signals.
	// risk_engine MUST apply an additional haircut when Conflict is true.
	Conflict bool `json:"conflict"`
	// Explain is a human-readable description of the M×S decision rationale.
	Explain string `json:"explain,omitempty"`
}

// Action is the regime-level positioning directive produced by the M×S matrix.
// A represents the most aggressive posture; E represents fully flat (cash only).
type Action string

const (
	// ActionA: aggressive — full risk budget, strong directional conviction.
	ActionA Action = "A"
	// ActionB: moderate — partial risk budget, measured exposure.
	ActionB Action = "B"
	// ActionC: conservative — minimal risk budget, defensive bias.
	ActionC Action = "C"
	// ActionD: defensive — hedge-oriented, capital preservation priority.
	ActionD Action = "D"
	// ActionE: flat — cash only, all positions closed or in hedge.
	ActionE Action = "E"
)

// ActionProfile refines Action with a sizing qualifier, allowing signal_factory
// to select the appropriate strategy parameter set.
type ActionProfile string

const (
	ProfileAggressive   ActionProfile = "aggressive"
	ProfileModerate     ActionProfile = "moderate"
	ProfileConservative ActionProfile = "conservative"
	ProfileDefensive    ActionProfile = "defensive"
	ProfileFlat         ActionProfile = "flat"
)

// StrategyTemplate selects the signal_factory strategy template for this regime.
type StrategyTemplate string

const (
	// TemplateTrendFollowing: directional momentum strategies.
	TemplateTrendFollowing StrategyTemplate = "trend_following"
	// TemplateRangeTrading: mean-reversion within defined bands.
	TemplateRangeTrading StrategyTemplate = "range_trading"
	// TemplateBreakout: volatility expansion and breakout strategies.
	TemplateBreakout StrategyTemplate = "breakout"
	// TemplateHedge: delta-neutral or inverse exposure strategies.
	TemplateHedge StrategyTemplate = "hedge"
	// TemplateCash: no active strategies; capital held in stable assets.
	TemplateCash StrategyTemplate = "cash"
)

// ExposureCaps constrains risk_engine position sizing derived from this DecisionCard.
type ExposureCaps struct {
	// MaxLeverage is the leverage ceiling (e.g. 2.0 = 2× notional).
	MaxLeverage float64 `json:"max_leverage"`
	// MaxExposurePct is the maximum single-exposure size as a fraction of NAV (e.g. 0.10 = 10%).
	MaxExposurePct float64 `json:"max_position_pct"`
}

// PositionCaps is a deprecated alias for ExposureCaps kept for Go source compatibility.
type PositionCaps = ExposureCaps

type decisionCardWire struct {
	CardID         string           `json:"card_id"`
	GeneratedAt    int64            `json:"generated_at"`
	MarketState    MarketState      `json:"market_state"`
	MacroState     MacroState       `json:"macro_state"`
	Action         Action           `json:"action"`
	Profile        ActionProfile    `json:"profile"`
	Template       StrategyTemplate `json:"template"`
	RiskTier       int              `json:"risk_tier"`
	ExposureCaps   ExposureCaps     `json:"position_caps"`
	RiskMultiplier float64          `json:"risk_multiplier"`
	Conflict       bool             `json:"conflict"`
	Explain        string           `json:"explain,omitempty"`
}

// MarshalJSON keeps the wire field name stable while accepting either field name in Go.
func (d DecisionCard) MarshalJSON() ([]byte, error) {
	wire := decisionCardWire{
		CardID:         d.CardID,
		GeneratedAt:    d.GeneratedAt,
		MarketState:    d.MarketState,
		MacroState:     d.MacroState,
		Action:         d.Action,
		Profile:        d.Profile,
		Template:       d.Template,
		RiskTier:       d.RiskTier,
		ExposureCaps:   d.ExposureCaps,
		RiskMultiplier: d.RiskMultiplier,
		Conflict:       d.Conflict,
		Explain:        d.Explain,
	}
	if wire.ExposureCaps == (ExposureCaps{}) && d.PositionCaps != (ExposureCaps{}) {
		wire.ExposureCaps = d.PositionCaps
	}
	return json.Marshal(wire)
}

// UnmarshalJSON keeps the canonical and compatibility fields in sync.
func (d *DecisionCard) UnmarshalJSON(data []byte) error {
	var wire decisionCardWire
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}

	d.CardID = wire.CardID
	d.GeneratedAt = wire.GeneratedAt
	d.MarketState = wire.MarketState
	d.MacroState = wire.MacroState
	d.Action = wire.Action
	d.Profile = wire.Profile
	d.Template = wire.Template
	d.RiskTier = wire.RiskTier
	d.ExposureCaps = wire.ExposureCaps
	d.PositionCaps = wire.ExposureCaps
	d.RiskMultiplier = wire.RiskMultiplier
	d.Conflict = wire.Conflict
	d.Explain = wire.Explain
	return nil
}
