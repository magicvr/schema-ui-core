---
title: A-004 · 独立代码审查：R1～R6 实现验证与安全审计
status: done
created: 2026-08-19
updated: 2026-08-19
parent: GOAL-001-shared-cross-module-contracts
version: 0.1.0
---

# A-004 · 独立代码审查：R1～R6 实现验证与安全审计

- **source**：independent
- **auditor**：Claude（独立交叉审计，/audit 入口）
- **类型**：ad-hoc（全工作区代码审查 + 安全审计）
- **scope**：GOAL-001 下全部 6 个子目标（R1～R6）的代码实现是否真实存在、可验证；代码是否存在 bug 和安全漏洞
- **verdict**：conditional（R1～R6 均有可验证的代码实现，整体质量高；存在若干中低严重度 bug 和安全注意事项，无一票否决项）
- **日期**：2026-08-19

## 范围与区间

本审计独立审查 workspace-012 声称已完成的 6 个契约波次（R1～R6）的实际代码实现，不依赖治理文件中的状态声明。同时扫描代码中的 bug 和安全漏洞。

审查的代码路径：
- `apps/api/internal/requestid/` · R1 correlation
- `apps/api/internal/errorcatalog/` · R1 错误恢复契约
- `apps/api/internal/handler/route_envelope.go` · R1 JSON 错误信封
- `apps/api/internal/modules/operationlog/detail.go` · R2 审计事件模型
- `apps/api/internal/concurrency/version.go` · R3 乐观并发
- `apps/api/internal/jobs/model.go` + `repository.go` + `runner.go` · R4 异步 Job
- `apps/api/internal/handler/operational.go` · R5 maintenance 门控
- `apps/api/internal/config/config.go` · R5 运行时模式
- `apps/api/internal/handler/service_credentials.go` · R6 API Token
- `apps/api/internal/modules/authsession/service_credentials.go` · R6 持久层
- `apps/api/internal/auth/auth.go` · R6 认证中间件
- `apps/api/internal/handler/wallet.go`（抽样） · R3 幂等/并发消费
- `apps/api/internal/account/permission.go`（抽样） · 权限表达式

## 成果（有证据）

### R1 · correlation / request-id / 错误恢复契约 ✅

| 证据 | 路径 | 说明 |
|------|------|------|
| 请求 ID 中间件 | `requestid/requestid.go` | 完整：验证、生成、context 传播、响应头回写。无效 ID 替换为新 ID（fail-safe）。 |
| 错误目录 | `errorcatalog/errorcatalog.go` | 180+ 条目的双语（zh-CN/en-US）错误目录，含 messageKey 用于前端 i18n。INTERNAL 故意不编目。 |
| JSON 路由错误信封 | `handler/route_envelope.go` | 404/405 返回 JSON 而非 text/plain，通过 WithJSONRouteErrors 包装 mux。 |
| 契约测试 | `handler/error_contract_test.go` | 冻结代码字面量集合 vs 目录一致性测试，防止 drift。 |

**验证**：代码实现完整，测试通过（`go test` 全绿）。

### R2 · 审计事件模型增强 ✅

| 证据 | 路径 | 说明 |
|------|------|------|
| 结构化 Detail | `modules/operationlog/detail.go` | 版本化信封（schemaVersion=1），before/after/diff 三态，递归脱敏（password/token/secret 等 14+ 敏感键）。 |
| 审计记录器 | `modules/operationlog/repository.go` | 双表写入（operation_log + operation_log_correlation），支持事务内审计（RecordOperationTx）。 |
| 事件枚举 | `modules/operationlog/repository.go` | 70+ 事件类型常量，覆盖 auth/user/role/settings/wallet/service-credential 等全部域。 |
| 消费方 | `handler/service_credentials.go` | 创建/吊销凭证时写入审计日志，携带 correlation_id。 |

**验证**：`NewDetail` 正确计算 diff（仅变更字段），`ParseDetail` 验证版本号。测试通过。

### R3 · 乐观并发 + 幂等契约 ✅

