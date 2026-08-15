---
id: A-003
goal: GOAL-018-mfa-manager-ui
title: S5 关门独立审计（MFA 个人中心管理组件 · renderer custom 节点）
date: 2026-08-15
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
scope: S5 关门（00-meta S1~S5 + D-001 §1~§4 + D-002 + E-001~E-003 + A-001/A-002 + renderer 契约 + MfaManager 流 + 验证证据 + GOAL-017 回归关门依赖链 + 安全/协议/go）
audit_type: close-out
verdict: fail
status: recorded
parent: GOAL-018-mfa-manager-ui
created: 2026-08-15
updated: 2026-08-15
version: 1.0.0
---

# A-003 · 独立关门审计意见（S5 · MFA 个人中心管理组件）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型** / **scope**：close-out · S5 关门（成功标准 S1~S5、D-001 §1~§4、renderer custom 契约、MfaManager 流、验证证据、用户裁决依赖链、安全/协议/go）
- **verdict**：**fail**

## 范围与区间

- **工作区**：`workspace-011-admin-functional-modules`（`workspace.md` 绑定 Root=`GOAL-001-admin-functional-modules`、`canonical_scope`、`plan_refs`/`primary_plan`=VP-011 已校验；`shared_materials_catalog: none`）。未读取或比较其他工作区目标状态。同区 GOAL-017 仅作用户裁决 / 回归关门依赖只读对照。
- **已通读**：本目标 `00-meta`、`01-decision.md`、D-001、D-002、`02-execution.md`、E-001～E-003、`03-audit.md`、A-001、A-002。父目标裁决上下文：GOAL-017 `03-audit/A-007`、`A-008`、`01-decision/D-005`。
- **代码核对**：`apps/web/src/renderer/render.ts`（`RenderCustomNode` / `WHITELISTED_NODE_TYPES` / `parseRenderNode`）、`render.tsx`（`case "custom"`）、`renderer/custom-components.ts`、`components/mfa-manager.tsx`、`components/mfa-manager.test.tsx`、`main.tsx`、`i18n/messages/{zh-CN,en-US}.json`、`i18n/s5-denominator-render.test.tsx`、`i18n/schema-keys.structural.test.ts`、`renderer/render.test.ts`、`protocol/load-page.ts`、`protocol/conformance/runtime-schema-validate.ts`、`docs/schemas/node.schema.json`、`apps/api/internal/modules/account/schema/account.json`、`handler/mfa.go`、`modules/mfa/service.go` `requireActiveSecondFactor`。
- **本轮复跑**（2026-08-15）：
  - `apps/web` `npx vitest run src/components/mfa-manager.test.tsx src/renderer/render.test.ts src/i18n/schema-keys.structural.test.ts` → **68/68 绿**（含同目录 `render.test.tsx`）。
  - 独立 Ajv 对 `account.json` 跑与运行时相同的 page/node/action/reaction 校验（`allErrors` / `strict:false` / `validateSchema:false`，与 `runtime-schema-validate.ts` 一致）：**`ok: false`**（见 F-001）。
  - 未复跑 web 全量 974 与 go `./...`（E-003 有记录；本 scope 以关键包 + D-VAL 复现为准。全绿不能覆盖缺失的 account 页 D-VAL，见 F-001 / F-003）。
