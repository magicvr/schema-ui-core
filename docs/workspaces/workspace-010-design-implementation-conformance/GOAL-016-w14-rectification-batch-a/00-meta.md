---
id: GOAL-016-w14-rectification-batch-a
title: W14 整改批 A · 功能面补全（F-01 定时任务 handler / F-02 数据权限范围设置 / F-03 审计日志结构化过滤与导出 / F-04 系统通知本地化 messageKey）
status: active
parent: GOAL-015-w14-user-perspective-review
created: 2026-08-17
updated: 2026-08-17
version: 0.1.0
progress: 0/4
---

# GOAL-016 · W14 整改批 A（F-01～F-04 功能面补全）

[GOAL-015](../GOAL-015-w14-user-perspective-review/00-meta.md)（W14）的**下级整改子目标（批 A）**：承接用户书面裁决（D-003）——F-01～F-14 **全部 in-scope、分批实施**。本子目标 = **批 A（功能面补全，F-01～F-04）**，为首个整改子目标；批 C/D/B 后续渐进添加为 GOAL-015 下级。

## 当前边界

- **范围（本波实施）**：F-01 定时任务可指定 handler（新增端点列出可用 handler，D-003 冻结）；F-02 数据权限页数据范围设置入口（updateScopes 接线）；F-03 审计/活动日志结构化过滤（事件/操作者/时间范围）与导出；F-04 系统通知本地化（存 messageKey）。
- **非范围**：批 C（F-08～F-10）、批 D（F-11～F-14）、批 B（F-05～F-07）在本子目标**不实施**（GOAL-015 后续渐进添加整改子目标）；不改 Profile 默认集 / 模块矩阵 / Manifest 装配语义（go 判定沿用 W14：无影响不暂挂，涉端点契约改动时须 go 复核）；不做视觉重设计。

## 成功标准与路线图（P-001）

- [ ] **S1 · 方案冻结**：F-01～F-04 设计（端点/schema/存储契约）+ 信息项登记
- [ ] **S2 · 实施**：F-01～F-04 代码与 schema/前端接线
- [ ] **S3 · 测试与回归**：单测/e2e + vitest/tsc + Go 全量（涉端点契约变更时 + go 复核）
- [ ] **S4 · 自审与关门**：审计 + 台账同步 + goal-tree/workspace 同步

progress: 由四个等权检查点派生（S1～S4）；当前 **0/4**（本子目标 2026-08-17 由 W14 用户裁决（D-003）+ GOAL-015 路线图立项，尚未开工）。

## 审计策略

| 阶段 / 项 | 默认模式 | 说明 |
|-----------|----------|------|
| S1 冻结 | self | 方案落盘 + 证据核对（F-01～F-04 来自 D-001 台账） |
| F-01/F-02 实施 | self（涉契约时升级 independent） | 端点/schema 契约变更时按兼容性升级 |
| F-04 实施 | independent | 通知文案存 messageKey（数据语义变化、旧文案迁移） |
| S4 关门 | self | 常规关门自审 |

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | non-blocking | F-01 handler 目录暴露端点路径与安全（鉴权/白名单） | S2 F-01 | S1 | 用户已裁决「新增端点」（D-003）；as-built 对照 | open | — | 待设计（S1） |
| I-002 | non-blocking | F-04 旧文案数据迁移（已存在英文硬编码通知是否重发/迁移） | S2 F-04 | S1 | 方案（存 messageKey 后旧记录处理） | open | — | 待设计（S1） |

## 父目标

- [GOAL-015-w14-user-perspective-review](../GOAL-015-w14-user-perspective-review/00-meta.md)

## 台账布局

- `01-decision/`：D-NNN；`02-execution/`：E-NNN；`03-audit/`：A-NNN。
- 跨区引用用 Q2 路径（workspace-protocol §2.6）。
