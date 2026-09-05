---
doc_type: goal-execution
id: E-008-a012-response
parent: GOAL-005-r4-evidence-closeout
date: 2026-09-05
status: done
version: 1.0.0
---

# E-008 · A-012 关门撤回与 F-007/F-008 证据补齐

## 事实（时间线）

- 2026-09-05 · **A-012 independent runtime closure re-audit**（codex local · verdict **conditional** · open required **2**）：主方案方向确认（不更换迁移/Manifest/装配方案、不做 Offer Delete）；F-007 缺 0070 多语句 Apply 中途失败零残留 + corrected reopen 的可执行证据；F-008 缺 Telegram-enabled 真实组合根及结构化 Manifest 验收（现有关闭声明把「计划可解析」当成了启用分支装配证据）。
- 2026-09-05 · **关门撤回**：按 P-003（未合法闭合 required 存在时不得保持关门/推导 done），GOAL-005 `done → active`（2/2 → 1/2）、Root GOAL-001 `done → active`（4/4 → 3/4，R4 重开）、VP-031 `closed → active`（v0.3.3，revision history 记录撤回）；goal-tree / workspace 投影同步。
- 2026-09-05 · **D-004 裁决落盘**：F-007/F-008 均走 fixed（不选 residual/overruled）；不改变生产迁移/装配语义。
- 2026-09-05 · **代码实现**：
  - `apps/api/internal/store/migrate_failure_reopen_test.go`：自定义两版本 catalog（v1 bootstrap、v2 多语句 Apply 首条 DDL 后注入错误）；`TestMigrateMidApplyFailureNoResidueThenReopen`（SQLite）与 `TestPostgresMigrateMidApplyFailureNoResidueThenReopen`（真实 PostgreSQL）——失败后 schema 对象 + ledger 行零残留、corrected 同版本/同 checksum catalog 重开恰一条 v2 记录。
  - `apps/api/internal/composition/composition.go`：`NewHTTPSender` 透传 `options.HTTPClient/APIBaseURL`（生产 options 零值，语义不变；使既有 fake seam 覆盖出站发送器）。
  - `apps/api/internal/composition/composition_digitaloffer_telegram_test.go`：`TestDigitalOfferTelegramCompositionRoot`——`biz.digital-offer + channel.telegram` 真实 Fx 组合根 + fake Bot-API seam；`/price` 经真实 webhook HTTP 面、`/entitlements` 与 `/buy` 经同一 Fx 注入 dispatcher 直驱；结构化 Manifest page route/schemaUrl + sidebar navigation 断言；schema/HTTP fail-closed 保留；`DELETE /offers/{id}` 405/404 负向固化。
- 2026-09-05 · **验证**：`go test -count=1 ./internal/store -run 'Test(Migrate|PostgresMigrate)MidApplyFailureNoResidueThenReopen$'` PASS（SQLite + 真实 PG）；`go test -count=1 ./internal/composition -run '^TestDigitalOfferTelegramCompositionRoot$'` PASS；`go build ./...` PASS；`go vet ./internal/store ./internal/composition` PASS；`go test -count=1 ./...` 全量 PASS（apps/api 全包 ok，含 store 77.9s / handler 54.7s / digitaloffer.service 18.7s 的 PG 集成路径）。
- 2026-09-05 · **A-013 响应落盘**：F-007/F-008 全部 closed（fixed ×2，证据路径逐条对应 A-012 必改要求）。

## 产物路径

- `01-decision/D-004-a012-response.md`
- `03-audit/A-013-self-response-a012.md`
- `apps/api/internal/store/migrate_failure_reopen_test.go`
- `apps/api/internal/composition/composition_digitaloffer_telegram_test.go`
- `apps/api/internal/composition/composition.go`（seam 透传 1 行 + 注释）

## 进度评估

- C2 重开进行中：A-012 两项 required 已 closed 并留痕，**未重新关门**——按 A-012 建议先做一次 focused independent finding-closure 复审（仅核对失败/reopen、Telegram-enabled 组合根、结构化 Manifest 三项），`pass · 0 required` 后再关门（GOAL-005 2/2、Root 4/4、VP-031 closed）。
