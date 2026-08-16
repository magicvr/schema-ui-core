---
id: A-005
goal: GOAL-020-wallet-auto-account
title: A-004 闭合复审 · 自动开户
date: 2026-08-16
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
scope: A-004 F-001 闭合 + F-003~F-005 核验
audit_type: finding-closure
verdict: pass
status: recorded
parent: GOAL-020-wallet-auto-account
created: 2026-08-16
updated: 2026-08-16
version: 1.0.0
---

# A-005 · A-004 闭合复审（independent）

## 范围与区间

- **工作区**：`workspace-011-admin-functional-modules`；`root_goal` = `GOAL-001-admin-functional-modules`；`canonical_scope` 匹配；`shared_materials_catalog: none`（无资料引用可当事实或闭合依据）。
- **目标**：`GOAL-020-wallet-auto-account`（parent = Root；区内短 id 合法）。
- **scope**：A-004 F-001 **finding-closure**（A-003 F-001 required 闭合核验）+ A-003 F-003~F-005 / A-004 F-002~F-003 核验。顺带复核 A-003 F-002（A-004 已判产物 `fixed`）。**不是**全目标重审。
- **对照原文**：`03-audit/A-003-s5-independent.md`（required：F-001、F-002）；`03-audit/A-004-s5-reaudit.md`（verdict conditional；A-004 F-001 required 仍 open：缺并发测试）。
- **保证等级**：L0。本意见不是第三方鉴证。
- **P-005**：`00-meta` / `01-decision` / 本目标 `03-audit.md` 信息表中 I-001、I-002（required，最晚 S1）与 I-003（non-blocking，最晚 S1）均为 **verified**；本复审无到期未闭合 required 信息项。
- **P-003**：合法闭合仅 `fixed` / `accepted-residual` / `user-overruled`。本轮未见用户书面 residual / overruled。`03-audit.md` 响应节仍为占位「（S5 独立审计后更新）」——**不构成**「关闭声明不实」（故非 fail）；台账标 closed 仍须 `/govern` 留痕。
- **本轮执行**（2026-08-16）：
  - `apps/api`：`go test ./internal/modules/wallet/store/ -run TestGetOrCreateUserAccountConcurrent -count=1` → **ok**（~1.054s，exit 0）。
  - `apps/api`：`go test ./internal/handler/ -run TestWalletByOwnerAdjustFailureStillAuditsCreate|TestWalletErrorCodesCataloged|TestWalletRoutesGates -count=1` → **ok**（~3.717s，exit 0）。
  - 未复跑 web；未做多连接 / 真实压测。

## 逐条闭合核验

### A-004 F-001（承接 A-003 F-001）· 冲突重读测试腿 + `isNew` 复位

| 子项 | A-004 / 本轮要求 | 本轮证据 | 判定 |
|------|------------------|----------|------|
| 冲突重读复位 | 冲突分支 `isNew = false` | `apps/api/internal/modules/wallet/store/repository.go` 264–268：`isUniqueViolation` 后显式 `isNew = false` 再重读 | **达成** |
| handler 审计条件 | 基于 `created` | `apps/api/internal/handler/wallet.go` 90–93（读）、132–135（调账）：均 `if created { recordWalletEvent(... EventWalletAccountCreate ...) }` | **达成** |
| 并发测试存在 | `TestGetOrCreateUserAccountConcurrent` | `apps/api/internal/modules/wallet/store/concurrent_test.go`（包 `store_test`） | **达成** |
| 8 goroutine | 同 owner 并行 | `const n = 8`；`sync.WaitGroup` + `go func` 调 `GetOrCreateUserAccount("u-concurrent", now())`（12–33 行） | **达成** |
| `created` 恰 1 | 败方不得报新建 | 46–48 行：`creates != 1` → Fatal | **达成** |
| 账户单行 | 同一 id + 表内一行 | 49–51 行 distinct ids = 1；53–60 行 `ListAccounts` `total != 1 \|\| len != 1` → Fatal | **达成** |
| 测试执行 | 可复跑 | 本轮 `go test … -run TestGetOrCreateUserAccountConcurrent -count=1` **通过** | **达成** |
| UNIQUE 重读分支被执行 | A-003/A-004「冲突重读测试可核对」的运行时腿 | 平台 `apps/api/internal/store/store.go` 49：`db.SetMaxOpenConns(1)`。`WithTx` 在唯一连接上串行 Begin/Commit，后到 goroutine 走 SELECT 命中（`repository.go` 243–246），**不进入** 264–278。本测试在去掉 `isNew = false` 时仍会绿。 | **未达成（残留）** |

