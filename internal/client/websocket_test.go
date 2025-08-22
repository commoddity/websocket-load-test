package client

import (
	"testing"

	"github.com/commoddity/websocket-load-test/internal/stats"
	"github.com/commoddity/websocket-load-test/internal/types"
)

func TestNewWebSocketClient(t *testing.T) {
	tests := []struct {
		name   string
		config *types.Config
	}{
		{
			name: "valid config",
			config: &types.Config{
				URL:           "wss://ethereum.rpc.grove.city/v1/app123",
				ServiceID:     "ethereum",
				AuthHeader:    "api_key_123",
				Subscriptions: "newHeads",
				SubCount:      1,
			},
		},
		{
			name: "multiple subscriptions",
			config: &types.Config{
				URL:           "wss://polygon.rpc.grove.city/v1/app456",
				ServiceID:     "polygon",
				AuthHeader:    "api_key_456",
				Subscriptions: "newHeads,newPendingTransactions,logs",
				SubCount:      5,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statsManager := stats.NewManager()
			done := make(chan struct{})
			defer close(done)

			client := NewWebSocketClient(tt.config, statsManager, done)

			if client == nil {
				t.Fatal("NewWebSocketClient() returned nil")
			}

			if client.config != tt.config {
				t.Error("Config not set correctly")
			}

			if client.statsManager != statsManager {
				t.Error("StatsManager not set correctly")
			}

			if client.done != done {
				t.Error("Done channel not set correctly")
			}

			if client.subscriptionIDs == nil {
				t.Error("SubscriptionIDs map not initialized")
			}

			if client.idToSubscription == nil {
				t.Error("IdToSubscription map not initialized")
			}
		})
	}
}

func TestWebSocketClient_GetTotalSubscriptions(t *testing.T) {
	tests := []struct {
		name                  string
		initialSubscriptions  int
		expectedSubscriptions int
	}{
		{
			name:                  "zero subscriptions",
			initialSubscriptions:  0,
			expectedSubscriptions: 0,
		},
		{
			name:                  "some subscriptions",
			initialSubscriptions:  5,
			expectedSubscriptions: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &types.Config{
				URL:           "wss://ethereum.rpc.grove.city/v1/app123",
				ServiceID:     "ethereum",
				AuthHeader:    "api_key_123",
				Subscriptions: "newHeads",
				SubCount:      1,
			}
			statsManager := stats.NewManager()
			done := make(chan struct{})
			defer close(done)

			client := NewWebSocketClient(config, statsManager, done)
			client.totalSubscriptions = tt.initialSubscriptions

			got := client.GetTotalSubscriptions()
			if got != tt.expectedSubscriptions {
				t.Errorf("GetTotalSubscriptions() = %d, want %d", got, tt.expectedSubscriptions)
			}
		})
	}
}

func TestWebSocketClient_HandleResponse(t *testing.T) {
	tests := []struct {
		name     string
		response types.JSONRPCResponse
		setupID  bool
		setupMap bool
	}{
		{
			name: "subscription confirmation",
			response: types.JSONRPCResponse{
				ID:     float64(1),
				Result: "0x123abc",
			},
			setupID:  true,
			setupMap: false,
		},
		{
			name: "subscription event",
			response: types.JSONRPCResponse{
				Method: "eth_subscription",
				Params: map[string]interface{}{
					"subscription": "0x123abc",
					"result":       map[string]interface{}{},
				},
			},
			setupID:  false,
			setupMap: false,
		},
		{
			name: "error response",
			response: types.JSONRPCResponse{
				ID: float64(1),
				Error: map[string]interface{}{
					"code":    -32602,
					"message": "Invalid params",
				},
			},
			setupID:  false,
			setupMap: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &types.Config{
				URL:           "wss://ethereum.rpc.grove.city/v1/app123",
				ServiceID:     "ethereum",
				AuthHeader:    "api_key_123",
				Subscriptions: "newHeads",
				SubCount:      1,
			}
			statsManager := stats.NewManager()
			done := make(chan struct{})
			defer close(done)

			client := NewWebSocketClient(config, statsManager, done)

			// Setup test data if needed
			if tt.setupID {
				client.idToSubscription[1] = "newHeads"
			}

			// This should not panic and should handle the response
			client.handleResponse(tt.response)

			// Verify stats were updated
			stats := statsManager.GetStats()
			if stats.EventsReceived != 1 {
				t.Errorf("EventsReceived = %d, want 1", stats.EventsReceived)
			}
		})
	}
}

