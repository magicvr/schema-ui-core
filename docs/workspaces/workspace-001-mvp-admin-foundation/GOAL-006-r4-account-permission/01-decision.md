---
id: GOAL-006-r4-account-permission
doc: decision
status: active
parent: GOAL-001-mvp-admin-foundation
created: 2026-07-31
updated: 2026-07-31
version: 0.2.0
---

# 决策记录 · GOAL-006

## 信息需求与阶段门禁

P-005 信息台账维护在 [00-meta.md](00-meta.md)。本目标的 `I-006-001` 在 R4 **方案冻结前**须验证；父目标 `I-PROTO-002` 在 R4 **实施前**须合法闭合。本目标不修改 Root 门禁状态。

## D-001 · 立项 R4 核心账号与权限子目标

**日期**：2026-07-31
**状态**：accepted

**决定**：

在 `GOAL-001-mvp-admin-foundation` 下创建 `GOAL-006-r4-account-permission`，将 R4 范围限定为账号权限最小 API 设计、`D-PERM` 映射冻结与前后端鉴权链路实现；Root 保持 `active`，纲领进度仍为 `3/6`。

**为什么**：

- Root 路线图把 R4 定义为「核心账号与权限」，且 `I-PROTO-002` 要求「设计最小 API + 对照 permissions-inheritance fixtures」。
- R3（GOAL-005）已完成 Admin 外壳/导航，其默认空 navigation context 需要真实身份/权限来源衔接——正是 R4 的产品边界。
- 协议清单将账号/权限映射为 `D-PERM` 核心能力，并有配套 behavioral fixture 可作验证来源。

**未选方案**：

- 在 R3 目标中补入鉴权：会越过 R3 已关门边界并绕过 `I-PROTO-002` 实施门禁。
- 吞并 R5 Renderer 全量与范例页：会提前越过 `I-PROTO-003` 验收/关门门禁。
- 以硬编码角色表代替 `D-PERM` 契约映射：无法证明与固定协议语义一致。

**影响**：

本目标进入 `active` 规划阶段；R4 方案冻结前验证 `I-006-001`，实施前闭合 `I-PROTO-002`。当前不修改 `apps/*`，不改变父目标或协议门禁状态。

## D-002 · 采用「信息就绪 → 方案冻结 → 实施 → 验证关门」的 R4 路线

日期：2026-07-31
**状态**：accepted

**决定**：

先处理 `I-006-001` 并闭合父目标 `I-PROTO-002`，再冻结账号权限最小 API、`D-PERM` 映射与前后端集成边界；在此之前不把待确认取舍写成实现事实。实现完成后必须补结构/行为/运行时证据和阶段自审，才可讨论 `done`。

**为什么**：

- 沿袭 R3（GOAL-005 D-002/D-005）已建立的信息门禁纪律：开放 required 信息项是方案冻结与受影响实施门禁。
- 账号/权限涉及安全边界，未验证的映射或伪造的「已验证」都会在验收时产生误导。

**未选方案**：

- 先写鉴权代码再补契约映射：会把未知项伪装成已决定行为。
- 仅以构建成功或登录页可打开作为 R4 关门证据：无法覆盖 `D-PERM`/权限继承契约。

**影响**：

`I-006-001` 与 `I-PROTO-002` 是 R4 方案冻结/实施门禁；必要时须按 P-004 由用户裁决 fixed、accepted-residual 或 user-overruled，不能静默放行。

## D-003 · 采用固定上游资料与明确的验证证据边界

日期：2026-07-31
**状态**：accepted

**决定**：

R4 规划以 `protocol-inventory-v2.7.0.md` 登记的 source commit `ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b` 为资料定位锚点，优先对照 `D-PERM` 与 `permissions-inheritance` 等固定 fixture；未来验收须记录 schema/fixture 或等价可核对证据。`conformance/reference-js`、`reference-python` 和 `runner` 仅可作为参考，不能单独证明兼容。

**为什么**：

