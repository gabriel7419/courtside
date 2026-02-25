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

// scoreboardResponse is returned by GET /stats/scoreboard
type scoreboardResponse struct {
	ResultSets []resultSet `json:"resultSets"`
}

// boxScoreSummaryResponse is returned by GET /stats/boxscoresummaryv2
type boxScoreSummaryResponse struct {
	ResultSets []resultSet `json:"resultSets"`
}

// boxScoreTraditionalResponse is returned by GET /stats/boxscoretraditionalv2
type boxScoreTraditionalResponse struct {
	ResultSets []resultSet `json:"resultSets"`
}

// playByPlayResponse is returned by GET /stats/playbyplayv2
type playByPlayResponse struct {
	ResultSets []resultSet `json:"resultSets"`
}

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
