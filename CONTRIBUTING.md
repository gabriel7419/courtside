# Contributing to Courtside

Courtside is a fork of [Golazo](https://github.com/0xjuanma/golazo) adapted for the NBA. The project is early-stage — there's a lot of ground to cover and all contributions are welcome.

## How to Contribute

1. Check if there's already an [issue](https://github.com/gabriel7419/courtside/issues) for what you want to work on — if not, create one before starting
2. Fork the repo and create a feature branch
3. Make your changes and run `go build` to confirm it compiles
4. Submit a pull request referencing the issue (e.g., `Fixes #12`)

Every PR must reference an existing issue.

## Development Workflow

```bash
git clone https://github.com/gabriel7419/courtside.git
cd courtside
git checkout -b feature/your-feature-name

# After making changes:
go build
./courtside

git commit -m "feat: add quarter-by-quarter scoring display"
git push origin feature/your-feature-name
```

## What Needs Help

**Foundational (highest priority):**

- NBA API integration — implement `internal/nba/client.go`
- Data mapping — convert NBA API responses to internal types
- Basic UI labels — update terminology (Quarter, Timeout, Clock, etc.)

**Medium priority:**

- Statistics parsing — parse NBA box scores
- Live game polling — real-time score updates
- Event timeline — display play-by-play events
- Conference filter in Settings

**Nice to have:**

- r/nba highlights integration
- Notifications for key plays (3-pointers, dunks)
- Playoff bracket view
- Player stats dialog

## Code Style

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Run `gofmt` before committing
- Add comments for exported functions
- Keep functions small and focused

**Example:**

```go
// GamesByDate retrieves all NBA games scheduled for the given date.
func (c *Client) GamesByDate(ctx context.Context, date time.Time) ([]api.Game, error) {
    dateStr := date.Format("2006-01-02")
    // ...
    return games, nil
}
```

## NBA API Notes

The NBA Stats API is not officially documented but is publicly accessible. Every request needs these headers or you'll get a 403:

```go
req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
req.Header.Set("Referer", "https://www.nba.com/")
req.Header.Set("Accept", "application/json")
```

Keep at least 200ms between requests and cache aggressively — finished games never change, so there's no reason to re-fetch them.

See [docs/API_REFERENCE.md](docs/API_REFERENCE.md) for full endpoint documentation.

## Key Differences from Golazo

| Component | Golazo | Courtside |
|---|---|---|
| API client | `internal/fotmob/` | `internal/nba/` |
| Main struct | `Match` | `Game` |
| Time display | `"45'"` | `"Q3 2:34"` |
| Score range | 0–5 typical | 90–120 typical |
| Leagues | 50+ worldwide | Eastern/Western conferences |
| Highlights | r/soccer | r/nba |

## Useful Scripts

```bash
# Test NBA API endpoints
go run scripts/test_nba_api.go --endpoint=scoreboard --date=2026-02-25

# Test a specific game
go run scripts/test_nba_api.go --endpoint=summary --game=0022300789

# Clear cache
go run scripts/clear_cache.go
```

## Getting Help

- Questions → [GitHub Discussions](https://github.com/gabriel7419/courtside/discussions)
- Bugs → [GitHub Issues](https://github.com/gabriel7419/courtside/issues)
- Read the plan → [docs/FORK_PLAN.md](docs/FORK_PLAN.md)

Thank you for contributing!
