<!-- scaffolded by spawn dev (https://github.com/LinkSpawnDev/linkspawn-create-context) on 2026-06-11 -->
# `.context/` — linkspawn-create-context

This folder is the source of truth for what this project is, why it exists, how it's built, and how AI agents should behave inside it. Read it before doing anything in the repo.

## Load order (Context Cascade)

Any agent (Claude, Gemini, future) starting a session reads these in order:

1. **`AGENTS.md`** — the operating contract: guardrails, definition of done, protocols
2. **`PROJECT_STATE.md`** — current phase, active work, blockers
3. **`ARCHITECTURE.md`** — system design, module boundaries
4. **`DECISIONS.md`** — locked decisions (append-only; do not relitigate)
5. **`LESSONS.md`** — rules learned from past corrections

## Update rules

- `PROJECT_STATE.md` is updated **at the end of every session** by whichever agent ran it. State must be honest — what works, what doesn't, what's stubbed.
- `DECISIONS.md` and `LESSONS.md` are append-only. Don't rewrite history; supersede with a new entry.
- The other docs change rarely. When they do, the change is a deliberate act, not a side-effect.

## Voice

This project is owned by Nicholas Pegler and is public-facing: the repo demonstrates the very method it scaffolds. Claims in the README and templates must be defensible against the real practice they were distilled from. If unsure, flag it explicitly — do not paper over.
