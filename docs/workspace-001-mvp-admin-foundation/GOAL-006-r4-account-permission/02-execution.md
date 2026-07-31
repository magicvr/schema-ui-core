---
id: GOAL-006-r4-account-permission
doc: execution
status: active
parent: GOAL-001-mvp-admin-foundation
created: 2026-07-31
updated: 2026-07-31
version: 0.3.0
---

# 执行记录 · GOAL-006

## R4 实施计划（2026-07-31 · 计划，非事实）

> 用户在 2026-07-31 经 `/govern` 确认进入 R4 实施阶段并选定完整推进路径。以下为实施计划（P-001 路线图细化）；实际完成情况按下方时间线记录，**计划不等于完成**。

| 步骤 | 内容 | 产物 |
|------|------|------|
| E1 | Go 侧会话模型：静态/注入会话（D-004 允许的最小闭环会话方案） | `apps/api/internal/account/`（session、permission、response 类型） |
| E2 | Go 侧 `GET /api/accounts/me`：返回 `{ user, features }`，无会话 fail-closed（401） | `apps/api/internal/handler/account.go` + 路由注册 |
| E3 | Go 侧独立鉴权模型：`$context.user/features` 表达式求值 + 授权检查，供业务路由调用（D-004：Go 独立鉴权，不信任前端结果） | `apps/api/internal/account/permission.go` |
| E4 | Web 侧 `$context` 挂载：`main.tsx` 注入 `loadAccountContext`（fetch `/api/accounts/me`，失败降级空上下文但可复核），`App` 消费 | `apps/web/src/account/context.ts` |
| E5 | Web 侧 D-PERM 求值引擎：`effectivePermission` AND 公式、结构边（children/tabs/table 操作边/默认 submit）、cascade 类型白名单、intent 挂载点矩阵、执行时序 fail-closed（ADR-0023 D2a/D3/D4a/D4b/D4c） | `apps/web/src/renderer/permissions.ts` + 单元测试 |
| E6 | L2 校验错误码（fail-closed）：`PROTOCOL_VERSION_TOO_LOW` / `CAPABILITY_REQUIRED` / `PERMISSION_CASCADE_TYPE_INVALID` / `PERMISSION_CASCADE_KEYS_INVALID` / `PERMISSION_CASCADE_SOURCE_MISSING` / `PERMISSION_INTENT_FORBIDDEN` / `PERMISSION_INTENT_INVALID` | 并入 E5 校验层 |
| E7 | `permissions-inheritance` 17 例 fixture 对照测试：13 valid（求值 targets；其中 5 例含 execution 时序断言）+ 4 invalid（错误码断言） | `apps/web/src/renderer/permissions-inheritance.test.ts` |
| E8 | 验证：Go `go test ./...` / `go build ./...`；Web `npm run test` / `npm run build`；运行时 HTTP 证据（启动 API，curl `/api/accounts/me` 有/无会话） | 命令输出与 02-execution 时间线事实 |
| E9 | 实施阶段自审（04，source=self）后调 grok `/audit GOAL-006` 独立审计，合并响应意见 | 03-audit A-00N 追加 |

边界：R4 实施范围仅含最小会话与权限求值链路（D-004）；不建账号 CRUD 管理页、不建 SSO/联邦、不做细粒度审计后台、不实现 R5 Renderer 全量。会话方案先采用静态/注入（D-004 明确允许），token 会话留待后续决策。

## 时间线

### 2026-07-31 · R4 实施：Go 会话与独立鉴权（E1–E3）

- 新建 `apps/api/internal/account/`：
  - `session.go`：`Session` / `User` 类型与 `StaticDevSession()`（静态/注入会话，D-004 允许的最小闭环方案；角色 `admin`+`editor`、`features.beta=true` 覆盖全部 fixture 场景）。
  - `permission.go`：`Evaluate`（`$context.user.*` / `$context.features.*`，`==` / `!=` / `contains`，literal 解析与 web `parseLiteral` 语义一致）与 `Allow`（fail-closed 授权入口：未解析表达式拒绝、未声明路径拒绝）。
  - `permission_test.go`：表达式求值 / 非法表达式 / fail-closed 共 12 断言。
