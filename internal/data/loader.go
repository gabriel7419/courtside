// Package data provides utilities for loading mock football match data.
package data

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/0xjuanma/golazo/internal/api"
)

// MockMatches loads matches from the mock data file.
// It searches for the mock data file in common locations relative to the working directory.
func MockMatches() ([]api.Match, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("get working directory: %w", err)
	}

	paths := []string{
		filepath.Join(wd, "assets", "mock_data.json"),
		filepath.Join(wd, "internal", "fotmob", "mock_data.json"),
	}

	var data []byte
	var readErr error
	for _, path := range paths {
		data, readErr = os.ReadFile(path)
		if readErr == nil {
			break
		}
	}

	if readErr != nil {
		return nil, fmt.Errorf("read mock data file: %w", readErr)
	}

	var mockData struct {
		Matches []struct {
			ID      int    `json:"id"`
			Round   string `json:"round"`
			Home    struct {
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
				UTCTime  string `json:"utcTime"`
				Started  bool   `json:"started"`
				Finished bool   `json:"finished"`
				Cancelled bool  `json:"cancelled"`
				LiveTime *struct {
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
		} `json:"matches"`
	}

	if err := json.Unmarshal(data, &mockData); err != nil {
		return nil, fmt.Errorf("unmarshal mock data: %w", err)
	}

	matches := make([]api.Match, 0, len(mockData.Matches))
	for _, m := range mockData.Matches {
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

	return matches, nil
}

