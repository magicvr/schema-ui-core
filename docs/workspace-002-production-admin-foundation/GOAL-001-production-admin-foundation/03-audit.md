---
title: 审计台账 · 生产级可用 Admin 基架
status: active
created: 2026-08-01
updated: 2026-08-04
parent: null
version: 0.9.0
---

# 审计台账 · GOAL-001

## 正式意见索引

| 编号 | source | 日期 | scope | verdict | 状态 |
|------|--------|------|-------|---------|------|
| A-001 | self | 2026-08-02 | R1 · 协议实施边界与 Schema Renderer 产品化 | pass | 已出具；无开放 R1 required finding |
| A-002 | independent | 2026-08-03 | apps/api + apps/web · VP-002 功能实现与产品意图交叉审计 | fail | 已响应（2026-08-03）；**F-002-001/002/003 全部 `fixed`**（F-002-001 于 2026-08-04 经 GOAL-010 关闭）；F-002-004~006 recommended 非阻断 |
| A-003 | independent | 2026-08-04 | finding-closure · Root A-002 F-002-001 | pass | 已响应（2026-08-04）：采纳 pass；F-002-001 `fixed` 维持；R-001/R-002 handled |
| A-004 | self | 2026-08-04 | Root close-out · 全目标关门审计 | pass | 已出具；其后 A-005 新 required 使 Root 回退 `active` |
| A-005 | independent | 2026-08-04 | apps/api + apps/web · VP-002 产品意图独立复审（无 skill） | fail | 已响应：F-001 → **fixed**（GOAL-012）；R-001～R-003 recommended 非阻断 |
| A-006 | independent | 2026-08-04 | apps/api + apps/web · VP-002 产品意图再审（无 skill） | pass | 已响应：R-001～R-004 → **fixed**；R-005 → residual-by-design / handled；无 required |

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
- F-002-004~006（recommended）非阻断；F-002-001 关闭证据的 `/audit` finding-closure 独立复审见 **A-003（independent · pass，2026-08-04）**——`fixed` 维持，无新增 required。
- Root `status: active`、派生进度 `5/5` 不变。

### 后续

- F-002-002/003 已于 2026-08-03 按 `fixed` 合法闭合（GOAL-009 S1/S2 + A-001 self close-out）。
- F-002-001 已于 2026-08-04 按 `fixed` 合法闭合（GOAL-010 S1～S5 + A-002 self close-out，见关闭证据表）。
- A-002 全部 required 已合法闭合；Root close-out 关门审计与 VP-002 关门为独立用户裁决。

## A-003 · finding-closure · Root A-002 F-002-001（2026-08-04）

- **source**：independent
- **auditor**：Grok Build · `/audit`
- **类型 / scope**：finding-closure；仅复核 Root A-002 **F-002-001**（Schema Renderer/CRUD 硬编码单一 records 实体，required / high）的 `fixed` 关闭证据是否充分、可重复核对。不审 F-002-002/003（已由 GOAL-009 闭合）、recommended F-002-004~006、Root close-out 关门、VP-002 关门，也不复判 GOAL-010 全目标五件套是否可另开新 finding。
- **verdict**：**pass**

### 范围与区间

- 工作区：`workspace-002-production-admin-foundation`；`workspace.md` 的 Root（`GOAL-001-production-admin-foundation`）、canonical root、`primary_plan: VP-002-production-admin-foundation` 一致；`shared_materials_catalog: none`，未使用共享资料作为事实或关闭证据。
- 已读：本文件 A-002 原文与关闭证据表、[GOAL-010](../GOAL-010-a002-schema-adapter/) 五件套（含 A-002 self close-out）、GOAL-011 作为 S4 双实体载体的既有关闭事实；代码与测试静态/运行核对见下。
- **证据边界（本轮独立）**：HEAD `a14ba36`，工作树 clean；本轮重跑 `apps/api` `go test ./internal/handler/ -count=1` + `go vet ./...`（exit 0），以及 web 聚焦 `vitest run` 于 `resource.test.ts` + `schema-table.test.tsx` + `schema-crud.test.tsx`（**60/60**）。未在本轮重跑全量 vitest/e2e/Docker Compose；全量四类回归以 GOAL-010 S5 执行记录为既有证据引用，不扩写为本轮独立全量验收。

### 成果（有证据）

对照 F-002-001 原文三项证伪点与建议关闭路径：

