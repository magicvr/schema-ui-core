---
id: A-002-r1-independent-semantics-freeze
doc: audit-opinion
status: recorded
source: independent
verdict: pass
scope: R1 information readiness and semantic freeze
goal_id: GOAL-002-r1-scope-semantics-freeze
auditor: grok-build (grok-4.6 · reasoning high)
created: 2026-09-17
updated: 2026-09-17
parent: GOAL-002-r1-scope-semantics-freeze
version: 0.1.0
---

# A-002 · R1 独立审计 · 列表分母与工作流语义冻结（2026-09-17）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型** / **scope**：design-plan / information-readiness；R1 C1～C2 信息冻结与语义决策，不含 R2～R4 实现验收、不含本目标关门
- **verdict**：pass
- **工作区**：`workspace-037-admin-workflow-continuity`（`root_goal`: `GOAL-001-admin-workflow-continuity`；`canonical_scope` 已核对；`shared_materials_catalog: none`）

## 范围与区间

独立核对对象仅为 `GOAL-002-r1-scope-semantics-freeze`。核对项：

1. I-037-001～004 是否有证据支持 **R1 信息冻结**（不是 R2～R4 阶段完成）。
2. 用户选择的 localStorage 方案 A 是否被 D-003 正确限定。
3. D-004 / D-005 是否清晰覆盖 dirty-state 与反馈/恢复边界。
4. 当前未提交的 Saved Views / dirty-state 实现是否被正确标注为后续 R2/R3 证据，而非 R1 完成证据。
5. 执行台账是否有断链或矛盾。
6. 是否存在 required / 必改 finding。

核对方式：只读现有五件套、ledger、附件，并用现有命令对照 `apps/api/modules/**/schema/*.json`、`apps/web/src/app/searchable-profile-matrix.test.ts`、当前工作树与 `git` 状态。未写入临时脚本或其他文件；未改 status / progress / 方案正文。

## 成果（有证据）

### I-037-001 · 分母矩阵可独立复算

- 附件 `attachments/r1-denominator-matrix.json`、`attachments/r1-form-matrix.json` 可解析。
- 对 `apps/api/modules/**/schema/*.json` 递归盘点：`type: table` **24** 个、`type: form` **58** 个；表/表单 id 集合与附件一致。
- 24 张表的 `columns`、`sortableFields` 与 Schema `props.columns` 一致；58 个表单的 `mode` / `targetTable` / `fieldIds` 与 Schema 一致。
- 4 个 profile 的 discoverable `page:` / `action:` 集合与 `apps/web/src/app/searchable-profile-matrix.test.ts` 的 oracle 一致（mvp 5/6、admin 18/16、demo 13/6、custom 22/18）。
- 因此 I-037-001 的 R1 信息冻结（列表/表单分母与 Profile 可发现覆盖）有可重复核对的证据。

### I-037-002 / D-003 · 方案 A 已被正确限定

- E-003 记录用户于 2026-09-17 选择 D-002 方案 A；D-003 `accepted` 承接该选择。
- D-003 将方案 A 限定为：浏览器 `localStorage`；键按 `user.id + pageId + tableId` 隔离；allowlist 为 `q/filters/sort/order/pageSize` 与列可见性；`activeViewId` 只作 UI 恢复指针、不进资源请求；不保存当前页/选择集/record view/modal 草稿；同浏览器同设备；malformed/越权/Schema 不匹配 fail closed；存储异常走统一反馈；不新增 API/表/迁移。
- 未选 B/C 已写明。R2 实现证据被明确留在后续阶段。I-037-002 作为 **R1 信息项** 可视为 verified；**不**等于 Saved View 退出判据已满足。

### I-037-003 / D-004 · dirty-state 边界已冻结

D-004 覆盖：查询/分页非业务 dirty；default-mode 以挂载快照为 baseline；App 内部导航（含面包屑/菜单）dirty 先确认；`beforeunload` 用原生合同；`popstate` 取消恢复已提交 URL、确认提交目标；成功提交清 dirty、失败保留；dirty modal 关闭/取消同一确认语义。矩阵对应行已从 `proposed` 标为 `已冻结（D-004）`。R3 浏览器/自动化证据仍待后续阶段。

### I-037-004 / D-005 · 反馈/恢复边界已冻结

