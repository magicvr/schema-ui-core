---
id: workspace-010-design-implementation-conformance
title: 设计意图与实现符合性工作区
status: active
root_goal: GOAL-001-design-implementation-conformance
canonical_scope: docs/workspaces/workspace-010-design-implementation-conformance/
shared_materials_catalog: none
vision_role: delivery
plan_refs:
  - VP-010-design-implementation-conformance
primary_plan: VP-010-design-implementation-conformance
created: 2026-08-11
updated: 2026-09-19
version: 0.59.0
parent: null
---

# 工作区上下文 · 设计意图与实现符合性

本工作区是 [VP-010-design-implementation-conformance](../../vision/plans/VP-010-design-implementation-conformance.md)（`active` · **长期设计意图—实现符合性程序**）的唯一 lead delivery workspace。

- **Root** 为长期程序容器（默认 `active`）。  
- **子目标** 为有界符合性审视/整改波次（可 `done`）。  
- 不因单波完成而关闭本区或 VP；不改变 Charter `primary_workspace`。  
- 与 [workspace-009-production-hardening](../workspace-009-production-hardening/workspace.md) **正交**：009 = 安全与健壮性；本区 = 架构/产品意图与 as-built 对齐。

## 绑定

| 字段 | 当前值 | 说明 |
|------|--------|------|
| 工作区 ID | `workspace-010-design-implementation-conformance` | 与本区目标及资料引用的 `workspace_id` 一致 |
| Root Goal | `GOAL-001-design-implementation-conformance` | `parent: null`；长期容器 |
| canonical 范围 | `docs/workspaces/workspace-010-design-implementation-conformance/` | 本区唯一目标状态范围 |
| 共享资料目录 | `none` | 暂无固定共享资料 |
| 愿景角色 | `delivery` | VP-010 lead；不改变 Charter primary workspace |
| 规划对齐 | `primary_plan` = `VP-010-design-implementation-conformance` | 持续程序意图 |

## 愿景对齐

Charter：`schema-ui-core-admin-foundation@0.4.0`。
VP-010 为设计意图—实现符合性持续程序；与 VP-008 `go` 消费有效性接口见该 VP。  
若本区波次改变 Profile 默认集 / 模块矩阵 / Manifest 装配语义，须按规则暂挂或重验证业务对 `go` 的消费。

## 波次（实现层指针）

