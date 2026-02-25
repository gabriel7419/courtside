// Script de teste para explorar a NBA Stats API
// Execute: go run scripts/test_nba_api.go

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
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
			Timeout: 15 * time.Second,
		},
	}
}

func (c *NBAClient) makeRequest(url string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// Headers obrigatórios para NBA API
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Referer", "https://www.nba.com/")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Origin", "https://www.nba.com")

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
	url := fmt.Sprintf("%s/scoreboard?GameDate=%s&LeagueID=00", baseURL, date)
	return c.makeRequest(url)
}

func (c *NBAClient) GetBoxScoreSummary(gameID string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/boxscoresummaryv2?GameID=%s", baseURL, gameID)
	return c.makeRequest(url)
}

func (c *NBAClient) GetBoxScoreTraditional(gameID string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/boxscoretraditionalv2?GameID=%s&StartPeriod=1&EndPeriod=10&StartRange=0&EndRange=28800&RangeType=0", baseURL, gameID)
	return c.makeRequest(url)
}

func (c *NBAClient) GetPlayByPlay(gameID string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/playbyplayv2?GameID=%s&StartPeriod=1&EndPeriod=10", baseURL, gameID)
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
	endpoint := flag.String("endpoint", "scoreboard", "Endpoint to test: scoreboard, summary, traditional, playbyplay")
	date := flag.String("date", time.Now().Format("2006-01-02"), "Date for scoreboard (YYYY-MM-DD)")
	gameID := flag.String("game", "", "Game ID for detailed endpoints")
	output := flag.String("output", "", "Save output to file")
	flag.Parse()

	client := NewNBAClient()

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
		fmt.Println("Valid endpoints: scoreboard, summary, traditional, playbyplay")
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
	fmt.Println("\n📚 Examples:")
	fmt.Println("   go run scripts/test_nba_api.go --endpoint=scoreboard --date=2026-02-25")
	fmt.Println("   go run scripts/test_nba_api.go --endpoint=summary --game=0022300789 --output=summary.json")
	fmt.Println("   go run scripts/test_nba_api.go --endpoint=traditional --game=0022300789")
}
