---
doc_type: goal-decision
id: D-006-r3-c2-a008-fixed
parent: GOAL-004-r3-session-operator-console
date: 2026-09-04
source: user
status: done
version: 0.1.0
---

# D-006 · R3 C2 A-008 required finding 闭合路径裁决

## 用户裁决

用户通过裁决工具选择：对 A-008 的 F-001（PostgreSQL 唯一冲突安全幂等）与 F-002（主体映射事务顺序）采用 **fixed** 路径，并按修正后重新进行 self + Grok independent 复审。

## 闭合边界

- F-001 以 D-005 的书面补全闭合：`INSERT ... ON CONFLICT DO NOTHING` + `RowsAffected()`；0 行是重复成功，1 行才 upsert 会话；不得在 PostgreSQL 已冲突的 Tx 内继续语句，并必须补 SQLite/PG 重复路径证据。
- F-002 以 D-005 的书面补全闭合：既有 `GetOrCreateSubject` 在唯一收据前执行并保持独立 `Store.Run`；主体映射失败时不铸造 inbox，重试不得因重复收据短路而吞掉分发。
- 本裁决不接受 residual、不作 user-overruled；在 A-008 required finding 经响应并独立复审确认前，C2 仍不得进入生产代码实施。

