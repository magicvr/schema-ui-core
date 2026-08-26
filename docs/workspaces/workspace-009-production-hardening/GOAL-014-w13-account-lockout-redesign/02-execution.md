---
status: active
created: 2026-08-26
updated: 2026-08-26
parent: GOAL-014-w13-account-lockout-redesign
version: 0.1.0
---

# 执行索引 · GOAL-014

## 时间线（事实）

### E-001 · 立项（2026-08-26，S1 完成）

1. 用户就 GOAL-013 A-001 F-007 完成三路径裁决：**fixed**，并指定以"当前目标的下级子目标"承载治理上下文。
2. 建立五件套（本目录）；D-001 落盘裁决；路线图 S1–S6 登记（progress 来源）。
3. goal-tree 同步：GOAL-013 树下新增本目标节点与状态表行。

**路线图状态**：S1 ✅；下一阶段 S2（锁定模型方案设计与冻结）。

### E-002 · S3 分层锁定模型实施完成（2026-08-26）

1. **迁移 0061 `login_failures`**（authsession compiled catalog；sqlite/postgres 成对 DDL，INTEGER/BIGINT 按 0055 先例）：`login_failures(user_id, ip, fail_count, locked_until, updated_at, PK(user_id,ip))` + `users.last_login_failure_at`（全局计数 24h 滑动重启载体）。
2. **authsession 仓储**（新文件 `accounts_lock_source.go`）：`RecordLoginFailureFor`（原子 UPSERT+滑动窗：计数器上次移动早于锁窗则重置为 1）、`LoginLockedFor`（存储错误 fail-closed）、`ResetLoginFailuresFor`（成功登录清除该账号全部来源行）。既有全局方法 `RecordLoginFailure`/`ResetLoginFailures` 接入 last_login_failure_at 滑动语义。`UnlockUser` 同步清除来源行（管理员解锁必须解除全部锁定面——首轮全量回归暴露的实现缺口，已修）。
3. **auth 核心**（`internal/auth/auth.go`）：`Login` 增加可选 clientIP（variadic，缺省落 "-" 单桶以保持既有测试语义）；校验顺序 = 用户存在（dummy burn）→ 来源对锁 → 全局锁 → disabled → 密码。失败路径双记账（来源桶阈值 `IPSourceLockThreshold=5` / 全局天花板 `LockThresholdFailures=100`，24h 滑动）；**移除失败触发的 RevokeAllRefreshTokensForUser**；OnLockOpened 仅在全局熔断开启时触发（来源锁静默——防止攻击者刷受害者的通知中心）。生产调用点 `handler/auth.go` 传入 `loginClientIP(r)`。
4. **测试**：新增 `lockout_source_test.go` 三组缺陷形状回归锁（来源隔离：他源失败不影响本源登录；全局天花板分布式制动 + OnLockOpened 恰一次 + 锁窗自愈；失败不再吊销刷新令牌——双令牌对照）。既有四条旧契约测试按 D-002 重写/修正：TestAccountLockLifecycle（过期断言改走来源对锁）、TestAccountLockRevokesSessions→TestAccountSourceLockKeepsSessions（不吊销为新契约）、两个通知测试改驱动全局路径。
5. **验证**：go vet ./... 0 输出；受影响包（handler/auth/authsession）全绿后全量复跑确认。**补充（同日）**：真实 Postgres 方言复核——`TestAuthsessionPostgresApplyIntegration`、`TestFullCatalogPostgresBootstrapIntegration` 等 PG 门控集成测试全绿（0061 成对 DDL 在 PG 上实际应用）；另对 `RecordLoginFailureFor` 首插并发竞态补兜底（唯一冲突→新事务重试，`cf5675f1`）。

**路线图状态**：S1 ✅ S2 ✅ S3 ✅；下一阶段 S4 全量回归 + checkpoint。

## 执行记录目录

| 编号 | 文件 | 内容 | 状态 |
|------|------|------|------|
| E-001 | （本文件时间线第 1 节） | 立项 + D-001 落盘 + 五件套建立 | done 2026-08-26 |
| E-002 | （本文件时间线第 2 节） | S3 分层锁定模型实施（迁移 0061 + auth 改造 + 回归锁 ×3 + 旧契约测试修正） | done 2026-08-26 · checkpoint `26655b55` |
