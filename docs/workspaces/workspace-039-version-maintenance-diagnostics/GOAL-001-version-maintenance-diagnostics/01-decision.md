---
id: GOAL-001-version-maintenance-diagnostics
doc: decision
status: done
parent: null
created: 2026-09-19
updated: 2026-09-19
version: 0.2.0
---

# 决策记录 · GOAL-001

## 信息需求与阶段门禁

> 本文件是稳定索引。信息台账正文在 `00-meta.md`（Root 层）维护；长决策与独立决策记录放在 `01-decision/D-NNN-<slug>.md`。

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-039-001 | required | 版本身份权威字段与升级说明入口 | R1、R3 | R1 | 扫描 `pkg/version` 与 system-monitoring | **verified** | — | GOAL-002 D-001 §2 |
| I-039-002 | required | 四种 runtime.mode 投影 | R1、R2 | R1 | 对照 bootstrap / operational / feedback-policy | **verified** | — | GOAL-002 D-001 §1 |
| I-039-003 | required | 诊断摘要字段分母 | R1、R3 | R1 | 对照 healthz/readyz/status | **verified** | — | GOAL-002 D-001 §3 |
| I-039-004 | required | 承载面与默认集 | 激活、`go` | 激活前 | 默认候选核验 | **verified** | 改新模块须复核 `go` | D-001；VRev-102 |
| I-039-005 | required | Admin 类 freshness | 激活与开区 | 激活前 | 五域 freshness | **verified** | 下次基线变更复核 | D-001；VRev-102 |
| I-039-006 | non-blocking | runtime.mode 热切换 | 不进首波 | — | `/vision` 复核 | deferred | 责任人 `/vision` | 首波不承诺 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-09-19 | 工作区与 Root 建立：VP-039 激活落盘、`I-039-004` 裁决与 freshness 记录 | accepted | `01-decision/D-001-workspace-root-establishment.md` |
| D-002 | 2026-09-19 | 用户书面关门确认 | accepted | `01-decision/D-002-user-close-confirmation.md` |
