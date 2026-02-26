package nba

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gabriel7419/courtside/internal/api"
)

const (
	baseURL = "https://stats.nba.com/stats"
)

// Client implements api.Client for the NBA Stats API.
type Client struct {
	httpClient  *http.Client
	baseURL     string
	rateLimiter *RateLimiter
	cache       *ResponseCache
}

// Cache returns the response cache for external access (e.g., pre-populating live matches).
func (c *Client) Cache() *ResponseCache {
	return c.cache
}

// NewClient creates a new NBA API client with default configuration.
func NewClient() *Client {
	return &Client{
		httpClient:  &http.Client{Timeout: 15 * time.Second},
		baseURL:     baseURL,
		rateLimiter: NewRateLimiter(250 * time.Millisecond),
		cache:       NewResponseCache(DefaultCacheConfig()),
	}
}

// do executes a GET request to the NBA Stats API with the required headers.
func (c *Client) do(ctx context.Context, url string, dst interface{}) error {
	c.rateLimiter.Wait()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	// These headers are required — without them the API returns 403.
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Referer", "https://www.nba.com/")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Origin", "https://www.nba.com")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("execute request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d for %s", resp.StatusCode, url)
	}

	if err := json.NewDecoder(resp.Body).Decode(dst); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}

// MatchesByDate retrieves all NBA games scheduled for the given date.
func (c *Client) MatchesByDate(ctx context.Context, date time.Time) ([]api.Match, error) {
	dateStr := date.UTC().Format("2006-01-02")

	if cached := c.cache.Matches(dateStr); cached != nil {
		return cached, nil
	}

	url := fmt.Sprintf("%s/scoreboard?GameDate=%s&LeagueID=00", c.baseURL, dateStr)

	var resp scoreboardResponse
	if err := c.do(ctx, url, &resp); err != nil {
		return nil, fmt.Errorf("fetch scoreboard for %s: %w", dateStr, err)
	}

	gameHeader := findResultSet(resp.ResultSets, "GameHeader")
	lineScore := findResultSet(resp.ResultSets, "LineScore")

	// Build a map of gameID → line scores (two rows per game: home + away)
	type lineScoreRow struct {
		teamID   int
		teamAbbr string
		teamName string
		pts      *int
		ptsQ1    *int
		ptsQ2    *int
		ptsQ3    *int
		ptsQ4    *int
		ptsOT1   *int
	}
	lineScores := make(map[string][]lineScoreRow)
	for _, row := range lineScore.RowSet {
		gameID := lineScore.colStr(row, "GAME_ID")
		ls := lineScoreRow{
			teamID:   lineScore.colInt(row, "TEAM_ID"),
			teamAbbr: lineScore.colStr(row, "TEAM_ABBREVIATION"),
			teamName: lineScore.colStr(row, "TEAM_CITY_NAME") + " " + lineScore.colStr(row, "TEAM_NICKNAME"),
			pts:      lineScore.colIntPtr(row, "PTS"),
			ptsQ1:    lineScore.colIntPtr(row, "PTS_QTR1"),
			ptsQ2:    lineScore.colIntPtr(row, "PTS_QTR2"),
			ptsQ3:    lineScore.colIntPtr(row, "PTS_QTR3"),
			ptsQ4:    lineScore.colIntPtr(row, "PTS_QTR4"),
			ptsOT1:   lineScore.colIntPtr(row, "PTS_OT1"),
		}
		lineScores[gameID] = append(lineScores[gameID], ls)
	}

	var matches []api.Match
	for _, row := range gameHeader.RowSet {
		gameID := gameHeader.colStr(row, "GAME_ID")
		statusID := gameHeader.colInt(row, "GAME_STATUS_ID")
		statusText := gameHeader.colStr(row, "GAME_STATUS_TEXT")
		homeTeamID := gameHeader.colInt(row, "HOME_TEAM_ID")
		livePeriod := gameHeader.colInt(row, "LIVE_PERIOD")
		liveBcast := gameHeader.colStr(row, "LIVE_PERIOD_TIME_BCAST") // "Q3 2:34"

		// Parse game time from GAME_DATE_EST
		var matchTime *time.Time
		if dateVal := gameHeader.colStr(row, "GAME_DATE_EST"); dateVal != "" {
			if t, err := time.Parse("2006-01-02T15:04:05", dateVal); err == nil {
				matchTime = &t
			}
		}

		// Map API status to internal status
		var status api.MatchStatus
		switch statusID {
		case gameStatusScheduled:
			status = api.MatchStatusNotStarted
		case gameStatusLive:
			status = api.MatchStatusLive
		case gameStatusFinal:
			status = api.MatchStatusFinished
		default:
			status = api.MatchStatusNotStarted
		}

		// Build team info from line scores
		var homeTeam, awayTeam api.Team
		var homeScore, awayScore *int
		rows := lineScores[gameID]
		for _, ls := range rows {
			team := api.Team{
				ID:        ls.teamID,
				Name:      strings.TrimSpace(ls.teamName),
				ShortName: ls.teamAbbr,
			}
			if ls.teamID == homeTeamID {
				homeTeam = team
				homeScore = ls.pts
			} else {
				awayTeam = team
				awayScore = ls.pts
			}
		}

		// Parse game ID to int (last 10 chars are the numeric ID)
		// We use a simple hash since game IDs are strings like "0022300789"
		numericID := gameHeader.colInt(row, "GAME_SEQUENCE")
		if numericID == 0 {
			numericID = simpleHash(gameID)
		}

		// Period + clock
		var quarter *int
		var clock *string
		if livePeriod > 0 {
			quarter = &livePeriod
		}
		if liveBcast != "" && status == api.MatchStatusLive {
			clock = &liveBcast
		}

		// LiveTime: human-readable game time
		var liveTime *string
		if liveBcast != "" {
			liveTime = &liveBcast
		} else if status == api.MatchStatusFinished {
			ft := "Final"
			liveTime = &ft
		}

		// Conference from GAMECODE (e.g., "20260225/LALGWS")
		league := api.League{
			Name: "NBA",
		}

		isOT := livePeriod > 4

		m := api.Match{
			ID:           numericID,
			League:       league,
			HomeTeam:     homeTeam,
			AwayTeam:     awayTeam,
			Status:       status,
			HomeScore:    homeScore,
			AwayScore:    awayScore,
			MatchTime:    matchTime,
			LiveTime:     liveTime,
			Quarter:      quarter,
			Clock:        clock,
			IsPlayoffs:   isPlayoffGame(gameID),
			SeriesStatus: nil,
		}
		_ = isOT // will be used in MatchDetails
		_ = statusText

		matches = append(matches, m)
	}

	c.cache.SetMatches(dateStr, matches)
	return matches, nil
}

