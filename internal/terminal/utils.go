package terminal

import (
	"syscall"
	"unsafe"

	"github.com/fatih/color"
)

var (
	// Colors for output formatting
	Green   = color.New(color.FgGreen, color.Bold)
	Red     = color.New(color.FgRed, color.Bold)
	Yellow  = color.New(color.FgYellow, color.Bold)
	Blue    = color.New(color.FgBlue, color.Bold)
	Magenta = color.New(color.FgMagenta, color.Bold)
	Cyan    = color.New(color.FgCyan, color.Bold)
)

// GetTerminalWidth returns the current terminal width
func GetTerminalWidth() int {
	type winsize struct {
		Row    uint16
		Col    uint16
		Xpixel uint16
		Ypixel uint16
	}

	ws := &winsize{}
	retCode, _, _ := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == 0 {
		return int(ws.Col)
	}
	return 80 // fallback to 80 columns
}

// GetSubscriptionEmoji returns the appropriate emoji for each subscription type
func GetSubscriptionEmoji(subscriptionType string) string {
	switch subscriptionType {
	case "newHeads":
		return "🧊" // Ice cube for blocks
	case "newPendingTransactions":
		return "⚡" // Lightning for fast pending transactions
	case "logs":
		return "📄" // Document for logs/events
	case "syncing":
		return "🔄" // Refresh for syncing
	default:
		return "📡" // Generic antenna for unknown types
	}
}