| 波次 | 子目标 | status |
|------|--------|--------|
| W1 | GOAL-002-w1-examples-optional-module | **done**（6/6 · 2026-08-11 关门；go 已恢复） |
| W2 | GOAL-003-demo-profile | **done**（6/6 · 2026-08-11 关门；go 无影响不暂挂） |
| W3 | GOAL-004-w3-schema-host-protocol-conformance | **done**（6/6 · 2026-08-13 关门；S6 cross 审计 A-007/A-008，BLOCKING 清零；用户 P-004 裁决 account-locked 实现生产源；go 无影响不暂挂） |
| W4 | GOAL-005-w4-long-content-presentation | **done**（6/6 · 2026-08-13 关门；S6 cross 审计 A-003 independent + A-004 self，BLOCKING 清零，F-1/F-2/F-3 全 fixed，E-004 浏览器点验；go 无影响不暂挂） |
| W5 | GOAL-006-w5-recordview-declared-fields | **done**（4/4 · 2026-08-14 关门；recordView 声明字段 + fail-open + dev 卫生；HEAD 回归 V-001～V-006 绿；**go 无影响不暂挂**；A-001 跨门禁 F-1 移交 W6） |
| W6 | GOAL-007-w6-container-smoke-reproducibility | **done**（3/3 · 2026-08-14 关门；F-1a claim GIT_COMMIT 接线、F-1b nginx upstream 作用域、F-1c SM-007 页面集；V-007 exit 8 + V-008 exit 0 完整绿；**go 恢复可消费**） |
| W7 | GOAL-008-w7-yaml-config | **done**（5/5 · 2026-08-14 关门：A-003 grok 审计 pass，F-001~F-005 fixed；configs/config.yaml 权威 + ${VAR} 敏感引用 + env 覆盖；workspace-11 导航排序覆盖载体已就位） |
| W8 | GOAL-009-w8-component-visual-style | **done**（5/5 · 2026-08-14 关门：语种下拉 / 明暗按钮统一 / 下拉暗色审计；self 审计；go 无影响不暂挂） |
| W9 | GOAL-010-w9-branding-asset-upload | **done**（6/6 · 2026-08-15 关门：品牌图标 URL 填写 → 上传控件 + 专用资产存储 + 自动图像处理；S6 cross 审计 A-001 self + A-002 grok independent pass，findings 全 fixed；go 无影响不暂挂） |
| W10 | GOAL-011-w10-account-page-conformance | **done**（4/4 · 2026-08-15 关门：数据权限页七层修复 + 翻页滚动稳定 + 表格样式刷新/时间格式化；参考样式 user-overruled；A-001/A-002 self pass；go 无影响不暂挂） |
| W11 | GOAL-012-w11-mfa-ux-review | **done**（5/5 · 2026-08-15 关门：MFA 三缺陷修复 + UX P0/P1 实施；A-001 self pass + A-002 grok independent conditional→resolved + A-003 closeout self pass；Go 全量 + Web 1002/1002；go 无影响不暂挂） |
| W12 | GOAL-013-w12-product-surface-intent | **done**（4/4 · 2026-08-16 关门：T-05/T-01/T-03/T-02/T-06 实施；T-04 移交 GOAL-022；回归 Go 0 FAIL + Web 1027/1027；A-001 self pass + A-002 grok conditional（F-001 fixed / F-003·F-004 fixed / F-005 accepted）；T-06 go 判定不暂挂） |
| W13 | GOAL-014-w13-settings-tabs-and-topbar | **done**（4/4 · 2026-08-16 五次关门：设置页功能单元 Tabs + 移动端品牌条 + 汉堡靠左 + 搜索框组贴合 + 顶栏明暗/语种按键对调 + 个人中心头像上传（T-05）+ 顶栏头像即时刷新修复（E-008）+ 通知中心交互修正（T-06）+ 列表筛选即时生效（T-07）；A-001～A-005 self pass；回归 Go 0 FAIL + vitest 1037/1037 + tsc 0 + e2e admin/mvp 全绿（含 W11/W12 遗留 e2e 断言修复）；go 各轮均无影响不暂挂） |
| W14 | GOAL-015-w14-user-perspective-review | **done**（8/8 · 2026-08-17 关门：S1 审视 + S2 台账落盘（F-01～F-14）+ S3 独立审计 A-002（grok-4.6 · pass）+ S4 审计响应/同步 + I-001 用户书面裁决（D-003）。多次关门尝试被用户否决/修正后，整改按 D-003 分批 A→C→D→B 作为 GOAL-015 下级子目标渐进添加；批 A/C/D/B 全部完成，S5 终审通过；go 无影响不暂挂） |
| W14-批A | GOAL-016-w14-rectification-batch-a（**GOAL-015 下级**） | **done**（4/4 · 2026-08-17 关门：F-01～F-04 功能面补全实施完成；S1 冻结 D-001、S2/S3 回归 Go 全量 + Web 1041/1041 + tsc + build、A-001 independent conditional（F-001 fixed）+ A-002 self pass） |
| W14-批C | GOAL-017-w14-rectification-batch-c（**GOAL-015 下级**） | **done**（4/4 · 2026-08-17 关门：F-08～F-10 调试痕迹清理实施完成；S1 冻结 D-001、S2/S3 回归 Web 全量 1041/1041 + tsc + build、A-001 self pass） |
| W14-批D | GOAL-018-w14-rectification-batch-d（**GOAL-015 下级**） | **done**（4/4 · 2026-08-17 关门：F-11～F-14 表单与无障碍实施完成；S1 冻结 D-001、S2/S3 回归 Web 全量 1041/1041 + tsc + build、A-001 self pass） |
| W14-批B | GOAL-019-w14-rectification-batch-b（**GOAL-015 下级**） | **done**（4/4 · 2026-08-17 关门：F-05～F-07 一致性硬化实施完成；S1 冻结 D-001、S2/S3 回归 Go 全量 + Web 全量 1041/1041 + tsc + build、A-001 independent fail→fixed + A-002 self pass） |
| W15 | GOAL-020-w15-user-perspective-findings | **done**（8/8 · 2026-08-17 关门：A-004 required 已 fixed；Root/VP 仍 active） |
| W15-批A | GOAL-021-w15-rectification-batch-a（**GOAL-020 下级**） | **done** 4/4 |
| W15-批B | GOAL-022-w15-rectification-batch-b（**GOAL-020 下级**） | **done** 4/4 |
| W15-批C | GOAL-023-w15-rectification-batch-c（**GOAL-020 下级**） | **done** 4/4 |
| W16 | GOAL-024-w16-user-perspective-improvements | **done**（8/8 · 2026-08-17：S1～S5 完成；批 A/B/C 全部 done，GOAL-024 已关门） |
| W16-批A | GOAL-025-w16-rectification-batch-a（**GOAL-024 下级**） | **done**（4/4 · 2026-08-17：F01/F07/F08 实施、Go/Web 全量回归、independent A-001 + 响应 A-002 + 关门 A-003） |
| W16-批B | GOAL-026-w16-rectification-batch-b（**GOAL-024 下级**） | **done**（4/4 · 2026-08-17：F02/F03/F04 实施、Go/Web 全量回归、关门 A-001） |
| W16-批C | GOAL-027-w16-rectification-batch-c（**GOAL-024 下级**） | **done**（4/4 · 2026-08-17：F05/F06/F09/F10 实施、Go/Web 全量回归、关门 A-001） |
| W17 | GOAL-028-w17-cron-preview-field-binding | **done**（4/4 · 2026-08-18：S1～S4；A-001 self pass；go 不暂挂） |
| W18 | GOAL-029-w18-preview-copy-and-import-modal | **done**（4/4 · 2026-08-18：S1～S4；A-001 self pass；go 不暂挂） |
| W19 | GOAL-030-w19-my-wallet-lazy-open-empty-state | **done**（4/4 · 2026-08-18：S1～S4；A-001 self pass；go 不暂挂） |
| W20 | GOAL-031-w20-notification-settings-in-account | **done**（4/4 · 2026-08-18：S1～S4；A-001 self pass；go 不暂挂） |
| W21 | GOAL-032-w21-startup-db-identity | **done**（5/5 · 2026-08-22 关门：Identify/Plan；A-003 F-001～F-003 fixed；A-004 self pass；go 不暂挂） |
| W28 | [GOAL-040-w28-admin-passwd-convention](GOAL-040-w28-admin-passwd-convention/00-meta.md) | **done**（4/4 · 2026-09-06 立项并当日关门：ADMIN_PASSWD 现有库 admin 凭据声明约定 + AI 可发现性 + TEST_ADMIN 机制退役；真实栈验证 SM-001~005 PASS；A-001 self pass 0 required；用户确认关门；Root 保持 active） |
| W29 | [GOAL-041-w29-api-web-protocol-conformance](GOAL-041-w29-api-web-protocol-conformance/00-meta.md) | **done**（6/6 · 2026-09-06 关门：S1 分母 + S2 分类/cross + S3 custom 裁决 + S4 整改 + S5 运行时验证 + S6 cross 关门（A-008 walker 跨平台 fixed；用户书面确认）。Root 保持 active） |
| W30 | [GOAL-042-w30-w29-followup-supplement](GOAL-042-w30-w29-followup-supplement/00-meta.md) | **done**（3/3 · 2026-09-06 关门：F1 legacy 能力全量审计（守卫 35/35 + 32 schema 双向修正）+ F2 claim↔host-support 单源一致性 + F3 10 页行为单测；回归 Web 1332/1332 + Go 0 FAIL；A-001 self pass） |
| W31 | [GOAL-043-w31-cross-workspace-residual-closeout](GOAL-043-w31-cross-workspace-residual-closeout/00-meta.md) | **done**（4/4 · 2026-09-18 关门：跨工作区残余统一收口 + roadmap 登记节） |
| W32 | [GOAL-044-w32-r4-residual-seams](GOAL-044-w32-r4-residual-seams/00-meta.md) | **done**（4/4 · 2026-09-19 立项并当日关门：承接 `[workspace-038]` GOAL-005 A-001 F-002/F-003/F-004 —— 通用列值本地化 `valueLabels` + `refreshTable` 定向刷新 seam + `activeStatuses` 空闲不轮询；三处变异验证；已回填 038 `A-003` 三条 `fixed`） |
| W33 | [GOAL-045-w33-list-actions-slot-and-roles-trigger](GOAL-045-w33-list-actions-slot-and-roles-trigger/00-meta.md) | **done**（4/4 · 2026-09-19 立项并当日关门：列表页 actions **左侧插槽**本地扩展（`props.slot = "list-page-actions"`）+ `users` 页接入 + `roles` 页补齐「导出所选」（多选 + 触发节点，后端零改动）；副产物：修复两页 table 同索引导致的**跨页列状态串扰**；vitest 121/1476 + e2e 双 profile 全绿） |


