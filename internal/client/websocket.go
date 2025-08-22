package client

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/commoddity/websocket-load-test/internal/stats"
	"github.com/commoddity/websocket-load-test/internal/terminal"
	"github.com/commoddity/websocket-load-test/internal/types"
	"github.com/gorilla/websocket"
)

// WebSocketClient manages WebSocket connections and subscriptions
type WebSocketClient struct {
	config             *types.Config
	statsManager       *stats.Manager
	subscriptionIDs    map[string]int
	idToSubscription   map[int]string
	totalSubscriptions int
	done               chan struct{}
}

// NewWebSocketClient creates a new WebSocket client
func NewWebSocketClient(config *types.Config, statsManager *stats.Manager, done chan struct{}) *WebSocketClient {
	return &WebSocketClient{
		config:           config,
		statsManager:     statsManager,
		subscriptionIDs:  make(map[string]int),
		idToSubscription: make(map[int]string),
		done:             done,
	}
}

// GetTotalSubscriptions returns the total number of subscriptions
func (c *WebSocketClient) GetTotalSubscriptions() int {
	return c.totalSubscriptions
}

// Start begins the connection loop
func (c *WebSocketClient) Start() {
	go c.connectionLoop()
}

// connectionLoop handles the main connection lifecycle
func (c *WebSocketClient) connectionLoop() {
	for {
		select {
		case <-c.done:
			return
		default:
			c.connectAndListen()
		}
	}
}

// connectAndListen establishes a WebSocket connection and listens for messages
func (c *WebSocketClient) connectAndListen() {
	// Parse the WebSocket URL
	u, err := url.Parse(c.config.URL)
	if err != nil {
		terminal.Red.Printf("❌ Invalid URL: %v\n", err)
		time.Sleep(5 * time.Second)
		return
	}

	// Convert to WebSocket scheme if needed
	switch u.Scheme {
	case "http":
		u.Scheme = "ws"
	case "https":
		u.Scheme = "wss"
	}

	headers := http.Header{}
	headers.Add("Target-Service-Id", c.config.ServiceID)

	// Add authorization header if provided
	if c.config.AuthHeader != "" {
		headers.Add("Authorization", c.config.AuthHeader)
	}

	c.statsManager.IncrementConnectionAttempts()

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), headers)
	if err != nil {
		c.statsManager.IncrementReconnections()
		time.Sleep(5 * time.Second)
		return
	}

	defer conn.Close()

	// Update stats
	c.statsManager.StartNewConnection()

	// Show initial stats display
	c.statsManager.DisplayRunningStats(c.totalSubscriptions)

	// Send subscription requests
	c.sendSubscriptions(conn)

	// Listen for messages
	c.listenForMessages(conn)
}

// sendSubscriptions sends all subscription requests to the WebSocket server
func (c *WebSocketClient) sendSubscriptions(conn *websocket.Conn) {
	subTypes := strings.Split(c.config.Subscriptions, ",")
	requestID := 1

	for _, sub := range subTypes {
		sub = strings.TrimSpace(sub)
		if sub == "" {
			continue
		}

		// Create multiple instances of each subscription type
		for instance := 1; instance <= c.config.SubCount; instance++ {
			var params interface{}
			switch sub {
			case "newHeads":
				params = []string{"newHeads"}
			case "newPendingTransactions":
				params = []string{"newPendingTransactions"}
			case "logs":
				params = []interface{}{"logs", map[string]interface{}{"topics": []interface{}{nil}}}
			default:
				params = []string{sub}
			}

			subscribeReq := types.JSONRPCRequest{
				JSONRPC: "2.0",
				ID:      requestID,
				Method:  "eth_subscribe",
				Params:  params,
			}

			if err := conn.WriteJSON(subscribeReq); err != nil {
				terminal.Red.Printf("❌ Failed to send subscription for %s #%d: %v\n", sub, instance, err)
				requestID++
				continue
			}

			// Store mapping for response tracking
			subKey := fmt.Sprintf("%s#%d", sub, instance)
			c.subscriptionIDs[subKey] = requestID
			c.idToSubscription[requestID] = sub

			c.totalSubscriptions++
			requestID++

			// Add small delay between subscriptions to avoid overwhelming the server
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// listenForMessages listens for incoming WebSocket messages
func (c *WebSocketClient) listenForMessages(conn *websocket.Conn) {
	for {
		select {
		case <-c.done:
			return
		default:
			var response types.JSONRPCResponse
			err := conn.ReadJSON(&response)
			if err != nil {
				c.statsManager.EndConnection()
				c.statsManager.IncrementReconnections()
				time.Sleep(2 * time.Second)
				return
			}

			// Handle the response
			c.handleResponse(response)
		}
	}
}

// handleResponse processes incoming WebSocket responses
func (c *WebSocketClient) handleResponse(response types.JSONRPCResponse) {
	c.statsManager.HandleResponse(response)

	// Handle subscription confirmation responses
	if response.Result != nil {
		if id, ok := response.ID.(float64); ok {
			if subType, exists := c.idToSubscription[int(id)]; exists {
				// Store the actual subscription ID returned by the server
				if resultStr, ok := response.Result.(string); ok {
					c.statsManager.SetSubscriptionMapping(resultStr, subType)
				}
			}
		}
	}
}
