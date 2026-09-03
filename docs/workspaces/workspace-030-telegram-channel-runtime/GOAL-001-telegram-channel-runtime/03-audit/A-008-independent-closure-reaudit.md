---
doc_type: goal-audit
id: A-008-independent-closure-reaudit
parent: GOAL-001-telegram-channel-runtime
date: 2026-09-03
source: independent
auditor: grok-4.6 (reasoning high)
scope: A-006 F-001/F-002 与 A-007 声称闭合的独立复审（不以 self 结论为证据）
audit_type: finding-closure
verdict: pass
open_required: 0
status: recorded
created: 2026-09-03
updated: 2026-09-03
version: 0.1.0
---

# A-008 · A-007 闭合复审（independent）

## 审计基本信息

| 字段 | 值 |
|------|-----|
| 被审目标 | [GOAL-001-telegram-channel-runtime](../00-meta.md) |
| 工作区 | `workspace-030-telegram-channel-runtime`（`root_goal` 匹配；`shared_materials_catalog: none`） |
| source | `independent` |
| auditor | grok-4.6（reasoning high） |
| 类型 | `finding-closure` |
| scope | A-006 **F-001 / F-002** 是否按源码与 **Fx 生产图** 闭合 |
| 对照 | [A-006](A-006-independent-closure-reaudit.md)；[A-007 self](A-007-a006-closure-response.md) 仅作声称清单 |
| 方法 | 只读现码。本会话复跑 `internal/channel/telegram`、`modules/channel/telegram`、`internal/composition`（含 `TestTelegramFxInjection_SameRuntime`，均 ok）。不改 status |
| verdict | **pass** |
| 开放 required | **0** |

---

## 范围与区间

用户：「再复审一下」。本审只核 A-006 两条 required 是否在当前仓库闭合。A-007 的 pass **不是**证据。不重开已降级的 Schema/Nav（A-006 已标 recommended）。

---

## F-001 · 闭合（经 Fx 图，不是手工传 tr）

A-006 要求：`tr *TelegramRuntime` 非 variadic 必选；删空 tr 时第二次 `newTelegramRuntime`；用 Fx 证明同一实例。

源码：

1. `newMux` / `newMuxWithExtraProviders` 最后一参为 `tr *TelegramRuntime`（`composition.go` L311、L338）。全仓无 `trs ...*TelegramRuntime`。
2. `newMuxWithExtraProviders` **不再** fallback 构造 runtime。启用时只用注入的 `tr.Webhook` / `tr.Manager`（L605–611）。
3. `NewApp` → `newAppWithOptions`：`fx.Provide(newTelegramRuntime, func(tr) Dispatcher, func(tr) Sender, newMux)`（L145–148）。dig 会注入非 variadic 的 `*TelegramRuntime`。
4. **`TestTelegramFxInjection_SameRuntime`**：`newAppWithOptions` + `fx.Populate(&injected)` + `fx.Populate(&mux)` + `app.Start`。在 **Populate 得到的** `injected.Dispatcher` 上 `RegisterCommand("status")`，再打 **Populate 得到的** mux webhook。断言 handler 被调用。没有把 tr 传进 `newMux`。若仍是双实例，此测试必红。本会话 composition 包测试 ok。

未启用路径：`newTelegramRuntime` 返回 DisabledDispatcher / DisabledSender（L855–858）；`Webhook` 为 nil，不挂路由。符合 D-002 §1。

`ResolveTelegramPorts` 仍是独立工厂（注释写明非 Fx 生产路径）。可接受。

**判定：F-001 fixed。**

---

## F-002 · 闭合（主密钥离开源码；init fail-closed）

A-006 要求：可旋转 master key；`initPersistence` 失败冒泡。

源码：

1. grep 无 `defaultMasterKey` / `schema-ui-core:channel:telegram:master-key`。
2. `NewRuntimeManager(..., masterKey []byte, ...)` 空 key 直接 error（`runtime.go` L45–47）。
3. `newTelegramRuntime` 走 `mail.LoadOrCreateMasterKey(cfg.TelegramMasterKey, path)`（L822–831）：`TELEGRAM_MASTER_KEY` 或密钥文件；缺省路径 `telegram-master.key` 在 DB 目录旁（与 mail 相同的默认，见下 recommended）。
4. `initPersistence` 返回 error；`NewRuntimeManager` 透传（L65–68、L77–120）。读库/解密失败不再 `_ =`。
5. `Update` 仍先写库再改内存。重启测试仍断言 `GetToken()` / `GetSecret()`（`composition_telegram_test.go` L205–241），并用 `TelegramMasterKey: "test-master-key"` 保证加解密同一把钥匙。

YAML 无 `telegram.master_key` 明文槽，口令只走 env（`config.go` L406–410、L738）。合理。

**判定：F-002 fixed。**

---

## 仍开放的 recommended（不阻断 A-006 闭合）

| ID | 状态 | 说明 |
|----|------|------|
| R-004 | **仍 open** | HTTP 200 且 body 非 JSON 时 `Unmarshal` 失败仍当成功（`http_sender.go` L141–147）。A-006/A-007 未升格。 |
| R-007 | residual | Allow/Record TOCTOU。需 `/govern` 书面确认用户接受。 |
| 默认密钥文件与 DB 同目录 | **recommended / 新 R-009** | `telegram-master.key` 默认 `filepath.Dir(DBPath)`。备份数据目录会同时带走密文与钥匙。Mail W13 F-017 同一形状；可用 `TELEGRAM_MASTER_KEY_PATH` 挪走。不把 F-002 打回 open。 |
| Schema/Nav/tab | recommended | A-006 已降级。判据 #5 仍是 API-only，除非用户书面收窄或补 UI。 |
| PATCH 空串 vs seed | informational | `initPersistence` 仅在解密结果非空时覆盖内存；清空 token 后重启可能回到 env seed。 |

---

## Findings

本轮 **无新 required**。

### R-004 · HTTP 200 非 JSON 仍当成功（维持）

- **严重度**：low · recommended · open
- **证据**：`http_sender.go` L141–147。

### R-009 · 默认 master key 文件与数据库同目录

- **严重度**：low · recommended · open
- **证据**：`composition.go` L824–826。缓解：设 `TELEGRAM_MASTER_KEY` 或 `TELEGRAM_MASTER_KEY_PATH`。

---

## 必改项汇总

| ID | 本审 |
|----|------|
| F-001 | **closed（fixed）** — 非 variadic；无 fallback；Fx Populate 同实例测试绿 |
| F-002 | **closed（fixed）** — `LoadOrCreateMasterKey`；init 失败 fail-closed |

开放 required = **0**。本意见仍不修改 `00-meta` status。

---

## 与既有意见的异同

| 既有 | 本审 |
|------|------|
| A-006 F-001/F-002 open | **闭合**。A-006 指出的 dig variadic 与源码常量钥匙，在当前文件中已不存在，且有经 `NewApp` 图的测试。 |
| A-007 self pass | 本审 **独立核对后同意** 两条 required 的 fixed。不把 A-007 当证据，把源码与本会话测试当证据。 |
| A-004/A-006 对「测试绿 ≠ 装配通」 | 本轮测试走 `app.Start` + Populate，不再是那类假绿。 |

---

## 结论 + 建议给编排器/用户的下一步

**verdict: pass。** A-006 的两条 required 在现码上成立。建议 `/govern` 把 F-001/F-002 记为 **fixed**。R-004、R-009、Schema/Nav、R-007 保持 recommended/residual，不阻断本次闭合。

---

## 声明

本意见不修改 status/progress；响应由 `/govern` 处理。
