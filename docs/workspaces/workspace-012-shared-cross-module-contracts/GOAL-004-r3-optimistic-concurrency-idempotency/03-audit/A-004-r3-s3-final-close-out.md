---
id: A-004-r3-s3-final-close-out
doc: audit-entry
record_id: A-004
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
scope: R3 S3 final close-out：D-002、checkpoint 08dcec8、E-003/A-003 与四条成功标准
audit_type: close-out
verdict: pass
status: recorded
parent: GOAL-004-r3-optimistic-concurrency-idempotency
created: 2026-08-18
updated: 2026-08-18
version: 0.1.0
---

# A-004 · R3 S3 independent final close-out（2026-08-18）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high；项目级路径见 `docs/architecture/independent-audit-execution.md`；本目标 D-001 决策 6 / I-002）
- **类型**：close-out / S3 关门独立审计
- **scope**：R3 S3 final close-out：核验 D-002、checkpoint `08dcec8`、E-003/A-003 与四条成功标准。重点：强 ETag 语法和拒绝集；If-Match / expectedVersion / legacy version 三来源存在性与一致性；428/400/409 error contract；ETag 出现面与 mutation 边界；payload 指纹含 `actor_id` 不含 `actor_name`；unique race 回滚后回读；`operationId`/`state`/`replayed`/`idempotencyKey`/`resourceVersion`；失败不铸造 operation；replay ledger 与 wallet 业务审计均不重复；兼容无 key / legacy version
- **verdict**：pass

## 范围与区间

- **工作区**：`workspace-012-shared-cross-module-contracts`（`workspace.md`：`root_goal` = `GOAL-001-shared-cross-module-contracts`；`canonical_scope` 与本目标路径一致；`shared_materials_catalog: none`；`vision_role: delivery`；`primary_plan` = `VP-012-shared-cross-module-contracts`）
- **covered**：D-002 冻结条款 1–5 + replay 1–6；`08dcec8` 九个 API 文件；E-003 / A-003 主张；GOAL-004 四条成功标准；I-001 / I-002；A-001 required 闭合是否仍真实；本轮独立复跑定向与 `apps/api` 全量测试
- **excluded**：将 GOAL-004 / Root R3 标为 `done`（本意见不改 status / progress / 路线图 / 方案正文 / `goal-tree.md` / `workspace.md`）；R4 Job 状态机；为 settings / scheduledtasks / dictionary 补 version CAS；新增 GET-by-id；其他工作区上下文；共享资料内容（目录为 `none`）
- **本轮复验**：
  - 定向：`go test ./internal/concurrency ./internal/modules/wallet/store -count=1` exit 0（concurrency 0.445s；wallet/store 1.212s）
  - 定向 handler：`go test ./internal/handler -count=1 -run "Wallet|ErrorContract|ErrorCatalog"` exit 0（12.927s）
  - 全量：`apps/api` `go test ./... -count=1` exit 0（handler 227.388s；operationlog 12.396s；wallet/store 1.223s；docscheck 0.755s）
- **HEAD**：`3c0fc50`（干净工作树）。实现停在 `08dcec8`；其后仅文档提交 `3c0fc50`（E-003 / A-003 / 索引），`git diff 08dcec8 HEAD -- apps/api` 为空

## 工作区与对齐（只读）

| 检查项 | 结论 | 证据 |
|--------|------|------|
| 工作区绑定 | 通过 | `workspace.md` Root / canonical / `plan_refs`+`primary_plan` 与 GOAL-004 `parent`、`primary_plan` 一致 |
| 共享资料引用 | 无引用，不构成关闭证据 | `shared_materials_catalog: none` |
| 对齐链 | 未见与 Root R3 / VP-012 有界切片冲突 | Root R3 = expectedVersion / ETag / 409 / idempotency_key；GOAL-004 以 wallet 为首个真实消费切片；R4 状态机仍排除 |
| Vision Review required | 本 scope 未见开放 required | `docs/vision/reviews.md` 声明 open required = 0；本意见不审 Vision Review 本身 |
| 既有 Goal 审计 | A-001 required 仍由 A-002 `fixed`；A-003 self `pass` 待本条独立复核 | `03-audit.md` 索引；本条之前开放 required = 0 |

## Checkpoint 核验

