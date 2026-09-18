---
doc_type: vision-plan
id: VP-037-admin-workflow-continuity
title: Admin 工作流连续性与安全反馈
status: active
vision_ref: schema-ui-core-admin-foundation@0.4.0
lead_workspace: workspace-037-admin-workflow-continuity
created: 2026-09-16
updated: 2026-09-18
parent: null
version: 1.2.0
---

# VP-037 · Admin 工作流连续性与安全反馈

## 意图

在现有 Admin 导航、发现入口和横切契约之上，补齐高频操作的连续性、安全反馈与通用列表体验：用户可以复用列表视图，离开有未保存修改的页面时得到可靠保护，在成功、失败、重试和恢复路径上获得一致反馈，并在列表页获得范例化的筛选、视图和分页布局。目标是降低重复配置、误操作和“操作到底有没有成功”的不确定性，而不是扩张为新的搜索或业务域平台。

本 VP 已于 2026-09-16 经用户确认从 `planned` 激活为 `active`，并绑定 `workspace-037-admin-workflow-continuity`；R1、R2、R3、R4 已分别完成信息/语义冻结、Saved View、dirty-state 与统一反馈实现及 Goal 审计，R5 当前承载组合验收与关门准备。

## 状态、激活与关门门禁

| 项 | 值 |
|-----|-----|
| status | **`active`**（2026-09-18 · v1.2.0；R1 C3、R2 C4、R3 C4、R4 C4 已完成；R5 `GOAL-006` 当前 `active · 3/4`；R6 `GOAL-007` 按用户选择 A 回开后 C5 实现/回归已完成、当前 `active · 5/6`；Root `active · 4/6`；lead `workspace-037-admin-workflow-continuity`） |
| 组合位置 | **Admin 功能分支 · 体验增强**；承接 VP-036 之后的工作流连续性下一拍 |
| Vision Review | 计划阶段 [VRev-094](../reviews/VRev-094-vp037-admin-workflow-continuity-planned.md) self `pass`；激活就绪 [VRev-095](../reviews/VRev-095-vp037-admin-workflow-continuity-activation.md) self `pass`；当前 open required = 0 |
| 激活门禁 | Admin 类 freshness PASS；I-037-006 verified；用户确认 workspace/Root 命名；V-F124 保持 recommended，不阻断激活 |
| 实现边界 | Root `GOAL-001-admin-workflow-continuity` active · 4/6；R1 已冻结、R2/R3/R4 已实现并通过阶段审计，R6 由 `GOAL-007-list-page-visual-alignment` 承载，原 C1～C4 已完成，C5 布局修订与回归已完成，C6 修订审计待完成 |

## 首波范围与边界

| 范围 | 首波决定 |
|------|----------|
| Saved Views | 面向现有已注册列表页，保存并恢复用户级筛选、排序、列配置等视图状态；具体字段与持久化形态由 R1 冻结 |
| 未保存变更保护 | 覆盖页面内导航、浏览器离开/刷新和提交/重置后的 dirty-state 语义；不替代业务表单自己的校验 |
| 统一 Toast / 错误恢复 | 收敛成功、失败、可重试、维护/不可用等反馈语义，并与既有 API 错误合同对齐 |
| 列表页视觉收敛 | 参考 `raw/new-table` 调整通用列表、筛选折叠、页面级按钮、视图位置/“全部{对象}”文案与始终可见分页；不改变 shell、查询/重置逻辑、Saved View 存储格式或未实装的多选能力 |
| 明确排除 | 实体全文检索、专用搜索引擎、批量结果中心、组织/部门/岗位与数据权限、新业务域、Redis/MQ/多实例 |
| 继承边界 | 复用 VP-034 导航分组、VP-036 发现入口、VP-005/007 体验基线、VP-012 横切契约；不重开既有 VP |

## 方向级退出判据

1. 已注册列表页的 Saved Views 分母、用户隔离、字段序列化、恢复/失效与权限边界形成可核验矩阵；不把共享视图或跨用户协作偷偷纳入首波。
2. Saved Views 在首波分母内完成保存、选择、恢复、更新、删除与空态/错误态闭环；无效或越权视图 fail closed，不污染当前页面状态。
3. dirty-state 状态机覆盖内部导航、浏览器刷新/关闭、提交成功、重置和取消等路径；误离开保护与确认文案可由浏览器/自动化回归验证。
4. 成功、失败、重试、维护/不可用反馈使用统一语义和可访问呈现；错误恢复不会重复提交、吞掉服务端错误或暴露敏感信息。
5. 首波不解除实体全文搜索、`RT-X01`/`RT-X02`、批量结果中心、组织/数据权限或架构分支 gated 项；边界证据与回归结果可追溯。
6. VP 工作区阶段链、Goal 审计、required 信息与 Vision Review 全部闭合，且关门前组合投影同步；开放 required finding = 0。
7. 通用列表页的视觉与交互基线可由共享实现覆盖，顶部功能栏/左侧导航不变，查询/重置合同不变，单页有效列表仍显示分页区域。

## 纲领路线图

