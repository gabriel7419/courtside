// api_results_render.go - Render fetched match data using Golazo's stats UI (right panel)
//
// Usage:
//   go run scripts/api_results_render.go <league_id> <date> [season]
//
// Examples:
//   go run scripts/api_results_render.go 47 2024-12-24           # Premier League matches on Dec 24
//   go run scripts/api_results_render.go 87 today                # La Liga matches today
//   go run scripts/api_results_render.go 77 2022-12-18 2022      # World Cup 2022 Final (Argentina vs France)
//
// This script fetches match data and renders it using the same UI components
// as the Golazo app's stats view (right panel with match details).
//
// League IDs (common):
//   47 - Premier League
//   87 - La Liga
//   54 - Bundesliga
//   55 - Serie A
//   53 - Ligue 1
//   42 - Champions League
//   73 - Europa League
//   77 - FIFA World Cup
//   530 - MLS
//   325 - Liga MX

package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gabriel7419/courtside/internal/debug"
	"github.com/gabriel7419/courtside/internal/ui"
	"golang.org/x/term"
)

func main() {
	if len(os.Args) < 3 {
		printUsage()
		os.Exit(1)
	}

	// Parse arguments
	leagueID, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("Error: Invalid league_id: %s (must be a number)\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}

	date, err := parseDate(os.Args[2])
	if err != nil {
		fmt.Printf("Error: Invalid date: %s (%v)\n", os.Args[2], err)
		printUsage()
		os.Exit(1)
	}

	// Optional season parameter
	season := ""
	if len(os.Args) >= 4 {
		season = os.Args[3]
	}

	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║          Golazo Stats Panel Renderer                         ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Printf("League ID: %d\n", leagueID)
	fmt.Printf("Date: %s\n", date.Format("2006-01-02 (Monday)"))
	if season != "" {
		fmt.Printf("Season: %s\n", season)
	}
	fmt.Println()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Fetch matches
	fmt.Println("Fetching matches...")
	_, matches, _ := debug.FetchLeagueData(ctx, leagueID, date, "results", season)
	if len(matches) == 0 {
		_, matches, _ = debug.FetchLeagueData(ctx, leagueID, date, "fixtures", season)
	}

	if len(matches) == 0 {
		fmt.Println("[!] No matches found for this date/league/season combination")
		os.Exit(0)
	}

	fmt.Printf("Found %d matches\n\n", len(matches))

	// Get terminal dimensions for rendering
	width, height := getTerminalSize()
	panelWidth := width - 4 // Leave some margin
	if panelWidth > 80 {
		panelWidth = 80 // Cap at reasonable width
	}
	panelHeight := height - 10
	if panelHeight < 30 {
		panelHeight = 30
	}

	// Render each match using the Golazo stats panel UI
	for i, match := range matches {
		fmt.Println(strings.Repeat("═", panelWidth))
		fmt.Printf("Match %d of %d: %s vs %s\n", i+1, len(matches), match.HomeTeam.Name, match.AwayTeam.Name)
		fmt.Println(strings.Repeat("═", panelWidth))
		fmt.Println()

		// Fetch match details
		fmt.Printf("Fetching details for match ID %d...\n", match.ID)
		_, details, err := debug.FetchMatchDetails(ctx, match.ID)
		if err != nil {
			fmt.Printf("Error fetching details: %v\n\n", err)
			continue
		}

		if details == nil {
			fmt.Println("[!] Could not retrieve match details\n")
			continue
		}

		// Render using Golazo's stats panel UI
		rendered := ui.RenderMatchDetailsPanel(panelWidth, panelHeight, details)
		fmt.Println(rendered)
		fmt.Println()

		// Only render first 3 matches to avoid overwhelming output
		if i >= 2 {
			if len(matches) > 3 {
				fmt.Printf("... and %d more matches (showing first 3)\n\n", len(matches)-3)
			}
			break
		}
	}

	fmt.Println("Done!")
}

// getTerminalSize returns the terminal width and height, with sensible defaults
func getTerminalSize() (int, int) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width <= 0 {
		width = 100 // Default width
	}
	if height <= 0 {
		height = 40 // Default height
	}
	return width, height
}

// parseDate parses date string in various formats
func parseDate(dateStr string) (time.Time, error) {
	dateStr = strings.ToLower(strings.TrimSpace(dateStr))

	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	switch dateStr {
	case "today":
		return today, nil
	case "yesterday":
		return today.AddDate(0, 0, -1), nil
	case "tomorrow":
		return today.AddDate(0, 0, 1), nil
	}

	// Try parsing as YYYY-MM-DD
	if t, err := time.Parse("2006-01-02", dateStr); err == nil {
		return t, nil
	}

	// Try parsing as MM-DD (assume current year)
	if t, err := time.Parse("01-02", dateStr); err == nil {
		return time.Date(now.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC), nil
	}

	return time.Time{}, fmt.Errorf("use format YYYY-MM-DD, or 'today', 'yesterday', 'tomorrow'")
}

func printUsage() {
	fmt.Println(`
Usage: go run scripts/api_results_render.go <league_id> <date> [season]

Renders fetched match data using Golazo's stats UI (right panel).

Arguments:
  league_id   The FotMob league ID (integer)
  date        Date in YYYY-MM-DD format, or 'today', 'yesterday', 'tomorrow'
  season      Optional: Season identifier (e.g., "2022" for World Cup, "2024/2025" for leagues)

Common League IDs:
  47  - Premier League (England)
  87  - La Liga (Spain)
  54  - Bundesliga (Germany)
  55  - Serie A (Italy)
  53  - Ligue 1 (France)
  42  - Champions League
  73  - Europa League
  77  - FIFA World Cup
  530 - MLS (USA)
  325 - Liga MX (Mexico)

Examples:
  go run scripts/api_results_render.go 47 today
  go run scripts/api_results_render.go 87 2024-12-24
  go run scripts/api_results_render.go 77 2022-12-18 2022    # World Cup 2022 Final
`)
}
