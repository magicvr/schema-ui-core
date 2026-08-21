---
id: A-003
goal: GOAL-014-form-experience
source: independent
auditor: grok-4.6
date: 2026-08-14
scope: S1~S4 全部记录 + 实现代码（关门前独立审计）
audit_type: close-out
verdict: fail
parent: GOAL-014-form-experience
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
status: recorded
---

# A-003 · independent 审计（S5 关门前 · R4 表单体验）

## 结论

**verdict: fail**

服务端可选 `fieldErrors` 信封、单列默认布局、以及 `validateFieldValues` 纯函数本身有证据。但 D-002 §3 的两条用户可见主路径在 schema 驱动表单上**均未接通**：

1. schema 上的 `required` / `pattern` / `minLength` / `maxLength` 在 `gateRenderFormFields` 被丢掉，提交前校验对示范页不生效；
2. `readResourceApiError` 已解析 `fieldErrors`，`runRequest` 组装 `ActionResult` 时丢弃，回显代码是死的类型断言。

叠加表单级 alert 优先 `messageKey`（无字段 params），用户仍只看到 catalog 通用文案。触发场景（数据字典编辑空缺项 → 仅 `INVALID_PATCH_FIELD`）在真实 UI 路径上**未修复**。A-002 `pass` 与 E-003「提交前校验阻止请求 / fieldErrors 回显」主张名不副实。

**不可关门。** 存在未闭合 required findings。

## 范围与区间

| 项 | 值 |
|----|-----|
| 工作区 | `workspace-011-admin-functional-modules`（Root `GOAL-001-admin-functional-modules`；`shared_materials_catalog: none`） |
| 目标 | `GOAL-014-form-experience` |
| 记录 | D-001/D-002；E-001~E-005；A-001/A-002（self） |
| 代码 | errorcatalog / localize / resources / error_contract_test；resource.ts / form-controls.ts/.tsx / render.ts/.tsx；data-dictionary + dictionary-entries schema |
| 信息项 | I-001~I-004 均标 closed（S1/S4 决策成立）；本审不重开信息项——缺口是实施接线，不是未知信息 |
| 共享资料 | 无 |

本意见 **不** 修改 `status` / `progress` / goal-tree。

## 成果（有证据）

### 1. 服务端信封扩展（D-002 §2）——HTTP 层成立

| 主张 | 证据 |
|------|------|
| `Body` 仅对 `INVALID_CREATE_FIELD` / `INVALID_PATCH_FIELD` 拼接 caller message；其它 catalog 码保持纯本地化文本 | `errorcatalog.go:182-197` `isFieldValidationCode` |
| `FieldError` + `BodyWithFields`：`fieldErrors` 可选；`len==0` 不写该键 | `errorcatalog.go:199-215` |
| factory create/patch 走 `writeLocalizedFieldError` | `resources.go:492-495`、`555-557` |
| 旧信封 `{error,message,messageKey}` 形状保持 | `Body` 仍写这三项；新增键可被旧客户端忽略 |
| frozen 码集合完整；正则覆盖 `writeLocalizedFieldError` | `error_contract_test.go:71`；本审 `go test ./internal/handler/ -run TestErrorCodeContractPinnedSet\|TestErrorCatalogCoversFrozenCodes` **绿** |

### 2. 校验器纯函数（D-002 §3.2）——单测层成立

`validateFieldValues`（`form-controls.ts:445-508`）语义与 D-002 / 审计要点一致：

- 布尔 `switch`/`checkbox` 跳过 `required`（`452`）
- 空值跳过 `pattern`（`462`）
- 非法正则 `try/catch` 降级，不拦提交（`474-477`）
- `undefined`/`null` 视为空（`isEmptyValue` `431-436`）

`form-controls.test.ts` +4 覆盖 required/布尔/pattern/长度/数值/非必填空值。本审 `npm test -- form-controls.test.ts render.test.tsx schema-crud.test.tsx`：**83/83 绿**。

### 3. 布局（D-002 §4 / 用户 brief：单列 + max-w-lg）

- 默认 `grid-cols-1` + `data-form-columns="1"`（`form-controls.tsx:721-731`）；硬编码两列已移除
- modal 保持 `max-w-lg`（`modal.tsx:91`）；与 E-004 / 本审计 brief 一致（D-002 的 schema `width` 未做，见 F-005）
- 本审 render / schema-crud 全绿：单列默认未破坏现有测试。仓库内**无**表单 schema 声明 `columns`，不存在「两列 fixture 依赖」被破坏的证据

### 4. 兼容性 / 解析边界

