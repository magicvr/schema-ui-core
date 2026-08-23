---
id: E-004-w9-s3-implementation
goal: GOAL-009-w9-api-web-security-audit
status: done
created: 2026-08-21
updated: 2026-08-21
parent: GOAL-009-w9-api-web-security-audit
version: 0.1.0
---

# E-004 · W9 S3 实施：D-003 范围 12 条 required 修复 + API/Web 回归（2026-08-21）

## 事实

按 [D-003](../01-decision/D-003-w9-scope-and-go-hold.md) 冻结范围实施全部 12 条 required；每条修复均为可核对的代码改动，未动 recommended/info 项。

| Finding | 修复位置与方式 |
|---------|----------------|
| F-001 | 新增 `internal/kernel/unique_violation.go` `IsUniqueViolation`（SQLite 文案 + pgx `SQLSTATE 23505`/duplicate-key 文案，走 unwrap 链）；`wallet/store/repository.go` 检测委托该助手；`GetOrCreateUserAccount` 重构为"快路径读 → 独立事务 INSERT → 唯一冲突后**新事务**重读"，消除 PG 上失败 INSERT 中止事务导致的死回退 |
| F-002 | `apps/web/nginx.conf` 增加 `location = /.well-known/schema-ui/host-bootstrap.json` 精确代理块（与 manifest 同型）；新增回归测试 `src/host/nginx-proxy.test.ts`（3 断言：两个 well-known 均存在且均 proxy_pass 到 api） |
| F-004 | `authsession/accounts.go` `RecordLoginFailure` 改为单条原子 `UPDATE ... failed_login_count + 1`（行锁串行并发失败），同事务内读计数判阈值、开锁窗口并清零——消除 PG READ COMMITTED 下的丢失更新 |
| F-005 | `mfa/store/repository.go` 新增 `AdvanceLastUsedStep`（守卫 UPDATE `WHERE last_used_step < ?` + RowsAffected 判定）；`mfa/service.go` `Verify` 以 CAS 结果作为重放门——并发同码仅一次成功 |
| F-006 | `mfa/store/repository.go` 新增 `UpdateRecoveryCodesIfUnchanged`（updated_at 乐观 CAS）；`service.go` `consumeRecoveryCode` 改为"匹配→CAS 写→失败则重读重试（≤4 次）"，并发兑换不同码不再互相复活、同码不可双用 |
| F-007 | `jobs/runner.go` 处理器 goroutine 加 `recover`，panic 转为 runnerOutcome 错误 → 走既有 `finish()` 记 JOB_HANDLER_FAILED；`scheduledtasks/scheduler.go` tick 循环加 recover（记 slog）、`Execute` 处理器调用加 recover → 记 failed task_run |
| F-008 | `web/src/renderer/permissions.ts` 新增 `actionGateTargetId`（key 优先、actionRef 兜底），四处 targetId 注册（row actions/toolbar × intent/local）与渲染端查找键对齐——声明 permissionIntent 但漏 key 的动作不再绕过 UI 门禁 |
| F-009 | `permissions.ts` `effectivePermission`：cascade 声明了 key 但缺 permissions 源时 **deny**（原 skip-and-continue 为 allow）；L2 校验器接线留作后续 hardening（运行时门禁已 fail-closed，finding 影响消除） |
| F-010 | `handler/resources.go` `delete()` 归属预检改为 fail-closed：Get 非 NotFound 错误 → `writeEntityError` 并中止（镜像 update()），不再跳过归属检查继续删除 |
| F-011 | `authsession/service_credentials.go` 重名判定改 `kernel.IsUniqueViolation` + 双方言约束名（sqlite `service_credentials.name` / PG `service_credentials_name_key`） |
| F-012 | `scheduledtasks/store/repository.go` `ListTasks` 与 `ListAllRuns` 的 q 子句加括号，enabled/status 过滤不再被 OR 优先级绕过 |
| F-025 | `scheduledtasks/store/cron.go` `Matches` 按 POSIX 语义：DOM/DOW 均受限（解析集未覆盖全值域）时取 OR；全值域映射行为等价于 `*`，子集检测与语法检测对所有表达式等价，`CronFields` 数组形状不变（describeCron 等调用方零改动） |

## 回归证据（全绿）

- API：`go build ./...` exit 0；`go vet ./...` exit 0；`go test ./...` **exit 0 全部包通过**（含 wallet/store、scheduledtasks/store、authsession、mfa、jobs、handler 132s）；F-010 收尾后 `internal/handler` + `internal/modules/recyclebin` 复跑 ok。
- Web：`npm test` **74 文件 / 1075 测试全部通过**（含新增 nginx-proxy 3 项、permissions-inheritance 18 项、row-action-bindings、schema-crud 22 项——无测试锚定旧缺陷行为）；`npm run build` exit 0（chunk 体积警告为既有非阻断项）。

## 边界说明

- F-013～F-023（recommended）/ F-024（info）本波未实施，符合 D-003 §3。
- VP-008 go 宣称维持暂挂（D-003 §4），恢复待 S4 required=0 后另写决策。
