---
title: 目标树 · workspace-010-design-implementation-conformance
status: active
created: 2026-08-11
updated: 2026-08-23
parent: null
version: 0.47.0
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
├── GOAL-015-w14-user-perspective-review [done] · W14 · 真实用户视角审视 API/Web 并落盘改进项台账（F-01～F-14）+ 整改承接（8/8）
│   ├── GOAL-016-w14-rectification-batch-a [done] · 整改批 A · 功能面补全（F-01～F-04）（4/4）
│   ├── GOAL-017-w14-rectification-batch-c [done] · 整改批 C · 调试痕迹清理（F-08～F-10）（4/4）
│   ├── GOAL-018-w14-rectification-batch-d [done] · 整改批 D · 表单与无障碍（F-11～F-14）（4/4）
│   └── GOAL-019-w14-rectification-batch-b [done] · 整改批 B · 一致性硬化（F-05～F-07）（4/4）
├── GOAL-020-w15-user-perspective-findings [done] · W15 · 真实用户视角二期审视与体验加固台账（W15-F01～W15-F14）（8/8）
│   ├── GOAL-021-w15-rectification-batch-a [done] · 整改批 A · F01/F02/F04/F05/F07（4/4）
│   ├── GOAL-022-w15-rectification-batch-b [done] · 整改批 B · F03/F11/F10/F12（4/4）
│   └── GOAL-023-w15-rectification-batch-c [done] · 整改批 C · F06/F08/F09/F13/F14（4/4）
├── GOAL-024-w16-user-perspective-improvements [done] · W16 · 真实用户视角未计划改进项台账与规划（W16-F01～W16-F10）+ 整改承接（8/8）
│   ├── GOAL-025-w16-rectification-batch-a [done] · W16 整改批 A · 安全与认证基线（F01 首次改密 / F07 一键下线其他 / F08 验证码与 MFA 备份）（4/4）
│   ├── GOAL-026-w16-rectification-batch-b [done] · W16 整改批 B · 核心资产与数据交互（F02 文件预览复制 / F03 导入模板与错误定位 / F04 金额格式化与调账警示）（4/4）
│   └── GOAL-027-w16-rectification-batch-c [done] · W16 整改批 C · 系统运维与通用外观（F05 Cron 预览 / F06 监控自动刷新 / F09 字典 Badge / F10 页脚版权）（4/4）
├── GOAL-028-w17-cron-preview-field-binding [done] · W17 · Cron 字段绑定与中文 describeCron（4/4）
├── GOAL-029-w18-preview-copy-and-import-modal [done] · W18 · 预览弹窗/复制链接与导入模态模板（4/4）
├── GOAL-030-w19-my-wallet-lazy-open-empty-state [done] · W19 · 我的钱包惰性开通与未开户空态（4/4）
├── GOAL-031-w20-notification-settings-in-account [done] · W20 · 通知设置迁入个人中心（4/4）
├── GOAL-032-w21-startup-db-identity [done] · W21 · 启动时数据库身份判定与迁移计划（5/5）
├── GOAL-033-w22-residual-closeout [done] · W22 · accepted-residual 残余全库清点收口（A 组修复 ×6 / B 组复核 ×6 / 台账卫生 ×3）（18/18）
└── GOAL-034-w23-admin-login-home-redirect [done] · W23 · admin 登录后 home 推导回归修复（N-001 承接）（4/4）
├── GOAL-035-w24-e2e-dual-dialect-matrix [done] · W24 · 浏览器 e2e 双数据库方言矩阵（收尾层双方言各测一次）（4/4）
├── GOAL-036-w25-page-performance-guardrails [active] · W25 · 页面性能问题全盘修复与防复发（钱包页 + 全局机制 + 防复发栅栏）（5/6 · 未闭门）
│   └── GOAL-037-w25-f008-wallet-reconcile-race [active] · W25 承接 · F-008 钱包对账竞态修复（池化+FK 时代偶发不一致）（0/4）
```

**W24（2026-08-23 关门，4/4）**：承接 GOAL-034 用户复审（强制 sqlite 属绕过；收尾层应双方言各测一次）。实现方言契约（默认 sqlite / pg 显式 opt-in）+ `cmd/e2e-pgset` scratch 库自动建/验/删 + `globalSetup` fail-fast 校验 + CI `profile×dialect` 矩阵；F-1 配置双载（双份 scratch 库）修复（E2E_PG_NAME 守卫 + DROP WITH FORCE）。回归：sqlite 9/9 + postgres 9/9（遗留 0）+ vitest 1088 + go 全绿 + tsc/build 0；A-001 self pass。I-001 实验先证（专用 pg 9/9 绿）closed。

**W25（2026-08-23 立项，GOAL-036 active 5/6，未闭门）**：我的钱包页面性能优化 → **用户升级为全盘修复与此类问题+防复发**（D-002 书面裁决：纳入 monitoring 定向刷新与 schema 注册校验；大表 COUNT 出局）。S1–S4：四因素诊断（SQLite 单连接串行 + fsync / 同 URL 重复请求 / 挂载即写整页重拉 / schema 重取）+ 钱包页实施（后端池 4 + WAL/busy_timeout/synchronous、前端 in-flight 合并 + 探活后写 + shell 级文档缓存）+ 26 页全盘扫描（system-monitoring 6×1、data-display 3×1 由全局机制覆盖）。S5 防复发（E-002）：`store_wal_test.go` 连接面白盒回归（池/WAL/超时，防回退 MaxOpenConns=1）+ 渲染层合并/定向刷新回归（statCard+chart、refreshList）+ `custom-components.schema.test.ts` 注册校验 + provider `refreshList` 定向刷新（monitoring tick 由整页 reloadList 改为只刷 /status，9→3 请求/tick）+ playbook §6 性能规范。**S6 待办（不闭门）**：I-001 e2e 已关闭（双 profile 全绿 9/9×2；暴露并修复后端「删用户遗留 user_roles 孤儿 → 角色永久不可删」缺陷 + 2 单元回归，E-004）；I-002 活栈计时复核已关闭（双栈实测：请求数 −47%～−86%、RTT150ms 呈现耗时 −1.4s、schema 缓存命中实证，E-005）；A-001（independent）响应完成（F-001～F-006 fixed，A-002/E-006）；F-007 fixed（E-007）；**F-008（wallet reconcile 竞态）移交下级子目标 GOAL-037 承接（2026-08-23 用户书面：GOAL-037 关门后再回归关门 GOAL-036）**。

**W23（2026-08-23 关门，4/4）**：N-001 根因 = e2e 挂具 store 隔离失效（本地 gitignored `configs/.env` 2026-08-21 建，`DB_DIALECT=postgres` 劫持临时 SQLite；全新种子 admin/admin 401，登录链第一步断开），**非路由回归**；W22 基线实验（git stash 无法移除 gitignored 文件）结论失效。修复：挂具钉死 `DB_DIALECT=sqlite`（playwright.config.ts）+ signInZh/sign-in fallback 等待硬化 + 连带 F-1（RowActionsMenu scroll-close 竞态，产品面）F-2（fallback 按钮等待）fixed。回归：go 全包 ok / vitest 1088 / tsc+build 0 / e2e admin 连续 5 轮 9/9 + mvp 9/9。A-001 self pass，required 0。I-001 closed。

**W21（2026-08-22 关门，5/5）**：启动 Identify→Plan→Execute。A-003 确认 F-001～F-003 fixed；A-004 self 关门 pass。Root/VP 保持 active。

**W22（2026-08-23 立项，GOAL-033）**：accepted-residual 全库清点收口。全库扫描命中约 420 行，存续 residual 23 项；用户 P-004 裁决执行：A 组一次性修复 ×6（W7 e2e admin M3 补跑 / W10 种子 admin must_change_password 迁移 / W11 导航边界单测 / W6 登录密码可见切换 / W9 上传扩展启发式 / W9 MFA verify 独立限流）+ B 组触发到期复核 ×6（W3 架构债 ×4：R4-I004 retention 过期、F-003b document 字节、C4-004 allowlist、C5-002 Start/Ready 矩阵；W8 F-007 freshness；W13 迁移器回写核实）+ 台账卫生 ×3（H1 GOAL-017 F-004 兑现回写、H2 W17 N-001 用词纠偏、H3 dogfood 夹具重命名）；C 组设计裁决类 11 项保留不动。安全面改动（A5/A6）关门需 independent 审计。

**W20（2026-08-18 关门，4/4）**：通知设置迁入个人中心。S1 D-001 + S2 实施 + S3 定向 + S4 A-001 self pass。Root/VP 保持 active。

**W19（2026-08-18 关门，4/4）**：进页 POST 惰性开通 + WALLET_NOT_FOUND 空态。S1 D-001 + S2 实施 + S3 定向 + S4 A-001 self pass。Root/VP 保持 active。

**W18（2026-08-18 关门，4/4）**：承接 GOAL-024 A-007 F-001/F-002。S1 D-001 + S2 实施 + S3 定向 + S4 A-001 self pass。Root/VP 保持 active。

**W17（2026-08-18 关门，4/4）**：承接 GOAL-024 A-005 F-004 / A-007 F-003。S1 D-001 + S2 实施 + S3 定向 + S4 A-001 self pass。Root/VP 保持 active。

**W16（2026-08-17 立项，8/8 关门）**：S1 台账建立（D-001）+ S2 技术方案（D-002，I-001 closed）+ S3 分批规划（D-003）+ S4 自审 A-002 pass；批 A GOAL-025、批 B GOAL-026、批 C GOAL-027 全部 done 4/4；S5 终审 A-003 pass；GOAL-024 已 done。Root/VP 保持 active。

**W15（2026-08-17 关门，8/8）**：S1～S5 + 批 A/B/C。I-001 = D-002。A-002/A-004 independent 的 required 已 fixed。Root/VP 仍 active。

**W6（2026-08-14 关门，3/3）**：F-1 修复——claim `GIT_COMMIT` 接线、nginx `upstream` 作用域、smoke.sh SM-007 按 profile 页面集；V-007 exit 8 + **V-008 exit 0 完整绿**（SM-006 PASS）；**go 判定：恢复可消费**（冻结命令全部可执行）。

**W5（2026-08-14 关门，4/4）**：recordView 按 schema 声明渲染字段（标题/顺序），缺失/异常 fail-open 兜底；users/roles/activity schema + i18n + 测试；dev 脚本与 QUICKSTART 卫生。HEAD 回归 V-001～V-006 绿；**go 判定：无影响、不暂挂**（未改 Profile 默认集/模块矩阵/Manifest 装配/协议 pin）。A-001 记录跨门禁 F-1（容器 smoke 复现性破损，W3 引入）移交 freshness review。

**W12（2026-08-16 关门，4/4）**：S1/S2 冻结（D-002～D-008）；S3 实施 T-05（回收站时间 ISO）/ T-01（顶栏用户下拉）/ T-03（个人中心 Tabs）/ T-02（列表搜索矩阵 12 页）/ T-06（模块启用只认 config.yaml，废除 APP_PROFILE/APP_MODULES_ENABLED）；T-04 移交 [workspace-011] GOAL-022；S4 回归 Go 全量 0 FAIL + Web 1027/1027 + tsc 0；A-001 self pass；A-002 grok independent conditional（required F-001 台账问题已 fixed，F-003/F-004 fixed，F-005 accepted）；T-06 go 判定「部署契约变化、默认集不变 → 不暂挂」。

**W11（2026-08-15 关门，5/5）**：S1 裁决（D-001/D-002/D-003）；S2 MFA 三缺陷修复（二维码、401→400 分轨、解绑成功提示+登出、错码重填）；S3 UX P0（optionsSource 上游对象形态 + /api/permissions、/api/menu-items 目录端点）；S4 UX P1（Toast 浮动、8 页搜索表单、行操作收纳、分页增强、空状态）；S5 回归 Go 全量 + Web 1002/1002 + tsc 0；审计 A-001 self pass + A-002 independent（grok）conditional→resolved（F-001~F-007 全 fixed）+ A-003 closeout self pass；go 判定：无影响不暂挂。

**W10（2026-08-15 关门，4/4）**：数据权限页（workspace-011 GOAL-016 交付）七层根因修复（view→body、table props 化、rowKey、PATCH resource 入 body、shield 图标、列表信封、capability 声明）+ 列表翻页滚动位置保持 + 通用表格组件样式刷新（列宽/通用截断/空值兜底/表头层级/ghost 按钮/悬停/padding）与时间本地化格式 + 页脚偏移；参考样式对齐裁决 user-overruled（实测不好看，撤销）；A-001/A-002 self 审计 pass，无 required findings；Go 全量 + Web 991/991 绿；go 判定：无影响不暂挂。

**W9（2026-08-15 关门，6/6）**：设置页【品牌】图标由 URL 填写改为上传——专用 brand-assets 存储（非文件库/非通用 uploads 仓）+ 公开 GET（nosniff/sandbox/immutable）+ 服务端重编码（PNG/JPEG/GIF/WebP→PNG/JPEG、512/64 限幅、q82、≤4MiB、8192px 解压炸弹防线）+ config.yaml 参数 + 替换/清空/重置/启动 GC 清理闭环 + schema 上传控件/移除按钮/i18n/错误码契约；Go 全量 + Web 967/967 + 活栈点验；S6 cross 审计 A-001 self + A-002 independent（grok-4.6 high）**pass**，全部 findings fixed；go 判定：无影响不暂挂。

Root **保持 active**。W1/W2/W3/W4 均关门；W4 六检查点全部完成（2026-08-13 关门：S6 cross 审计
A-003 independent + A-004 self，BLOCKING 清零，F-1/F-2/F-3 全 fixed，E-004 浏览器点验）；不推导 Root/VP done。


**W13（2026-08-16 五次关门，4/4）**：设置页功能单元 Tabs（恢复默认常驻 Tabs 外）；移动端品牌条独立一行（<lg）+ 汉堡按键靠工具栏最左；搜索框组【文本框+搜索键】零间隙贴合恒同行；顶栏亮暗/语种按键对调；T-05 个人中心头像上传（共享 RasterAssetStore、迁移 0035/0036、/me 快照 + 用户菜单展示）+ 顶栏头像即时刷新修复；T-06 通知中心交互修正（铃铛条目点击跳转展开并标已读、列表点击即读+行内展开、移除行内已读 action、未读数即时刷新——notification-center 自定义组件 + notifications.read 响应头）；T-07 列表筛选即时生效（下拉等筛选项变动立即重新筛选、文本框+搜索键提交式、筛选记录以已提交查询为真相源）；A-001～A-005 self pass；回归 Go 0 FAIL + vitest 1037/1037 + tsc 0 + e2e admin/mvp 全绿（含 W11/W12 遗留 e2e 断言修复）；go 判定各轮均无影响、不暂挂。

**W14（2026-08-17 关门，done 8/8）**：真实用户视角审视 API/Web 已实现功能；S1 审视完成（E-001）；S2 台账落盘（D-001：F-01～F-14；I-002 collecting）；S3 独立审计 A-002（grok-4.6 · pass，三条 non-blocking 全 fixed）；S4 审计响应 + 台账同步 + I-001 用户书面裁决（D-003）。多次关门尝试（E-002/A-003 违规、E-004/A-005）均被用户否决/修正（E-003/A-004、E-005/A-006）。**整改按 D-003 分批 A→C→D→B 作为 GOAL-015 下级子目标渐进添加**；F-01 新增端点 / F-04 存 messageKey / F-08 直接移除。全部整改子目标（批 A/C/D/B）已完成，S5 终审通过。

**W14-批A（2026-08-17 关门，done 4/4）**：GOAL-015 下级整改子目标 GOAL-016-w14-rectification-batch-a——功能面补全 F-01（定时任务 handler 新增端点）/ F-02（数据权限范围设置）/ F-03（审计结构化过滤与导出）/ F-04（通知本地化 messageKey）。S1 冻结（D-001）；S2/S3 实施回归（Go 全量 + Web 1041/1041 + tsc + build）；独立审计 A-001 conditional（F-001 fixed、F-002/F-003 响应）+ 自审 A-002 pass。批 C 已渐进添加。

**W14-批C（2026-08-17 关门，done 4/4）**：GOAL-015 下级整改子目标 GOAL-017-w14-rectification-batch-c——调试痕迹清理 F-08（移除 pageId/route 技术框）/ F-09（反馈文案本地化与去错误码前缀）/ F-10（Schema 加载失败友好化）。S1 冻结（D-001）；S2/S3 实施回归（Web 全量 1041/1041 + tsc + build）；自审 A-001 pass。批 D 已渐进添加。

**W14-批D（2026-08-17 关门，done 4/4）**：GOAL-015 下级整改子目标 GOAL-018-w14-rectification-batch-d——表单与无障碍 F-11（必填标记）/ F-12（确认对话框焦点）/ F-13（桌面表格键盘选中）/ F-14（小缺口）。S1 冻结（D-001）；S2/S3 实施回归（Web 全量 1041/1041 + tsc + build）；自审 A-001 pass。批 B 已渐进添加。

**W14-批B（2026-08-17 关门，done 4/4）**：GOAL-015 下级整改子目标 GOAL-019-w14-rectification-batch-b——一致性硬化 F-05（列表端点校验）/ F-06（错误码与目录）/ F-07（搜索排序一致性）。S1 冻结（D-001）；S2/S3 实施回归（Go 全量 + Web 全量 1041/1041 + tsc + build）；独立审计 A-001 fail→required 全 fixed + 自审 A-002 pass。
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
| GOAL-015-w14-user-perspective-review | W14 · 真实用户视角审视 API/Web 并落盘改进项台账（F-01～F-14）+ 整改承接 | GOAL-001-design-implementation-conformance | done | 8/8 | 2026-08-17 |
| GOAL-016-w14-rectification-batch-a | W14 整改批 A：功能面补全（F-01 定时任务 handler / F-02 数据权限范围设置 / F-03 审计结构化过滤与导出 / F-04 通知本地化 messageKey） | GOAL-015-w14-user-perspective-review | done | 4/4 | 2026-08-17 |
| GOAL-017-w14-rectification-batch-c | W14 整改批 C：调试痕迹清理（F-08 移除调试框 / F-09 反馈文案本地化与去错误码前缀 / F-10 Schema 加载失败友好化） | GOAL-015-w14-user-perspective-review | done | 4/4 | 2026-08-17 |
| GOAL-018-w14-rectification-batch-d | W14 整改批 D：表单与无障碍（F-11 必填标记 / F-12 确认对话框焦点 / F-13 桌面表格键盘选中 / F-14 小缺口） | GOAL-015-w14-user-perspective-review | done | 4/4 | 2026-08-17 |
| GOAL-019-w14-rectification-batch-b | W14 整改批 B：一致性硬化（F-05 列表端点校验 / F-06 错误码与目录 / F-07 搜索排序一致性） | GOAL-015-w14-user-perspective-review | done | 4/4 | 2026-08-17 |
| GOAL-020-w15-user-perspective-findings | W15 · 真实用户视角二期审视与体验加固台账（W15-F01～W15-F14）+ 整改承接 | GOAL-001-design-implementation-conformance | done | 8/8 | 2026-08-17 |
| GOAL-021-w15-rectification-batch-a | W15 整改批 A：F01 会话容灾 / F02 表格重试 / F04 JSON 信封 / F05 CORS / F07 refresh 错误码 | GOAL-020-w15-user-perspective-findings | done | 4/4 | 2026-08-17 |
| GOAL-022-w15-rectification-batch-b | W15 整改批 B：F03 时间格式 / F11 GET 只读 / F10 Retry-After / F12 分页 | GOAL-020-w15-user-perspective-findings | done | 4/4 | 2026-08-17 |
| GOAL-023-w15-rectification-batch-c | W15 整改批 C：F06 改密提示 / F08 校验码 / F09 Toast / F13 当前会话 / F14 细节 | GOAL-020-w15-user-perspective-findings | done | 4/4 | 2026-08-17 |
| GOAL-024-w16-user-perspective-improvements | W16 · 真实用户视角未计划改进项台账与规划（W16-F01～W16-F10）+ 整改承接 | GOAL-001-design-implementation-conformance | done | 8/8 | 2026-08-17 |
| GOAL-025-w16-rectification-batch-a | W16 整改批 A：安全与认证基线（F01 首次改密 / F07 一键下线其他 / F08 验证码与 MFA 备份） | GOAL-024-w16-user-perspective-improvements | done | 4/4 | 2026-08-17 |
| GOAL-026-w16-rectification-batch-b | W16 整改批 B：核心资产与数据交互（F02 文件预览复制 / F03 导入模板与错误定位 / F04 金额格式化与调账警示） | GOAL-024-w16-user-perspective-improvements | done | 4/4 | 2026-08-17 |
| GOAL-027-w16-rectification-batch-c | W16 整改批 C：系统运维与通用外观（F05 Cron 预览 / F06 监控自动刷新 / F09 字典 Badge / F10 页脚版权） | GOAL-024-w16-user-perspective-improvements | done | 4/4 | 2026-08-17 |
| GOAL-028-w17-cron-preview-field-binding | W17 · Cron 字段绑定与中文 describeCron | GOAL-001-design-implementation-conformance | done | 4/4 | 2026-08-18 |
| GOAL-029-w18-preview-copy-and-import-modal | W18 · 预览弹窗/复制链接与导入模态模板 | GOAL-001-design-implementation-conformance | done | 4/4 | 2026-08-18 |
| GOAL-030-w19-my-wallet-lazy-open-empty-state | W19 · 我的钱包惰性开通与未开户空态 | GOAL-001-design-implementation-conformance | done | 4/4 | 2026-08-18 |
| GOAL-031-w20-notification-settings-in-account | W20 · 通知设置迁入个人中心 | GOAL-001-design-implementation-conformance | done | 4/4 | 2026-08-18 |
| GOAL-032-w21-startup-db-identity | W21 · 启动时数据库身份判定与迁移计划 | GOAL-001-design-implementation-conformance | done | 5/5 | 2026-08-22 |
| GOAL-033-w22-residual-closeout | W22 · accepted-residual 残余全库清点收口（A 组修复 ×6 + B 组触发复核 ×6 + 台账卫生 ×3） | GOAL-001-design-implementation-conformance | done | 18/18 | 2026-08-23 |
| GOAL-034-w23-admin-login-home-redirect | W23 · admin 登录后 home 推导回归修复（N-001 承接） | GOAL-001-design-implementation-conformance | done | 4/4 | 2026-08-23 |
| GOAL-035-w24-e2e-dual-dialect-matrix | W24 · 浏览器 e2e 双数据库方言矩阵（收尾层双方言各测一次） | GOAL-001-design-implementation-conformance | done | 4/4 | 2026-08-23 |
| GOAL-036-w25-page-performance-guardrails | W25 · 页面性能问题全盘修复与防复发（钱包页 + 全局机制 + 防复发栅栏） | GOAL-001-design-implementation-conformance | active | 5/6 | 2026-08-23 |
| GOAL-037-w25-f008-wallet-reconcile-race | W25 承接 · F-008 钱包对账竞态修复（池化+FK 时代偶发不一致） | GOAL-036-w25-page-performance-guardrails | active | 0/4 | 2026-08-23 |


## 维护说明

- Root 是长期能力容器；`status: done` 仅在程序废弃或 `primary_plan` 迁移且用户确认时使用。
- 波次 progress 只写在子目标；不得用波次完成数推导 Root done。
- 层级唯一来源是目标 `00-meta.md` 的 `parent`。
