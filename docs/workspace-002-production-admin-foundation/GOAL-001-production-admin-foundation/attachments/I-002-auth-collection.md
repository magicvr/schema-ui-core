---
title: I-002 · 认证/会话/部署/依赖现状收集与可选方案、验收矩阵
status: active
doc_type: info-collection
created: 2026-08-02
updated: 2026-08-02
parent: GOAL-001-production-admin-foundation
version: 0.1.0
related_info: I-002
related_decision: D-006
---

# I-002 · 认证方案信息收集（D-006 边界）

> **性质**：回答「认证/会话机制、凭据边界与安全配置采用什么最小方案」所需的事实，形成候选方案与验收矩阵供用户裁决。
> **不是**：已选定实现方案；R2 方案冻结；任何认证机制被表述为已选定、已验证。**`I-002` 保持 `collecting`。**
> **扫描日期**：2026-08-02（仓库只读对照；无产品代码变更）。共享资料目录 `none`，本附件全部事实来自本工作区代码与配置。

## 0. 总览结论

| 维度 | 结论 |
|------|------|
| 前端认证现状 | 启动时 GET `/api/accounts/me` 取 `{ user, features }` 作为 `$context`；**无登录/登出流程、无受保护路由、无 token/cookie 存取**；权限为渲染层门控（非安全边界） |
| 后端认证现状 | 进程内静态开发会话 `StaticDevSession`（dev-001, admin+editor）；**无请求级身份解析**；写路由 gate 绑注入的静态会话，HTTP 匿名客户端在默认进程配置下仍可写 |
| 配置现状 | 后端环境变量仅 APP_*/HTTP_*/LOG_LEVEL；**无任何认证/密钥/会话配置键** |
| 依赖现状 | 前端无 auth/cookie/router 库；后端 `go.mod` **零第三方依赖**、无 DB |
| 部署现状 | **无 Dockerfile / docker-compose / 生产静态托管 / 生产 /api 反代**；Vite `/api` 代理仅 dev；CI 仅测试+构建 |
| R2 关键开放前提 | **同源 vs 跨源托管未定**（决定 cookie 系 vs bearer/JWT 系）；会话存储（内存 vs 持久化接 R3）未定；是否接受新增依赖未定 |

## 1. 前端现状事实（apps/web）

| 事实 | 证据 |
|------|------|
| 技术栈：React 19 + Vite 6 + TS 5.8 + Tailwind 4；**无路由库**（History API + manifest route match）、无状态管理库 | `apps/web/package.json`；`apps/web/src/app/App.tsx` `useEffect`/`popstate` |
| runtime 依赖仅 UI 基础件：`@radix-ui/react-slot`、`class-variance-authority`、`clsx`、`lucide-react`、`tailwind-merge`；dev 依赖 playwright / ajv / jsdom / vitest / vite / tailwindcss | `apps/web/package.json` |
| 启动即 `loadAccountContext()` → GET `/api/accounts/me` → 返回 `{ user, features }` → 作为 `NavigationContext` 注入 `App`；失败 → `console.error` + 非阻塞 banner，导航/动作 fail-closed | `apps/web/src/main.tsx` L25-42；`apps/web/src/account/context.ts` L11-34；`apps/web/src/app/App.tsx` L401-408 |
| **无登录/登出 UI 或流程**：iconRegistry 有 `logout` icon，但仅 manifest 声明存在时渲染；无 login route、无 token/cookie 存储、无受保护路由 | `apps/web/src/app/App.tsx` L9/L44；全库 grep 无 login/logout/auth 流 |
| 权限为**渲染层门控**（导航投影 + row-action/表单禁用），客户端可见，非安全边界 | `apps/web/src/app/navigation.ts`；`apps/web/src/renderer/permissions.ts` |
| 开发代理：`server.proxy` 把 `/api` 转发到 `http://127.0.0.1:8080`（`changeOrigin`，**仅 dev**） | `apps/web/vite.config.ts` L14-23 |
| Web 无 `.env.example`；无 `import.meta.env` 认证配置读取；唯一 localStorage 用途是 `theme` | `apps/web`（无 env 文件）；`main.tsx` `applyStoredTheme` |
| 页面文档运行时源为 `GET /api/schema/{pageId}`（dev 经代理） | `apps/web/src/app/App.tsx` `SchemaPageSurface`；`apps/api/internal/handler/schema.go` |

