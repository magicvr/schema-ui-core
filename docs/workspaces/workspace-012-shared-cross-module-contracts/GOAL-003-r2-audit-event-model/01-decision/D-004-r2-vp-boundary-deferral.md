---
id: D-004-r2-vp-boundary-deferral
doc: decision-entry
status: accepted
created: 2026-08-18
updated: 2026-08-18
parent: GOAL-003-r2-audit-event-model
version: 0.1.0
---

# D-004 · R2 VP-012 边界延期

为保持 R2 的可验证交付边界，本轮只关闭 D-003 指定的结构化 detail、递归脱敏、correlation，以及 auth/settings/users 三类真实 mutation。VP-012 方向中的 session/effective actor、保留/归档触发，以及 `users_state`、MFA、wallet 等未列入 D-003 的写路径，不作为 R2 的完成标准，保留给后续 R3-R6 波次或单独目标承载。

这不是已交付或已验证的声明；后续波次仍需重新登记信息项、冻结范围并提供对应实现与审计证据。
