---
id: A-002-independent-cross-audit-s4-s5
title: S4+S5 independent cross-audit — implementation solid, closeout ledger not yet honest
date: 2026-08-09
source: independent
scope: GOAL-004 (S4) + GOAL-005 (S5) combined
verdict: conditional
parent: GOAL-005-s5-regression-fork-example-and-closeout
provider: grok build CLI (model grok-4.5, reasoning effort high) — run by the user directly in a terminal outside this session; pasted verbatim into the session by the user for ledger recording (P-003 "写入：交叉工具直接追加，或代贴并保留 `source: independent`").
---

# A-002 · Independent cross-audit — S4+S5

## Evidence reviewed

- Root: `docs/workspace-006-design-system-and-ui-experience/GOAL-001-design-system-and-ui-experience/00-meta.md`
- S4: `GOAL-004-s4-state-and-feedback/{00-meta,01-decision,02-execution,03-audit}.md`, `01-decision/D-001-s4-baseline-and-plan.md`, `02-execution/E-001-s4-implementation.md`, `03-audit/A-001-s4-self-audit.md`
- S5: `GOAL-005-s5-regression-fork-example-and-closeout/{00-meta,01-decision,02-execution}.md`, `01-decision/D-001-fork-example-and-regression-scope.md`, `02-execution/E-001-s5-implementation.md` (no `03-audit*` / `attachments/` present on disk)
- Implementation: `apps/web/src/components/ui/async-state.ts`, `async-state.test.ts`, `data-table.tsx`, `data-table.test.tsx`, `renderer/render.tsx` (StatCardView / ChartView / `useDisplayData`), `renderer/render.test.tsx` (Skeleton status tests), `theme/brand.example.css`, `theme/brand-example.test.ts`, `theme/README.md`, `index.css` token declarations, `app/shell.test.ts` (`_state` rename)
- Git: commit `6ce76f4` (S4 code) — 7 files; S5 theme files uncommitted staged alongside GOAL-005 docs
- Commands (read-only, `apps/web`): `npm run test` → **29 files / 616 tests passed**; `npm run build` → **exit 0** (`tsc -b && vite build`). Playwright e2e was **not** re-run in this audit (out of prompt command list).

## Findings

1. **required** — GOAL-005 closeout prose asserts facts that are not true yet. `GOAL-005/.../00-meta.md` claims Root `progress: 5/5`, S1–S5 all complete, and that independent audit already lives at `03-audit/A-002-independent-cross-audit-s4-s5.md`. On disk: Root still has S5 unchecked and `progress: 4/5`; the A-002 path does not exist (this opinion is chat-only per audit prompt). Must be rewritten to factual present-tense after ledger write + Root sync, or closeout will launder anticipatory claims as evidence.

2. **required** — GOAL-005 five-piece scaffold is incomplete vs AGENTS §3. Folder currently has only `00-meta` / `01-decision*` / `02-execution*`; missing `03-audit.md`, `03-audit/`, and `attachments/`. Cannot lawfully index A-002 (or any formal opinion) under the required audit ledger layout until scaffolded.

3. **non-blocking** — `useDisplayData` never clears `error` on refetch start or on success (`render.tsx` ~1113–1134). A failed fetch followed by a successful reload (same component instance / `reloadToken`) can keep showing the stale error because S4 correctly prefers `error` over `loading`/`ready`. Pre-existing hook hygiene, not introduced by the pure function, but S4's centralized error-first path makes the sticky-error surface explicit; worth fixing when retry UX matters.

4. **non-blocking** — DataTable unit coverage for the new `role="alert"` error branch is weak. `data-table.test.tsx` asserts error **text** only, not `role="alert"` (unlike loading's `role="status"` / `.animate-pulse`). Production code does set `role="alert"` on the error `<td>` (intentional a11y upgrade vs pre-S4 loading-first table body); low regression risk given text assertions + schema-table paths, but the unit test does not lock the a11y contract D-001 claimed.

5. **non-blocking** — S5 structural test is meaningful (non-empty override set, strict subset of `index.css` declared names, requires primary/chart/radius families) and `brand.example.css` only declares tokens already present on the base (`--primary`, `--primary-foreground`, `--chart-1..5`, `--radius`). Residual gap: the regex does not verify **values** are valid CSS or that dark-mode pairs stay consistent; acceptable for "minimal fork example," not a full brand contract.

6. **non-blocking (observation / pass on this axis)** — `resolveAsyncDisplayState` is genuinely consumed by all three claimed sites (`DataTable`, `StatCardView`, `ChartView`). New tests drive real render paths (pure-function unit tests; real `DataTable` mount; `RenderPage` + pending fetcher capturing pre-settle Skeleton `role="status"` then post-settle content). Precedence `error > loading > empty > ready` matches D-001 and the stale-`loading`+error case. Observable error/empty copy for statCard/chart preserved; DataTable emptyMessage preserved; DataTable error gained `role="alert"` (improvement). Scope discipline holds: S4 commit is the seven claimed files; S5 is theme example only; no second Token system; no `I-PROTO-FULL-001` disposition expansion. `shell.test.ts` `_state` rename is minimal and behavior-preserving (unblocks `noUnusedParameters` only).

## Verdict rationale

**conditional**: S4 implementation quality is real (not test theater), precedence is correct, regression suite is green at **616/616** + build exit 0 under this auditor's re-run, and S5 fork example + structural subset test are adequate for a minimal brand-Token sample. Closeout / S5 governance evidence is **not** yet pass-grade: incomplete five-piece, and `00-meta` pre-claims Root 5/5 plus a non-existent A-002 path. Fix findings 1–2 before any Root S5 checkmark or `/govern` closeout proposal; findings 3–5 are residual quality items, not blockers for treating the **code** deliverables as acceptable.
