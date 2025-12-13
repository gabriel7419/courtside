package data

import (
	"fmt"
	"math/rand"
	"time"
)

// LiveUpdateGenerator generates mock live updates for a match.
type LiveUpdateGenerator struct {
	matchID int
	updates []string
	index   int
}

// NewLiveUpdateGenerator creates a new live update generator for a match.
func NewLiveUpdateGenerator(matchID int) *LiveUpdateGenerator {
	updates := getMockLiveUpdates(matchID)
	return &LiveUpdateGenerator{
		matchID: matchID,
		updates: updates,
		index:   0,
	}
}

// GetNextUpdate returns the next live update string.
func (g *LiveUpdateGenerator) GetNextUpdate() (string, bool) {
	if g.index >= len(g.updates) {
		return "", false
	}
	update := g.updates[g.index]
	g.index++
	return update, true
}

// HasMore returns true if there are more updates.
func (g *LiveUpdateGenerator) HasMore() bool {
	return g.index < len(g.updates)
}

// getMockLiveUpdates returns mock live update strings for a match.
func getMockLiveUpdates(matchID int) []string {
	switch matchID {
	case 1: // Man Utd vs Liverpool
		return []string{
			"68' Corner kick for Man Utd",
			"69' Shot on target by Fernandes, saved by goalkeeper",
			"70' Substitution: Antony → Sancho (Man Utd)",
			"71' Free kick for Liverpool",
			"72' Yellow card: Van Dijk (Liverpool)",
			"73' Corner kick for Liverpool",
			"74' Shot wide by Nunez",
		}
	case 2: // Real Madrid vs Barcelona
		return []string{
			"24' Corner kick for Real Madrid",
			"25' Shot blocked by Barcelona defense",
			"26' Free kick for Barcelona",
			"27' Yellow card: Modric (Real Madrid)",
			"28' Substitution: Gavi → Pedri (Barcelona)",
			"29' Corner kick for Barcelona",
		}
	case 3: // AC Milan vs Inter (finished)
		return []string{
			"Match finished: AC Milan 3-2 Inter",
		}
	case 4: // Arsenal vs Chelsea (not started)
		return []string{
			"Match not started yet",
		}
	default:
		return []string{
			fmt.Sprintf("Live update for match %d", matchID),
		}
	}
}

// GenerateRandomUpdate generates a random live update string.
func GenerateRandomUpdate(matchID int) string {
	updates := []string{
		"Corner kick",
		"Free kick",
		"Shot on target",
		"Shot wide",
		"Yellow card",
		"Substitution",
		"Offside",
		"Foul",
		"Throw-in",
		"Goal kick",
	}

	rand.Seed(time.Now().UnixNano())
	return updates[rand.Intn(len(updates))]
}
