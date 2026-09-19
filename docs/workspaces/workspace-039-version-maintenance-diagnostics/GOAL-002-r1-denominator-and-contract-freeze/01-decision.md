---
id: GOAL-002-r1-denominator-and-contract-freeze
doc: decision
status: done
parent: GOAL-001-version-maintenance-diagnostics
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# 决策记录 · GOAL-002

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-039-001 | required | 版本身份权威与升级入口 | C1、R3 | R1 | 侦察 + P-004 | **verified** | — | D-001 §2 |
| I-039-002 | required | 四模式投影 | C2、R2 | R1 | 侦察 + P-004 | **verified** | — | D-001 §1 |
| I-039-003 | required | 诊断字段分母 | C3、R3 | R1 | 侦察 + 可见性裁决 | **verified** | — | D-001 §3 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-09-19 | R1 分母与契约冻结（C1/C2/C3；用户 P-004） | accepted | `01-decision/D-001-r1-contract-and-denominator-freeze.md` |

## P-004 裁决记录（2026-09-19）

| 决策点 | 用户裁决 | 未选 |
|--------|---------|------|
| maintenance 呈现（I-039-002） | **B** · 改 Host 生产者投影使 Shell 仍加载 | A HostFailureScreen 终态；C 混合 |
| 升级入口（I-039-001） | **A** · QUICKSTART 升级说明链接 | B 仅版本；C 站内 changelog |
| Shell 可见性 | **A** · 横幅所有已登录；版本仅 monitoring.read | B 版本也全员；C 仅 monitoring.read |
