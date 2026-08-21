---
id: GOAL-011-w10-account-page-conformance
title: W10 · 个人中心页面层符合性（参考样式对齐 + 数据权限页去留）
status: done
parent: GOAL-001-design-implementation-conformance
created: 2026-08-15
updated: 2026-08-15
version: 0.2.0
progress: 4/4
---

# GOAL-011 · W10 · 个人中心页面层符合性（参考样式对齐 + 数据权限页去留）

VP-010 / workspace-010 的**第十波**（用户 2026-08-15 裁决立项）：承载「个人中心页（/account）修改」的治理上下文。该修改已在 workspace-011 语境以提交 65d2a74 完成首轮实现（会话列表标题 + 状态筛选 + 翻页控件 + MFA 上移）；本波处理两项用户问题 + 通用表格组件样式刷新，**2026-08-15 关门（4/4）**：

1. **数据权限页去留**（用户提出）→ **裁决：保留并修复**。页面存在七层缺陷（view→body、table props 化、rowKey、PATCH resource 入 body、shield 图标、列表信封、capability 声明），全部修复并验证（E-002）。
2. **参考样式对齐**（用户提出）→ **裁决：不采纳（user-overruled）**。曾按 VP-005 参考样式实现（一体化页脚 + chevron + Showing X to Y of Z + 副标题），用户实测后裁决「没有之前的好看，撤销修改要求」，改动已回退（E-004 记录）；原样式保留并做增量体验优化。
3. **增量体验优化**（用户确认）：列表翻页滚动位置保持（E-003）、通用表格组件样式刷新（列宽/通用截断/空值兜底/表头层级/ghost 操作按钮/行悬停/padding）+ 时间显示本地化格式 + 页脚文字偏移（E-004）。

## 当前边界

- 范围：/account 页面会话列表的视觉符合性（标题行 / 表格页脚翻页行 / 相关 token 映射）、data-permission 页面的面层处置（保留 / 隐藏 / 移除）及对应回归。
- 非范围：不改 /account 的功能语义（筛选/翻页/吊销已交付）；不重做 VP-005 全量视觉基线；data-permission 模块内部功能（策略/范围语义）不在本波重写。

## 成功标准与路线图（P-001）

- [x] **S1 · 方案冻结**：I-001 / I-002 / I-003 三项用户裁决落盘（D-001），信息表 closed
- [x] **S2 · 实施**：data-permission 七层修复（E-002）+ 翻页滚动稳定（E-003）+ 表格组件样式刷新与时间格式化（E-004）；参考样式裁决不采纳并回退
- [x] **S3 · 验证与回归**：Go 全量全绿 + Web 991/991 + tsc 0；相关测试更新
- [x] **S4 · 审计与关门**：A-001 scaffold self（pass）+ A-002 关门 self 审计（pass，无 required findings）；go 判定：无影响不暂挂；goal-tree / workspace 同步

progress: 由四个等权检查点派生（S1～S4 全勾后 4/4）；当前 **4/4**（2026-08-15 关门）。

## 审计策略

本目标为**产品面/视觉符合性**波次（非 security/data/migration 门禁；若 I-002 裁决为模块级移除则升级为 cross——涉及 migration 与契约）。默认按 P-003 采用 **self**；触发模块移除时在 S4 改为 **cross**（provider 惯例 grok，执行时确认）。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 |
|----|------|-----------------|----------|--------------|-----------------|------|
| I-001 | required | 参考样式对齐范围与判定：标题行与翻页控件行是否按参考改；差异清单见 E-001 | 方案 | S1 | 用户 2026-08-15 书面裁决：**不采纳**（实测参考样式「没有之前的好看」，撤销修改要求，改动已回退） | **closed**（user-overruled） |
| I-002 | required | 数据权限页去留：保留 / 仅隐藏菜单 / 模块级移除；来源与影响面见 E-001 | 方案 | S1 | 用户 2026-08-15 书面裁决：**明确设立的交付物，保留并修复**（页面报错 + 菜单图标） | **closed**（裁决：保留 + 修复；修复见 E-002） |
| I-003 | required | 样式对齐的 token 策略：引入参考 token 还是映射现有主题 token | 方案 | S1 | 随 I-001 裁决：不采纳参考样式，保持现有语义 token（明暗自适应） | **closed**（user-overruled） |

## 父目标

- [GOAL-001-design-implementation-conformance](../GOAL-001-design-implementation-conformance/00-meta.md)

## 台账布局

- `01-decision/`：D-NNN 决策（S1 起）；`02-execution/`：E-NNN 事实记录；`03-audit/`：A-NNN 审计意见。
- 跨区引用（workspace-011 的 GOAL-016 等）用 Q2 路径，见 workspace-protocol §2.6。
