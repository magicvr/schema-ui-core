---
title: 审计台账 · 生产级可用 Admin 基架
status: active
created: 2026-08-01
updated: 2026-08-04
parent: null
version: 0.4.0
---

# 审计台账 · GOAL-001

## 正式意见索引

| 编号 | source | 日期 | scope | verdict | 状态 |
|------|--------|------|-------|---------|------|
| A-001 | self | 2026-08-02 | R1 · 协议实施边界与 Schema Renderer 产品化 | pass | 已出具；无开放 R1 required finding |
| A-002 | independent | 2026-08-03 | apps/api + apps/web · VP-002 功能实现与产品意图交叉审计 | fail | 已响应（2026-08-03）；**F-002-001/002/003 全部 `fixed`**（F-002-001 于 2026-08-04 经 GOAL-010 关闭）；F-002-004~006 recommended 非阻断 |

## A-001 · Root R1 阶段自审（2026-08-02）

- **source**：self
- **auditor**：GitHub Copilot · `/govern`
- **类型 / scope**：stage；R1 · 协议实施边界与 Schema Renderer 产品化，包括 R1 三个已关门子目标的阶段退出证据、`I-001` 与 Root 路线图检查点。
- **verdict**：pass
- **audit_type**：execution-facts

### 范围与区间

- 当前工作区为 `workspace-002-production-admin-foundation`，canonical root、Root、`goal-tree.md` 与 VP-002 绑定一致；`shared_materials_catalog: none`，未使用共享资料作为本意见事实或关闭证据。
- 本意见只核对 R1；真实认证（R2）、持久化权限（R3）、Schema CRUD（R4）和工程化关门（R5）不在本阶段审计范围。

### 成果（有证据）

- `I-001` 已由 D-004 的实现差量矩阵验证并冻结 R1 方案边界。
- `GOAL-002-r1-schema-load-validate` 已关门：加载、结构校验和统一错误面由 A-001 independent 与 A-002 self 关门审计确认，均无开放 required finding。
- `GOAL-003-r1-default-render-path` 已关门：默认 `schemaUrl` → 加载、校验、`RenderPage` 主路径和手写示例迁移由 A-003 self 与 A-004 independent 关门审计确认，均为 pass。
- `GOAL-004-r1-representative-node-pages` 已关门：代表性列表、表单、组合 Node 页面以及成功/失败路径回归由 A-001 self 与 A-002 independent 关门审计确认，均为 pass。
- 两份 2026-08-02 独立审计均记录 Web `425/425` 测试、Web 生产构建、Go `test` 和 `vet` 成功。该命令结果为既有阶段证据，本次未重新执行。

### 对照 R1 退出证据

| 检查项 | 状态 | 证据 |
|--------|------|------|
| R1 方案边界 | 通过 | `I-001` = verified；D-004 实现差量矩阵 |
| Schema 加载、校验与错误面 | 通过 | GOAL-002 A-001 / A-002；目标已 done |
| 默认 Schema Renderer 主路径 | 通过 | GOAL-003 A-003 / A-004；目标已 done |
| 代表性 Node 页面与回归 | 通过 | GOAL-004 A-001 / A-002；目标已 done |
| R1 Root 检查点 | 通过 | 三个 R1 子目标均为 done；`00-meta.md` 与 `goal-tree.md` 均为 `1/5` |

### Findings

- 无开放 R1 required finding。
- GOAL-003 的 Schema 行操作应用级断言与 GOAL-004 的 `recordView` 真实数据联动仍为 recommended R4 follow-up；它们不扩大或否定 R1 的已核对范围，不阻断 R1 阶段结论。

### 必改项汇总

- 无开放 required finding。
- `I-002` 为 R2 的 required 信息项，现为 `collecting`；它不追溯阻断 R1，但阻断 R2 的方案冻结与实施。

### 结论 + 建议下一步

