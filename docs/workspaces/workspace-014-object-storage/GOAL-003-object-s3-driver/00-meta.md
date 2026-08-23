---
id: GOAL-003-object-s3-driver
title: R2 S3 兼容接入（驱动 + readyz 扩依赖）
status: done
parent: GOAL-001-object-storage
created: 2026-08-21
updated: 2026-08-21
version: 0.2.0
progress: 3/3
plan_refs:
  - VP-014-object-storage
primary_plan: VP-014-object-storage
---

# GOAL-003 · R2 S3 兼容接入

## 意图

承载 Root [GOAL-001](../GOAL-001-object-storage/00-meta.md) 纲领路线图 **R2**：按 D-004/D-005 交付 S3 兼容适配器（aws-sdk-go-v2，API 子集见 Root D-004），并在显式配置 driver=s3 时把 readyz 扩展到对象存储后端。前置信息项 I-001 / I-003 已由 Root D-004 / D-005 关闭。

## 范围

- `internal/objectstore/s3.go`：实现冻结的 kernel.ObjectStore 端口（key=`<ns>/<id>`，meta=user metadata，NoSuchKey→ErrObjectNotFound）。
- go.mod：引入 aws-sdk-go-v2 service/s3 及其依赖模块。
- readyz：driver=s3 显式配置时经 HeadBucket 探针扩依赖；local/缺省语义零变化。
- 测试：离线 stub 合同测试 + 可选 live 集成测试（S3_TEST_* 未设则干净 skip，沿 pgtest 先例）。
- 不含：三类落盘调用方改线（R3）、公共面去 os.File（R4）。

## 检查点

1. s3 适配器编译并满足端口合同（stub 测试覆盖 put/get/stat/delete/exists/list + 错误映射 + key 构造 fail-closed）。
2. readyz 扩依赖生效且 local 缺省不变（composition 接线证据）。
3. 全量 go build / go test 绿。

## 父目标

- GOAL-001-object-storage（R2 阶段；VP-014 架构 A2）
