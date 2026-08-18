---
id: GOAL-005-r4-async-job-contract
doc: audit
status: active
parent: GOAL-001-shared-cross-module-contracts
created: 2026-08-18
updated: 2026-08-18
version: 0.1.1
---

# 审计记录 · GOAL-005

## 审计索引

| A-ID | 日期 | source | scope | verdict | 开放 required | 文件 |
|------|------|--------|-------|---------|---------------|------|
| A-001 | 2026-08-18 | self | R4 S0：D-001/E-001/I-001～I-004；migration owner、状态机、wallet 消费边界 | pass | 0 | [A-001-r4-s0-design-self.md](03-audit/A-001-r4-s0-design-self.md) |
| A-002 | 2026-08-18 | independent | R4 S0 design-plan+execution-facts：migration/Profile、状态机、lease、wallet、actor、过期 | conditional | 6（F-001～F-006） | [A-002-r4-s0-design-independent.md](03-audit/A-002-r4-s0-design-independent.md) |
| A-003 | 2026-08-18 | self | response：A-002 F-001～F-008；D-002 精确契约 | conditional | 6（候选 fixed，待复核） | [A-003-r4-a002-response-self.md](03-audit/A-003-r4-a002-response-self.md) |
| A-004 | 2026-08-18 | independent | finding-closure：A-002 F-001～F-006；D-002；I-002/I-003；S0→S1 | conditional | 1（F-009；F-001～F-006 可 fixed） | [A-004-r4-a002-closure-independent.md](03-audit/A-004-r4-a002-closure-independent.md) |
| A-005 | 2026-08-18 | self | response：A-004 F-009；recover-cancel + CompleteWithCommit | conditional | 1（F-009 候选 fixed） | [A-005-r4-a004-response-self.md](03-audit/A-005-r4-a004-response-self.md) |

## 当前状态

S0 self A-001 = pass；A-004 已确认 A-002 F-001～F-008 fixed。A-005/D-002 v0.2.0 为 F-009 增加 recover-cancel + CompleteWithCommit；复核前 F-009 仍开放，I-003 保持 `collecting`，不进入 S1。
