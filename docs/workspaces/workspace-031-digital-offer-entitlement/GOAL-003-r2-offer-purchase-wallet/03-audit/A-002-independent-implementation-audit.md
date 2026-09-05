---
doc_type: goal-audit
id: A-002-independent-implementation-audit
status: closed
parent: GOAL-003-r2-offer-purchase-wallet
created: 2026-09-05
updated: 2026-09-05
date: 2026-09-05
source: independent
auditor: codex (gpt-5.6-sol, medium)
audit_type: execution-facts
scope: R2 实施审计，资金路径门禁；对照 D-002 v1.1.0 与 VP-031 判据 1/2
verdict: fail
open_required: 2
version: 1.0.0
---

# A-002 · R2 实施独立审计（execution-facts）

- **source**：independent
- **auditor**：codex (gpt-5.6-sol, medium)
- **类型 / scope**：execution-facts · R2 资金路径门禁；审查 `apps/api/modules/digitaloffer/`、`apps/api/internal/handler/digitaloffer.go`、composition/钱包 subject 增量/error-contract/v70 钉死点，对照 GOAL-002 D-002 v1.1.0（§1/§2/§4/§7/§8/§9/§10）与 VP-031 判据 1/2
- **verdict**：**fail**
- **open required**：**2**

## 范围与区间

- 工作区绑定核对：`workspace-031-digital-offer-entitlement` 的 `root_goal`、`canonical_scope`、`primary_plan` 与 VP-031 lead workspace 一致；本轮未读取或混入其他工作区状态。
- 实施区间：migration 0070、digitaloffer store/service/provider/manifest/schema、handler、composition、wallet subject 增量、errorcatalog/error-contract、compiled persistence/store identity，以及指定 purchase/handler 验收测试。
- 合同分母：D-002 标题与修订说明为 v1.1.0；本轮按 v1.1.0 审计。`00-meta.md`、workspace/goal-tree 中仍显示 v1.0.0，作为投影问题记录，不改变本次代码符合性分母。
- 当前工作树不是提交快照：审计前 `git status --short` 显示 `apps/api/modules/digitaloffer/store/store.go` 已修改、A-001 为未跟踪文件；本意见评价 2026-09-05 当前工作树事实，不把 commit `0e65f392` 当作完整证据分母。

## 成果（有证据）

- migration 0070 在 SQLite/PostgreSQL 两套 DDL 中均建立 `digital_offers`、`digital_purchases`、`digital_entitlements`，落实价格正数、权益形态互斥、购买 `UNIQUE(subject_id, request_id)`、`fulfilled` 单终态及权益形态 CHECK；v70 同时进入 compiled persistence 与 lost-ledger identity 指纹（`apps/api/modules/digitaloffer/migration/migration.go:16-164`、`apps/api/modules/compiled/persistence.go:26-52`、`apps/api/internal/store/identity.go:90-120`）。
- Offer 更新使用 `WHERE id=? AND version=?` 并递增 version，0 行更新区分 not-found/version-conflict；购买唯一冲突经 `kernel.IsUniqueViolation` 分类为本 attempt 失败（`apps/api/modules/digitaloffer/store/store.go:309-390`）。
- 购买重试循环位于事务 runner 外且最多 3 次；每次 attempt 从 `(subject_id, request_id)` 回读开始，随后校验 offer/subject/account/currency，再按 freeze → deduct_frozen → purchase → entitlement 顺序执行（`apps/api/modules/digitaloffer/service/service.go:249-417`）。两笔 wallet mutation 使用正的 `offer.PriceAmount`、互异 `:freeze`/`:deduct` 幂等键，并共享 `biz_offer_purchase` / `purchaseID` 反链（同文件 `359-392`）。
- Admin 写路径将域写与 `RecordOperationTx` 放在同一 caller-owned transaction；recorder 缺失/失败返回 error。重复 void 在检测到已作废后直接成功返回，不追加审计（`apps/api/modules/digitaloffer/service/service.go:120-211,418-487`）。
- Admin 路由权限键与 provider 权限贡献一致；公开目录无认证，按 `bizoffer|list|<ip>` 每请求 `AllowRecord`，拒绝返回 429 + `Retry-After`，审计范围内未发现 key-wide `Clear`（`apps/api/internal/handler/digitaloffer.go:35-267`、`apps/api/modules/digitaloffer/provider.go:65-127`）。
- `biz.digital-offer` 仅在 `plan.HasModule("biz.digital-offer")` 时装配，不在默认 Profile；购买与权益模型只关联 `subject_id`，未创建 `admin.users` 关系（`apps/api/internal/composition/composition.go:619-632`、`apps/api/kernel/profile.go:25-113`、`apps/api/modules/digitaloffer/store/store.go:56-99`）。
- 本轮实际执行：在 `apps/api/`、临时 `GOCACHE` 下，`go build ./...` 成功；`go test ./modules/digitaloffer/... ./internal/handler/` 成功。测试输出包含 digitaloffer service `ok` 与 internal/handler `ok`，真 PostgreSQL `TestPurchasePostgresAcceptance` 在本轮目标测试中执行通过。

