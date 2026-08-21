---
id: A-006-r6-s1-s2-self
goal: GOAL-007-r6-api-token-service-credential
source: self
verdict: pass
status: recorded
created: 2026-08-19
updated: 2026-08-19
parent: GOAL-007-r6-api-token-service-credential
version: 0.1.0
---

# A-006 · R6 S1/S2 实施自审

## 审计头

| 项 | 值 |
|----|----|
| source | self |
| scope | R6 S1/S2：D-003；0044/0045；hash-only repository；service principal；human management API；scope ceiling；create/use/revoke audit；R5 gate；Profile/Manifest/protocol 不变式 |
| verdict | pass |
| required findings | 0 |

## 核对结论

1. secret 为 `sui_sc_` + 32 random bytes Base64URL，raw 只在 create 201 返回；repository 与 API/audit response 均无 raw/hash 暴露，数据库测试只命中 SHA-256 hex 与 15-char prefix。
2. 0044/0045 顺序、旧 operation correlation 保留、新三事件可写和未知事件 CHECK 由迁移/operationlog 测试覆盖；`created_by` 无 users FK，删除用户契约未改变。
3. create/revoke mutation 与 audit 使用同一 Store transaction；forced operationlog failure 的 API 测试证明 credential mutation 回滚。并发大小写重名只成功一次，并发 revoke 只发生一次状态转换和一次 audit callback。
4. service prefix 在 JWT 与 dev fallback 前分流；synthetic id、空 roles 和冻结 permission scopes 注入通用 identity。未知/过期/吊销统一 401；service credential 不回落 dev session。
5. 管理 API 仅接受 human principal；scope 必须存在、属于创建者 effective permissions 且排除两枚 credential management keys。list/detail 无 secret/hash，包含 active/expired/revoked metadata 与固定分页/排序。
6. account/account-self/avatar/MFA/notifications/wallet-self 已切换 `UserIdentityFrom`；permission-gated resource 路径仍使用通用 identity，data-permission self 因 synthetic actor id 不继承 creator。
7. R5 operational gate 在 server 层覆盖新 POST 路由，三种非 normal 模式均返回稳定 503 envelope；Profile/module IDs、Manifest 与协议资产无交付变化。

## Findings

无 required finding。S3 必须由指定 independent provider 复跑关键证据并给出关门 verdict；本条不代替 independent。

## 结论

S1/S2 self `pass`，开放 required=0；路线图完成 3/4，S3 放行至 independent 关门审计。
