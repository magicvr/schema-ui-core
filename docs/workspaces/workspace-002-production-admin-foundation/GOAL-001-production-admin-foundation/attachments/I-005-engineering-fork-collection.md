---
title: I-005 · 工程化部署基线、15 分钟 fork 计时口径与可复现实验环境信息收集
status: active
doc_type: info-collection
created: 2026-08-02
updated: 2026-08-02
parent: GOAL-001-production-admin-foundation
version: 0.2.2
related_info: I-005
related_info_b: I-006
related_decisions: D-012, D-013
---

# I-005 · R5 工程化 / fork 信息收集（含 I-006 复核）

> **性质**：回答「目标部署基线、15 分钟 fork 计时口径与可复现实验环境是什么」所需的当前工程面事实、候选方案、验收口径与独立复现方法；同时复核 `I-006`（是否在本波次纳入最小操作日志）。
> **裁决状态**：用户于 2026-08-02 按 P-004 书面裁决附件 §6 四类边界——**部署基线 A**（文档双进程 + Docker Compose）、**建议计时口径**（终点=登录+后台首页可交互、不含依赖下载、≥1 次独立复现）、**建议复现方法**（文档步骤 + smoke 清单 + 独立复现记录，R5 落地 `scripts/smoke.sh`）、**I-006 方案甲**（操作日志为 R5 可选加分 checkpoint）。Root **D-013** 已将 `I-005` 置为 `verified`、`I-006` 置为 `closed`。**F-001 澄清（v0.2.1 · 响应 GOAL-008 A-001）**：候选 A 的「可选」= **fork 使用者可选择本地双进程或 Compose 两条启动路径**；**Compose 本身是 R5 必须交付和验收的第二启动路径，不是可选加分项**（详见 GOAL-008 `01-decision` D-002）。
> **不是**：R5 详细实施方案冻结、R5 子目标立项、Docker 或部署实现、15 分钟计时实测、产品代码变更。本轮**未修改任何产品代码**。
> **扫描日期**：2026-08-02。工作区 `shared_materials_catalog: none`，全部事实来自本仓库代码、README、CI 配置与 `git` 探测；本轮**未运行**应用/测试（仅静态核对）。

## 0. 总览结论

| 维度 | 当前事实 | 对 R5 的含义 |
|------|----------|---------------|
| 运行命令 | API `make run` / `go run ./cmd/server`（:8080）；Web `npm run dev`（:5173，`WEB_PORT` 可覆盖） | 本地双进程可启动；但无统一一键启动，fork 需两条命令与两个窗口 |
| 配置面 | API 以 env 配置（`.env.example` 齐全），生产缺 `AUTH_JWT_SECRET` / `ADMIN_INITIAL_PASSWORD` fail-closed | dev/prod 区分已有雏形；生产部署所需环境变量无端到端文档 |
| 容器 / 部署 | **全仓无 Dockerfile / docker-compose / 容器路径**；无生产静态托管与生产 `/api` 反代配置 | VP-002 #7 的「Docker 一键启动」与 #6 的 fork 体验尚未具备，是 R5 核心差量 |
| 健康检查 | `GET /healthz` 公开返回 `{status,timestamp,version,commit}`；API 有优雅停机 | 容器/编排探针可用；Web 侧无等价启动验证 |
| CI | `r6-basic-matrix.yml`：web（npm ci/test/build）、api（go test/build）、browser-e2e（Playwright chromium） | 有测试/构建/E2E，**无部署/发布步骤** |
| 仓库卫生 | `apps/api/data/` 运行时 DB 与迁移快照已 gitignore（`/data/`） | 运行时产物不落库；R5 复现实验应保持同一边界 |
| fork 文档 | 根 `README.md` 仅骨架级「本地运行」两段；无环境变量清单、无 15 分钟口径、无可复现实验清单 | 「按文档配置环境变量并启动 ≤15 分钟」尚无口径与证据方法 |
| 操作日志（I-006） | 代码库无 operation/audit log 能力；仅 `LOG_LEVEL` 服务日志 | VP-002 将最小操作日志列为加分项、非硬关门条件；纳入与否需裁决 |

## 1. 当前工程 / 部署面事实

### 1.1 API（Go 1.26+ · `apps/api`）

