---
id: A-004
goal: GOAL-020-wallet-auto-account
title: A-003 闭合复审 · 自动开户
date: 2026-08-16
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
scope: A-003 F-001~F-005 闭合核验
audit_type: finding-closure
verdict: conditional
status: recorded
parent: GOAL-020-wallet-auto-account
created: 2026-08-16
updated: 2026-08-16
version: 1.0.0
---

# A-004 · A-003 闭合复审（independent）

## 范围与区间

- **工作区**：`workspace-011-admin-functional-modules`；`root_goal` = `GOAL-001-admin-functional-modules`；`canonical_scope` 匹配；`shared_materials_catalog: none`（无资料引用可当事实或闭合依据）。
- **目标**：`GOAL-020-wallet-auto-account`（parent = Root；区内短 id 合法）。
- **scope**：A-003 F-001~F-005 **finding-closure**——只核对闭合证据，**不是**全目标重审。
- **对照原文**：`03-audit/A-003-s5-independent.md`（verdict conditional；required：F-001、F-002）。
- **保证等级**：L0。本意见不是第三方鉴证。
- **P-005**：`00-meta` / `01-decision` 中 I-001、I-002（required，最晚 S1）与 I-003（non-blocking，最晚 S1）均为 **verified**；本复审无到期未闭合 required 信息项。`03-audit.md` 信息表仍写 I-001~I-003 open（A-003 F-005 轻残留，见下）。
- **P-003**：合法闭合仅 `fixed` / `accepted-residual` / `user-overruled`。本轮未见用户书面 residual / overruled。`03-audit.md` 响应节仍为占位「（S5 独立审计后更新）」，无闭合声明——**不构成**「关闭声明不实」（故非 fail）。
- **本轮执行**：通读 A-003 与目标五件套/ledger；逐条对照代码与测试。复跑 `apps/api`：`./internal/modules/wallet/store/` + `./internal/handler/` **通过**（2026-08-16，exit 0）。未复跑 web；未做真实多连接压测。用户反馈同批（schema protocolVersion / breadcrumb / i18n）非 A-003 主题，仅略读，不纳入本 scope 判定。

## F-001 / F-002 闭合核验

### A-003 F-001 · 并发 UNIQUE 冲突重读后 `created` 仍为 true

| 子项 | A-003 要求 | 本轮证据 | 判定 |
|------|------------|----------|------|
| 冲突重读复位 | 冲突分支 `isNew = false` | `apps/api/internal/modules/wallet/store/repository.go` 264–268：`isUniqueViolation` 后显式 `isNew = false` 再重读 | **达成** |
| handler 审计条件 | 基于 `created` | `apps/api/internal/handler/wallet.go` 90–93（读）、132–135（调账）：均 `if created { recordWalletEvent(... EventWalletAccountCreate ...) }` | **达成** |
| 冲突重读测试 | 可复现 UNIQUE 冲突重读；断言 `created=false` 且不写第二笔 `wallet.account-create` | 全仓库 **不存在** `TestGetOrCreateUserAccountConcurrent`。现有 `repository_test.go` `TestGetOrCreateUserAccount`（282–318）仍为**顺序**二次调用：第二次走 SELECT 命中（243–246），**未进入** 264–278 冲突重读分支。注释仍写「concurrent-style」但未执行并发。store 包仅此一条 GetOrCreate 测试。 | **未达成** |

A-003 放行条件原文：「改代码 + 冲突重读测试可核对」。代码缺陷（败方仍报新建）已可静态核对；**测试腿未落地**。顺序二次调用不能代替 UNIQUE 冲突分支的运行时核对。

- **本条 A-004 F-001**：A-003 F-001 **未合法闭合**（`fixed` 的「可核对修正」在测试腿上不成立；亦无 residual / overruled 留痕）。
- **严重度**：med · **建议**：required · **影响门禁**：S5 关门 / data（账户创建审计幂等）· **状态**：open

### A-003 F-002 · by-owner 调账 Mutate 失败已开户但不写 account-create

| 子项 | A-003 要求 | 本轮证据 | 判定 |
|------|------------|----------|------|
| 审计时点 | 开户成功立即记 auto 审计（或等价：先校验 / 失败补偿） | `handler/wallet.go` 121–145：`GetOrCreateUserAccount` 之后、`Mutate` **之前**（127–135）按 `created` 写 `wallet.account-create`；Mutate 失败走 142–144 `writeWalletError` return，**不再**把创建审计绑在调账成功上 | **达成** |
| 失败路径测试 | 非法调账后账户已存在 + 审计条数 | `handler/wallet_auto_f2_test.go` `TestWalletByOwnerAdjustFailureStillAuditsCreate`：`amountDelta:0` + 空 memo → 4xx（20–22）；随后 `GET /api/wallet/by-owner/u300` = 200（24–28）；operations 中 `wallet.account-create` **≥ 1**（43–51） | **达成** |

