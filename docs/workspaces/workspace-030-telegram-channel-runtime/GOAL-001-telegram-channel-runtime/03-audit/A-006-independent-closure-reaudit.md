---
doc_type: goal-audit
id: A-006-independent-closure-reaudit
parent: GOAL-001-telegram-channel-runtime
date: 2026-09-03
source: independent
auditor: grok-4.6 (reasoning high)
scope: A-004 F-001/F-002 与 A-005 声称闭合的独立复审（不以 self 结论为证据）
audit_type: finding-closure
verdict: fail
open_required: 2
status: recorded
created: 2026-09-03
updated: 2026-09-03
version: 0.1.0
---

# A-006 · A-005 闭合复审（independent）

## 审计基本信息

| 字段 | 值 |
|------|-----|
| 被审目标 | [GOAL-001-telegram-channel-runtime](../00-meta.md) |
| 工作区 | `workspace-030-telegram-channel-runtime`（`root_goal` 匹配；`shared_materials_catalog: none`） |
| source | `independent` |
| auditor | grok-4.6（reasoning high） |
| 类型 | `finding-closure` |
| scope | A-004 **F-001 / F-002** 是否在 **Fx 生产图** 与 catalog/加密上真正闭合 |
| 对照 | [A-004](A-004-independent-closure-reaudit.md)；[A-005 self](A-005-a004-closure-response.md) 仅作声称清单 |
| 方法 | 只读现码。对照 `go.uber.org/dig` v1.19.0（fx v1.24.0 依赖）对 variadic 的处理。本会话复跑 telegram 三包、`internal/store`、`cmd/schema-ui`（均 ok）。不改 status |
| verdict | **fail** |
| 开放 required | **2**（F-001、F-002 仍 open） |

---

## 成果（源码，不是 A-005）

相对 A-004，下列是真进展：

| 主张 | 证据 |
|------|------|
| 存在 `TelegramRuntime` + `newTelegramRuntime`；未启用返回 Disabled* | `composition.go` L808–843 |
| `fx.Provide(newTelegramRuntime)` 以及从该对象取出 `kernel.TelegramDispatcher` / `TelegramSender` | `composition.go` L137–139 |
| 手工把 **同一个** `tr` 传入 `newMuxWithExtraProviders` 时，Register 的命令会被真实 webhook 分到 | `TestTelegramChannelComposition_RealWebhookMount` |
| 运行时 `CREATE TABLE IF NOT EXISTS` 已删除 | grep `channel/telegram` 无命中 |
| 迁移 v66 `telegram_config` 进 compiled catalog；fingerprint head=66 | `modules/channel/telegram/migration/migration.go`；`compiled/persistence.go` L21/44；`identity.go` L93/113；`identity_test.go` L117；`migrate_test.go` 基线 66 |
| `Update` 先加密写库成功再改内存 | `runtime.go` L123–161 |
| 重启测试走 `Update` 后断言 `GetToken()` / `GetSecret()` | `composition_telegram_test.go` L205–241 |
| GET 不再回显末 4 位 | `runtime.go` L20–26、L175–181 |
| export `sensitiveFields` 登记 telegram 键 | `configpkg.go` L243–244 |

---

## F-001 · 仍 open：Fx **忽略** variadic，生产图仍是两套 runtime

A-005 声称：`newMux` 接收 `tr *TelegramRuntime`，「直接将 Fx 注入的同一 `tr.Dispatcher` 挂到 webhook」。

源码：

```go
func newMux(..., trs ...*TelegramRuntime)
func newMuxWithExtraProviders(..., trs ...*TelegramRuntime)
```

`NewApp` 里 `fx.Provide(newMux)`。本仓 `go.uber.org/fx v1.24.0` → `dig v1.19.0`。dig 文档（`doc.go` L87–97）与实现（`param.go` L114–122）：

> Constructors that accept a variadic number of arguments are treated as if they don't have those arguments.  
> The constructor will be called with all other dependencies and **no variadic arguments**.

因此生产路径：

1. `newTelegramRuntime` 构造 **实例 A**，导出给其它模块的 `kernel.TelegramDispatcher` / `TelegramSender`。
2. `newMux` 的 `trs` **恒为空**（variadic 被 dig 丢掉）。
3. `newMuxWithExtraProviders` 走 `else { tr = newTelegramRuntime(...) }`（`composition.go` L598–603），再构造 **实例 B** 给 webhook/settings。

业务模块 `RegisterCommand` 打在 A 上；入站 `Dispatch` 走 B。这与 A-004 的双 dispatcher **同一缺陷**，只是多了一层看起来像单例的 `fx.Provide`。

`TestTelegramChannelComposition_RealWebhookMount` **显式传入** `tr`（L319），不经过 Fx，所以绿。它不能证明 `NewApp`。

`ResolveTelegramPorts` 仍是第三条 `newTelegramRuntime()` 工厂（L845–848）。

**判定：F-001 未闭合。** 修复必须把 `*TelegramRuntime` 写成 **非 variadic 必选参数**（未启用也提供 Disabled stub），禁止 `newMux` 在空 tr 时再 `newTelegramRuntime`。