- 新建 `apps/api/internal/handler/account.go`：`GET /api/accounts/me` 返回 `{ user, features }`；`sessionProvider` 可注入（测试用），无会话 → 401 `UNAUTHENTICATED`（fail-closed）。路由经 `handler.Register` 挂载。
- 新建 `apps/api/internal/handler/account_test.go`：有会话 200 结构断言 + 无会话 401 fail-closed。
- 验证：`go vet ./...`、`go test ./...`、`go build ./...` 全部通过；`HTTP_ADDR=127.0.0.1:18091` 启动后 `GET /api/accounts/me` 返回 200 与静态会话 JSON、`/healthz` 200（运行时证据，见下）。
- 本步不建 token 会话/登录端点；静态/注入方案为 D-004 明确允许，token 留待后续决策。

### 2026-07-31 · R4 实施：Web $context 挂载与 D-PERM 求值引擎（E4–E7）

- 新建 `apps/web/src/account/context.ts`：`loadAccountContext`（GET `/api/accounts/me` → `NavigationContext`；失败返回空 context + error，fail-closed 渲染姿势）。
- 新建 `apps/web/src/account/context.test.ts`：成功映射 / HTTP 错误降级 / 网络错误降级 3 例。
- `apps/web/src/main.tsx`：manifest 加载后先 `loadAccountContext` 再渲染 `App`，`navigationContext` 注入真实身份/权限来源（R3 默认空 context 得到衔接）。
- `apps/web/vite.config.ts`：`/api` 代理 → `http://127.0.0.1:8080`（dev 联调契约）。
- 新建 `apps/web/src/renderer/permissions.ts`（D-PERM 求值引擎，对照 ADR-0023）：
  - `validatePermissions`：L2 校验（`PROTOCOL_VERSION_TOO_LOW` / `CAPABILITY_REQUIRED` / `PERMISSION_CASCADE_TYPE_INVALID` / `PERMISSION_CASCADE_KEYS_INVALID` / `PERMISSION_CASCADE_SOURCE_MISSING` / `PERMISSION_INTENT_FORBIDDEN` / `PERMISSION_INTENT_INVALID`）；cascade 类型白名单 `section/grid/form/tabs/table`、keys 非空去重仅 `edit|delete`、同名源权限、intent 挂载点矩阵（RowAction / toolbar Trigger / actionButton）、columns 禁 intent。
  - `evaluatePermissionTargets`：`effectivePermission` AND 公式（祖先 cascade 边界 + 目标本地，未声明按 true，只能收紧）；结构边（`children[]`、`tabs.props.items[].content`、table `actions/toolbar` 挂载边、default form 隐式 submitAction）；modal content / `navigatedPage` 新根；D4a 表单 edit 白名单（`input/inputNumber/datePicker/dateRangePicker/select/upload`，仅 default form）；未标注意图只走本地权限；columns 仅本地。
  - `executeAction`：D4c 统一时序 fail-closed（visible → permission → disabled/requiresSelection → confirm → action；拒绝后不展示 confirm、不构造动作）。
- 新建 `apps/web/src/renderer/permissions-inheritance.test.ts`：固定 `attachments/dperm/cases.json`（SHA-256 核验 = D-004 记录值）17 例全部对照通过（13 valid 求值 + 4 invalid 错误码；13 valid 中含 5 例 execution 时序断言）。
- 验证：`npm run test` 94 例全过（含既有 76 例回归）、`npm run build`（tsc + vite）通过。
- 端到端联调：API `127.0.0.1:8080` + vite dev 代理链路验证 `GET /api/accounts/me` 经 `/api` 代理返回 200 与静态会话 JSON（运行时证据）。

## 完成后边界

### 2026-07-31 · R4 目标立项与范围登记