| 原主张 / 证伪点 | 本轮核对 | 证据 |
|-----------------|----------|------|
| `schema-table.tsx` 无论 `dataSource` 均调用 `fetchRecords` | **已消除** | 现用 `schemaTableDataSource` + `isValidDataSource` 解析端点；合法时 `fetchResourceList(...)`；非法/缺失 fail-closed 不请求（`apps/web/src/renderer/schema-table.tsx`） |
| `records.ts` 强制 `id/name/status/owner/updatedAt` 与固定分页形状 | **已消除** | `records.ts` **已删除**；`resource.ts` 提供 `ResourceItem`/`ResourceList` + 统一 envelope `{items,total,page,pageSize}`，items 为任意 JSON 对象；`rowKey` 不变量在 table 层强制（`resource.ts` / `schema-table.tsx`） |
| `health.go` 仅注册 `/api/records`，无通用 CRUD | **已消除** | `Register` 仅 `registerResource(... usersResource / rolesResource)`；通用工厂在 `resources.go`；`records.go` **已删除**（`health.go` L32–33） |
| 建议路径：Schema 驱动通用适配层 + 后端资源契约；records 降为实例后产品面可退场 | **已落实** | 契约 I-010-001（GOAL-010 S1）；工厂 + users/roles 注册（GOAL-010 S2 + GOAL-011）；前端泛化（S3）；users/roles Schema-only 接入 + records `0006 records_retire`（S4/GOAL-011）；self close-out A-002 + 本意见 |
| 产品 fixture 仍绑 records | **已消除** | fixtures：`users.json`/`roles.json`/`data-table.json`/`search-form-table.json` 的 `dataSource`/`url` 指向 `/api/users` 或 `/api/roles`；产品 `apps/web/src` 无 `fetchRecords`/`RecordItem`/`/api/records` |
| 台账单义性（索引 ↔ 关闭表 ↔ 载体状态） | **一致** | 本文件索引 A-002 行、关闭证据表、仍开放项/后续均写 F-002-001 **`fixed`（2026-08-04）**；GOAL-010 `done / 5/5`；与 GOAL-009 曾出现的「索引仍 open / 表已 fixed」类矛盾不同 |

### 对照关闭判据（F-002-001）

| 判据 | 判断 | 说明 |
|------|------|------|
| 硬编码 records 传输/字段模型已去除 | ✅ | 见上表三项证伪点 |
| 存在可重复使用的通用适配层 + 后端注册契约 | ✅ | `resource.ts` + `resources.go` + I-010-001 |
| 至少一个非 records 业务资源仅靠 Schema 接入可证明 | ✅ | users + roles 双资源；fixture 与 e2e（GOAL-010 S5 既有）面向 users/roles |
| 关闭路径与用户裁决（D-014 通用适配层，不降级 VP）一致 | ✅ | 未采用「降级为单一 records 示例」路径 |
| `fixed` 声明与证据路径可指回 | ✅ | 关闭证据表 → GOAL-010 S1～S5 + GOAL-010 A-002 self；本轮代码/测试可独立复核核心链 |
| 信息门禁（载体 GOAL-010） | ✅ 不阻断本 scope | `I-010-001`/`I-010-002` 均为 verified（载体台账）；本 finding-closure 不重开信息项 |

### Findings

- **无开放 required finding。**
- **R-001（recommended / low · 观察项，非阻断）**：`apps/web/src/protocol/conformance/stage3-fixtures.test.ts` 仍含协议历史 fixture 的 `/api/records` URL。属协议一致性测试数据，非产品运行面注册/路由/fixture；与 GOAL-010 A-002 / GOAL-011 既有观察边界一致，**不否定** F-002-001 的产品面闭合。
- **R-002（recommended / low · 观察项，非阻断）**：GOAL-010 S5 执行记录引用的 HEAD 为 `21e6bd7`；本轮核对时 HEAD 为 clean `a14ba36`。关闭证据链的**产品语义**在本快照仍成立；进入 Root/VP-002 关门前建议编排器在 close-out 收据中固定当前 revision 身份（不要求为此重开 F-002-001）。

### 必改项汇总

- 无。F-002-001 的 `fixed` 闭合**维持**；本意见不新增 required，不要求回退关闭状态。

### 与既有意见的异同

