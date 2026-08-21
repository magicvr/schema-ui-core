---
id: GOAL-004-object-families-migration
title: R3 三类落盘收口走端口
status: active
parent: GOAL-001-object-storage
created: 2026-08-21
updated: 2026-08-21
version: 0.1.0
plan_refs:
  - VP-014-object-storage
primary_plan: VP-014-object-storage
---

# GOAL-004 · R3 三类落盘收口走端口

## 意图

承载 Root [GOAL-001](../GOAL-001-object-storage/00-meta.md) 纲领路线图 **R3**：avatars / brand-assets / uploads（含 file-library 与 data-transfer 共享上传目录）全部改走同一 `kernel.ObjectStore` 端口；composition 构造唯一适配器实例（缺省本地盘 / 显式 s3）。

## 范围

- handler 内部：`uploadStore`、`RasterAssetStore` 持久化改端口委托；HTTP 面形状不变。
- 模块公共契约：filelibrary / datatransfer 的 `uploadDir string` 参数改 `kernel.ObjectStore`（VP 意图：模块只消费端口）。
- composition：单一适配器实例注入全部消费方；main.go 的"s3 未接线"警告随接线完成移除。
- 行为保持：GC / 配额 / owner 门禁 / 遗留磁盘兼容逐点对应（见 D-001）。

## 检查点

1. 三类落盘持久化全部经端口（grep 无直接 os.ReadDir/WriteFile 于三类路径）。
2. 全量测试绿（handler/filelibrary/datatransfer/composition 套件原语义通过）。
3. driver=s3 时三类落盘真实经 S3Store（接线锁测试）。

## 父目标

- GOAL-001-object-storage（R3 阶段；VP-014 架构 A2）
