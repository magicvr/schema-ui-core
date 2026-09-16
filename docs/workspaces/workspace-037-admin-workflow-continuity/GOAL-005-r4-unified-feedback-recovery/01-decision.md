---
id: GOAL-005-r4-unified-feedback-recovery
doc: decision
status: active
parent: GOAL-001-admin-workflow-continuity
created: 2026-09-17
updated: 2026-09-17
version: 0.1.0
---

# 决策台账 · GOAL-005 R4

## 信息需求与阶段门禁

| ID | 级别 | 影响门禁 | 状态 | 证据 / 结论 |
|----|------|----------|------|-------------|
| R4-I-001 | required | C1/C3 | collecting | 需建立 resource/Host 分类与安全文案映射 |
| R4-I-002 | required | C2/C3 | collecting | 需核对读 retry 与写入不重试边界 |
| R4-I-003 | required | C2/C3 | collecting | 需核对共享反馈表面的 role/focus/dismiss/retry 去重 |
| R4-I-004 | required | C3/C4 | collecting | 需核对维护/不可用/离线/超时的 resource/Host 边界 |
| R4-I-005 | non-blocking | UX/支持体验 | deferred | 仅保留安全 code/correlation；真实支持需求出现时再由 `/vision` 复核 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|--------|------|
| D-001 | 2026-09-17 | R4 反馈分类、重试与无障碍合同 | accepted | [D-001-r4-feedback-recovery-contract.md](01-decision/D-001-r4-feedback-recovery-contract.md) |

## 当前投影

- R4 不新增用户待裁决的架构选择，直接执行 R1 D-005 已冻结的反馈/恢复边界；实现上优先复用既有 `readResourceApiError`、`FeedbackRegion`、`DataTable` retry 与 Host failure 合同。
- R4 先完成信息登记与分类合同，再改共享反馈表面；在 R4 C1～C3 门禁关闭前，不把跨页面一致性写成已完成事实。
