// Package notify provides desktop notification functionality for match events.
// Currently supports macOS, Linux, and Windows via the beeep library.
package notify

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/gabriel7419/courtside/internal/api"
	"github.com/gabriel7419/courtside/internal/assets"
	"github.com/gabriel7419/courtside/internal/constants"
	"github.com/gabriel7419/courtside/internal/data"
	"github.com/gen2brain/beeep"
)

var (
	iconPath     string
	iconPathOnce sync.Once
)

// getIconPath returns the path to the cached notification icon.
// The icon is embedded in the binary and written to the cache directory on first use.
// Returns empty string if caching fails (notification will work without icon).
func getIconPath() string {
	iconPathOnce.Do(func() {
		cacheDir, err := data.CacheDir()
		if err != nil {
			return
		}

		iconPath = filepath.Join(cacheDir, "icon.png")

		// Only write if file doesn't exist
		if _, err := os.Stat(iconPath); os.IsNotExist(err) {
			if err := os.WriteFile(iconPath, assets.Logo, 0644); err != nil {
				iconPath = "" // Reset on write failure
			}
		}
	})
	return iconPath
}

// Notifier defines the interface for sending desktop notifications.
// This allows for easy mocking in tests and potential future implementations.
type Notifier interface {
	// Goal sends a notification for a new goal event.
	Goal(event api.MatchEvent, homeTeam, awayTeam api.Team, homeScore, awayScore int) error
}

// DesktopNotifier implements Notifier using native desktop notifications.
type DesktopNotifier struct {
	enabled bool
}

// NewDesktopNotifier creates a new desktop notifier.
// Notifications are enabled by default.
func NewDesktopNotifier() *DesktopNotifier {
	return &DesktopNotifier{
		enabled: true,
	}
}

// SetEnabled enables or disables notifications.
func (n *DesktopNotifier) SetEnabled(enabled bool) {
	n.enabled = enabled
}

// Enabled returns whether notifications are currently enabled.
func (n *DesktopNotifier) Enabled() bool {
	return n.enabled
}

// Goal sends a desktop notification for a new scoring event.
// Works for both football goals and NBA field goals / free throws.
func (n *DesktopNotifier) Goal(event api.MatchEvent, homeTeam, awayTeam api.Team, homeScore, awayScore int) error {
	if !n.enabled {
		return nil
	}

	// Play terminal beep via stderr
	_, _ = os.Stderr.WriteString("\a")

	// Build notification content
	title := constants.NotificationTitleGoal
	message := formatGoalMessage(event, homeTeam, awayTeam, homeScore, awayScore)

	_ = beeep.Notify(title, message, getIconPath())
	return nil
}

// formatGoalMessage creates the notification message for a scoring event.
// Handles football goals and NBA field goals / free throws.
func formatGoalMessage(event api.MatchEvent, homeTeam, awayTeam api.Team, homeScore, awayScore int) string {
	scorer := "Unknown"
	if event.Player != nil {
		scorer = *event.Player
	}

	teamName := event.Team.ShortName
	if teamName == "" {
		teamName = event.Team.Name
	}

	// Use DisplayMinute if set (e.g. "Q3 2:34"), otherwise fall back to Minute'
	timeStr := ""
	if event.DisplayMinute != "" {
		timeStr = event.DisplayMinute
	} else {
		timeStr = fmt.Sprintf("%d'", event.Minute)
	}

	// Build event label
	var label string
	switch event.Type {
	case "field_goal":
		if event.IsThree != nil && *event.IsThree {
			label = "3PT"
		} else {
			label = "BASKET"
		}
		if event.Points != nil {
			label = fmt.Sprintf("%s +%d", label, *event.Points)
		}
	case "free_throw":
		label = "FT +1"
	default:
		label = "GOAL"
		if event.Assist != nil && *event.Assist != "" {
			scorer = fmt.Sprintf("%s (%s)", scorer, *event.Assist)
		}
	}

	return fmt.Sprintf("%s %s [%s · %s]\n%s %d - %d %s",
		scorer,
		timeStr,
		label,
		teamName,
		homeTeam.ShortName,
		homeScore,
		awayScore,
		awayTeam.ShortName,
	)
}
