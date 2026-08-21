---
id: D-001-s3-adapter-freeze
title: R2 S3 适配器方案冻结
status: accepted
created: 2026-08-21
updated: 2026-08-21
parent: GOAL-003-object-s3-driver
version: 0.1.0
---

# D-001 · R2 S3 适配器方案冻结

依据 Root D-004（SDK 与 API 子集）/ D-005（凭证注入）。R1 冻结的端口合同（GOAL-002 D-001）**不改一字**。

## 实现要点

1. **客户端**：`config.LoadDefaultConfig` 仅注入 static credentials + region；`s3.NewFromConfig` 带 `BaseEndpoint` 与 `UsePathStyle`（cfg.ObjectsS3UsePathStyle，缺省 true）。不触碰默认凭证链/共享配置文件（D-005）。
2. **方法映射**：
   - Put → PutObject（body 整块 bytes；meta → Metadata map，空值键省略）
   - Get → GetObject 读全量 body；NoSuchKey/`NotFound`/`StatusCode 404` → `kernel.ErrObjectNotFound`
   - Stat → HeadObject（ContentLength + Metadata）；404 → 哨兵
   - Delete → DeleteObject；删除不存在的 key 在 S3 上本就幂等成功 → nil
   - Exists → HeadObject 判 404
   - List → ListObjectsV2（prefix=<ns>/，分页聚合），从 key TrimPrefix 得 id，升序返回
3. **Ping(ctx)**：HeadBucket——readyz 探针用；以 `interface{ Ping(context.Context) error }` 断言消费，**不给端口加方法**（保 R1 冻结面稳定）。
4. **校验复用**：namespace/id 校验与 LocalStore 同源（kernel.ValidObjectNamespace/ValidObjectID），fail-closed 先于任何网络调用。
5. **超时**：操作沿用调用方 ctx；readyz 探针由既有 1s readyz ctx 约束。

## 测试策略

- **离线**：对内部 `s3API` 接口打 stub，覆盖全部端口行为 + 错误映射 + 校验拒绝（无网络）。
- **Live（可选）**：`S3_TEST_ENDPOINT/S3_TEST_BUCKET/S3_TEST_ACCESS_KEY/S3_TEST_SECRET` 全设时跑真实 round-trip（MinIO 即可）；未设干净 skip（pgtest 先例）。CI 不依赖 live。

## 边界

- 不建桶、不管 lifecycle/policy/CORS；桶缺失 = Ping 报错 = readyz unavailable（运维责任，符合"显式配置才扩依赖"）。
- 分片上传不在 API 子集（VP 非目标 RT-S04）。
