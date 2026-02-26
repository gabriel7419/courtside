package reddit

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// NBA r/nba highlight post titles come in many formats, e.g.:
//   - "Jayson Tatum 31 PTS - Full Game Highlights | Celtics vs Heat"
//   - "LeBron James [32/8/7] - Game Highlights | Lakers vs Warriors"
//   - "Boston Celtics vs Miami Heat - Full Game Highlights | Feb 25, 2026"
//   - "Celtics [87] - Heat [79] 3rd Quarter Highlights"
//
// Unlike soccer, NBA highlights don't encode a "minute" — they encode quarter
// and total points. The matcher therefore focuses on:
//   1. Team names in the title  (primary signal)
//   2. Scorer name              (high value for scoring-event highlights)
//   3. Score proximity          (exact score match is a good tiebreaker)
//   4. Date proximity           (post must be within ±24h of the game)

// findBestMatch finds the best matching Reddit search result for a basketball event.
func findBestMatch(results []SearchResult, goal GoalInfo) *SearchResult {
	if len(results) == 0 {
		return nil
	}

	homeNorm := normalizeTeamName(goal.HomeTeam)
	awayNorm := normalizeTeamName(goal.AwayTeam)

	var bestMatch *SearchResult
	bestScore := 0

	for i := range results {
		result := &results[i]
		titleLower := strings.ToLower(result.Title)

		score := 0

		// ── 1. Date filter ─────────────────────────────────────────────
		if !goal.MatchTime.IsZero() {
			postDate := result.CreatedAt
			matchStart := goal.MatchTime.Add(-24 * time.Hour)
			matchEnd := goal.MatchTime.Add(48 * time.Hour)

			if postDate.Before(matchStart) || postDate.After(matchEnd) {
				continue // outside valid window
			}
			// Bonus for posts very close to game time
			if postDate.After(goal.MatchTime.Add(-6*time.Hour)) && postDate.Before(goal.MatchTime.Add(12*time.Hour)) {
				score += 5
			}
		}

		// ── 2. Team names ──────────────────────────────────────────────
		homeFound := containsTeamName(titleLower, homeNorm)
		awayFound := containsTeamName(titleLower, awayNorm)

		if !homeFound && !awayFound {
			continue // at least one team required
		}
		if homeFound {
			score += 10
		}
		if awayFound {
			score += 10
		}
		if homeFound && awayFound {
			score += 5 // bonus for both teams present
		}

		// ── 3. Scorer / player name ────────────────────────────────────
		// For NBA, the scorer name (e.g., "Tatum", "LeBron") is often the
		// most prominent signal in the title.
		if goal.ScorerName != "" {
			scorerNorm := normalizeName(goal.ScorerName)
			if containsName(titleLower, scorerNorm) {
				score += 20 // High bonus — player name in title is very specific
			}
		}

		// ── 4. Score proximity ─────────────────────────────────────────
		// NBA highlight titles sometimes embed the final score like
		// "Celtics [87] - Heat [79]" or "BOS 112, PHI 104".
		// We check for both the home and away score as integers.
		if scoreFoundInTitle(titleLower, goal.HomeScore, goal.AwayScore) {
			score += 15
		}

		// ── 5. Keyword bonuses ─────────────────────────────────────────
		if strings.Contains(titleLower, "highlight") {
			score += 3
		}
		if strings.Contains(titleLower, "game winner") || strings.Contains(titleLower, "buzzer") {
			score += 5
		}

		// ── 6. Upvote tiebreaker ───────────────────────────────────────
		score += min(result.Score/200, 5) // max 5 pts from upvotes

		if score > bestScore {
			bestScore = score
			bestMatch = result
		}
	}

	// NBA highlight matching is softer than soccer goals — we only need
	// at least one team name plus either a player or score match.
	const minScore = 20
	if bestScore < minScore {
		return nil
	}
	return bestMatch
}

// scoreFoundInTitle checks if both the home and away score appear in the title
// as standalone integers (e.g., "87" and "79" in "Celtics [87] - Heat [79]").
func scoreFoundInTitle(title string, home, away int) bool {
	if home <= 0 && away <= 0 {
		return false
	}
	homePattern := fmt.Sprintf(`\b%d\b`, home)
	awayPattern := fmt.Sprintf(`\b%d\b`, away)
	homeRe := regexp.MustCompile(homePattern)
	awayRe := regexp.MustCompile(awayPattern)
	return homeRe.MatchString(title) && awayRe.MatchString(title)
}

// normalizeTeamName converts an NBA team name to a normalized form for matching.
// Handles abbreviations ("BOS", "LAL"), city names, and nicknames.
func normalizeTeamName(name string) string {
	norm := strings.ToLower(name)

	// Remove generic prefixes/suffixes not useful for matching
	prefixes := []string{"the "}
	for _, p := range prefixes {
		norm = strings.TrimPrefix(norm, p)
	}

	// Keep only alphanumeric + spaces
	norm = regexp.MustCompile(`[^a-z0-9\s]`).ReplaceAllString(norm, "")
	return strings.TrimSpace(norm)
}

// normalizeName converts a player name to a normalized form for matching.
func normalizeName(name string) string {
	norm := strings.ToLower(name)
	norm = regexp.MustCompile(`[^a-z\s]`).ReplaceAllString(norm, "")
	return strings.TrimSpace(norm)
}

// containsTeamName checks if a title contains a team name or any significant word from it.
func containsTeamName(title, teamNorm string) bool {
	titleNorm := normalizeTeamName(title)

	if strings.Contains(titleNorm, teamNorm) {
		return true
	}

	// Match individual words (e.g., "celtics", "warriors", "lakers")
	words := strings.Fields(teamNorm)
	for _, word := range words {
		if len(word) > 3 && strings.Contains(titleNorm, word) {
			return true
		}
	}

	return false
}

// containsName checks if a title contains a player name (full name or last name).
func containsName(title, nameNorm string) bool {
	if strings.Contains(title, nameNorm) {
		return true
	}
	// Try matching just the last name (most recognizable for NBA stars)
	parts := strings.Fields(nameNorm)
	if len(parts) > 0 {
		lastName := parts[len(parts)-1]
		if len(lastName) > 2 && strings.Contains(title, lastName) {
			return true
		}
	}
	return false
}

// MatchConfidence represents how confident we are in a match.
type MatchConfidence int

const (
	ConfidenceNone   MatchConfidence = 0
	ConfidenceLow    MatchConfidence = 1
	ConfidenceMedium MatchConfidence = 2
	ConfidenceHigh   MatchConfidence = 3
)

// CalculateConfidence returns the confidence level for a Reddit result match.
func CalculateConfidence(result SearchResult, goal GoalInfo) MatchConfidence {
	titleLower := strings.ToLower(result.Title)
	homeNorm := normalizeTeamName(goal.HomeTeam)
	awayNorm := normalizeTeamName(goal.AwayTeam)

	hasHome := containsTeamName(titleLower, homeNorm)
	hasAway := containsTeamName(titleLower, awayNorm)

	hasPlayer := false
	if goal.ScorerName != "" {
		hasPlayer = containsName(titleLower, normalizeName(goal.ScorerName))
	}

	hasScore := scoreFoundInTitle(titleLower, goal.HomeScore, goal.AwayScore)

	if hasHome && hasAway && (hasPlayer || hasScore) {
		return ConfidenceHigh
	}
	if (hasHome || hasAway) && (hasPlayer || hasScore) {
		return ConfidenceMedium
	}
	if hasHome || hasAway {
		return ConfidenceLow
	}
	return ConfidenceNone
}
