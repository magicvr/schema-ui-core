---
status: done
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-002-r1-contract-freeze
version: 0.1.0
---

# D-002 · R1 关门（用户确认）

## 决策

用户 2026-08-29 书面确认（P-004）：**「确认冻结面，关门 R1，开 R2」**。

1. 契约冻结面清单 **v0.1.0 → v1.0.0 生效**（`attachments/freeze-face-v0.1.0.md`，标题/状态同步）；semver/breaking 流程与 changelog 模板随清单生效（semver-breaking-policy §5 门禁自即日起适用）。
2. R1 检查点关闭，Root progress 0/5 → 1/5。
3. 未决项明确交接 R2：A-001 F-001（C 层泄漏收敛裁定）、A-001 F-002（B 层全量符号回填）、I-001（collecting → 随 R2 收尾）。
4. 后续阶段 R2/R3：`KernelAPIVersion` 主号不变（当前 2.0.0），打包验证不得引入契约破坏。