---
id: D-001
title: 实现计划：shutdown_timeout 配置键 / main 接线 / compose 对齐 / 测试矩阵
date: 2026-08-27
status: planned
---

# D-001 · R2 实现计划（2026-08-27 · 合同 v0.1.0 §2/§6 + A-001 F-001）

## 变更面（全部在合同「本波允许」内）

| # | 文件 | 变更 | 合同条款 |
|---|------|------|----------|
| 1 | `internal/config/config.go` | `Config` 增加 `HTTPShutdownTimeout time.Duration`；`yamlFile.http` 增加 `shutdown_timeout *string`；默认 `10s`；YAML/env 严格解析（解析失败 → `LoadError`）；`<=0` → `LoadError`（fail-closed，任何环境） | §6 |
| 2 | `internal/config/config.default.yaml` | `http:` 段增加 `shutdown_timeout: 10s`（含注释） | §6 |
| 3 | `apps/api/configs/config.yaml` | 同键（operator 副本） | §6 |
| 4 | `cmd/server/main.go` | `shutdownCtx` 超时改用 `cfg.HTTPShutdownTimeout`（去硬编码 10s） | §1/§3/§6 |
| 5 | `compose.yaml` | api 服务增加 `stop_grace_period: 15s`（≥ 默认预算；编排宽限不截断排空） | §2 部署对齐 |
| 6 | `internal/config/config_test.go` | 测试锁：默认 10s / YAML 覆盖 / env 覆盖 / 非法 YAML 值 fail-closed / 非法 env 值 fail-closed / `0s` fail-closed | §6 |

## 校验语义

- 解析失败（YAML 或 env）：`cfg.LoadError` 置错 → `ValidateProd` 第一道即拒绝启动（任何环境）。
- 语义非法（`<=0`，如显式 `0s` / `-1s`）：`Load` 末尾置 `LoadError`；零值 Config（测试绕过 Load）不受影响（与 validateDB 的加载器旁路约定一致）。
- env 优先级不变：env 覆盖 YAML；均缺省时 10s。

## 越界（禁止）

不改 Job runner / Store / 迁移台账 / Profile 默认集 / Manifest；不新增其他配置键；不动 `server.New` 的请求级超时。

## 验收

- `go vet ./...` 与 `go test ./internal/config/...`（及受影响包）全绿；
- 新增测试全部通过；
- 手工核对（可选）：`HTTP_SHUTDOWN_TIMEOUT=bogus` 时 `ValidateProd` 报错。