| 阶段 | 目的 | 进入条件 / 产物 |
|------|------|-----------------|
| R1 · 分母与语义冻结 | 盘点列表页、筛选/排序/列配置、表单 dirty-state、反馈类型与权限/Profile 覆盖 | I-037-001～004 verified；矩阵与 D-003～D-005 取舍决策落盘；R1 Goal C3 self/independent 审计与响应完成（GOAL-002 `done · 3/3`） |
| R2 · Saved Views | 实现用户级保存/恢复与失效边界 | R1 冻结；由 `GOAL-003-r2-saved-views` 承载并已完成（`done · 4/4`）；A-001/A-002/A-003 pass，checkpoint `39c744ef` |
| R3 · 未保存保护 | 实现并验证 dirty-state 与离开确认 | R1 冻结；由 `GOAL-004-r3-unsaved-change-protection` 承载并已完成（`done · 4/4`）；A-003 independent recheck、A-004 self 与 checkpoint `d2b39189` 已记录 |
| R4 · 统一反馈与恢复 | 收敛 Toast、错误分类、重试/恢复和可访问状态 | R1 冻结；由 `GOAL-005-r4-unified-feedback-recovery` 承载并已完成（`done · 4/4`）；A-003 independent recheck、A-004 self 与 checkpoint `89666e5c` 已记录 |
| R5 · 组合验收与关门 | 核对非目标、审计链、残余与愿景投影 | R2～R4 已完成；由 `GOAL-006-r5-composition-acceptance` 承载（`active · 3/4`），C1～C3 已完成，继续核对退出判据 5～6，用户书面确认后才关 Root/VP |
| R6 · 列表页视觉与筛选体验收敛 | 参考范例页收敛通用列表、筛选折叠、页面级操作/视图布局、对象语义文案与分页展示 | 由 `GOAL-007-list-page-visual-alignment` 承载；原 C1～C4 已完成，用户选择 A 后 C5 布局修订与回归已完成，当前 `active · 5/6`，C6 修订审计待完成；不改变 R5-I-004 或既有查询/重置合同 |

## P-005 信息需求

| 编号 | 所需信息 | 级别 | 影响 | 最晚需要 | 收集/验证动作 | 状态 | 延期/复核与证据 |
|------|----------|------|------|----------|----------------|------|------------------|
| I-037-001 | 现有列表页、状态字段、Profile/权限覆盖的精确分母 | required | R1 范围冻结、R2/R5 验收 | R1 | 扫描页面/路由/列表注册表，形成机器可核对矩阵 | verified | 2026-09-17：`GOAL-002.../attachments/r1-denominator-matrix.json`、`r1-form-matrix.json` |
| I-037-002 | Saved View 所有权、持久化/序列化、权限变化后的失效语义 | required | R1 方案冻结、R2 实施/验收 | R1 | 用户确认方案 A；矩阵/D-003 冻结 localStorage 序列化 allowlist、Schema/权限失效与异常边界 | verified | 2026-09-17；R2 完成实现与测试后留存阶段证据；D-003 |
| I-037-003 | dirty-state 在内部路由、浏览器离开、提交、重置、取消中的统一语义 | required | R1 方案冻结、R3 实施/验收 | R1 | 盘点表单与路由守卫；矩阵/D-004 冻结状态机 | verified | 2026-09-17；R3 留存浏览器/自动化实现证据 |
| I-037-004 | Toast、API 错误、重试、维护/不可用反馈的分类与可访问呈现 | required | R1 方案冻结、R4 实施/验收 | R1 | 对照现有错误 envelope、反馈组件与 maintenance 门控；矩阵/D-005 冻结映射 | verified | 2026-09-17；R4 E-002/E-003 与 A-001 已补齐跨页面回归实现证据 |
| I-037-005 | 跨用户共享视图、最近使用/收藏与协作权限是否进入后续波次 | non-blocking | 后续 UX 波次边界 | 关门后或出现协作触发 | 不纳入首波；出现明确多用户协作需求时由 `/vision` 复核 | deferred | 延期理由：首波聚焦个人工作流；责任人：`/vision`；复核触发：真实协作需求出现 |
| I-037-006 | 激活前 Admin freshness 与 VP-008 `go` 消费有效性 | required | 激活与开区 | 激活前 | 执行 Admin 类 freshness review，并核对当前 Charter/VP 引用与区间变更 | verified | 2026-09-16：当前 HEAD `0c29c08`；`apps/**` 无 staged/unstaged 区间变更；VRev-095 |
| I-037-007 | R6 范例布局、通用列表调用链、shell 边界与查询/分页合同 | required | R6 方案与实施 | R6-C2 | 读取 raw 范例、盘点 renderer/components/app/test；对象语义实现由 GOAL-007 承接 | verified | 2026-09-18；D-001/E-001 已记录，C3 前复核对象文案 | `GOAL-007-list-page-visual-alignment/01-decision/D-001-list-page-visual-contract.md` |

