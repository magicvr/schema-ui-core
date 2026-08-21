# 附件 · A-002 全文 · W9 A-001 合理性复审证据（2026-08-21）

> auditor：grok-build（grok-4.6 · reasoning high · `/audit`）。  
> 本附件是 A-002 的逐条证据；结论以 [A-002](../03-audit/A-002-w9-a001-reasonableness.md) 为准。  
> **不是** S4 关门审；不闭合 A-001 findings。

## 方法与边界

- 只读 `workspace-009-production-hardening` 与 `apps/api`、`apps/web` 现行源码。
- 对 A-001 每条 required 回读被引文件；对 P3 抽样。
- 未动态 exploit、未跑测试套件、未起 compose。路径以本会话工作树为准。

## A-001 台账自洽性（支撑 A-002 F-001）

| 来源 | 主张 | 独立观察 |
|------|------|----------|
| A-001 结论 / I-001 | 12 required：F-001/F-002 high + F-003～F-012 med | F-003 **无标题、无正文** |
| A-001 required 节 | 列出 11 个 `### F-0xx` | F-001, F-002, F-012, F-004, F-005, F-006, F-008, F-009, F-010, F-011, F-007 |
| A-001 必改项 | 「F-003～F-012 对应 F-012/F-004/F-005/F-006/F-007/F-008/F-009/F-010/F-011」 | 9 个编号，不是 10 |
| `03-audit.md` 索引 | 开放 required **22** | 与 12 或 11 都不符 |
| 全文 P2 | 10 项（API 7 + Web 2 + 潜伏 1） | 第 7 项 cron DOM/DOW **未映射到任何 F-ID** |

编号错位：A-001 把全文 P2-1（scheduledtasks WHERE）标成 **F-012**，而 F-003 空缺。全文 P2 项数（10）与「10 条 med required」的意图能对上，但落盘 F 清单对不上。

## 逐条源码复核

### A-001 F-001 · 钱包 unique 仅 SQLite — **成立 · 保留 high required**

`apps/api/internal/modules/wallet/store/repository.go:845-847`：

```go
func isUniqueViolation(err error) bool {
	return err != nil && (contains(err.Error(), "UNIQUE constraint failed") || contains(err.Error(), "constraint failed: UNIQUE"))
}
```

pgx `SQLSTATE 23505` 文案是 `duplicate key value violates unique constraint`，两处子串都不命中。调用点：

- `CreateAccount:231-232` → 本应 `ErrOwnerTaken`，PG 上包成 500。
- `GetOrCreateUserAccount:299-313` → 冲突回退死代码。A-001 关于「PG 事务在 unique 失败后中止，同事务重读也不可行」成立：即便修好匹配，仍须新事务重读。
- `Mutate:490-493` → `errIdempotencyRace` 在 PG 永不触发；随后 `:510-511` 的 replay 走不到。余额仍受 version CAS（`:479-480` `ErrVersionConflict`）保护，**无双花/错账**。

`kernel/store.go:16` 将 postgres 标为 production-authoritative dialect。仓库内无共享 `23505` 助手（仅此 `isUniqueViolation`）。compose 默认 SQLite 不削弱本条：PG 是一等方言。

### A-001 F-002 · nginx 未代理 host-bootstrap — **成立 · 保留 high required**（措辞见 A-002 F-002）

| 核对 | 路径 | 结果 |
|------|------|------|
| 仅代理 manifest | `apps/web/nginx.conf:33-42` | 只有 `location = /.well-known/schema-ui/app-manifest.json` |
| SPA fallback | `nginx.conf:60-62` | `try_files $uri $uri/ /index.html` |
| 404/410 才算未提供 | `apps/web/src/host/bootstrap.ts:172-188` | 200 + 非 JSON content-type → `protocol` |
| boot 终态 | `boot.ts:169-180` + `failure.ts` | `BOOTSTRAP_DOCUMENT_FAILED` → `HOST_PROTOCOL_REJECTED`，protocol-rejected 的 retry 为 `none` |
| 镜像打入 conf | `apps/web/Dockerfile:27` | `COPY apps/web/nginx.conf /etc/nginx/conf.d/default.conf` |
| compose 用该镜像 | `compose.yaml:60-67` | web dockerfile = `apps/web/Dockerfile` |
| API 确有路由 | `composition.go:530-534`、`handler/bootstrap.go:67` | `RegisterBootstrapWithAvailability` |
| 开发/preview 已代理 | `vite.config.ts:25-28,46-48` | 缺陷限于生产 nginx |

