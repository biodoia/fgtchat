---
feature: fgtchat
status: docs-aligned
updated: 2026-08-02
standard: fgt-sdk-v3.0
---

# SPEC — fgtchat

## S1: Ruolo

Repository: `github.com/biodoia/fgtchat` · Go `1.21`.

**Modulo**: `github.com/biodoia/fgtchat`

## S2: Design

### Architettura

- **Pattern**: fgt-sdk v3.0 — 1 Body (systemd) + N Heads (TUI / Web / Chat / MCP) dove applicabile.
- **Pacchetti osservati**: n/d

### Networking

- Listener su **`127.0.0.1:<porta>`** (mai `0.0.0.0` per servizi app).
- Hostname pubblici: **`fgtchat.braigo.dev`** (o sottodomini dedicati) tramite **aigoproxy**.
- Registrazione rotte: `POST http://127.0.0.1:80/api/routes`.
- **Tailscale / tsnet / `*.ts.net`**: vietati (direttiva 2026-08-01).

### Storage & secrets

- Persistenza relazionale/vettoriale: PostgreSQL + pgvector, database **provisionato da memogo** (DSN in goleciave, driver nativo). memogo non è mai nel data path.
- Credenziali: **goleciave** (mai secret in chiaro in repo).
- Cache locale: **PebbleDB** (pure Go).
- Vietato: SQLite/BoltDB/BadgerDB come store autoritativo; CGO non abilitato per binari di produzione.

### LLM & memoria agenti

- Router model-first: **gogatewai**.
- Memoria condivisa: **mem0** MCP `http://127.0.0.1:12000`, `user_id=biodoia`.

### Stack rilevato dal codice

- **Rete**: listener `127.0.0.1:<porta>` · host pubblici `*.braigo.dev` via **aigoproxy**
- **Memoria agenti**: mem0 MCP `http://127.0.0.1:12000` (`user_id=biodoia`)
- **Tailscale**: **rimosso** dall'ecosistema (2026-08-01) — vietato in docs e codice nuovo

## S3: Out of scope

- Non reintrodurre Tailscale come requisito di rete.
- Non documentare bind pubblici diretti senza aigoproxy.

## Verifica

Allineato dopo ispezione di `go.mod` / albero sorgente, `nocodaigo explode --no-llm`,
e `godocai scan`.

---
*docs(align): 2026-08-02*
