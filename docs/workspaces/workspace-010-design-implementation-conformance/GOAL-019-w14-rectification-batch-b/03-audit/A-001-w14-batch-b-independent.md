---
id: GOAL-019-w14-rectification-batch-b
doc: audit
status: active
parent: GOAL-015-w14-user-perspective-review
created: 2026-08-17
updated: 2026-08-17
version: 0.1.0
---

# A-001 · W14 整改批 B 独立审计（F-05～F-07）

- source: independent
- auditor: independent subagent
- date: 2026-08-17
- scope: GOAL-019 S1-S3（F-05 列表端点校验/分页 · F-06 错误目录契约 · F-07 搜索/排序/过滤一致性）
- verdict: **fail**

## 成果（有证据）

- F-05 分页校验：
  - recycle-bin 列表 page/pageSize 校验与上限：`apps/api/internal/handler/recyclebin.go:53-62`，sort/order 校验 `:63-78`，repository 分页生效 `apps/api/internal/modules/recyclebin/store/repository.go:105-121`。
  - wallet accounts/entries/reconcile-runs 分页校验与生效：`apps/api/internal/handler/wallet.go:57-66, 284-303, 319-328`；entries entryType 过滤传递 `:329-330`。
  - per-task runs 分页参数校验并生效：`apps/api/internal/handler/scheduledtasks.go:300-319`。
  - data-permission policies 提供资源分页包裹（内存切片分页）：`apps/api/internal/handler/datapermission.go:65-95`。
- F-06 错误目录/契约：
  - `OPERATION_NOT_FOUND` 已入目录并带 messageKey：`apps/api/internal/errorcatalog/errorcatalog.go:156`。
  - 契约测试纳入 OPERATION_NOT_FOUND 且要求双语+messageKey：`apps/api/internal/handler/error_contract_test.go:82-84, 147-162`。
  - i18n key 存在：`apps/web/src/i18n/messages/en-US.json:705`、`zh-CN.json:705`。
- F-07 搜索/排序/过滤：
  - 通知 q 大小写不敏感：`apps/api/internal/modules/authsession/notifications_repository.go:141-144`。
  - recycle-bin 暴露 sort/order 并校验：`recyclebin.go:63-78`，repository 排序映射 `repository.go:98-104`。
  - wallet 账户搜索扩展（owner_id/owner_type/currency）：`apps/api/internal/modules/wallet/store/repository.go:148-151`。
  - wallet ledger entry-type 过滤：`wallet.go:329-330`、`repository.go:510-513`；schema 选项 `wallet-entries.json:34-55`；i18n keys `en-US.json:654-660`、`zh-CN.json:654-660`。

## Findings

### F-001 · F-06 错误码细分未完成（required · 阻断）

- 严重度：required / blocking
- 证据：
  - `apps/api/internal/handler/datapermission.go:114-127` 仍用 `INVALID_SCOPE` 表达 policies 端点请求体 JSON/缺失字段；而 `apps/api/internal/errorcatalog/errorcatalog.go:137` 中 `INVALID_SCOPE` 的目录文案是 “scope must be all or self（数据范围必须是 all 或 self）”，两类语义仍复用同一目录码。
  - `apps/api/internal/handler/wallet.go:315-317` 在 GET entries 缺少 `accountId` 时仍用 `INVALID_WALLET_BODY`；`errorcatalog.go:142` 中该码文案为 “body must be JSON with ownerType and ownerId”，对 query 参数缺失是误导性消息。
- 影响：F-06「错误码复用误导消息（INVALID_SCOPE / INVALID_WALLET_BODY 细分）」未完整落地；目录码与真实错误语义不匹配，违背错误契约目标。
- blocking：是。需在 S1 决策冻结具体细分码/映射并完成 handler 与目录、契约测试、i18n 对齐后才可推进关门。

### F-002 · F-07 wallet ledger 搜索 q 未实现（required · 阻断）

- 严重度：required / blocking
- 证据：
  - `apps/api/internal/modules/wallet/schema/wallet-entries.json:25-29` 暴露 `q` 搜索框，`schema.walletEntries.search.q`（en-US “Memo / reference”）两语种 key 均存在（`en-US.json:654`、`zh-CN.json:654`）。
  - 但 `apps/api/internal/handler/wallet.go:311-344` 的 `walletListEntries` 未读取/传递 `q`；`WalletService.ListEntries` 签名无 q；`apps/api/internal/modules/wallet/store/repository.go:504-548` 的 `ListEntries` 只应用 entryType，q 为空。
- 影响：schema 宣称的 ledger 搜索（备注/关联单）在端点与存储层完全未实现，F-07「wallet 搜索扩展」对 ledger 缺口。
- blocking：是。

### F-003 · S1-S3 实施未在目标台账记录（required · 阻断）

- 严重度：required / blocking（过程/证据门禁）
- 证据：
  - `00-meta.md` 仍为 `progress: 0/4`、S1～S4 全未勾选；`01-decision.md` 决策索引为空，`I-001/I-002` 仍 open；`02-execution.md` 仅 `E-001`（立项），无 S1/S2/S3 实施记录。
  - 同时代码已包含大量 W14 批 B 改动（如 `wallet.go:1-4` 注释引用 GOAL-019 D-002、`errorcatalog.go:155-156`、`notifications_repository.go:27-33` 等）。
- 影响：审计范围 S1-S3 的实际实施没有决策冻结记录、执行事实记录与证据回指，违反 `03-audit`/五件套台账纪律，无法可靠关门。
- blocking：是（对推进/关门）。

### F-004 · data-permission policies 分页为内存分页，真实分页未定（recommended · 非阻断）

- 严重度：recommended / non-blocking
- 证据：`datapermission.go:75-84` 先取全量 `ListPolicies` 再做切片分页，未做 SQL LIMIT/OFFSET 下推。
- 影响：若「真分页」要求存储层 pushdown，则尚未满足；若契约仅要求端点分页参数生效，则当前实现可接受。
- 建议：S1 决策记录 I-001 结论（真分页或明确不分页），并在 03-audit/01-decision 留痕。

## 必改项汇总

1. F-001：按 F-06 拆分 INVALID_SCOPE / INVALID_WALLET_BODY 的复用语义（含 GET entries 缺 accountId、policies 请求体错误），更新 handler、errorcatalog、契约测试与 i18n。
2. F-002：实现 wallet ledger `q` 搜索（备注/关联单），贯通 schema → handler → repository，并补测试。
3. F-003：回填 S1-S3 决策/执行/审计台账与 goal-tree 状态（本意见不改状态，交由编排器响应）。

## 结论

S1-S3 中大部分 F-05 分页校验与 F-07 的 entryType/排序/通知 q 已落地并有代码证据；F-06 的 OPERATION_NOT_FOUND 与契约测试、i18n 亦已补齐。但存在 2 个实质性代码缺口（F-001 错误码细分未完成、F-002 wallet ledger q 未实现）及 1 个台账证据缺口（F-003），故整体判定 **fail**；修复并回填台账后方可推进关门。

## 声明

This opinion does not modify status/progress; it records an independent audit finding for the orchestrator to respond (P-003). It does not touch goal-tree, status, or progress.