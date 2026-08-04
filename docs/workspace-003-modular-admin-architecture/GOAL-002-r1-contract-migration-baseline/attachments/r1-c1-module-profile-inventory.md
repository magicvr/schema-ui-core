---
id: R1-C1-EVIDENCE
title: R1 C1 模块、Profile 候选与注册路径基线
status: recorded
created: 2026-08-04
updated: 2026-08-04
parent: GOAL-002-r1-contract-migration-baseline
version: 0.1.0
source: self
---

# R1 C1 · 模块、Profile 候选与注册路径基线

## 证据口径

本记录只描述当前仓库可核对的入口、注册路径和候选依赖闭包。`present` 表示当前源码存在该路径，`candidate` 表示可以承接 R1 Profile 设计的候选边界，`absent` 表示本轮检索没有发现实现。候选矩阵不是 R2 的精确 Profile 配置，也不证明模块化重构已经实施。

## 当前注册与消费链

| 领域 | 当前路径 | 证据 | 当前判定 |
|------|----------|------|----------|
| API 启动与中央注册 | `main()` 创建 `http.ServeMux` 并调用 `handler.Register`；`Register` 集中挂载 health/readiness、auth/accounts、users、roles、operations、settings、schema | `apps/api/cmd/server/main.go:23-27,66-69`; `apps/api/internal/handler/health.go:27-37` | present；集中式注册，不是模块 Registry |
| API 资源注册 | `registerResource` 统一挂载 list/detail/create/update/delete，并通过 `Authenticator` middleware 包装 | `apps/api/internal/handler/resources.go:156-182` | present；资源契约可抽取为候选模块边界 |
| API 权限 | `requirePermission` 对无身份返回 401、无权限返回 403 | `apps/api/internal/handler/resources.go:185-199` | present；fail-closed 权限入口 |
| API Schema | `schemasHandler` 挂载 `GET /api/schema/{pageId}`；嵌入 fixture 文件名映射为 pageId；未知 page 返回 `SCHEMA_NOT_FOUND` | `apps/api/internal/handler/schema.go:15-27,29-68` | present；静态 fixture 消费路径 |
| Web 启动与认证门 | `loadAppManifest()` 成功后渲染 `AuthProvider` → `AuthGate` → `App` | `apps/web/src/main.tsx:81-103` | present；Web 消费根 |
| Web route/schema 消费 | `App` 对 `manifest.pages` 调用 `matchRoute`，再以 `page.schemaUrl` 进入 Schema 页面 surface | `apps/web/src/app/App.tsx:313-317,354-364` | present；Manifest 驱动，不是独立中央 route registry |
| Web navigation | `projectNavigation` 将 Manifest navigation 投影为 `top`、`sidebar`、`user` | `apps/web/src/app/navigation.ts:137-152` | present；按 context 做可见性投影 |
| Web capability/permission | 页面 `meta.requiredCapabilities` 由权限验证器校验，缺少 permission capability 时返回 `CAPABILITY_REQUIRED` | `apps/web/src/renderer/permissions.ts:94-110` | present；能力错误可追踪 |
| Web icon registry | `App.tsx` 内有集中 `iconRegistry` | `apps/web/src/app/App.tsx:43-59` | present；仍是应用内注册表 |
| 协议基线 | I-PROTO-001 v0.1.3 固定 D-NODE、D-EXPR、D-DATA、D-PERM、D-APP、D-VER、D-VAL 为 include，D-COMP/D-ACT/D-TABLE/D-FORM 为 include-partial，D-UPLOAD 为 exclude | `docs/workspace-001-mvp-admin-foundation/GOAL-001-mvp-admin-foundation/attachments/I-PROTO-001-coverage-draft.md:35-80` | present；只继承 protocol range |
| Shell | 未发现 `apps/shell`、Shell 注册入口或独立 Shell 模块 API | 本轮 `rg` 检索结果；仓库当前路径清单 | absent；作为 R1 target boundary，不能写成现状实现 |
| Profile registry | 未发现 `mvp`/`admin` Profile 定义、配置文件、注册函数或现成依赖矩阵 | 本轮 `rg` 检索结果 | absent；本记录只提供候选，不冻结精确集合/precedence |

## `mvp`/`admin` Profile 候选与依赖闭包

| 候选模块 | `mvp` | `admin` | 当前依赖闭包 | 现状与边界 |
|----------|--------|---------|--------------|------------|
| `core.server-registration` | required candidate | required candidate | `cmd/server/main.go` → `handler.Register` → Store/Auth/HTTP handlers | 当前集中注册；未来 Registry 归属待 R1/C3 决策 |
| `core.auth-session` | required candidate | required candidate | API `authsHandler`/`accountsHandler` + Web `AuthProvider`/`AuthGate` + request identity | 现有路径 present；Profile 配置尚不存在 |
| `core.manifest-route` | required candidate | required candidate | `loadAppManifest` → `AppManifest` validation → `matchRoute` → `page.schemaUrl`/`loadPageDocument` | Manifest 消费链 present；聚合 API 不存在 |
| `core.navigation-capability` | required candidate | required candidate | `projectNavigation` + navigation context + `requiredCapabilities` validation + API `requirePermission` | 当前分属 Web/API；统一模块契约待 C3 |
| `core.schema-render` | required candidate | required candidate | API embedded schema documents + Web `SchemaPageSurface`/renderer | schema fixture present；动态模块供给不成立 |
| `admin.users` | optional candidate | required candidate | `usersResource` + users schema/menu + Store users + permissions | API 注册 present；精确 Profile 归属待 I-004/R2 |
| `admin.roles` | optional candidate | required candidate | `rolesResource` + roles schema/menu + Store RBAC + permissions | API 注册 present；精确 Profile 归属待 I-004/R2 |
| `admin.settings` | optional candidate | required candidate | `settingsHandler` + settings schema/menu + Store settings + permissions | API 注册 present；精确 Profile 归属待 I-004/R2 |
| `admin.activity` | optional candidate | required candidate | `registerOperations` + operation log + activity schema/menu + permissions | `activity` page/menu 与 operations handler 的映射需后续冻结，当前不能假定等价 |
| `shell.registry` | target candidate | target candidate | Shell runtime + module registry + Profile selection | 当前 absent；不计入现状实现 |

## 检查点结论

C1 的“现状路径 + 候选闭包矩阵”已形成，可作为 R1 C3/C4 继续决策的输入。未发现 Shell/Profile registry 和现成矩阵已明确记录为缺口；这不会被解释为实现完成。Root I-001 仍为 `open`，直至 R1 方案冻结、阶段审计和 `/govern` 响应完成。
