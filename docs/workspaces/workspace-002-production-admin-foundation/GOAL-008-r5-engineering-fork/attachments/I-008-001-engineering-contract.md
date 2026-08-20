---
title: I-008-001 · 环境/配置与容器部署契约
status: active
doc_type: contract
created: 2026-08-02
updated: 2026-08-03
parent: GOAL-008-r5-engineering-fork
version: 1.0.1
related_info: I-008-001
related_decisions: D-003, D-011
---

# I-008-001 · 环境 / 配置与容器部署契约（冻结）

> **性质**：回答「环境/配置与容器部署的精确契约是什么」——把 Root D-013 部署基线 A 与 A-001 R-002 最低清单固化为可实施、可验收的版本化契约。冻结后 `I-008-001` 由 GOAL-008 **D-003** 置为 `verified`，解除 S1（环境/配置基线）与 S2（容器一键启动）的方案/立项目门禁。
> **不是**：S1/S2 的实现成品（Dockerfile / compose.yaml / README 正文 / smoke.sh 属实施），也不是完整生产运维 / CI-CD 部署流水线（仍为非目标）。
> **依据**：Root D-012 / D-013；A-001 R-002 最低清单；本仓库静态核对（`apps/api/internal/config/config.go`、`.env.example`、`handler/health.go`、`apps/web/vite.config.ts`、`apps/web/README.md`、`.github/workflows/r6-basic-matrix.yml`）。

## 1. 环境变量契约（API · dev/prod 行为）

| 键 | 默认（dev） | 生产行为 | 说明 |
|----|------------|----------|------|
| `APP_NAME` | `schema-ui-core-api` | — | 应用名 |
| `APP_ENV` | `development` | `production` | 驱动 dev/prod 密钥与种子分支 |
| `HTTP_ADDR` | `:25080` | `:25080`（容器内固定；W7 F-008 起 compose 不发布宿主端口） | API 监听地址 |
| `HTTP_READ_TIMEOUT` | `5s` | 保持 | 读超时 |
| `HTTP_WRITE_TIMEOUT` | `10s` | 保持 | 写超时 |
| `HTTP_IDLE_TIMEOUT` | `60s` | 保持 | 空闲超时 |
| `LOG_LEVEL` | `info` | `info`（可 `debug`/`warn`/`error`） | slog 级别 |
| `AUTH_JWT_SECRET` | dev 内建密钥 + 警告 | **必填，缺失 fail-closed** | access 签发密钥；生产缺失 → 启动失败 |
| `AUTH_ACCESS_TTL` | `15m` | 保持 | access 时效 |
| `AUTH_REFRESH_TTL` | `720h`（30d） | 保持 | refresh 时效 |
| `DB_PATH` | `./data/schema-ui.db` | `/app/data/schema-ui.db`（volume 挂载点） | SQLite 路径 |
| `ADMIN_INITIAL_PASSWORD` | dev 兜底 `admin` | **必填，缺失 fail-closed** | 首次种子 admin 密码；生产缺失 → 启动失败 |
| `AUTH_DEV_SESSION_ENABLED` | `false` | `false`（生产禁止） | 静态开发会话兜底，显式 opt-in |

**Web**：仅 dev 用 `WEB_PORT`（默认 25173，`strictPort`）；SPA 使用**相对路径** `/api/*`（`auth-client.ts` 硬编码，无 `VITE_` API base），因此生产**必须**同源反代（见 §4），无需 CORS。

## 2. 健康检查 / 启动验证契约

| 项 | 契约 |
|----|------|
| API liveness | `GET /healthz` → `200` + `{"status":"ok",...}`；进程存活探针，**不访问数据库** |
| API readiness | `GET /readyz` → `200` + `{"status":"ok",...}`；在 liveness 之上执行轻量 SQLite `SELECT 1`，数据库不可读时 `503 {"status":"unavailable",...}`（A-002 F-002-006 / GOAL-009 S5） |
| Web readiness | 静态服务返回 `index.html`（SPA fallback）；`/api` 经反代可达 `api` 的 `/readyz` |
| Compose healthcheck | `api`：`wget -qO- http://127.0.0.1:25080/readyz`（或容器内 `curl`）作为 `service_healthy` 就绪判据；`web`：`wget -qO- http://127.0.0.1/` 200 |
| 启动验证口径 | 本地双进程与 `docker compose up` 两条路径，均以「`/healthz` ok + 登录种子 admin + 后台首页可交互」为终态（S3 计时终点，见 `I-008-002`；smoke 判据沿用 `/healthz` liveness） |

## 3. 容器 / Compose 契约（部署基线 A）

| 项 | 契约 |
|----|------|
| 清单位置 | 仓库根 `compose.yaml`；服务名 `api`、`web` |
| `api` 服务 | 构建 `apps/api/Dockerfile`（**多阶段**：`golang:1.26` 构建 → 精简运行镜像）；暴露 `25080`；`DB_PATH` 指向 `/app/data/schema-ui.db`；secret 经 `.env` / compose env（`AUTH_JWT_SECRET`、`ADMIN_INITIAL_PASSWORD` 生产必填）；healthcheck `/readyz`；`restart: on-failure` |
| `web` 服务 | 构建 `apps/web/Dockerfile`（**多阶段**：`node:22` `npm ci` + `npm run build` → `nginx:alpine` 服务 `dist/`）；暴露 `80`；healthcheck 静态页 200；`depends_on: api: condition: service_healthy` |
| DB volume | 命名卷（如 `db-data`）挂载到 `api` 的 `/app/data`；`docker compose down`/重启后数据保持 |
| 探针 | `/healthz`（api liveness）、`/readyz`（api readiness，Compose `service_healthy`）与静态页 200（web），见 §2 |

