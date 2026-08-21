# apps/api · schema-ui-core API 骨架

MVP Admin 基架的 Go 服务（workspace-001 历史目标编号：GOAL-003 骨架 + GOAL-006/007 账号权限与语义资源域）。**R2 起提供真实认证（GOAL-005）**：短 JWT Access + Opaque Refresh + SQLite；**R3 起持久化 RBAC（GOAL-006）**：资源路由按 `users.read`/`users.write`、`roles.read`/`roles.write` 权限键授权，`/api/accounts/me` 返回持久化菜单投影；**GOAL-011 起 users/roles 为语义资源（records 已按 0006 退场）**：`/api/users`、`/api/roles` 由 owner module Provider 贡献，store 只负责平台迁移/执行与恢复边界。（注：上述 GOAL-00N 均为 workspace-001 mvp-admin-foundation 的历史目标编号，与本仓库 workspace-003 的同名编号无关。）

## 要求

- Go **1.26+**（本仓 R1 在 Windows 实测 `go1.26.0`）
- Module：`github.com/magicvr/schema-ui-core/apps/api`

## 布局

```text
cmd/server/          # 进程入口（配置解析、store 打开、认证接线）
internal/config/     # 环境配置（含 Profile/认证配置键）
internal/server/     # http.Server 包装
internal/kernel/     # 框架无关模块契约、Profile 与依赖图
internal/composition/# Fx 组合根与模块 Provider 接线
internal/handler/    # core HTTP 路由（healthz / auth / accounts / schema / manifest）
internal/auth/       # R2 认证核心：JWT、refresh 轮换、bcrypt、请求身份中间件
internal/account/    # 会话模型与权限求值库（D-004 / D-PERM）
internal/store/      # SQLite 平台执行、全局迁移台账、快照与恢复边界
internal/modules/    # users/roles/settings/activity 与 core Schema owner modules
pkg/version/         # 构建版本变量
```

## 运行

```bash
# 配置权威是 configs/config.yaml（W7）：非敏感值直接写 YAML；敏感值写 ${VAR}
# 占位符，真实值来自 configs/.env（开发，gitignored）或进程 env（生产）。
# 已设置的进程 env 总是覆盖 YAML。Compose 路径由仓库根 .env 提供插值。
# 模块启用集只认 configs/config.yaml（T-06）：app.profile 或 app.modules（preset / list）

make run
# 或
go run ./cmd/server
```

默认监听 `:25080`（`HTTP_ADDR`）。首次启动在 `DB_PATH`（默认 `./data/schema-ui.db`）建表并种子 admin：
- dev 缺省 `ADMIN_INITIAL_PASSWORD=admin`；生产必须显式设置。
- 生产缺少 `AUTH_JWT_SECRET` → 启动失败（fail-closed）；dev 使用开发密钥并打警告。

## 配置键（R2 · GOAL-005 D-004）

| 键 | 默认 | 说明 |
|----|------|------|
| `AUTH_JWT_SECRET` | dev 开发密钥 | access token 签发密钥；生产必填，缺失 fail-closed |
| `AUTH_ACCESS_TTL` | `15m` | access token 时效 |
| `AUTH_REFRESH_TTL` | `720h` (30d) | refresh token 时效 |
| `DB_PATH` | `./data/schema-ui.db` | SQLite 路径 |
| `app.profile` (YAML) | `mvp` | `mvp`、`admin`、`demo`；无 `app.modules` 时选内置预设 |
| `app.modules` (YAML) | 无 | `preset`（内置名或预设文件路径）或内联 `list`；互斥；覆盖 Profile 默认集合 |
| `ADMIN_INITIAL_PASSWORD` | dev `admin` | 首次种子 admin 密码；生产必填 |
| `AUTH_DEV_SESSION_ENABLED` | `false` | 显式本地开发静态会话兜底；**生产禁止启用** |

### 开发 vs 生产（GOAL-008 S1 / I-008-001）

