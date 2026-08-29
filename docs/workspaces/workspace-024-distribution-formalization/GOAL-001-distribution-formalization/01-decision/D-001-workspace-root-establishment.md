---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-001-distribution-formalization
version: 0.1.0
---

# D-001 · 工作区 Root 建立（2026-08-29）

1. **绑定**：lead workspace = `workspace-024-distribution-formalization`（`vision_role: delivery`，单一 lead）；`root_goal` = `GOAL-001-distribution-formalization`（`parent: null`）；`primary_plan` = `VP-024-distribution-formalization`（active v0.2.0）；不改变 Charter `primary_workspace`。
2. **freshness 三字段（V-F087 · VRev-052 recommended → fixed）**：
   - `consumer_vp` = VP-024（消费候选基线：HEAD `c9122478` · `apps/api/v0.3.0` · 六包 GH Packages）
   - `last_freshness_review_at` = 2026-08-29（VRev-052：`041744b3` → `c9122478` PASS）
   - `next_freshness_review_trigger` = 下一次 VP 激活/消费前主动核对；协议 pin / 依赖锁 / 迁移台账 / Profile 默认集 / provenance 任一变更即触发重验证（VP-008 `go` 消费有效性规则）
3. **审计模式**：阶段关门 default self；R2（公开发布 · release 门禁）与 R4（fork 对照 · 实验数据门禁）与 Root 关门 = independent（grok build 先例）。
4. **信息门禁（P-005）**：I-024-001（R2 前置 · npmjs 授权）required / I-024-002（R3 · CI 环境）required / I-024-003（R4 · fork 对照样本）required / I-024-004（R1 · serve 接线）non-blocking——登记于 Root `00-meta`。
5. **外部动作边界**：npmjs 正式 scope/凭据（判据 #2）不随激活授权；R2 到达前由用户书面授权或裁决降级（同 I-003 / VP-023 先例；V-F088 → fixed）。