D-005 覆盖：成功 `role=status` / 错误 `role=alert`；幂等读允许显式 retry、写失败不自动 retry；catalog 文案优先且禁止回显敏感载荷；401 Host / 403 fail closed；维护/不可用/离线/超时为可恢复读取失败；Saved View 失效丢弃且不污染当前 query；toast/列表错误/retry/确认须可键盘到达。R4 跨页面回归证据仍待后续阶段。

### 未提交实现未被当作 R1 完成证据

当前 HEAD 为 `92cf8582`。工作树含未提交 Saved Views / dirty-state 切片：

- 已修改：`apps/web/src/app/App.tsx`、`apps/web/src/renderer/render.tsx`、`apps/web/src/renderer/schema-table.tsx`、i18n
- 未跟踪：`apps/web/src/renderer/saved-views.ts`、`saved-views.test.ts`、`saved-views.ui.test.tsx`、`dirty-state.ts`、`dirty-state.test.ts`

`02-execution.md` 当前事实与 A-001 均写明：这些是语义记录之后的未提交实现，**不作为 R1 C1/C2 完成依据**，分别由 R2/R3 承接。本独立审计同意该标注，也未把这些文件或任何测试输出当作 R1 完成证据。未复跑测试套件。

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| C1 分母矩阵（24 列表表面、58 表单） | 达成（R1 信息） | 附件 + 对 Schema 的独立复算；见 F-002/F-003 的非阻断精度缺口 |
| C2 Saved View / dirty-state / 反馈语义已确认 | 达成（R1 信息） | D-003、D-004、D-005 + 矩阵冻结列 + E-003/E-004 |
| C3 I-037-001～004 关闭、self+independent、Root 投影 | 未完成 | A-001 self `pass`；本条 independent `pass`；Root 投影与 C3 勾选仍须 `/govern` 响应；本意见不改 status/progress |
| R2 Saved Views 实现/回归 | 未开始（本 scope 外） | 未提交代码不得充当本阶段完成证据 |
| R3 dirty-state 自动化证据 | 未开始（本 scope 外） | 同上 |
| R4 跨页面反馈回归 | 未开始（本 scope 外） | D-005 只冻结语义 |

## Findings

### F-001 · 执行台账存在过期句子，与后续 verified 记录并列

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：权威信息登记（`00-meta.md` / `01-decision.md`）与 E-003/E-004 将 I-037-001～004 标为 R1 `verified`，但较早执行条目未改写完毕。
- 证据：
  - `02-execution/E-002-r1-matrix-record.md` 仍写「I-037-002 现为 `collecting`，等待 R2 实现证据」。
  - `02-execution/E-001-initial-r1-scan.md`（`updated: 2026-09-17`）仍写「I-037-001～004 尚未全部关闭」。
  - `02-execution.md` 仍写「当前 Git HEAD 为 `1e823416`」；独立核对时 HEAD 已是 `92cf8582`（`1e823416` 的后代；`1e823416..HEAD` 的 `apps/**` 提交为空，工作树另有未提交实现）。
  - 决策索引将 D-002 列为 `proposed`，文件 frontmatter 为 `draft`。
- 影响：读者若只看 E-001/E-002，会把已冻结的 R1 信息项误判为仍 collecting，或把 R2 实现证据缺口当成 R1 未冻结。不推翻 D-003～D-005 与 00-meta 的冻结结论。
- 建议：由 `/govern` 在响应本意见时更正过期句子，区分「R1 信息 verified」与「R2～R4 实现证据未到」；不要用 R2 未完成把 I-037-002 退回 collecting。

### F-002 · custom profile 隐藏路由集合漏列 `telegram-operator`

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：`r1-denominator-matrix.json` 的 custom `registeredPageIds` 含 `telegram-operator`，`discoverablePageIds` 不含，但 `hiddenRoutePageIds` 只有 `users-invites` / `dictionary-entries` / `task-runs` / `wallet-entries`。`registered \ discoverable` 实际还有 `telegram-operator`。该页是 custom 组件、无 table/form，不影响 24/58 分母。
- 证据：矩阵 `profileDenominator.custom`；`apps/api/modules/channel/telegram/schema/telegram-operator.json`；`searchable-profile-matrix.test.ts` 的 `OPTIONAL_PAGES` 含该页且 discoverable oracle 不含 `page:telegram-operator`。
- 影响：Profile 隐藏路由分类不自洽；不改变列表/表单分母，也不构成 Saved View 漏表。
- 建议：R2 使用 profile 覆盖作 oracle 前，把 `telegram-operator` 补进 custom `hiddenRoutePageIds` 并修正 `derivedCounts.hiddenRoutePages`。

