package contracts

import "context"

// MarketDataProvider is the P0 port interface exposed by market_data to downstream
// consumers (market_regime, regime_engine).
// Implementations live in market_data; this contract is stable across module versions.
type MarketDataProvider interface {
	// LatestRegimeSnapshot returns the most recent RegimeSnapshot for the given symbol.
	// Returns an error if no snapshot is available or the symbol is unknown.
	LatestRegimeSnapshot(ctx context.Context, symbol string) (RegimeSnapshot, error)

	// SubscribeRegimeSnapshots streams RegimeSnapshot updates for the given symbols.
	// The returned channel is closed when ctx is cancelled or an unrecoverable error occurs.
	SubscribeRegimeSnapshots(ctx context.Context, symbols []string) (<-chan RegimeSnapshot, error)
}

// MacroDataProvider is the P0 port interface exposed by macro_data to downstream
// consumers (macro_regime, regime_engine).
// Implementations live in macro_data; this contract is stable across module versions.
type MacroDataProvider interface {
	// LatestRegimeCard returns the current macro regime card.
	// Returns an error if no card has been computed yet.
	LatestRegimeCard(ctx context.Context) (RegimeCard, error)

	// SubscribeRegimeCards streams RegimeCard updates.
	// The returned channel is closed when ctx is cancelled or an unrecoverable error occurs.
	SubscribeRegimeCards(ctx context.Context) (<-chan RegimeCard, error)
}

// DecisionCardProvider is the P0 port interface exposed by regime_engine to downstream
// consumers (signal_factory, risk_engine).
type DecisionCardProvider interface {
	// LatestDecisionCard returns the most recent DecisionCard.
	// Returns an error if the engine has not yet produced a card.
	LatestDecisionCard(ctx context.Context) (DecisionCard, error)

	// SubscribeDecisionCards streams DecisionCard updates.
	// The returned channel is closed when ctx is cancelled or an unrecoverable error occurs.
	SubscribeDecisionCards(ctx context.Context) (<-chan DecisionCard, error)
}
