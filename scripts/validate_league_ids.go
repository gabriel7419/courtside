package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gabriel7419/courtside/internal/data"
)

// LeagueInfo and region constants are imported from the data package

// AllSupportedLeagues is imported from the data package to avoid duplication

func fetchLeagueName(leagueID int) (string, error) {
	url := fmt.Sprintf("https://www.fotmob.com/leagues/%d/overview", leagueID)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return extractLeagueName(string(body))
}

func extractLeagueName(html string) (string, error) {
	// Look for the league name in the title tag
	titleRegex := regexp.MustCompile(`<title[^>]*>([^<]+)</title>`)
	matches := titleRegex.FindStringSubmatch(html)
	if len(matches) > 1 {
		title := matches[1]
		// Extract league name from title (remove "matches, tables and news" suffix)
		parts := strings.Split(title, " matches, tables and news")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0]), nil
		}
	}

	return "", fmt.Errorf("league name not found in HTML")
}

func validateLeague(league data.LeagueInfo) string {
	fetchedName, err := fetchLeagueName(league.ID)
	if err != nil {
		return fmt.Sprintf("ERROR: %v", err)
	}

	if strings.TrimSpace(fetchedName) == strings.TrimSpace(league.Name) {
		return "VALID"
	}

	return fmt.Sprintf("INVALID (got: %s)", fetchedName)
}

func main() {
	fmt.Printf("%-8s %-30s %-15s %s\n", "ID", "TRACKED NAME", "COUNTRY", "RESULT")
	fmt.Println(strings.Repeat("-", 100))

	for region, leagues := range data.AllSupportedLeagues {
		fmt.Printf("\n%s:\n", region)
		for _, league := range leagues {
			result := validateLeague(league)
			fmt.Printf("%-8d %-30s %-15s %s\n", league.ID, league.Name, league.Country, result)

			// Small delay to be respectful to the server
			time.Sleep(100 * time.Millisecond)
		}
	}
}
