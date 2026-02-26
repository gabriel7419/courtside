# Courtside — Fork Summary

Courtside is a fork of [Golazo](https://github.com/0xjuanma/golazo) adapted for NBA. This document summarises what changed and what remains the same.

## What Changed

| Area | Change |
|---|---|
| **Data source** | FotMob API → NBA Stats API (`stats.nba.com`) |
| **Client package** | `internal/fotmob/` → `internal/nba/` |
| **Events** | Goals/cards → Field goals (2pt/3pt), free throws, fouls, timeouts |
| **Live time** | `"45+2'"` → `"Q3 2:34"` |
| **Score display** | Football low scores → NBA high scores (90–130) |
| **Standings** | League table (P/W/D/L/GD/Pts) → Conference standings (W/L/PCT/GB/Streak) |
| **Statistics dialog** | Possession/shots → FG%/3P%/FT%/REB/AST/STL/BLK/TO |
| **Box score** | — (new) → Two-column player stats per team (PTS/REB/AST/FG) |
| **Highlights** | `r/soccer` Media posts → `r/nba` Highlight posts |
| **Notifications** | "GOLAZO!" → "🏀 Courtside!" with event label (BASKET/3PT/FT) |
| **CLI name** | `golazo` → `courtside` |
| **Module path** | `github.com/0xjuanma/golazo` → `github.com/gabriel7419/courtside` |

## What Stayed the Same

- TUI framework: [Bubble Tea](https://github.com/charmbracelet/bubbletea) + Lip Gloss
- `api.Client` interface (extended with NBA fields, not replaced)
- Cache + rate-limiting infrastructure
- Desktop notifications backend (`beeep`)
- Reddit client structure (only URL/flair changed)
- Settings storage and version checking

## New in Courtside (not in Golazo)

- `--mock` flag for fully offline development (no network required)
- Player-level box score from `boxscoretraditionalv2`
- Quarter-by-quarter score breakdown
- Playoff series status support (`SeriesStatus`)

---

*Based on [Golazo](https://github.com/0xjuanma/golazo) by [@0xjuanma](https://github.com/0xjuanma) — MIT License*
