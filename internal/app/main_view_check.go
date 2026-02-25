package app

import (
	"time"

	"github.com/gabriel7419/courtside/internal/constants"
	tea "github.com/charmbracelet/bubbletea"
)

// mainViewCheckMsg is sent after the check delay completes.
type mainViewCheckMsg struct {
	selection int // 0 for Stats, 1 for Live Matches
}

// performMainViewCheck performs a delay check before navigating.
func performMainViewCheck(selection int) tea.Cmd {
	return tea.Tick(constants.MainViewCheckDelay, func(t time.Time) tea.Msg {
		return mainViewCheckMsg{selection: selection}
	})
}