| 提交 | 内容 | 与 E-003 / A-003 主张对照 |
|------|------|--------------------------|
| `08dcec8` `feat(workspace-012): implement R3 concurrency idempotency contract` | `concurrency/version.go`+test；wallet handler/provider/store + tests；errorcatalog + frozen error contract | **属实**。九个 API 文件，+401/−103。本轮工作树相对该提交的 API 差异为空 |
| `3c0fc50` `docs(workspace-012): record R3 implementation audit` | E-003、A-003、GOAL-004/Root/workspace/goal-tree 索引元数据 | **属实**。纯文档；不改变 `08dcec8` 实现 |

A-003 写「实现 checkpoint `08dcec8`」与「`go test ./... -count=1` exit 0」可重复核对。本轮独立复跑再次 exit 0。

## 成果（有证据）

### D-002 版本前置条件

| 条款 | 核验 | 证据 |
|------|------|------|
| 强 ETag 唯一合法形式 `"vN"`，N 非负十进制 int64；允许整个 header 前后空白 | 通过 | `concurrency.ETag(12) == "v12"`；`parseETag` 要求首尾 `"`、前缀 `v`、digits 不含空白/`,`/`+`/`-`，再 `ParseInt` 十进制 int64 且 `>=0`。`TestResolveExpectedVersion/header` 覆盖 `" v2 "` |
| 拒绝弱标签、`*`、列表、多余引号、未加引号、负数、溢出、标签内部空白 | 通过（实现）；表测覆盖除「多余引号」外的拒绝集 | 弱标签 `W/"v1"` 首字符不是 `"`；`*` 长度/引号失败；列表 `"v1", "v2"` 因 `,` 拒绝；未加引号 `v1` 拒绝；`"v-1"` 含 `-`；溢出 `ParseInt` 失败；`"v 1"` 内部空白。多余引号 `""v1""` 因 `HasPrefix(tag[1:], "v")` 失败，实现 fail closed，但表测未列该案（F-001） |
| 三来源：`If-Match`、JSON number `expectedVersion`、legacy JSON number `version`；指针区分缺失与显式 `0` | 通过 | PATCH body `ExpectedVersion *int64` / `Version *int64`；`ResolveExpectedVersion` 对 nil 跳过、对 `*0` 计入。helper：`expected zero`、`legacy`、`all agree`、`missing`、`mismatch` |
| 字符串或其他 JSON 类型 400 | 通过 | `json.Decoder` 解到 `*int64` 失败 → `INVALID_WALLET_BODY` 400。D-002 要求 400，未钉死必须是 `INVALID_PRECONDITION` |
| 皆缺失 428 `PRECONDITION_REQUIRED`；header 非法或来源矛盾 400 `INVALID_PRECONDITION`；CAS stale 409 `LEDGER_VERSION_CONFLICT` | 通过 | handler 映射与 catalog 双语 + frozen 字面量。HTTP：缺字段 428；`expectedVersion:1` + `If-Match: "v0"` → 400；legacy `version:0` 在已 bump 账户上 → 409 |
| 单账户成功响应设 ETag；列表不设；不新增 GET-by-id | 通过 | `writeWalletAccount` / `writeWalletMutation` 设 `concurrency.ETag(account.Version)`。出现面：GET/POST by-owner、POST account、PATCH status、by-owner adjust、account-id mutation。`GET /api/wallet/accounts` 走 `writeJSON`（只设 Content-Type）。`/me` 与 `/me` POST 仍是单元素 list 包络，D-002 未列入。provider 无 `GET /api/wallet/accounts/{id}` |
| ledger mutation 不消费 If-Match / expectedVersion | 通过 | `walletMutate` / by-owner adjust body 无 version 字段，不调用 `ResolveExpectedVersion`。资源属性更新用前置条件；余额写用 durable key + 仓库内 CAS |

### D-002 幂等 replay

