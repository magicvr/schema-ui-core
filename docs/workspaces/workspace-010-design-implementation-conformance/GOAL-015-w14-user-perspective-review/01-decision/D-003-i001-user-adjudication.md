---
id: GOAL-015-w14-user-perspective-review
doc: decision
status: active
parent: GOAL-001-design-implementation-conformance
created: 2026-08-17
updated: 2026-08-17
version: 0.1.0
---

# D-003 · W14 I-001 用户书面裁决：F-01～F-14 全部 in-scope（分批整改）

## 背景

用户裁决撤销前次未经裁决的关门（E-003/A-004）后，I-001（required）重新开放为 GOAL-015 本波关门门禁：须用户对 F-01～F-14 的 in-scope / defer / 优先级作出书面裁决。2026-08-17 用户经 GUI 问询作出裁决（见本文件）。

## 裁决（用户书面，2026-08-17）

1. **I-001 · F-01～F-14 in-scope/优先级**：**分批全部 in-scope**，按 D-001 建议实施顺序执行：
   - **批次 A**：F-01～F-04（功能面补全）→
   - **批次 C**：F-08～F-10（调试痕迹清理）→
   - **批次 D**：F-11～F-14（表单与无障碍）→
   - **批次 B**：F-05～F-07（一致性硬化）。
   - 实施方式：另起整改波次（workspace-010 子目标）按批次执行。
2. **F-01 handler 目录暴露方式**：**新增端点**列出可用 handler（如 `GET /api/scheduled-tasks/handlers`），前端字段动态加载选项。
3. **F-04 通知本地化方案**：**存 messageKey**（通知存 key，前端按语言渲染文案）。
4. **F-08 调试框处理**：**直接移除** pageId + route 技术信息框。

## 影响

- I-001 角度：本波关门所需用户裁决已取得（P-004 ✓）。GOAL-015 的 S4 可据此合法性关门（S1～S3 已完成，审计 A-001/A-002 pass、A-003 步骤事实、整改 deferred 决策由本裁决覆盖）。
- 整改实施**不在 GOAL-015 范围**（GOAL-015 交付 = 审视 + 台账 + 审计 + 同步）；按裁决另起整改波次。批次 A 为首个可执行整改波次。
- 三个方案选择（F-01/F-04/F-08）作为整改波次输入，冻结。
