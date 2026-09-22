---
id: A-002-w18-clean-context-reviewer
doc: audit-entry
goal: GOAL-019-w18-api-web-security-hardening
date: 2026-09-22
source: independent
auditor: gpt-5.6-sol / high / clean-context REVIEWER
type: cross-audit-independent
scope: W18 API/Web security fixes, production assembly, regression evidence, governance consistency
verdict: fail
open_required: 3
parent: GOAL-001-production-hardening
version: 0.1.0
---

# A-002 · W18 clean-context independent review

## 结论

独立 Reviewer 判定 `fail`。已核对 API JWT 启动校验、Web 鉴权修复、现有回归命令与 diff 检查，但发现 3 个开放 required findings；S6 不得放行，目标不得关门。

## Findings

### F-001 · required / MAJOR · `serve` 对 active MFA enrollment 未 fail-closed

`compiled.PersistenceCatalog()` 包含 `user_mfa` 迁移，而 `server.Run` 的生产装配对 MFA verifier 传入 nil。若下游连接启用 MFA 的数据库，`Authenticator.Login` 会退化为密码成功即签发 token。建议在打开 store 后、listener/runtime 启动前查询是否存在 active enrollment；查询错误或发现 active enrollment 均拒绝启动，并补回归测试。

### F-002 · required / MAJOR · Web `Request` URL 判断跨 realm 不安全

当前 `input instanceof Request` 只识别当前 realm。来自 iframe 或其他浏览器 realm 的 `Request` 会落入字符串化路径，被解析为当前源的 `[object Request]`，从而可能向真实跨源请求注入 Bearer 与 refresh token。应改为跨 realm 安全的 URL 提取，并加入真实跨 realm（或等价浏览器级）同源/跨源回归覆盖。

### F-003 · required / MAJOR · 治理台账不一致

`goal-tree.md` 的正式树缺少 GOAL-019，文件末尾残留 `@@` 与过期的 `(0/6)` 片段；其 frontmatter 日期、`03-audit.md` 的 I-001～I-003 状态、E-003 的阶段/进度叙述也与 `00-meta.md` 的 5/6、83% 不一致。必须修正 canonical tree、执行记录和审计索引，并保持 S6/I-004 未完成。

## 已核对证据

- API：JWT secret 非 development 缺失/弱值在 `Run` 启动前被拒绝；`go test ./... -count=1` 通过。
- Web：现有同源/跨源字符串与 Request 回归通过；Web 全量 `124 files / 1510 tests`、类型检查和 `git diff --check` 通过。
- Reviewer 复现了跨 realm `Request` 风险，并核对了 persistence/MFA migration 与 nil verifier 的装配链。

## 未能替代生产环境确认的事项

实际生产数据库内容、外部消费者和真实 PostgreSQL 启动未在本审计会话中执行；这些不改变 F-001 的 fail-closed 要求。

## 响应要求

修复 F-001～F-003 后重新执行 S5/S6 所需验证，并追加独立审计响应条目；在 required findings 合法闭合前不得将 `I-004` 标为 verified 或将目标标记为 `done`。
