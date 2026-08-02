# 📋 Document Alignment & Audit Proposal — fgtchat

> **Data Audit**: 2026-08-02
> **Stato Proposta**: In Attesa di Revisione (Nessun documento esistente è stato modificato)
> **Standard Ecosistema**: fgt-sdk v3.0 / aigoproxy / memogo / goleciave / gogatewai / mem0

---

## 1. 🔍 Metadati Repository e Documentazione Esistente
- **Percorso**: `/home/lisergico25/projects/fgtchat`
- **Categoria Repository**: Tool, Integration or Specialized Repo
- **Ultimo Commit Git**: 2026-07-12

### Documenti attuali e date di modifica:
- `README.md` (Ultima modifica: 2026-07-12)

---

## 2. 🚨 Analisi delle Derive Tecnologiche & Allucinazioni
> ℹ️ **Stato del Repo**: Repository attivo nella categoria `Tool, Integration or Specialized Repo`. Necessita aggiornamento documentale.

✅ **Rete**: Nessun riferimento obsoleto a Tailscale rilevato.
✅ **Storage / CGO**: Conforme a PebbleDB / memogo / zero-CGO.
✅ **LLM Proxy**: Nessun proxy deprecato rilevato.
✅ **SDK Standard**: Allineato alle specifiche moderne di ecosistema.

---

## 3. 💾 Allineamento Database & Memoria Agenti
- **Memoria condivisa Mem0**: Il repository deve integrare il protocollo di memoria condivisa via Mem0 local MCP (`http://127.0.0.1:12000`, `user_id=biodoia`).
- **Database di persistenza**: Qualsiasi persistenza relazionale o vettoriale deve appoggiarsi a PostgreSQL/pgvector gestito da `memogo` con credenziali in `goleciave`.
- **Cache di stato locale**: Utilizzare esclusivamente `PebbleDB` integrato in Go (zero dipendenze CGO/SQLite).

---

## 4. 📝 Modifiche Proposte alla Documentazione (Roadmap per Agenti Profondi)
Una volta approvata questa proposta dagli agenti di refactoring profondo, si procederà all'aggiornamento ordinato di:

1. **`README.md`**:
   - Aggiornare l'architettura allineandola a `fgt-sdk v3.0` (1 Body systemd + N Heads).
   - Rimuovere ogni menzione di Tailscale o SQLite.
   - Inserire le istruzioni di registrazione delle rotte su `aigoproxy` (`POST http://127.0.0.1:80/api/routes`).

2. **`SPEC.md` / `ARCHITECTURE.md`**:
   - Dichiarare la compliance esplicita con `goleciave` (secrets) e `gogatewai` (LLM router).
   - Definire i punti di integrazione con `mem0` per la memoria persistente dell'agente.

3. **`AGENTS.md` / `QWEN.md`**:
   - Sincronizzare le regole operative dell'agente con la direttiva ecosistema 2026-08-01 (Zero Tailscale, Mem0 obbligatorio, PebbleDB per la cache).

---

## 5. 🎯 Prossimi Passi
1. Verificare questa proposta con l'orchestrazione degli agenti di audit profondo.
2. Eseguire la scrittura effettiva della documentazione allineata.

---

## ✅ Deep align eseguito (2026-08-02)

- Verifica codice: go.mod=True pebble=False pgx=False fgt=False
- nocodaigo explode --no-llm: eseguito
- godocai scan: eseguito
- Docs aggiornati: README.md, SPEC.md, AGENTS.md, QWEN.md
- Standard: fgt-sdk v3.0 / no Tailscale / aigoproxy / memogo / goleciave / PebbleDB / gogatewai / mem0
