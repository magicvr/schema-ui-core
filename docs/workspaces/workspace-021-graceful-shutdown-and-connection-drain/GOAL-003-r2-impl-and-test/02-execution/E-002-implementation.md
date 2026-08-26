---
id: E-002
title: 实施落地（配置键 / main 接线 / compose 对齐 / 测试锁）
date: 2026-08-27
status: done
---

# E-002 · R2 实施（2026-08-27）

## 事实（code diff 全部落地并验证）

| 文件 | 变更 | 验证 |
|------|------|------|
| `internal/config/config.go` | `HTTPShutdownTimeout` 字段 + `http.shutdown_timeout` YAML 键（严格解析 `strictDurationPtr`，空/非法 → LoadError）+ env `HTTP_SHUTDOWN_TIMEOUT`（非法 → LoadError）+ `<=0` → LoadError（fail-closed）；默认 10s | `go test ./internal/config/...` ✓ |
| `internal/config/config.default.yaml` | `http.shutdown_timeout: 10s`（注释 §6） | 嵌入默认路径 ✓ |
| `apps/api/configs/config.yaml` | 同键（operator 副本） | — |
| `apps/api/configs/.env.example` | `HTTP_SHUTDOWN_TIMEOUT` 登记（TestCanonicalEnvExample 要求） | ✓ |
| `cmd/server/main.go` | `shutdownCtx` 改用 `cfg.HTTPShutdownTimeout`（去硬编码 10s） | vet ✓ |
| `compose.yaml` | api 服务 `stop_grace_period: 15s`（≥ 默认预算，SIGKILL 不截断排空） | compose 语法 ✓ |
| `internal/config/config_test.go` | `TestHTTPShutdownTimeout`：默认/YAML/env 覆盖/非法 YAML/空 YAML/非法 env/0s/-1s 全 fail-closed | 7/7 subtests ✓ |

**越界核对**：未改 Job runner / Store / 迁移 / Profile；未动请求级超时。合同 §2/§6 + A-001 F-001 责任履约。

## 验证 / 后续

- 全量 `go test ./...` 后台运行中；绿后 C3 自审并关门（E-003）。