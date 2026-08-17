---
title: 目标树 · workspace-010-design-implementation-conformance
status: active
created: 2026-08-11
updated: 2026-08-17
parent: null
version: 0.20.0
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
├── GOAL-011-w10-account-page-conformance [done] · W10 个人中心页面层符合性（数据权限页修复 + 表格样式刷新）（4/4）
├── GOAL-012-w11-mfa-ux-review [done] · W11 个人中心 MFA 缺陷修复与全局 UX 审视整改（5/5）
├── GOAL-013-w12-product-surface-intent [done] · W12 产品面交互意图对齐（顶栏下拉 / 列表搜索 / 个人中心 Tabs / 我的钱包 / 回收站时间 / YAML 模块配置）（4/4）
├── GOAL-014-w13-settings-tabs-and-topbar [done] · W13 · 设置页 Tabs 化与顶栏/搜索交互打磨 + 个人中心头像上传 + 通知中心交互修正 + 列表筛选即时生效（设置页功能单元 Tabs / 移动端品牌条 / 搜索框组贴合 / 明暗-语种按键对调 / 移动端汉堡靠左 / 头像上传 / 通知点击即读 / 筛选即时生效）（4/4）
├── GOAL-015-w14-user-perspective-review [active] · W14 · 真实用户视角审视 API/Web 并落盘改进项台账（F-01～F-14）+ 整改承接（7/8）
│   ├── GOAL-016-w14-rectification-batch-a [done] · 整改批 A · 功能面补全（F-01～F-04）（4/4）
│   ├── GOAL-017-w14-rectification-batch-c [done] · 整改批 C · 调试痕迹清理（F-08～F-10）（4/4）
│   ├── GOAL-018-w14-rectification-batch-d [done] · 整改批 D · 表单与无障碍（F-11～F-14）（4/4）
│   └── GOAL-019-w14-rectification-batch-b [active] · 整改批 B · 一致性硬化（F-05～F-07）（0/4）
```

**W6（2026-08-14 关门，3/3）**：F-1 修复——claim `GIT_COMMIT` 接线、nginx `upstream` 作用域、smoke.sh SM-007 按 profile 页面集；V-007 exit 8 + **V-008 exit 0 完整绿**（SM-006 PASS）；**go 判定：恢复可消费**（冻结命令全部可执行）。

**W5（2026-08-14 关门，4/4）**：recordView 按 schema 声明渲染字段（标题/顺序），缺失/异常 fail-open 兜底；users/roles/activity schema + i18n + 测试；dev 脚本与 QUICKSTART 卫生。HEAD 回归 V-001～V-006 绿；**go 判定：无影响、不暂挂**（未改 Profile 默认集/模块矩阵/Manifest 装配/协议 pin）。A-001 记录跨门禁 F-1（容器 smoke 复现性破损，W3 引入）移交 freshness review。

**W12（2026-08-16 关门，4/4）**：S1/S2 冻结（D-002～D-008）；S3 实施 T-05（回收站时间 ISO）/ T-01（顶栏用户下拉）/ T-03（个人中心 Tabs）/ T-02（列表搜索矩阵 12 页）/ T-06（模块启用只认 config.yaml，废除 APP_PROFILE/APP_MODULES_ENABLED）；T-04 移交 [workspace-011] GOAL-022；S4 回归 Go 全量 0 FAIL + Web 1027/1027 + tsc 0；A-001 self pass；A-002 grok independent conditional（required F-001 台账问题已 fixed，F-003/F-004 fixed，F-005 accepted）；T-06 go 判定「部署契约变化、默认集不变 → 不暂挂」。

**W11（2026-08-15 关门，5/5）**：S1 裁决（D-001/D-002/D-003）；S2 MFA 三缺陷修复（二维码、401→400 分轨、解绑成功提示+登出、错码重填）；S3 UX P0（optionsSource 上游对象形态 + /api/permissions、/api/menu-items 目录端点）；S4 UX P1（Toast 浮动、8 页搜索表单、行操作收纳、分页增强、空状态）；S5 回归 Go 全量 + Web 1002/1002 + tsc 0；审计 A-001 self pass + A-002 independent（grok）conditional→resolved（F-001~F-007 全 fixed）+ A-003 closeout self pass；go 判定：无影响不暂挂。

**W10（2026-08-15 关门，4/4）**：数据权限页（workspace-011 GOAL-016 交付）七层根因修复（view→body、table props 化、rowKey、PATCH resource 入 body、shield 图标、列表信封、capability 声明）+ 列表翻页滚动位置保持 + 通用表格组件样式刷新（列宽/通用截断/空值兜底/表头层级/ghost 按钮/悬停/padding）与时间本地化格式 + 页脚偏移；参考样式对齐裁决 user-overruled（实测不好看，撤销）；A-001/A-002 self 审计 pass，无 required findings；Go 全量 + Web 991/991 绿；go 判定：无影响不暂挂。

**W9（2026-08-15 关门，6/6）**：设置页【品牌】图标由 URL 填写改为上传——专用 brand-assets 存储（非文件库/非通用 uploads 仓）+ 公开 GET（nosniff/sandbox/immutable）+ 服务端重编码（PNG/JPEG/GIF/WebP→PNG/JPEG、512/64 限幅、q82、≤4MiB、8192px 解压炸弹防线）+ config.yaml 参数 + 替换/清空/重置/启动 GC 清理闭环 + schema 上传控件/移除按钮/i18n/错误码契约；Go 全量 + Web 967/967 + 活栈点验；S6 cross 审计 A-001 self + A-002 independent（grok-4.6 high）**pass**，全部 findings fixed；go 判定：无影响不暂挂。

Root **保持 active**。W1/W2/W3/W4 均关门；W4 六检查点全部完成（2026-08-13 关门：S6 cross 审计
A-003 independent + A-004 self，BLOCKING 清零，F-1/F-2/F-3 全 fixed，E-004 浏览器点验）；不推导 Root/VP done。


**W13（2026-08-16 五次关门，4/4）**：设置页功能单元 Tabs（恢复默认常驻 Tabs 外）；移动端品牌条独立一行（<lg）+ 汉堡按键靠工具栏最左；搜索框组【文本框+搜索键】零间隙贴合恒同行；顶栏亮暗/语种按键对调；T-05 个人中心头像上传（共享 RasterAssetStore、迁移 0035/0036、/me 快照 + 用户菜单展示）+ 顶栏头像即时刷新修复；T-06 通知中心交互修正（铃铛条目点击跳转展开并标已读、列表点击即读+行内展开、移除行内已读 action、未读数即时刷新——notification-center 自定义组件 + notifications.read 响应头）；T-07 列表筛选即时生效（下拉等筛选项变动立即重新筛选、文本框+搜索键提交式、筛选记录以已提交查询为真相源）；A-001～A-005 self pass；回归 Go 0 FAIL + vitest 1037/1037 + tsc 0 + e2e admin/mvp 全绿（含 W11/W12 遗留 e2e 断言修复）；go 判定各轮均无影响、不暂挂。

**W14（2026-08-17 用户结构裁决：active · 6/8，整改完成前不得 done）**：真实用户视角审视 API/Web 已实现功能；S1 审视完成（E-001）；S2 台账落盘（D-001：F-01～F-14；I-002 collecting）；S3 独立审计 A-002（grok-4.6 · pass，三条 non-blocking 全 fixed）；S4 审计响应 + 台账同步 + I-001 用户书面裁决（D-003）。多次关门尝试（E-002/A-003 违规、E-004/A-005）均被用户否决/修正（E-003/A-004、E-005/A-006）。**整改按 D-003 分批 A→C→D→B 作为 GOAL-015 下级子目标渐进添加**；F-01 新增端点 / F-04 存 messageKey / F-08 直接移除。整改完成前 GOAL-015 保持 active。批 A（R1）、批 C（R2）、批 D（R3）已完成；批 B 已渐进添加。

**W14-批A（2026-08-17 关门，done 4/4）**：GOAL-015 下级整改子目标 GOAL-016-w14-rectification-batch-a——功能面补全 F-01（定时任务 handler 新增端点）/ F-02（数据权限范围设置）/ F-03（审计结构化过滤与导出）/ F-04（通知本地化 messageKey）。S1 冻结（D-001）；S2/S3 实施回归（Go 全量 + Web 1041/1041 + tsc + build）；独立审计 A-001 conditional（F-001 fixed、F-002/F-003 响应）+ 自审 A-002 pass。批 C 已渐进添加。

**W14-批C（2026-08-17 关门，done 4/4）**：GOAL-015 下级整改子目标 GOAL-017-w14-rectification-batch-c——调试痕迹清理 F-08（移除 pageId/route 技术框）/ F-09（反馈文案本地化与去错误码前缀）/ F-10（Schema 加载失败友好化）。S1 冻结（D-001）；S2/S3 实施回归（Web 全量 1041/1041 + tsc + build）；自审 A-001 pass。批 D 已渐进添加。

**W14-批D（2026-08-17 关门，done 4/4）**：GOAL-015 下级整改子目标 GOAL-018-w14-rectification-batch-d——表单与无障碍 F-11（必填标记）/ F-12（确认对话框焦点）/ F-13（桌面表格键盘选中）/ F-14（小缺口）。S1 冻结（D-001）；S2/S3 实施回归（Web 全量 1041/1041 + tsc + build）；自审 A-001 pass。批 B 已渐进添加。

**W14-批B（2026-08-17 立项，active 0/4）**：GOAL-015 下级整改子目标 GOAL-019-w14-rectification-batch-b——一致性硬化 F-05（列表端点校验）/ F-06（错误码与目录）/ F-07（搜索排序一致性）。五件套落盘；尚未开工（S1 冻结待推进）。
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
| GOAL-011-w10-account-page-conformance | W10 · 个人中心页面层符合性（数据权限页修复 + 表格样式刷新） | GOAL-001-design-implementation-conformance | done | 4/4 | 2026-08-15 |
| GOAL-012-w11-mfa-ux-review | W11 · 个人中心 MFA 缺陷修复与全局 UX 审视整改（M-01～M-03 + U-01～U-14 落盘） | GOAL-001-design-implementation-conformance | done | 5/5 | 2026-08-15 |
| GOAL-013-w12-product-surface-intent | W12 · 产品面交互意图对齐（顶栏下拉 / 列表搜索 / 个人中心 Tabs / 我的钱包 / 回收站时间 / YAML 模块配置） | GOAL-001-design-implementation-conformance | done | 4/4 | 2026-08-16 |
| GOAL-014-w13-settings-tabs-and-topbar | W13 · 设置页 Tabs 化与顶栏/搜索交互打磨 + 个人中心头像上传 + 通知中心交互修正 + 列表筛选即时生效（设置页功能单元 Tabs / 移动端品牌条 / 搜索框组贴合 / 明暗-语种按键对调 / 移动端汉堡靠左 / 头像上传 / 通知点击即读 / 筛选即时生效） | GOAL-001-design-implementation-conformance | done | 4/4 | 2026-08-16 |
| GOAL-015-w14-user-perspective-review | W14 · 真实用户视角审视 API/Web 并落盘改进项台账（F-01～F-14）+ 整改承接（整改完成前不 done） | GOAL-001-design-implementation-conformance | active | 7/8 | 2026-08-17 |
| GOAL-016-w14-rectification-batch-a | W14 整改批 A：功能面补全（F-01 定时任务 handler / F-02 数据权限范围设置 / F-03 审计结构化过滤与导出 / F-04 通知本地化 messageKey） | GOAL-015-w14-user-perspective-review | done | 4/4 | 2026-08-17 |
| GOAL-017-w14-rectification-batch-c | W14 整改批 C：调试痕迹清理（F-08 移除调试框 / F-09 反馈文案本地化与去错误码前缀 / F-10 Schema 加载失败友好化） | GOAL-015-w14-user-perspective-review | done | 4/4 | 2026-08-17 |
| GOAL-018-w14-rectification-batch-d | W14 整改批 D：表单与无障碍（F-11 必填标记 / F-12 确认对话框焦点 / F-13 桌面表格键盘选中 / F-14 小缺口） | GOAL-015-w14-user-perspective-review | done | 4/4 | 2026-08-17 |
| GOAL-019-w14-rectification-batch-b | W14 整改批 B：一致性硬化（F-05 列表端点校验 / F-06 错误码与目录 / F-07 搜索排序一致性） | GOAL-015-w14-user-perspective-review | active | 0/4 | 2026-08-17 |

## 维护说明

- Root 是长期能力容器；`status: done` 仅在程序废弃或 `primary_plan` 迁移且用户确认时使用。
- 波次 progress 只写在子目标；不得用波次完成数推导 Root done。
- 层级唯一来源是目标 `00-meta.md` 的 `parent`。