## 4. Web SPA fallback 与 `/api` 反代（nginx）

- `location /` → `try_files $uri $uri/ /index.html;`（SPA fallback，直接刷新路由可回退）。
- `location /api` → `proxy_pass http://api:25080;`（含 `/api` 前缀透传；API 路由均在 `/api` 下）。
- 单源同源部署：SPA 相对路径 `/api/*` 与反代天然匹配，**无需 CORS 配置**。
- 静态资源：`dist/` 产物与 `public/.well-known/schema-ui/app-manifest.json` 由 nginx 直接服务。

## 5. 服务依赖 / 超时 / 失败行为

| 行为 | 契约 |
|------|------|
| 启动依赖 | `web` 等待 `api` healthy 后再启动（compose `depends_on` + healthcheck）；本地双进程以文档顺序启动（api 先、web 后） |
| fail-closed | 生产缺 `AUTH_JWT_SECRET` / `ADMIN_INITIAL_PASSWORD` → api 拒绝启动；`AUTH_DEV_SESSION_ENABLED` 生产禁止启用 |
| 优雅停机 | api 捕获 SIGINT/SIGTERM，10s 宽限内 `srv.Shutdown` |
| 失败行为 | api 不可达时 web 反代返回 502；DB volume 独立于容器生命周期，重启不丢数据；迁移快照机制在容器内同样生效（`/app/data` 预迁移快照） |
| 并发/单写者 | 沿用现有 SQLite 单写者 + last-write-wins；容器单实例 |

## 6. CI 入口

- `.github/workflows/r6-basic-matrix.yml`（或后续新增 job）增加**容器/smoke 入口**：构建 `api`/`web` 镜像 + `docker compose up` + 执行 smoke（`/healthz` → 登录 → `/me` → 代表页）作为回归；精确 job 与 smoke 判据随 S2/S4 落地（`I-008-002` 冻结 smoke 判据）。
- 现状 CI 另有 web（npm ci/test/build）、api（go test/build）、browser-e2e（Playwright）；无完整部署/CD 流水线（保持非目标）。

## 7. 验收清单（S1/S2 可核对）

| # | 检查项 | 对应契约 |
|---|--------|----------|
| C-001 | `.env.example` 与 config.go 键一致；dev/prod 行为注释齐全 | §1 |
| C-002 | `GET /healthz`（liveness）200 + `{"status":"ok"}`（本地与容器内均可） | §2 |
| C-003 | `docker compose up` 后 api healthy、web 200，登录种子 admin 成功 | §2/§3 |
| C-004 | `apps/api/Dockerfile`、`apps/web/Dockerfile`、根 `compose.yaml` 存在且与 §3 一致 | §3 |
| C-005 | nginx 配置含 SPA fallback + `/api` 反代；`docker compose up` 后直接刷新 `/list-edit-lifecycle` 可回退 | §4 |
| C-006 | 重启后（`down`/`up` 或 `restart`）DB 数据保持；删除记录不复活 | §3/§5 |
| C-007 | CI 增加容器/smoke 入口 job（本地 smoke 通过） | §6 |

## 8. 边界（非目标）

- 完整生产运维 / CI-CD 部署流水线、TLS 终止、多实例水平扩展、对象存储、监控告警：**非目标**（维持 Root D-013 边界）。
- 精确 nginx.conf 内容、镜像 tag、资源限制：属 S2 实施细节，由实施留痕，不在此枚举。
- 15 分钟计时与 smoke.sh 退出码判据：由 `I-008-002` 冻结，不在本契约重复。

## 8a. 修订记录

| 版本 | 日期 | 变更 |
|------|------|------|
| 1.0.0 | 2026-08-02 | 冻结（GOAL-008 D-003；S1/S2 方案门禁解除；C-001～C-007 随 S1/S2 验收） |
| 1.0.1 | 2026-08-03 | 响应 GOAL-009 A-003 R-002：探针语义同步 A-002 F-002-006 / GOAL-009 S5——§2/§3 明确 `/healthz` = liveness（不访问 DB）、`/readyz` = readiness（liveness + SQLite `SELECT 1`，故障 503），Compose `service_healthy` 用 `/readyz`。C-001～C-007 为 S1/S2 历史验收事实（当时仅 `/healthz` liveness），不因本次语义同步重写；`/readyz` 由 GOAL-009 S5 引入并经 `health_test.go` 正常/故障注入覆盖。`I-008-001` 维持 `verified`（语义同步不改变信息结论）。 |

## 9. 证据索引

- `apps/api/internal/config/config.go`、`apps/api/.env.example`、`apps/api/internal/handler/health.go`
- `apps/web/vite.config.ts`、`apps/web/README.md`、`apps/web/src/account/auth-client.ts`（相对 `/api` 路径）
- `.github/workflows/r6-basic-matrix.yml`
- Root D-012 / D-013、[I-005-engineering-fork-collection.md](../../GOAL-001-production-admin-foundation/attachments/I-005-engineering-fork-collection.md) v0.2.2
- A-001 R-002 最低清单（GOAL-008 `00-meta` 信息表）
| 1.0.2 | 2026-08-20 | 对齐现行端口与 smoke 入口（D-011）：§1 `HTTP_ADDR` 默认 `:8080`→`:25080`、Web `WEB_PORT` 默认 `5173`→`25173`；§2 Compose healthcheck `127.0.0.1:8080`→`:25080`；§3 `api` 暴露 `8080`→`25080`（compose 不发布宿主端口）；§4 `proxy_pass http://api:8080`→`:25080`；§6 CI `container-smoke` 改为调用 `scripts/pre-release-smoke.sh`（隔离 + CSP + W16-F01 + C-006）。C-001～C-007 为 S1/S2 历史验收事实，不因本次端口/CI 文字同步重写。
