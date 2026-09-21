---
id: E-002-c1-c3-freeze
doc: execution-entry
parent: GOAL-002-r1-denominator-and-contract-freeze
status: recorded
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# E-002 · C1～C3 冻结（用户 P-004 裁决落盘）

## 事实（2026-09-19）

### 1. 用户 P-004 裁决（本轮唯一决策来源）

| 决策点 | 用户裁决 |
|--------|---------|
| C2 契约形态（`I-038-002`） | **方案 B** · 另立本地模块自有异步契约；ADR-0022 同步语义完全冻结；协议 pin 零改动 |
| C3 首波分母（`I-038-003`） | **仅「新建批量导出所选」**（1 条） |
| C1 作用域模型（`I-038-001`） | **管理作用域 + 新增 `jobs.read` 权限**（`PolicyAdmin`） |

裁决以结构化提问形式取得（三个决策点各含备选与取舍说明）；未选方案记录于 `01-decision/D-001-…` §4。

### 2. 落盘产物

| 产物 | 路径 | 对应检查点 |
|------|------|-----------|
| R1 冻结决策 | `01-decision/D-001-r1-contract-and-denominator-freeze.md` | C1/C2/C3 |
| Job 种类×作用域矩阵 | `attachments/r1-job-kind-scope-matrix.md` | C1 |
| 首波分母矩阵 | `attachments/r1-first-wave-denominator-matrix.md` | C3 |

### 3. 信息项状态变更

| ID | 变更前 | 变更后 | 依据 |
|----|--------|--------|------|
| `I-038-001` | open | **verified** | 侦察 + C1 矩阵（D-001 §2） |
| `I-038-002` | open | **verified** | 侦察 + 用户裁决方案 B（D-001 §1） |
| `I-038-003` | open | **verified** | 侦察 + 用户裁决首波分母（D-001 §3） |

`GOAL-002` 检查点：C1/C2/C3 → 完成，`progress: 0/4 → 3/4`。C4（审计与投影）仍待执行。

### 4. 冻结要点摘要

- **C2**：异步批量操作由 `admin.jobs` 自有本地端点承载（202 + jobId → 轮询 → 结果）；ADR-0022 的 `batchMapping` / `$selection.keys` / `runBatchRequest` / `200 {"deleted": n}` → reload 清选路径**逐字保持**；不触碰任何 pinned 工件。
- **C1**：`admin.jobs` 提供跨 actor 管理读面，新权限 `jobs.read`（+ 写权限，键名 R2 冻结）；**既有 `GetForActor` actor 隔离语义与其冻结测试不动**，通用读面走新方法 + 新路由。
- **C3**：首波 = 新建「批量导出所选」1 条；同步 `batch-delete`（S-1/S-2）保持；排除清单 X-1～X-11 明确登记。

### 5. 派生结论 vs 用户裁决（诚实标注）

`D-001` §3.3 把「`users`/`roles` 同步 `batch-delete` 保持同步」记为**派生结论**（依据 VP-038 显式非目标 `VP-038:51`），**不**冒充本轮用户裁决——本轮 P-004 提问中该条为确认项，用户未勾选。若用户意图改变，须回到 VP-038 层修订非目标后再改 R1。

### 6. 本轮**未**做的事（边界）

- **未**改动 `apps/**`：本轮写入全部为 `docs/workspaces/workspace-038-…/GOAL-002-…/` 下的治理文档。
- **未**执行审计：C4 待执行（模式 `cross`：self + 本地 grok build 4.6 high independent）。
- **未**创建 R2 子目标：待 C4 通过后按 P-001 立项。
- **未**把未定项写成已裁决：O-1（前端触发机制）/ O-2（capability 声明口径）/ O-3（管理列表索引）三项显式登记为 R2 方案项（D-001 §1.3）。
- **未**在愿景层改判 `V-F126`：本目标只完成承接动作（分母明确化），闭合登记属 `/vision`。

### 7. Git checkpoint

见 `E-003`（R1 冻结交付物提交记录）。
