---
id: A-004
goal: GOAL-019-r3-s14-wallet-ledger
title: S1 方案冻结独立审计 · admin.wallet
date: 2026-08-16
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
scope: S1 方案冻结（D-002 + I-001~I-004 + F-001/F-002 响应 + data 门禁）
audit_type: design-plan
verdict: conditional
status: recorded
parent: GOAL-019-r3-s14-wallet-ledger
created: 2026-08-16
updated: 2026-08-16
version: 1.0.0
---

# A-004 · 独立审计意见（S1 方案冻结 · S-14 钱包/账务）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型** / **scope**：design-plan · S1 方案冻结（D-002 全文、I-001~I-004 闭合证据、019-F-001/019-F-002 响应、data 门禁：金额语义 / 流水不可变与快照链 / 双层审计 / 迁移 0031/0032 / 权限键拆分 / 协议独立对照）
- **verdict**：**conditional**

## 范围与区间

- **工作区**：`workspace-011-admin-functional-modules`（`root_goal: GOAL-001-admin-functional-modules`；`canonical_scope` 匹配；`plan_refs`/`primary_plan` = `VP-011-admin-functional-modules`；`shared_materials_catalog: none`）。未读取或比较其他工作区目标状态。
- **已通读**：本目标 `00-meta`、`01-decision.md`、D-001、D-002、E-001、E-002、`03-audit.md`、A-001～A-003；Root `00-meta` R3 行；`goal-tree.md` GOAL-019 行；I-011-001 §4 S-14 / §7 / §8；protocol-inventory-v2.7.0 §2.5。
- **代码核对**（按 D-002 引用验伪）：`datapermission/provider.go` + `manifest/fragment.json` + `schema/data-permission.json` + `migration/migration.go`；`recyclebin/migration/migration.go`；`operationlog/repository.go` EventXxx；`operationlog/migration/migration.go`（max Version 30）；`store/store.go` `WithTx`；`kernel/profile.go` ProfileAdmin；`kernel/provider.go` DefaultNavigationOrder；`composition/composition.go` L321–328；`composition/composition_test.go` L471；`handler/error_contract_test.go`；`handler/mfa.go` EventMFAEnroll 记录位置；`modules/compiled/persistence.go`。
- **covered**：方案可实施性、data 门禁（余额变动 + 迁移）、I-001~I-004 是否真实闭合、019-F-001/F-002 响应、协议独立对照、迁移编号、无越界、S2 放行条件。
- **excluded**：S2 实现、S3～S5、不改 status/progress/goal-tree/方案正文。本条 **不是** 立项复审；A-002 不代替本条。
- **保证等级**：L0（入口分离）。不得解读为第三方鉴证。

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| 模块五段式先例存在：provider / fragment / schema / migration / compiled 注册 | `datapermission/provider.go` L116–186；`manifest/fragment.json`；`schema/data-permission.json`；`datapermission/migration/migration.go` L19–46（Version 27）；`compiled/persistence.go` L12 / L30 |
| 领域表 DDL 先例（无 FK、CHECK/UNIQUE/索引） | `recyclebin/migration/migration.go` L19–43（Version 25） |
| `store.WithTx` 存在且为「BeginTx → fn → Commit/Rollback」平台事务边界 | `store/store.go` L67–81；`SetMaxOpenConns(1)` 在 L49 |
| operationlog EventXxx 命名风格：`{domain}.{action}`，多词 kebab | `operationlog/repository.go` L14–59（如 `EventMFAEnroll = "mfa.enroll"`、`EventDataPermissionPolicyUpdate = "data-permission.policy-update"`） |
| 六事件名 `wallet.account-create` / `account-update` / `adjust` / `freeze` / `unfreeze` / `reconcile` 与既有风格同构 | D-002 §2 L28；对照 repository.go L25–59 |
| 双层审计「handler 层独立记录、非同一事务」与 MFA 先例一致 | D-002 §2 L29；`handler/mfa.go` L119–127（`Enroll` 成功后 `recordMFAEvent`，`_ = operations.RecordOperation`） |
| 迁移 0031/0032 相对当前最大 **30** 不冲突 | `operationlog/migration/migration.go` L262–266 `Version: 30`（MFA CHECK 重建）；全仓无 Version 31/32；0028/0030 为 CHECK 超集重建先例（L127–155） |
| ProfileAdmin 为 admin 默认集内容扩展先例；mvp/demo 不含 data-permission/mfa 等同款 | `kernel/profile.go` L46–89（ProfileAdmin 段注释）；mvp L26–45 / demo L95–109 无 S 系列领域模块 |
| DefaultNavigationOrder 现长 **13**，尾部追加 `menu_wallet` 与「13→14」一致 | `kernel/provider.go` L403–418（13 项，末项 `menu_data_permission`） |
| composition 接线位置：datapermission 在 L326 附近按 `plan.HasModule` 追加 | `composition/composition.go` L321–328 |
| 错误码契约分类：frozen literals + frozen domain | `handler/error_contract_test.go` L19–77；D-002 §6 L110 候选码须在 S2 收敛进该测试 |
| 页面 requiredCapabilities 与 data-permission schema 同类（非 fragment 短列表） | D-002 §5 L103；`schema/data-permission.json` L6–14 八键；`fragment.json` 仅 `app.manifest`/`app.navigation`（与「页 schema」主张一致） |
| 协议对照独立口径：§2.5 无钱包专属面；样例为上游 `_samples/` | protocol-inventory L139–157（`order-list-batch` / `order-detail-lifecycle`）；I-011-001 L97–101（不得外推 9/0）；D-002 §5 L100–104 未声称协议覆盖、不改 pin v2.8.0、不接入 data-transfer |
| I-011-001 §4 S-14 = 余额/流水/对账 + 余额变动审计 + 迁移基建 | I-011-001 L71；D-002 标题与 §1–§2 覆盖该最小面 |
| 019-F-001 已裁定：`ref_type`/`ref_id` 可选空引用，对账仅账本内部勾稽 | A-002 F-001；D-002 §1 L23；`03-audit.md` L38 |
| 019-F-002 已落实：`wallet.adjust` 与 `wallet.write` 拆分；金额变动三端点走 adjust | A-002 F-002；D-002 §3 L40–46；`00-meta` I-004 L45 |
| 整数最小单位 + 账户表恒等式 CHECK + version 乐观锁已写入方案 | D-002 §1 L20–21；§4 L59–67 |
| 流水无 UPDATE/DELETE 路径（应用层不可变）已声明 | D-002 §1 L22 |
| 无越界：不引入转账/外部支付/多租户；不改装配语义/协议 pin | D-002 §5 L104、§7 L116–121；D-001 L20 |

