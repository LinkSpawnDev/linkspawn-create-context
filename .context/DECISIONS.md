<!-- scaffolded by spawn dev (https://github.com/LinkSpawnDev/linkspawn-create-context) on 2026-06-11 -->
# DECISIONS.md — linkspawn-create-context

Append-only log of locked decisions. Don't rewrite history — supersede with a new entry that references the old one. Agents: do not relitigate a decision recorded here; if you believe one is wrong, say so explicitly and wait for the owner to re-open it.

Entry format:

```markdown
## DR-NNN — <title> (YYYY-MM-DD)

**Status:** Accepted | Superseded by DR-MMM
**Context:** the situation and the options that were on the table.
**Decision:** what was chosen.
**Trade-offs accepted:** what this costs us, stated honestly.
```

---

## DR-001 — Adopt the `.context/` operating system (2026-06-11)

**Status:** Accepted
**Context:** AI agents need durable, tool-agnostic project memory: rules, state, architecture, and history that survive between sessions and across different agent tools.
**Decision:** All agent-facing project knowledge lives as plain markdown in `.context/`, scaffolded by [linkspawn-create-context](https://github.com/LinkSpawnDev/linkspawn-create-context). Root `CLAUDE.md`/`GEMINI.md` are thin pointers into it.
**Trade-offs accepted:** The docs are only as good as their upkeep — `PROJECT_STATE.md` must be updated every session or the system rots.

## DR-002 — One engine, two distributions (2026-06-11)

**Status:** Accepted
**Context:** A full private template set (the LinkSpawn way: skills, audit module, stack invariants) and a public stripped set need to coexist without maintaining two CLIs. Options: two codebases; one repo with build tags; public engine + private overlay.
**Decision:** This public repo ships the engine (`pkg/scaffold`, `pkg/cli`) and the core/minimal presets, built as `spawn` from `cmd/spawn`. The private repo (`linkspawn-context-kit`) embeds its own template overlay and preset list, imports the engine, and builds its own `spawn` superset binary.
**Trade-offs accepted:** `pkg/*` is public API surface — engine changes must respect the private consumer. The public repo must stay honest: it runs the identical engine the private build uses.

## DR-003 — CLI emitting plain markdown, not an MCP server (2026-06-11)

**Status:** Accepted
**Context:** The scaffolder could have been an MCP server or agent plugin. But the whole point of `.context/` is tool-mortality hedging: knowledge in plain files that survive any vendor's product decisions.
**Decision:** A standalone CLI with zero runtime coupling to any agent product. Generated output is plain markdown readable by every tool. Agent-specific integrations, if ever added, wrap the CLI — never replace it.
**Trade-offs accepted:** No in-agent discoverability; users must find and install the binary themselves.
