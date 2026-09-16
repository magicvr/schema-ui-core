---
id: D-005-feedback-recovery-semantics-frozen
doc: decision
goal_id: GOAL-002-r1-scope-semantics-freeze
status: accepted
created: 2026-09-17
updated: 2026-09-17
parent: GOAL-002-r1-scope-semantics-freeze
version: 0.1.0
---

# D-005 · 反馈与恢复语义冻结

本条沿用现有 `readResourceApiError`、`FeedbackRegion`、`DataTable` 与 Host failure 合同，将 R1 基线固定为以下实现边界：

- 成功操作使用短暂 `role=status` feedback；错误使用持久 `role=alert` feedback；同一操作不重复发出成功事件。
- 幂等读取（列表、页面 Schema、Saved View storage read）允许显式 retry；retry 只发起新的读取，不重复写入。
- 写入失败不自动 retry；表单保留当前值、字段错误和 dirty 状态，用户再次提交才产生新的写请求。
- 用户文案优先使用 catalog `messageKey/params`；诊断 code/correlation 可用于定位，但不得回显 token、密码或原始敏感载荷。
- 401/重新认证继续由 Host failure/认证层处理；403 继续 fail closed；维护、不可用、离线和超时呈现可恢复的读取失败，不伪装成业务成功。
- Saved View 读取到 malformed、未知、越权或已不匹配当前 Schema 的记录时丢弃该记录并通过统一反馈提示；当前 query 不被污染。
- Toast、列表错误、retry 控件和确认路径必须可键盘到达并维持既有 `status/alert` 无障碍角色。

## 证据与边界

当前错误合同和待实现差距见 `attachments/r1-state-feedback-matrix.md`；R4 负责把已冻结语义补齐到跨页面回归，不把本条当作 R4 已完成事实。
