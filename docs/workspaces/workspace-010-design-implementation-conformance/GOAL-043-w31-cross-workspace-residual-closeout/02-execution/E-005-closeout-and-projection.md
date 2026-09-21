---
id: E-005-closeout-and-projection
doc: execution-entry
status: recorded
goal_id: GOAL-043-w31-cross-workspace-residual-closeout
created: 2026-09-18
updated: 2026-09-18
parent: GOAL-001-design-implementation-conformance
version: 1.0.0
---

# E-005 · 回归、自审、回填与投影（C4）

## 回归

| 检查 | 结果 |
|------|------|
| `npx vitest run src/host/resource-feedback-parity.test.ts` | **3 passed**（变异验证：删 host 条目 / 改 locale key 均被指名报出） |
| `npx vitest run src/typecheck-convention.guard.test.ts` 等受影响面 | 通过（见下全量） |
| `npm test`（全量 Vitest） | **114 文件 / 1437 测试全部通过**（较 GOAL-011 时点 113/1434 增 1 文件 + 3 用例） |
| `npm run typecheck`（`tsc -b` + e2e 工程） | exit **0** |
| `npm run test:e2e -- list-visual-surface`（`APP_PROFILE=admin`） | **4 passed**（1.2m；含 new 暗色断言与双页面分页用例） |
| 同上（`APP_PROFILE=mvp`） | **4 passed**（1.3m） |
| `git diff --check` | 通过 |
| 产品行为变更 | 无：仅测试新增 + `HostFailureScreen` 导出文案表（无行为差异）；变异实验均已还原 |

## 自审与回填

- `A-001`（self）verdict `pass`，开放 required = 0。
- 原 finding 回填（按 P-003 只加闭合注记，不改其 `status`/`progress`/历史正文）：
  - `GOAL-009` `A-001 F-001`、`F-002` → **fixed**（`03-audit/A-001-*.md` 条目 + `03-audit.md` 索引与结论段）；
  - `GOAL-005` `A-002 F-002` → **fixed**（`03-audit.md` 响应表 + 结论段）；
  - `GOAL-008` `A-002 F-002` → **bounded residual**（`03-audit.md` 索引 + 结论文）；
  - `GOAL-006` `D-002` 追加后续说明；workspace-037 Root `00-meta` 备注、VP-037 计划残余段、`docs/vision/workspaces.md`、`roadmap.md` 同步指针。

## 投影

| 层 | 变化 |
|----|------|
| `GOAL-043-w31-cross-workspace-residual-closeout` | `active · 0/4 → done · 4/4` |
| workspace-010 Root `GOAL-001-design-implementation-conformance` | 不变（长期程序容器保持 `active`）；`goal-tree.md` 与 `workspace.md` 追加 W31 行 |
| `docs/vision/roadmap.md` | v0.91.0；新增「未决项统一登记」节 |
| `docs/vision/reviews.md` / `revisions.md` | `VRev-097`（self `pass`，`V-F124` fixed）；`VR-083`（editorial） |
| VP-037 / workspace-037 | **不重开**：其 `closed` v1.7.0 与 Root `done · 6/6` 不变；仅残余状态与指针更新 |
