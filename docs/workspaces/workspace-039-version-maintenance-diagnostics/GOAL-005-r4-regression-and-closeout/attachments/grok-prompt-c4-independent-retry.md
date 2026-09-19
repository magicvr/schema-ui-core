# Grok Build · R4 独立关门审计重试（精简）

这是 GOAL-005 R4 的独立审计重试。上一次 grok 因重跑大套件超时且未产出 A-002；本次请**不要重跑全量 Go/Vitest/Playwright**，只读证据与代码边界，最多 18 turns，直接产出并落盘 A-002。

执行 `/audit`，模型 grok-4.6 · reasoning high，`source: independent`。不改 status/progress/goal-tree。

## 必读文件

- `docs/workspaces/workspace-039-version-maintenance-diagnostics/workspace.md`
- `goal-tree.md`
- `GOAL-001.../00-meta.md`
- `GOAL-005.../00-meta.md`
- `GOAL-005.../D-001-r4-regression-denominator.md`
- `GOAL-005.../02-execution/E-002-r4-regression.md`
- `GOAL-005.../02-execution/E-003-r4-exit-matrix.md`
- `GOAL-005.../03-audit/A-001-r4-self.md`
- `attachments/r4-regression-matrix.md`
- VP-039 plan and GOAL-002/003/004 audit indexes

## 必查代码边界（不要跑大套件）

- `apps/api/internal/handler/bootstrap.go` / `operational.go` / `account.go`
- `apps/web/src/app/runtime-banner.tsx` / `version-chip.tsx` / `AuthGate.tsx` / `App.tsx`
- `apps/api/kernel/profile.go`
- `apps/web/src/protocol/upstream/` 是否未改（git diff 可查）

## 判定

1. 退出矩阵是否诚实区分 pass、既有 bounded harness residual、用户 close gate。
2. 17/5/1 mvp full + isolated 1 pass + admin 7/1/0 是否足以作为可运行范围证据；不把既有 fresh-seed 排序 residual 升格为 VP-039 产品 required。
3. Go/Vitest/typecheck 结果是否有证据路径。
4. R1–R3 cross chains 是否 required=0。
5. 是否存在越界改动。
6. Root/VP 仍 active 且用户书面确认仍是开放门禁。

## 输出/落盘

verdict + findings，比较 A-001。若无 required，写 `03-audit/A-002-r4-independent.md` 并更新 `03-audit.md` 索引。推荐项可记录，但不得阻止提请用户确认关闭。