## 对照成功标准（S1）

| 标准 | 状态 | 证据 |
|------|------|------|
| 账务模型（余额口径 / 精度 / 并发）可实施 | 部分 | 三余额恒等式、INTEGER 分、version + WithTx 已写清。`amount_delta` 对 freeze/unfreeze 与「语义 = 总额变动」+ `CHECK != 0` 互否（F-001） |
| 流水不可变 + 快照链对账可验证 | 部分 | 不可变声明 + 三余额快照列存在。缺按 `entry_type` 的可执行重放表；链序仅 `created_at`（秒级，`now.Unix()` 先例）未定义并列键（F-001 / F-003） |
| 幂等键自洽 | 部分 | 可选 + UNIQUE +「返回既有流水」已写；范围为全局、同键异载荷未裁定（F-002） |
| 双层审计真实（operationlog 模式存在） | 满足 | EventXxx + 0028/0030 CHECK 重建 + handler 非同事务残余已文档化 |
| 迁移 0031/0032 不与现序列冲突 | 满足 | 当前 max = 30 |
| 权限键拆分落实 019-F-002 | 满足 | `wallet.read` / `write` / `adjust`；金额变动专用 adjust |
| 协议对照独立（不沿用 9/0） | 满足 | D-002 §5 + inventory §2.5；本地领域模块 |
| I-001 required 设计层闭合、未伪装实现已落地 | 部分 | 闭合对象是方案不是运行账本；对账/金额按类型语义未写满即标 `verified`（F-001） |
| I-002/I-003 required 设计层闭合 | 满足 | D-002 §2 / §5 可重复核对；未把未实现迁移/协议面写成已落地 |
| I-004 non-blocking（Profile/命名）闭合；写路径键已拆出 | 满足 | D-002 §3；admin 默认集 + `admin.wallet` / `menu_wallet` |
| 019-F-001 / 019-F-002 已响应 | 满足 | 空引用裁定 + 三键冻结 |

## Findings

### F-001 · `amount_delta` 语义与 freeze/unfreeze / CHECK 互否；快照链缺可执行重放规则

