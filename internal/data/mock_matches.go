package data

import (
	"encoding/json"
	"os"

	"github.com/0xjuanma/golazo/internal/api"
)

// MockMatchesData contains multiple matches for testing.
type MockMatchesData struct {
	Matches []MockMatchData `json:"matches"`
}

type MockMatchData struct {
	ID    int    `json:"id"`
	Round string `json:"round"`
	Home  struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		ShortName string `json:"shortName"`
	} `json:"home"`
	Away struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		ShortName string `json:"shortName"`
	} `json:"away"`
	Status struct {
		UTCTime   string `json:"utcTime"`
		Started   bool   `json:"started"`
		Finished  bool   `json:"finished"`
		Cancelled bool   `json:"cancelled"`
		LiveTime  *struct {
			Short string `json:"short"`
		} `json:"liveTime,omitempty"`
		Score *struct {
			Home int `json:"home"`
			Away int `json:"away"`
		} `json:"score,omitempty"`
	} `json:"status"`
	League struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Country     string `json:"country"`
		CountryCode string `json:"countryCode"`
	} `json:"league"`
}

// MockMatches returns multiple mock matches for testing.
func MockMatches() ([]api.Match, error) {
	// Try to load from config directory first
	configPath, err := MockDataPath()
	if err == nil {
		if data, err := os.ReadFile(configPath); err == nil {
			var mockData MockMatchesData
			if err := json.Unmarshal(data, &mockData); err == nil {
				return convertMockMatches(mockData.Matches), nil
			}
		}
	}

	// Fallback to embedded mock data
	return getDefaultMockMatches(), nil
}

// getDefaultMockMatches returns default mock matches.
func getDefaultMockMatches() []api.Match {
	matches := []api.Match{
		{
			ID: 1,
			League: api.League{
				ID:          1,
				Name:        "Premier League",
				Country:     "England",
				CountryCode: "GB",
			},
			HomeTeam: api.Team{
				ID:        1,
				Name:      "Manchester United",
				ShortName: "Man Utd",
			},
			AwayTeam: api.Team{
				ID:        2,
				Name:      "Liverpool",
				ShortName: "Liverpool",
			},
			Status:    api.MatchStatusLive,
			HomeScore: intPtr(2),
			AwayScore: intPtr(1),
			LiveTime:  stringPtr("67'"),
			Round:     "Matchday 15",
		},
		{
			ID: 2,
			League: api.League{
				ID:          2,
				Name:        "La Liga",
				Country:     "Spain",
				CountryCode: "ES",
			},
			HomeTeam: api.Team{
				ID:        3,
				Name:      "Real Madrid",
				ShortName: "Real Madrid",
			},
			AwayTeam: api.Team{
				ID:        4,
				Name:      "Barcelona",
				ShortName: "Barcelona",
			},
			Status:    api.MatchStatusLive,
			HomeScore: intPtr(1),
			AwayScore: intPtr(1),
			LiveTime:  stringPtr("23'"),
			Round:     "Matchday 12",
		},
		{
			ID: 3,
			League: api.League{
				ID:          3,
				Name:        "Serie A",
				Country:     "Italy",
				CountryCode: "IT",
			},
			HomeTeam: api.Team{
				ID:        5,
				Name:      "AC Milan",
				ShortName: "AC Milan",
			},
			AwayTeam: api.Team{
				ID:        6,
				Name:      "Inter Milan",
				ShortName: "Inter",
			},
			Status:    api.MatchStatusFinished,
			HomeScore: intPtr(3),
			AwayScore: intPtr(2),
			LiveTime:  stringPtr("FT"),
			Round:     "Matchday 10",
		},
		{
			ID: 4,
			League: api.League{
				ID:          1,
				Name:        "Premier League",
				Country:     "England",
				CountryCode: "GB",
			},
			HomeTeam: api.Team{
				ID:        7,
				Name:      "Arsenal",
				ShortName: "Arsenal",
			},
			AwayTeam: api.Team{
				ID:        8,
				Name:      "Chelsea",
				ShortName: "Chelsea",
			},
			Status: api.MatchStatusNotStarted,
			Round:  "Matchday 15",
		},
	}

	// Save to config directory for persistence
	if configPath, err := MockDataPath(); err == nil {
		mockData := MockMatchesData{Matches: convertToMockData(matches)}
		if data, err := json.MarshalIndent(mockData, "", "  "); err == nil {
			os.WriteFile(configPath, data, 0644)
		}
	}

	return matches
}

func convertMockMatches(mockMatches []MockMatchData) []api.Match {
	matches := make([]api.Match, 0, len(mockMatches))
	for _, m := range mockMatches {
		match := api.Match{
			ID: m.ID,
			League: api.League{
				ID:          m.League.ID,
				Name:        m.League.Name,
				Country:     m.League.Country,
				CountryCode: m.League.CountryCode,
			},
			HomeTeam: api.Team{
				ID:        m.Home.ID,
				Name:      m.Home.Name,
				ShortName: m.Home.ShortName,
			},
			AwayTeam: api.Team{
				ID:        m.Away.ID,
				Name:      m.Away.Name,
				ShortName: m.Away.ShortName,
			},
			Round: m.Round,
		}

		// Determine status
		if m.Status.Cancelled {
			match.Status = api.MatchStatusCancelled
		} else if m.Status.Finished {
			match.Status = api.MatchStatusFinished
			if m.Status.LiveTime != nil {
				match.LiveTime = &m.Status.LiveTime.Short
			}
		} else if m.Status.Started {
			match.Status = api.MatchStatusLive
			if m.Status.LiveTime != nil {
				match.LiveTime = &m.Status.LiveTime.Short
			}
		} else {
			match.Status = api.MatchStatusNotStarted
		}

		// Set scores if available
		if m.Status.Score != nil {
			match.HomeScore = &m.Status.Score.Home
			match.AwayScore = &m.Status.Score.Away
		}

		matches = append(matches, match)
	}
	return matches
}

func convertToMockData(matches []api.Match) []MockMatchData {
	result := make([]MockMatchData, 0, len(matches))
	for _, m := range matches {
		mm := MockMatchData{
			ID:    m.ID,
			Round: m.Round,
		}
		mm.Home.ID = m.HomeTeam.ID
		mm.Home.Name = m.HomeTeam.Name
		mm.Home.ShortName = m.HomeTeam.ShortName
		mm.Away.ID = m.AwayTeam.ID
		mm.Away.Name = m.AwayTeam.Name
		mm.Away.ShortName = m.AwayTeam.ShortName
		mm.League.ID = m.League.ID
		mm.League.Name = m.League.Name
		mm.League.Country = m.League.Country
		mm.League.CountryCode = m.League.CountryCode

		mm.Status.Started = m.Status == api.MatchStatusLive
		mm.Status.Finished = m.Status == api.MatchStatusFinished
		mm.Status.Cancelled = m.Status == api.MatchStatusCancelled

		if m.LiveTime != nil {
			mm.Status.LiveTime = &struct {
				Short string `json:"short"`
			}{Short: *m.LiveTime}
		}

		if m.HomeScore != nil && m.AwayScore != nil {
			mm.Status.Score = &struct {
				Home int `json:"home"`
				Away int `json:"away"`
			}{
				Home: *m.HomeScore,
				Away: *m.AwayScore,
			}
		}

		result = append(result, mm)
	}
	return result
}

func intPtr(i int) *int {
	return &i
}