| 事实 | 证据 |
|------|------|
| 入口 `cmd/server`；`make run` / `go run ./cmd/server`；`make build` 产 `bin/schema-ui-core-api`（version/commit/builtAt ldflags） | `apps/api/Makefile` |
| env 配置键：`APP_NAME`、`APP_ENV`（development/production）、`HTTP_ADDR`（:8080）、读/写/空闲超时、`LOG_LEVEL`、`AUTH_JWT_SECRET`、`AUTH_ACCESS_TTL`（15m）、`AUTH_REFRESH_TTL`（720h）、`DB_PATH`（./data/schema-ui.db）、`ADMIN_INITIAL_PASSWORD`、`AUTH_DEV_SESSION_ENABLED`（默认 false） | `apps/api/.env.example`；`apps/api/internal/config/config.go` |
| 非 development 缺 `AUTH_JWT_SECRET` → 启动失败；缺 `ADMIN_INITIAL_PASSWORD` → 启动失败（fail-closed） | `apps/api/cmd/server/main.go` `resolveJWTSecret` / `resolveSeedHash` |
| 健康探活 `GET /healthz` 公开，返回 `{status:"ok",timestamp,version,commit}` | `apps/api/internal/handler/health.go` |
| SQLite 存储（DB_PATH），迁移 runner + 迁移前快照，启动种子 admin/seedRecords（空表才 seed）；SIGINT/SIGTERM 优雅停机 | `apps/api/cmd/server/main.go`；`apps/api/internal/store/` |
| 认证：短 JWT access + opaque refresh（SHA-256 哈希存 SQLite）+ bcrypt；记录写门禁 `records.read` / `records.write`（401/403） | `apps/api/internal/auth/`；`handler/records.go` |
| **无 Dockerfile / 容器 / 生产托管配置** | 仓库全局 `find`/`grep` 无 docker、nginx、caddy 引用（排除 node_modules/dist） |

### 1.2 Web（React 19 + Vite 6 + TS + Tailwind 4 · `apps/web`）

| 事实 | 证据 |
|------|------|
| 要求 Node 20+；`npm install`；`npm run dev`（默认 :5173，`WEB_PORT` 可覆盖，`strictPort`）；`npm run build` = `tsc -b && vite build`；`npm test`；`npm run test:e2e`（Playwright Chromium） | `apps/web/package.json`；`apps/web/README.md`；`apps/web/vite.config.ts` |
| Vite dev 代理 `/api` → `http://127.0.0.1:8080`（**仅 dev**） | `apps/web/vite.config.ts` |
| shell 拉取 app manifest、加载 `/api/accounts/me` 会话快照、Schema 驱动渲染代表页（R1/R4） | `apps/web/src/app/`、`src/renderer/` |
| **无生产静态托管配置**（nginx/Caddy/托管平台）；生产 `/api` 同源反代未文档化 | 仓库探测 |

### 1.3 CI / 仓库卫生

| 事实 | 证据 |
|------|------|
| `.github/workflows/r6-basic-matrix.yml`：web（npm ci / test / build）、api（go test / go build）、browser-e2e（npm ci + playwright chromium + test:e2e） | `.github/workflows/r6-basic-matrix.yml` |
| 无部署/发布 job；无容器构建 job | 同文件 |
| `apps/api/.gitignore` 忽略 `/data/`（运行时 DB 与迁移快照）、`bin/`、`.env`（保留 `.env.example`） | `apps/api/.gitignore` |
| 根 README 有「本地运行（骨架阶段）」API/Web 两段与探活示例 | `README.md` |

### 1.4 现有文档缺口（fork 视角）

- 无统一的「clone → 配置 → 启动 → 登录」快速上手文档；README 仅两个 app 的骨架命令。
- 无环境变量生产配置清单（哪些必填、fail-closed、TTL、种子密码策略）。
- 无容器/部署边界说明（dev/prod、同源反代、DB volume）。
- 无「15 分钟 fork」计时口径或独立复现实验方法。
- 无操作日志/审计日志现状记录（见 §5）。

## 2. 候选部署基线（D-013 前历史候选 · 已裁决方案 A）

