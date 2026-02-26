package app

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gabriel7419/courtside/internal/api"
	"github.com/gabriel7419/courtside/internal/data"
	"github.com/gabriel7419/courtside/internal/nba"
	"github.com/gabriel7419/courtside/internal/reddit"
)

// LiveRefreshInterval is the interval between automatic live game list refreshes.
const LiveRefreshInterval = 30 * time.Second

// LiveBatchSize is kept for compatibility with the progressive load message types.
// For NBA, the scoreboard returns all games in a single call, so we always use 1 batch.
const LiveBatchSize = 1

// fetchLiveBatchData fetches all live NBA games in a single scoreboard call.
// The batch concept is kept for message compatibility — batchIndex 0 is always the last.
func fetchLiveBatchData(client *nba.Client, useMockData bool, batchIndex int) tea.Cmd {
	return func() tea.Msg {
		isLast := true // NBA: always one batch

		if useMockData {
			if batchIndex == 0 {
				return liveBatchDataMsg{
					batchIndex: batchIndex,
					isLast:     isLast,
					matches:    data.MockNBALiveMatches(),
				}
			}
			return liveBatchDataMsg{batchIndex: batchIndex, isLast: isLast}
		}

		if client == nil {
			return liveBatchDataMsg{batchIndex: batchIndex, isLast: isLast}
		}

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		matches, err := client.LiveMatches(ctx)
		if err != nil {
			return liveBatchDataMsg{batchIndex: batchIndex, isLast: isLast}
		}

		return liveBatchDataMsg{
			batchIndex: batchIndex,
			isLast:     isLast,
			matches:    matches,
		}
	}
}

// scheduleLiveRefresh schedules the next live game list refresh.
func scheduleLiveRefresh(client *nba.Client, useMockData bool) tea.Cmd {
	return tea.Tick(LiveRefreshInterval, func(t time.Time) tea.Msg {
		if useMockData {
			return liveRefreshMsg{matches: data.MockNBALiveMatches()}
		}
		if client == nil {
			return liveRefreshMsg{}
		}

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		matches, err := client.LiveMatchesForceRefresh(ctx)
		if err != nil {
			return liveRefreshMsg{}
		}
		return liveRefreshMsg{matches: matches}
	})
}

// fetchMatchDetails fetches game details from the NBA API.
func fetchMatchDetails(client *nba.Client, matchID int, useMockData bool) tea.Cmd {
	return func() tea.Msg {
		if useMockData {
			details, _ := data.MockMatchDetails(matchID)
			return matchDetailsMsg{details: details}
		}

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		details, err := client.MatchDetails(ctx, matchID)
		if err != nil {
			return matchDetailsMsg{}
		}
		return matchDetailsMsg{details: details}
	}
}

// fetchMatchDetailsForceRefresh fetches game details bypassing the cache.
func fetchMatchDetailsForceRefresh(client *nba.Client, matchID int, useMockData bool) tea.Cmd {
	return func() tea.Msg {
		if useMockData {
			details, _ := data.MockNBAMatchDetails(matchID)
			return matchDetailsMsg{details: details}
		}

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		details, err := client.MatchDetailsForceRefresh(ctx, matchID)
		if err != nil {
			return matchDetailsMsg{}
		}
		return matchDetailsMsg{details: details}
	}
}

// schedulePollTick schedules the next poll after 30 seconds (NBA games update quickly).
func schedulePollTick(matchID int) tea.Cmd {
	return tea.Tick(30*time.Second, func(t time.Time) tea.Msg {
		return pollTickMsg{matchID: matchID}
	})
}

// PollSpinnerDuration is how long to show the "Updating..." spinner.
const PollSpinnerDuration = 1 * time.Second

// schedulePollSpinnerHide schedules hiding the spinner after the display duration.
func schedulePollSpinnerHide() tea.Cmd {
	return tea.Tick(PollSpinnerDuration, func(t time.Time) tea.Msg {
		return pollDisplayCompleteMsg{}
	})
}

// fetchPollMatchDetails fetches game details for a live-polling refresh.
func fetchPollMatchDetails(client *nba.Client, matchID int, useMockData bool) tea.Cmd {
	return func() tea.Msg {
		if useMockData {
			details, _ := data.MockNBAMatchDetails(matchID)
			return matchDetailsMsg{details: details}
		}

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		details, err := client.MatchDetailsForceRefresh(ctx, matchID)
		if err != nil {
			return matchDetailsMsg{}
		}
		return matchDetailsMsg{details: details}
	}
}

