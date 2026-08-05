---
id: GOAL-007-r4-records-retirement-closure
title: R4 · Records 历史演示实体退场核验与残余清理
status: done
parent: GOAL-005-r4-full-module-migration
created: 2026-08-05
updated: 2026-08-05
version: 0.1.0
progress: 4/4
plan_refs:
  - VP-003-modular-admin-architecture
primary_plan: VP-003-modular-admin-architecture
serves_summary: 承接 D-003 的 historical-only 裁决，证明 Records 当前产品面已退场，清理误导性演示命名，并保留迁移账本、历史 operation-log 与通用协议兼容边界。
---

# GOAL-007 · Records 历史演示实体退场核验与残余清理

## 概述

本子目标承接 GOAL-005 D-003 和 GOAL-006 D-003。它只处理历史 Records 范例实体
在当前运行面中的残余核验与误导性命名清理，不恢复 Records CRUD，不改写已应用迁移，
也不删除历史治理证据或通用 `recordView`/`recordSource` 协议能力。

## 愿景对齐

| 字段 | 值 |
|------|----|
| `parent` | `GOAL-005-r4-full-module-migration` |
| `plan_refs` | `VP-003-modular-admin-architecture` |
| `primary_plan` | `VP-003-modular-admin-architecture` |
| Charter | `schema-ui-core-admin-foundation@0.2.0`（经 VP-003 间接对齐） |
| 审计模式 | `cross`；使用 Grok Build `grok-4.5` / `high` 独立复审 |

## 成功标准

- [x] **C7.1 / 范围继承**：D-003 已接受，保留/删除边界和父子 handoff 已落盘。
- [x] **C7.2 / 运行面清理**：API/Web/manifest/fixture/CI 无 Records 专属产品面；
  通用测试不再用 Records 演示命名误导读者。
- [x] **C7.3 / 验证与交叉审计**：相关测试通过，self + Grok independent 无开放
  required finding。
- [x] **C7.4 / 子目标关门**：更新父目标与 goal-tree，提交本子目标 close checkpoint。

四个检查点等权；`progress: 4/4` 表示范围继承、运行面清理、Grok A-003 `pass` 复审
与 C7.4 关门均完成。完成本子目标不关闭 GOAL-005、Root 或 VP-003；Records
historical-only 证据已回传 GOAL-005 C1/C4/C5。

## 信息门禁

| 编号 | 级别 | 所需信息 | 影响 | 最晚阶段 | 收集动作 | 状态 | 证据 |
|------|------|----------|------|----------|----------|------|------|
| R7-I001 | required | 当前 API/Web/manifest/fixture/CI 是否仍有 Records 专属运行面？ | C7.2/C7.3 | C7.2 | 全仓静态扫描、目标代码核对、测试 | verified | `attachments/r7-records-scan.md`；Grok A-003；`TestRetiredRecordsRoutesUnregistered` |
| R7-I002 | required | 哪些 Records 相关代码必须保留为迁移/审计兼容？ | C7.2/C7.3 | C7.2 | 核对迁移 checksum、operation-log 和协议消费者 | verified | GOAL-005/006 D-003；I-011-002；Grok A-003 |
| R7-I003 | non-blocking | 历史治理文档中的 Records 文字是否都已明确为历史证据？ | C7.4 | C7.3 | 扫描并记录历史路径，不改写已关闭目标 | open | 历史 workspace-001/002 台账；不阻断关门 |

## 阶段路线图

1. 继承 D-003 和历史退场契约，确定保留边界（已完成）。
2. 扫描并清理当前代码中的专属 Records 残余和误导性演示命名（已完成）。
3. 运行 API/Web 定向测试，完成 self + Grok independent 复审（已完成，Grok A-003
   `pass`，A-004 闭合 F-IND-007-001/002）。
4. 关闭本子目标并把 Records historical-only evidence 回传 GOAL-005 C4/C5（已完成，
   见 E-004）。

## 范围与非目标

范围包括 API/Web 当前代码、运行时 manifest/fixture、测试命名、CI residue guard 和
与 D-003 直接相关的治理证据。非目标包括改写 `0003`/`0006`、删除历史 `records.*`
operation-log 合法值、删除通用 `recordView`/`recordSource`/`RecordID`，以及恢复任何
Records 产品资源。
