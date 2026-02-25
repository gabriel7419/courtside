package nba

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gabriel7419/courtside/internal/api"
)

// LiveUpdateParser converts game events into human-readable update strings.
// It mirrors the interface of fotmob.LiveUpdateParser so the rest of the app
// can use it without changes.
type LiveUpdateParser struct{}

// NewLiveUpdateParser creates a new live update parser.
func NewLiveUpdateParser() *LiveUpdateParser {
	return &LiveUpdateParser{}
}

// Event type prefixes (used by the UI for color coding)
const (
	EventPrefixScore        = "●" // field goal / free throw
	EventPrefixFoul         = "▪" // foul
	EventPrefixTimeout      = "⏸" // timeout
	EventPrefixSubstitution = "↔" // substitution
	EventPrefixOther        = "·" // other events
)

// ParseEvents converts a list of game events into readable update strings.
// Events are sorted by minute descending (most recent first).
func (p *LiveUpdateParser) ParseEvents(events []api.MatchEvent, homeTeam, awayTeam api.Team) []string {
	sorted := make([]api.MatchEvent, len(events))
	copy(sorted, events)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Minute > sorted[j].Minute
	})

	updates := make([]string, 0, len(sorted))
	for _, event := range sorted {
		s := p.formatEvent(event, homeTeam, awayTeam)
		if s != "" {
			updates = append(updates, s)
		}
	}
	return updates
}

// formatEvent formats a single game event into a readable string.
func (p *LiveUpdateParser) formatEvent(event api.MatchEvent, homeTeam, awayTeam api.Team) string {
	isHome := event.Team.ID == homeTeam.ID
	if event.Team.ID == 0 && event.Team.ShortName != "" {
		isHome = event.Team.ShortName == homeTeam.ShortName
	}
	teamMarker := "[A]"
	if isHome {
		teamMarker = "[H]"
	}

	player := ""
	if event.Player != nil {
		player = *event.Player
	}

	switch strings.ToLower(event.Type) {
	case "field_goal":
		label := "[BASKET]"
		if event.IsThree != nil && *event.IsThree {
			label = "[3PT]"
		}
		pts := ""
		if event.Points != nil {
			pts = fmt.Sprintf("+%d", *event.Points)
		}
		return fmt.Sprintf("%s %s %s %s %s", EventPrefixScore, event.DisplayMinute, label, player+pts, teamMarker)

	case "free_throw":
		return fmt.Sprintf("%s %s [FT] %s %s", EventPrefixScore, event.DisplayMinute, player, teamMarker)

	case "foul":
		foulType := "Foul"
		if event.EventSubtype != nil {
			foulType = strings.Title(*event.EventSubtype) + " Foul"
		}
		return fmt.Sprintf("%s %s [%s] %s %s", EventPrefixFoul, event.DisplayMinute, foulType, player, teamMarker)

	case "timeout":
		return fmt.Sprintf("%s %s [Timeout] %s", EventPrefixTimeout, event.DisplayMinute, teamMarker)

	case "substitution":
		playerIn := ""
		if event.Assist != nil {
			playerIn = *event.Assist
		}
		return fmt.Sprintf("%s %s [SUB] {OUT}%s {IN}%s %s", EventPrefixSubstitution, event.DisplayMinute, player, playerIn, teamMarker)

	default:
		if player != "" {
			return fmt.Sprintf("%s %s %s %s", EventPrefixOther, event.DisplayMinute, player, teamMarker)
		}
		return ""
	}
}

// NewEvents returns events present in newEvents but not in oldEvents.
func (p *LiveUpdateParser) NewEvents(oldEvents, newEvents []api.MatchEvent) []api.MatchEvent {
	oldMap := make(map[int]bool, len(oldEvents))
	for _, e := range oldEvents {
		oldMap[e.ID] = true
	}
	var fresh []api.MatchEvent
	for _, e := range newEvents {
		if !oldMap[e.ID] {
			fresh = append(fresh, e)
		}
	}
	return fresh
}

// --- StatsData (mirrors fotmob.StatsData for app/commands.go compatibility) ---

// StatsDataDays is the number of days of data fetched for the stats view.
const StatsDataDays = 5

// StatsData holds aggregated game data for the stats view.
type StatsData struct {
	AllFinished   []api.Match
	TodayFinished []api.Match
	TodayUpcoming []api.Match
}