## 2. 后端现状事实（apps/api）

| 事实 | 证据 |
|------|------|
| Go 1.26，stdlib `net/http` ServeMux（method patterns）；**`go.mod` 零 requires（无第三方依赖、无 DB/ORM/迁移）** | `apps/api/go.mod`；`apps/api/cmd/server/main.go` |
| 环境配置仅：APP_NAME / APP_ENV / HTTP_ADDR(`:8080`) / HTTP_READ·WRITE·IDLE_TIMEOUT / LOG_LEVEL。**无 SESSION_* / TOKEN_* / SECRET / JWT / cookie / 密码 配置键** | `apps/api/internal/config/config.go` L10-32；`apps/api/.env.example` |
| 会话模型 `Session{user, features}`；`StaticDevSession()` = dev-001（roles admin+editor）；provider 可注入，nil → fail-closed 401 | `apps/api/internal/account/session.go` L5-33；`apps/api/internal/handler/account.go` |
| **无请求级身份**：sessionProvider 为进程注入（生产接线 = 恒含 admin 的 `StaticDevSession`），不按请求凭证解析；HTTP 客户端无凭证在默认配置下仍可 PATCH/DELETE | `apps/api/internal/handler/account.go` L16-18；`apps/api/internal/handler/records.go` L83-114；`apps/api/README.md`「鉴权边界」 |
| 写路由 gate：有效会话 + `admin` 角色；无会话 → `401 UNAUTHENTICATED`，非 admin → `403 FORBIDDEN`；GET 只读开放。`account.Allow` 为 D-PERM 表达式求值（fail-closed） | `apps/api/internal/handler/records.go` L95-114；`apps/api/internal/account/permission.go` L146-154 |
| 无 CORS、无 cookie、无 CSRF、无安全头、无 TLS 终止、无限流；仅 PATCH body ≤4 KiB（`MaxBytesReader`）、`pageSize` ≤100 | `apps/api/internal/handler/records.go` L18-21/L214；`apps/api/internal/handler/health.go`（无中间件） |
| 路由全集：`GET /healthz`、`GET /api/accounts/me`、`GET /api/records`、`GET /api/records/{id}`、`PATCH/DELETE /api/records/{id}`、`GET /api/schema/{pageId}` | `apps/api/internal/handler/health.go` L20-25 |
| 服务器带 Read/Write/Idle timeout；SIGINT/SIGTERM 优雅关停（10s） | `apps/api/internal/server/server.go`；`apps/api/cmd/server/main.go` L44-55 |

## 3. 部署 / 工程化现状事实

| 事实 | 证据 |
|------|------|
| **无 Dockerfile、无 docker-compose、无容器/编排配置**（仓库级扫描仅 CI workflow + Makefile） | 仓库 `find`（Dockerfile/compose/yaml 命中仅 `.github/workflows/r6-basic-matrix.yml`） |
| CI 仅测试+构建+E2E：web `npm ci/test/build`；api `go test/build`；browser E2E 由 playwright webServer 拉起 Go API + Vite dev server。**无部署步骤、无制品/镜像推送、无环境变量注入** | `.github/workflows/r6-basic-matrix.yml` |
| **生产静态托管缺失**：web `dist` 无托管故事（无 Nginx / 无 Go 嵌入 / 无 CDN）；生产 `/api` 无反向代理/同源配置。Vite `/api` proxy 仅 dev | `apps/web/vite.config.ts`；`apps/web/README.md`；`apps/api`（无静态服务） |
| API 构建经 Makefile + ldflags 注入 `pkg/version`；健康检查 `GET /healthz` 返回 status/version/commit（R5 可复用） | `apps/api/Makefile`；`apps/api/internal/handler/health.go` |

## 4. 安全与约束观察（事实 + 影响判断）

