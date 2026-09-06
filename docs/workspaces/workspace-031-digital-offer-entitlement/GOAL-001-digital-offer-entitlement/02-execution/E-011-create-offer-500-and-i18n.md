---
doc_type: goal-execution
id: E-011-create-offer-500-and-i18n
parent: GOAL-001-digital-offer-entitlement
date: 2026-09-06
status: done
version: 1.0.0
---

# E-011 · 关门后维护：新建 Offer 500 + 中英文对照优化

## 事实（时间线）

- 2026-09-06 · **用户报告**：新建数字 Offer（「新建 Offer」表单提交）报错；要求顺带修正数字 Offer 相关中英文对照（尤其中文描述不贴日常直觉）。
- 2026-09-06 · **根因排查**：
  - 对 `POST /api/digitaloffer/offers` 以合法负载实测 → **HTTP 500 INTERNAL**（"digital offer operation failed"）。`CreateOffer` 在事务内先 `InsertOfferInTx` 再 `recordOperationTx`（D-002 §7 fail-closed 配对），审计行写入 `operation_log` 时触发 **`event` CHECK 约束**——`operation_log` 的 CHECK 白名单从未加入 `bizoffer.*` 事件（`bizoffer.offer.create/update/status`、`bizoffer.entitlement.void`），插入被拒导致整事务回滚 → 500。
  - 为何此前测试全绿：handler 测试使用 `stubTxRecorder`（不走真实 `operation_log`），composition 测试只做 GET 未做 POST，故未触达真实 CHECK。该缺口属 R2 关门时的遗漏（A-008 等审计未覆盖真实 audit 写入路径）。
- 2026-09-06 · **修复（migration 0071）**：`core.operationlog` 新增 `operation_log_digitaloffer_events` 迁移，按 0053 模式重建 `operation_log`（保留 correlation + session 侧表），把四个 `bizoffer.*` 事件并入 CHECK 枚举。同步 `identity.go` catalog head=71、`lockedHeadExtraTables[71]={}`、store 迁移尾部断言与 catalog 校验表（checksum `4520ca0c…cb246`）；`TestMigrateFreshDB` 增补 `bizoffer.offer.create` 写入选通断言。
- 2026-09-06 · **i18n 优化（zh-CN 为主，贴合日常直觉）**：
  - `schema.digitaloffer.field.priceAmount`：标价（最小货币单位）→ **标价（单位：分）**
  - `schema.digitaloffer.field.entitlementForm`：权益形态 → **权益类型**（`column.form` 形态 → 类型，`error.bizOfferFormConflict` 同步）
  - `schema.digitaloffer.field.durationSeconds`：时长（秒）→ **有效时长（秒）**
  - `schema.digitaloffer.field.countPerPurchase`：每次购买次数 → **单次购买获得次数**
  - `schema.digitaloffer.field.subjectId` / `column.subjectId`：主体（ID）→ **归属主体（ID）**
  - `error.invalidBizOfferRequest`：数字 Offer 请求无效 → **提交的 Offer 信息无效，请检查必填项和数值范围**
  - en-US：`Price (min units)` → `Price (minor units)`（其余 en-US 已准确，不改）
- 2026-09-06 · **验证**：
  - `apps/api` 全包 `go test ./...` 全绿；web vitest 全包 92 文件 / 1218 测试通过（含 i18n 键集一致性）。
  - 重启 dev API 后实测：duration 与 count 两种表单新建均 **HTTP 201**；`operation_log` 出现 `bizoffer.offer.create` 审计行（fail-closed 配对闭合）。

## 产物路径

- `apps/api/modules/operationlog/migration/migration.go`（0071 迁移 + 共享 rebuild helper）
- `apps/api/internal/store/identity.go`、`identity_test.go`、`migrate_test.go`、`restart_test.go`、`operations_test.go`
- `apps/web/src/i18n/messages/zh-CN.json`、`en-US.json`

## 进度评估

- 关门状态不变（Root done 4/4 · VP-031 closed v0.3.4）；低风险可逆维护，审计模式 none。
- 遗留提示：`bizoffer.offer.create/update/status`、`bizoffer.entitlement.void` 现已在 CHECK 白名单内；后续若新增本域事件需再次走 operationlog 迁移。
