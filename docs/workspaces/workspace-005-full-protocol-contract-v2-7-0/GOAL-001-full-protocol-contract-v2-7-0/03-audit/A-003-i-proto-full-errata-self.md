---
id: A-003-i-proto-full-errata-self
doc: audit-entry
record_id: A-003
source: self
auditor: govern orchestrator
audit_type: finding-response
scope: I-PROTO-FULL-001 v1.0.1 勘误
status: recorded
verdict: pass
parent: null
created: 2026-08-10
updated: 2026-08-10
version: 1.0.0
---

# A-003 · I-PROTO-FULL-001 勘误自审

## 核对

1. 当前权威覆盖表为 v1.0.1，明确 320 total = 318 executed + 2 local adapter excluded。
2. 两项 exclusion id 与 `CAPABILITY_REQUIRED` / `MISSING_REQUIRED_CAPABILITY` 差异已逐项写入覆盖表与 D-003。
3. 域级 12/12 include、registry 24/24、suite 16/16 include 未改变；I-002 仍为 N/A，未产生协议范围 residual。
4. workspace-008 A-002 F-001 已通过 A-003 响应以 `fixed` 路径闭合；A-002 原 verdict 与 finding 原文未改写。

## 证据

- `attachments/I-PROTO-FULL-001-coverage-v2-7-0.md` v1.0.1
- `01-decision/D-003-i-proto-full-errata.md`
- `02-execution/E-007-i-proto-full-errata.md`
- `apps/web/src/protocol/upstream-fixtures.test.ts`
- `apps/web/src/protocol/conformance/stage3-fixtures.test.ts`（260/260 passed）
- `apps/api/internal/docscheck`（`go test ./internal/docscheck` passed）

**verdict：pass**。本自审只确认勘误内容与证据链，不推进 Root/VP 状态或用户 go/no-go。
