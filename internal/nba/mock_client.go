package nba

import (
	"context"
	"fmt"
	"time"

	"github.com/gabriel7419/courtside/internal/api"
	"github.com/gabriel7419/courtside/internal/data"
)

// MockClient implements api.Client using hard-coded NBA fixture data.
// Use it when the NBA Stats API is unavailable (network issues, development, CI).
// Switch between clients by passing -mock flag or setting NBA_MOCK=1 env var.
type MockClient struct{}

// NewMockClient creates a mock NBA client that returns fixture data without
// making any network requests.
func NewMockClient() *MockClient {
	return &MockClient{}
}

// MatchesByDate returns mock NBA games for the given date.
// All games use today's date, so any date returns the same fixture set.
func (c *MockClient) MatchesByDate(_ context.Context, _ time.Time) ([]api.Match, error) {
	matches := append(data.MockNBALiveMatches(), data.MockNBAUpcomingMatches()...)
	return matches, nil
}

// MatchDetails returns mock game details for the given matchID.
func (c *MockClient) MatchDetails(_ context.Context, matchID int) (*api.MatchDetails, error) {
	details, err := data.MockNBAMatchDetails(matchID)
	if err != nil {
		return nil, err
	}
	if details == nil {
		return nil, fmt.Errorf("mock: no details for match %d", matchID)
	}
	return details, nil
}

// MatchDetailsForceRefresh is identical to MatchDetails for the mock client
// (there is no cache to bypass).
func (c *MockClient) MatchDetailsForceRefresh(ctx context.Context, matchID int) (*api.MatchDetails, error) {
	return c.MatchDetails(ctx, matchID)
}

// LiveMatches returns mock live games.
func (c *MockClient) LiveMatches(_ context.Context) ([]api.Match, error) {
	return data.MockNBALiveMatches(), nil
}

// LiveMatchesForceRefresh is identical to LiveMatches for the mock client.
func (c *MockClient) LiveMatchesForceRefresh(ctx context.Context) ([]api.Match, error) {
	return c.LiveMatches(ctx)
}

// Leagues returns an empty slice (NBA uses conferences, not leagues).
func (c *MockClient) Leagues(_ context.Context) ([]api.League, error) {
	return []api.League{}, nil
}

// LeagueMatches returns an empty slice.
func (c *MockClient) LeagueMatches(_ context.Context, _ int) ([]api.Match, error) {
	return []api.Match{}, nil
}

// LeagueTable returns mock NBA conference standings.
// leagueID 0 = both, 1 = East, 2 = West.
func (c *MockClient) LeagueTable(_ context.Context, leagueID int, _ string) ([]api.LeagueTableEntry, error) {
	return mockNBAStandings(leagueID), nil
}

// Cache returns nil; the mock client has no cache.
// This satisfies any caller that does a nil check before using the cache.
func (c *MockClient) Cache() *ResponseCache {
	return nil
}

// mockNBAStandings returns a condensed set of standings fixture data.
func mockNBAStandings(leagueID int) []api.LeagueTableEntry {
	type row struct {
		id   int
		name string
		abbr string
		conf string
		pos  int
		w    int
		l    int
		pct  float64
		gb   string
		strk string
	}

	rows := []row{
		// East
		{1610612738, "Boston Celtics", "BOS", "East", 1, 48, 11, .814, "—", "W4"},
		{1610612739, "Cleveland Cavaliers", "CLE", "East", 2, 44, 16, .733, "4.5", "W2"},
		{1610612741, "Chicago Bulls", "CHI", "East", 3, 35, 25, .583, "13.5", "L1"},
		{1610612752, "New York Knicks", "NYK", "East", 4, 33, 27, .550, "15.5", "W1"},
		{1610612748, "Miami Heat", "MIA", "East", 5, 30, 30, .500, "18.5", "L2"},
		{1610612749, "Milwaukee Bucks", "MIL", "East", 6, 29, 31, .483, "19.5", "W1"},
		{1610612755, "Philadelphia 76ers", "PHI", "East", 7, 28, 32, .467, "20.5", "L3"},
		// West
		{1610612760, "Oklahoma City Thunder", "OKC", "West", 1, 46, 13, .780, "—", "W5"},
		{1610612743, "Denver Nuggets", "DEN", "West", 2, 40, 20, .667, "6.5", "L1"},
		{1610612750, "Minnesota Timberwolves", "MIN", "West", 3, 36, 24, .600, "10.5", "W2"},
		{1610612744, "Golden State Warriors", "GSW", "West", 4, 32, 28, .533, "14.5", "W1"},
		{1610612747, "Los Angeles Lakers", "LAL", "West", 5, 29, 31, .483, "17.5", "L2"},
		{1610612742, "Dallas Mavericks", "DAL", "West", 6, 27, 33, .450, "19.5", "L1"},
		{1610612756, "Phoenix Suns", "PHX", "West", 7, 25, 35, .417, "21.5", "W1"},
	}

	var entries []api.LeagueTableEntry
	for _, r := range rows {
		if leagueID == 1 && r.conf != "East" {
			continue
		}
		if leagueID == 2 && r.conf != "West" {
			continue
		}
		entries = append(entries, api.LeagueTableEntry{
			Team:      api.Team{ID: r.id, Name: r.name, ShortName: r.abbr},
			Position:  r.pos,
			Played:    r.w + r.l,
			Won:       r.w,
			Lost:      r.l,
			Points:    r.w,
			PointsFor: int(r.pct * 1000),
			Form:      r.strk,
			Note:      fmt.Sprintf("%s | GB: %s", r.conf, r.gb),
		})
	}
	return entries
}
