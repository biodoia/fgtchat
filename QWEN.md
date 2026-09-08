# QWEN.md — fgtchat

Vincoli locali per agent (Qwen / Codex / Grok / Claude) su questo repo.

## REGOLA #1 — Ecosistema Biodoia 2026-08-01

- Vietato Tailscale (`tailscaled`, `tsnet`, `*.ts.net`).
- Rete = **aigoproxy** + bind `127.0.0.1`.
- DB = PostgreSQL+pgvector **provisionato da memogo** (memogo installa il DB e mette il DSN in goleciave; NON è un datastore, mai nel data path) + secrets **goleciave**.
- Cache = **PebbleDB**. CGO_ENABLED=0.
- LLM = **gogatewai**. Memoria agent = **mem0** (`biodoia`).

## REGOLA #2 — Verità = codice

Non fidarti di docs vecchie se contraddicono `go.mod` / `internal/`.
Se la doc parla di Tailscale e il codice no: correggi la doc, non il codice verso Tailscale.

## REGOLA #3 — Docs tools

- **nocodaigo**: `explode <path> --no-llm` per snippet strutturali.
- **godocai**: `scan` / `generate` con `NINEROUTER_URL` + `GODOCAI_MODEL=top` se LLM ok.

## Quick facts

| Campo | Valore |
|-------|--------|
| Name | fgtchat |
| Module | github.com/biodoia/fgtchat |
| Go | 1.21 |
| Pebble | False |
| pgx | False |
| fgt-sdk | False |

---
*aligned 2026-08-02*
