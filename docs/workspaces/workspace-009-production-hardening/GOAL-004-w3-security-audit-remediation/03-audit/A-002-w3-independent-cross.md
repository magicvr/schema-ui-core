---
id: A-002
goal: GOAL-004-w3-security-audit-remediation
title: W3 八项修复独立交叉审计
source: independent
auditor: grok-build · grok-4.5 · high · audit skill
date: 2026-08-11
verdict: conditional
status: recorded
---

# A-002 · W3 八项修复独立交叉审计（independent）

## 条目头

| 字段 | 值 |
|------|-----|
| **source** | independent |
| **auditor** | grok-build · grok-4.5 · high · 执行 `audit`（VP-009 provider） |
| **类型** | execution-facts / close-out 前交叉复核 |
| **scope** | GOAL-004 W3 八项修复实现正确性与回归风险（batch-delete 原子、recordSource `//`、nginx、登录限流、admin 委托边界、logo 反斜杠 + refresh 竞态、Serve fail-closed） |
| **verdict** | **conditional** |
| **工作区** | `workspace-009-production-hardening`（Root `GOAL-001-production-hardening`；canonical 已校验；`shared_materials_catalog: none`） |

## 范围与区间

- 代码：`apps/api`（authsession 仓库、handler users/roles/resources/auth/rate_limit、composition、settings logo 校验）+ `apps/web`（request-construction、auth-client、nginx.conf）。
- 过程：`02-execution/E-001-w3-remediation.md`、`03-audit/A-001-w3-self.md`、`01-decision/D-001-w3-remediation-scope.md`、`00-meta.md`。
- 信息项：I-001/I-002/I-003 均已 verified（本 scope 无到期 open required 信息门禁）。
- **只读审计**：不改目标 status/progress/goal-tree；探针测试文件审后删除，未留生产改动。

## 成果（有证据 · 本人核对）

| # | 项 | 独立项 verdict | 证据摘要 |
|---|----|----------------|----------|
| 1 | batch-delete 原子化 | **fail**（见 F-001） | 单事务 + 先守卫后删除路径存在；**多 admin 同批** last-admin 守卫失效 |
| 2 | recordSource 拒 `//` | **pass** | `buildRecordSource` 仅 `isRelativeProtocolUrl`；本地 5 用例 |
| 3 | nginx 8m + 头 + `/api/` | **pass** | 与 `maxUploadBytes = 8<<20` 对齐；安全头与前缀匹配到位 |
| 4 | 登录限流 | **pass**（边界 note） | 可信 peer 才信 X-Real-IP；IP\|username；成功 clear；容量驱逐 |
| 5 | admin 委托边界 | **pass**（范围外遗漏 note） | 改密/demote → `ADMIN_ACCOUNT_FORBIDDEN`；契约/双语已钉 |
| 6 | logo `\` + refresh 竞态 | **pass** | API 两处拒 `\`；`refreshGeneration` + 竞态测 |
| 7 | Serve fail-closed | **pass** | 非 `ErrServerClosed` → `os.Exit(1)` |
| 8 | 契约/回归不破坏 | **pass**（F-001 测缺口） | `ADMIN_ACCOUNT_FORBIDDEN` 入 frozen；users/roles 既有测绿；**缺「整批清空 admin」用例** |

## 逐项审计

### 1 · batch-delete 原子化 — **fail**

**实现核对**

- `DeleteUsersBatch`（`users_repository.go:266-316`）：`withTx` 内先循环存在性 / self / per-id last-admin，再统一删 refresh + users。
- `DeleteRolesBatch`（`roles_repository.go:184-218`）：同样两阶段；system / in-use / not-found 任一失败整批回滚 — **roles 侧正确**。
- `BatchDeleter`（`resources.go:62-68`）+ `batchDelete`（`resources.go:641-668`）：users/roles 实体实现 `DeleteBatch` 委托仓库；未实现实体仍顺序删（D-001 明确范围）。
- 既有测：`TestDeleteUsersBatchAtomicRollback`、`TestRolesRepositoryBatchDeleteAtomicRollback`、`TestUsersBatchDeleteAtomicRollbackHTTP` 覆盖单 last-admin、self、not-found、in-use、system、dedupe、HTTP 回滚。

**缺陷（TOCTOU / 守卫语义 · F-001）**

`countAdminUsersExcluding(tx, id)`（`users_repository.go:318-327`）只排除**当前** id。当批内包含多个 admin 且合计覆盖全部 admin 时：

1. 对 admin A：库内仍有 B → `other > 0` → 通过  
2. 对 admin B：库内仍有 A → `other > 0` → 通过  
3. 随后两阶段删除把 A、B 一并删掉 → **系统 0 个 admin**

本人探针（临时 `TestAuditProbeBatchDeleteAllAdmins`，审后删除）：

- 输入：仅 `user-admin` + `user-admin2` 两个 admin，`DeleteUsersBatch([both], "user-external")`
- 结果：`err=<nil> deleted=2`，`remaining admin links=0`

单条 `DeleteUser` 无此问题（第二次调用时 other=0）。**整批原子修复引入了单删不具备的 last-admin 集体绕过。**  
任意持 `users.write` 的非 admin（HTTP 层 batch 已用 um 角色测过）即可在 ≥2 admin 时清空管理员，属可用性/特权完整性破坏。

**非问题**

- 事务回滚本身：单 id last-admin / self / not-found 路径原子性成立（测绿）。
- roles 批删：无「至少保留一个 X」类跨 id 不变量，守卫可组合。

### 2 · recordSource.url 拒绝 `//` — **pass**

