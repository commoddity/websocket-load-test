# WebSocket Load Test

A simple WebSocket client designed for load testing and monitoring Grove Portal's WebSocket endpoints. 

This tool provides real-time statistics, subscription management, and detailed connection monitoring for Ethereum-compatible blockchain WebSocket services.

<p align="center">
<a href="https://github.com/buildwithgrove/path">
  <img src="https://storage.googleapis.com/grove-brand-assets/Presskit/Logo%20Joined-2.png" alt="Grove Logo" title="Uses PATH for RPC" width="500">
  <br>
  🌿 Load tests WebSocket endpoints using Grove's PATH API.
  <br>
</a>
</p>

## Prerequisites

**Grove Portal Account Required** - This tool is designed to work with Grove Portal's WebSocket endpoints. You'll need:

1. **Grove Portal Account**: Sign up at [https://www.portal.grove.city/](https://www.portal.grove.city/)
2. **Application**: Create an application in your Grove Portal dashboard
3. **API Credentials**: Get your Application ID and API Key from the portal

## Features

- 🚀 **Multiple Subscription Types**: Support for `newHeads`, `newPendingTransactions`, `logs`, and custom subscriptions
- 📊 **Real-time Statistics**: Live dashboard with connection metrics, message rates, and performance data
- 🔄 **Automatic Reconnection**: Robust reconnection logic with detailed connection history
- 📈 **Performance Monitoring**: Track message rates, success rates, and connection reliability
- 🎨 **Colorized Output**: Beautiful terminal interface with emojis and colored output
- ⚡ **Multiple Instances**: Create multiple subscription instances for load testing
- 📋 **Connection History**: Detailed tracking of all connection sessions

## Installation

```bash
go install github.com/commoddity/websocket-load-test@latest
```

## Usage

```bash
# Load test XRPLEVM WebSocket endpoint
websocket-load-test \                                                  
    --service "xrplevm" \
    --app-id $GROVE_PORTAL_APP_ID \
    --api-key $GROVE_PORTAL_API_KEY \
    --subs "newHeads,newPendingTransactions" \
    --count 10
```

This example will:
- Connect to Grove Portal's XRPL EVM WebSocket endpoint
- Authenticate using your API key
- Create 10 `newHeads` subscriptions and 10 `newPendingTransactions` subscriptions
- Display real-time statistics and performance metrics

### Command Line Options

| Flag        | Short  | Description                        | Default                  | Example                            |
| ----------- | ------ | ---------------------------------- | ------------------------ | ---------------------------------- |
| `--url`     | `-u`   | WebSocket URL to connect to        | `ws://localhost:3069/v1` | `--url "wss://api.example.com/v1"` |
| `--service` | `-s`   | Target service ID header value     | `ethereum`               | `--service "polygon"`              |
| `--auth`    | `-a`   | Authorization header value         | _(empty)_                | `--auth "token123"`                |
| `--subs`    | _none_ | Comma-separated subscription types | `newHeads`               | `--subs "newHeads,logs"`           |
| `--count`   | `-c`   | Number of subscriptions per type   | `1`                      | `--count 10`                       |
| `--help`    | `-h`   | Show detailed help and examples    | _none_                   | `--help`                           |

Use `websocket-load-test --help` for detailed usage examples and feature descriptions.

### Supported Subscription Types

- **`newHeads`** 🧊 - New block headers
- **`newPendingTransactions`** ⚡ - Pending transactions

## Dashboard Features

The real-time dashboard displays:

### Connection Metrics
- Total connections and reconnections
- Connection attempts and current duration
- Average, longest, and shortest connection times
- Connection reliability percentage

### Subscription Metrics
- Total subscriptions created
- Confirmation and subscription events
- Error events and success rate
- Messages categorized by subscription type

### Message Metrics
- Total messages received
- Current connection messages
- Messages per second (current and overall)
- Time since last event

### Performance Metrics
- Success rate percentage
- Events per subscription ratio
- Connection event rates
- Overall reliability statistics

### Connection History
- Last 5 connection sessions
- Duration and message count per session
- Connection start/end timestamps

## Project Structure

```
websocket-load-test/
├── main.go                          # Application entry point
├── internal/
│   ├── client/
│   │   └── websocket.go            # WebSocket client implementation
│   ├── stats/
│   │   └── manager.go              # Statistics collection and display
│   ├── terminal/
│   │   └── utils.go                # Terminal utilities and colors
│   └── types/
│       └── models.go               # Data models and types
├── go.mod                          # Go module definition
├── go.sum                          # Go module checksums
├── .golangci.yml                   # Linter configuration
└── README.md                       # This file
```

## Architecture

The application is organized into focused packages:

- **`main`**: Application entry point and configuration
- **`client`**: WebSocket connection management and message handling
- **`stats`**: Statistics collection, calculation, and display
- **`terminal`**: Terminal utilities, colors, and UI helpers
- **`types`**: Data models and type definitions

## Configuration

### Grove Portal Configuration

This tool is designed for Grove Portal. Required configuration:

1. **Grove Portal Account**: Create account at [https://www.portal.grove.city/](https://www.portal.grove.city/)
2. **Application Setup**: Create an application in your Grove Portal dashboard
3. **Get Your Credentials**: Copy your Application ID and API Key from the dashboard

### Best Practices

1. **Authentication**: Always provide your Grove Portal API key with `--api-key` or `-k`
2. **Service Selection**: Choose the correct service (ethereum, polygon, xrplevm, arbitrum, optimism, base)
3. **URL Auto-Construction**: URLs are automatically built as `wss://[service].rpc.grove.city/v1/[app-id]`
4. **Load Testing**: Start with small `--count` values (1-10) and increase gradually
5. **Rate Limits**: Monitor your Grove Portal dashboard for usage and limits

## Troubleshooting

### Common Issues

1. **Connection Failed**
   - Verify your Grove Portal Application ID is correct in the URL
   - Check that your API key is valid and set correctly
   - Ensure you have an active Grove Portal subscription
   - Verify the chain service ID matches (ethereum, polygon, xrplevm, etc.)

2. **Authentication Errors**
   - Check that `GROVE_PORTAL_API_KEY` environment variable is set
   - Verify your API key hasn't expired in the Grove Portal dashboard
   - Ensure your application has WebSocket access enabled

3. **No Messages Received**
   - Verify the subscription types are supported by the target chain
   - Check your Grove Portal dashboard for any service disruptions
   - Monitor error events in the real-time dashboard

4. **Rate Limiting**
   - Check your Grove Portal dashboard for rate limit status
   - Reduce the `-count` parameter for load testing
   - Monitor your application's usage in the Grove Portal dashboard

5. **High Reconnection Rate**
   - Check network stability
   - Verify you haven't exceeded Grove Portal rate limits
   - Consider reducing subscription count or frequency
