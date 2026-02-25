# Quick Start

This guide gets you from zero to a running development environment.

## Prerequisites

- Go 1.21+
- Git
- A terminal with 256-color support

## 1. Set up the repo

```bash
git clone https://github.com/gabriel7419/courtside.git
cd courtside
go mod download
go build
```

## 2. Test the NBA API

Before writing any code, make sure the API is accessible:

```bash
# Using the included test script
go run scripts/test_nba_api.go --endpoint=scoreboard --date=2026-02-25

# Or with curl
curl -H "User-Agent: Mozilla/5.0" \
     -H "Referer: https://www.nba.com/" \
     "https://stats.nba.com/stats/scoreboard?GameDate=2026-02-25&LeagueID=00"
```

If you get a valid JSON response, you're good. If you get 403, double-check the `User-Agent` and `Referer` headers.

## 3. Rename the Go module

Edit `go.mod` to point to your repository:

```
module github.com/gabriel7419/courtside
```

Then update all imports:

```powershell
# Windows PowerShell
Get-ChildItem -Recurse -Include *.go | ForEach-Object {
    (Get-Content $_.FullName) -replace 'github.com/0xjuanma/golazo', 'github.com/gabriel7419/courtside' | Set-Content $_.FullName
}
```

```bash
# Linux/macOS
find . -name "*.go" -exec sed -i 's|github.com/0xjuanma/golazo|github.com/gabriel7419/courtside|g' {} +
```

Run `go mod tidy && go build` to confirm everything compiles.

## 4. Create a branch and start coding

```bash
git checkout -b feature/nba-adaptation
```

The first concrete task is creating `internal/nba/client.go`. Use `internal/fotmob/client.go` as a reference — the structure is nearly identical, just pointing at a different API.

See [docs/FORK_PLAN.md](docs/FORK_PLAN.md) for the full implementation roadmap.

## Troubleshooting

**403 on the API** — Check that `User-Agent` and `Referer` headers are set. See [docs/API_REFERENCE.md](docs/API_REFERENCE.md).

**Module not found** — Run `go mod tidy`.

**Build fails after renaming** — Search for leftover old imports:
```bash
grep -r "github.com/0xjuanma/golazo" .
```