- R1 阶段退出证据完整，本阶段 self audit verdict = **pass**；Root 保持 `active`，路线图检查点保持 `1/5`。
- 下一步仅进行 D-006 定义的 `I-002` 信息收集；在形成可核对认证方案并获得后续决策前，不冻结或实施 R2。

## 当前审计边界

- A-001（self）覆盖 R1 阶段退出证据，verdict = pass；无开放 R1 required finding。
- VP-002 的 Vision Review 只覆盖愿景与组合边界，不替代本 Root 的 Goal Audit。
- 后续 self / independent 意见从 `A-003` 起共用序列，required finding 只能按 `fixed`、`accepted-residual` 或 `user-overruled` 合法闭合。

## A-002 · apps/api + apps/web VP-002 独立功能审计（2026-08-03）

- **source**：independent
- **auditor**：Codex（独立功能审计；未调用 `audit` skill）
- **类型 / scope**：execution-facts + product-fit；核对 `apps/api`、`apps/web` 当前代码相对 VP-002 产品级成功标准（完整 Schema Renderer、真实认证、最小权限、Schema CRUD、工程化交付）的实现事实、跨层链路和缺陷。不修改目标状态、进度、方案正文或其他目标台账。
- **verdict**：**fail**。当前实现可证明“认证 + 持久化 RBAC 种子/records 代表性 CRUD + 容器/浏览器基线”这一交付切片，但不能证明 VP-002 所要求的可直接接业务的通用 Schema Admin 基架；同时存在会绕过前端门禁或错误保持认证状态的功能缺陷。

### 已确认符合的范围

- `apps/api` 已提供真实登录、refresh 轮换、登出、请求级身份与 `401/403` 权限门禁；生产 `AUTH_DEV_SESSION_ENABLED=true` 会被启动守卫拒绝。
- SQLite 迁移/种子已创建用户、角色、权限、菜单及关系表，records API 具备持久化 list/search/detail/create/update/delete 闭环；Schema fixture 能被 API 只读加载。当前没有用户/角色/权限/菜单管理 API，这在 VP-002 允许“初期只靠种子机制”的范围内不单独判为 required，但限制了后续业务运营能力。
- `apps/web` 默认入口将 `authFetch` 注入 Schema 页面和 records transport；`npm test -- --run` 为 23 个文件 / 458 个测试通过，`npm run build` 通过；Chromium E2E 2/2 通过（真实登录/反向代理与 SQLite records CRUD）。这些证据不扩展为通用业务实体或 VP 关门证据。

### Required findings

#### F-002-001 · Schema Renderer/CRUD 仍硬编码单一 records 实体（required / high）

- `apps/web/src/renderer/schema-table.tsx:81-132` 无论 `dataSource` 值为何都调用 `fetchRecords`；`apps/web/src/renderer/records.ts:6-19,78-88,140-155` 强制 `id/name/status/owner/updatedAt` 和固定分页形状。`apps/api/internal/handler/health.go:25-30` 仅注册 `/api/records` CRUD 与只读 `/api/schema/{pageId}`，没有按 Schema/资源注册的通用 CRUD 或 Schema 持久化入口。
- 影响：当前 fixture 的 `list-edit-lifecycle` 代表页可运行，但通过修改 Schema 无法在不改 Renderer 主路径的前提下接入任意业务实体；自定义表格数据会因固定 RecordList/RecordItem 解析失败。这不满足 VP-002 §产品级成功标准 1、4、6 及阶段 3 的业务接入意图。
- **建议关闭路径**：将表格/表单 transport、字段模型和 response mapping 提升为 Schema 驱动的通用适配层，并为业务资源提供明确后端契约；或在 VP/Root 明确把本阶段降级为单一代表性 records 示例并重新获得用户裁决。当前不能将现状宣称为完整生产级基架。

#### F-002-002 · 表单 Schema/reaction 错误只展示、不阻断提交（required / high）

