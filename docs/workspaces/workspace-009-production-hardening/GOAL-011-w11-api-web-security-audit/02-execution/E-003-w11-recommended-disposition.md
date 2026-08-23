---
id: E-003-w11-recommended-disposition
goal: GOAL-011-w11-api-web-security-audit
status: done
created: 2026-08-22
updated: 2026-08-22
parent: GOAL-011-w11-api-web-security-audit
version: 0.1.0
---

# E-003 · W11 recommended 处置（F-007～F-019）（2026-08-22）

## 事实

按 [D-002](../01-decision/D-002-w11-scope-and-go-hold.md) §4：真实缺口就地修正；纯设计取舍/误报逐条留痕（fixed/overruled 有据）。与 W9 E-005 / W10 D-003 先例一致：**作废必须给出可核对的源码/契约依据，不得"假修复"**。

## fixed（11 条）

### F-007 · 锁定/禁用跳过 dummy bcrypt 且不计入限流（fixed）

- `internal/auth/auth.go` `Login`：locked/disabled 分支在返回前执行 `VerifyPassword(timingDummyHash, password)`——与未知用户/密码错消耗相同 bcrypt 时间，用户名枚举时序通道关闭（D2 残余）。
- `handler/auth.go` login：`ErrAccountLocked` / `ErrAccountDisabled` 分支改为 `rateLimiter.record(limiterKey)` 后再回 401——锁定/禁用探测同样计入 IP|username 失败桶。
- 回归：既有 `TestRefreshRejectsLockedAccount` 与 handler locked-login 用例不变（行为面 = 401 UNAUTHORIZED 不变）。

### F-008 · 回收站 Restore 业务 INSERT 与 MarkRestored 非同一事务（fixed）

- 复用 F-002 的 caller-owned-tx seam（先例：wallet `ReconcileOnceTx`）：
  - `recyclebin/store`：`MarkRestoredTx(ctx, tx, id, now)`。
  - `datadictionary/store`：`CreateTypeTx(ctx, tx, t)` / `CreateEntryTx(ctx, tx, e)`。
  - `scheduledtasks/store`：`CreateTaskTx(ctx, tx, t)`。
  - `recyclebin.Service`：`NewService(..., runner)`（runner = kernel store，结构满足 `recyclestore.TxRunner`）；`Restore` 在**一个事务**内执行 restoreRowTx + MarkRestoredTx。
- 崩溃窗口消除：不再存在「行已恢复、快照仍可恢复」的中间态。
- 回归：`service_test.go` `TestRestoreAtomicityRollsBackOnFailedMark`——同一事务内行 INSERT 成功 + mark 失败 → 整体回滚（无行、快照仍 unrestored）。既有 restore 冲突/DICT_KEY_NOT_FOUND 语义保持不变。

### F-009 · 调度器未知 handler 静默 noop 且记 `ran`（部分 fixed + 保留设计）

- fixed：`scheduler.go` `Execute`：未知 handler 不再静默降级 `system.noop` 并记 `ran`——改为记录 **failed** run（detail 含 handler 名）并返回错误；运行历史不再掩盖配置错误。需要 handler 目录的 UI 已存在（`/api/scheduled-tasks/handlers`）。
- 保留设计（留痕）：`lastRun` 仍为进程内存去重——scheduler 文件头与 compose 文档已写明 best-effort 单实例语义（无多副本高可用承诺）；做成 DB 持久化 lastRun 需要迁移 + 租约语义，属超出本波的特性变更，作为 residual 移交（见 A-002）。

### F-010 · web `restoreSession` 网络失败映射 reauth-required（fixed）

- `auth-client.ts` `restoreSession`：刷新失败后**再试一次**（瞬时失败时 token 仍在——W15-F01 语义），第二次仍失败才返回 `reauth`；真正的 401/403 首次即清 token 语义不变。
- 回归：`auth-client.test.ts`——双失败仍 reauth + 保留 token；新增「瞬时抖动后恢复 session」（reject 一次 → rotate 成功 → kind session）。

### F-011 · 空 `inputNumber` 经 coerce 变 0（fixed）

- `form-controls.ts` `coerceToKind("number")`：`""`/`undefined`/`null` → `undefined`（提交层缺失字段 = PATCH 不动原值）；非数值垃圾串仍 fail-closed 到 0（既有 pinned 测试保持）。
- 效果：清空钱包 amountDelta 等数字字段不再提交 0。既有 defaultValue 路径（空 → defaultValue）不变。