### F-003 · `data-permission/policies` 的 `tableFilterFields` 与 Schema 不符

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：矩阵将该表 `tableFilterFields` 记为 `["resource"]`，但 Schema 无 `props.filters`；`resource` 是 `rowKey`。其余 23 张表的 columns / sortable / table filters 与 Schema 一致。
- 证据：`attachments/r1-denominator-matrix.json` 的 `listSurfaces`；`apps/api/modules/datapermission/schema/data-permission.json`。
- 影响：若 R2 盲信该字段作为 allowlist，可能持久化不存在的 table filter。未提交的 renderer 实现是从活 Schema 收集 filter fields，不把本矩阵当 runtime 权威。
- 建议：R2 以 Schema/renderer 活配置为序列化 allowlist 来源；矩阵该字段改为 `[]` 或注明派生规则。

### F-004 · 自定义列表型表面未写入 24 张表分母，需在 R2 明确排除

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：`notifications`（custom `notification-center`，search form 的 `targetTable` 为 `notifications-table` 但无 `type: table`）、`mail`（`mail-admin-tab`）、`telegram-operator` 不是 Schema table。按「已注册 table 节点」口径排除是对的，但矩阵未显式写「自定义列表表面不在 Saved View 首波分母」。
- 证据：`apps/api/modules/notifications/schema/notifications.json`；`apps/api/modules/settings/schema/mail.json`；telegram-operator schema。
- 影响：R2 若把「已注册列表页」扩到自定义表面，会超出 R1 冻结的 24 张表。
- 建议：R2 方案冻结时写明 Saved View 只覆盖矩阵 24 个 `type: table` 表面。

无其他 finding。409 的「reload/再编辑入口」写在矩阵 D-005 行，决策正文用写失败不自动 retry 概括；视为 R4 实现细化，不单列。

## 必改项汇总

无 required / 必改 finding。本意见 **不** 将上述 recommended 项接受为 residual，也 **不** overrule。

I-037-001～004 作为 R1 最晚阶段的 required 信息项：有矩阵 + 用户选择 + D-003～D-005 证据，足以支持 **R1 信息冻结**。它们同时影响 R2/R3/R4 实施/验收，那些门禁必须分别用阶段实现证据关闭，不能用本次 verified 代替。

## 与既有意见的异同

| 项 | A-001 self | 本条 independent |
|----|------------|-----------------|
| verdict | pass | pass |
| I-037-001～004 R1 冻结 | 同意 | 同意，并独立复算 24/58 与 profile oracle |
| 未提交实现不当作 R1 完成 | 同意 | 同意；补充当前 HEAD `92cf8582` 与未提交路径清单 |
| required finding | 无 | 无 |
| 额外意见 | 未列矩阵/台账精度缺口 | F-001～F-004 recommended；不阻断 R1 信息冻结 |

无结论冲突，无需 P-004 裁决。

## 结论 + 建议给编排器/用户的下一步

R1 信息就绪与语义冻结独立结论为 **pass**。方案 A 已被 D-003 正确限定；D-004/D-005 足以作为 R3/R4 的实现合同；当前未提交 Saved Views/dirty-state 代码被正确降为后续阶段证据。

**R2～R4 仍需各自阶段实现证据**：序列化/失效/异常路径（R2）、导航/beforeunload/popstate 自动化（R3）、跨页面反馈回归（R4）。不得因本 `pass` 勾选 Root R2～R4 或宣称 Saved View / dirty guard / 统一反馈已完成。

建议 `/govern`：

1. 响应 A-001 + A-002；不关闭 recommended finding 也可推进 C3，因无开放必改项。
2. 若勾选 C3，须把 Root R1 检查点投影为完成，并处理 Root `01-decision.md` 仍把 I-037-002～004 写成 `collecting`、Root D-003 仍写「I-037-002 保持 collecting」的投影不一致（Root 文件不在本意见写入范围）。
3. 按 F-001 清理 GOAL-002 执行台账过期句子；F-002～F-004 留给 R2 方案/矩阵修订。
4. 在开设 R2/R3/R4 子目标前，保持本目标 C3 与 Root R1 投影的先决顺序。

## 声明

本意见 `source: independent`，不修改 status / progress / 检查点 / goal-tree / 方案正文；不关闭 finding；不接受 residual；不 overrule。响应由 `/govern` 处理。
