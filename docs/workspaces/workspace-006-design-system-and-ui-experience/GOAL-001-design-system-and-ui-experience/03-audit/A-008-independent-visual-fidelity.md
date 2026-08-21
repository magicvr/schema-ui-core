---
id: GOAL-001-design-system-and-ui-experience
doc: audit-entry
record_id: A-008
source: independent
scope: Root S2/S3 视觉 fidelity 复审（对照 D-004；重做 commit f16dc9f / 5716df9；开放 F-VUI-001/002）
verdict: pass
status: recorded
parent: GOAL-001-design-system-and-ui-experience
created: 2026-08-09
updated: 2026-08-09
version: 0.1.0
auditor: grok build (independent cross-audit)
---

# A-008 · Independent visual fidelity audit (Root · post-rework)

## Range

| Field | Value |
|-------|-------|
| type | independent cross-audit |
| covered | Root open required **F-VUI-001 / F-VUI-002** (from A-006); related recommended **F-VUI-004**; code under claimed rework commits `f16dc9f` (S2) / `5716df9` (S3) |
| inputs | Root D-004；`attachments/visual-direction-stitch-summary.md`；A-006 / A-007；GOAL-003 E-002 / A-002 / A-003（self — **not** rubber-stamped）；`data-table.tsx`、`render.tsx` RecordView、`form-controls.tsx`、`App.tsx`、`LoginPage.tsx`；tests listed below |
| excluded | Pixel-diff against gitignored Stitch PNGs（D-001 / E-002 explicitly do not require）；full e2e re-run in this session；status/progress/goal-tree edits（auditor does not change them） |

## Anti-pattern check (A-006)

A-006 fail was **process theater**: S2 reduced to chart pie Token color; S3 reduced to mobile hamburger drawer; Login / form / table dual-end / recordView Drawer **zero structural change**.

This re-audit asks: did E-002 replace that with **observable structure on shipped paths**, or only docs/data-attributes?

| A-006 gap | Post-rework evidence | Theater? |
|-----------|----------------------|----------|
| No mobile card list | `data-table.tsx` dual-end: desktop `table` (`md:block`) + `MobileCardList` (`md:hidden`, title + 1–2 secondary + actions/⋯)；`SchemaTable` mounts `DataTable` | **No** — two DOM presentations |
| No recordView Drawer/Sheet | `render.tsx` `RecordView`: selection-driven fixed right drawer + backdrop；`max-sm` bottom sheet；`role="dialog"`；static fixtures keep inline panel | **No** — layout chrome, not color swap |
| Form/login unused primitives | `form-controls.tsx` → `Input`/`Label`/`Textarea`；`LoginPage` → `Card`/`Input`/`Label`/`Button`；`StatCardView` → `Card` | **No** — real imports on main paths |
| Shell = drawer only | `App.tsx`: sticky topbar + permanent `w-64` sidenav (`lg:block`，`data-shell-sidenav-width="256"`) + mobile drawer retained as sub-capability；nav active language | **No** — shell layout language beyond hamburger |
| Login zero diff | `LoginPage.tsx` Card surface, branding row, radial backdrop, design-system marker | **No** |

## Prior findings (re-evaluated)

### F-VUI-001 · S2 将「Token 接线」偷换为「Renderer 视觉重构」完成

| Field | Value |
|-------|-------|
| level | required |
| status | **fixed** |
| evidence | (1) Dual-end list: `apps/web/src/components/data-table.tsx` (`data-table-presentation=dual-end` / `desktop-table` / `mobile-cards`)；consumed by `apps/web/src/renderer/schema-table.tsx` (`<DataTable …/>`). (2) recordView Drawer/Sheet: `apps/web/src/renderer/render.tsx` `RecordView` (`data-record-view=panel|backdrop`, drawer mode when selection-driven；production schemas include `recordView`, e.g. `apps/api/internal/modules/users/schema/users.json`). (3) Form + display surfaces: `form-controls.tsx` design-system primitives；`StatCardView` uses `Card`. (4) Tests: `data-table.test.tsx` dual-end assertion；`visual-fidelity.test.tsx` FormControls + static recordView + dual-end markers. Meets A-006 closure minimum: dense desktop table + mobile cards + recordView Drawer/Sheet + at least one of form/login/display upgraded. |
| impact | Unblocks honest Root S2 checkmark **only after** orchestrator sync；does **not** itself set Root `done` |
| note | Still not pixel-perfect Stitch；structure matches D-004 §4–5 product rules. Residual recommended findings below do not re-open this required item. |

### F-VUI-002 · S3 未满足 D-004 壳与工作流呈现分母

