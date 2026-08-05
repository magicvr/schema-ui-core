---
id: GOAL-005-r4-full-module-migration
title: R4 · 全量一方模块迁移
status: active
parent: GOAL-001-modular-admin-architecture
created: 2026-08-05
updated: 2026-08-05
version: 0.1.1
progress: 1/5
plan_refs:
  - VP-003-modular-admin-architecture
primary_plan: VP-003-modular-admin-architecture
serves_summary: 在 R3 试点门闩通过后，将 Users、Roles、仍存在的 Schema-driven Admin 能力及其他现有 Admin 能力迁入统一模块契约；Records 仅保留历史迁移与审计兼容，不恢复无语义演示实体。
---

# GOAL-005 · R4 全量一方模块迁移

## 概述

本子目标承接 Root 的 R4 阶段。R3 已通过并提交；R4 现在只建立范围、信息
门禁和迁移计划，不把已有 Users/Roles CRUD 代码或 R3 模块代码直接宣称为
R4 迁移完成。R4 的最终目标是统一模块契约下的完整一方 Admin 能力覆盖，
并清除 Shell/中央注册特例；R5 运维/升级恢复与 R6 终态删除不在本目标内。

## 愿景对齐

| 字段 | 值 |
|------|----|
| `parent` | `GOAL-001-modular-admin-architecture` |
| `plan_refs` | `VP-003-modular-admin-architecture` |
| `primary_plan` | `VP-003-modular-admin-architecture` |
| Charter | `schema-ui-core-admin-foundation@0.2.0`（经 VP-003 间接对齐） |
| 审计模式 | `independent`；关键阶段使用 Grok Build `grok-4.5` / `high` |

## 成功标准

- [x] **C1 / 范围与信息冻结**：完整一方能力清单、六项模块能力映射、
  contribution/provider 方案、Records/Schema CRUD 范围和 operationlog 边界
  均已由证据或用户裁决关闭 required 信息门禁。
- [ ] **C2 / 模块契约扩展**：模块能以结构化方式贡献 HTTP、Schema、授权、
  Navigation、Manifest 和 Persistence；冲突保持 fail closed。
- [ ] **C3 / Users 与 Roles 迁移**：旧中心注册/Schema 所有权清除，保留现有
  CRUD、授权、角色分配、最后管理员保护、密码和 operationlog 行为。
- [ ] **C4 / Schema 与其他能力迁移**：严格按 C1 冻结的范围执行；Records
  保持 historical-only，不恢复不存在的 CRUD 代码或产品面。
- [ ] **C5 / 验收与关门**：同一 Web 构建双 Profile、API/Schema/Manifest/授权/
  持久化和行为矩阵通过；self + Grok independent 审计无开放 required finding，
  并形成进入 R5 的结论。

五个检查点等权；当前为 `progress: 1/5`（C1 已由 GOAL-006/007 承接并以 done 收束）。
完成本子目标只表示 R4 关闭，不代表 Root、VP-003、R5 或 R6 关闭。

## 信息门禁

| 编号 | 级别 | 必须回答的问题 | 影响 | 最晚阶段 | 收集动作 | 状态 | 证据/延期 |
|------|------|----------------|------|----------|----------|------|-----------|
| R4-I001 | required | 当前一方 Admin 能力、中央注册点和模块所有权的完整清单是什么？ | C1 范围冻结、C3/C4 实施 | C1 关闭前 | 扫描 API/Web/Schema/Navigation/迁移/测试并建立映射 | verified | `attachments/r4-c1-capability-inventory.md`；D-002/E-005；A-002 F-GROK-R4-001 fixed |
| R4-I002 | required | 现有 Kernel contribution metadata 如何扩展为可验证的 HTTP/Schema/Authorization/Navigation/Manifest/Persistence provider？ | C2 方案冻结、所有迁移 | C1 关闭前 | 用户整包接受冻结包为 D-003 契约正文（含字段/双检/顺序/`CompiledPersistence()`/owner matrix）；最终审计核对实现边界 | verified | GOAL-006 D-003；freeze package `status: accepted`；最终独立复审 |
| R4-I003 | required | VP-003 的 `records/Schema CRUD` 是恢复 Records 产品 CRUD，还是只迁移当前仍存在的 Schema-driven Admin 能力并保留 Records 历史事件？ | C1 范围冻结、C4 实施、R4 验收 | C1 关闭前 | 用户书面裁决、D-003、退场核验子目标 | verified | 用户裁决 historical-only；GOAL-005/006 D-003；GOAL-007 |
| R4-I004 | required | Users/Roles 业务写入与 operationlog 的一致性、失败语义和 retention 边界是否保持当前契约？ | C3 行为验收、C5 数据门禁 | C1 关闭前 | 用户接受 Option A 并书面接受 bounded residual；最终审计核对实现边界 | accepted-residual | GOAL-006 D-003；owner `magicvr`；review `2026-08-05 08:32:22 +08:00` |
| R4-I005 | non-blocking | R4 期间是否提供 hosted E2E/容器运行环境作为补充证据？ | C5 证据强度 | C5 | 记录可用环境和本地替代边界 | open | 不阻断本地 C1/C3，但不得升级本地结果为 hosted CI |

## 阶段计划

1. C1：完成能力清单、Provider、Records 范围和 operationlog 边界；三项 P-004
   裁决已落盘，Provider 精确契约经 GOAL-006 整包冻结，Grok A-006 `pass` 完成最终
   freeze review；C1 required 信息门禁已闭合，仍待 GOAL-007 运行面关门后由
   `/govern` 放行 C2。
2. C2：按冻结方案扩展 Kernel/组合根模块贡献边界，先用冲突和最小模块测试
   验证 provider contract，再迁移业务。由 [GOAL-008-r4-c2-module-contract-extension](../GOAL-008-r4-c2-module-contract-extension/00-meta.md)
   承接。
3. C3：按模块切片迁移 Users、Roles，逐项核对现有 API/Schema/授权/持久化/日志。
4. C4：按 C1 对仍存在的 Schema-driven Admin 能力和其他能力作实施；Records
   保持 historical-only，不恢复产品 CRUD。
5. C5：运行完整矩阵、独立审计、required finding closure，提交 R4 close-out。

## 范围与非目标

范围包括 Users、Roles、VP-003 经 D-003 收敛后的 Schema-driven Admin 能力、其他
当前一方 Admin 能力的模块化迁移，以及旧 Shell/中央注册/Schema 特例清除。Records
仅作为历史迁移与 operation-log 兼容边界存在，不作为当前产品资源。

非目标包括 R5 Profile 运维与配置收敛、fresh/reconcile、readyz/诊断、Docker/
代理/fork 文档和 R6 生产旧路径删除；这些边界不得在 R4 中偷渡。
