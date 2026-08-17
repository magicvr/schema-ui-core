---
id: GOAL-015-w14-user-perspective-review
title: W14 · 真实用户视角审视 API/Web 并落盘改进项台账（F-01～F-14）+ 整改承接（子目标分批实施，全部完成后才关门）
status: active
parent: GOAL-001-design-implementation-conformance
created: 2026-08-17
updated: 2026-08-17
version: 0.9.0
progress: 6/8
---

# GOAL-015 · W14 · 真实用户视角审视 API/Web（改进项落盘）

VP-010 / workspace-010 的**第十四波**（用户 2026-08-17 点名立项）：以**真实用户/管理员视角**审视 `apps/api` 与 `apps/web` 已实现功能，找出「并非很小」的改进点并**在本波落盘**。本波交付 = **审视报告 + 改进项台账（F-01～F-14）+ 审计会签 + 台账同步 + 整改承接（子目标）**；**整改由 GOAL-015 的下级子目标分批实施，全部整改完成前本波保持 active，不得关门**。

## 当前边界

- **已完成（S1～S4 · 审视/台账/审计/同步）**：S1 审视执行（E-001）；S2 台账与待决项落盘（D-001：F-01～F-14 + I-001/I-002 登记）；S3 独立交叉审计（A-002，grok-4.6 · pass）；S4 审计响应 + 台账同步 + I-001 用户书面裁决（D-003）——本波曾两次尝试关门（E-002/A-003 违规 → E-004/A-005 曾视为合法），但**用户最终裁决：整改完成前 GOAL-015 不得 done**（E-005/A-006），GOAL-015 保持 **active**。
- **整改实施（用户裁决 D-003 · GOAL-015 子目标分批）**：F-01～F-14 **全部 in-scope**，分批 **A（F-01～F-04）→ C（F-08～F-10）→ D（F-11～F-14）→ B（F-05～F-07）**作为 **GOAL-015 的下级子目标**（可渐进添加；已建 [GOAL-016-w14-rectification-batch-a](../GOAL-016-w14-rectification-batch-a/00-meta.md)）；F-01 新增端点 / F-04 存 messageKey / F-08 直接移除已冻结为整改输入。**全部整改子目标完成后 GOAL-015 方可关门**。不改 Profile 默认集 / 模块矩阵 / Manifest 装配语义；不新增业务域模块；不做视觉重设计。

## 成功标准与路线图（P-001）

- [x] **S1 · 审视执行**：api/web 全量功能面走查 + 关键发现证据复核（E-001）
- [x] **S2 · 台账与待决项落盘**：改进项台账 F-01～F-14（D-001）+ 信息项登记（I-001/I-002）
- [x] **S3 · 独立交叉审计**：grok-4.6（reasoning high）对 S1/S2 证据做 independent 审计（A-002，verdict pass）
- [x] **S4 · 审计响应与台账同步**：响应 A-002（F-001~F-003 已处理）+ **I-001 用户书面裁决（D-003）** + goal-tree/workspace 同步 + git 提交（**不关门**；曾两次尝试关门均被用户否决/修正）
- [x] **R1 · 整改批 A（F-01～F-04）完成**——子目标 [GOAL-016-w14-rectification-batch-a](../GOAL-016-w14-rectification-batch-a/00-meta.md)（done · 4/4）
- [x] **R2 · 整改批 C（F-08～F-10）完成**——子目标 [GOAL-017-w14-rectification-batch-c](../GOAL-017-w14-rectification-batch-c/00-meta.md)（done · 4/4）
- [ ] **R3 · 整改批 D（F-11～F-14）完成**——子目标（渐进添加）
- [ ] **R4 · 整改批 B（F-05～F-07）完成**——子目标（渐进添加）
- [ ] **S5 · 关门**：全部整改子目标 done + 终审 + goal-tree/workspace 同步为 done

progress: 由八个等权检查点派生（S1～S4 + R1～R4）；当前 **6/8**（S1～S4 完成；R1 批 A、R2 批 C 完成；R3～R4 整改待后续批次推进，S5 关门待整改全部完成）。GOAL-015 在整改完成前保持 active。

## 审计策略

| 阶段 / 项 | 模式 | 说明 |
|-----------|------|------|
| S1 审视 | self | 只读审视，可逆；证据逐条核对（A-001） |
| S2 台账 | self | 只读落盘；无代码改动 |
| S3 关门前置 | independent | 用户书面偏好 grok-4.6 · reasoning high；审核 S1/S2 证据与「本波不实施整改」诚实性（A-002，pass） |
| S4 审计响应/同步 | self | 审计响应 + 台账同步（不做关门结论）——曾两次关门（A-003/A-005）均被用户否决/修正（A-004/A-006） |

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
| I-001 | required → **已裁决** | F-01～F-14 的 in-scope / defer / 优先级 | 整改子目标范围（R1～R4） | R1～R4 | 用户裁决（P-004） | **closed** | 用户 2026-08-17 书面裁决（D-003）：F-01～F-14 **全部 in-scope**、分批 A→C→D→B 作为 GOAL-015 子目标实施；F-01 新增端点 / F-04 存 messageKey / F-08 直接移除。裁决即整改范围输入；GOAL-015 待整改全部完成后关门 |
| I-002 | non-blocking | F-01 定时任务 handler 目录暴露方式（新增端点 / 静态选项 / fork 扩展点） | 未来整改波次（F-01） | 未来整改波次 S3 | as-built + 方案 | **collecting** | 现 `HandlerKeys()` 仅 `system.noop`；需方案（D-001 §3） |

I-001 已由用户书面裁决（D-003）关闭（整改范围确定）。GOAL-015 **active · 4/8**：S1～S4 完成；R1～R4 整改经子目标推进；S5 关门待整改全部完成。

## 父目标

- [GOAL-001-design-implementation-conformance](../GOAL-001-design-implementation-conformance/00-meta.md)

## 子目标（整改承接，渐进添加）

| 批 | 子目标 | status | 范围 |
|----|--------|--------|------|
| A | [GOAL-016-w14-rectification-batch-a](../GOAL-016-w14-rectification-batch-a/00-meta.md) | done · 4/4 | F-01 定时任务 handler / F-02 数据权限范围设置 / F-03 审计结构化过滤与导出 / F-04 通知本地化 messageKey |
| C | [GOAL-017-w14-rectification-batch-c](../GOAL-017-w14-rectification-batch-c/00-meta.md) | done · 4/4 | F-08 移除调试框 / F-09 反馈文案本地化与去错误码前缀 / F-10 Schema 加载失败友好化 |
| D | [GOAL-018-w14-rectification-batch-d](../GOAL-018-w14-rectification-batch-d/00-meta.md) | active · 0/4 | F-11 必填标记 / F-12 确认对话框焦点 / F-13 桌面表格键盘选中 / F-14 小缺口 |
| B | 渐进添加 | — | F-05～F-07（一致性硬化） |

## 台账布局

- `01-decision/`：D-NNN；`02-execution/`：E-NNN；`03-audit/`：A-NNN。
- 跨区引用用 Q2 路径（workspace-protocol §2.6）。