A-001 写终态 `HOST_PROTOCOL_REJECTED`，全文写 `BOOTSTRAP_DOCUMENT_FAILED`：两者是映射关系，不是矛盾。  
「栈必然启动失败」过述：web healthcheck 是 wget `/`（`compose.yaml:73-74`），容器可 healthy。

### A-001 F-012（全文 P2-1）· scheduledtasks OR/AND — **成立**

`scheduledtasks/store/repository.go:86-97`：

```
WHERE lower(key) LIKE … OR lower(name) LIKE … AND enabled = ?
```

AND 优先于 OR → `(key LIKE) OR (name LIKE AND enabled)`。`q+enabled=true` 会返回 key 命中的 **disabled** 行。`ListAllRuns:300-310` 同型，status 过滤只绑在第二支。  
对照 `datadictionary/store/repository.go:251` 已对 OR 加括号。  
handler 为 admin `tasks.read`（`scheduledtasks.go` 注释），不是匿名越权。

### A-001 F-004 · RecordLoginFailure RMW — **代码成立 · 影响过述**

`authsession/accounts.go:173-196`：事务内读 `failed_login_count`，应用层 `+1`，再 UPDATE 写回。无 `failed_login_count = failed_login_count + 1`，无 `FOR UPDATE`。PG 并发丢失更新可推迟 `LockThresholdFailures`（`auth.go:55` 为 5）。  
补偿：`auth.go:60` `newLoginRateLimiter(15*time.Minute, 20, …)`，login handler `:102-111` 先限流；captcha 在限流之后、验密之前。见 A-002 F-004。

### A-001 F-005 · TOTP 跨事务重放 — **成立**

`mfa/service.go:123-137`：`GetState`（独立 Run）→ `ValidateTotp(..., st.LastUsedStep)` → `SetLastUsedStep`（另一次 Run，`repository.go:161-170` 无 CAS，`WHERE user_id = ?` 无 `last_used_step` 谓词）。  
`DeleteProof:233-239` 不检查 RowsAffected；`mfa_proofs` 无 `user_id` 唯一约束（`migration.go:30-36`），`BeginChallenge` 每次新 proof。两条并发 Verify（同码、两 proof 或同 proof）均可接受。SQLite 亦存在（A-001 此点正确，与 F-004 不同）。

### A-001 F-006 · 恢复码丢失更新/双花 — **成立**

`consumeRecoveryCode`（`service.go:300-319`）：反序列化全表哈希 → bcrypt 命中 → `UpdateRecoveryCodes` 整数组回写（`repository.go:144-157`，按 `user_id` 覆盖，无乐观版本）。并发兑两枚不同码：后写覆盖前写，已消费码可复活。同码并发：两次都对原始切片匹配，双成功。

### A-001 F-007 · 无 panic recover — **代码成立 · 可用性**

- 实际路径：`apps/api/internal/jobs/runner.go:278-281`（A-001 漏 `internal/`）。handler 在无 recover 的 goroutine 里执行；未捕获 panic 终止整个进程。`:247-255` 的 defer 只清理 `active` map。
- `scheduledtasks/scheduler.go:64-75` loop 与 `:131` `handler(...)` 无 recover；`Start` 把 loop 放在 goroutine 里，handler panic 同样击穿进程。

### A-001 F-008 · 缺 key 后门禁跳过 — **部分成立**

`permissions.ts:353-356`：`targetId: stringValue(action.key)` → 缺 key 则为 `""`。  
`render.tsx:619` / `:914`：`gateTargetId = item.key || actionRef`。与空 targetId 对不上 → `hasPermissionEntry` 为 false → `:464-468` 跳过 `executeAction`。  
`schema-table.tsx:770,866`（`src/renderer/`）：`crud?.effectivePermission(key) ?? true`，且 `effectivePermission`（`render.tsx:750-754`）对 absent target **显式 return true**，注释写明引擎默认 allow。  
结论：`permissionIntent` 已声明但漏 `key` 是失配 bug；未声明 intent 的默认放行是协议，不是「门禁整体失效」。服务端鉴权仍是硬边界（A-001 已写）。L2 `validatePermissions` **不检查** action.key 是否存在。

### A-001 F-009 · cascade 缺源 fail-open — **运行时成立**

`permissions.ts:512-515`：source `undefined` 时不 `return false`。`:520` 默认 `true`。  
`validatePermissions:202-205` 会报 `PERMISSION_CASCADE_SOURCE_MISSING`。全仓调用方仅 `permissions-inheritance.test.ts:58`。生产 `load-page` / render 路径不跑 L2。

