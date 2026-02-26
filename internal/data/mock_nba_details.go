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

	case 9004: // DEN 98 - OKC 115 (Final)
		d.QuarterScores = []int{
			24, 31, // Q1
			22, 28, // Q2
			28, 26, // Q3
			24, 30, // Q4
		}
		d.Events = nbaMockFinishedEvents(m)
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
