---
id: GOAL-004-r3-optimistic-concurrency-idempotency
doc: audit-entry
record_id: A-001
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
scope: R3 S0 · 乐观并发/幂等契约范围、I-001/I-002 信息门禁、D-001 决策与 E-001 扫描事实
audit_type: design-plan+execution-facts
verdict: conditional
status: recorded
parent: GOAL-001-shared-cross-module-contracts
created: 2026-08-18
updated: 2026-08-18
version: 0.1.0
---

# A-001 · R3 S0 设计/计划 + 扫描事实独立审计（2026-08-18）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high；项目级路径见 `docs/architecture/independent-audit-execution.md`）
- **类型**：stage（S0）/ design-plan + execution-facts
- **scope**：R3 S0：wallet 消费切片、版本前置条件（强 ETag / If-Match / expectedVersion / legacy version / 428·400·409）、ledger `operationId` 与 replay 语义、无 key 兼容、重复审计风险、I-001/I-002、对照 wallet repository/service/handler/migrations/tests
- **verdict**：conditional

## 范围与区间

- **工作区**：`workspace-012-shared-cross-module-contracts`（`workspace.md`：`root_goal` = `GOAL-001-shared-cross-module-contracts`；`canonical_scope` 与本目标路径一致；`shared_materials_catalog: none`；`vision_role: delivery`；`primary_plan` = `VP-012-shared-cross-module-contracts`）
- **covered**：GOAL-004 定义与非目标、S0 路线图、D-001 契约边界、E-001 已记录扫描事实、I-001/I-002、对照 `admin.wallet` 持久化/服务/HTTP/迁移/测试与 Manifest ETag、抽样 settings/dictionary/scheduledtasks 无版本 CAS
- **excluded**：S1/S2 实施、S3 关门、未重新执行测试套件、其他工作区上下文、共享资料内容（目录为 `none`）、R4 Job 状态机设计
- **本轮未复验**：未运行 `go test` / Web 测试；「测试覆盖」仅核到测试代码与断言存在，运行结果标为证据不足

## 工作区与对齐（只读）

| 检查项 | 结论 | 证据 |
|--------|------|------|
| 工作区绑定 | 通过 | `workspace.md` Root / canonical / `plan_refs`+`primary_plan` 与 GOAL-004 `parent`、`primary_plan` 一致 |
| 共享资料引用 | 无引用，不构成关闭证据 | `shared_materials_catalog: none`；GOAL-004 未引用 `material_id`/`sha256` |
| 对齐链 | 未发现与 Root R3 / VP-012 方向的明显冲突 | Root R3 = expectedVersion / ETag / 409 / idempotency_key；GOAL-004 以 wallet 为有界真实消费切片。R4 状态机已显式排除。非目标排除 Profile / Manifest 装配 / 模块矩阵 / Tier D |
| Vision Review required | 本 scope 未见开放 required | `docs/vision/reviews.md` 索引声明 open required = 0；本意见不审 Vision Review 本身 |
| 既有 Goal 审计 | 无 | `03-audit.md` 索引与 `03-audit/` 在本条之前为空 |

## 成果（有证据）