| 字段 | 值 |
|------|-----|
| level | **required**（med） |
| status | open |
| evidence | D-002 §1 L22：`amount_delta`「**语义 = 总额变动**」且「freeze/unfreeze 变动 available/frozen **且总额不变**」。D-002 §4 L73：`CHECK (amount_delta != 0)`。若总额变动语义对 freeze/unfreeze 成立，则 `amount_delta` 必为 0，与 CHECK 互否。§6 L109 只写「恒等式 + 快照链校验（连续条目 balance_after 衔接）」，未给出按 `entry_type` 的 apply 表（符号、作用列、拒绝条件、首笔相对 (0,0,0) 的衔接）。I-001 证据栏指向 D-002 §1（`00-meta` L42；`01-decision.md` L17）并标 **verified** |
| closure | — |
| 影响门禁 | S1 方案冻结 / S2 实施（data · I-001 对账语义） |

整数分、账户表恒等式、乐观锁成立，**不是**整份模型名不副实。缺的是金额变动原语的唯一解释：S2 若按「一律总额变动」实现则无法落 freeze；若按「freeze 转移额度」实现则须自造符号与重放，对账测试无法对照冻结稿。

S1 必须补一张可核对表后再放行 S2（或用户书面 residual，不推荐）：

| entry_type | amount_delta 含义 | 符号 | total | available | frozen | 拒绝 |
|------------|-------------------|------|-------|-----------|--------|------|
| adjust | 总额变动 | ≠ 0 | += d | += d | 不变 | 结果任一口径 < 0 或恒等式破（§6 已写 available 不足） |
| freeze | available→frozen 额度 | 须裁定（建议 > 0） | 不变 | -= \|d\| | += \|d\| | available 不足；disabled |
| unfreeze | frozen→available 额度 | 须裁定（建议 > 0） | 不变 | += \|d\| | -= \|d\| | frozen 不足；disabled？ |

快照链验收应同时写死：按该表 `apply(prev, entry) == (after_*)`；末笔快照 == 账户三余额；每笔快照自身满足恒等式。

### F-002 · `idempotency_key` 全局 UNIQUE +「返回既有流水」未限定本账户 / 同载荷

| 字段 | 值 |
|------|-----|
| level | **required**（med） |
| status | open |
| evidence | D-002 §1 L21：「可选 `idempotency_key`（**UNIQUE**；**重复提交返回既有流水**）」；§4 L79：`idempotency_key TEXT UNIQUE`（表级、无 `account_id` 复合）。§6 L110 同时列 `LEDGER_IDEMPOTENCY_CONFLICT` 候选，与「一律返回既有」未分流。I-001 含「幂等与并发控制」且已 `verified` |
| closure | — |
| 影响门禁 | S1 方案冻结 / S2 实施（data · I-001 幂等） |

按字面实现：账户 B 复用账户 A 的 key → 全局 UNIQUE 命中后若「返回既有流水」会跨账户泄露流水；若直接 409，则两账户不能共用同一客户端 key。同 key 不同金额（重放篡改）也未规定是返回旧单还是 CONFLICT。

S1 须裁定并写入 D-002：

1. UNIQUE 范围 = `(account_id, idempotency_key)`（建议）或显式全局并禁止跨账户返回；
2. 同账户 + 同 key + 同载荷 → 返回既有；同账户 + 同 key + 异载荷 → `LEDGER_IDEMPOTENCY_CONFLICT`；
3. 查找必须带 `account_id`，禁止按裸 key 取他户流水。

### F-003 · 快照链排序键未定义；ledger 快照无恒等式 CHECK

| 字段 | 值 |
|------|-----|
| level | recommended（med） |
| status | open |
| evidence | D-002 §4 L85 索引 `(account_id, created_at DESC)`；`created_at INTEGER`。既有模块写秒级 Unix：`datapermission/store/repository.go` L135 / L210 `now.Unix()`。同秒两笔（同事务连续写或并发重试）时「连续条目」无并列键。账户表有 `CHECK (balance_total = available + frozen)`（D-002 L67）；`wallet_ledger_entries` 三 `balance_after_*` **无**对应 CHECK（L74–76 仅各自 `>= 0`） |
| closure | — |

建议：链序 `(created_at ASC, id ASC)`（ULID 可排）或只按 `id`；ledger 快照加恒等式 CHECK。不单独阻断——若 F-001 重放表已写序，本条可并入勘误。

### F-004 · disabled 拒调账与冻结，未写 unfreeze

| 字段 | 值 |
|------|-----|
| level | recommended（low） |
| status | open |
| evidence | D-002 §1 L24：「disabled 时拒绝调账与冻结，流水只读」；unfreeze 端点在 §3 L42。停用时若仍有 `balance_frozen > 0`，解冻是否放行未写 |
| closure | — |

建议：disabled 同时拒绝 unfreeze（冻结资金随停用锁定），或明确「只读账户仍可解冻回 available」。S2 前一句即可。

