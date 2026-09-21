---
id: E-003-r1-projection
doc: execution-entry
parent: GOAL-001-batch-operations-and-job-center
status: recorded
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# E-003 · R1 关门与 Root 投影

## 事实（2026-09-19）

### 1. R1 纲领阶段完成

Root `GOAL-001` 的纲领检查点 **R1（分母与契约冻结）** 由子目标 `GOAL-002-r1-denominator-and-contract-freeze` 交付并已于同日关门：

| 项 | 值 |
|----|-----|
| 子目标状态 | `active` → **`done · 4/4`** |
| 检查点 | C1 分母与作用域矩阵 / C2 契约形态冻结 / C3 首波分母冻结 / C4 R1 审计与投影 —— **全部完成** |
| 开放 required | **0** |
| 审计链 | `A-001` self `pass` → `A-002` independent（grok-build · grok-4.6 · high · `/audit`）`conditional`（3 required）→ `A-003` 响应 required 全 `fixed` |
| 关门依据 | 4 检查点全达成 + 开放 required = 0 + 无 P-004 待裁冲突 |

> 子目标关门属**非关键决策**，按用户目标轮次指令「子目标关门等非关键决策可经交叉审计后静默执行」——已经 `cross` 模式审计后执行，未额外请求用户裁决。

### 2. Root R1 检查点投影

| Root 检查点 | 投影前 | 投影后 |
|-------------|--------|--------|
| R1 分母与契约冻结 | `[ ]` | **`[x]`** |

Root `progress: 0/5` → **`1/5`**（`00-meta.md` frontmatter 与 `goal-tree.md` 树/表已同步）。

### 3. R1 冻结结论（用户 P-004 裁决 2026-09-19）

| 决策点 | 裁决 | 落盘 |
|--------|------|------|
| C2 契约形态 | **方案 B**：另立本地模块自有异步契约；ADR-0022 同步语义逐字冻结；协议 pin 零改动 | `GOAL-002/01-decision/D-001-…` §1 |
| C3 首波分母 | **仅「新建批量导出所选」1 条**；同步 `batch-delete` 保持；排除清单 X-1～X-12 | `D-001` §3；`attachments/r1-first-wave-denominator-matrix.md` |
| C1 作用域模型 | **管理作用域 + 新增 `jobs.read` 权限**（`PolicyAdmin`）；既有 actor 隔离语义与测试不动 | `D-001` §2；`attachments/r1-job-kind-scope-matrix.md` |

### 4. 同步文件

| 文件 | 变更 |
|------|------|
| `GOAL-001/00-meta.md` | R1 检查点勾选；`progress: 1/5`；`version: 0.2.0`；子目标表更新为 `done` |
| `GOAL-001/02-execution.md` | 索引补 E-003 |
| `goal-tree.md` | 树/纲领路线图/状态表/说明节同步（Root `1/5`、GOAL-002 `done · 4/4`）；`version: 0.3.0` |
| `workspace.md` | Root 投影行与纲领阶段表更新；补 R1 关门记录 |

### 5. 移交 R2

`D-001` §5 的 T-1～T-8 与未定项 O-1（前端触发机制）/ O-2（capability 声明口径）/ O-3（管理列表索引）随 R1 关门移交 R2 子目标（按 P-001 立项）。三项未定项**必须在 R2 方案中冻结**，不得默认沿用。

### 6. 边界

- **未**改动 `apps/**`：R1 全程为治理文档写入。
- **未**触碰 pinned 协议工件（`docs/schemas/**`、`apps/web/src/protocol/upstream/**`）。
- **未**在愿景层改判 `V-F126`（承接动作已完成，闭合登记属 `/vision`，已登记为交接项）。
- **未**重开 VP-012/011/037/036；**未**解除任何 gated 基础设施条件。

### 7. Git checkpoint

见 `E-004`（提交记录，覆盖路径 `docs/workspaces/workspace-038-batch-operations-and-job-center/`）。
