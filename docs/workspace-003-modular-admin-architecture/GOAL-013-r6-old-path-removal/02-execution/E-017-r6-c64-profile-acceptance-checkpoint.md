---
id: E-017-r6-c64-profile-acceptance-checkpoint
doc: execution-entry
goal: GOAL-013-r6-old-path-removal
source: orchestrator
date: 2026-08-06
status: recorded
---

# E-017 · R6 C6.4 双 Profile 验收接线 checkpoint

## 已发生事实

- 提交 `88a3840`（`test(r6): add dual-profile acceptance matrix`）完成第二个 C6.4
  实现切片：
  - Playwright 显式只接受 `mvp` / `admin`，向 Go API 与 Vite 传递同一 `APP_PROFILE`，
    并等待 `/readyz`；Shell E2E 动态断言 Manifest source、Profile 页面集合、Schema
    200/404 与受保护路由 401/404；
  - `compose.yaml` 显式透传 `APP_PROFILE` / `APP_MODULES_ENABLED`，支持可配置宿主端口；
    `scripts/smoke.sh` 新增 `SMOKE_EXPECTED_PROFILE`，比较 API/Web Manifest bytes 与
    `X-Schema-UI-Manifest-Source: api`，覆盖 Profile 页面、Schema、路由和重启 readiness；
  - `.github/workflows/r6-basic-matrix.yml` 增加 API vet、退休路径存在性/引用扫描、
    browser E2E `mvp/admin` matrix 与 container smoke `mvp/admin` matrix；
  - README、QUICKSTART、API/Web README 与 `.env.example` 统一为真实环境键
    `APP_PROFILE` / `APP_MODULES_ENABLED`，明确 Go 不自动加载 `.env`、owner module
    接入方式以及升级/恢复边界。
- checkpoint 前完成动态验证：
  - workflow YAML parse、Git for Windows Bash `bash -n scripts/smoke.sh`、`mvp` 与
    `admin` 两次 `docker compose config --quiet` 均成功；
  - `apps/api`: `go test -count=1 ./...`、`go vet ./...`、`go build ./...` 均退出 0；
  - `apps/web`: `npm test` 为 24 files / `495/495`，`npm run build` 成功；
  - Playwright Chromium 在 `APP_PROFILE=mvp` 与 `APP_PROFILE=admin` 下各为 `2/2`；
  - 非治理范围宽搜未发现旧 handler Schema 路径、Web public Manifest 路径或旧
    `MODULES_ENABLED` 键；本地遗留的空 `handler/fixtures/` 目录已删除；
  - `git diff --cached --check` 通过，提交只含 10 个 owned paths；任务开始前的
    `account_test.go`、`auth_test.go`、`health_test.go` 换行状态未暂存。

## 状态边界

- 本条证明 C64-V01～V03 的实现接线与本地基础回归成立；workflow matrix 尚未形成
  hosted CI 结果，本地 Compose 解析也不等于镜像构建、容器运行或重启持久化证据。
- C64-V04～V07 的完整升级/恢复、双 Profile 容器、custom/fail-closed 与 clean committed
  clone/fork 复现仍开放。
- R6-I004 保持 `collecting`，C6.4 不勾选，GOAL-013 保持 `active / 3/4`，Root 保持
  `active / 5/6`。

## 下一步（计划）

按 D-004 执行 C64-V01～V07 全部终态验证，形成逐条 evidence attachment 后再进入
GOAL-013 self + Grok Build independent 关门审计。