// MatchDetails retrieves detailed information about a specific game.
// gameID is the numeric ID stored in api.Match.ID.
func (c *Client) MatchDetails(ctx context.Context, matchID int) (*api.MatchDetails, error) {
	if cached := c.cache.Details(matchID); cached != nil {
		return cached, nil
	}

	// We need the string game ID — store it when fetching MatchesByDate.
	// For now, reconstruct from matchID. In a full implementation the
	// string game ID would be stored similarly to fotmob's pageURLs.
	gameIDStr := storedGameID(matchID)
	if gameIDStr == "" {
		return nil, fmt.Errorf("game ID not found for match %d; call MatchesByDate first", matchID)
	}

	url := fmt.Sprintf("%s/boxscoresummaryv2?GameID=%s", c.baseURL, gameIDStr)
	var summaryResp boxScoreSummaryResponse
	if err := c.do(ctx, url, &summaryResp); err != nil {
		return nil, fmt.Errorf("fetch box score summary for game %s: %w", gameIDStr, err)
	}

	details := parseSummary(summaryResp, matchID)

	// Fetch play-by-play for live games
	if details.Status == api.MatchStatusLive {
		pbpURL := fmt.Sprintf("%s/playbyplayv2?GameID=%s&StartPeriod=1&EndPeriod=10", c.baseURL, gameIDStr)
		var pbpResp playByPlayResponse
		if err := c.do(ctx, pbpURL, &pbpResp); err == nil {
			details.Events = parsePlayByPlay(pbpResp)
		}
	}

	// Fetch team box score stats (FG%, REB, AST, etc.) for all game states
	statsURL := fmt.Sprintf("%s/boxscoretraditionalv2?GameID=%s&StartPeriod=1&EndPeriod=10&StartRange=0&EndRange=28800&RangeType=0", c.baseURL, gameIDStr)
	var statsResp boxScoreTraditionalResponse
	if err := c.do(ctx, statsURL, &statsResp); err == nil {
		details.Statistics = parseTeamStats(statsResp, details.HomeTeam.ID)
	}

	c.cache.SetDetails(matchID, details)
	return details, nil
}

