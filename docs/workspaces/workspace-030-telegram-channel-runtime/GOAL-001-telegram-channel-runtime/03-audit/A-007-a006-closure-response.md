---
doc_type: goal-audit
id: A-007-a006-closure-response
parent: GOAL-001-telegram-channel-runtime
date: 2026-09-03
source: self
scope: GOAL-001 A-006 独立复审意见响应（F-001 非 variadic 必选参数 + F-002 主密钥离开源码）
audit_type: finding-closure
verdict: pass
open_required: 0
---

# A-007 · Root GOAL-001 A-006 独立复审意见响应（合并响应）

## 1. 响应背景

编排器（/govern）响应独立复审意见 [A-006-independent-closure-reaudit.md](A-006-independent-closure-reaudit.md)（grok-4.6 · reasoning high · `verdict: fail`，2 required F-001/F-002 维持 open）。依据用户明确指令（2026-09-03）：

> F-001 把 `*TelegramRuntime` 改成非 variadic 必选参数并删 fallback；F-002 主密钥离开源码。不要再用「测试里传入了 tr」证明 Fx 已注入。

按 **fixed** 路径整改，不以 self 声称代替可核对代码证据。

## 2. 必改项彻底闭合台账（Required Findings）

| ID | 严重度 | 闭合路径 | 闭合事实与核验代码 | 状态 |
|----|--------|----------|-------------------|------|
| **F-001** | med / required | **fixed** | 1. **非 variadic 必选参数**：`newMux` 与 `newMuxWithExtraProviders` 最后一参由 `trs ...*TelegramRuntime` 改为 `tr *TelegramRuntime`（`apps/api/internal/composition/composition.go`），删除 `len(trs)` 取值逻辑与 `tr = newTelegramRuntime(...)` fallback 分支。<br>2. **单实例来源**：Fx 图仅经 `fx.Provide(newTelegramRuntime)` 构造一次 `*TelegramRuntime`；`newMux` 直接消费注入的 `tr`，webhook/settings 用同一实例。<br>3. **经 NewApp/fx 的同一实例测试**：新增 `TestTelegramFxInjection_SameRuntime`（`composition_telegram_test.go`）——用 `fx.Populate` 从真实 `newAppWithOptions`（= NewApp 同一 fx 图 + 测试探针）取出注入的 `*TelegramRuntime` 与 `*http.ServeMux`，在注入 dispatcher 上 `RegisterCommand("status")`，经注入 mux 打真实 webhook，断言命令被分发。**不再手工把 tr 传进 newMux**。<br>4. `ResolveTelegramPorts` 收敛为独立 helper（仅测试/外部无 Fx 消费者使用），返回 error 传播（不再静默吞错）。 | **closed** |
| **F-002** | med / required | **fixed** | 1. **主密钥离开源码**：删除 `runtime.go` 中 `defaultMasterKey = sha256(...)` 编译期常量与 `crypto/sha256` import；`NewRuntimeManager(..., masterKey []byte, ...)` 改为必选参数，空 key 直接构造错误（fail-closed）。<br>2. **组合根解析密钥**：`newTelegramRuntime` 复用 `mail.LoadOrCreateMasterKey`——`cfg.TelegramMasterKey`（`TELEGRAM_MASTER_KEY` env）或 `cfg.TelegramMasterKeyPath`（`TELEGRAM_MASTER_KEY_PATH` env / `telegram.master_key_path` yaml）指向的密钥文件（缺省 `filepath.Join(filepath.Dir(DBPath), "telegram-master.key")`）。<br>3. **initPersistence fail-closed**：`initPersistence` 返回 error，`NewRuntimeManager` 透传；读库/解密失败不再 `_ =` 吞错，组合根启动失败。<br>4. **配置面**：`config.go` 新增 `TelegramMasterKey` / `TelegramMasterKeyPath` 字段与 env/yaml 接线；`.env.example` 登记 `TELEGRAM_MASTER_KEY` / `TELEGRAM_MASTER_KEY_PATH`（`TestCanonicalEnvExample` 覆盖）。<br>5. **测试适配**：channel 包测试统一走 `testMasterKey()`/`newTestRuntimeManager`，全量测试绿（见 §4）。 | **closed** |

## 3. Recommended 项状态核销

- **R-004（HTTP 200 非 JSON 仍当成功）**：`http_sender.go` L141–147 —— **维持 open**（recommended，非必改；A-006 未升格）。本响应不关闭。
- **R-007（Allow/Record TOCTOU 残余）**：A-006 同意残余、标注"需 /govern 书面确认"。本响应随 F-001/F-002 整改不重新开启；残余范围仍为 VP-027 端口形状，留待用户书面确认（不阻塞 required 门禁）。

## 4. 验证证据

- `go build ./...`：通过。
- `go test ./internal/channel/telegram/... ./modules/channel/telegram/...`：ok。
- `go test ./internal/composition/...`：ok（含 `TestTelegramFxInjection_SameRuntime`、`TestTelegramRuntime_PersistenceAcrossRestart`、`TestTelegramChannelComposition_RealWebhookMount`、`TestResolveTelegramPorts_EnabledAndDisabled`）。
- `go test ./internal/config/...`：ok（`TestCanonicalEnvExample` 覆盖新增 env）。
- `go test ./...`（apps/api）：全部 ok，exit 0。

## 5. 结论与关门事实

- A-006 指出的 F-001 与 F-002 必改项已通过真实代码修改闭合，证据可核对：
  - F-001：签名不再 variadic、fallback 删除、Fx 同一实例测试经 `fx.Populate` 走真实组合根。
  - F-002：主密钥不再出现在源码；`initPersistence` 失败在组合根 fail-closed。
- 开放 required findings：**0**。
- R-004（recommended）保持 open，不阻塞关门；R-007 残余留待用户书面确认。
- Root 目标 `GOAL-001-telegram-channel-runtime` 维持已关门状态（required 门禁重新闭合）。
