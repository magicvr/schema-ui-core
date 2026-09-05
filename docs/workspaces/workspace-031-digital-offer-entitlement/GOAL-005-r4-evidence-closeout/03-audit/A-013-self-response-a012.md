---
doc_type: goal-audit
id: A-013-self-response-a012
parent: GOAL-005-r4-evidence-closeout
date: 2026-09-05
status: closed
version: 1.0.0
---

# A-013 · 响应 A-012（self · response）

## A-013 · 响应 A-012 focused runtime closure re-audit（2026-09-05）

- **source**：self（编排器响应记录，非独立审）
- **模式**：response · 响应 A-012（independent · codex · verdict **conditional** · open required **2**：F-007/F-008，均 medium）
- **verdict**：conditional → 两项 required 全部 closed（fixed ×2）

### Finding 闭合

**F-007（迁移 Apply 中途失败零残留 + reopen）→ fixed**

- 证据：`apps/api/internal/store/migrate_failure_reopen_test.go`
  - `TestMigrateMidApplyFailureNoResidueThenReopen`（SQLite）：多语句 v2 首条 DDL 后注入错误 → 断言 `f007_fail_items` 表与索引、v2 ledger 行均不存在、ledger 恰为 `{1}`；修正后的同版本/同 checksum catalog 重开 → 迁移成功、v2 恰 1 行、`verifyIntegrity` 通过。**PASS**。
  - `TestPostgresMigrateMidApplyFailureNoResidueThenReopen`（真实 PostgreSQL，`pgtest` 环境）：同语义断言（`to_regclass` 为空、ledger=1 → 重开后 =2、`string_agg` 版本 = `1,2`）。**PASS**。
- 与 A-012 必改要求逐条对应：①首条创建对象后返回注入错误、断言 schema 对象与 ledger 行均不存在；②修正后同版本/同目录重开、断言成功且仅一条 ledger 记录；③SQLite 必测 + PG 同语义进集成测试。未加 Down、未改 profile-gated。

**F-008（Telegram-enabled 真实组合根 + 结构化 Manifest）→ fixed**

- 证据：`apps/api/internal/composition/composition_digitaloffer_telegram_test.go` + `composition.go` seam 透传（`NewHTTPSender(rt, options.HTTPClient, options.APIBaseURL)`；生产 options 零值，行为不变）。
  - `TestDigitalOfferTelegramCompositionRoot`（真实 `newAppWithOptions`/Fx + fake Bot-API seam，零外网依赖）：**PASS**。
  - 必改要求 ①：`biz.digital-offer + channel.telegram` 真实组合测试成立，fake/runtime seam 即既有 `telegramRuntimeOptions`。
  - 必改要求 ②：`/price` 经真实 webhook HTTP 面（`POST /api/channel/telegram/webhook` + secret 头）可达，`/entitlements`、`/buy` 经同一 Fx 注入 dispatcher 直驱可达——三命令均注册在同一 dispatcher（webhook 挂载实例），回包均经 fake client 捕获断言。
  - 必改要求 ③：Manifest 解析为结构，断言两 page 的 route/schemaUrl 与 sidebar navigation pageRef；schema/HTTP fail-closed 断言（401/200/404）保留。
  - 负向固化：`DELETE /api/digitaloffer/offers/{id}` → 405/404（无物理删除，A-012 非 required 建议 1 一并落实）。
- A-012 非 required 建议 2（`digitaloffer.offer.manage` / `digitaloffer.entitlement.void` 权限分离测试）、建议 3（registry/Provider descriptor 等价专门测试）为 recommended 非 required，不阻塞关门；当前真实 app 启动已含 descriptor 一致性校验，留作后续可选加固。

### 状态声明

- GOAL-005 / Root / VP-031 关门状态已撤回（见 E-008 与各自状态文件），**未**重新关门：A-012 建议完成 F-007/F-008 后做一次 **focused independent finding-closure 复审**再复用于后继 VP / 重新关门。复审只需核对：①SQLite/PG Apply 中途失败零残留 + corrected reopen；②Telegram-enabled 真实组合根与同一 runtime 命令注册；③结构化 Manifest page/route/schemaUrl/navigation 断言。
- 本响应记录不修改 A-012 原文；F-007/F-008 的 closed 以本条登记 + E-008 事实为依据，独立复审通过前不推导 `done`。
