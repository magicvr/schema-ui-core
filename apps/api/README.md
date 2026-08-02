# apps/api · schema-ui-core API 骨架

MVP Admin 基架的 Go 服务（GOAL-003 骨架 + GOAL-006/007 账号权限与记录示例域）。**R2 起提供真实认证（GOAL-005）**：短 JWT Access + Opaque Refresh + SQLite；**R3 起持久化 RBAC（GOAL-006）**：records 路由按 `records.read` / `records.write` 权限键授权，`/api/accounts/me` 返回持久化菜单投影；**R4 起 records 数据源为 SQLite（GOAL-007）**：`0003 records_persist` 迁移 + `seedRecords` 空表种子，`POST` 新增，进程内切片不再作为生产路径。

## 要求

- Go **1.26+**（本仓 R1 在 Windows 实测 `go1.26.0`）
- Module：`github.com/magicvr/schema-ui-core/apps/api`

## 布局

```text
cmd/server/          # 进程入口（配置解析、store 打开、认证接线）
internal/config/     # 环境配置（含 R2 认证配置键）
internal/server/     # http.Server 包装
internal/handler/    # HTTP 路由（healthz / auth / accounts / records / schema）
internal/auth/       # R2 认证核心：JWT、refresh 轮换、bcrypt、请求身份中间件
internal/account/    # 会话模型与权限求值库（D-004 / D-PERM）
internal/store/      # SQLite 认证 + R3 RBAC / 迁移存储（users、refresh_tokens、schema_migrations、roles/…）
pkg/version/         # 构建版本变量
```

## 运行

```bash
# 可选
cp .env.example .env

make run
# 或
go run ./cmd/server
```

默认监听 `:8080`（`HTTP_ADDR`）。首次启动在 `DB_PATH`（默认 `./data/schema-ui.db`）建表并种子 admin：
- dev 缺省 `ADMIN_INITIAL_PASSWORD=admin`；生产必须显式设置。
- 生产缺少 `AUTH_JWT_SECRET` → 启动失败（fail-closed）；dev 使用开发密钥并打警告。

## 配置键（R2 · GOAL-005 D-004）

| 键 | 默认 | 说明 |
|----|------|------|
| `AUTH_JWT_SECRET` | dev 开发密钥 | access token 签发密钥；生产必填，缺失 fail-closed |
| `AUTH_ACCESS_TTL` | `15m` | access token 时效 |
| `AUTH_REFRESH_TTL` | `720h` (30d) | refresh token 时效 |
| `DB_PATH` | `./data/schema-ui.db` | SQLite 路径 |
| `ADMIN_INITIAL_PASSWORD` | dev `admin` | 首次种子 admin 密码；生产必填 |
| `AUTH_DEV_SESSION_ENABLED` | `false` | 显式本地开发静态会话兜底；**生产禁止启用** |

## 端点

| 方法 | 路径 | 鉴权 | 说明 |
|------|------|------|------|
| GET | `/healthz` | 公开 | 探活 |
| POST | `/api/auth/login` | 公开 | 校验用户名/密码，签发 access + refresh |
| POST | `/api/auth/refresh` | 公开（需有效 refresh） | 轮换 refresh，签发新 access + refresh |
| POST | `/api/auth/logout` | 公开 | 撤销 refresh（幂等） |
| GET | `/api/accounts/me` | **Bearer** | 返回请求身份的会话快照 `{ user, features }` |
| GET | `/api/records` | **Bearer + records.read** | R5 D-DATA 列表：`q` / `sort` / `order` / `page` / `pageSize`（需 `records.read` 权限） |
| POST | `/api/records` | **Bearer + records.write** | R4 新建（body `name`/`status`/`owner` 必填）→ 201 + 完整 record；`INVALID_CREATE_*` 400 |
| GET | `/api/records/{id}` | **Bearer + records.read** | R5 D-DATA 详情（需 `records.read` 权限） |
| PATCH | `/api/records/{id}` | **Bearer + records.write** | R5 D-ACT 编辑（name/status/owner；需 `records.write` 权限） |
| DELETE | `/api/records/{id}` | **Bearer + records.write** | R5 D-ACT 删除（需 `records.write` 权限） |
| GET | `/api/schema/{pageId}` | 公开（只读） | 页面 Schema 文档 |

## 鉴权边界（R3 · 真实认证 + 权限键）

- **请求级身份**：受保护路由经 `Authorization: Bearer <access>` 由中间件解析身份；不再依赖进程注入的 `StaticDevSession`。
  - 无 / 无效 / 过期 access → `401 UNAUTHENTICATED`；已认证但缺少所需权限键 → `403 FORBIDDEN`。
- **记录权限（R3 · GOAL-006）**：`GET /api/records` 与 `GET /api/records/{id}` 门禁 `records.read`；`POST` / `PATCH` / `DELETE`（含 `/api/records/{id}`）门禁 `records.write`。种子 admin 持 read + write，editor / viewer 仅 read。
- **记录数据源（R4 · GOAL-007）**：`records` 存于 SQLite（`0003 records_persist`，`updated_at` Unix 毫秒，API 序列化为 RFC3339 含毫秒）；首次 `seedAdmin=true` 且表为空时 `seedRecords` 插入 8 条演示数据，非空则跳过。写并发为单写者 last-write-wins，同一毫秒内连续 update 由单调钳制保证 `updatedAt` 严格递增。
- **Access**：短时效 JWT（`AUTH_ACCESS_TTL`，默认 15m），负载为 `sub`（用户 id）。
- **Refresh**：opaque 随机串，**SHA-256 哈希存 SQLite**；登出/刷新即撤销（轮换）；过期/撤销 → `401`。
- **静态开发会话**：仅 `AUTH_DEV_SESSION_ENABLED=true` 时作为显式本地兜底（替换 401）；生产默认关闭（M9）。
- **凭据**：bcrypt（cost 10）哈希存储；密码/密钥不落仓库；生产 `AUTH_JWT_SECRET` 缺失 fail-closed。
- **请求上限（F-009-007）**：auth/records 写 body ≤ 4 KiB（`MaxBytesReader`）；`pageSize` ≤ 100。

## 测试

```bash
make test
# 或
go test ./...
```

覆盖：auth 生命周期（登录/刷新轮换/登出/过期/撤销）、请求身份 401/403、store（种子幂等、token 生命周期、0003 迁移与 records 持久化/毫秒往返）、records 读/写权限门禁与 create/list/detail/PATCH/DELETE、schema 文档读取。

## 非目标（当前 R4 边界）

- R2 阶段的用户/角色/菜单持久化非目标已由 **R3（GOAL-006）** 落地为持久化 RBAC（版本迁移 + 权限键授权 + 菜单投影）；本 README 的鉴权描述均按当前 R3 边界。
- 完整 IAM（SSO / SCIM / 复杂策略）；`D-UPLOAD`；完整协议兼容主张（见仓库愿景与 Root 信息门禁）。
