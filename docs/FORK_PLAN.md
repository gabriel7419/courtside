# NBA Fork — Implementation Plan

Reference document for the Courtside NBA adaptation of Golazo. For a shorter overview, see [FORK_SUMMARY.md](FORK_SUMMARY.md).

---

## Original Architecture (Golazo)

```
golazo/
├── cmd/                    # CLI (Cobra)
├── internal/
│   ├── api/               # Sport-agnostic interface
│   ├── fotmob/            # FotMob API client
│   ├── reddit/            # r/soccer highlight search
│   ├── ui/                # TUI (Bubble Tea)
│   ├── data/              # Settings, storage, mock data
│   ├── notify/            # Desktop notifications
│   ├── app/               # Core application logic
│   ├── constants/
│   ├── debug/
│   └── version/
├── assets/
├── scripts/
└── docs/
```

Key design advantages: well-abstracted `api.Client` interface, TTL-based cache, configurable rate limiting, and UI fully decoupled from data logic.

---

## Target Architecture (Courtside)

```
courtside/
├── cmd/                    # Renamed: courtside CLI
├── internal/
│   ├── api/               # Extended with NBA fields (Quarter, Clock, PlayerStatLine)
│   ├── nba/               # NEW — NBA Stats API client + mock client
│   ├── reddit/            # Adapted: r/nba Highlight posts
│   ├── ui/                # Adapted: quarter/clock display, box score, standings
│   ├── data/              # NBA mock data (matches, details, player stats)
│   ├── notify/            # NBA event labels (BASKET, 3PT, FT)
│   ├── app/               # NBA client wired in, mock data branches
│   ├── constants/         # NBA terminology (Games, Final, Conference, Arena)
│   ├── debug/             # Unchanged
│   └── version/           # Unchanged
├── assets/
├── scripts/
└── docs/
```

---

## API

### NBA Stats API (primary)

- **Base URL:** `https://stats.nba.com/stats/`
- **Cost:** Free
- **Auth:** None (but requires specific request headers)
- **Rate limiting:** Undocumented — use 200–300ms between requests

**Endpoints used:**

```
GET /scoreboard?GameDate=YYYY-MM-DD&LeagueID=00        → daily scoreboard
GET /boxscoresummaryv2?GameID=<id>                      → game summary
GET /boxscoretraditionalv2?GameID=<id>&...              → player + team stats
GET /playbyplayv2?GameID=<id>&StartPeriod=1&EndPeriod=10 → play-by-play events
GET /leaguestandingsv3?LeagueID=00&Season=<year>&...    → standings
```

See full details in [API_REFERENCE.md](API_REFERENCE.md).

---

## Football → NBA Mapping

### Data Structures

| Football (Golazo) | NBA (Courtside) |
|---|---|
| `Match` | `Match` (reused, extended) |
| `League` | Conference (`"NBA"`) |
| `Round` | Intentionally unused |
| `LiveTime` (`"45+2"`) | `Quarter` + `Clock` (`"Q3 2:34"`) |
| Half-time score | Quarter-by-quarter scores (`QuarterScores []int`) |
| `GoalsFor` / `GoalsAgainst` | `PointsFor` / `PointsAgainst` (win%) |

### Events

| Football | NBA |
|---|---|
| Goal | Field Goal (2pt/3pt), Free Throw |
| Yellow Card | Personal Foul / Technical Foul |
| Red Card | Ejection / Flagrant Foul |
| Substitution | Substitution |
| — | Timeout |

### Statistics

| Football | NBA |
|---|---|
| Possession | Time of Possession |
| Shots / Shots on Target | FGA / FGM |
| Passes | Assists |
| — | Rebounds (OREB/DREB) |
| — | Steals, Blocks, Turnovers |
| — | FG%, 3P%, FT% |

---

## Implementation Checklist

### Phase 0 — Preparation ✅
- [x] Implementation plan
- [x] README and CONTRIBUTING adapted  
- [x] API reference documentation
- [x] API test script

### Phase 1 — Setup ✅
- [x] Rename Go module (`go.mod`)
- [x] Update all internal imports
- [x] Confirm clean `go build`

### Phase 2 — Data Layer ✅
- [x] Extend `internal/api/types.go` (Quarter, Clock, PlayerStatLine, LeagueTableEntry)
- [x] Create `internal/nba/` package (client, types, cache, ratelimit, live parser)
- [x] Implement `MatchesByDate`, `MatchDetails`, `LiveMatches`, `LeagueTable`
- [x] Parse box score stats (`BoxScoreTraditionalV2` → team + player stats)
- [x] Adapt `internal/data/settings.go` (NBA conferences and teams)
- [x] Add NBA mock data for offline development (`internal/nba/mock_client.go`)

### Phase 3 — Core Functionality ✅
- [x] `MatchesByDate` — daily scoreboard
- [x] `MatchDetails` — box score summary + play-by-play
- [x] Cache and rate limiting (TTL: live 10s, finished 24h)

### Phase 4 — Live Data ✅
- [x] 30-second polling for live games
- [x] Map NBA events (field goals, fouls, timeouts, free throws)
- [x] Real-time score and quarter updates
- [x] `LiveUpdateParser` for play-by-play streaming

### Phase 5 — UI ✅
- [x] Quarter/clock display (`"Q3 2:34"`)
- [x] High-score formatting (NBA scores 90–130)
- [x] Standings dialog: W/L/PCT/GB/Streak columns, East/West sub-headers
- [x] Box score section: two-column player stats (PTS/REB/AST/FG)
- [x] Statistics dialog: FG%, 3P%, FT%, REB, AST, STL, BLK, TO, PF

### Phase 6 — Extra Features ✅
- [x] Highlights via r/nba (Highlight flair, NBA title matching)
- [x] NBA scoring notifications (BASKET +2, 3PT +3, FT +1, DisplayMinute)
- [x] Standings (conference standings from `leaguestandingsv3`)
- [x] Offline / API-unavailable mode (`--mock` flag)

### Phase 7 — Release 🔜
- [ ] README final screenshots
- [ ] Build and install scripts for Courtside
- [ ] Release v1.0.0

---

## Testing the API

```bash
# Daily scoreboard
go run scripts/test_nba_api.go --endpoint=scoreboard --date=2026-02-25

# Game summary
go run scripts/test_nba_api.go --endpoint=summary --game=0022300789

# Box score with player stats
go run scripts/test_nba_api.go --endpoint=traditional --game=0022300789

# Play-by-play
go run scripts/test_nba_api.go --endpoint=playbyplay --game=0022300789

# Conference standings
go run scripts/test_nba_api.go --endpoint=standings --season=2025-26

# Offline mode (no network required)
go run ./cmd/courtside --mock
```

---

## Resources

- [swar/nba_api](https://github.com/swar/nba_api) — comprehensive Python reference for NBA Stats API
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI framework used by Courtside
- [Original Golazo](https://github.com/0xjuanma/golazo) — base project this was forked from

---

*Based on [Golazo](https://github.com/0xjuanma/golazo) by [@0xjuanma](https://github.com/0xjuanma) — MIT License*
