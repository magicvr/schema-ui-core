---
id: GOAL-002-object-port-freeze
title: R1 内核对象存储端口与配置面冻结
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

# GOAL-002 · R1 内核对象存储端口与配置面冻结

## 意图

承载 Root [GOAL-001](../GOAL-001-object-storage/00-meta.md) 纲领路线图 **R1**：冻结内核对象存储端口公共类型（无本地路径 / `os.File`）与配置面（缺省本地盘、S3 显式配置），交付本地盘适配器为缺省实现。前置信息项 I-002 / I-005 已由 Root D-002 / D-003 关闭。

## 范围

- `internal/kernel`：`ObjectStore` 端口类型 + 命名空间 / id 校验规则（合同冻结）。
- `internal/objectstore`：本地盘适配器（缺省实现，磁盘布局与现状逐点兼容）。
- `internal/config`：`storage.objects` 配置面（driver / local.root / s3.*）+ fail-closed 校验。
- 不含：S3 驱动接入（R2）、三类落盘调用方改线（R3）、公共面去 `os.File`（R4）、readyz 扩依赖（R2，仅 S3 显式配置后）。

## 检查点

1. 端口类型冻结并编译（kernel.ObjectStore + 校验规则）。
2. 本地适配器 round-trip 证据（put/get/stat/list/delete/exists + 遗留 meta 兼容 + fail-closed 校验）。
3. 配置面加载与校验证据（缺省 local；driver=s3 缺凭证 fail-closed；s3 键误配 driver=local fail-closed）。
4. 全量 `go build ./...` + `go test ./...` 绿。

## 父目标

- GOAL-001-object-storage（R1 阶段；VP-014 架构 A2）
