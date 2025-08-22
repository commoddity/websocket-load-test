package stats

import (
	"fmt"
	"strings"
	"time"

	"github.com/commoddity/websocket-load-test/internal/terminal"
	"github.com/commoddity/websocket-load-test/internal/types"
)

// Manager handles statistics collection and display
type Manager struct {
	stats             *types.Stats
	connectionHistory []types.ConnectionHistory
	messagesByType    map[string]int
	subIDToType       map[string]string
	spinnerChars      []string
	spinnerIndex      int
	needFullClear     bool
}

// NewManager creates a new statistics manager
func NewManager() *Manager {
	return &Manager{
		stats:          &types.Stats{ClientStartTime: time.Now()},
		messagesByType: make(map[string]int),
		subIDToType:    make(map[string]string),
		spinnerChars:   []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		needFullClear:  true,
	}
}

// GetStats returns the current stats
func (m *Manager) GetStats() *types.Stats {
	return m.stats
}

// IncrementConnectionAttempts increments the connection attempts counter
func (m *Manager) IncrementConnectionAttempts() {
	m.stats.ConnectionAttempts++
}

// StartNewConnection starts tracking a new connection
func (m *Manager) StartNewConnection() {
	m.stats.TotalConnections++
	m.stats.CurrentConnStart = time.Now()
	m.stats.CurrentConnMessages = 0
	m.needFullClear = true
}

// IncrementReconnections increments the reconnection counter
func (m *Manager) IncrementReconnections() {
	if m.stats.TotalConnections > 0 {
		m.stats.TotalReconnections++
	}
}

// EndConnection records the end of a connection
func (m *Manager) EndConnection() {
	if m.stats.TotalConnections > 0 {
		connectionDuration := time.Since(m.stats.CurrentConnStart)
		m.stats.TotalUptime += connectionDuration

		// Record connection history
		m.connectionHistory = append(m.connectionHistory, types.ConnectionHistory{
			ConnectionNum: m.stats.TotalConnections,
			StartTime:     m.stats.CurrentConnStart,
			EndTime:       time.Now(),
			Duration:      connectionDuration,
			Messages:      m.stats.CurrentConnMessages,
		})

		// Update longest/shortest connection times
		if m.stats.LongestConnection == 0 || connectionDuration > m.stats.LongestConnection {
			m.stats.LongestConnection = connectionDuration
		}
		if m.stats.ShortestConnection == 0 || connectionDuration < m.stats.ShortestConnection {
			m.stats.ShortestConnection = connectionDuration
		}

		m.needFullClear = true
	}
}

// HandleResponse processes a WebSocket response and updates statistics
func (m *Manager) HandleResponse(response types.JSONRPCResponse) {
	m.stats.EventsReceived++
	m.stats.CurrentConnMessages++
	m.stats.LastEventTime = time.Now()

	if response.Method == "eth_subscription" {
		m.stats.SubscriptionEvents++

		// Extract subscription type from the subscription event
		if params, ok := response.Params.(map[string]interface{}); ok {
			if subscription, exists := params["subscription"]; exists {
				subscriptionType := m.getSubscriptionTypeFromID(fmt.Sprintf("%v", subscription))
				if subscriptionType != "" {
					m.messagesByType[subscriptionType]++
				} else {
					m.messagesByType["unknown"]++
				}
			}
		}
	} else if response.Result != nil {
		// Check if this is a subscription confirmation response
		if _, ok := response.ID.(float64); ok {
			m.stats.ConfirmationEvents++
		}
	} else if response.Error != nil {
		m.stats.ErrorEvents++
	}
}

// SetSubscriptionMapping sets the mapping between subscription ID and type
func (m *Manager) SetSubscriptionMapping(subscriptionID, subscriptionType string) {
	m.subIDToType[subscriptionID] = subscriptionType
}

// getSubscriptionTypeFromID attempts to determine subscription type from subscription ID
func (m *Manager) getSubscriptionTypeFromID(subscriptionID string) string {
	if subType, exists := m.subIDToType[subscriptionID]; exists {
		return subType
	}
	return ""
}

