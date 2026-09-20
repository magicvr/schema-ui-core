---
id: A-041-r1-self-response-to-a040
doc_type: goal-audit-entry
source: self
auditor: /govern
date: 2026-09-20
scope: GOAL-002-r1-contract-and-denominator-freeze · response to A-040 / C3 seven-item closure, F-I-027 hygiene, F-I-025 count wording
verdict: conditional
open_required: 3
status: recorded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-001-timestamptz-persistence-contract
version: 0.1.0
---

# A-041 · R1 self response to A-040

- **source**：self（编排器响应，**不**冒充 independent）
- **verdict**：conditional
- **开放 required**：3（F-I-002、F-I-004、F-I-005）

## 1. 接受 A-040 的判定

| 项 | A-040 verdict | 本响应动作 |
|----|---------------|------------|
| F-I-004 | **收窄、不能闭合**；给出 §G 的 7 项剩余 | 本轮**逐项收口**（§2） |
| 已接受子项 | 新文件内部三类区分与双 token；`kernel`+`internal/backup` 包路径；SQLite 行号与「禁止事务内调用」；harness 规格（含文件名）；§1 四条现状事实成立；§7 与 `D-019` 自洽 | 保持不动 |
| §D 末段（断言 4 的「负值非法」不贴切） | 不升格 required，但指正成立 | 已改（§2.6） |
| **F-I-027**（recommended） | 新增：runbook/Port/conversion §0 三源未收口 + 无限定 `D-010`/`D-007` | 已修（§2.2、§2.8） |
| **F-I-025**（recommended） | 仍 open：「20 张时间列表」表述 | 已改（§2.9） |
| 新 required | **无** | — |

## 2. 逐项收口（A-040 §G 1–7 + F-I-027 + F-I-025）

### 2.1 §G1 · 旧 runbook 的 `<artifact>` 矛盾 —— 已标 `superseded`

`r1-c3-backup-restore-runbook-v0.1.md`：`status: proposed → superseded`（v0.2.0）。正文**保留为历史**（按审计记录不得改写原则），但：

- 顶部新增 superseded 声明，移交唯一权威给 C3 边界文件，并说明移交原因；
- 就地标注**两处错误点**：① 第 1 步的 per-migration snapshot 是「批次边界产物」而非常量旧合同产物；② 第 1 步 `<artifact>` 是预转换 dump 却被第 5 步 restore 后由第 6 步断言 `timestamptz(6)`，**必然失败**；
- Open 清单标注已由权威文件 §4/§5/§6 具体化。

### 2.2 §G2 / F-I-027 · 权威收口

| 载体 | 修改 |
|------|------|
| `r1-backup-port-contract-draft-v0.1.md` L60 | 删「`or a C3-specific native snapshot`」式并列，改为**明令** SQLite provider 必须用**转换提交后**的 C3 专属快照；**明确写出** `snapshotBeforePending` **不是** RecoveryPoint 源 |
| 同文件 L71（Open 项） | 标**已解决**并指向 C3 边界文件 §2/§3 |
| 同文件尾 | 新增权威范围声明：本文件只辖 kernel Port 类型表面与后置条件；C3 其余以边界文件为唯一权威 |
| `r1-c2-per-column-conversion-contract-v1.0-fc.md` §0 | 「该 C3 附件**尚未落盘**」→ 已落盘，且为 C3 **唯一权威**；原 runbook 已 `superseded` |

### 2.3 §G3 · PG 对称调用点 —— 已写出（并标注现状缺口）

C3 边界文件 §4.2 拆为 **SQLite 侧**与 **PG 侧**两张表：

- PG：`postgres.migrate` → `applyPendingPG`（`postgres.go:120-127`）循环**之前**生成 A′；循环**内**生成 C′；**成功之后**调 `CreateRecoveryPoint` 产出 B′；`applyMigrationPG`（`:154-173`）事务内**禁止**调用。
- **明确标注现状缺口**：`applyPendingPG` 现行**既无** snapshot 步骤，**也无可对称的 `verifyIntegrity()`**（SQLite 侧有 `migrate.go:345-367`，PG 侧缺失）→ 这三处是**新增设计**，R2 落码须同时补 PG 形状校验入口。

### 2.4 §G4 · A 与 C 分开 —— 已修正

