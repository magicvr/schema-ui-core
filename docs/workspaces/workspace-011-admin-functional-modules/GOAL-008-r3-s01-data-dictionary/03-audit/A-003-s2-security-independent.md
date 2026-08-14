---
id: A-003
goal: GOAL-008-r3-s01-data-dictionary
source: independent
date: 2026-08-14
scope: S-01 实现安全/数据门禁（admin.data-dictionary vs D-002 冻结方案）
verdict: conditional
auditor: grok-build
audit_type: execution-facts
status: recorded
parent: GOAL-008-r3-s01-data-dictionary
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# A-003 · independent 安全/数据审计（S-01 实现）

## 范围与区间

- **auditor**：grok-build（independent cross-audit）
- **type**：execution-facts / security-data gate
- **workspace**：`workspace-011-admin-functional-modules`（`root_goal` = `GOAL-001-admin-functional-modules`；`canonical_scope` 已核对；`shared_materials_catalog: none`）
- **covered**：
  - `apps/api/internal/modules/datadictionary/store/repository.go`
  - `apps/api/internal/handler/dictionary.go`（及工厂 `resources.go` 门禁/PATCH 解码/排序白名单）
  - `apps/api/internal/modules/datadictionary/`（provider、schema ×2、manifest fragment）
  - `apps/api/internal/modules/datadictionary/migration/migration.go`（0019）
  - `apps/api/internal/modules/operationlog/migration/migration.go`（0020）
  - `apps/api/internal/kernel/profile.go`
  - 计划契约 `01-decision/D-002-s1-plan-freeze.md`（及 D-001/D-003 Profile 声明）
  - 页面文档对照 `docs/schemas/page.schema.json` / `action.schema.json` / form-control 白名单与 L2 能力门
- **excluded**：端到端浏览器手测、生产部署、其它工作区上下文、非本模块的资源工厂实现变更
- **信息项**：I-001 / I-002 均已 closed（D-002 / D-001）；本 scope 无到期未关闭 required 信息项；无共享资料引用

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| 读/写门禁 fail-closed（anon 401 / 无 key 403） | 两资源均挂 `dictionary.read` / `dictionary.write`（dictionary.go:228–229、243–244）；工厂 `requirePermission`（resources.go:227–241）匿名 401、缺权限 403；全部 CRUD + batch-delete 经该包装。`TestDictionaryPermissionGates`；provider 端到端匿名 GET 401 |
| 权限贡献 admin-only | provider.go:83–85 `PolicyID: PolicyAdmin`；导航 `Permission: dictionary.read`（provider.go:91–99） |
| UNIQUE(dict_key, entry_key) | 0019 DDL（migration.go:30–40）；CreateEntry/UpdateEntry 将 unique 错误映射为 `ErrEntryKeyTaken` → 409 `DICT_ENTRY_KEY_TAKEN`（repository.go:288–290、313–315；dictionary.go:274–275） |
| PATCH 缺字段不覆盖存值 | 工厂 `decodeResourcePatch` 缺 key 不写入 body（resources.go:415–462）；类型/条目 Update 对 name/description/dictKey/label/sort/remark/enabled 做 present-or-keep（dictionary.go:85–103、173–195） |
| 类型删除级联条目 | 0019 `REFERENCES dict_types(key) ON DELETE CASCADE`（migration.go:32）；store 启动 `PRAGMA foreign_keys = ON` 且断言（store/migrate.go:213–224）；`DeleteType` 仅删类型行（repository.go:184–196）；`TestDictionaryLifecycle` 断言级联后 entries total=0 |
| 未知 dict_key 拒绝 | CreateEntry/UpdateEntry 先 `COUNT(*) FROM dict_types WHERE key = ?`（repository.go:276–281、302–308）；`ErrDictKeyNotFound` → 400 `DICT_KEY_NOT_FOUND`（dictionary.go:276–277）；生命周期测试覆盖 create 路径 |
| SQL 参数化 + 排序白名单 | 用户输入均 `?` 绑定；`sortCol` 仅来自白名单 map，未知回落固定列；`order` 仅 `desc` 否则 `asc`（repository.go:90–101、214–225）。工厂另将 sort/order 校验为 400（resources.go:291–313） |
| 审计事件与 0020 CHECK 一致 | 常量 `dictionary.create/update/delete`（operationlog/repository.go:39–41）；handler 只发射这三项（dictionary.go:77/107/120/165/199/212）；0020 CHECK 在 0018 超集上仅追加这三项（migration.go:86）；写失败 `slog.Error`（dictionary.go:297–307） |
| 0020 重建保行 + 1..20 连续 | `rebuildOperationLog` rename/copy/drop/index（migration.go:200–219）；版本 1–2/9/11/12 auth、3/6 corepersistence、4/5/8/14/15/18/20 operationlog、7/10 settings、13 account、16/17 notifications、19 dictionary，无空洞 |
| Profile 仅为声明的内容扩展 | `profileDefaults[ProfileAdmin]` 追加 `"admin.data-dictionary"` 并注释 D-001（profile.go:68–70）；`ResolveProfile` 逻辑未改；mvp/demo 未加入；D-001 §2 / D-003 声明一致 |
| NavigateAction 文档形状合法 | `openEntries` 为 `{type:navigate,url:/dictionary-entries}`（data-dictionary.json:157–159），符合 `action.schema.json` NavigateAction |

