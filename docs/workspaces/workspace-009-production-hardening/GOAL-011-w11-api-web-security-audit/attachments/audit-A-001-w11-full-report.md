---
title: W11 api/web 独立审计全文
status: recorded
created: 2026-08-22
updated: 2026-08-22
parent: GOAL-011-w11-api-web-security-audit
version: 0.1.0
---

# W11 · api + web 独立代码审计全文（2026-08-22）

> 正式摘要与 finding 台账：[03-audit/A-001-w11-independent.md](../03-audit/A-001-w11-independent.md)。本附件保留方法、对照旧审查、证据摘录与未入账项。

## 方法

- **指令**：独立审计 `apps/api` 与 `apps/web`；用户明确「不要加载任何 skills」。
- **主线**：grok-4.6 深读鉴权、上传、钱包、MFA、Postgres 方言、调度器、回收站、渲染层提交/搜索/路由。
- **并行**：3 个只读 explore 子代理（API auth；API handlers/modules；web renderer/auth）。
- **交叉**：HIGH/MEDIUM 由主线重读源码后入账；子代理补遗条目核实后再入账。
- **基线**：工作树 `dev`（相对 `origin/dev` ahead 2）；未把测试套件再跑一遍作为本意见的放行依据。
- **auditor 偏差**：非工作区默认 `grok build · grok-4.5 · high · /audit`。

## 对照 2026-08-10（C1–C8 / D1–D8）

| ID | 当时 | 现在 |
|----|------|------|
| C1 上传存储型 XSS | 客户端 MIME 落盘、无 Content-Disposition | **已修**：嗅探 + 活性内容门禁、`files.write`、owner-only GET、`attachment` + CSP sandbox |
| C2 refresh 竞态 | check-then-act + 前端并发 refresh | **双发已修**（原子 UPDATE + in-flight Promise）；家族吊销见 F-015 |
| C3 `APP_ENV` 默认 development | 裸跑 fail-open | **已修**：空 `APP_ENV` 启动失败 |
| C4 Bootstrap 只在 WasFresh | 半失败锁死 | **已修**：`NeedsBootstrap()` |
| C5 提交失败按钮卡死 | 无 try/finally | **已修** |
| C6 清空搜索不清除过滤 | 空 q 不写入 | **已修** |
| C7 未声明权限被拦 | 与默认放行矛盾 | **已修** |
| C8 路由 query 未达渲染层 | 硬编码空 query | **首屏已修**；同页 query 变化见 F-012 |
| D1 `roles: null` 清空角色 | PATCH 语义错 | **已修** |
| D2 登录无限流 + 用户名时序 | 无 limiter / 无 dummy bcrypt | **部分修**；锁定/禁用残余 = F-007 |
| D3 schema 匿名可读 | 信息泄露 | **仍在** = F-020 |
| D4 URIError | decodeURIComponent 无 catch | **已修**（`safeDecode`） |
| D5 迁移快照秒级粒度 | — | 本轮未复审为开放缺陷 |
| D6 defaultLocale `""` | 两层校验不一致 | **残余** = F-021 |
| D7 inputNumber 空→0；同文件无法重传 | — | 文件重传已修；空数字 **未清** = F-011 |
| D8 theme color-scheme | — | **已修** |

## Required 证据摘录

### F-001 Postgres EXISTS → int

```go
// users_repository.go CreateUserManagement
var exists int
if err := tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)`, user.Username).Scan(&exists); err != nil {
    return fmt.Errorf("check username: %w", err)
}
if exists == 1 {
    return ErrUsernameTaken
}
```

同文件 `DeleteUser` 已用 `var exists bool`。`scanUserListRow` 注释（R6）：postgres 出 native bool，sqlite 出 0/1，int 目标会让 users 列表 500。创建路径是漏网。调用方：`usersEntity.Create`、`importUsersCSV`。

### F-002 删除 204 / 快照失败

`resources.go` delete：`Entity.Delete` 成功后 `Trash.Record` 失败只记日志，仍 `WriteHeader(204)`。batch-delete 同类。命中回收站接入的资源（字典、定时任务）。

### F-003 MFA 第二因子

- 登录 MFA 分支（密码已过）不调用 `rateLimiter.record`。
- 每次 `BeginChallenge` → `CreateProof`；`mfa_proofs` 无过期 GC。
- `IncrementProofFailures`：`UPDATE ... SET fail_count = fail_count + 1 WHERE id = ?`，无 `fail_count < 5`。
- TOTP：6 位、`totpWindow = 1`。

### F-004 JWT × MFA 密钥

`NewService` HKDF 材料 = 当前 JWT secret。composition 只传入 `[]byte(secret)`，不传 `AuthJWTSecretPrevious`。轮换后密文解不开，`Required()` 仍 true。

### F-005 验证码并发

`ConsumeChallenge`：同一事务内 SELECT 再 DELETE。第二事务 DELETE 0 行不报错，`matched` 仍可在第一事务为 true。Postgres READ COMMITTED。

### F-006 钱包对账

`POST /api/wallet/reconcile`：`_ = Decode(&body)`；`wallet.read`。store：`accountID == ""` → `SELECT id FROM wallet_accounts`。jobs 测试用空字符串表示全库，属有意哨兵；垃圾 JSON 走同一路径不是。

## Recommended 要点

- **F-007**：未知用户 dummy bcrypt + 计入限流；锁定/禁用在密码前返回且不 record。约 6 次后快 401 vs 慢至 429。
- **F-008**：Restore = Get + restoreRow INSERT + MarkRestored CAS，三笔事务。
- **F-009**：`lastRun` map；手动 `Execute` 不写入；未知 handler → noop 记 `ran`。Compose 文档写明非多实例，Postgres 方言已一等。
- **F-010**：`refreshAccess` 网络错误返回 false 且不清 token；`restoreSession` 映射 `reauth`；AuthGate 画 `HOST_REAUTH_REQUIRED`。
- **F-011**：`NumberField` 清空写 `undefined`；`displayValue` → `coerceToKind("number")` 非有限 → 0。
- **F-012**：`useRecordSourcePrefill` effect 依赖无 `crud.route`。
- **F-013**：`urlEscape` 只替换空格与 `:`。
- **F-015**：已吊销 refresh 再刷新只返回 `ErrTokenRevoked`，不 bump `token_version`。
- **F-017**：`formulaSafe` 只认 `= + - @` 首字节。

## 抽查后扎实的部分

Refresh 旋转原子性与 JWT HMAC 方法校验；token_version 吊销；must-change-password 门；上传 XSS 加固与品牌/头像转码；钱包 Apply 不变量、乐观锁、幂等键；SMTP 主题控制字符拒绝；对象存储 ID 形态；CORS 默认关；metrics token `subtle.ConstantTimeCompare`；渲染层提交 try/finally；搜索清空。

## 建议实施顺序（供 S2 裁决，非本意见放行）

1. F-001 Postgres 创建用户
2. F-002 删除与快照同事务
3. F-003 MFA 第二因子限流 + 原子失败上限
4. F-004 MFA 密钥与 JWT 解耦
5. F-005 验证码单语句消费
6. F-006 对账 Decode 400 + 写权限