### F-012 · `useRecordSourcePrefill` deps 缺 route（fixed）

- `render.tsx`：effect 依赖改为序列化 route key（`JSON.stringify(crud?.route)`）——App 每次渲染重建 context 对象，直接依赖对象身份会无限重取；内容级 key 只在 query/params 真正变化时重跑 prefill，同页只改 query 不再残留上一份记录（保存写错行风险消除）。

### F-013 · otpauth URL 只转义空格/冒号（fixed）

- `mfa/totp.go` `urlEscape`：改为对非 RFC3986 unreserved 字符全部百分号编码（`?`/`#`/`&` 等不再能注入 URI query/fragment）。

### F-014 · 超限上传 500 STORAGE_UNAVAILABLE（fixed）

- `handler/upload.go`：`io.ReadAll` 错误经 `errors.As(*http.MaxBytesError)` → 413 `FILE_TOO_LARGE`；其余仍 500。

### F-016 · 恢复字典条目丢 badgeStyle（fixed）

- `recyclebin/service.go` `dictEntryFromPayload`：补 `BadgeStyle` 字段（快照含该字段，恢复同列写回）。

### F-017 · `formulaSafe` 漏 `\t`/`\r`（fixed）

- `handler/export.go`：函数前缀表加 `'\t'`、`'\r'`（OWASP CSV 注入变体前缀）。

### F-018 · 头像配额 check-then-act（fixed）

- `handler/account_avatar.go`：配额检查与 store 写入在包级 `avatarQuotaMu` 临界区内——并发上传不能同时通过计数检查而超帽。

## overruled（有据，2 条）

### F-015 · Refresh 重放不灭会话族（overruled——实现会造成已设计流程回归）

- **主张**：已吊销 refresh 重放时不吊销兄弟 token，赢者会话存活至 720h TTL。
- **依据**：`apps/web/src/account/auth-client.ts` `doRefresh` 的非 2xx 分支显式依赖「重放只返回 ErrTokenRevoked 而不产生连带」：跨标签轮换（A-002 F-003）中，落败标签用旧 token 重试刷新，胜者标签刚写入新 token；若重放即烧家族，正常双标签并发刷新会把两边的会话全部杀掉（含刚旋转出的新对）——把设计的原子旋转流程变成登出触发器。服务端已具备的防线：旋转原子（guarded UPDATE，双发已修）、重放只影响自身、refresh TTL 为有意会话寿命设计（720h 上限）。
- **处置**：不改代码（auth.go 保持单 token 吊销语义）；残余登记在案：若未来引入「重放 = 窃取」的威胁模型升级，需先改客户端跨标签重试协议再启用家族吊销（复审触发 = 安全模型变更）。

### F-019 · 前端 custom action 绕过 `executeAction`（overruled——UI 对齐项，硬门禁在 API）

- **主张**：下载/导出等 custom action 不经前端 executeAction 权限面。
- **依据**：`renderer/render.tsx` `runCustomAction` 的 handler 全部映射到受服务端权限门禁的端点（`export.users/roles` → data.export；`library.*` → files.read/delete），API 返回 403/401 直接透传前端错误面；custom action 是协议白名单 handler 名，页面 schema 不为其声明 permission target（D5 契约形状），前端 executeAction 语义不适用。审计原文亦注明「硬门禁仍在 API」。
- **处置**：不改代码；如需对齐 UX（按钮级禁用），属下波 UI 工作项，非正确性缺陷。

## 处置汇总

| 范围 | 条数 | fixed | overruled（有据） |
|------|------|-------|-------------------|
| required | 6 | 6（E-002） | 0 |
| recommended | 13 | 11（F-007/008/009/010/011/012/013/014/016/017/018） | 2（F-015/F-019） |

## 回归验证

- API：`go build ./...` 0；`go vet ./...` 0；`go test ./...` 全绿（含新增回归锁；PG 门控用例无 PG 环境自动 skip）。
- Web：`npx vitest run` 76 files / **1085 tests** 全绿；`npx tsc -b` 0 错误。
- Git checkpoint：`72a5397`（37 files；S3 全部实施 + 回归锁 + 台账）。