// fetchStatsDayData fetches games for a single day (progressive loading).
// dayIndex: 0 = today, 1 = yesterday, etc.
func fetchStatsDayData(client *nba.Client, useMockData bool, dayIndex int, totalDays int) tea.Cmd {
	return func() tea.Msg {
		isToday := dayIndex == 0
		isLast := dayIndex == totalDays-1

		if useMockData {
			if isToday {
				// Use NBA live + finished games as "today's finished games"
				var finished []api.Match
				for _, m := range data.MockNBALiveMatches() {
					if m.Status == api.MatchStatusFinished {
						finished = append(finished, m)
					}
				}
				return statsDayDataMsg{
					dayIndex: dayIndex,
					isToday:  true,
					isLast:   isLast,
					finished: finished,
				}
			}
			return statsDayDataMsg{dayIndex: dayIndex, isToday: false, isLast: isLast}
		}

		if client == nil {
			return statsDayDataMsg{dayIndex: dayIndex, isToday: isToday, isLast: isLast}
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		date := time.Now().UTC().AddDate(0, 0, -dayIndex)
		matches, err := client.MatchesByDate(ctx, date)
		if err != nil {
			return statsDayDataMsg{dayIndex: dayIndex, isToday: isToday, isLast: isLast}
		}

		var finished, upcoming []api.Match
		for _, m := range matches {
			switch m.Status {
			case api.MatchStatusFinished:
				finished = append(finished, m)
			case api.MatchStatusNotStarted:
				if isToday {
					upcoming = append(upcoming, m)
				}
			}
		}

		return statsDayDataMsg{
			dayIndex: dayIndex,
			isToday:  isToday,
			isLast:   isLast,
			finished: finished,
			upcoming: upcoming,
		}
	}
}

// fetchStatsMatchDetails fetches game details for the stats (finished games) view.
func fetchStatsMatchDetails(client *nba.Client, matchID int, useMockData bool) tea.Cmd {
	return func() tea.Msg {
		if useMockData {
			details, _ := data.MockFinishedMatchDetails(matchID)
			return matchDetailsMsg{details: details}
		}
		if client == nil {
			return matchDetailsMsg{}
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		details, err := client.MatchDetails(ctx, matchID)
		if err != nil {
			return matchDetailsMsg{}
		}
		return matchDetailsMsg{details: details}
	}
}

// fetchGoalLinks fetches highlight links from Reddit for a game's scoring events.
func fetchGoalLinks(redditClient *reddit.Client, details *api.MatchDetails) tea.Cmd {
	return func() tea.Msg {
		if redditClient == nil || details == nil {
			return goalLinksMsg{matchID: 0}
		}

		var goals []reddit.GoalInfo
		for _, event := range details.Events {
			if event.Type != "field_goal" && event.Type != "goal" {
				continue
			}

			scorer := ""
			if event.Player != nil {
				scorer = *event.Player
			}

			isHome := event.Team.ID == details.HomeTeam.ID

			homeScore := 0
			if details.HomeScore != nil {
				homeScore = *details.HomeScore
			}
			awayScore := 0
			if details.AwayScore != nil {
				awayScore = *details.AwayScore
			}

			matchTime := time.Now()
			if details.MatchTime != nil {
				matchTime = *details.MatchTime
			}

			goals = append(goals, reddit.GoalInfo{
				MatchID:       details.ID,
				HomeTeam:      details.HomeTeam.Name,
				AwayTeam:      details.AwayTeam.Name,
				HomeTeamShort: details.HomeTeam.ShortName,
				AwayTeamShort: details.AwayTeam.ShortName,
				ScorerName:    scorer,
				Minute:        event.Minute,
				DisplayMinute: event.DisplayMinute,
				HomeScore:     homeScore,
				AwayScore:     awayScore,
				IsHomeTeam:    isHome,
				MatchTime:     matchTime,
			})
		}

		if len(goals) == 0 {
			return goalLinksMsg{matchID: details.ID}
		}

		links := redditClient.GoalLinks(goals)
		return goalLinksMsg{matchID: details.ID, links: links}
	}
}

// fetchStandings fetches standings from the NBA API (stub — always returns empty).
func fetchStandings(client *nba.Client, leagueID int, leagueName string, parentLeagueID int, homeTeamID, awayTeamID int) tea.Cmd {
	return func() tea.Msg {
		if client == nil {
			return standingsMsg{leagueID: leagueID}
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		standings, err := client.LeagueTable(ctx, leagueID, leagueName)
		if err != nil {
			return standingsMsg{leagueID: leagueID}
		}

		return standingsMsg{
			leagueID:   leagueID,
			leagueName: leagueName,
			standings:  standings,
			homeTeamID: homeTeamID,
			awayTeamID: awayTeamID,
		}
	}
}