- 当前「静态开发会话 + 渲染层权限」**不是生产身份源**，与 VP-002「真实认证、请求级身份」成功边界冲突——这是 R2 要消除的差量。
- 无 secret 管理、无密钥轮换、无 cookie/密码哈希逻辑。`StaticDevSession` 恒含 `admin`，默认进程配置下等于无鉴权写路径（API README 已声明）。
- **同源 vs 跨源托管未定**：决定 cookie 系（HttpOnly SameSite）与 bearer/JWT 系的取舍，也是 R5 部署基线的输入。
- 会话存储（进程内 vs 持久化）未定；R3 持久化身份（I-003 open）未实现，R2 方案必须保留与 R3 的依赖边界（D-006 边界 5），不把 R3 数据模型写成既定事实。

## 5. 可选认证方案（候选，非选定）

> 三方案均需满足 §6 验收矩阵；下述差异仅在机制与前提。**用户裁决前不冻结任何方案，`I-002` 保持 `collecting`。**

### 方案 A · HttpOnly 会话 Cookie（同源）+ 服务端会话
- **机制**：`POST /api/auth/login`（用户名+密码）→ 服务端校验 → 建立会话（进程内或 R3 存储）→ `Set-Cookie: HttpOnly; SameSite=Lax|Strict; Secure`；请求经中间件从 cookie 解析身份；登出 = 作废会话。
- **前提**：web 与 api **同源托管**（Go 嵌入 `dist` 或反代）；会话存储；CSRF 防护（SameSite + CSRF token）。
- **优点**：HttpOnly 阻 XSS 窃取；浏览器原生过期/恢复；登录/登出/过期/撤销语义直接；Go 标准库可做到（依赖可保持零）。
- **风险**：需要 CSRF 防护；会话存储与 R3 关系待定；同源托管是 R5 部署前提。
- **对照 D-006 边界**：1（生命周期齐全）/2（cookie 存储与 Secure、SameSite 约束）/3（中间件解析）/4（dev 兜底独立开关）/5（身份对象字段与 R3 兼容）。

### 方案 B · Opaque Bearer Token（`Authorization` 头）+ 服务端会话
- **机制**：登录 → 签发随机不透明 token（哈希后存内存/DB）→ 前端存储（memory/localStorage/sessionStorage，需明示取舍）→ 每次请求带 `Bearer`；中间件校验；登出=撤销；过期=TTL。
- **前提**：SPA 与 API 同源或 CORS+credentials 明确配置；token 存储安全取舍落盘。
- **优点**：SPA 友好、可跨源；撤销直接；与 R3 解耦简单；服务端仍可无依赖。
- **风险**：XSS 可窃取存储 token；localStorage 注入面；需 CORS/预检与刷新/过期补充。
- **对照 D-006 边界**：1（生命周期齐全，刷新需补充）/2（token 存储与传递边界）/3（中间件）/4（dev 兜底）/5（同 A）。

### 方案 C · 签名 JWT（无状态 access + 可选 refresh）
- **机制**：登录 → 签发 JWT（HMAC/RSA）→ 前端存储 → `Bearer`；服务端无会话存储；`exp` 控制过期；撤销需黑名单或短 TTL+refresh 轮换。
- **前提**：密钥管理（签发/轮换）；refresh 轮换策略；**Go stdlib 无 JWT，需引入第三方依赖或手写 HMAC（手写易错，不推荐）**。
- **优点**：无会话存储、水平扩展友好。
- **风险**：撤销难；密钥管理复杂；XSS 窃取；引入新依赖与当前零依赖现状冲突。
- **对照 D-006 边界**：1（过期/撤销需补充机制）/2（token 边界）/3（中间件验签）/4（dev 兜底）/5（同 A）。

## 6. 验收矩阵（R2 阶段验收建议 · 方案无关，均须满足）

> 对照 VP-002 阶段 2 验收建议 与 Root 成功边界。判定均为**可自动/手动核对的通过条件**；每项须在对应子目标实施后留证据。

