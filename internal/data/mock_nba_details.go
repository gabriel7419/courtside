package data

import (
	"github.com/gabriel7419/courtside/internal/api"
)

// MockNBAMatchDetails returns realistic NBA game details for development/offline use.
// Covers the games defined in MockNBALiveMatches.
func MockNBAMatchDetails(matchID int) (*api.MatchDetails, error) {
	allMatches := append(MockNBALiveMatches(), MockNBAUpcomingMatches()...)
	for _, m := range allMatches {
		if m.ID == matchID {
			return buildNBADetails(m), nil
		}
	}
	return nil, nil
}

func buildNBADetails(m api.Match) *api.MatchDetails {
	d := &api.MatchDetails{
		Match:      m,
		Statistics: nbaMockStats(m.ID),
		Venue:      nbaMockArena(m.ID),
		Attendance: nbaMockAttendance(m.ID),
	}

	switch m.ID {
	case 9001: // BOS 87 - MIA 79  (Q3 live)
		three := true
		two := false
		pts3 := 3
		pts2 := 2
		pts1 := 1
		q := 3
		d.Quarter = &q
		d.QuarterScores = []int{
			28, 22, // Q1: BOS 28, MIA 22
			26, 24, // Q2: BOS 26, MIA 24
			33, 33, // Q3 (partial)
			0, 0,
		}
		d.Events = []api.MatchEvent{
			{ID: 1, DisplayMinute: "Q1 9:14", Type: "field_goal", Team: m.HomeTeam, Player: strp("J. Brown"), IsThree: &two, Points: &pts2},
			{ID: 2, DisplayMinute: "Q1 7:02", Type: "field_goal", Team: m.AwayTeam, Player: strp("B. Adebayo"), Points: &pts2},
			{ID: 3, DisplayMinute: "Q1 4:33", Type: "field_goal", Team: m.HomeTeam, Player: strp("J. Tatum"), IsThree: &three, Points: &pts3},
			{ID: 4, DisplayMinute: "Q1 2:10", Type: "free_throw", Team: m.AwayTeam, Player: strp("J. Butler"), Points: &pts1},
			{ID: 5, DisplayMinute: "Q2 8:51", Type: "field_goal", Team: m.HomeTeam, Player: strp("A. Horford"), Points: &pts2},
			{ID: 6, DisplayMinute: "Q2 5:23", Type: "field_goal", Team: m.AwayTeam, Player: strp("T. Herro"), IsThree: &three, Points: &pts3},
			{ID: 7, DisplayMinute: "Q2 2:44", Type: "timeout", Team: m.AwayTeam},
			{ID: 8, DisplayMinute: "Q3 10:00", Type: "field_goal", Team: m.HomeTeam, Player: strp("J. Tatum"), Points: &pts2},
			{ID: 9, DisplayMinute: "Q3 7:34", Type: "foul", Team: m.AwayTeam, Player: strp("B. Adebayo")},
			{ID: 10, DisplayMinute: "Q3 4:52", Type: "field_goal", Team: m.HomeTeam, Player: strp("J. Brown"), IsThree: &three, Points: &pts3},
		}

	case 9002: // LAL 51 - GSW 58  (Q2 live)
		q := 2
		pts3 := 3
		pts2 := 2
		pts1 := 1
		three := true
		two := false
		d.Quarter = &q
		d.QuarterScores = []int{
			28, 31, // Q1
			23, 27, // Q2 (partial)
			0, 0, 0, 0,
		}
		d.Events = []api.MatchEvent{
			{ID: 11, DisplayMinute: "Q1 10:12", Type: "field_goal", Team: m.AwayTeam, Player: strp("S. Curry"), IsThree: &three, Points: &pts3},
			{ID: 12, DisplayMinute: "Q1 8:00", Type: "field_goal", Team: m.HomeTeam, Player: strp("A. Davis"), Points: &pts2},
			{ID: 13, DisplayMinute: "Q1 5:44", Type: "free_throw", Team: m.HomeTeam, Player: strp("L. James"), Points: &pts1},
			{ID: 14, DisplayMinute: "Q1 3:20", Type: "field_goal", Team: m.AwayTeam, Player: strp("K. Thompson"), IsThree: &three, Points: &pts3},
			{ID: 15, DisplayMinute: "Q2 9:05", Type: "field_goal", Team: m.AwayTeam, Player: strp("D. Green"), Points: &pts2},
			{ID: 16, DisplayMinute: "Q2 6:30", Type: "timeout", Team: m.HomeTeam},
			{ID: 17, DisplayMinute: "Q2 4:10", Type: "field_goal", Team: m.HomeTeam, Player: strp("L. James"), Points: &pts2},
			{ID: 18, DisplayMinute: "Q2 1:18", Type: "foul", Team: m.HomeTeam, Player: strp("R. Westbrook")},
		}
		_ = two

	case 9003: // MIL 112 - PHI 104 (Final)
		d.QuarterScores = []int{
			32, 28, // Q1
			28, 30, // Q2
			26, 24, // Q3
			26, 22, // Q4
		}
		d.Events = nbaMockFinishedEvents(m)
		d.HomePlayerStats = nbaMockPlayers(m.HomeTeam, []playerSeed{
			{"D. Lillard", "G", "36:12", 34, 4, 7, 3, 9, 10, 3, 1, 8},
			{"K. Middleton", "F", "30:44", 18, 5, 4, 6, 8, 0, 0, 2, 2},
			{"B. Lopez", "C", "28:55", 14, 9, 2, 5, 8, 2, 0, 4, 0},
			{"K. Portis", "F", "24:18", 12, 8, 1, 4, 9, 1, 1, 2, 0},
			{"M. Beauchamp", "G", "22:00", 10, 3, 3, 4, 7, 1, 0, 2, 0},
			{"J. Crowder", "F", "18:30", 8, 4, 2, 3, 6, 0, 1, 1, 0},
		})
		d.AwayPlayerStats = nbaMockPlayers(m.AwayTeam, []playerSeed{
			{"J. Embiid", "C", "34:02", 32, 11, 5, 11, 20, 2, 2, 4, 0},
			{"T. Maxey", "G", "36:18", 24, 3, 8, 8, 18, 4, 0, 2, 0},
			{"O. Melton", "G", "28:00", 16, 4, 3, 6, 13, 3, 1, 2, 0},
			{"K. Lowry", "G", "24:44", 10, 2, 7, 3, 8, 1, 0, 1, 0},
			{"P. Reed", "F", "20:10", 8, 6, 1, 3, 7, 0, 1, 2, 0},
		})

	case 9004: // DEN 98 - OKC 115 (Final)
		d.QuarterScores = []int{
			24, 31, // Q1
			22, 28, // Q2
			28, 26, // Q3
			24, 30, // Q4
		}
		d.Events = nbaMockFinishedEvents(m)
		d.HomePlayerStats = nbaMockPlayers(m.HomeTeam, []playerSeed{
			{"N. Jokic", "C", "35:30", 28, 14, 9, 10, 18, 2, 2, 3, 0},
			{"J. Murray", "G", "33:00", 22, 5, 7, 8, 16, 3, 0, 4, 0},
			{"M. Porter Jr.", "F", "28:10", 16, 7, 2, 5, 11, 4, 0, 2, 0},
			{"K. Caldwell-Pope", "G", "24:00", 12, 3, 3, 4, 8, 2, 0, 2, 0},
			{"A. Gordon", "F", "26:20", 10, 6, 2, 3, 7, 0, 1, 2, 0},
		})
		d.AwayPlayerStats = nbaMockPlayers(m.AwayTeam, []playerSeed{
			{"S. Gilgeous-Alexander", "G", "36:00", 35, 5, 6, 13, 24, 5, 1, 4, 0},
			{"J. Williams", "F", "32:14", 22, 8, 4, 7, 14, 4, 1, 2, 0},
			{"C. Holmgren", "C", "30:00", 18, 10, 3, 6, 12, 2, 4, 2, 0},
			{"L. Dort", "G", "28:00", 14, 4, 2, 5, 12, 3, 0, 2, 0},
			{"I. Joe", "G", "22:10", 12, 3, 4, 4, 9, 2, 0, 1, 0},
		})
	}

	return d
}

