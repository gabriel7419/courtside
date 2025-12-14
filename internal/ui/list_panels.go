package ui

import (
	"github.com/0xjuanma/golazo/internal/api"
	"github.com/0xjuanma/golazo/internal/constants"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

// RenderLiveMatchesListPanel renders the left panel using bubbletea list component.
func RenderLiveMatchesListPanel(width, height int, listModel list.Model) string {
	// Calculate available space for list (accounting for panel border, padding, and title)
	h, v := panelStyle.GetFrameSize()
	titleHeight := 3 // Title + spacing
	availableWidth := width - h
	availableHeight := height - v - titleHeight

	// Set list size
	listModel.SetSize(availableWidth, availableHeight)

	// Wrap list in panel
	title := panelTitleStyle.Width(width - 6).Render(constants.PanelLiveMatches)
	listView := listModel.View()

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		listView,
	)

	panel := panelStyle.
		Width(width).
		Height(height).
		Render(content)

	return panel
}

// RenderStatsListPanel renders the left panel for stats view using bubbletea list component.
func RenderStatsListPanel(width, height int, listModel list.Model) string {
	// Calculate available space for list (accounting for panel border, padding, and title)
	h, v := panelStyle.GetFrameSize()
	titleHeight := 3 // Title + spacing
	availableWidth := width - h
	availableHeight := height - v - titleHeight

	// Set list size
	listModel.SetSize(availableWidth, availableHeight)

	// Wrap list in panel
	title := panelTitleStyle.Width(width - 6).Render(constants.PanelFinishedMatches)
	listView := listModel.View()

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		listView,
	)

	panel := panelStyle.
		Width(width).
		Height(height).
		Render(content)

	return panel
}

// RenderMultiPanelViewWithList renders the live matches view with list component.
func RenderMultiPanelViewWithList(width, height int, listModel list.Model, details *api.MatchDetails, liveUpdates []string, sp spinner.Model, loading bool) string {
	// Calculate panel dimensions
	leftWidth := width * 35 / 100
	if leftWidth < 25 {
		leftWidth = 25
	}
	rightWidth := width - leftWidth - 1
	if rightWidth < 35 {
		rightWidth = 35
		leftWidth = width - rightWidth - 1
	}

	// Render left panel (matches list)
	leftPanel := RenderLiveMatchesListPanel(leftWidth, height, listModel)

	// Render right panel (match details with live updates)
	rightPanel := renderMatchDetailsPanel(rightWidth, height, details, liveUpdates, sp, loading)

	// Create separator
	separatorStyle := lipgloss.NewStyle().
		Foreground(borderColor).
		Height(height).
		Padding(0, 1)
	separator := separatorStyle.Render("│")

	// Combine panels
	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPanel,
		separator,
		rightPanel,
	)

	return content
}

// RenderStatsViewWithList renders the stats view with list component.
func RenderStatsViewWithList(width, height int, listModel list.Model, details *api.MatchDetails) string {
	// Calculate panel dimensions
	leftWidth := width * 40 / 100
	if leftWidth < 30 {
		leftWidth = 30
	}
	rightWidth := width - leftWidth - 1
	if rightWidth < 40 {
		rightWidth = 40
		leftWidth = width - rightWidth - 1
	}

	panelHeight := height - 2

	// Render left panel (finished matches list)
	leftPanel := RenderStatsListPanel(leftWidth, panelHeight, listModel)

	// Render right panel (match stats)
	rightPanel := renderMatchStatsPanel(rightWidth, panelHeight, details)

	// Create separator
	separatorStyle := lipgloss.NewStyle().
		Foreground(borderColor).
		Height(panelHeight).
		Padding(0, 1)
	separator := separatorStyle.Render("│")

	// Combine panels
	combined := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPanel,
		separator,
		rightPanel,
	)

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, combined)
}
