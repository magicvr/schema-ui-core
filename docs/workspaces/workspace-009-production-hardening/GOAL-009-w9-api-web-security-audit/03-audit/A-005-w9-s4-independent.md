---
id: A-005-w9-s4-independent
doc: audit-entry
goal: GOAL-009-w9-api-web-security-audit
title: W9 S4 关门前 independent 复核（D-003 范围 12 条 required）
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
date: 2026-08-21
scope: D-003 冻结范围 12 条 required 是否 genuine fixed，且未引入新缺陷（S4 / finding-closure）
verdict: pass
status: recorded
parent: GOAL-001-production-hardening
created: 2026-08-21
updated: 2026-08-21
version: 0.1.0
---

# A-005 · W9 S4 关门前 independent 复核（2026-08-21）

## 条目头

| 字段 | 值 |
|------|-----|
| **source** | independent |
| **auditor** | grok-build (grok-4.6 · reasoning high · `/audit`) |
| **类型** | close-out / finding-closure |
| **scope** | D-003 冻结范围 12 条 required（F-001、F-002、F-004～F-012、F-025；F-003 作废）是否 genuine fixed，且未引入新缺陷。消费清单权威 = [D-002](../01-decision/D-002-w9-finding-inventory.md)；实施 = [E-004](../02-execution/E-004-w9-s3-implementation.md)；self = [A-004](A-004-w9-self.md) |
| **verdict** | **pass** |
| **工作区** | `workspace-009-production-hardening`（Root `GOAL-001-production-hardening`；canonical `docs/workspaces/workspace-009-production-hardening/`；`plan_refs`/`primary_plan` = `VP-009-production-hardening`；`vision_ref` = `schema-ui-core-admin-foundation@0.2.0`；`shared_materials_catalog: none`） |

## 范围与区间

- **覆盖**：D-002 消费 ID 的 12 条 required 在现行代码中是否对齐 A-001/D-002 定义并 genuine fixed；回归证据是否可重复核对；修复是否引入并发/方言/行为回归；A-004 self 结论是否与代码事实相符。
- **方法**：工作区绑定核对 → GOAL-009 五件套通读 → 12 条点名路径源码抽验（含调用点与消费者对齐）→ 本会话重跑回归（见下）。未做动态 exploit、未起 compose、未接 live Postgres。
- **不覆盖**：不改 `status` / `progress` / 检查点 / 方案正文 / goal-tree。不把 F-013～F-024 升格为本波 required。不自行恢复 VP-008 `go` 宣称，不把本意见当作已关门。
- **排除**：F-003 作废不审；I-003 仍为 non-blocking 用户裁决（本条即 D-003 §6 所要求的 grok 复核，但不代 `/govern` 关闭该信息项）。

## P-005 / 工作区核对

| 核对项 | 结论 |
|--------|------|
| 工作区绑定 | `workspace.md`：`id=workspace-009-production-hardening`；Root `GOAL-001-production-hardening`；canonical 与 `goal-tree.md` 一致；`vision_role: delivery`；`primary_plan` = `VP-009-production-hardening`。Charter `schema-ui-core-admin-foundation@0.2.0`；VP-009 `vision_ref` 精确匹配。共享资料目录 `none`，本 scope 未把资料引用当关闭证据。未读取其他工作区。 |
| I-001（finding 清单） | **verified**（D-002）；本条不重开。消费 12 = F-001/F-002 + F-004～F-012 + F-025。 |
| I-002（范围 + VP-008 go） | **verified**（D-003：整单 12 + 闭合前暂挂 go）。本条确认代码闭合条件已满足；**恢复对外 go 宣称仍须 `/govern` 书面决策**，不得由本意见直接改宣称。 |
| I-003（provider 偏差 / 是否追加 grok 复核） | **open / non-blocking**；最晚阶段为关门前。本条满足 D-003 §6 的 grok independent 腿，不替代用户对该信息项的书面关闭。 |
| 到期 required 信息项 | 无到期未关闭 required 信息项阻断本 S4 代码闭合 scope。 |
| 共享资料 | 无 |

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| 12/12 required 均有可核对代码改动，位置与 D-002 定义对齐 | 见下表；无偷换范围、无以 F-013～F-024 冒充 required |
| A-004 回归主张可重复核对 | 本会话：`apps/api` `go test ./...` **exit 0**；`go vet ./...` **exit 0**；受影响包 `-count=1` 全绿（kernel / wallet/store / authsession / mfa / jobs / scheduledtasks / scheduledtasks/store / handler）；`apps/web` `npm test` **74 文件 / 1075 测试**；`npm run build` **exit 0**（chunk 体积警告为既有非阻断项） |
| A-004 self 结论与代码事实相符 | 12 条 genuine fixed；S3 边界（recommended 未实施、go 暂挂）成立；A-004 关于 F-009 L2 校验器的 recommended 残余本条确认 |
| 未引入阻断性新缺陷 | 源码抽验未见把 fail-closed 改回 fail-open、未见钱包 isNew 语义回退、未见 cron 数组形状破坏；残余见本条 recommended |

