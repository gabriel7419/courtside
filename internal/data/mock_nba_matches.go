package data

import (
	"time"

	"github.com/gabriel7419/courtside/internal/api"
)

// NBA team IDs (aligned with the NBA Stats API)
const (
	teamBOS = 1610612738 // Boston Celtics
	teamBKN = 1610612751 // Brooklyn Nets
	teamNYK = 1610612752 // New York Knicks
	teamPHI = 1610612755 // Philadelphia 76ers
	teamTOR = 1610612761 // Toronto Raptors

	teamCHI = 1610612741 // Chicago Bulls
	teamCLE = 1610612739 // Cleveland Cavaliers
	teamDET = 1610612765 // Detroit Pistons
	teamIND = 1610612754 // Indiana Pacers
	teamMIL = 1610612749 // Milwaukee Bucks

	teamATL = 1610612737 // Atlanta Hawks
	teamCHA = 1610612766 // Charlotte Hornets
	teamMIA = 1610612748 // Miami Heat
	teamORL = 1610612753 // Orlando Magic
	teamWAS = 1610612764 // Washington Wizards

	teamDEN = 1610612743 // Denver Nuggets
	teamMIN = 1610612750 // Minnesota Timberwolves
	teamOKC = 1610612760 // Oklahoma City Thunder
	teamPOR = 1610612757 // Portland Trail Blazers
	teamUTA = 1610612762 // Utah Jazz

	teamGSW = 1610612744 // Golden State Warriors
	teamLAC = 1610612746 // LA Clippers
	teamLAL = 1610612747 // LA Lakers
	teamPHX = 1610612756 // Phoenix Suns
	teamSAC = 1610612758 // Sacramento Kings

	teamDAL = 1610612742 // Dallas Mavericks
	teamHOU = 1610612745 // Houston Rockets
	teamMEM = 1610612763 // Memphis Grizzlies
	teamNOP = 1610612740 // New Orleans Pelicans
	teamSAS = 1610612759 // San Antonio Spurs
)

// nbaTeam is a helper to build an api.Team for NBA.
func nbaTeam(id int, name, abbr string) api.Team {
	return api.Team{ID: id, Name: name, ShortName: abbr}
}

// MockNBALiveMatches returns realistic live NBA game data for development/offline use.
// Includes 3 live games (various quarters) and 2 finished games.
func MockNBALiveMatches() []api.Match {
	now := time.Now()
	q3Time := "Q3 4:52"
	q2Time := "Q2 1:18"
	finalStr := "Final"

	three := 3
	two := 2

	return []api.Match{
		// Game 1: BOS vs MIA — Q3 live
		{
			ID:        9001,
			League:    api.League{ID: 1, Name: "NBA"},
			HomeTeam:  nbaTeam(teamBOS, "Boston Celtics", "BOS"),
			AwayTeam:  nbaTeam(teamMIA, "Miami Heat", "MIA"),
			Status:    api.MatchStatusLive,
			HomeScore: intPtr(87),
			AwayScore: intPtr(79),
			MatchTime: &now,
			LiveTime:  &q3Time,
			Quarter:   &three,
		},
		// Game 2: LAL vs GSW — Q2 live
		{
			ID:        9002,
			League:    api.League{ID: 2, Name: "NBA"},
			HomeTeam:  nbaTeam(teamLAL, "Los Angeles Lakers", "LAL"),
			AwayTeam:  nbaTeam(teamGSW, "Golden State Warriors", "GSW"),
			Status:    api.MatchStatusLive,
			HomeScore: intPtr(51),
			AwayScore: intPtr(58),
			MatchTime: &now,
			LiveTime:  &q2Time,
			Quarter:   &two,
		},
		// Game 3: MIL vs PHI — Final
		{
			ID:        9003,
			League:    api.League{ID: 1, Name: "NBA"},
			HomeTeam:  nbaTeam(teamMIL, "Milwaukee Bucks", "MIL"),
			AwayTeam:  nbaTeam(teamPHI, "Philadelphia 76ers", "PHI"),
			Status:    api.MatchStatusFinished,
			HomeScore: intPtr(112),
			AwayScore: intPtr(104),
			MatchTime: &now,
			LiveTime:  &finalStr,
		},
		// Game 4: DEN vs OKC — Final
		{
			ID:        9004,
			League:    api.League{ID: 2, Name: "NBA"},
			HomeTeam:  nbaTeam(teamDEN, "Denver Nuggets", "DEN"),
			AwayTeam:  nbaTeam(teamOKC, "Oklahoma City Thunder", "OKC"),
			Status:    api.MatchStatusFinished,
			HomeScore: intPtr(98),
			AwayScore: intPtr(115),
			MatchTime: &now,
			LiveTime:  &finalStr,
		},
	}
}

// MockNBAUpcomingMatches returns scheduled NBA games for tonight.
func MockNBAUpcomingMatches() []api.Match {
	tonight := time.Now().Truncate(24 * time.Hour).Add(19 * time.Hour) // 7 PM tipoff
	later := tonight.Add(30 * time.Minute)

	return []api.Match{
		{
			ID:        9010,
			League:    api.League{ID: 1, Name: "NBA"},
			HomeTeam:  nbaTeam(teamNYK, "New York Knicks", "NYK"),
			AwayTeam:  nbaTeam(teamCLE, "Cleveland Cavaliers", "CLE"),
			Status:    api.MatchStatusNotStarted,
			MatchTime: &tonight,
		},
		{
			ID:        9011,
			League:    api.League{ID: 2, Name: "NBA"},
			HomeTeam:  nbaTeam(teamPHX, "Phoenix Suns", "PHX"),
			AwayTeam:  nbaTeam(teamSAC, "Sacramento Kings", "SAC"),
			Status:    api.MatchStatusNotStarted,
			MatchTime: &later,
		},
	}
}
