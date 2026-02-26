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

	// These headers are required — without them the API returns 403 or times out.
	// Mirrors the headers used by the nba_api Python library.
	req.Header.Set("Host", "stats.nba.com")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	// Note: do NOT set Accept-Encoding — Go's http.Client handles gzip transparently
	req.Header.Set("Referer", "https://www.nba.com/")
	req.Header.Set("Origin", "https://www.nba.com")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("x-nba-stats-origin", "stats")
	req.Header.Set("x-nba-stats-token", "true")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")

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

	url := fmt.Sprintf("%s/scoreboardv3?GameDate=%s&LeagueID=00", c.baseURL, dateStr)

	var resp scoreboardV3Response
	if err := c.do(ctx, url, &resp); err != nil {
		return nil, fmt.Errorf("fetch scoreboard for %s: %w", dateStr, err)
	}

	var matches []api.Match
	for _, g := range resp.Scoreboard.Games {
		// Map API status to internal status
		var status api.MatchStatus
		switch g.GameStatus {
		case gameStatusScheduled:
			status = api.MatchStatusNotStarted
		case gameStatusLive:
			status = api.MatchStatusLive
		case gameStatusFinal:
			status = api.MatchStatusFinished
		default:
			status = api.MatchStatusNotStarted
		}

		// Parse game time
		var matchTime *time.Time
		if g.GameTimeUTC != "" {
			// try RFC3339 first, then "2006-01-02T15:04:05Z"
			for _, layout := range []string{time.RFC3339, "2006-01-02T15:04:05Z"} {
				if t, err := time.Parse(layout, g.GameTimeUTC); err == nil {
					matchTime = &t
					break
				}
			}
		}

		// LiveTime: human-readable game time
		var liveTime *string
		if g.GameStatusText != "" {
			s := strings.TrimSpace(g.GameStatusText)
			liveTime = &s
		}

		// Quarter + clock
		var quarter *int
		var clock *string
		if g.Period > 0 {
			q := g.Period
			quarter = &q
		}
		if g.GameClock != "" && status == api.MatchStatusLive {
			clock = &g.GameClock
		}

		// Teams and scores
		homeScore := g.HomeTeam.Score
		awayScore := g.AwayTeam.Score
		homeTeam := api.Team{
			ID:        g.HomeTeam.TeamID,
			Name:      g.HomeTeam.TeamCity + " " + g.HomeTeam.TeamName,
			ShortName: g.HomeTeam.TeamTricode,
		}
		awayTeam := api.Team{
			ID:        g.AwayTeam.TeamID,
			Name:      g.AwayTeam.TeamCity + " " + g.AwayTeam.TeamName,
			ShortName: g.AwayTeam.TeamTricode,
		}

		numericID := simpleHash(g.GameID)
		storeGameID(numericID, g.GameID)

		// Series status for playoffs
		var seriesStatus *string
		if g.SeriesText != "" {
			seriesStatus = &g.SeriesText
		}

		// Quick build of quarter scores from Periods (v3 includes them inline)
		var quarterScores []int
		if len(g.HomeTeam.Periods) > 0 {
			for i := 0; i < len(g.HomeTeam.Periods) && i < len(g.AwayTeam.Periods); i++ {
				quarterScores = append(quarterScores, g.HomeTeam.Periods[i].Score, g.AwayTeam.Periods[i].Score)
			}
		}
		_ = quarterScores // used when building MatchDetails

		m := api.Match{
			ID:           numericID,
			League:       api.League{Name: "NBA"},
			HomeTeam:     homeTeam,
			AwayTeam:     awayTeam,
			Status:       status,
			HomeScore:    &homeScore,
			AwayScore:    &awayScore,
			MatchTime:    matchTime,
			LiveTime:     liveTime,
			Quarter:      quarter,
			Clock:        clock,
			IsPlayoffs:   isPlayoffGame(g.GameID),
			SeriesStatus: seriesStatus,
		}
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

	// Fetch play-by-play for live games (v3)
	if details.Status == api.MatchStatusLive {
		pbpURL := fmt.Sprintf("%s/playbyplayv3?GameID=%s&StartPeriod=1&EndPeriod=10", c.baseURL, gameIDStr)
		var pbpResp playByPlayV3Response
		if err := c.do(ctx, pbpURL, &pbpResp); err == nil {
			details.Events = parsePlayByPlayV3(pbpResp)
		}
	}

	// Fetch team + player box score stats (v3)
	statsURL := fmt.Sprintf("%s/boxscoretraditionalv3?GameID=%s&StartPeriod=1&EndPeriod=10&StartRange=0&EndRange=28800&RangeType=0", c.baseURL, gameIDStr)
	var statsResp boxScoreTraditionalV3Response
	if err := c.do(ctx, statsURL, &statsResp); err == nil {
		details.Statistics = parseTeamStatsV3(statsResp, details.HomeTeam.ID)
		home, away := parsePlayerStatsV3(statsResp, details.HomeTeam.ID)
		details.HomePlayerStats = home
		details.AwayPlayerStats = away
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
		statusText := gs.colStr(row, "GAME_STATUS_TEXT") // "Final", "Q3 2:34", "7:30 pm ET"
		livePeriod := gs.colInt(row, "LIVE_PERIOD")
		livePcTime := gs.colStr(row, "LIVE_PC_TIME") // remaining clock e.g "PT02M34.00S"

		switch statusID {
		case gameStatusScheduled:
			details.Status = api.MatchStatusNotStarted
		case gameStatusLive:
			details.Status = api.MatchStatusLive
		case gameStatusFinal:
			details.Status = api.MatchStatusFinished
		}

		// LiveTime — shown in the header below the score
		statusText = strings.TrimSpace(statusText)
		if statusText != "" {
			details.LiveTime = &statusText
		}

		// Live clock from LIVE_PC_TIME ("PT02M34.00S" → "2:34")
		if livePcTime != "" && details.Status == api.MatchStatusLive {
			clockFmt := formatISODuration(livePcTime)
			details.Clock = &clockFmt
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

	// Attendance and arena from GameInfo
	if len(gi.RowSet) > 0 {
		giRow := gi.RowSet[0]
		details.Attendance = gi.colInt(giRow, "ATTENDANCE")
		arena := gi.colStr(giRow, "GAME_TIME") // GameInfo has GAME_TIME; arena comes from GameSummary
		_ = arena
	}

	// Arena name is in GameSummary
	if len(gs.RowSet) > 0 {
		arena := gs.colStr(gs.RowSet[0], "ARENA_NAME")
		if arena != "" {
			details.Venue = arena
		}
	}

	return details
}

// parseTeamStatsV3 converts boxscoretraditionalv3 TeamStats resultSet → []api.MatchStatistic.
// v3 uses camelCase field names (fieldGoalsPercentage, reboundsTotal, etc.).
func parseTeamStatsV3(resp boxScoreTraditionalV3Response, homeTeamID int) []api.MatchStatistic {
	ts := findResultSet(resp.ResultSets, "TeamStats")
	if len(ts.RowSet) < 2 {
		return nil
	}

	var homeRow, awayRow []interface{}
	for _, row := range ts.RowSet {
		teamID := ts.colInt(row, "teamId")
		if teamID == homeTeamID {
			homeRow = row
		} else {
			awayRow = row
		}
	}
	if homeRow == nil || awayRow == nil {
		return nil
	}

	formatPct := func(row []interface{}, field string) string {
		v := ts.colFloat(row, field)
		if v == 0 {
			return "—"
		}
		return fmt.Sprintf("%.1f%%", v*100)
	}
	formatInt := func(row []interface{}, field string) string {
		return fmt.Sprintf("%d", ts.colInt(row, field))
	}
	stat := func(key, label, homeVal, awayVal string) api.MatchStatistic {
		return api.MatchStatistic{Key: key, Label: label, HomeValue: homeVal, AwayValue: awayVal}
	}

	return []api.MatchStatistic{
		stat("fg_pct", "FG%", formatPct(homeRow, "fieldGoalsPercentage"), formatPct(awayRow, "fieldGoalsPercentage")),
		stat("fg3_pct", "3P%", formatPct(homeRow, "threePointersPercentage"), formatPct(awayRow, "threePointersPercentage")),
		stat("ft_pct", "FT%", formatPct(homeRow, "freeThrowsPercentage"), formatPct(awayRow, "freeThrowsPercentage")),
		stat("reb", "Rebounds", formatInt(homeRow, "reboundsTotal"), formatInt(awayRow, "reboundsTotal")),
		stat("oreb", "Off. Rebounds", formatInt(homeRow, "reboundsOffensive"), formatInt(awayRow, "reboundsOffensive")),
		stat("dreb", "Def. Rebounds", formatInt(homeRow, "reboundsDefensive"), formatInt(awayRow, "reboundsDefensive")),
		stat("ast", "Assists", formatInt(homeRow, "assists"), formatInt(awayRow, "assists")),
		stat("stl", "Steals", formatInt(homeRow, "steals"), formatInt(awayRow, "steals")),
		stat("blk", "Blocks", formatInt(homeRow, "blocks"), formatInt(awayRow, "blocks")),
		stat("tov", "Turnovers", formatInt(homeRow, "turnovers"), formatInt(awayRow, "turnovers")),
		stat("pf", "Personal Fouls", formatInt(homeRow, "foulsPersonal"), formatInt(awayRow, "foulsPersonal")),
		stat("fga", "FG Attempted", formatInt(homeRow, "fieldGoalsAttempted"), formatInt(awayRow, "fieldGoalsAttempted")),
		stat("fg3a", "3P Attempted", formatInt(homeRow, "threePointersAttempted"), formatInt(awayRow, "threePointersAttempted")),
		stat("fta", "FT Attempted", formatInt(homeRow, "freeThrowsAttempted"), formatInt(awayRow, "freeThrowsAttempted")),
	}
}

// parsePlayerStatsV3 converts boxscoretraditionalv3 PlayerStats → home and away []api.PlayerStatLine.
// v3 uses camelCase field names. Players with no minutes (DNP) are excluded.
// Results sorted by points descending.
func parsePlayerStatsV3(resp boxScoreTraditionalV3Response, homeTeamID int) (home, away []api.PlayerStatLine) {
	ps := findResultSet(resp.ResultSets, "PlayerStats")
	if len(ps.RowSet) == 0 {
		return nil, nil
	}

	sortByPts := func(players []api.PlayerStatLine) {
		for i := 1; i < len(players); i++ {
			p := players[i]
			j := i - 1
			for j >= 0 && players[j].Points < p.Points {
				players[j+1] = players[j]
				j--
			}
			players[j+1] = p
		}
	}

	for _, row := range ps.RowSet {
		// v3 uses "minutes" field ("PT36M24.00S" or "" for DNP)
		mins := ps.colStr(row, "minutes")
		if mins == "" || mins == "PT00M00.00S" {
			continue
		}
		// Convert ISO duration "PT36M24.00S" → "36:24"
		if mins != "" {
			mins = formatISODuration(mins)
		}

		teamID := ps.colInt(row, "teamId")
		// v3 has firstName + familyName or nameI ("J. Tatum")
		name := ps.colStr(row, "nameI")
		if name == "" {
			name = ps.colStr(row, "firstName") + " " + ps.colStr(row, "familyName")
		}

		line := api.PlayerStatLine{
			Name:      name,
			Position:  ps.colStr(row, "position"),
			Minutes:   mins,
			Points:    ps.colInt(row, "points"),
			Rebounds:  ps.colInt(row, "reboundsTotal"),
			Assists:   ps.colInt(row, "assists"),
			Steals:    ps.colInt(row, "steals"),
			Blocks:    ps.colInt(row, "blocks"),
			Turnovers: ps.colInt(row, "turnovers"),
			FGM:       ps.colInt(row, "fieldGoalsMade"),
			FGA:       ps.colInt(row, "fieldGoalsAttempted"),
			FG3M:      ps.colInt(row, "threePointersMade"),
			FTM:       ps.colInt(row, "freeThrowsMade"),
			FTA:       ps.colInt(row, "freeThrowsAttempted"),
			PlusMinus: ps.colInt(row, "plusMinusPoints"),
		}

		if teamID == homeTeamID {
			home = append(home, line)
		} else {
			away = append(away, line)
		}
	}

	sortByPts(home)
	sortByPts(away)
	return home, away
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

// parsePlayByPlayV3 converts play-by-play events (v3, camelCase) to api.MatchEvent slice.
// v3 fields: actionType ("Made Shot"/"Foul"/"Timeout"), subType ("3PT"?), description,
// teamId, period, clock ("PT02M34.00S"), scoreHome, scoreAway, playerName.
func parsePlayByPlayV3(resp playByPlayV3Response) []api.MatchEvent {
	pbp := findResultSet(resp.ResultSets, "PlayByPlay")

	var events []api.MatchEvent
	for _, row := range pbp.RowSet {
		actionType := pbp.colStr(row, "actionType")

		// Skip non-displayable actions
		switch actionType {
		case "period", "game", "rebound", "violation", "":
			continue
		}

		period := pbp.colInt(row, "period")
		clock := pbp.colStr(row, "clock")     // "PT02M34.00S"
		clockFmt := formatISODuration(clock)  // → "2:34"
		subType := pbp.colStr(row, "subType") // e.g. "3pt" for 3-pointers
		desc := pbp.colStr(row, "description")
		playerName := pbp.colStr(row, "playerNameI") // "J. Tatum"
		teamIDVal := pbp.colInt(row, "teamId")
		scoreHome := pbp.colStr(row, "scoreHome")
		scoreAway := pbp.colStr(row, "scoreAway")

		displayMinute := fmt.Sprintf("Q%d %s", period, clockFmt)

		eventType := actionTypeToEventType(actionType)

		event := api.MatchEvent{
			ID:            pbp.colInt(row, "actionId"),
			Minute:        period * 12,
			DisplayMinute: displayMinute,
			Type:          eventType,
			Team:          api.Team{ID: teamIDVal},
		}
		if playerName != "" {
			event.Player = &playerName
		}

		if actionType == "Made Shot" || actionType == "field_goal" {
			isThree := strings.EqualFold(subType, "3pt")
			event.IsThree = &isThree
			pts := 2
			if isThree {
				pts = 3
			}
			event.Points = &pts
		}
		if actionType == "Free Throw" && strings.Contains(desc, "MADE") {
			pts := 1
			event.Points = &pts
		}

		_ = scoreHome
		_ = scoreAway

		events = append(events, event)
	}

	return events
}

// actionTypeToEventType maps v3 actionType strings to internal event types.
func actionTypeToEventType(actionType string) string {
	switch actionType {
	case "Made Shot", "field_goal":
		return "field_goal"
	case "Missed Shot":
		return "field_goal_missed"
	case "Free Throw":
		return "free_throw"
	case "Foul":
		return "foul"
	case "Substitution":
		return "substitution"
	case "Timeout":
		return "timeout"
	case "Jump Ball":
		return "jump_ball"
	case "Ejection":
		return "ejection"
	default:
		return "other"
	}
}

// formatISODuration converts "PT02M34.00S" → "2:34".
// Returns the input unchanged if it does not match the expected format.
func formatISODuration(s string) string {
	var m, sec float64
	if _, err := fmt.Sscanf(s, "PT%fM%fS", &m, &sec); err == nil {
		return fmt.Sprintf("%d:%02d", int(m), int(sec))
	}
	return s
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

// storeGameID is a package-private alias used by MatchesByDate.
func storeGameID(numericID int, stringID string) {
	StoreGameID(numericID, stringID)
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
