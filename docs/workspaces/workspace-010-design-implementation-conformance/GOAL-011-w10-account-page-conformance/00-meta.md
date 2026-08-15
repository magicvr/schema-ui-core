---
id: GOAL-011-w10-account-page-conformance
title: W10 · 个人中心页面层符合性（参考样式对齐 + 数据权限页去留）
status: active
parent: GOAL-001-design-implementation-conformance
created: 2026-08-15
updated: 2026-08-15
version: 0.1.0
progress: —
---

# GOAL-011 · W10 · 个人中心页面层符合性（参考样式对齐 + 数据权限页去留）

VP-010 / workspace-010 的**第十波**（用户 2026-08-15 裁决立项）：承载「个人中心页（/account）修改」的治理上下文。该修改已在 workspace-011 语境以提交 65d2a74 完成首轮实现（会话列表标题 + 状态筛选 + 翻页控件 + MFA 上移）；**本目标尚未到变门（closeout）条件**，两项用户提出的问题仍需考虑并裁决后方可关门：

1. **参考样式对齐**（用户提出）：会话列表的标题行与翻页控件行与 VP-005 视觉参考 `raw/stitch-vp005-visual-refs/exports/stitch_schema_ui_core_admin_console/schema_ui_core_data_table` 存在差异，参考样式明显更好看——是否应按参考改过来（I-001 / I-003）。
2. **数据权限页去留**（用户提出）：菜单中出现「数据权限」页（路由 /data-permission）——该页从何而来、是否应该删除（I-002）。**2026-08-15 用户裁决：明确设立的交付物（GOAL-016），保留并修复**——页面渲染报错（D-VAL 失败）与菜单无图标已修复（E-002）；样式对齐（标题行/翻页控件行）仍待 I-001/I-003 裁决。

## 当前边界

- 范围：/account 页面会话列表的视觉符合性（标题行 / 表格页脚翻页行 / 相关 token 映射）、data-permission 页面的面层处置（保留 / 隐藏 / 移除）及对应回归。
- 非范围：不改 /account 的功能语义（筛选/翻页/吊销已交付）；不重做 VP-005 全量视觉基线；data-permission 模块内部功能（策略/范围语义）不在本波重写。

## 成功标准与路线图（P-001）

- [ ] **S1 · 方案冻结**：I-001 / I-002 / I-003 三项用户裁决落盘（D-001），信息表 closed
- [ ] **S2 · 实施**：参考样式对齐（标题行 + 翻页控件行 + token 映射）与 data-permission 处置按裁决执行
- [ ] **S3 · 验证与回归**：Go 全量 + Web 全量 + tsc；相关测试更新
- [ ] **S4 · 审计与关门**：self 审计（影响面大则 cross + grok independent）；go 判定；goal-tree 同步

progress: 由四个等权检查点派生（S1～S4 全勾后 4/4）；当前 **—**（未开始，待用户裁决）。

## 审计策略

本目标为**产品面/视觉符合性**波次（非 security/data/migration 门禁；若 I-002 裁决为模块级移除则升级为 cross——涉及 migration 与契约）。默认按 P-003 采用 **self**；触发模块移除时在 S4 改为 **cross**（provider 惯例 grok，执行时确认）。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 |
|----|------|-----------------|----------|--------------|-----------------|------|
| I-001 | required | 参考样式对齐范围与判定：标题行（页面级 h1+副标题 vs 表格内 h2）与翻页控件行（一体化卡片页脚 + chevron-only vs 当前页码按钮导航）是否按参考改；差异清单见 E-001 | 方案 | S1 | 用户裁决（E-001 已列事实与建议） | **open** |
| I-002 | required | 数据权限页去留：保留 / 仅隐藏菜单 / 模块级移除；来源与影响面见 E-001 | 方案 | S1 | 用户 2026-08-15 书面裁决：**明确设立的交付物，保留并修复**（页面报错 + 菜单图标） | **closed**（裁决：保留 + 修复；修复见 E-002） |
| I-003 | required | 样式对齐的 token 策略：引入参考 token（surface-container-lowest / label-caps / body-sm 等）还是映射现有主题 token | 方案 | S1 | 用户裁决（E-001 已列差异） | **open** |

## 父目标

- [GOAL-001-design-implementation-conformance](../GOAL-001-design-implementation-conformance/00-meta.md)

## 台账布局

- `01-decision/`：D-NNN 决策（S1 起）；`02-execution/`：E-NNN 事实记录；`03-audit/`：A-NNN 审计意见。
- 跨区引用（workspace-011 的 GOAL-016 等）用 Q2 路径，见 workspace-protocol §2.6。
