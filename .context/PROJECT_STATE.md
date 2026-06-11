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

1. Owner: review + push branch `ci/release-pipeline`, merge to main (agent never pushes).
2. **Before tagging** — create the public `LinkSpawnDev/homebrew-tap` repo (empty, default branch is enough) and add a fine-grained PAT with write access to it as the `HOMEBREW_TAP_GITHUB_TOKEN` actions secret on this repo. The release workflow fails at the cask-publish step without both.
3. Tag `v0.1.0` (`git tag v0.1.0 && git push --tags`) — **target Mon 2026-06-15**. The release workflow gates on gofmt/vet/test, then GoReleaser builds darwin/linux/windows × amd64/arm64, publishes the GitHub Release (version auto-stamped from the tag), and pushes the `spawn` cask to the tap.
4. Manually verify the interactive `spawn init` form in a real terminal (huh path is untested by automation).
5. Phase 1: `spawn doctor` — checks list already specified in the private `context-review/CLI_BLUEPRINT.md`.

## Blockers & open questions

- None.

## Session status

**Last session (2026-06-11):** Repo bootstrapped from scratch (prior attempt's engine patterns carried over, code rewritten as exported `pkg/`). Engine, CLI, both presets, README, CI workflow written; tests green; dogfooded on itself.

**Session 2 (2026-06-11):** Repo pushed public; CI green. On branch `ci/release-pipeline`: bumped `actions/checkout` v4→v5 and `actions/setup-go` v5→v6 (Node 20 deprecation; Node 24 forced default 2026-06-16); added `.goreleaser.yaml` (v2, `homebrew_casks` — the `brews` section is deprecated) and `.github/workflows/release.yml` (tag-triggered, gofmt/vet/test gate before GoReleaser); README install section now lists brew/binaries/go install. Verified locally: `goreleaser check` + full `--snapshot` build (6 targets, cask generated, version stamp confirmed), gofmt/vet/test/build all clean. Release blocked on owner creating the tap repo + token secret (Next actions 2).

**Next steps:** see Next actions above.
