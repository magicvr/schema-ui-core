---
id: GOAL-001-store-dialects
title: Store 双方言（PostgreSQL 生产权威 + SQLite 内嵌）
status: done
parent: null
created: 2026-08-20
updated: 2026-08-20
version: 0.9.0
progress: 5/5
plan_refs:
  - VP-013-store-dialects
primary_plan: VP-013-store-dialects
serves_summary: 交付架构 A1：内核持久化端口 + PostgreSQL 实现 + 现有迁移台账对写；SQLite 保留为 dev/mvp/快测默认；无 ORM。不承载 Admin 功能或业务域。
---

# GOAL-001 · Store 双方言（PostgreSQL 生产权威 + SQLite 内嵌）

## 概述

本 Root 承载 [VP-013-store-dialects](../../../vision/plans/VP-013-store-dialects.md)（**`active`**）的实现：把现行 SQLite-only `store` 收成内核持久化端口，补齐 PostgreSQL 方言，并把开区时 compiled-global 台账对写到两方言。

**边界**：不引入 ORM；不强制本地默认改成 PostgreSQL；不承接对象存储 / Redis / 队列；不加 Admin 页面或业务域表。安全 finding → VP-009；符合性 gap → VP-010。

## 纲领路线图（P-001）

| 阶段 | 内容 | 先后 | 状态 |
|------|------|------|------|
| R1 | **端口与配置面冻结**：Tx 公共类型 ≠ `*sql.Tx`；方言由配置选择；缺省 `db.path` SQLite；PG DSN 键名冻结；v1.4 补 path 扩展名谓词、`COLLATE NOCASE`、checksum 输入与嵌套 Run 检测 | 起点 | ✅ GOAL-002（D-001 / A-001；D-002 / A-003 闭合 A-002；D-003 / A-005 闭合 A-004；D-004 / A-007 闭合 A-006 required；**D-005 / A-009** 闭合 A-008 recommended；合同 **v1.4**） |
| R2 | **PostgreSQL 接入**：驱动、连接池、`readyz` 扩依赖 | 依赖 R1 | ✅ GOAL-003（pgx v5 stdlib；Open/Ping/WasFresh；self A-001 + independent A-002 pass；A-003 关闭 F-001～F-005；**done**，2026-08-20） |
| R3 | **台账对写**：开区时全部 compiled 迁移两方言 apply + checksum | 依赖 R2 | ✅ GOAL-004（48 迁移双写；live PG 全量 boot + 台账 + 系统级合规；self A-001~A-006 + independent A-005；**done 5/5**，2026-08-20） |
| R4 | **仓库公共面收口**：Handler / 模块公共契约去掉 `*sql.Tx` 与驱动类型 | 依赖 R1；可与 R3 部分并行 | ✅ GOAL-005（全仓 kernel.Store/kernel.Tx + D 链；postgres 完整启动 live；LIKE/COLLATE/instr 等价改写；self A-001~A-005 + independent A-004；**done 6/6**，2026-08-20） |
| R5 | **双路径证据**：SQLite 默认路径回归 + PostgreSQL 生产向验收（迁移、共事务、备份合同） | 依赖 R3/R4 | ✅ GOAL-006（升级策略 I-001 + 备份合同 I-004 verified；跨模块共事务 + 数据迁移原型 live；independent A-001 → A-003 fixed；**done 5/5**，2026-08-20） |

`progress` = 已完成阶段数 / 5。当前 `5/5`（R1–R5 全部完成；**Root done，2026-08-20**）。progress 不放行、不关门。

## 成功标准（方向级）

1. 内核持久化端口落地；公共契约无 `*sql.Tx` / 驱动类型。
2. 开区时 compiled 台账在 SQLite 与 PostgreSQL 上均可 apply + checksum。
3. SQLite 仍为本地/Compose 默认；无 PG 仍能开发与快测。
4. 生产向验收以 PostgreSQL 为准（迁移、`readyz`、共事务、备份/恢复合同之一可核对）。
5. 未引入 ORM；未改 Charter；未进 Admin 功能/业务域。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 存量 SQLite 文件库到 PostgreSQL：in-place 升级是否可行，还是只支持 dump/restore / fresh bootstrap | R5 验收；退出判据 2 | R5 开始前 | D-002 + 数据迁移最小原型 | **verified** | 2026-08-20（GOAL-006 D-002/E-002） | in-place 跨引擎不可行（有界 residual）；fresh bootstrap（live 证明）+ sqlite→PG 逻辑迁移原型 round-trip PASS |
| I-002 | required | PostgreSQL 驱动选型（`database/sql` + pgx stdlib / 其他）；须兼容内核端口且禁止 ORM | R2 方案冻结 | R2 实施前 | Root D-002 决策 + `go get` 编译证据 | **verified** | 闭合于 D-002（2026-08-20） | pgx v5 stdlib（驱动名 `pgx`）；lib/pq 维护模式排除 |
| I-003 | non-blocking | 哪些模块公共 API / 内核类型泄漏 `*sql.Tx` | R4 范围 | R4 方案 | 代码扫描清单 | **verified** | 2026-08-20（R4 收口） | 全仓公共面无 `*sql.Tx`（GOAL-005）；含 instr/LIKE/COLLATE 运行时债改写 |
| I-004 | required | PG 备份/恢复合同（替代 `VACUUM INTO` 的生产路径）具体形态 | R5；退出判据 4 | R5 开始前 | D-002 + pg_dump/pg_restore 验证 | **verified** | 2026-08-20（GOAL-006 D-002/E-002） | `pg_dump -F c`→`pg_restore` round-trip 实跑通过（catalog 48 迁移/35 表 checksum 一致；独立审复核） |

## 父目标

- null（Root；Charter `schema-ui-core-admin-foundation@0.2.0` / VP-013）

## 台账布局

新目标为三个可追加台账创建同名平铺目录：`01-decision/`、`02-execution/`、`03-audit/`。索引文件保留 frontmatter、摘要和条目索引；独立记录使用 `D-NNN-*`、`E-NNN-*`、`A-NNN-*` 文件。
