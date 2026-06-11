<!-- scaffolded by spawn dev (https://github.com/LinkSpawnDev/linkspawn-create-context) on 2026-06-11 -->
# ARCHITECTURE.md — linkspawn-create-context

System design, module boundaries, and data flow. Change this file rarely and deliberately — it describes decisions, not aspirations.

---

## Overview

A single Go binary (`spawn`) that renders embedded markdown template trees into a target repository. Templates are organised one directory per preset; each tree's layout mirrors the output layout, so adding a file to a preset requires no code change. The engine is exported (`pkg/`) so a private sibling repo can embed its own template overlay and reuse the command tree unchanged — one engine, two distributions.

## Stack

- **Go** + `go:embed` — single static binary, templates compiled in, `go install` distribution (DR-002).
- **spf13/cobra** — command tree.
- **charmbracelet/huh** — interactive `init` prompts; every prompt is also answerable by flag for scripting.
- **text/template** — variable substitution (`ProjectName`, `Mission`, `Author`, `Date`, `CLIVersion`).
- No config files, no network calls, no state outside the target directory.

## Module boundaries

- Dependency direction: `cmd/spawn → pkg/cli → pkg/scaffold`. `templates/` is data, imported only by `cmd/spawn`. No upward imports, no cycles.
- `pkg/scaffold` is pure file mechanics: render an `fs.FS` into a directory with conflict detection. It knows nothing about cobra, presets, or prompts.
- `pkg/cli` owns the command surface and preset selection; it receives templates as an opaque `fs.FS` via `cli.Config`.
- `pkg/*` is public API surface for the private overlay (`linkspawn-context-kit`).

## Data flow

`spawn init` → resolve flags/prompts into `scaffold.Data` → `fs.Sub` the chosen preset out of the embedded FS → render every file (pass 1, in memory, all-or-nothing) → conflict scan against the target directory → write (pass 2). `.gitignore` is the one special case: appended with a managed marker block, idempotent on re-runs.

## Locked decisions

Architectural decisions are recorded in `DECISIONS.md` (append-only). Do not propose alternatives to a locked decision without explicitly re-opening it with the owner.
