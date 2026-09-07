---
doc_type: goal-audit
id: A-004-independent-closure-reaudit
parent: GOAL-001-telegram-channel-runtime
date: 2026-09-03
source: independent
auditor: grok-4.6 (reasoning high)
scope: A-002 F-001/F-002 与 R-001～R-008 声称闭合的独立复审（不以 A-003 self 结论为证据）
audit_type: finding-closure
verdict: fail
open_required: 2
status: recorded
created: 2026-09-03
updated: 2026-09-03
version: 0.1.0
---

# A-004 · A-002 闭合复审（independent）

## 审计基本信息

| 字段 | 值 |
|------|-----|
| 被审目标 | [GOAL-001-telegram-channel-runtime](../00-meta.md) |
| 工作区 | `workspace-030-telegram-channel-runtime`（`root_goal` 匹配；`shared_materials_catalog: none`） |
| source | `independent` |
| auditor | grok-4.6（reasoning high） |
| 类型 | `finding-closure` |
| scope | A-002 **F-001 / F-002** 是否按源码闭合；R-001～R-008 声称整改是否名实相符 |
| 对照 | [A-002](A-002-independent-design-code-audit.md)；[A-003 self](A-003-independent-audit-response.md) 仅作「声称清单」，不作证据 |
| 方法 | 只读现码与测试。本会话复跑 `internal/channel/telegram`、`modules/channel/telegram`、`internal/composition`（`-count=1`，均 ok）。**不**改 status / progress / goal-tree |
| verdict | **fail** |
| 开放 required | **2**（F-001、F-002 仍 open） |

### 范围与区间

- **在 scope**：A-002 两条 required 的闭合证据；recommended 项是否兑现；A-003 台账与代码是否一致。
- **不在 scope**：改 Charter/VP/status；把 recommended 升格为关门阻断（除非用户书面升级）。
- **P-005**：无新到期 required 信息项阻断本意见本身。
- **共享资料**：无。

---

## 范围与区间

用户：「已做修改，复审一下」。本审只核 **当前仓库文件**。A-003 写「F-001/F-002 fixed、开放 required = 0」**不构成闭合**。

---

## 成果（有证据 · 来自源码）

相对 A-002，下列改动**真实存在**，且包级测试为绿：

| 主张 | 证据 |
|------|------|
| `DisabledSender.Send` 返回 `kernel.ErrTelegramDisabled` | `disabled.go` L19–22；`http_sender_test.go` `TestDisabledSenderAndDispatcher` |
| `DisabledDispatcher.Register*` 成功空操作 | `disabled.go` L34–48 |
| `ResolveTelegramPorts` 在**被直接调用时**能返回 disabled 或 live 实现 | `composition.go` L808–817；`TestResolveTelegramPorts_EnabledAndDisabled` |
| 启用模块时 live `NewRuntimeManager(..., st)` 把 store 传入 | `composition.go` L595 |
| `RuntimeManager.Update` 尝试 `INSERT … ON CONFLICT DO UPDATE` | `runtime.go` L120–147 |
| webhook secret 先 SHA-256 再 32 字节 `ConstantTimeCompare` | `webhook.go` L125–127 |
| Dispatch 失败打 `slog.Warn` | `webhook.go` L221–223 |
| HTTP 200 + `"ok": false` 当失败，有测试 | `http_sender.go` L139–144；`TestHTTPSender_Status200_ButOKFalse` |
| `RuntimeStatus` 含 `token_set` / `secret_set` 与 `captured_messages_count` | `runtime.go` L15–22、L160–167 |

这些是**局部正确的补丁**。它们**不等于** F-001/F-002 在进程装配上已闭合。

---

## 对照 A-002 必改项（闭合判定）

### F-001 · 进程级端口装配 — **仍 open（假闭合）**

A-003 声称：`ResolveTelegramPorts` 「供后续业务模块装配」。

源码事实：

1. **进程入口未调用该函数。** `NewApp` 的 `fx.Provide` 列表（`composition.go` L118–139）没有 Telegram 端口。`newMux` 启用分支（L591–607）仍**就地** `NewDispatcher()` / `NewHTTPSender(...)`，局部变量只注入 webhook。
2. **全仓唯一调用点是测试。** grep `ResolveTelegramPorts`：定义在 `composition.go` L808，调用仅 `composition_telegram_test.go`。
3. **即使业务模块自己调用 helper，也不是「同一实例」。** helper 在 enabled 路径再 `NewDispatcher()` 一次。webhook 用的是 L593 的实例；`RegisterCommand` 打在 helper 返回值上，webhook `Dispatch` 看不见。这比 A-002 时更危险：看起来有进程级 API，实际是第二套 dispatcher。
4. 模块未启用时，live `newMux` **仍然不**安装 no-op 到任何可注入槽。disabled stub 只在有人显式调用 helper 时存在。