## 对照成功标准

| 标准 | 本 scope | 状态 | 证据 |
|------|----------|------|------|
| S1 独立审计落盘 | 前序 | 达成 | A-001 |
| S2 用户范围/go 裁决 | 前序 | 达成 | D-002 + D-003 |
| S3 按范围实施 + API/Web 回归 | 前序 + 本条复跑 | 达成 | E-004 + 本会话回归 |
| S4 self/independent 复核 required 合法闭合 | **本条判定代码闭合条件已满足** | 达成（意见层） | A-004 self pass + 本条 independent pass。本意见不改 status；合法闭合与关门由 `/govern` 响应 |

## 逐条闭合判定（D-002 消费 ID）

### F-001 · 钱包唯一冲突检测仅 SQLite → **fixed**

| 核对项 | 路径 | 结论 |
|--------|------|------|
| 双方言谓词 | `apps/api/internal/kernel/unique_violation.go` | `IsUniqueViolation` 走 unwrap 链；匹配 SQLite `UNIQUE constraint failed` / `constraint failed: UNIQUE` 与 PG `duplicate key value violates unique constraint` / `SQLSTATE 23505` |
| pgx 文案属实 | `github.com/jackc/pgx/v5@v5.10.0/pgconn/errors.go` `PgError.Error()` | `Severity + ": " + Message + " (SQLSTATE " + Code + ")"`；23505 与 duplicate-key 文案均命中。生产驱动为 `pgx/stdlib`（`store/postgres.go`） |
| wallet 委托 | `wallet/store/repository.go:829-835` | `isUniqueViolation` → `kernel.IsUniqueViolation`；CreateAccount / Mutate 幂等冲突 / GetOrCreate 共用 |
| PG 中止事务死回退 | `GetOrCreateUserAccount` | 快路径独立读 → **独立事务 INSERT** → 冲突后 **新事务** `GetUserAccountByOwner`。`store.Run` / `postgres.Run` 把回调错误原样返回（先 rollback），失败者不再在已中止事务里重读 |
| 并发语义 | `wallet/store/concurrent_test.go` | SQLite 下 8 路 get-or-create：恰好 1 个 `created=true`、1 行。isNew 语义保持（落败方不报 create） |

原缺陷（PG 上 `errIdempotencyRace` / 冲突回退永不触发）不再成立。无 `IsUniqueViolation` 单测、无 live PG 对 23505 的集成锁——见本条 recommended F-003。

### F-002 · 生产 nginx 未代理 host-bootstrap → **fixed**

