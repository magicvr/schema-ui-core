---
id: D-001
doc: decision-entry
goal: GOAL-007-w6-container-smoke-reproducibility
status: accepted
created: 2026-08-14
updated: 2026-08-14
version: 0.1.0
---

# D-001 · W6 范围：claim GIT_COMMIT 接线修复

## 背景（F-1 = F-1a + F-1b，均为 post-go）

- **F-1a**：`generate-claim.mjs` 自 W3（`5e4c384`，2026-08-13，post-go）强制 `git rev-parse HEAD`；`.dockerignore` 排除 `.git`。→ compose web 镜像构建必然失败（实测 `fatal: not a git repository`）。
- **F-1b**（W6 验证中发现）：`apps/web/nginx.conf` 自 VP-009 W3（`7dbc3b5`，安全加固）把 `upstream api_upstream` 块**嵌进 `server {}` 内部**（nginx 只允许 http 层）→ nginx `[emerg]` 启动崩溃，web 容器不可达，登录代理失败。

→ V-007/V-008 与 CI `container-smoke` 不可复现（构建 + 运行双断点）。按 VP-008 量尺属「使冻结证据不可复现 + 阻止标准 Admin 模块容器构建/运行」；用户 2026-08-14 裁决 **A：以 VP-010 W6 波次修复**。

## 决策

1. `generate-claim.mjs`：buildId 取 `process.env.GIT_COMMIT`（trim 非空）否则 git 回退；本地行为不变。
2. `apps/web/Dockerfile`：build 阶段 `ARG GIT_COMMIT` + `ENV GIT_COMMIT=${GIT_COMMIT}`。
3. `compose.yaml`：web build `args: GIT_COMMIT: ${GIT_COMMIT:-unknown}`（未传时为 `git:unknown`，确定性可复现）。
4. CI `r6-basic-matrix.yml`：container-smoke env 加 `GIT_COMMIT: ${{ github.sha }}`。
5. `apps/web/nginx.conf`（F-1b）：`upstream api_upstream` 移到 `server {}` 之外（http 层）；`location /api/` 的 `proxy_pass http://api_upstream` 与安全头保持不动。

**边界**：不改 claim 内容 / 校验语义、协议 pin、Profile/模块矩阵 / 运行时行为；不改代理语义（仅修正配置作用域）。

## 审计模式

**self**（小改动 + 容器构建/smoke 实测证据）。

## 未选方案

- `COPY .git` 进构建上下文：膨胀镜像、泄露历史；不选。
- 跳过容器内 claim 生成（预生成拷贝）：改变 claim 与构建解耦语义、弱化证据链；不选。
