# QUICKSTART · 15 分钟 fork 上手（R5 · GOAL-008 S3）

> 面向 fork 使用者的最小上手段。目标：**按本文档从零配置并启动，≤15 分钟（不含依赖下载/镜像拉取）进入系统**——登录成功 + 后台可交互。
> 本文件的终点与计时口径按 [I-008-002 fork 复现协议 v0.1.1](docs/workspaces/workspace-002-production-admin-foundation/GOAL-008-r5-engineering-fork/attachments/I-008-002-fork-reproduction-protocol.md) 冻结；独立复现记录见 GOAL-008 `02-execution.md`。

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

配置环境（本地 fork 开发）：

```bash
# apps/api/.env.example 只是配置参考；Go API 不会自动加载该文件。
# Compose 才会读取仓库根 .env（gitignored）。本地进程请 export：
export APP_PROFILE=mvp                 # 或 admin
# custom 时还必须提供完整的显式模块列表：
# export APP_PROFILE=custom
# export APP_MODULES_ENABLED=core.server-registration,...
```

- 开发（`APP_ENV=development` 显式设置）不要求显式密钥；生产（compose）必须提供 `AUTH_JWT_SECRET` 与 `ADMIN_INITIAL_PASSWORD`（缺省 fail-closed 启动失败）。**`APP_ENV` 必须显式设置**——未设置时启动失败（C3：不静默回退到公开开发密钥/密码）。
- `APP_PROFILE` 只接受 `mvp`、`admin`、`custom`；`APP_MODULES_ENABLED` 非空时覆盖 Profile 默认集合。
- PowerShell 等价写法为 `$env:APP_PROFILE="mvp"`；每个本地 API/Web 进程都必须继承同一 Profile。
- 首次启动自动建表并种子 `admin` 用户与系统角色（GOAL-011：users/roles 语义资源；records 已按版本化迁移 `0006` 退场）。

## 2. 启动（两条路径选一）

### 路径 A · Docker Compose（推荐，无需本机 Go/Node）

```bash
# 仓库根 .env（gitignored）写入，避免新 shell 重复 export：
#   AUTH_JWT_SECRET=<强随机串>
#   ADMIN_INITIAL_PASSWORD=<初始 admin 密码>
#   APP_PROFILE=mvp                 # 或 admin
#   APP_MODULES_ENABLED=            # 可选，逗号分隔
docker compose up -d --build
```

- API：`http://localhost:25080`（`GET /healthz` 探活）
- Web：`http://localhost:25081`（nginx 服务 SPA，`/api` 同源反代到 API）
- 停止：`docker compose down`（SQLite 数据由命名卷 `db-data` 保持）

### 路径 B · 本地双进程（开发默认）

```bash
# 终端 1 —— API
cd apps/api && APP_ENV=development APP_PROFILE=mvp go run ./cmd/server  # 监听 :25080；或改为 admin

# 终端 2 —— Web
cd apps/web && npm ci && npm run dev      # 监听 ${WEB_PORT:-25173}
```

## 3. 验收四终点（≤15 分钟达标判据）

| 终点 | 检查 | 达标 |
|------|------|------|
| 1 | `GET ${API_BASE_URL}/healthz` 与 `/readyz` | HTTP 200，JSON `status: "ok"` |
| 2 | `POST ${WEB_BASE_URL}/api/auth/login`（admin / `ADMIN_INITIAL_PASSWORD`） | HTTP 200，响应含非空 `accessToken` |
| 3 | 携带 token `GET ${WEB_BASE_URL}/api/accounts/me` | HTTP 200，含 `user` 与 `features` |
| 4 | **浏览器**登录后打开 `${WEB_BASE_URL}/users` | 页面标题 `Users`，列表已加载 `admin` 种子用户（users 资源 CRUD） |

> **默认 base URL**：Compose → API `http://localhost:25080`、Web `http://localhost:25081`；本地双进程 → API `:25080`、Web `http://localhost:${WEB_PORT:-25173}`。以实际端口为准，不得用默认值覆盖实测端口。

### 命令行冒烟（终点 1–3 快速验证）

```bash
curl -fsS http://localhost:25080/healthz
TOKEN=$(curl -fsS -X POST http://localhost:25081/api/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"<ADMIN_INITIAL_PASSWORD>"}'
  | node -e 'process.stdin.on("data",d=>process.stdout.write(JSON.parse(d).accessToken))')
curl -fsS http://localhost:25081/api/accounts/me -H "Authorization: Bearer $TOKEN"
```

### 完整 smoke（S4 机器可判定）

