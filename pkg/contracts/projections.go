package contracts

// ponytail: keep legacy projection names as aliases; no wrapper types needed.
type (
	RegimeSnapshotEvent = RegimeSnapshot
	RegimeCardEvent     = RegimeCard
	DecisionCardEvent   = DecisionCard

	MarketRegimePort = MarketDataProvider
	MacroRegimePort  = MacroDataProvider
	RegimeEnginePort = DecisionCardProvider
)