### A-001 F-010 · delete 归属预检 fail-open — **潜伏成立**

`handler/resources.go:716-723`：仅 `gerr == nil` 时做 `scopeOwned`；Get 出错则继续 `Delete`。`update():680-683` 对非 `errResourceNotFound` fail-closed。  
`Scoper:` 赋值全仓只在 `resources_test.go:356`。无生产资源接线时，self-scope 预检本身不跑（`resolveScope` 在 Scoper nil 时返回 nil，`:304-307`）。缺陷是模式级潜伏。

### A-001 F-011 · 凭据重名方言 — **成立**

`service_credentials.go:59`：`strings.Contains(..., "service_credentials.name")`（SQLite 索引/列文案）。  
PG DDL `migration.go:254`：`name CITEXT NOT NULL UNIQUE` → 约束名典型为 `service_credentials_name_key`。pgx 23505 不含 `service_credentials.name` → 400 `ErrCredentialNameTaken` 走不到，变 500。

### 全文 P2-7 · cron DOM/DOW AND — **成立 · 无 F-ID**

`cron.go:99-107`：`Matches` 要求 DOM **与** DOW 都命中。Vixie/POSIX 在 DOM、DOW 均非 `*` 时用 OR。`0 0 1 * 1` 只在「1 号且周一」触发。属调度少跑，不是多跑/越权。A-001 required 节未收录。

## P3 抽样

| A-001 | 抽样结果 |
|-------|----------|
| F-013 头像配额 TOCTOU | `account_avatar.go:62-71` CountOwner 后 `storeUploadForOwner`，无 `quotaMu`。`upload.go:120-123,341-348` 确有配额锁先例。成立。 |
| F-016 otpauth 转义 | `totp.go:77-84` 只替换空格/冒号。成立。 |
| F-020 `Contains("unique")` 4 处 | `datadictionary` ×3 + `scheduledtasks` ×1。`RequestCancel`（`jobs/repository.go:124-148`）忽略 RowsAffected。`CreateUserManagement:87-103` 先 EXISTS 再 INSERT，并发 PG unique → 500。成立。 |
| F-021 recycle 两事务 | `recyclebin/service.go:70-85`：`restoreRow` 成功后再 `MarkRestored`。崩溃可留「行已在、快照未标记」→ 重试 `RECYCLE_RESTORE_CONFLICT`。成立。 |
| F-022 rebind / WithTx / 自动建角色 | `rebind.go:13-28` 逐 rune 替换 `?`，字面量中的 `?` 会被改写（注释承认为 contract violation）。`Store.Run:91-96` **有** recover+rollback 再 re-panic；`WithTx:112-124` 无 recover，但是 sqlite 测试缝，生产走 Run。`CreateUser`→`linkUserRole`→`ensureRole` 可 `ON CONFLICT DO NOTHING` 建任意合法 key；管理面 `CreateUserManagement:78-84` 未知角色 `ErrInvalidRole`。A-001「当前仅测试调用」对 CreateUser 路径基本对；把 WithTx 写成生产死锁过严。 |
| F-023 `//host` / 负 maxSize | `App.tsx:830-834`：`href.startsWith("/")` 放行 `//host`。`upload-orchestration.ts:123`：`maxSize >= 0` 才校验，负值跳过客户端上限（服务端仍有 8MiB）。成立。 |
| F-024 `.env` | `git check-ignore -v apps/api/configs/.env` → `apps/api/.gitignore:9:.env`；`git ls-files` 未跟踪。成立。本附件不摘录内容。 |

## 路径/行号偏差（不单独开 finding）

| A-001 写法 | 实际 |
|------------|------|
| `jobs/runner.go:278-281` | `apps/api/internal/jobs/runner.go:278-281` |
| `schema-table.tsx` 暗示 components | `apps/web/src/renderer/schema-table.tsx:770,866` |
| `render.tsx:914` 门禁查找 | 当前 `gateTargetId` 在 `:914`；`hasPermissionEntry` 在 `:464-467` / 批处理 `:619` |
| `composition.go:532` | 文件在 `apps/api/internal/composition/composition.go:532`，内容确实是 RegisterBootstrap |

行号整体仍可定位，不构成误报。

## 排除（同意 A-001「已排除的疑点」方向）

本复审未做全量注入/路径穿越重扫。抽查与 A-001 排除清单无矛盾：排序白名单、上传 ID 形态、钱包整数+CAS、nginx `script-src 'self'` 等不在本条 required 范围。未发现新的 P0。