// DisplayRunningStats shows a constantly updating dashboard of statistics
func (m *Manager) DisplayRunningStats(totalSubscriptions int) {
	terminalWidth := terminal.GetTerminalWidth()

	if m.needFullClear {
		fmt.Print("\033[2J\033[H")
		m.needFullClear = false
	} else {
		fmt.Print("\033[H\033[0J")
	}

	// Update spinner
	m.spinnerIndex = (m.spinnerIndex + 1) % len(m.spinnerChars)

	// Calculate timing stats
	currentConnDuration := time.Since(m.stats.CurrentConnStart)
	totalClientRuntime := time.Since(m.stats.ClientStartTime)

	// Calculate rates
	var messagesPerSecond, overallRate float64
	if currentConnDuration.Seconds() > 0 {
		messagesPerSecond = float64(m.stats.CurrentConnMessages) / currentConnDuration.Seconds()
	}
	if totalClientRuntime.Seconds() > 0 {
		overallRate = float64(m.stats.EventsReceived) / totalClientRuntime.Seconds()
	}

	// Calculate time since last event
	var timeSinceLastEvent time.Duration
	if !m.stats.LastEventTime.IsZero() {
		timeSinceLastEvent = time.Since(m.stats.LastEventTime)
	}

	// Header with spinner
	headerText := fmt.Sprintf("%s WebSocket Client Dashboard - Live Stats", m.spinnerChars[m.spinnerIndex])
	if len(headerText) > terminalWidth {
		headerText = headerText[:terminalWidth-3] + "..."
	}
	terminal.Green.Println(headerText)

	// Separator line
	separatorWidth := terminalWidth
	if separatorWidth > 100 {
		separatorWidth = 100
	}
	if separatorWidth < 20 {
		separatorWidth = 20
	}
	fmt.Println(strings.Repeat("═", separatorWidth))

	// Connection Stats
	terminal.Cyan.Println("📡 CONNECTION METRICS")
	fmt.Printf("🔗 Total Connections:     %s%d%s\n", terminal.Green.Sprint(""), m.stats.TotalConnections, "")
	fmt.Printf("🔄 Reconnections:         %s%d%s\n", terminal.Yellow.Sprint(""), m.stats.TotalReconnections, "")
	fmt.Printf("🎯 Connection Attempts:   %s%d%s\n", terminal.Blue.Sprint(""), m.stats.ConnectionAttempts, "")
	fmt.Printf("⏱️  Current Conn Duration: %s%v%s\n", terminal.Green.Sprint(""), currentConnDuration.Round(time.Second), "")
	fmt.Printf("🏃 Total Runtime:         %s%v%s\n", terminal.Cyan.Sprint(""), totalClientRuntime.Round(time.Second), "")

	// Calculate and show average connection duration
	if len(m.connectionHistory) > 0 {
		var totalDuration time.Duration
		for _, conn := range m.connectionHistory {
			totalDuration += conn.Duration
		}
		avgDuration := totalDuration / time.Duration(len(m.connectionHistory))
		fmt.Printf("📊 Avg Connection Time:   %s%v%s\n", terminal.Blue.Sprint(""), avgDuration.Round(time.Second), "")
	}

	// Subscription Stats
	fmt.Println()
	terminal.Magenta.Println("📡 SUBSCRIPTION METRICS")
	fmt.Printf("📊 Total Subscriptions:   %s%d%s\n", terminal.Magenta.Sprint(""), totalSubscriptions, "")
	fmt.Printf("✅ Confirmations:         %s%d%s\n", terminal.Green.Sprint(""), m.stats.ConfirmationEvents, "")
	fmt.Printf("🧊 Subscription Events:   %s%d%s\n", terminal.Cyan.Sprint(""), m.stats.SubscriptionEvents, "")
	fmt.Printf("❌ Error Events:          %s%d%s\n", terminal.Red.Sprint(""), m.stats.ErrorEvents, "")

	// Show messages by subscription type
	if len(m.messagesByType) > 0 {
		fmt.Println()
		terminal.Blue.Println("📊 MESSAGES BY TYPE")
		for subType, count := range m.messagesByType {
			emoji := terminal.GetSubscriptionEmoji(subType)
			fmt.Printf("%s %s: %s%d%s msgs\n", emoji, subType, terminal.Cyan.Sprint(""), count, "")
		}
	}

	// Message Stats
	fmt.Println()
	terminal.Blue.Println("📨 MESSAGE METRICS")
	fmt.Printf("📈 Total Messages:        %s%d%s\n", terminal.Blue.Sprint(""), m.stats.EventsReceived, "")
	fmt.Printf("📨 Current Conn Messages: %s%d%s\n", terminal.Cyan.Sprint(""), m.stats.CurrentConnMessages, "")
	fmt.Printf("⚡ Messages/Second:       %s%.2f%s\n", terminal.Yellow.Sprint(""), messagesPerSecond, "")
	fmt.Printf("📊 Overall Rate:          %s%.2f%s/sec\n", terminal.Cyan.Sprint(""), overallRate, "")
	fmt.Printf("⏰ Last Event:            %s%v%s ago\n", terminal.Green.Sprint(""), timeSinceLastEvent.Round(time.Second), "")

	// Performance Stats
	fmt.Println()
	terminal.Yellow.Println("⚡ PERFORMANCE METRICS")

	// Success rate
	if m.stats.EventsReceived > 0 {
		successRate := float64(m.stats.EventsReceived-m.stats.ErrorEvents) / float64(m.stats.EventsReceived) * 100
		fmt.Printf("✅ Success Rate:          %s%.1f%%%s\n", terminal.Green.Sprint(""), successRate, "")
	}

	// Events per subscription
	if totalSubscriptions > 0 {
		eventsPerSub := float64(m.stats.SubscriptionEvents) / float64(totalSubscriptions)
		fmt.Printf("📊 Events/Subscription:   %s%.1f%s\n", terminal.Cyan.Sprint(""), eventsPerSub, "")
	}

	// Connection duration metrics
	if m.stats.LongestConnection > 0 {
		fmt.Printf("🏆 Longest Connection:    %s%v%s\n", terminal.Green.Sprint(""), m.stats.LongestConnection.Round(time.Second), "")
	}
	if m.stats.ShortestConnection > 0 {
		fmt.Printf("⚡ Shortest Connection:   %s%v%s\n", terminal.Yellow.Sprint(""), m.stats.ShortestConnection.Round(time.Second), "")
	}

	// Connection History Section
	if len(m.connectionHistory) > 0 {
		fmt.Println()
		fmt.Println()
		terminal.Yellow.Println("📋 CONNECTION HISTORY")

		// Show last 5 connections
		start := 0
		if len(m.connectionHistory) > 5 {
			start = len(m.connectionHistory) - 5
		}

		for i := start; i < len(m.connectionHistory); i++ {
			conn := m.connectionHistory[i]
			fmt.Printf("🔗 Connection #%s%d%s: %s%d%s msgs in %s%v%s (%s to %s)\n",
				terminal.Green.Sprint(""), conn.ConnectionNum, "",
				terminal.Cyan.Sprint(""), conn.Messages, "",
				terminal.Blue.Sprint(""), conn.Duration.Round(time.Second), "",
				conn.StartTime.Format("15:04:05"),
				conn.EndTime.Format("15:04:05"))
		}
	}

	// Footer
	fmt.Println()
	fmt.Println(strings.Repeat("═", separatorWidth))
	fmt.Printf("🕐 Last Updated: %s\n", time.Now().Format("15:04:05"))
}

