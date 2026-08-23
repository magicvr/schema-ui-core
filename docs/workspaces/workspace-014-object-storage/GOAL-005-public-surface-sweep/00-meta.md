---
id: GOAL-005-public-surface-sweep
title: R4 公共面收尾核查（无本地路径 / os.File）
status: done
parent: GOAL-001-object-storage
created: 2026-08-21
updated: 2026-08-21
version: 0.2.0
progress: 4/4
plan_refs:
  - VP-014-object-storage
primary_plan: VP-014-object-storage
---

# GOAL-005 · R4 公共面收尾核查

## 意图

承载 Root [GOAL-001](../GOAL-001-object-storage/00-meta.md) 纲领路线图 **R4**：以证据化扫描确认 Handler / 模块公共契约不再把本地路径或 `os.File` 当作存储合同。R3 已完成主要迁移，本目标为收尾核查门禁 + 少量加固。

## 范围与检查点

1. 扫描证据：导出函数无存储语义的本地路径参数；`*os.File` 零引用；`uploadDir` 仅存于测试直盘断言。
2. 边界声明：SQL 库路径（store.OpenWithCatalog / systemmonitoring dbPath）属 Store 方言（VP-013），不在本 VP 对象存储合同内。
3. 加固：composition `newObjectStore` 对未知 driver 显式拒绝（A-002 N-005 可选加固落地）。
4. 全量测试绿。

## 父目标

- GOAL-001-object-storage（R4 阶段；VP-014 架构 A2）