- **covered**：成功标准对照、renderer 契约、MfaManager 流与 i18n、验证证据充分性、GOAL-017 回归关门依赖链、安全/协议/go、I-001。
- **excluded**：不改 `status` / `progress` / goal-tree / `00-meta` / D-001 / D-002 正文。
- **保证等级**：L0。

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| 五件套 + ledger 齐全；编号 018；parent=`GOAL-017-r3-s10-mfa-2fa` | `00-meta`；`01-decision/` D-001/D-002；`02-execution/` E-001～E-003；`03-audit/` A-001/A-002 |
| 用户裁决链可追踪：阻断 GOAL-017 关门 → 本子目标承接 → 本目标关门后回归关闭 GOAL-017 | GOAL-017 D-005 §3；本目标 E-001；`00-meta` 概述 / 父目标节 |
| I-001（S1）已由 D-001 §1 闭合：模块级注册表替代 props 线程，E-002 留痕 | D-001 §1；`custom-components.ts` L18–36；E-002 L17 |
| renderer 白名单 + 解析 fail-closed：缺 `component` → `RENDER_INVALID_BODY` | `render.ts` L302–315 / L430–437；`render.test.ts` L66–74 |
| 分发：未注册组件渲染 fallback 文本，不抛 | `render.tsx` L2297–2308 |
| 注册表 `register` / `get` / `resetCustomComponentsForTests`；`mfa-manager` 自注册 + `main.tsx` 引入 | `custom-components.ts`；`mfa-manager.tsx` L252–254；`main.tsx` L11–13 |
| MfaManager 组件存在 status / enroll 一次性载荷 / confirm / disable / rotate / unavailable 降级 / AuthError i18n 映射 | `mfa-manager.tsx` L29–48 / L84–160 / L172–241 |
| 无新权限键；仅消费既有 `/api/mfa/*` | `account.json` L5–14（capabilities 未增 MFA 权限）；组件只打 `/api/mfa/status|enroll|confirm|disable|recovery/rotate` |
| i18n zh/en 17 键成对存在；catalog 键集一致 | `zh-CN.json` / `en-US.json` L500–514；本轮 `schema-keys.structural.test.ts` 4/4 绿 |
| 本轮 mfa-manager 两流 + render.test custom 解析绿 | enroll→confirm；status 失败占位；`parseRenderNode` 缺 component fail-closed |
| D-002 go 不暂挂：无新后端能力、无权限键、无协议 capability / pin 变更 | D-002 L17–19；`conformance-local-report.json` `artifactVersion: 2.8.0`；本审未见 MFA API 语义改动 |
| A-001 / A-002 self 无开放 required | 索引 |

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| S1 方案冻结（D-001） | **满足**（方案层） | D-001 + A-001 pass；I-001 由模块级注册表闭合。`01-decision.md` 信息表仍写 I-001 `open`（F-004，不重开信息门禁） |
| S2 / D-001 §1 renderer 契约 | **部分** | 白名单 / 解析 / 分发 / 注册表落地；**生产页加载走 D-VAL，custom 节点顶层 `component` 被拒**（F-001） |
| S2 / D-001 §2 MfaManager 流 | **部分** | status→enroll 一次性展示→confirm→disable(TOTP)→rotate(recovery) 代码存在；**disable/rotate 未按「code 或 recoveryCode」组请求体**（F-002） |
| S2 / D-001 §3 account.json 接入 | **不满足** | 节点已写入 `account.json` L148–151，但运行时 `loadPageDocument` 校验失败，整页不可达（F-001） |
| S3 / D-001 §4 测试策略 | **部分** | 两流 + parse 白名单有；**无** account D-VAL、**无** render 分发/未注册 fallback 测试、**无** disable/rotate 单测、s5 **未**访问 `/account`（F-003） |
| S4 go 判定 | **满足** | D-002 与「本地扩展、不改协议 pin」一致；缺陷是文档未过 D-VAL，不是 go 失效 |
| S5 可关门 | **不满足** | 开放 required F-001（high）、F-002（med） |
| 用户裁决：本目标关门后回归关闭 GOAL-017 | **不可执行** | 依赖链台账可追踪（D-005）；本目标未过 S5，**不得**据此把 GOAL-017 标 `done` |

## Findings

### F-001 · `account.json` custom 节点未过 D-VAL：生产个人中心整页 `PAGE_SCHEMA_INVALID`