## 对照成功标准（D-002 安全/数据相关）

| 标准 | 结论 |
|------|------|
| dictionary.read/write 门禁 401/403 | **满足** |
| UNIQUE(dict_key, entry_key) | **满足** |
| PATCH 缺字段不擦除 | **满足**（API 层） |
| 类型删除级联条目 | **满足** |
| 未知 dict_key → 400 | **满足** |
| SQL 参数化 + 排序白名单 | **满足** |
| 审计事件 ↔ 0020 CHECK；写失败留痕 | **满足** |
| 0019 建表 + 0020 保行 + 版本 1..20 | **满足** |
| Profile 仅内容扩展 | **满足** |
| 页面 schema / 表单控件 / 无新 renderer 扩展 | **不满足**（F-001、F-002） |

## Findings

### F-001 · 表单控件使用未登记类型 `number`，且 `defaultValue` 未声明 `form.controls.advanced`

| 字段 | 值 |
|------|-----|
| level | required |
| status | open |
| evidence | `schema/data-dictionary.json:97,145` 与 `schema/dictionary-entries.json:99,153`：`type: "number"`；同文件 `enabled`/`sort` 带 `defaultValue`（data-dictionary.json:85,98,133,146；dictionary-entries.json:93,100,147,154）；两页 `requiredCapabilities` 仅有 `form.controls.extended`，无 `form.controls.advanced`（data-dictionary.json:6–14；dictionary-entries.json:6–14） |
| severity | med–high（协议门禁；Host 对非法控件 fail-closed，schema 驱动写 UI 不可提交） |

**说明**：冻结表单白名单是 `inputNumber` 而非 `number`（`form-controls.ts:19–33`；`component-registry.json` 登记 `inputNumber`）。`type: "number"` 被 `gateRenderFormFields` 记为 `FORM_TYPE_NOT_WHITELISTED` 并丢弃该字段（`render.ts:607–613`）。另：任何 `defaultValue` 要求 protocol ≥ 2.7 **且** `form.controls.advanced`（`form-controls.ts:367–382`；`page.schema.json` L2 说明）。缺能力时 `switch` 的 `defaultValue` 同样被挡。`hasBlockingErrors` 为真时提交直接 return，按钮 disabled（`render.tsx:1281–1300,1382`）。

这不是 AJV `page.schema.json` 结构失败（字段类型不在 JSON Schema 里枚举），而是 D-002 §1「既有 form controls / 无新 renderer 扩展」与 Host L2 门的失败。A-002 self 写「页面 schema 通过 page.schema.json AJV 校验」不足以覆盖本门。

API 层 CRUD 不受影响（生命周期测试走 HTTP，不经 renderer）。

**建议修复**：将四处 `type: "number"` 改为 `inputNumber`；两页 `requiredCapabilities` 增加 `form.controls.advanced`（与 `admin.users` / `admin.roles` 先例一致）。用 `checkFormCapabilities` / 打开新建/编辑弹窗确认无 `FORM_*` 告警且可提交。

### F-002 · 行操作 `navigate` 未声明 `actions.row.navigate`；当前 Host 行分发不会执行该类型

