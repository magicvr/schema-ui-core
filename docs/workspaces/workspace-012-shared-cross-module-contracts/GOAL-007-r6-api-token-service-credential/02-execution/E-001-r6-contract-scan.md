---
id: E-001-r6-contract-scan
goal: GOAL-007-r6-api-token-service-credential
status: recorded
created: 2026-08-19
updated: 2026-08-19
parent: GOAL-007-r6-api-token-service-credential
version: 0.1.0
---

# E-001 · R6 认证、权限、审计与协议现状扫描

## 已核对事实

- 用户 access token 是 HS256 JWT，仅携带用户 `Subject` 与 `token_version`；middleware 加载用户、校验 enabled/must-change-password 后把 `account.User` 写入 context。
- refresh token 已采用 256-bit URL-safe opaque raw、SHA-256 hash-only 持久化、过期和原子撤销；可复用生成/比对模式，但不复用用户 session 记录。
- resource permission gate 只读取 `account.User.Permissions`；permission key 来自持久化 roles/permissions 关系，没有第二套 capability grant 解释器。
- operation log 已有 actor/record/detail/correlation 字段与 before/after/diff 脱敏；当前没有 service credential 事件。
- composition 把核心和 Provider routes 汇入共同 handler，R5 operational gate 会覆盖 credential 管理 mutation。
- Profile 默认模块中没有独立 credential 模块；public Manifest 对 token/secret/api-key 类 key fail closed，Web auth client 仅处理用户 JWT/refresh。
- VP-012 将 R6 限定为“机器凭据管理面、作用域、吊销、审计；与用户会话分离”，并要求 Profile/Manifest/协议装配语义不变。

## 信息台账结论

- I-001、I-005、I-006 已验证。
- I-002～I-004 仍为 `collecting`；须由 D-002 精确契约与 cross 设计审计关闭，当前不放行 S1/S2。

## 下一步（计划）

读取即将修改的认证、authsession migration/repository、handler 和 composition 代码，形成 D-002 精确契约并执行 self + independent 设计审计。
