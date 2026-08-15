---
id: A-004
goal: GOAL-018-mfa-manager-ui
title: S5 关门 A-003 required 闭合复审
date: 2026-08-15
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
scope: A-003 F-001/F-002 required 闭合验证（account.json 运行时 D-VAL + splitMFAInput disable/rotate + 契约/回归 + 本地扩展标注 / 协议 pin）
audit_type: finding-closure
verdict: pass
status: recorded
parent: GOAL-018-mfa-manager-ui
created: 2026-08-15
updated: 2026-08-15
version: 1.0.0
---

# A-004 · 独立复审意见（S5 · A-003 required 闭合）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型** / **scope**：finding-closure · A-003 F-001（account.json D-VAL / `/account` 可达 renderer）、F-002（disable/rotate「code 或 recoveryCode」）闭合验证；顺带核对 F-003 测试缺口、F-004 I-001 索引、本地扩展标注与协议 pin
- **verdict**：**pass**

## 范围与区间

- **工作区**：`workspace-011-admin-functional-modules`（`workspace.md` 绑定 Root=`GOAL-001-admin-functional-modules`、`canonical_scope`、`plan_refs`/`primary_plan`=VP-011 已校验；`shared_materials_catalog: none`）。未读取或比较其他工作区目标状态。同区 GOAL-017 仅作回归关门依赖只读对照。
- **已通读**：本目标 `00-meta`、`01-decision.md`、D-001、D-002、`02-execution.md`、E-001～E-003、`03-audit.md`、A-001～A-003。父目标裁决：GOAL-017 D-005、`03-audit.md`、A-007/A-008。
- **代码核对**：`runtime-schema-validate.ts`（运行时 overlay）、`docs/schemas/node.schema.json`（properties 无 `component`，`additionalProperties: false`）、`docs/schemas/component-registry.json`（无 `custom` 条目）、`account.json` L148–151、`load-page.ts` L105–112、`load-page.test.ts` L67–82、`mfa-manager.tsx` `splitMFAInput` / disable / rotate、`mfa-manager.test.tsx`、`handler/mfa.go` L159–169 / L187–196、`service.go` `requireActiveSecondFactor` L257–276、`render.ts` custom 解析、`01-decision.md` I-001。
- **本轮复跑**（2026-08-15）：
  - 独立 Ajv（`allErrors` / `strict:false` / `validateSchema:false`，与 `runtime-schema-validate.ts` 一致）对**真实** `apps/api/internal/modules/account/schema/account.json`：未叠加 overlay → `ok:false`，`instancePath=/body/children/3`，`additionalProperty=component`；叠加与运行时相同的 `component` overlay → **`ok:true`，0 errors**。
  - `apps/web` `vitest run` → **56 文件 976/976 绿**（相对 A-003 记录的 974 增 2：custom D-VAL + disable/rotate 用例）。
  - `apps/web` `tsc --noEmit -p tsconfig.json` → **0**。
  - `apps/api` `go test -p 1 ./...` → **全绿**。
- **covered**：A-003 F-001/F-002 关闭证据是否真实、充分、可重复核对；真实 account.json D-VAL；splitMFAInput 与 handler/service 二选一一致；契约测试是否钉住；本地扩展是否改 pin；F-003/F-004；GOAL-017 回归依赖。
- **excluded**：不改 `status` / `progress` / goal-tree / `00-meta` / D-001 / D-002 正文。
- **保证等级**：L0。

## 成果（有证据）

