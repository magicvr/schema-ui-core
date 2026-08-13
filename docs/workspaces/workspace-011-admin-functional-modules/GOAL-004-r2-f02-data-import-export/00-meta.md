---
id: GOAL-004-r2-f02-data-import-export
title: R2-F02 · 数据导入/导出（schema 驱动 · 共享能力）
status: active
parent: GOAL-001-admin-functional-modules
created: 2026-08-14
updated: 2026-08-14
version: 0.3.0
progress: 4/5
---

# GOAL 004-r2-f02-data-import-export · 数据导入/导出（schema 驱动 

## 概述

一等公民 F-02（I-011-001 `3）：新建**共享能力模块**（admin.data-transfer 候选名），提供 schema 驱动的列表**导出**（CSV/Excel）与**批量导入**（校验/预览/错误报告）。基架现无导出端点（仅 multipart 单文件上传 C-09）。导出权限键 + 操作审计。

## 当前边界

- 导出：任意 schema 列表 → CSV/Excel（列映射/全量或筛选集）
- 导入：CSV/Excel → 校验 → 预览 → 执行 → 错误报告（不回滚语义按方案）
- 权限键（导出/导入）与操作审计（复用 operationlog）
- 协议对照按 `8 必办-1 独立做（export 扩展动作键 ≠ 导出契约 → 本地契约 + fail-open 留痕）

## 成功标准与路线图（P-001）

- [x] **S1 · 方案冻结**：方案冻结：必办-1（协议对照）+ 导入校验/错误报告模型 + 权限键 + 审计设计；方案级 self 审视（D-002 / A-001，2026-08-14）
- [x] **S2 · 实现**：实现：共享服务 + 导出/导入流程 + 页面（导出按钮、导入向导/错误报告）+ 测试（E-003 · 39a1671）
- [x] **S3 · 验证**：验证：单元/集成 + 代表资源（users/roles）实测 + 全量回归（E-003 · go 全绿 + web 893/893 + 冒烟）
- [x] **S4 · go 影响判定 + 自审**：go 影响判定 + self 审计（D-003 无影响不暂挂；A-002 pass）
- [ ] **S5 · 关门**：关门：全部 required 闭合 + 关门审计 + goal-tree 同步

progress: 0/5 由五个等权检查点派生（S1 完成后更新）。

## 审计策略

self + 关门 cross 候选（共享数据能力、权限/审计面 → 独立审计建议 grok build，按 P-004 在关门门禁确认）。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 |
|----|------|-----------------|----------|--------------|-----------------|------|
| I-001 | required | 协议面（v2.8.0）对导出/导入是否有契约？ | S1 方案 | 对照 protocol-inventory（export 扩展动作键）、node.schema.json | **closed**（D-002 §2） |
| I-002 | required | 导入校验/错误报告模型与权限键设计（含大小/类型限制、审计） | S1 方案 | 复用 C-09 上传加固经验（VP-009 W2/W4） | **closed**（D-002 §3/§4） |
| I-003 | non-blocking | Excel 二进制格式依赖（xlsx 库）评估 | S2 实现 | 方案冻结时定默认（CSV 优先） | **closed**（D-002 §8） |

## 父目标

- [GOAL-001-admin-functional-modules](../GOAL-001-admin-functional-modules/00-meta.md)

## 台账布局

本目标从首条记录起使用 `01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger；索引与目录条目共同构成正式记录。