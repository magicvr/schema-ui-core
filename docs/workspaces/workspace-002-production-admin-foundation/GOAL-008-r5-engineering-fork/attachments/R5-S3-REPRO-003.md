---
title: S3 独立复现记录 · R5-S3-REPRO-003 · 无项目编译缓存 15 分钟 fork 体验（响应 A-012 F-005）
status: active
doc_type: reproduction-record
created: 2026-08-03
updated: 2026-08-03
parent: GOAL-008-r5-engineering-fork
version: 0.1.0
related_info: I-008-002
related_decisions: D-008
protocol: I-008-002-fork-reproduction-protocol.md v0.1.2
attempt: R5-S3-REPRO-003
supersedes: R5-S3-REPRO-002
---

# S3 · 独立复现记录 R5-S3-REPRO-003（无 BuildKit 结果缓存 · 响应 A-012 F-005）

## attempt

| 字段 | 值 |
|------|----|
| 编号 | `R5-S3-REPRO-003`（正式计时尝试 = retry #3；前两次失败尝试见 result，均按协议 §3.3 留痕） |
| 日期 | 2026-08-03（UTC+0） |
| 操作者 | /govern（AI 助手）；`same-operator-clean-session`（协议 §3.3 允许） |
| 独立性声明 | **clean-checkout worktree**：`git worktree add --detach` 于仓库 `1961e5a`（detached HEAD，干净 checkout）；**预 T0 的 `git status --short` 输出为空**（run log retry #3 节 L369）；容器栈启动前 `docker ps` 为空；未复用任何已启动服务、既有数据库或既有卷。工作树在正式尝试全程无未提交改动（run log 末尾两次 `git status --short` 均为空）。 |

## source

| 字段 | 值 |
|------|----|
| 仓库 | https://github.com/magicvr/schema-ui-core |
| commit/ref | `1961e5ae7e122d6263846f0ff673fbca074eb03f`（当前 HEAD；相对 A-012 审计时的 `df913a5` 仅追加台账/证据文档，`apps/`、`scripts/`、`compose.yaml`、workflow 均未改动） |
| 工作树状态 | clean（detached checkout，预 T0 与运行后 `git status --short` 均为空） |
| diff | 无 |

## path

| 字段 | 值 |
|------|----|
| 路径 | `compose` |
| API base URL | `http://localhost:8080` |
| Web base URL | `http://localhost:8081`（compose.yaml `web` 映射 `8081:80`，nginx `/api` 反代） |

## platform

| 项 | 值 |
|----|----|
| OS/架构 | Windows 11 Pro `10.0.26200` / x64 |
| Git | 2.47.0.windows.2 |
| Docker / Compose | 29.6.2 / v5.3.1（BuildKit） |
| curl | 8.21.0（Windows） |
| node | v22.17.0（JSON 解析、单调计时与终点 4 浏览器脚本） |
| 浏览器 | Playwright Chromium（headless，真实浏览器登录与列表加载） |

## cache precondition（响应 A-012 F-005：禁用/隔离项目编译缓存）

- **BuildKit 结果缓存已整体禁用（A-012 F-005 核心措施）**：正式计时前执行 `docker rmi schema-ui-core-api:local schema-ui-core-web:local` + `docker builder prune -a -f`（run log 中 BuildKit 输出可核对：`go build` 与 `npm run build` 均以 `RUN` 步骤实际执行并各自 `DONE 12.8s` / `DONE 6.1s`；**正式 retry #3 的项目编译层均非 `CACHED`，且该次仅有一条非编译 `WORKDIR` 缓存**——不再存在 A-012 所指的「项目编译层 `CACHED`」。注：完整归档还包含 retry #2 的同类 `WORKDIR` 缓存，均不涉及编译层；响应 A-013 R-013 表述收窄）。
- **保留的排除项（协议 §3.1 允许）**：工具已安装；基础镜像（`golang:1.26-alpine`/`node:22`/`nginx:1.27-alpine`/`alpine:3.20`）已在本地（镜像拉取/层获取排除）；Go 模块与 npm 依赖经 Dockerfile 内 `--mount=type=cache` 缓存挂载复用（语言依赖缓存排除）。注：因结果缓存整体禁用，`go mod download` 与 `npm ci` 亦在窗口内重新执行（run log 可见 `#23 DONE 11.7s`、`#35 DONE 7.7s`）——这是**保守**口径（依赖步骤本可排除却计入计时），不影响达标结论。
- **计入计时**：`.env` 写入、项目自身编译（`go build` / `npm run build` 在 T0 后**实际执行**，非 CACHED）、配置、服务启动、数据库迁移、种子初始化、登录与页面加载全部在计时起点后执行。

## timing（UTC · 单调原始读数；run log 见 `r5-repro-003-run.txt`）

