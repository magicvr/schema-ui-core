---
id: D-004-s3-driver-subset
title: S3 驱动与 API 子集裁定（I-001 闭合）
status: accepted
created: 2026-08-21
updated: 2026-08-21
parent: GOAL-001-object-storage
version: 0.1.0
---

# D-004 · S3 驱动与 API 子集（闭合 I-001）

## 决定

- **SDK**：AWS SDK Go v2 的 S3 客户端（`github.com/aws/aws-sdk-go-v2/service/s3`）。理由：MinIO / R2 / AWS 三者的公约数客户端——三者均以 S3 兼容 API 为第一接口；SDK 以 `BaseEndpoint` + `UsePathStyle` 覆盖自建端点场景。不引入 minio-go 等第二客户端（避免同一 VP 内两套 S3 表面）。
- **API 子集冻结**（超出即越界）：`PutObject / GetObject / DeleteObject / HeadObject / ListObjectsV2`，外加 readyz 探针用 `HeadBucket`。错误映射只认 `NoSuchKey`/`NotFound`（→ `kernel.ErrObjectNotFound`）。
- **Key 方案**：单桶内 `<namespace>/<id>`（Root D-002）；桶必须预先存在——适配器**不**建桶、不配 lifecycle/CORS/policy。
- **Meta 表示**：S3 user metadata（x-amz-meta-name/type/kind/owner），随 body 原子写入——无本地边车的双对象部分失败窗口（对照 A-001 N-001）。List 只需 id，不需要 meta。
- **明确不用**：multipart、presigned URL、分页 token 以外的高级列举特性、IAM role / IMDS 凭证链（见 D-005）。

## 未选方案

- **minio-go**：更轻，但属第二方言面，且 AWS/R2 兼容性由第三方库转译背书，弱于 SDK 直连。
- **手写 REST 客户端**：SigV4 正确性与维护成本不自担。

## 影响

R2 交付 `internal/objectstore/s3.go`；go.mod 引入 aws-sdk-go-v2 相关模块（存储 SDK，非对象存储"方言"，VP 边界不受影响）。
