---
id: E-001-r5-exit-matrix
doc: execution-entry
parent: GOAL-006-r5-evidence-and-closeout
status: recorded
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# E-001 · R5 退出矩阵与浏览器回归

## 事实（2026-09-19）

### 1. VP-038 方向级退出判据逐条核对

> 判据原文见 `docs/vision/plans/VP-038-batch-operations-and-job-center.md` §方向级退出判据。**每条只引可核对证据**（工作区台账条目 / 命令与结果 / 代码位置）；无法核对的一律标「未达成」或「有界」，不写叙事结论。

| # | 判据 | 结论 | 证据 |
|---|------|------|------|
| 1 | **分母与契约冻结**：首波 Job 种类、可见作用域、六态与进度/结果/过期对外语义、批量异步契约（请求形状 / 202+jobId / 读取路径）形成**可机器核对**的矩阵；排除项显式点名 | **达成** | `GOAL-002` `D-001`（C2 = 方案 B 本地异步契约；C3 首波仅「批量导出所选」；C1 管理作用域 + `jobs.read`；K-1～K-7；R-1）；矩阵 `GOAL-002/attachments/r1-job-kind-scope-matrix.md`、`r1-first-wave-denominator-matrix.md`；排除项见该 `D-001` §5 与 Root `00-meta` 非目标 |
| 2 | **通用作业可见性**：按权限列出并查看详情（状态、进度、attempt/max_attempts、错误码与消息、时间戳、correlation）；无权限/越作用域/不存在 fail-closed **不泄露他人作业** | **达成** | `GET /api/jobs`、`/{id}`、`/{id}/result`（`jobs.read` = `PolicyAdmin`）；投影含 `attempt`/`maxAttempts`/`errorCode`/`errorMessage`/`createdAt`/`updatedAt`/`finishedAt`/`resultExpiresAt`/`correlationId`；测试 `internal/handler/jobs_test.go`（管理作用域列表、404 不泄露、409/410 三段语义）+ R2 独立审计 `GOAL-003 A-002`（`pass`，含「actor 隔离未被放宽」复证） |
| 3 | **批量操作异步承接**：至少一条真实批量操作走 Job（queued→running→终态 + 进度 + 终态结果），**同时**既有 ADR-0022 同步 `batch-delete` 语义与回归不退化 | **达成** | `jobs.batch-export`（`POST /api/jobs/batch-export` → 202 + jobId，双重门禁 `jobs.write`+`data.export`，真实细粒度进度按行上报）；行为测试 `modules/jobs/export_test.go`（CSV 结果、进度非硬编码、校验先于建行）；同步路径零 diff + 回归锚点（`GOAL-004 A-002` 独立腿以**空 diff** 复核 `resources.go`/协议 fixture 未动） |
| 4 | **结果中心体验**：进行中/成功/失败/取消/过期/结果已过期六类呈现；导出类结果可下载；失败可重试；中英文、浅色/深色、加载/空态/错误态与既有约定一致 | **达成**（三条 low 级残余已跨区修复） | `jobs` 页：六态徽标 + 六态**逐值本地化**（`valueLabels`，en-US/zh-CN 双断言）、进度列、`recordView` 15 字段详情、结果下载（`jobs.downloadResult` → 共享助手，CSV 用服务端 `fileName`）、取消/重试（`jobs.write`，行级可用性由服务端派生字段驱动）、自动刷新（off/5/10/30s，**空闲不轮询**）；测试 `renderer/jobs-result-center.test.tsx`（13 例）、`components/jobs-batch-export.test.tsx`（4 例）、`lib/job-result-download.test.ts`（4 例）、`components/jobs-auto-refresh.test.tsx`（4 例）、`renderer/table-refresh-seam.test.tsx`（3 例）；残余由 `[workspace-010] GOAL-044`（`done · 4/4`）交付并回填 `fixed` |
| 5 | **权限与 Profile 安全**：按 Profile 与权限过滤，不绕过路由守卫与数据范围；mvp/admin（及适用 demo/custom）覆盖**有矩阵证据**；直接 URL 行为与既有守卫一致 | **达成** | `admin.jobs` 仅进 admin 默认集（用户 P-004 裁决，`kernel/profile.go`）；只读组合下两条写路由**不声明也不挂载**（`TestJobActionRoutesAbsentWithoutActions`）；写门禁以**自定义只读角色真反例**钉住（`TestJobWriteGateRequiresJobsWriteNotJobsRead`：持 `jobs.read` 无 `jobs.write` → 读 200 / 写 403 且行未变）；`descriptor = 计划描述符`（`TestPlanDescriptorMatchesTheFullProvider` + kernel `MODULES_API_MISMATCH` 装配检查） |
| 6 | **基础设施与范围保持**：未实现/解除 Redis、外部队列、多实例、跨进程索引或专用搜索引擎；未重开 VP-012/011/037；未改 Charter 目的/边界/非目标；未混入新业务域 | **达成** | R1～R5 区间 `git diff --stat` 无 `docs/schemas/**`、`apps/web/src/protocol/upstream/**`、`apps/api/modules/jobs/migration/**`（迁移仍 v72，`async_jobs` v42 checksum `55e1d3f8…` 未变）；Jobs 六态合同与同步 `batch-delete` 零改动（多次以空 diff 复核）；无 Redis/broker/多实例代码；未新开 VP |
| 7 | **证据与审计**：退出矩阵、浏览器/自动化回归、必要的独立意见已落盘；开放 required = 0；组合投影同步且**经用户书面确认**关门 | **进行中** | 本文件（退出矩阵 + 回归，见 §2）；`GOAL-005` cross 审计两腿 `pass` 且 8/8 recommended `fixed`；R5 独立意见与用户确认见后续条目（C3/C4） |

