<!-- scaffolded by spawn dev (https://github.com/LinkSpawnDev/linkspawn-create-context) on 2026-06-11 -->
# linkspawn-create-context — Lessons Learned

> Updated after every correction from the user. Read at session start.
> Diagnostic test for every rule: "Would removing this line cause a specific mistake?"
> If the answer is no, the rule is noise — delete it.

Entry format:

```markdown
## [YYYY-MM-DD] Short title

**Mistake:** what went wrong, one sentence — the class of error, not the symptom.
**Rule:** the specific behaviour that prevents recurrence.
**Context:** where this applies (paths, layers, situations).
```

---

## [bootstrap] Verification skipped under time pressure

**Mistake:** Tasks marked complete before output was proven correct, leading to rework.
**Rule:** Never mark a task complete without a concrete verification step (test run, diff, constraints check, or explicit user sign-off).
**Context:** All layers.