## 对照成功标准

| 对照项 | 结论 | 证据 / 说明 |
|--------|------|-------------|
| D-002 §1/§2/§10：模型、DDL、迁移、Profile、subject-only 红线 | 达成 | 0070 双方言 CHECK/唯一约束；v70 钉死点一致；按 plan 显式装配；未见越界域或 admin.users 关联 |
| D-002 §4.2/§4.4：单事务顺序、回读起步、3 attempts、互异幂等键、ref 反链 | 代码结构达成 | service/store 的事务与错误分类符合冻结协议；SQLite 并发测试与 PG 同 request 测试通过 |
| VP-031 判据 2：失败零残留、并发可证明 | **未达到门禁证据要求** | 合同明确要求双数据库覆盖同 request、不同 request 余额竞争及 retry exhaustion 无残留；现有 PG 仅覆盖同 request，且两库均无可控 retry exhaustion/中间步骤失败注入 |
| D-002 §7 / VP-031 判据 1：Admin 写审计 fail-closed | 代码结构达成，测试部分覆盖 | create 审计失败回滚有测试；update/void 成功审计与重复 void 无第二条审计有测试；未见 update/void 审计失败注入，但同一事务结构可静态核对 |
| D-002 §8：公开目录桶 | 达成 | 1 分钟/60、IP key、AllowRecord、429/Retry-After、无 Clear |
| D-002 §8：购买桶 | **未实现** | 合同覆盖 Service API 购买调用且 R2 切片明确含 purchase 桶；Service 无 limiter 字段/构造参数/AllowRecord 调用 |
| D-002 §9：R2 错误码 | 当前 R2 面基本达成 | version/invalid request 已登记并映射；`BIZOFFER_ENTITLEMENT_INVALID` 按合同修订说明随 R3 Check/Consume 登记，不据此阻断 R2 |
| VP-031 判据 1：协议页面、权限、审计、公开上架列表 | 达成 | provider/schema/manifest/handler 接线存在；公开列表只投影 on_sale 且不泄漏 version |

## Findings

### F-001 · R2 强制 purchase 请求计数桶未接入，A-001 将其延期到 R3 与冻结合同冲突

- **严重度**：high
- **建议**：required
- **状态**：open
- **证据**：
  - D-002 §8 明确 `bizoffer|purchase|<subject_id>` 覆盖“Telegram `buy` 与 Service API 购买调用”，冻结为 1 分钟/10 次、每次请求 `AllowRecord`、永不 key-wide `Clear`（`.../D-002-digital-offer-contract.md:271-283`）。
  - D-002 R2 切片明确包含“§8（purchase 桶）”（同文件 `311-314`）。
  - `Service` 仅持有 repo/wallet/subjects/operations/now，`NewService` 不接收 limiter；`Purchase` 中无 `AllowRecord`（`apps/api/modules/digitaloffer/service/service.go:53-66,249-291`）。limiter 只进入 Provider/HTTP handler，当前仅公开目录创建 list limiter（`apps/api/modules/digitaloffer/provider.go:25-34,65-72`；`apps/api/internal/handler/digitaloffer.go:234-257`）。
  - A-001 F-001 把 purchase 桶称为 R3 接线项并认为“Service API 测试面无处可挂”（`03-audit/A-001-self-implementation-review.md:40-44`），但这与合同对 Service API 的显式覆盖及 R2 切片相反，且未见先行合同修订。
