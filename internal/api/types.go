package api

import "time"

// League represents a football league
type League struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
	Logo        string `json:"logo,omitempty"`
}

// Team represents a football team
type Team struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"short_name"`
	Logo      string `json:"logo,omitempty"`
}

// MatchStatus represents the status of a match
type MatchStatus string

const (
	MatchStatusNotStarted MatchStatus = "not_started"
	MatchStatusLive       MatchStatus = "live"
	MatchStatusFinished   MatchStatus = "finished"
	MatchStatusPostponed  MatchStatus = "postponed"
	MatchStatusCancelled  MatchStatus = "cancelled"
)

// Match represents a football match
type Match struct {
	ID        int         `json:"id"`
	League    League      `json:"league"`
	HomeTeam  Team        `json:"home_team"`
	AwayTeam  Team        `json:"away_team"`
	Status    MatchStatus `json:"status"`
	HomeScore *int        `json:"home_score,omitempty"`
	AwayScore *int        `json:"away_score,omitempty"`
	MatchTime *time.Time  `json:"match_time,omitempty"`
	LiveTime  *string     `json:"live_time,omitempty"` // e.g., "45+2", "HT", "FT"
	Round     string      `json:"round,omitempty"`
}

// MatchEvent represents an event in a match (goal, card, substitution, etc.)
type MatchEvent struct {
	ID        int       `json:"id"`
	Minute    int       `json:"minute"`
	Type      string    `json:"type"` // "goal", "card", "substitution", etc.
	Team      Team      `json:"team"`
	Player    *string   `json:"player,omitempty"`
	Assist    *string   `json:"assist,omitempty"`
	EventType *string   `json:"event_type,omitempty"` // "yellow", "red", "in", "out", etc.
	Timestamp time.Time `json:"timestamp"`
}

// MatchDetails contains detailed information about a match
type MatchDetails struct {
	Match
	Events     []MatchEvent `json:"events"`
	HomeLineup []string     `json:"home_lineup,omitempty"`
	AwayLineup []string     `json:"away_lineup,omitempty"`

	// Additional match information
	HalfTimeScore *struct {
		Home *int `json:"home,omitempty"`
		Away *int `json:"away,omitempty"`
	} `json:"half_time_score,omitempty"`
	Venue         string  `json:"venue,omitempty"`          // Stadium name
	Winner        *string `json:"winner,omitempty"`         // "home" or "away"
	MatchDuration int     `json:"match_duration,omitempty"` // 90, 120, etc.
	ExtraTime     bool    `json:"extra_time,omitempty"`     // If match went to extra time
	Penalties     *struct {
		Home *int `json:"home,omitempty"`
		Away *int `json:"away,omitempty"`
	} `json:"penalties,omitempty"`
}

// LeagueTableEntry represents a team's position in the league table
type LeagueTableEntry struct {
	Position       int  `json:"position"`
	Team           Team `json:"team"`
	Played         int  `json:"played"`
	Won            int  `json:"won"`
	Drawn          int  `json:"drawn"`
	Lost           int  `json:"lost"`
	GoalsFor       int  `json:"goals_for"`
	GoalsAgainst   int  `json:"goals_against"`
	GoalDifference int  `json:"goal_difference"`
	Points         int  `json:"points"`
}
