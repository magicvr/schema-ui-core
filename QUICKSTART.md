# QUICKSTART · 15 分钟 fork 上手（R5 · GOAL-008 S3）

> 面向 fork 使用者的最小上手段。目标：**按本文档从零配置并启动，≤15 分钟（不含依赖下载/镜像拉取）进入系统**——登录成功 + 后台可交互。
> 本文件的终点与计时口径按 [I-008-002 fork 复现协议 v0.1.1](docs/workspace-002-production-admin-foundation/GOAL-008-r5-engineering-fork/attachments/I-008-002-fork-reproduction-protocol.md) 冻结；独立复现记录见 GOAL-008 `02-execution.md`。

## 0. 前置

| 项 | 要求 |
|----|------|
| 工具 | Go **1.26+**、Node **22+**、npm **10+**；或 Docker **20.10+** + Compose **v2.20+**（二选一启动路径） |
| 数据库 | 无需手动安装；SQLite 内嵌于 API |
| 依赖缓存（不计时） | `go mod download`、`npm ci`、Compose 镜像 pull/build 缓存**提前完成**，否则耗时计入仍 ≤15 分钟 |

## 1. 获取并准备

```bash
git clone <your-fork-url> && cd schema-ui-core
git checkout <待测 ref>        # 记录实际 ref；工作树保持 clean
```

创建密钥（本地 fork 开发）：

```bash
cp apps/api/.env.example apps/api/.env   # 或按需写入仓库根 .env（gitignored）
```

- 开发（`APP_ENV=development`）不要求显式密钥；生产（compose）必须提供 `AUTH_JWT_SECRET` 与 `ADMIN_INITIAL_PASSWORD`（缺省 fail-closed 启动失败）。
- 首次启动自动建表并种子 `admin` 用户与 8 条演示 records。

## 2. 启动（两条路径选一）

### 路径 A · Docker Compose（推荐，无需本机 Go/Node）

```bash
# 仓库根 .env（gitignored）写入，避免新 shell 重复 export：
#   AUTH_JWT_SECRET=<强随机串>
#   ADMIN_INITIAL_PASSWORD=<初始 admin 密码>
docker compose up -d --build
```

- API：`http://localhost:8080`（`GET /healthz` 探活）
- Web：`http://localhost:8081`（nginx 服务 SPA，`/api` 同源反代到 API）
- 停止：`docker compose down`（SQLite 数据由命名卷 `db-data` 保持）

### 路径 B · 本地双进程（开发默认）

```bash
# 终端 1 —— API
cd apps/api && go run ./cmd/server        # 监听 :8080

# 终端 2 —— Web
cd apps/web && npm ci && npm run dev      # 监听 ${WEB_PORT:-5173}
```

## 3. 验收四终点（≤15 分钟达标判据）

| 终点 | 检查 | 达标 |
|------|------|------|
| 1 | `GET ${API_BASE_URL}/healthz` | HTTP 200，JSON `status: "ok"` |
| 2 | `POST ${WEB_BASE_URL}/api/auth/login`（admin / `ADMIN_INITIAL_PASSWORD`） | HTTP 200，响应含非空 `accessToken` |
| 3 | 携带 token `GET ${WEB_BASE_URL}/api/accounts/me` | HTTP 200，含 `user` 与 `features` |
| 4 | **浏览器**登录后打开 `${WEB_BASE_URL}/list-edit-lifecycle` | 页面标题 `List + edit lifecycle`，列表已加载 `Acme Console`（种子记录） |

> **默认 base URL**：Compose → API `http://localhost:8080`、Web `http://localhost:8081`；本地双进程 → API `:8080`、Web `http://localhost:${WEB_PORT:-5173}`。以实际端口为准，不得用默认值覆盖实测端口。

### 命令行冒烟（终点 1–3 快速验证）

```bash
curl -fsS http://localhost:8080/healthz
TOKEN=$(curl -fsS -X POST http://localhost:8081/api/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"<ADMIN_INITIAL_PASSWORD>"}'
  | node -e 'process.stdin.on("data",d=>process.stdout.write(JSON.parse(d).accessToken))')
curl -fsS http://localhost:8081/api/accounts/me -H "Authorization: Bearer $TOKEN"
```

### 完整 smoke（S4 机器可判定）

```bash
# 对已启动实例做非破坏性检查（SM-001～005）
SMOKE_USERNAME=admin SMOKE_PASSWORD=<ADMIN_INITIAL_PASSWORD> bash scripts/smoke.sh
# 含种子可重复性（SM-006，disposable/隔离环境）
SMOKE_USERNAME=admin SMOKE_PASSWORD=<ADMIN_INITIAL_PASSWORD> bash scripts/smoke.sh --disposable
```

## 4. 下一步：接业务

- 新增业务页面：在 `docs/schemas` 添加页面 Schema，并在 Web `app-manifest.json` 挂载路由——**无需修改 Renderer 主路径**。
- 权限：编辑持久化 RBAC 种子（角色/权限键），见 `apps/api/internal/store` 种子与迁移。
- 参考：[README.md](README.md) 工程化段；`apps/api/README.md` / `apps/web/README.md` 端点与配置表。
