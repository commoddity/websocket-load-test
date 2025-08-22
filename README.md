# WebSocket Load Test

A simple WebSocket client designed for load testing and monitoring Grove Portal's WebSocket endpoints. 

This tool provides real-time statistics, subscription management, per-type message logging, and detailed connection monitoring for XRPL EVM blockchain WebSocket services.

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

- 🚀 **Multiple Subscription Types**: Support for `newHeads`, `newPendingTransactions`
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
    --service xrplevm \
    --app-id $GROVE_PORTAL_APP_ID \
    --api-key $GROVE_PORTAL_API_KEY \
    --subs "newHeads,newPendingTransactions" \
    --count 10
    --log
```

This example will:
- Connect to Grove Portal's XRPL EVM WebSocket endpoint
- Authenticate using your API key
- Create 10 `newHeads` subscriptions and 10 `newPendingTransactions` subscriptions
- Display real-time statistics and performance metrics

### Command Line Options

| Flag        | Short  | Description                         | Default      | Example                  |
| ----------- | ------ | ----------------------------------- | ------------ | ------------------------ |
| `--service` | `-s`   | Grove Portal service (only xrplevm) | `xrplevm`    | `--service "xrplevm"`    |
| `--app-id`  | `-a`   | Grove Portal Application ID         | _(required)_ | `--app-id "app123"`      |
| `--api-key` | `-k`   | Grove Portal API Key                | _(required)_ | `--api-key "key456"`     |
| `--subs`    | _none_ | Comma-separated subscription types  | `newHeads`   | `--subs "newHeads,logs"` |
| `--count`   | `-c`   | Number of subscriptions per type    | `1`          | `--count 10`             |
| `--log`     | `-l`   | Display latest WebSocket message    | `false`      | `--log`                  |
| `--help`    | `-h`   | Show detailed help and examples     | _none_       | `--help`                 |

Use `websocket-load-test --help` for detailed usage examples and feature descriptions.

### Supported Subscription Types

- **`newHeads`** 🧊 - New block headers
- **`newPendingTransactions`** ⚡ - Pending transactions

## Message Logging

Use the `--log` or `-l` flag to enable real-time message logging. When enabled, the tool displays the latest received WebSocket message for each subscription type in formatted JSON below the dashboard:

```bash
# Enable message logging
websocket-load-test \
    --app-id "your_app_id" \
    --api-key "your_api_key" \
    --log
```

**Features:**
- 📝 **Per-Type Display**: Shows the latest message for each subscription type (newHeads, newPendingTransactions, logs, etc.)
- 🕐 **Timestamps**: Displays when each message was received
- 🎨 **JSON Formatting**: Pretty-printed JSON with proper indentation
- 🔄 **Live Updates**: Automatically replaces with newer messages per type
- 📊 **Categorization**: Separates subscription events, confirmations, and errors
- 🎯 **Organized Layout**: Groups messages by subscription type with appropriate emojis

**Example Output:**
```
🕐 Last Updated: 19:57:32

📝 LATEST MESSAGES BY TYPE

✅ confirmations - Received at 19:57:25:
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": "0x123abc..."
}

🧊 newHeads - Received at 19:57:31:
{
  "jsonrpc": "2.0",
  "method": "eth_subscription",
  "params": {
    "subscription": "0x123abc...",
    "result": {
      "number": "0x1234567",
      "hash": "0xabcdef...",
      "parentHash": "0x987654...",
      "timestamp": "0x64a8b2f0"
    }
  }
}

⚡ newPendingTransactions - Received at 19:57:30:
{
  "jsonrpc": "2.0",
  "method": "eth_subscription",
  "params": {
    "subscription": "0x456def...",
    "result": {
      "hash": "0x987654...",
      "from": "0xabc123...",
      "to": "0xdef456..."
    }
  }
}
```