| 主张 | 证据 | 核验 |
|------|------|------|
| 目标范围/非目标已写清，且挂 VP-012 | `00-meta.md` 范围/非目标/`plan_refs` | 通过 |
| D-001 选定 `admin.wallet` 为首个真实消费模块，理由可核对 | `D-001` 决策 1；migration 0031 `version INTEGER NOT NULL DEFAULT 0` + `UNIQUE (account_id, idempotency_key)`；`store/repository.go` `UpdateStatus`/`Mutate` 均 `WHERE id=? AND version=?`，0 rows → `ErrVersionConflict`；`WithTx` 单事务 | 通过：wallet **足以**作为 R3 首个真实切片，不必新造 demo |
| 账户级 CAS 与 409 码已存在 | `repository.go` `ErrVersionConflict`；`handler/wallet.go` `writeWalletError` → 409 `LEDGER_VERSION_CONFLICT`；`errorcatalog.go` 已收录；`TestWalletIdempotencyAndStatus` 断言 stale PATCH → 409 | 通过；**本轮未跑测试** |
| 同账户同 key 同 payload 回放既有 entry，异载荷冲突 | `Mutate` 按 `(account_id, idempotency_key)` 查找，比对 `EntryType`/`AmountDelta`/`Memo`/`RefType`/`RefID`；`TestMutateIdempotency`、`TestMutateIdempotencyRefCompare`、`TestWalletIdempotencyAndStatus`（ledger `total==1`） | 通过；payload 字段未写入 D-001（F-002） |
| 无 key 旧调用保持可用 | `IdempotencyKey==""` 跳过查找；`nullIfEmpty` 把空串存成 SQL NULL；SQLite `UNIQUE` 允许多个 NULL | 通过；E-001「保留无 key」与实现一致 |
| HTTP status PATCH 只吃 JSON `version`，缺字段即 Go 零值 0 | `wallet.go` `var body struct { Status string; Version int64 }`；无 `If-Match`/`expectedVersion`；新账户 `version` 默认 0 | 通过：E-001 该条属实。缺字段在从未变更的账户上会被当成 version 0 **静默接受** |
| Manifest ETag 仅服务 `If-None-Match` 304 | `handler/manifest.go`：内容 sha256 强 ETag，只比较相等后 304；无写前置条件 | 通过：不能直接复用为 version precondition |
| Replay HTTP 无 `operationId`/`state`/`replayed` | `accountToMap`/`entryToMap`；mutate 响应为 `{account, entry}`；`entryToMap` **亦不输出** `idempotencyKey` | 通过；E-001 未点名缺 `idempotencyKey`（F-004） |
| Replay 会再写一条 wallet 业务审计 | `walletMutate` 与 `POST .../by-owner/{ownerId}/adjust` 在 `Mutate` 成功后无条件 `recordWalletEvent`；`Mutate`/`Service.Mutate` 不返回 replay 标志；`TestWalletIdempotencyAndStatus` 只验流水条数，不验 operations 计数 | 通过：E-001 风险成立，D-001 未冻结禁止重复审计（F-003） |
| settings / dictionary / scheduledtasks 无旧值条件 | `settings/repository/repository.go` `ON CONFLICT DO UPDATE` 只改字段与 `updated_at`；`datadictionary/store/repository.go` 与 `scheduledtasks/store/repository.go` 的 `UPDATE ... SET ... updated_at=? WHERE id=?` 无 version | 通过：不作为 R3 首切片成立 |
| I-002 模式已唯一为 `independent`，provider 与项目级路径一致 | `D-001` 决策 6；`00-meta` I-002 `verified`；`docs/architecture/independent-audit-execution.md` | 通过；本条不能代替 S3 关门独立审 |
| R3 未把 R4 状态机写进方案 | `D-001` 决策 5：只定义同步终态 `succeeded` 与 replay identity | 通过；与 Root「R4 依赖 R1/R3、尚未开始」一致 |

## 对照成功标准（S0 适用部分）

GOAL-004 四条成功标准均属 S1–S3 交付物。S0 只评估「是否已具备冻结契约 / 进入实施的信息」。

| 标准 | S0 状态 | 证据 |
|------|---------|------|
| 1. 单资源稳定 ETag；`If-Match` 与 `expectedVersion` 可互换且不一致时拒绝 | 未开始；格式与解析规则未冻到可实现 | 无 wallet `ETag` 头；无 GET `/api/wallet/accounts/{id}`；D-001 只写 `"v<non-negative-int>"`（F-001） |
| 2. stale 写 409；缺失/非法前置条件不被当作 version 0 | 基线部分存在（stale→409）；缺失仍是零值 0 | `Version int64`；D-001 已选 428/400，但无错误码与存在性判定（F-001） |
| 3. 同 key 同 payload 同 `operationId` 且 `replayed=true`，不重复写账本；异载荷 409 | 账本层部分存在；HTTP/审计/竞争路径未冻 | 回放返回既有 `entry.ID`；无 `replayed`；唯一约束竞争直接 409（F-002/F-003） |
| 4. shared contract、wallet 各层与兼容路径有测试；API 全量验证 | 未开始 | 现有测试覆盖 CAS/幂等/stale PATCH，不覆盖缺 version、ETag、428、重复审计 |

## 信息门禁核对（P-005）

| ID | 级别 | 最晚阶段 | 状态 | 是否到期 | 本轮结论 |
|----|------|----------|------|----------|----------|
| I-001 | required | S0 结束前 | verified | 本轮审的是「切片是否足够冻结」 | **切片判断成立**，E-001 主干事实可重复核对。`verified` 不得被读成「全部 wire 细节已冻」（F-004）。本条不要求把 I-001 改回 `collecting` |
| I-002 | required | S1 实施前 | verified | **未阻断 S0** | data/compatibility 按 P-003 可唯一判定 `independent`（无需再问模式）。provider = 项目级 grok-build（grok-4.6 reasoning high）。本条 A-001 **不能**关闭 S3 关门独立审 |

无 `deferred` 项。无用户书面 `accepted-residual`。

## Findings

### F-001 · S1 版本前置条件的 wire 语义未冻到可实现粒度

