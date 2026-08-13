---
id: GOAL-004-w3-schema-host-protocol-conformance
doc: execution
status: active
parent: GOAL-001-design-implementation-conformance
created: 2026-08-12
updated: 2026-08-13
version: 0.1.1
---

# 执行记录 · GOAL-004

## 执行索引

| E-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| E-001 | 2026-08-12 | 立项与 Host/App 协议候选基线 | recorded | `02-execution/E-001-goal-and-host-gap-baseline.md` |
| E-002 | 2026-08-13 | 上游 H0 处置同步与 cross 审计闭环 | recorded | `02-execution/E-002-h0-disposition-sync-and-cross-audit.md` |
| E-003 | 2026-08-13 | 上游 H0 闭合与进入 H1 accept 设计阶段 | recorded | `02-execution/E-003-upstream-h0-closed-h1-entered.md` |
| E-004 | 2026-08-13 | S4 生产实现 — 上游 2.8 候选机器契约 pin、Host 模块与浏览器级证据 | recorded | `02-execution/E-004-s4-host-interop-implementation.md` |
| E-005 | 2026-08-13 | v2.8.0 正式身份纠偏 — 重 pin 521cff8 / content 4fae4605 并重生成 claim | recorded | `02-execution/E-005-formal-v2-8-0-identity-repin.md` |
| E-006 | 2026-08-13 | S2 方案冻结 + I-005 台账收尾 + cross 审视闭环（A-005/A-006） | recorded | `02-execution/E-006-s2-plan-freeze-and-i005-ledger.md` |
| E-007 | 2026-08-13 | S4 残余整改完成 + S5 符合性验证证据（claim 重生成 buildId git:5e4c384） | recorded | `02-execution/E-007-s4-remediation-and-s5-verification.md` |

## 事实边界

> 上游 ADR-0034～0037 accepted 且 H2 机器契约落地（`453008d`），本仓按维护者裁定以固定
> commit 的候选契约满足 I-003 的 S4 开始条件并完成 S4 生产实现（E-004）：三个 host capability
> 的生产代码路径（`src/host/*` + main.tsx/App shell 集成 + Go API bootstrap 端点）消费上游
> 96 fixtures 零排除，浏览器级 focus/live-region/恢复测试全过，候选 claim 绑定 artifact digest /
> fixture digest / build ID。页面协议 2.7 的已登记 residual 与候选绑定性质见 E-004；上游
> 2.8.0 正式发布（H4）后重 pin 并重生成 claim。**身份纠偏（E-005，2026-08-13）**：上游
> v2.8.1 审计 0080 V379 裁定正式 v2.8.0 = tag @ `521cff8` / content `4fae4605…`；本仓
> 已按正式身份重 pin（commit `fd641c6`）并重生成 claim（buildId `git:fd641c6…`），vendored
> 工件与正式 tag 字节级核验一致，无生产代码改动。S2 出口门禁（I-001/I-002/I-005/I-006）
> 与 cross 方案审视仍按计划进行。
