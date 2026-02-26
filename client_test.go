package main

import (
	"context"
	"testing"
	"time"

	"github.com/gabriel7419/courtside/internal/nba"
)

func TestMatchDetailsFinished(t *testing.T) {
	client := nba.NewClient()
	ctx := context.Background()

	date := time.Now().UTC().AddDate(0, 0, -1)
	t.Logf("Fetching MatchesByDate for %s...", date.Format("2006-01-02"))

	matches, err := client.MatchesByDate(ctx, date)
	if err != nil {
		t.Fatalf("ERROR MatchesByDate: %v", err)
	}

	if len(matches) == 0 {
		t.Fatalf("No matches found for yesterday.")
	}

	m := matches[0]
	t.Logf("Found match: %s vs %s (ID: %d). Fetching MatchDetails...", m.HomeTeam.ShortName, m.AwayTeam.ShortName, m.ID)

	details, err := client.MatchDetails(ctx, m.ID, &m)
	if err != nil {
		t.Fatalf("ERROR MatchDetails: %v", err)
	}

	t.Logf("SUCCESS: Scores: %v - %v, Status: %v", details.HomeScore, details.AwayScore, details.Status)
}
