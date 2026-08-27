---
title: A-002 · A-001 响应（self · 编排响应）
source: self
status: recorded
created: 2026-08-23
updated: 2026-08-23
parent: GOAL-036-w25-page-performance-guardrails
version: 0.1.0
scope: 响应 A-001（independent · conditional · F-001～F-005）；含 self 新增 F-006/F-007
verdict: conditional（原 verdict 不回溯；本响应在 P-003 三路径下闭合全部 required）
---

# A-002 · 响应 A-001 独立复审意见（2026-08-23，self）

响应对象：`A-001-correction-review-independent.md`（grok-4.6，`conditional`，open required：F-001 high / F-002 med）。采纳全部 findings；F-001/F-002 按 `fixed` 闭合，F-003～F-005 按 `fixed` 闭合，self 补充 F-006（fixed）/ F-007（记录移交）。原 verdict 与 finding 原文不改写。

## 关闭证据表

| Finding | 状态 | 证据路径 |
|---------|------|----------|
| **F-001（high · required）** SQLite 池化未把 `foreign_keys` 迁入每条连接，CASCADE 在 3/4 连接失效 | **fixed** | `store.go` `sqliteDSNParams` 增补 `_foreign_keys=on`（连接级 PRAGMA 随 DSN 逐连接生效；注释引用 A-001 归因）；e2e admin 复跑 **9 通过 / 1 profile 专属跳过 / 0 失败（exit 0）**，`schema-crud`（删用户→删角色）通过 |
| **F-002（med · required）** S5 防复发测试未钉死 per-conn `foreign_keys` 与 CASCADE/RESTRICT | **fixed** | `store_wal_test.go` 新增 `TestFileStoreEveryConnectionEnforcesForeignKeys`：同时持有 `sqlitePoolDefault` 条 Conn 断言每条 `foreign_keys=1`；在**非首条连接**上验证 `ON DELETE CASCADE`（删用户 → `user_roles` 归零）、`ON DELETE RESTRICT`（在用角色删除被拒）、`refresh_tokens` 指向缺失用户插入被拒；`TestSQLiteDSNPragmas` 断言补 `_foreign_keys=on`；`go test ./internal/store/` ok |
| **F-003（med · recommended）** `refreshList` 不丢 in-flight，与 `reloadList` 不对称 | **fixed** | `render.tsx` `refreshList` 与 `reloadList` 对称：丢弃目标 URL 的 in-flight 后 bump token（注释引用 A-001 F-003）；`render.test.tsx` 新增「挂起期间 refreshList 不 join 旧请求」用例（挂起 fetcher 双 gate，断言第二次调用发出）；vitest 全量 **1097/1097** |
| **F-004（low · recommended）** I-002 测量脚本不入库；`attachments/README.md` 仍写"当前为空"；`00-meta` 路线图 S6 未勾 I-002 | **fixed** | `attachments/README.md` 重写（证据清单 + 脚本位置说明；脚本入库存为后续可选）；`00-meta.md` 路线图 S6 行勾选 I-002 |
| **F-005（low · recommended）** 批删回归名义含 MFA、断言只有 `user_roles` | **fixed** | `TestDeleteUsersBatchCleansRoleAndMfaLinks` 补 `user_mfa` 播种与删除后计数断言；`go test ./internal/modules/authsession/` ok |
| **F-006（self 新增）** `TestOperationLogAuthEvents` 的顺序断言依赖单连接串行写序，池化后同毫秒写入顺序不再由提交序决定（随机 id 后缀决胜），间歇失败 | **fixed** | `operations_test.go`：位置顺序断言改为**事件集合断言** + 保留逐条属性断言（Actor/Detail/SessionID/敏感键）；单跑 ×1 与全量通过 |
| **F-007（self 新增 · 记录）** `TestScheduledTaskRunsPagination` 偶发 500（`POST /run` 第 2/3 次） | **记录移交** | **预存 flake，非本波引入**：`git stash` 至基线（HEAD，单连接、无 FK 改动）后 `-count=20` 复现（`run 2 = 500`）。候选机制：`scheduler.go` `newRunID()` 的 crypto/rand 失败回落路径（时间串 hex，同毫秒重复 → 主键冲突 500）。本响应不并入；建议小修（回落路径加计数/多源随机）并单独跟踪，不阻断本响应 |

## 叙事修正（I-001 归因，P-003 留痕、不回溯 E-004 原文）

E-004 将 admin e2e 失败归因于「`DeleteUser` 从不清理 `user_roles`」——SQL 字面成立；A-001 补充的**主要机制**为：W25 池化后 `foreign_keys` 未随 DSN 进入每条连接 → `ON DELETE CASCADE` 在 3/4 连接失效（探针实证）→ 孤儿行出现 → `deletable=false` 永久化；`DeleteUser` 无显式清理为**第二层**。两者均已修复：FK 入 DSN（每条连接）+ 显式删 `user_roles`/`user_mfa`（纵深，`user_mfa` 无 CASCADE，显式删为必需）。I-001 维持 **closed**（material + 本响应复跑 e2e 9/9）；I-002 维持 closed。

## 仍开放项

- F-007（预存 flake）移交独立小修（建议在后续波次/govern 立项；候选修法：`newRunID` 回落路径保证唯一性，或该测试容忍重试）。
- 关门自审与关门决定仍待用户（A-001 编号已被独立审占用，后续自审将顺延编号 A-003）。

## 结论

A-001 的 2 条 open required（F-001 high / F-002 med）已按 `fixed` 合法闭合，3 条 recommended 已 fixed，self 新增 1 fixed + 1 记录移交。验证证据：go 定向套件（store + authsession + handler除F-007预存flake）全绿、vitest 1097/1097、tsc 0、e2e admin 9/9（F-001 关闭直接证据）。建议下一步由用户决定：接受响应 → 复跑一次 mvp e2e 作伴行证据 → 走关门路径。