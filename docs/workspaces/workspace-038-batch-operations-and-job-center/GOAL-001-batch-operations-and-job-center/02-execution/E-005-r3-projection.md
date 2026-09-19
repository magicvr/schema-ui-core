---
id: E-005-r3-projection
doc: execution-entry
parent: GOAL-001-batch-operations-and-job-center
status: recorded
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# E-005 · R3 关门与 Root 投影

## 事实（2026-09-19）

### 1. R3 纲领阶段完成

Root `GOAL-001` 的纲领检查点 **R3（批量操作异步承接）** 由子目标 `GOAL-004-r3-async-batch-operation` 交付并已于同日关门：

| 项 | 值 |
|----|-----|
| 子目标状态 | `active` → **`done · 4/4`** |
| 检查点 | C1 异步写面与进度 / C2 前端触发 / C3 同步路径回归 / C4 R3 审计与投影 —— **全部完成** |
| 开放 required | **0** |
| 审计链 | `A-001` self `pass` → `A-002` independent（grok-build · grok-4.6 · high · `/audit`）**`pass`** → `A-003` 响应（8 条 recommended 全 `fixed`） |
| 关门依据 | 4 检查点全达成 + 开放 required = 0 + 无 P-004 待裁冲突 |

> 子目标关门属**非关键决策**，按用户目标轮次指令「子目标关门等非关键决策可经交叉审计后静默执行」——已经 `cross` 模式审计后执行，未额外请求用户裁决。

### 2. Root R3 检查点投影

| Root 检查点 | 投影前 | 投影后 |
|-------------|--------|--------|
| R3 批量操作异步承接 | `[ ]` | **`[x]`** |

Root `progress: 2/5` → **`3/5`**（`00-meta.md` frontmatter 与 `goal-tree.md` 树/表已同步）。

**VP-038 判据 3 达成**：至少一条真实批量操作以异步 Job 承接——选中行 → `202 + jobId` → 可观察 `queued→running→终态` 与**进度** → 终态读取结果；同时既有 ADR-0022 同步 `batch-delete` 的已交付语义与回归**未退化**（独立审计以空 diff 与全绿回归核实）。

### 3. R3 交付结论（可核对的实现事实）

| 项 | 内容 |
|----|------|
| 新 Job kind | `jobs.batch-export`（runner 启动前注册） |
| 新路由 | `POST /api/jobs/batch-export` → **202 + job 投影** |
| 新权限 | `jobs.write`（`PolicyAdmin`）；与既有 `data.export` **并列门禁** |
| 数据面 | `users` / `roles`（与同步导出同分母）；列集/转义/BOM **复用同步导出**；上限 500 id |
| 进度 | 按选中行线性上报（`done*85/total` → 90 → 99），100 由运行时置 |
| 结果 | `{resource,rowCount,fileName,csv}`，经 R2 的 `GET /api/jobs/{id}/result` 读取 |
| 前端 | `jobs-batch-export` 自定义组件；users 页 `props.selection.mode=multiple` + custom 节点；**未**声明 `table.selection`/`actions.batch.request`；**不调用** `reloadList()` |
| 安全口径 | 异步路径**不扩大**数据外带面（双重门禁）；第二道门当前为 defence-in-depth（`PolicyAdmin ⊂ PolicyAdminEditor`），已加**变异测试验证过的嵌套守卫** |
| 回归 | Go `go test ./...` 全绿；web `1440/1440`；`tsc -b` exit 0；同步批量回归锚点 48/48；unknown-custom 警告 **0** |

### 4. 同步文件

| 文件 | 变更 |
|------|------|
| `GOAL-001/00-meta.md` | R3 检查点勾选；`progress: 3/5`；`version: 0.4.0`；子目标表更新为 `done` |
| `GOAL-001/02-execution.md` | 索引补 E-005 |
| `goal-tree.md` | 树/纲领路线图/状态表/说明节同步（Root `3/5`、GOAL-004 `done · 4/4`） |
| `workspace.md` | Root 投影行与纲领阶段表更新；补 R3 关门记录与交付要点 |

### 5. 移交 R4

| 项 | 内容 |
|----|------|
| 结果中心完整体验 | 列表/详情/进度/终态/过期/重试/取消/下载（Root 路线图 R4） |
| 前端组件的**交互级**测试 | R3 已闭合"注册缺口"（警告归零）；交互覆盖（空选择禁用/只收 202/不 reloadList/终态下载）随 R4 的呈现改动一并补测 |
| `I-038-013` | 非阻塞（最晚 R4）：导出文件名/格式与两个导出入口的 UI 文案一致性 |
| `jobs.write` 的写操作面 | 重试/取消等结果中心动作（R2 未声明，R3 仅声明提交写键） |

### 6. 边界

- **未**触碰任何 pinned 协议工件（`docs/schemas/**`、`apps/web/src/protocol/upstream/**`）——独立审计核实五个 checkpoint 的 `--stat`。
- **未**改同步 `batch-delete` 的路由、权限、原子语义或其协议 fixture。
- **未**经 ADR-0022 `batchMapping`/`runBatchRequest` 提交（K-1 落地）。
- **未**重开 VP-012（R2 的 wallet 结果地址改动为字节等价重构，R2 独立审计已复核）。
- **未**在愿景层改判 `V-F126`（承接动作已完成，闭合登记属 `/vision`）。

### 7. Git checkpoint

见提交记录（R3 关门 + 审计响应 + Root 投影），覆盖路径 `docs/workspaces/workspace-038-batch-operations-and-job-center/`。
