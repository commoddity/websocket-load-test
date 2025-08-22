package types

import "time"

// Stats contains all statistics for the WebSocket client
type Stats struct {
	TotalConnections    int
	TotalReconnections  int
	CurrentConnStart    time.Time
	TotalUptime         time.Duration
	EventsReceived      int
	ClientStartTime     time.Time
	SubscriptionEvents  int
	ConfirmationEvents  int
	ErrorEvents         int
	LastEventTime       time.Time
	ConnectionAttempts  int
	CurrentConnMessages int
	LongestConnection   time.Duration
	ShortestConnection  time.Duration
}

// ConnectionHistory tracks individual connection sessions
type ConnectionHistory struct {
	ConnectionNum int
	StartTime     time.Time
	EndTime       time.Time
	Duration      time.Duration
	Messages      int
}

// JSONRPCRequest represents a JSON-RPC request
type JSONRPCRequest struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Method  string `json:"method"`
	Params  any    `json:"params,omitempty"`
}

// JSONRPCResponse represents a JSON-RPC response
type JSONRPCResponse struct {
	JSONRPC string `json:"jsonrpc"`
	ID      any    `json:"id,omitempty"`
	Result  any    `json:"result,omitempty"`
	Error   any    `json:"error,omitempty"`
	Method  string `json:"method,omitempty"`
	Params  any    `json:"params,omitempty"`
}

// Config holds the configuration for the WebSocket client
type Config struct {
	URL           string
	ServiceID     string
	AuthHeader    string
	Subscriptions string
	SubCount      int
}
