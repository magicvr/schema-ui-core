---
id: D-006-reopen-after-visual-fidelity-audit
title: 废止 D-005 关门 · 回退 Root/工作区完成状态
date: 2026-08-09
status: accepted
parent: GOAL-001-design-system-and-ui-experience
supersedes: D-005-root-closeout-user-confirmed
related_audit: A-006
---

# D-006 · 重开决策（视觉 fidelity 复审后）

## 决策

1. **废止** [D-005](./D-005-root-closeout-user-confirmed.md) 将 Root / 工作区标为 `done` 的效力（D-005 文件保留为历史记录，`status` 改为 `superseded`）。
2. Root `GOAL-001-design-system-and-ui-experience`：`status` → **`active`**。
3. 工作区 `workspace-006-design-system-and-ui-experience`：`status` → **`active`**。
4. Root 阶段检查点诚实回写：
   - **保持勾选**：S1（Token/主题/FOUC/F-002）、S4（Skeleton / async-state）。
   - **取消勾选**：S2、S3、S5（S5 过程关门前提已失效；fork 示例仍记在 GOAL-005）。
   - 派生 `progress: **2/5**`。
5. 子目标 GOAL-003（S2+S3）：`status` → **`active`**；其先前过窄 C1/C2 勾选取消；后续成功标准须对齐 D-004，不得再以「仅 chart 色 / 仅移动抽屉」代表 Root S2/S3 完成。
6. GOAL-002 / GOAL-004 / GOAL-005 保持 `done`（局部交付真实）；**禁止**仅凭子目标 done 再次推导 Root done。

## 依据

- A-006（fail）：用户可观察 UI 与 D-004 / Stitch 参考脱节；S2/S3 完成被偷换。
- 用户本会话书面要求：落盘审计意见并回退工作区完成状态。
- AGENTS / P-002：实施事实须可验证；未合法闭合 required 不得 `done`（回退后 F-VUI-001/002 为开放 required）。

## 未选方案

| 方案 | 未选原因 |
|------|----------|
| 仅改文档措辞、保持 `done` | 与用户可观察事实冲突；继续过程不诚实 |
| 回退全部 S1–S5 与 GOAL-002/004/005 | 过度否定真实基建与状态面交付；A-006 未否决 S1/S4 技术证据 |
| 静默改 status 不写审计 | 违反 P-003 台账与可追溯 |

## 闭合

| 项 | 结果 |
|----|------|
| F-VUI-003（过早 done） | **fixed**（本决策 + E-005 状态回写） |
| F-VUI-001 / F-VUI-002 | 仍 **open required** |
| 再次 Root `done` | 须 S2/S3 按 D-004 交付 + 开放 required 闭合 + 用户书面确认 |
