# NBA Stats API — Referência

Referência rápida para os endpoints usados no Courtside. Para contexto geral, veja [FORK_PLAN.md](FORK_PLAN.md).

## Headers Obrigatórios

Todas as requisições precisam destes headers — sem eles a API retorna 403:

```http
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36
Referer: https://www.nba.com/
Accept: application/json
Accept-Language: en-US,en;q=0.9
Origin: https://www.nba.com
```

---

## Endpoints

### 1. Scoreboard — Jogos do Dia

```
GET https://stats.nba.com/stats/scoreboard?GameDate=YYYY-MM-DD&LeagueID=00
```

**Parâmetros:**
- `GameDate` — data no formato `YYYY-MM-DD`
- `LeagueID` — sempre `00` para NBA

**Resposta (simplificada):**

```json
{
  "resultSets": [
    {
      "name": "GameHeader",
      "headers": [
        "GAME_DATE_EST", "GAME_SEQUENCE", "GAME_ID",
        "GAME_STATUS_ID", "GAME_STATUS_TEXT", "GAMECODE",
        "HOME_TEAM_ID", "VISITOR_TEAM_ID", "SEASON",
        "LIVE_PERIOD", "LIVE_PC_TIME", "NATL_TV_BROADCASTER_ABBREVIATION",
        "LIVE_PERIOD_TIME_BCAST", "WH_STATUS"
      ],
      "rowSet": [
        ["2026-02-25T00:00:00", 1, "0022300789", 2, "Live",
         "20260225/LALGWS", 1610612744, 1610612747, "2023",
         3, "PT2M34.00S", "TNT", "Q3 2:34", 1]
      ]
    },
    {
      "name": "LineScore",
      "headers": [
        "GAME_DATE_EST", "GAME_ID", "TEAM_ID", "TEAM_ABBREVIATION",
        "TEAM_CITY_NAME", "TEAM_WINS_LOSSES",
        "PTS_QTR1", "PTS_QTR2", "PTS_QTR3", "PTS_QTR4",
        "PTS_OT1", "PTS", "FG_PCT", "FT_PCT", "FG3_PCT",
        "AST", "REB", "TOV"
      ],
      "rowSet": [
        ["2026-02-25T00:00:00", "0022300789", 1610612747,
         "LAL", "Los Angeles", "35-20",
         28, 24, 22, null, null, 74,
         0.488, 0.800, 0.400, 18, 32, 8]
      ]
    }
  ]
}
```

**Campos relevantes do GameHeader:**

| Campo | Tipo | Descrição |
|---|---|---|
| `GAME_ID` | string | ID único, ex: `"0022300789"` |
| `GAME_STATUS_ID` | int | `1` = agendado, `2` = ao vivo, `3` = finalizado |
| `GAME_STATUS_TEXT` | string | `"Live"`, `"Final"`, `"9:00 PM ET"` |
| `LIVE_PERIOD` | int | Quarter atual (1–4, 5+ = OT) |
| `LIVE_PC_TIME` | string | Tempo restante no formato ISO, ex: `"PT2M34.00S"` |
| `LIVE_PERIOD_TIME_BCAST` | string | Formato broadcast, ex: `"Q3 2:34"` |

**Campos relevantes do LineScore:**

| Campo | Descrição |
|---|---|
| `PTS_QTR1..4` | Pontos por quarter (null se ainda não jogado) |
| `PTS` | Total de pontos |
| `FG_PCT` | % Field Goals (0.488 = 48.8%) |
| `FG3_PCT` | % 3-Pointers |
| `FT_PCT` | % Free Throws |
| `AST` | Assistências |
| `REB` | Rebotes totais |
| `TOV` | Turnovers |

---

### 2. Box Score Summary — Resumo do Jogo

```
GET https://stats.nba.com/stats/boxscoresummaryv2?GameID=0022300789
```

Retorna `resultSets` com: `GameSummary`, `LineScore`, `LastMeeting`, `SeasonSeries`, `Officials`, `GameInfo`.

**Campos úteis do GameInfo:**

```json
{
  "name": "GameInfo",
  "headers": ["GAME_DATE", "ATTENDANCE", "GAME_TIME"],
  "rowSet": [["2026-02-25T00:00:00", "18064", "2:23"]]
}
```

---

### 3. Box Score Traditional — Estatísticas Completas

```
GET https://stats.nba.com/stats/boxscoretraditionalv2?GameID=0022300789&StartPeriod=1&EndPeriod=10&StartRange=0&EndRange=28800&RangeType=0
```

Retorna `PlayerStats` e `TeamStats`.

**Exemplo — um jogador:**

