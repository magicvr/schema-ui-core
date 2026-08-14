---
id: A-003
goal: GOAL-012-r3-s12-recycle-bin
source: independent
date: 2026-08-14
scope: S5 关门 · 数据门禁（admin.recycle-bin vs D-002 冻结方案）
verdict: fail
auditor: grok-build
audit_type: close-out
status: recorded
parent: GOAL-012-r3-s12-recycle-bin
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# A-003 · independent 数据审计（S-12 实现 · 首轮）

## 结论

**verdict: fail**。3 个 required findings（F-001 批量快照 ID 碰撞、F-002 恢复冲突真服务 HTTP 测试缺失、F-003 冲突 wire 格式 + web i18n）+ 5 个 recommended（F-004 批量清除缺失、F-005 快照 ID 格式、F-006 错误码区分、F-007 文档不一致、F-008 测试强化）。原始意见：`attachments/grok-audit-s12-closeout.txt`。

## Findings

- F-001（required, high）：批量删除用同一 now 生成快照 ID（hexID(UnixNano)）→ 同秒多条 PK 冲突，Record 失败仅 slog → 批量无快照。
- F-002（required, med）：恢复冲突（键被占用 + 未恢复）无真服务 HTTP 测试（handler 用 fake）。
- F-003（required, med）：冲突走 writeError（无 messageKey）；web 缺 error.recycle* 键。
- F-004（recommended）：D-002 §4 页面批量清除未实现。
- F-005（recommended）：快照 ID 非 D-002 的 "recycle-" + 16 hex（实际 38 hex）。
- F-006（recommended）：已恢复与键冲突共用 RECYCLE_RESTORE_CONFLICT 语义混淆（已分离）。
- F-007（recommended）：01-decision.md I-001 open vs meta closed；S3/S4 勾选缺失。
- F-008（recommended）：Trash nil 测试只断言 204 未断言无快照；dict-entries 恢复无单测。

## 响应

- F-001/F-005：Record 用 crypto/rand 8 字节（"recycle-" + 16 hex）；TestRecordDistinctIDsSameSecond。
- F-002：provider_test TestRecycleRealServiceRestoreConflictHTTP（真服务：冲突 409 + 快照保留 + 解决后恢复 200 + messageKey 断言）。
- F-003：DomainError 路径改 writeLocalizedError；web i18n +2 键。
- F-004：POST /api/recycle-bin/purge-all + schema toolbar purgeAll + store PurgeAllUnrestored。
- F-006：service 先查 RestoredAt → ErrItemAlreadyRestored；键冲突才 RECYCLE_RESTORE_CONFLICT。
- F-007：01-decision I-001 → closed；meta S3/S4 勾选。
- F-008：nil 测试断言 recycle total==0；TestRestoreDictEntryRoundTrip。

闭合状态见 A-004 复审（grok 再审 + 全量回归）。