A-004 放行条件原文：「补可核对并发/冲突测试」。用户本轮指定核验项（测试真实存在且执行、`created` 恰 1、单行、`isNew=false` 复位）均有证据。A-004 建议修正把「真实并发」与「等价强制 UNIQUE 冲突」并列；交付的是前者。因连接池上限，该「真实并发」在 SQLite 事务边界上等价于顺序 SELECT 命中，**不能**运行时核对 264–278。

- **A-003 F-001 / A-004 F-001（required）**：代码缺陷 + A-004 点名缺失的测试产物，本复审视为产物侧可核对（`fixed` 证据充分）。编排器须在 `/govern` 响应节或 `02-execution` 留痕后，方可在台账上标 closed。
- **UNIQUE 分支运行时未覆盖**：不把 A-004 F-001 required 整条重新打开（避免在已交付指定测试后平移门槛）；记为本条 **recommended** 残留（见下 F-001）。

### A-003 F-002 · by-owner 调账 Mutate 失败已开户但不写 account-create

| 子项 | A-003 要求 | 本轮证据 | 判定 |
|------|------------|----------|------|
| 审计时点 | 开户成功立即记 auto 审计 | `handler/wallet.go` 121–145：`GetOrCreateUserAccount` 之后、`Mutate` **之前**（127–135）按 `created` 写 `wallet.account-create` | **达成**（同 A-004） |
| 失败路径测试 | 非法调账后账户已存在 + 审计条数 | `handler/wallet_auto_f2_test.go` `TestWalletByOwnerAdjustFailureStillAuditsCreate`：`amountDelta:0` + 空 memo → 4xx；随后 GET by-owner = 200；`wallet.account-create` ≥ 1 | **达成** |
| 测试执行 | 可复跑 | 本轮 handler 包该用例随 `-run` **通过** | **达成** |

本复审不新开 F-002。台账标 closed 仍待 `/govern` 留痕（同 A-004）。

### A-003 F-003 · 自动开户 id 未复用 `newID` 随机后缀

- **判定**：维持闭合（fixed）。
- **证据**：`repository.go` 252–258：`crypto/rand` 12 字节 + `fmt.Sprintf("%016x%s", now.UnixMilli(), hex.EncodeToString(randBytes))`，与 `provider.go` `newID`（42–47）同构。`crypto/rand` / `encoding/hex` 已导入。

### A-003 F-004 · D-001 冻结路径与实现不一致

- **判定**：维持闭合（fixed）。
- **证据**：`01-decision/D-001-auto-account-plan.md` `version: 1.1.0`；§1 读/写端点为 `GET/POST /api/wallet/by-owner/{ownerId}[/adjust]`，并注明 A-003 F-004 勘误。与 `handler/wallet.go` 74、99 一致。

### A-003 F-005 / A-004 F-002 / A-004 F-003 · 验证与台账卫生

| 子项 | A-004 当时 | 本轮 | 判定 |
|------|------------|------|------|
| 403 循环含 by-owner | 已闭合 | `handler/wallet_test.go` 115–116：`GET /api/wallet/by-owner/u1`、`POST /api/wallet/by-owner/u1/adjust`；本轮 `TestWalletRoutesGates` 通过 | **维持闭合** |
| `TestWalletErrorCodesCataloged` 含 `WALLET_USER_AUTO_ONLY` | A-004 F-002 未列入（320 行） | `wallet_test.go` 320 码表**已含** `WALLET_USER_AUTO_ONLY`；本轮该测试通过 | **闭合**（A-004 F-002） |
| `03-audit.md` 信息表 vs 权威 verified | 当时仍写 I-001~I-003 open | 索引 13–20 行现写 I-001/I-002 required + I-003 non-blocking 均 **verified**；与 `00-meta` / `01-decision` 一致 | **闭合** |
| goal-tree `progress` vs meta | 当时 0/5 vs 4/5 | `goal-tree.md` 40、69 行 GOAL-020 = **4/5**；`00-meta.md` = **4/5**。progress 非门禁、非闭合依据 | **闭合** |

A-003 F-005 与 A-004 两条 recommended 残留（码表测试、信息表/progress 漂移）本轮均有证据闭合。

## Findings（本条 A-005）

### F-001 · 并发测试因 `MaxOpenConns(1)` 未打到 UNIQUE 冲突重读