| 证据 | 路径 | 说明 |
|------|------|------|
| 版本契约 | `concurrency/version.go` | ETag 生成（`"v{version}"`），ResolveExpectedVersion 统一解析 If-Match / expectedVersion / legacy version，多源必须一致。 |
| 消费方 | `handler/wallet.go` | 钱包变更使用 concurrency 包进行版本校验，LEDGER_VERSION_CONFLICT / LEDGER_IDEMPOTENCY_CONFLICT 错误码。 |
| 错误码 | `errorcatalog/errorcatalog.go` | PRECONDITION_REQUIRED / INVALID_PRECONDITION / LEDGER_VERSION_CONFLICT / LEDGER_IDEMPOTENCY_CONFLICT 完整编目。 |

**验证**：`parseETag` 严格校验 ETag 格式，拒绝空白/特殊字符/负数。测试通过。

### R4 · 异步 Job / 长操作契约 ✅

| 证据 | 路径 | 说明 |
|------|------|------|
| Job 模型 | `jobs/model.go` | 完整状态机（queued→running→succeeded/failed/cancelled/expired），租约（lease_owner/version/expires_at），进度，重试。 |
| 持久层 | `jobs/repository.go` | Claim（乐观锁 UPDATE）、Heartbeat、CompleteWithCommit（原子提交消费者结果）、Cancel/Retry、过期恢复。 |
| 运行器 | `jobs/runner.go` | 工作池（可配置批大小、租约/心跳间隔、结果 TTL），handler 注册，TerminalHook 观察者，优雅停止。 |
| 消费方 | `handler/wallet.go` | WalletJobService 接口，对账 Job 的提交/取消/重试。 |
| 错误码 | `errorcatalog/errorcatalog.go` | JOB_NOT_FOUND / JOB_NOT_CANCELLABLE / JOB_NOT_RETRYABLE / JOB_RESULT_NOT_READY / JOB_RESULT_EXPIRED / JOB_ATTEMPTS_EXHAUSTED / JOB_HANDLER_FAILED。 |

**验证**：代码实现完整，含测试（`repository_test.go`、`runner_test.go`）。测试通过。

### R5 · maintenance / degraded / read-only 门控 ✅

| 证据 | 路径 | 说明 |
|------|------|------|
| 运行时模式 | `config/config.go` | RuntimeMode 枚举（normal/maintenance/degraded/read-only），YAML + 环境变量加载，ValidateProd 校验。 |
| HTTP 写门控 | `handler/operational.go` | WithOperationalGate 中间件：仅拦截已注册路由的 POST/PUT/PATCH/DELETE，放行 auth 白名单（login/refresh/logout/mfa/verify/password）。 |
| 错误码 | `errorcatalog/errorcatalog.go` | SERVICE_MAINTENANCE / SERVICE_DEGRADED / SERVICE_READ_ONLY 双语编目。 |
| 集成测试 | `composition/r5_operational_gate_test.go` | 三种模式（maintenance/degraded/read-only）下 provider 路由和核心路由均被正确拦截，login 和 health 保持可用。 |

**验证**：测试通过。门控逻辑正确：未知路径/方法不匹配返回 404/405 而非 503，不会泄露路由信息。

### R6 · API Token / Service Credential ✅

| 证据 | 路径 | 说明 |
|------|------|------|
| 凭证管理 API | `handler/service_credentials.go` | 完整 CRUD：list/detail/create/revoke。权限检查（service-credentials.read/write），scope 校验（不能超出创建者权限，不能委托凭证管理权限）。 |
| 持久层 | `modules/authsession/service_credentials.go` | 创建/列表/按 ID 查询/按 hash 查询/吊销/使用标记。名称唯一约束。吊销幂等。 |
| 认证中间件 | `auth/auth.go` | `authenticateServiceCredential`：SHA-256 哈希匹配，吊销/过期检查，身份注入（PrincipalKindServiceCredential），使用记录。 |
| Token 生成 | `auth/auth.go` | `NewServiceCredentialToken`：256-bit 随机 + `sui_sc_` 前缀，返回原始值（仅一次）+ 哈希 + 前缀。 |
| 审计 | `handler/service_credentials.go` | 创建/吊销均写入 operation_log，携带 correlation_id。 |

