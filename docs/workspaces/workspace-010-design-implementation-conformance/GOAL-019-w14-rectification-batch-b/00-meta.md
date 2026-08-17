---
id: GOAL-019-w14-rectification-batch-b
title: W14 整改批 B · 数据一致性与健壮性（F-05 列表端点校验 / F-06 错误码与目录 / F-07 搜索排序一致性）
status: done
parent: GOAL-015-w14-user-perspective-review
created: 2026-08-17
updated: 2026-08-17
version: 0.2.0
progress: 4/4
---

# GOAL-019 · W14 整改批 B（F-05～F-07 一致性硬化）

[GOAL-015](../GOAL-015-w14-user-perspective-review/00-meta.md)（W14）的**下级整改子目标（批 B）**：承接用户书面裁决（D-003）——F-01～F-14 全部 in-scope、分批实施；批 A/C/D 已关门。本子目标 = **批 B（数据一致性与健壮性，F-05～F-07）**。

## 当前边界

- **范围（本波实施）**：
  - F-05：手写列表端点校验一致性（recycle-bin/wallet 非法分页、wallet pageSize 上限、data-permission policies 真分页或明确不分页、per-task runs 分页参数生效）。
  - F-06：错误码复用误导消息（INVALID_SCOPE/INVALID_WALLET_BODY 细分）、`OPERATION_NOT_FOUND` 进错误目录并带 messageKey。
  - F-07：搜索/排序一致性（通知 q 大小写不敏感、wallet 搜索扩展、recycle-bin 暴露 sort/order、wallet ledger 增加 entry-type 过滤）。
- **非范围**：批 A/C/D（已完成）；不改 Profile 默认集 / 模块矩阵 / Manifest 装配语义（涉端点契约变更时须 go 复核）。

## 成功标准与路线图（P-001）

- [x] **S1 · 方案冻结**：F-05～F-07 具体契约变更与兼容性评估
- [x] **S2 · 实施**：API handler/repository/schema 修改
- [x] **S3 · 测试与回归**：Go 全量 + Web 相关 + tsc（如涉前端）
- [x] **S4 · 自审与关门**：审计 + 台账同步 + goal-tree/workspace 同步

progress: 由四个等权检查点派生（S1～S4）；当前 **4/4**（S1～S4 完成，2026-08-17 关门）。

## 审计策略

| 阶段 / 项 | 默认模式 | 说明 |
|-----------|----------|------|
| S1 冻结 | self | 方案落盘 + 证据核对（来自 D-001 台账） |
| F-05/F-06/F-07 实施 | self（涉契约时升级 independent） | 分页/错误码/搜索契约变化需回归 |
| S4 关门 | self | 常规关门自审 |

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | non-blocking | F-05 data-permission policies 是否改为真分页 | S2 F-05 | S1 | as-built + 方案 | **closed** | — | D-001：内存分页足够 |
| I-002 | non-blocking | F-07 wallet ledger entry-type 取值集合 | S2 F-07 | S1 | as-built + 方案 | **closed** | — | D-001：adjust/freeze/unfreeze/deduct_frozen |

## 父目标

- [GOAL-015-w14-user-perspective-review](../GOAL-015-w14-user-perspective-review/00-meta.md)

## 台账布局

- `01-decision/`：D-NNN；`02-execution/`：E-NNN；`03-audit/`：A-NNN。
- 跨区引用用 Q2 路径（workspace-protocol §2.6）。