A-002 要求的三件事（组合根持有端口、未启用 no-op、启用时与 webhook **同一**实例）**一条都没在 Fx 图里落地**。disabled 类型本身写对了，装配没接通。

**判定：F-001 未闭合。** 测试绿只证明 helper 的单元行为。

### F-002 · Admin 设置持久化 — **仍 open（半套且绕开台账）**

A-003 声称：`telegram_config` 表 + 启动重载 + 跨重启测试。

源码事实（正向）：

- live 路径 `NewRuntimeManager(..., st)`（L595），`Update` 会写库。原先「只改内存」在**这一条调用链**上不再成立。

源码事实（不足，故不能标 fixed）：

1. **绕开 compiled persistence catalog。** `channel.telegram` 的 `CompiledPersistence()` 仍返回 `nil, nil`（`provider.go` L50–52）。表由 `initPersistence` `CREATE TABLE IF NOT EXISTS` 在运行时创建（`runtime.go` L69–74）。本仓冻结：模块表只经 `compiled.PersistenceCatalog()` / `schema_migrations`。运行时 DDL 不进台账、无 checksum、无版本号。后续若补 catalog 迁移，会撞上「表已在、台账未记」。
2. **密钥明文落库。** `bot_token` / `webhook_secret` 为 `TEXT NOT NULL`，无加密。Mail 先例是加密列 + master key。A-002 建议「对标 mail」的持久化形态未兑现。
3. **失败语义 fail-open。** `initPersistence` 对 `runner.Run` 用 `_ =`（L68）：建表/读库失败时静默停留在 seed 内存，调用方不知道没持久化。`Update` 先改内存再写库（L125–144）：写库失败时 handler 500，但进程已用新 token/secret，与 DB 分叉。
4. **「跨重启」测试未证明声称。** `TestTelegramRuntime_PersistenceAcrossRestart`：用 SQL `UPDATE` 改行，不用 `RuntimeManager.Update`；重启后只再 `SELECT` 同一行，**不断言**新 `RuntimeManager.GetToken()` / `GetSecret()` 读到持久值。任何「表还在」的 SQLite 都能绿。
5. **设置面仍无 Schema/Nav/tab。** A-002 把这一点并进 F-002。代码与前端仍无 telegram 页。可降为 residual，但不能在未书面收窄时宣称 #5 完整。

**判定：F-002 未闭合。** 有写库意图，没有合法持久化面。

---

## Recommended 项（独立核对 A-003 台账）

| ID | A-003 声称 | 本审 | 说明 |
|----|------------|------|------|
| R-001 恒时比较 | fixed | **可接受闭合** | SHA-256 后再 32 字节 compare，长度侧信道被压到 hash。残留：hash 输入长度仍有微小时序，不作必改。 |
| R-002 真 webhook 装配测试 | fixed | **仍 open** | `TestTelegramChannelComposition` 仍是 dummy 200（L80–100）。新增测试不经过 `newMux`/`newServer` 真 `WebhookHandler`。 |
| R-003 Dispatch 吞错 | fixed | **可接受闭合** | 仍 200（合同合法 Update），但有 `slog.Warn`。运维盲区降为有日志。 |
| R-004 `ok` 字段 | fixed | **部分** | 200+`ok:false` 已测。`json.Unmarshal` 失败时 `err == nil` 才检查，`ok` 缺省 false 反而是安全的；非 JSON 且 Unmarshal 失败则当成功（L141–147）。留 recommended。 |
| R-005 不回显片段 | fixed | **未闭合** | 加了 `token_set`/`secret_set`，**仍输出** `token_masked`/`secret_masked`（末 4 位）。A-003「防止秘密片段回显」与 `maskSecret` 代码矛盾。 |
| R-006 YAML/export | fixed | **未闭合** | `.env.example` 增加注释键不等于 YAML 拒绝明文，也不等于 `sensitiveFields` 登记 telegram 键。 |
| R-007 Allow/Record TOCTOU | accepted-residual | **同意残余** | 属 VP-027 端口形状。本审不升格。闭合合法性由 `/govern` 确认用户书面接受。 |
| R-008 字段名 | fixed | **可接受闭合** | `captured_messages_count` 已有；多余别名 `captured_count` 无害。 |

---

## Findings

### F-001 · 进程级 Telegram 端口仍未进入组合根（A-002 未闭合）

- **严重度**：med
- **建议**：required
- **状态**：open
- **描述**：`ResolveTelegramPorts` 是测试用工厂，不是进程单例。live `newMux` 另造 dispatcher/sender。业务模块无法 Register 到 webhook 使用的那一份。disabled stub 未挂入未启用时的组合图。
- **证据**：`composition.go` L118–139、L591–607、L805–817；grep 调用仅测试文件。
- **建议修复**：`newMux`（或 `fx.Provide`）只构造**一次** dispatcher/sender/runtime；启用与否都返回同一对端口；webhook 与后续模块构造函数吃同一指针。禁止第二条 `NewDispatcher()`。