| # | 验收项 | 可验证行为 | pass 判定 |
|---|--------|-----------|-----------|
| M1 | 登录成功 | 正确凭据提交后 | 会话建立、前端进入已认证态、`/api/accounts/me` 返回请求身份 |
| M2 | 登录失败 | 错误/缺失凭据提交后 | 返回认证失败（401 语义）、不建立会话、错误可观察且不泄露细节 |
| M3 | 登出 | 登出动作后 | 会话被作废，后续请求未认证（401），前端回未登录态 |
| M4 | 会话恢复 | 刷新/重开页面 | 凭据（cookie/token）存在时自动恢复身份；无凭据 → 未登录态，不误入已认证 |
| M5 | 过期 | 超过会话/Token TTL | 请求返回 401，前端转未登录或重新登录路径 |
| M6 | 撤销 | 服务端显式撤销后 | 该会话/token 立即失效（A 撤销会话；B 撤销 token；C 走黑名单或 refresh 轮换） |
| M7 | 重启后身份 | 重启 API 服务 | 身份符合既定持久化策略（进程内会话重启失效 vs 持久化恢复），并有明确断言 |
| M8 | 未登录 → 401 / 无权限 → 403 | 匿名与无权限请求 | 后端**按请求身份**（非静态进程注入）返回 `401` / `403`，业务路由统一经身份中间件 |
| M9 | 静态开发会话仅 dev | 生产配置下启动 | 静态 dev session 只作为显式 `opt-in` 本地兜底，生产默认**不**启用 |
| M10 | 请求级身份中间件 | 业务 handler 读取请求 | 身份解析自请求（cookie/token），不再是 `StaticDevSession` 进程注入 |
| M11 | 凭据边界 | 传输与存储 | 密码/令牌经 HTTPS 传输；token/cookie 存储策略落盘；无硬编码 secret；config 仅经 env 注入 |
| M12 | 与 R3 依赖边界 | 认证产物 → 身份模型 | 认证返回的身份对象字段（id/name/roles）与 R3 用户—角色模型兼容，但**不**实现 R3 数据模型 |
| M13 | 安全配置边界 | dev/prod 差异 | 环境区分明确；secret 与密钥从 env 注入且不落仓库；生产不默认开放静态会话 |
| M14 | 自动化可复现 | CI/E2E | login→me→401/403 路径有自动化测试且 CI 可复现（如浏览器 E2E 扩展） |

> 备注：`401`/`403` 的**责任边界**（谁返回、错误码语义、前端如何转跳）由后续 R2 方案决策固化；本矩阵只固定「按请求身份返回」的可核对条件。

## 7. 待用户裁决点（下一步，P-004 信息取舍）

1. **同源 vs 跨源托管**：决定方案 A 可行性与 B/C 的 CORS 处理，同时是 R5 部署基线输入（I-005）。
2. **会话存储**：进程内（MVP，重启失效）vs 持久化（接 R3）；若接 R3 需明确不越权实施 R3。
3. **依赖策略**：是否接受新增认证依赖（方案 C 的 JWT 库）vs 保持 Go 零依赖（A/B）。
4. **凭据与 secret 管理边界**：密码存储、secret 注入、密钥轮换的 M11 落地口径。
5. 裁决后：记录决策（02）→ `I-002` 置 `verified` → 冻结 R2 方案边界 → 再创建 R2 子目标。本附件不自行冻结。

## 8. 证据索引（只读核对路径）

- 前端：`apps/web/package.json`、`src/main.tsx`、`src/account/context.ts`、`src/app/App.tsx`、`src/app/navigation.ts`、`src/renderer/permissions.ts`、`vite.config.ts`、`README.md`
- 后端：`apps/api/go.mod`、`internal/config/config.go`、`internal/account/session.go`、`internal/account/permission.go`、`internal/handler/account.go`、`internal/handler/records.go`、`internal/handler/health.go`、`internal/server/server.go`、`cmd/server/main.go`、`.env.example`、`README.md`、`Makefile`
- 部署/CI：`.github/workflows/r6-basic-matrix.yml`；仓库级 `find`（无 Dockerfile/compose）