- 与 **GOAL-010 A-002（self · close-out · pass）** 一致：三项证伪点已消除，通用适配层 + 双实体 + records 退场构成充分关闭证据。
- 与 **Root A-002 响应关闭证据表** 一致：F-002-001 标为 `fixed（2026-08-04）` 有独立代码与聚焦测试支撑，非仅文档宣称。
- 与 **GOAL-009 A-002 finding-closure（conditional）** 的差异：彼时问题在 Root 索引与关闭表单义性；本轮 F-002-001 索引/关闭表/载体状态三处一致，故本 scope 可 **pass**。
- 本意见**不**把 F-002-001 闭合扩写为 Root 已可关门或 VP-002 已可关门；那是独立 close-out / Vision 流程。

### 结论 + 建议给编排器 / 用户的下一步

- **pass**：F-002-001 关闭证据充分、可重复核对；维持 `fixed` 合法闭合。A-002 三条 required 在独立复审下仍全部成立为已闭合（本 scope 仅复审 F-002-001；002/003 不在本轮重审）。
- 建议 `/govern`：采纳本意见（索引已含 A-003；无需改 F-002-001 状态）；可选处理 R-001/R-002 为 handled/留痕；进入 **Root close-out 关门审计** 与/或 **VP-002 关门** 的用户裁决（二者均为独立流程，不由本 finding-closure 自动放行）。

### 声明

本意见仅追加独立审计记录（`source: independent`），**不**修改目标 `status` / 检查点 / 派生 `progress`、关闭证据表结论、`goal-tree.md` 或 VP 状态；响应与后续推进由 `/govern` 处理。

## A-003 响应（编排器 · 2026-08-04）

- **响应编号**：A-003（independent · finding-closure · pass）
- **响应来源**：`/govern` 编排器（self 侧记录，非独立审）
- **用户裁决**（P-004）：**采纳 `pass`**；F-002-001 的 `fixed` 合法闭合**维持**（不重开、不改状态）；按 A-003 建议进入 **Root close-out 关门审计**与 **VP-002 关门**的用户裁决流程——本响应只采纳意见并维持已闭合状态，不自动放行关门。
- **冲突检查**：A-003（independent · pass）与 GOAL-010 A-002（self · close-out · pass）同向一致，无 verdict / 必改项冲突（P-004 §3.2 不触发）。本 finding-closure scope 的 self 覆盖已由载体 GOAL-010 A-002 承担；Root 全范围 close-out 的 self 覆盖按 P-004 §3.1 在关门裁决中询问，不在此静默跳过。

### 响应内容

1. **采纳 A-003 verdict = `pass`**：F-002-001 关闭证据充分、可重复核对。编排器静态复核与 A-003 一致——`schema-table.tsx` 现按 `dataSource` 解析端点并 fail-closed；`records.ts` 已删除（`resource.ts` 统一 envelope + `rowKey` 不变量）；`health.go` 仅注册 users/roles 通用资源；产品 fixtures 指向 `/api/users` / `/api/roles`；`apps/web/src` 无 `fetchRecords`/`RecordItem`/`/api/records` 残留。
2. **F-002-001 维持 `fixed`（2026-08-04）**：正式意见索引 ↔ 关闭证据表 ↔ 载体状态（GOAL-010 `done / 5/5`）三处一致；不重开。
3. **R-001 → handled**：`apps/web/src/protocol/conformance/stage3-fixtures.test.ts` 的 `/api/records` URL 属协议一致性测试的历史数据，非产品运行面注册/路由/fixture；与 GOAL-010 A-002 / GOAL-011 既有观察边界一致，不否定 F-002-001 产品面闭合；不升级为 required。
4. **R-002 → handled**：revision 身份固定——本轮响应时 HEAD `a14ba36`（与 A-003 陈述一致），工作树仅本 `03-audit.md`（A-003 意见 + 本响应）未提交；Root close-out 收据将以最终 committed revision 固定并回填本节，不要求为此重开 F-002-001。

### 关闭证据表

| finding | 状态 | 证据路径 |
|---------|------|----------|
| F-002-001（Renderer 硬编码 records 实体） | **fixed 维持（2026-08-04）** | A-003（independent · pass）独立复核 + GOAL-010 A-002（self · close-out · pass）；关闭证据表 → GOAL-010 S1～S5 + GOAL-011 双实体载体；HEAD `a14ba36` 聚焦核对（api handler 60/60、go vet exit 0，见 A-003） |
| A-003 R-001（协议 fixture 残留 records URL） | **handled** | `stage3-fixtures.test.ts` 为协议一致性测试历史数据；与 GOAL-010 A-002 / GOAL-011 既有观察边界一致 |
| A-003 R-002（关闭证据 revision 身份） | **handled** | HEAD `a14ba36` + 工作树仅 `03-audit.md` 未提交；close-out 收据固定 committed revision |

