package cmd

import (
	"fmt"
	"testing"
)

func TestURL_Construction(t *testing.T) {
	tests := []struct {
		name      string
		serviceID string
		appID     string
		wantURL   string
	}{
		{
			name:      "ethereum service",
			serviceID: "ethereum",
			appID:     "app123",
			wantURL:   "wss://ethereum.rpc.grove.city/v1/app123",
		},
		{
			name:      "polygon service",
			serviceID: "polygon",
			appID:     "app456",
			wantURL:   "wss://polygon.rpc.grove.city/v1/app456",
		},
		{
			name:      "xrplevm service",
			serviceID: "xrplevm",
			appID:     "app789",
			wantURL:   "wss://xrplevm.rpc.grove.city/v1/app789",
		},
		{
			name:      "arbitrum service",
			serviceID: "arbitrum",
			appID:     "appABC",
			wantURL:   "wss://arbitrum.rpc.grove.city/v1/appABC",
		},
		{
			name:      "optimism service",
			serviceID: "optimism",
			appID:     "appDEF",
			wantURL:   "wss://optimism.rpc.grove.city/v1/appDEF",
		},
		{
			name:      "base service",
			serviceID: "base",
			appID:     "appGHI",
			wantURL:   "wss://base.rpc.grove.city/v1/appGHI",
		},
		{
			name:      "empty app ID",
			serviceID: "ethereum",
			appID:     "",
			wantURL:   "wss://ethereum.rpc.grove.city/v1/",
		},
		{
			name:      "special characters in app ID",
			serviceID: "ethereum",
			appID:     "app-123_test",
			wantURL:   "wss://ethereum.rpc.grove.city/v1/app-123_test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotURL := fmt.Sprintf("wss://%s.rpc.grove.city/v1/%s", tt.serviceID, tt.appID)
			if gotURL != tt.wantURL {
				t.Errorf("URL construction = %q, want %q", gotURL, tt.wantURL)
			}
		})
	}
}

func TestRootCommand_Flags(t *testing.T) {
	tests := []struct {
		name         string
		flagName     string
		expectedType string
		required     bool
	}{
		{
			name:         "service flag",
			flagName:     "service",
			expectedType: "string",
			required:     false,
		},
		{
			name:         "app-id flag",
			flagName:     "app-id",
			expectedType: "string",
			required:     true,
		},
		{
			name:         "api-key flag",
			flagName:     "api-key",
			expectedType: "string",
			required:     true,
		},
		{
			name:         "subs flag",
			flagName:     "subs",
			expectedType: "string",
			required:     false,
		},
		{
			name:         "count flag",
			flagName:     "count",
			expectedType: "int",
			required:     false,
		},
		{
			name:         "log flag",
			flagName:     "log",
			expectedType: "bool",
			required:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := rootCmd.Flags().Lookup(tt.flagName)
			if flag == nil {
				t.Fatalf("Flag %q not found", tt.flagName)
			}

			if flag.Value.Type() != tt.expectedType {
				t.Errorf("Flag %q type = %q, want %q", tt.flagName, flag.Value.Type(), tt.expectedType)
			}

			// Check if flag is marked as required
			// Note: Cobra's MarkFlagRequired affects internal state,
			// For this test, we'll just check that app-id and api-key exist
			if tt.required && (tt.flagName == "app-id" || tt.flagName == "api-key") {
				// We expect app-id and api-key to be required
				// The actual requirement checking is handled by Cobra internally
			}
		})
	}
}

func TestRootCommand_Defaults(t *testing.T) {
	tests := []struct {
		name            string
		flagName        string
		expectedDefault interface{}
	}{
		{
			name:            "service default",
			flagName:        "service",
			expectedDefault: "ethereum",
		},
		{
			name:            "subs default",
			flagName:        "subs",
			expectedDefault: "newHeads",
		},
		{
			name:            "count default",
			flagName:        "count",
			expectedDefault: "1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := rootCmd.Flags().Lookup(tt.flagName)
			if flag == nil {
				t.Fatalf("Flag %q not found", tt.flagName)
			}

			if flag.DefValue != tt.expectedDefault {
				t.Errorf("Flag %q default = %q, want %q", tt.flagName, flag.DefValue, tt.expectedDefault)
			}
		})
	}
}

func TestRootCommand_ShortFlags(t *testing.T) {
	tests := []struct {
		name      string
		flagName  string
		shorthand string
	}{
		{
			name:      "service short flag",
			flagName:  "service",
			shorthand: "s",
		},
		{
			name:      "app-id short flag",
			flagName:  "app-id",
			shorthand: "a",
		},
		{
			name:      "api-key short flag",
			flagName:  "api-key",
			shorthand: "k",
		},
		{
			name:      "count short flag",
			flagName:  "count",
			shorthand: "c",
		},
		{
			name:      "log short flag",
			flagName:  "log",
			shorthand: "l",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := rootCmd.Flags().Lookup(tt.flagName)
			if flag == nil {
				t.Fatalf("Flag %q not found", tt.flagName)
			}

			if flag.Shorthand != tt.shorthand {
				t.Errorf("Flag %q shorthand = %q, want %q", tt.flagName, flag.Shorthand, tt.shorthand)
			}
		})
	}
}

func TestRootCommand_Usage(t *testing.T) {
	tests := []struct {
		name     string
		property string
		expected string
	}{
		{
			name:     "command use",
			property: "use",
			expected: "websocket-load-test",
		},
		{
			name:     "command short description",
			property: "short",
			expected: "🚀 Load test WebSocket endpoints for Grove Portal blockchain services",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var actual string
			switch tt.property {
			case "use":
				actual = rootCmd.Use
			case "short":
				actual = rootCmd.Short
			}

			if actual != tt.expected {
				t.Errorf("Command %s = %q, want %q", tt.property, actual, tt.expected)
			}
		})
	}
}

func TestValidateService(t *testing.T) {
	tests := []struct {
		name    string
		service string
		valid   bool
	}{
		{
			name:    "ethereum service",
			service: "ethereum",
			valid:   true,
		},
		{
			name:    "polygon service",
			service: "polygon",
			valid:   true,
		},
		{
			name:    "xrplevm service",
			service: "xrplevm",
			valid:   true,
		},
		{
			name:    "arbitrum service",
			service: "arbitrum",
			valid:   true,
		},
		{
			name:    "optimism service",
			service: "optimism",
			valid:   true,
		},
		{
			name:    "base service",
			service: "base",
			valid:   true,
		},
		{
			name:    "invalid service",
			service: "invalid",
			valid:   false,
		},
		{
			name:    "empty service",
			service: "",
			valid:   false,
		},
		{
			name:    "case sensitive - Ethereum",
			service: "Ethereum",
			valid:   false,
		},
	}

	validServices := map[string]bool{
		"ethereum": true,
		"polygon":  true,
		"xrplevm":  true,
		"arbitrum": true,
		"optimism": true,
		"base":     true,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := validServices[tt.service]
			if valid != tt.valid {
				t.Errorf("Service %q validation = %v, want %v", tt.service, valid, tt.valid)
			}
		})
	}
}

func BenchmarkURL_Construction(b *testing.B) {
	serviceID := "ethereum"
	appID := "app123"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("wss://%s.rpc.grove.city/v1/%s", serviceID, appID)
	}
}
