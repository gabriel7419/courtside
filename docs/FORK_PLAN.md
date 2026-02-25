# NBA Fork — Plano de Implementação

Documento de referência para adaptar o Golazo (futebol) para NBA. Consulte o [FORK_SUMMARY.md](FORK_SUMMARY.md) para uma visão mais curta.

---

## Arquitetura Atual (Golazo)

```
golazo/
├── cmd/                    # CLI (Cobra)
├── internal/
│   ├── api/               # Interface sport-agnostic
│   ├── fotmob/            # Cliente FotMob API
│   ├── reddit/            # Busca de highlights
│   ├── ui/                # Interface TUI (Bubble Tea)
│   ├── data/              # Settings e storage
│   ├── notify/            # Notificações desktop
│   ├── app/               # Lógica principal
│   ├── constants/
│   ├── debug/
│   └── version/
├── assets/
├── scripts/
└── docs/
```

Pontos positivos da arquitetura: interface `api.Client` bem abstraída, cache com TTL, rate limiting configurável, UI desacoplada da lógica de dados.

---

## Estrutura Alvo (Courtside)

```
courtside/
├── cmd/                    # Manter, só renomear comandos
├── internal/
│   ├── api/               # Adaptar tipos para NBA (Game, Quarter, etc.)
│   ├── nba/               # NOVO — cliente NBA Stats API
│   ├── reddit/            # Adaptar (r/nba em vez de r/soccer)
│   ├── ui/                # Adaptar labels (Quarter, Timeout, Clock)
│   ├── data/              # Adaptar (conferências e divisões)
│   ├── notify/            # Manter sem alterações
│   ├── app/               # Adaptar lógica para NBA
│   ├── constants/         # Adaptar
│   ├── debug/             # Manter
│   └── version/           # Manter
├── assets/
├── scripts/
└── docs/
```

---

## API

### Opção recomendada: NBA Stats API

- **Base URL:** `https://stats.nba.com/stats/`
- **Custo:** Gratuita
- **Autenticação:** Nenhuma (mas requer headers específicos)
- **Rate limiting:** Não documentado — usar 200ms entre requests

**Endpoints principais:**

```
GET /scoreboard?GameDate=YYYY-MM-DD&LeagueID=00       → jogos do dia
GET /boxscoresummaryv2?GameID=<id>                     → resumo do jogo
GET /boxscoretraditionalv2?GameID=<id>                 → estatísticas completas
GET /playbyplayv2?GameID=<id>&StartPeriod=1&EndPeriod=10 → play-by-play
```

Ver detalhes em [API_REFERENCE.md](API_REFERENCE.md).

### Alternativa: ESPN API (não oficial)

`https://site.api.espn.com/apis/site/v2/sports/basketball/nba/`

JSON mais limpo, mas pode mudar sem aviso.

---

## Mapeamento: Futebol → NBA

### Estrutura de Dados

| Futebol (Golazo) | NBA (Courtside) |
|---|---|
| `Match` | `Game` |
| `League` | `Conference` / `Division` |
| `Round` | `GameNumber` / `Date` |
| `LiveTime` (`"45+2"`) | `Quarter` + `Clock` (`"Q3 2:34"`) |
| `MatchStatus` | `GameStatus` |

### Eventos

| Futebol | NBA |
|---|---|
| Goal | Field Goal (2pt/3pt), Free Throw |
| Yellow Card | Personal Foul / Technical Foul |
| Red Card | Ejection / Flagrant Foul |
| Substitution | Substitution |
| N/A | Timeout |

### Estatísticas

| Futebol | NBA |
|---|---|
| Possession | Time of Possession |
| Shots | FGA (Field Goal Attempts) |
| Shots on Target | FGM (Field Goals Made) |
| Passes | Assists |
| N/A | Rebounds (OFF/DEF) |
| N/A | Steals, Blocks, Turnovers |
| N/A | FG%, 3PT%, FT% |

---

## Adaptações de Código

### `internal/api/types.go`

```go
// Antes (futebol)
type Match struct {
    ID        int
    League    League
    HomeTeam  Team
    AwayTeam  Team
    Status    MatchStatus
    HomeScore *int
    AwayScore *int
    MatchTime *time.Time
    LiveTime  *string   // "45+2", "HT", "FT"
    Round     string
}

// Depois (NBA)
type Game struct {
    ID           int
    Conference   string
    HomeTeam     Team
    AwayTeam     Team
    Status       GameStatus
    HomeScore    *int
    AwayScore    *int
    GameTime     *time.Time
    Quarter      *int      // 1-4, 5+ para OT
    Clock        *string   // "2:34"
    IsPlayoffs   bool
    SeriesStatus *string   // "Series tied 2-2"
}
```

