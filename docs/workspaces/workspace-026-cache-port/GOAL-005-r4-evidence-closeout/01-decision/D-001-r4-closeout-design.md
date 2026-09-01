---
doc_type: goal-decision
id: D-001-r4-closeout-design
parent: GOAL-005-r4-evidence-closeout
date: 2026-09-01
status: accepted
version: 0.1.0
---

# D-001 · R4 关门设计（2026-09-01）

## 上下文

R4 为 VP-026 收官阶段（Root 3/4 · GOAL-002/003/004 done · I-026-001～004 全部 verified）。设计依 workspace-020/025 R4 先例与项目级独立审计执行路径。

## 冻结内容

| 项 | 决定 |
|----|------|
| 证据矩阵 | 判据 #1～#8 逐条 ↔ 阶段产物/测试/审计证据（`attachments/r4-evidence-matrix.md`）；每项 verified 或说明 |
| 越界核账 | `git diff --name-only 54fb57e7..HEAD`（VP-026 规划 commit → R3 关门 commit）全量路径分类；红线 = Charter / `go.mod`+`go.sum` / Profile 默认集 / Manifest 装配 / 迁移台账，零触碰 |
| Root 关门双审 | **cross**：A-001 self（Root 03-audit）+ A-002 本地 grok build（grok-4.6 · high · headless）independent（Root 03-audit；原始输出附录件）；A-003 编排器合并响应 |
| VRev-061 | /vision 层关门审视（self）：VP-026 八条退出判据 + 边界 + 台账一致性；落 `docs/vision/reviews/VRev-061-vp026-cache-port-close-out.md` + `reviews.md` 索引 |
| VP-026 closed 呈报 | 双审 pass + VRev-061 pass → 呈报**用户书面确认关门**（P-004 最终裁决点 · workspace-025 D-003 先例）→ VP-026 `active → closed`（关门记录表 + 修订史 + roadmap/workspaces/revisions 同步） |
| 工作区结项 | Root `done` 4/4（goal-tree 收官）· workspace.md 结项记录（root_goal 状态 + 结项说明） |
| 提交 | 关门 checkpoint 单次提交（owned paths only） |

## 未选方案

| 项 | 未选 | 理由 |
|----|------|------|
| 仅 self 关门（无 independent） | 未选 | Root 关门 = 实证门禁，项目惯例要求 grok build independent 腿（Root D-001「实证门禁可按需 independent」；workspace-020/025 先例双审） |
| 分多次提交 | 未选 | 关门为单一台账事务；单 checkpoint 提交保证可追溯（先例一致） |

## 影响

- C1～C3 按本合同执行；关门不重开历史阶段；VP closed 不改变 Charter / primary workspace。