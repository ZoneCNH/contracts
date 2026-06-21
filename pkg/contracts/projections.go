package contracts

// ponytail: keep legacy projection names as aliases; no wrapper types needed.
type (
	RegimeSnapshotEvent = RegimeSnapshot
	RegimeCardEvent     = RegimeCard
	DecisionCardEvent   = DecisionCard

	MarketDataProviderPort = MarketDataProvider
	MacroDataProviderPort  = MacroDataProvider
	RegimeEnginePort       = DecisionCardProvider
)
