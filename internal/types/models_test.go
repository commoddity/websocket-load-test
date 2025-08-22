package types

import (
	"testing"
	"time"
)

func TestConfig_Validation(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		valid  bool
	}{
		{
			name: "valid config",
			config: Config{
				URL:           "wss://ethereum.rpc.grove.city/v1/app123",
				ServiceID:     "ethereum",
				AuthHeader:    "api_key_123",
				Subscriptions: "newHeads",
				SubCount:      1,
			},
			valid: true,
		},
		{
			name: "empty URL",
			config: Config{
				URL:           "",
				ServiceID:     "ethereum",
				AuthHeader:    "api_key_123",
				Subscriptions: "newHeads",
				SubCount:      1,
			},
			valid: false,
		},
		{
			name: "zero sub count",
			config: Config{
				URL:           "wss://ethereum.rpc.grove.city/v1/app123",
				ServiceID:     "ethereum",
				AuthHeader:    "api_key_123",
				Subscriptions: "newHeads",
				SubCount:      0,
			},
			valid: false,
		},
		{
			name: "negative sub count",
			config: Config{
				URL:           "wss://ethereum.rpc.grove.city/v1/app123",
				ServiceID:     "ethereum",
				AuthHeader:    "api_key_123",
				Subscriptions: "newHeads",
				SubCount:      -1,
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := tt.config.URL != "" && tt.config.SubCount > 0
			if valid != tt.valid {
				t.Errorf("Config validation = %v, want %v", valid, tt.valid)
			}
		})
	}
}

func TestStats_Initialization(t *testing.T) {
	stats := &Stats{
		ClientStartTime: time.Now(),
	}

	if stats.TotalConnections != 0 {
		t.Errorf("TotalConnections = %v, want 0", stats.TotalConnections)
	}
	if stats.EventsReceived != 0 {
		t.Errorf("EventsReceived = %v, want 0", stats.EventsReceived)
	}
	if stats.ClientStartTime.IsZero() {
		t.Error("ClientStartTime should not be zero")
	}
}

func TestConnectionHistory_Validation(t *testing.T) {
	tests := []struct {
		name    string
		history ConnectionHistory
		valid   bool
	}{
		{
			name: "valid history",
			history: ConnectionHistory{
				ConnectionNum: 1,
				StartTime:     time.Now().Add(-time.Hour),
				EndTime:       time.Now(),
				Duration:      time.Hour,
				Messages:      100,
			},
			valid: true,
		},
		{
			name: "negative connection number",
			history: ConnectionHistory{
				ConnectionNum: -1,
				StartTime:     time.Now().Add(-time.Hour),
				EndTime:       time.Now(),
				Duration:      time.Hour,
				Messages:      100,
			},
			valid: false,
		},
		{
			name: "end time before start time",
			history: ConnectionHistory{
				ConnectionNum: 1,
				StartTime:     time.Now(),
				EndTime:       time.Now().Add(-time.Hour),
				Duration:      time.Hour,
				Messages:      100,
			},
			valid: false,
		},
		{
			name: "negative message count",
			history: ConnectionHistory{
				ConnectionNum: 1,
				StartTime:     time.Now().Add(-time.Hour),
				EndTime:       time.Now(),
				Duration:      time.Hour,
				Messages:      -1,
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := tt.history.ConnectionNum > 0 &&
				tt.history.EndTime.After(tt.history.StartTime) &&
				tt.history.Messages >= 0
			if valid != tt.valid {
				t.Errorf("ConnectionHistory validation = %v, want %v", valid, tt.valid)
			}
		})
	}
}

func TestJSONRPCRequest_Structure(t *testing.T) {
	tests := []struct {
		name    string
		request JSONRPCRequest
		wantID  int
	}{
		{
			name: "subscription request",
			request: JSONRPCRequest{
				JSONRPC: "2.0",
				ID:      1,
				Method:  "eth_subscribe",
				Params:  []string{"newHeads"},
			},
			wantID: 1,
		},
		{
			name: "unsubscribe request",
			request: JSONRPCRequest{
				JSONRPC: "2.0",
				ID:      2,
				Method:  "eth_unsubscribe",
				Params:  "0x123",
			},
			wantID: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.request.ID != tt.wantID {
				t.Errorf("JSONRPCRequest.ID = %v, want %v", tt.request.ID, tt.wantID)
			}
			if tt.request.JSONRPC != "2.0" {
				t.Errorf("JSONRPCRequest.JSONRPC = %v, want 2.0", tt.request.JSONRPC)
			}
		})
	}
}

func TestJSONRPCResponse_Structure(t *testing.T) {
	tests := []struct {
		name     string
		response JSONRPCResponse
		wantType string
	}{
		{
			name: "subscription confirmation",
			response: JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      1,
				Result:  "0x123abc",
			},
			wantType: "result",
		},
		{
			name: "subscription event",
			response: JSONRPCResponse{
				JSONRPC: "2.0",
				Method:  "eth_subscription",
				Params: map[string]interface{}{
					"subscription": "0x123abc",
					"result":       map[string]interface{}{},
				},
			},
			wantType: "subscription",
		},
		{
			name: "error response",
			response: JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      1,
				Error: map[string]interface{}{
					"code":    -32602,
					"message": "Invalid params",
				},
			},
			wantType: "error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var responseType string
			if tt.response.Result != nil {
				responseType = "result"
			} else if tt.response.Method == "eth_subscription" {
				responseType = "subscription"
			} else if tt.response.Error != nil {
				responseType = "error"
			}

			if responseType != tt.wantType {
				t.Errorf("Response type = %v, want %v", responseType, tt.wantType)
			}
		})
	}
}