## 固定共享资料引用

> `shared-materials/index.json` 只能提供候选路径与摘要。缺完整引用字段的行无效。

| reference_id | workspace_id | material_id | source | version | sha256 | purpose | local_record | status |
|--------------|--------------|-------------|--------|---------|--------|---------|--------------|--------|
| — | — | — | — | — | — | — | — | — |

## 波次补充 · W22（2026-08-23）

GOAL-033-w22-residual-closeout done 18/18（accepted-residual 全库清点收口；详见该目标五件套与本区 goal-tree）。移交跟踪槽 **N-001**：admin 登录后停留 `/` 未跳 `/dashboard`（先于 W22 存在的既有回归，基线实验证实；疑似 W14–W21 home 推导/路由漂移）→ 建议下一符合性波次承接。*（归因更正：N-001 实为挂具 store 隔离失效，非路由回归——见下方 W23 记录与 GOAL-034 D-001；W22 基线实验无法移除 gitignored 文件。）*

**W23（2026-08-23 关门，GOAL-034 done 4/4）**：承接 W22 移交槽 N-001——admin 登录后停留 `/` 未跳 `/dashboard`。根因 = e2e 挂具 store 隔离失效（本机 gitignored `configs/.env` 2026-08-21 建，`DB_DIALECT=postgres` 劫持挂具临时 SQLite，全新种子 admin/admin 401），**非 home 推导/路由回归**；W22「先于 W22 基线实验」结论因 git stash 无法移除 gitignored 文件而失效。修复：`playwright.config.ts` 钉死 `DB_DIALECT=sqlite` + signInZh/sign-in fallback 等待硬化 + 连带 F-1（RowActionsMenu scroll-close 竞态 → 仅触发钮位移才关闭）F-2 fixed。回归：go 全包 ok / vitest 1088 / tsc+build 0 / e2e admin 连续 5 轮 9/9 + mvp 9/9；A-001 self pass（required 0）；I-001 closed。N-001 移交槽闭合。（2026-08-23 用户复审：钉方言属绕过，正确形态 = 双方言矩阵，承接 GOAL-035。）