| 字段 | 值 |
|------|-----|
| level | **required**（high） |
| status | open |
| evidence | 生产路径是 `page.route → schemaUrl → loadPageDocument → validatePageDocument → RenderPage`（`App.tsx` L309 / L343–345；`load-page.ts` L105–112：校验失败抛 `PAGE_SCHEMA_INVALID`，**不进入 renderer**）。`docs/schemas/node.schema.json` L71 `additionalProperties: false`，properties 无 `component`。`account.json` L148–151 `body.children[3]` 为 `{type:"custom", id:"mfa-manager-block", component:"mfa-manager"}`。本轮以与 `runtime-schema-validate.ts` 相同 Ajv 选项校验该文档：`ok:false`，`instancePath=/body/children/3`，`additionalProperty=component`。`type:"custom"` 作为自由字符串可通过；顶层 `component` 不能。结果：用户打开个人中心时 **profile / 改密 / 会话表 / MFA 整页失败**，是对已关门 GOAL-005 账户页的回归，且 D-001 §3「account.json 接入」在运行时名不副实。E-002 L19、E-003 L17、A-002 L21 将接入写成已交付事实。代表页 / s5 分母均不访问 `/account`（`representative-pages.test.tsx` `MIGRATED_PAGE_IDS` 无 account；`s5-denominator-render.test.tsx` 只把 account.json 放进 document map，用例路由为 `/roles` `/settings` 等），故 web 974/974 与本缺陷共存。 |
| closure | — |
| 影响门禁 | S5 关门 / S2 account 接入 / S3 验证充分性 / GOAL-017 回归关门 |

D-002「本地扩展、不改协议」成立，但不能把协议非法文档挂进走 D-VAL 的页面。最小充分修复（保持 D-002：不改 pin / capability）：把注册键放进已允许的 `props`（`props` 无 `additionalProperties: false`），并让 `parseRenderNode` 读 `props.component`（或等价协议内字段）；**补** `validatePageDocument(account.json) === ok` 与 `/account` 经 `loadPageDocument` 的渲染断言。不要静默绕过 D-VAL。不要把 `component` 写进 `node.schema.json` 却仍声称无协议变更。

### F-002 · Disable / rotate 请求体与 D-001 §2、i18n「动态码或恢复码」不一致

| 字段 | 值 |
|------|-----|
| level | **required**（med） |
| status | open |
| evidence | D-001 §2：Disable 为「code **或** recoveryCode 输入」；API `handler/mfa.go` L159–169 / L187–196 与 `service.go` `requireActiveSecondFactor` L267–275：非空 `recoveryCode` 走私钥恢复路径，否则校验 TOTP `code`。组件 disable 只发 `{code: disableCode.trim()}`（`mfa-manager.tsx` L137）；rotate 只发 `{recoveryCode: rotateRecovery.trim()}`（L146–148）。i18n 两处占位均为「当前动态码或恢复码」（`zh-CN.json` / `en-US.json` L510、L512）。用户按文案用恢复码停用 MFA → 服务端当 TOTP 验 → `MFA_INVALID`；用 6 位动态码轮换恢复码 → 当恢复码验 → 同样失败。`mfa-manager.test.tsx` 无 disable/rotate 用例，未钉住该契约。 |
| closure | — |
| 影响门禁 | S5 关门 / D-001 §2 自助解绑 |

丢失验证器后的自助停用是本目标存在的理由之一（用户裁决阻断 GOAL-017 关门，就是要个人中心 enroll/confirm/**disable**/rotate）。管理员 Reset MFA 仍在，不能代替本 UI 契约。最小充分：同一输入按形状分流（如 `/^\d{6}$/` → `code`，否则 → `recoveryCode`）并同时用于 disable 与 rotate；补单测（恢复码 disable、TOTP rotate）。

### F-003 · D-001 §4 测试缺口：生产加载链与分发/解绑流未钉住