func nbaMockFinishedEvents(m api.Match) []api.MatchEvent {
	pts3 := 3
	pts2 := 2
	pts1 := 1
	three := true
	two := false

	players := []string{"Player A", "Player B", "Player C"}
	teams := []api.Team{m.HomeTeam, m.AwayTeam, m.HomeTeam, m.AwayTeam}

	var events []api.MatchEvent
	quarters := []string{"Q1", "Q2", "Q3", "Q4"}
	times := []string{"9:30", "7:15", "4:48", "2:00"}
	id := 100 + m.ID

	for qi, q := range quarters {
		for ti, t := range times {
			team := teams[(qi+ti)%2]
			pl := players[(qi+ti)%3]
			isThree := (qi+ti)%3 == 0
			eventType := "field_goal"
			pts := &pts2
			isT := &two
			if isThree {
				pts = &pts3
				isT = &three
			}
			if (qi+ti)%5 == 0 {
				eventType = "free_throw"
				pts = &pts1
				isT = &two
			}
			events = append(events, api.MatchEvent{
				ID:            id,
				DisplayMinute: q + " " + t,
				Type:          eventType,
				Team:          team,
				Player:        &pl,
				IsThree:       isT,
				Points:        pts,
			})
			id++
		}
	}
	_ = three
	return events
}