### 2. 浏览器/自动化回归（C2）

| 项 | 结果 |
|----|------|
| 命令 | `cd apps/web && npm run test:e2e`（Playwright chromium · 默认 sqlite scratch 库 · 单 worker） |
| 首次运行 | **15 passed / 1 failed / 4 skipped**；失败项 = `force-password-change.spec.ts`（fresh seed 强制改密） |
| 失败根因（证据） | 失败页面快照显示 `admin / admin` 登录返回 **"invalid username or password"**，说明该次运行所连 scratch 库中 admin 密码已被改过。`workers: 1` + `fullyParallel: false` 下用例按文件名顺序串行，而 **VP-036 新增的 `command-palette.spec.ts` 排在 `force-password-change.spec.ts` 之前**，其 sign-in helper **会自行完成强制改密**，从而消费掉「fresh seed」前提。 |
| 隔离复验 | `npx playwright test force-password-change.spec.ts` → **1 passed (10.3s)**（fresh 库下前提成立）→ 确认为**共用库 + 文件顺序假设被后续波次破坏**的既有挂具缺陷，与 R1～R5 实现无关 |
| 修复 | 按 R5 发现处置：将该用例重命名为 `00-force-password-change.spec.ts` 使其**真的**先跑（文件内注释记录原因与两向证据），不改变任何断言 |
| 修复后复跑 | **`16 passed / 4 skipped / 0 failed`（exit 0，3.9m）** —— 含 `00-force-password-change.spec.ts`（1.9s，先跑）与全部其余 spec；4 个 skip 为既有的 profile/配置条件跳过（`localization` 的 admin-only 段 + `telegram-operator-layout` 三例按频道路由条件跳过） |

### 3. e2e 覆盖与判据 4/5 的关系（诚实边界）

- e2e 分母（11 个 spec 文件）覆盖：登录/强制改密、命令面板、host 失败面、列表视觉、本地化、schema 鉴权传输、schema CRUD、shell、Telegram 运营台布局、长内容抽查。**没有一个 spec 直接驱动 jobs 结果中心**（无「提交批量导出 → 观察进度 → 下载」的端到端用例）。
- 因此判据 4/5 的浏览器侧证据来自**渲染/交互级测试**（真实 `jobs.json` + 生产渲染链 + 真实路由断言的 24 例，见 §1 判据 4）与 HTTP 契约测试，而不是 e2e。
- **该缺口已登记**（R5 `A-001` F-002）：是否新增一条 jobs 结果中心的 e2e 用例，属「后续波次 or bounded residual」的选择，交由用户/审计裁量，不在本波次静默扩大范围。

### 4. 边界

- 本条目只登记**已发生**的核对与命令结果；未达成或有界项在 §1 判据 7、§3 显式标注。
- 未改任何实现代码（除 §2 的 e2e 文件重命名 + 注释）；未触碰 pinned 工件。
