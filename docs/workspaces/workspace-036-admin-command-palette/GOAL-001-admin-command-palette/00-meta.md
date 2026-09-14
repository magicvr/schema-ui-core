---
id: GOAL-001-admin-command-palette
title: Admin 全局检索与 Command Palette 交付
status: done
parent: null
created: 2026-09-14
updated: 2026-09-14
version: 0.7.0
progress: 4/4
plan_refs:
  - VP-036-admin-command-palette
primary_plan: VP-036-admin-command-palette
vision_ref: schema-ui-core-admin-foundation@0.4.0
---

# GOAL-001 · Admin 全局检索与 Command Palette 交付

## 概述

在现有 Schema 驱动 Admin Shell、导航分组、模块贡献、权限与 Profile 体系之上，交付 VP-036 首波：对已注册页面、导航项和声明式动作提供权限安全的跨模块发现入口与 Command Palette。Root 只承接 VP-036 的实现层路线图，不把实体记录全文搜索或专用搜索基础设施写入本目标。

## 愿景对齐

- 工作区：`workspace-036-admin-command-palette`
- Charter：`schema-ui-core-admin-foundation@0.4.0`
- VP：`VP-036-admin-command-palette`
- `plan_refs` / `primary_plan`：均为 `VP-036-admin-command-palette`
- `serves_summary`：实现页面/导航/声明式动作的权限安全发现与键盘优先 Command Palette，复用已交付导航分组与模块贡献语义；不新增业务域或解除搜索基础设施 trigger。

## 范围与非目标

### 本目标范围

- 已注册模块、页面、导航项和声明式 UI 动作的可核对分母。
- `SearchableItem` / provider 契约、跨模块聚合、排序与去重。
- Profile/权限过滤、直接 URL 与既有路由守卫语义保持。
- Command Palette 键盘交互、焦点/ARIA、i18n/theme、加载/空态/错误态。
- 与 VP-034 导航分组和 top/user slot 的联动与回归。

### 明确非目标

- 业务实体全文索引、任意数据库表扫描、Saved Views 或最近/固定项持久化。
- 专用搜索引擎 `RT-X01`、DB 全文检索 `RT-X02`、Redis、MQ、跨进程索引、多实例。
- 新权限模型、组织权限、SSO/多租户、远程命令执行或既有守卫绕过。
- 新业务域、批量结果中心、未保存保护或 Toast 全局重做。

## 成功标准与纲领路线图

以下 4 个检查点构成 Root 的派生 progress 来源；纲领阶段按顺序推进。

- [x] **R1 范围与信息冻结**：页面/导航/动作分母、Profile/权限语义、排序/去重、快捷键与实体搜索排除有可核对决策与矩阵；self A-001 `pass`，grok independent A-002 `conditional` 的 F-001/F-002 已由 A-003 `fixed` 响应。
- [x] **R2 契约与聚合**：`SearchableItem` / provider v1 已实现；Manifest provider 与注入 provider 可稳定聚合、去重和排序，不增加 Shell 中央业务注册分支；A-004 self `pass`，完整 Web Vitest 1355/1355。
- [x] **R3 Palette 实现与体验**：Command Palette、键盘/ARIA、焦点、i18n/theme、直接 URL 与导航分组联动按既有语义可用；A-005 self `pass`，A-006 grok independent `pass`，A-007 已响应 4 条 recommended，mvp/admin SQLite/Postgres browser smoke 已通过。
- [x] **R4 回归与关门准备**：Profile×权限×路由矩阵、浏览器/自动化回归、边界复核、Goal 审计与必要的独立意见落盘；A-008 self `pass` + A-009 grok independent `pass`，A-010 已闭合 2 条 recommended（ID oracle 与 confirm 分支），四 Profile matrix 与 mvp/admin SQLite/Postgres smoke 已通过；开放 required = 0，等待用户确认 Root 关门。

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-036-001 | required | 当前所有可检索页面、导航项和声明式动作的精确分母与 Profile/optional module 覆盖是什么？ | R1 范围冻结、R2 聚合、R4 回归 | R1 | 扫描 provider、Manifest fragment、导航/动作注册，建立 item→module→profile 矩阵 | verified | R1 冻结前复核已完成 | `attachments/r1-searchable-item-matrix.md`；首波只承诺注册页面/导航/声明式动作，实体记录不进分母 |
| I-036-002 | required | Palette 结果的权限、Profile、直接 URL 与动作守卫如何保持现有语义？ | R1 方案冻结、R3 实施、R4 回归 | R1 | 对照权限过滤、路由守卫、Profile 装配与 action binding，建立允许/拒绝矩阵 | verified | R1 语义冻结；R2/R3 已复核实现，R4 继续回归 | D-002；不得新增权限绕过，programmatic gate 已由 A-004/A-005/A-006 及测试证据固定 |
| I-036-003 | required | 结果字段、标签本地化、排序/去重、快捷键、结果上限和焦点/ARIA 语义是什么？ | R1 方案冻结、R3 实施、R4 验收 | R1 | UX 方案 + 设计系统/locale 约定 + 浏览器可访问性验证 | verified | R1 口径冻结；R3 已有双语/主题/ARIA 与浏览器 smoke，R4 继续矩阵回归 | D-002；CommandPalette targeted/full tests、mvp/admin 双方言 smoke；不引入 recent/pinned |
| I-036-004 | required | 首波是否承诺实体级全局搜索，以及是否触发 `RT-X01` / `RT-X02`？ | VP 范围与基础设施门禁 | R1 | 用户书面确认的 VP-036 边界；新增实体搜索须另行 `/vision` 复核 | verified | 实体搜索需求或规模证据出现时复核 | 首波不承诺实体全文搜索；`RT-X01` / `RT-X02` 保持 gated；见 VRev-090/VRev-092 |
| I-036-005 | non-blocking | 最近搜索、固定项或持久化偏好是否进入首波？ | R3 体验扩展 | R3 | 实现前评估；如需要，另立 UX VP 或追加有界范围决策 | deferred | 不进入首波；下一 UX VP 规划时复核 | Saved Views / 最近项不作为本 VP 退出条件 |
| I-036-006 | required | 激活时当前代码候选是否仍满足 Admin 类 freshness 与 VP-008 `go` 消费有效性？ | 激活 | 激活前 | 复核协议 pin、依赖锁、迁移、Profile 默认集、provenance 与区间变更 | verified | 下次涉及 Admin 类基线/默认集/协议身份的区间变更时复核 | `5c341ec7`→`97aefe8c`；`apps/**` 无差异；五域 freshness PASS；见 VRev-092 |

## 父目标

- Root 目标，`parent: null`。

## 台账布局

本目标从第一条记录起使用平铺 ledger：`01-decision/`、`02-execution/`、`03-audit/`。索引文件保留 frontmatter、摘要与条目链接；独立记录使用 `D-NNN-*`、`E-NNN-*`、`A-NNN-*` 文件。

## 备注

- 工作区建立与 Root 设立是已发生事实；R1～R4 全部完成，Root 已由用户书面确认（2026-09-14）关门为 `done`；VP-036 `closed` 记录由 `/vision` 另行执行。
- `progress: 4/4` 只由上方 4 个显式检查点派生；`status: done` 由审计链（A-001～A-010，open required = 0）与用户确认放行。
