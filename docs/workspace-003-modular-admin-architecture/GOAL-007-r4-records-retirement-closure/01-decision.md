---
id: GOAL-007-r4-records-retirement-closure
doc: decision
status: active
parent: GOAL-005-r4-full-module-migration
created: 2026-08-05
updated: 2026-08-05
version: 0.1.0
---

# 决策记录 · GOAL-007

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 | 影响门禁 | 状态 | 证据 |
|----|------|----------|------|------|------|
| R7-I001 | required | 当前运行面是否有 Records 专属实现 | C7.2/C7.3 | verified | 扫描附件、定向测试 |
| R7-I002 | required | 迁移/审计/通用协议保留边界 | C7.2/C7.3 | verified | D-003、`0003`/`0006`、operation-log |
| R7-I003 | non-blocking | 历史文档是否显式标注为历史证据 | C7.4 | open | 历史治理目录 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-05 | Records historical-only 退场核验范围 | accepted | [01-decision/D-001-r4-records-retirement-scope.md](01-decision/D-001-r4-records-retirement-scope.md) |

## 当前约束

- 不删除已应用迁移定义或历史 operation-log 合法值；这不是 Records 产品面。
- 不把通用 `Record` 类型、`recordView`、`recordSource`、`record_id` 或测试用例中的
  协议 URL 误判为 Records 实体实现。
- C7.3/C7.4 未完成前，本目标不能 `done`，GOAL-005 C4/C5 不能由本目标单独放行。
