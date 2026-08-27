---
id: D-001
title: R3 证据 harness 方案（合同 §8 · A/B/C + 双方言矩阵）
date: 2026-08-27
status: accepted
---

# D-001 · R3 harness 方案（2026-08-27）

## 决策

按合同 v0.1.0 §8 实施三层证据（前台为本地可跑、真断言；无新增 API 面）：

| # | harness | 形态 | 断言 | 运行位 |
|---|---------|------|------|--------|
| A | clean drain（进程内） | `internal/composition/shutdown_drain_test.go`：真实 app（mvp + SQLite）→ 拖住 body 的 in-flight POST `/api/auth/login`（探测读证明 in-flight）→ `Stop(5s)` → 客户端补完 body | Stop = nil；in-flight 请求收到 HTTP 响应（§1 步骤 2/3） | 所有 OS |
| B | budget hole（进程内） | 同上但 body 永不补完 → `Stop(2s)` | Stop = deadline error（§2 强制退出语义） | 所有 OS |
| A′ | clean drain（进程级 · linux/CI） | `cmd/server/shutdown_harness_test.go`：`go build` → 真实进程 + 探测 in-flight → **SIGTERM** → 补完 body | exit 0 + `shutdown.complete` 日志（§1/§3/§7） | linux/macOS（Windows 不支持 SIGTERM→子进程，skip） |
| B′ | budget hole（进程级 · linux/CI） | 同上 + `HTTP_SHUTDOWN_TIMEOUT=1s`，body 永不补完 | **exit 1** + `shutdown.timeout` 日志（§3） | linux/macOS |
| C | 中断重跑（仓库内） | `internal/jobs/shutdown_reclaim_test.go`：运行中 Job 被 `Runner.Stop` 中断（无终态转移）→ 租约过期 → 新 Runner reclaim（attempt+1）→ succeeded | attempt=2；interrupt 后 status 仍 running（§4） | 所有 OS |
| 双方言 | A 的 PG 变体 | `TestShutdownDrainHarnessPostgres`（PG_TEST_* 门控，CI） | 与 A 同断言（§5 双方言一致） | CI（有 PG） |
| 迁移台账 | 回归锁 | 全量 suite 的 `migration_catalog_test`（checksum 冻结断言）+ store 双方言重启测试 | checksum 不变 | 所有 OS |

## 合同 §1/§7 日志事件

main.go 增加三事件：`shutdown.starting`（信号后首行）、`shutdown.complete`（干净收尾）、`shutdown.timeout|error`（预算耗尽/错误，exit 1）——进程级 harness 直接断言。

## 验收判据

- 本地（Windows）：A/B/C + store 双测锁绿；进程级 A′/B′ skip（记录为 CI 责任）。
- CI（linux + PG）：A′/B′ + PG 变体绿。
- 越界：无新增 API 路由；无 runner/Store/迁移改动；无 Profile 变更。

## 未选方案

- 新增测试专用慢端点（污染 API 面）→ 改掷「in-flight 请求 = 未完成 body + 探测读」，零 API 变更。
- Windows 进程级 SIGTERM（Go 不支持向子进程发 SIGTERM）→ 进程内等价 harness + linux/CI 进程级。