| Field | Value |
|-------|-------|
| level | required |
| status | **fixed** |
| evidence | (1) Shell: `apps/web/src/app/App.tsx` — `data-shell=admin` / `topbar-sidenav`；sticky topbar (`data-shell-region=topbar`)；desktop sticky sidenav `w-64` ≈256px；rounded active nav；user chip `border`+`bg-card`；mobile hamburger drawer retained. (2) Sign-in: `apps/web/src/app/LoginPage.tsx` — Card/Input/Label/Button；`data-login-surface=design-system`. (3) Tests: `shell.test.ts` structural source checks for topbar/sidenav width markers；`LoginPage.test.tsx` primitive consumption. Meets A-006 closure: shell + login observable upgrade beyond drawer-only. |
| impact | Unblocks honest Root S3 checkmark after orchestrator sync；not Root `done` by itself |
| note | Dialog/Confirm hosts remain hand-rolled buttons with Token classes (S1-era); language is consistent enough with shell (overlay/card/primary) for C2 — not a required re-open. |

### F-VUI-004 · S1 primitives「可发现」≠「主路径已消费」（recommended）

| Field | Value |
|-------|-------|
| level | recommended |
| status | **fixed** |
| evidence | Main paths now import and render `Card` / `Input` / `Label` / `Textarea` (`LoginPage`, `form-controls`, `StatCardView`). Closure intent from A-006 (“S2 主表面接入 Card/Input 等”) satisfied. `Badge` still unused — non-blocking. |
| impact | Explains post-rework user-visible change on login/forms/stats; no longer explains “looks unchanged” as missing primitive consumption |

### F-VUI-003

Already **fixed** in A-007 (status rollback). Not re-opened.

## New findings (this audit)

### F-VUI-005 · Selection-driven recordView drawer lacks dedicated runtime test

| Field | Value |
|-------|-------|
| level | recommended |
| status | open |
| evidence | `visual-fidelity.test.tsx` only mounts **static** `props.record` panel（asserts backdrop **absent**）. Selection path (`canClose` + backdrop + close clears `crud.selectRow`) is implemented in `render.tsx` + `SchemaTable` `onRowClick`, but no fidelity test opens drawer from row select. |
| impact | Regression risk on the primary production chrome；does **not** erase code evidence for F-VUI-001 |
| closure | Add jsdom/integration test: select row → `[data-record-view=backdrop]` + panel `aria-modal` → close clears selection |

### F-VUI-006 · Mobile Sheet breakpoint uses `max-sm` (&lt;640) vs D-004 mobile band (&lt;768)

| Field | Value |
|-------|-------|
| level | recommended |
| status | open |
| evidence | D-004 §3/§5: mobile &lt;768 → full-height Sheet. `RecordView` drawer classes use `max-sm:…` (Tailwind sm=640). Between 640–767px user may still see right drawer, not sheet. List dual-end correctly uses `md` (768). |
| impact | Minor fidelity gap at tablet-narrow；not a D-004 §4 list failure |
| closure | Align sheet breakpoint with `md` / documented mobile band, or document intentional sheet breakpoint with rationale |

### F-VUI-007 · `shell.test.ts` pure drawer state helpers are re-implementation theater

| Field | Value |
|-------|-------|
| level | recommended |
| status | open |
| evidence | First describe block defines local `openDrawer`/`closeDrawer` not imported from `App.tsx` — tautological. **Real** S3 evidence is App source structure checks (`w-64`, sticky topbar, data-shell markers) + actual `App.tsx` implementation. |
| impact | Overstates test strength if cited alone； does not invalidate shell code |
| closure | Prefer mounting shell fragments or integration coverage; keep/expand structural + integration assertions |

## Verdict

**pass**

Relative to A-006 fail and D-004 product rules (as operationalized by GOAL-003 D-001: structure + real-module tests, not mandatory Stitch pixel diff):

- Open required **F-VUI-001** → **fixed** (evidence above).
- Open required **F-VUI-002** → **fixed** (evidence above).
- Recommended **F-VUI-004** → **fixed**.
- New recommended residual: F-VUI-005 / F-VUI-006 / F-VUI-007 — **do not** block S2/S3 honesty claims.

**Orchestrator note (not executed here):** May propose Root S2/S3 checkmarks and GOAL-003 C1/C2 after responding to this opinion. **Still require** user written confirmation before Root/`workspace` `done` again (D-006 / A-006 process lesson). Do not treat A-003 self-pass alone as independent confirmation — this A-008 is the independent ledger entry.

**Fail-closed confirmation:** Delivery is **not** token-only or docs-only. Dual-end list, Drawer/Sheet recordView, shell topbar+sidenav, and login Card surface are structural code changes on shipped modules.
