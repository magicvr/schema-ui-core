---
id: A-002-c62-slice1-2-wiring-evidence
doc: audit-entry
goal: GOAL-013-r6-old-path-removal
source: independent
auditor: Grok Build / grok-4.5
date: 2026-08-05
scope: C6.2 切片 1–2 证据审阅（ownership 登记 + CollectPersistence 生产接线）· 切片 3 放行闸门
audit_type: execution-facts | finding-closure
verdict: conditional
status: recorded
parent: GOAL-001-modular-admin-architecture
created: 2026-08-05
updated: 2026-08-05
version: 0.1.0
---

# A-002 · C6.2 切片 1–2 接线证据独立审计（2026-08-05）

- **source**：independent
- **auditor**：Grok Build / grok-4.5
- **类型**：execution-facts（切片 1–2 是否贴 D-002；可否开切片 3）
- **scope**：GOAL-013 C6.2 过渡实现 · [D-002](../01-decision/D-002-r6-persistence-ownership.md) · [设计附件](../attachments/r6-persistence-ownership-design.md) · [E-005](../02-execution/E-005-r6-c62-persistence-wiring.md) · `apps/api/internal/{store,composition,kernel}`
- **verdict**：conditional
- **工作区**：`workspace-003-modular-admin-architecture` · Root `GOAL-001` · `shared_materials_catalog: none`

## 范围与区间

### 覆盖

| 项 | 核对点 |
|----|--------|
| 归属 | 0001–0008 moduleID 是否与 D-002 §2 表一致 |
| 接线 | composition 是否 CollectPersistence → OpenWithCatalog；无生产旁路 |
| 契约 | catalog 与 ledger 元数据一致；fail-closed；相关单测 |
| 诚实边界 | 是否误宣称 F-001/F-002/F-005 或 C6.2/VP 退出完成 |
| 切片 3 准备 | 粒度是否仍沿用 D-002 module 包 |

### 排除

- 不审切片 3 物理迁出实现（尚未完成）。
- 不审 F-005 seed、C6.3 Schema 字节、C6.4 验收。
- 不改 status / progress / goal-tree / 代码。

## 成果（有证据）

### 1. Ownership 登记 vs D-002（通过）

| Version | Name | D-002 moduleID | 代码 `compiledMigrations[].moduleID` |
|---------|------|----------------|--------------------------------------|
| 0001 | r2_baseline | core.auth-session | core.auth-session |
| 0002 | rbac_expand | core.auth-session | core.auth-session |
| 0003 | records_persist | core.persistence | core.persistence |
| 0004 | operation_log | core.operationlog | core.operationlog |
| 0005 | operation_log_expand | core.operationlog | core.operationlog |
| 0006 | records_retire | core.persistence | core.persistence |
| 0007 | site_settings | admin.settings | admin.settings |
| 0008 | operation_log_settings | core.operationlog | core.operationlog |

证据：`apps/api/internal/store/migrate.go` L68–131；`MigrationCatalog()` 将 moduleID 写入 `ContributionIdentity.ModuleID`。

### 2. CollectPersistence 生产接线（条件通过）

| 主张 | 证据 | 判定 |
|------|------|------|
| composition 收集 catalog | `composition.openStore`：`providers := []kernel.Provider{storePersistenceProvider{}}` → `kernel.CollectPersistence` → `store.OpenWithCatalog` | **生产路径成立** |
| 元数据 fail-closed | `validateCatalogMatchesLedger` 比对 version/name/checksum/moduleID；`TestOpenWithCatalogRejectsDivergentCatalog` | **成立** |
| Collect 校验 | `TestMigrationCatalogMatchesLedgerWithOwnership` 调 `CollectPersistence` | **成立** |
| 包测试 | 本审计执行：`go test ./internal/store/ ./internal/composition/ ./internal/kernel/` 通过 | **成立** |
| catalog **驱动执行** Apply | `OpenWithCatalog` 校验后调用 `open()` → `migrate()` 仍遍历 **内部** `compiledMigrations`；**未**执行传入 catalog 的 `Apply` | **过渡态缺口**（F-C62-001） |
| 终态「各模块 CompiledPersistence」 | 仅 `storePersistenceProvider` 聚合返回 `MigrationCatalog()`；admin.* / core.* 模块 Provider 仍 `return nil, nil` | **切片 3 待做**（诚实） |

### 3. 旁路与测试（条件通过）

- 生产 composition：**仅** `OpenWithCatalog`。
- `store.Open`（无 catalog）仍存在，被 handler/auth/modules 测试广泛使用 → 测试旁路可接受，但须标明 **非生产权威路径**（F-C62-002）。

### 4. 叙事诚实（大体通过）

- E-005 与 `OpenWithCatalog` 注释写明：slice 2 接线、**slice 3 才迁 Apply/DDL**。
- C6.2 检查点未勾选；未宣称 VP 退出 #2/#3/#5 取证。
- 注意：R6-I002 在 `00-meta` 标 `verified` 的证据是 **设计可实施**（D-002），不是「接线+迁出完成」——可接受，但 03-audit 索引仍写 I 全 collecting，与 meta 不一致（F-C62-003）。

### 5. 切片 3 粒度（通过 · 沿用 D-002）

无需重开粒度讨论。迁出包：

```text
core.auth-session   → 0001, 0002
core.operationlog   → 0004, 0005, 0008
admin.settings      → 0007
core.persistence    → 0003, 0006（历史）
store               → 平台 runner（消费 catalog，不再持有领域 DDL 权威源）
```

## 对照闸门问题

