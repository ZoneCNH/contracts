package contracts

import (
	"encoding/json"
	"testing"
	"time"
)

// ---- RegimeSnapshot ----

func TestRegimeSnapshot_JSONRoundTrip(t *testing.T) {
	snap := RegimeSnapshot{
		SnapshotID:      "550e8400-e29b-41d4-a716-446655440000",
		Symbol:          "sym-001",
		RegimeState:     MarketStateS1Bull,
		Bias:            BiasLong,
		TradePermission: TradePermissionAllowNormal,
		Confidence:      0.87,
		FiveDimScores: FiveDimScores{
			Trend:      0.75,
			Leverage:   0.60,
			Heat:       0.45,
			Deleverage: 0.20,
			Volatility: 0.35,
		},
		FreshnessState: FreshnessHealthy,
		RiskTags:       []string{"high_oi_divergence"},
		EventTime:      1718870400000,
		ObservedAt:     1718870400500,
		PITVintage:     "2026-06-20T00:00:00Z",
	}

	b, err := json.Marshal(snap)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got RegimeSnapshot
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got.SnapshotID != snap.SnapshotID {
		t.Errorf("SnapshotID: got %q want %q", got.SnapshotID, snap.SnapshotID)
	}
	if got.RegimeState != MarketStateS1Bull {
		t.Errorf("RegimeState: got %q want S1_BULL", got.RegimeState)
	}
	if got.TradePermission != TradePermissionAllowNormal {
		t.Errorf("TradePermission: got %q want ALLOW_NORMAL", got.TradePermission)
	}
	if got.Confidence != 0.87 {
		t.Errorf("Confidence: got %v want 0.87", got.Confidence)
	}
	if len(got.RiskTags) != 1 || got.RiskTags[0] != "high_oi_divergence" {
		t.Errorf("RiskTags: got %v", got.RiskTags)
	}
}

func TestMarketState_AllValues(t *testing.T) {
	states := []MarketState{
		MarketStateS1Bull, MarketStateS2ShortSqueeze, MarketStateS3Bear,
		MarketStateS4Crash, MarketStateS5Choppy, MarketStateS6LowVol,
		MarketStateS7Compression, MarketStateUnknown, MarketStateDislocated,
	}
	if len(states) != 9 {
		t.Errorf("expected 9 MarketState values, got %d", len(states))
	}
}

func TestTradePermission_AllValues(t *testing.T) {
	perms := []TradePermission{
		TradePermissionAllowNormal, TradePermissionAllowReduced,
		TradePermissionReduceOnly, TradePermissionForbidOpen,
	}
	if len(perms) != 4 {
		t.Errorf("expected 4 TradePermission values, got %d", len(perms))
	}
}

// ---- RegimeCard ----

func TestRegimeCard_JSONRoundTrip(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	card := RegimeCard{
		MState: MacroStateM1LiqBull,
		LGIP: LGIPScore{
			Liquidity: 0.72,
			Growth:    0.65,
			Inflation: -0.30,
			Pressure:  -0.15,
		},
		Confidence:    0.91,
		Timestamp:     now,
		DataFreshness: now.Add(-2 * time.Hour),
		InputSource:   "macro_data_py/v2",
	}

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got RegimeCard
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got.MState != MacroStateM1LiqBull {
		t.Errorf("MState: got %q want M1", got.MState)
	}
	if got.LGIP.Liquidity != 0.72 {
		t.Errorf("LGIP.Liquidity: got %v want 0.72", got.LGIP.Liquidity)
	}
	if got.InputSource != "macro_data_py/v2" {
		t.Errorf("InputSource: got %q want macro_data_py/v2", got.InputSource)
	}
}

func TestMacroState_AllValues(t *testing.T) {
	states := []MacroState{
		MacroStateM0Unknown, MacroStateM1LiqBull, MacroStateM2Recovery,
		MacroStateM3SoftLanding, MacroStateM4Hawkish, MacroStateM5RecessionCut,
		MacroStateM6CreditDeleverage, MacroStateM7Stagflation,
	}
	if len(states) != 8 {
		t.Errorf("expected 8 MacroState values (M0+M1-M7), got %d", len(states))
	}
}

// ---- DecisionCard ----