| 事件 | UTC 时间戳 | 单调原始读数（ns） | 相对秒（单调） |
|------|-----------|--------------------|----------------|
| 起点 T0：写入 `.env` + `docker compose up -d --build` | 2026-08-03T02:09:49.734Z | `403981233142700` | 0.0 |
| 终点 1：`/healthz` 200 `{"status":"ok"}` | 2026-08-03T02:10:51.279Z | `404042779240500` | 61.546 |
| 终点 2：`POST /api/auth/login` → 200 + accessToken（176 字符，脱敏） | 2026-08-03T02:10:51.482Z | `404042981760800` | 61.749 |
| 终点 3：`GET /api/accounts/me` → 200 + user + features | 2026-08-03T02:10:51.627Z | `404043127422900` | 61.894 |
| 终点 4：浏览器登录 → `/list-edit-lifecycle` 标题 + `Acme Console` 列表加载 | 2026-08-03T02:10:54.565Z | `404046066135300` | 64.833 |

**达标**：四个终点均在窗口内满足，单次计时 `64.833s ≤ 900s` → **PASS**（项目编译在计时内实际执行，仍远低于 15 分钟上限）。

## checks

| 终点 | 检查项 | 结果 | 证据 |
|------|--------|------|------|
| 1 | `GET http://localhost:8080/healthz` | **PASS** | HTTP 200，JSON 含 `status: "ok"`（run log T1 行完整响应体） |
| 2 | `POST http://localhost:8081/api/auth/login`（admin / 本地 `.env` 生成密码） | **PASS** | HTTP 200，`accessToken` 非空（len 176；响应体未输出，run log T2 行） |
| 3 | `GET http://localhost:8081/api/accounts/me`（Bearer） | **PASS** | 200，`user{name:Admin, roles:[admin,editor], permissions:[records.read,records.write]}` + `features{menu_list_edit_lifecycle:true}`（run log T3 行完整响应） |
| 4 | 浏览器（Playwright Chromium）登录后打开 `/list-edit-lifecycle` | **PASS** | `ENDPOINT4=PASS title=list-edit-lifecycle cell=Acme Console`；截图 `r5-repro-003-endpoint4.png`（sha256 `89171fb1e43393d7714001dae30aa8732cd5514af85b23c03969d242289809f8`）；脚本复用提交内 `r5-repro-endpoint4.mjs`（逐字节一致） |
| smoke（S4 关联） | `bash scripts/smoke.sh --disposable`（隔离 project） | **PASS** | SM-001～006 全 PASS + `SMOKE RESULT: PASS` + `exit=0`（见 `r5-smoke-disposable-local-v0.1.2.txt`；S4 证据，本轮不改动） |

完整命令与输出：`attachments/r5-repro-003-run.txt`（含 `docker compose up -d --build` 完整 BuildKit 输出——`go build`/`npm run build` 均实际执行非 `CACHED`、`compose_up_exit=0`、各终点响应体、预 T0 与运行后 `git status --short` 干净核验、teardown `down -v`）。

## secrets

- `AUTH_JWT_SECRET` 与 `ADMIN_INITIAL_PASSWORD` 来源：复现 worktree 仓库根 `.env`（gitignored，本次生成）；值未输出到日志、记录或提交内容（run log 仅记 `token_len=176`，响应体脱敏）。

## result

- **pass**：四终点全满足，`64.833s ≤ 900s`，且项目编译（`go build` 12.8s、`npm run build` 6.1s）在计时窗口内实际执行（BuildKit 输出非 `CACHED`）。
- 日志/截图路径：`attachments/r5-repro-003-run.txt`；`attachments/r5-repro-003-endpoint4.png`。
- 失败/重试记录（协议 §3.3，run log 内按 §4 字段完整留痕）：
  1. **attempt #1（runner 工具链故障）**：驱动脚本在 `docker compose up -d --build` 处因 PowerShell NativeCommandError（stderr 管道）中断，BuildKit 输出未被归档、终点未测量；compose 实际已完成且容器 healthy。环境未因修复而改动——以 `down -v` + `docker rmi` + `docker builder prune -a` 完整复位后重试。
  2. **attempt #2（runner 工具链故障）**：T1/T4 在窗口内 PASS（浏览器真实登录 + 列表加载），但 T2/T3 因 PowerShell 5.1 向 `curl.exe` 传内联 `-d` JSON 时引号被吞、登录体未送达而失败（`login_http_ok=0`、`/me` UNAUTHENTICATED）——**非产品/被测栈故障**；随后对同一已构建栈以 `--data @file` 修正调用复验：`token_len=176`、`/me` 返回完整 user+features。按 §3.3「登录失败应立即使该尝试失败」，该尝试记为失败，并复位环境后重试。
  3. **attempt #3（正式尝试）**：单次通过，无失败（本记录采用）。
- 补充：`R5-S3-REPRO-002`（13.5s，BuildKit 项目编译层 `CACHED`）不再作为 F-005 的 S3 计时证据，由本 clean-checkout、无编译缓存记录取代（历史保留）。
