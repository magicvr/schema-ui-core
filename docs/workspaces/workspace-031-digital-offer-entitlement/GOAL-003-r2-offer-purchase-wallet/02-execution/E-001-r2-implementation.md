---
doc_type: goal-execution
id: E-001-r2-implementation
parent: GOAL-003-r2-offer-purchase-wallet
date: 2026-09-05
status: done
version: 1.0.0
---

# E-001 · R2 实施（C1/C2/C3）

## 事实（时间线）

- 2026-09-05 · **C1 切片方案**：文件/接口清单按 D-002 v1.0.0 条款映射落定——`apps/api/modules/digitaloffer/`（migration v0070 / store / service 子包 / provider / manifest / schema 两页）+ `internal/handler/digitaloffer.go`（路由构建，wallet 先例）+ errorcatalog / compiled 持久化注册 / composition `plan.HasModule("biz.digital-offer")` 装配 / web i18n 错误码与页面文案。service 放子包（`digitaloffer/service`）以避免 handler→模块根导入环（wallet/voucher 先例）。
- 2026-09-05 · **C2 模块骨架落地**：
  - 迁移 `0070 digital_offers`：`digital_offers` / `digital_purchases` / `digital_entitlements` 三表 + 形态互斥 CHECK + 唯一约束与索引；SQLite INTEGER / PG BIGINT 双变体；经 `modules/compiled/persistence.go` 注册进全局目录（每次 `testsupport.OpenStore` 均应用）。
  - store：Offer/Purchase/Entitlement 域模型；插入/查询/乐观锁更新；`InsertPurchaseInTx` 以 `kernel.IsUniqueViolation` 分类自表唯一竞争；`VoidEntitlementInTx` 幂等（RowsAffected=0 + 行存在 = 已作废）；消耗候选与五条件原子扣减。
  - service：`Purchase` 按 §4.2/§4.4——attempt 循环在 `store.Run` 外（≤3 次）、attempt 起点回读（幂等重放 / `ErrRequestIdConflict` 终态）、单事务内 offer/subject/账户/币种校验 → `MutateInTx(freeze)` → `MutateInTx(deduct_frozen)`（均带 `ref_type=biz_offer_purchase` 反链 + 互异幂等键 `<request_id>:freeze`/`:deduct`）→ 凭证 → 权益；Admin 写（create/update/status/void）走 `TransactionalRecorder.RecordOperationTx` 同事务 fail-closed 审计，重复 void 不追加审计。
  - wallet/subject 增加只读 `SubjectExistsInTx`（购买 tx 内主体存在性门控，增量无行为变更）。
  - provider：`biz.digital-offer` 权限键（read / offer.manage / entitlement.void）、6 条 Admin 路由 + 公开 `GET /api/biz/offers`（IP 桶 AllowRecord、无内部字段）、两个 schema 页、manifest fragment、导航贡献；不进默认 Profile。
  - errorcatalog + web i18n（en-US/zh-CN）登记 BIZOFFER_* 码与 schema.digitaloffer.* 文案；错误契约钉死测试同步。
- 2026-09-05 · **D-001 合同附录**：§9 加法修订 `BIZOFFER_VERSION_CONFLICT` / `INVALID_BIZOFFER_REQUEST`（D-002 → v1.1.0）；`BIZOFFER_ENTITLEMENT_INVALID` 目录登记随 R3。
- 2026-09-05 · **C3 购买链路与验收**：
  - `service/purchase_test.go`：单事务可见状态（freeze+deduct 恰好各一笔、ref 反链、无冻结残留、duration/count 权益正确）、五条终态失败路径（无凭证/权益/冻结残留）、幂等重放与跨 offer 冲突、同 request 并发双发收敛（恰一凭证/一次扣款/一份权益）、余额恰够一次的并发竞争（恰一成功）、币种不匹配、审计 fail-closed（注入审计失败 → 域写回滚）与 void 幂等。
  - **PostgreSQL 验收（真库）**：`TestPurchasePostgresAcceptance` 经 pgtest DSN 建库 + 全量目录迁移后通过（双数据库要求达成）。
  - `internal/handler/digitaloffer_test.go`：未认证 401 门控、创建/状态变更（乐观锁）、形态不可变 409、公开目录只含 on_sale 且不泄漏内部字段、审计失败 → 域行回滚、购买后凭证/权益经 Admin 面可见、公开目录限流 429 + Retry-After。
  - 全仓 `go test ./...` 回归通过（唯一改动文件为错误契约钉死清单的预期同步）。

## 产物路径

- `apps/api/modules/digitaloffer/`（migration/store/service/provider/manifest/schema）
- `apps/api/internal/handler/digitaloffer.go` + `digitaloffer_test.go`
- `apps/api/internal/errorcatalog/errorcatalog.go`、`apps/api/modules/compiled/persistence.go`、`apps/api/internal/composition/composition.go`
- `apps/web/src/i18n/messages/{en-US,zh-CN}.json`
- `apps/api/modules/wallet/subject/subject.go`（SubjectExistsInTx 增量）

## 进度评估

- C1/C2/C3 关门；C4（self + codex independent 审计与关门）待执行。