`I-037-001`～`I-037-004` 的 R1 信息冻结已关闭；R2/R3/R4 已分别通过实现/回归证据关闭阶段门禁，R5 当前继续核对组合退出判据与关门投影。`I-037-006` 已 verified，不再阻断本次激活；`I-037-005` 是有界延期，不代表已验证或承诺后续实现。

## Workspaces

| workspace | role | scope | lead | 状态 |
|-----------|------|-------|------|------|
| workspace-037-admin-workflow-continuity | delivery | VP-037 首波与 R6 列表页视觉实现层范围 | workspace-037-admin-workflow-continuity | **active · Root 4/6**；R1 `done · 3/3`、R2 `done · 4/4`、R3 `done · 4/4`、R4 `done · 4/4`，R5 `active · 3/4`，R6 `active · 5/6`（C5 布局修订与回归已完成，C6 修订审计待完成） |

## 关系与结构选型

- 这是同一 Charter 下的下一波 Admin 体验意图，不是 Charter 修订。
- 当前唯一 delivery workspace 为 `workspace-037-admin-workflow-continuity`；Root 为 `GOAL-001-admin-workflow-continuity`，二者均只承接 VP-037。
- 选择新 VP，不塞入 VP-010：本 VP 新增用户工作流能力，不是 as-designed / as-built 符合性整改。
- 不以 VP-036 的已关闭分母重新定义“全局搜索”；Saved Views 等连续性能力单独冻结，避免把批量结果中心和实体搜索混入同一波。
- 不新建业务域，不解除架构分支的 Redis、MQ、多实例或搜索 gated 条件。

## 规划短史

- 2026-09-16：用户确认沿用总路线图推荐，建立 `VP-037-admin-workflow-continuity` `planned` v0.1.0；0 区；计划阶段 self Vision Review = VRev-094 `pass`。
- 2026-09-16：V-F124 作为非阻断 recommended 留给激活/R1：需把页面、状态、权限与持久化边界形成机器可核对矩阵。
- 2026-09-16：用户确认 workspace slug `workspace-037-admin-workflow-continuity` 与 Root slug `GOAL-001-admin-workflow-continuity`；VP-037 `planned → active` v0.2.0；Admin freshness PASS；VRev-095 self `pass`；交 `/govern` 建立 delivery workspace 与 Root。
- 2026-09-17：R2 `GOAL-003-r2-saved-views` 完成 C1～C4，A-001 self、A-002 independent、A-003 close-out 均 `pass`；commit `39c744ef`；Root 投影为 `active · 2/5`。
- 2026-09-17：按 R2 关门后的既定路线开设 R3 `GOAL-004-r3-unsaved-change-protection`，初始 `active · 0/4`，承接 D-004 dirty-state 合同与 `39c744ef` 基础切片。
- 2026-09-17：R3 完成 C1～C4；A-002 independent 的 F-001 required finding 经 E-003 修正并由 A-003 independent recheck 确认 `fixed`，A-004 self pass；Git checkpoint `d2b39189`；Root 投影为 `active · 3/5`。
- 2026-09-17：按既定路线开设 R4 `GOAL-005-r4-unified-feedback-recovery`，初始 `active · 0/4`；承接 R1 D-005，不引入新的 API 或 Host 终态语义。
- 2026-09-17：R4 完成 C1～C4；A-002 grok independent `conditional` 的 F-001 按 E-004 `fixed` 响应，由 A-003 independent recheck `pass` 确认；E-005 补 recommended 回归，A-004 self close-out `pass`；checkpoint `89666e5c`；R4 `done · 4/4`，Root 投影为 `active · 4/5`，下一阶段为 R5 组合验收。
- 2026-09-17：按 D-010 开设 R5 `GOAL-006-r5-composition-acceptance`（`active · 0/4`），承载 VP-037 组合验收、非目标/对齐核对、最终验证和 Root/VP 用户确认门禁。
- 2026-09-18：用户明确要求 Root 暂不关门并追加列表页视觉收敛；按 Root D-011 在同一 VP/workspace 下开设 R6 `GOAL-007-list-page-visual-alignment`，Root 路线图扩展为 6 个检查点（`active · 4/6`），R5-I-004 用户确认门禁保持开放；本次不改 Charter strategic 方向、vision_ref 或新建 VP/workspace。
- 2026-09-18：R6 `GOAL-007-list-page-visual-alignment` 原完成 C1～C4 与 A-001 self 审计，曾投影为 `done · 4/4`、Root `active · 5/6`；用户随后选择方案 A 回开 R6，追加 C5/C6 布局修订，当前 Root `active · 4/6`、R6 `active · 4/6`。R5-I-004 用户确认门禁仍开放，VP 继续保持 `active`。
- 2026-09-18：R6 C5 布局修订实现与回归完成，C6 修订审计待完成；Root 继续 `active · 4/6`，R5-I-004 用户确认门禁仍开放。

## Closeout placeholder

本节在进入关门审计前填写：方向级退出判据逐条证据、workspace/Root 状态、Goal 审计与独立意见（如风险要求）、required finding 闭合、用户书面确认、`roadmap.md` / `workspaces.md` / revisions / reviews 投影同步。当前不代表已完成。
