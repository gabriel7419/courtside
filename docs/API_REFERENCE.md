# NBA Stats API — Reference

Quick reference for the endpoints used in Courtside. For architectural context, see [FORK_PLAN.md](FORK_PLAN.md).

## Required Headers

All requests to `stats.nba.com` require these headers — without them the API returns 403:

```http
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36
Referer: https://www.nba.com/
Accept: application/json
Accept-Language: en-US,en;q=0.9
Origin: https://www.nba.com
```

---

## Endpoints

### 1. Scoreboard — Daily Games

```
GET https://stats.nba.com/stats/scoreboard?GameDate=YYYY-MM-DD&LeagueID=00
```

**Parameters:**
- `GameDate` — date in `YYYY-MM-DD` format
- `LeagueID` — always `00` for NBA

**Key result sets:** `GameHeader`, `LineScore`

**GameHeader fields:**

| Field | Type | Description |
|---|---|---|
| `GAME_ID` | string | Unique game ID, e.g. `"0022300789"` |
| `GAME_STATUS_ID` | int | `1` = scheduled, `2` = live, `3` = final |
| `GAME_STATUS_TEXT` | string | `"Live"`, `"Final"`, `"9:00 PM ET"` |
| `HOME_TEAM_ID` | int | Home team ID |
| `VISITOR_TEAM_ID` | int | Away team ID |
| `LIVE_PERIOD` | int | Current quarter (1–4, 5+ = OT) |
| `LIVE_PC_TIME` | string | Remaining time in ISO format, e.g. `"PT2M34.00S"` |
| `LIVE_PERIOD_TIME_BCAST` | string | Broadcast format, e.g. `"Q3 2:34"` |

**LineScore fields:**

| Field | Description |
|---|---|
| `TEAM_ID` | Team identifier |
| `PTS_QTR1..4` | Points per quarter (null if not yet played) |
| `PTS` | Total points |
| `FG_PCT` | Field goal % (e.g. `0.488` = 48.8%) |
| `FG3_PCT` | 3-pointer % |
| `FT_PCT` | Free throw % |
| `AST` | Assists |
| `REB` | Total rebounds |
| `TOV` | Turnovers |

---

### 2. Box Score Summary

```
GET https://stats.nba.com/stats/boxscoresummaryv2?GameID=0022300789
```

Returns: `GameSummary`, `LineScore`, `LastMeeting`, `SeasonSeries`, `Officials`, `GameInfo`.

**GameInfo fields:** `GAME_DATE`, `ATTENDANCE`, `GAME_TIME`

---

### 3. Box Score Traditional — Full Stats

```
GET https://stats.nba.com/stats/boxscoretraditionalv2?GameID=0022300789&StartPeriod=1&EndPeriod=10&StartRange=0&EndRange=28800&RangeType=0
```

Returns `PlayerStats` and `TeamStats` result sets.

**Player row example:** LeBron James — 36 min, 12/20 FG, 3/7 3P, 5/6 FT, 9 REB, 7 AST, 32 PTS, +5

**Field glossary:**

| Field | Meaning |
|---|---|
| `MIN` | Minutes played (`"MM:SS"`) |
| `FGM` / `FGA` | Field Goals Made / Attempted |
| `FG_PCT` | Field goal percentage |
| `FG3M` / `FG3A` | 3-Pointers Made / Attempted |
| `FG3_PCT` | 3-point percentage |
| `FTM` / `FTA` | Free Throws Made / Attempted |
| `FT_PCT` | Free throw percentage |
| `OREB` / `DREB` | Offensive / Defensive Rebounds |
| `REB` | Total rebounds |
| `AST` | Assists |
| `STL` | Steals |
| `BLK` | Blocks |
| `TO` | Turnovers |
| `PF` | Personal fouls |
| `PTS` | Points |
| `PLUS_MINUS` | Point differential while on court |

---

### 4. Play-by-Play

