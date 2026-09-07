---
doc_type: goal-decision
id: D-001-r4-closeout-design
parent: GOAL-005-r4-evidence-closeout
date: 2026-09-01
status: accepted
version: 0.1.0
---

# D-001 · R4 关门设计

## 决策

| 项 | 决定 |
|----|------|
| 证据矩阵 | 七条方向级判据逐条映射：判据 → 证据文件（阶段目标五件套 / 代码 / 短文）→ 验证命令。分母 = VP-027 退出判据 1～7 |
| 越界核账区间 | 全波次 `889a80bb^..HEAD`（激活 → R1 → R2 → R3 四连 commit）；禁区 = go.mod / go.sum / kernel/profile.go / internal/manifest / config.default.yaml / docs/vision/charter.md / redis 依赖 |
| 关门流程 | ① C1 矩阵 + 核账 + 最终回归复跑 → ② C2 Root 双审（A-001 self + A-002 grok build independent）→ VRev-063（vision 层关门就绪 · self）→ ③ **P-004 用户书面确认** VP-027 `active → closed`（v0.3.0）→ vision 台账原子同步（roadmap 行 27 / workspaces.md 027 行 / reviews.md / revisions VR-057）→ Root `done` 4/4 · 工作区结项 · 最终 checkpoint |
| 审计模式 | cross（self + grok build grok-4.6 · high independent · 项目级默认执行路径） |

## 未选方案

| 项 | 未选 | 理由 |
|----|------|------|
| 静默关门（不询问） | 未选 | P-003/P-004：VP 关门必须用户书面确认 + 工作区证据链接；子目标静默授权不覆盖 VP/组合层关门 |

## 影响

- C3 为用户裁决点（P-004），确认前不执行任何 VP/Root 状态变更。