| 维度 | 开发（`APP_ENV=development`） | 生产（`APP_ENV=production` / compose） |
|------|-------------------------------|-----------------------------------------|
| `AUTH_JWT_SECRET` | 未设则用内建 dev 密钥并打警告 | **必填，缺失 fail-closed** |
| `ADMIN_INITIAL_PASSWORD` | 未设则兜底 `admin` | **必填，缺失 fail-closed** |
| `AUTH_DEV_SESSION_ENABLED` | 显式 opt-in 可选 | **必须 `false`** |
| `DB_PATH` | `./data/schema-ui.db` | compose 挂载 `/app/data/schema-ui.db`（命名卷） |
| 启动形态 | 本地双进程（api + web） | `docker compose up`（第二启动路径；fork 用户二者可选） |

Profile 选择只影响启动时模块集合，不改变编译产物或全局迁移台账：`mvp` 包含 users/roles，
`admin` 另外包含 settings/activity，`demo`（W2 · **非生产向演示 Profile**）= mvp 集 + `dev.examples`
（启动即展示 8 个协议范例页 + Examples 导航，home 指向 `overview`），`custom` 必须显式提供完整
依赖闭包。`app.modules` 的优先级高于 Profile 默认值；未知、重复或缺依赖模块会 fail-closed。
生产只应使用 `mvp` / `admin`；`demo` 用于开发/演示，不得作为生产默认。

完整契约见 GOAL-008 `attachments/I-008-001-engineering-contract.md`。

### 启动与健康验证（C-002）

```bash
# 1) 启动 API（compose 或本地）
docker compose up -d api   # 或：make run / go run ./cmd/server
# 2) 探活
curl -fsS http://localhost:25080/healthz
# -> {"status":"ok","timestamp":"...","version":"...","commit":"..."}
# 3) 登录种子 admin（首次启动按 ADMIN_INITIAL_PASSWORD 种子）
curl -fsS -X POST http://localhost:25080/api/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"<ADMIN_INITIAL_PASSWORD>"}'
# -> {"accessToken":"...","refreshToken":"..."}
# 4) 会话
TOKEN=$(...); curl -fsS http://localhost:25080/api/accounts/me -H "Authorization: Bearer $TOKEN"
# -> {"user":{...},"features":{...}}
```

- `GET /healthz` 公开返回 `200 {"status":"ok",...}`，作为 liveness 探活与启动验证判据（不访问数据库）。
- `GET /readyz` 公开返回 `200 {"status":"ok",...}`，为 readiness 就绪探针：在 liveness 之上执行轻量 SQLite 读，**并**仅在模块图 Start+Ready 全部成功后返回 `200`（R5 真实模块图 readiness；未就绪返回 `503 {"status":"not-ready",...}`，数据库不可读返回 `503 {"status":"unavailable",...}`）；Compose 以它作为 `service_healthy`。
- API 优雅停机：`SIGINT`/`SIGTERM` → 10s 宽限内 `Shutdown`。

## 端点

| 方法 | 路径 | 鉴权 | 说明 |
|------|------|------|------|
| GET | `/healthz` | 公开 | 探活（liveness） |
| GET | `/readyz` | 公开 | 就绪（readiness，含 SQLite 检查；容器探针） |
| POST | `/api/auth/login` | 公开 | 校验用户名/密码，签发 access + refresh |
| POST | `/api/auth/refresh` | 公开（需有效 refresh） | 轮换 refresh，签发新 access + refresh |
| POST | `/api/auth/logout` | 公开 | 撤销 refresh（幂等） |
| GET | `/api/accounts/me` | **Bearer** | 返回请求身份的会话快照 `{ user, features }` |
| GET | `/api/users` | **Bearer + users.read** | 用户列表：`q` / `sort` / `order` / `page` / `pageSize` |
| POST | `/api/users` | **Bearer + users.write** | 新建用户（`username`/`name`/`password` 必填；`roles` 可选数组）→ 201；`INVALID_CREATE_*` 400 / `USERNAME_TAKEN` 409 |
| GET | `/api/users/{id}` | **Bearer + users.read** | 用户详情（永不返回 `password_hash`） |
| PATCH | `/api/users/{id}` | **Bearer + users.write** | 编辑（`name`/`password`/`roles`）；self/last-admin 保护 409 |
| DELETE | `/api/users/{id}` | **Bearer + users.write** | 删除用户（self → 409；原子撤销 refresh token） |
| GET | `/api/roles` | **Bearer + roles.read** | 角色列表 |
| POST | `/api/roles` | **Bearer + roles.write** | 新建角色（`key`/`name`；key 格式校验）；`ROLE_KEY_TAKEN` 409 |
| GET | `/api/roles/{id}` | **Bearer + roles.read** | 角色详情（`system` 布尔） |
| PATCH | `/api/roles/{id}` | **Bearer + roles.write** | 编辑 `name`（key 不可变；system 角色 409） |
| DELETE | `/api/roles/{id}` | **Bearer + roles.write** | 删除角色（system / 被用户使用 → 409） |
| GET | `/api/schema/{pageId}` | 公开（只读） | 页面 Schema 文档 |
| GET | `/.well-known/schema-ui/app-manifest.json` | 公开（只读） | API 聚合的 Profile Manifest（生产唯一来源） |
| GET | `/api/settings`、`/api/branding` | **admin Profile** | Settings 模块路由；未启用时为 404 |
| GET | `/api/operations` | **admin Profile** | Activity 查询路由；未启用时为 404 |

