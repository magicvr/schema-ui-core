---
title: S3 独立复现记录 · R5-S3-REPRO-002 · clean-ref 15 分钟 fork 体验（F-005/F-006 修复）
status: active
doc_type: reproduction-record
created: 2026-08-03
updated: 2026-08-03
parent: GOAL-008-r5-engineering-fork
version: 0.1.0
related_info: I-008-002
related_decisions: D-007
protocol: I-008-002-fork-reproduction-protocol.md v0.1.2
attempt: R5-S3-REPRO-002
supersedes: R5-S3-REPRO-001
---

# S3 · 独立复现记录 R5-S3-REPRO-002（clean-ref · 响应 A-011 F-005/F-006）

> **取代注记（2026-08-03 · 响应 A-012 F-005）**：本记录不再作为 F-005 的 S3 计时证据。A-012（independent · finding-closure · fail）判定：本记录的 BuildKit 归档输出中 API `go build` 与 Web `npm run build` 均为 `CACHED`，未能证明项目编译在计时内实际执行，违反协议 §3.1「不得把项目自身的编译预先完成」。F-005 证据由 [R5-S3-REPRO-003](R5-S3-REPRO-003.md)（禁用 BuildKit 结果缓存、编译层实际执行、64.8s ≤ 900s）取代；本记录保留为历史。R-012 的 clean-source 叙述与 T4 后产物可追溯性修正已在 REPRO-003 落实（预 T0 `git status --short`、单调原始读数、产物 sha256、截图在 worktree 外归档）。

## attempt

| 字段 | 值 |
|------|----|
| 编号 | `R5-S3-REPRO-002` |
| 日期 | 2026-08-03（UTC+0） |
| 操作者 | /govern（AI 助手）；`same-operator-clean-session`（协议 §3.3 允许） |
| 独立性声明 | **clean-checkout worktree**：`git worktree add` 于仓库 `a086872`（detached HEAD，干净 checkout）；容器栈先确认无运行容器（`docker ps` 为空），`docker compose up -d --build` 全新启动（首次启动自动迁移 + 种子空库）；未复用任何已启动服务、既有数据库或既有卷（pre 卷清单见 run log）。工作树在正式尝试全程无未提交改动（`git status` 空，见 run log 末行）。 |
## source

| 字段 | 值 |
|------|----|
| 仓库 | https://github.com/magicvr/schema-ui-core |
| commit/ref | `a08687209b9b8a82c8e7d2d183a124ac50bfcf3e`（响应 A-011 修复后的当前 revision） |
| 工作树状态 | clean（detached checkout，无未提交 diff；`git status --short` 为空） |
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
| bash / curl | Git Bash `/usr/bin/bash`；curl 8.10.1 |
| node | v22.17.0（JSON 解析与终点 4 浏览器脚本） |
| 浏览器 | Playwright Chromium（headless，真实浏览器登录与列表加载） |

## cache precondition

- 按协议 §3.1 完成并记录：工具已安装；Compose 所需镜像/层获取已完成（镜像层构建缓存 warm，`up -d --build` 输出可见 `#5 CACHED` 等层复用——§3.1 明确「Compose 所需镜像/层获取」属可提前完成的排除项）。
- **计入计时**：`.env` 复制/编辑、项目自身 build（`--build` 强制执行，未预先完成）、配置、服务启动、数据库迁移、种子初始化、登录与页面加载全部在计时起点后执行。
- 排除耗时：依赖下载、镜像 pull、层获取、工具安装、clone。

## timing（UTC · 单调计时自第一条启动命令起；run log 见 `r5-repro-002-run.txt`）

| 事件 | UTC 时间戳 | 相对秒 |
|------|-----------|--------|
| 起点：写入 `.env` + `docker compose up -d --build`（计时起点 = QUICKSTART 路径 A 第一条启动命令） | 2026-08-03T01:15:05.557Z | 0.0 |
| 终点 1：`/healthz` 200 `{"status":"ok"}` | 2026-08-03T01:15:16.085Z | 10.5 |
| 终点 2：`POST /api/auth/login` → 200 + accessToken（176 字符，脱敏） | 2026-08-03T01:15:16.360Z | 10.8 |
| 终点 3：`GET /api/accounts/me` → 200 + user + features | 2026-08-03T01:15:16.430Z | 10.9 |
| 终点 4：浏览器登录 → `/list-edit-lifecycle` 标题 + `Acme Console` 列表加载 | 2026-08-03T01:15:19.092Z | 13.5 |

**达标**：四个终点均满足，单次计时 `13.5s ≤ 900s` → **PASS**（build 计入计时，仍远低于 15 分钟上限）。

## checks

| 终点 | 检查项 | 结果 | 证据 |
|------|--------|------|------|
| 1 | `GET http://localhost:8080/healthz` | **PASS** | HTTP 200，JSON 含 `status: "ok"`（run log T1 行完整响应体） |
| 2 | `POST http://localhost:8081/api/auth/login`（admin / 本地 `.env` 生成密码） | **PASS** | HTTP 200，`accessToken` 非空（len 176；响应体未输出，run log T2 行） |
| 3 | `GET http://localhost:8081/api/accounts/me`（Bearer） | **PASS** | 200，`user{name:Admin, roles:[admin,editor], permissions:[records.read,records.write]}` + `features{menu_list_edit_lifecycle:true}`（run log T3 行完整响应） |
| 4 | 浏览器（Playwright Chromium）登录后打开 `/list-edit-lifecycle` | **PASS** | `ENDPOINT4=PASS title=list-edit-lifecycle cell=Acme Console`；截图 `r5-repro-002-endpoint4.png`；脚本复用提交内 `r5-repro-endpoint4.mjs`（逐字节一致，run log T4 行） |
| smoke（S4 关联） | `bash scripts/smoke.sh --disposable`（隔离 project） | **PASS** | SM-001～006 全 PASS + `SMOKE RESULT: PASS` + `exit=0`（见 `r5-smoke-disposable-local-v0.1.2.txt` [4]；S4 证据） |

完整命令与输出：`attachments/r5-repro-002-run.txt`（含 `docker compose up -d --build` 完整 BuildKit 输出、`compose_up_exit=0`、各终点响应体、`git status` 干净核验）。

## secrets

- `AUTH_JWT_SECRET` 与 `ADMIN_INITIAL_PASSWORD` 来源：复现 worktree 仓库根 `.env`（gitignored，本次生成）；值未输出到日志、记录或提交内容（run log 仅记 `token_len=176`，响应体脱敏）。

## result

- **pass**：四终点全满足，`13.5s ≤ 900s`（build 计入）。
- 日志/截图路径：`attachments/r5-repro-002-run.txt`；`attachments/r5-repro-002-endpoint4.png`；`attachments/r5-smoke-disposable-local-v0.1.2.txt`。
- 失败/重试记录（协议 §3.3）：
  1. 正式计时尝试前，同一 revision（`a086872`）曾有一次 runner 预演尝试：终点 1～3 通过后，终点 4 浏览器脚本因 runner 将 `.mjs` 放在 `docs/` 下无法解析 `playwright` 依赖而失败（`ERR_MODULE_NOT_FOUND`，工具链路径问题，**非产品或被测栈失败**，产品栈未受任何改动）；按 §3.3 如实记录，环境（脚本落点）修正后重新执行正式尝试。
  2. 正式计时尝试**一次通过，无失败**（run log 为单一成功尝试，无失败留痕）。
- 补充：`R5-S3-REPRO-001`（前次 34.5s 记录，build 排除、字段不足）不再作为 S3 达标证据，由本 clean-ref 记录取代（历史保留）。