| 核对项 | 路径 | 结论 |
|--------|------|------|
| 精确代理块 | `apps/web/nginx.conf:44-58` | `location = /.well-known/schema-ui/host-bootstrap.json` 与 manifest 同型：`proxy_pass http://api:25080` + 同源头/超时 |
| 镜像打入 | `apps/web/Dockerfile:27`；`compose.yaml` web `dockerfile: apps/web/Dockerfile` | 生产 conf 仍是该文件；SPA fallback 不再抢该路径 |
| 回归锁 | `apps/web/src/host/nginx-proxy.test.ts` | 3 断言：两 well-known 精确 location 均在，且均 `proxy_pass` 到 api |
| 原触发链 | `host/bootstrap.ts:172-188` | 仅 404/410 视为未提供；代理后 API 可返回 JSON 或真实 404，不再 200 `text/html` |

原「compose 容器 healthy、浏览器 `HOST_PROTOCOL_REJECTED`」路径被切断。未起 compose 做真实 HTTP 探测；静态 conf + 镜像 COPY + 单测足以核对 finding 定义。

### F-004 · RecordLoginFailure 非原子 → **fixed**

`authsession/accounts.go:178-210`：单条 `UPDATE users SET failed_login_count = failed_login_count + 1`，同事务 `SELECT` 计数，达阈值再 `failed_login_count = 0` 并写 `locked_until`。行锁串行并发失败，消除 PG READ COMMITTED 丢失更新。`affected==0` → `ErrNotFound`。阈值/清零语义与原先一致（仅并发正确性）。无专门并发单测（recommended F-003）。

### F-005 · TOTP 跨事务 check-then-act → **fixed**

| 核对项 | 路径 | 结论 |
|--------|------|------|
| CAS | `mfa/store/repository.go:179-199` `AdvanceLastUsedStep` | `UPDATE … WHERE user_id = ? AND last_used_step < ?`，`RowsAffected==1` 才算赢 |
| Verify 以 CAS 为门 | `mfa/service.go:131-140` | `ValidateTotp` 通过后必须 `advanced` 才 `ok`；同码并发第二路 0 行 → `ErrMFAInvalid` |
| 顺序重放仍锁 | `mfa/service_test.go` `TestServiceLifecycle` | 同 step 新 proof 拒绝 |

原 GetState→Validate→`SetLastUsedStep` 跨事务双接受不再是 Verify 路径。`SetLastUsedStep` 仍用于 Confirm 激活（非本 finding 的 login proof 面）。`requireActiveSecondFactor` 的 TOTP 支路仍只校验不推进水位——disable/rotate 与 login 并发同码属本波外残余，不重开 F-005。

### F-006 · 恢复码整表回写丢失更新 → **fixed**（同秒 OCC 残余见 recommended F-002）

| 核对项 | 路径 | 结论 |
|--------|------|------|
| 乐观 CAS | `UpdateRecoveryCodesIfUnchanged` | `UPDATE … WHERE user_id = ? AND updated_at = ?` |
| 重试 | `service.go:309-343` `consumeRecoveryCode` | 匹配 → CAS；失败则 `GetState` 重读，最多 4 次。同码在 **令牌推进** 时第二路重读后 bcrypt 不再命中 |
| 顺序一次性 | `TestServiceLifecycle` | 同一 recovery code 第二次 Verify → `ErrMFAInvalid` |

原「跨事务整表回写、不同码互相复活 / 同码双成功」在 `now.Unix() != prevUpdatedAt.Unix()`（常见双提交：上次 MFA 写发生在更早一秒）下已消除。`updated_at` 为 Unix 秒：若首次 CAS 时 `now` 与上次写入同一秒，WHERE 令牌不推进，两路都可 `affected=1`。见本条 recommended F-002；**不把 F-006 判回 open**——原无界丢失更新已有 CAS+重试，残余窗口有界。

### F-007 · job/task 无 panic recover → **fixed**

| 核对项 | 路径 | 结论 |
|--------|------|------|
| jobs | `internal/jobs/runner.go:278-288` | handler goroutine `defer recover` → `runnerOutcome{err}` → 既有 `finish()` → `JOB_HANDLER_FAILED` |
| scheduler tick | `scheduledtasks/scheduler.go:73-82` | tick 内 IIFE recover，记 slog，loop 继续 |
| Execute | `scheduler.go:141-165` | handler recover → `detail` 非空 → `Status=failed` 记 `task_run` |

