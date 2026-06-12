# linkspawn-create-context

**`spawn` scaffolds a `.context/` operating system for AI-agent-driven development** — agent rules, project state, architecture notes, decision logs, and a lessons loop, as plain markdown that any coding agent (Claude Code, Gemini CLI, Cursor, or the next one) can load.

```text
$ spawn init

  created .context/AGENTS.md
  created .context/ARCHITECTURE.md
  created .context/DECISIONS.md
  created .context/LESSONS.md
  created .context/PROJECT_STATE.md
  created .context/README.md
  created .gitignore
  created CLAUDE.md
  created GEMINI.md

✓ core preset scaffolded for "my-project".
```

## Why

AI coding agents are stateless. Every session starts cold: the agent doesn't know your guardrails, your locked decisions, your definition of done, or the mistake it made last Tuesday. Most people re-explain these every session — or worse, don't, and pay for it.

The `.context/` system externalises that knowledge into durable, tool-agnostic plain text:

| File | Job |
|---|---|
| `AGENTS.md` | The operating contract: hard guardrails, definition of done, what to do when things go wrong |
| `PROJECT_STATE.md` | Shared memory across sessions and agents — phase, active work, blockers. Updated every session |
| `ARCHITECTURE.md` | System design and module boundaries, so agents stop guessing |
| `DECISIONS.md` | Append-only log of locked decisions, so agents stop relitigating them |
| `LESSONS.md` | Rules learned from corrections, read at session start — the agent stops repeating its mistakes |
| `CLAUDE.md` / `GEMINI.md` | Thin root pointers so every tool lands in the same cascade |

These templates weren't written for this repo — they were **distilled from years of running this system across real commercial and personal projects**, then extracted into a scaffolder. The patterns that survived are the ones that earn their keep.

## Install

Homebrew (macOS / Linux):

```sh
brew install LinkSpawnDev/tap/spawn
```

Go toolchain:

```sh
go install github.com/LinkSpawnDev/linkspawn-create-context/cmd/spawn@latest
```

Or grab a prebuilt binary for macOS, Linux, or Windows from the [releases page](https://github.com/LinkSpawnDev/linkspawn-create-context/releases).

## Use

```sh
cd your-project
spawn init                 # interactive: name, mission, owner, preset
spawn init --yes           # accept defaults, no prompts (scripting/CI)
spawn init --preset minimal  # just a root AGENTS.md + CLAUDE.md pointer
spawn init --force         # overwrite existing files
```

`spawn` never overwrites your files unless you pass `--force`, and it appends to an existing `.gitignore` rather than replacing it. Re-runs are idempotent.

### Presets

- **`core`** (default) — the full `.context/` tier plus root pointers. For any project an agent will work in more than once.
- **`minimal`** — a single condensed `AGENTS.md` at the root plus a `CLAUDE.md` containing `@AGENTS.md`. For repos too small to justify a folder.

## The method, in five rules

1. **Guardrails are absolute.** The agent never commits, never pushes, never publishes, never invents facts to fill gaps. Drafting is the agent's job; signing is yours.
2. **State is updated every session.** `PROJECT_STATE.md` is the handoff between today's agent and tomorrow's — honest about what works, what doesn't, what's stubbed.
3. **Decisions are append-only.** A locked decision is superseded, never rewritten — and never silently "fixed" by a helpful agent.
4. **Stop-rules beat cleverness.** When a test fails unexpectedly or data is missing, the right move is *stop and raise*, not improvise. The templates pre-authorise that.
5. **Corrections become rules.** Every "no, don't do that" gets written to `LESSONS.md` and re-read at session start. The quality test for every rule: *would removing this line cause a specific mistake?*

Everything is markdown. Nothing is coupled to one vendor's tool. When a tool dies — and tools die — nothing of yours dies with it.

## Roadmap

- `spawn doctor [--fix]` — validate an existing `.context/`: broken links, stale state, naming drift, junk files
- `spawn add <module>` — optional modules (ADR folder, session log, skills)
- `spawn lesson "<title>"` — append a correctly-formatted lesson without opening the file

## Licence

MIT. Built by [LinkSpawn](https://github.com/LinkSpawnDev).
