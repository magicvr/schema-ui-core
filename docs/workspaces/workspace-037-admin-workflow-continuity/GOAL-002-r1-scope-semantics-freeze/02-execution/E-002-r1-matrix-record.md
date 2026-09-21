---
id: E-002-r1-matrix-record
doc: execution
goal_id: GOAL-002-r1-scope-semantics-freeze
status: recorded
created: 2026-09-17
updated: 2026-09-17
parent: GOAL-002-r1-scope-semantics-freeze
version: 0.2.0
---

# E-002 · R1 分母与状态/反馈矩阵

## 已发生事实

- 通过 `searchable-profile-matrix.test.ts`、manifest、Schema 和 renderer 代码完成 profile/page/action 基线核对。
- 形成 `attachments/r1-denominator-matrix.json`：4 个 profile 的注册/可发现/隐藏路由、稳定 ID 合同、24 个列表表面、查询/权限基线。
- 形成 `attachments/r1-form-matrix.json`：递归盘点得到 58 个 form 节点，包含 modal/action content form。
- 形成 `attachments/r1-state-feedback-matrix.md`：记录当前 dirty-state、导航、反馈、错误合同与待冻结语义。

## 当前边界

上述附件记录的是当前代码事实与明确标注的候选语义；没有把 Saved View 持久化、dirty guard 或统一反馈实现写成已完成。此前记录的“等待用户在 D-002 的方案 A/B 中裁决”已由 2026-09-17 的 D-003 取舍记录承接；I-037-002～004 的 R1 信息冻结现为 `verified`，R2～R4 仍分别等待实现证据。