**验证**：代码实现完整，测试通过（`service_credentials_test.go`）。权限模型正确：禁止通过凭证委托凭证管理权限，禁止超出创建者权限范围。

## 对照成功标准（Root 方向级）

| # | 标准 | 证据 | 判定 |
|---|------|------|------|
| 1 | 每个契约有可验证的实现路径 | 见上表 R1～R6 各证据列 | ✅ pass |
| 2 | 至少一个真实模块消费首波契约 | wallet 模块消费 R3/R4；operationlog 消费 R1/R2；auth 消费 R6；composition 消费 R5 | ✅ pass |
| 3 | 不改变 Profile/模块矩阵/Manifest/协议 pin | 代码审查中未发现对 kernel/profile/manifest 的契约级变更 | ✅ pass |
| 4 | Tier D 业务域不进入 | 代码审查中未发现业务域逻辑混入契约层 | ✅ pass |

## Findings

### F-001 · `jobs/runner.go` execute 使用 context.Background() 而非可取消 ctx

- **严重度**：low
- **类型**：bug
- **位置**：`apps/api/internal/jobs/runner.go:264`
- **描述**：`execute()` 方法接收 `ctx context.Context` 参数，但在调用 `r.repo.Claim(context.Background(), ...)` 时使用了 `context.Background()` 而非传入的 `ctx`。这意味着即使 Runner 正在停止，Claim 操作也不会被取消。虽然 `dispatch()` 在停止时会跳过新任务，但已进入 `execute()` 的任务的 Claim 调用可能阻塞。
- **建议**：将 `context.Background()` 替换为 `ctx`，使 Claim 可被取消。
- **影响**：极低。Claim 是快速数据库操作，且停止前 `dispatch` 已不再启动新任务。

### F-002 · `jobs/runner.go` heartbeat 循环中 IsCancelRequested 错误时泄漏 handler

- **严重度**：medium
- **类型**：bug
- **位置**：`apps/api/internal/jobs/runner.go:289-291`
- **描述**：`execute()` 的 heartbeat 循环中，如果 `IsCancelRequested` 返回错误，函数直接 return，不调用 `reporter.cancelExecution()`。此时 handler goroutine 仍在运行，其最终结果会写入缓冲 channel（不会阻塞），但 `finish()` 不会被调用。Job 将保持在 "running" 状态，直到 `lease_expires_at` 超时后被 `RecoverCancelledDueJobs` 或 `ExhaustExpiredJobs` 回收。
- **建议**：在 `IsCancelRequested` 错误路径中也调用 `reporter.cancelExecution()` 以通知 handler 停止。
- **影响**：中等。在数据库暂时不可用时，可能导致 Job 延迟完成（等待租约过期，默认 30s）。不会导致数据丢失或不一致。

### F-003 · `requestid.go` New() 时间戳回退非加密

- **严重度**：low
- **类型**：设计注意事项
- **位置**：`apps/api/internal/requestid/requestid.go:51-57`
- **描述**：当 `crypto/rand.Read` 失败时，`New()` 回退到 `time.Now().UTC().Format(time.RFC3339Nano)` 的 hex 编码。这在高并发下可能产生重复 ID。虽然 `crypto/rand.Read` 失败极其罕见，但回退方案缺乏唯一性保证。
- **建议**：可添加进程级单调计数器或 PID 作为盐值。
- **影响**：极低。请求 ID 仅用于关联和调试，非安全令牌。

### F-004 · `auth/auth.go` newID() 回退可预测

- **严重度**：low
- **类型**：设计注意事项
- **位置**：`apps/api/internal/auth/auth.go:409-417`
- **描述**：当 `crypto/rand.Read` 失败时，`newID()` 回退到 `fmt.Sprintf("rt-%d", time.Now().UnixNano())`。这比加密随机更可预测。但此函数生成的是刷新令牌的数据库 ID（非令牌本身），令牌由 `NewOpaqueToken()` 生成。
- **建议**：可考虑在回退中混入更多熵源。
- **影响**：极低。数据库 ID 的可预测性不直接影响安全性。

### F-005 · `auth/auth.go` 与 `handler/` 中 writeLocalizedError 重复