**W24（2026-08-23 关门，GOAL-035 done 4/4）**：承接 GOAL-034 用户书面复审（「强制 sqlite 属绕过；收尾层 e2e 应双方言各测一次」；实验先证专用 pg 全量 9/9 绿）。实现：`DB_DIALECT` 方言契约（默认 sqlite / pg 显式 opt-in，`.env` 无法再改道）+ `apps/api/cmd/e2e-pgset`（scratch 库 create/verify/drop/list，凭据与 API/pgtest 同源）+ `globalSetup/Teardown` fail-fast 校验与清理 + `npm run test:e2e:postgres` + CI `profile×dialect [sqlite,postgres]` 矩阵。F-1（Playwright 配置双载 → 双份 scratch 库）修复：`E2E_PG_NAME` 守卫复用 + `DROP WITH (FORCE)` + teardown 可见。回归：sqlite 9/9 + postgres 9/9（遗留 0）+ vitest 1088 + go 全绿 + tsc/build 0；A-001 self pass（required 0）；I-001 closed。

**W25（2026-08-23 立项，GOAL-036 active 5/6，未闭门）**：我的钱包页面性能优化 → **用户升级为全盘修复此类问题 + 防复发**（D-002 书面裁决；文件夹改名 `GOAL-036-w25-page-performance-guardrails`）。S1–S4：四因素诊断（SQLite `MaxOpenConns=1` 全局串行 + 逐提交 fsync；同 URL 展示节点重复请求；wallet-ensure 挂载即写 + 整页重拉；schema 每次导航重取）+ 钱包页实施（后端文件库池 4 + `_busy_timeout=5000&_journal_mode=WAL&_synchronous=NORMAL`，`:memory:` 保单连接，pg 零改动；前端 provider `fetchList` in-flight 合并 + 探活后写 + shell 级文档缓存）+ 26 页全盘扫描（system-monitoring 6×同 URL、data-display 3×同 URL 由全局机制覆盖；无第二例挂载即写组件）。S5 防复发（E-002）：store 连接面白盒回归测试（池/WAL/超时，防回退单连接）+ 渲染层合并/定向刷新回归 + schema 组件注册校验测试 + provider `refreshList` 定向刷新（monitoring 轮询 tick 由整页 `reloadList` 改为只刷 `/status`，9→3 请求/tick；事件表随手动刷新）+ `module-contribution-playbook.md` §6 页面数据面性能规范（1.1.0）。回归：go store 全绿 + build 0；vitest 定向全绿；全量回归于 E-003 补跑。**S6 进展**：**I-001 closed（2026-08-23）**——e2e 双 profile 全绿（admin 9/9 + mvp 9/9，另各 1 profile 专属跳过）；回归暴露后端缺陷「删用户遗留 `user_roles` 孤儿行 → 角色 `deletable=false` 永久化」（`DeleteUser`/`DeleteUsersBatch` 不清理关联；探针实证），已修复（同事务补删 `user_roles` + `user_mfa`）并加 2 项单元回归（E-004）。**I-002 closed（2026-08-23）**——双栈活栈实测（基线 `0878d7f` vs 当前，Playwright）：页面相关请求数 −47%～−86%（SPA 二次回访 14→2、schema 缓存命中 1→0）；呈现耗时本机 −17%～−51%、RTT150ms −25%（−1.4s/次）（E-005）。A-001（independent · conditional）响应完成：F-001～F-006 fixed（FK 每连接入 DSN + 多连接回归 + refreshList 对称 + 台账卫生 + 批删 MFA + 顺序断言）、F-007 fixed（E-007：测试替身 run id 时钟量化；产品 newRunID 防御加固）；**F-008（wallet reconcile 竞态，偶发 inconsistent）由下级 GOAL-037 承接根治**：机制 = 同毫秒流水 id 随机后缀乱序 → 回放先遇 freeze 未入账；修复 = 产品/替身 id 同毫秒单调计数 + **0050 数据修复迁移**（既库乱序重排，fail-closed）+ **成功审计原子化**（job 事务内 RecordOperationTx）+ 失败路径可观测（E-002/E-003，GOAL-037 done 4/4）。**GOAL-036 回归关门（用户书面约定）：A-003 self pass → done 6/6（2026-08-23）**。Root/VP 保持 active 程序容器。（承接 W19/GOAL-030 开通语义不变；与 workspace-009 正交。）

