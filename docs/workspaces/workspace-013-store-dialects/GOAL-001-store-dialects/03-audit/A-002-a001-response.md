---
id: GOAL-001-store-dialects
doc: audit-entry
record_id: A-002
source: self
scope: 响应 independent A-001（F-001~F-005 全部 recommended）
verdict: pass
status: recorded
parent: null
created: 2026-08-21
updated: 2026-08-21
version: 1.0.0
---

# A-002 · 编排器响应 independent A-001（Root 闭门依据 + recommended 处置）

## 响应对象

- independent **A-001**（grok-4.6 · 思考强度 high · `pass`）：Root close-out 以**代码 + 本机复跑 + HEAD CI** 支撑，开放 required = 0。
- 编排器判定：**接受 A-001 作为 Root 闭门依据**；Root 维持 `done 5/5`（此前 F-005 指出的「标 done 时 Root 03-audit 为空」程序缺口，由本条 A-001 + 本 A-002 补齐）。

## Finding 处置表

| A-001 Finding | level | 处置 | 证据 |
|---------------|-------|------|------|
| F-001 sqlite `Store.WithTx(*sql.Tx)` 残留 + 过时注释 | recommended/med | **fixed（注释）+ 保留为文档化测试适配器**：`store/store.go` 更新过时注释（R4 已收口模块公共面；`WithTx` 是 sqlite 测试/适配接缝，不进生产契约）；方法按卫生债保留（测试与 testsupport 依赖）。不属于生产公共契约，不阻断闭门 | store.go L70-72/L106+ |
| F-002 连接池旋钮未接到配置/composition | recommended/med | **fixed**：新增 `db.pool_max_open` / `db.pool_max_idle` / `db.conn_max_lifetime`（env `DB_POOL_MAX_OPEN/DB_POOL_MAX_IDLE/DB_CONN_MAX_LIFETIME`，0=默认）；`composition.openStore` 透传 `store.OpenOptions` 池字段；`TestDBPoolKnobs`（yaml+env+0 默认）+ 文档（config.yaml/config.default.yaml/env.example） | config.go / composition.go / config_test.go `TestDBPoolKnobs` / yaml |
| F-003 公共面 `sql.ErrNoRows` 残留 | recommended/low | **fixed（局部）**：jobs `repository.go`/`model.go` `sql.ErrNoRows` → `kernel.ErrNoRows`（去掉 repository 未用 `database/sql` 导入）；模块内部 `sql.Null*` 留在 store 适配层（记录为已知卫生渠，不阻断，不追改） | jobs/*.go |
| F-004 pg_dump 17 client ↔ PG 15.4 `transaction_timeout` SET 告警 | recommended/low | **fixed**：备份/恢复合同（GOAL-006 D-002）补一句——dump 客户端主版本应 ≤ 服务器，或恢复时接受/允许未知 GUC 的 SET 告警（功能与 checksum 已证成立） | GOAL-006 D-002 |
| F-005 标 done 时 Root 03-audit 为空（程序缺口） | recommended/low | **fixed（台账补齐）**：Root close-out 依据改为 本 A-001（independent 代码/复跑/CI）+ 本 A-002；本索引登记 A-001/A-002 | GOAL-001 03-audit.md |

## 仍开放（非门禁）

- `sql.Null*` 模块内部扫描类型（建议留 store 适配层；低优先级卫生项）。
- `Store.WithTx` 测试适配器（按文档化卫生债保留）。

以上均 **recommended**，不构成闭门 required；Root `done 5/5` 保持成立。

## 结论与下一步

A-001 `pass`（无 required）→ 响应 `pass`；Root 闭门依据落地为「独立审计 A-001 + 本响应」在 `03-audit` 台账留痕。建议后续把 F-002/F-003(partial) 已修内容随下一批提交进入 CI；无需重开目标。
