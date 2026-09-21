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
| 6 | **基础设施与范围保持**：未实现/解除 Redis、外部队列、多实例、跨进程索引或专用搜索引擎；未重开 VP-012/011/037；未改 Charter 目的/边界/非目标；未混入新业务域 | **达成** | `git diff --stat e1893a1a..HEAD -- docs/schemas apps/web/src/protocol/upstream apps/api/modules/jobs/migration apps/api/internal/jobs/repository.go apps/api/internal/jobs/runner.go` **为空**（R5 独立腿复验一致）；`async_jobs` v42 checksum 仍为 `55e1d3f8…`（`internal/store/migrate_test.go` frozen catalog，独立腿复跑 `TestCompiledMigrationCatalogOwnership` + `modules/jobs/migration` 通过）；Jobs 六态合同与同步 `batch-delete` 零改动（多次以空 diff 复核）；无 Redis/broker/多实例代码；未新开 VP。**口径精确化（R5 独立腿 F-004 附带指出）**：自 VP-038 激活（`e125d902`）起 `apps/api/modules/jobs/migration/**` **有** R2 授权的 v72 索引增量，故「migration 全区间零改动」的说法不成立；**准确表述**是「相对 R3 关门基线 `e1893a1a` 为空，且 v42 描述符/checksum 与 v72 之前的历史未动」 |
| 7 | **证据与审计**：退出矩阵、浏览器/自动化回归、必要的独立意见已落盘；开放 required = 0；组合投影同步且**经用户书面确认**关门 | **进行中** | 本文件（退出矩阵 + 回归，见 §2）；`GOAL-005` cross 审计两腿 `pass` 且 8/8 recommended `fixed`；R5 独立意见与用户确认见后续条目（C3/C4） |

### 2. 浏览器/自动化回归（C2）

| 项 | 结果 |
|----|------|
| 命令 | `cd apps/web && npm run test:e2e`（Playwright chromium · 单 worker · 默认 sqlite scratch 库） |
| **profile 口径（R5 独立腿 F-002 纠正）** | 该命令**未设 `APP_PROFILE`**，`playwright.config.ts` 默认 **`mvp`**；`admin.jobs` 不在 mvp/demo 默认集，故本 VP 新增面在**默认 e2e 中模块级缺席**。此前本表「在 admin profile 上跑通」的措辞与所记录命令不符，已按事实更正；4 个 skip 与 mvp 条件跳过一致（`localization` 的 admin-only 段 + `telegram-operator-layout` 三例按频道路由条件跳过） |
| 首次运行 | **15 passed / 1 failed / 4 skipped**；失败项 = `force-password-change.spec.ts`（fresh seed 强制改密） |
| 失败根因（证据） | 失败页面快照显示 `admin / admin` 登录返回 **"invalid username or password"**，说明该次运行所连 scratch 库中 admin 密码已被改过。`workers: 1` + `fullyParallel: false` 下用例按文件名顺序串行，而 **VP-036 新增的 `command-palette.spec.ts` 排在 `force-password-change.spec.ts` 之前**，其 sign-in helper **会自行完成强制改密**，从而消费掉「fresh seed」前提。 |
| 隔离复验 | `npx playwright test force-password-change.spec.ts` → **1 passed (10.3s)**（fresh 库下前提成立）→ 确认为**共用库 + 文件顺序假设被后续波次破坏**的既有挂具缺陷，与 R1～R5 实现无关 |
| 修复 | 按 R5 发现处置：将该用例重命名为 `00-force-password-change.spec.ts` 使其**真的**先跑（文件内注释记录原因与两向证据），不改变任何断言 |
| 修复后复跑 | **`16 passed / 4 skipped / 0 failed`（exit 0，3.9m）** —— 含 `00-force-password-change.spec.ts`（1.9s，先跑）与全部其余 spec；4 个 skip 为既有的 profile/配置条件跳过（`localization` 的 admin-only 段 + `telegram-operator-layout` 三例按频道路由条件跳过） |
| **补测（用户 2026-09-19 指令：先补 jobs e2e 再关门）** | 新增 `e2e/jobs-result-center.spec.ts`：真实浏览器路径 **选择行 → 提交异步批量导出 → 观察进度 → 取回 CSV → 到结果中心读同一作业（本地化终态）→ 用行操作再次下载**；断言下载文件名 = 服务端 `users-selection.csv`（两次下载同源）。以 admin profile 运行（`admin.jobs` 仅在该 profile 挂载；spec 内 `test.skip(appProfile !== "admin")` 保证 mvp 下显式跳过而非假绿） |
| **双 profile 全量复跑** | 默认（mvp）：**16 passed / 5 skipped / 0 failed**（第 5 个 skip = 新 spec 的 profile 守卫）；`APP_PROFILE=admin`：**17 passed / 4 skipped / 0 failed**（新 spec 实跑通过，20.8s 单跑）。两条命令均 exit 0 |

### 3. e2e 覆盖与判据 4/5 的关系（诚实边界）

- **补测前**（R5 审计 A-001/A-002 F-002 登记的缺口）：11 个 spec 无「提交批量导出 → 观察进度 → 下载」路径，且默认 `APP_PROFILE=mvp` 不含 `admin.jobs` → 本 VP 新增面在默认 e2e 中模块级缺席。
- **补测后（2026-09-19，用户指令）**：新增 `e2e/jobs-result-center.spec.ts`（admin profile 专用，mvp 下显式 skip），把该路径变成**真实浏览器端到端**证据：选择行 → `导出所选` → **观察真实进度**（`[data-jobs-batch-export-progress]`）→ 终态下载（断言文件名 = 服务端 `users-selection.csv`）→ 结果中心列出同一作业且状态为**本地化**终态（`Succeeded`）→ 行操作再次下载同一文件。
- 判据 4/5 的浏览器侧证据因此由「间接（交互级 + 契约测试）」升级为「**端到端 + 交互级 + 契约测试**」三层；`A-002` F-002 的缺口闭合（`A-004` 登记）。默认 mvp 运行下该 spec 以 profile 守卫跳过——这是**显式跳过**，不会伪装成本 VP 已被 mvp 覆盖。
- 仍未覆盖（保持登记）：并行/重试/取消的浏览器端到端路径（当前由交互级与 HTTP 契约测试覆盖）；`GOAL-006 A-001/A-002 F-001`（e2e fresh-seed 顺序契约）仍为 roadmap 上的 bounded residual。

### 4. 边界

- 本条目只登记**已发生**的核对与命令结果；未达成或有界项在 §1 判据 7、§3 显式标注。
- 未改任何实现代码（除 §2 的 e2e 文件重命名 + 注释）；未触碰 pinned 工件。
