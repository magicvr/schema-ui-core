---
title: D-001 · W15 范围收录与暂不推进
status: active
created: 2026-08-30
updated: 2026-08-30
parent: null
version: 0.1.0
---

# D-001 · W15 范围收录与暂不推进

日期：2026-08-30 · scope：`[workspace-009-production-hardening] GOAL-016-w15-api-web-audit-remediation` · 决策性质：立项与治理上下文

## §1 用户指令

用户要求：在工作区 9 开设一个新的子目标承载本次 `apps/api` 与 `apps/web` 审计问题的治理上下文，准备修正这些问题，但暂时不要推进。

## §2 决定

1. 在同一 Root 下开设 W15 子目标 `GOAL-016-w15-api-web-audit-remediation`，不新建工作区、不新建 VP。
2. 将本次独立审计的 F-001～F-007 纳入该目标的修正分母；F-001～F-006 作为 required，F-007 作为本波评估的 recommended；既有 refresh-token `localStorage` residual 不重开。
3. 当前目标保持 `status: draft`，仅完成目标五件套、审计意见台账、范围信息和高层路线图；不执行代码修改、方案冻结、回归验收、`go` 暂挂/恢复或关门。
4. 后续实施前按本目标 I-001～I-005 补齐信息；security/production 影响范围的实施默认走 `cross`（self + 项目约定 independent），但本条不视为已启动审计。

## §3 范围与不变项

- API：`schema-ui serve` 默认配置/生产配置校验、bootstrap 初始密码、MFA TOTP step-up。
- Web：邀请 token URL 清理与 Vitest fixture 根路径修复。
- 主机侧：LocalStore Unix 权限作为有界 hardening 评估。
- 不修改 Charter、VP-009、Root 状态；不将本波立项写成任何 finding 已 fixed。
- 旧 `GOAL-*` 目标保持历史状态；本目标编号按工作区最大编号 +1，未复用既有编号。

## §4 未选路径

| 路径 | 不采用原因 |
|------|------------|
| 直接修改代码而不建波次目标 | 无法承载 required findings、信息门禁和后续独立复核，不符合 VP-009 波次模型 |
| 重开或新建 VP / Charter | 本次是同一持续安全程序内的有界修正，不改变愿景目的、边界或意图 |
| 先把条件性发现标为 residual/overruled | 用户已要求准备修正；任何 residual/overruled 需后续 P-004 书面裁决，当前不得代裁 |

## §5 当前状态

本决策只记录立项和暂停实施的边界。S2～S6 尚未开始，代码与运行时状态保持原样。