| 字段 | 值 |
|------|-----|
| level | required |
| severity | med |
| status | open |
| 影响门禁 | S0「冻结边界」完成；S1 实施 |
| evidence | `D-001` 决策 2；`00-meta` 成功标准 1–2；`handler/wallet.go` PATCH body；`handler/manifest.go`；`errorcatalog.go`（无 428/前置条件码）；`error_contract_test.go` 冻结集；`provider.go` 路由表（无 GET `/api/wallet/accounts/{id}`）；`wallet.json` `updateStatus.bodyMapping.version` |

D-001 正确选定了方向：共享 helper、强 ETag `"v<non-negative-int>"`、三来源一致、缺失 428、非法/矛盾 400、stale 409 `LEDGER_VERSION_CONFLICT`、保留 legacy `version`。独立核验后，S1 仍必须**发明**下列语义才能写代码：

1. **解析**：`If-Match` 是否允许弱标签 `W/"v1"`、列表、`*`、引号外空白；非法时一律 400 还是部分忽略。`expectedVersion` 是 JSON number、字符串 `"1"`，还是 `"v1"`。
2. **缺失 vs 0**：当前 `Version int64` 无法区分缺字段与显式 `0`。新账户默认 version 0，因此 `PATCH {"status":"disabled"}` 在从未变更的账户上会成功。D-001 要求 428，但未写「三来源皆不存在才算缺失；显式 `0`/`"v0"` 合法」。
3. **错误码**：R1 已有错误包络。409 复用 `LEDGER_VERSION_CONFLICT` 成立；428/400 **没有** catalog / 冻结码。S1 若临时起名会冲击错误契约。
4. **ETag 出现面**：成功标准写「单资源响应具有稳定 ETag」。现状 version 只在 JSON body；无 GET-by-id；`GET /api/wallet/me` 是单元素 list 包络；列表不能用账户 version 做强 ETag。须写明：哪些读/写设 `ETag` 头、列表是否不做强 ETag、mutation 响应是否带账户 ETag。
5. **mutation 不收前置条件**：D-001 把 `expectedVersion` 限于 status PATCH。账本写走服务端 CAS + 可选 key。这与 VP-012「资源更新用版本、操作用幂等」可并存，但必须写死，避免 S1 把 If-Match 扩到 adjust/freeze。

在补 D-001 修订或 D-002 之前，**不得开始 S1 实现**。

### F-002 · Replay 契约未冻：payload 指纹、唯一约束竞争、resourceVersion、失败不落 operationId

| 字段 | 值 |
|------|-----|
| level | required |
| severity | med |
| status | open |
| 影响门禁 | S0 契约边界；S2 实施 |
| evidence | `repository.go` `Mutate` 比对字段与 unique-violation 分支；`Service.Mutate` 每次预生成新 `entryID`；`entryToMap`；`00-meta` 成功标准 3；`D-001` 决策 4–5 |

D-001 用 ledger `entry.id` 作为 durable `operationId`、成功态仅 `succeeded`，**足以满足 R3 且不侵入 R4**（不必引入 queued/running/failed/cancelled 或 Job 表）。但下列边界仍会让 S2 静默选型：

1. **「同 payload」字段**：实现比对 `entry_type` / `amount_delta` / `memo` / `ref_type` / `ref_id`，**不含** `actor_id`/`actor_name`。不同操作者复用同一 key 会被当成 replay。共享契约必须点名指纹，不能只写「同 payload」。
2. **同 key 同 payload 并发**：先 `SELECT` 再 `INSERT`；唯一约束冲突直接 `ErrIdempotencyConflict`，**不回读比对**。两个并发相同请求可能 409，而不是同一 `operationId` + `replayed=true`。E-001 只写「竞争路径保持回滚、不声称跨库锁」，未裁定：客户端可先 409 再重试得到 replay，还是实现必须在 unique 冲突后回读。这与成功标准 3 字面冲突。
3. **`resourceVersion`**：账本行没有 version 快照；replay 返回的是**当时**账户 version。后续 mutation 之后回放，该值 ≠ 原写入时的 version。须冻结「响应时刻的当前账户 version」，并说明不能拿历史 replay 的 version 去 If-Match。
4. **失败不铸造 operation**：余额不足 / disabled / invalid 在 INSERT 前返回，key 不落库；同 key 重试会再执行。R3 同步终态只有 `succeeded` 时这是正确的，但必须写明：失败不是 durable operation，不占用 key。
5. **可观测性**：`entryToMap` 不输出 `idempotencyKey`；HTTP 也无 `operationId`/`state`/`replayed`/`resourceVersion`。S2 若只加字段名而不改 handler 映射，成功标准 3 仍不可测。

