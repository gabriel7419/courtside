package ui

import (
	"math/rand"
	"strings"
	"time"

	"github.com/0xjuanma/golazo/internal/constants"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
)

// SpinnerTickInterval is the unified tick rate for all spinners (100ms = 10 fps).
// This balances smooth animation with keyboard responsiveness.
const SpinnerTickInterval = 100 * time.Millisecond

// TickMsg is the unified message type for all spinner updates.
// Only ONE tick chain should exist at any time to prevent message queue flooding.
type TickMsg struct{}

// SpinnerTick returns a command that generates a TickMsg after the standard interval.
// This is the ONLY function that should create spinner ticks - ensures single tick chain.
func SpinnerTick() tea.Cmd {
	return tea.Tick(SpinnerTickInterval, func(time.Time) tea.Msg {
		return TickMsg{}
	})
}

// RandomCharSpinner is a custom spinner that cycles through random characters.
// Note: Spinners do NOT self-tick. The app manages the tick chain centrally.
type RandomCharSpinner struct {
	chars      []rune
	currentIdx int
	width      int
	startColor colorful.Color // Gradient start color (cyan)
	endColor   colorful.Color // Gradient end color (red)
}

// NewRandomCharSpinner creates a new random character spinner.
func NewRandomCharSpinner() *RandomCharSpinner {
	// Random characters similar to the image: alphanumeric, symbols, special chars
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()_+-=[]{}|;:,.<>?/~`£€¥")

	// Create gradient: cyan to red (high energy theme)
	startColor, _ := colorful.Hex(constants.GradientStartColor) // Bright cyan
	endColor, _ := colorful.Hex(constants.GradientEndColor)     // Bright red

	return &RandomCharSpinner{
		chars:      chars,
		currentIdx: rand.Intn(len(chars)),
		width:      20, // Default width for spinner
		startColor: startColor,
		endColor:   endColor,
	}
}

// Tick advances the spinner animation state.
// Does NOT return a tick command - the app manages the tick chain.
func (r *RandomCharSpinner) Tick() {
	r.currentIdx = rand.Intn(len(r.chars))
}

// View renders the spinner with gradient colors.
func (r *RandomCharSpinner) View() string {
	if r.width <= 0 {
		r.width = 20
	}

	// Create a string of characters for the spinner
	spinnerChars := make([]rune, r.width)
	for i := range spinnerChars {
		charIdx := (r.currentIdx + i) % len(r.chars)
		spinnerChars[i] = r.chars[charIdx]
	}

	// Apply gradient to each character
	var result strings.Builder
	for i, char := range spinnerChars {
		ratio := float64(i) / float64(r.width-1)
		color := r.startColor.BlendLab(r.endColor, ratio)
		hexColor := color.Hex()
		charStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(hexColor))
		result.WriteString(charStyle.Render(string(char)))
	}

	return result.String()
}

// SetWidth sets the width of the spinner.
func (r *RandomCharSpinner) SetWidth(width int) {
	r.width = width
}
