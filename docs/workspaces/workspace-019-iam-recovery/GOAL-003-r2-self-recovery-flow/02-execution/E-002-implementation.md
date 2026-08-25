---
id: E-002
doc: execution-entry
goal: GOAL-003-r2-self-recovery-flow
status: recorded
created: 2026-08-25
updated: 2026-08-25
version: 1.0.0
---

# E-002 · R2 实施切片落地（C2/C3/C4 满）

2026-08-25 完成（git：`299f8f52` 后端 + `9628ca8f` Web）：

- **C2**：迁移 0056 `password_recovery_challenges`（双方言成对体；checksum `e19db1a2…` 入冻结台账）；identity.go catalog head 55→56、两份表清单补行；identity_test/migrate_test/operations_test/restart_test 四处黄金断言同步；store 包测试全绿。
- **C3**：
  - 域逻辑 `authsession/recovery.go`：ResolveRecoveryTarget（用户名精确 → verified lower(email)）、StartRecovery 两阶段派发+失败补偿、EvaluateRecoveryCode 只读分类、ConsumeRecoveryAttempt（错码与第二因子失败共用预算，≤5 作废）、CompleteRecovery（先消费挑战再走 UpdateUser）。
  - HTTP `handler/recovery.go`：`POST /api/auth/recovery/start|complete` 中央公开注册；start 无路径账号同形静默（防枚举）；complete 统一 RECOVERY_CODE_INVALID 口径 + MFA 门序（码匹配→第二因子→密码基线）；loginRateLimiter IP|identifier 桶；错误码 4 个入 errorcatalog 并加入契约测试冻结集。
  - MFA 门 `mfa.Service.VerifySecondFactor`（requireActiveSecondFactor 薄封装，nil fail-closed）；operational allowlist 扩 recovery 两路径；composition 以 true-nil 接口装配。
  - 测试：authsession 域 6 组 + handler 面 4 组全绿。
- **C4**：`recovery_e2e_test.go` 经 mock OutboxSink 取码完成闭环；LoginPage 两步恢复流（第二因子字段按服务端要求展示、过期回步骤一保账号值）；zh/en i18n 各 19 键；web tsc 干净 + vitest **1105/1105** 绿（顺修 W25 守卫既有红测：custom-components.schema.test.ts 漏 import email-identity 组件）。
- 未改 Profile 默认集 / 迁移 checksum / 既有 API 契约。

后续：C5 —— grok build（grok-4.6 · high）independent 审计进行中；意见落盘后由编排器合并响应，self 关门审随后。
