package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/0xjuanma/golazo/internal/data"
)

type LeagueInfo struct {
	ID      int
	Name    string
	Country string
	Count   int
}

func main() {
	// Get API key
	apiKey, err := data.FootballDataAPIKey()
	if err != nil {
		fmt.Printf("❌ Error getting API key: %v\n", err)
		os.Exit(1)
	}

	baseURL := "https://v3.football.api-sports.io"
	today := time.Now().Format("2006-01-02")

	fmt.Printf("Fetching all finished matches for %s...\n\n", today)
	url := fmt.Sprintf("%s/fixtures?date=%s&status=FT", baseURL, today)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("x-apisports-key", apiKey)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Printf("❌ Status %d: %s\n", resp.StatusCode, string(bodyBytes))
		os.Exit(1)
	}

	var result map[string]interface{}
	bodyBytes, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		fmt.Printf("❌ Parse error: %v\n", err)
		os.Exit(1)
	}

	response, ok := result["response"].([]interface{})
	if !ok {
		fmt.Printf("❌ Unexpected response structure\n")
		os.Exit(1)
	}

	fmt.Printf("✓ Found %d finished matches\n\n", len(response))

	// Extract league information
	leagueMap := make(map[int]*LeagueInfo)

	for _, match := range response {
		if m, ok := match.(map[string]interface{}); ok {
			if league, ok := m["league"].(map[string]interface{}); ok {
				if leagueID, ok := league["id"].(float64); ok {
					id := int(leagueID)
					if leagueMap[id] == nil {
						leagueMap[id] = &LeagueInfo{ID: id}
					}
					leagueMap[id].Count++

					if name, ok := league["name"].(string); ok && leagueMap[id].Name == "" {
						leagueMap[id].Name = name
					}
					if country, ok := league["country"].(string); ok && leagueMap[id].Country == "" {
						leagueMap[id].Country = country
					}
				}
			}
		}
	}

	// Convert to slice and sort by count
	leagues := make([]*LeagueInfo, 0, len(leagueMap))
	for _, league := range leagueMap {
		leagues = append(leagues, league)
	}

	sort.Slice(leagues, func(i, j int) bool {
		return leagues[i].Count > leagues[j].Count
	})

	// Print all leagues
	fmt.Println("═══════════════════════════════════════════════════════════════════════════════")
	fmt.Println("ALL LEAGUES FOUND IN TODAY'S FINISHED MATCHES")
	fmt.Println("═══════════════════════════════════════════════════════════════════════════════")
	fmt.Printf("%-6s | %-50s | %-20s | %s\n", "ID", "League Name", "Country", "Matches")
	fmt.Println("───────────────────────────────────────────────────────────────────────────────")
	for _, league := range leagues {
		fmt.Printf("%-6d | %-50s | %-20s | %d\n", league.ID, league.Name, league.Country, league.Count)
	}
	fmt.Println("═══════════════════════════════════════════════════════════════════════════════\n")

	// Check our supported leagues
	fmt.Println("CHECKING OUR SUPPORTED LEAGUES:")
	fmt.Println("───────────────────────────────────────────────────────────────────────────────")
	supportedLeagues := map[int]string{
		39:  "Premier League",
		140: "La Liga",
		78:  "Bundesliga",
		135: "Serie A",
		61:  "Ligue 1",
	}

	for id, expectedName := range supportedLeagues {
		found := false

		// Check by ID
		if league, exists := leagueMap[id]; exists {
			found = true
			if league.Name == expectedName {
				fmt.Printf("✓ ID %3d: %s (%s) - %d matches [EXACT MATCH]\n", id, league.Name, league.Country, league.Count)
			} else {
				fmt.Printf("⚠ ID %3d: Found '%s' (%s) - %d matches [Expected: %s]\n", id, league.Name, league.Country, league.Count, expectedName)
			}
		} else {
			// Check by name
			for _, league := range leagues {
				if league.Name == expectedName {
					found = true
					fmt.Printf("✗ ID %3d: NOT FOUND, but found '%s' with ID %d (%s) - %d matches\n",
						id, league.Name, league.ID, league.Country, league.Count)
					break
				}
			}
		}

		if !found {
			fmt.Printf("✗ ID %3d: %s - NOT FOUND in today's matches\n", id, expectedName)
		}
	}

	fmt.Println("\n═══════════════════════════════════════════════════════════════════════════════")
	fmt.Println("RECOMMENDED LEAGUE IDs FOR API-SPORTS.IO:")
	fmt.Println("═══════════════════════════════════════════════════════════════════════════════")

	// Find the correct IDs for our supported leagues
	foundCount := 0
	for _, expectedName := range supportedLeagues {
		for _, league := range leagues {
			if league.Name == expectedName {
				fmt.Printf("%d, // %s (%s) - %d matches today\n", league.ID, league.Name, league.Country, league.Count)
				foundCount++
				break
			}
		}
	}

	if foundCount == 0 {
		fmt.Println("\n⚠ None of our expected leagues found. Checking for similar names...")
		searchTerms := []string{"Premier", "Liga", "Bundesliga", "Serie", "Ligue"}
		for _, term := range searchTerms {
			for _, league := range leagues {
				if len(league.Name) > 0 &&
					(league.Name == term ||
						(len(term) <= len(league.Name) &&
							(league.Name[:min(len(term), len(league.Name))] == term ||
								league.Name[len(league.Name)-min(len(term), len(league.Name)):] == term))) {
					fmt.Printf("  Similar: ID %d - %s (%s) - %d matches\n",
						league.ID, league.Name, league.Country, league.Count)
				}
			}
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