### 仍开放项

- **无开放 required finding**。A-002 三条 required（F-002-001/002/003）全部合法闭合（`fixed`）；recommended F-002-004~006 已实施（非 required 闭合）。
- R-001/R-002 已 handled；无其他 open recommended 阻断项。
- **Root close-out 关门审计与 VP-002 关门为独立用户裁决**：Root 全范围尚无 self close-out 审计，是否补自审待用户按 P-004 §3.1 裁决（本响应不代答）。

### 状态

Root 当时 `status: active`、派生进度 `5/5`；`goal-tree.md` 与 `00-meta.md` 已于 2026-08-04 同步本响应与 close-out 裁决入口。后续见 **A-004** self close-out。

## A-004 · Root close-out 关门审计（2026-08-04）

- **source**：self
- **auditor**：Grok Build · `/govern`
- **类型 / scope**：close-out；Root 全目标成功边界、五个纲领检查点 R1～R5、子目标交付链、意见台账（A-001～A-003）、信息台账（I-001～I-006）、A-002 整改闭合与当前 revision 回归收据。不审 VP-002 关门（独立 `/vision` 流程）。
- **verdict**：**pass**
- **audit_type**：close-out

### 范围与区间

- 工作区：`workspace-002-production-admin-foundation`；`workspace.md` 的 `root_goal`、`canonical_scope`、`vision_role: delivery`、`primary_plan: VP-002-production-admin-foundation` 与 Root `plan_refs`/`primary_plan` 一致；`shared_materials_catalog: none`，未使用共享资料作为事实或关闭证据。
- 愿景链：Charter `schema-ui-core-admin-foundation@0.1.0` active；VP-002 `vision_ref` 精确匹配；Vision Review **0 open required**（仅 recommended `F-V003` 仍 open，不阻断本 Goal 关门）。
- 证据边界：子目标关门审计与既有独立意见为历史证据；本轮 close-out **现时回归**见下「本轮回归收据」。产品面静态抽查：`health.go` 仅注册 users/roles；`apps/web/src` 除协议一致性测试外无 `fetchRecords`/`RecordItem`/`/api/records` 产品残留。

### 成果（有证据）

| 纲领阶段 | 状态 | 主要载体 / 证据 |
|----------|------|-----------------|
| R1 · Schema Renderer 产品化 | 通过 | GOAL-002/003/004 `done`；A-001 self R1 pass；I-001 verified |
| R2 · 真实认证与请求级身份 | 通过 | GOAL-005 `done`；I-002 verified；browser E2E / CI 既有证据 |
| R3 · 持久化 RBAC 与菜单投影 | 通过 | GOAL-006 `done`；I-003 verified；迁移/种子/permission/features |
| R4 · Schema 驱动 CRUD 闭环 | 通过 | GOAL-007 `done`（历史 records 代表实体）→ GOAL-010/011 演进为 users/roles 通用资源；I-004 verified |
| R5 · 工程化与 fork 体验 | 通过 | GOAL-008 `done`；I-005 verified、I-006 closed（操作日志加分已实施）；REPRO-003 / smoke / compose |
| A-002 整改波次 | 通过 | GOAL-009 `done`（F-002-002/003 + recommended 004～006）；GOAL-010 `done` + GOAL-011 `done`（F-002-001）；A-003 independent finding-closure pass |

子目标树：GOAL-002～GOAL-011 **全部 `done`**；无未关门子目标阻断 Root。

### 对照成功边界

