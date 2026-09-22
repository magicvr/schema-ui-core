---
id: A-007-w18-final-s6-review
doc: audit-entry
goal: GOAL-019-w18-api-web-security-hardening
date: 2026-09-22
source: independent
auditor: Codex REVIEWER / clean-context / read-only
type: cross-audit-independent
scope: W18 final S6 gate, A-005 finding closure, API/Web regression evidence and governance/source boundaries
verdict: pass
open_required: 0
parent: GOAL-001-production-hardening
version: 0.1.0
---

# A-007 · W18 final clean-context independent review

## 结论

`verdict: pass`，开放 required = 0。API JWT 与 active MFA fail-closed、Web 跨 realm 鉴权边界、回归分母及治理来源边界均与当前证据自洽；A-005 F-001/F-002 均可标记 `fixed`，未发现新的 required finding。

本意见只通过 S6 独立审计门禁，不直接修改 status、progress、I-004、路线图或 goal-tree。

## 既有 required finding 闭合

- A-005 F-001 fixed：I-003 已明确记录 Web `124/1516`、typecheck、真实 Chromium iframe E2E `1 passed`；I-004、S6 和 `active · 5/6 · 83%` 在审计返回时未被提前修改。
- A-005 F-002 fixed：A-002/A-004/A-005 均保留 `source: independent` 与原 verdict/open_required；A-004 不再混入 Supervisor response；E-005/A-006 单独承载响应，且 A-002/A-004 SHA256 与登记值一致。

## Verified

### API

- `server.Run` 在 store/listener/authentication 前执行完整 `Config.validate()`；非 development 的缺失/弱 JWT secret fail closed。
- 未装配 MFA verifier 时，active enrollment 命中或查询失败均关闭 store 并拒绝启动；empty/pending 不误判；相关 repository、资源清理和不得继续启动测试通过。

### Web

真实 Chromium iframe realm 已验证：同源 session-list Request 保留 Authorization 与 X-Refresh-Token；跨源 Request/URL 不带两项凭据；同源 auth endpoint 401 不调用 refresh。实现已移除 realm-sensitive `instanceof Request/URL` 与 `String(Request)`，对不可解析对象 fail closed。

### 治理

- `00-meta.md`、goal-tree 和审计台账在响应前保持 active/5/6/83% 且 S6/I-004 open；响应后才进入完成投影。
- 独立意见与 Supervisor response 分离，未改写历史 independent verdict。

## Command results

- API `go test ./... -count=1`：exit 0。
- Web `npm run test`：124 files、1516 tests passed。
- Web `npm run typecheck`：exit 0。
- Chromium `auth-client-cross-realm.spec.ts`：1 passed。
- `git diff --check`：exit 0（仅 LF/CRLF 提示）。

## Unable to verify

未连接真实生产 PostgreSQL、外部消费者或生产部署；未运行完整 Playwright 套件。以上不形成开放 required。

## 声明

本意见为只读 `source: independent` 审计，不修改代码、status、progress、检查点、I-004 或 goal-tree。
