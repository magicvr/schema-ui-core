---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-004-r3-compose-cicd
version: 0.1.0
---

# E-002 · S1–S3 实跑证据（2026-08-29）

## 交付物

- golden-field `consumer-regression.yml` 重构（commit `c4d14ea`）：免凭据（npmjs 公开 · 项目 .npmrc 钉 npmjs）· setup-node pnpm cache · Go `go get @latest`+`go build` · serve 后台 + healthz/readyz 轮询 · 四探针 · **SIGTERM 收尾断言 `shutdown.complete`**（workflow 即 harness A 主机形态）。
- 主仓 compose 实跑（`docker compose up -d --build` 全服务）——**首次本环境实跑**（核销 GOAL-005 F-001）。

## 证据

| 判据 | 证据 | 结果 |
|------|------|------|
| C1 compose 全服务 | `up --build` exit 0 · api `healthy`（readyz `{"status":"ok"}` via exec）· web `healthy` + 反代 :25081 → 200 · `compose stop`（grace 15s）→ api **ExitCode 0** · web ExitCode 0；日志 `shutdown.starting`（terminated）→ `shutdown.complete` | ✅ |
| C2 workflow 重构 + 本地等价 | workflow 文件重构提交；本地等价（go build → serve → healthz 200 · readyz 200 → 四探针全绿：protocol 2.9 / render 1573B / six-package / token 覆盖） | ✅ |
| C3 drain harness A/B（linux 容器 = CI 等价） | **A**：compose stop → exit 0 + `shutdown.starting`→`shutdown.complete`；**B**：`HTTP_SHUTDOWN_TIMEOUT=1s` + 宿主 TcpClient 在途慢请求 → `shutdown.timeout`（context deadline exceeded）→ **ExitCode 1** | ✅ |
| C4 I-024-002 | CI 环境 = 本地等价 + linux 容器实跑（workflow 文件 = 交付物）；**hosted runner 实触发登记 R7 复核**（不主张 hosted acceptance · workspace-002 先例） | ✅（有界） |

## 残余与登记

1. **hosted runner 实触发**（workflow_dispatch / repository_dispatch）：登记 R7 收口复核项（本地/容器等价已证；hosted 触发 = 外部网络动作，随 R7 报告或用户指令）。
2. R1 残余 1（信号级 drain harness）→ **核销**（harness A/B 本次容器实跑完成）。
3. workspace-023 A-002 F-007（workflow 无 pnpm setup / 跨仓 token）→ **核销**（重构后 setup-node cache pnpm + 免 token）。