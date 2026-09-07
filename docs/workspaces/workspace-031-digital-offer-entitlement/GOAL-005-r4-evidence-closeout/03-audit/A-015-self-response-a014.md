---
doc_type: goal-audit
id: A-015-self-response-a014
parent: GOAL-005-r4-evidence-closeout
date: 2026-09-05
status: closed
version: 1.0.0
---

# A-015 · 响应 A-014 F-001（self · response）

## A-015 · 响应 A-014 recommended F-001 并执行关门前置加固（2026-09-05）

- **source**：self（编排器响应记录，非独立审）
- **模式**：response · 响应 A-014（independent · finding-closure · verdict **pass** · open required 0 · 1 条 recommended：F-001）
- **verdict**：pass（F-001 recommended → **fixed**；经用户确认「先处理 F-001 再关门」后执行）

### Finding 闭合

**F-001（recommended · low · 两条组合根测试 Manifest 断言强度不对齐；同一 dispatcher 可加显式断言）→ fixed**

1. **disabled 分支结构化 Manifest 断言**：`apps/api/internal/composition/composition_digitaloffer_test.go` 的 `TestDigitalOfferCompositionRoot` 原 `strings.Contains` 子串断言（仅查两个 page id）升级为 `json.Unmarshal` 结构化解析——断言 `digitaloffer-offers` / `digitaloffer-entitlements` 两 page 的 `route`（`/digitaloffer-offers`、`/digitaloffer-entitlements`）与 `schemaUrl`（`/api/schema/digitaloffer-offers`、`/api/schema/digitaloffer-entitlements`）精确匹配，并断言 sidebar navigation 含两 `pageRef`；与 Telegram-enabled 测试（A-012 F-008 必改要求 ③）断言强度对齐。
2. **dispatcher 指针同一显式断言**：`apps/api/internal/composition/composition_digitaloffer_telegram_test.go` 的 `TestDigitalOfferTelegramCompositionRoot` 增加两层证明——
   - 显式指针相等：`tr.Dispatcher == tr.DispatcherState`（同一 `*Dispatcher` 实例同时注入 digital-offer 装配与 kernel 端口）；
   - 行为探针（observable identity）：启动后向 Fx 注入的 `tr.DispatcherState` 注册 `probe_f008` 命令，再经**真实 webhook HTTP 面**（`POST /api/channel/telegram/webhook` + secret 头）驱动 `/probe_f008`，断言探针 handler 实际执行——webhook 分发路径与直驱使用同一实例。

### 验证（apps/api · Go 1.26.0 · 2026-09-05）

- `go test -count=1 ./internal/composition -run 'TestDigitalOffer(PlanResolution|CompositionRoot|TelegramCompositionRoot)$' -v`：三项 **PASS**。
- `go test -count=1 ./internal/composition`：**PASS**（21.1s，全包）。
- `go vet ./internal/composition`：PASS。

### 状态声明

- A-014 的 `pass · open required 0` 保持成立；F-001（recommended）已按用户选择在关门前置处理完成（fixed）。**无 open required / 无开放必改**。
- 关门执行（GOAL-005 2/2 · Root 4/4 · VP-031 closed v0.3.4）记录见 E-009 与各状态文件；本响应不修改 A-014 原文。
