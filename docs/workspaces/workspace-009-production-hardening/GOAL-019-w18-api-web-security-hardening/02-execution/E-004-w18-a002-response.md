---
id: E-004-w18-a002-response
doc: execution-entry
goal: GOAL-019-w18-api-web-security-hardening
date: 2026-09-22
status: recorded
parent: GOAL-001-production-hardening
version: 0.1.0
---

# E-004 · A-002 required findings 响应事实

## 已完成事实

- F-001：`apps/api/modules/mfa/store/repository.go` 新增 `HasActiveEnrollment()`；`server.Run` 在 listener/runtime 之前执行检查。查询失败或发现 active enrollment 时关闭 store 并返回明确错误。repository/server 测试覆盖 no-active、active、查询失败、资源清理和未启动断言。
- F-002：`apps/web/src/account/auth-client.ts` 的目标 URL 解析不再使用 realm-sensitive `instanceof`，按 Request `url` / URL `href` 解析并对不可解析对象 fail closed。测试覆盖跨 realm 同源/跨源 Request、跨 realm URL、auth endpoint 401 及原有字符串行为。
- F-003：已修正 canonical `goal-tree.md` 的正式树、状态表、frontmatter、残留 `@@` / 过期 `(0/6)`；勘误 E-002 的阶段算术为 2/6；`00-meta.md` 与 `03-audit.md` 的 I-003 分母统一为 Web 124/1516、并纳入浏览器 E2E；I-004 仍 open。

## 验证事实

- API：`go test ./...` 通过。
- Web：`npm run test` 为 124 个文件、1516 个测试通过；`npm run typecheck` 通过；新增真实 Chromium iframe realm E2E 通过（1 passed），覆盖跨 realm Request/URL 的同源与跨源鉴权头，以及 auth endpoint 401 不触发 refresh。
- 工作树：`git diff --check` 通过。

## 状态边界

以上是修复事实，不等同于 independent finding 已闭合。等待新的 clean-context Reviewer 对 F-001～F-003 复核；在其通过前 I-004、S6 和目标 `done` 均保持未完成。