| A-003 主张 / 本复审核对项 | 闭合证据 |
|---------------------------|----------|
| F-001：生产路径仍是 `loadPageDocument` → `validatePageDocument` → 失败抛 `PAGE_SCHEMA_INVALID` | `load-page.ts` L105–112 未改；修复点在校验器，不是绕过 D-VAL |
| F-001：真实 `account.json` 现可通过运行时 D-VAL | 本轮独立 Ajv：overlay 后 `ok:true`。custom 节点仍为 `{type:"custom", id:"mfa-manager-block", component:"mfa-manager"}`（`account.json` L148–151） |
| F-001：未改上游 pin 字节 | `git status`：`docs/schemas/node.schema.json` 与 `component-registry.json` **未改**。`node.schema.json` L16 起仍无 `component`；`additionalProperties: false` 仍在 L71 |
| F-001：本地扩展标注充分 | `runtime-schema-validate.ts` L35–53：校验期 spread 增 `component`，注释写明 GOAL-018 本地扩展、**不是**上游协议变更、pin 工件保持 byte-identical |
| F-001：additionalProperties 仍拒未知键 | 未叠加 overlay 时同一文档仍只报 `component`；overlay **只**放行该一键。其它顶层未知属性仍会被拒 |
| F-001：契约测试钉住 custom 节点 | `load-page.test.ts` L67–82：`validatePageDocument` 对含 `{type:"custom", component:"mfa-manager"}` 的合成文档 `ok:true`（本轮 11/11 绿） |
| F-002：`splitMFAInput` 6 位→`code`，其余→`recoveryCode` | `mfa-manager.tsx` L72–79 |
| F-002：disable / rotate 均走该识别 | disable L147 `postJSON(..., splitMFAInput(disableCode))`；rotate L156 `splitMFAInput(rotateRecovery)` |
| F-002：与 handler/service 二选一一致 | `handler/mfa.go` L160–161 / L188–189 解 `code`+`recoveryCode`；`requireActiveSecondFactor` L267–275：非空 `recoveryCode` 走私钥恢复，否则验 TOTP `code` |
| F-002：i18n 仍为「动态码或恢复码」 | `zh-CN.json` / `en-US.json` L510、L512 |
| F-002：恢复码停用用例存在 | `mfa-manager.test.tsx` L110–121：输入 `ABCDEFGH`（非 6 位）后 UI 回到 Enable MFA；本轮 3/3 绿 |
| F-004：`01-decision.md` I-001 | L17 现为 **verified**（与 `00-meta` / D-001 §1 一致） |
| 回归 | web 976/976；`tsc --noEmit` 0；go `./...` 全绿 |
| D-002 go 不暂挂 / pin 不变 | 成立：无新 capability、无 pin 文件改动；扩展在运行时校验层 |

## 对照 A-003 required 闭合标准

| 标准 | 状态 | 证据 |
|------|------|------|
| `account.json` 通过运行时 `validatePageDocument` | **闭合** | 本轮对真实文档独立 Ajv + overlay → `ok:true` |
| 个人中心经 `loadPageDocument` 到达 renderer（不再整页 `PAGE_SCHEMA_INVALID`） | **闭合**（门禁） | D-VAL 是该失败的唯一闸门；真实文档现 `ok`。套件未用真实 `account.json` 走 `loadPageDocument` / `renderAt("/account")`，见 F-001 recommended |
| 不要静默绕过 D-VAL | **闭合** | 仍走 `validatePageDocument`；只把 `component` 登记为本地允许属性 |
| 不要把 `component` 写进 `node.schema.json` 却声称无协议变更 | **闭合** | **未**改 `node.schema.json`。编排器口头「properties 增 component」与落盘不符；实际是校验期 overlay，与 D-002 / 注释一致，且优于改 pin 文件 |
| disable / rotate 兑现「code 或 recoveryCode」 | **闭合** | 同一 `splitMFAInput`；形状与 handler 二选一一致 |
| i18n 与请求体一致 | **闭合** | 两处占位仍为双因素文案；请求体按形状分流 |
| 补契约测试（恢复码 disable、TOTP rotate） | **部分**（不恢复 required） | 恢复码 disable 有 UI 断言；rotate 半段被 `if (rotateInput !== null)` 吞掉，且两流都不断言 fetch body。实现已对，见 F-001 recommended |
| I-001 无到期未闭环 required | **闭合** | `00-meta` + `01-decision.md` 均为 verified |

## Findings

### F-001 · 契约测试未钉真实 `account.json` / `/account` 加载链；rotate 半段为空跑

| 字段 | 值 |
|------|-----|
| level | recommended（med） |
| status | open |
| evidence | A-003 F-003 写明：修 F-001 时须补 `validatePageDocument(account.json) === ok` 与 `/account` 经 `loadPageDocument` 的渲染断言。现 `load-page.test.ts` L67–82 只用**合成**文档（无真实 `account.json` 的 form/table 兄弟节点）。`representative-pages.test.tsx` `MIGRATED_PAGE_IDS` 仍无 `account`；`s5-denominator-render.test.tsx` 仍把 account.json 放进 map，用例路由为 `/roles` `/settings`。`mfa-manager.test.tsx` L95–137：disable 后 `refresh` 把 status 置 `enabled:false`，rotate 表单卸掉；L124 `if (rotateInput !== null)` 使 6 位轮换断言不执行，测试仍绿；两流均未断言 `fetch` body 为 `{recoveryCode:"ABCDEFGH"}` / `{code:"123456"}`。`render.test.ts` 仍只钉白名单/解析，无分发/未注册 fallback。 |
| closure | — |
| 影响门禁 | **不阻断 S5**。本复审已对真实 `account.json` 独立跑通运行时 D-VAL；F-002 实现与 handler 一致。本条是回归钉不充分，不是生产闸门仍失败。 |

