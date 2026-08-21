---
id: A-001-self-r2-s3
title: 自审 · R2 S3 兼容接入
source: self
date: 2026-08-21
scope: GOAL-003 全部交付（S3 适配器 / readyz 扩依赖 / 测试 / 依赖引入）
verdict: pass
parent: GOAL-003-object-s3-driver
version: 0.1.0
---

# A-001 · 自审：R2 S3 兼容接入（verdict: pass）

## 合同与证据指回

- **端口零改动**：`kernel/objectstore.go` 无 diff（commit 1545134 未触碰）；Ping 经可选能力断言消费。
- **语义等价**：与 LocalStore 对照——Put upsert / Delete 幂等 / Get·Stat 缺失→ErrObjectNotFound / List 升序+缺失命名空间空切片 / 校验先于 IO（stub 断言零后端调用）。stub 测试与 local_test.go 用例一一对应。
- **API 子集**：仅 PutObject/GetObject/HeadObject/DeleteObject/ListObjectsV2+HeadBucket；无 multipart/presigned/建桶/lifecycle（grep 可核）。
- **凭证**：static-only 构造；错误串不带 key 值；user metadata 不含明文 secret（secret 只用于 SigV4，不落 metadata——meta 四键为 name/type/kind/owner）。
- **readyz**：仅 driver=s3 显式构造探针；nil 条目忽略；local 缺省行为不变（composition 仅在 s3 分支构造）。

## Findings

| 编号 | 级别 | 内容 | 处置 |
|------|------|------|------|
| N-101 | note | ListObjectsV2 分页以 ContinuationToken 手工循环（不用 paginator），token=上一页最后 key 的 fake 语义与 S3 真实 token 不完全同形——真实后端只透传 token，适配器不解析其内容，故无影响。 | 无需动作 |
| N-102 | note | live 集成测试仅在 S3_TEST_* 全设时运行；CI 默认离线。真实 MinIO round-trip 证据留待 R5 双路径验收时补强（属 R5 分母）。 | R5 补 |
| — | required | 无 | 开放 required = 0 |

## 结论

R2 自审 pass，开放必改 0。已并行发起独立审计（grok build · grok-4.6 · high），意见落盘后由编排器响应。
