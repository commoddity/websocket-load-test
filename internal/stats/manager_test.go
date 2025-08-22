package stats

import (
	"testing"
	"time"

	"github.com/commoddity/websocket-load-test/internal/types"
)

func TestNewManager(t *testing.T) {
	manager := NewManager()

	if manager == nil {
		t.Fatal("NewManager() returned nil")
	}

	stats := manager.GetStats()
	if stats == nil {
		t.Fatal("GetStats() returned nil")
	}

	if stats.ClientStartTime.IsZero() {
		t.Error("ClientStartTime should be set")
	}

	if stats.TotalConnections != 0 {
		t.Errorf("TotalConnections = %d, want 0", stats.TotalConnections)
	}
}

func TestManager_IncrementConnectionAttempts(t *testing.T) {
	tests := []struct {
		name       string
		increments int
		want       int
	}{
		{
			name:       "single increment",
			increments: 1,
			want:       1,
		},
		{
			name:       "multiple increments",
			increments: 5,
			want:       5,
		},
		{
			name:       "zero increments",
			increments: 0,
			want:       0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewManager()
			for i := 0; i < tt.increments; i++ {
				manager.IncrementConnectionAttempts()
			}

			if manager.GetStats().ConnectionAttempts != tt.want {
				t.Errorf("ConnectionAttempts = %d, want %d", manager.GetStats().ConnectionAttempts, tt.want)
			}
		})
	}
}

func TestManager_StartNewConnection(t *testing.T) {
	manager := NewManager()

	// Record time before starting connection
	beforeStart := time.Now()

	manager.StartNewConnection()

	stats := manager.GetStats()

	if stats.TotalConnections != 1 {
		t.Errorf("TotalConnections = %d, want 1", stats.TotalConnections)
	}

	if stats.CurrentConnMessages != 0 {
		t.Errorf("CurrentConnMessages = %d, want 0", stats.CurrentConnMessages)
	}

	if stats.CurrentConnStart.Before(beforeStart) {
		t.Error("CurrentConnStart should be set to recent time")
	}
}

func TestManager_HandleResponse(t *testing.T) {
	tests := []struct {
		name                   string
		response               types.JSONRPCResponse
		wantSubscriptionEvents int
		wantConfirmationEvents int
		wantErrorEvents        int
	}{
		{
			name: "subscription event",
			response: types.JSONRPCResponse{
				Method: "eth_subscription",
				Params: map[string]interface{}{
					"subscription": "0x123",
					"result": map[string]interface{}{
						"number": "0x1",
					},
				},
			},
			wantSubscriptionEvents: 1,
			wantConfirmationEvents: 0,
			wantErrorEvents:        0,
		},
		{
			name: "confirmation response",
			response: types.JSONRPCResponse{
				ID:     float64(1),
				Result: "0x123abc",
			},
			wantSubscriptionEvents: 0,
			wantConfirmationEvents: 1,
			wantErrorEvents:        0,
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
			wantSubscriptionEvents: 0,
			wantConfirmationEvents: 0,
			wantErrorEvents:        1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewManager()
			manager.HandleResponse(tt.response)

			stats := manager.GetStats()

			if stats.SubscriptionEvents != tt.wantSubscriptionEvents {
				t.Errorf("SubscriptionEvents = %d, want %d", stats.SubscriptionEvents, tt.wantSubscriptionEvents)
			}

			if stats.ConfirmationEvents != tt.wantConfirmationEvents {
				t.Errorf("ConfirmationEvents = %d, want %d", stats.ConfirmationEvents, tt.wantConfirmationEvents)
			}

			if stats.ErrorEvents != tt.wantErrorEvents {
				t.Errorf("ErrorEvents = %d, want %d", stats.ErrorEvents, tt.wantErrorEvents)
			}

			if stats.EventsReceived != 1 {
				t.Errorf("EventsReceived = %d, want 1", stats.EventsReceived)
			}

			if stats.CurrentConnMessages != 1 {
				t.Errorf("CurrentConnMessages = %d, want 1", stats.CurrentConnMessages)
			}

			if stats.LastEventTime.IsZero() {
				t.Error("LastEventTime should be set")
			}
		})
	}
}

func TestManager_EndConnection(t *testing.T) {
	manager := NewManager()

	// Start a connection first
	manager.StartNewConnection()

	// Wait a bit to ensure duration > 0
	time.Sleep(1 * time.Millisecond)

	// End the connection
	manager.EndConnection()

	stats := manager.GetStats()

	if stats.TotalUptime == 0 {
		t.Error("TotalUptime should be greater than 0")
	}

	if stats.LongestConnection == 0 {
		t.Error("LongestConnection should be set")
	}

	if stats.ShortestConnection == 0 {
		t.Error("ShortestConnection should be set")
	}
}

func TestManager_SetSubscriptionMapping(t *testing.T) {
	tests := []struct {
		name             string
		subscriptionID   string
		subscriptionType string
	}{
		{
			name:             "newHeads mapping",
			subscriptionID:   "0x123abc",
			subscriptionType: "newHeads",
		},
		{
			name:             "logs mapping",
			subscriptionID:   "0x456def",
			subscriptionType: "logs",
		},
		{
			name:             "empty subscription ID",
			subscriptionID:   "",
			subscriptionType: "newHeads",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewManager()
			manager.SetSubscriptionMapping(tt.subscriptionID, tt.subscriptionType)

			// Test the mapping by checking if it's retrievable
			retrievedType := manager.getSubscriptionTypeFromID(tt.subscriptionID)
			if retrievedType != tt.subscriptionType {
				t.Errorf("getSubscriptionTypeFromID(%q) = %q, want %q",
					tt.subscriptionID, retrievedType, tt.subscriptionType)
			}
		})
	}
}

func TestManager_IncrementReconnections(t *testing.T) {
	tests := []struct {
		name                  string
		totalConnections      int
		expectedReconnections int
	}{
		{
			name:                  "no previous connections",
			totalConnections:      0,
			expectedReconnections: 0,
		},
		{
			name:                  "with previous connections",
			totalConnections:      1,
			expectedReconnections: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewManager()

			// Set up the total connections
			for i := 0; i < tt.totalConnections; i++ {
				manager.StartNewConnection()
			}

			manager.IncrementReconnections()

			if manager.GetStats().TotalReconnections != tt.expectedReconnections {
				t.Errorf("TotalReconnections = %d, want %d",
					manager.GetStats().TotalReconnections, tt.expectedReconnections)
			}
		})
	}
}

func BenchmarkHandleResponse(b *testing.B) {
	manager := NewManager()
	response := types.JSONRPCResponse{
		Method: "eth_subscription",
		Params: map[string]interface{}{
			"subscription": "0x123",
			"result": map[string]interface{}{
				"number": "0x1",
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.HandleResponse(response)
	}
}

func BenchmarkStartNewConnection(b *testing.B) {
	manager := NewManager()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.StartNewConnection()
	}
}
