# NBA Fork — Overview

This document summarizes how Courtside relates to Golazo and what the key differences are between the two projects.

## What is Courtside?

Courtside is a fork of [Golazo](https://github.com/0xjuanma/golazo) adapted for following NBA games in the terminal. The codebase is essentially the same — the difference is in the data layer (`internal/nba/` instead of `internal/fotmob/`) and some UI adjustments (labels, score formatting, etc.).

Most of the work involved implementing the NBA Stats API client and mapping its responses to the project's internal types.

## Key Differences

| Component | Golazo | Courtside |
|---|---|---|
| API | FotMob | NBA Stats API |
| Client package | `internal/fotmob/` | `internal/nba/` |
| Core struct | `Match` | `Game` |
| Time display | `"45'"` | `"Q3 2:34"` |
| Score range | Low (0–5) | High (90–120) |
| Competitions | 50+ worldwide leagues | Eastern/Western + playoffs |
| Highlights | r/soccer | r/nba |

## Architecture Advantages

Golazo was designed with clear separation of concerns:

- `internal/api/` defines a sport-agnostic interface — swap the client without touching the UI
- The TTL-based cache already exists and works
- Rate limiting is already configurable
- About 80% of the UI can be reused (only labels and formatting need updating)

## Expected Challenges

- The NBA Stats API is not officially documented and may change without notice
- Responses use a tabular format (`headers[]` + `rowSet[][]`), requiring manual index mapping
- Basketball has more event types than football (field goals, fouls, timeouts, jump balls, etc.)
- High scores (e.g., 98–105) need layout adjustments in the UI

## What's Already Done

- Full documentation (this file, FORK_PLAN, API_REFERENCE)
- API test script (`scripts/test_nba_api.go`)
- Football-to-NBA concept mapping
- Proposed structure for `internal/nba/`
- Full NBA client implementation (`internal/nba/client.go`)
- UI adaptation for basketball (quarter/clock display, scoring events)

## Next Steps

Implement the `LiveUpdateParser` for NBA-specific events and refine the play-by-play display. See [FORK_PLAN.md](FORK_PLAN.md) for the full roadmap.