---

## F-002 · 仍 open：catalog/重启测试已落地；「加密」密钥写死在源码

A-005 声称 AES-GCM at-rest、对标 mail。

已成立：v66 进台账、无运行时 DDL、`bot_token_enc`/`webhook_secret_enc`、`Update` 先写库、重启测 `GetToken()`。

未成立：

1. **主密钥是编译期常量。** `runtime.go` L15–16、L61：`sha256("schema-ui-core:channel:telegram:master-key:v1")`。任何拿到二进制或本仓库的人都能解密 `telegram_config`。Mail 用 `LoadOrCreateMasterKey` + 文件路径（W13 F-017）。这不是 at-rest 加密，是带固定钥匙的编码。DB 备份泄漏模型下等同明文。
2. **`initPersistence` 仍 `_ = runner.Run`（L73）。** 读库/解密失败静默停留在 seed。`Update` 已 fail-closed，启动重载没有。
3. **Schema/Nav/tab 仍无。** A-004 要求做或书面收窄判据 #5。A-005 两者都没有。本审把此项保持 recommended，不单独升格。

**判定：F-002 未闭合。** 持久化机制本身可用；闭合「对标 mail 的加密」不成立。

---

## Recommended（对 A-005 台账）

| ID | A-005 | 本审 |
|----|-------|------|
| R-001 / R-003 / R-008 | closed | **同意闭合** |
| R-002 | closed | **并入 F-001**。有真实 mux 测试，但未走 `NewApp`/Fx；生产双实例未覆盖 |
| R-004 | closed | **仍部分**：`Unmarshal` 失败当成功（`http_sender.go` L141–147） |
| R-005 | closed | **同意闭合**（无 masked 字段） |
| R-006 | closed | **同意闭合（export 路径）**。`cfgTree` 仍无 telegram 段，值不会进包；`sensitiveFields` 已登记。YAML 本地明文仍可存在，不作必改 |
| R-007 | residual | **同意**；需 `/govern` 书面确认 |

---

## Findings

### F-001 · Fx 未把 TelegramRuntime 注入 newMux（variadic 被 dig 丢弃）

- **严重度**：med（相对进程正确性：高影响；维持 A-002 编号）
- **建议**：required
- **状态**：open
- **证据**：`composition.go` L137–140、L285–307、L598–603、L816–848；`go.uber.org/dig@v1.19.0/doc.go` L87–97、`param.go` L114–122；测试 L302–320 手工传 tr。
- **建议修复**：`newMux` / `newMuxWithExtraProviders` 的最后一参改为 `tr *TelegramRuntime`（非 `...`）。未启用也注入 Disabled stub。删除空 tr 时的第二次 `newTelegramRuntime`。加一条 **经 `NewApp`/`fx`** 的指针相等或命令分发测试。

### F-002 · 加密主密钥硬编码；启动重载吞错

- **严重度**：med
- **建议**：required
- **状态**：open
- **证据**：`runtime.go` L15–16、L61、L71–73；对比 `composition.go` `newMailRuntime` 的 `LoadOrCreateMasterKey`。
- **建议修复**：复用 mail master key 或独立 `TELEGRAM_MASTER_KEY` / 文件；禁止源码常量。`initPersistence` 失败要返回 error 并在组合根 fail-closed。Schema/Nav 做或用户书面收窄 #5。

### R-004 · HTTP 200 非 JSON 仍当成功

- **严重度**：low · recommended · open
- **证据**：`http_sender.go` L141–147。

---

## 必改项汇总

| ID | 建议 | 一句话 |
|----|------|--------|
| **F-001** | required · **仍 open** | 去掉 variadic；Fx 注入的那一份必须就是 webhook 的那一份 |
| **F-002** | required · **仍 open** | 可旋转的 master key；init 失败 fail-closed |

A-005 的 `closed` / required=0 **不能成立**。catalog 与重启测试不必推倒重来。

---

## 与既有意见的异同

| 既有 | 本审 |
|------|------|
| A-004 F-001 双 dispatcher | **维持**。工厂已统一，**接线被 dig 丢掉**。手工测试不能代替 Fx。 |
| A-004 F-002 运行时 DDL / 明文 / 假测试 | DDL 与假测试 **已修**；明文变成固定钥匙加密，**仍不算闭合**。 |
| A-005 self pass | **驳回为证据**。 |

---

## 结论 + 建议给编排器/用户的下一步

**verdict: fail。** 本轮比 A-004 实在：v66、persist-first、GetToken 重启、真实 webhook 手工装配都在。不能放行的是：生产 `NewApp` 仍构造两套 TelegramRuntime；「加密」钥匙在 Git 里。

建议 `/govern`：

1. F-001：`tr *TelegramRuntime` 必选参数，删 fallback `newTelegramRuntime`；用 Fx 测同一指针。
2. F-002：主密钥离开源码；init 失败要冒泡。
3. 不要再把「测试里传入 tr」写成「Fx 已注入」。

---

## 声明

本意见不修改 status/progress；响应由 `/govern` 处理。
