---
id: A-001-r6-s0-design-self
goal: GOAL-007-r6-api-token-service-credential
source: self
verdict: pass
status: recorded
created: 2026-08-19
updated: 2026-08-19
parent: GOAL-007-r6-api-token-service-credential
version: 0.1.0
---

# A-001 · R6 S0 设计自审

## 审计头

| 项 | 值 |
|----|----|
| source | self |
| scope | R6 S0：D-001/D-002、I-001～I-006；机器 principal、生命周期、权限上限、审计与装配不变式 |
| verdict | pass |
| required findings | 0 |

## 核对结论

1. `sui_sc_` 前缀使 service credential 与 JWT 可确定性分流；raw 一次性返回、hash-only 存储与现有 refresh token 安全模式一致，同时生命周期完全独立。
2. hidden principal kind + `UserIdentityFrom` 保留现有 permission-gated handler 兼容性，并把账户/MFA/通知/wallet self-service 从机器主体显式隔离。
3. scope 同时受 catalog existence、创建者 effective permission 子集与 reserved management permission 三重约束；service credential 不能自举管理凭据。
4. 管理路由、system-data permissions 与 repository 均附着既有 `core.auth-session`，不新增 Profile/module/page/navigation/fragment，Manifest 与协议边界可保持不变。
5. guarded revoke、唯一约束和必填有界 expiry 形成可核对的并发与生命周期门禁；首版无 PATCH，因此不需要伪造 ETag 语义。
6. create/use/revoke 审计字段排除 secret/hash/header；correlation 与既有 operation log 可复用，best-effort 失败不改变认证结果。

## Findings

| ID | 等级 | finding | disposition |
|----|------|---------|-------------|
| F-001 | recommended | S1 repository tests 必须直接查询数据库，证明 raw secret 未落盘、scope 排序稳定且 guarded revoke 并发只有一个状态转换。 | implementation gate：S1 |
| F-002 | recommended | S2 必须覆盖所有 user-only identity consumer，防止 service credential 误入 self-service handler。 | implementation gate：S2 |
| F-003 | recommended | S3 composition 测试必须固定 Profile/module IDs 与 Manifest bytes 不变，并验证 operational gate 对 credential mutation 生效。 | close-out gate：S3 |

## 结论

Self 设计审计通过，开放 required=0；cross 模式尚缺 independent，I-002～I-004 与 D-002 不在本条中提前关闭。
