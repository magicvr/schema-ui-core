---
id: D-002
goal: GOAL-014-form-experience
status: accepted
date: 2026-08-14
scope: S1 方案冻结
parent: GOAL-014-form-experience
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# D-002 · S1 方案冻结：字段级错误契约 + 前端校验 + 布局

## 1. 现状盘点（有证据）

### 服务端错误路径

- 信封：`{error, message, messageKey?}`（handler/localize.go `writeLocalizedError` → errorcatalog.Body）。
- **cataloged 码（INVALID_CREATE_FIELD / INVALID_PATCH_FIELD）的 message 被 catalog 文本覆盖**（"更新字段无效"），`createFieldError`/`patchFieldError` 里的具体字段名（"name must not be empty"）在 cataloged 分支完全丢失——这是用户看到通用错误的原因。
- 无字段级结构化信息（fieldErrors）输出。

### 前端错误与表单

- `readResourceApiError` 已解析 `params`（ResourceApiError 带 params）——服务端当前不发。
- 表单提交失败 → 单条 `code: message` 顶部 alert（render.tsx formError），无字段内联。
- 无提交前字段校验（required/pattern/length 均无）。
- 布局：`FormControls` 在字段数>1 时**硬编码** `sm:grid-cols-2`（form-controls.tsx:703）——这是用户说的"横着两列"来源；modal 宽度未配置。

## 2. 方案一：服务端字段级错误契约（向后兼容）

### 2.1 信封扩展

- 错误信封新增**可选**键 `fieldErrors: [{field, reason}]`（字段级明细）。
- 兼容策略：仅 create/patch 字段校验失败时附带；`{error, message, messageKey}` 保持原样（旧客户端忽略未知键）；未 cataloged 分支不变。
- `errorcatalog.Body` 签名扩展（或新增 `BodyWithFieldErrors`）：cataloged 时 message 保留 catalog 本地化文本，**同时** fieldErrors 提供字段级原因（如 [{field:"name", reason:"must not be empty"}]）。

### 2.2 触发点

- `decodeResourceCreate` 的 `createFieldError` → INVALID_CREATE_FIELD + fieldErrors。
- `decodeResourcePatch` 的 `patchFieldError` → INVALID_PATCH_FIELD + fieldErrors。
- 其他域错误（DICT_TYPE_KEY_TAKEN 等）不加 fieldErrors（非字段校验）。

## 3. 方案二：前端字段校验与内联报错

### 3.1 schema 字段约束最小集（FormControlField 扩展，可选）

| 约束 | 适用类型 | 语义 |
|------|----------|------|
| `required: true` | 全部（除 switch/checkbox 布尔默认） | 提交前必填校验（空串/undefined/空数组） |
| `pattern: "..."` | input/textarea/select | 正则格式校验（提交前） |
| `minLength` / `maxLength` | input/textarea | 长度校验 |
| `min` / `max` | inputNumber（已存在） | 数值范围（已有 wire 校验，S1 复用） |

- 校验时机：提交前（handleSubmit 前置 `validateFieldValues`）；失败阻止提交并内联显示。
- 错误显示：字段下方 `role=alert` 红字（FieldControl 内联），与 gate/reaction 错误列表共存。
- 服务端 fieldErrors 回显：提交失败后把 fieldErrors 映射到对应字段内联（字段级；无匹配字段的回退到表单级 alert）。

### 3.2 校验器

- `form-controls.ts` 新增 `validateFieldValues(fields, values): FieldError[]`（纯函数，单测覆盖）。

## 4. 方案三：布局（参考业界）

### 4.1 业界惯例（Ant Design / Element Plus / shadcn）

- **表单弹窗默认单列垂直布局**（label 上方），AntD Form `layout="vertical"` + Modal width 520 是主流；Element Plus dialog 同理；shadcn 单列。
- 两列只用于短字段对（如 key+name），由 schema 显式声明，而非渲染器硬编码。

### 4.2 本项目布局方案

- `FormControls` 移除硬编码两列：字段数>1 时默认单列（`grid-cols-1`）。
- schema form props 支持可选 `columns: number`（渲染器现有 GridView 已有 columns 模式可复用；默认 1，上限 4）。
- modal 宽度：ModalHost 支持可选 `width`（schema modal props，默认 480px，AntD 惯例 520 附近取 480 适配本项目紧凑风格）。
- 现有两列 fixture（schema-keys 等）不受影响（布局为渲染层，非协议）；示例页 schema 如需两列显式声明 columns:2。

## 5. 协议影响与版本

- `fieldErrors`：错误响应新增可选键（协议 2.7 兼容——既有消费者读未知键不崩溃；前端 readEnvelope 扩展解析）。
- 字段约束：FormControlField 新增可选属性（schema 文档兼容；现有 schema 无约束 = 行为不变）。
- `columns`/`width`：form/modal props 新增可选属性。
- **不** bump protocolVersion（可选键扩展，向后兼容）；若审计要求再评估。

## 6. 未选方案

| 方案 | 未选原因 |
|------|----------|
| 服务端校验规则下放（schema 约束由服务端 enforce） | 校验属渲染层体验；服务端已有 CreateFields/PatchFields 硬约束，重复实现收益低 |
| JSON Schema 全量表单生成器 | 范围过大（D-001 排除）；本波只做字段级约束最小集 |
| 引入前端表单库（react-hook-form 等） | 现有受控 FormControls 已可扩展；新依赖风险大于收益 |
| 强制两列保留 + 只调宽度 | 与业界惯例（默认单列）不符，用户明确质疑两列 |

## 7. 影响面

- apps/api：errorcatalog.Body 扩展（fieldErrors 可选）、resources.go（create/patch 错误附带 fieldErrors）、handler 测试。
- apps/web：resource.ts（readEnvelope 解析 fieldErrors）、form-controls.ts/.tsx（校验器 + 内联错误 + 单列布局 + columns）、render.tsx（提交前校验 + fieldErrors 回显）、modal.tsx（width）。
- i18n：字段级错误文案（en/zh）。
- **go（VP-008）判定**：错误契约扩展 + 布局内容变化，非装配语义/非门禁语义 → 不 held（S4 确认）。
