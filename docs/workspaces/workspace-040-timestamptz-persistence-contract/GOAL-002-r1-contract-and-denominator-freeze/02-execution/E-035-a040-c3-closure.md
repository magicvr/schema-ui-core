---
id: E-035-a040-c3-closure
doc_type: goal-execution-entry
status: recorded
date: 2026-09-20
parent: GOAL-002-r1-contract-and-denominator-freeze
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-035 · A-040 响应：C3 七项收口与 F-I-025/F-I-027 卫生

## 事实

1. **A-040 判定**：**F-I-004 收窄、不能闭合**；开放 required 仍为 **3**（F-I-002 / F-I-004 / F-I-005）；**无新 required**；新增 recommended **F-I-027**。
   - 已接受的子项：新文件内部三类区分与双 token；`kernel` + `internal/backup` 包路径；SQLite 行号与「禁止事务内调用」；harness 规格（含文件名）作为 R1 设计足够；§1 四条现状事实成立；§7 与 `D-019` 自洽。
   - §G 给出 **7 项剩余**；另在 §D 末段指正 §5 断言 4 的「负值非法」放错位置（不升格 required）。
2. **本轮逐项收口**：
   - **§G1**：`r1-c3-backup-restore-runbook-v0.1.md` 标 `status: superseded`（v0.2.0）；正文保留为历史并就地标注**两处错误点**（per-migration snapshot 是批次边界产物；预转换 `<artifact>` 被当目标形状 restore 输入必失败）。
   - **§G2 / F-I-027**：Port draft **L60** 删「`snapshotBeforePending` or」并列并**明令**其不是 RecoveryPoint 源；**L71** 标已解决；文件尾新增权威范围声明。conversion contract §0「尚未落盘」改为已落盘且为 C3 唯一权威。
   - **§G3**：C3 边界文件 §4.2 拆为 SQLite 侧与 **PG 侧**两张调用点表（`postgres.go:120-127` / `:154-173`），并**标注现状缺口**——`applyPendingPG` 既无 snapshot 步骤也无可对称的 `verifyIntegrity()`。
   - **§G4**：§2 三行精确化 A/B/C 时机；**纠正** C 的形状表述为「**混合形状**，随批次推进」；新增「A 与 C 必须分开」。
   - **§G5**：新增 §4.3 重试要求四条（可检测状态 → 显式动作 → 形状判定后补创 B 或走回滚 → 有界重试；禁止把「无 B」当门禁已满足）。
   - **§G6 + §D末段**：新增 §5.1 **6 类错误分类表**；`TestLegacyArtifactMustFail` 必须断言分类属于 `{TimeContractMismatch, TemporalColumnSetIncomplete}` 且**不属于** `{ArtifactNotFound, ArtifactUnreadable, ToolFailure}`；§5 断言 4 删除 B 上的「负值非法」。
   - **§G7**：新增 §3.1 C→B **机械身份**——`contract_shape` 由对 restore 目标的**实测**得出，判定**实测驱动而非命名驱动**。
   - **§D 编号限定**：§4.1 的 `D-010`/`D-007` 改为 **Root** `D-007` + **Root** `D-010`（child `D-006` + child `D-009`）。
   - **F-I-025**：rebuild 附件 §0 的「20 张时间列表」改为「20 张**带时间列**的表（分母 90 列）+ 3 张**无时间列**联接表 + 3 张跨 descriptor 子表」。
3. **本轮 `apps/**` 未修改**（上轮新增的 `apps/api/internal/w040contracttest/` 仅测试包，已在上轮提交）。

## 证据

- A-040 全文：`03-audit/A-040-r1-independent-e033-c3-boundary-fi004.md`。
- A-041 响应：`03-audit/A-041-r1-self-response-to-a040.md`。
- 附件变更：`attachments/r1-c3-backup-restore-runbook-v0.1.md`（superseded）、`attachments/r1-backup-port-contract-draft-v0.1.md`（L60/L71 + 权威声明）、`attachments/r1-c3-backup-recovery-boundary-v1.0-fc.md`（§2/§3.1/§4.1–4.3/§5.1/§8）、`attachments/r1-c2-per-column-conversion-contract-v1.0-fc.md`（§0）、`attachments/r1-c2-per-table-rebuild-ddl-v1.0-fc.md`（§0 计数）。
- 代码对位：`internal/store/migrate.go:52-53,81-103,108-132,345-367`；`internal/store/postgres.go:94-95,120-127,154-173`；`kernel/store.go:30-47`。

## 状态评估

- **开放 required = 3**；**本条未闭合任何 required**（7 项收口的接受与否待 independent 复审）。
- recommended **F-I-025 / F-I-027** 的对应修改已落盘，**是否闭合待复审**。
- 下一步：`/audit` 复审本轮收口（重点 §4.2 PG 调用点与 §4.3 重试的可执行性）+ F-I-002 测试文件。