协议清单明确区分 structural contract、behavioral fixture 和 excluded reference/runner；沿用 R3（GOAL-005 D-003）已确立的证据边界，避免把「找到路径」或「调用参考实现」误写成已验证。

**影响**：

该决定固定证据方向，不代表当前已有本地 schema、fixture、运行时或 conformance 结果。

## D-004 · 冻结 R4 方案：账号权限最小 API 与 D-PERM 映射

**日期**：2026-07-31
**状态**：accepted（用户按 `/govern` 建议确认「按此方向冻结」）

**决定**：

按固定协议 `schema-ui-docs@v2.7.0`（commit `ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b`）冻结 R4 最小 API 与 `D-PERM` 映射，闭合本目标 `I-006-001` 与父目标 `I-PROTO-002`（二者均为 required，现 `verified`）。

**1. 最小 API（账号会话最小闭环）**

| 端点 | 语义 | 输入 → 输出 |
|------|------|-------------|
| `GET /api/accounts/me` | 返回当前账号会话与权限上下文（`$context` 快照来源） | 会话 → `{ user, features }` |
| `POST /api/accounts/login`（可选） | 会话建立（本仓脚手架可先于 e2e 用静态/注入会话） | 凭据 → 会话 token |
| `POST /api/accounts/logout`（可选） | 会话失效 | 会话 → 空 |

R4 实施范围仅包含最小会话与权限求值链路；不建账号 CRUD 管理页、不建 SSO/联邦、不做细粒度审计后台。会话方案（静态/注入 vs token）与端点形态在实施计划（`02-execution` 路线图）落地时记录事实，不改变本冻结的契约语义。

**2. D-PERM 结构/行为映射（冻结）**

- **`$context` 快照**：前端 Renderer 只读消费 `$context.user.*` / `$context.features.*`（页面快照语义，ADR-0003），由 `GET /api/accounts/me` 提供；不是安全边界。
- **`permissions`**：`view` / `edit` / `delete` 三个标准键（`node.schema.json` `Permissions`）；表达式仅允许 `$context.*`，禁止 `$deps.*`（`01-node-protocol.md` §3.9）。
- **`permissionCascade`**：仅 `section` / `grid` / `form` / `tabs` / `table` 可声明；`keys` 非空、去重、仅 `edit` / `delete`，且须同名 `permissions.<key>`（ADR-0023 / `node.schema.json` `PermissionCascade`）。
- **有效权限公式**（参与 cascade 的目标）：`effectivePermission(t, k) = AND(根到 t 路径上声明 k ∈ permissionCascade.keys 的祖先的 permissions[k]) AND (t.permissions[k]，若 t 自身声明)`；未声明按 `true`；只能收紧、不能放宽（ADR-0023 D3 / `01-node-protocol.md` §3.9.1）。
- **权限结构边**：Node `children[]`、`tabs.props.items[].content`、`table.props.actions[]`、`table.props.toolbar[]`、default form 隐式 `submitAction`；`table.props.columns[]` 不在树上（只本地 `permissions`，不吃 cascade）；modal `content` 与 navigate 新页面是新根（ADR-0023 D2a）。
- **`permissionIntent`**：仅 RowAction、toolbar Trigger、`actionButton.props` 可声明 `edit` / `delete`；未标注意图不参与 `edit` / `delete` 继承；不得从 `key` / `actionRef` / HTTP method / 文案推断意图（ADR-0023 D4b / component-registry）。
- **表单 edit 目标白名单**：default form（`mode` 缺省或 `default`）的 `input` / `inputNumber` / `datePicker` / `dateRangePicker` / `select` / `upload`；`submitAction` 为隐式 edit 目标（不写 intent）；search form 不参与（ADR-0023 D4a）。
- **执行时序**（Renderer）：`visibleWhen` → effective/local permission → `disabled` OR `requiresSelection` → fail-closed stop → `confirm` → Action/form submit；拒绝后不得展示 confirm、不得构造/发送 request/navigate/modal/submit（ADR-0023 D4c / renderer-spec §7.1）。
- **门控**：出现 `permissionCascade` 或 `permissionIntent` 须 `meta.protocolVersion ≥ "2.3"` 且 `requiredCapabilities` 含 `permissions.inheritance`；缺失由 L2 fail-closed（错误码 `PROTOCOL_VERSION_TOO_LOW` / `CAPABILITY_REQUIRED`）。
- **L2 校验错误码**：`PERMISSION_CASCADE_TYPE_INVALID`、`PERMISSION_CASCADE_KEYS_INVALID`、`PERMISSION_CASCADE_SOURCE_MISSING`、`PERMISSION_INTENT_FORBIDDEN`、`PERMISSION_INTENT_INVALID`（fixture 断言）。
- **前后端职责**：React 主责 UI 显隐/禁用与 intent 求值（renderer-spec）；Go 主责鉴权与权限模型、账号会话；后端**独立鉴权**，不把前端 `$context` 权限结果当安全边界（协议总纲 / ADR-0023 背景）。
- **行为权威**：`conformance/fixtures/permissions-inheritance/cases.json`（fixtureVersion 1.0，17 cases：13 valid + 4 invalid；target kinds：formField 7 / formSubmit 7 / rowAction 5 / actionButton 5 / toolbarTrigger 2 / column 1）。本仓不 vendor 上游 reference-js/python/runner，仅作参考。

