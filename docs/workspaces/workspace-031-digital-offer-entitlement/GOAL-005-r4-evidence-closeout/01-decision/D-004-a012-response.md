---
doc_type: goal-decision
id: D-004-a012-response
parent: GOAL-005-r4-evidence-closeout
date: 2026-09-05
status: accepted
version: 1.0.0
---

# D-004 · 响应 A-012：F-007/F-008 证据缺口以 fixed 路径闭合

## 背景

A-012（independent · codex · runtime closure re-audit · verdict **conditional** · open required **2**，均 medium）在 A-006 整改后对 workspace-031 Root/GOAL-005 运行时关门证据做了 code-first 复审。审计确认主方案方向正确（**不建议**更换 profile-gated migration、新增 public assembly API 或实现 Offer Delete），但指出两项关门证据仍把「实现推断」当成「失败路径/启用分支事实」：

- **F-007**（required · medium）：0070 多语句 Apply 中途失败后「无 schema/ledger 残留 + corrected catalog reopen 成功」没有可执行切片——事务实现成立（`migrate.go` Apply 与 ledger 同事务）但缺 runner 级失败/恢复测试。
- **F-008**（required · medium）：`channel.telegram` 启用的真实组合根未被装配——现有组合根测试明确排除 Telegram；Manifest 断言仍是子串搜索，未对 route/schemaUrl/navigation refs 做结构化断言。

A-012 明确：是否以书面 `accepted-residual` 限定范围后先立项后继 VP，由 `/govern` 提交用户裁决；本独立意见不替用户接受残余。

## 裁决

**采用 fixed 路径**（不选 accepted-residual / user-overruled）：

### F-007 → fixed（补 runner 级失败/reopen 测试，SQLite + PG）

在 `apps/api/internal/store` 新增迁移 runner 级测试（不新增生产代码）：

1. 自定义两版本 catalog：v1 bootstrap（建 ledger + 空基表）、v2 多语句 Apply（`CREATE TABLE` + `CREATE INDEX` 后注入错误）。
2. 首开失败 → 断言该迁移的 schema 对象（表 + 索引）与 ledger 行**均不存在**、v1 已提交（零残留）。
3. 用修正后的**同版本/同 checksum** catalog 重开 → 断言迁移成功、v2 ledger 恰好 1 行。
4. SQLite 必测；本仓真实 PG 环境可用 → PG 同语义进集成测试（`TestPostgresMigrateMidApplyFailureNoResidueThenReopen`）。

不新增 Down migration，不改 profile-gated，不改 compiled-global 语义（A-012 明确无需）。

### F-008 → fixed（Telegram-enabled 真实组合根 + 结构化 Manifest 断言）

1. **语义保持的 seam 透传**：`composition.buildTelegramRuntime` 中 `NewHTTPSender(rt, nil, "")` 改为 `NewHTTPSender(rt, options.HTTPClient, options.APIBaseURL)`——生产路径 `options` 恒为零值，行为不变；仅让既有 `telegramRuntimeOptions` 测试 seam（lifecycle 测试已在用）覆盖出站发送器，从而 Telegram-enabled 组合根测试可用 fake client 捕获回包、**零真实外网 Bot API 依赖**。
2. 新增 `internal/composition/composition_digitaloffer_telegram_test.go`：
   - `biz.digital-offer + channel.telegram` 真实 `newAppWithOptions`/Fx 组合测试，用 fake HTTP seam（`getMe`/`deleteWebhook`/`sendMessage`）。
   - 断言 Fx 注入的 `*TelegramRuntime` 为**同一实例**（webhook 挂载的 dispatcher），`HasBusinessHandlers()` 为真（price/buy/entitlements 已注册）。
   - 经**真实 webhook HTTP 面**驱动 `/price`（secret 校验 → bot identity → subject 映射 → inbound 持久化 → dispatcher → 出站回包），断言回包经 fake client 捕获。
   - dispatcher 直驱 `/entitlements`、`/buy`，断言回包（三条冻结命令均可达）。
   - Manifest 响应**解析为结构**，断言两 page 的 route/schemaUrl 与 sidebar navigation pageRef；保留 schema/HTTP fail-closed 断言（401/200/404）。
   - 负向：`DELETE /api/digitaloffer/offers/{id}` 405/404（固化「无物理删除」边界，A-012 非 required 建议 1 一并落实）。

### 未采用的替代

- **accepted-residual（限定后继 VP 范围）**：未采用。F-007/F-008 是**本目标**运行时关门证据的组成部分（非仅后继 VP 前置）；修复边界明确（测试 + 一条 seam 透传），不改变生产语义，无需把残余风险留给后继。
- **user-overruled**：不适用；审计事实（缺失败路径/启用分支证据）成立。
- **更换迁移/装配方案**：A-012 明确不建议，且会扩大核心语义变更面，需 cross audit——不选。

## 影响

- Root / GOAL-005 / VP-031 关门状态**撤回**（GOAL-005 → active、Root → active、VP-031 → active），待两项 required 闭合后经 **focused independent closure 复审**（A-012 推荐的唯一复审）重新关门。
- 无合同语义变更、无迁移/Manifest/装配方案变更；F-008 的 seam 透传在生产路径零值不变。
- A-012 非 required 建议 2/3（权限分离测试、descriptor 等价专门测试）不阻塞关门，标记为后续可选；建议 1（DELETE 负向测试）已随 F-008 一并落实。