- `apps/web/src/renderer/render.tsx:509-510` 计算 `gate.errors`，`540-557` 的 `handleSubmit` 不检查该错误或 `reaction.errors`，`565-587` 仅渲染提示；`594-600` 的提交按钮只由 `submitting` 控制。
- 影响：无效/不兼容的表单字段或 reaction 仍可触发 POST/PATCH，违反 VP-002 阶段 3 的输入校验与“无效 Schema 确定性拒绝”边界。
- **建议关闭路径**：任何 gate/reaction 错误时禁用提交并在 handler 层再次拒绝；增加“错误显示后请求未发出”的回归测试。

#### F-002-003 · refresh 重试仍返回 401 时认证状态不丢失（required / medium）

- `apps/web/src/account/auth-client.ts:117-128` 在首次 401 后 refresh 成功只重试一次；重试仍为 401 时不 `clearTokens()`，也不触发 `onAuthLost`。`152-164` 中 `login` 的 `/me` 失败被吞掉并返回带 token 的登录快照与空 `features`，`AuthContext` 随即进入 `authenticated`。
- 影响：撤销/过期或 `/me` 故障可能留下“已登录但所有请求持续失败”的状态，不能可靠满足 VP-002 真实认证失效与错误处理标准（§2、阶段 2）。
- **建议关闭路径**：二次 401 必须清理 token 并通知 auth-lost；login 后 `/me` 失败应回滚 token 并以登录失败呈现，而不是静默认证降级；补充二次 401 和 `/me` 失败测试。

### Recommended findings

#### F-002-004 · 生产登录页无条件展示 `admin / admin`（recommended / medium）

- `apps/web/src/app/LoginPage.tsx:94-96` 固定渲染 `Local development seed: admin / admin`，无环境开关；Compose 生产路径要求外部 `ADMIN_INITIAL_PASSWORD`，实际密码不必为 `admin`。
- 影响：生产用户会得到错误凭据提示并形成默认密码/开发会话仍可用的错误安全信号，违反 VP-002 对生产身份边界与可配置启动的要求。
- **建议关闭路径**：仅在显式 development 配置下显示本地 seed 提示；生产构建隐藏该文案，文档改为读取部署时配置的管理员密码。

#### F-002-005 · 生产 JWT secret 仅做非空校验（recommended / medium）

- `apps/api/cmd/server/main.go:101-110` 的 `resolveJWTSecret` 在非 development 只判断 `AUTH_JWT_SECRET != ""`；`apps/api/internal/config/config.go:51-60` 未施加长度/熵/格式门禁。
- 影响：生产可用极短、可猜 secret 启动，HS256 会话安全性完全依赖部署者自律，不符合 VP-002 “可替换认证实现与安全配置”的 production-grade 边界。
- **建议关闭路径**：在 production/staging 对 secret 施加可核对的最小长度/熵规则并增加反例测试；开发环境可保留显式低门槛。

#### F-002-006 · `/healthz` 被 Compose 当作 readiness，但不检查 SQLite（recommended / medium）

- `apps/api/internal/handler/health.go:33-40` 始终返回 200，不访问 Store；`compose.yaml:30-35` 以该端点作为 `service_healthy` 条件。
- 影响：进程存活但数据库损坏、只读或迁移后不可用时，Compose 仍会启动 Web 并报告 API healthy；这是运行诊断语义不准确，当前不否定已有启动/smoke 证据。
- **建议**：区分 liveness/readiness，或让 readiness 执行轻量 SQLite 查询/迁移状态检查并补故障注入测试。

### 证据与边界