- `/govern` 复核显式工作区、Charter/VP 对齐、Root 路线图、R3 关门（GOAL-005 A-007/A-008）与父目标信息门禁。
- 创建本目标五件套和 `attachments/` 目录，并将 `GOAL-006-r4-account-permission` 挂到 `GOAL-001-mvp-admin-foundation`。
- 将 R4 范围记录为账号权限最小 API 设计、`D-PERM` 映射冻结与前后端鉴权链路；明确排除 R5 Renderer/业务范例、完整权限继承产品化与完整协议支持。
- 登记 `I-006-001` 为 required/open（R4 方案冻结前验证）；父目标 `I-PROTO-002` 保持 open、作为 R4 **实施**门禁。
- 同步工作区 `goal-tree.md`；Root `status` 保持 `active`、`progress` 保持 `3/6`。
- 本次没有修改 `apps/*`，没有收集或验证任何 `I-006-*`，没有放行方案冻结、实施或 `done`；`I-PROTO-002` / `I-PROTO-003` 未改变。

### 2026-07-31 · R4 契约收集与方案冻结（D-004）

- 从固定 commit `ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b`（artifact `2.7.0`）拉取并落盘 D-PERM 资料至 `attachments/dperm/`：
  - `permissions-inheritance/cases.json`（fixtureVersion 1.0，17 cases：13 valid + 4 invalid；target kinds：formField 7 / formSubmit 7 / rowAction 5 / actionButton 5 / toolbarTrigger 2 / column 1）
  - `node.schema.json`（`Permissions` view/edit/delete；`PermissionCascade` keys enum edit|delete、unique、minItems 1）
  - `0023-container-permission-inheritance.md`（ADR-0023：effectivePermission AND 公式、4 条结构边、5 类 cascade type 白名单、columns 不参与、执行时序 fail-closed）
  - `permission-inheritance.md`（v2.7 场景：编辑/删除继承扩展示例，R5 范例页候选）
  - SHA-256 已核验（见 D-004 表）。
- 对照语义规范原文（固定 commit 下 `01-node-protocol.md` §3.9/§3.9.1、`03-component-registry.md` intent 矩阵与 D4a 表单 edit 目标、`08-renderer-spec.md` §7.1 执行时序）与覆盖表 v0.1.3（`D-PERM=include`、`permissions-inheritance=include`），完成最小 API 设计与映射结论。
- 用户确认「按此方向冻结」且「I-PROTO-002 在方案冻结时一并闭合」；D-004 落盘。
- `I-006-001` → `verified`；父目标 `I-PROTO-002` → `verified`（Root meta 同步留痕）。
- 本阶段**未**实施代码（`apps/*` 未修改）；R4 实施仍需用户指令并记实施事实；`I-PROTO-003` / `I-PROTO-004` 未改变。

## 完成后边界

1. R4 方案已冻结（D-004）；`I-006-001` 与父目标 `I-PROTO-002` 已 `verified`。
2. R4 实施（前端权限求值/显隐禁用、Go 会话与鉴权模型、fixture 对照测试）按 02-execution 后续时间线推进并记事实；不把「方案冻结」写成「已实现」。
3. 开放 required 信息项到期前不得越过对应门禁；`I-PROTO-003`（R5 验收/关门）不属本目标处理。

## 进度评估

R4 完成「契约发现与信息就绪」「方案冻结」两阶段（`I-006-001` 与父目标 `I-PROTO-002` 均 `verified`），并完成 R4 **实施**（本日）：Go 会话与 `/api/accounts/me`、Go 独立鉴权表达式求值、Web `$context` 挂载、D-PERM 求值引擎与 17 例 fixture 对照（13 valid 求值 + 4 invalid 错误码；13 valid 中含 5 例 execution 时序断言）全部落地；`go test`/`go build`、web 94 项测试、`npm run build`、HTTP 运行时与代理联调证据均已记录。Root `progress` 仍为 `3/6`（R4 为纲领检查点 4/6，须在关门自审与用户确认后按检查点重算）；「验证与关门」阶段未开始。