**W26（2026-08-26 关门，GOAL-038 done 4/4）**：用户点名三项符合性对齐——① 用户邮箱身份绑定管理端读面；② 发送邮件控制台与出站记录移出设置页为独立页面并注册左侧导航，出站记录覆盖全部渠道（含 mock），权限沿用 `settings.read` 不新设；③ 邀请管理「撤销」MISSING_PATH_BINDING 修复。S1 D-001 方案冻结（I-001～I-003 required 全 closed：0060 加列迁移 / admin.settings 两页 + menu_mail·menu_mail_outbox / ListUsers 同查询投影无 N+1）。S2 实施（E-002）：users 读面 email/emailStatus/emailStatusStyle + users.json 邮箱 badge 列与详情字段；`mail`（控制台）+ `mail-outbox`（声明式 table 六列 + recordView 含正文）两页贡献、设置页移除 tab-mail、DefaultNavigationOrder/BuiltinModules/快照测试 lockstep；Switcher 对 resend/smtp 落 sent/failed 记录、mock 单记 delivered；users-invites 撤销绑定补齐 + row-action-bindings 防复发登记。S3 回归（E-003）：Go 全量 0 FAIL + vitest 1116/1116 + tsc 0 + build ok；go 判定 = additive 产品面，**无影响不暂挂**；同日用户指示**补跑 e2e 双方言矩阵全绿（sqlite 9/9 + postgres 9/9，各 1 预期 skip；scratch 库闭环回收）**。S4 A-001 self 关门审计 **pass**（required 0；F-001/F-002 non-blocking 留痕）。Root/VP 保持 active 程序容器。

