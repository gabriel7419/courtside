package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gabriel7419/courtside/internal/api"
	"github.com/gabriel7419/courtside/internal/constants"
)

const standingsDialogID = "standings"

// StandingsDialog displays the NBA conference standings table.
type StandingsDialog struct {
	leagueName  string
	standings   []api.LeagueTableEntry
	homeTeamID  int
	awayTeamID  int
	scrollIndex int
}

// NewStandingsDialog creates a new standings dialog.
func NewStandingsDialog(leagueName string, standings []api.LeagueTableEntry, homeTeamID, awayTeamID int) *StandingsDialog {
	return &StandingsDialog{
		leagueName:  leagueName,
		standings:   standings,
		homeTeamID:  homeTeamID,
		awayTeamID:  awayTeamID,
		scrollIndex: 0,
	}
}

// ID returns the dialog identifier.
func (d *StandingsDialog) ID() string {
	return standingsDialogID
}

// Update handles input for the standings dialog.
func (d *StandingsDialog) Update(msg tea.Msg) (Dialog, DialogAction) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "s", "q":
			return d, DialogActionClose{}
		case "j", "down":
			if d.scrollIndex < len(d.standings)-1 {
				d.scrollIndex++
			}
		case "k", "up":
			if d.scrollIndex > 0 {
				d.scrollIndex--
			}
		}
	}
	return d, nil
}

// View renders the standings table.
func (d *StandingsDialog) View(width, height int) string {
	dialogWidth, dialogHeight := DialogSize(width, height, 90, 36)
	content := d.renderTable(dialogWidth - 6)
	return RenderDialogFrameWithHelp(d.leagueName+" Standings", content, constants.HelpStandingsDialog, dialogWidth, dialogHeight)
}

// NBA standings column widths
const (
	nbColPos    = 4  // "#"  rank
	nbColTeam   = 20 // team abbreviation/name
	nbColW      = 4  // W
	nbColL      = 4  // L
	nbColPct    = 6  // .xxx %
	nbColGB     = 6  // games back
	nbColStreak = 6  // W3 / L2
)

func (d *StandingsDialog) renderTable(width int) string {
	if len(d.standings) == 0 {
		return dialogDimStyle.Render("No standings data available")
	}

	var lines []string

	// Conference group headers (East / West)
	var prevConf string
	first := true

	lines = append(lines, d.renderHeaderRow(width))
	lines = append(lines, dialogSeparatorStyle.Render(strings.Repeat("─", width)))

	for _, entry := range d.standings {
		// Parse conference from Note field ("East | GB: 3.5")
		conf := ""
		if parts := strings.SplitN(entry.Note, " | ", 2); len(parts) == 2 {
			conf = parts[0]
		}

		// Insert conference sub-header on change
		if conf != "" && conf != prevConf {
			if !first {
				lines = append(lines, "")
			}
			confLabel := lipgloss.NewStyle().
				Foreground(neonCyan).
				Bold(true).
				Width(width).
				Render("  " + conf + "ern Conference")
			lines = append(lines, confLabel)
			lines = append(lines, dialogSeparatorStyle.Render(strings.Repeat("─", width)))
			prevConf = conf
			first = false
		}

		lines = append(lines, d.renderTeamRow(entry, width))
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

// renderHeaderRow renders the table header with NBA columns.
func (d *StandingsDialog) renderHeaderRow(width int) string {
	return lipgloss.JoinHorizontal(lipgloss.Top,
		dialogHeaderStyle.Width(nbColPos).Align(lipgloss.Right).Render("#"),
		"  ",
		dialogHeaderStyle.Width(nbColTeam).Align(lipgloss.Left).Render("Team"),
		dialogHeaderStyle.Width(nbColW).Align(lipgloss.Right).Render("W"),
		dialogHeaderStyle.Width(nbColL).Align(lipgloss.Right).Render("L"),
		dialogHeaderStyle.Width(nbColPct).Align(lipgloss.Right).Render("PCT"),
		dialogHeaderStyle.Width(nbColGB).Align(lipgloss.Right).Render("GB"),
		dialogHeaderStyle.Width(nbColStreak).Align(lipgloss.Right).Render("Strk"),
	)
}

// renderTeamRow renders a single team row with NBA columns.
func (d *StandingsDialog) renderTeamRow(entry api.LeagueTableEntry, width int) string {
	isHighlighted := entry.Team.ID == d.homeTeamID || entry.Team.ID == d.awayTeamID

	// Team display: prefer abbreviation
	teamName := entry.Team.ShortName
	if teamName == "" {
		teamName = entry.Team.Name
	}
	if len(teamName) > nbColTeam-1 {
		teamName = teamName[:nbColTeam-2] + "…"
	}

	// Win percentage from PointsFor (stored as win% × 1000)
	pctStr := "—"
	if entry.Played > 0 {
		pct := float64(entry.PointsFor) / 1000.0
		pctStr = fmt.Sprintf(".%03d", int(pct*1000)%1000)
		if pct >= 1.0 {
			pctStr = "1.000"
		}
	}

	// Games behind
	gbStr := "—"
	if parts := strings.SplitN(entry.Note, "GB: ", 2); len(parts) == 2 {
		gbStr = parts[1]
		if gbStr == "0" || gbStr == "0.0" {
			gbStr = "—"
		}
	}

	// Streak
	streak := entry.Form
	if streak == "" {
		streak = "—"
	}

	rowContent := lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.NewStyle().Width(nbColPos).Align(lipgloss.Right).Render(fmt.Sprintf("%d", entry.Position)),
		"  ",
		lipgloss.NewStyle().Width(nbColTeam).Align(lipgloss.Left).Render(teamName),
		lipgloss.NewStyle().Width(nbColW).Align(lipgloss.Right).Render(fmt.Sprintf("%d", entry.Won)),
		lipgloss.NewStyle().Width(nbColL).Align(lipgloss.Right).Render(fmt.Sprintf("%d", entry.Lost)),
		lipgloss.NewStyle().Width(nbColPct).Align(lipgloss.Right).Render(pctStr),
		lipgloss.NewStyle().Width(nbColGB).Align(lipgloss.Right).Render(gbStr),
		lipgloss.NewStyle().Width(nbColStreak).Align(lipgloss.Right).Render(streak),
	)

	if isHighlighted {
		return lipgloss.NewStyle().
			Background(neonDark).
			Foreground(neonCyan).
			Bold(true).
			Width(width).
			Render(rowContent)
	}

	return dialogValueStyle.Render(rowContent)
}

// formatGoalDifference formats goal difference with +/- sign (kept for football compatibility).
func formatGoalDifference(gd int) string {
	if gd > 0 {
		return fmt.Sprintf("+%d", gd)
	}
	return fmt.Sprintf("%d", gd)
}
