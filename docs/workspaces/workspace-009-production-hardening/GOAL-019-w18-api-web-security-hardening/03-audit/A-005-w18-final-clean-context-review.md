---
id: A-005-w18-final-clean-context-review
doc: audit-entry
goal: GOAL-019-w18-api-web-security-hardening
date: 2026-09-22
source: independent
auditor: Codex REVIEWER / clean-context / read-only
type: cross-audit-independent
scope: S6 final gate, A-002/A-004 finding closure, API/Web production boundaries, regression evidence and governance consistency
verdict: conditional
open_required: 2
parent: GOAL-001-production-hardening
version: 0.1.0
---

# A-005 · W18 final clean-context independent review

## 结论

`verdict: conditional`，开放 required = 2。API 与 Web 安全修复均有代码、针对性测试及独立复跑证据；A-002 的代码 findings 与 A-004 的真实浏览器 finding 已实质修复。但审计投影和 independent 文件来源边界仍需响应，S6 当前不得通过，I-004 保持 open，目标保持 `active · 5/6 · 83%`。

## Required findings

### F-001 · required / MEDIUM · 当前审计投影与 E-004 不一致

独立复审时 `03-audit.md` 的 I-003 虽已改为 Web `124/1516` 并引用 E-004，但同一单元格仍称真实浏览器 E2E“待补验证”；E-004 与复审命令均证明 Chromium spec 已 `1 passed`。I-004 仍只描述 A-002 的 3 个 required，未反映 A-004/A-005 的当前开放状态。需只更新当前投影与新增 response 条目，不改写 A-002/A-004：I-003 明确记录 Chromium `1 passed`，I-004 反映现行开放 finding，并保持 S6/I-004 open 直到复审。

### F-002 · required / MEDIUM · independent 文件来源边界不可直接证明

GOAL-019 目录当前为 untracked，缺少可比较的 Git 基线，无法从工作树证明 A-002/A-004 与原始 Reviewer 输出逐字一致；且 A-004 中曾混入 Supervisor 的响应事实段，使 independent 与 response 边界不清。应与本会话保留的原始 Reviewer 输出比对，记录核对结果与文件哈希；后续响应单独写入 self/response 条目，不再修改 independent 意见。若无法证明，则需用户书面选择 `accepted-residual` 或 `user-overruled`。

## 已验证

- API `server.Run` 的 JWT 与 active MFA enrollment fail-closed 路径、资源清理、空/pending/active 测试均充分；A-002 F-001 genuine fixed。
- Web 已移除 realm-sensitive `instanceof` / `String(Request)`；真实 Chromium iframe realm spec 观察到同源保留 Authorization 与 session-list refresh，跨源 Request/URL 不带凭据，auth endpoint 401 不触发 refresh；A-002 F-002 与 A-004 F-001 genuine fixed。
- `go test ./... -count=1`、Web 124/1516、typecheck、指定 Chromium 1 passed、diff check 均通过。
- E-002 阶段算术应为 S1/S2=`2/6`；goal-tree 主树/状态表和 `00-meta` 的 active/5/6/83% 一致。

## 未验证

未运行真实生产 PostgreSQL、外部消费者或生产部署冒烟；未执行完整 Playwright 套件。该 opinion 只读，未修改代码或状态。