func TestDecisionCard_JSONRoundTrip(t *testing.T) {
	dc := DecisionCard{
		CardID:      "card-001",
		GeneratedAt: 1718870400000,
		MarketState: MarketStateS1Bull,
		MacroState:  MacroStateM1LiqBull,
		Action:      ActionA,
		Profile:     ProfileAggressive,
		Template:    TemplateTrendFollowing,
		RiskTier:    5,
		ExposureCaps: ExposureCaps{
			MaxLeverage:    3.0,
			MaxExposurePct: 0.20,
		},
		RiskMultiplier: 1.0,
		Conflict:       false,
		Explain:        "M1(bull macro) × S1(bull market) → A aggressive/trend_following",
	}

	b, err := json.Marshal(dc)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got DecisionCard
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got.Action != ActionA {
		t.Errorf("Action: got %q want A", got.Action)
	}
	if got.RiskTier != 5 {
		t.Errorf("RiskTier: got %d want 5", got.RiskTier)
	}
	if got.ExposureCaps.MaxLeverage != 3.0 {
		t.Errorf("MaxLeverage: got %v want 3.0", got.ExposureCaps.MaxLeverage)
	}
	if got.Conflict {
		t.Error("Conflict: expected false")
	}
}

func TestDecisionCard_LegacyPositionCapsMarshalAndUnmarshal(t *testing.T) {
	dc := DecisionCard{
		CardID:         "card-legacy",
		GeneratedAt:    1718870400001,
		MarketState:    MarketStateS2ShortSqueeze,
		MacroState:     MacroStateM2Recovery,
		Action:         ActionB,
		Profile:        ProfileModerate,
		Template:       TemplateBreakout,
		RiskTier:       4,
		PositionCaps:   ExposureCaps{MaxLeverage: 1.75, MaxExposurePct: 0.12},
		RiskMultiplier: 0.8,
		Conflict:       false,
		Explain:        "legacy caller still populates PositionCaps",
	}

	b, err := json.Marshal(dc)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got DecisionCard
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got.ExposureCaps != dc.PositionCaps {
		t.Errorf("ExposureCaps: got %+v want %+v", got.ExposureCaps, dc.PositionCaps)
	}
	if got.PositionCaps != dc.PositionCaps {
		t.Errorf("PositionCaps: got %+v want %+v", got.PositionCaps, dc.PositionCaps)
	}
}

func TestDecisionCard_UnmarshalJSONError(t *testing.T) {
	var got DecisionCard
	if err := json.Unmarshal([]byte(`{"card_id":{}}`), &got); err == nil {
		t.Fatal("expected type-mismatched JSON to fail")
	}
}

func TestDecisionCard_Conflict(t *testing.T) {
	// M6 (credit deleveraging) × S1 (bull) = conflict scenario
	dc := DecisionCard{
		CardID:         "card-conflict",
		GeneratedAt:    1718870400000,
		MarketState:    MarketStateS1Bull,
		MacroState:     MacroStateM6CreditDeleverage,
		Action:         ActionD,
		Profile:        ProfileDefensive,
		Template:       TemplateHedge,
		RiskTier:       1,
		ExposureCaps:   ExposureCaps{MaxLeverage: 0.5, MaxExposurePct: 0.05},
		RiskMultiplier: 0.3,
		Conflict:       true,
		Explain:        "M6(credit deleveraging) conflicts with S1(bull); forced defensive",
	}

	if !dc.Conflict {
		t.Error("expected Conflict=true for M6×S1")
	}
	if dc.Action != ActionD {
		t.Errorf("expected defensive action D for conflict, got %q", dc.Action)
	}
}

func TestAction_AllValues(t *testing.T) {
	actions := []Action{ActionA, ActionB, ActionC, ActionD, ActionE}
	if len(actions) != 5 {
		t.Errorf("expected 5 Action values (A-E), got %d", len(actions))
	}
}

func TestStrategyTemplate_AllValues(t *testing.T) {
	templates := []StrategyTemplate{
		TemplateTrendFollowing, TemplateRangeTrading,
		TemplateBreakout, TemplateHedge, TemplateCash,
	}
	if len(templates) != 5 {
		t.Errorf("expected 5 StrategyTemplate values, got %d", len(templates))
	}
}