| 边界 | 行为 | 判定 |
|------|------|------|
| 旧客户端忽略未知键 | `fieldErrors` 可选；缺省不写 | 安全 |
| `fieldErrors` 非数组 | `parseFieldErrors` 返回 `[]`（`resource.ts:113-114`） | fail-closed，不抛 |
| `fieldErrors: []` | 同上；不内联 | 可接受 |
| `readEnvelope` 未知键 | 只抽 `error/message/messageKey/params/fieldErrors` | 容忍 |
| `validateFieldValues` 对 `undefined` | 当空；非 required 不报 | 正确 |

### 5. go（VP-008）——同意不 held

对照 workspace.md VP-008 接口（Profile 默认集 / 模块矩阵 / Manifest 装配 / 协议 pin / 共同门禁语义）：本波为可选错误键 + 宿主字段约束 + 渲染层布局，**不**改装配顺序、Profile、pin 或门禁语义。协议不 bump 合理（可选键；`required` 已是 registry 控件属性）。I-004 closed 维持。

### 6. 字段名对齐（若回显接通）

factory 使用 `CreateFields`/`PatchFields` 字面量（字典：`key`/`name`/`dictKey`/`entryKey`/`label`，`dictionary.go:241-262`），与 schema `fields[].id` 一致。**不是** `id` vs 业务键错位。回显失败原因是 F-002，不是命名。

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| S1 方案冻结（可选 fieldErrors + 约束最小集 + 单列） | 达成 | D-002；A-001 |
| S2 服务端 fieldErrors 信封 | 达成（HTTP） | errorcatalog + resources + E-004 冒烟记录 |
| S2 前端 schema 约束提交前校验 | **未达成** | F-001：解析丢约束 |
| S2 服务端 fieldErrors 回显内联 | **未达成** | F-002：ActionResult 丢字段 |
| S2 单列默认 + modal max-w-lg | 达成 | form-controls.tsx；modal.tsx |
| S3 单测绿 / error_contract | 达成（覆盖面不足） | 本审 83 + handler 契约绿；未锁 F-001/F-002 |
| S4 go 不 held | 达成 | E-005；本审复核 |
| 触发：字典编辑空缺项有字段级提示 | **未达成** | F-001+F-002；见下 |

## Findings

### F-001 · schema 约束未进入 `gateRenderFormFields`（提交前校验对示范页无效）

- **严重度**：high
- **建议**：**required**
- **状态**：open
- **描述**：

  D-002 §3 / E-003 / A-002 主张 schema `required`/`pattern`/`minLength`/`maxLength` 在提交前由 `validateFieldValues` 拦截。`FormInner` 校验的是 `gate.fields`（`render.tsx:1309-1332`），来自 `gateRenderFormFields`（`render.ts:575-651`）。

  该解析器**拷贝** `min`/`max`/`defaultValue` 等既有键，**不拷贝** GOAL-014 新增的四类约束：

  ```615:651:apps/web/src/renderer/render.ts
      fields.push({
        id: entry.id,
        type: entry.type,
        // ... label / min / max / defaultValue ...
        // 无 required / pattern / minLength / maxLength
      });
  ```

  示范页已写约束（`data-dictionary.json:71-84,128-134`；`dictionary-entries.json:73-91`），运行时 `field.required === undefined`，空 `name` **不会**被客户端拦住，请求仍打到 PATCH。

  `form-controls.test.ts` 直接构造带 `required` 的对象，测不到这条接线。`schema-crud` 空提交用例仍期望服务端 `INVALID_CREATE_FIELD`（`schema-crud.test.tsx:458-463`），与「提交前阻止」相反，回归锁的是旧路径。

- **证据**：`render.ts:615-651`；`render.tsx:1309-1332`；`data-dictionary.json:128-134`
- **闭合**：`fixed` = 解析透传四类约束 + 至少一条 schema-crud/render 测：带 `required` 的空提交不发请求且字段内联。`accepted-residual` / `user-overruled` 须书面范围（例如「只做服务端回显」）——当前回显也未接通（F-002）。

### F-002 · `fieldErrors` 未进入 `ActionResult`；表单级 alert 仍走通用 catalog