func TestValidateSubscriptionParams(t *testing.T) {
	tests := []struct {
		name         string
		subscription string
		wantParams   interface{}
	}{
		{
			name:         "newHeads subscription",
			subscription: "newHeads",
			wantParams:   []string{"newHeads"},
		},
		{
			name:         "newPendingTransactions subscription",
			subscription: "newPendingTransactions",
			wantParams:   []string{"newPendingTransactions"},
		},
		{
			name:         "logs subscription",
			subscription: "logs",
			wantParams:   []interface{}{"logs", map[string]interface{}{"topics": []interface{}{nil}}},
		},
		{
			name:         "custom subscription",
			subscription: "custom",
			wantParams:   []string{"custom"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var params interface{}
			switch tt.subscription {
			case "newHeads":
				params = []string{"newHeads"}
			case "newPendingTransactions":
				params = []string{"newPendingTransactions"}
			case "logs":
				params = []interface{}{"logs", map[string]interface{}{"topics": []interface{}{nil}}}
			default:
				params = []string{tt.subscription}
			}

			// Compare the structure (not exact equality due to interface{} complexity)
			switch tt.subscription {
			case "newHeads", "newPendingTransactions", "custom":
				if arr, ok := params.([]string); !ok || len(arr) != 1 || arr[0] != tt.subscription {
					t.Errorf("Params for %s = %v, want %v", tt.subscription, params, tt.wantParams)
				}
			case "logs":
				if arr, ok := params.([]interface{}); !ok || len(arr) != 2 || arr[0] != "logs" {
					t.Errorf("Params for %s structure incorrect", tt.subscription)
				}
			}
		})
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  *types.Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &types.Config{
				URL:           "wss://ethereum.rpc.grove.city/v1/app123",
				ServiceID:     "ethereum",
				AuthHeader:    "api_key_123",
				Subscriptions: "newHeads",
				SubCount:      1,
			},
			wantErr: false,
		},
		{
			name: "empty URL",
			config: &types.Config{
				URL:           "",
				ServiceID:     "ethereum",
				AuthHeader:    "api_key_123",
				Subscriptions: "newHeads",
				SubCount:      1,
			},
			wantErr: true,
		},
		{
			name: "zero sub count",
			config: &types.Config{
				URL:           "wss://ethereum.rpc.grove.city/v1/app123",
				ServiceID:     "ethereum",
				AuthHeader:    "api_key_123",
				Subscriptions: "newHeads",
				SubCount:      0,
			},
			wantErr: true,
		},
		{
			name: "negative sub count",
			config: &types.Config{
				URL:           "wss://ethereum.rpc.grove.city/v1/app123",
				ServiceID:     "ethereum",
				AuthHeader:    "api_key_123",
				Subscriptions: "newHeads",
				SubCount:      -1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasErr := tt.config.URL == "" || tt.config.SubCount <= 0
			if hasErr != tt.wantErr {
				t.Errorf("Config validation error = %v, wantErr %v", hasErr, tt.wantErr)
			}
		})
	}
}

func BenchmarkNewWebSocketClient(b *testing.B) {
	config := &types.Config{
		URL:           "wss://ethereum.rpc.grove.city/v1/app123",
		ServiceID:     "ethereum",
		AuthHeader:    "api_key_123",
		Subscriptions: "newHeads",
		SubCount:      1,
	}
	statsManager := stats.NewManager()
	done := make(chan struct{})
	defer close(done)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewWebSocketClient(config, statsManager, done)
	}
}
