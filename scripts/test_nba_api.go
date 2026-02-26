// NBA Stats API test script — hit endpoints and inspect raw responses.
// Usage: go run scripts/test_nba_api.go [flags]

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gabriel7419/courtside/internal/nba"
)

const (
	baseURL = "https://stats.nba.com/stats"
)

type NBAClient struct {
	httpClient *http.Client
}

func NewNBAClient() *NBAClient {
	return &NBAClient{
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (c *NBAClient) makeRequest(url string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// Required headers — without these the API returns 403 or times out.
	// Mirrors the headers used by the nba_api Python library.
	req.Header.Set("Host", "stats.nba.com")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	// Note: do NOT set Accept-Encoding — Go's http.Client handles gzip transparently
	req.Header.Set("Referer", "https://www.nba.com/")
	req.Header.Set("Origin", "https://www.nba.com")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("x-nba-stats-origin", "stats")
	req.Header.Set("x-nba-stats-token", "true")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")

	fmt.Printf("🔄 Making request to: %s\n", url)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	fmt.Printf("✅ Response status: %s\n", resp.Status)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parse JSON: %w", err)
	}

	return result, nil
}

func (c *NBAClient) GetScoreboard(date string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/scoreboardv3?GameDate=%s&LeagueID=00", baseURL, date)
	return c.makeRequest(url)
}

func (c *NBAClient) GetBoxScoreSummary(gameID string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/boxscoresummaryv2?GameID=%s", baseURL, gameID)
	return c.makeRequest(url)
}

func (c *NBAClient) GetBoxScoreTraditional(gameID string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/boxscoretraditionalv3?GameID=%s&StartPeriod=1&EndPeriod=10&StartRange=0&EndRange=28800&RangeType=0", baseURL, gameID)
	return c.makeRequest(url)
}

func (c *NBAClient) GetStandings(season string) (map[string]interface{}, error) {
	// season format: "2025-26"
	url := fmt.Sprintf("%s/leaguestandingsv3?LeagueID=00&Season=%s&SeasonType=Regular+Season", baseURL, season)
	return c.makeRequest(url)
}

func (c *NBAClient) GetPlayByPlay(gameID string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/playbyplayv3?GameID=%s&StartPeriod=1&EndPeriod=10", baseURL, gameID)
	return c.makeRequest(url)
}

func printJSON(data interface{}) {
	pretty, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("Error formatting JSON: %v\n", err)
		return
	}
	fmt.Println(string(pretty))
}

func saveJSON(filename string, data interface{}) error {
	pretty, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, pretty, 0644)
}

