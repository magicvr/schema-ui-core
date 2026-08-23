---
id: A-001-w11-independent
goal: GOAL-011-w11-api-web-security-audit
doc: audit-entry
record_id: A-001
source: independent
scope: apps/api + apps/web 当前实现
verdict: fail
status: recorded
created: 2026-08-22
updated: 2026-08-22
parent: GOAL-011-w11-api-web-security-audit
version: 0.1.0
---

# A-001 · W11 api/web 独立安全审计（2026-08-22）

- **source**：independent
- **auditor**：grok-4.6 会话主线 + 3 并行 explore 子代理（API auth / API handlers-modules / web renderer-auth）+ 主线逐条源码交叉复核。用户本轮明确禁止加载 skills，**未走** `/audit` skill。**非工作区默认 grok-4.5 `/audit` provider**，偏差见 `00-meta.md` I-003。
- **类型 / scope**：ad-hoc（用户直接指令）· `apps/api`（Go）与 `apps/web`（React/TS）当前实现的 bug 与安全问题
- **verdict**：**fail**（3 条 HIGH required + 3 条 MEDIUM required 开放；13 recommended；6 informational）

## 范围与区间

2026-08-22 工作树快照（`dev`，相对 `origin/dev` ahead 2）。方法与完整证据见全文附件：[attachments/audit-A-001-w11-full-report.md](../attachments/audit-A-001-w11-full-report.md)。

对照 2026-08-10 审查（C1–C8 / D1–D8）：C1–C8 与 D1、D4–D8 **已修**；D2 部分修（锁定/禁用时序残余 = F-007）；D3 仍在（F-020）；D6 残余（F-021）。

## 结论摘要

经 W1–W10 多轮加固后，认证旋转、上传 XSS、钱包账本不变量、CORS 默认关闭、metrics 恒定时间比较等基础面扎实。本轮问题集中在：**Postgres 方言漏网**（创建用户）、**回收站删除/快照非原子**、**MFA 第二因子可在线穷举**、以及 JWT 轮换与 MFA 密钥耦合、验证码并发一次性、钱包对账写操作挂在 `wallet.read`。

**P0=0 · P1=3 · P2 required=3 · recommended=13 · informational=6。开放 required = 6。**

## Findings（required）

### F-001 · Postgres 上 `CreateUserManagement` 把 `EXISTS` 扫进 `int`，创建/导入用户 500
- 严重度：**high** ｜ 建议：**required** ｜ 状态：open
- **文件**：`apps/api/internal/modules/authsession/users_repository.go:86-91`
- **问题**：`var exists int` + `SELECT EXISTS(...)`。Postgres/`pgx` 的 `EXISTS` 返回 bool；R6 已把列表路径改为 `bool`（`scanUserListRow` 注释写明），删除路径也已是 `bool`，**创建路径漏了**。Scan 失败发生在 INSERT 之前，所有创建/CSV 导入在 Postgres 上 500。
- **建议**：改为 `var exists bool`（与 `DeleteUser` 一致），并加 Postgres 集成测试覆盖 create + import。

### F-002 · 删除成功但回收站快照失败仍返回 204
- 严重度：**high** ｜ 建议：**required** ｜ 状态：open
- **文件**：`apps/api/internal/handler/resources.go:741-754`（batch-delete 同类）
- **问题**：`Entity.Delete` 提交后 `Trash.Record` 失败只 `slog.Error`，HTTP 仍 204。字典类型/条目、定时任务会不可恢复地丢失。
- **建议**：删除与快照同一事务；快照失败则整次删除失败（或补偿回滚）。

### F-003 · 偷到密码后可在线穷举 TOTP（无限 proof、无限流、失败计数非原子）
- 严重度：**high** ｜ 建议：**required** ｜ 状态：open
- **文件**：`apps/api/internal/handler/auth.go:144-156`（MFA 成功路径不 `record` 限流）；`apps/api/internal/modules/mfa/service.go:110-149`；`apps/api/internal/modules/mfa/store/repository.go:282-291`（`fail_count = fail_count + 1` 无 `fail_count < 5` 守卫）；`mfa_proofs` 无过期清理（对比 captcha 懒清理）
- **问题**：密码通过后每次签发新的 5 分钟 proof（每条 5 次猜测），proof 数量无上限；`/api/auth/mfa/verify` 无限流。6 位 TOTP ±1 窗口约 3 个有效码。并发失败猜测都能通过 check-then-act。验证码默认关。
- **建议**：第二因子限流；创建 proof 前清理过期行并对「密码已过、等待 MFA」限流；失败计数做成 `UPDATE ... AND fail_count < 5`。

### F-004 · JWT 轮换使已存 TOTP 密文不可解，MFA 用户全锁
- 严重度：**medium**（运营/可用性，命中已落地的 VP-016 轮换）｜ 建议：**required** ｜ 状态：open
- **文件**：`apps/api/internal/modules/mfa/service.go:61-67`（HKDF 只用当前 JWT secret）；`apps/api/internal/composition/composition.go:357-359`（不传入 previous）
- **问题**：VP-016 已做 JWT 双密钥验证；MFA 仍只用当前 secret。轮换 `AUTH_JWT_SECRET` 后 `Required()` 仍为 true，登录永远卡在第二因子。
- **建议**：MFA 用独立密钥；或轮换窗口内同时尝试 previous secret 解密。

### F-005 · 验证码一次性消费在 Postgres 上不抗并发
- 严重度：**medium** ｜ 建议：**required** ｜ 状态：open
- **文件**：`apps/api/internal/modules/logincaptcha/store/repository.go:57-78`
- **问题**：`SELECT` 再 `DELETE`，无行锁/`DELETE ... RETURNING`。READ COMMITTED 下两事务可都读到同一行并都算验证成功。一条算术题可放行两次登录。
- **建议**：单语句消费（`DELETE ... WHERE id=? AND expires_at>? AND answer_hash=? RETURNING`）或 `SELECT ... FOR UPDATE`。