最小充分（非关门前置）：测试直接 `validatePageDocument(JSON.parse(readFileSync(account.json)))`；拆成两个用例分别 disable（恢复码）与 rotate（6 位），并断言请求 JSON；可选 `renderAt("/account")`。

### F-002 · 声称的 `component-registry.json` `custom` 条目未落盘

| 字段 | 值 |
|------|-----|
| level | recommended（low） |
| status | open |
| evidence | 编排器称 registry 增 `custom`。`docs/schemas/component-registry.json` 收于 L2398–2400，**无** `"custom"` 键；`git status` 该文件未改。`node.schema.json` L10 描述「type 需存在于 component-registry.json」，但 D-VAL 的 `type` 是自由字符串，registry 缺失不导致 `PAGE_SCHEMA_INVALID`。 |
| closure | — |
| 影响门禁 | 无。工具/文档完备性，不是运行时门禁。 |

## 必改项汇总

无。A-003 required（F-001 high、F-002 med）均已合法闭合（`fixed`，可重复核对）。

A-003 F-004（recommended）已随 `01-decision.md` I-001=`verified` 闭合。A-003 F-003 的残余测试缺口记入本意见 F-001 recommended，**不**升级为 required。

## 与既有意见的异同

- A-003（independent · S5 fail）：同意当时复现（未 overlay 时真实 account.json 仍 `additionalProperty=component`）。**不同意**修复必须把键放进 `props`：校验期 overlay 同样不改 pin、不绕过 D-VAL，且与 D-002「本地扩展」一致。A-003「不要改 `node.schema.json` 却声称无协议变更」——本轮**没有**改该文件。
- A-001 / A-002（self）：不重开。A-002 当时未跑 D-VAL 的缺口已由本轮独立 Ajv 补上。
- 不与 GOAL-017 A-007/A-008 冲突：A-007 F-001/F-002 仍为 fixed（A-008 pass）。用户 D-005 把个人中心 UI 升级为阻断 GOAL-017 关门并开本子目标；本目标 A-003 required 现已闭合，**可以**回归关闭 GOAL-017。`status: done` 仍由 `/govern` 执行。
- 不把 D-002 go 判定改为失效。协议 pin 文件未动。

## 结论 + 建议给编排器/用户的下一步

**verdict: pass**。A-003 required 已合法闭合：真实 `account.json` 经与运行时相同的 Ajv 选项 + 本地 `component` overlay 为 `ok:true`；disable/rotate 均按 6 位/`recoveryCode` 分流，与 `requireActiveSecondFactor` 一致。web 976/976、`tsc --noEmit`、go 全量本轮全绿。

**可关门。** 无未合法闭合的 required findings。无到期 required 信息项。

**GOAL-017**：依赖链在 D-005 / E-001 / 本 `00-meta` 可追踪；A-007 required 已由 A-008 闭合；本目标 A-003 required 现已闭合。**可以**回归关门 GOAL-017。`status: done` 由 `/govern` 执行（含该目标已登记的波次级 e2e/冒烟，若编排器仍将其列为该目标 S5 事实，须在关门记录中写明已跑或明确不在本复审范围）。

建议 `/govern`：

1. 将本目标 A-003 F-001/F-002 标 `fixed`（闭合证据 = 本 A-004）；本目标 `status: done`，同步 goal-tree。
2. 回归关闭 GOAL-017（`status: done` + goal-tree）。
3. F-001 recommended（真实 account.json / `/account` 渲染 / rotate body 断言）可另排，**不**阻断本次关门。

## 声明

本意见不修改 `status` / `progress` / goal-tree / `00-meta` / D-001 / D-002 正文。响应由 `/govern` 处理。保证等级 L0。
