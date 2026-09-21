---
id: E-004-r2-projection
doc: execution-entry
parent: GOAL-001-batch-operations-and-job-center
status: recorded
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# E-004 · R2 关门与 Root 投影

## 事实（2026-09-19）

### 1. R2 纲领阶段完成

Root `GOAL-001` 的纲领检查点 **R2（通用作业读面）** 由子目标 `GOAL-003-r2-generic-job-read-surface` 交付并已于同日关门：

| 项 | 值 |
|----|-----|
| 子目标状态 | `active` → **`done · 4/4`** |
| 检查点 | C1 模块与权限接线 / C2 查询与索引 / C3 读面 API 与作用域 / C4 R2 审计与投影 —— **全部完成** |
| 开放 required | **0** |
| 审计链 | `A-001` self `pass` → `A-002` independent（grok-build · grok-4.6 · high · `/audit`）**`pass`** → `A-003` 响应（6 条 recommended 全 `fixed`） |
| 关门依据 | 4 检查点全达成 + 开放 required = 0 + 无 P-004 待裁冲突 |

> 子目标关门属**非关键决策**，按用户目标轮次指令「子目标关门等非关键决策可经交叉审计后静默执行」——已经 `cross` 模式审计后执行，未额外请求用户裁决。

### 2. Root R2 检查点投影

| Root 检查点 | 投影前 | 投影后 |
|-------------|--------|--------|
| R2 通用作业读面 | `[ ]` | **`[x]`** |

Root `progress: 1/5` → **`2/5`**（`00-meta.md` frontmatter 与 `goal-tree.md` 树/表已同步）。

### 3. R2 交付结论（可核对的实现事实）

| 项 | 内容 |
|----|------|
| 新模块 | `admin.jobs`（进 admin 默认集；Profile 内容扩展） |
| 新权限 | `jobs.read`（`PolicyAdmin`）；写权限 `jobs.write` 归 R4 |
| 新路由 | `GET /api/jobs`、`GET /api/jobs/{id}`、`GET /api/jobs/{id}/result` |
| 作用域 | **管理作用域**（跨 actor）；既有 `GetForActor` actor 隔离**未放宽**（独立审计以空 diff 核实） |
| 查询 | `ListJobs`（kind/status/actor/时间过滤 + 排序白名单 + COUNT 同 WHERE + 分页）、`GetJob` |
| 迁移 | v72 `jobs_management_indexes`（`idx_jobs_created_at ON jobs(created_at DESC, id DESC)`）；`async_jobs`(42) 行**字节不变** |
| 跨 VP 触碰 | `wallet.go` 的结果地址改经共享 `jobs.ResultURL`；**输出字符串逐字不变**（独立审计对任意 id 复算成立） |
| R-1 修复 | `jobRuntime.enabled` 改为 `admin.jobs` 或 `admin.wallet` 任一存在即启用；**带变异测试验证的回归守卫** |
| 回归 | Go `go test ./...` 全绿；web `1440/1440`；`tsc -b` exit 0 |

### 4. 同步文件

| 文件 | 变更 |
|------|------|
| `GOAL-001/00-meta.md` | R2 检查点勾选；`progress: 2/5`；`version: 0.3.0`；子目标表更新为 `done`；备注修订 |
| `GOAL-001/02-execution.md` | 索引补 E-004 |
| `goal-tree.md` | 树/纲领路线图/状态表/说明节同步（Root `2/5`、GOAL-003 `done · 4/4`） |
| `workspace.md` | Root 投影行与纲领阶段表更新；补 R2 关门记录 |

### 5. 移交 R3

`D-001` §5 的 **T-5**（异步批量导出端点 + Job kind + 细粒度 `reporter.Progress`）与 **T-6**（前端触发机制 O-1 与 capability 口径 O-2 的**实现**）随 R2 关门移交 R3 子目标（按 P-001 立项）。O-1/O-2 的**方案**已在 R2 冻结（`D-001` §3）。

### 6. 边界

- **未**触碰任何 pinned 协议工件（`docs/schemas/**`、`apps/web/src/protocol/upstream/**`）——独立审计已核实五个 checkpoint 的 `--stat`。
- **未**把同步 `batch-delete` 改成异步；**未**重开 VP-012（对 `wallet.go` 的改动为字节等价重构，独立审计复核通过）。
- **未**声明 `jobs.write`（写操作归 R4）。
- **未**在愿景层改判 `V-F126`（承接动作已完成，闭合登记属 `/vision`；交接项见 R1 首波矩阵 §7）。

### 7. Git checkpoint

见提交记录（R2 关门 + 审计响应 + Root 投影），覆盖路径 `docs/workspaces/workspace-038-batch-operations-and-job-center/`。