### F-006 · 钱包对账：坏 JSON 静默变成全库对账；提交/取消/重试只需 `wallet.read`
- 严重度：**medium** ｜ 建议：**required** ｜ 状态：open
- **文件**：`apps/api/internal/handler/wallet.go:295-353`；`apps/api/internal/modules/wallet/store/repository.go:659-674`（空 `accountID` 列出全部账户）
- **问题**：`_ = json.NewDecoder(...).Decode(&body)` 丢弃错误。空字符串全库对账在 jobs 测试中是有意哨兵；**解码失败静默当全库**不是。submit/cancel/retry 挂在 `wallet.read`，只读角色可排队跑全账本并取消/重试（jobs 按 actor 隔离，非跨用户 IDOR）。
- **建议**：Decode 失败返回 400；写操作改挂 `wallet.write` / `wallet.adjust`；全库对账用显式哨兵字段。

## Findings（recommended）

| F-ID | 严重度 | 摘要 | 位置 |
|------|--------|------|------|
| F-007 | med | 锁定/禁用账户跳过 dummy bcrypt 且不计入限流（D2 残余；约 6 次可枚举用户名） | `auth.go:175-190`；`handler/auth.go:129-137` |
| F-008 | med | 回收站 Restore：业务 INSERT 与 `MarkRestored` 非同一事务；崩溃留下「行已在、快照仍可恢复」 | `recyclebin/service.go:63-88` |
| F-009 | med | 调度器 `lastRun` 仅内存；重启同分钟、手动 run+tick、多 API 副本会双跑；未知 handler 静默 noop 且记 `ran` | `scheduledtasks/scheduler.go:34-36, 123-137` |
| F-010 | med | `restoreSession` 把网络失败的 refresh 映射为 `reauth-required`（token 未清） | `auth-client.ts` + `AuthContext.tsx` + `main.tsx` AuthGate |
| F-011 | med | 空 `inputNumber` 经 `coerceToKind("number")` 变成 0（D7 未清干净）；钱包 `amountDelta` 空值会提交 0 | `form-controls.ts:218-221` |
| F-012 | med | `recordSource` 预填 effect 不含 `crud.route`；同页只改 query 会留下上一份记录，保存可能写错行 | `render.tsx:useRecordSourcePrefill` |
| F-013 | low–med | otpauth URL 只转义空格和冒号；用户可改的 `name` 含 `?&#` 会破坏 authenticator URI | `mfa/totp.go:76-84` |
| F-014 | low | 超限上传 `MaxBytesReader` 失败映射为 500 `STORAGE_UNAVAILABLE` 而非 413 `FILE_TOO_LARGE` | `handler/upload.go` |
| F-015 | low–med | Refresh 重放不灭会话族（C2 残余）：旋转原子，但已吊销 refresh 再用不吊销兄弟 token；赢者会话默认 720h | `auth.go:266-278` |
| F-016 | low | 回收站恢复字典条目丢掉 `badgeStyle` | `recyclebin/service.go:182-194` |
| F-017 | med | `formulaSafe` 漏 `\t`/`\r` 前缀（OWASP CSV 注入变体） | `handler/export.go:231-239` |
| F-018 | low–med | 头像配额 check-then-act，无上传那把 `quotaMu` | `handler/account_avatar.go:62-74` |
| F-019 | low | 前端 custom action（下载/导出）绕过 `executeAction`；硬门禁仍在 API | `render.tsx:459-461` |

## Findings（informational）

| F-ID | 摘要 |
|------|------|
| F-020 | `GET /api/schema/{pageId}`、manifest、bootstrap 仍匿名（D3；可能为产品选择） |
| F-021 | repository 允许 `defaultLocale: ""`，configuration 层拒绝（D6 残余） |
| F-022 | 受信代理 CIDR 默认仅 IPv4 loopback，无 `::1` |
| F-023 | 前端仍处理 423 `ACCOUNT_LOCKED`；后端锁定已改为 401（W7 F-009），该 UI 路径为死代码 |
| F-024 | `admin.data-permission` 的 `Scoper` 未接到生产资源（v1 文档化未完成，非回归） |
| F-025 | 用户 `roles` JSON 与 `user_roles` 不一致时 `UserByUsername` 失败，该用户登录 500 |

## 已核实非缺陷 / 已修（不入账为开放 finding）

C1 上传 XSS；C2 旋转双发（原子 UPDATE + 前端 in-flight 去重；家族吊销为 F-015 残余）；C3 `APP_ENV` fail-closed；C4 `NeedsBootstrap`；C5 提交 try/finally；C6 清空搜索；C7 未声明权限默认放行；C8 路由 query 首屏注入（同页 query 变化为 F-012）；D1 `roles: null`；钱包 `Apply` 不变量与乐观锁；SMTP 主题 CR/LF 拒绝；对象存储 ID 形态；metrics token 恒定时间比较。

## 必改项汇总（required）

F-001、F-002、F-003、F-004、F-005、F-006 共 **6** 条，全部 **open**。

## 结论 + 建议下一步

- verdict **fail**：6 条 required 开放，其中 3 条 HIGH。独立意见不改动目标 status/progress（P-003）。
- 建议下一步：用户书面裁决 I-002（修复范围；是否暂挂 VP-008 go）。优先顺序建议：F-001 → F-002 → F-003 → F-004 → F-005 → F-006。
- I-003：本会话不是工作区默认 grok `/audit` provider；关门前是否追加该腿由用户决定。
