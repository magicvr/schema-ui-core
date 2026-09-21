---
id: D-009-open-r4-unified-feedback
doc: decision
status: accepted
goal_id: GOAL-001-admin-workflow-continuity
created: 2026-09-17
updated: 2026-09-17
parent: GOAL-001-admin-workflow-continuity
version: 1.0.0
---

# D-009 · 开设 R4 统一反馈与恢复

接受按 VP-037 路线开设 `GOAL-005-r4-unified-feedback-recovery`，范围承接 R1 D-005：统一成功/错误反馈的用户文案与无障碍角色，明确幂等读取的显式重试边界，保持写入失败草稿与 dirty 状态，不自动重复提交，并区分普通资源反馈与 Host 层维护/不可用/认证终态。

R4 复用 `readResourceApiError`、现有 `FeedbackRegion`、`DataTable` retry 和 Host failure 合同，不新增 API、错误 envelope、第二套全局状态管理或 R5 组合验收。Root 保持 `active · 3/5`，R4 初始为 `active · 0/4`。