// PrintFinalStats displays the final session summary
func (m *Manager) PrintFinalStats(totalSubscriptions int) {
	if m.stats.CurrentConnStart != (time.Time{}) {
		m.stats.TotalUptime += time.Since(m.stats.CurrentConnStart)
	}

	totalClientRuntime := time.Since(m.stats.ClientStartTime)

	// Clear screen and show final summary
	fmt.Print("\033[2J\033[H")

	terminal.Cyan.Println("🏁 FINAL SESSION SUMMARY")
	fmt.Println(strings.Repeat("═", 60))

	// Connection Summary
	terminal.Cyan.Println("📡 CONNECTION SUMMARY")
	fmt.Printf("🔗 Total Connections:     %s%d%s\n", terminal.Green.Sprint(""), m.stats.TotalConnections, "")
	fmt.Printf("🔄 Total Reconnections:   %s%d%s\n", terminal.Yellow.Sprint(""), m.stats.TotalReconnections, "")
	fmt.Printf("🎯 Connection Attempts:   %s%d%s\n", terminal.Blue.Sprint(""), m.stats.ConnectionAttempts, "")
	fmt.Printf("📡 Total Subscriptions:   %s%d%s\n", terminal.Magenta.Sprint(""), totalSubscriptions, "")
	fmt.Printf("⏱️  Total Uptime:         %s%v%s\n", terminal.Green.Sprint(""), m.stats.TotalUptime.Round(time.Second), "")
	fmt.Printf("🏃 Total Runtime:         %s%v%s\n", terminal.Cyan.Sprint(""), totalClientRuntime.Round(time.Second), "")

	// Message Summary
	fmt.Println()
	terminal.Blue.Println("📨 MESSAGE SUMMARY")
	fmt.Printf("📈 Total Messages:        %s%d%s\n", terminal.Blue.Sprint(""), m.stats.EventsReceived, "")
	fmt.Printf("🧊 Subscription Events:   %s%d%s\n", terminal.Cyan.Sprint(""), m.stats.SubscriptionEvents, "")
	fmt.Printf("✅ Confirmations:         %s%d%s\n", terminal.Green.Sprint(""), m.stats.ConfirmationEvents, "")
	fmt.Printf("❌ Error Events:          %s%d%s\n", terminal.Red.Sprint(""), m.stats.ErrorEvents, "")

	// Performance Summary
	fmt.Println()
	terminal.Yellow.Println("⚡ PERFORMANCE SUMMARY")

	if m.stats.EventsReceived > 0 && m.stats.TotalUptime > 0 {
		connectionRate := float64(m.stats.EventsReceived) / m.stats.TotalUptime.Seconds()
		fmt.Printf("📈 Connection Event Rate: %s%.2f%s events/sec\n", terminal.Yellow.Sprint(""), connectionRate, "")
	}

	if m.stats.EventsReceived > 0 && totalClientRuntime > 0 {
		overallRate := float64(m.stats.EventsReceived) / totalClientRuntime.Seconds()
		fmt.Printf("📊 Overall Event Rate:    %s%.2f%s events/sec\n", terminal.Cyan.Sprint(""), overallRate, "")
	}

	if totalClientRuntime > 0 {
		reliability := (m.stats.TotalUptime.Seconds() / totalClientRuntime.Seconds()) * 100
		fmt.Printf("📡 Connection Reliability: %s%.1f%%%s\n", terminal.Green.Sprint(""), reliability, "")
	}

	if m.stats.EventsReceived > 0 {
		successRate := float64(m.stats.EventsReceived-m.stats.ErrorEvents) / float64(m.stats.EventsReceived) * 100
		fmt.Printf("✅ Success Rate:          %s%.1f%%%s\n", terminal.Green.Sprint(""), successRate, "")
	}

	if m.stats.TotalConnections > 1 {
		avgConnectionTime := m.stats.TotalUptime / time.Duration(m.stats.TotalConnections)
		fmt.Printf("⏳ Avg Connection Time:   %s%v%s\n", terminal.Blue.Sprint(""), avgConnectionTime.Round(time.Second), "")
	}

	fmt.Println()
	fmt.Println(strings.Repeat("═", 60))
	terminal.Green.Println("👋 Session Complete - Thanks for using WebSocket Client!")
}
