# Courtside — Documentation Index

Navigation guide for the project docs.

## Getting Started

| Document | Description |
|---|---|
| [QUICKSTART.md](QUICKSTART.md) | Set up your dev environment and run your first build |
| [CONTRIBUTING.md](CONTRIBUTING.md) | Workflow, code style, and what needs help |
| [README.md](README.md) | Project overview and feature list |

## Technical Reference

| Document | Description |
|---|---|
| [docs/API_REFERENCE.md](docs/API_REFERENCE.md) | NBA Stats API endpoints, headers, response format |
| [docs/FORK_PLAN.md](docs/FORK_PLAN.md) | Full implementation roadmap (7 phases), data mapping, architecture |
| [docs/FORK_SUMMARY.md](docs/FORK_SUMMARY.md) | High-level overview of the Golazo → Courtside adaptation |

## Reference

| Document | Description |
|---|---|
| [docs/SUPPORTED_TEAMS.md](docs/SUPPORTED_TEAMS.md) | All 30 NBA teams by conference and division |
| [docs/NOTIFICATIONS.md](docs/NOTIFICATIONS.md) | Desktop notification setup by OS |

## Development Tools

```bash
# Test NBA API endpoints
go run scripts/test_nba_api.go --endpoint=scoreboard --date=2026-02-25
go run scripts/test_nba_api.go --endpoint=summary --game=0022300789
go run scripts/test_nba_api.go --endpoint=traditional --game=0022300789
go run scripts/test_nba_api.go --endpoint=playbyplay --game=0022300789

# Clear local cache
go run scripts/clear_cache.go
```

## Suggested Reading Order

**Before writing code:** QUICKSTART → docs/API_REFERENCE → docs/FORK_PLAN

**To contribute:** CONTRIBUTING → docs/FORK_PLAN (pick a phase)

**To understand the architecture:** docs/FORK_SUMMARY → `internal/fotmob/` (the reference implementation)