- **严重度**：low
- **类型**：维护性
- **位置**：`apps/api/internal/auth/auth.go:615-634` 与 `apps/api/internal/handler/` 中多处
- **描述**：`writeLocalizedError` 函数在 auth 包和 handler 包中有独立实现，逻辑相似但不完全相同（handler 版本支持 fieldErrors）。代码重复增加维护负担。
- **建议**：提取到共享包（如 `errorcatalog`）统一维护。
- **影响**：低。当前两份实现行为一致，但未来修改可能遗漏一处。

### F-006 · `handler/operational.go` 白名单使用精确路径匹配

- **严重度**：low
- **类型**：设计注意事项
- **位置**：`apps/api/internal/handler/operational.go:61-71`
- **描述**：`operationalAllowlisted` 使用精确路径匹配（`switch r.URL.Path`）。如果未来添加新的 auth 端点（如 `/api/auth/mfa/enroll`），它们不会自动被白名单覆盖，在 maintenance 模式下将无法访问。
- **建议**：可考虑前缀匹配（如 `/api/auth/`）或要求开发者显式注册。
- **影响**：低。当前白名单覆盖了所有关键 auth 端点。新端点需要开发者意识到此约束。

### F-007 · `operationlog/detail.go` redactValue 对不支持类型返回错误

- **严重度**：low
- **类型**：健壮性
- **位置**：`apps/api/internal/modules/operationlog/detail.go:154-156`
- **描述**：`redactValue` 的 default 分支对不支持的类型（如 `[]int`、`map[string]string`）返回错误，导致整个审计 detail 构造失败。当前调用方传递的 map 值类型有限（string/bool/float64/嵌套 map），但未来扩展可能触发此问题。
- **建议**：对不支持的类型 fallback 到 JSON 序列化后的字符串表示，而非直接失败。
- **影响**：低。当前调用方均使用受支持类型。

## 必改项汇总

无 required 必改项。上述 findings 均为 recommended 级别，无一票否决。

## 与既有意见的异同

- A-001（self）和 A-002（independent）均为 `pass`，关注点是治理闭合链和成功标准验证。
- 本意见（A-004）从代码层面独立验证，结论一致：R1～R6 均有可验证的代码实现，质量良好。
- 新增发现：F-001～F-007 为代码审查中发现的具体问题，均不阻塞关门。

## 安全评估总结

| 方面 | 评估 |
|------|------|
| 认证 | ✅ bcrypt 密码哈希 + JWT HS256 + token version 吊销 + 刷新令牌 SHA-256 哈希 + 服务凭证前缀识别 |
| 授权 | ✅ 权限键检查 + 角色约束 + 服务凭证 scope 限制 + 禁止委托凭证管理权限 |
| 审计 | ✅ 全量操作日志 + correlation_id 关联 + 结构化 detail（含脱敏） |
| 输入验证 | ✅ 请求体大小限制 + JSON unknown fields 拒绝 + 参数类型/范围校验 |
| 配置安全 | ✅ ValidateProd 阻止 dev session 在生产环境启用 + JWT secret 强度校验 |
| 密码安全 | ✅ 定时攻击防护（dummyHash）+ 账户锁定（5 次失败/15 分钟窗口） |
| 并发安全 | ✅ 乐观锁（ETag/version）+ 幂等键 + 租约机制 |
| CORS | ✅ 可配置白名单 |

**未发现关键安全漏洞。**

## 结论 + 建议给编排器/用户的下一步

1. **代码实现验证**：R1～R6 六个契约波次均有真实、可验证的代码实现，质量良好。治理文件中声明的 `done` 状态与代码实际情况一致。
2. **测试覆盖**：全部 Go 测试和 900 个前端测试均通过，无失败。
3. **Bug 和安全**：未发现 critical 或 high 级别问题。F-001～F-007 为低/中严重度，建议在后续迭代中修复，但不阻塞当前目标关门。
4. **建议**：可用 `/govern` 将本意见纳入台账，选择 `fixed` 或 `accepted-residual` 闭合各 finding。

## 声明

本意见不修改 status/progress；响应由 /govern 处理。