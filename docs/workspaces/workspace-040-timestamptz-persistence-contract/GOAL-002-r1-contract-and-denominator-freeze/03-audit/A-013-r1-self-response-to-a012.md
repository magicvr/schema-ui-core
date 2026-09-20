---
id: A-013-r1-self-response-to-a012
doc_type: goal-audit-entry
source: self
auditor: /govern
date: 2026-09-20
scope: GOAL-002-r1-contract-and-denominator-freeze · response to A-012 / freeze package alignment
verdict: conditional
open_required: 6
status: recorded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-001-timestamptz-persistence-contract
version: 0.1.0
---

# A-013 · R1 self response to A-012

## 响应

### F-I-014 → fixed（freeze-package alignment）

已对齐 C2 冻结载体：

- `r1-c2-column-contract-draft-v0.1.md` 已删除 `NULL/backfill` 替代项，固定 D-008 的 config 0→NULL；补写新值/迁移统一向零截断到微秒；固定 D-009 的 filelibrary/configpkg include；固定 `pg_dump -F c` + `pg_restore`；固定 D-007 最小 kernel Port/internal Service 边界。
- `r1-c2-c3-guardrails-v0.1.md` 头注不再把已决方向写成待 P-004 选择。
- child D-005 仍为 proposed 草案，尚未被静默升格为 accepted。

## 仍开放 required

- F-I-002：逐列 codec/USING/rebuild、非法值、排序与 round-trip。
- F-I-003：90 列 old→new→read/write mapping、voucher <=0、login/task predicates。
- F-I-004：Port 方法、metadata、前后 backup、restore-to-new-db 与 rollback evidence。
- F-I-005：72-entry append-only test/leftover list/runner ModuleID。
- F-I-006：CHECK/index/WHERE/ORDER old/new rebuild list。
- F-I-010：planning coverage 已经独立认可，但 formatter/parser/fixture 实施仍未完成；C2/R3 门禁保持开放。

## 放行

A-012 的新 required F-I-014 已响应；C2/C3 仍不可冻结，R2 不启动。

## 追加响应（2026-09-20 · D-010/D-011）

- Backup Port 方法已由用户冻结为仅 `CreateRecoveryPoint`，成功返回含最低验证后置条件；Verify/RestoreTo/provider orchestration 留在 internal（Root D-010 / child D-009）。
- `schema_migrations.applied_at` conversion owner 已冻结为 `core.persistence`；Store runner 仍负责写入（Root D-011 / child D-010）。
- C2 column-contract/guardrails 已同步这些 owner/Port 边界；具体 Port metadata、restore procedure、逐列 codec/predicate/checksum 测试仍开放。

下一步继续补逐列 mapping、predicate/check、catalog test owner、Backup Port metadata/restore script，再进行 grok independent 复审。