```
GET https://stats.nba.com/stats/playbyplayv2?GameID=0022300789&StartPeriod=1&EndPeriod=10
```

Returns a `PlayByPlay` result set. One row per event.

**EVENTMSGTYPE values:**

| Value | Event |
|---|---|
| 1 | Field Goal Made |
| 2 | Field Goal Missed |
| 3 | Free Throw Made |
| 4 | Free Throw Missed |
| 5 | Rebound |
| 6 | Personal Foul |
| 7 | Violation |
| 8 | Substitution |
| 9 | Timeout |
| 10 | Jump Ball |
| 11 | Ejection |
| 12 | Start of Period |
| 13 | End of Period |

The description is in `HOMEDESCRIPTION` or `VISITORDESCRIPTION` depending on which team the event belongs to.

---

### 5. League Standings

```
GET https://stats.nba.com/stats/leaguestandingsv3?LeagueID=00&Season=2025-26&SeasonType=Regular+Season
```

Returns a `Standings` result set with one row per team.

**Key fields:** `TeamAbbreviation`, `TeamCity`, `TeamName`, `Conference`, `ConferenceRank`, `WINS`, `LOSSES`, `WinPCT`, `ConferenceGamesBack`, `CurrentStreak`

---

## Best Practices

**Rate limiting:** The API does not document limits. Use 200–300ms between requests to avoid throttling.

**Caching strategy:**

| Data | TTL |
|---|---|
| Live game data | 10 seconds |
| Daily scoreboard | 30 seconds |
| Finished games | 24 hours (scores never change) |
| Player stats | 1 minute (live), 24 hours (final) |

---

## Game ID Format

```
0022300789
│││││└───── Sequential game number
││││└────── Type: 2 = Regular Season, 4 = Playoffs
│││└─────── Season: 23 = 2023-24
││└──────── Century (0 = 2000s)
│└───────── Always 0
└────────── Always 0
```

---

## All 30 Team IDs

| Team | ID |
|---|---|
| Atlanta Hawks | 1610612737 |
| Boston Celtics | 1610612738 |
| Cleveland Cavaliers | 1610612739 |
| New Orleans Pelicans | 1610612740 |
| Chicago Bulls | 1610612741 |
| Dallas Mavericks | 1610612742 |
| Denver Nuggets | 1610612743 |
| Golden State Warriors | 1610612744 |
| Houston Rockets | 1610612745 |
| LA Clippers | 1610612746 |
| Los Angeles Lakers | 1610612747 |
| Miami Heat | 1610612748 |
| Milwaukee Bucks | 1610612749 |
| Minnesota Timberwolves | 1610612750 |
| Brooklyn Nets | 1610612751 |
| New York Knicks | 1610612752 |
| Orlando Magic | 1610612753 |
| Indiana Pacers | 1610612754 |
| Philadelphia 76ers | 1610612755 |
| Phoenix Suns | 1610612756 |
| Portland Trail Blazers | 1610612757 |
| Sacramento Kings | 1610612758 |
| San Antonio Spurs | 1610612759 |
| Oklahoma City Thunder | 1610612760 |
| Toronto Raptors | 1610612761 |
| Utah Jazz | 1610612762 |
| Memphis Grizzlies | 1610612763 |
| Washington Wizards | 1610612764 |
| Detroit Pistons | 1610612765 |
| Charlotte Hornets | 1610612766 |

Full reference: [swar/nba_api](https://github.com/swar/nba_api/blob/master/src/nba_api/stats/library/data.py)

---

## Testing

```bash
go run scripts/test_nba_api.go --endpoint=scoreboard --date=2026-02-25
go run scripts/test_nba_api.go --endpoint=summary --game=0022300789
go run scripts/test_nba_api.go --endpoint=traditional --game=0022300789
go run scripts/test_nba_api.go --endpoint=playbyplay --game=0022300789
go run scripts/test_nba_api.go --endpoint=standings --season=2025-26
```
