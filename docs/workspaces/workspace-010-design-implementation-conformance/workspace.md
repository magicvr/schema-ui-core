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
updated: 2026-08-23
version: 0.47.0
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

Charter：`schema-ui-core-admin-foundation@0.2.0`。  
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


## 固定共享资料引用

> `shared-materials/index.json` 只能提供候选路径与摘要。缺完整引用字段的行无效。

| reference_id | workspace_id | material_id | source | version | sha256 | purpose | local_record | status |
|--------------|--------------|-------------|--------|---------|--------|---------|--------------|--------|
| — | — | — | — | — | — | — | — | — |

## 波次补充 · W22（2026-08-23）

GOAL-033-w22-residual-closeout done 18/18（accepted-residual 全库清点收口；详见该目标五件套与本区 goal-tree）。移交跟踪槽 **N-001**：admin 登录后停留 `/` 未跳 `/dashboard`（先于 W22 存在的既有回归，基线实验证实；疑似 W14–W21 home 推导/路由漂移）→ 建议下一符合性波次承接。*（归因更正：N-001 实为挂具 store 隔离失效，非路由回归——见下方 W23 记录与 GOAL-034 D-001；W22 基线实验无法移除 gitignored 文件。）*

**W23（2026-08-23 关门，GOAL-034 done 4/4）**：承接 W22 移交槽 N-001——admin 登录后停留 `/` 未跳 `/dashboard`。根因 = e2e 挂具 store 隔离失效（本机 gitignored `configs/.env` 2026-08-21 建，`DB_DIALECT=postgres` 劫持挂具临时 SQLite，全新种子 admin/admin 401），**非 home 推导/路由回归**；W22「先于 W22 基线实验」结论因 git stash 无法移除 gitignored 文件而失效。修复：`playwright.config.ts` 钉死 `DB_DIALECT=sqlite` + signInZh/sign-in fallback 等待硬化 + 连带 F-1（RowActionsMenu scroll-close 竞态 → 仅触发钮位移才关闭）F-2 fixed。回归：go 全包 ok / vitest 1088 / tsc+build 0 / e2e admin 连续 5 轮 9/9 + mvp 9/9；A-001 self pass（required 0）；I-001 closed。N-001 移交槽闭合。（2026-08-23 用户复审：钉方言属绕过，正确形态 = 双方言矩阵，承接 GOAL-035。）

**W24（2026-08-23 关门，GOAL-035 done 4/4）**：承接 GOAL-034 用户书面复审（「强制 sqlite 属绕过；收尾层 e2e 应双方言各测一次」；实验先证专用 pg 全量 9/9 绿）。实现：`DB_DIALECT` 方言契约（默认 sqlite / pg 显式 opt-in，`.env` 无法再改道）+ `apps/api/cmd/e2e-pgset`（scratch 库 create/verify/drop/list，凭据与 API/pgtest 同源）+ `globalSetup/Teardown` fail-fast 校验与清理 + `npm run test:e2e:postgres` + CI `profile×dialect [sqlite,postgres]` 矩阵。F-1（Playwright 配置双载 → 双份 scratch 库）修复：`E2E_PG_NAME` 守卫复用 + `DROP WITH (FORCE)` + teardown 可见。回归：sqlite 9/9 + postgres 9/9（遗留 0）+ vitest 1088 + go 全绿 + tsc/build 0；A-001 self pass（required 0）；I-001 closed。

**W25（2026-08-23 立项，GOAL-036 active 5/6，未闭门）**：我的钱包页面性能优化 → **用户升级为全盘修复此类问题 + 防复发**（D-002 书面裁决；文件夹改名 `GOAL-036-w25-page-performance-guardrails`）。S1–S4：四因素诊断（SQLite `MaxOpenConns=1` 全局串行 + 逐提交 fsync；同 URL 展示节点重复请求；wallet-ensure 挂载即写 + 整页重拉；schema 每次导航重取）+ 钱包页实施（后端文件库池 4 + `_busy_timeout=5000&_journal_mode=WAL&_synchronous=NORMAL`，`:memory:` 保单连接，pg 零改动；前端 provider `fetchList` in-flight 合并 + 探活后写 + shell 级文档缓存）+ 26 页全盘扫描（system-monitoring 6×同 URL、data-display 3×同 URL 由全局机制覆盖；无第二例挂载即写组件）。S5 防复发（E-002）：store 连接面白盒回归测试（池/WAL/超时，防回退单连接）+ 渲染层合并/定向刷新回归 + schema 组件注册校验测试 + provider `refreshList` 定向刷新（monitoring 轮询 tick 由整页 `reloadList` 改为只刷 `/status`，9→3 请求/tick；事件表随手动刷新）+ `module-contribution-playbook.md` §6 页面数据面性能规范（1.1.0）。回归：go store 全绿 + build 0；vitest 定向全绿；全量回归于 E-003 补跑。**S6 进展**：**I-001 closed（2026-08-23）**——e2e 双 profile 全绿（admin 9/9 + mvp 9/9，另各 1 profile 专属跳过）；回归暴露后端缺陷「删用户遗留 `user_roles` 孤儿行 → 角色 `deletable=false` 永久化」（`DeleteUser`/`DeleteUsersBatch` 不清理关联；探针实证），已修复（同事务补删 `user_roles` + `user_mfa`）并加 2 项单元回归（E-004）。**I-002 closed（2026-08-23）**——双栈活栈实测（基线 `0878d7f` vs 当前，Playwright）：页面相关请求数 −47%～−86%（SPA 二次回访 14→2、schema 缓存命中 1→0）；呈现耗时本机 −17%～−51%、RTT150ms −25%（−1.4s/次）（E-005）。自审 A-001 待办，不闭门。（承接 W19/GOAL-030 开通语义不变；与 workspace-009 正交。）