本轮 handler 包测试通过，该用例随包执行。GET 本身也是 get-or-create，单独不能证明「是调账请求开的户」；但在「先 GetOrCreate 再 Mutate」实现下，原缺陷（调账已开户、无审计，随后 GET `created=false` 亦不补审计）会使 `creates` 为 0，测试会红。与代码时点一并，**可重复核对**。

- **闭合路径**：`fixed`（产物可核对）。编排器须在 `/govern` 响应节或 `02-execution` 留痕后，方可在台账上标 A-003 F-002 closed。
- **本复审不新开 F-002。**

## F-003 ~ F-005 核验（recommended）

### A-003 F-003 · 自动开户 id 未复用 `newID` 随机后缀

- **判定**：闭合（fixed）。
- **证据**：`repository.go` 252–258：`rand.Read` 12 字节 + `fmt.Sprintf("%016x%s", now.UnixMilli(), hex.EncodeToString(randBytes))`，与 `provider.go` `newID`（42–47）同构。`crypto/rand` / `encoding/hex` 已导入。`isUniqueViolation`（701–702）仍不区分 PK 与 owner UNIQUE；F-003 建议的「插入前走 newID」已落地，碰撞误判路径不再是同 `now` 确定性 id。

### A-003 F-004 · D-001 冻结路径与实现不一致

- **判定**：闭合（fixed）。
- **证据**：`01-decision/D-001-auto-account-plan.md` `version: 1.1.0`；§1 读/写端点已写 `GET/POST /api/wallet/by-owner/{ownerId}[/adjust]`，并注明 A-003 F-004 勘误（Go 1.22 ServeMux 与 `.../accounts/{id}/entries` 重叠）。与 `handler/wallet.go` 74、99 一致。

### A-003 F-005 · 验证与台账卫生

| 子项 | 判定 | 证据 |
|------|------|------|
| 403 循环含 by-owner | **闭合** | `handler/wallet_test.go` 115–116：`GET /api/wallet/by-owner/u1`、`POST /api/wallet/by-owner/u1/adjust` 已进 403 圈（与 91–92 的 401 圈对齐） |
| 错误码表含 `WALLET_USER_AUTO_ONLY` | **部分** | `errorcatalog.go` 150 双语 + `error_contract_test.go` 60 冻结集已列入；`wallet_auto_test.go` 117–118 断言 409 正文含该码。`wallet_test.go` `TestWalletErrorCodesCataloged`（320）**仍未**列入该码——轻残留 |
| `03-audit.md` 信息表 vs 权威 verified | **未同步** | 索引 13–20 行仍写 I-001~I-003 open；`00-meta` / `01-decision` 为 verified |
| goal-tree `progress` vs meta | **未同步** | `goal-tree.md` 40、69 行 GOAL-020 = **0/5**；`00-meta.md` = **4/5**。progress 非门禁、非闭合依据 |

F-005 整体：**部分闭合**。403 与冻结集已补；索引信息表与 goal-tree progress 仍为 A-003 已声明的轻残留，**不阻断**本复审 required 集合。本条不升格为 required。

### 本轮 recommended 残留（A-004）

#### F-002 · `TestWalletErrorCodesCataloged` 仍漏 `WALLET_USER_AUTO_ONLY`（F-005 子项）

- **严重度**：low · **建议**：recommended · **状态**：open
- **证据**：`handler/wallet_test.go` 320 码表未含 `WALLET_USER_AUTO_ONLY`；冻结集已覆盖（`error_contract_test.go` 60）。

#### F-003 · 台账卫生：信息表与 goal-tree progress 未同步（F-005 子项）

- **严重度**：low · **建议**：recommended · **状态**：open
- **证据**：`03-audit.md` 13–20；`goal-tree.md` 40、69；`00-meta.md` `progress: 4/5`。

## Findings（本条 A-004）

### F-001 · A-003 F-001 冲突重读测试腿未落地，required 未闭合

