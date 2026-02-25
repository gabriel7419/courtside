# NBA Fork — Visão Geral

Este documento resume como o Courtside se relaciona com o Golazo e quais as diferenças chave entre os dois projetos.

## O que é o Courtside?

Courtside é um fork do [Golazo](https://github.com/0xjuanma/golazo) adaptado para acompanhar jogos da NBA no terminal. A base de código é essencialmente a mesma — a diferença está na camada de dados (`internal/nba/` em vez de `internal/fotmob/`) e em alguns ajustes de UI (labels, formatação de placar, etc.).

A maior parte do trabalho é implementar o cliente da NBA Stats API e mapear as respostas para os tipos internos do projeto.

## Diferenças Chave

| Componente | Golazo | Courtside |
|---|---|---|
| API | FotMob | NBA Stats API |
| Pacote do cliente | `internal/fotmob/` | `internal/nba/` |
| Struct principal | `Match` | `Game` |
| Exibição de tempo | `"45'"` | `"Q3 2:34"` |
| Placar | Baixo (0–5) | Alto (90–120) |
| Competições | 50+ ligas mundiais | Eastern/Western + playoffs |
| Highlights | r/soccer | r/nba |

## Vantagens da Arquitetura

O Golazo foi projetado com boa separação de responsabilidades:

- `internal/api/` define uma interface agnóstica de esporte — troca-se o cliente sem mexer na UI
- O sistema de cache com TTL já existe e funciona
- O rate limiting já é configurável
- Cerca de 80% da UI pode ser reaproveitada (só mudar labels e formatações)

## Desafios Esperados

- A NBA Stats API não é documentada oficialmente e pode mudar
- As respostas são em formato tabular (`headers[]` + `rowSet[][]`), o que exige mapeamento manual dos índices
- A NBA tem mais tipos de eventos que futebol (field goals, fouls, timeouts, jump balls...)
- Scores altos (ex: 98-105) precisam de ajuste no layout da UI

## O que já está pronto

- Documentação completa (este arquivo, FORK_PLAN, API_REFERENCE)
- Script de teste da API (`scripts/test_nba_api.go`)
- Mapeamento de conceitos de futebol → NBA
- Estrutura proposta para `internal/nba/`

## Próximo passo

Implementar `internal/nba/client.go`, começando pelo método `GamesByDate()`. O arquivo `internal/fotmob/client.go` é a melhor referência de como estruturar isso.

Ver [FORK_PLAN.md](FORK_PLAN.md) para o roadmap completo.
