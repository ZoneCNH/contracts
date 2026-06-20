package contracts

import "time"

// RegimeCard is the P0 output contract of macro_engine → regime_engine.
// It captures the current macroeconomic regime state classified via the LGIP framework.
//
// Consumer:
//   - regime_engine: combined with RegimeSnapshot to produce DecisionCard via M×S matrix.
type RegimeCard struct {
	// MState is the classified macroeconomic regime (M1–M7 or M0_UNKNOWN).
	MState MacroState `json:"m_state"`
	// LGIP holds the four-factor scores (Liquidity, Growth, Inflation, Pressure).
	LGIP LGIPScore `json:"lgip"`
	// Confidence is the classification model's confidence in [0, 1].
	Confidence float64 `json:"confidence"`
	// Timestamp is the time at which this card was computed.
	Timestamp time.Time `json:"timestamp"`
	// DataFreshness is the staleness boundary of the underlying macro input data.
	DataFreshness time.Time `json:"data_freshness"`
	// InputSource identifies the macro data product that produced this card (e.g. "macro_data_py/v2").
	InputSource string `json:"input_source"`
}

// MacroState classifies the current macroeconomic regime using the LGIP four-factor framework.
//
// Historical validation anchors:
//   - COVID crash 2020: M6 → M5
//   - Rate-hike cycle 2022: M4
//   - Recovery 2023: M2 → M3
type MacroState string

const (
	// MacroStateM0Unknown: insufficient data to classify the macro regime.
	MacroStateM0Unknown MacroState = "M0_UNKNOWN"
	// MacroStateM1LiqBull: 流动牛市 — L↑ G↑ I↓ (ample liquidity, rising growth, low inflation).
	MacroStateM1LiqBull MacroState = "M1"
	// MacroStateM2Recovery: 再通复苏 — L↓ G↑ I↓ (tightening liquidity, recovering growth).
	MacroStateM2Recovery MacroState = "M2"
	// MacroStateM3SoftLanding: 软着繁荣 — L→ G↑ I→ (neutral liquidity, strong growth, stable inflation).
	MacroStateM3SoftLanding MacroState = "M3"
	// MacroStateM4Hawkish: 鹰派通胀 — L↓ G↓ I↑ (tightening, slowing growth, rising inflation).
	MacroStateM4Hawkish MacroState = "M4"
	// MacroStateM5RecessionCut: 衰退降息 — L↑ G↓ I↓ (easing, falling growth and inflation).
	MacroStateM5RecessionCut MacroState = "M5"
	// MacroStateM6CreditDeleverage: 信用去杠 — L↓↓ G↓ P↑↑ (liquidity crisis, systemic stress).
	MacroStateM6CreditDeleverage MacroState = "M6"
	// MacroStateM7Stagflation: 滞胀冲击 — L→ G↓ I↑ (stagnation with persistent inflation).
	MacroStateM7Stagflation MacroState = "M7"
)

// LGIPScore holds the four-factor macro scoring vector.
// Each score is signed: positive = expansionary, negative = contractionary.
// Absolute magnitude indicates the strength of the signal.
type LGIPScore struct {
	// Liquidity (L): derived from M2 growth, central bank balance sheet, and SOFR.
	Liquidity float64 `json:"liquidity"`
	// Growth (G): derived from GDP, PMI, and employment data.
	Growth float64 `json:"growth"`
	// Inflation (I): derived from CPI, PCE, and PPI.
	Inflation float64 `json:"inflation"`
	// Pressure (P): derived from VIX, credit spreads, and DXY.
	Pressure float64 `json:"pressure"`
}