**3. 固定资料证据（本地落盘，SHA-256 核验）**

资料落盘于 `attachments/dperm/`，均取自固定 commit `ca9e5fe…`：

| 文件 | SHA-256 |
|------|---------|
| `permissions-inheritance/cases.json` | `ac124fa1d831d0aa2544b7544b1e177c3498c8c3b36ee4d535e8c3f2f5b8849e` |
| `docs/schemas/node.schema.json` | `967c0d4fa9068eb0ebd09456905f17087624b54150dafd5d53b28c31bd04fb16` |
| `docs/decisions/0023-container-permission-inheritance.md` | `1a82bdf0e39747eb200f1c55682160c297d10820084a95b0e33febd14765209c` |
| `docs/05-scenarios/permission-inheritance.md` | `94df1dae0d9f49c7eef7bcba7da66219592a311c1ddf9a1ec1c653d660311c6b` |

语义规范原文见固定 commit：`docs/01-node-protocol.md` §3.9 / §3.9.1（`permissions` / `permissionCascade`）、`docs/03-component-registry.md`（intent 挂载点矩阵、表单 edit 目标）、`docs/08-renderer-spec.md` §7.1（执行时序与门禁）；inventory §2.4 ADR-0023 索引一致。

**未选方案**：

- 以硬编码角色表代替 `D-PERM` 契约映射：无法证明与固定协议语义一致（D-001 沿袭）。
- 在 R3 目标中补鉴权：越过 R3 关门边界与 `I-PROTO-002` 门禁（D-001 沿袭）。
- 将 Go 侧做成「按前端权限结果过滤」：把展示层权限当安全边界，违反 ADR-0023「后端仍必须独立鉴权」。
- 在方案冻结阶段就实施代码：`I-PROTO-002` 仅闭合设计/映射门禁，R4 实施仍需用户指令与实施事实记录（见 02-execution 完成后边界）。

**影响**：

- `I-006-001`（required）→ **verified**：证据 = 本 D-004 + `attachments/dperm/` 固定资料（SHA-256）+ 覆盖表 v0.1.3 `D-PERM=include`、`permissions-inheritance=include`。
- 父目标 `I-PROTO-002`（required，R4 实施门禁）→ **verified**：证据 = 本 D-004 + `attachments/dperm/`；闭合范围仅限「账号权限最小 API 与 D-PERM 映射」设计结论，**不**放行 R4 实施本身，也**不**改变 `I-PROTO-003`（R5 验收/关门）或 `I-PROTO-004`（vendor vs pin，仍 open）。
- 方案冻结不等于实施：R4 实施、前后端代码、fixture 测试与运行时证据仍须按 02-execution 路线图推进并记事实；实施完成且阶段自审通过前不讨论 `done`。
