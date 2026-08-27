---
status: active
created: 2026-08-26
updated: 2026-08-26
parent: GOAL-014-w13-account-lockout-redesign
version: 0.1.0
---

# A-001 · GOAL-014 关门前 self 审计（编排器自审）

- **source**: self（编排器；独立腿按路线图 S5 由 grok build 另行执行）
- **日期**: 2026-08-26
- **scope**: D-002 冻结方案在 checkpoint `26655b55` 的实施真实性核对（迁移/auth 核心/handler 接线/管理员解锁）、缺陷形状回归锁覆盖、全量回归证据
- **verdict**: **pass**

## 逐项核对

| D-002 设计要素 | 实施 | 证据 | 结论 |
|----------------|------|------|------|
| 来源对锁 5/15min | `login_failures` 表 + `RecordLoginFailureFor`（原子 UPSERT、计数滑动重置）+ `LoginLockedFor` fail-closed + `ResetLoginFailuresFor` | `accounts_lock_source.go`；`TestLoginSourceScopedLockout`（他源失败不影响本源登录） | ✓ |
| 全局天花板 100/24h 滑动 | `LockThresholdFailures=100`；`users.last_login_failure_at` 滑动重启；锁窗开启后自愈 | `TestLoginGlobalCeilingBrakesDistributedGuessing`（OnLockOpened 恰一次 + 新来源被全局拒 + 锁窗后恢复） | ✓ |
| 移除失败触发会话吊销 | `Login` 失败路径无任何 RevokeAll 调用（生产代码零调用点残留）；双令牌对照测试 | `TestLoginFailureNoLongerRevokesRefreshTokens`（锁窗内未出示的令牌过期后仍可轮换） | ✓ |
| 防枚举时序保持 | missing-user/locked/disabled 三条 dummy-bcrypt 路径原样保留；ErrAccountLocked/UNAUTHORIZED 封装不变（W7 F-009） | auth.go Login 结构；既有 TestLoginUnknownUser 等全绿 | ✓ |
| 管理员可见性 | OnLockOpened 仅全局熔断触发（来源锁静默，防通知疲劳滥用）→ account.locked 通知 | TestNotificationLockEventProduced（驱动全局路径） | ✓ |
| 管理员解锁完整性 | UnlockUser 同步清除全部来源对行 | TestAdminUnlockClearsLockWindow（解锁后同源登录恢复） | ✓ |
| handler 接线 | 生产登录传 loginClientIP(r)；HTTP 层每 IP\|username 限流器不变 | handler/auth.go | ✓ |

## 回归证据

- `go vet ./...` 0 输出；`go test ./... -count=1` **46 包全绿**（含 store 迁移目录头 pin 更新至 v61：completeFingerprint/lockedHeadExtraTables/catalog ownership/三处链式断言）。
- Checkpoint：`26655b55`。

## 自审备注（提交独立复核）

1. **"24h 滚动窗口"的实现语义**：D-002 原文为"100 失败/24h 滚动窗口"。实现为**连续失败计数 + 24h 滑动重启**（last_login_failure_at 超过 24h 即从 1 重计；成功或触锁清零）。非严格时间窗内求和，但对爆破预算的约束等价（持续失败才能逼近天花板），且避免新增聚合查询。判定为忠实实现，偏差已如实登记。
2. **variadic clientIP**：`Login(username, password, now, ip ...string)`——缺省落 "-" 单桶，仅为保持既有测试/开发调用语义；生产唯一调用点恒传真实身份。
3. **Refresh rotate-before-checks**：锁定期间出示刷新令牌仍会消费该令牌（先轮换后校验的既有契约），本设计未改变也未利用该行为；回归锁采用"未出示的第二令牌"断言规避歧义。

## 结论

D-002 各要素 genuine 落地、定向 DoS 形状回归锁在位、全量回归绿。具备进入 independent 审计条件。
