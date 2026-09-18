---
id: GOAL-001-admin-workflow-continuity
title: Admin 工作流连续性与安全反馈交付
status: active
parent: null
created: 2026-09-16
updated: 2026-09-18
version: 1.7.0
progress: 5/6
plan_refs:
  - VP-037-admin-workflow-continuity
primary_plan: VP-037-admin-workflow-continuity
vision_ref: schema-ui-core-admin-foundation@0.4.0
---

# GOAL-001 · Admin 工作流连续性与安全反馈交付

## 概述

在现有 Admin 导航、发现入口、设计系统、locale/settings 与横切契约之上，交付 VP-037 的工作流连续性与列表体验收敛：用户级 Saved Views、未保存变更保护、统一 Toast/错误恢复，以及用户明确追加的通用列表页视觉与筛选布局改造。Root 只承接 VP-037 的实现层路线图，不把实体全文搜索、批量结果中心、组织/数据权限、新业务域或架构 gated 项写入本目标。

工作区与 Root 已建立；R1 信息与语义冻结、R2 Saved Views、R3 未保存变更保护、R4 统一反馈与恢复以及 R6 列表页视觉收敛已完成并通过各自 Goal 审计。R5 组合验收仍有用户书面确认门禁；整改子目标 `GOAL-008-typecheck-evidence-convention`（非纲领阶段）已 `done · 4/4`，使 R6 的 A-002 F-005 跨区部分闭环。当前 `progress: 5/6` 是显式检查点的派生展示。Root/VP 继续保持 active。

## 愿景对齐

- 工作区：`workspace-037-admin-workflow-continuity`
- Charter：`schema-ui-core-admin-foundation@0.4.0`
- VP：`VP-037-admin-workflow-continuity`
- `plan_refs` / `primary_plan`：均为 `VP-037-admin-workflow-continuity`
- `serves_summary`：实现个人工作流视图复用、离开前安全保护与一致成功/失败反馈，保持用户隔离与既有 API/权限语义；不新增业务域或解除搜索/Redis/MQ/多实例 gated 条件。

## 范围与非目标

### 本目标范围

- 现有已注册列表页的 Saved Views 分母、用户所有权、筛选/排序/列配置状态与恢复/失效语义。
- 页面内导航、浏览器刷新/关闭/离开与提交/重置/取消相关的 dirty-state 与确认路径。
- 成功、失败、重试、维护/不可用反馈的统一 Toast、可访问呈现与恢复语义。
- 已注册通用列表页的范例化视觉、筛选折叠、页面级操作布局、视图位置/语义文案与始终可见分页。
- Profile/权限边界、API 错误合同、maintenance 门控与既有导航/发现入口的回归。

### 明确非目标

- 实体全文索引、专用搜索引擎 `RT-X01` / `RT-X02`、跨进程索引。
- 批量结果中心、组织/部门/岗位、数据权限 `org`、新业务域。
- Redis、MQ、多实例、第二持久化栈或 API 错误合同重写。
- 跨用户共享视图、协作权限、最近/收藏功能；如出现真实协作触发，另行 `/vision` 复核。

## 成功标准与纲领路线图

以下 6 个检查点构成 Root 的派生 progress 来源；纲领阶段按 R1 →（R2/R3/R4 可在 R1 后并行）→ R5 → R6 推进。用户明确要求在 R5 用户确认门禁仍开放时先开设 R6；这不代表 R5 或 Root/VP 已完成。

