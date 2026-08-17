---
id: GOAL-015-w14-user-perspective-review
title: W14 · 真实用户视角审视 API/Web 并落盘改进项台账（F-01～F-14）；整改 deferred（待用户裁决）
status: done
parent: GOAL-001-design-implementation-conformance
created: 2026-08-17
updated: 2026-08-17
version: 0.3.0
progress: 4/4
---

# GOAL-015 · W14 · 真实用户视角审视 API/Web（改进项落盘）

VP-010 / workspace-010 的**第十四波**（用户 2026-08-17 点名立项）：以**真实用户/管理员视角**审视 `apps/api` 与 `apps/web` 已实现功能，找出「并非很小」的改进点并**在本波落盘**。本波交付 = **审视报告 + 改进项台账（F-01～F-14）+ 审计会签 + 台账同步**；**不实施整改**（整改为 deferred backlog，需用户裁决后另起波次）。

## 当前边界

- **本波交付（已完成）**：S1 审视执行（E-001）；S2 台账与待决项落盘（D-001：F-01～F-14 + I-001/I-002 登记）；S3 独立交叉审计（A-002，grok-4.6 · pass）；S4 审计响应 + 台账同步 + 关门自审（A-003；E-002）。
- **非范围 / deferred**：F-01～F-14 的**整改实施**一律 deferred，需用户对 in-scope/优先级给出指示后另起波次执行（I-001 门禁面向未来整改波次）。不改协议 schema 结构；不改 Profile 默认集 / 模块矩阵 / Manifest 装配语义；不新增业务域模块；不做视觉重设计。本波无任何业务代码改动。

## 成功标准与路线图（P-001）

- [x] **S1 · 审视执行**：api/web 全量功能面走查 + 关键发现证据复核（E-001）
- [x] **S2 · 台账与待决项落盘**：改进项台账 F-01～F-14（D-001）+ 信息项登记（I-001/I-002）
- [x] **S3 · 独立交叉审计**：grok-4.6（reasoning high）对 S1/S2 证据与关门边界做 independent 审计（A-002，verdict pass）
- [x] **S4 · 台账同步与关门**：响应 A-002（F-001~F-003 全部处理）+ goal-tree/workspace 同步 + git 提交 + 关门自审（A-003）

progress: 由四个等权检查点派生（S1～S4）；当前 **4/4**（2026-08-17 关门；A-003 self pass；A-002 independent pass）。

## 审计策略

| 阶段 / 项 | 模式 | 说明 |
|-----------|------|------|
| S1 审视 | self | 只读审视，可逆；证据逐条核对（A-001） |
| S2 台账 | self | 只读落盘；无代码改动 |
| S3 关门前置 | independent | 用户书面偏好 grok-4.6 · reasoning high；审核 S1/S2 证据与「本波不实施整改」诚实性（A-002，pass） |
| S4 关门 | self | 常规关门自审（A-003，pass） |

## 审计响应（P-003 · A-002）

| finding | 级别 | 响应 |
|---------|------|------|
| F-001 `00-meta` S3/S4 预勾与 2/4 不一致 | non-blocking | **fixed**：S3/S4 已实际完成，检查点全勾，progress=4/4，当前边界「已完成/进行中」措辞修正 |
| F-002 D-001 §3 用本波检查点号描述未来整改阶段 | non-blocking | **fixed**：D-001 §3 加「以下阶段号属于未来整改波次」标注 |
| F-003 F-14 通知空态子项过述为英文硬编码 | non-blocking | **fixed**：F-14 移除该子项；改为「空收件箱语义文案欠佳（本地化已齐，`feedback.noItemsMatch`）」 |

A-002 无 required finding，结论「可放行 W14 关门（S4）」。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|
| I-001 | required（对整改波次） | F-01～F-14 的 in-scope / defer / 优先级 | **未来整改波次**（非本波） | 未来整改波次 S2 | 用户裁决（P-004） | **deferred** | 延期理由：用户原始请求仅要求「审视并落盘」，未指示实施整改；责任人：用户；下一复核触发：用户指示开始整改。A-002 核验三字段齐备、非借重组绕过门禁。本波关门不受其阻断 |
| I-002 | non-blocking | F-01 定时任务 handler 目录暴露方式（新增端点 / 静态选项 / fork 扩展点） | 未来整改波次（F-01） | 未来整改波次 S3 | as-built + 方案 | **collecting** | 现 `HandlerKeys()` 仅 `system.noop`；需方案（D-001 §3） |

无对本波关门构成阻断的开放 required 信息项。

## 父目标

- [GOAL-001-design-implementation-conformance](../GOAL-001-design-implementation-conformance/00-meta.md)

## 台账布局

- `01-decision/`：D-NNN；`02-execution/`：E-NNN；`03-audit/`：A-NNN。
- 跨区引用用 Q2 路径（workspace-protocol §2.6）。