---
id: GOAL-026-w16-rectification-batch-b
title: W16 整改批 B · 核心资产与数据交互（F02 文件预览复制 / F03 导入模板与错误定位 / F04 金额格式化与调账警示）
status: done
parent: GOAL-024-w16-user-perspective-improvements
created: 2026-08-17
updated: 2026-08-17
version: 0.2.0
progress: 4/4
---

# GOAL-026 · W16 整改批 B（核心资产与数据交互）

[GOAL-024](../GOAL-024-w16-user-perspective-improvements/00-meta.md)（W16）的**下级整改子目标（批 B）**：承接 D-003 分批规划，批 A 已关门后渐进添加。本子目标实施 **W16-F02 文件库在线预览与一键复制直链**、**W16-F03 数据导入 CSV 模板与逐行错误定位**、**W16-F04 钱包金额分转元展示与调账警示**。

## 当前边界

- **范围**：F02 文件预览/复制链接；F03 导入模板下载/fieldErrors 明细；F04 金额格式化/调账二次确认。
- **非范围**：批 A（已 done）、批 C（F05/F06/F09/F10）不在此实施；不改 Profile 默认集 / 模块矩阵 / Manifest 装配语义。

## 成功标准与路线图（P-001）

- [x] **S1 · 方案冻结**：F02/F03/F04 设计（端点/schema/前端交互）+ 信息项登记（D-001）。
- [x] **S2 · 实施**：F02/F03/F04 代码与 schema/前端接线。
- [x] **S3 · 测试与回归**：Go 全量 + Web vitest/tsc + 相关 e2e。
- [x] **S4 · 自审与关门**：审计 + 台账同步 + goal-tree/workspace 同步。

progress: 由四个等权检查点派生（S1～S4）；当前 **4/4**。

## 审计策略

| 阶段 / 项 | 默认模式 | 说明 |
|-----------|----------|------|
| S1 冻结 | self | 方案落盘 + 证据核对 |
| F03 导入契约调整 | self（涉 API 响应兼容时升级 independent） | 错误明细结构可能影响既有导入调用方 |
| S4 关门 | self | 常规关门自审 |

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|
| I-001 | required | F03 导入错误响应升级为 `fieldErrors` 时如何保持既有 `errors` 兼容 | S1 方案冻结 | S1 | 检查导入 handler 与前端调用方 | **closed** | D-001 §3：新增 `fieldErrors` + 保留 `errors` |
| I-002 | non-blocking | F02 文件库预览 URL 是否可直接用于新标签页（鉴权/Content-Disposition） | S1 方案冻结 | S1 | 检查 filelibrary 路由 | **closed** | D-001 §2：使用现有 `downloadUrl` 端点 |

## 父目标

- [GOAL-024-w16-user-perspective-improvements](../GOAL-024-w16-user-perspective-improvements/00-meta.md)

## 台账布局

- `01-decision/`：D-NNN；`02-execution/`：E-NNN；`03-audit/`：A-NNN。
