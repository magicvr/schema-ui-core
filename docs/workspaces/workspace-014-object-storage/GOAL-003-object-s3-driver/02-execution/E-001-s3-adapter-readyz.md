---
id: E-001-s3-adapter-readyz
title: R2 实施——S3 适配器 + readyz 扩依赖
status: recorded
created: 2026-08-21
updated: 2026-08-21
parent: GOAL-003-object-s3-driver
version: 0.1.0
---

# E-001 · R2 实施：S3 适配器 + readyz 扩依赖

## 事实

按 [D-001](../01-decision/D-001-s3-adapter-freeze.md)（依据 Root D-004/D-005）实施：

1. **适配器**：`apps/api/internal/objectstore/s3.go` —— `S3Store` 实现冻结端口；key=`<ns>/<id>`；meta 走 S3 user metadata（空值省略）；NoSuchKey/NotFound/404 → `kernel.ErrObjectNotFound`；`NewS3` 用 static credentials 显式构造（不触碰默认链）；namespace/id 校验先于任何网络调用。
2. **Ping**：HeadBucket 探针，经可选能力断言消费——**端口合同零改动**。
3. **readyz 扩依赖**：`handler.RegisterWithMFAProbes` + `readyz` 额外探针参数（nil 忽略）；composition 在 `driver=s3` 显式配置时构造适配器并传入 Ping 探针；local 缺省路径零变化。
4. **依赖**：go.mod 引入 aws-sdk-go-v2（service/s3 v1.107.3 等）。
5. **测试**：
   - 离线 stub 合同测试 `s3_test.go`：round-trip/upsert、NotFound 哨兵、校验拒绝且**零后端调用**、List 分页聚合+跨命名空间隔离+缺失命名空间空列表、Ping、构造器拒绝空桶。
   - 可选 live 集成 `s3_live_test.go`：S3_TEST_* 全设才跑真实 MinIO round-trip，否则干净 skip。
   - readyz 探针测试 `health_probes_test.go`：探针通过=200 / 失败=503 / nil 条目忽略。

## 验证证据

- go build ./... exit 0；go vet exit 0。
- go test ./internal/{objectstore,handler,composition}/ 全绿（handler 127s / composition 16s）。
- 全量 go test ./... 见 E-002 补记。

## 边界

三类落盘调用方仍未接线（R3）；main.go 的 driver=s3 启动警告（A-002 R-002）保持有效。
