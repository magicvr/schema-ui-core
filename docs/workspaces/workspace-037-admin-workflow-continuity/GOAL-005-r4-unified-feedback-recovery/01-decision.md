---
id: GOAL-005-r4-unified-feedback-recovery
doc: decision
status: done
parent: GOAL-001-admin-workflow-continuity
created: 2026-09-17
updated: 2026-09-17
version: 0.3.0
---

# 决策台账 · GOAL-005 R4

## 信息需求与阶段门禁

| ID | 级别 | 影响门禁 | 状态 | 证据 / 结论 |
|----|------|----------|------|-------------|
| R4-I-001 | required | C1/C3 | verified | `feedback-policy.ts` 与 `transportFailureResult` 保留 AbortError → timeout、网络 TypeError → offline；E-004 已补回归，A-003 independent recheck 待完成 |
| R4-I-002 | required | C2/C3 | verified | 仅读取路径提供显式 retry；写入/action 不自动重试且 F-001 修复未放宽该边界 |
| R4-I-003 | required | C2/C3 | verified | 共享反馈表面的 role/focus/dismiss/retry 去重与键盘路径已有 E-002/A-001/A-002 证据 |
| R4-I-004 | required | C3/C4 | verified | 普通 resource/Host 边界与 timeout/offline policy 已核对；E-004 补齐 transport catch，A-003 independent recheck 待完成 |
| R4-I-005 | non-blocking | UX/支持体验 | deferred | 仅保留安全 code/correlation；真实支持需求出现时再由 `/vision` 复核 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|--------|------|
| D-001 | 2026-09-17 | R4 反馈分类、重试与无障碍合同 | accepted | [D-001-r4-feedback-recovery-contract.md](01-decision/D-001-r4-feedback-recovery-contract.md) |

## 当前投影

- R4 不新增用户待裁决的架构选择，直接执行 R1 D-005 已冻结的反馈/恢复边界；实现上优先复用既有 `readResourceApiError`、`FeedbackRegion`、`DataTable` retry 与 Host failure 合同。
- R4 C1～C3 的实现与回归证据已形成；A-002 的 F-001 required finding 已按 `fixed` 路径完成代码响应，等待 A-003 independent recheck 后再执行 C4 checkpoint 与关门。
