package terminal

import (
	"testing"
)

func TestGetSubscriptionEmoji(t *testing.T) {
	tests := []struct {
		name             string
		subscriptionType string
		wantEmoji        string
	}{
		{
			name:             "newHeads subscription",
			subscriptionType: "newHeads",
			wantEmoji:        "🧊",
		},
		{
			name:             "newPendingTransactions subscription",
			subscriptionType: "newPendingTransactions",
			wantEmoji:        "⚡",
		},
		{
			name:             "logs subscription",
			subscriptionType: "logs",
			wantEmoji:        "📄",
		},
		{
			name:             "syncing subscription",
			subscriptionType: "syncing",
			wantEmoji:        "🔄",
		},
		{
			name:             "unknown subscription type",
			subscriptionType: "unknownType",
			wantEmoji:        "📡",
		},
		{
			name:             "empty string",
			subscriptionType: "",
			wantEmoji:        "📡",
		},
		{
			name:             "case sensitive - NewHeads",
			subscriptionType: "NewHeads",
			wantEmoji:        "📡",
		},
		{
			name:             "case sensitive - NEWHEADS",
			subscriptionType: "NEWHEADS",
			wantEmoji:        "📡",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetSubscriptionEmoji(tt.subscriptionType)
			if got != tt.wantEmoji {
				t.Errorf("GetSubscriptionEmoji(%q) = %q, want %q", tt.subscriptionType, got, tt.wantEmoji)
			}
		})
	}
}

func TestGetTerminalWidth(t *testing.T) {
	tests := []struct {
		name    string
		wantMin int
		wantMax int
	}{
		{
			name:    "terminal width should be reasonable",
			wantMin: 20,   // minimum expected width
			wantMax: 1000, // maximum reasonable width
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			width := GetTerminalWidth()
			if width < tt.wantMin || width > tt.wantMax {
				t.Errorf("GetTerminalWidth() = %d, want between %d and %d", width, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestColorVariables(t *testing.T) {
	tests := []struct {
		name  string
		color interface{}
	}{
		{"Green color exists", Green},
		{"Red color exists", Red},
		{"Yellow color exists", Yellow},
		{"Blue color exists", Blue},
		{"Magenta color exists", Magenta},
		{"Cyan color exists", Cyan},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.color == nil {
				t.Errorf("Color variable %s is nil", tt.name)
			}
		})
	}
}

func BenchmarkGetSubscriptionEmoji(b *testing.B) {
	subscriptionTypes := []string{
		"newHeads",
		"newPendingTransactions",
		"logs",
		"syncing",
		"unknownType",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, subType := range subscriptionTypes {
			GetSubscriptionEmoji(subType)
		}
	}
}

func BenchmarkGetTerminalWidth(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetTerminalWidth()
	}
}