| 成功边界（00-meta） | 状态 | 证据摘要 |
|--------------------|------|----------|
| Schema Renderer 默认页面能力 + 结构/运行时/失败路径 | 通过 | R1 三子目标 + 默认 `schemaUrl` 主路径 + 回归 |
| 真实认证链路（登录/登出/会话/请求身份） | 通过 | GOAL-005；JWT+refresh；F-002-003 清理路径 GOAL-009 |
| 用户/角色/菜单持久化 + 后端授权 | 通过 | GOAL-006 RBAC；GOAL-011 users/roles API + 授权 |
| 至少一个实体 Schema 驱动 CRUD 闭环 | 通过 | 现时 users/roles Schema-only；通用工厂 + 前端适配层（GOAL-010/011） |
| 种子/环境/健康/容器/fork ≤15 分钟可复现 | 通过 | GOAL-008 S1～S5；QUICKSTART + REPRO-003 + smoke.sh + compose |
| 阶段事实与审计；关门前无开放 required 信息项/必改 finding | 通过 | 见意见/信息台账；本意见确认 |

### 意见台账（相关 · 关门 scope）

| 意见 | source | verdict | 与关门关系 |
|------|--------|---------|------------|
| A-001 | self | pass | R1 阶段退出；无开放 R1 required |
| A-002 | independent | fail → 已响应 | F-002-001/002/003 **全部 `fixed`**；004～006 recommended 已实施非阻断 |
| A-003 | independent | pass | F-002-001 关闭证据独立复核；R-001/R-002 handled |
| A-004 | self | pass | 本 close-out |

- **开放 required finding：0**
- **冲突：无**（A-003 与 GOAL-010 self close-out 同向；A-002 失败结论已由 fixed 整改吸收）

### 信息台账（关门 scope）

| ID | 级别 | 状态 | 说明 |
|----|------|------|------|
| I-001 | required | **verified** | R1 方案边界 |
| I-002 | required | **verified** | R2 认证方案 |
| I-003 | required | **verified** | R3 持久化/权限模型 |
| I-004 | required | **verified** | R4 代表实体方向 |
| I-005 | required | **verified** | R5 工程/fork 边界 |
| I-006 | non-blocking | **closed** | 操作日志可选加分（GOAL-008 S6 已实施） |

- **到期开放 required 信息项：0**
- 无 `accepted-residual` 依赖；无信息冲突阻断关门。

### 本轮回归收据（2026-08-04 · close-out 现时）

| 检查 | 结果 |
|------|------|
| 基线 revision | HEAD `a14ba36`（工作树含本 close-out 文档变更，未声称 clean tree CI） |
| `apps/api` `go vet ./...` | exit 0 |
| `apps/api` `go test ./... -count=1` | 7 有测试包全绿（cmd/server、account、auth、config、handler、store 等） |
| `apps/web` `vitest run` | **23 files / 491 tests** 全绿 |
| `apps/web` `tsc -b` | 干净 |
| 产品面 `/api/records` 残留 | 仅 `protocol/conformance/stage3-fixtures.test.ts` 历史协议数据（A-003 R-001 handled 边界） |
| 后端资源注册 | `registerResource(... usersResource / rolesResource)` only |

未在本轮重跑 Docker Compose 全路径或 Playwright E2E；工程/浏览器证据以 GOAL-008/009/010/011 既有关门与 S5 记录为引用，不扩写为「本轮独立容器验收」。不主张 GitHub-hosted Actions 本快照已绿。

### Findings

- **无开放 required finding。**
- **无新增 recommended finding。** 既有 A-003 R-001/R-002 已 handled，不重开。

### 必改项汇总

- 无。

### 结论 + 建议下一步

- Root 五个纲领检查点 `5/5` 全勾选；子目标 GOAL-002～011 全部 `done`；A-002 三条 required 合法闭合；信息门禁无到期开放 required；愿景/工作区绑定一致；本轮 API/Web 回归通过。
- **verdict = pass**。按用户本轮指令（补 Root self 关门审计，通过后置 `done`）：**Root 置 `status: done`**，派生进度保持 `5/5`；同步 `goal-tree.md` / `00-meta` / `02-execution`。
- **VP-002 关门**不由本 Goal 审计自动放行，建议下一拍 `/vision` 在工作区证据链上评估 VP-002 是否可关闭。

## A-005 · apps/api + apps/web · VP-002 产品意图独立复审（2026-08-04）

