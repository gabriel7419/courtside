// api_results_raw.go - Debug script to inspect raw FotMob API responses and struct conversions
//
// Usage:
//   go run scripts/api_results_raw.go <league_id> <date> [season]
//
// Examples:
//   go run scripts/api_results_raw.go 47 2024-12-24           # Premier League matches on Dec 24
//   go run scripts/api_results_raw.go 87 today                # La Liga matches today
//   go run scripts/api_results_raw.go 54 yesterday            # Bundesliga matches yesterday
//   go run scripts/api_results_raw.go 77 2022-12-18 2022      # World Cup 2022 Final (Argentina vs France)
//   go run scripts/api_results_raw.go 77 2022-11-22 2022      # World Cup 2022 group stage
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
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gabriel7419/courtside/internal/debug"
	"github.com/goforj/godump"
)

const (
	maxMatches = 3 // Limit matches to display for readability
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

	// Optional season parameter (e.g., "2022" for World Cup 2022, "2024/2025" for leagues)
	season := ""
	if len(os.Args) >= 4 {
		season = os.Args[3]
	}

	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║          FotMob API Response Inspector (Raw)                 ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Printf("League ID: %d\n", leagueID)
	fmt.Printf("Date: %s\n", date.Format("2006-01-02 (Monday)"))
	if season != "" {
		fmt.Printf("Season: %s\n", season)
	}
	fmt.Printf("Max matches to display: %d\n", maxMatches)
	fmt.Println()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// We'll make 2 API calls: fixtures (upcoming/scheduled) and results (finished)
	tabs := []string{"fixtures", "results"}

	for _, tab := range tabs {
		fmt.Println(strings.Repeat("═", 70))
		fmt.Printf("API Call: tab=%s\n", tab)
		fmt.Println(strings.Repeat("═", 70))
		fmt.Println()

		// Make raw API call to get the full response
		rawResponse, matches, err := debug.FetchLeagueData(ctx, leagueID, date, tab, season)
		if err != nil {
			fmt.Printf("Error fetching %s: %v\n\n", tab, err)
			continue
		}

		// 1. Print raw API response (pretty JSON)
		fmt.Println("┌─────────────────────────────────────────────────────────────────┐")
		fmt.Println("│ 1. RAW API RESPONSE (JSON)                                      │")
		fmt.Println("└─────────────────────────────────────────────────────────────────┘")
		fmt.Println()

		// Pretty print the raw JSON
		prettyJSON, err := json.MarshalIndent(rawResponse, "", "  ")
		if err != nil {
			fmt.Printf("Error formatting JSON: %v\n", err)
		} else {
			// Truncate if too long, but show structure
			jsonStr := string(prettyJSON)
			if len(jsonStr) > 5000 {
				// Show first 2500 chars and last 500
				fmt.Println(jsonStr[:2500])
				fmt.Println("\n... [truncated for readability] ...\n")
				fmt.Println(jsonStr[len(jsonStr)-500:])
			} else {
				fmt.Println(jsonStr)
			}
		}
		fmt.Println()

		// 2. Print converted internal struct using godump
		fmt.Println("┌─────────────────────────────────────────────────────────────────┐")
		fmt.Printf("│ 2. CONVERTED MATCHES (api.Match structs) - %d found             \n", len(matches))
		fmt.Println("└─────────────────────────────────────────────────────────────────┘")
		fmt.Println()

		if len(matches) == 0 {
			fmt.Printf("   [!] No matches found for this date in the %s tab\n", tab)
		} else {
			// Limit matches for readability
			displayMatches := matches
			if len(displayMatches) > maxMatches {
				displayMatches = displayMatches[:maxMatches]
				fmt.Printf("   (Showing %d of %d matches)\n\n", maxMatches, len(matches))
			}

			for i, match := range displayMatches {
				fmt.Printf("   ╔═══ Match %d ═══╗\n", i+1)
				// Use godump to pretty-print the struct
				godump.Dump(match)
				fmt.Println()
			}
		}
		fmt.Println()
	}

	// Bonus: Fetch match details for first available match
	fmt.Println(strings.Repeat("═", 70))
	fmt.Println("BONUS: Match Details API Call")
	fmt.Println(strings.Repeat("═", 70))
	fmt.Println()

	// Try to get a match to fetch details for (use our season-aware function)
	_, allMatches, _ := debug.FetchLeagueData(ctx, leagueID, date, "results", season)
	if len(allMatches) == 0 {
		_, allMatches, _ = debug.FetchLeagueData(ctx, leagueID, date, "fixtures", season)
	}

	if len(allMatches) == 0 {
		fmt.Println("   [!] No matches found to fetch details for")
	} else {
		matchID := allMatches[0].ID
		fmt.Printf("   Fetching details for match ID: %d\n", matchID)
		fmt.Printf("   (%s vs %s)\n\n", allMatches[0].HomeTeam.Name, allMatches[0].AwayTeam.Name)

		// Raw API call for match details
		rawDetails, details, err := debug.FetchMatchDetails(ctx, matchID)
		if err != nil {
			fmt.Printf("   Error fetching match details: %v\n", err)
		} else {
			// Print raw response
			fmt.Println("┌─────────────────────────────────────────────────────────────────┐")
			fmt.Println("│ 1. RAW MATCH DETAILS RESPONSE (JSON - truncated)                │")
			fmt.Println("└─────────────────────────────────────────────────────────────────┘")
			fmt.Println()

			prettyJSON, _ := json.MarshalIndent(rawDetails, "", "  ")
			jsonStr := string(prettyJSON)
			if len(jsonStr) > 4000 {
				fmt.Println(jsonStr[:2000])
				fmt.Println("\n... [truncated for readability] ...\n")
				fmt.Println(jsonStr[len(jsonStr)-500:])
			} else {
				fmt.Println(jsonStr)
			}
			fmt.Println()

			// Print converted struct
			fmt.Println("┌─────────────────────────────────────────────────────────────────┐")
			fmt.Println("│ 2. CONVERTED MATCH DETAILS (api.MatchDetails struct)            │")
			fmt.Println("└─────────────────────────────────────────────────────────────────┘")
			fmt.Println()
			godump.Dump(details)
		}
	}

	fmt.Println()
	fmt.Println("Done!")
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
Usage: go run scripts/api_results_raw.go <league_id> <date> [season]

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
  go run scripts/api_results_raw.go 47 today
  go run scripts/api_results_raw.go 87 2024-12-24
  go run scripts/api_results_raw.go 54 yesterday
  go run scripts/api_results_raw.go 77 2022-12-18 2022    # World Cup 2022 Final
  go run scripts/api_results_raw.go 47 2024-12-21 2024/2025  # Premier League specific season
`)
}
