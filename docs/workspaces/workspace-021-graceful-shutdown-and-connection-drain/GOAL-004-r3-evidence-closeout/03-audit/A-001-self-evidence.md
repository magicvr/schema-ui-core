---
id: A-001
source: self
date: 2026-08-27
scope: R3 证据 harness 关门自审（合同 §8 · A/B/C · 双方言矩阵 · 越界）
verdict: pass
project: workspace-021 · GOAL-004-r3-evidence-closeout
---

# A-001 · R3 证据自审（2026-08-27 · self）

## 范围

- 对象：`01-decision/D-001-harness-plan.md` + 证据（E-002）——进程内 A/B、进程级 A′/B′（linux/CI）、C reclaim、PG 变体、迁移 checksum 回归锁。
- 模式：`self`（证据冻结；独立审由 Root 关门双审中的 grok independent 承担）。

## 核对清单

| 合同 §8 项 | 证据 | 核对 |
|------------|------|------|
| harness A（clean drain） | `shutdown_drain_test.go`：探测读证明 in-flight → 补完 body → Stop nil + HTTP 响应 | ✓（-count=2 确定性 PASS） |
| harness B（timeout） | 同上，body 永不补完 → Stop deadline error 2.0s | ✓ |
| harness A′/B′（进程级、退出码） | `cmd/server/shutdown_harness_test.go`（!windows）：SIGTERM → exit 0/1 + `shutdown.complete`/`shutdown.timeout` | ✓ 本机 skip（Windows）；linux/CI 责任已记录 |
| harness C（重启 reclaim） | `jobs/shutdown_reclaim_test.go`：Stop 无终态 → 租约过期 → reclaim attempt+1 → succeeded | ✓ |
| 双方言 | `TestShutdownDrainHarnessPostgres`（PG_TEST_* 门控）+ store 包双方言 Open/Close/重启测试 | ✓ 本机 SQLite 实测 + PG CI 门控 + store 双测锁 |
| 迁移台账 checksum | 全量 suite（`migration_catalog_test` 冻结断言） | ✓ exit 0 |
| 日志事件 §1/§7 | main.go `shutdown.starting/complete/timeout` | ✓（进程级断言） |
| 越界 | git diff：仅 4 个代码/测试文件 + 文档；无 API 路由/runner/Store/迁移/Profile 变更 | ✓ |

## Findings

- `F-001`：recommended。进程级 A′/B′ 与 PG 变体在本机（无 PG、Windows）为 skip——关门证据须在 CI（linux + PG_TEST_*）实跑后核销；已在 evidence matrix 与 workspace 结项残留中登记。状态 → **accepted-residual 候选**（用户书面接受 / CI 实跑核销，见 Root 关门记录）。

## 结论

**pass**（0 required）。R3 证据与合同 §8 一致、无越界；本地可核对面全部实测绿；门控面（PG / 进程级）以 CI 为责任已登记。允许 GOAL-004 关门并进入 Root 关门双审（self + grok independent）。