这些应与 F-001 一并写入 D-002（或修订 D-001），再声称 S0 边界已冻。

### F-003 · Replay 禁止重复业务审计未写入 D-001

| 字段 | 值 |
|------|-----|
| level | required |
| severity | med |
| status | open |
| 影响门禁 | S0 契约边界；S2；R2 审计模型消费 |
| evidence | `handler/wallet.go` `walletMutate` L403–405、by-owner adjust L172–173；`Service.Mutate` 无 replay 标志；`E-001` 风险第 3 条；`TestWalletIdempotencyAndStatus` 只断言 entries `total==1` |

E-001 已写：「不能生成第二个成功 operation 或重复业务审计」。D-001 只要求返回原 entry id 并标 `replayed=true`，没有把「replay **不得**再写 `wallet.adjust`/`freeze`/… 业务事件」冻进契约。当前 handler 在 `Mutate` 成功后无条件 `recordWalletEvent`；仓库回放与首次成功对 handler 不可区分。

R3 若宣称 operation replay 可观测，重复审计会让同一 `operationId` 对应多条业务事件，削弱 R2 结构化审计的「一次成功一次事件」假设。S2 必须：仓库/服务返回 `replayed`；handler 在 `replayed=true` 时跳过业务审计（或写明确的 replay 标记且不新增成功事件）；测试断言 operations 计数不增加。

### F-004 · I-001 `verified` 只覆盖「切片足够」；扫描未登记若干可测缺口（recommended）

| 字段 | 值 |
|------|-----|
| level | recommended |
| severity | low |
| status | open |
| 影响门禁 | 台账读法；不单独阻断 I-001 |
| evidence | `00-meta` I-001 证据栏；`E-001`；本条成果表 |

I-001 问的是现有实现与真实切片是否足以冻结契约。答案对「选 wallet、不要 settings/tasks/dictionary、不要 Manifest ETag」为是；E-001 主干也属实。编排器**不要**把 `verified` 读成 F-001～F-003 已关闭。建议在响应里把 I-001 证据改成「切片足够，见 E-001；wire 细节见 A-001 / 后续 D-002」，状态可保持 `verified`。

E-001 还可补记（不构成 I-001 重开）：无 GET-by-id；`entryToMap` 无 key；缺 version / 重复审计 / 同 key 并发无测试。

## 必改项汇总

1. **F-001**：在 S1 开工前落盘可实现的版本前置条件（ETag 解析与拒绝集、三来源存在性、428/400 错误码、ETag 出现面、mutation 不收 If-Match）。
2. **F-002**：冻 replay 契约（payload 指纹、unique 竞争是回读还是允许 409 后重试、`resourceVersion`=响应时刻账户 version、失败不落 operationId、HTTP 结果字段）。
3. **F-003**：冻「同 key 同 payload replay 不得再写第二条成功业务审计」，并列为 S2 必测。

## 与既有意见的异同

无既有 self / independent 条目。本条为 GOAL-004 第一条意见。不与同区 GOAL-002 / GOAL-003 已关闭意见冲突。

与同区 GOAL-003 A-001 的对照（仅方法论，不引入其 finding 状态）：R2 S0 时 I-002 模式未唯一；本目标 D-001 已唯一为 `independent`，故 I-002 不再作为本条 required。

## 结论 + 建议给编排器/用户的下一步

**verdict = conditional。** wallet 作为真实消费切片成立；E-001 对 CAS / 409 / 账户作用域 key / 事务 / Manifest 只读 ETag / 缺字段零值 / replay 再审计的扫描可核对；I-002 模式与 provider 选择正确；ledger `entry.id` + 单一 `succeeded` 能承载 R3 且不侵入 R4。但 S0「冻结边界」在 S1 wire 与 S2 replay/审计上仍有实施者必须发明的语义，不能无条件放行 S1。

建议 `/govern`：

1. 响应本条 A-001；**不要**把 S0 标完成，**不要**进入 S1 实现。
2. 先写 D-002（或修订 D-001）覆盖 F-001～F-003；F-004 仅回写 I-001 证据栏口径。
3. I-002 保持 `verified`；S3 仍须另一次 independent 关门审。本条不是关门审。
4. 本意见不改 `status`/`progress`/goal-tree。

## 声明

本意见 `source: independent`，不修改目标 `status` / `progress` / 检查点 / 方案正文 / `goal-tree.md` / Root / `workspace.md`。响应、finding 闭合与阶段推进由 `/govern` 处理。保证等级为框架默认 **L0**（入口分离），不得表述为第三方鉴证。