| 条款 | 核验 | 证据 |
|------|------|------|
| 指纹 = `entry_type`/`amount_delta`/`memo`/`ref_type`/`ref_id`/`actor_id`；不含 `actor_name` | 通过 | `sameIdempotencyPayload`；`TestMutateIdempotency`：换 `ActorID` → `ErrIdempotencyConflict`；只换 `ActorName` → replay 成功 |
| unique 竞争先回滚本事务，再按 `(account_id,idempotency_key)` 回读；同指纹 replay，不同指纹或无法回读 409 | 实现通过；无直接可执行测试 | INSERT unique → `errIdempotencyRace`；`WithTx` 原样返回错误并 `Rollback`；随后 `replayAfterIdempotencyRace` 新事务回读。`errors.Is` 不被包装。无注入 unique-violation / 并发测试（F-001） |
| HTTP `operation`：`operationId`=ledger entry id，`state=succeeded`，`replayed`，可选 `idempotencyKey`，`resourceVersion`=响应时刻账户 version | 通过 | `writeWalletMutation`；`TestWalletIdempotencyAndStatus` 两次 adjust 断言同一 `operationId`、`state`、key、`resourceVersion==1`、第二次 `replayed=true`、ETag 仍 `"v1"` |
| 失败在 insert 前退出：不铸造 durable operation、不占用 key | 通过 | `Mutate` 在 disabled / `Apply`（invalid/insufficient）之后才 UPDATE+INSERT。handler 失败走 `writeWalletError`，无 `operation`。既有 over-freeze / disabled 测试不增 ledger 行。Service 预生成的 `entryID` 失败时丢弃 |
| 同 key 同 payload replay 不再写第二条成功 wallet 业务审计 | 通过 | handler `if !replayed { recordWalletEvent }`。`TestWalletIdempotencyAndStatus`：`ListOperationsFiltered(EventWalletAdjust)` `opTotal==1` 且 ledger `total==1` |
| 无 key 仍返回本次 succeeded operation，`replayed=false`，不承诺跨请求 replay | 通过 | 空 key 跳过查找；`nullIfEmpty` 存 NULL；SQLite UNIQUE 允许多个 NULL。`writeWalletMutation` 省略空 `idempotencyKey`。`TestWalletLifecycleAndAdjustFlow` 无 key adjust/freeze/unfreeze 成功（未钉死 operation 字段，见 F-002） |

### A-003 / E-003 主张复核

| 主张 | 本轮独立核验 | 结论 |
|------|--------------|------|
| 共享 `"vN"` + 三来源 + 428/400/409 已进 catalog | 代码、表测、HTTP、frozen 字面量与 `TestWalletErrorCodesCataloged` | **属实** |
| 指定单账户响应设 ETag；列表不设 | `writeWalletAccount`/`writeWalletMutation` vs list `writeJSON` | **属实** |
| replay 标志、指纹含 actor_id、竞争回滚回读、失败无 operation | Service `entry.ID != entryID`；指纹测试；race 代码；失败路径在 INSERT 前 | **属实**（race 无单测，不推翻实现） |
| replay 审计与 ledger 双计数 = 1 | HTTP 测试硬断言 | **属实** |
| A-001 F-001～F-004 实现层已覆盖 | 对照 D-002 + 本轮代码/测试 | **属实**；S0 闭合不重开 |
| checkpoint `08dcec8`；API 全量通过 | 提交存在；本轮 `go test ./... -count=1` exit 0 | **属实** |

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| 1. wallet 单资源响应具有稳定 ETag；`If-Match` 与 `expectedVersion` 可互换且不一致时拒绝 | pass | POST by-owner ETag `"v0"`；adjust `"v1"`；legacy PATCH `"v2"`；`expectedVersion:2` → `"v3"`；仅 If-Match `"v3"` → `"v4"`；矛盾来源 400 `INVALID_PRECONDITION` |
| 2. stale 写 409；缺失或非法前置条件不被当作 version 0 静默接受 | pass | 缺字段 428（不再把缺省 0 当合法期望）；非法/矛盾 400；stale legacy `version:0` 409；显式 `0` 在 helper 中合法 |
| 3. 同 key 同 payload replay 返回相同 `operationId` 且 `replayed=true`，不重复写账本；异载荷 409 | pass | 双请求同一 `operationId`；ledger `total==1`；`wallet.adjust` 审计 `==1`；异载荷 409 `LEDGER_IDEMPOTENCY_CONFLICT` |
| 4. shared contract、wallet repository/service/HTTP 与兼容路径均有测试；API 全量验证通过 | pass | concurrency 表测；store 幂等/CAS/actor 指纹；handler 428/400/409/replay/legacy；无 key 生命周期仍成功；本轮全量 exit 0 |

## 信息门禁核对（P-005）

