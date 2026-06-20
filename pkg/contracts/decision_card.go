package contracts

// DecisionCard is the P0 output contract of regime_engine → signal_factory / risk_engine.
// It is the authoritative cross-domain action directive produced by the M×S joint
// decision matrix (MacroState × MarketState).
//
// Consumers:
//   - signal_factory: uses Action + Profile + Template for signal strength modulation.
//   - risk_engine:    uses RiskTier + PositionCaps + RiskMultiplier for position limits.
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
	// PositionCaps constrains maximum leverage and position size relative to NAV.
	PositionCaps PositionCaps `json:"position_caps"`
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

// PositionCaps constrains risk_engine position sizing derived from this DecisionCard.
type PositionCaps struct {
	// MaxLeverage is the leverage ceiling (e.g. 2.0 = 2× notional).
	MaxLeverage float64 `json:"max_leverage"`
	// MaxPositionPct is the maximum single-position size as a fraction of NAV (e.g. 0.10 = 10%).
	MaxPositionPct float64 `json:"max_position_pct"`
}
