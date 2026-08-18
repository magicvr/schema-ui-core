---
id: D-001-r6-initial-boundary
goal: GOAL-007-r6-api-token-service-credential
status: proposed
created: 2026-08-19
updated: 2026-08-19
parent: GOAL-007-r6-api-token-service-credential
version: 0.1.0
---

# D-001 · R6 机器凭据初始边界

## 已确认方向

1. 机器凭据与用户 JWT/refresh session 分离；认证后投影独立 principal，不借用用户 `token_version` 撤销语义。
2. 复用现有 256-bit opaque token 生成和 SHA-256 hash-only 持久化模式；raw secret 只在创建响应返回一次。
3. scope 使用现有 permission key，不建立第二套授权命名空间；管理端和调用端均 fail closed。
4. 管理面和认证面保留 request-id/error envelope/R5 operational gate；审计 detail 不包含 raw secret、hash 或 Authorization header。
5. 实现归入现有 `core.auth-session` 候选能力，避免新增默认模块、Profile 或 Manifest contribution。

## 待冻结契约

- credential principal 的 context 形状，以及与现有 `account.User` permission gate 的组合方式。
- 创建者可授予 scope 的权限上限、系统管理员边界和未知 permission 的处理。
- secret 前缀、最大 TTL、默认 TTL、过期/吊销/重复吊销与 last-used 写入语义。
- 管理路由、稳定错误码、ETag/expectedVersion 是否适用，以及审计事件集合。
- Bearer 解析如何无歧义区分 JWT 与 service credential，并保持用户认证错误非回归。

## 未选方案

- 把 service credential 签成用户 JWT：混合用户会话与机器生命周期，不能独立吊销或表达机器主体。
- 明文或可逆加密保存 secret：扩大泄露面，且不需要再次展示原文。
- 新增默认 Profile/module：会改变 VP-012 明确要求保持不变的装配语义。
- 在 Web 保存或自动使用 service credential：机器凭据不应成为浏览器用户会话。

## 门禁

I-002～I-004 须由精确契约与 cross 设计审计关闭；D-001 在此之前保持 `proposed`，不得进入 S1/S2 实施。
