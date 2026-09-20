---
id: A-011-r1-self-response-to-a010
doc_type: goal-audit-entry
source: self
auditor: /govern
date: 2026-09-20
scope: GOAL-002-r1-contract-and-denominator-freeze · response to A-010 / C2/C3 guardrails
verdict: conditional
open_required: 6
status: recorded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-001-timestamptz-persistence-contract
version: 0.1.0
---

# A-011 · R1 self response to A-010

## 响应

### F-I-013 → fixed（recommended）

Child D-005 guardrails 已修正：不再把 child D-006 与 Root D-006 的 Backup SPI/Service 口径混用，明确 Root D-007 的最小 kernel Backup/RecoveryPoint Port 与 internal orchestration/provider 边界；wire 交叉引用补 D-009。

### F-I-010 → fixed（planning coverage，等待 independent 复审）

公共 wire inventory 已补齐 A-010 点名的 `account_self.go:391`、service credential audit detail 与 fixtures、`dictionary.go:329/337`、`scheduledtasks.go:506/513/516`；filelibrary ModTime 与 configpkg metadata 已由用户 D-009 明确纳入；API input 与 Web display parser 的边界按 D-005 记录；RFC3339Nano 输出禁令写入 guardrails。当前 formatter/fixtures 仍未修改，故仅关闭 planning coverage，不宣称实现完成。

### C2/C3 用户裁决响应

- 精度：统一向零截断到微秒（D-008）。
- mail/telegram config D0：legacy 0 → NULL，移除 default 0（D-008）。
- PG provider：固定 `pg_dump -F c` + `pg_restore`（D-008）。
- Backup Port：最小 kernel port、Service/providers internal（D-007）。

## 仍开放 required

- F-I-002：逐列 codec/DDL/非法值/排序/round-trip。
- F-I-003：90 列 old→new→read/write mapping、voucher <=0、predicates。
- F-I-004：Port methods、metadata、before/after backup、restore-to-new-db、rollback evidence。
- F-I-005：72-entry append-only test/leftover list/runner ModuleID。
- F-I-006：CHECK/index/WHERE/ORDER old/new rebuild list。
- F-I-010：独立复审确认 planning coverage；实施仍未开始。

C2/C3 与 R2 继续阻断；下一步是 grok independent 复审本轮修正，然后继续补具体 C2/C3 证据。
