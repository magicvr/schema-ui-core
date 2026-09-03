---
doc_type: goal-execution
id: E-009-a006-response
parent: GOAL-001-telegram-channel-runtime
date: 2026-09-03
status: recorded
version: 1.0.0
---

# E-009 · A-006 复审整改（F-001 非 variadic / F-002 主密钥离开源码）

## 背景

grok independent [A-006](../03-audit/A-006-independent-closure-reaudit.md) 复审 fail：F-001（`*TelegramRuntime` variadic 被 dig 丢弃，生产仍双 runtime）、F-002（主密钥 `sha256(...)` 编译期常量写死源码、`initPersistence` 吞错）维持 open。用户明确指令（2026-09-03）：F-001 改非 variadic 必选参数并删 fallback；F-002 主密钥离开源码；不再用「测试里传入了 tr」证明 Fx 已注入。

## 事实（代码改动）

1. **F-001**（`apps/api/internal/composition/composition.go`）：
   - `newMux` 与 `newMuxWithExtraProviders` 最后一参 `trs ...*TelegramRuntime` → `tr *TelegramRuntime`（非 variadic 必选）。
   - 删除空 `trs` 时的 fallback `tr = newTelegramRuntime(plan, cfg, st, rateLimiters)`。
   - 新增 `newAppWithOptions`（= NewApp 同一 Fx 图 + 测试探针选项）作为 composition-root 测试 seam。
   - `ResolveTelegramPorts` 返回 error 传播。
2. **F-002**（`apps/api/internal/channel/telegram/runtime.go` + `config.go`）：
   - 删除 `defaultMasterKey`（源码常量）与 `crypto/sha256` import；`NewRuntimeManager` 增必选 `masterKey []byte`，返回 `(*RuntimeManager, error)`，空 key fail-closed。
   - `initPersistence` 返回 error（读库/解密失败不再吞错），组合根透传。
   - `newTelegramRuntime` 经 `mail.LoadOrCreateMasterKey(cfg.TelegramMasterKey, path)` 解析密钥；新增 `TelegramMasterKey`/`TelegramMasterKeyPath` 配置字段与 env（`TELEGRAM_MASTER_KEY`/`TELEGRAM_MASTER_KEY_PATH`）/yaml（`telegram.master_key_path`）接线；`.env.example` 登记。
3. **测试**：
   - 新增 `TestTelegramFxInjection_SameRuntime`：经 `newAppWithOptions` + `fx.Populate` 从真实 Fx 图取出注入的 `*TelegramRuntime` 与 `*http.ServeMux`，在注入 dispatcher 注册命令、经注入 mux 打真实 webhook 断言分发——证明同一实例，不手工传 tr。
   - `channel/telegram` 包测试统一走 `testMasterKey()`/`newTestRuntimeManager`；composition/s2/metrics 等调用点适配非 variadic 签名（无 channel.telegram 的 plan 传 nil）。

## 验证

- `go build ./...` 通过。
- `go test ./internal/channel/telegram/... ./modules/channel/telegram/...` ok。
- `go test ./internal/composition/...` ok（含 `TestTelegramFxInjection_SameRuntime`、`TestTelegramRuntime_PersistenceAcrossRestart`、`TestTelegramChannelComposition_RealWebhookMount`）。
- `go test ./internal/config/...` ok（`TestCanonicalEnvExample` 覆盖新增 env）。
- `go test ./...`（apps/api）全部 ok。

## 评估

A-006 两条 required 已按用户指令 fixed 闭合；开放 required = 0。R-004（recommended）维持 open；R-007 残余待用户书面确认。详见 [A-007](../03-audit/A-007-a006-closure-response.md)。
