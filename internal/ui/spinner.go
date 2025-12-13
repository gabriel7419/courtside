package ui

import "github.com/charmbracelet/lipgloss"

// SpinnerStyle returns the style for the spinner.
func SpinnerStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(accentColor)
}
