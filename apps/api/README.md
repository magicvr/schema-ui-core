# apps/api · schema-ui-core API 骨架

MVP Admin 基架的 Go 服务（GOAL-003 骨架 + GOAL-006/007 账号权限与记录示例域）。**演示 API：无生产鉴权，见下方「鉴权边界」。**

## 要求

- Go **1.26+**（本仓 R1 在 Windows 实测 `go1.26.0`）
- Module：`github.com/magicvr/schema-ui-core/apps/api`

## 布局

```text
cmd/server/          # 进程入口
internal/config/     # 环境配置
internal/server/     # http.Server 包装
internal/handler/    # HTTP 路由（healthz / accounts / records）
internal/account/    # 会话模型与权限求值库（D-004 / D-PERM）
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

默认监听 `:8080`（`HTTP_ADDR`）。

## 端点

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/healthz` | 探活 |
| GET | `/api/accounts/me` | R4 会话快照（静态开发会话 `dev-001`，`roles: [admin, editor]`） |
| GET | `/api/records` | R5 D-DATA 列表：`q` / `sort`（name·status·owner·updatedAt）/ `order` / `page` / `pageSize` |
| GET | `/api/records/{id}` | R5 D-DATA 详情 |
| PATCH | `/api/records/{id}` | R5 D-ACT 编辑（name/status/owner；成功刷新 `updatedAt`；需 admin 会话） |
| DELETE | `/api/records/{id}` | R5 D-ACT 删除（需 admin 会话） |

## 鉴权边界（MVP 声明，非生产）

- `/api/accounts/me` 返回**静态开发会话**（`account.StaticDevSession`，roles: admin+editor）；nil 会话提供者按 fail-closed 返回 `401 UNAUTHENTICATED`。
- **`/api/records` 写路由（PATCH/DELETE）fail-closed 鉴权**：需要有效会话且会话须含 `admin` 角色（`account.Allow`）；无会话 → `401 UNAUTHENTICATED`，非 admin → `403 FORBIDDEN`。GET 只读路由保持开放。
- **范围说明**：该 gate 绑定**进程内注入的会话提供者**（生产接线为 `StaticDevSession`，恒含 admin），**非**按 HTTP 请求凭证/令牌鉴权——匿名 HTTP 客户端在默认进程配置下仍可 PATCH/DELETE 成功。这是 MVP 静态会话的边界，不是网络侧身份鉴权；生产化需真实登录/令牌。
- **请求上限（F-009-007）**：PATCH body ≤ 4 KiB（`MaxBytesReader`）；`pageSize` ≤ 100，超限返回 `400 INVALID_PAGE_SIZE`。

## 测试

```bash
make test
# 或
go test ./...
```

## 非目标（MVP）

- 订单 / 钱包 / 通知等业务 API
- 真实 login/logout/token/IAM（R4 可选，非当前必做）
- 完整协议兼容主张（见仓库愿景与 Root 信息门禁）