- 旧旁路：`startsWith("/")` 会放行 `//host`。
- 现：`request-construction.ts:488-490` 仅 `isRelativeProtocolUrl`；`286-292` 显式拒 `//`，path 再走 `PROTOCOL_URL_RE`（拒 `\`）。
- 与 rowAction（`298`）一致；`//evil` / `/\evil` / `https://` 本地测拒绝（`request-construction.test.ts:35-58`）。
- **Bearer 外泄路径对本入口关闭。** 未发现绕过。

### 3 · nginx — **pass**

- `apps/web/nginx.conf:21-28`：`client_max_body_size 8m`（= `upload.go` `maxUploadBytes = 8<<20`）；`server_tokens off`；nosniff / DENY / Referrer / 保守 CSP。
- `location /api/`（`45-54`）前缀匹配，不再吃 `/apix`；`X-Real-IP $remote_addr` 与限流可信 peer 设计一致。
- residual（CSP 非生产精细、无 HSTS）与 00-meta / D-001 一致，不升 required。

### 4 · 登录限流 — **pass**（N-001/N-002 note）

- `loginClientIP`（`rate_limit.go:102-110`）：仅 `trustedReverseProxy`（loopback/private）时信 `X-Real-IP`，否则 peer — 防客户端伪造。
- 键 `IP|username`（`auth.go:71`）；失败 `record`；成功 `clear`（`91-92`）；capacity 驱逐最旧（`65-76`）。
- 单测：`TestLoginRateLimiterUnit`、`TestLoginClientIPTrustsXRealIPOnlyFromTrustedPeer`、既有 `TestLoginRateLimit` 绿。
- **边界**：`allow` 与 `record` 非原子，并发可略超 max（best-effort 可接受）；全 private 网段 peer 可信依赖「前置反代」部署假设（N-002）。多实例不共享 — 00-meta residual。

### 5 · 委托边界 — **pass**（N-003 recommended）

- `authorizeAdminTargetBoundary`（`users.go:302-321`）：非 admin 且非 self → admin 目标改密 / demote 拒 `ADMIN_ACCOUNT_FORBIDDEN`。
- 顺序：先 `authorizeRoleAssignment`，再 target 边界，再 `UpdateUser`（store 仍有 last-admin/self）。
- `error_contract_test.go:45` frozen + `errorcatalog.go:72` 中英条目。
- `TestUsersAdminTargetBoundary`：改密 403、无 assign demote、有 assign demote、非 admin 目标可管。
- **范围外遗漏（recommended）**：D-001 仅冻结改密/demote；非 admin 仍可 **DELETE** admin（非 last）及改 name 等，未在本波封。不构成对本波声明的 fail，但边界不对称。

### 6 · logo 反斜杠 + refresh 竞态 — **pass**

- `normalizeLogoURL`（`repository.go:287`）与 `validateLogoURL`（`configuration.go:112`）同源拒 `\` / `//` / 空白；`repository_test.go` 含 `/\evil.com` 等。
- `auth-client.ts`：`logout` 先 `bumpRefreshGeneration`（`247`）再清 token；`doRefresh` 写 token 前多次 `isCurrentGeneration`（`118`/`130`/`151`）；`auth-client.test.ts:241-271` 覆盖 in-flight 不写回。
- 未发现 generation 绕过写回。

### 7 · Serve fail-closed — **pass**

- `composition.go:265-271`：`Serve` 错误且非 `ErrServerClosed` → log + `os.Exit(1)`。
- 语义正确：半死实例不继续挂；依赖 compose restart。`os.Exit` 跳过 fx 优雅清理属已知权衡，可接受。