```bash
# 对已启动实例做非破坏性部分检查（SM-001～005）
# 注意：SM-006 未运行 → 退出码 8（部分绿），**不是** S4 完整绿
SMOKE_USERNAME=admin SMOKE_PASSWORD=<ADMIN_INITIAL_PASSWORD> SMOKE_EXPECTED_PROFILE=mvp bash scripts/smoke.sh

# S4 完整绿（含 SM-006 种子可重复性）——必须运行在显式隔离环境：
#   1) 用独立 compose project 启动（不得指向普通开发库）：
#      docker compose -p ci-smoke-local down -v && docker compose -p ci-smoke-local up -d
#   2) 提供隔离身份 + 书面确认标记（脚本机器校验 project/卷绑定，不满足 → exit 2）
SMOKE_USERNAME=admin SMOKE_PASSWORD=<ADMIN_INITIAL_PASSWORD> SMOKE_EXPECTED_PROFILE=mvp \
SMOKE_ISOLATION_ID=ci-smoke-local SMOKE_DISPOSABLE_CONFIRM=yes \
bash scripts/smoke.sh --disposable
```

> 退出码：`0`=完整绿（含 disposable SM-006）｜`2`=参数/工具/安全前提（隔离校验失败等）｜`3`=readiness 超时｜`4`=登录/身份｜`5`=路由/数据｜`6`=种子断言｜`8`=部分绿（非 disposable）｜`70`=内部错误。判据见 [I-008-002 协议 v0.1.2](docs/workspaces/workspace-002-production-admin-foundation/GOAL-008-r5-engineering-fork/attachments/I-008-002-fork-reproduction-protocol.md) §5.3。

## 4. 升级与恢复边界

- API 启动时会为每个待执行的数据迁移在非空 SQLite 文件旁创建
  `schema-ui.db.pre-vNNNN-<UTC>.sqlite` 快照，并执行完整性检查；新库不会生成快照。
- 升级前停止 API/Compose 写入并保留数据库副本。若迁移失败，先保留失败数据库，再在
  停机状态将选定的 `pre-vNNNN` 快照复制回 `DB_PATH`，使用已验证的旧二进制/镜像重启，
  不手工编辑 `schema_migrations`。
- Compose 使用命名卷 `db-data`；恢复前先 `docker compose stop api`，把快照导出/复制回
  `/app/data/schema-ui.db`，再启动 API 并检查 `/readyz`。Profile 切换不会删除禁用模块的表或数据。

## 5. 下一步：接业务

> **协议覆盖权威**：本仓对 `schema-ui-docs@v2.7.0` 的整份契约覆盖由 **`I-PROTO-FULL-001`** 定义（[workspace-005 Root attachments](docs/workspaces/workspace-005-full-protocol-contract-v2-7-0/GOAL-001-full-protocol-contract-v2-7-0/attachments/I-PROTO-FULL-001-coverage-v2-7-0.md)：12/12 能力域、24/24 registry type、16/16 conformance 套件 include）。历史 `I-PROTO-001 v0.1.3` 仅为 MVP 回归基线（只读）。协议能力清单见 [protocol-inventory-v2.7.0.md](docs/vision/protocol-inventory-v2.7.0.md)；任何「已支持 v2.7.0」声明必须以覆盖表 + 实现证据背书。

- **完整一方标准 Admin 功能模块**（必须 / 禁止 / 归属判定、组合根与 Profile、全局迁移）：见操作契约  
  **[docs/architecture/module-contribution-playbook.md](docs/architecture/module-contribution-playbook.md)**  
  （架构边界：[docs/architecture/module-architecture.md](docs/architecture/module-architecture.md)；概览入口：[docs/architecture/overview.md](docs/architecture/overview.md)）
- 新增业务页面（**无需修改 Renderer 主路径**；从属于上表 MUST）：
  1. 在对应 owner module 的 `apps/api/internal/modules/<module>/schema/` 添加页面 Schema 文档，并由该模块 Provider 贡献字节（**需重建/重启 API** 后生效）；core 示例位于 `apps/api/internal/modules/schemarender/schema/`；标准 Admin 正例：`apps/api/internal/modules/users/`；
  2. 在模块 Provider 的 Manifest/Navigation contribution 中登记 `pages[]` 与 `navigation`；不要在 `apps/web/public/` 放置生产 Manifest。Manifest 由 API 聚合并经 `/.well-known/schema-ui/app-manifest.json` 发布。
  - 注意：`docs/schemas/` 是上游 **协议 JSON Schema**（node/page/action…），**不是**业务页面文档目录。
- 权限：通过模块的 Authorization/Persistence contribution 声明权限键与 system-data reconcile；全局迁移/快照执行仍由 `apps/api/internal/store` 负责；模块迁移须进入全局台账（见 playbook M5）。
- 参考：[README.md](README.md) 工程化段；`apps/api/README.md` / `apps/web/README.md` 端点与配置表。
