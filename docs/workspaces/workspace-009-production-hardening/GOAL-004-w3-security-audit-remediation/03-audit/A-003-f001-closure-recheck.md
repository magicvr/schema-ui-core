---
id: A-003
goal: GOAL-004-w3-security-audit-remediation
title: A-002 F-001 闭合复核（批级 last-admin）
source: independent
auditor: grok-build · grok-4.5 · high · audit skill
date: 2026-08-11
verdict: pass
status: recorded
---

# A-003 · A-002 F-001 闭合复核（independent）

## 条目头

| 字段 | 值 |
|------|-----|
| **source** | independent |
| **auditor** | grok-build · grok-4.5 · high · 执行 `audit`（finding-closure） |
| **类型** | finding-closure |
| **scope** | A-002 **F-001** 修复复核：`DeleteUsersBatch` 批级 last-admin 判定是否堵住全 admin 清空；绕过面；与单条 `DeleteUser` 语义一致性 |
| **verdict** | **pass** |
| **工作区** | `workspace-009-production-hardening`（Root `GOAL-001-production-hardening`；canonical 已校验；`shared_materials_catalog: none`） |

## 范围与区间

- 代码：`apps/api/internal/modules/authsession/users_repository.go`（`DeleteUsersBatch`、`countAdminUsersExcludingBatch`、`DeleteUser`、`countAdminUsersExcluding`）
- 回归：`users_repository_test.go`（`TestDeleteUsersBatchRejectsRemovingAllAdmins`、`TestDeleteUsersBatchAtomicRollback`）；`handler/users_batch_test.go`（`TestUsersBatchDeleteRejectsRemovingAllAdminsHTTP` 等）
- 过程：`02-execution/E-002-f001-fix.md`；对照 `03-audit/A-002-w3-independent-cross.md` F-001
- **只读审计**：不改目标 status/progress/goal-tree；本地 `go test` 仅验证。

## 成果（有证据 · 本人核对）

| # | 核对项 | verdict | 证据摘要 |
|---|--------|---------|----------|
| a | 批级判定堵住全 admin 清空 | **pass** | `adminsInBatch > 0` 时 `countAdminUsersExcludingBatch`：`admin AND user_id NOT IN (keys)`；`other == 0` → `ErrLastAdmin` 整批回滚（`users_repository.go:305-312`、`328-343`） |
| b | 新绕过（空批 / NOT IN 空集 / actor 在批内） | **pass**（无 required） | 见下「绕过面」 |
| c | 与单条 `DeleteUser` 语义一致 | **pass** | 单 id：`countAdminUsersExcluding(id)` ≡ 批集大小 1 的批级排除；self / last-admin / not-found 守卫序一致 |

**本地复跑（2026-08-11）**

- `go test ./internal/modules/authsession/ -run 'TestDeleteUsersBatch|TestUsersLastAdmin|TestDeleteUserSerializes'` → **PASS**
- `go test ./internal/handler/ -run TestUsersBatchDelete` → **PASS**（含 `TestUsersBatchDeleteRejectsRemovingAllAdminsHTTP`）

## 逐项复核

### a · 批级 last-admin 是否彻底堵住全 admin 清空 — **pass**

A-002 F-001 根因：per-id `countAdminUsersExcluding(tx, id)` 在「批内含全部 admin」时，每个 id 仍看到其他 admin 在库 → 守卫全过 → 统一删除后 0 admin。

现实现：

1. 去重 `keys := dedupeKeys(ids)`（`266-267`）
2. 守卫循环：exists → self → 累计 `adminsInBatch`（不再 per-id last-admin）
3. **批级**：`adminsInBatch > 0` 且 `countAdminUsersExcludingBatch(tx, keys) == 0` → `ErrLastAdmin`（fail closed，事务回滚）
4. 仅守卫全过后才删 refresh + users

不变量：若批会删除任意 admin，则删除后剩余 admin 数 =「不在选择集内的 admin 数」；该数为 0 即拒绝。覆盖：

| 场景 | 期望 | 覆盖 |
|------|------|------|
| 批 = 全部 admin（≥2） | `ErrLastAdmin`，零删除 | 仓库 + HTTP 测 |
| 批 = 全部 admin + viewer | 同上 + viewer 回滚 | 仓库测 |
| 存活 admin3 后同批 | 允许删除 | 仓库 + HTTP 反例 |
| 单 last admin 在批内 | `ErrLastAdmin` | 既有 atomic 测 |
| 双 admin 时删其一 | 允许 | 既有 atomic 测 |