## 鉴权边界（R3 · 真实认证 + 权限键）

- **请求级身份**：受保护路由经 `Authorization: Bearer <access>` 由中间件解析身份；不再依赖进程注入的 `StaticDevSession`。
  - 无 / 无效 / 过期 access → `401 UNAUTHENTICATED`；已认证但缺少所需权限键 → `403 FORBIDDEN`。
- **资源权限（GOAL-011）**：users/roles 五路由经通用工厂 `requirePermission` 门禁 `users.read`/`users.write`、`roles.read`/`roles.write`。种子 admin 持 read + write，editor / viewer 仅 read。
- **语义资源（GOAL-011 S2/S3）**：`/api/users`、`/api/roles` 走通用资源工厂（`resources.go`）；users 敏感字段隔离（`password_hash` 永不出响应）、角色分配双写（legacy JSON + `user_roles`，不隐式建角色）、self/last-admin 保护；roles system 角色（种子 admin/editor/viewer）不可改删、in-use 保护。records 已由 `0006 records_retire` 从产品运行面退场，仅保留历史迁移账本与 `records.*` operation-log 兼容值（`historical 0006`，不恢复产品 CRUD）。
- **Access**：短时效 JWT（`AUTH_ACCESS_TTL`，默认 15m），负载为 `sub`（用户 id）。
- **Refresh**：opaque 随机串，**SHA-256 哈希存 SQLite**；登出/刷新即撤销（轮换）；过期/撤销 → `401`。
- **静态开发会话**：仅 `AUTH_DEV_SESSION_ENABLED=true` 时作为显式本地兜底（替换 401）；生产默认关闭（M9）。
- **凭据**：bcrypt（cost 10）哈希存储；密码/密钥不落仓库；生产 `AUTH_JWT_SECRET` 缺失 fail-closed。
- **请求上限（F-009-007）**：资源写 body ≤ 4 KiB（`MaxBytesReader`）；`pageSize` ≤ 100。

## 测试

```bash
make test
# 或
go test ./...
go vet ./...
go build ./...
```

覆盖：auth 生命周期（登录/刷新轮换/登出/过期/撤销）、请求身份 401/403、store（种子幂等、token 生命周期、迁移链 0001～0006 与 users/roles 持久化/毫秒往返）、users/roles 读/写权限门禁与 create/list/detail/PATCH/DELETE + 领域保护（self/last-admin/system/in-use）、schema 文档读取。

### 升级与恢复

非空数据库在有待执行的数据迁移时，会在 `DB_PATH` 旁创建并完整性校验
`<db>.pre-vNNNN-<UTC>.sqlite` 快照；fresh bootstrap 不创建快照。升级前应停止写入并保留
数据库副本。迁移失败时保留失败文件，停机后将选定快照复制回 `DB_PATH`，再使用已验证的
旧二进制/镜像启动；不要手工编辑 `schema_migrations`。Profile 切换不删除禁用模块的表或数据。

## 非目标（当前 R4 边界）

- R2 阶段的用户/角色/菜单持久化非目标已由 **R3（GOAL-006）** 落地为持久化 RBAC（版本迁移 + 权限键授权 + 菜单投影）；本 README 的鉴权描述均按当前 R3 边界。
- 完整 IAM（SSO / SCIM / 复杂策略）；`D-UPLOAD`；完整协议兼容主张（见仓库愿景与 Root 信息门禁）。