§2 表新增 **A / B / C 三行的精确时机**：A = 批次**循环外一次**（旧合同）；C = **循环内 per-migration**。并**纠正** A-040 指出的我方表述错误：C 的合同形状**不是**「旧合同」而是**混合形状**——随批次推进，中途含**部分**新合同列，批次末为全部新合同。为此在 §2 末新增「A 与 C 必须分开，不可合并为『批次前一次』，也不可互替」。

### 2.5 §G5 · 调用点失败后的重试 —— 新增 §4.3

问题（A-040 指出）：B 生成失败后 SQLite 走 `actionNoop`（`migrate.go:52-53`，只跑 `verifyIntegrity()`）**不重试**；PG（`postgres.go:94-95`）**连 `verifyIntegrity` 都没跑** → **B 可能永不补创**。

冻结要求四条：①「catalog 已到目标版本但无 B 记录」必须是**可检测状态**并产出**显式**动作（非 `actionNoop`）；② 该动作须跑形状校验：形状已是新合同且无 B → **补创 B**；形状仍旧合同 → 走**回滚**（A/C）；③ 重试**有界**（每次启动最多一次），失败留可核对记录，不因补创失败拒绝启动（沿用现行 fail-closed 立场）；④ **禁止**把「无 B」当作已满足 RecoveryPoint 门禁。

### 2.6 §G6 + §D末段 · 反向断言钉错误分类 —— 新增 §5.1

- 新增 **6 类错误分类表**：`TimeContractMismatch` / `TemporalColumnSetIncomplete` / `ChecksumMismatch` / `ArtifactNotFound` / `ArtifactUnreadable` / `ToolFailure`。
- **硬要求**：`TestLegacyArtifactMustFail` 必须断言分类**属于** `{TimeContractMismatch, TemporalColumnSetIncomplete}`，并**明确断言不属于** `{ArtifactNotFound, ArtifactUnreadable, ToolFailure}`——否则「缺文件也会绿」（A-040 原话）。
- 证据落点 = **测试名 + 观测分类**，不另发明附件格式（A-040 §E 已接受）。
- **§5 断言 4 已更正**：删除「负值非法」于 B 的正向 round-trip；负值断言改落 `TestLegacyArtifactMustFail`（A/C）或转换预检（`m0`）。

### 2.7 §G7 · C → B 的机械身份 —— 新增 §3.1

仅靠文件名政策不足。新增：A/C **不得**生成 `RecoveryPoint`（记录只含 `artifact_path` + `catalog_version` + `contract_shape = legacy`）；B 元数据必须含 `TimeContract` + `CatalogVersion` + `ChecksumSet` + `Verification` + 批次末 version，且 `contract_shape` 由**对 restore 目标的实测**得出。判定规则：`CreateRecoveryPoint` 在形状校验阶段**实测** restore 目标的时间列物理类型，形状不符一律归类 `TimeContractMismatch` → **实测驱动而非命名驱动**。

### 2.8 §D · `D-010`/`D-007` 编号限定 —— 已修

§4.1 改为「**Root** `D-007` + **Root** `D-010`（child `D-006` + child `D-009`）」，符合本文件 §9 自己的限定要求（A-040 指出原句无限定，违反自身规则）。

### 2.9 F-I-025 · 计数表述 —— 已改

rebuild 附件 §0 覆盖说明由「15 个 descriptor / **20 张时间列表**」改为「15 个 descriptor，含 **20 张带时间列的表**（分母 90 列），加 3 张**无时间列**联接表 + 3 张跨 descriptor 子表（后者时间列由 v80/v81 另行转换）」。

## 3. 仍开放（不得放行）

- **F-I-002**：可执行边界测试已落地（`D-020`、`apps/api/internal/w040contracttest/`，三测试全绿）；是否据此闭合**待 independent 复审该测试文件**。
- **F-I-004**：本轮 7 项已逐项收口（§2.1–2.7）；**接受与否待 independent 复审**。
- **F-I-005**：需 R2 落码才能记录真实 `MigrationChecksum` 哈希，结构性依赖；测试改写同属 R2。

**本响应不闭合任何 required。** C2/C3 未冻结，R2 未放行。

## 4. 下一步

`/audit` 复审本轮收口：重点 F-I-004 的 7 项（尤其 PG 调用点与 §4.3 重试是否可执行）、F-I-002 的测试文件、以及 F-I-025/F-I-027 两条 recommended 是否可闭合。