- **严重度**：low
- **建议**：recommended
- **影响门禁**：S3 证据完整度（不阻断本复审 required 集合；不单独阻断 S5）
- **状态**：open
- **描述**：`TestGetOrCreateUserAccountConcurrent` 已落地、断言 `created` 恰 1 + 单行、本轮执行通过；冲突分支 `isNew = false` 可静态核对。但 `store.Open` 固定 `SetMaxOpenConns(1)`，8 个 goroutine 在唯一连接上串行事务：胜者 INSERT，其余走 SELECT 命中，**不会**触发 `isUniqueViolation` / 264–278。因此「败方在 UNIQUE 冲突后 `created=false`」仍只能静态阅读，不能运行时重复核对。这与 A-004 批评顺序测试「未进入冲突重读」属同一机制缺口，只是外包了一层 goroutine。
- **证据**：
  - `apps/api/internal/modules/wallet/store/concurrent_test.go`（测试形态与执行：本轮 exit 0）
  - `apps/api/internal/store/store.go` 49（`SetMaxOpenConns(1)`）
  - `apps/api/internal/modules/wallet/store/repository.go` 237–278（SELECT 命中 vs 冲突重读）
- **建议修正**：补一条**强制** UNIQUE 冲突用例（例如可注入的 TxRunner / 双连接 shared-cache / 在同一事务可见窗口插入冲突行），断言重读分支 `created=false`；或用户书面 `accepted-residual`（写清「仅静态核对冲突分支」与复审触发）。不要求本条升格 required。

本复审**不**新开 required finding。

## 必改项汇总

| ID | 级别 | 一句话 | 放行条件 |
|----|------|--------|----------|
| — | — | 本条 scope 内无未闭合 required | — |

A-003 F-001 / A-004 F-001、A-003 F-002：产物侧视为可核对 `fixed`，待 `/govern` 响应节或 `02-execution` 留痕。A-005 F-001 为 recommended，不单独阻断。

## 与既有意见的异同

| 条目 | 异同 |
|------|------|
| A-004 independent · finding-closure conditional | **同意**其当时判定：测试腿缺失则 required 未闭合。本轮该测试已存在且执行通过，故不再维持 A-004 F-001 required。**部分收紧**：指出 `MaxOpenConns(1)` 使该测试仍不进入 264–278，记 recommended，不平移 required 门槛。 |
| A-003 independent · S5 conditional | 同意其 F-001/F-002 定性。本轮核验：两条 required 的代码 + 指定测试均齐。不把「响应节仍为占位」写成关闭声明不实。 |
| A-002 self · 无 finding | 对「顺序测试 = 冲突重读已覆盖」仍不同意；并发测试已补，剩余仅为 UNIQUE 分支运行时覆盖（recommended）。 |
| A-001 self · 立项 pass | 本复审不回溯立项。 |
| P-004 | 无 verdict 相反。无用户 residual/overruled。本条 `pass` 只覆盖本 finding-closure scope，**不**自行改 `status` / 关门。 |

## 结论 + S5 关门放行条件

**verdict: pass。** A-004 点名的 required 缺口（A-003 F-001 并发测试腿）已有可重复执行的产物：`TestGetOrCreateUserAccountConcurrent` 本轮通过，冲突分支 `isNew = false` 与 handler 按 `created` 审计可静态核对。A-003 F-002 审计前置 + 失败路径测试本轮复跑通过。A-003 F-003/F-004 维持闭合。A-003 F-005 与 A-004 recommended（码表、信息表、goal-tree 4/5）本轮均已同步。scope 内无未关闭 high/med required，无到期 required 信息项。

UNIQUE 冲突重读的**运行时**覆盖仍弱（A-005 F-001 recommended），不构成「关闭声明不实」，也不单独阻断本复审。

**S5 关门**不由本条直接放行。编排器最低条件：

1. 在本目标 `03-audit` 响应节或 `02-execution` 写明 A-003 F-001 / A-004 F-001 与 A-003 F-002 的 `fixed` 与证据路径（本复审已核验产物；口头不算）。
2. 无新的未闭合 required；A-005 F-001 可同批修、列为残余，或关门后补。
3. 本意见不修改 `status` / `progress` / 方案正文 / goal-tree。响应与是否关门走 **/govern**。

## 声明

本意见 `source: independent`，保证等级 **L0**。不修改目标 `status`/`progress`，不关闭检查点，不改方案或代码。响应、finding 闭合与关门由用户通过 `/govern` 处理。
