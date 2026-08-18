---
id: A-003-r6-a002-response
goal: GOAL-007-r6-api-token-service-credential
source: self
verdict: pass
status: recorded
created: 2026-08-19
updated: 2026-08-19
parent: GOAL-007-r6-api-token-service-credential
version: 0.1.0
responds_to: A-002
---

# A-003 · R6 A-002 response

## 审计头

| 项 | 值 |
|----|----|
| source | self |
| scope | response：A-002 F-001～F-007；D-003 修订契约 |
| verdict | pass |
| required findings | 3（proposed fixed；待 independent closure） |

## 关闭证据

| finding | 响应状态 | 证据 |
|---------|----------|------|
| F-001 required | proposed fixed | D-003 §1：0044 credential + 0045 operationlog 三事件、correlation-safe rebuild 与测试 |
| F-002 required | proposed fixed | D-003 §2：`created_by` 无 FK，保留历史 id 且用户删除行为不变 |
| F-003 required | proposed fixed | D-003 §3/§8：NOCASE UNIQUE、trim、`ErrCredentialNameTaken` → `INVALID_CREATE_FIELD` |
| F-004 recommended | proposed fixed | D-003 §4：service prefix 先于 JWT/devSession，失败禁止 fallback |
| F-005 recommended | proposed fixed | D-003 §5：create/revoke audit 同 transaction；use/last-used best-effort |
| F-006 recommended | proposed fixed | D-003 §6：六类 user-only consumer 与 machine `self` scope 明示 |
| F-007 recommended | proposed fixed | D-003 §7：32-char id、含 revoked list/detail、分页与排序 |

## 仍开放项

在 independent finding-closure 前，A-002 F-001～F-003 的权威开放投影仍为 3；D-003 保持 `proposed`，I-002～I-004 保持 `collecting`，不放行 S1/S2。

## 结论

响应不存在 residual/overrule，也不与 A-002 冲突；全部意见均选择可核对的 `fixed` 路径，下一步由同 provider 复审关闭证据。
