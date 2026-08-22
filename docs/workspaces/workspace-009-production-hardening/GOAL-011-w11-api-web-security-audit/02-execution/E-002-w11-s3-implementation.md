---
id: E-002-w11-s3-implementation
goal: GOAL-011-w11-api-web-security-audit
status: done
created: 2026-08-22
updated: 2026-08-22
parent: GOAL-011-w11-api-web-security-audit
version: 0.1.0
---

# E-002 · W11 S3 实施：6 条 required 修复（2026-08-22）

## 事实

对照 [D-002](../01-decision/D-002-w11-scope-and-go-hold.md) 整单采纳范围实施。逐条修复均保留原失败语义的 fail-closed 方向，未改任何 wire 错误码。

### F-001 · Postgres `CreateUserManagement` EXISTS 扫进 int（HIGH · fixed）

- `apps/api/internal/modules/authsession/users_repository.go` `CreateUserManagement`：`var exists int` → `var exists bool`，`exists == 1` → `exists`。与同文件 `DeleteUser` / `scanUserListRow`（R6 已改 bool）一致——postgres 出 native bool，sqlite 出 0/1，int 目标使创建/CSV 导入在 postgres 500。
- 调用方 `usersEntity.Create` 与 `importUsersCSV` 共用此路径，一改全修。
- 回归锁：`internal/store/postgres_test.go` 新增 `TestPostgresCreateUserManagementExistsScan`（PG_TEST_* 门控）：postgres 上创建成功 + 重复用户名 `ErrUsernameTaken` + 读回一致。

### F-002 · 删除成功但回收站快照失败仍返 204（HIGH · fixed）

- 事务化 seam（handler + 两个模块）：
  - `handler/resources.go`：新增 `TrashTxRecorder`（`RecordTx(ctx, tx, …)`）与 `TrashTxDeleter`（`DeleteTrashTx(ctx, id, user, now, record func(ctx, tx) error)`）两个可选接口；`delete()` 与 `batchDelete()` 顺序分支在**双方都实现**时走同一事务路径，快照失败 → 删除整体回滚 → 失败响应；未实现的实体保留原 delete-then-snapshot 语义（字节一致）。
  - `modules/recyclebin/store/repository.go`：`RecordTx(ctx, tx, item)`；`Record` 复用。
  - `modules/recyclebin/service.go`：`RecordTx`（满足 `handler.TrashTxRecorder`）。
  - `modules/datadictionary/store/repository.go`：`DeleteTypeTx(ctx, id, record)` / `DeleteEntryTx(ctx, id, record)`；原 `DeleteType`/`DeleteEntry` 转为薄包装（record=nil）。
  - `modules/scheduledtasks/store/repository.go`：`DeleteTaskTx(ctx, id, record)`；`DeleteTask` 薄包装。
  - `handler/dictionary.go` / `handler/scheduledtasks.go`：dictType/dictEntry/task 三个实体实现 `DeleteTrashTx`（事务成功后记 audit；回滚不记 audit；对真实仓库做类型断言，断言失败 fail-closed 报错而非静默降级）。
- 回归锁：
  - `handler/recyclebin_test.go` 新增 `txRecordingTrash`（实现双接口、可注入失败）+ `TestRecycleFactoryDeleteAndSnapshotSameTx`（204 + dict 行消失 + recycle_items=1 同事务提交）+ `TestRecycleFactorySnapshotFailureRollsBackDelete`（500 + dict 行仍在 + recycle_items=0 回滚证明）。
  - 既有 `recordingTrash`（仅 `Record`）驱动的 factory 测试继续走 legacy 路径，语义不变。

### F-003 · MFA 第二因子可在线穷举（HIGH · fixed）

- `handler/auth.go` login MFA 分支：签发 proof 前同一 IP|username 限流桶 `allow` 检查（超限 429 + Retry-After），签发成功后 `record`——「密码已过、等待 MFA」状态被限流，每 15 分钟每个桶最多 20 个 proof（每 proof 5 次猜测）。
- `modules/mfa/store/repository.go`：
  - `CreateProof` 在同一事务内先清理该用户已过期 proof 行（懒清理，captcha 先例）——proof 数不能无界增长。
  - `IncrementProofFailures` 加 `AND fail_count < ?`（上限 5，store 侧常量 `proofFailLimit`）——并发失败猜测不再能通过 check-then-act 突破穷举预算。