**W27（2026-08-26 立项并关门，GOAL-039 done 4/4）**：用户点名两页补强——邀请管理与邮件出站记录页面加上合理的筛选和排序。S1 D-001 冻结（无 required 信息项；白名单与默认序由既有约定唯一判定）：invites 后端 ListInvites 扩展 q/sort/order（LOWER+LIKE over email/id/invited_by，排序白名单 createdAt 默认/expiresAt × asc/desc，二级 id 稳定分页）；outbox 读面切换 `mail.OutboxListQuery` 契约（page/pageSize 归一化默认 50 上限 200，替换无消费方的 limit/offset；q/channel/delivery_status 筛选未知值 fallback-all；created_at × asc/desc）。S2 实施（E-001）：users-invites.json 搜索表单加 q + 两列 sortable；mail-outbox.json 插入搜索表单三控件 + created_at sortable + table.sort 能力；i18n 双目录新键。S3 回归：Go 全量 0 FAIL + vitest 1116/1116 + tsc 0 + build ok；go 判定 = additive 无影响不暂挂。S4 A-001 self **pass**（0 开放 required；status 列不可排序等两条 non-blocking 留痕）。Root 保持 active。

**W31（2026-09-18 立项并当日关门，GOAL-043 done 4/4）**：承接用户指令「能现在处理的直接处理掉……暂时不需要处理的确保在路线图中被正确统一登记」，承载 VP-037 关门后残余的**治理上下文**。**修复（均经变异验证）**：① 暗色下开关计算背景断言（4 条不变量；把类名改成 `bg-control dark:bg-[oklch(0.955_0_0)]` 时浅色断言全过、新断言失败）② 分页契约改 roles + users 双页面参数化 ③ Host 终态与普通 resource 反馈跨表直接对照（新增 `resource-feedback-parity.test.ts`，8 共有条件 × 分类一致/不回退通用文案/命名空间不混用/双语 key 齐备）。**收口**：`GOAL-008 A-002 F-002` → bounded residual（可执行面由守卫 + CI 锁死；文档侧 354 行形态不可唯一确定）。**愿景层**：`V-F124` 经 VRev-097 self `pass` 转 `fixed`；`VR-083` 记录 editorial 变更。**登记**：`roadmap.md` 新增「未决项统一登记」节（有界残余 / 悬置决策 / trigger-gated 能力三类表 + 维护约定）。回归：Vitest 114/1437 + typecheck 0 + e2e 4 passed × admin/mvp；无产品行为变更；A-001 self `pass`（0 required）。Root/VP-010 保持 active 程序容器；VP-037/workspace-037 不重开。

