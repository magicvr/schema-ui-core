---
id: GOAL-004-r3-bounded-pilot
doc: decision
status: active
parent: GOAL-001-modular-admin-architecture
created: 2026-08-05
updated: 2026-08-05
version: 0.2.0
---

# 决策记录 · GOAL-004

## 信息与门禁规则

Root I-006 仍由 Root `00-meta.md` 维护；本子目标只承接三项收集证据。C1
的源码盘点、构建/代理边界和回滚策略不能被解释为已经完成 C2/C3 实施。

## 决策索引

| 编号 | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-05 | 建立 R3 有界试点与四项检查点 | accepted | [01-decision/D-001-r3-stage-scope.md](01-decision/D-001-r3-stage-scope.md) |
| D-002 | 2026-08-05 | R3 I-006 信息收集计划 | accepted | [01-decision/D-002-r3-i006-collection-plan.md](01-decision/D-002-r3-i006-collection-plan.md) |
| D-003 | 2026-08-05 | I-006 静态入口、兼容和回滚边界（C1 草案） | accepted-for-audit | [01-decision/D-003-r3-i006-boundary.md](01-decision/D-003-r3-i006-boundary.md) |

## 当前约束

- R3 只验证 operationlog、Activity、Settings 的模块化切口；不得借试点代码
  推导 VP-003 已完成。
- operationlog 必须 always-on；Activity 和 Settings 才能按 Profile 启停。
- 生产 Manifest 只能走 API 聚合端点；Web 静态文件即使保留为开发/测试输入，
  也不得成为生产静默兜底。
- C1 的兼容窗口、告警、移除触发和回滚策略在独立审计及演练证据前不得标记
  为 `verified`。