func nbaMockStats(id int) []api.MatchStatistic {
	base := []api.MatchStatistic{
		{Key: "fg_pct", Label: "FG%", HomeValue: "48.2%", AwayValue: "44.7%"},
		{Key: "fg3_pct", Label: "3P%", HomeValue: "38.5%", AwayValue: "35.2%"},
		{Key: "ft_pct", Label: "FT%", HomeValue: "82.0%", AwayValue: "78.6%"},
		{Key: "reb", Label: "Rebounds", HomeValue: "44", AwayValue: "38"},
		{Key: "ast", Label: "Assists", HomeValue: "27", AwayValue: "22"},
		{Key: "stl", Label: "Steals", HomeValue: "8", AwayValue: "6"},
		{Key: "blk", Label: "Blocks", HomeValue: "5", AwayValue: "4"},
		{Key: "tov", Label: "Turnovers", HomeValue: "12", AwayValue: "15"},
		{Key: "pf", Label: "Fouls", HomeValue: "18", AwayValue: "21"},
	}
	// Vary slightly by game ID for visual interest
	if id%2 == 0 {
		base[0].HomeValue, base[0].AwayValue = "45.1%", "51.3%"
		base[3].HomeValue, base[3].AwayValue = "39", "47"
	}
	return base
}

func nbaMockArena(id int) string {
	arenas := map[int]string{
		9001: "TD Garden",
		9002: "Crypto.com Arena",
		9003: "Fiserv Forum",
		9004: "Ball Arena",
		9010: "Madison Square Garden",
		9011: "Footprint Center",
	}
	if a, ok := arenas[id]; ok {
		return a
	}
	return "NBA Arena"
}

func nbaMockAttendance(id int) int {
	m := map[int]int{
		9001: 19156,
		9002: 18997,
		9003: 17341,
		9004: 19520,
	}
	if a, ok := m[id]; ok {
		return a
	}
	return 18000
}

func strp(s string) *string { return &s }

// playerSeed holds the raw data for one mock player row.
// Fields: name, pos, min, pts, reb, ast, fgm, fga, fg3m, ftm, fta, plusMinus
type playerSeed struct {
	name, pos, min                          string
	pts, reb, ast, fgm, fga, fg3m, ftm, fta int
	pm                                      int
}

// nbaMockPlayers converts a slice of playerSeeds to []api.PlayerStatLine
// assigned to the given team.
func nbaMockPlayers(team api.Team, seeds []playerSeed) []api.PlayerStatLine {
	lines := make([]api.PlayerStatLine, 0, len(seeds))
	for _, s := range seeds {
		lines = append(lines, api.PlayerStatLine{
			Name:      s.name,
			Position:  s.pos,
			Minutes:   s.min,
			Points:    s.pts,
			Rebounds:  s.reb,
			Assists:   s.ast,
			FGM:       s.fgm,
			FGA:       s.fga,
			FG3M:      s.fg3m,
			FTM:       s.ftm,
			FTA:       s.fta,
			PlusMinus: s.pm,
		})
		_ = team // team association is implicit via HomePlayerStats / AwayPlayerStats
	}
	return lines
}
