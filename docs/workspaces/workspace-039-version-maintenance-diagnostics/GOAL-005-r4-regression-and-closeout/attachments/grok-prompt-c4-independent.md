# Grok Build · 独立关门审计（GOAL-005 R4）

在仓库根执行 `/audit`。模型 grok-4.6 · reasoning high。`source: independent`。只出意见并落盘，不改 status/progress/goal-tree。

## 目标

`[workspace-039-version-maintenance-diagnostics] GOAL-005-r4-regression-and-closeout`

scope：VP-039 判据 1–7、R4 C1–C3、回归矩阵、边界与 Root/VP 关门提请。用户书面确认仍是独立门禁，不由本审计代替。

## 必读

1. workspace-039 `workspace.md`、`goal-tree.md`
2. Root `00-meta.md`
3. GOAL-005 `00-meta.md`、`D-001-r4-regression-denominator.md`
4. `attachments/r4-regression-matrix.md`
5. `02-execution/E-002-r4-regression.md`、`E-003-r4-exit-matrix.md`
6. `03-audit/A-001-r4-self.md`
7. VP-039 plan + GOAL-002/003/004 audit indexes
8. Code boundaries: `apps/api/internal/handler/bootstrap.go`, `operational.go`, `account.go`, `apps/web/src/app/runtime-banner.tsx`, `version-chip.tsx`, `AuthGate.tsx`, `App.tsx`, `kernel/profile.go`, `protocol/upstream/**`

## 核验

- R1–R3 cross audit chains and required closure.
- Go full, Web full Vitest, forced TypeScript, browser results and whether claims match evidence.
- mvp/admin/profile permissions, mode mapping, zh/en, light/dark evidence.
- Existing full-browser fresh-seed residual: determine whether it is unrelated bounded residual or VP-039 blocker; isolated spec and admin slice evidence.
- Product boundary: no default Profile change, no pinned protocol change, no VP-012 gate change, no new gated infra.
- User close gate must remain open; do not mark VP/Root closed.

## 输出

verdict + required/recommended findings, comparison with A-001, and explicit recommendation whether R4 evidence is sufficient to ask the user for written Root/VP close confirmation. If possible, write `03-audit/A-002-r4-independent.md` and update `03-audit.md`.
