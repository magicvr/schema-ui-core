---
doc_type: goal-execution
id: E-010-a008-response-and-r5
parent: GOAL-001-telegram-channel-runtime
date: 2026-09-03
status: recorded
version: 1.0.0
---

# E-010 · A-008 响应 + 遗留 recommended 处置 + R5 判据 #5 补做

## 背景

用户指令（2026-09-03）：响应 GOAL-001 A-008，并处理所有遗留 recommended 项目。

## 事实

1. **A-008 响应（A-009）**：A-008（independent）`pass`，F-001/F-002 closed（required=0）。A-009 self 接受并全量处置遗留。
2. **R-004（HTTP 200 非 JSON 仍当成功）→ fixed**：`http_sender.go` 非 JSON 200 body 返回 error（fail-closed）；新增 `TestHTTPSender_Status200_NonJSONBodyFailsClosed`。
3. **informational（PATCH 空串 vs seed）→ fixed**：`runtime.go initPersistence` 行存在即权威（含空值无条件应用）；新增 `TestTelegramRuntime_ClearSurvivesRestart`。
4. **R-009（默认密钥文件与 DB 同目录）→ accepted-residual（用户书面）**：与 mail W13 F-017 对齐；`TELEGRAM_MASTER_KEY`/`TELEGRAM_MASTER_KEY_PATH` 分离路径已在 `.env.example` 与代码注释强化。
5. **R-007（Allow/Record TOCTOU）→ 用户裁决新建 VP**：登记 **VP-032-rate-limiter-atomic-port**（planned v0.1.0 · VRev-071 self `pass` · roadmap 已落盘）。
6. **判据 #5 Schema/Nav/tab → 用户裁决补做 Admin UI tab**：新建 GOAL-006 并交付（C1 后端 Schema/Page/Nav/Manifest + C2 前端 telegram-admin-tab/i18n + C3 self A-001 `pass`），判据 #5 恢复完整交付口径。

## 验证

- `go test ./...`（apps/api）全部 ok。
- 前端：`telegram-admin-tab.test.tsx` 2/2、`schema-keys.structural.test.ts` 4/4 ok；改动文件 `tsc` 无错误。

## 评估

A-008（pass）已响应；所有遗留 recommended/residual/informational 项已按用户裁决处置完毕；R5 判据 #5 补做交付完成。Root 回归 done。