- **source**：independent
- **auditor**：Grok Build（完全独立代码审计；**未**调用 audit/govern skill 出意见）
- **类型 / scope**：product-fit + execution-facts；核对 `apps/api`、`apps/web` 当前代码相对 [VP-002](../../../vision/plans/VP-002-production-admin-foundation.md) 七条产品级成功标准与「最终判断标准」（改 Schema 接业务、不重写 Renderer 主路径）。不修改业务代码；意见落盘后由 `/govern` 响应。
- **verdict**：**fail**（1 required open）。核心 Admin 能力（Schema 主路径、真实认证、持久化 RBAC、users/roles Schema CRUD、种子、工程化）经静态抽查 + 全量回归（`go test ./...` 全绿；web vitest **491/491**）**未发现目标级漂移**；但默认 Shell **仍存在产品面阻断死链**，不满足「可直接使用的生产级 Admin 基架」洁净交付。

### 范围与方法

- 工作区：`workspace-002-production-admin-foundation`；Root 当时为 `done / 5/5`；primary plan = VP-002。
- 对照权威：VP-002 意图与产品级成功标准 1～7；非目标（完整 IAM、全量协议等）不抬高验收。
- 方法：只读核对 api/web 关键路径（auth、resources 工厂、users/roles、seed、schema embed、App 主路径、users/roles fixtures、manifest）；本机复跑 `go test ./...` 与 `vitest run`；**未**重跑 Compose/Playwright 本轮。

### 对照 VP-002 七条（摘要）

| # | 标准 | 结论 | 证据要点 |
|---|------|------|----------|
| 1 | Schema Renderer 主路径 | **满足** | `App.tsx` → `loadPageDocument` → `RenderPage`；通用 `resource.ts` + `schema-table.tsx`；无 `fetchRecords`/`RecordItem` 产品路径 |
| 2 | 真实认证 | **满足** | JWT access + opaque refresh + 中间件；dev session 仅 opt-in；`ValidateProd` 拦生产 dev-session 与弱 JWT secret |
| 3 | 持久化身份/最小权限 | **满足** | SQLite users/roles/permissions/menus；种子 admin/editor/viewer；后端 `requirePermission`；角色委派与 grant 子集约束 |
| 4 | Schema 驱动 CRUD | **满足** | `users.json`/`roles.json` 完整 list/create/edit/delete + 权限/密码表单；e2e 覆盖 users+roles 管理链（既有） |
| 5 | 可重复种子 | **满足** | `seedRBAC` 幂等；admin 种子 |
| 6 | Fork 接业务 | **部分** | 登录与 Schema 页路径可用；**default Shell 死链**阻断「完整可用后台」观感；新增页真实落点为 embed fixture（见 R-001） |
| 7 | 工程化 | **满足** | healthz/readyz、compose、CI smoke、Dockerfile；活动面无 `/api/records` 产品残留（CI 门禁） |

### Findings

#### F-001 · required · medium — default Shell 导航死链（manifest 有 page、fixture 无文档）

- **severity**：`apps/web/public/.well-known/schema-ui/app-manifest.json` 声明并导航到 `activity`、`settings`；`apps/api/internal/handler/schema.go` 仅 `//go:embed fixtures/schema/*.json` 服务现有 7 份文档（overview/data-table/search-form-table/form-controls/form-with-reactions/users/roles）。**无** `activity.json` / `settings.json`。
- **复现**：登录后点击 sidebar「Activity」或 user 区「Settings」→ `GET /api/schema/activity|settings` → `404 SCHEMA_NOT_FOUND`（fail-closed 正确，但导航不应广告必失败页）。
- **影响门禁**：VP-002 成功标准 1/6 的「可用 Admin 页面 / fork 后进入系统使用后台」产品洁净度；Root 重新关门与 VP-002 关闭证据诚实性。
- **建议关闭路径**：从 manifest pages+navigation **移除**占位入口（推荐），或补 **最小** Schema 文档使路由可渲染；并加 **manifest pageId ⊆ embed fixtures** 回归测试。

#### R-001 · recommended · low — fork 文档「新增页面」路径错误

- **描述**：根 `QUICKSTART.md` §4 写「在 `docs/schemas` 添加页面 Schema」；`docs/schemas/` 实为上游 **协议 JSON Schema**（node/page/action…），不是页面文档。权威落点为 `apps/api/internal/handler/fixtures/schema/*.json` + 重建 API（`apps/web/README.md` 已正确）。
- **影响**：fork 用户按 QUICKSTART 无法完成 VP「改 Schema 新增页」路径。
- **边界**：本 finding 非 apps 运行时阻断；作 GOAL-012 可选 S5。

#### R-002 · recommended · low — AuthUser 类型省略 permissions