- 本轮重跑：`apps/api` `go test ./... -count=1`、`go vet ./...`；`apps/web` `npm test -- --run`（23 files / 458 tests）、`npm run build`；`WEB_PORT=9999 npm run test:e2e`（Chromium 2 passed）。
- 未执行 Docker Compose 故障注入、生产部署平台验证或完整 fork 计时；因此未把本地测试扩写为生产发布验收。
- 本意见仅追加独立审计记录，不修改 Root `status: active`、`progress: 5/5`、`goal-tree.md`、VP 状态或其他目标；required finding 的响应与是否关闭/重定义范围由 `/govern`/用户裁决处理。

## A-002 响应（编排器 · 2026-08-03）

- **响应编号**：A-002（independent · fail · apps/api + apps/web product-fit）
- **响应来源**：`/govern` 编排器（self 侧记录，非独立审）
- **用户裁决**（P-004，决策 [D-014](01-decision.md)）：三条 required 全部走 `fixed` 关闭路径；F-002-001 → 通用适配层改造（`GOAL-010-a002-schema-adapter`）；F-002-002/003 → 缺陷修复（`GOAL-009-a002-auth-form-fixes`）；recommended F-002-004~006 → GOAL-009 可选加分（是否实施待用户决定）；A-002 同 scope self 审计延后至修复完成后随关门补（P-004 §3.1）。

### 关闭证据表

| finding | 状态 | 证据路径 |
|---------|------|----------|
| F-002-001（Renderer 硬编码 records 实体） | **fixed（2026-08-04）** | [GOAL-010-a002-schema-adapter](../GOAL-010-a002-schema-adapter/00-meta.md) S1～S5（`done / 5/5`）：通用资源契约 + 后端注册表/工厂 + 前端泛化 + users/roles 双实体 Schema-only 接入（Renderer 零 diff）+ records 0006 退场 + S5 全量回归（go test 7 包全绿 + vet、vitest 491/491、tsc/build、e2e 2/2）+ [GOAL-010 A-002 self close-out](../GOAL-010-a002-schema-adapter/03-audit.md#a-002--s5-关门审计goal-010-全目标-close-out--root-f-002-001-关闭证据2026-08-04)（pass） |
| F-002-002（表单校验错误不阻断提交） | **fixed（2026-08-03）** | [GOAL-009 S1](../GOAL-009-a002-auth-form-fixes/02-execution.md)（`render.tsx` 门禁 + 3 回归）+ [A-001 self close-out](../GOAL-009-a002-auth-form-fixes/03-audit.md)（pass） |
| F-002-003（认证失效状态不清理） | **fixed（2026-08-03）** | [GOAL-009 S2](../GOAL-009-a002-auth-form-fixes/02-execution.md)（`auth-client.ts` 清理 + 3 回归）+ [A-001 self close-out](../GOAL-009-a002-auth-form-fixes/03-audit.md)（pass） |
| F-002-004~006（recommended） | **已实施（2026-08-03 · GOAL-009 S5 加分，非 required 闭合）** | [GOAL-009 S5](../GOAL-009-a002-auth-form-fixes/02-execution.md)：seed 文案环境门控、生产 JWT secret 校验、`/readyz` readiness |

### 仍开放项

- ~~F-002-001 仍 open~~：已按 `fixed` 合法闭合（2026-08-04，GOAL-010 `done / 5/5` + self close-out A-002 pass，见关闭证据表）。**A-002 三条 required 全部合法闭合**，Root 关门与 VP-002 关门阻断解除（进入独立用户裁决流程）。
- F-002-004~006（recommended）非阻断；F-002-001 关闭证据的 `/audit` finding-closure 独立复审为可选加固（未执行）。
- Root `status: active`、派生进度 `5/5` 不变。

### 后续

- F-002-002/003 已于 2026-08-03 按 `fixed` 合法闭合（GOAL-009 S1/S2 + A-001 self close-out）。
- F-002-001 已于 2026-08-04 按 `fixed` 合法闭合（GOAL-010 S1～S5 + A-002 self close-out，见关闭证据表）。
- A-002 全部 required 已合法闭合；Root close-out 关门审计与 VP-002 关门为独立用户裁决。