- **严重度**：high
- **建议**：**required**
- **状态**：open
- **描述**：

  `readResourceApiError` 已解析 `fieldErrors`（`resource.ts:113-160`）。`runRequest` 在 `!response.ok` 时只转发 `code/message/messageKey/params`（`render.tsx:461-469`）。`ActionResult` 类型无 `fieldErrors`（`render.tsx:164-166`）。

  `handleSubmit` 用断言读取不存在的键（`render.tsx:1350`），恒为 `undefined`，`setFieldErrors` 恒为空对象。D-002 §3.1「映射到字段内联；无匹配回退表单级」未发生。

  表单级 alert 优先 `messageKey` 且无 params（`render.tsx:1420-1426`）。服务端 `Body` 拼进 `message` 的字段原因（`errorcatalog.go:185-188`）**不会显示**。zh-CN 用户仍看到 `INVALID_PATCH_FIELD: 更新字段无效`——与立项触发句同一症状。

  E-004 HTTP 冒烟只证明 API 信封，不证明 UI 回显。A-002「服务端 fieldErrors 回显字段内联」无端到端证据。

- **证据**：`render.tsx:164-166,461-469,1348-1361,1420-1426`；`resource.ts:127-160`
- **闭合**：`fixed` = `ActionResult`/`runRequest` 透传 `fieldErrors`；匹配字段内联；无匹配才表单级；有 `fieldErrors` 时勿用无 params 的 `messageKey` 盖掉具体原因。补一条 mock 带 `fieldErrors` 的提交测。

### F-003 · 内联错误文案未 catalog 化（en/zh）

- **严重度**：low
- **建议**：non-blocking（recommended）
- **状态**：open
- **描述**：D-002 §7 将「字段级错误文案（en/zh）」列入影响面。`validateFieldValues` 六条 message 为英文字面量（`form-controls.ts:458-504`）；服务端 `reason` 为英文（`not be empty` / `must not be empty` / `be a string`）。无 `error.field.*` catalog。即便 F-001/F-002 修好，zh-CN 内联仍是英文。建议 catalog 化客户端码；服务端 `reason` 用稳定码或 `messageKey`+params，避免把英文诊断当 UI 文案。

### F-004 · 其它 `INVALID_*_FIELD` 路径未附 `fieldErrors`

- **严重度**：low
- **建议**：non-blocking
- **状态**：open
- **描述**：D-002 §2.2 只冻 factory decode。仍有同码出口不走 `writeLocalizedFieldError`：

  | 路径 | 行为 |
  |------|------|
  | `account_self.go:131` | `writeLocalizedError(..., INVALID_PATCH_FIELD, "name must not be empty")` — `Body` 会拼接 message，**无** `fieldErrors` |
  | `dictionary.go:68,165`、`scheduledtasks.go:80` | `DomainError` + 同码 → `writeEntityError` → `writeLocalizedError`（`resources.go:275-277`）。在 CreateFields 已强制后多半不可达 |

  非误拼接（其它码仍不拼诊断）。个人中心改名空值仍无字段数组。建议个人中心补 `fieldErrors`；DomainError 死路径可删或改专用码。

### F-005 · `columns` 可配的 CSS 与 D-002 width 未落地

- **严重度**：low
- **建议**：non-blocking
- **状态**：open
- **描述**：

  1. `form-controls.tsx:727` 动态拼接 `"sm:grid-cols-" + cols`。仓库内无完整字面量 `sm:grid-cols-2/3/4`；Tailwind v4 扫描通常**不会**生成这些 utility。`data-form-columns` 会变，视觉仍可能单列。当前无 schema 使用 `columns`，默认单列不受影响。
  2. D-002 写「上限 4 与 GridView 一致」；GridView 上限是 **6**（`render.tsx:1503-1505`）。
  3. D-002 ModalHost 可选 `width` 默认 480px 未实现；本审按用户 brief 接受现成 `max-w-lg`。

  建议：完整 class 映射或 safelist；上限与 GridView 对齐或改 D-002 表述。

### F-006 · `date-range` 空对象不算空；switch `required` 与 registry 不一致

- **严重度**：low
- **建议**：non-blocking
- **状态**：open
- **描述**：`isEmptyValue` 对 `date-range` 只把 `undefined`/`null` 当空（`form-controls.ts:431-436`），`{start:"",end:""}` 不算空。`required` 的 dateRangePicker 会**漏拦**，不是误拦合法表单。`upload` 空串/空数组行为正确。当前无 schema 给这两类加 `required`。

  registry ADR-0028：`switch.required:true` 要求值为 `true`（`component-registry.json:1097-1116`）。D-002 / 实现一律跳过布尔。属方案取舍，不是误拦。

### F-007 · 无自动化锁 `fieldErrors` 信封与回显

- **严重度**：med
- **建议**：non-blocking
- **状态**：open
- **描述**：仓库内无 `BodyWithFields` / `fieldErrors` JSON 的 Go 测；无 render/schema-crud 测断言内联回显或「required 阻止 fetch」。E-004 手测 HTTP 未防止 F-001/F-002。S3「911/911 全绿」不覆盖本波用户路径。本审未复跑全量 911（只跑相关 83 + error_contract）。