**W32（2026-09-19 立项并当日关门，GOAL-044 done 4/4）**：承接 `[workspace-038-batch-operations-and-job-center]` `GOAL-005`（R4 结果中心）cross 审计后留下的 3 条 low 级残余——用户 2026-09-19 指令「三条都修，先在 workspace-010 开承载子目标，038 侧有界接受后转由本区新子目标执行修正」。三项均以**跨页面通用能力**落地：① 列级 **`valueLabels`**（值 → i18n 键，fail-open 回落原始值；与 pinned `tagMap` 划界，不改 `docs/schemas/**`）+ jobs 六态接入；② **`refreshTable(tableId)` 定向刷新 seam**（用当前查询重取该表格且**保留选择**；`reloadList()` 的 ADR-0022 D2 语义以对照测试钉住未变；只读轮询不删 in-flight 键的理由见 `D-001` §2）；③ **空闲不轮询**（`activeStatuses` + 行注册表 `publishTableRows`/`tableRows`；行不可得时保守刷新）。**变异验证**：去 `valueLabels`、让 `refreshTable` 清空选择、禁用空闲判定——三处均实测变红后还原。回归：Go 全量 0 FAIL + vitest 119 files / 1468 tests + typecheck 0。审计 `A-001` self `pass`（0 required）+ `A-002` 三条 recommended 全 `fixed`；已按 P-003 回填 038 `A-003` 三条为 `fixed`，并在 roadmap「未决项统一登记」登记「本地扩展登记」项。Root/VP-010 保持 active 程序容器；不重开 VP-038/workspace-038。

**W33（2026-09-19 立项并当日关门，GOAL-045 done 4/4）**：用户在 VP-038 关门后提出两问——① roles 列表页（同样具备导出）缺「导出所选」；② 该按钮应位于**列表控件行**并**靠左**，而非筛选栏上方。选型（用户 P-004）：位置方案 **A**（表格 page-actions 行**左侧插槽**）、载体 **本区 W33**（VP-038 保持 `closed`）。交付：渲染器**本地扩展** `props.slot = "list-page-actions"` + `props.targetTable`——`data-list-page-actions` 行拆左右两段，左段为插槽宿主（经 CRUD seam 发布，`registerListActionsSlot` 决定该行是否渲染），未声明 slot 的 custom 节点走**逐字原路径**，目标表缺失或 slot 值未知时 **fail-open 原地渲染**（布局能力不得让操作入口静默消失）；`users` 页接入插槽；`roles` 页增 `selection.mode=multiple` + `roles-batch-export`（`resource: roles`；**后端零改动**，R3 冻结分母本就含 roles）。**实施中发现并修复一处既有隐患**：roles 页新增 custom 节点后两页 table 落在同一子节点索引，React 复用同一 `SchemaTable` 实例，把上一张表的 `visibleColumns` 渗入另一张（浏览器 e2e 捕获：users 往返后仅剩两 schema 交集列）→ 修复 = **在 SchemaTable 内按表格身份重置表级状态**（曾试「按节点 id 作 key」，因属无键子数组中的混键改动、收益不明确而撤回）+ 回归锁。**变异验证**：`roles.json` slot 拼写错 → 2 例红。**用户报告运行期控制台循环**（`Maximum update depth exceeded` 刷屏）→ 浏览器探针 + 插桩计数 + **W33 前基线对照**确认**先于本波次存在**（保存视图 page-actions 宿主 claim 效应把 `availabilityVersion` 当依赖、其 cleanup 又 `release()` 递增该值 → 自激循环）；修复 = 拆成「claim（只依赖宿主元素 + tableId + 重试计数）」+「重试（宿主空闲时 +1）」，并新增**永久守卫** `apps/web/e2e/console-health.spec.ts`（真实外壳逐页走查，任何 console.error/pageerror 即失败；变异验证：把旧依赖加回即让守卫失败）。**既有 flake 归因与加固**：`s5-denominator-render` 2 例在**干净树**同样失败（负载敏感），为 3 个真实 App 渲染用例补显式超时后，全套件在默认超时下稳定绿。回归：vitest **121 files / 1476 tests** + typecheck/build exit 0 + e2e mvp **17 passed / 5 skipped / 0 failed**、admin **18 passed / 4 skipped / 0 failed** + Go 相关包与 docscheck 全绿。审计 `A-001` self `pass`（0 required + 4 recommended）→ `A-002` 响应（3 fixed + 1 accepted-residual）。登记：`roadmap.md` §一 roles 触发面 `fixed`，「本地扩展清单」补入左侧插槽与取值约定。Root/VP-010 保持 active 程序容器；不重开 VP-038/workspace-038。