单处理器 panic 不再击穿进程。无 panic 注入单测（recommended F-003）。

### F-008 · 权限门禁缺 key 时 fail-open → **fixed**

| 核对项 | 路径 | 结论 |
|--------|------|------|
| 注册键 | `permissions.ts:312-318` `actionGateTargetId` | `key` 优先，否则 `actionRef` |
| 四处注册 | `permissions.ts` row intent / row local / toolbar intent / toolbar local | 均用 `actionGateTargetId`（E-004「四处」即此，不含 actionButton） |
| 消费者对齐 | `render.tsx:914`；`schema-table.tsx:770,866`（及 427） | 同一 `key \|\| actionRef` 查找 |

声明 `permissionIntent` 但漏 `key` 的动作不再以 `targetId=""` 注册、被 `actionRef` 查找 miss 后跳过门禁。`actionButton` 仍用 `props.key \|\| node.id`（协议 D4b，非本 finding 点名路径）。`schema-table` 对未注册 target 仍 `?? true`；本修复使 intent 动作被注册，原击穿不再发生。

### F-009 · cascade 缺源 fail-open → **fixed**

`permissions.ts:527-533`：cascade 声明了 key 但 `permissions[key]` 缺失或求值为 false → **deny**。末端 `return true` 仅当无 cascade 命中且无 local——与 A-002「未标记默认 allow 是协议」一致。`validatePermissions` 仍仅 `permissions-inheritance.test.ts` 调用（本条 recommended F-001 / 沿用 A-004）。运行时门禁已 fail-closed，finding 安全影响消除。

### F-010 · delete() 归属预检 fail-open → **fixed**

`handler/resources.go:715-731`：Get 非 `errResourceNotFound` → `writeEntityError` 并 **return**（镜像 `update()`:680-683），不再跳过归属检查继续 `Delete`。`scopeOwned(nil constraint)` 为 true，Trash 预读不会把全局 scope 误 404。batch-delete 在 Get 失败时 **跳过该 id 不删**（fail-closed 另一形态），非本 finding 点名路径。A-001 已注明无生产 Scoper；潜伏模式已对齐。

### F-011 · 凭据重名仅 SQLite 文案 → **fixed**

`authsession/service_credentials.go:59-67`：`kernel.IsUniqueViolation(err)` **且** 约束名匹配 `service_credentials.name`（SQLite）或 `service_credentials_name_key`（PG 列级 UNIQUE 默认名；`migration.go:254` CITEXT UNIQUE）。其它唯一冲突（如 `token_hash`）不误映射为 `ErrCredentialNameTaken`。`service_credentials_test.go` 在 SQLite 锁定重名 → sentinel。未跑 live PG。

### F-012 · scheduledtasks WHERE OR/AND → **fixed**

`scheduledtasks/store/repository.go` `ListTasks`（约 86-99 行）与 `ListAllRuns`（约 302-316 行）：`q` 的 OR 组已加括号，随后 `AND enabled` / `AND status` 不再被 OR 优先级绕过。对照 datadictionary entries 的加括号先例。LIKE `%/_` 未转义属 F-017 recommended，本波未修。

### F-025 · cron DOM/DOW 用 AND → **fixed**

`scheduledtasks/store/cron.go:107-118` `Matches`：`domRestricted && dowRestricted` 时 **OR**，否则 AND。受限判定为解析集未覆盖全值域（DOM `<31`、DOW `<7`）。finding 例 `0 0 1 * 1`：两集均为真子集 → OR（1 号 **或** 周一）。`*` / `*/1` 全集走 AND，与「另一方受限」的 POSIX 行为一致。`CronFields` 形状不变，`Next`/`describeCron` 调用方零改动。

