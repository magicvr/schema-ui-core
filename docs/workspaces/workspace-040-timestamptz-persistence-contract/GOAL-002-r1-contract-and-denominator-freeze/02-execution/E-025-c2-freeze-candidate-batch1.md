---
id: E-025-c2-freeze-candidate-batch1
doc_type: goal-execution-entry
status: recorded
date: 2026-09-20
parent: GOAL-002-r1-contract-and-denominator-freeze
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-025 · C2 冻结候选第一批

## 事实

1. 用户 2026-09-20 经 P-004 裁决：(a) **C2/C3 冻结口径 = 可直接执行的逐列 DDL/编解码合同 + 已记录 checksum**；(b) **v74 表范围 = `system_data_reconcile`（不含 `schema_migrations`）**；(c) 本轮交付 C2 冻结候选第一批。
2. 新增三份 `freeze-candidate` 附件（**设计合同材料，非实施事实**）：
   - `attachments/r1-c2-per-column-conversion-contract-v1.0-fc.md`：90 行逐列合同（`#N` ↔ inventory v0.3 一一对应）；冻结表达式仅三式 E1/E2/E3；自检 `sec 80 + ms 10` = `NN 71 + N 14 + D0 5` = 90，owner 分配 v73 3 / v74 31 / v75 3 / v76 5 / v77 4 / v78 2 / v79 4 / v80 4 / v81 2 / v82 2 / v83 5 / v84 1 / v85 11 / v86 7 / v87 6。经机械校验：90 行、行号 1–90 无重无缺。
   - `attachments/r1-c2-predicate-exact-sql-v1.0-fc.md`：**一张** exact old/new SQL 单表（收口 predicate-index-matrix 与 readwrite-predicate-spec 两源），含 `m0–m5` 迁移顺序、不可逆点清单、`P-*`/`T-*` 测试 ID；§5 处置清单并集经机械校验覆盖全部 90 列。
   - `attachments/r1-c2-descriptor-ledger-v1.0-fc.md`：v73–v87 的 `Name` / `transform_id`（按仓库既有 `"0073:vp040-temporal-core-persistence:v1"` 约定）、唯一表范围、`MigrationChecksum` 计算输入与 `m0–m5` 语句序位、append-only 边界、金额列断言拆分要求。
3. 三份附件均**明确声明未闭合项**：已迁移 checksum 值、可执行测试改写、双方言 checksum 约定二选一，均需 R2 落码后闭合。

## 证据

- 90 行与分配自检：本次落盘时以脚本对附件正文机械计数（`^\|\s*\d+\s*\|` 命中 90 行，行号 1–90 无重复、无缺失；`sec`/`ms`、`E1`/`E2`/`E3` 计数与声明一致）。
- 谓词覆盖自检：§5 处置清单六类并集 = 90 列，无遗漏。
- `old` 列 SQL 与行号核对来源：`apps/api/modules/authsession/accounts_lock_source.go:67-128`、`authsession/users_repository.go:232-234`、`authsession/invites.go:200,206,297`、`authsession/service_credentials.go:165,185,198`、`authsession/notifications_repository.go:106,150,156,164,220,240`、`wallet/voucher/service.go:330,340,347,391`、`jobs/migration/migration.go:36-43,45-47,75-82`、`internal/jobs/repository.go:89,244,256,266,271,282,288,300-301`、`scheduledtasks/store/repository.go:289-290,343-344`、`recyclebin/migration/migration.go:30,48`、`recyclebin/store/repository.go:100,215,230`、`logincaptcha/store/repository.go:40,69`、`mfa/store/repository.go:273`、`authsession/recovery.go:287`、`authsession/email_identity.go:259`、`kernel/persistence.go:14-17`。
- `apps/` 代码本条目**未修改**（`git status` 无 `apps/**` 变更）。

## 状态评估

- 本条目只记录**设计合同材料的产出**。A-029 的 F-I-002～F-I-006 **仍未闭合**；C2/C3 **未冻结**；R2 **未放行**。
- 上述三份附件是否构成合法闭合，须由 independent 复审判定（不得由本条目自证）。