// MatchDetailsForceRefresh bypasses the cache and fetches fresh game data.
func (c *Client) MatchDetailsForceRefresh(ctx context.Context, matchID int) (*api.MatchDetails, error) {
	c.cache.ClearDetails(matchID)
	return c.MatchDetails(ctx, matchID)
}

// LiveMatches returns currently live NBA games.
func (c *Client) LiveMatches(ctx context.Context) ([]api.Match, error) {
	if cached := c.cache.LiveMatches(); cached != nil {
		return cached, nil
	}
	today := time.Now()
	all, err := c.MatchesByDate(ctx, today)
	if err != nil {
		return nil, err
	}
	var live []api.Match
	for _, m := range all {
		if m.Status == api.MatchStatusLive {
			live = append(live, m)
		}
	}
	c.cache.SetLiveMatches(live)
	return live, nil
}

// LiveMatchesForceRefresh fetches live games bypassing the cache.
func (c *Client) LiveMatchesForceRefresh(ctx context.Context) ([]api.Match, error) {
	c.cache.ClearLive()
	return c.LiveMatches(ctx)
}

// Leagues returns an empty slice (NBA uses conferences, not leagues).
func (c *Client) Leagues(_ context.Context) ([]api.League, error) {
	return []api.League{}, nil
}

// LeagueMatches returns an empty slice (use MatchesByDate instead).
func (c *Client) LeagueMatches(_ context.Context, _ int) ([]api.Match, error) {
	return []api.Match{}, nil
}

