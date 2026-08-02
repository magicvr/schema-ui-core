---
title: 审计台账 · R3 · 持久化 RBAC、菜单投影与版本迁移
status: active
created: 2026-08-02
updated: 2026-08-02
parent: GOAL-001-production-admin-foundation
version: 0.2.0
---

# 审计台账 · GOAL-006

## 正式意见索引

| 编号 | source | 日期 | scope | verdict | 状态 |
|------|--------|------|-------|---------|------|
| A-001 | independent | 2026-08-02 | S1 execution-facts + 目标定义/信息门禁 | pass | open（无 required finding；recommended 待 `/govern` 择要响应） |

## 当前审计边界

- A-001 覆盖：目标定义与 D-001～D-004、`I-006-001/002` 门禁状态，以及 **S1**（版本迁移与可恢复起点）执行事实对照。
- **不**覆盖 S2～S6 实现或 Root R3 阶段关门；S2～S6 仍未勾选，本意见不构成其放行。
- `I-006-001/002` 保持 `verified`；当前无到期未关的 required 信息项。
- 审计意见不直接修改 `status` / `progress`；响应和推进由 `/govern` 与用户裁决维护。

---

## A-001 · S1 执行事实与目标定义交叉审计（2026-08-02）

- **source**：independent
- **auditor**：GitHub Copilot（Grok 4.5）· `/audit`
- **类型**：execution-facts（主）+ goal-definition / info-gate（辅）
- **scope**：`workspace-002-production-admin-foundation` / `GOAL-006-r3-persistent-rbac-menu` · S1 检查点主张；D-001～D-004 与 `I-006-001/002`；不审 S2～S6 产品行为
- **verdict**：**pass**
- **工作区**：`workspace-002-production-admin-foundation`；`root_goal=GOAL-001-production-admin-foundation`；`canonical_scope` 匹配；`shared_materials_catalog: none`（本审未将共享资料作证据）

### 范围与区间

| 项 | 结论 |
|----|------|
| 工作区绑定 | 显式工作区已校验；目标 `parent` 指向 Root；未读取其他工作区上下文 |
| 信息门禁 | `I-006-001`（S1/S2/S3 实施前）与 `I-006-002`（S5 前）均为 `verified`，附件与 D-002/D-003 可互指；无开放 required |
| S1 成功标准 | 「schema_migrations + 顺序/校验和 + 事务化 fail-closed + 升级前可恢复副本」 |
| 代码证据根 | `apps/api/internal/store/migrate.go`、`migrate_test.go`、`store.go`、`store_test.go` |
| 文档证据 | `00-meta` S1 勾选与 `1/6`；`01-decision` D-001～D-004；`02-execution` S1 时间线；附件 I-006-001/002 |

### 成果（有证据）

1. **版本台账与 runner**：`schema_migrations(version,name,checksum,applied_at)`；编译期 `0001 r2_baseline` / `0002 rbac_expand`；checksum = 规范化 SQL + `transformID` 的 SHA-256 小写 hex；未知版本 / 缺号（台账不自 1 起或中间断层）/ name 不符 / checksum 漂移均 fail closed（`validateApplied`）。
2. **事务化迁移**：`applyMigration` 单事务执行 `up` + 台账插入；0002 非法 roles 时 0001 保留、0002 整步回滚（`TestMigrateFailClosedInvalidRoles`）。
3. **0001 指纹**：空库建 R2 表+台账；既有 R2 经 table/column/FK/index 指纹后登记；部分基线拒绝且不留空台账（`TestMigrateFailClosedPartialBaseline`）。
4. **0002 DDL + 回填**：六张 RBAC 表及反向索引与 I-006-001 §4 / D-002 一致；role key 正则与用户内去重（`TestMigrateExistingR2DedupeRoles`）。
5. **pre-v0002 快照**：非空文件库在首个 `version>=2` 待迁移前 `VACUUM INTO <db>.pre-v0002-<UTC>.sqlite`，快照 `integrity_check=ok`；空库不产生快照；迁移后主库 `integrity_check` + `foreign_key_check`（`TestMigrateExistingR2DB` / `verifyIntegrity`）。
6. **Open 契约**：`Open(path, adminUsername, adminPasswordHash, seedAdmin)` 签名保留；旧 `CREATE TABLE IF NOT EXISTS` 路径已移除；`seedAdmin` 仍仅写 `users` JSON（S2/S3 范围，D-004 已界定中间态）。
7. **回归（本轮独立复跑，2026-08-02）**：
   - `go test ./internal/store/ -count=1` → PASS（含 9 个迁移相关 + 3 个既有 store 用例）
   - `go test ./... -count=1`（`apps/api`）→ PASS
   - `go vet ./...` 干净；`gofmt -l ./internal/store/` 无输出
