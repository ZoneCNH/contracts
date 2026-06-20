package contracts_test

import (
	"encoding/json"
	"testing"

	"github.com/ZoneCNH/contracts/pkg/contracts"
)

func TestSignalIntent_JSONRoundTrip(t *testing.T) {
	intent := contracts.SignalIntent{
		ID:          "uuid-001",
		GeneratedAt: 1718900000000,
		CardID:      "card-001",
		Symbol:      "BTCUSDT",
		Action:      contracts.ActionA,
		Template:    contracts.TemplateTrendFollowing,
		Strength:    0.85,
		SizePct:     0.765,
		Conflict:    false,
		Explain:     "action=A template=trend_following strength=0.85 risk_mul=0.90",
	}

	b, err := json.Marshal(intent)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got contracts.SignalIntent
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got.ID != intent.ID {
		t.Errorf("ID: got %q want %q", got.ID, intent.ID)
	}
	if got.Symbol != intent.Symbol {
		t.Errorf("Symbol: got %q want %q", got.Symbol, intent.Symbol)
	}
	if got.Action != intent.Action {
		t.Errorf("Action: got %q want %q", got.Action, intent.Action)
	}
	if got.Template != intent.Template {
		t.Errorf("Template: got %q want %q", got.Template, intent.Template)
	}
	if got.Strength != intent.Strength {
		t.Errorf("Strength: got %f want %f", got.Strength, intent.Strength)
	}
}

func TestSignalFactoryProvider_Interface(t *testing.T) {
	// 编译期断言：contracts.SignalFactoryProvider 接口可声明
	var _ contracts.SignalFactoryProvider = (*mockFactory)(nil)
}

type mockFactory struct{}

func (m *mockFactory) Generate(card contracts.DecisionCard, symbols []string) ([]contracts.SignalIntent, error) {
	return nil, nil
}
