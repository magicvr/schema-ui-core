---
id: GOAL-001-wallet-prepaid-instrument
doc: audit-entry
record_id: A-006
status: recorded
parent: null
created: 2026-09-02
updated: 2026-09-02
version: 0.1.0
---

# A-006 · A-005 独立交叉审计合并响应与必改项闭合记录

- **source**：self（编排器响应，非 independent）
- **date**：2026-09-02
- **scope**：对 Root `[workspace-029-wallet-prepaid-instrument/GOAL-001-wallet-prepaid-instrument]` A-005（grok-build independent · conditional）全部 findings 的响应、处置与闭合核验
- **verdict**：**pass**（本 scope open required = 0；F-001 已按 `fixed` 路径合法闭合；F-002～F-005 recommended 维持 open、不阻断，处置见下）

## Findings 逐条响应与闭合证据

| Finding | 级别 | 内容 | 闭合路径 | 闭合证据与说明 |
|---------|------|------|----------|----------------|
| **F-001** | required / med | `OwnerExistsFunc` 把 `SubjectExists` OR 进 user 自动开户 HTTP 门禁：已登记 `subject_id` 可过 by-owner 门禁，开出无 `admin.users` 行的 `owner_type=user` 孤儿账本 | **fixed** | 1. `apps/api/internal/composition/composition.go`：`walletOwnerExists` 改为**只查** `authRepository.UserByID`，删除 `SubjectExists` OR 分支；注释更新为 by-owner 面仅开 USER 户、主体存在性留在 `CreateAccount(subject)`/`Redeem`；<br>2. 新增 `apps/api/internal/composition/wallet_owner_gate_subject_test.go`（`TestWalletByOwnerGateRejectsRegisteredSubject`）：已登记 subject id 调 `POST /api/wallet/by-owner/{id}` 与 `.../adjust` → **404 USER_NOT_FOUND**，且 `wallet_accounts` 无 `owner_type=user` 行；阳性对照 `user-admin` → 200 正常开户；<br>3. 阴性对照：还原修复前代码复跑 → **FAIL**（200 + 铸造 `owner_type:user` 账本，精确复现 F-001），修复后 PASS——测试对回归敏感；<br>4. 全包回归 exit 0：`go test ./internal/composition ./internal/handler ./modules/wallet/... -count=1`（composition 23.0s / handler 38.4s / wallet 模块组全部 ok）。 |
| **F-002** | recommended / med | 导出是渲染器对 `items[].code` 的启发式，非协议声明 download | **open（backlog）** | A-005 自判不阻断判据 #5（本 GUI 同手势可导出已实证）；协议级 download outcome 属 UI 合同演进，留待后续 UI/协议波次，不在本必改闭环内处置。 |
| **F-003** | recommended / low | PostgreSQL 上无 Redeem / 并发双花 e2e | **open（backlog）** | A-005 已实测 PG `ON CONFLICT DO NOTHING` 不 abort 同事务；剩余为覆盖缺口（CAS SQL 双方言通用）。本会话无 `PG_TEST_*` 环境，未复跑；维持 recommended open，不阻断。 |
| **F-004** | recommended / low | `batch_id` 无唯一约束 | **open（backlog）** | 非资金不变式（A-005 自判）；留待后续批次管理加固波次。 |
| **F-005** | recommended / low | 过期字段是 Unix 秒 `inputNumber` | **open（backlog）** | 运营 UX（模块层过期拒绝已有测试）；字段已存在且可用，维持 recommended open。 |

## 附带发现（本次验证暴露的既有缺口）

- `TestPublishedManifestNavigationOrder/admin_default_order_matches_the_frozen_list` 在 HEAD（67939f1a）即 FAIL：R3（GOAL-003）注册的 `menu_wallet_vouchers` 未进 admin 冻结导航序列表。已把 `"Prepaid vouchers"` 追加进 `composition_test.go` 的 want 列表（以实际注册序 Order 11 为准，附 VP-029 R3 注释）；composition 全包转绿。该缺口与 F-001 无因果关系，随本响应一并修复留痕（见 E-005）。

## 关闭判定

1. **Required 归零**：A-005 唯一 required（F-001）已按 `fixed` 路径闭合，证据 = 代码 diff + HTTP 回归测试 + 阴性对照 + 全包回归 exit 0。**open required = 0**。
2. **A-005 各判据行状态**：判据 #1 的 conditional（F-001）解除——by-owner 门禁回到只查 `admin.users`，主体校验保持在 `CreateAccount(subject)`/`Redeem`；判据 #7 恢复成立（本条意见范围内 open required = 0）。其余判据行维持 A-005 原判定（#2/#4/#5/#6 pass 或本 GUI 口径 pass；#3 SQLite pass + PG 覆盖缺口 recommended）。
3. **状态**：本响应不改 `status`/`progress`（Root 维持 `done` · 4/4；此前 `done` 与开放 required 并存的不一致随本闭合消除）。F-002～F-005 为 recommended，不构成关门或放行门禁，登记 backlog 处置，不静默标为已修。
4. **可选复审**：建议用户按需调用 `/audit`（grok build · independent）对 A-006 的 F-001 闭合证据做独立复审（对齐 A-002→A-003→A-004 的既有闭环模式）；非强制门禁。

## 声明

本条目为编排器 self 响应，不伪装 `source: independent`；A-005 意见原文与证据保留在 `03-audit/A-005-root-design-code-closeout-independent.md`。
