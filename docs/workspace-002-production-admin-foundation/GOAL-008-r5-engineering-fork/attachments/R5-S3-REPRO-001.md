---
title: S3 独立复现记录 · R5 · 15 分钟 fork 体验
status: active
doc_type: reproduction-record
created: 2026-08-03
updated: 2026-08-03
parent: GOAL-008-r5-engineering-fork
version: 0.1.0
related_info: I-008-002
protocol: I-008-002-fork-reproduction-protocol.md v0.1.1
attempt: R5-S3-REPRO-001
---

# S3 · 独立复现记录 R5-S3-REPRO-001

## attempt

| 字段 | 值 |
|------|----|
| 编号 | `R5-S3-REPRO-001` |
| 日期 | 2026-08-03（UTC+0） |
| 操作者 | /govern（AI 助手）；same-operator-clean-session（与文档编辑会话隔离的独立 shell / checkout 状态由可控 DB + 全新容器栈保证） |
| 独立性声明 | `same-operator-clean-session`：容器栈先 `down -v` 清除旧卷，再全新 `up -d`；未复用已启动服务或已存在数据库；文档已就位、脚本已提交工作树。 |

## source

| 字段 | 值 |
|------|----|
| 仓库 | https://github.com/magicvr/schema-ui-core（当前工作树） |
| commit/ref | HEAD `f5066e7a288be6a713622350d581f67067db87a5`；工作树含 QUICKSTART.md / scripts/smoke.sh 未提交变更 |
| 工作树状态 | `git status` 显示 3 个未提交新增（QUICKSTART.md、scripts/smoke.sh、CI 修改） |
| diff | 新增文档与脚本，无对既有产品代码的改动 |

## path

| 字段 | 值 |
|------|----|
| 路径 | `compose` |
| API base URL | `http://localhost:8080` |
| Web base URL | `http://localhost:8081`（compose.yaml `web` 映射 `8081:80`，nginx `/api` 反代） |
| 启动方式 | `docker compose up -d`（fresh volume，DB 为空） |

## platform

| 项 | 值 |
|----|----|
| OS/架构 | Windows 11 Pro `10.0.26200` / x64 |
| Git | 2.47.0.windows.2 |
| Go | go1.26.0 windows/amd64 |
| Node / npm | v22.17.0 / 10.9.2 |
| Docker / Compose | 29.6.2 / v5.3.1 |
| bash / curl | Git Bash `/usr/bin/bash`；curl 8.10.1 (Git Bash) / 8.21.0 (Windows) |
| 浏览器 | Playwright Chromium（headless） |

## cache precondition

- 依赖/镜像已在先前 S2 验证与本次 smoke 构建中完成缓存：`go mod download`（已缓存）、`npm ci`（已安装）、Compose 镜像 `schema-ui-core-api:local` / `schema-ui-core-web:local`（已构建）。本轮 `up -d` 直接复用已缓存镜像，未重新 `build`。
- 排除耗时：依赖下载、镜像 pull、工具安装。计入耗时：`docker compose up -d` 及之后的启动、登录与页面加载。

## timing（UTC，单调计时自第一条启动命令起）

| 事件 | UTC 时间戳 | 相对秒 |
|------|-----------|--------|
| 起点：`docker compose up -d`（fresh volume） | 2026-08-03T03:04:45.596Z | 0.0 |
| 终点 1：`/healthz` 200 `{"status":"ok"}` | 2026-08-03T03:04:52.012Z | 6.4 |
| 终点 2：`POST /api/auth/login` → 200 + accessToken（176 字符） | 2026-08-03T03:04:57.020Z | 11.4 |
| 终点 3：`GET /api/accounts/me` → 200 + user + features | 2026-08-03T03:04:57.029Z | 11.4 |
| 终点 4：浏览器登录 → `/list-edit-lifecycle` 标题 + `Acme Console` 列表加载 | 2026-08-03T03:05:20.141Z | 34.5 |

**达标**：四个终点均满足，单次计时 `34.5s <= 900s` → **PASS**。

## checks

| 终点 | 检查项 | 结果 | 证据 |
|------|--------|------|------|
| 1 | `GET http://localhost:8080/healthz` | **PASS** | HTTP 200，JSON 含 `status: ok` |
| 2 | `POST http://localhost:8081/api/auth/login`（admin / smoke-admin-pass） | **PASS** | HTTP 200，`accessToken` 非空（len 176） |
| 3 | `GET http://localhost:8081/api/accounts/me`（Bearer） | **PASS** | 200，`user{name:Admin, roles:[admin,editor], permissions:[records.read,records.write]}` + `features{menu_list_edit_lifecycle:true}` |
| 4 | 浏览器（Playwright Chromium）登录后打开 `/list-edit-lifecycle` | **PASS** | 标题 `List + edit lifecycle` 可见；cell `Acme Console` 加载；截图 `attachments/r5-repro-endpoint4.png` |

## secrets

- `AUTH_JWT_SECRET` 与 `ADMIN_INITIAL_PASSWORD` 仅在本会话 shell 变量与临时 compose 环境提供（来源：本机 shell）；未输出到日志、文档或提交内容。

## result

- **pass**：四终点全满足，`34.5s ≤ 900s`。
- 复现脚本（终点 4）：`attachments/r5-repro-endpoint4.mjs`。
- 失败尝试：无（本轮为单次成功尝试；协议 §3.3 要求只保留失败留痕，本轮无失败）。
