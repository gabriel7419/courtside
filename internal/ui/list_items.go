package ui

import (
	"github.com/0xjuanma/golazo/internal/api"
	"github.com/charmbracelet/bubbles/list"
)

// MatchListItem implements the list.Item interface for matches.
type MatchListItem struct {
	Match   api.Match
	Display MatchDisplay
}

// Title returns the match title for the list item.
func (m MatchListItem) Title() string {
	return m.Display.Title()
}

// Description returns the match description for the list item.
func (m MatchListItem) Description() string {
	return m.Display.Description()
}

// FilterValue returns the value to use for filtering.
// Returns team names for searching (e.g., "Arsenal vs Chelsea").
func (m MatchListItem) FilterValue() string {
	return m.Title()
}

// ToMatchListItems converts a slice of MatchDisplay to list items.
func ToMatchListItems(matches []MatchDisplay) []list.Item {
	items := make([]list.Item, len(matches))
	for i, match := range matches {
		items[i] = MatchListItem{
			Display: match,
			Match:   match.Match,
		}
	}
	return items
}