func main() {
	// Flags
	endpoint := flag.String("endpoint", "scoreboard", "Endpoint: scoreboard, summary, traditional, playbyplay, standings")
	date := flag.String("date", time.Now().Format("2006-01-02"), "Date for scoreboard (YYYY-MM-DD)")
	gameID := flag.String("game", "", "Game ID for detailed endpoints")
	season := flag.String("season", currentSeason(), "Season for standings (e.g. 2025-26)")
	output := flag.String("output", "", "Save output to file")
	flag.Parse()

	client := NewNBAClient()

	if *endpoint == "test_internal" {
		client := nba.NewClient()
		ctx := context.Background()
		date := time.Now().UTC().AddDate(0, 0, -1)
		fmt.Printf("Fetching MatchesByDate for %v...\n", date.Format("2006-01-02"))
		matches, err := client.MatchesByDate(ctx, date)
		if err != nil {
			fmt.Printf("ERROR MatchesByDate: %v\n", err)
			return
		}
		if len(matches) > 0 {
			m := matches[0]
			fmt.Printf("Found match: %s vs %s (ID: %d). Fetching MatchDetails...\n", m.HomeTeam.ShortName, m.AwayTeam.ShortName, m.ID)
			details, err := client.MatchDetails(ctx, m.ID, &m)
			if err != nil {
				fmt.Printf("ERROR MatchDetails: %v\n", err)
			} else {
				fmt.Printf("SUCCESS:\n")
				fmt.Printf("  Status: %v\n", details.Status)
				fmt.Printf("  LiveTime: %v\n", *details.LiveTime)
				fmt.Printf("  Home: %s, Score: %v\n", details.HomeTeam.ShortName, details.HomeScore)
				fmt.Printf("  Away: %s, Score: %v\n", details.AwayTeam.ShortName, details.AwayScore)
				if len(details.QuarterScores) > 0 {
					fmt.Printf("  Quarters: %v\n", details.QuarterScores)
				}
			}
		} else {
			fmt.Println("No matches found for yesterday.")
		}

		fmt.Println("\n--- Fetching Raw boxscoretraditionalv3 ---")
		if len(matches) > 0 {
			gameIDStr := fmt.Sprintf("%010d", matches[0].ID)
			url := fmt.Sprintf("https://stats.nba.com/stats/boxscoretraditionalv3?GameID=%s&StartPeriod=1&EndPeriod=10&StartRange=0&EndRange=28800&RangeType=0", gameIDStr)
			req, _ := http.NewRequest("GET", url, nil)
			req.Header.Set("User-Agent", "Mozilla/5.0")
			req.Header.Set("Referer", "https://www.nba.com/")
			req.Header.Set("x-nba-stats-origin", "stats")
			req.Header.Set("x-nba-stats-token", "true")
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				defer resp.Body.Close()
				body, _ := io.ReadAll(resp.Body)
				fmt.Println(string(body)[:2000]) // print first 2000 chars
			}
		}

		return
	}

	var result map[string]interface{}
	var err error

	fmt.Println("🏀 NBA Stats API Tester")
	fmt.Println("========================\n")

	switch *endpoint {
	case "scoreboard":
		fmt.Printf("📅 Fetching scoreboard for: %s\n\n", *date)
		result, err = client.GetScoreboard(*date)

	case "summary":
		if *gameID == "" {
			fmt.Println("❌ Error: --game flag is required for summary endpoint")
			fmt.Println("Example: go run scripts/test_nba_api.go --endpoint=summary --game=0022300789")
			os.Exit(1)
		}
		fmt.Printf("📊 Fetching box score summary for game: %s\n\n", *gameID)
		result, err = client.GetBoxScoreSummary(*gameID)

	case "traditional":
		if *gameID == "" {
			fmt.Println("❌ Error: --game flag is required for traditional endpoint")
			fmt.Println("Example: go run scripts/test_nba_api.go --endpoint=traditional --game=0022300789")
			os.Exit(1)
		}
		fmt.Printf("📈 Fetching traditional box score for game: %s\n\n", *gameID)
		result, err = client.GetBoxScoreTraditional(*gameID)

	case "standings":
		fmt.Printf("📊 Fetching standings for season: %s\n\n", *season)
		result, err = client.GetStandings(*season)

	case "playbyplay":
		if *gameID == "" {
			fmt.Println("❌ Error: --game flag is required for playbyplay endpoint")
			fmt.Println("Example: go run scripts/test_nba_api.go --endpoint=playbyplay --game=0022300789")
			os.Exit(1)
		}
		fmt.Printf("⏱️  Fetching play-by-play for game: %s\n\n", *gameID)
		result, err = client.GetPlayByPlay(*gameID)

	default:
		fmt.Printf("❌ Unknown endpoint: %s\n", *endpoint)
		fmt.Println("Valid endpoints: scoreboard, summary, traditional, playbyplay, standings")
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n✅ Success! Data retrieved.\n")

	// Print summary of response
	if resultSets, ok := result["resultSets"].([]interface{}); ok {
		fmt.Printf("📦 Result sets found: %d\n", len(resultSets))
		for i, rs := range resultSets {
			if rsMap, ok := rs.(map[string]interface{}); ok {
				name := rsMap["name"]
				var rowCount int
				if rows, ok := rsMap["rowSet"].([]interface{}); ok {
					rowCount = len(rows)
				}
				fmt.Printf("   %d. %s (%d rows)\n", i+1, name, rowCount)
			}
		}
		fmt.Println()
	}

	// Save to file if requested
	if *output != "" {
		if err := saveJSON(*output, result); err != nil {
			fmt.Printf("❌ Error saving to file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("💾 Output saved to: %s\n\n", *output)
	} else {
		// Print to stdout
		fmt.Println("📄 Full Response:")
		fmt.Println("==================\n")
		printJSON(result)
	}

	fmt.Println("\n✨ Done!")
	fmt.Println("\n💡 Tips:")
	fmt.Println("   - Use --output=file.json to save response")
	fmt.Println("   - Use --date=2026-02-25 to specify date")
	fmt.Println("   - Use --game=0022300789 for game-specific data")
	fmt.Println("   - Use --season=2025-26 for standings")
	fmt.Println("\n📚 Examples:")
	fmt.Println("   go run scripts/test_nba_api.go --endpoint=scoreboard --date=2026-02-25")
	fmt.Println("   go run scripts/test_nba_api.go --endpoint=summary --game=0022300789 --output=summary.json")
	fmt.Println("   go run scripts/test_nba_api.go --endpoint=traditional --game=0022300789")
	fmt.Println("   go run scripts/test_nba_api.go --endpoint=standings --season=2025-26")
}

// currentSeason returns the NBA season string for today's date.
// e.g. if today is Feb 2026 → "2025-26"
func currentSeason() string {
	now := time.Now()
	year := now.Year()
	// NBA season starts in October; before October we're still in previous season
	if now.Month() < 10 {
		year--
	}
	return fmt.Sprintf("%d-%02d", year, (year+1)%100)
}
