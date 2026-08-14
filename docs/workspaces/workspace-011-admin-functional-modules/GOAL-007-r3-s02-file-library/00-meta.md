---
id: GOAL-007-r3-s02-file-library
title: R3-S02 · 文件/附件库（统一文件管理、引用、清理；复用 upload 基建）
status: done
parent: GOAL-001-admin-functional-modules
created: 2026-08-14
updated: 2026-08-14
version: 0.5.0
progress: 5/5
---

# GOAL-007-r3-s02-file-library · 文件/附件库

## 概述

常用档 S-02（I-011-001 §4）：新建**共享能力模块**（`admin.file-library` 候选名），提供统一文件/附件库——文件列表（schema 驱动代表页）、上传（复用 C-09 form-with-upload 基建与 VP-009 授权/配额/所有权加固）、下载、删除、引用与清理。基架现仅控件级上传（C-09），无统一文件管理面。

## 当前边界

- 统一文件面：文件列表 CRUD 代表页 + 上传/下载/删除 + 元数据（名称/类型/大小/上传者/时间）
- 授权/配额/所有权加固复用 VP-009 上传经验（授权键、配额、文件所有权校验）
- 引用与清理：被引用文件的删除保护与清理策略（S1 方案冻结时定）
- Profile 归属：进入 admin 默认集（S1 方案确认；Profile 内容扩展，不改装配语义——F-01 先例）

## 成功标准与路线图（P-001）

- [x] **S1 · 方案冻结**：边界（引用/清理策略、配额语义）、权限键、审计设计、协议对照（上传/文件面）、Profile 归属；方案级 self 审视（D-001/D-002/A-001，2026-08-14）
- [x] **S2 · 实现**：模块 provider + 迁移 + schema 代表页 + 文件端点（列表/上传/下载/删除）+ 前端（页面/上传交互）+ 测试（E-002）
- [x] **S3 · 验证**：单元/集成 + 代表场景实测 + 全量回归（go test / web suite / e2e 双 profile 8/8；冒烟留 S5，E-003）
- [x] **S4 · go 影响判定 + 自审**：go 影响判定（D-003 内容扩展不触发失效，不暂挂）+ self 审计（A-002 pass）
- [x] **S5 · 关门**：独立审计（A-003 grok pass，0 required）+ 4 条 recommended 全 fixed（D-004）+ 关门（E-004）+ goal-tree 同步

progress: 0/5 由五个等权检查点派生（S1 完成后更新）。

## 审计策略

文件库涉上传授权/配额/所有权 → security/data 高影响门禁 → 独立审计（grok build；P-004 在关门门禁确认 provider）。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 |
|----|------|-----------------|----------|--------------|-----------------|------|
| I-001 | required | 文件面协议契约（v2.8.0 upload 相关 + node.schema.json 扩展动作键）与呈现自由边界 | S1 方案 | 对照 protocol-inventory + node.schema.json | **closed**（D-002 §1） |
| I-002 | required | 引用/清理策略与配额语义（删除保护、孤儿清理、配额上限） | S1 方案 | 复用 VP-009 W2/W4 上传加固 + R2 通知保留策略先例 | **closed**（D-002 §6） |
| I-003 | required | Profile 归属：S-02 进入 admin 默认集？mvp 保持精简？ | S1 方案 | F-01 先例（Profile 内容扩展 + adminFunctionalOrder） | **closed**（D-001 §2） |

## 父目标

- [GOAL-001-admin-functional-modules](../GOAL-001-admin-functional-modules/00-meta.md)

## 台账布局

本目标从首条记录起使用 `01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger；索引与目录条目共同构成正式记录。