| 字段 | 值 |
|------|-----|
| level | recommended（med） |
| status | open |
| evidence | D-001 §4 写了三类：render 分发（注册/未注册 fallback）、s5/代表页渲染 account 含 custom、MfaManager enroll→confirm→**disable**。实际：`render.test.ts` 只钉白名单/解析（L66–74），**没有** `render.tsx` 分发或「unknown custom component」断言；`resetCustomComponentsForTests`（`custom-components.ts` L34–36）无任何测试调用；`mfa-manager.test.tsx` 仅两流（enroll→confirm、unavailable）；s5 未 `renderAt("/account", …)`。E-003 L17 写「s5 分母渲染 account 页含 custom 节点降级安全」——document map 有 account.json，**没有**访问该页的断言。schema-keys 只保证 zh/en 键集相等，**不断言** `schema.account.mfa.*` 存在。 |
| closure | — |

不单独升级为 required（已由 F-001 / F-002 阻断）。修复 F-001 时必须补 D-VAL + `/account` 加载链，否则复审仍应 fail closed。

### F-004 · `01-decision.md` 信息表 I-001 仍为 `open`

| 字段 | 值 |
|------|-----|
| level | recommended（low） |
| status | open |
| evidence | `00-meta` 信息表 I-001 = **verified**（D-001 §1）；`01-decision.md` L17 仍 `open` / 「待确认」。权威以 `00-meta` + D-001 为准，**不**视为到期未闭环 required 信息项。 |
| closure | — |

## 必改项汇总

1. **F-001（required · high）**：`account.json` 必须通过运行时 `validatePageDocument`；个人中心须经 `loadPageDocument` 到达 renderer。未修前 **不得** S5 关门，**不得**据此回归关闭 GOAL-017。
2. **F-002（required · med）**：disable / rotate 必须兑现「code 或 recoveryCode」；i18n 与请求体一致；补契约测试。

无其他 required。I-001 无到期未闭环项（`00-meta` verified；F-004 仅为索引不同步）。

## 与既有意见的异同

- A-001（self · 立项/S1 pass）：同意方案层与 I-001 闭合。D-001 未写 D-VAL 约束；本意见不追溯方案为 fail，把缺口记在实施接入。
- A-002（self · S2–S4 pass）：同意组件、注册表、i18n、go 不暂挂、未改 MFA API。**不同意**可进 S5：self 未把 `account.json` 送进 `validatePageDocument`，也未核对 disable/rotate 请求体相对 D-001 §2 与 i18n。
- 不与 GOAL-017 A-007/A-008 冲突：A-007 F-001/F-002（required）仍为 fixed（本 scope 未重开）；A-007/A-008 将个人中心 UI 标 non-blocking 后，用户 D-005 **升级为阻断 GOAL-017 关门**并开本子目标。本目标按 D-001 验收，不能再把「UI 接不上」降回 non-blocking。
- 不把 D-002 go 判定改为失效。

## 结论 + 建议给编排器/用户的下一步

**verdict: fail**。renderer 契约与 MfaManager 组件在隔离测试中可核对；S5 不能过：custom 节点使生产 `/account` 整页 D-VAL 失败（F-001），且 disable/rotate 未兑现冻结的双因素输入契约（F-002）。E-002/E-003/A-002「account.json 接入 / 可进 S5」名不副实。

**不可关门。** 存在未合法闭合的 required findings。

**GOAL-017**：依赖链在 D-005 / E-001 / 本 `00-meta` 可追踪，但 **不能**根据本意见把 GOAL-017 标 `done`。`status: done` 仍由 `/govern` 在本目标 required 闭合且复审通过后执行。

建议 `/govern`：

1. 修 F-001 + F-002（含 D-VAL 与 disable/rotate 契约测试）；勿改本意见原文。
2. 再开独立复审（本 A-003 required 闭合）。F-003 随 F-001 测试一并关。
3. 未闭合前 **不得** 本目标 `status: done`，**不得** 回归关闭 GOAL-017。勿用 `progress: 3/5` 放行。

## 声明

本意见不修改 `status` / `progress` / goal-tree / `00-meta` / D-001 / D-002 正文。响应由 `/govern` 处理。保证等级 L0。
