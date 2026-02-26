<div align="center">
  <h1>🏀 Courtside</h1>
</div>

<div align="center">

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/gabriel7419/courtside)](https://goreportcard.com/report/github.com/gabriel7419/courtside)
![macOS](https://img.shields.io/badge/macOS-000000?logo=apple&logoColor=white)
![Linux](https://img.shields.io/badge/Linux-FCC624?logo=linux&logoColor=black)
![Windows](https://img.shields.io/badge/Windows-0078D6?logo=windows&logoColor=white)

A minimalist terminal UI (TUI) for following NBA games in real-time — live scores, play-by-play events, and box score stats, all without leaving your terminal.

*Inspired by [Golazo](https://github.com/0xjuanma/golazo) — adapted for basketball.*

</div>

> [!NOTE]
> This project is under active development. See the [roadmap](#roadmap) for what's coming.

## Features

- **Live updates** — scores, fouls, timeouts, and substitutions with automatic polling
- **Box score stats** — FG%, rebounds, assists, steals, blocks, turnovers in a focused dialog
- **Finished games** — results from today, last 3 days, or last 5 days
- **Conference filtering** — Eastern and Western, with playoff series support
- **Highlight links** — links to r/nba highlights
- **Desktop notifications** — for key moments during live games

## What's Different from Golazo?

| | Golazo | Courtside |
|---|---|---|
| Sport | Football/Soccer | Basketball (NBA) |
| Data source | FotMob API | NBA Stats API |
| Game structure | 2 halves (45 min) | 4 quarters (12 min) + OT |
| Score display | `2-1` | `105-98` |
| Live time | `45+2'` | `Q3 2:34` |
| Events | Goals, cards, subs | Field goals, fouls, timeouts |
| Competitions | 50+ leagues worldwide | NBA (East/West) + playoffs |
| Highlights | r/soccer | r/nba |

## Installation

### Build from source

```bash
git clone https://github.com/gabriel7419/courtside.git
cd courtside

# Build it
make build

# Run it
./courtside
```

Alternatively, you can install it globally via `make install` or `go install ./cmd/...`.

## Usage

```bash
courtside

# Or use the make command
make run

# Run with mock data (useful during the off-season or when no games are live)
courtside --mock
# or: make mock
```

**Navigation:** `↑`/`↓` or `j`/`k` to move, `Enter` to select, `/` to filter, `Tab` to switch pane, `Esc` to go back, `q` to quit.

**Views:**
- **Today's games** — live and upcoming games
- **Finished games** — recent results (last 3 or 5 days)
- **Settings** — filter by conference, toggle notifications

## Docs

- [Quick Start](QUICKSTART.md) — set up your development environment
- [Supported Teams](docs/SUPPORTED_TEAMS.md) — all 30 NBA teams by conference and division
- [Notifications](docs/NOTIFICATIONS.md) — desktop notification setup
- [API Reference](docs/API_REFERENCE.md) — NBA Stats API endpoints and response format
- [Implementation Plan](docs/FORK_PLAN.md) — full roadmap across 7 phases

## Roadmap

- [x] Project planning and documentation
- [x] NBA API integration (`internal/nba/`)
- [x] List today's games
- [x] Game details view
- [x] Live game polling
- [x] Statistics dialog
- [ ] r/nba highlights integration
- [ ] Desktop notifications
- [x] Playoff bracket view
- [ ] WNBA support

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for workflow, code guidelines, and what currently needs help.

## Credits

Courtside is a fork/adaptation of [Golazo](https://github.com/0xjuanma/golazo) by [@0xjuanma](https://github.com/0xjuanma).

Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea), [Lip Gloss](https://github.com/charmbracelet/lipgloss) and [Bubbles](https://github.com/charmbracelet/bubbles) by [Charm](https://charm.sh).

## License

MIT — see [LICENSE](LICENSE) for details.

---

<div align="center">

*Made with ❤️ for basketball and terminal enthusiasts*

</div>
