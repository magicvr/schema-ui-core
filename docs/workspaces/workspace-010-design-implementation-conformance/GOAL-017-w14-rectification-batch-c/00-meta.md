---
id: GOAL-017-w14-rectification-batch-c
title: W14 整改批 C · 调试痕迹清理（F-08 移除 pageId/route 技术框 / F-09 反馈文案本地化与去错误码前缀 / F-10 Schema 加载失败友好化）
status: done
parent: GOAL-015-w14-user-perspective-review
created: 2026-08-17
updated: 2026-08-17
version: 0.2.0
progress: 4/4
---

# GOAL-017 · W14 整改批 C（F-08～F-10 调试痕迹清理）

[GOAL-015](../GOAL-015-w14-user-perspective-review/00-meta.md)（W14）的**下级整改子目标（批 C）**：承接用户书面裁决（D-003）——F-01～F-14 全部 in-scope、分批实施；批 A 已关门（GOAL-016）。本子目标 = **批 C（真实用户可见的调试痕迹与本地化，F-08～F-10）**。

## 当前边界

- **范围（本波实施）**：
  - F-08：每页标题右侧无条件显示 `pageId` + `route` 技术信息框 → **直接移除**（D-003 冻结）。
  - F-09：反馈 toast 前缀技术错误码 + 多处英文硬编码消息 → 改为 i18n 友好文案（保留错误码在可访问/调试位置或去掉前缀）。
  - F-10：页面 schema 加载失败时把错误码当主标题 + 显示原始 URL/issue.path → 改为人类可读解释与恢复引导。
- **非范围**：批 A（已完成）、批 D（F-11～F-14）、批 B（F-05～F-07）；不改 Profile 默认集 / 模块矩阵 / Manifest 装配语义；不做视觉重设计。

## 成功标准与路线图（P-001）

- [x] **S1 · 方案冻结**：F-08～F-10 具体修改点与 i18n 键清单
- [x] **S2 · 实施**：前端/文案修改
- [x] **S3 · 测试与回归**：vitest/tsc + Web 全量（如涉 Go 则 Go 全量）
- [x] **S4 · 自审与关门**：审计 + 台账同步 + goal-tree/workspace 同步

progress: 由四个等权检查点派生（S1～S4）；当前 **4/4**（S1～S4 完成，2026-08-17 关门）。

## 审计策略

| 阶段 / 项 | 默认模式 | 说明 |
|-----------|----------|------|
| S1 冻结 | self | 方案落盘 + 证据核对（来自 D-001 台账） |
| S2 实施 | self | 呈现/文案类可逆修改 |
| S4 关门 | self | 常规关门自审 |

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | non-blocking | F-09 错误码保留位置（去掉 toast 前缀后是否仍要可访问调试信息） | S2 F-09 | S1 | as-built + 方案 | **closed** | — | D-001：`data-feedback-code` + title |
| I-002 | non-blocking | F-10 失败页恢复动作（reload / 返回 / 支持链接） | S2 F-10 | S1 | as-built + 方案 | **closed** | — | D-001：重新加载按钮 + 技术详情折叠 |

## 父目标

- [GOAL-015-w14-user-perspective-review](../GOAL-015-w14-user-perspective-review/00-meta.md)

## 台账布局

- `01-decision/`：D-NNN；`02-execution/`：E-NNN；`03-audit/`：A-NNN。
- 跨区引用用 Q2 路径（workspace-protocol §2.6）。
