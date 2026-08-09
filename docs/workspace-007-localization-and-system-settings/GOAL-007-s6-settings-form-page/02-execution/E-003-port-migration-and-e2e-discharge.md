---
id: E-003
doc: execution
title: S6 · F-002 解除 — 端口迁移 25080/25173 + e2e 双 profile 补跑通过
status: recorded
parent: GOAL-007-s6-settings-form-page
created: 2026-08-09
updated: 2026-08-09
version: 0.1.0
---

# E-003 · F-002 解除：端口迁移 + e2e 补跑（2026-08-09）

## 事实

- **端口迁移（用户指令，commit `01ed50d`）**：默认端口迁移到 25000+ —— API `HTTP_ADDR` `:8080`→`:25080`、web dev/preview `5173`→`25173`、compose web 宿主 `8081`→`25081`、nginx `proxy_pass api:25080`；同步 `.env.example`、Dockerfile EXPOSE、playwright readyz、CI `r6-basic-matrix.yml`、`scripts/smoke.sh`、README×4。`25080/25173/25081` 均不在 Windows Hyper-V 排除区间（`netsh` 复核）。
- **shell.spec 既有跨用例 DB 污染修复（commit `c7f63d9`）**：localization M3 把 `siteTitle` 写进共享 playwright SQLite，shell.spec 硬编码断言默认品牌失败（S5 只跑单 spec 未暴露）；改为从公开 `/api/branding` 读当前 `siteTitle` 断言，对 prior 测试状态健壮。
- **e2e 双 profile 补跑（F-002 复审触发解除）**：`APP_PROFILE=admin` 与 `mvp` 各 **3 passed / 1 skipped（对应另一 profile 的测试）**，M3 设置表单保存→投影在真实栈上通过。日志：`attachments/s6-e2e-admin.log`、`attachments/s6-e2e-mvp.log`。
- 验证：vitest **728/728**、`npm run build`、`go test ./apps/api/...` 全绿。

## 产物

| 路径 | 说明 |
|------|------|
| `apps/api/internal/config/config.go`、`.env.example`、`Dockerfile` | API 默认 `:25080` |
| `apps/web/{vite.config.ts, playwright.config.ts}` | web `25173`、proxy/readyz `:25080` |
| `compose.yaml`、`apps/web/nginx.conf` | compose/nginx 端口同步 |
| `.github/workflows/r6-basic-matrix.yml`、`scripts/smoke.sh` | CI/smoke 端口同步 |
| `apps/web/e2e/shell.spec.ts` | 品牌断言读当前 siteTitle |
| `attachments/s6-e2e-{admin,mvp}.log` | e2e 补跑证据（F-002 解除） |

## 里程碑 checkpoint

- commit：`01ed50d`（端口迁移）、`c7f63d9`（shell.spec 修复）。