- 回归锁：`modules/mfa/store/repository_test.go` `TestCreateProofPurgesExpired` + `TestIncrementProofFailuresCapped`。

### F-004 · JWT 轮换使已存 TOTP 密文不可解（MEDIUM · fixed）

- `modules/mfa/service.go`：`NewService(repo, serverSecret, previousSecret []byte)`——previous 非空时派生前一把 HKDF 密钥，进入轮换窗口；`decryptSecret` 先试当前密钥、失败再试 previous；`maybeRewrap` 在任何**成功**第二因子验证（verify / confirm / disable / recovery rotate，仅 CAS 赢家）后用当前密钥重封密文。previous 为空时行为与旧版逐字节一致。
- `modules/mfa/store/repository.go`：新增 `UpdateSecretCiphertext(userID, ciphertext, now)`。
- `composition/composition.go`：`NewService(..., []byte(secret), []byte(cfg.AuthJWTSecretPrevious))`——与 VP-016 JWT 双密钥同构；无需新配置面。
- 未选方案记录（见 D-002）：独立 MFA 专用密钥需新增部署配置面与生命周期管理，本波不采用。
- 回归锁：`modules/mfa/service_test.go` `TestServiceRotationWindow`——A 下注册 → 轮换为 current=B/previous=A 可登录 → 成功后重封 → 仅 current=B（previous 丢弃）仍可登录 → 仅 A 无法再解密（返回 `ErrMFAInvalid` 面）。

### F-005 · 验证码一次性消费在 Postgres 不抗并发（MEDIUM · fixed）

- `modules/logincaptcha/store/repository.go` `ConsumeChallenge`：改为**单条守卫生成 DELETE**（`DELETE … WHERE id=? AND expires_at>? AND answer_hash=?` + RowsAffected==1 判定胜者），同事务内补一次无守卫 DELETE 保持「任意尝试即消费」契约（错误/过期首试也销毁 challenge，防暴破）。READ COMMITTED 下并发双正确答案恰好一方胜出。
- 回归锁：`internal/store/postgres_test.go` `TestPostgresCaptchaConsumeConcurrent`（PG_TEST_* 门控）：两并发消费恰好 1 胜，随后消费 false；sqlite 既有 lifecycle/expiry 用例不变。

### F-006 · 钱包对账坏 JSON 静默全库 + 写操作权限过松（MEDIUM · fixed）

- `handler/wallet.go` POST `/api/wallet/reconcile`：Decode 失败（`io.EOF` 空 body 除外——空 body 是文档化全库对账哨兵，jobs 测试锁定）→ 400 `INVALID_BODY`；垃圾 JSON 不再静默变全库对账。
- submit / cancel / retry 写操作权限从 `wallet.read` 收紧为 `wallet.write`（GET 维持 `wallet.read`）。`admin.wallet` manifest 的 wallet.write 键已存在（provider.go:186），无权限面泄漏。
- 显式哨兵字段未改（`accountId: ""` 语义被 jobs 测试与协议锁定；解码 400 已消除误触发面），决策与理由见 D-002 未选方案。
- 回归锁：`handler/wallet_test.go` `TestWalletReconcileBadBodyAndWriteGate`——垃圾 body 400；无 wallet.write 角色对 submit/cancel/retry 全 403；既有 reconcile 202 流程（空 body）保持。

## 回归验证

- API：`go vet ./...` 0 告警（见 A-002 复跑记录）；`go test ./...` 全绿（含本波新增用例；PG 门控用例在无 PG_TEST_* 环境自动 skip）。
- Web：本波未改任何 web 文件；`npm test` / `npm run build` 基线不变（A-002 记录复跑结果）。

## 产物

- 修复：上面 6 条对应的 13 个源码文件 + 2 个测试文件（mfa store/service、recyclebin handler、wallet handler）。
- 新增测试：`TestPostgresCreateUserManagementExistsScan`、`TestPostgresCaptchaConsumeConcurrent`、`TestCreateProofPurgesExpired`、`TestIncrementProofFailuresCapped`、`TestServiceRotationWindow`、`TestRecycleFactoryDeleteAndSnapshotSameTx`、`TestRecycleFactorySnapshotFailureRollsBackDelete`、`TestWalletReconcileBadBodyAndWriteGate`。

## 后续（计划，非事实）

- A-002 self 审计 → S4 independent 复核（grok build · grok-4.6）→ recommended 处置（F-007～F-019）→ 闭合记录 + go 恢复 + 关门。