---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-004-r3-compose-cicd
version: 0.1.0
---

# D-001 · R3 实跑范围与方式（2026-08-29）

## 决策

1. **compose 实跑 = 全服务**：`docker compose up -d --build`（api + web）→ api `service_healthy`（readyz）→ web 反代可达；`docker compose stop`（`stop_grace_period: 15s` ≥ 10s 预算）→ 断言 api 日志含 `shutdown.starting` / `shutdown.complete` 且 **ExitCode 0**。AUTH_JWT_SECRET / ADMIN_INITIAL_PASSWORD 由命令环境注入（不落盘、不写 .env）。
2. **consumer-regression workflow 重构**：删除 GH Packages token 步骤（npmjs 公开免凭据）；`setup-node` 挂 pnpm cache；Go 侧 = `go get @latest` + `go build` + serve 后台 + healthz/readyz 轮询 + 四探针 + **SIGTERM 收尾并断言 `shutdown.complete`**（RT-D02 出口，workflow 即 harness A 的主机形态）。
3. **信号级 drain harness（linux 容器 = CI 等价）**：
   - A：compose stop 路径（上）；
   - B：`docker run` api 镜像 + `HTTP_SHUTDOWN_TIMEOUT=1s`，宿主 PowerShell TcpClient 保持**在途慢请求**（部分头部不闭合），`docker stop` → 断言容器 **ExitCode 1** + 日志 `shutdown.timeout`。
4. **Windows 本地不伪造 SIGTERM**（信号面以容器/CI 证据为准）；本地等价实跑 = build/安装/探针/serve 起停。
5. **hosted runner 实触发登记**为 R7 复核项（GitHub-hosted acceptance 不主张——workspace-002 先例「远端 CI 尚未触发不主张 hosted acceptance」）；workflow 文件本身 = 交付物。

## 未选方案

| 项 | 未选 | 理由 |
|----|------|------|
| hosted 触发 | 本轮推送触发 workflow_dispatch | 私有仓 Actions 需 runner 可用 + 触发副作用；等价证据（容器/本地）充分，实触发登记 R7 |
| harness B 在容器内施压 | busybox/alpine 无 nc/python | 宿主慢客户端更简单可控 |
| 双栈 web 冒烟 | 浏览器 e2e 进本波 | 判据 #3 只要求 compose/CI 可运行性；e2e 非本波分母 |