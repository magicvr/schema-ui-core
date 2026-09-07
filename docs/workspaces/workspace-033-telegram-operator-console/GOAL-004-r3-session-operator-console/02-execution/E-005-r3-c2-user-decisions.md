---
doc_type: goal-execution
id: E-005-r3-c2-user-decisions
parent: GOAL-004-r3-session-operator-console
date: 2026-09-04
status: done
version: 0.1.0
---

# E-005 · R3 C2 用户方案裁决事实

## 已发生事实

- 用户通过裁决工具选择 C2 的三项方案：双表最小持久化对象面、规范化字段不留 raw JSON、兼容现有 handler 失败语义。
- D-004 已原样记录上述选择，并明确承接 D-002 的 `chat_id` 会话边界与 D-003 的持久化先于确认/offset 合同。
- C2 仍无迁移、repository、会话/消息持久化或业务代码；schema 列级参数、共同入站接线和测试矩阵待 self/independent 审视后实施。

## 门禁事实

- C1 已关闭，R3 保持 `active · 1/4`；C2 进入合同审视阶段，未宣称完成。
- 该次裁决未选择出站表、raw JSON 或 inbound dispatch 状态机；若实施中需要越过这些边界，须重新进入 `/govern` 裁决。