| 问题 | 结论 |
|------|------|
| 切片 1 ownership 登记是否贴 D-002？ | **是** |
| 切片 2 生产 CollectPersistence 元数据接线是否成立？ | **是**（metadata gate） |
| 是否已「catalog 驱动 Apply」终态？ | **否**（过渡；F-C62-001） |
| 可否开切片 3？ | **可以（conditional）** — 先响应 F-C62-001 边界写清；不要求先做 F-005/C6.3 |
| 可否勾选 C6.2 / 闭合 F-001/F-002/F-005？ | **否** |
| 可否宣称 VP 退出取证？ | **否** |

## Findings

### F-C62-001 · catalog 接线为元数据门禁，非 catalog 驱动执行（须诚实冻结）

- **严重度**：med
- **建议**：required（叙事 / 切片 3 入口边界）
- **状态**：open
- **描述**：`OpenWithCatalog` 仅 `validateCatalogMatchesLedger` 后调用与 `Open` 相同的 `open()`→`migrate()`，执行体仍是包内 `compiledMigrations`。传入 catalog 的 `Apply` **不参与**运行时升级。这与 D-002「改收 catalog **应用**」的终态表述有差距，但作为切片 2 过渡可接受——**前提**是不得把「CollectPersistence 已接线」写成 F-002 已闭合或 C6.2 完成。
- **证据**：`store/store.go` `OpenWithCatalog` / `open` / `migrate`；`composition.openStore`；E-005「切片 3 物理迁出」
- **建议修复**：在 E-005 或 D 补丁一句冻结：  
  - 切片 2 = 元数据收集 + 与 ledger 一致 fail-closed；  
  - 切片 3 = Apply/DDL 迁入模块 `CompiledPersistence`，runner **只执行 catalog.Apply**（或等价：权威源仅 catalog，store 无第二套 compiledMigrations）。  
  响应后本 finding 可 `fixed`（边界澄清）而不要求立刻做完切片 3。

### F-C62-002 · `store.Open` 无 catalog 旁路仍可用于测试

- **严重度**：low
- **建议**：recommended
- **状态**：open
- **描述**：测试与部分 composition 测试仍 `store.Open`，不经 CollectPersistence。生产 composition 已走 OpenWithCatalog，故不阻断切片 3；切片 3 后应逐步让关键升级测试走 catalog 路径，避免双权威。
- **证据**：`handler/testhelpers_test.go`、`modules/*/provider_test.go`、`auth/auth_test.go` 等
- **建议修复**：注释 `Open` 为 test/legacy；或测试改用 `OpenWithCatalog(MigrationCatalog())`。

### F-C62-003 · GOAL-013 审计索引与 meta 信息门禁不同步

- **严重度**：low
- **建议**：required（process-hygiene）
- **状态**：open
- **描述**：`00-meta`：R6-I001/I002 `verified`、C6.1 已勾选；`03-audit.md` 仍写 I 全 collecting、C6.1–C6.4 均未勾选。影响后续编排对门禁状态的读取。
- **证据**：`GOAL-013/00-meta.md` L38–60；`03-audit.md` L17–30
- **建议修复**：编排刷新 03-audit 信息就绪表与结论段，与 meta/E-004/E-005 对齐。

### F-C62-004 · C6.2 / F-001/F-002/F-005 未完成（继承确认）

- **严重度**：high（实现）
- **建议**：required（继承 Root A-010）
- **状态**：open（确认）
- **描述**：切片 1–2 **不**闭合：F-001 上帝对象、F-002 物理迁出、F-005 seed 贡献驱动。本审计确认代码与台账均未假闭合。
- **证据**：store 仍含领域 SQL/seed；模块 `CompiledPersistence` 空；seed 仍中心

## 必改项汇总

| ID | 级别 | 阻断 |
|----|------|------|
| **F-C62-001** | required | 切片 3 **叙事**入口：澄清 metadata vs apply；**不**阻断开始编码切片 3 |
| **F-C62-003** | required | ledger 同步；不阻断切片 3 编码 |
| F-C62-002 | recommended | 测试旁路收敛 |
| F-C62-004 | required（继承） | C6.2 / VP 退出取证前 |

## 切片 3 放行意见

**允许进入切片 3**，条件：

1. 接受 F-C62-001：切片 2 ≠ catalog 驱动执行终态；切片 3 成功标准含「权威迁移源离开 store.compiledMigrations / runner 消费模块 catalog」。  
2. 粒度 **固定为 D-002 module 包**（上表），不重开设计。  
3. 切片 3 **范围**：Apply/DDL + 各 owner `CompiledPersistence` + store 收窄；**不含** F-005、C6.3、领域仓储大搬家（可另切片）。  
4. 完成后跑既有 migrate/recovery 矩阵 + CollectPersistence 多 provider 收集测试。

## 与既有意见

| 意见 | 关系 |
|------|------|
| Root A-010 F-002 | 切片 2 推进「接线」；物理迁出仍 open（F-C62-004） |
| D-002 / 设计附件 | 归属表实现一致；「收 catalog 应用」终态未齐（F-C62-001） |
| E-005 | 事实基本可复核；建议补 F-C62-001 边界句 |

## 结论与下一步

**verdict: conditional**

- 切片 1–2 **证据充分可进入切片 3**。  
- **不得**勾选 C6.2、闭合 F-001/F-002/F-005、宣称 VP 退出。  
- `/govern`：响应 F-C62-001/003（短文档）后推进切片 3；F-005 另切片。

### 声明

本意见 `source: independent`，只写 GOAL-013 `03-audit` ledger。  
**未**改 status/progress/goal-tree/代码。响应归 **`/govern`**。
