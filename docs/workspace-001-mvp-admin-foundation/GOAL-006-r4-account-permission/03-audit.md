---
id: GOAL-006-r4-account-permission
doc: audit
status: active
parent: GOAL-001-mvp-admin-foundation
created: 2026-07-31
updated: 2026-07-31
version: 0.1.0
---

# 审计 · GOAL-006

> 本文件是目标的唯一正式意见台账（P-003）。self / independent 意见共用 `A-00N` 序列。

## 信息就绪核对

| 核对项 | 状态 | 备注 |
|--------|------|------|
| 影响本 scope 的 I-00N | `I-006-001` 为 required/open | 见 [00-meta.md](00-meta.md)；R4 方案冻结前须验证 |
| 到期 required 是否已 verified / residual | 无到期项 | `I-006-001` 最晚阶段为方案冻结前，本轮尚在规划立项 |
| 固定资料引用 | 沿用 Root 冻结基线 | `I-PROTO-001` v0.1.3 冻结；协议固定 commit `ca9e5fe…`（artifact `2.7.0`） |
| 当前实现证据 | 无 R4 实现事实 | 见 [02-execution.md](02-execution.md)；本次未修改 `apps/*` |

## 意见台账索引

| ID | 日期 | source | scope | verdict | 开放 required |
|----|------|--------|-------|---------|---------------|
| — | — | — | — | — | — |

## 当前开放门禁

- `I-006-001`（required/open）：R4 方案冻结前验证账号权限最小 API 与 `D-PERM` 映射。
- 父目标 `I-PROTO-002`（required/open）：R4 **实施**门禁，实施前须合法闭合。
- 父目标 `I-PROTO-003`（required/open）：R5 验收/关门门禁，本目标不处理。
- 父目标 `I-PROTO-004`（non-blocking/open）：关闭时须补 schema-conformance 等价性校验或显式记录等价范围（GOAL-005 A-007 F-002 跟进）。

## 备注

本目标尚未到复盘节点；规划立项阶段无正式审计意见。
