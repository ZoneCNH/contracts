package contracts

// RegimeSnapshot is the P0 output contract of market_engine → regime_engine.
// It captures the current market state classification for a single symbol.
//
// PIT constraint: EventTime ≤ ObservedAt. Snapshots violating this MUST be rejected.
//
// Consumer:
//   - regime_engine: combines with RegimeCard to produce DecisionCard via M×S matrix.
type RegimeSnapshot struct {
	// SnapshotID is a UUID v4 uniquely identifying this snapshot.
	SnapshotID string `json:"snapshot_id"`
	// Symbol is the normalized instrument identifier (e.g. "BTCUSDT").
	Symbol string `json:"symbol"`
	// RegimeState is the classified market state (S1–S7, UNKNOWN, or DISLOCATED).
	RegimeState MarketState `json:"regime_state"`
	// Bias is the directional orientation of the current regime.
	Bias MarketBias `json:"bias"`
	// TradePermission is the trading gate derived from RegimeState.
	TradePermission TradePermission `json:"trade_permission"`
	// Confidence is the model's classification confidence in [0, 1] (DI-001).
	Confidence float64 `json:"confidence"`
	// FiveDimScores holds the five-dimensional scoring vector used to classify RegimeState.
	FiveDimScores FiveDimScores `json:"five_dim_scores"`
	// FreshnessState describes the quality of the underlying data feed.
	FreshnessState FreshnessState `json:"freshness_state"`
	// RiskTags lists active risk labels attached by the quality gate (e.g. "high_oi_divergence").
	RiskTags []string `json:"risk_tags,omitempty"`
	// EventTime is the data event timestamp in UTC milliseconds.
	// Must be ≤ ObservedAt (PIT constraint).
	EventTime int64 `json:"event_time"`
	// ObservedAt is the system processing timestamp in UTC milliseconds.
	ObservedAt int64 `json:"observed_at"`
	// PITVintage is the point-in-time snapshot version label for backtesting reproducibility.
	PITVintage string `json:"pit_vintage,omitempty"`
}

// MarketState classifies the current market regime as one of seven canonical states
// or two special states (UNKNOWN, DISLOCATED).
type MarketState string

const (
	// MarketStateS1Bull: 多头 — trending up, strong momentum.
	MarketStateS1Bull MarketState = "S1_BULL"
	// MarketStateS2ShortSqueeze: 挤空 — short-covering driven rally.
	MarketStateS2ShortSqueeze MarketState = "S2_SHORT_SQUEEZE"
	// MarketStateS3Bear: 空头 — trending down, sustained selling.
	MarketStateS3Bear MarketState = "S3_BEAR"
	// MarketStateS4Crash: 踩踏 — rapid deleveraging cascade.
	MarketStateS4Crash MarketState = "S4_CRASH"
	// MarketStateS5Choppy: 震荡 — directionless, mean-reverting.
	MarketStateS5Choppy MarketState = "S5_CHOPPY"
	// MarketStateS6LowVol: 低波 — compressed volatility, range-bound.
	MarketStateS6LowVol MarketState = "S6_LOW_VOL"
	// MarketStateS7Compression: 压缩 — pre-breakout contraction.
	MarketStateS7Compression MarketState = "S7_COMPRESSION"
	// MarketStateUnknown: insufficient data to classify.
	MarketStateUnknown MarketState = "UNKNOWN"
	// MarketStateDislocated: market structure broken; normal rules suspended.
	MarketStateDislocated MarketState = "DISLOCATED"
)

// MarketBias is the directional orientation of the market regime.
type MarketBias string

const (
	BiasLong    MarketBias = "LONG"
	BiasShort   MarketBias = "SHORT"
	BiasNeutral MarketBias = "NEUTRAL"
)

// TradePermission defines the set of allowed trading operations for this regime.
type TradePermission string

const (
	// TradePermissionAllowNormal: full open and close allowed.
	TradePermissionAllowNormal TradePermission = "ALLOW_NORMAL"
	// TradePermissionAllowReduced: open allowed at reduced size.
	TradePermissionAllowReduced TradePermission = "ALLOW_REDUCED"
	// TradePermissionReduceOnly: closing existing positions only.
	TradePermissionReduceOnly TradePermission = "REDUCE_ONLY"
	// TradePermissionForbidOpen: no new positions; risk system override.
	TradePermissionForbidOpen TradePermission = "FORBID_OPEN"
)

// FreshnessState describes the data feed quality of the upstream market data.
type FreshnessState string

const (
	FreshnessHealthy    FreshnessState = "healthy"
	FreshnessDegraded   FreshnessState = "degraded"
	FreshnessStale      FreshnessState = "stale"
	FreshnessRecovering FreshnessState = "recovering"
)

// FiveDimScores holds the five-dimensional market scoring vector used by market_engine.
// Each score is normalized to [0, 1]; weights are: Trend 30%, Leverage 25%,
// Heat 20%, Deleverage 15%, Volatility 10%.
type FiveDimScores struct {
	// Trend: price momentum and moving average structure.
	Trend float64 `json:"trend"`
	// Leverage: OI changes and funding rate.
	Leverage float64 `json:"leverage"`
	// Heat: volume deviation and sentiment indicators.
	Heat float64 `json:"heat"`
	// Deleverage: liquidation volume and cascade risk.
	Deleverage float64 `json:"deleverage"`
	// Volatility: realized and implied volatility.
	Volatility float64 `json:"volatility"`
}