### F-002 · Admin 设置持久化未进入 compiled catalog，且明文、失败静默、测试未证重载

- **严重度**：med
- **建议**：required
- **状态**：open
- **描述**：写库存在，但不构成合法闭合：运行时 DDL 绕开全局迁移台账；token/secret 明文；init 吞错；Update 内存/库可分叉；跨重启测试不断言 `GetToken`。Schema/Nav 仍缺。
- **证据**：`runtime.go` L66–104、L120–147；`provider.go` L50–52；`compiled/persistence.go` 无 telegram；`composition_telegram_test.go` L186–239。
- **建议修复**：
  1. 用 `CompiledPersistence` 增加带 version/checksum 的 `telegram_config` 迁移（sqlite+postgres）；删除 `CREATE TABLE IF NOT EXISTS`。
  2. 密钥加密 at-rest（对标 mail），GET 只给 `token_set`/`secret_set`。
  3. init/Update 失败 fail-closed；先持久化成功再改内存，或写失败回滚内存。
  4. 测试：`Update` → close store → 新 `RuntimeManager` → `GetToken()`/`GetSecret()` 等于写入值；seed 不得覆盖非空 DB。
  5. Schema/Nav 要么做，要么用户书面把判据 #5 收窄为 API-only。

### R-002 · composition 仍未驱动真实 WebhookHandler（仍 open）

- **严重度**：low · recommended · open
- **证据**：`composition_telegram_test.go` L80–100。

### R-004 · HTTP 200 但 body 非 JSON 仍当成功（部分残留）

- **严重度**：low · recommended · open
- **证据**：`http_sender.go` L141–147。

### R-005 · GET 仍回显密钥末 4 位（仍 open）

- **严重度**：low · recommended · open
- **证据**：`runtime.go` L164–165、L176–184。A-003 声称与代码不符。

### R-006 · YAML 明文与 export 登记表未改（仍 open）

- **严重度**：low · recommended · open
- **证据**：`config.go` telegram 字段；`configpkg.go` `sensitiveFields` 仍只 jwt/initial_password；`.env.example` L176–180 仅为注释。

---

## 必改项汇总

| ID | 严重度 | 建议 | 一句话 |
|----|--------|------|--------|
| **F-001** | med | **required · 仍 open** | 组合根单例端口；webhook 与业务模块同一 dispatcher；未启用走 disabled |
| **F-002** | med | **required · 仍 open** | catalog 迁移 + 加密 + fail-closed + 用 `GetToken` 测重启；或书面收窄 #5 |

A-003 将两条标为 closed **不能成立**。未合法闭合前，本意见**不支持**「开放 required = 0」或对判据 #2/#5 无条件完成。

---

## 与既有意见的异同

| 既有 | 本审 |
|------|------|
| A-002 F-001/F-002 open | **维持 open**。出现了 stub 类型与写库代码，但装配与台账要求未满足。 |
| A-003 self「全部 fixed、required = 0」 | **驳回为证据**。台账与 live 图不一致；持久化测试未测重载。 |
| A-002 R-001/R-003/R-008 | 本审可接受闭合。 |
| A-002 R-002/R-005/R-006 | 仍 open；A-003 标 fixed 与代码不符。 |
| A-002 R-004 | 主路径已修，留非 JSON body 残余。 |
| A-002 R-007 | 同意残余；需 `/govern` 确认用户书面 accepted-residual。 |

---

## 结论 + 建议给编排器/用户的下一步

**verdict: fail。** 整改有真实补丁（disabled 类型、SHA-256 compare、`ok:false`、写库意图），但 A-003 对 F-001/F-002 的闭合是**文档闭合、代码未接通**。helper 与 live `newMux` 双路径尤其不能放行：后续 VP-031 若调用 `ResolveTelegramPorts`，命令会注册到 webhook 看不到的 dispatcher。

建议 `/govern`：

1. 把 A-003 对 F-001/F-002 的 `closed` 改回 **open**，或等本意见后重开整改。
2. F-001：**先接线再谈测试**——`newMux` 只 new 一次，把同一实例注入 webhook 与 `fx.Provide`/`kernel` 槽。
3. F-002：不要再叠 `CREATE TABLE IF NOT EXISTS`；走 catalog + 加密，或用户书面接受「明文旁路表 + API-only」残余。
4. 未闭合前不要把 Root `done` 解释成 A-002 已消化。

---

## 声明

本意见不修改 status/progress；响应由 `/govern` 处理。