- [x] **R1 分母与语义冻结**：列表页、Saved View 所有权/持久化、dirty-state、反馈分类与 Profile/权限覆盖形成可核对矩阵；由 `GOAL-002-r1-scope-semantics-freeze` 承载并以 `done · 3/3` 完成，I-037-001～004 verified，A-001 self + A-002 independent pass，A-003 已响应。
- [x] **R2 Saved Views**：在首波分母内完成保存、选择、恢复、更新、删除、空态/错误态与无效/越权 fail-closed 闭环；由 `GOAL-003-r2-saved-views` 承载并以 `done · 4/4` 完成，A-001/A-002/A-003 pass，Git checkpoint `39c744ef` 已记录。
- [x] **R3 未保存变更保护**：内部导航、浏览器离开/刷新、提交成功、重置和取消路径可验证，确认不会丢失修改或绕过提交结果；由 `GOAL-004-r3-unsaved-change-protection` 承载并以 `done · 4/4` 完成，A-003 independent recheck、A-004 self 与 checkpoint `d2b39189` 已记录。
- [x] **R4 统一反馈与恢复**：成功/失败/重试/维护反馈统一、可访问且不重复提交、不吞服务端错误、不泄露敏感信息；由 `GOAL-005-r4-unified-feedback-recovery` 承载并以 `done · 4/4` 完成，A-003 independent recheck、A-004 self 与 checkpoint `89666e5c` 已记录。
- [ ] **R5 组合验收与关门准备**：非目标边界、阶段事实、Goal 审计、必要独立意见与 VP 投影闭合；当前由 `GOAL-006-r5-composition-acceptance` 承载（`active · 3/4`，C1～C3 已完成），用户确认后才可将 Root/VP 关门。
- [x] **R6 列表页视觉与筛选体验收敛**：参考 `raw/new-table` 收敛通用列表、筛选、页面级按钮、视图布局与分页展示；由 `GOAL-007-list-page-visual-alignment` 承载并以 `done · 8/8` 完成（C1～C8 全部完成；C6 审计 A-002 `conditional`，其 F-005 跨区部分经整改子目标 `GOAL-008` 闭环）；不改变 shell 与既有查询/重置合同。

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-037-001 | required | 当前列表页、状态字段、Profile/权限覆盖的精确分母是什么？ | R1 范围冻结、R2/R5 验收 | R1 | 扫描页面/路由/列表注册表，建立机器可核对矩阵 | verified | 2026-09-17 已完成矩阵核对 | `GOAL-002.../attachments/r1-denominator-matrix.json`、`r1-form-matrix.json` |
| I-037-002 | required | Saved View 所有权、持久化/序列化、权限变化后的失效语义是什么？ | R1 方案冻结、R2 实施/验收 | R1 | 用户确认 localStorage 方案 A；矩阵与 D-003 冻结 allowlist/失效边界 | verified | 2026-09-17；R2 已复核实现/回归证据 | `GOAL-002.../01-decision/D-003-saved-view-localstorage-accepted.md` + GOAL-003 A-003 |
| I-037-003 | required | dirty-state 在内部路由、浏览器离开、提交、重置、取消中的统一语义是什么？ | R1 方案冻结、R3 实施/验收 | R1 | 盘点表单与路由守卫，D-004 冻结状态机 | verified | 2026-09-17；R3 已完成浏览器/自动化证据 | `GOAL-002.../01-decision/D-004-workflow-semantics-frozen.md` + GOAL-004 A-004 |
| I-037-004 | required | Toast、API 错误、重试、维护/不可用反馈如何分类并可访问呈现？ | R1 方案冻结、R4 实施/验收 | R1 | 对照现有错误 envelope 与 FeedbackRegion，D-005 冻结映射 | verified | 2026-09-17；R4 E-002～E-006、A-003/A-004 已闭合阶段实现/回归与审计证据 | `GOAL-002.../01-decision/D-005-feedback-recovery-semantics-frozen.md` + GOAL-005 R4 `03-audit` |
| I-037-005 | non-blocking | 跨用户共享视图、最近使用/收藏与协作权限是否进入后续波次？ | 后续 UX 波次边界 | 关门后或出现协作触发 | 不纳入首波；真实协作需求出现时由 `/vision` 复核 | deferred | 理由：首波聚焦个人工作流；责任人：`/vision`；复核触发：真实协作需求 | 待确认 |
| I-037-006 | required | 激活前 Admin freshness 与 VP-008 `go` 消费有效性是否仍成立？ | 激活与开区 | 激活前 | 执行 Admin 类 freshness review，并核对当前 Charter/VP 引用与区间变更 | verified | 2026-09-16 已核对 | VRev-095；激活基线 HEAD `0c29c08`，后续治理提交 `d9440e12`、`1e823416` 未修改 `apps/**` |
| I-037-007 | required | R6 的范例布局、通用列表调用链、shell 边界、查询/重置与分页合同是什么？ | R6 方案与实施 | R6-C2 | 读取 `raw/new-table`，盘点 `App.tsx`/renderer/components/test 并记录对象文案未知 | verified | 2026-09-18 已完成基线；对象名实现验证由 GOAL-007 I-007-004 承接 | `GOAL-007-list-page-visual-alignment/01-decision/D-001-list-page-visual-contract.md`、E-001 |

I-037-001～004 的 R1 信息冻结已关闭；R2、R3、R4 已分别以实现/回归与 Goal 审计证据关闭阶段门禁。I-037-005 是有界延期，不代表已验证或承诺后续实现；I-037-006 已 verified，不再阻断本次激活与开区；I-037-007 已由 R6 基线盘点验证，C5/C7/C8 继续沿用该信息合同（C7/C8 未新增未知项）。R5 仍需组合验收、Root 审计与用户确认后才能关门，R6 不替代该门禁。

## 父目标

- Root 目标，`parent: null`。

## 台账布局

本目标从第一条记录起使用平铺 ledger：`01-decision/`、`02-execution/`、`03-audit/`。索引文件保留 frontmatter、摘要与条目链接；独立记录使用 `D-NNN-*`、`E-NNN-*`、`A-NNN-*` 文件。

## 备注

- workspace/Root scaffold、R1、R2、R3、R4 与 R6 子目标已完成是已发生事实；R5 尚未完成。
- `progress: 5/6` 只由上方 6 个显式检查点派生；整改子目标 `GOAL-008` 为非纲领阶段、不计入分母，也不放行阶段、不关闭 finding、不覆盖 status。R5 的 `R5-I-004` 用户书面确认未闭合，故 Root/VP 仍为 `active`。
- Vision Review `VRev-094`/`VRev-095` 属愿景层；Goal 审计须写入本目标 `03-audit/`，不能用 Vision Review 代替。