**W34（2026-09-19 立项并当日关门，GOAL-046 done 4/4）**：用户在 VP-038 关门后报告：「无论是用户列表页还是角色列表页，选中列表项并点击『导出所选』时都会报错宣示『未找到』，而并不会正常导出。」诊断出**两个独立缺陷**（均属已交付面的可用性缺口，非导出逻辑错误）：① **operator config 漂移（用户环境直接原因）**——`configs/config.yaml` 声明 `profile: custom` 且注释自称「= the full admin preset PLUS channel.telegram」，但 VP-038 将 `admin.jobs` 加入 admin preset 时漏改这份内联列表，模块未装配 → `POST /api/jobs/batch-export` 未挂载 → 路由回落 `{"error":"NOT_FOUND","messageKey":"error.notFound"}`（zh「未找到」），而节点属于 `admin.users`/`admin.roles` 在所有 profile 都照常渲染；② **触发面缺可用性门禁（mvp/demo 同类形态）**——这两个 preset 既无 `admin.jobs`（`jobs.write`）也无 `admin.data-transfer`（`data.export`），入口同样只有 404/403 一条路。修复：operator config 补 `admin.jobs`（含成因注释）；组件按**路由真实的两道门禁**判断可用性，不可用时渲染为**可见 + disabled + 说明文案**而非隐藏（依据 = 入口不撒谎 + 空插槽宿主的高度契约：隐藏会留下高度 0 的空宿主并破坏 page-actions 高度契约，该回归被 `list-visual-surface` e2e 实测捕获后在本波次修正；W33 `D-001` §3 的 fail-open 是**插槽布局**约定，不构成权限语义要求——据独立审计澄清）；提交命中 404 时改报「当前部署未启用批量导出（未装配 admin.jobs 模块）」，不再把服务端裸 `NOT_FOUND` 当成数据缺失。**守卫**：`internal/config` 的 `TestOperatorConfigCoversAdminPreset`（operator config 内联列表必须是 admin preset 的**超集**且含 `admin.jobs`）+ **对称判别性**可用性用例（先选行、断言 `unavailable` 标记）+ `apps/web/e2e/batch-export-availability.spec.ts`（**双 profile** 可用性契约：admin 可用 / mvp 可见但禁用；并同步修正 harness `customE2EModules` 同样漏掉 `admin.jobs` 的漂移）。**变异验证 6 处**：删/注释 `- admin.jobs` → config 守卫指名失败；判据只查 `jobs.write` / 只查 `data.export` → 对应用例各自变红（双向判别性）；去掉可用性判断、404 文案回退 → 对应用例红；e2e 去判据 → mvp 分支 `toBeDisabled` 失败。**真实端到端**：补齐 operator config 后 users/roles 均 `202 → succeeded` 且结果可下载（2155 字节 CSV），两权限齐备 → 入口在用户环境保持可用。回归：vitest **121 files / 1482 tests** + typecheck/build exit 0 + e2e mvp **18 passed / 5 skipped / 0 failed**、admin **19 passed / 4 skipped / 0 failed**（数字已含 always-on 的可用性 spec）+ `go test ./...` 全绿。**审计**：self `A-001` `pass` → **independent `A-002`（grok build · grok-4.6 · high · `/audit`）`conditional`**：独立复核确认根因、修复面未越界、config 守卫变异有效，但指出两处 required（XOR 单测不具判别性；可用性 e2e 当时**未入仓**且 harness `custom` 列表不一致）并判定 self 关闭过早 → `A-003` 响应两条均 `fixed`，并 append-only 更正 self 的三处高估（含 W33 §3 引用过度延伸），**开放 required = 0**。登记：`roadmap.md`「未决项统一登记」新增该缺陷行（`fixed`）。Root/VP-010 保持 active 程序容器；不重开 VP-038/workspace-038。