- **影响**：所有模块内 `Purchase` 调用当前均可绕过冻结的 10/min 请求计数门禁；“§8 已符合、R2 可关门”的主张不属实。即使当前尚无生产 HTTP/Telegram 购买入口，也不能以调用方缺失静默删除已冻结的 Service API 限流义务。
- **要求**：在合同不变前提下，将 `kernel.RateLimiterProvider`/purchase limiter 接到 R2 Service API 的统一购买入口，使用冻结 key/阈值与 `AllowRecord`，不得调用 `Clear`；补预算内、耗尽、subject 隔离、窗口恢复及无 Clear 的测试。若确需延期，必须先按治理流程修订 D-002 与 R2 分母，不能以审计意见代替合同变更。

### F-002 · §4.4 强制双数据库并发与 retry-exhaustion 零残留验收未完成

- **严重度**：high
- **建议**：required
- **状态**：open
- **证据**：
  - D-002 §4.4 冻结的 R2 验收矩阵要求双数据库验证：同 request 双发、异 offer 冲突、不同 request 竞争仅够一次的余额、retry exhaustion 无残留（`.../D-002-digital-offer-contract.md:153-163`）；R2 切片再次要求“§4.4 并发验收测试（双数据库）”（同文件 `311-314`）。
  - SQLite 测试覆盖单事务、终态失败、幂等、同 request 并发、不同 request 余额竞争、currency mismatch（`apps/api/modules/digitaloffer/service/purchase_test.go:129-369`）。
  - PostgreSQL 仅有 `TestPurchasePostgresAcceptance`，覆盖同 request 4-way 并发后的一凭证/一次扣款/一份权益/无冻结（同文件 `443-512`）；没有 PG 不同 request 余额竞争与异 offer conflict。
  - 全测试文件未提供可控 transient failure/attempt 计数，未验证终态错误不重试、最多 3 次、`purchase attempts exhausted`，也未注入 freeze 后的 deduct/purchase/entitlement 失败来直接证明 retry exhaustion 与中间失败后 ledger/purchase/entitlement 全部零残留（测试入口清单见同文件 `129,197,247,277,319,356,372,443`）。
- **影响**：核心实现的同事务结构有较强静态证据，现有 SQLite/PG 测试也证明了部分收敛；但合同要求的是资金路径跨方言、重试耗尽情况下可重复核对的证据。当前证据不足以放行该门禁，A-001 关于“双数据库验证余额竞争、分母覆盖判据 2”的表述超出实际测试覆盖。
- **要求**：至少补齐 PostgreSQL 的异 offer conflict 与不同 request 余额竞争；为 transaction runner/repository 增加可控失败点或等价测试 seam，覆盖 deduct、purchase insert、entitlement insert 失败及 3 次可重试错误耗尽，逐项断言无 wallet ledger/冻结、无 purchase、无 entitlement。相同关键用例应在 SQLite/PostgreSQL 双跑。

### F-003 · A-001 F-002 的“100 字符截断”仅按字节截断，Unicode/PG 行为未闭合

- **严重度**：medium
- **建议**：recommended
- **状态**：open
- **证据**：A-001 声明 `searchQ` 为“trim + 100 字符截断”且据此 fixed（`03-audit/A-001-self-implementation-review.md:46-49`）；当前实现使用 `len(q) > 100` 与 `q[:100]`（`apps/api/modules/digitaloffer/store/store.go:670-677`），这是字节上限，不是字符上限，可能切断多字节 UTF-8 rune。该补丁当前还是工作树修改，指定测试中未见 ASCII/多字节边界用例。
- **影响**：性能上限对 ASCII 有效，但“100 字符”声明不准确；多字节输入可能形成非法 UTF-8，PostgreSQL 路径可能拒绝参数并转为 500。故 A-001 F-002 只能视为部分修复，不能按现有表述确认完全 fixed。
- **建议**：用 rune/grapheme 明确口径并安全截断，增加 100/101 ASCII、中文/emoji、trim 与三列表共用入口测试；在提交证据落定前不要把该修复描述为已提交事实。

