---
id: GOAL-027-w16-rectification-batch-c
title: W16 整改批 C · 系统运维与通用外观（F05 Cron 预览 / F06 监控自动刷新 / F09 字典 Badge / F10 页脚版权）
status: done
parent: GOAL-024-w16-user-perspective-improvements
created: 2026-08-17
updated: 2026-08-17
version: 0.2.0
progress: 4/4
---

# GOAL-027 · W16 整改批 C（系统运维与通用外观）

[GOAL-024](../GOAL-024-w16-user-perspective-improvements/00-meta.md)（W16）的**下级整改子目标（批 C）**：承接 D-003 分批规划，批 A/B 已关门后渐进添加。本子目标实施 **W16-F05 Cron 自然语言与未来运行预览**、**W16-F06 系统监控自动轮询刷新**、**W16-F09 数据字典 Badge 颜色/标签风格**、**W16-F10 系统设置页脚版权与备案号**。

## 当前边界

- **范围**：F05 Cron 预览；F06 监控自动刷新；F09 字典 Badge；F10 页脚版权。
- **非范围**：批 A/B（已 done）；不改 Profile 默认集 / 模块矩阵 / Manifest 装配语义。

## 成功标准与路线图（P-001）

- [x] **S1 · 方案冻结**：F05/F06/F09/F10 设计（端点/schema/前端交互）+ 信息项登记（D-001）。
- [x] **S2 · 实施**：F05/F06/F09/F10 代码与 schema/前端接线。
- [x] **S3 · 测试与回归**：Go 全量 + Web vitest/tsc + 相关 e2e。
- [x] **S4 · 自审与关门**：审计 + 台账同步 + goal-tree/workspace 同步。

progress: 由四个等权检查点派生（S1～S4）；当前 **4/4**。

## 审计策略

| 阶段 / 项 | 默认模式 | 说明 |
|-----------|----------|------|
| S1 冻结 | self | 方案落盘 + 证据核对 |
| F05 新增 API | self | Cron preview 端点涉及新契约，但低风险 |
| S4 关门 | self | 常规关门自审 |

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|
| I-001 | required | F05 Cron 预览端点响应结构（description/nextRuns）与前端调用方式 | S1 方案冻结 | S1 | 检查 scheduledtasks 包解析能力 | **closed** | D-001 §2：`{ description, nextRuns }` |
| I-002 | non-blocking | F10 页脚字段写入 settings 资源的 schema/API 路径 | S1 方案冻结 | S1 | 检查 settings schema | **closed** | D-001 §5：现有 settings schema/API 扩展字段 |

## 父目标

- [GOAL-024-w16-user-perspective-improvements](../GOAL-024-w16-user-perspective-improvements/00-meta.md)

## 台账布局

- `01-decision/`：D-NNN；`02-execution/`：E-NNN；`03-audit/`：A-NNN。
