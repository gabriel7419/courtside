package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/0xjuanma/golazo/internal/data"
)

func main() {
	// Get API key
	apiKey, err := data.FootballDataAPIKey()
	if err != nil {
		fmt.Printf("❌ Error getting API key: %v\n", err)
		fmt.Println("\nMake sure FOOTBALL_DATA_API_KEY is set:")
		fmt.Println("  export FOOTBALL_DATA_API_KEY=\"your-key-here\"")
		os.Exit(1)
	}

	fmt.Printf("✓ API key found (length: %d)\n\n", len(apiKey))

	// Test the API endpoint directly
	baseURL := "https://v3.football.api-sports.io"

	// Test 1: Get all fixtures for today (no filters)
	fmt.Println("Test 1: Fetching all fixtures for today...")
	today := time.Now().Format("2006-01-02")
	url1 := fmt.Sprintf("%s/fixtures?date=%s", baseURL, today)
	testEndpoint(url1, apiKey)

	// Test 2: Get finished matches for today (all leagues)
	fmt.Println("\nTest 2: Fetching finished matches for today (all leagues)...")
	url2 := fmt.Sprintf("%s/fixtures?date=%s&status=FT", baseURL, today)
	testEndpoint(url2, apiKey)

	// Test 3: Get finished matches for today - Premier League only
	fmt.Println("\nTest 3: Fetching finished Premier League matches for today...")
	url3 := fmt.Sprintf("%s/fixtures?date=%s&status=FT&league=39", baseURL, today)
	testEndpoint(url3, apiKey)

	// Test 4: Get finished matches for today - La Liga only
	fmt.Println("\nTest 4: Fetching finished La Liga matches for today...")
	url4 := fmt.Sprintf("%s/fixtures?date=%s&status=FT&league=140", baseURL, today)
	testEndpoint(url4, apiKey)

	// Test 5: Get finished matches for today - Bundesliga only
	fmt.Println("\nTest 5: Fetching finished Bundesliga matches for today...")
	url5 := fmt.Sprintf("%s/fixtures?date=%s&status=FT&league=78", baseURL, today)
	testEndpoint(url5, apiKey)

	// Test 6: Test all supported leagues in one go (simulating actual app behavior)
	fmt.Println("\nTest 6: Testing all supported leagues for today (simulating app behavior)...")
	supportedLeagues := []int{39, 140, 78, 135, 61} // Premier League, La Liga, Bundesliga, Serie A, Ligue 1
	totalMatches := 0
	for _, leagueID := range supportedLeagues {
		url := fmt.Sprintf("%s/fixtures?date=%s&status=FT&league=%d", baseURL, today, leagueID)
		fmt.Printf("  League %d: ", leagueID)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("x-apisports-key", apiKey)
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			var result map[string]interface{}
			bodyBytes, _ := io.ReadAll(resp.Body)
			json.Unmarshal(bodyBytes, &result)
			resp.Body.Close()
			if response, ok := result["response"].([]interface{}); ok {
				fmt.Printf("%d matches\n", len(response))
				totalMatches += len(response)
			} else {
				fmt.Printf("0 matches\n")
			}
		} else {
			fmt.Printf("error\n")
		}
	}
	fmt.Printf("  Total finished matches across all supported leagues: %d\n", totalMatches)

	// Test 7: Test date range (1 day, 3 days, 7 days) - simulating actual app behavior
	fmt.Println("\nTest 7: Testing date ranges (1d, 3d, 7d) - simulating app behavior...")
	for _, days := range []int{1, 3, 7} {
		fmt.Printf("\n  Testing %d day(s):\n", days)
		today := time.Now()
		dateFrom := today.AddDate(0, 0, -(days - 1))
		fmt.Printf("    Date range: %s to %s\n", dateFrom.Format("2006-01-02"), today.Format("2006-01-02"))

		totalMatches := 0
		currentDate := dateFrom
		for !currentDate.After(today) {
			dateStr := currentDate.Format("2006-01-02")
			fmt.Printf("    Date %s:\n", dateStr)

			for _, leagueID := range supportedLeagues {
				url := fmt.Sprintf("%s/fixtures?date=%s&league=%d&status=FT", baseURL, dateStr, leagueID)
				req, _ := http.NewRequest("GET", url, nil)
				req.Header.Set("x-apisports-key", apiKey)
				client := &http.Client{Timeout: 10 * time.Second}
				resp, err := client.Do(req)
				if err == nil && resp.StatusCode == http.StatusOK {
					var result map[string]interface{}
					bodyBytes, _ := io.ReadAll(resp.Body)
					json.Unmarshal(bodyBytes, &result)
					resp.Body.Close()
					if response, ok := result["response"].([]interface{}); ok {
						if len(response) > 0 {
							fmt.Printf("      League %d: %d matches\n", leagueID, len(response))
							totalMatches += len(response)
						}
					}
				}
			}

			currentDate = currentDate.AddDate(0, 0, 1)
		}
		fmt.Printf("    Total matches for %d day(s): %d\n", days, totalMatches)
	}

	// Test 8: Try without status filter (maybe matches aren't marked as FT yet)
	fmt.Println("\nTest 8: Testing without status filter (all statuses)...")
	for _, leagueID := range supportedLeagues {
		url := fmt.Sprintf("%s/fixtures?date=%s&league=%d", baseURL, today, leagueID)
		fmt.Printf("  League %d: ", leagueID)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("x-apisports-key", apiKey)
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			var result map[string]interface{}
			bodyBytes, _ := io.ReadAll(resp.Body)
			json.Unmarshal(bodyBytes, &result)
			resp.Body.Close()
			if response, ok := result["response"].([]interface{}); ok {
				fmt.Printf("%d matches (all statuses)\n", len(response))
				if len(response) > 0 {
					// Show status breakdown
					statusCount := make(map[string]int)
					for _, match := range response {
						if m, ok := match.(map[string]interface{}); ok {
							if fixture, ok := m["fixture"].(map[string]interface{}); ok {
								if status, ok := fixture["status"].(map[string]interface{}); ok {
									if short, ok := status["short"].(string); ok {
										statusCount[short]++
									}
								}
							}
						}
					}
					fmt.Printf("    Status breakdown: %v\n", statusCount)
				}
			}
		}
	}

	// Test 9: Check yesterday and day before (matches might finish late)
	fmt.Println("\nTest 9: Checking yesterday and day before for finished matches...")
	for i := 1; i <= 2; i++ {
		checkDate := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		fmt.Printf("\n  Date: %s\n", checkDate)
		totalMatches := 0
		for _, leagueID := range supportedLeagues {
			url := fmt.Sprintf("%s/fixtures?date=%s&league=%d&status=FT", baseURL, checkDate, leagueID)
			req, _ := http.NewRequest("GET", url, nil)
			req.Header.Set("x-apisports-key", apiKey)
			client := &http.Client{Timeout: 10 * time.Second}
			resp, err := client.Do(req)
			if err == nil && resp.StatusCode == http.StatusOK {
				var result map[string]interface{}
				bodyBytes, _ := io.ReadAll(resp.Body)
				json.Unmarshal(bodyBytes, &result)
				resp.Body.Close()
				if response, ok := result["response"].([]interface{}); ok {
					if len(response) > 0 {
						fmt.Printf("    League %d: %d finished matches\n", leagueID, len(response))
						totalMatches += len(response)
					}
				}
			}
		}
		fmt.Printf("    Total finished matches: %d\n", totalMatches)
	}

	// Test 10: Try using season parameter (some APIs require season)
	fmt.Println("\nTest 10: Testing with season parameter (2024)...")
	currentYear := time.Now().Year()
	for _, leagueID := range supportedLeagues {
		url := fmt.Sprintf("%s/fixtures?date=%s&league=%d&season=%d&status=FT", baseURL, today, leagueID, currentYear)
		fmt.Printf("  League %d: ", leagueID)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("x-apisports-key", apiKey)
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			var result map[string]interface{}
			bodyBytes, _ := io.ReadAll(resp.Body)
			json.Unmarshal(bodyBytes, &result)
			resp.Body.Close()
			if response, ok := result["response"].([]interface{}); ok {
				fmt.Printf("%d matches\n", len(response))
			}
		}
	}

	// Test 11: Analyze what leagues are actually in the 669 matches
	fmt.Println("\nTest 11: Analyzing leagues in today's finished matches...")
	url := fmt.Sprintf("%s/fixtures?date=%s&status=FT", baseURL, today)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("x-apisports-key", apiKey)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err == nil && resp.StatusCode == http.StatusOK {
		var result map[string]interface{}
		bodyBytes, _ := io.ReadAll(resp.Body)
		json.Unmarshal(bodyBytes, &result)
		resp.Body.Close()
		if response, ok := result["response"].([]interface{}); ok {
			leagueCount := make(map[int]int)        // league ID -> count
			leagueNames := make(map[int]string)     // league ID -> name
			leagueCountries := make(map[int]string) // league ID -> country

			for _, match := range response {
				if m, ok := match.(map[string]interface{}); ok {
					if league, ok := m["league"].(map[string]interface{}); ok {
						if leagueID, ok := league["id"].(float64); ok {
							id := int(leagueID)
							leagueCount[id]++

							if name, ok := league["name"].(string); ok {
								leagueNames[id] = name
							}
							if country, ok := league["country"].(string); ok {
								leagueCountries[id] = country
							}
						}
					}
				}
			}

			fmt.Printf("  Found %d unique leagues\n", len(leagueCount))
			fmt.Println("\n  All leagues with their IDs and match counts:")

			// Sort by count
			type leagueInfo struct {
				id      int
				name    string
				country string
				count   int
			}
			leagues := make([]leagueInfo, 0, len(leagueCount))
			for id, count := range leagueCount {
				leagues = append(leagues, leagueInfo{
					id:      id,
					name:    leagueNames[id],
					country: leagueCountries[id],
					count:   count,
				})
			}
			// Sort by count (descending)
			for i := 0; i < len(leagues); i++ {
				maxIdx := i
				for j := i + 1; j < len(leagues); j++ {
					if leagues[j].count > leagues[maxIdx].count {
						maxIdx = j
					}
				}
				leagues[i], leagues[maxIdx] = leagues[maxIdx], leagues[i]
			}

			// Show all leagues (or top 50 if too many)
			maxShow := len(leagues)
			if maxShow > 50 {
				maxShow = 50
			}
			for i := 0; i < maxShow; i++ {
				fmt.Printf("    ID: %3d | %-40s | %-20s | %d matches\n",
					leagues[i].id,
					leagues[i].name,
					leagues[i].country,
					leagues[i].count)
			}
			if len(leagues) > 50 {
				fmt.Printf("    ... and %d more leagues\n", len(leagues)-50)
			}

			// Check if our supported league IDs exist in the response
			fmt.Println("\n  Checking our supported league IDs:")
			supportedLeagueIDs := []int{39, 140, 78, 135, 61}
			supportedLeagueNames := map[int]string{
				39:  "Premier League",
				140: "La Liga",
				78:  "Bundesliga",
				135: "Serie A",
				61:  "Ligue 1",
			}

			for _, id := range supportedLeagueIDs {
				if count, exists := leagueCount[id]; exists {
					name := leagueNames[id]
					country := leagueCountries[id]
					fmt.Printf("    ✓ ID %3d (%s, %s): %d matches found\n", id, name, country, count)
				} else {
					fmt.Printf("    ✗ ID %3d (%s): NOT FOUND in today's matches\n", id, supportedLeagueNames[id])

					// Try to find leagues with similar names
					fmt.Printf("      Searching for similar league names...\n")
					searchTerms := []string{
						"Premier League", "Premier", "EPL",
						"La Liga", "Liga",
						"Bundesliga",
						"Serie A", "Serie",
						"Ligue 1", "Ligue",
					}
					for _, term := range searchTerms {
						for foundID, foundName := range leagueNames {
							if foundName != "" &&
								(foundName == term ||
									(len(term) > 5 && len(foundName) > 5 &&
										(foundName[:5] == term[:5] || foundName[len(foundName)-5:] == term[len(term)-5:]))) {
								fmt.Printf("      → Found: ID %d - %s (%s) - %d matches\n",
									foundID, foundName, leagueCountries[foundID], leagueCount[foundID])
							}
						}
					}
				}
			}

			// Also check by name matching
			fmt.Println("\n  Checking by league name (case-insensitive):")
			for id, expectedName := range supportedLeagueNames {
				found := false
				for foundID, foundName := range leagueNames {
					if foundName == expectedName {
						fmt.Printf("    ✓ Found '%s' with ID %d (we use %d): %d matches\n",
							expectedName, foundID, id, leagueCount[foundID])
						if foundID != id {
							fmt.Printf("      ⚠ ID MISMATCH! API uses %d, we use %d\n", foundID, id)
						}
						found = true
						break
					}
				}
				if !found {
					fmt.Printf("    ✗ '%s' not found by exact name match\n", expectedName)
				}
			}
		}
	}
}