- **严重度**：med
- **建议**：required
- **影响门禁**：S5 关门 / data（账户创建审计幂等的可重复核对）
- **状态**：open
- **描述**：A-003 F-001 的代码缺陷已修（冲突分支 `isNew = false`；handler 仍按 `created` 写 `wallet.account-create`）。A-003 放行条件第二腿——「可复现的 UNIQUE 冲突重读测试」——在仓库中不存在。用户材料所称 `TestGetOrCreateUserAccountConcurrent`（8 并发、`created` 恰 1、单行）全库无匹配。现有 `TestGetOrCreateUserAccount` 与 A-003 当时指出的缺口相同：顺序二次调用不进入冲突重读。无该测试则冲突路径仅能静态阅读，不能运行时重复核对「败方 `created=false` 且不双写审计」。
- **证据**：
  - `store/repository.go` 264–268（代码已修）
  - `store/repository_test.go` 282–318（仍仅顺序）；全库 grep `TestGetOrCreateUserAccountConcurrent` / `GetOrCreateUserAccountConcurrent`：**0 命中**
  - `03-audit/A-003-s5-independent.md` 必改项表：F-001 放行条件 =「改代码 + 冲突重读测试可核对」
  - 无 `accepted-residual` / `user-overruled` 书面留痕
- **建议修正**：补一条真实并发（或等价强制 UNIQUE 冲突）测试：同 owner 并行 `GetOrCreateUserAccount`，断言恰好 1 个 `created=true`、账户单行、败方 `created=false`；可选再断言 handler 侧 `wallet.account-create` 恰 1。或由用户书面 `accepted-residual` / `user-overruled`（写清范围与复审触发）后闭合。

### F-002 · wallet 码表测试仍漏 `WALLET_USER_AUTO_ONLY`

- **严重度**：low · **建议**：recommended · **状态**：open
- **证据**：`handler/wallet_test.go` 320；对照 `error_contract_test.go` 60。

### F-003 · `03-audit` 信息表与 goal-tree progress 仍漂移

- **严重度**：low · **建议**：recommended · **状态**：open
- **证据**：`03-audit.md` 13–20；`goal-tree.md` 40、69 vs `00-meta.md` `progress: 4/5`。

## 必改项汇总

| ID | 级别 | 一句话 | 放行条件 |
|----|------|--------|----------|
| A-004 F-001（承接 A-003 F-001） | required / med | 冲突重读代码已修，但 UNIQUE 冲突重读测试仍缺，A-003 F-001 未合法闭合 | 补可核对并发/冲突测试；或用户书面 residual / overruled |

F-002（A-003）本复审视为产物已 `fixed`，待 `/govern` 留痕。A-004 F-002 / F-003 为 recommended，不单独阻断。

## 与既有意见的异同

| 条目 | 异同 |
|------|------|
| A-003 independent · S5 conditional | **同意**其 F-001/F-002 定性。本轮核验：F-002 代码+失败路径测试已齐；F-001 代码已齐、测试腿未齐。不把「代码已改」静默当成 F-001 整条 `fixed`。 |
| A-002 self · 无 finding | 仍部分不同意（同 A-003）：顺序测试不能覆盖冲突分支。本轮该缺口依旧。 |
| A-001 self · 立项 pass | 本复审不回溯立项。 |
| P-004 | 无 verdict 相反（本条仍 `conditional`）。无用户 residual/overruled。不得因 F-002 已修或 A-002 pass 放行 S5。 |

## 结论 + S5 关门放行条件

**verdict: conditional。** A-003 F-002（调账失败开户必须有 `account-create`）闭合证据充分可重复核对。A-003 F-003/F-004 已修。A-003 F-005 大部已修，台账/progress/`TestWalletErrorCodesCataloged` 为轻残留。A-003 F-001 的**代码**缺陷已修，但放行条件要求的冲突重读测试不存在——required 未按 P-003 三路径合法闭合。

**不得无条件将本目标标 `done`。** 关门最低条件：

1. A-003 F-001 / 本条 F-001 按 P-003 合法闭合：补上可核对的 UNIQUE 冲突重读测试（`fixed`），**或**用户书面 `accepted-residual` / `user-overruled`（范围 + 复审触发）。
2. A-003 F-002 在 `/govern` 响应节或 `02-execution` 写明 `fixed` 与证据路径（本复审已核验产物）。
3. 本意见不修改 `status` / `progress` / 方案正文 / goal-tree。响应与是否关门走 **/govern**。

建议编排器下一步：先补 F-001 并发/冲突测试（或提请用户裁决 residual），再响应 A-003 + 本 A-004；F-005 卫生可同批或关门后补。

## 声明

本意见 `source: independent`，保证等级 **L0**。不修改目标 `status`/`progress`，不关闭检查点，不改方案或代码。响应、finding 闭合与关门由用户通过 `/govern` 处理。
