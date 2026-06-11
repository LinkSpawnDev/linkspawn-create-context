<!-- scaffolded by spawn dev (https://github.com/LinkSpawnDev/linkspawn-create-context) on 2026-06-11 -->
# PROJECT_STATE.md — linkspawn-create-context

**Last updated:** 2026-06-11
**Current phase:** Phase 0 — v0.1 public release
**Owner:** Nicholas Pegler

This file is updated at the **end of every session** by whichever agent ran it. It is the single source of truth for current state and blockers. State must be honest — what works, what doesn't, what's stubbed.

---

## Phase status

| Phase | Description | Status |
|---|---|---|
| Phase 0 | Engine + core/minimal presets + `spawn init`; public v0.1 release | **In progress** — target **Mon 2026-06-15** |
| Phase 1 | `spawn doctor [--fix]` — validate existing `.context/` folders | Not started |
| Phase 2 | `spawn add <module>`, `spawn lesson`, private overlay (`linkspawn-context-kit`) | Not started |

## Active work

- Engine (`pkg/scaffold`): rendering, nested dirs, conflict refusal, `--force`, idempotent `.gitignore` append — **done, tested** (6 tests green).
- CLI (`pkg/cli`): `init` with flags + huh interactive form, presets, version stamping — **done**; interactive path manually unverified (flags path verified end-to-end).
- Templates: `core` (6 `.context/` files + CLAUDE/GEMINI pointers + gitignore block) and `minimal` (root AGENTS.md + `@AGENTS.md` pointer) — **done**, distilled from focalmomentum-tube's AGENTS.md lineage.
- Dogfooded on this repo itself (this very `.context/` is `spawn init` output, sections filled in).

## Next actions

1. Owner: review, commit, and push the v0.1 tree (agent never pushes).
2. Verify CI goes green on GitHub Actions, then tag `v0.1.0` with `-ldflags "-X main.version=v0.1.0"` documented in the release notes — **target Mon 2026-06-15**.
3. Manually verify the interactive `spawn init` form in a real terminal (huh path is untested by automation).
4. Phase 1: `spawn doctor` — checks list already specified in the private `context-review/CLI_BLUEPRINT.md`.

## Blockers & open questions

- None.

## Session status

**Last session (2026-06-11):** Repo bootstrapped from scratch (prior attempt's engine patterns carried over, code rewritten as exported `pkg/`). Engine, CLI, both presets, README, CI workflow written; tests green; dogfooded on itself.

**Next steps:** see Next actions above.