| 候选 | 内容 | 对 VP-002 #6/#7 的匹配 | 成本 |
|------|------|------------------------|------|
| **A · 双进程文档 + Docker Compose 一键启动（已裁决；D-013 前称「可选」）** | 本地双进程为默认（现状命令 + 文档补齐）；另交付 `compose.yaml`：api 多阶段构建（Go → 精简镜像，暴露 :8080，DB volume）+ web 静态构建（nginx/Caddy 服务 `dist/` 并把 `/api` 反代到 api service），单源同源、免 CORS | 直接满足 #7「Docker 一键启动」与 #6「按文档配置并启动」；健康探针复用 `/healthz` | 中：需 compose + 镜像 + 文档；web 反代与 api 探针有现成基础 |
| **B · 仅文档化本地双进程（无容器）** | R5 只补齐文档、健康检查说明与复现清单；Docker 记为 R5 非目标或后续加分项 | 部分满足 #6；不满足 #7 的 Docker 一键启动（VP-002 明列为可交付项） | 低 |
| **C · API 直接托管 web 静态产物（单进程）** | 扩展 api 提供静态文件 + SPA fallback + `/api`（同源单端口） | 满足 #6/#7 的「单进程」形态，省反代 | 中高：需改 api 静态服务与 SPA fallback，改动产品主路径 |

> 建议：**A**。B 会把 VP-002 #7 的 Docker 一键启动直接留在 R5 未交付；C 需要改动 api 主路径（与「Schema Renderer 主路径不改写」目标交叉）。A 保留现有本地路径、新增容器路径，最贴近 VP-002 成功边界。**F-001 澄清（2026-08-02）**：本候选「可选」指 **fork 使用者可选择本地双进程或 Compose 两条启动路径**；Compose 本身是 R5 **必须交付和验收**的第二启动路径，不是「可选加分项」（详见 GOAL-008 `01-decision` D-002 / 03-audit A-001 响应）。最终形态与镜像方案已由 D-013（部署基线 A）决定；精确镜像 / Compose 契约由 `I-008-001` 在 GOAL-008 冻结。

## 3. 15 分钟 fork 计时口径（D-013 前历史候选 · 已采纳建议口径）

VP-002 #6 的体验指标为「按文档完成环境配置并成功登录进入系统，目标时间 ≤15 分钟」。建议口径（候选，未冻结）：

| 项 | 建议口径 | 待确认点 |
|----|----------|----------|
| 起点 | 干净开发机（Go 1.26+、Node 20+ 已安装），`git clone` 完成后、文档第一条命令执行前 | 是否含 `git clone` / 依赖下载时间？（受网络波动影响；建议**不含**首次依赖下载，以「依赖就绪后按文档命令到登录成功」计时，并在复现记录中注明依赖预装） |
| 终点 | ① `GET /healthz` 返回 `ok`；② 浏览器登录种子 admin（`admin` / `ADMIN_INITIAL_PASSWORD`）成功；③ 看到 schema 驱动的可用后台页面（菜单投影 + 至少一个可交互代表页） | 是否以「看到后台首页」为终点，还是「完成一次记录列表/详情」？建议终点 = 登录成功并看到后台首页（可交互列表加载）；不把 CRUD 完成计入「启动体验」 |
| 测量方法 | 由文档读者独立执行一次完整启动并记录：日期、仓库 ref（commit）、`git`/依赖版本、耗时（分钟）、步骤完成时刻；可与 CI smoke 对照 | 是否要求 ≥1 次独立复现才算达标 |
| 判定 | 耗时 ≤15 分钟且三步终点全部满足 → 达标 | — |

> 建议：终点 = 登录成功 + 后台首页可交互（列表加载）；不含依赖下载；记录一次完整复现。是否把「记录一次 CRUD」纳入终点由用户裁决（建议不纳入，区分「启动体验」与「能力验证」）。

## 4. 可复现实验方法（D-013 前历史候选 · 已采纳建议方法）

R5 的「可重复验收」建议由「文档步骤 + smoke 清单」构成：

1. **文档步骤（README/QUICKSTART）**：clone → `cp apps/api/.env.example .env` 并按需设生产键 → 启动 api（`make run`）→ 启动 web（`npm run dev`）→ 登录 → 说明数据初始化（DB_PATH、种子）。
2. **smoke 清单（机器可核对）**：
   - `curl -f http://localhost:8080/healthz` → `{"status":"ok",...}`；
   - 登录：`POST /api/auth/login`（admin + 初始密码）→ 200 + access/refresh；
   - 会话：`GET /api/accounts/me`（Bearer）→ `{user, features}`；
   - 代表页：浏览器访问 `data-table` / `list-edit-lifecycle`，列表加载且写入口按角色呈现；
   - 数据初始化可重复：清空 DB_PATH 后重启，种子行为一致（已有 seed 幂等测试基础）。
