---
status: active
created: 2026-08-26
updated: 2026-08-26
parent: GOAL-014-w13-account-lockout-redesign
version: 0.1.0
---

# D-002 · 账号锁定模型重设计方案冻结（2026-08-26）

**目标**：消除 A-001 F-007 定向 DoS 面——"知道用户名的攻击者可用 5 次失败锁定任意账号 15 分钟并吊销其全部刷新令牌"，同时不放弃在线爆破防线与防枚举时序语义。

## 候选模型对比

| 模型 | 防爆破 | 抗定向 DoS | 主要代价 |
|------|--------|-----------|----------|
| A · 纯 per-(account\|IP) 锁定 | 仅每 IP（HTTP 层本就存在） | 最优 | 跨 IP 聚合防线完全消失 |
| B · 全局指数退避（无锁定，next-attempt 上限 ~30s） | ~2880 次/日/账号（bcrypt 下可忽略破解率） | 好（DoS 退化为秒级等待） | 放弃硬锁；与 W4 语义断裂最大 |
| **C · 分层（选定）** | 每 IP 锁 + 高阈值全局熔断 | 好（成本 ×20 且可见） | 需新迁移表 |

## 选定：模型 C · 分层锁定

1. **IP 维度锁**：新表 `login_failures(user_id, ip, fail_count, locked_until, updated_at, PRIMARY KEY(user_id, ip))`。单 IP 对单账号连续 5 次失败 → 该 (user,ip) 对锁定 15 分钟；合法用户自身 IP 不受第三方失败影响。ip 记录为客户端身份字符串（loginClientIP 既有代理信任规则）。
2. **全局熔断（高阈值、低频）**：users 表既有 failed_login_count/locked_until 保留为全局计数器，阈值提升 5 → **100**（24h 滚动窗口内），触发全局锁 15 分钟。定向 DoS 成本 ×20 且触发即经既有 OnLockOpened 钩子产生管理员通知（滥用可见）。
3. **移除失败触发的全量会话吊销**：RevokeAllRefreshTokensForUser 不再由登录失败路径调用（无论 IP 锁还是全局锁）。理由：强制登出是本 finding 中最锐利的武器化面；真实入侵响应走密码修改（token_version 吊销，W4 P0-3 机制保留）与管理员禁用路径。W4"locked account must not keep rotating"的假设（锁=入侵信号）在新模型下不再成立——失败≠入侵，此为有意取舍并提交独立审计复核。
4. **成功清零**：登录成功清除该用户全部 IP 对行 + 全局计数器（Reset 语义扩展）。
5. **契约保持**：ErrAccountLocked / ACCOUNT_LOCKED wire code 不变（语义变为"该来源对该账号处于锁窗"）；missing-user/disabled 的 dummy-bcrypt 防枚举时序不变；HTTP 层每 IP|username 限流器不动。

## 实施要点（S3）

- authsession 迁移新增 login_failures 表（compiled catalog 追加 descriptor）。
- auth.Repository 接口扩展：RecordLoginFailure(userID, ip, ...) / ResetLoginFailures(userID, ip, ...) / LockedUntilFor(userID, ip)；Authenticator.Login 增加 clientIP 参数（调用点：handler/auth.go 登录面传 loginClientIP(r)）。
- 测试：定向 DoS 形状回归锁（他 IP 失败不影响受害者登录）、分布式聚合熔断用例、防枚举时序保持用例、迁移升级用例。

## 未选方案留痕

- 用户裁决仅固定"fixed"；模型 A/B/C 为编排器工程裁量，最终接受以 S5 independent 审计 + S6 用户关门为准。
