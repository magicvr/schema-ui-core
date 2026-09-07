---
doc_type: goal-audit
id: A-009-a008-response
parent: GOAL-001-telegram-channel-runtime
date: 2026-09-03
source: self
scope: GOAL-001 A-008 独立复审意见响应 + 遗留 recommended 全量处置
audit_type: finding-closure
verdict: pass
open_required: 0
---

# A-009 · Root GOAL-001 A-008 独立复审意见响应 + 遗留 recommended 处置

## 1. 响应背景

编排器（/govern）响应独立复审意见 [A-008-independent-closure-reaudit.md](A-008-independent-closure-reaudit.md)（grok-4.6 · reasoning high · `verdict: pass` · 开放 required = 0：A-006 F-001/F-002 已按源码与 Fx 生产图闭合）。用户指令（2026-09-03）：响应 A-008，并**处理所有遗留 recommended 项目**。

## 2. A-008 接受（independent pass）

A-008 独立核对后确认：

- **F-001 fixed**：`newMux`/`newMuxWithExtraProviders` 非 variadic 必选 `tr *TelegramRuntime`；无 fallback 第二工厂；`TestTelegramFxInjection_SameRuntime` 经 `fx.Populate` 走真实 `NewApp` 图证明同一实例。
- **F-002 fixed**：grep 无 `defaultMasterKey`；`LoadOrCreateMasterKey`（env/文件）解析主密钥；`initPersistence` 失败冒泡 fail-closed。
- 开放 required = **0**。

本响应接受该 pass，两条 required 维持 closed。

## 3. 遗留 recommended / residual / informational 处置台账

| ID | 级别 | 处置 | 用户裁决 / 证据 |
|----|------|------|-----------------|
| **R-004**（HTTP 200 非 JSON 仍当成功） | recommended | **fixed** | 代码：`http_sender.go` — `json.Unmarshal` 失败即返回 error（fail-closed），不再静默当成功。测试：`TestHTTPSender_Status200_NonJSONBodyFailsClosed`。 |
| **informational**（PATCH 空串 vs seed） | informational | **fixed** | 代码：`runtime.go initPersistence` — 行存在即权威，解密值（含空）无条件应用；清空 token/secret 跨重启保持空，不回 env seed。测试：`TestTelegramRuntime_ClearSurvivesRestart`。 |
| **R-009**（默认 master key 文件与 DB 同目录） | recommended | **accepted-residual**（用户书面） | 用户裁决（2026-09-03）：与 mail W13 F-017 完全对齐——默认同目录零配置；生产用 `TELEGRAM_MASTER_KEY` 或 `TELEGRAM_MASTER_KEY_PATH` 分离；代码注释 + `.env.example` 已强化注明。复审触发：KMS/HSM（RT-K02）落地或密钥管理波次。 |
| **R-007**（Allow/Record TOCTOU） | residual | **新建 VP 下一波做端口原子化**（用户书面） | 用户裁决（2026-09-03）：新建 VP（planned）承接 `kernel.RateLimiter` 原子化（`AllowRecord`）并迁移全部消费者。本 VP 已登记 = **VP-032**（见 roadmap；`/vision` 正式冻结退出分母后再交 `/govern` 开区）。复审触发：VP-032 激活。 |
| **Schema/Nav/tab（判据 #5）** | recommended（A-006 降级） | **补做 Admin UI tab**（用户书面） | 用户裁决（2026-09-03）：不做 API-only 收窄，**补做 Admin UI tab**。已立项 = workspace-030 新子目标（GOAL-006），后端 Schema/Page/Nav + 前端 custom component 落地。判据 #5 恢复为完整交付口径。 |

## 4. 验证证据

- `go test ./internal/channel/telegram/... ./internal/composition/... ./internal/config/...`：全部 ok（含新增 `TestHTTPSender_Status200_NonJSONBodyFailsClosed`、`TestTelegramRuntime_ClearSurvivesRestart`）。
- `go build ./...`：通过。

## 5. 结论

- A-008 independent pass 接受；开放 required = **0**。
- 遗留 recommended/residual/informational 全量处置：R-004 / informational **fixed**；R-009 **accepted-residual**（用户书面）；R-007 **登记 VP-032**（用户书面裁决新建 VP）；判据 #5 **立项补做 UI tab**（用户书面）。
- Root `GOAL-001-telegram-channel-runtime` 维持 done；判据 #5 补做由新子目标承接（见 goal-tree）。