3. **（若候选 A）容器 smoke**：`docker compose up` 后同一组 `/healthz` + 登录 + 代表页检查；DB volume 重启保持。
4. **验收证据**：一次独立复现记录（日期、ref、版本、耗时）+ smoke 输出；建议后续落地为一个可运行脚本（如 `scripts/smoke.sh`）与 R5 子目标实施项。

> 说明：本轮只固定**方法口径**；smoke 脚本与容器实现属 R5 子目标实施，不在本收集内实现。

## 5. I-006 复核 · 是否在本波次纳入最小操作日志

- **现状**：`apps/api` 与 `apps/web` 无 operation/audit log 能力（`grep` 无命中）；服务端仅 `slog` 结构化日志（`LOG_LEVEL`），不含「谁在何时对哪个 record 做了什么写操作」的可审计记录。
- **VP-002 边界**：阶段 3 验收明确「最小操作日志（记录关键写操作）作为加分项，不作为 VP 硬性关门条件」；Charter 非目标未排除。
- **价值**：对生产 Admin 基架，记录 records create/update/delete 与认证关键事件（登录/登出/刷新/失败）可支撑审计与排障；fork 用户开箱即有可追溯写路径。
- **成本**：需要定义事件 schema、SQLite 存储（表 + repository + 迁移）、写路径接线（records handler + auth 事件）、查询/清理边界与测试；属独立工作块，与容器/fork 正交。
- **建议**：
  - **方案甲（推荐）**：纳入 R5 为**可选加分检查点**——若实施，记入 SQLite（如 `operation_log` 迁移 + repository），覆盖 records 写 + auth 关键事件；**不阻断** R5 核心工程化/容器/fork 验收；若工时紧张可降级为「记录决策 + 留作后续」。
  - **方案乙**：本波次不纳入；在 R5 决策中明确「最小操作日志为非目标/后续加分项」，维持 VP-002 加分项定位。
- **裁决点**：甲（纳入，作为可选加分 checkpoint）还是乙（本波次排除）。升级为 required 需另记决策并更新 VP/Root 信息表。**（已裁决 · D-013：方案甲——纳入 R5 为可选加分 checkpoint，S6 不进成功标准分母；是否实施由用户另行决定。）**

## 6. 信息裁决点（P-004 · 已于 D-013 裁决）

> **已裁决（2026-08-02 · Root D-013）**：① 部署基线 = **A**；② 计时口径 = **建议口径**；③ 复现方法 = **建议方法**（含 R5 落地 `scripts/smoke.sh`）；④ I-006 = **方案甲**。`I-005` → `verified`，`I-006` → `closed`。

下列内容为 **D-013 裁决前**的候选与建议口径（保留供 R5 子目标方案冻结作为输入参考，不再是开放裁决点）：

1. **部署基线候选**：A（双进程文档 + Compose，已裁决）/ B（仅文档双进程，未选）/ C（API 托管静态单进程，未选）。
2. **计时口径候选**：终点 = 登录 + 后台首页可交互？是否含依赖下载？（已裁决：不含依赖下载；终点 = 登录 + 后台首页可交互；要求 ≥1 次独立复现并记录。）
3. **复现方法候选**：「文档步骤 + smoke 清单 + 独立复现记录」口径（已裁决采纳；`scripts/smoke.sh` 在 R5 子目标实施）。
4. **I-006 候选**：方案甲 / 方案乙（已裁决：方案甲——R5 可选加分 checkpoint，S6 不进成功标准分母）。

R5 子目标（GOAL-008）据此承接实施；`I-008-001/002/003` 作为实施前 required 契约继续 `open`（见 GOAL-008 `00-meta` 信息需求表）。R1～R4 结论与本收集互不推翻。

## 7. 证据索引

- `apps/api/Makefile`、`apps/api/.env.example`、`apps/api/README.md`
- `apps/api/internal/config/config.go`、`apps/api/cmd/server/main.go`
- `apps/api/internal/handler/health.go`
- `apps/api/.gitignore`（`/data/` 忽略运行时 DB）
- `apps/web/package.json`、`apps/web/README.md`、`apps/web/vite.config.ts`
- `.github/workflows/r6-basic-matrix.yml`
- `README.md`（根 · 本地运行骨架）
- `docs/vision/plans/VP-002-production-admin-foundation.md`（#6/#7 验收方向）
- 仓库全局探测：无 Dockerfile / docker-compose / nginx / caddy 引用；无 operation/audit log 命中