### `internal/data/settings.go`

```go
// Antes (50+ ligas de futebol)
var AllSupportedLeagues = map[string][]LeagueInfo{
    RegionEurope:  {{ID: 47, Name: "Premier League"}, ...},
    RegionAmerica: {{ID: 268, Name: "Brasileirão"}, ...},
}

// Depois (conferências NBA)
const (
    ConferenceEastern = "Eastern"
    ConferenceWestern = "Western"
)

var AllSupportedConferences = map[string][]ConferenceInfo{
    ConferenceEastern: {
        {Division: "Atlantic", Teams: []string{"Celtics", "Nets", "Knicks", "76ers", "Raptors"}},
        {Division: "Central", Teams: []string{"Bulls", "Cavaliers", "Pistons", "Pacers", "Bucks"}},
        {Division: "Southeast", Teams: []string{"Hawks", "Hornets", "Heat", "Magic", "Wizards"}},
    },
    ConferenceWestern: {
        {Division: "Northwest", Teams: []string{"Nuggets", "Timberwolves", "Thunder", "Blazers", "Jazz"}},
        {Division: "Pacific", Teams: []string{"Warriors", "Clippers", "Lakers", "Suns", "Kings"}},
        {Division: "Southwest", Teams: []string{"Mavericks", "Rockets", "Grizzlies", "Pelicans", "Spurs"}},
    },
}
```

### Reddit Integration

```
Antes: subreddit r/soccer, keywords: "goal", "GOAL"
Depois: subreddit r/nba, keywords: "highlight", "HIGHLIGHT", "clutch"
```

---

## Checklist de Implementação

### Fase 0: Preparação (completo)
- [x] Plano de implementação
- [x] README e CONTRIBUTING adaptados
- [x] Documentação da API
- [x] Script de teste da API

### Fase 1: Setup Inicial
- [ ] Renomear módulo em `go.mod`
- [ ] Atualizar imports em todos os arquivos `.go`
- [ ] Confirmar `go build` sem erros

### Fase 2: Camada de Dados
- [ ] Adaptar `internal/api/types.go` (Game, Quarter, Clock)
- [ ] Criar estrutura `internal/nba/`
- [ ] Implementar `internal/nba/client.go`
- [ ] Implementar `internal/nba/types.go` (mapeamento API → tipos Go)
- [ ] Adaptar `internal/data/settings.go` (conferências NBA)

### Fase 3: Funcionalidade Básica
- [ ] Implementar `GamesByDate()` — jogos do dia
- [ ] Implementar `GameDetails()` — detalhes de um jogo
- [ ] Testar cache e rate limiting
- [ ] Testes unitários básicos

### Fase 4: Dados ao Vivo
- [ ] Polling para jogos ao vivo
- [ ] Mapear eventos (field goals, faltas, timeouts)
- [ ] Atualização de placar em tempo real
- [ ] Testar com jogos reais

### Fase 5: UI
- [ ] Atualizar labels (Quarter, Clock, etc.)
- [ ] Adaptar formatação de placar (scores altos: 98-105)
- [ ] Atualizar dialogs de estatísticas e standings

### Fase 6: Features Extras
- [ ] Highlights via r/nba
- [ ] Notificações adaptadas para NBA
- [ ] Tabela de playoffs
- [ ] Filtros por conferência

### Fase 7: Release
- [ ] README final com screenshots
- [ ] Scripts de instalação adaptados
- [ ] Release v1.0.0

---

## Testando a API Manualmente

```bash
# Jogos de hoje
curl -H "User-Agent: Mozilla/5.0" \
     -H "Referer: https://www.nba.com/" \
     "https://stats.nba.com/stats/scoreboard?GameDate=2026-02-25&LeagueID=00"

# Detalhes de um jogo específico
curl -H "User-Agent: Mozilla/5.0" \
     -H "Referer: https://www.nba.com/" \
     "https://stats.nba.com/stats/boxscoresummaryv2?GameID=0022300789"
```

Ou use o script incluído:

```bash
go run scripts/test_nba_api.go --endpoint=scoreboard --date=2026-02-25
```

---

## Recursos

- [swar/nba_api](https://github.com/swar/nba_api) — referência Python muito completa
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — framework TUI
- [Golazo original](https://github.com/0xjuanma/golazo) — base deste projeto

---

*Baseado em [Golazo](https://github.com/0xjuanma/golazo) por [@0xjuanma](https://github.com/0xjuanma) — Licença MIT*