### F-004 · 审计索引与合同版本投影未同步

- **严重度**：low
- **建议**：recommended
- **状态**：open
- **证据**：A-001 文件已存在，但本轮前 `03-audit.md` 只有 C4 占位行，未登记 A-001；D-002 标题/修订说明与 A-001 已使用 v1.1.0，而 workspace.md、goal-tree.md、GOAL-003 `00-meta.md` 仍写 v1.0.0（`D-002-digital-offer-contract.md:10-15`；`A-001-self-implementation-review.md:12-17`；`workspace.md:36-38`；`goal-tree.md:17-19`；`GOAL-003.../00-meta.md:13`）。
- **影响**：不直接改变当前代码行为，但降低审计意见枚举完整性与实施分母可追踪性。
- **建议**：由 `/govern` 响应本意见时补登记 A-001，并把权威目标/工作区投影统一到实际采用的 D-002 v1.1.0；本 independent 审计不修改这些状态或正文投影。

## 必改项汇总

1. **F-001**：落实 D-002 §8 的 R2 purchase Service API 请求计数桶，或先合法修订合同与阶段分母。
2. **F-002**：补齐 §4.4 强制的 SQLite/PostgreSQL 双数据库并发、重试耗尽与零残留证据矩阵。

在 F-001/F-002 未按 `fixed`、用户书面 `accepted-residual` 或 `user-overruled` 合法闭合前，不建议放行 R2 资金路径门禁或将 GOAL-003 标为 done。

## 与既有意见的异同（对照 A-001）

- **一致**：认可 migration/store/service 主体结构、单事务 freeze→deduct→凭证→权益顺序、幂等回读/ref 反链、Admin 权限/审计事务、公开目录 list 桶、Profile/subject-only 红线均有真实实现；认可 `BIZOFFER_ENTITLEMENT_INVALID` 属 R3 Check/Consume 登记项。
- **不同**：A-001 将 purchase 桶视为 R3 low/recommended；本意见依据 D-002 §8 与 §11 判定其为 R2 明文义务，定为 high/required。
- **不同**：A-001 声明 §4.4 双数据库并发分母已覆盖；本意见核对测试后确认 PG 只覆盖同 request，且两库均缺 retry-exhaustion/中间失败注入，故资金 fail-closed 仍缺 required 证据。
- **不同**：A-001 把 F-002 标为 fixed；本意见确认三列表已接入 `searchQ`，但实现是 100 bytes 而非 100 characters，且无多字节边界测试，因此仅认可为部分修复。
- **补充**：A-001 尚未进入 `03-audit.md` 正式索引，且 v1.1.0 分母未同步到若干权威投影；这些不替代代码 finding，但应由编排器修复台账一致性。

## 结论 + 建议下一步

**结论：fail。** 核心购买事务代码与多项测试是真实的，但 R2 冻结合同明确要求的 purchase Service API 限流未实现，且 §4.4 双数据库/retry-exhaustion 零残留验收分母未完成；因此当前不能证明资金路径门禁已满足，也不能采纳 A-001 的无条件 pass。

建议使用 `/govern` 响应 A-002：优先修复 F-001、F-002；随后对 F-003/F-004 作修复或明确处置，并在重跑目标 build/test 与双数据库矩阵后发起 finding-closure 独立复审。

## 声明

本意见只追加 `source: independent` 审计记录与索引，不修改 GOAL-003 的 `status`、`progress`、检查点、goal-tree、合同正文、执行台账或业务代码；意见响应与状态推进由 `/govern` 处理。
