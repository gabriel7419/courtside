package ui

import (
	"fmt"
	"strings"

	"github.com/gabriel7419/courtside/internal/api"
)

// MatchDisplay wraps a match with display information for rendering.
type MatchDisplay struct {
	api.Match
}

// Title returns a formatted title for the match.
func (m MatchDisplay) Title() string {
	home := m.HomeTeam.ShortName
	if home == "" {
		home = m.HomeTeam.Name
	}
	away := m.AwayTeam.ShortName
	if away == "" {
		away = m.AwayTeam.Name
	}
	return home + " vs " + away
}

// Description returns a formatted description for the match.
// Shows score, league, live time on first line; KO time on second line.
func (m MatchDisplay) Description() string {
	var parts []string

	// Add score if available
	if m.HomeScore != nil && m.AwayScore != nil {
		parts = append(parts, fmt.Sprintf("%d - %d", *m.HomeScore, *m.AwayScore))
	}

	// Add league name
	if m.League.Name != "" {
		parts = append(parts, m.League.Name)
	}

	// Add live time
	if m.LiveTime != nil {
		parts = append(parts, *m.LiveTime)
	}

	line1 := strings.Join(parts, " • ")

	// Add start time (kick-off time) on second line
	if m.MatchTime != nil {
		return line1 + "\nKO " + m.MatchTime.Local().Format("15:04")
	}

	return line1
}