| 字段 | 值 |
|------|-----|
| level | required |
| status | open |
| evidence | `schema/data-dictionary.json:157–159` `openEntries`；表行 `actionRef: "openEntries"`（220–226）；`requiredCapabilities` 无 `actions.row.navigate`（6–14）。登记表：`component-registry.json:561–563`「navigate 须声明 actions.row.navigate」。Host：`invokeAction` 只特判 `modal`，其余走 `runRowAction` → `runRequest`（`render.tsx:702–729`）；`runRequest` 对非 `request`/`custom` 返回 `ACTION_NOT_REQUEST`（`render.tsx:369–373`） |
| severity | med（协议声明缺失 + 计划中的「条目」行操作在当前 Host 上 fail-closed；无授权绕过） |

**说明**：`openEntries` 的 JSON 形状是合法 NavigateAction（`type` + 单斜杠 `url`），AJV 通过。但行级 `actionRef` 指向 navigate 时协议要求声明 `actions.row.navigate`（ADR-0021）。本仓现有页面中，data-dictionary 是唯一使用 `"type": "navigate"` 的业务页；Host 行分发尚未实现该类型，点击「条目」会得到 `ACTION_NOT_REQUEST`，而不是路由到 `/dictionary-entries`。这与 D-002 §1「既有导航路径、无新 renderer 扩展」的可执行含义不一致：模块消费了一条 Host 尚未接线的协议路径。

不构成授权绕过（fail-closed）。用户仍可手输 `/dictionary-entries`（manifest 已登记该页）。

**建议修复**：两选一并留痕——(1) 声明 `actions.row.navigate`，并在 `invokeAction` 对 `type: "navigate"` 走应用路由（这是 ADR-0021 的既有契约，属 Host 补齐而非新控件）；或 (2) 改用当前 Host 已能执行的导航方式（并修订 D-002）。不要在未声明能力的情况下依赖行 navigate。

### F-003 · 类型级联删除不写条目级 `dictionary.delete`

| 字段 | 值 |
|------|-----|
| level | recommended |
| status | open |
| evidence | `DeleteType` 只 `DELETE FROM dict_types WHERE id = ?`（repository.go:187），条目靠 FK CASCADE 消失；`recordDictionaryEvent(..., EventDictionaryDelete, ..., id)` 仅用类型 id（dictionary.go:115–121）。无按条目 id 的 `dictionary.delete` |
| severity | low（D-002 §5 只要求 record_id = 被删行 id；类型删除语义已文档化为级联。取证时无法从 operation_log 还原被级联销毁的 entry id） |

**建议修复**：若审计要可还原从属行，在同一事务内先选出将被级联的 entry id 并各写一条 `dictionary.delete`（或在类型事件 `detail` 中记录 entry id 列表）。v1 接受残余则可 `accepted-residual`（范围：级联条目无独立审计行；复审触发：合规/取证要求按条目追溯）。

## 必改项汇总

- **required / 必改**：F-001、F-002
- **recommended**：F-003（不阻断本安全/数据门禁的授权、完整性、SQL、迁移项）

## 与既有意见的异同

| 条目 | 关系 |
|------|------|
| A-002 self（S2–S4，verdict pass） | **同意**授权、UNIQUE、PATCH 合并、级联、DICT_KEY_NOT_FOUND、0020 CHECK、迁移连续、Profile 内容扩展。**不同意**「页面 schema / 无新 renderer 扩展」已满足——A-002 以 AJV 结构校验代替 L2 表单白名单与行 navigate 能力门 |
| A-001 self（S1） | 方案级安全意图与 API 实现一致；S1 对 NavigateAction 的选用在 Host 行分发上尚未落地 |

## 结论 + 建议给编排器/用户的下一步

**verdict: conditional** — 安全/数据核心（授权 fail-closed、完整性约束、PATCH 不擦除、SQL 参数化、审计事件与 CHECK 对齐、迁移保行、Profile 装配语义未改）有代码与测试证据。不可无条件放行本门禁：两页 schema 的表单控件与行 navigate 能力声明未通过协议/Host 门（F-001、F-002）。

建议 `/govern`：响应本意见；先修 F-001（`inputNumber` + `form.controls.advanced`）；F-002 在「补 Host navigate 分发」与「改导航方式」之间按 P-004 取舍并留痕；F-003 可维护波次或 residual。闭合 required 前不要把本目标标为 `done`。

### 声明

本意见 `source: independent`，**不修改**目标 `status` / `progress` / goal-tree / 方案正文或实现代码；响应与状态变更由 `/govern` 与用户裁决处理。
