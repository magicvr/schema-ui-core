---
id: D-001-r3-stage-scope
doc: decision
goal: GOAL-004-r3-bounded-pilot
date: 2026-08-05
status: accepted
---

# D-001 · 建立 R3 有界试点

## 决定

建立 `GOAL-004-r3-bounded-pilot`，以四个等权检查点承接 Root R3：C1 信息
门禁、C2 A+B 实施、C3 V-1～V-4 验证、C4 D 门审计。

## 边界

试点只覆盖 operationlog、Activity、Settings 及其模块注册、Schema、Manifest、
导航和 Host/Shell 边界。users/roles 全量迁移、R4 批量迁移、VP 退出和 R6
终态清理不提前纳入本子目标的完成条件。

实施边界进一步固定为：`core.operationlog` 始终启用，Activity 和 Settings
通过 Profile/模块组合路径按需启停；当前集中式 `handler.Register` 尚不能被
视为已满足模块契约；本目标不启动 R4 全量迁移。

## 依据

VP-003 R3 明确要求五项 A 交付、四类病灶切除、V-1～V-4 和 D 门全部通过；
试点通过只允许进入 R4，不是 VP 关闭依据。
