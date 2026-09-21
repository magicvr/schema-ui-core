---
id: E-003-tsc-evidence-residue-closeout
doc: execution-entry
status: recorded
goal_id: GOAL-043-w31-cross-workspace-residual-closeout
created: 2026-09-18
updated: 2026-09-18
parent: GOAL-001-design-implementation-conformance
version: 1.0.0
---

# E-003 · `tsc` 文档证据余项收口（C2 · 承接 `GOAL-008 A-002 F-002`）

## 事实

2026-09-18，对「全仓 `tsc` 简写未逐条裁定」按 `D-001` §3 的口径收口（不逐条考古）。

### 1. 可执行面（决定性）

`package.json`、`.github/workflows/*.yml`、仓库 `scripts/` 与 `apps/web/scripts/` 中的 `tsc` 调用由 `apps/web/src/typecheck-convention.guard.test.ts`（`GOAL-008` 建立、`GOAL-010` 加固）**每次 `npm test` 断言**：不得出现非检查型调用（`-b`，或 `-p` 指向自身选择源文件的配置）。CI 另有显式 `npm run typecheck` 门禁。因此「明日再有失实类型检查证据」的通道已被关闭，正确性不再依赖历史文档形态。

### 2. 文档侧量化（本次扫描，`git grep "tsc" -- docs`）

| 分类 | 行数 | 说明 |
|------|------|------|
| 命令形态可确定（`tsc -b` / `tsc -p <含 include 的配置>`） | 主体 | 正确口径，无需处理 |
| 指 `npm run build`（= `tsc -b && vite build`） | 45 | 形态由脚本确定，实质为真实检查 |
| 字面 `tsc --noEmit` | 40 | workspace-037 已更正或已加勘误注记；跨区条目由 `GOAL-008 E-006` 加注 |
| **叙述式、形态不可从文本唯一确定**（如「tsc clean / tsc 0 / tsc 未受影响」） | **269** | 本次口径下的 bounded residual |
| 口径下「提到 tsc 但未限定形态」合计 | 354 | 含上两类中的部分重叠计数 |

### 3. 处置（bounded residual）

- **不逐条考古**：文本无法唯一确定当时的命令形态，逐条追溯既不可核对也无收益。
- **触发条件**：某条历史记录被**再次当作类型检查证据引用**时，按 `GOAL-008 D-001` 的正确口径复核，并在引用处注明形态或指向本目标。
- **登记**：数量、口径与触发条件写入 `docs/vision/roadmap.md` 的「未决项统一登记」节（`E-004`），避免散落遗忘。

C2 完成。回归数据见 `E-004`。
