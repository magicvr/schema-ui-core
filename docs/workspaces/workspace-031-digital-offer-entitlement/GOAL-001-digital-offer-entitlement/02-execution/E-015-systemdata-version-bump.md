---
doc_type: goal-execution
id: E-015-systemdata-version-bump
parent: GOAL-001-digital-offer-entitlement
date: 2026-09-06
status: done
version: 1.0.0
---

# E-015 · 关门后维护：dev.cmd start 启动失败（系统数据 checksum 漂移）修复

## 事实（时间线）

- 2026-09-06 · **用户报告**：`.\dev.cmd start` 启动失败。
- 2026-09-06 · **复现与根因**：`go run ./cmd/server`（configs/config.yaml）直接复现 —— `LIFECYCLE_START_FAILED [core.auth-session]: reconcile system data: system-data navigation/menu_digitaloffer_offers checksum drift`。
  - `systemdata.Reconcile.checkLedger` 的防篡改语义：**同版本 + checksum 不一致 = 漂移 → fail-closed 拒启**（防对系统数据表未声明篡改）。
  - 漂移来源：E-014 把 `menu_digitaloffer_offers` 导航贡献 Label 改为 "Digital products"，但 `SystemDataVersion` 未递增 → 代码 checksum 变了、ledger 还是 v1 旧值。
  - 为何测试全绿：测试用内存库每次从零登记，无旧 ledger；用户的持久化 `data/schema-ui.db` 里存有 v1 旧 checksum，才暴露。
- 2026-09-06 · **修复（设计机制）**：`apps/api/modules/authsession/systemdata/policy.go` `SystemDataVersion` 1 → 2。reconcile 对 `ledger v1 < code v2` 接受并按新 checksum 刷新 ledger；下次启动 v2 一致即干净。新增 `menu_digitaloffer_purchases`（E-013）为全新条目，本无旧 ledger，不触发漂移。
- 2026-09-06 · **验证**：
  - `go test ./modules/authsession/... ./kernel/... ./internal/composition/... ./internal/handler/... ./modules/digitaloffer/... ./internal/store/...` 全绿（含 reconcile drift/upgrade 测试，硬编码 v1 的 fixture 不受常量影响）。
  - 真实开发库启动成功（readyz 200），manifest 侧栏三入口齐全、顺序正确（预付凭证 → 数字商品 → 数字权益 → 数字订单 → 操作日志），既有 2 条 offer 数据完好。
  - 完整模拟用户路径：`dev.cmd start --no-browser` → API listening + /readyz 200 + Web listening :25173，exit 0；随后 `dev.cmd stop` 清理交还环境。

## 产物路径

- `apps/api/modules/authsession/systemdata/policy.go`（SystemDataVersion 2 + 注释）

## 进度评估

- 关门状态不变（Root done 4/4 · VP-031 closed v0.3.4）。低风险可逆维护。
- 教训留痕：变更系统数据贡献内容（导航 label/order/policy 等）必须同步递增 `SystemDataVersion`，否则既有库启动 fail-closed；内存库测试无法覆盖该场景（新增/变更后应至少对持久化库做一次启动冒烟）。
