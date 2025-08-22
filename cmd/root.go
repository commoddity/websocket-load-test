/*
Copyright © 2025 Grove Technologies

Licensed under MIT License
*/
package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/commoddity/websocket-load-test/internal/client"
	"github.com/commoddity/websocket-load-test/internal/stats"
	"github.com/commoddity/websocket-load-test/internal/terminal"
	"github.com/commoddity/websocket-load-test/internal/types"
	"github.com/spf13/cobra"
)

var (
	// Configuration flags
	serviceID     string
	appID         string
	apiKey        string
	subscriptions string
	subCount      int
	enableLogging bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "websocket-load-test",
	Short: "🚀 Load test WebSocket endpoints for Grove Portal blockchain services",
	Long: `🌿 WebSocket Load Test - Built for Grove Portal

A robust, feature-rich WebSocket client designed for load testing and monitoring 
Grove Portal's WebSocket endpoints. This tool provides real-time statistics, 
subscription management, and detailed connection monitoring for Ethereum-compatible 
blockchain WebSocket services.

🔗 Grove Portal: https://www.portal.grove.city/

Features:
📊 Real-time statistics dashboard with live metrics
⚡ Multiple subscription types (newHeads, newPendingTransactions, logs)
🔄 Automatic reconnection with detailed connection history  
📈 Performance monitoring (message rates, success rates, reliability)
🎨 Beautiful terminal interface with emojis and colored output
📋 Multiple subscription instances for comprehensive load testing

Prerequisites:
• Grove Portal account at https://www.portal.grove.city/
• Application created in Grove Portal dashboard
• Valid Application ID and API Key`,

	Example: `🌿 Grove Portal Examples:

  # Basic connection test (defaults to xrplevm)
  websocket-load-test \
    --app-id "your_app_id_here" \
    --api-key "your_api_key_here"

  # XRPL EVM with multiple subscriptions
  websocket-load-test \
    -a "your_app_id_here" \
    -k "your_api_key_here" \
    --subs "newHeads,newPendingTransactions" \
    --count 10

  # High-load testing with message logging
  websocket-load-test \
    --app-id "your_app_id_here" \
    --api-key "your_api_key_here" \
    --count 50 \
    --log

  # Only XRPL EVM service is supported

URLs are automatically constructed as:
  wss://xrplevm.rpc.grove.city/v1/[app-id]`,

	Run: runWebSocketLoadTest,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Grove Portal connection flags
	rootCmd.Flags().StringVarP(&serviceID, "service", "s", "xrplevm",
		"🎯 Grove Portal service (only xrplevm supported)")

	rootCmd.Flags().StringVarP(&appID, "app-id", "a", "",
		"🆔 Grove Portal Application ID")

	rootCmd.Flags().StringVarP(&apiKey, "api-key", "k", "",
		"🔐 Grove Portal API Key")

	// Subscription flags
	rootCmd.Flags().StringVar(&subscriptions, "subs", "newHeads",
		"📡 Comma-separated subscription types (newHeads,newPendingTransactions,logs)")

	rootCmd.Flags().IntVarP(&subCount, "count", "c", 1,
		"📊 Number of subscriptions to create for each type")

	rootCmd.Flags().BoolVarP(&enableLogging, "log", "l", false,
		"📝 Display latest WebSocket message in formatted JSON")

	// Mark required flags
	_ = rootCmd.MarkFlagRequired("app-id")
	_ = rootCmd.MarkFlagRequired("api-key")
}

// runWebSocketLoadTest is the main application logic
func runWebSocketLoadTest(cmd *cobra.Command, args []string) {
	// Validate service
	if serviceID != "xrplevm" {
		fmt.Printf("❌ Error: Only 'xrplevm' service is supported, got '%s'\n", serviceID)
		os.Exit(1)
	}

	// Construct Grove Portal WebSocket URL
	wsURL := fmt.Sprintf("wss://%s.rpc.grove.city/v1/%s", serviceID, appID)

	// Create configuration from flags
	config := &types.Config{
		URL:           wsURL,
		ServiceID:     serviceID,
		AuthHeader:    apiKey,
		Subscriptions: subscriptions,
		SubCount:      subCount,
		EnableLogging: enableLogging,
	}

	// Setup interrupt handler
	done := make(chan struct{})
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Initialize components
	statsManager := stats.NewManager()
	if enableLogging {
		statsManager.EnableLogging()
		statsManager.SetConfig(config)
	}
	wsClient := client.NewWebSocketClient(config, statsManager, done)

	// Display startup information
	displayStartupInfo(config)

	// Start the WebSocket client
	wsClient.Start()

	// Start automatic display updates
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				if statsManager.GetStats().TotalConnections > 0 {
					statsManager.DisplayRunningStats(wsClient.GetTotalSubscriptions())
				}
			}
		}
	}()

	// Wait for interrupt
	<-interrupt
	terminal.Cyan.Println("\n🛑 Received interrupt signal, shutting down...")
	close(done)

	// Print final statistics
	statsManager.PrintFinalStats(wsClient.GetTotalSubscriptions())
}

// displayStartupInfo shows the initial startup information
func displayStartupInfo(config *types.Config) {
	terminal.Green.Println("🚀 Starting WebSocket Load Test...")
	terminal.Green.Printf("📊 Target: %s\n", config.URL)
	terminal.Green.Printf("🎯 Service: %s\n", config.ServiceID)

	// Parse subscriptions
	subTypes := strings.Split(config.Subscriptions, ",")
	totalSubsToCreate := len(subTypes) * config.SubCount
	terminal.Green.Printf("📡 Subscriptions (%d types × %d instances = %d total):\n", len(subTypes), config.SubCount, totalSubsToCreate)
	for _, sub := range subTypes {
		sub = strings.TrimSpace(sub)
		emoji := terminal.GetSubscriptionEmoji(sub)
		terminal.Green.Printf("  %s %s (×%d)\n", emoji, sub, config.SubCount)
	}

	if config.AuthHeader != "" {
		authDisplay := config.AuthHeader
		if len(authDisplay) > 20 {
			authDisplay = authDisplay[:20]
		}
		terminal.Green.Printf("🔐 Auth: %s...\n", authDisplay)
	}
	fmt.Println()
}
