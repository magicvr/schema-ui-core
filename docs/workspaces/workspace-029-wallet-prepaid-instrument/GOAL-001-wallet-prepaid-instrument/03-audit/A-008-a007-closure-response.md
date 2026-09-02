---
id: GOAL-001-wallet-prepaid-instrument
doc: audit-entry
record_id: A-008
status: recorded
parent: null
created: 2026-09-02
updated: 2026-09-02
version: 0.1.0
---

# A-008 · A-007 复审登记与 A-005 recommended（F-002～F-005）闭合记录

- **source**：self（编排器响应，非 independent）
- **date**：2026-09-02
- **scope**：对 Root `[workspace-029-wallet-prepaid-instrument/GOAL-001-wallet-prepaid-instrument]` A-007（independent · pass）结果登记 + A-005 F-002～F-005（recommended）处置核验
- **verdict**：**pass**（A-007 open required = 0 维持；F-002～F-005 recommended 全部按 `fixed` 闭合；无新增 required）

## Findings 逐条响应与闭合证据

| Finding | 级别 | 内容 | 闭合路径 | 闭合证据与说明 |
|---------|------|------|----------|----------------|
| **A-005 F-002** | recommended / med | 导出是渲染器对 `items[].code` 的启发式，协议未声明 download；换客户端/其它表单会误下载 | **fixed** | `wallet-vouchers.json` `generateBatch.onSuccess` 显式声明 `downloadCsv`（列序 + `vouchers_{batchId}.csv`）；`render.tsx` 删除启发式，仅按声明（含 `code` 列 + 文件名模板）在成功手势导出 CSV；vitest：声明化 fixture 文件名断言 PASS + 新增阴性用例「无关表单带 `items[].code` 且无声明 → 不下载」。renderer 全目录 290 tests PASS。 |
| **A-005 F-003** | recommended / low | PostgreSQL 上无 Redeem / 并发双花 e2e | **fixed** | 新增 `internal/store/postgres_voucher_redeem_test.go`（`PG_TEST_*` 门控，CI postgres job 点亮）：同一新主体并发两张不同卡（500+1500）→ 单 subject 户、余额 2000、账本 2 条；重复核销 fail-closed 不双记。本会话 docker postgres:15 实证 **PASS 2.18s**（0065 PG 方言随全目录应用，间接验证）。 |
| **A-005 F-004** | recommended / low | `batch_id` 无唯一约束，重复生成混批 | **fixed** | 迁移 0065 `voucher_batches`（batch_id 主键注册表 + 回填；历史重复批次不拆解、约束防 NEW 混批，代码注释留痕）；`GenerateBatch` 同事务先登记、0 影响行即 `ErrVoucherBatchExists` 拒绝整批（并发 fail-closed）；HTTP 409 `VOUCHER_BATCH_EXISTS`（errorcatalog + 冻结清单 + i18n messageKey 同步）；迁移目录台账测试同步 v65（migrate/operations/restart/identity/completeLostLedgerTables）。服务层 + HTTP 层新测试 PASS。 |
| **A-005 F-005** | recommended / low | 过期字段是 Unix 秒 `inputNumber`，易填毫秒或留空 | **fixed** | 生成端点 fail-closed 范围校验：`expiresAt` 须为 Unix **秒** ∈ [1e9, 4102444800]（2001-09-09～2100-01-01），毫秒/越界 → 400 `INVALID_VOUCHER_PARAMS`（目录 en/zh 文案含范围）；≤0/缺省保持无过期语义；schema 标签与 en/zh i18n 提示范围。`TestVouchersGenerateExpiresAtValidation` PASS。 |

## 附带说明（非 finding）

- A-007 为 independent `pass`（F-001 闭合复审），随本条目登记入索引；F-002～F-005 在 A-007 中维持 open recommended，本条目为其正式处置留痕。
- F-002 的 CSV 正文形状与历史一致（列序/引号），协议文档新增声明不改变上游 pin 与 Manifest。
- 0065 回填对「0064 时代已存在的重复 batch 数据」不拆解（非资金不变式；约束只防新混批），限制已写入迁移注释与本记录，供运营知悉。

## 关闭判定

1. A-005 全部 findings 处置完毕：F-001 required（A-006/A-007 闭合）+ F-002～F-005 recommended（本条目 `fixed`）。**open required = 0**，recommended 亦清零（原 A-005 台账 backlog 移除）。
2. 验证证据：全后端回归 exit 0（composition/handler/store/modules/wallet）+ PG e2e PASS（docker postgres:15）+ renderer vitest 290 PASS。
3. 状态：不改 `status`/`progress`；Root 维持 `done` · 4/4。
4. 复审建议：如需第四只眼，可再调 `/audit` 复审本条目（非强制——本轮 recommended 均非门禁项）。

## 声明

本条目为编排器 self 响应，不伪装 `source: independent`；A-005/A-007 原文保留原文件。
