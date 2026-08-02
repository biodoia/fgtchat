# AGENTS.md — fgtchat

Regole operative per coding agent su questo repository.
Estende il workspace `/home/lisergico25/projects/AGENTS.md`.

## Standard ecosistema (inviolabili)

1. **No Tailscale** — non aggiungere `tailscaled`, `tsnet`, `*.ts.net`, auth tailnet.
2. **Rete** — bind `127.0.0.1:<porta>`; esporre solo via **aigoproxy** (`*.braigo.dev`).
3. **DB** — PostgreSQL/pgvector tramite **memogo**; secrets in **goleciave**.
4. **Cache** — **PebbleDB** pure Go; zero CGO; no SQLite autoritativo.
5. **LLM** — **gogatewai** (non LiteLLM legacy).
6. **Memoria** — mem0 locale `user_id=biodoia` (search_memories / add_memory).
7. **Shape** — fgt-sdk v3.0: 1 Body systemd + N Heads; single binary dove possibile.
8. **Commit** — granulari conventional (`feat`/`fix`/`docs`/`refactor`/`test`).

## Prima di modificare codice

1. Leggi `SPEC.md` e questo file.
2. Verifica dipendenze in `go.mod` (se presente).
3. Consulta catalogo **nocodaigo** se serve contesto snippato.
4. Usa **godocai** per scan/generate docs dopo cambi architetturali.

## Progetto

- **Path**: `/home/lisergico25/projects/fgtchat`
- **Modulo**: `github.com/biodoia/fgtchat`
- **Descrizione**: Repository: `github.com/biodoia/fgtchat` · Go `1.21`.

## Documentazione

Dopo cambi rilevanti aggiorna README / SPEC e riallinea con:

```bash
/home/lisergico25/projects/nocodaigo/nocodaigo explode /home/lisergico25/projects/fgtchat --no-llm
godocai scan -d /home/lisergico25/projects/fgtchat
```

---
*aligned 2026-08-02*