```json
{
  "name": "PlayerStats",
  "headers": [
    "GAME_ID", "TEAM_ID", "TEAM_ABBREVIATION", "TEAM_CITY",
    "PLAYER_ID", "PLAYER_NAME", "START_POSITION", "COMMENT",
    "MIN", "FGM", "FGA", "FG_PCT", "FG3M", "FG3A", "FG3_PCT",
    "FTM", "FTA", "FT_PCT", "OREB", "DREB", "REB",
    "AST", "STL", "BLK", "TO", "PF", "PTS", "PLUS_MINUS"
  ],
  "rowSet": [
    ["0022300789", 1610612747, "LAL", "Los Angeles",
     2544, "LeBron James", "F", "",
     "36:24", 12, 20, 0.600, 3, 7, 0.429,
     5, 6, 0.833, 1, 8, 9, 7, 2, 1, 3, 2, 32, 5]
  ]
}
```

**Siglas:**

| Sigla | Significado |
|---|---|
| `MIN` | Minutos jogados (`"MM:SS"`) |
| `FGM` / `FGA` | Field Goals Made / Attempted |
| `FG_PCT` | % Field Goals |
| `FG3M` / `FG3A` | 3-Pointers Made / Attempted |
| `FG3_PCT` | % 3-Pointers |
| `FTM` / `FTA` | Free Throws Made / Attempted |
| `FT_PCT` | % Free Throws |
| `OREB` / `DREB` | Rebotes Ofensivos / Defensivos |
| `REB` | Rebotes totais |
| `AST` | Assistências |
| `STL` | Steals (roubos de bola) |
| `BLK` | Blocks (tocos) |
| `TO` | Turnovers |
| `PF` | Faltas pessoais |
| `PTS` | Pontos |
| `PLUS_MINUS` | Diferença de pontos quando em quadra |

---

### 4. Play-by-Play — Eventos ao Vivo

```
GET https://stats.nba.com/stats/playbyplayv2?GameID=0022300789&StartPeriod=1&EndPeriod=10
```

**Exemplo de resposta:**

```json
{
  "resultSets": [{
    "name": "PlayByPlay",
    "headers": [
      "GAME_ID", "EVENTNUM", "EVENTMSGTYPE", "EVENTMSGACTIONTYPE",
      "PERIOD", "WCTIMESTRING", "PCTIMESTRING",
      "HOMEDESCRIPTION", "NEUTRALDESCRIPTION", "VISITORDESCRIPTION",
      "SCORE", "SCOREMARGIN"
    ],
    "rowSet": [
      ["0022300789", 2, 2, 1, 1, "7:00 PM", "11:42",
       null, null, "Curry 25' 3PT Jump Shot (3 PTS)", "0-3", "-3"],
      ["0022300789", 3, 1, 1, 1, "7:00 PM", "11:18",
       "James 2' Driving Layup (2 PTS) (Russell assists)", null, null, "2-3", "-1"],
      ["0022300789", 4, 6, 1, 1, "7:01 PM", "10:55",
       null, null, "Curry Personal Foul (1 PF)", null, null],
      ["0022300789", 5, 9, 0, 1, "7:02 PM", "9:30",
       null, "Warriors Timeout: Regular", null, null, null]
    ]
  }]
}
```

**EVENTMSGTYPE:**

| Valor | Evento |
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

---

## Boas Práticas

**Rate limiting:** A API não documenta limites, mas 200–300ms entre requests é razoável.

**Cache:**
```
Jogos ao vivo:      5–10s
Scoreboard geral:   30s
Jogos finalizados:  permanente (nunca mudam)
Estatísticas:       1min (jogo ao vivo)
```

---

## Game ID Format

```
0022300789
│││││└───── Número sequencial
││││└────── Tipo: 2 = Regular Season, 4 = Playoffs
│││└─────── Temporada: 23 = 2023-24
││└──────── Século (0 = 2000s)
│└───────── Sempre 0
└────────── Sempre 0
```

## Team IDs Principais

| Time | ID |
|---|---|
| Lakers | 1610612747 |
| Warriors | 1610612744 |
| Celtics | 1610612738 |
| Heat | 1610612748 |
| Bulls | 1610612741 |
| Knicks | 1610612752 |
| 76ers | 1610612755 |

Lista completa: [swar/nba_api](https://github.com/swar/nba_api/blob/master/src/nba_api/stats/library/data.py)

---

## Testando

```bash
# Script de teste incluído
go run scripts/test_nba_api.go --endpoint=scoreboard --date=2026-02-25
go run scripts/test_nba_api.go --endpoint=summary --game=0022300789
go run scripts/test_nba_api.go --endpoint=traditional --game=0022300789
go run scripts/test_nba_api.go --endpoint=playbyplay --game=0022300789

# curl direto
curl -H "User-Agent: Mozilla/5.0" \
     -H "Referer: https://www.nba.com/" \
     "https://stats.nba.com/stats/scoreboard?GameDate=2026-02-25&LeagueID=00"
```