**非重开残余**：Vixie 以句法 `*` 判定受限；本实现把 `1-31` / `0-6` 视为全集 ≈ `*`，故 `0 0 1-31 * 1` 只在周一触发（POSIX 会每天触发）。方向仍是少调度，不是多跑/越权。无 `0 0 1 * 1` 单测（recommended F-003）。

## Findings（本条独立意见）

### F-001 · F-009 的 L2 校验器仍未接入生产渲染路径

- 严重度：low ｜ 建议：recommended ｜ 状态：open
- 与 [A-004](A-004-w9-self.md) F-001 同事实：`validatePermissions` 唯一调用方仍是 `permissions-inheritance.test.ts`。运行时 `effectivePermission` 已对 sourceless cascade deny，安全影响已消除。malformed 结构不会在加载期显式报错。不阻断本波闭合。

### F-002 · 恢复码 CAS 以秒级 `updated_at` 为令牌，同秒首次并发仍可能碰撞

- 严重度：low ｜ 建议：recommended ｜ 状态：open
- `UpdateRecoveryCodesIfUnchanged` 的 OCC 键是 `updated_at`（Unix 秒）。当 `now.Unix() == prevUpdatedAt.Unix()` 时 SET 不推进令牌，两路 `WHERE updated_at = ?` 都可成功 → 同码双花或不同码丢失更新可在该 1 秒窗口重现。常见双提交（上次写入在更早一秒）已被 CAS+重读覆盖，故 **不重开 D-002 F-006**。更稳的令牌是 `recovery_codes_hash` 本身或单调 version。不阻断本波。

### F-003 · 若干修复缺少针对原缺陷形状的回归锁

- 严重度：low ｜ 建议：recommended ｜ 状态：open
- 本会话全量/定向回归为绿，但未锁原缺陷形状的测试包括：`kernel.IsUniqueViolation` 双方言字符串（含 wrap）；`RecordLoginFailure` 并发计数；`AdvanceLastUsedStep` / 恢复码并发；job/scheduler panic recover；cron `0 0 1 * 1` 的 OR。F-002 nginx 与 F-008/F-009 运行时有测试/夹具覆盖。证据缺口，不是关闭声明不实。不阻断本波。

## 必改项汇总（required 列表）

无。D-002 的 12 条 required 本条均判定 **fixed**。本条 3 条 recommended 不阻断 S4 代码闭合。

## 与既有意见的异同

| 意见 | 关系 |
|------|------|
| A-001 fail | 原 12 条代码事实本条确认已被针对性改动；不改写 A-001 历史正文 |
| A-002 conditional | 清单调和与分级说明已由 D-002/D-003 消费；本条不重开清单争议 |
| A-004 self pass | **同意** 12/12 genuine fixed、回归全绿、F-009 L2 残余 recommended、go 暂挂纪律。补充：F-006 同秒 OCC（recommended F-002）；针对性回归锁不足（recommended F-003）。A-004「同码不可双用」在令牌推进时成立，同秒窗口需加注，不足以否定 F-006 闭合 |
| D-003 §6 | 本条即 cross 模式的 independent 腿（grok-4.6 · high · `/audit`） |

## 结论 + 建议给编排器/用户的下一步

- **verdict pass**：scope 内无未关闭 high required；12 条消费 ID 均为 genuine fixed；无到期 required 信息项阻断本闭合；回归与 A-004 主张一致。
- 本意见 **不修改** 目标 status/progress，**不恢复** VP-008 go 宣称。
- 建议 `/govern`：响应本条（将 A-001 消费 12 条标为 `fixed`）→ 书面决策是否恢复 VP-008 go（另写 D-00N）→ 勾选 S4 / 关门。3 条 recommended 可留后续波次或就地记下不修。
- I-003 仍 open non-blocking：本条已执行 grok 复核；是否将该信息项标 verified 由编排器询问用户。

## 声明

本意见不修改 status/progress；响应由 `/govern` 处理。