8. **目标定义**：单一端到端子目标 + 六顺序检查点；D-004 明确 S1 含 0002 链与 pre-v0002、S2 独有读路径切换，避免范围抢跑；进度 `1/6` 与仅 S1 勾选一致。

### 对照成功标准（S1）

| 标准要素 | 判定 | 证据 |
|----------|------|------|
| `schema_migrations` | 满足 | DDL + 台账读写 + fresh/existing 测试 |
| 顺序 / 校验和检查 | 满足 | `validateApplied` + Unknown/Missing/Checksum 测试 |
| 事务化 fail-closed | 满足 | `applyMigration` 回滚 + InvalidRoles/PartialBaseline |
| 升级前可恢复副本 | 满足 | `snapshotPreV0002` + ExistingR2 快照查询复现 |
| 未越权勾选 S2～S6 | 满足 | meta 仅 S1；读路径/seed grants/授权/菜单仍属后续 |

### Findings

#### F-001 · recommended · 低

- **主题**：`TestMigrateFailClosedMissingIntermediate` 通过删除 version=1、仅留 version=2 触发「台账不自 1 起」，并未构造真正的中间缺号（如已应用 1 与 3）。
- **证据**：`migrate_test.go` `TestMigrateFailClosedMissingIntermediate`；对比 `validateApplied` 中 `a.version != applied[i-1].version+1` 分支缺少对应用例。
- **影响**：fail-closed 行为仍被覆盖；「缺中间版本」分支的回归保护偏弱。
- **建议**：S6 或下一迁移相关改动前，补一条 ledger=`(1,3)`（或等价）的用例。

#### F-002 · recommended · 低

- **主题**：I-006-001 `V-MIG-04` 中 unique / ON DELETE CASCADE|RESTRICT / 反向索引存在性未在 S1 测试中逐项覆盖；当前仅有 `TestForeignKeyEnabled`（refresh_tokens → 缺失 user）。
- **证据**：`migrate_test.go`；`rbacExpandDDL` 已声明约束但无对应正反测试。
- **影响**：不否定 S1 成功标准（标准未要求完整 V-MIG-04 矩阵）；S2～S4 开始依赖这些删除/唯一语义时证据不足。
- **建议**：在 S2/S3 或 S6 回归中补 RBAC 表约束与索引断言，勿把「DDL 已写出」当作已验证行为。

#### F-003 · recommended · 低

- **主题**：`02-execution` 写「路径经驱动绑定」；实现为 SQL 字符串字面量 + 单引号转义（`VACUUM INTO '…'`），非 `?` 参数绑定。不经 shell，路径来自 `s.path`，安全性方向正确，表述略过满。
- **证据**：`migrate.go` `snapshotPreV0002`；`02-execution.md` S1 节。
- **影响**：无功能缺陷；后续若有路径注入审查，应以实现为准修正行文。
- **建议**：`/govern` 响应执行记录时改为「驱动内 SQL 字面量转义，不经 shell」。

### 必改项汇总

- **required / 必改**：无
- **recommended**：F-001、F-002、F-003（均不阻断 S1 事实认定，也不阻断进入 S2 实施；F-002 建议在依赖 RBAC 约束的阶段前闭合）

### 与既有意见的异同

- 本目标此前无正式 A-00N 意见（索引原为「尚无正式审计意见」）。本条为首条 `independent` 意见。

### 结论 + 建议给编排器/用户的下一步

- **结论**：在声明的 S1 scope 内，实现、测试与文档主张一致；S1 勾选与 `progress: 1/6` **有可重复证据支撑**；信息门禁关闭方式合规；未发现 high/required 缺口或名不副实的完成声明。
- **建议 `/govern`**：
  1. 响应 A-001（可记录 recommended 的接纳/延期）。
  2. 在用户确认后推进 **S2**（阶段 A/B 读路径切换、规范化双写与集合比对）——注意全新库 `seedAdmin` 后 `user_roles` 仍可能为空，属 D-004 已承认的中间态，应由 S2/S3 闭合。
  3. **不要**将本 pass 解读为 Root R3 或 S2～S6 放行。

### 声明

本意见 `source: independent`，**不修改**目标 `status` / 检查点 / 派生 `progress` / 方案正文 / `goal-tree`。  
响应、finding 闭合与阶段推进由用户通过 **`/govern`** 处理。