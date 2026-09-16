---
id: D-001-r5-composition-closeout-contract
doc: decision
status: accepted
goal_id: GOAL-006-r5-composition-acceptance
created: 2026-09-17
updated: 2026-09-17
parent: GOAL-006-r5-composition-acceptance
version: 1.0.0
---

# D-001 · R5 组合验收与关门合同

R5 只承载 VP-037 首波的组合验收与关门准备，不新增产品功能、不重写 R1～R4 已接受合同，也不把建议性覆盖提升为 Root/VP 的硬门禁。

- C1 核对 R1～R4 的阶段状态、五件套、检查点、信息项与 required finding 台账；前序阶段审计意见必须可回指，不能用 Root 或 Vision Review 记录替代 Goal 审计。
- C2 核对 Saved Views、dirty-state、统一反馈的已确认范围，以及实体全文检索、批量结果中心、组织/数据权限、新业务域、Redis/MQ/多实例、跨用户协作等非目标边界；发现 residual 或冲突时按 P-004 记录，不静默放行。
- C3 复跑最终测试、类型检查与 diff check，记录 Git checkpoint，排除用户工作树中的 `.claude/settings.local.json` 等非本目标文件。
- C4 先完成 R5 self + Grok independent 审计并闭合所有 required finding，再向用户请求 Root/VP 关门书面确认；确认前 Root/VP/workspace 保持 `active`，确认后才允许投影 `done`/`closed`。
- R5-I-005 的 Host/resource 直接对照与协作需求保持 non-blocking deferred/recommended；只有真实需求或支持场景触发时，才由 `/vision` 另行复核。

本决定不引入新的 API、持久化层、跨用户权限、全局状态管理或新愿景计划。
