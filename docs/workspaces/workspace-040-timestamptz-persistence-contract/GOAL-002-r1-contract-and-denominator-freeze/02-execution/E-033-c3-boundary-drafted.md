---
id: E-033-c3-boundary-drafted
doc_type: goal-execution-entry
status: recorded
date: 2026-09-20
parent: GOAL-002-r1-contract-and-denominator-freeze
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-033 · A-038 响应与 C3 备份/回滚边界首次落盘

## 事实

1. **独立审计 A-038 完成并落盘**（`source: independent`，`auditor: grok-build (grok-4.6 · reasoning high)`）：verdict `conditional`，**开放 required 由 4 降至 3**（F-I-002 / F-I-004 / F-I-005）。
   - **closed**：**F-I-026**（PG 附件已补 v78 两列；v73–v87 的 **90 列全部被点名，无缺失**）。
   - **fixed**：F-I-002.1 第 1 项（`dict_entries` 可粘贴 10 列 CREATE，与 live cid 一致）与第 3 项（`#72/#73` 单路径，文件内无其它双分支）。
   - **无新 finding**；A-034 接受的形式均**保持**。
   - **补充**：`Descriptors()` 的 `Version` 字段序与 live cid 同序 → 列序交叉核对现有**三重**同序证据（cid / `sqlite_master` 折入文本 / `Descriptors()` 版本序），唯一会误导的是模块文件**物理行号序**。
   - **F-I-002 整条仍 open**，剩余实质项为「用例仍是 ID；非法/越界可执行测试未发生」。
2. **A-038 的证据窗口处置正确**：它明确声明 `attachments/r1-c3-backup-recovery-boundary-v1.0-fc.md` 当时为**未跟踪文件、不在 `d1fdb4cc`、不在其 scope**，不作为 F-I-004 证据。该文件在**本轮**才提交，**尚未被任何 independent 复审**。
3. **C3 首次落盘**：新增 `attachments/r1-c3-backup-recovery-boundary-v1.0-fc.md`（`status: freeze-candidate`），逐条对位 F-I-004 的五项关闭要求：
   - §2 三类产物唯一区分（A pre-conversion rollback / B post-conversion RecoveryPoint / C batch boundary）+ 三条硬规则（A/C 永不作为目标形状校验输入；B 只能由转换后目标库产生；三者身份必须可区分）；
   - §3 PG 双 token（`<rollback-artifact>` vs `<recovery-artifact>`），修正 runbook L34–L39 把预转换 `<artifact>` 当目标形状校验输入的矛盾；
   - §4 包路径（kernel 仅接口 + `internal/backup` 编排/provider/verify）与三个 before/after 调用点（**禁止**在 `applyMigration` 事务内调用）；
   - §5 restore-to-new-db harness 规格（两侧同构 + 6 项断言 + PG 版本兼容显式记录，不可静默 skip 通过）；
   - §6 旧 dump 必须**被拒**（`TestLegacyArtifactMustFail`），且该失败须作为正向证据收录。
4. **C3 现状事实（实测）**：`apps/` 中 `CreateRecoveryPoint` / `BackupService` / `RecoveryPoint` **0 匹配**；`kernel.Store`（`kernel/store.go:30-47`）无 Backup 接口；`internal/temporal` 不存在；SQLite 现有 `snapshotBeforePending`（`migrate.go:279-300`）是 per-migration rollback 点。→ **C3 全部为设计，无任何实现**。

## 证据

- A-038 全文：`03-audit/A-038-r1-independent-e032-a036-response-fi026.md`。
- A-039 响应：`03-audit/A-039-r1-self-response-to-a038.md`。
- 新增附件：`attachments/r1-c3-backup-recovery-boundary-v1.0-fc.md`。
- 代码对位：`kernel/store.go:30-47`；`internal/store/migrate.go:81-103,108-132,279-300,345-367`；`internal/store/postgres.go:154-171`。
- 本轮 `apps/**` **未修改**。

## 状态评估

- **开放 required = 3**（F-I-002、F-I-004、F-I-005）；**本条未闭合任何 required**。
- **F-I-002 的结构性观察（已在 A-039 §3 记录，下一轮需向用户取 P-004 口径）**：其唯一实质剩余是「可执行测试」——若允许在**一次性验证库**上按合同 DDL 跑负值/越界用例，可在 R1 内闭合且不触碰生产 schema；若必须绑定真实迁移代码，则属 R2，将构成「F-I-002 挡 C2 冻结、C2 冻结又挡跑测试」的死锁。
- 下一步：`/audit` 复审 C3 边界文件（从未复审）+ 就 F-I-002 测试口径取 P-004。