// LeagueTable returns NBA standings for the requested conference.
// leagueID: 0 = all teams, 1 = Eastern Conference, 2 = Western Conference.
// leagueName is ignored for NBA (kept for interface compatibility).
func (c *Client) LeagueTable(ctx context.Context, leagueID int, _ string) ([]api.LeagueTableEntry, error) {
	season := currentNBASeason()
	url := fmt.Sprintf("%s/leaguestandingsv3?LeagueID=00&Season=%s&SeasonType=Regular+Season", c.baseURL, season)

	var resp standingsResponse
	if err := c.do(ctx, url, &resp); err != nil {
		return nil, fmt.Errorf("fetch standings for %s: %w", season, err)
	}

	standings := findResultSet(resp.ResultSets, "Standings")

	var entries []api.LeagueTableEntry
	for _, row := range standings.RowSet {
		conf := standings.colStr(row, "Conference")

		// Filter by conference if requested
		switch leagueID {
		case 1: // Eastern
			if conf != "East" {
				continue
			}
		case 2: // Western
			if conf != "West" {
				continue
			}
		}

		teamName := strings.TrimSpace(standings.colStr(row, "TeamCity") + " " + standings.colStr(row, "TeamName"))
		wins := standings.colInt(row, "WINS")
		losses := standings.colInt(row, "LOSSES")

		// Win percentage
		total := wins + losses
		var pct float64
		if total > 0 {
			pct = float64(wins) / float64(total)
		}

		// Games behind (float stored as string like "3.5")
		gbStr := standings.colStr(row, "ConferenceGamesBack")
		if gbStr == "" {
			gbStr = standings.colStr(row, "LeagueGamesBack")
		}

		// Current streak ("W3" or "L2")
		streakStr := standings.colStr(row, "CurrentStreak")

		// Points for/against — use as proxy for stats
		// NBA doesn't expose PF/PA in standings v3, so we leave them 0
		// and store meaningful data in the label fields
		confRank := standings.colInt(row, "ConferenceRank")

		entry := api.LeagueTableEntry{
			Team:           api.Team{Name: teamName, ShortName: standings.colStr(row, "TeamAbbreviation")},
			Position:       confRank,
			Played:         total,
			Won:            wins,
			Lost:           losses,
			Points:         wins,            // NBA uses wins as points
			PointsFor:      int(pct * 1000), // win% × 1000 as a sortable int
			PointsAgainst:  0,
			GoalDifference: 0,
			Form:           streakStr,
			Note:           fmt.Sprintf("%s | GB: %s", conf, gbStr),
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

// currentNBASeason returns the NBA season string for the current date.
// e.g. Feb 2026 → "2025-26"
func currentNBASeason() string {
	now := time.Now()
	year := now.Year()
	if now.Month() < 10 {
		year-- // season starts in October
	}
	short := (year + 1) % 100
	return fmt.Sprintf("%d-%02d", year, short)
}

// --- Parsing helpers ---

// parseSummary converts a boxScoreSummaryResponse to api.MatchDetails.
func parseSummary(resp boxScoreSummaryResponse, matchID int) *api.MatchDetails {
	gs := findResultSet(resp.ResultSets, "GameSummary")
	ls := findResultSet(resp.ResultSets, "LineScore")
	gi := findResultSet(resp.ResultSets, "GameInfo")

	details := &api.MatchDetails{}
	details.ID = matchID

	if len(gs.RowSet) > 0 {
		row := gs.RowSet[0]
		statusID := gs.colInt(row, "GAME_STATUS_ID")
		livePeriod := gs.colInt(row, "LIVE_PERIOD")

		switch statusID {
		case gameStatusScheduled:
			details.Status = api.MatchStatusNotStarted
		case gameStatusLive:
			details.Status = api.MatchStatusLive
		case gameStatusFinal:
			details.Status = api.MatchStatusFinished
		}

		homeTeamID := gs.colInt(row, "HOME_TEAM_ID")
		details.League = api.League{Name: "NBA"}

		// Build teams from LineScore
		for _, lsRow := range ls.RowSet {
			teamID := ls.colInt(lsRow, "TEAM_ID")
			team := api.Team{
				ID:        teamID,
				Name:      strings.TrimSpace(ls.colStr(lsRow, "TEAM_CITY_NAME") + " " + ls.colStr(lsRow, "TEAM_NICKNAME")),
				ShortName: ls.colStr(lsRow, "TEAM_ABBREVIATION"),
			}
			score := ls.colIntPtr(lsRow, "PTS")

			// Quarter scores
			q1 := ls.colIntPtr(lsRow, "PTS_QTR1")
			q2 := ls.colIntPtr(lsRow, "PTS_QTR2")
			q3 := ls.colIntPtr(lsRow, "PTS_QTR3")
			q4 := ls.colIntPtr(lsRow, "PTS_QTR4")

			if teamID == homeTeamID {
				details.HomeTeam = team
				details.HomeScore = score
				appendQuarterScores(&details.QuarterScores, 0, q1, q2, q3, q4)
			} else {
				details.AwayTeam = team
				details.AwayScore = score
				appendQuarterScores(&details.QuarterScores, 1, q1, q2, q3, q4)
			}
		}

		if livePeriod > 0 {
			details.Quarter = &livePeriod
		}
		details.Overtime = livePeriod > 4
		if details.Overtime {
			details.ExtraTime = true
		}
	}

	// Attendance from GameInfo
	if len(gi.RowSet) > 0 {
		details.Attendance = gi.colInt(gi.RowSet[0], "ATTENDANCE")
	}

	return details
}

// parseTeamStats converts a boxScoreTraditionalResponse to []api.MatchStatistic.
// It reads the TeamStats resultSet which has one row per team, and produces
// home/away labelled comparison stats for the UI statistics dialog.
func parseTeamStats(resp boxScoreTraditionalResponse, homeTeamID int) []api.MatchStatistic {
	ts := findResultSet(resp.ResultSets, "TeamStats")
	if len(ts.RowSet) < 2 {
		return nil // need at least 2 rows (home + away)
	}

	// Identify home vs away rows
	var homeRow, awayRow []interface{}
	for _, row := range ts.RowSet {
		teamID := ts.colInt(row, "TEAM_ID")
		if teamID == homeTeamID {
			homeRow = row
		} else {
			awayRow = row
		}
	}
	if homeRow == nil || awayRow == nil {
		return nil
	}

	// formatPct renders a decimal like 0.512 as "51.2%"
	formatPct := func(row []interface{}, field string) string {
		v := ts.colFloat(row, field)
		if v == 0 {
			return "—"
		}
		return fmt.Sprintf("%.1f%%", v*100)
	}

	// formatInt renders an integer column
	formatInt := func(row []interface{}, field string) string {
		v := ts.colInt(row, field)
		return fmt.Sprintf("%d", v)
	}

	stat := func(key, label, homeVal, awayVal string) api.MatchStatistic {
		return api.MatchStatistic{Key: key, Label: label, HomeValue: homeVal, AwayValue: awayVal}
	}

	return []api.MatchStatistic{
		stat("fg_pct", "FG%", formatPct(homeRow, "FG_PCT"), formatPct(awayRow, "FG_PCT")),
		stat("fg3_pct", "3P%", formatPct(homeRow, "FG3_PCT"), formatPct(awayRow, "FG3_PCT")),
		stat("ft_pct", "FT%", formatPct(homeRow, "FT_PCT"), formatPct(awayRow, "FT_PCT")),
		stat("reb", "Rebounds", formatInt(homeRow, "REB"), formatInt(awayRow, "REB")),
		stat("oreb", "Off. Rebounds", formatInt(homeRow, "OREB"), formatInt(awayRow, "OREB")),
		stat("dreb", "Def. Rebounds", formatInt(homeRow, "DREB"), formatInt(awayRow, "DREB")),
		stat("ast", "Assists", formatInt(homeRow, "AST"), formatInt(awayRow, "AST")),
		stat("stl", "Steals", formatInt(homeRow, "STL"), formatInt(awayRow, "STL")),
		stat("blk", "Blocks", formatInt(homeRow, "BLK"), formatInt(awayRow, "BLK")),
		stat("tov", "Turnovers", formatInt(homeRow, "TO"), formatInt(awayRow, "TO")),
		stat("pf", "Personal Fouls", formatInt(homeRow, "PF"), formatInt(awayRow, "PF")),
		stat("fga", "FG Attempted", formatInt(homeRow, "FGA"), formatInt(awayRow, "FGA")),
		stat("fg3a", "3P Attempted", formatInt(homeRow, "FG3A"), formatInt(awayRow, "FG3A")),
		stat("fta", "FT Attempted", formatInt(homeRow, "FTA"), formatInt(awayRow, "FTA")),
	}
}

// appendQuarterScores grows the slice to store Q1..Q4 for the given team slot (0=home, 1=away).
// Layout: index 0=Q1home, 1=Q1away, 2=Q2home, 3=Q2away, ...
func appendQuarterScores(scores *[]int, slot int, q1, q2, q3, q4 *int) {
	getVal := func(v *int) int {
		if v == nil {
			return 0
		}
		return *v
	}
	quarters := []int{getVal(q1), getVal(q2), getVal(q3), getVal(q4)}
	for i, v := range quarters {
		idx := i*2 + slot
		for len(*scores) <= idx {
			*scores = append(*scores, 0)
		}
		(*scores)[idx] = v
	}
}

// parsePlayByPlay converts play-by-play events to api.MatchEvent slice.
func parsePlayByPlay(resp playByPlayResponse) []api.MatchEvent {
	pbp := findResultSet(resp.ResultSets, "PlayByPlay")

	var events []api.MatchEvent
	for _, row := range pbp.RowSet {
		msgType := pbp.colInt(row, "EVENTMSGTYPE")

		// Skip events we don't display
		if msgType == eventStartPeriod || msgType == eventEndPeriod ||
			msgType == eventRebound || msgType == eventViolation {
			continue
		}

		period := pbp.colInt(row, "PERIOD")
		clock := pbp.colStr(row, "PCTIMESTRING")
		homeDesc := pbp.colStr(row, "HOMEDESCRIPTION")
		awayDesc := pbp.colStr(row, "VISITORDESCRIPTION")
		neutralDesc := pbp.colStr(row, "NEUTRALDESCRIPTION")
		score := pbp.colStr(row, "SCORE")

		desc := homeDesc
		if desc == "" {
			desc = awayDesc
		}
		if desc == "" {
			desc = neutralDesc
		}

		eventType := msgTypeToString(msgType)
		displayMinute := fmt.Sprintf("Q%d %s", period, clock)

		event := api.MatchEvent{
			ID:            pbp.colInt(row, "EVENTNUM"),
			Minute:        period * 12,
			DisplayMinute: displayMinute,
			Type:          eventType,
		}

		// Detect 3-pointer and score events
		if msgType == eventFieldGoalMade {
			isThree := strings.Contains(desc, "3PT")
			event.IsThree = &isThree
			pts := 2
			if isThree {
				pts = 3
			}
			event.Points = &pts
		}
		if msgType == eventFreeThrowMade {
			pts := 1
			event.Points = &pts
		}

		_ = score // could be parsed to update running score

		events = append(events, event)
	}

	return events
}

// msgTypeToString maps EVENTMSGTYPE to a human-readable event type.
func msgTypeToString(t int) string {
	switch t {
	case eventFieldGoalMade:
		return "field_goal"
	case eventFieldGoalMissed:
		return "field_goal_missed"
	case eventFreeThrowMade:
		return "free_throw"
	case eventFreeThrowMissed:
		return "free_throw_missed"
	case eventFoul:
		return "foul"
	case eventSubstitution:
		return "substitution"
	case eventTimeout:
		return "timeout"
	case eventJumpBall:
		return "jump_ball"
	case eventEjection:
		return "ejection"
	default:
		return "other"
	}
}

// isPlayoffGame returns true if the game ID indicates a playoff game.
// NBA game IDs: position 2 is the game type (2=regular season, 4=playoffs).
func isPlayoffGame(gameID string) bool {
	if len(gameID) >= 3 {
		return gameID[2] == '4'
	}
	return false
}

// --- Game ID registry ---
// The NBA API uses string game IDs ("0022300789") while our api.Client interface
// uses int IDs. We maintain a registry to map between them.

var (
	gameIDRegistry = make(map[int]string) // numericID → string gameID
)

// StoreGameID stores the string game ID for later retrieval by numeric ID.
func StoreGameID(numericID int, stringID string) {
	gameIDRegistry[numericID] = stringID
}

// storedGameID retrieves the string game ID for a numeric ID.
func storedGameID(numericID int) string {
	return gameIDRegistry[numericID]
}

// simpleHash converts a string game ID to a stable int for use as a map key.
// Uses the last 6 digits of the game ID string.
func simpleHash(gameID string) int {
	if len(gameID) == 0 {
		return 0
	}
	// Use the last 8 characters as a number
	start := 0
	if len(gameID) > 8 {
		start = len(gameID) - 8
	}
	n := 0
	for _, c := range gameID[start:] {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		}
	}
	return n
}
