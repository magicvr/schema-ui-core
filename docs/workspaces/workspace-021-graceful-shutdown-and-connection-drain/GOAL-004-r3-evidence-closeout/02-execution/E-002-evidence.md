---
id: E-002
title: harness 实施与证据（A/B/C · 进程级 · 双方言矩阵）
date: 2026-08-27
status: done
---

# E-002 · R3 证据（2026-08-27）

## 事实

| 证据 | 文件 | 结果 |
|------|------|------|
| A clean drain（进程内） | `internal/composition/shutdown_drain_test.go` | **PASS**（`-count=2` 确定性；in-flight 请求在预算内完成，Stop nil，响应 HTTP/1.1 到达） |
| B budget hole（进程内） | 同上 | **PASS**（Stop = `context deadline exceeded` 2.0s；探测读证明请求确在途中） |
| A′/B′ 进程级 SIGTERM | `cmd/server/shutdown_harness_test.go`（`//go:build !windows`） | 本机 skip（Windows）；linux/CI 责任：exit 0 + `shutdown.complete` / exit 1 + `shutdown.timeout` |
| C 中断重跑 | `internal/jobs/shutdown_reclaim_test.go` | **PASS**（Stop 后 running 无终态；新 Runner reclaim attempt+1 → succeeded） |
| PG 变体（A 的双方言版） | `TestShutdownDrainHarnessPostgres` | 本机 skip（PG_TEST_* 未置）；CI 责任 |
| 迁移台账 checksum | 全量 suite（`store` 包 `migration_catalog_test` checksum 冻结 + 双方言重启/Open 契约测试） | **PASS**（`go test ./...` exit 0 全绿） |
| 合同 §1/§7 日志事件 | `cmd/server/main.go`：`shutdown.starting` / `shutdown.complete` / `shutdown.timeout`（error） | committed `117f0486`；进程级 harness 断言 |

**调试留痕**：初版 harness 两次 FAIL 根因 = Stop 与 accept 竞态（请求尚未被服务器挂起时 Shutdown 即返回 nil），非契约缺陷；以「探测读（400ms deadline 无响应 ⇒ in-flight）后再 Stop」修复，`-count=2` 复跑确定性 PASS。

## 验证 / 后续

- 全量回归 `go test ./...`（apps/api）exit 0。C3：关门审计（self + grok independent）→ Root 关门。