- **描述**：`apps/web/src/account/auth-client.ts` `AuthUser` 仅 `id/name/roles`；运行时 JSON 仍含 `permissions` 且 `$context` 表达式依赖它。类型/文档漂移，运行未破。
- **建议**：扩展 `AuthUser` 与测试断言，避免后续重构剥掉 permissions。

#### R-003 · recommended · low — 改密只撤销 refresh，不吊销 access JWT

- **描述**：`UpdateUser` 在密码变更时撤销 refresh；access 在 `AUTH_ACCESS_TTL`（默认 15m）内仍有效。属常见短 access 模型；若要更严会话边界，需 access 黑名单或极短 TTL。
- **级别**：recommended（非 VP 硬缺口）。

### 必改项汇总

| ID | 级别 | 状态 | 载体 |
|----|------|------|------|
| F-001 | required | **open** | `GOAL-012-a005-shell-nav-fixtures` |
| R-001～R-003 | recommended | open / non-blocking | GOAL-012 S5 或后续 |

### 关闭证据表（A-005）

| Finding | 关闭路径 | 状态 | 证据 |
|---------|----------|------|------|
| F-001 | fixed（待实施） | **open** | — |
| R-001 | optional | open | — |
| R-002 | optional | open | — |
| R-003 | residual-acceptable by design | open | 短 access TTL |

### 结论

- **无**「用 records 演示冒充生产 Admin」类目标漂移（records 已 0006 退场；users/roles 语义资源在位）。
- **有**偷工减料式产品洁净度缺口：Shell 遗留死链未在 Root 关门前清零。
- **verdict = fail**；Root 不得维持 `done` 直至 F-001 合法闭合。

---

## A-005 编排响应（2026-08-04 · `/govern`）

- **用户指令**：阻断问题 → 工作区 2 新设子目标修正 + 回退 Root 关门。
- **裁决**：F-001 走 `fixed`；立项 `GOAL-012-a005-shell-nav-fixtures`（parent Root）；**Root `done → active`**，progress 保持 `5/5`。
- **本响应不**将 F-001 标为 fixed（待 GOAL-012 实施与关门证据）。
- **同步**：`goal-tree.md`、Root `00-meta`/`03-audit`、GOAL-012 五件套。

### A-005 关闭响应（2026-08-04 · GOAL-012 完成）

- **F-001 → fixed**：checked-in `app-manifest.json` 已移除 `activity`/`settings` 及 Workspace/Settings 导航；manifest 仅保留 7 个与 embed fixture 对齐的 pageId。`schema_test` 新增「manifest pageIds all have embed fixtures」门禁。证据：GOAL-012 02-execution S1～S3；A-001 self close-out pass。
- **R-001 → fixed（可选）**：`QUICKSTART.md` §4 已改为 `fixtures/schema` + rebuild API + manifest（GOAL-012 S5）。
- **R-002 / R-003**：保持 recommended open/non-blocking（AuthUser 类型；改密 access TTL 设计）。
- **治理投影**：GOAL-012 `done / 4/4`；Root 保持 `active / 5/5`（重新关门须用户独立裁决）；A-005 无开放 required。

---

## A-006 · apps/api + apps/web · VP-002 产品意图再审（2026-08-04）

- **source**：independent
- **auditor**：Grok Build（完全独立代码审计；**未**调用 audit skill 出意见）
- **类型 / scope**：product-fit + execution-facts；核对 `apps/api`、`apps/web` 相对 [VP-002](../../../vision/plans/VP-002-production-admin-foundation.md) 七条产品级成功标准与「最终判断标准」（改 Schema 接业务、不重写 Renderer 主路径）；评估是否存在目标漂移或偷工减料。
- **verdict**：**pass**。无 required finding；5 条 recommended（R-001～R-005）非阻断。

### 范围与方法

- 工作区：`workspace-002-production-admin-foundation`；Root `active / 5/5`；primary plan = VP-002。
- 方法：只读核对 auth、resources 工厂、users/roles、seed、schema embed、App 主路径、manifest↔fixture、settings/activity；本机复跑 `go test ./...` 全绿 + `vitest` **491/491**（审计当时）。
- 对照权威：VP-002 意图与成功标准 1～7；non-goals 不抬高验收。

### 对照 VP-002 七条（摘要）

