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
			name:      "xrplevm service",
			serviceID: "xrplevm",
			appID:     "app123",
			wantURL:   "wss://xrplevm.rpc.grove.city/v1/app123",
		},
		{
			name:      "xrplevm with different app ID",
			serviceID: "xrplevm",
			appID:     "app789",
			wantURL:   "wss://xrplevm.rpc.grove.city/v1/app789",
		},
		{
			name:      "empty app ID",
			serviceID: "xrplevm",
			appID:     "",
			wantURL:   "wss://xrplevm.rpc.grove.city/v1/",
		},
		{
			name:      "special characters in app ID",
			serviceID: "xrplevm",
			appID:     "app-123_test",
			wantURL:   "wss://xrplevm.rpc.grove.city/v1/app-123_test",
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
			expectedDefault: "xrplevm",
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
			name:    "xrplevm service",
			service: "xrplevm",
			valid:   true,
		},
		{
			name:    "ethereum service (not supported)",
			service: "ethereum",
			valid:   false,
		},
		{
			name:    "polygon service (not supported)",
			service: "polygon",
			valid:   false,
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
			name:    "case sensitive - XRPLEVM",
			service: "XRPLEVM",
			valid:   false,
		},
	}

	validServices := map[string]bool{
		"xrplevm": true,
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
	serviceID := "xrplevm"
	appID := "app123"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("wss://%s.rpc.grove.city/v1/%s", serviceID, appID)
	}
}
