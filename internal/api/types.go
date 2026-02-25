package api

import "time"

// League represents a sports league or conference.
type League struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Country        string `json:"country"`
	CountryCode    string `json:"country_code"`
	Logo           string `json:"logo,omitempty"`
	ParentLeagueID int    `json:"parent_league_id,omitempty"`
}

// Team represents a team.
type Team struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"short_name"`
	Logo      string `json:"logo,omitempty"`
}

// MatchStatus represents the status of a match or game.
type MatchStatus string

const (
	MatchStatusNotStarted MatchStatus = "not_started"
	MatchStatusLive       MatchStatus = "live"
	MatchStatusFinished   MatchStatus = "finished"
	MatchStatusPostponed  MatchStatus = "postponed"
	MatchStatusCancelled  MatchStatus = "cancelled"
)

// Match represents a match or game.
type Match struct {
	ID        int         `json:"id"`
	League    League      `json:"league"`
	HomeTeam  Team        `json:"home_team"`
	AwayTeam  Team        `json:"away_team"`
	Status    MatchStatus `json:"status"`
	HomeScore *int        `json:"home_score,omitempty"`
	AwayScore *int        `json:"away_score,omitempty"`
	MatchTime *time.Time  `json:"match_time,omitempty"`
	LiveTime  *string     `json:"live_time,omitempty"` // football: "45+2", "HT", "FT"; NBA: "Q3 2:34"
	Round     string      `json:"round,omitempty"`
	PageURL   string      `json:"page_url,omitempty"`

	// NBA-specific fields
	Quarter      *int    `json:"quarter,omitempty"` // 1-4, 5+ = OT
	Clock        *string `json:"clock,omitempty"`   // "2:34"
	IsPlayoffs   bool    `json:"is_playoffs,omitempty"`
	SeriesStatus *string `json:"series_status,omitempty"` // "Series tied 2-2"
}

// MatchEvent represents an event during a match (goal, card, field goal, foul, etc.).
type MatchEvent struct {
	ID            int       `json:"id"`
	Minute        int       `json:"minute"`
	DisplayMinute string    `json:"display_minute,omitempty"` // "45+2'" or "Q3 2:34"
	Type          string    `json:"type"`                     // "goal", "card", "substitution", "field_goal", "foul", "timeout"
	Team          Team      `json:"team"`
	Player        *string   `json:"player,omitempty"`
	Assist        *string   `json:"assist,omitempty"`
	EventType     *string   `json:"event_type,omitempty"` // "yellow", "red", "personal", "technical", "flagrant"
	OwnGoal       *bool     `json:"own_goal,omitempty"`
	Timestamp     time.Time `json:"timestamp"`

	// NBA-specific fields
	Points       *int    `json:"points,omitempty"`        // 1 (free throw), 2, or 3 (field goal)
	IsThree      *bool   `json:"is_three,omitempty"`      // whether it was a 3-pointer
	EventSubtype *string `json:"event_subtype,omitempty"` // "personal", "technical", "flagrant"
}

// MatchStatistic represents a single statistic entry (possession, FG%, rebounds, etc.).
type MatchStatistic struct {
	Key       string `json:"key"`   // e.g., "possession", "fg_pct", "rebounds"
	Label     string `json:"label"` // e.g., "Possession", "FG%", "Rebounds"
	HomeValue string `json:"home_value"`
	AwayValue string `json:"away_value"`
}

// PlayerInfo represents basic player information for lineups/rosters.
type PlayerInfo struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Number   int    `json:"number,omitempty"`
	Position string `json:"position,omitempty"`
	Rating   string `json:"rating,omitempty"`
}

// MatchDetails contains detailed information about a match or game.
type MatchDetails struct {
	Match
	Events     []MatchEvent `json:"events"`
	HomeLineup []string     `json:"home_lineup,omitempty"`
	AwayLineup []string     `json:"away_lineup,omitempty"`

	// Score at half time / end of regulation
	HalfTimeScore *struct {
		Home *int `json:"home,omitempty"`
		Away *int `json:"away,omitempty"`
	} `json:"half_time_score,omitempty"`
	Venue         string  `json:"venue,omitempty"`
	Winner        *string `json:"winner,omitempty"` // "home" or "away"
	MatchDuration int     `json:"match_duration,omitempty"`
	ExtraTime     bool    `json:"extra_time,omitempty"`
	Penalties     *struct {
		Home *int `json:"home,omitempty"`
		Away *int `json:"away,omitempty"`
	} `json:"penalties,omitempty"`

	// Statistics
	Statistics []MatchStatistic `json:"statistics,omitempty"`

	// Context
	Referee    string `json:"referee,omitempty"`
	Attendance int    `json:"attendance,omitempty"`

	// Lineups (football)
	HomeFormation   string       `json:"home_formation,omitempty"`
	AwayFormation   string       `json:"away_formation,omitempty"`
	HomeStarting    []PlayerInfo `json:"home_starting,omitempty"`
	AwayStarting    []PlayerInfo `json:"away_starting,omitempty"`
	HomeSubstitutes []PlayerInfo `json:"home_substitutes,omitempty"`
	AwaySubstitutes []PlayerInfo `json:"away_substitutes,omitempty"`

	// xG (football)
	HomeXG *float64 `json:"home_xg,omitempty"`
	AwayXG *float64 `json:"away_xg,omitempty"`

	// Highlights
	Highlight *MatchHighlight `json:"highlight,omitempty"`

	// NBA-specific fields
	// QuarterScores stores scores per period as alternating pairs: [Q1home, Q1away, Q2home, Q2away, ...]
	QuarterScores []int `json:"quarter_scores,omitempty"`
	Overtime      bool  `json:"overtime,omitempty"`
}

// MatchHighlight represents a highlight video link.
type MatchHighlight struct {
	URL    string `json:"url"`
	Image  string `json:"image,omitempty"`
	Source string `json:"source,omitempty"`
	Title  string `json:"title,omitempty"`
}

// LeagueTableEntry represents a team's standing in a league or conference.
type LeagueTableEntry struct {
	Position       int    `json:"position"`
	Team           Team   `json:"team"`
	Played         int    `json:"played"`
	Won            int    `json:"won"`
	Drawn          int    `json:"drawn,omitempty"` // football only
	Lost           int    `json:"lost"`
	GoalsFor       int    `json:"goals_for,omitempty"`     // football only
	GoalsAgainst   int    `json:"goals_against,omitempty"` // football only
	PointsFor      int    `json:"points_for,omitempty"`    // NBA: win% × 1000
	PointsAgainst  int    `json:"points_against,omitempty"`
	GoalDifference int    `json:"goal_difference,omitempty"`
	Points         int    `json:"points"`         // football: pts; NBA: wins
	Form           string `json:"form,omitempty"` // e.g. "W3", "L2" for NBA; "WWDLL" for football
	Note           string `json:"note,omitempty"` // e.g. "East | GB: 3.5"
}
