// Package nba provides a client for the NBA Stats API.
// The API is publicly accessible but requires specific headers on every request.
// See docs/API_REFERENCE.md for endpoint documentation.
package nba

import (
	"fmt"
)

// The NBA Stats API returns data in a tabular format: each resultSet has
// a list of column names (headers) and a list of rows (rowSet).
// Use the col() helper to extract values by column name.

// resultSet is the common tabular format used by all NBA Stats API responses.
type resultSet struct {
	Name    string          `json:"name"`
	Headers []string        `json:"headers"`
	RowSet  [][]interface{} `json:"rowSet"`
}

// col returns the value of the given column for the given row, or nil.
func (rs resultSet) col(row []interface{}, field string) interface{} {
	for i, h := range rs.Headers {
		if h == field && i < len(row) {
			return row[i]
		}
	}
	return nil
}

// colStr returns the string value of a column, or "".
func (rs resultSet) colStr(row []interface{}, field string) string {
	v := rs.col(row, field)
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}

// colInt returns the int value of a column, or 0.
func (rs resultSet) colInt(row []interface{}, field string) int {
	v := rs.col(row, field)
	if v == nil {
		return 0
	}
	switch val := v.(type) {
	case float64:
		return int(val)
	case int:
		return val
	}
	return 0
}

// colIntPtr returns a *int, or nil if the value is null.
func (rs resultSet) colIntPtr(row []interface{}, field string) *int {
	v := rs.col(row, field)
	if v == nil {
		return nil
	}
	switch val := v.(type) {
	case float64:
		n := int(val)
		return &n
	case int:
		return &val
	}
	return nil
}

// colFloat returns the float64 value of a column, or 0.
func (rs resultSet) colFloat(row []interface{}, field string) float64 {
	v := rs.col(row, field)
	if v == nil {
		return 0
	}
	if f, ok := v.(float64); ok {
		return f
	}
	return 0
}

// findResultSet finds the named resultSet in a slice, or returns an empty one.
func findResultSet(sets []resultSet, name string) resultSet {
	for _, rs := range sets {
		if rs.Name == name {
			return rs
		}
	}
	return resultSet{}
}

// --- Response types ---

// scoreboardV3Response is returned by GET /stats/scoreboardv3
// NOTE: v3 returns a nested JSON structure, NOT the resultSets pattern.
// Schema: scoreboard > { gameHeader[], lineScore[], gameLeaders[] }
type scoreboardV3Response struct {
	Scoreboard scoreboardV3Scoreboard `json:"scoreboard"`
}

type scoreboardV3Scoreboard struct {
	GameDate   string             `json:"gameDate"`
	LeagueID   string             `json:"leagueId"`
	LeagueName string             `json:"leagueName"`
	Games      []scoreboardV3Game `json:"games"`
}

type scoreboardV3Game struct {
	GameID           string           `json:"gameId"`
	GameCode         string           `json:"gameCode"`
	GameStatus       int              `json:"gameStatus"`     // 1=scheduled, 2=live, 3=final
	GameStatusText   string           `json:"gameStatusText"` // "Final", "Q3 2:34"
	Period           int              `json:"period"`         // current quarter
	GameClock        string           `json:"gameClock"`      // "PT02M34.00S" or ""
	GameTimeUTC      string           `json:"gameTimeUTC"`    // ISO8601
	SeriesGameNumber string           `json:"seriesGameNumber,omitempty"`
	SeriesText       string           `json:"seriesText,omitempty"` // "Celtics lead 2-1"
	HomeTeam         scoreboardV3Team `json:"homeTeam"`
	AwayTeam         scoreboardV3Team `json:"awayTeam"`
}

type scoreboardV3Team struct {
	TeamID            int                  `json:"teamId"`
	TeamCity          string               `json:"teamCity"`
	TeamName          string               `json:"teamName"`
	TeamTricode       string               `json:"teamTricode"`
	TeamSlug          string               `json:"teamSlug"`
	Wins              int                  `json:"wins"`
	Losses            int                  `json:"losses"`
	Score             int                  `json:"score"`
	Seed              int                  `json:"seed,omitempty"`
	InBonus           bool                 `json:"inBonus,omitempty"`
	TimeoutsRemaining int                  `json:"timeoutsRemaining,omitempty"`
	Periods           []scoreboardV3Period `json:"periods,omitempty"`
}

type scoreboardV3Period struct {
	Period     int    `json:"period"`
	PeriodType string `json:"periodType,omitempty"`
	Score      int    `json:"score"`
}

// boxScoreSummaryResponse is returned by GET /stats/boxscoresummaryv2 (still v2, not deprecated)
type boxScoreSummaryResponse struct {
	ResultSets []resultSet `json:"resultSets"`
}

// boxScoreTraditionalV3Response is returned by GET /stats/boxscoretraditionalv3
// v3 uses camelCase fields inside resultSets (same envelope, different column names).
type boxScoreTraditionalV3Response struct {
	ResultSets []resultSet `json:"resultSets"`
}

// playByPlayV3Response is returned by GET /stats/playbyplayv3
// v3 uses camelCase fields; returns empty JSON in v2 since 2024-25.
type playByPlayV3Response struct {
	ResultSets []resultSet `json:"resultSets"`
}

// standingsResponse is returned by GET /stats/leaguestandingsv3
type standingsResponse struct {
	ResultSets []resultSet `json:"resultSets"`
}

// scoreboardResponse kept for backward-compat while we still parse scoreboardv2 in test script
// Remove once scoreboardv3 migration is complete.
type scoreboardResponse = scoreboardV3Response

// --- Status constants ---

// GAME_STATUS_ID values in the NBA API.
const (
	gameStatusScheduled = 1
	gameStatusLive      = 2
	gameStatusFinal     = 3
)

// EVENTMSGTYPE values in play-by-play.
const (
	eventFieldGoalMade   = 1
	eventFieldGoalMissed = 2
	eventFreeThrowMade   = 3
	eventFreeThrowMissed = 4
	eventRebound         = 5
	eventFoul            = 6
	eventViolation       = 7
	eventSubstitution    = 8
	eventTimeout         = 9
	eventJumpBall        = 10
	eventEjection        = 11
	eventStartPeriod     = 12
	eventEndPeriod       = 13
)