| ID | 级别 | 最晚阶段 | 状态 | 是否到期 | 本轮结论 |
|----|------|----------|------|----------|----------|
| I-001 | required | S0 结束前 | verified | 否（S0 已过且已 verified） | 切片判断仍成立。本条不重开 |
| I-002 | required | S1 实施前 | verified | 否 | 模式 `independent`；provider = 项目级 grok-build（grok-4.6 reasoning high）。**本条即 S3 关门独立审**，不能由 A-003 self 替代 |

无 `deferred` 项。无用户书面 `accepted-residual`。无到期且影响本 scope 的开放 required 信息项。

## Findings

### F-001 · unique 竞争回读路径没有可执行回归

| 字段 | 值 |
|------|-----|
| level | recommended |
| severity | low |
| status | open |
| 影响门禁 | 不阻断 S3 关门；回归锁定 |
| evidence | `repository.go` `errIdempotencyRace` / `replayAfterIdempotencyRace`；`store.go` `WithTx` 错误原样返回并 Rollback；`repository_test.go` 无 unique-violation 注入或并发用例 |

D-002 要求「先回滚本事务，再按 `(account_id,idempotency_key)` 回读」。实现结构正确：unique 失败不会提交已执行的余额 UPDATE；同指纹回读为 replay，读失败或指纹不同 fail closed 为 409。顺序 replay（先 SELECT 命中）已覆盖成功标准 3。缺少的是把 unique 冲突分支锁进测试，后续改 `WithTx` 包装或 `isUniqueViolation` 匹配串时可能静默退化成 500 或误 409。不构成实现名不副实。

### F-002 · 若干 D-002 出现面只经代码复核，HTTP/表测未钉死

| 字段 | 值 |
|------|-----|
| level | recommended |
| severity | low |
| status | open |
| 影响门禁 | 不阻断 S3 关门 |
| evidence | `version_test.go` 拒绝集未含多余引号；`wallet_test.go` 只断言 POST by-owner / account-id adjust 的 ETag，未断言 GET by-owner、list 无 ETag、mutation 忽略 If-Match；`TestWalletLifecycleAndAdjustFlow` 无 key 成功但未断言 `operation.replayed==false` 且省略 `idempotencyKey`；Web `en-US.json`/`zh-CN.json` 有既有 wallet 错误键，缺 `error.preconditionRequired` / `error.invalidPrecondition` |

实现与 D-002 一致，API catalog 已双语覆盖 428/400，HTTP 包络走 server-side locale。这些是回归与 Web messageKey 回退缺口，不是契约违反。API 成功标准 4 已由 shared/store/handler/legacy 路径与全量测试满足。

## 必改项汇总

无。本条开放 required = 0。

## 与既有意见的异同

- **A-001**（independent，conditional）：S0 wire/replay 未冻到可实现。A-002 以 D-002 `fixed` F-001～F-003。本条独立核验实现与 D-002 一致，**不重开** A-001 required。
- **A-002**（self，pass）：只放行 S1/S2 实施。闭合声明仍真实。
- **A-003**（self，pass）：S1/S2 实现 close-out。本条同意其四条成功标准结论与 checkpoint 主张；补记 recommended 测试/i18n 缺口，不把 A-003 降为不实。
- 无与 A-003 相反的 verdict，无需 P-004 冲突裁决。

## 结论 + 建议给编排器/用户的下一步

**verdict = pass。** D-002 冻结的强 ETag、三来源、428/400/409、ETag 出现面、mutation 边界、指纹、失败不落 operation、replay 双计数与无 key / legacy 兼容均能在 `08dcec8` 代码与本轮测试中重复核对。A-003 主张属实。开放 required = 0；无到期 required 信息项。

建议 `/govern`：

1. 响应本条 A-004。F-001 / F-002 为 recommended，不阻断关门；若本波要补测/补 Web i18n，可在响应里选做或显式留下一波。
2. 用户确认后可将 GOAL-004 标 `done`，并按 Root 纲领把 R3 标完成。本意见**不**改 status / 路线图 / goal-tree。
3. 不要把本条 recommended 当成未闭合必改。

## 声明

本意见 `source: independent`，不修改目标 `status` / `progress` / 检查点 / 方案正文 / `goal-tree.md` / Root / `workspace.md`。响应、finding 闭合与阶段推进由 `/govern` 处理。保证等级为框架默认 **L0**（入口分离），不得表述为第三方鉴证。
