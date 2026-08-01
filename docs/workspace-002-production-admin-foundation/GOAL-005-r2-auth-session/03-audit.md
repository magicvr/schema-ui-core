---
title: 审计台账 · R2 · 真实认证与请求级身份
status: active
created: 2026-08-02
updated: 2026-08-02
parent: null
version: 0.3.0
---

# 审计台账 · GOAL-005

## 正式意见索引

| 编号 | source | 日期 | scope | verdict | 状态 |
|------|--------|------|-------|---------|------|
| A-001 | independent | 2026-08-02 | close-out | conditional | F-001 **fixed**（2026-08-02 编排器响应，证据见下） |

## 当前审计边界

- 已记录独立 close-out 审计 A-001 及其编排器响应；目标的 status/progress 与关门仍由 `/govern` 维护（用户确认后置 `done`）。
- required finding 只能按 `fixed`、`accepted-residual` 或 `user-overruled` 合法闭合；`F-001` 已按 `fixed` 闭合（2026-08-02）。

## A-001 · R2 真实认证与请求级身份关门审计（2026-08-02）

- **source**：independent
- **auditor**：GitHub Copilot
- **类型** / **scope**：close-out；`GOAL-005-r2-auth-session` 的成功标准、I-005 信息门禁、实现事实与可复现验证
- **verdict**：conditional

### 范围与区间

- 工作区绑定已核对：`workspace-002-production-admin-foundation` 的 `root_goal`、canonical root、`plan_refs` / `primary_plan` 与 Root、`VP-002-production-admin-foundation`、`schema-ui-core-admin-foundation@0.1.0` 对齐；`shared_materials_catalog: none`，本审计未将外部或跨工作区资料作为关闭证据。
- 本审计仅评估 R2 关门，不评估 R3 的持久化用户—角色—菜单模型、R5 的跨源/部署基线或 VP/Root 关门。
- 已逐项核对 `I-005-001` 至 `I-005-005`：三项 required 均为 `verified`，两项 non-blocking 也已记录范围决策；没有到期的本目标 required 信息项。

### 成果（有证据）

- 后端认证路径存在明确实现：登录、刷新轮换、登出撤销、JWT Bearer 解析与 SQLite 哈希刷新令牌见 `apps/api/internal/auth/auth.go`、`apps/api/internal/handler/auth.go`；`/api/accounts/me` 与 records 写路由经请求级身份中间件，未认证为 `401`、非 admin 为 `403`，见 `apps/api/internal/handler/account.go`、`apps/api/internal/handler/records.go`。
- 静态开发会话由 `AUTH_DEV_SESSION_ENABLED` 显式控制且默认 false，配置见 `apps/api/internal/config/config.go`；前端 access 内存、refresh localStorage、刷新重试及会话丢失处理见 `apps/web/src/account/auth-client.ts`。
- 独立复跑结果：`apps/api` 的 `go test ./...` 通过；`apps/web` 的 `npm test` 为 22 个测试文件、441 项通过；`npm run build` 通过。
- CI 定义了 Linux 的 API、Web 和 Chromium browser E2E job，认证浏览器脚本覆盖未认证登录门禁、登录、`/me` 与 `/api/records` 代理链，见 `.github/workflows/r6-basic-matrix.yml`、`apps/web/e2e/shell.spec.ts`。

### 对照成功标准

- 登录/登出、刷新/撤销、请求级身份、SQLite 与依赖、前端认证闭环、显式 dev 兜底均有对应实现与单测或构建证据；本审计未发现把 R3 身份持久化模型误记为 R2 已完成的范围越界。
- D-002 对 refresh token localStorage 的 XSS 残余已记录为用户书面取舍及缓解边界，未作为未声明风险处理。

### Findings

#### F-001 · 浏览器端到端通过结果缺失，M14 的 CI 可复现主张尚不能关门