## 必改项汇总

| ID | 级别 | 一句话 |
|----|------|--------|
| **F-001** | **required** | `gateRenderFormFields` 丢弃 schema 约束，提交前校验对数据字典等示范页无效 |
| **F-002** | **required** | `runRequest` 丢弃 `fieldErrors`；表单级 alert 用无 params 的 `messageKey`，用户仍只见通用 INVALID_*_FIELD |
| F-003 | non-blocking | 内联文案未 i18n catalog |
| F-004 | non-blocking | account_self / DomainError 同码无 fieldErrors |
| F-005 | non-blocking | columns 动态 class / 上限 4 vs GridView 6 / width 未做 |
| F-006 | non-blocking | date-range 空对象；switch required 与 registry |
| F-007 | non-blocking | 缺信封与回显自动化锁 |

## 与既有意见的异同

| 条目 | 关系 |
|------|------|
| A-001 (self, S1, pass) | 同意方案冻结。协议不 bump：本审同意（可选键 + 已有 `required` 语义）。A-001 备注仍成立但非 required。 |
| A-002 (self, S2/S3, pass) | **分歧**：A-002 称提交前校验阻止请求、fieldErrors 回显、回归足以关门前自审。本审证实纯函数与 HTTP 信封存在，**schema→校验**与 **API→内联**两条接线断裂；911 绿未覆盖用户路径。故 S2/S3 主张不成立，verdict **fail**。 |

A-002 备注「其余模块 required 可后续补」：在 F-001 未修前，即便补 schema 也不生效。dictionary-entries **编辑**表单未给 `label`/`dictKey` 加 `required`（仅创建有），即使 F-001 修好，条目编辑空 label 仍依赖 F-002。

## 结论与建议下一步

1. **禁止**在 F-001、F-002 未按 `fixed` / `accepted-residual` / `user-overruled` 合法闭合前将 GOAL-014 标 `done`（P-003）。
2. 建议 `/govern` 响应 A-003：优先 **fixed** 两条接线（解析透传约束 + ActionResult 透传 fieldErrors + 有字段错误时不要用裸 messageKey 盖具体原因），并补 F-007 所述最小回归后再请复审。
3. F-003~F-006 不阻断关门（在 F-001/F-002 闭合后）。F-003 建议与 VP-007 双语一并排期。
4. I-001~I-004 无需重开；go 不 held 维持。
5. 无 required 信息门禁到期；共享资料目录为 none。

## 声明

本意见 `source: independent`，不修改目标 `status` / `progress` / goal-tree。响应、finding 闭合与是否关门由 **`/govern`** 处理。


---

## 编排器响应（2026-08-14）

### F-001 · fixed（schema 约束解析透传）

- `gateRenderFormFields`（render.ts）补 required/pattern/minLength/maxLength 透传；新增回归测试（render.test.ts「constraint passthrough」）。
- 示范 schema（data-dictionary / dictionary-entries）的 required 现在运行时生效：空 name 提交前被拦截。

### F-002 · fixed（fieldErrors 进 ActionResult）

- `ActionResult` 增加可选 `fieldErrors`；`runRequest` 从 `readResourceApiError` 转发；`handleSubmit` 用真实字段访问（去掉恒 undefined 的类型断言）。
- 新增回归测试（error-localization.test.tsx「parses fieldErrors from the GOAL-014 envelope」）。

### F-003 · fixed（内联文案 catalog 化）

- en/zh 增加 `form.validation.*` 6 键；`validateFieldValues` 每条错误带 messageKey；`handleSubmit` 用 t() 翻译（缺键回退 message）。

### F-004 · fixed（account_self 同码路径）

- account_self.go 的 name 空校验改走 `writeLocalizedFieldError`（带 fieldErrors）；errorcatalog import 补齐。

### F-005 · fixed（Tailwind 动态类）

- `GRID_COL_CLASSES` 静态查找（sm:grid-cols-2/3/4），避免 JIT 无法提取拼接类名；modal width 保持 max-w-lg（本次 brief 可接受）。

### F-006 · accepted-residual（dateRange 空值）

- `{start:"",end:""}` 不算空为既有语义（dateRange 两个输出字段独立提交）；如需必填可在后续波次声明 required 于 startField/endField。

### F-007 · fixed（自动化锁）

- F-001/F-002 回归测试已加；全量 web 913/913 + go 全绿复跑通过。