func testEndpoint(url string, apiKey string) {
	fmt.Printf("  URL: %s\n", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("  ❌ Failed to create request: %v\n", err)
		return
	}

	req.Header.Set("x-apisports-key", apiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("  ❌ Request failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyStr := string(bodyBytes)

	fmt.Printf("  Status: %d\n", resp.StatusCode)
	fmt.Printf("  Headers: %+v\n", resp.Header)

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("  ❌ Error response: %s\n", bodyStr[:min(1000, len(bodyStr))])
		return
	}

	// Try to parse JSON
	var result map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		fmt.Printf("  ❌ Failed to parse JSON: %v\n", err)
		fmt.Printf("  Response: %s\n", bodyStr[:min(1000, len(bodyStr))])
		return
	}

	// Print full response structure for debugging
	fmt.Printf("  Response keys: %v\n", getKeys(result))

	// Check for errors in response
	if errors, ok := result["errors"].([]interface{}); ok && len(errors) > 0 {
		fmt.Printf("  ⚠ API Errors: %+v\n", errors)
	}

	// Check response structure
	if response, ok := result["response"].([]interface{}); ok {
		fmt.Printf("  ✓ Found %d matches\n", len(response))
		if len(response) > 0 {
			// Show first match summary
			if firstMatch, ok := response[0].(map[string]interface{}); ok {
				fmt.Printf("  Sample match ID: %v\n", firstMatch["fixture"])
				if teams, ok := firstMatch["teams"].(map[string]interface{}); ok {
					if home, ok := teams["home"].(map[string]interface{}); ok {
						if away, ok := teams["away"].(map[string]interface{}); ok {
							fmt.Printf("  Sample: %v vs %v\n", home["name"], away["name"])
						}
					}
				}
			}
		} else {
			fmt.Printf("  ⚠ No matches found in response\n")
		}
	} else {
		fmt.Printf("  ⚠ Unexpected response structure\n")
		fmt.Printf("  Full response: %s\n", bodyStr[:min(500, len(bodyStr))])
	}
}

func getKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
