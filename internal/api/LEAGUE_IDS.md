# League ID Reference

This document tracks the league IDs used by different APIs in this application.

## Supported Leagues

The application currently supports 5 major leagues:
- Premier League (England)
- La Liga (Spain)
- Bundesliga (Germany)
- Serie A (Italy)
- Ligue 1 (France)

## API-Specific League IDs

### FotMob API
**Location:** `internal/fotmob/client.go`

```go
SupportedLeagues = []int{
    47, // Premier League
    87, // La Liga
    54, // Bundesliga
    55, // Serie A
    53, // Ligue 1
}
```

**API Endpoint:** `https://www.fotmob.com/api/leagues?id={leagueID}&tab=fixtures`

### API-Sports.io (Stats View)
**Location:** `internal/footballdata/client.go`

```go
SupportedLeagues = []int{
    39,  // Premier League
    140, // La Liga
    78,  // Bundesliga
    135, // Serie A
    61,  // Ligue 1
}
```

**API Endpoint:** `https://v3.football.api-sports.io/fixtures?league={leagueID}&date={date}`

## League ID Mapping

| League | FotMob ID | API-Sports.io ID |
|--------|-----------|------------------|
| Premier League | 47 | 39 |
| La Liga | 87 | 140 |
| Bundesliga | 54 | 78 |
| Serie A | 55 | 135 |
| Ligue 1 | 53 | 61 |

## Notes

- **FotMob** is used for the **Live Matches** view
- **API-Sports.io** is used for the **Stats** view (finished matches)
- League IDs are **different** between APIs - do not mix them up
- When adding new leagues, update both files and this document

