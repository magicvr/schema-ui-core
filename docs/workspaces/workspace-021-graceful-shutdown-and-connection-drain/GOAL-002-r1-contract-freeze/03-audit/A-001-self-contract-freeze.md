---
id: A-001
source: self
date: 2026-08-27
scope: R1 合同冻结关门自审（合同 v0.1.0 vs 基线证据 / 信息裁决 / 边界）
verdict: pass
project: workspace-021 · GOAL-002-r1-contract-freeze
---

# A-001 · R1 合同冻结自审（2026-08-27 · self）

## 范围

- 对象：`01-decision/D-002-contract-freeze.md`（合同 v0.1.0）+ C1 信息裁决（D-001 accepted）。
- 模式：`self`（方案冻结，低风险、可逆、无代码变更；R2 实施前门禁）。

## 核对清单（合同条款 ↔ 基线证据）

| 合同项 | 证据/事实 | 核对 |
|--------|-----------|------|
| §1 停机顺序步骤 2/3（Shutdown 拒绝新连接 + 存量排空） | `cmd/server/main.go` `app.Stop` → `composition.go` OnStop `srv.Shutdown(ctx)` | ✓ |
| §1 步骤 5/6/7（Job → runtime → Store）与 §4 | `OnStop` 顺序：`srv.Shutdown` → retention/metrics → `jobs.Stop` → `runtime.Stop` → `st.Close` → tracing；`runner.go` Stop cancel + `finish()` stopping 放弃转移 + lease-reclaim（`repository.go` Claim/ListRunnable/ExhaustExpiredJobs） | ✓ |
| §3 退出码 | `main.go`：Stop err → `os.Exit(1)`，否则 exit 0 | ✓ |
| §5 迁移仅启动期 / fail-closed | `store/open.go` + `migrate.go`（startup plan：fresh/noop/apply-pending/restore-ledger；一迁移一事务 + checksum）；`postgres.go` `migrate` 同合同；服务期无迁移入口 | ✓ |
| §6 配置键命名与 10s 默认 | `config.go` `http:` 段 + `HTTP_*` env 前缀（如 `HTTP_ADDR`）；`main.go` 现硬编码 10s；server.go timeout 配置驱动 | ✓ |
| §2 部署对齐 / compose | `compose.yaml`：`restart: on-failure`、无 `stop_grace_period`（默认 ~10s）——合同建议显式 `15s`（R2 落地） | ✓ |
| 边界不越（A3 / Profile / 迁移台账不改） | 合同 §0 / §5 / 未选方案；VP-021 首波冻结投影一致 | ✓ |

## Findings

- `F-001`：recommended。R2 实现必须以本合同 §1–§7 为验收分母；`http.shutdown_timeout` 落地时须带 fail-closed（非法值拒绝启动）测试锁；compose `stop_grace_period: 15s` 显式化。→ **fixed**（R2 目标 00-meta 已承接，2026-08-27）。

## 结论

**pass**（0 required）。合同 v0.1.0 与基线事实一致、边界无越界、信息门禁已由用户裁决关闭。允许 R1 关门；R2 以本合同实施。