### F-005 · 组合根权限基数过时（26→29；实测 27）

| 字段 | 值 |
|------|-----|
| level | recommended（low） |
| status | open |
| evidence | D-002 §6 L111：「admin 权限键 **26→29**」。`composition/composition_test.go` L466–471：admin **27** 权限 / **13** 导航（含 S-10 `users.mfa-reset`）。增量应为 **27→30**（+`wallet.read`/`write`/`adjust`）、导航 13→14（与 DefaultNavigationOrder 现长一致）。迁移 30→32 与 max 30 **自洽** |
| closure | — |

与 GOAL-016 A-004 F-004 同类。不阻断 S2；实施按 live snapshot 改断言。

### F-006 · 台账投影未跟：信息表仍写 I-00N open；S1 勾选含未执行 independent；goal-tree 仍 0/5

| 字段 | 值 |
|------|-----|
| level | recommended（low） |
| status | open |
| evidence | `03-audit.md` L17：「I-001~I-004 均 **open**（最晚需要阶段 S1）」vs `00-meta` L42–45 / `01-decision.md` L17–20 已 **verified**。`00-meta` L30 S1 检查点 `[x]` 正文含「self + **grok build independent**」，本条落盘前 independent 尚未存在；E-002 L21 与 `03-audit.md` L33 正确写「independent 待执行」。`00-meta` L8 `progress: 1/5` vs `goal-tree.md` L39 / L68 **0/5**（AGENTS §7：改 progress 必须同步树+表） |
| closure | — |

不得以 `progress: 1/5` 或 S1 勾选证明本条通过或放行 S2。建议 `/govern` 勘误：S1 勾选拆成「冻结稿 + self」与「independent」；goal-tree 在独立审闭合后再投影。本审计员不代改。

## 必改项汇总

1. **F-001（required · med）**：补 `amount_delta` 按 `entry_type` 的语义/符号/作用列/拒绝条件，以及快照链可执行重放规则。未闭合前**不可无条件放行 S2**。
2. **F-002（required · med）**：裁定幂等键唯一范围与「返回既有 / CONFLICT」分流，禁止跨账户按裸 key 返回流水。未闭合前**不可无条件放行 S2**。

无 high required。I-002/I-003/I-004 设计层证据充分；I-001 主体（分、恒等式、乐观锁、不可变流水、空引用）成立，但金额原语与幂等范围未冻满，**不能**把 I-001 `verified` 当作 S2 放行依据。

## 与既有意见的异同

- **A-003**（self · pass，0 findings）认为冻结稿证据充分、可进 independent。本意见同意：五段式/WithTx/EventXxx/0031-32/权限三键/协议独立/F-001·F-002 响应均成立。**不同意**「无缺口、可无条件进 S2」：self 未覆盖 `amount_delta` 互否与幂等键隔离（F-001/F-002）。
- **A-002** 立项 F-001/F-002（recommended）已由 D-002 §1/§3 响应，本 scope **不再开放**那两条。A-002 明确本条不能代替 S1 data 门禁独立审——成立。
- **A-001** 为立项 self，与本 scope 无冲突。
- 不与 A-001/A-002 verdict 冲突（立项 vs 方案是不同 scope）。A-003 与本条对 **S2 放行** 门禁互否（self 暗示稿已齐、本条 required 未闭）——属 P-004.2 门禁互否，须 `/govern` 响应 required，不得取乐观侧。

## 结论 + 建议给编排器/用户的下一步

**verdict: conditional**。S1 方案主体可实施，双层审计与迁移编号、权限拆分、协议独立对照诚实；**开放 F-001 / F-002 required**，不可无条件放行 S2。

建议 `/govern`：

1. 展示 F-001 / F-002；建议在 D-002 勘误补 apply 表 + 幂等范围（路径 = **fixed**）。不建议对金额原语走 residual。
2. F-003～F-006 随勘误或 S2 清单补记，不单独阻断。
3. 响应并闭合两条 required 后，再开 S2。**勿**把 `progress: 1/5`、S1 勾选或 I-001 `verified` 当作放行依据。
4. 本条闭合前 goal-tree 保持 0/5 与「独立审未放行」一致，优于抢先投 1/5。

**S2 放行条件（本意见）**：F-001 + F-002 按 P-003 三路径之一合法闭合（建议 fixed + D-002 补节）；I-002/I-003 维持设计层 verified；迁移仍从 31/32 起且实施前复核 max version。

## 声明

本意见不修改 `status` / `progress` / goal-tree / 方案正文。响应由 `/govern` 处理。保证等级 L0。
