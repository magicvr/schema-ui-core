---
id: GOAL-018-w14-rectification-batch-d
title: W14 整改批 D · 表单与无障碍（F-11 必填标记 / F-12 确认对话框焦点 / F-13 桌面表格键盘选中 / F-14 小缺口）
status: done
parent: GOAL-015-w14-user-perspective-review
created: 2026-08-17
updated: 2026-08-17
version: 0.2.0
progress: 4/4
---

# GOAL-018 · W14 整改批 D（F-11～F-14 表单与无障碍）

[GOAL-015](../GOAL-015-w14-user-perspective-review/00-meta.md)（W14）的**下级整改子目标（批 D）**：承接用户书面裁决（D-003）——F-01～F-14 全部 in-scope、分批实施；批 A（GOAL-016）、批 C（GOAL-017）已关门。本子目标 = **批 D（表单与无障碍，F-11～F-14）**。

## 当前边界

- **范围（本波实施）**：
  - F-11：必填字段标记（`*` / `aria-required`）与提交前可见提示。
  - F-12：确认对话框焦点圈 / ESC / 初始焦点（与 ModalHost 一致）。
  - F-13：桌面表格行键盘可选（tabIndex/role/onKeyDown；移动卡片已支持 Enter/Space）。
  - F-14：移动卡片列丢弃 / 单选空值误导 / 列排序无法清除 / Tabs 方向键与 aria-controls / 禁用按钮原因 / 字段错误 aria-describedby / <sm 语言切换 / 通知空收件箱语义文案。
- **非范围**：批 A、批 C（已完成）、批 B（F-05～F-07）；不改 Profile 默认集 / 模块矩阵 / Manifest 装配语义；不做视觉重设计。

## 成功标准与路线图（P-001）

- [x] **S1 · 方案冻结**：F-11～F-14 具体修改点与可访问性契约
- [x] **S2 · 实施**：前端组件/样式/键盘交互修改
- [x] **S3 · 测试与回归**：vitest/tsc + Web 全量（如涉 Go 则 Go 全量）
- [x] **S4 · 自审与关门**：审计 + 台账同步 + goal-tree/workspace 同步

progress: 由四个等权检查点派生（S1～S4）；当前 **4/4**（S1～S4 完成，2026-08-17 关门）。

## 审计策略

| 阶段 / 项 | 默认模式 | 说明 |
|-----------|----------|------|
| S1 冻结 | self | 方案落盘 + 证据核对（来自 D-001 台账） |
| S2 实施 | self | 可访问性/表单类可逆修改 |
| S4 关门 | self | 常规关门自审 |

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | non-blocking | F-12 确认对话框是否复用 ModalHost 的焦点陷阱实现 | S2 F-12 | S1 | as-built + 方案 | open | — | 待 S1 |
| I-002 | non-blocking | F-14 列排序清除交互（再次点击升/降后进入无排序态？） | S2 F-14 | S1 | as-built + 方案 | open | — | 待 S1 |

## 父目标

- [GOAL-015-w14-user-perspective-review](../GOAL-015-w14-user-perspective-review/00-meta.md)

## 台账布局

- `01-decision/`：D-NNN；`02-execution/`：E-NNN；`03-audit/`：A-NNN。
- 跨区引用用 Q2 路径（workspace-protocol §2.6）。