### 8 · 契约与回归 — **pass with gap**

- 聚焦 `go test`：`authsession` + `handler`（含 batch/admin/rate/error contract）本人复跑 **ok**。
- F-001 场景**无**回归用例 → 现有绿不能证明 last-admin 整批不变量。

## Findings

### F-001 · required · high · batch-delete 可一次删光全部 admin（last-admin 集体绕过）

| 项 | 内容 |
|----|------|
| **文件:行号** | `apps/api/internal/modules/authsession/users_repository.go:294-301`（per-id `countAdminUsersExcluding`）；`318-327`（计数语义）；对照 `304-311` 第二阶段统一删除 |
| **失败场景** | ≥2 个 admin 时，一次 `DeleteUsersBatch`/`POST /api/users/batch-delete` 包含全部 admin id → 守卫全过 → 提交后 0 admin。持 `users.write` 的委托账号即可触发。 |
| **与成功标准** | P0「整批原子 + last-admin 守卫」原子性成立，**不变量不成立**；E-001/A-001 写「同一 existence/self/last-admin 守卫」对**批级 last-admin** 不充分 |
| **确认度** | **高**（源码推理 + 本地探针：`err=nil deleted=2 remaining admins=0`） |
| **建议闭合** | **fixed**：批级计算「删除后剩余 admin 数」——例如 `totalAdmins - adminsInBatch >= 1`，或守卫阶段模拟删除集合后再计数；补 `TestDeleteUsersBatchRejectsRemovingAllAdmins` + HTTP 层等价用例。**禁止**仅口头关闭。 |

### F-002 · recommended · low · 委托边界未覆盖 admin 删除/改名

| 项 | 内容 |
|----|------|
| **文件:行号** | `users.go:188-190`（`Delete` 无 target 边界）；`302-321`（仅密码/roles） |
| **说明** | 与 D-001 范围一致；对称性上非 admin 仍可删非 last admin。后续波次可扩 `ADMIN_ACCOUNT_FORBIDDEN` 到 Delete/敏感字段。 |

### F-003 · recommended · low · 限流依赖「private peer = 可信反代」

| 项 | 内容 |
|----|------|
| **文件:行号** | `rate_limit.go:114-119`（`IsPrivate` 全信任） |
| **说明** | API 若无反代直暴露于 RFC1918，同网段可伪造 `X-Real-IP` 分片/绕过桶。compose+nginx 路径正确；部署文档宜写死「必须前置反代或收紧 trust」。 |

### N-001 · note

限流 `allow`/`record` 并发可略超 max；进程内 best-effort — 与 residual 一致。

### N-002 · note

CSP `'unsafe-inline'` style、无 HSTS — 00-meta residual。

## 必改项汇总

| ID | 级别 | 一句话 |
|----|------|--------|
| **F-001** | **required** | 批删 last-admin 必须按**集合**判定，禁止一次清空全部 admin；补回归测 |
| F-002 | recommended | 是否将 Delete/改名纳入 admin 目标边界（下波） |
| F-003 | recommended | 部署层收紧 trusted proxy 假设 |

**开放 required 数：1**（F-001）

## 与既有意见的异同（A-001 self）

| 点 | A-001 self | A-002 independent |
|----|------------|-------------------|
| verdict | pass | **conditional** |
| batch-delete | ✅ 原子 + 守卫 | 原子 ✅；**批级 last-admin ❌**（F-001） |
| recordSource / nginx / 限流 / 委托改密 demote / logo+refresh / Serve | ✅ | **同意 pass** |
| 开放 required | 0 | **1（F-001）** |
| residual notes | N-001～N-004 | 部分吸收为 F-003 / N-*；不重复升 required |

## 结论 + 建议给编排器/用户

**verdict: conditional** — 在 F-001 以 `fixed` / `accepted-residual` / `user-overruled` 合法闭合前：

- **不得**将 GOAL-004 标 `done`（P-003 开放必改门禁）；
- **不得**仅凭 A-001 无条件放行 W3 关门。

建议 `/govern`：

1. 优先 **fixed** F-001（批级 last-admin + 测）；
2. F-002/F-003 可记入下波或 residual；
3. 复审可再开 `/audit` 关闭 F-001。

Root `GOAL-001` / VP-009 保持长期 active 程序容器语义，不随本波关门。

## 声明

本意见 **source: independent**，**不修改** status / progress / 方案正文 / goal-tree。响应与闭合由 **`/govern`** 处理。
