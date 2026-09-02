---
id: GOAL-001-wallet-prepaid-instrument
doc: audit-entry
record_id: A-007
status: recorded
parent: null
created: 2026-09-02
updated: 2026-09-02
version: 0.1.0
---

# A-007 · A-005 F-001 闭合独立复审（2026-09-02）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型**：finding-closure
- **scope**：`[workspace-029-wallet-prepaid-instrument/GOAL-001-wallet-prepaid-instrument]` A-005 F-001 关闭证据；不以 A-006 台账声明、阴性对照叙述或全包回归口头结果作为闭合依据
- **verdict**：**pass**
- **工作区**：`workspace-029-wallet-prepaid-instrument` · `shared_materials_catalog: none`

### 范围与区间

复审 A-005 F-001 是否已按 `fixed` 路径用代码与可重复测试闭合。A-005 闭合方向原文：by-owner 存在性门禁**只**查 `admin.users`；主体存在性只留在 `CreateAccount(subject)` / `Redeem`；已登记 `subject_id` 调 `POST /api/wallet/by-owner/{id}` → 404 `USER_NOT_FOUND`，且无 `owner_type=user` 行。

本会话实证：

- 通读：`composition.go` `walletOwnerExists`；`handler/wallet.go` by-owner 两写路径；`provider.go` `CreateAccount`；`voucher/service.go` `Redeem`；`wallet_owner_gate_subject_test.go`。
- 复跑（exit 0）：
  - `go test ./internal/composition -run TestWalletByOwnerGateRejectsRegisteredSubject -v` — **PASS 0.07s**
  - `go test ./modules/wallet ./modules/wallet/voucher ./internal/handler -run "TestWalletServiceSubjectAccountLifecycle|TestRedeemSuccess|TestRedeemDuplicate|TestVouchersBatchGeneratePermission"` — PASS
- 未独立复跑 A-006 所述「还原修复前代码必 FAIL」的阴性对照；闭合不依赖该叙述。

### 成果（有证据）

| Finding | 闭合判定 | 代码与测试 |
|---------|----------|------------|
| **A-005 F-001** | **fixed** | `walletOwnerExists` 仅 `authRepository.UserByID`，无 `SubjectExists` OR。by-owner create/adjust 仍在门禁失败时 404 `USER_NOT_FOUND` 并 return，不会落到 `GetOrCreateUserAccount`。`CreateAccount(owner_type=subject)` 与 `Redeem` 仍各自 `SubjectExists`。回归测试走 **生产 composition**（`testMux` + admin profile），对已登记 subject：`POST /by-owner/{id}` 与 `.../adjust` 均为 404 `USER_NOT_FOUND`，`wallet_accounts` 中 `owner_type=user AND owner_id=subject_id` 计数为 0；阳性对照真实 `user-admin` → 200 且 user 行 = 1。未登记主体开户 / 核销成功路径本会话仍绿。 |

### 对照成功标准

A-005 开放 required 仅 F-001，关闭证据可重复核对。判据 #1 的 by-owner 切片：已登记主体不能经 user 自动开户 HTTP 铸造平行 user 账本。本条不重审 A-002 已闭合的三条 required。

### Findings

本条无新 required。

A-005 F-002～F-005（recommended：协议未声明 download、缺 PG Redeem e2e、`batch_id` 非 UNIQUE、过期字段 Unix 秒）本轮**未改闭合状态**。抽查 `wallet-vouchers.json` `generateBatch.onSuccess` 仍为 `reload`，与 F-002 开放一致。不阻断 F-001 闭合。

### 必改项汇总

无。open required = 0。

### 与既有意见的异同

A-006 self 对 F-001 的 `fixed` 主张，经本轮代码与测试核验成立。A-006 的阴性对照与「全包回归 exit 0」本条不采信为证据；采用的是当前源码 + 本会话复跑的 composition HTTP 回归与主体/核销切片。A-005 原文不改写。

### 结论 + 建议给编排器/用户的下一步

A-005 F-001 已合法闭合。可用 `/govern` 登记本条独立复审结果；F-002～F-005 维持 recommended backlog，不构成 required 门禁。

### 声明

本意见不修改 status/progress；响应由 /govern 处理。