### b · 绕过面 — **pass**（无新 required）

| 候选绕过 | 分析 | 结论 |
|----------|------|------|
| **去重后空批** | `len(keys)==0` → `return 0, nil`，不进事务、不调 `NOT IN`（`267-269`）。HTTP 另拒 empty selection（`TestUsersBatchDeleteFailClosed`） | **无安全绕过**（空操作） |
| **`NOT IN` 空集** | `countAdminUsersExcludingBatch` 仅在 `adminsInBatch > 0` 时调用；`adminsInBatch>0` ⇒ `len(keys)≥1`，占位符非空。参数绑定，非字符串拼接 id | **不可达** |
| **actor 在批内** | 守卫序 exists → **self** → admin 统计；`id == actorID` → `ErrSelfOperation`（`282-284`），先于 last-admin。与单条一致 | **无绕过**；非 admin actor 批含全 admin → 仍 `ErrLastAdmin`（HTTP um 用例） |
| **重复 id** | `dedupeKeys` 后计一次；合法 | 既有 dedupe 测 |
| **并发两批各删一名 admin** | 与单条并发同类；`withTx` + SQLite 写串行；`TestDeleteUserSerializesLastAdminCheck` 覆盖该类不变量 | **非本修复引入**；不升 required |

### c · 与 `DeleteUser` 语义 — **pass**

| 守卫 | `DeleteUser` | `DeleteUsersBatch` |
|------|--------------|-------------------|
| not-found | `ErrNotFound` | 任 id 不存在 → `ErrNotFound` 整批回滚 |
| self | `id == actorID` | 任一 key == actor → `ErrSelfOperation` |
| last-admin | 目标是 admin 且 `countAdminUsersExcluding(tx, id)==0` | 批内含 admin 且 `countAdminUsersExcludingBatch(tx, keys)==0` |
| 事务 | 单事务 | 单事务整批 |

单 id 批删与单删 last-admin 等价（排除集大小 1）。多 id 时批级集合语义是单删的正确推广，修复了 A-002 指出的集体绕过。

## 对照成功标准（F-001 建议闭合）

| A-002 建议 | 状态 |
|------------|------|
| 批级「删除后剩余 admin ≥ 1」 | **满足**（`NOT IN` 计数） |
| `TestDeleteUsersBatchRejectsRemovingAllAdmins` | **存在且绿** |
| HTTP 层等价用例 | **`TestUsersBatchDeleteRejectsRemovingAllAdminsHTTP` 绿**（含 admin3 反例） |

## Findings

本 scope（F-001 闭合）**无新 required / recommended finding**。

### 相对 A-002 的 finding 状态建议（供 /govern 闭合，本意见不改 status）

| ID | 原严重度 | 本复核 | 建议 |
|----|----------|--------|------|
| **A-002 F-001** | required · high | 修复 + 回归 + 绕过面核对成立 | **可标 `fixed`**（闭合证据：E-002 + 本 A-003 + 上述测试绿） |
| A-002 F-002 | recommended · low | 未重审范围 | 仍 recommended（委托删 admin / 改名） |
| A-002 F-003 | recommended · low | 未重审范围 | 仍 recommended（限流 private peer） |

## 必改项汇总

- **开放 required（本 scope）**：**0**
- 无新必改项。F-001 闭合手续由 **`/govern`** 写入响应/台账。

## 与既有意见的异同

- A-002 **conditional**，开放 required=1（F-001）。
- 本 A-003 为 **finding-closure** 窄 scope：确认 F-001 实现与测试充分，**不**重开 W3 全八项。
- A-001 self pass 不替代本独立闭合复核。

## 结论 + 建议下一步

- **verdict: pass**（F-001 闭合复核）
- **开放 required（本 scope）: 0**
- 建议用户：`/govern` 响应 A-002 F-001 → **fixed**（引用 E-002 + A-003 + 测试路径），再评估 GOAL-004 是否可进入关门路径（仍须处理其余 recommended 是否接受 residual，及整目标关门标准）。

## 声明

本意见 `source: independent`，**不**修改 status/progress/goal-tree；响应与 finding 合法闭合由 **`/govern`** 处理。
