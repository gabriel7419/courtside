package ui

import (
	"fmt"
	"strings"

	"github.com/0xjuanma/golazo/internal/constants"
	"github.com/0xjuanma/golazo/internal/data"
	"github.com/charmbracelet/lipgloss"
)

var (
	// Settings view styles
	settingsTitleStyle = lipgloss.NewStyle().
				Foreground(accentColor).
				Bold(true).
				Align(lipgloss.Center).
				MarginBottom(1)

	settingsItemStyle = lipgloss.NewStyle().
				Foreground(textColor).
				Padding(0, 1)

	settingsSelectedStyle = lipgloss.NewStyle().
				Foreground(highlightColor).
				Bold(true).
				Padding(0, 1)

	settingsCheckboxChecked = lipgloss.NewStyle().
				Foreground(accentColor).
				Bold(true)

	settingsCheckboxUnchecked = lipgloss.NewStyle().
					Foreground(dimColor)

	settingsCountryStyle = lipgloss.NewStyle().
				Foreground(dimColor).
				Italic(true)

	settingsHelpStyle = lipgloss.NewStyle().
				Foreground(dimColor).
				Align(lipgloss.Center).
				MarginTop(1)

	settingsBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(borderColor).
				Padding(1, 2)

	settingsInfoStyle = lipgloss.NewStyle().
				Foreground(secondaryColor).
				Align(lipgloss.Center).
				MarginTop(1)
)

// SettingsState holds the state for the settings view.
type SettingsState struct {
	Cursor     int          // Currently highlighted item
	Selected   map[int]bool // Map of league ID -> selected
	Leagues    []data.LeagueInfo
	HasChanges bool // Whether there are unsaved changes
}

// NewSettingsState creates a new settings state with current saved preferences.
func NewSettingsState() *SettingsState {
	settings, _ := data.LoadSettings()

	selected := make(map[int]bool)

	// If no leagues are selected in settings, none are checked (default = all)
	// User sees unchecked = will use all leagues
	if len(settings.SelectedLeagues) > 0 {
		for _, id := range settings.SelectedLeagues {
			selected[id] = true
		}
	}

	return &SettingsState{
		Cursor:   0,
		Selected: selected,
		Leagues:  data.AllSupportedLeagues,
	}
}

// Toggle toggles the selection state of the currently highlighted league.
func (s *SettingsState) Toggle() {
	if s.Cursor >= 0 && s.Cursor < len(s.Leagues) {
		leagueID := s.Leagues[s.Cursor].ID
		s.Selected[leagueID] = !s.Selected[leagueID]
		s.HasChanges = true
	}
}

// MoveUp moves the cursor up.
func (s *SettingsState) MoveUp() {
	if s.Cursor > 0 {
		s.Cursor--
	}
}

// MoveDown moves the cursor down.
func (s *SettingsState) MoveDown() {
	if s.Cursor < len(s.Leagues)-1 {
		s.Cursor++
	}
}

// Save persists the current selection to settings.yaml.
func (s *SettingsState) Save() error {
	var selectedIDs []int
	for _, league := range s.Leagues {
		if s.Selected[league.ID] {
			selectedIDs = append(selectedIDs, league.ID)
		}
	}

	settings := &data.Settings{
		SelectedLeagues: selectedIDs,
	}

	err := data.SaveSettings(settings)
	if err == nil {
		s.HasChanges = false
	}
	return err
}

// SelectedCount returns the number of selected leagues.
func (s *SettingsState) SelectedCount() int {
	count := 0
	for _, isSelected := range s.Selected {
		if isSelected {
			count++
		}
	}
	return count
}

// Fixed width for settings box to prevent UI shifting when selections change
const settingsBoxWidth = 52

// RenderSettingsView renders the settings view for league customization.
func RenderSettingsView(width, height int, state *SettingsState) string {
	if state == nil {
		return ""
	}

	// Title - centered within fixed width
	titleStyle := settingsTitleStyle.Width(settingsBoxWidth).Align(lipgloss.Center)
	title := titleStyle.Render("League Preferences")

	// Build the list of leagues
	var items []string
	for i, league := range state.Leagues {
		// Checkbox
		var checkbox string
		if state.Selected[league.ID] {
			checkbox = settingsCheckboxChecked.Render("[x]")
		} else {
			checkbox = settingsCheckboxUnchecked.Render("[ ]")
		}

		// League name and country
		leagueName := league.Name
		country := settingsCountryStyle.Render(fmt.Sprintf("(%s)", league.Country))

		line := fmt.Sprintf("%s %s %s", checkbox, leagueName, country)

		// Apply cursor styling
		if i == state.Cursor {
			items = append(items, settingsSelectedStyle.Render("▸ "+line))
		} else {
			items = append(items, settingsItemStyle.Render("  "+line))
		}
	}

	listContent := strings.Join(items, "\n")

	// Selection info - fixed width and centered to prevent shifting
	selectedCount := state.SelectedCount()
	var infoText string
	if selectedCount == 0 {
		infoText = "No leagues selected - All leagues will be shown"
	} else {
		infoText = fmt.Sprintf("%d of %d leagues selected", selectedCount, len(state.Leagues))
	}
	infoStyle := settingsInfoStyle.Width(settingsBoxWidth).Align(lipgloss.Center)
	info := infoStyle.Render(infoText)

	// Help text
	help := settingsHelpStyle.Render(constants.HelpSettingsView)

	// Combine content with fixed width
	innerContent := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		listContent,
		"",
		info,
	)

	// Add border with fixed width
	borderedStyle := settingsBorderStyle.Width(settingsBoxWidth + 6) // +6 for padding and border
	borderedContent := borderedStyle.Render(innerContent)

	// Final layout with help at bottom
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		borderedContent,
		help,
	)

	// Center in the terminal
	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}