| # | 标准 | 结论 | 证据要点 |
|---|------|------|----------|
| 1 | Schema Renderer 主路径 | **满足** | `App.tsx` → `loadPageDocument` → `RenderPage`；通用 `resource.ts` + `schema-table.tsx`；无 records 产品路径 |
| 2 | 真实认证 | **满足** | JWT + opaque refresh；dev session opt-in；`ValidateProd`；auth-lost 清会话 |
| 3 | 持久化身份/最小权限 | **满足** | SQLite RBAC + 幂等种子；`requirePermission`；角色委派/grant |
| 4 | Schema 驱动 CRUD | **满足** | users/roles Schema 闭环 + e2e 源码覆盖 |
| 5 | 可重复种子 | **满足** | `seedRBAC` 幂等；生产 `ADMIN_INITIAL_PASSWORD` |
| 6 | Fork 接业务 | **满足** | QUICKSTART → fixtures/schema；manifest 9 pageId 与 embed 1:1（含 settings/activity） |
| 7 | 工程化 | **满足** | healthz/readyz、compose、Docker、smoke 路径 |

### Findings（审计时）

#### Required

**无。**

#### R-001 · recommended · low — `AuthUser` 类型省略 `permissions`

- **位置**：`apps/web/src/account/auth-client.ts`
- **影响**：类型/文档漂移；运行时 `/me` JSON 仍含 `permissions`，表达式可工作。

#### R-002 · recommended · low — Renderer 主路径含 settings 域名副作用

- **位置**：`apps/web/src/renderer/render.tsx` `runRequest` 内硬编码 `/api/settings` → branding 事件
- **影响**：轻微污染 Renderer 纯度。

#### R-003 · recommended · low — Settings 写不进 operation_log

- **影响**：Activity 页看不到品牌变更；操作日志为 VP 加分项。

#### R-004 · recommended · low — 用户角色授予 UX 偏「能用即可」

- **事实**：roles 行值为 `string[]`，textarea 预填未走 wire coerce；标签未提示系统/自定义角色。
- **边界**：完整 IAM 为 non-goal；自定义角色仍须自由文本。

#### R-005 · recommended · low — 改密只撤销 refresh，不吊销 access JWT

- **事实**：短 access TTL 模型（默认 15m）；与历史 A-005 R-003 同向。

### 结论（审计时）

- **无**目标漂移；**无**阻断级偷工减料。
- **verdict = pass**；不阻断 Root/VP-002 关门门禁（关门仍为独立用户裁决）。

---

## A-006 编排响应（2026-08-04 · `/govern`）

- **用户指令**：在工作区 2 Root 落盘审计意见，并响应非阻断问题（修正）。
- **P-004**：无 required / 无冲突；独立意见已有，用户明确要求响应与修正，不强制另开 self 审计。
- **裁决**：R-001～R-004 走 **fixed**；R-005 走 **residual-by-design / handled**（短 access 模型，不引入 access 黑名单）。

### 关闭证据表（A-006）

| Finding | 关闭路径 | 状态 | 证据 |
|---------|----------|------|------|
| R-001 | fixed | **fixed** | `AuthUser.permissions?` + `parseAuthUser`；`auth-client.test.ts` 断言 `/me` 保留 permissions |
| R-002 | fixed | **fixed** | `render.tsx` 移除 settings 硬编码；`main.tsx` `useResourceFetcher` 在成功 PATCH `/api/settings` 后 `notifyBrandingChanged()` |
| R-003 | fixed | **fixed** | migration `0008 operation_log_settings` + `EventSettingsUpdate`；`settingsPatch` best-effort 写操作日志；`settings_test` 断言 |
| R-004 | fixed | **fixed** | 表单初值统一 `coerceFieldValue`；`coerceToKind` 将 `string[]`→textarea 逗号串；fixture 标签说明系统/自定义角色 |
| R-005 | residual-by-design | **handled** | 用户书面采纳短 access TTL 模型；不作为本波次代码变更 |

### 回归收据（响应实施后）

- `apps/api` `go test ./... -count=1` 全绿；`go vet` 路径随测试包覆盖
- `apps/web` `vitest run` **492/492**；`tsc -b` 干净
- Root 保持 `active / 5/5`；无开放 required finding
- **未做**：未自动 Root/VP-002 关门；未重跑 Compose / Playwright 全路径

### 治理投影

- 本响应写入 Root `03-audit` / `02-execution` / `goal-tree` 注记
- A-006 无开放 required；recommended 已闭合或 residual 留痕
