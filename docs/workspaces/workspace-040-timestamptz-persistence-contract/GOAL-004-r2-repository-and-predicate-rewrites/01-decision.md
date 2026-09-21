---
id: GOAL-004-r2-repository-and-predicate-rewrites
doc: decision
status: active
parent: null
created: 2026-09-20
updated: 2026-09-21
version: 0.2.0
---

# 决策记录 · GOAL-004（R2 M3）

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 | 影响门禁 | 状态 | 证据 / 决策 |
|----|------|----------|----------|------|-------------|
| I-041-003 | required（继承） | PG 侧可执行验证环境 | A/B/C | verified | Root `D-017` §3（用户 2026-09-20）：常驻 PostgreSQL 15.4 实测执行；破坏性 migration 只可作用于一次性/专用测试 database |
| I-041-005 | required | 双方言绑定载体（SQLite canonical 字符串 / PG `time.Time`）与统一入口 | **verified**（`D-001` §1；判据由 independent 复审） | 见 `01-decision/D-001-temporal-binding-and-predicate-scope.md` |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|--------|------|
| D-001 | 2026-09-20 | M3 绑定载体、谓词改造范围与回归载体 | accepted | `01-decision/D-001-temporal-binding-and-predicate-scope.md` |
| D-002 | 2026-09-21 | `F-I-002` 的 `user-overruled` 正式载体（用户 2026-09-20 裁决 + 2026-09-21 授权补写；结论不变） | accepted | `01-decision/D-002-fi002-user-overruled-record.md` |

> 新决策从 `01-decision/D-NNN-<slug>.md` 写入；编号在本目标内单调不复用。