- **级别**：required
- **严重度**：medium
- **证据**：本次独立执行 `apps/web` 的 `npm run test:e2e` 失败于 Playwright `webServer` 启动 Vite：`listen EACCES: permission denied 127.0.0.1:5173`。`02-execution.md` 也只记录了本机绑定限制和“CI 将覆盖”，而未提供已通过的 CI run、日志或制品链接。
- **影响**：M14 要求 login → me → 401/403 路径具备自动化测试且 CI 可复现。现有 CI 配置与脚本证明测试意图和 Linux 执行路径存在，但不能替代一次可复核的 browser E2E 通过结果；当前 E2E 脚本也未直接断言 401/403。
- **关闭条件建议**：优先在 Linux CI 或可绑定回环端口的环境跑通并留存 browser E2E 结果，同时把匿名 `401` 与非 admin `403` 加入该 E2E 或明确由已通过的 API 自动化测试承担并在关门响应中指向证据。若要接受本机平台残余，需用户书面选择 `accepted-residual`，注明仅限何种环境、风险范围与下一次 CI 复核触发条件。
- **状态**：open

### 必改项汇总

- `F-001`（required）：未关闭前不得将本目标置为 `done` 或勾选 Root R2 检查点。

### 与既有意见的异同

- 此为本目标首条正式意见；原台账“刚立项”描述已不再适用于当前已记录的实施事实。

### 结论 + 建议给编排器/用户的下一步

- 结论为 `conditional`：R2 实现、信息门禁和单元/构建验证总体具备关门基础，但 `F-001` 阻断无条件关门。
- 由 `/govern` 汇总本意见并请用户选择修复、接受有界 residual 或推翻该 finding；独立审计不修改 `status`、`progress`、成功标准或 `goal-tree`。

### 声明

本意见不修改 status/progress；响应由 `/govern` 处理。

## 响应 · A-001 / F-001（2026-08-02 · /govern）

- **source**：orchestrator（编排器响应；**不**冒充 `source: independent`）
- **用户裁决**：选择 `fixed` 路径（在 Linux CI 跑通 browser E2E 并补充 401/403 证据后复审关门），不接受本机平台 `accepted-residual`。
- **证据**：
  1. **browser E2E 增强**：`apps/web/e2e/shell.spec.ts` 追加匿名 `401` 断言（`GET /api/accounts/me` → `401 UNAUTHENTICATED`；匿名 `PATCH /api/records/rec-1` → `401`）与 admin 写门禁 `200` 对照（提交 `32d8486`）。
  2. **Linux CI 通过**：push `dev` 触发 run [#30711903555](https://github.com/magicvr/schema-ui-core/actions/runs/30711903555)（2026-08-01T18:09Z = 本地 2026-08-02）：
     - `browser E2E (Linux, Node 22)` → **success**；Playwright `1 passed (31.8s)`（`e2e/shell.spec.ts:9:1`，含新增 401/200 断言）
     - `web (Linux, Node 22)` → **success**（`npm test` 441 passed、`npm run build`）
     - `api (Linux, Go 1.26)` → **success**（`go test ./...` 全绿）
  3. **非 admin 403 证据**（由 API 自动化测试承担，分配决策 D-007）：`apps/api/internal/handler/records_test.go::TestRecordsWriteDeniedWithoutAdminRole`（editor 角色 PATCH/DELETE → `403 FORBIDDEN`）与 `TestRecordsWriteRequiresAuth`（匿名写 → `401`）、`account_test.go`（`/me` 匿名 → `401`）。
  4. **本地平台受限记录**：本机 Windows `127.0.0.1:5173` EACCES 依旧；以临时端口 `9999` 复跑脚本 `1 passed` 后还原，未进入提交；Linux CI 为权威可复现证据。
- **闭合**：`F-001` → **`fixed`**（可核对：CI run 结论 + 日志 `1 passed (31.8s)` + `shell.spec.ts` 断言源码 + `records_test.go` 403 测试）。
- **仍开放**：无。GOAL-005 满足 close-out 条件，进入关门复审（是否自审、置 `done` 与勾选 Root R2 检查点由用户裁决）。
