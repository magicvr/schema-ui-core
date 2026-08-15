---
title: 目标树 · workspace-010-design-implementation-conformance
status: active
created: 2026-08-11
updated: 2026-08-15
parent: null
version: 0.8.0
workspace_id: workspace-010-design-implementation-conformance
---

# 目标树 · 设计意图与实现符合性（持续对齐程序）

> 工作区：`workspace-010-design-implementation-conformance`
> canonical：`docs/workspaces/workspace-010-design-implementation-conformance/`
> Root：`GOAL-001-design-implementation-conformance`（**长期程序容器 · active**）
> primary_plan：`VP-010-design-implementation-conformance`（**active**）

## 树

```text
GOAL-001-design-implementation-conformance [active]  · 持续符合性程序
├── GOAL-002-w1-examples-optional-module [done]       · W1 范例面可选化
├── GOAL-003-demo-profile [done]                      · W2 demo Profile：mvp + 范例
├── GOAL-004-w3-schema-host-protocol-conformance [done] · W3 协议优先的 Host/App 符合性整改
├── GOAL-005-w4-long-content-presentation [done]      · W4 长内容列截断与详情换行
├── GOAL-006-w5-recordview-declared-fields [done]     · W5 recordView 声明字段（declared-fields 契约 + dev 卫生）
├── GOAL-007-w6-container-smoke-reproducibility [done] · W6 容器 smoke 复现性修复（F-1a/b/c）
├── GOAL-008-w7-yaml-config [done]  · W7 YAML 主配置体系（config.yaml + env 仅敏感信息）（5/5）
├── GOAL-009-w8-component-visual-style [done] · W8 组件视觉样式优化（语种下拉 / 明暗按钮 / 下拉暗色审计）（5/5）
├── GOAL-010-w9-branding-asset-upload [done] · W9 品牌图标上传（专用资产存储 + 自动图像处理）（6/6）
└── GOAL-011-w10-account-page-conformance [active] · W10 个人中心页面层符合性（参考样式对齐 + 数据权限页去留）（0/4 · 待用户裁决）
```

**W6（2026-08-14 关门，3/3）**：F-1 修复——claim `GIT_COMMIT` 接线、nginx `upstream` 作用域、smoke.sh SM-007 按 profile 页面集；V-007 exit 8 + **V-008 exit 0 完整绿**（SM-006 PASS）；**go 判定：恢复可消费**（冻结命令全部可执行）。

**W5（2026-08-14 关门，4/4）**：recordView 按 schema 声明渲染字段（标题/顺序），缺失/异常 fail-open 兜底；users/roles/activity schema + i18n + 测试；dev 脚本与 QUICKSTART 卫生。HEAD 回归 V-001～V-006 绿；**go 判定：无影响、不暂挂**（未改 Profile 默认集/模块矩阵/Manifest 装配/协议 pin）。A-001 记录跨门禁 F-1（容器 smoke 复现性破损，W3 引入）移交 freshness review。

**W10（2026-08-15 建立，0/4 · 未关门）**：个人中心页首轮修改（会话列表标题/状态筛选/翻页控件/MFA 上移，提交 65d2a74）治理上下文承载到本区；两项用户问题开放待裁决——(1) 标题行与翻页控件行是否对齐 VP-005 参考样式（schema_ui_core_data_table，E-001 差异清单）；(2) 数据权限页（/data-permission，workspace-011 GOAL-016 交付）去留。I-001/I-002/I-003 为 required 信息门禁，裁决闭合前不得推进 S2 与关门。

**W9（2026-08-15 关门，6/6）**：设置页【品牌】图标由 URL 填写改为上传——专用 brand-assets 存储（非文件库/非通用 uploads 仓）+ 公开 GET（nosniff/sandbox/immutable）+ 服务端重编码（PNG/JPEG/GIF/WebP→PNG/JPEG、512/64 限幅、q82、≤4MiB、8192px 解压炸弹防线）+ config.yaml 参数 + 替换/清空/重置/启动 GC 清理闭环 + schema 上传控件/移除按钮/i18n/错误码契约；Go 全量 + Web 967/967 + 活栈点验；S6 cross 审计 A-001 self + A-002 independent（grok-4.6 high）**pass**，全部 findings fixed；go 判定：无影响不暂挂。

Root **保持 active**。W1/W2/W3/W4 均关门；W4 六检查点全部完成（2026-08-13 关门：S6 cross 审计
A-003 independent + A-004 self，BLOCKING 清零，F-1/F-2/F-3 全 fixed，E-004 浏览器点验）；不推导 Root/VP done。

## 状态表

| id | title | parent | status | progress | updated |
|----|-------|--------|--------|----------|---------|
| GOAL-001-design-implementation-conformance | 设计意图与实现符合性（持续对齐程序） | null | active | —（程序容器，不用 n/n→done） | 2026-08-13 |
| GOAL-002-w1-examples-optional-module | W1 · 范例/演示产品面可选模块化 | GOAL-001-design-implementation-conformance | done | 6/6 | 2026-08-11 |
| GOAL-003-demo-profile | W2 · `demo` Profile：mvp + 范例页面 | GOAL-001-design-implementation-conformance | done | 6/6 | 2026-08-11 |
| GOAL-004-w3-schema-host-protocol-conformance | W3 · Schema-UI 语义对齐与 Host/App 协议增补 | GOAL-001-design-implementation-conformance | done | 6/6 | 2026-08-13 |
| GOAL-005-w4-long-content-presentation | W4 · 长内容列的列表截断与详情换行（以角色页权限/菜单为代表） | GOAL-001-design-implementation-conformance | done | 6/6 | 2026-08-13 |
| GOAL-006-w5-recordview-declared-fields | W5 · recordView 声明字段符合性（declared-fields 契约 + dev/文档卫生） | GOAL-001-design-implementation-conformance | done | 4/4 | 2026-08-14 |
| GOAL-007-w6-container-smoke-reproducibility | W6 · 容器 smoke 复现性修复（claim GIT_COMMIT / nginx upstream / SM-007 页面集） | GOAL-001-design-implementation-conformance | done | 3/3 | 2026-08-14 |
| GOAL-008-w7-yaml-config | W7 · YAML 主配置体系（config.yaml + env 仅敏感信息） | GOAL-001-design-implementation-conformance | done | 5/5 | 2026-08-14 |
| GOAL-009-w8-component-visual-style | W8 · 组件视觉样式优化（语种下拉 / 明暗按钮 / 下拉暗色审计） | GOAL-001-design-implementation-conformance | done | 5/5 | 2026-08-14 |
| GOAL-010-w9-branding-asset-upload | W9 · 品牌图标上传（专用资产存储 + 自动图像处理） | GOAL-001-design-implementation-conformance | done | 6/6 | 2026-08-15 |
| GOAL-011-w10-account-page-conformance | W10 · 个人中心页面层符合性（参考样式对齐 + 数据权限页去留） | GOAL-001-design-implementation-conformance | active | —（待用户裁决） | 2026-08-15 |
| GOAL-009-w8-component-visual-style | W8 · 组件视觉样式优化（语种下拉 / 明暗按钮 / 下拉暗色审计） | GOAL-001-design-implementation-conformance | done | 5/5 | 2026-08-14 |

## 维护说明

- Root 是长期能力容器；`status: done` 仅在程序废弃或 `primary_plan` 迁移且用户确认时使用。
- 波次 progress 只写在子目标；不得用波次完成数推导 Root done。
- 层级唯一来源是目标 `00-meta.md` 的 `parent`。
