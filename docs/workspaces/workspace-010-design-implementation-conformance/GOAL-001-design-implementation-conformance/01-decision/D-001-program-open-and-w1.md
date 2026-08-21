---
id: D-001-program-open-and-w1
doc: decision-entry
goal: GOAL-001-design-implementation-conformance
status: accepted
created: 2026-08-11
updated: 2026-08-11
version: 1.0.0
---

# D-001 · 开区：长期符合性程序与首波 W1

## 决定

1. 新建并激活 [VP-010-design-implementation-conformance](../../../../vision/plans/VP-010-design-implementation-conformance.md)，语义对齐 VP-009：**长期开放程序 + 周期回顾 + 波次子目标**，焦点为**设计意图—实现符合性**（非安全漏洞审计）。  
2. 唯一 lead 工作区：`workspace-010-design-implementation-conformance`；Root 为长期容器，不因单波 `done`。  
3. 首波 W1 = [GOAL-002-w1-examples-optional-module](../../GOAL-002-w1-examples-optional-module/00-meta.md)：纠正范例/演示产品面被抬成伪 core、生产 Profile 无法配置注销的符合性缺口。

## 为什么

- 用户确认：类 VP-009 的持续程序；009 管安全，新程序管「意图 vs 实现」。  
- post-go 审视表明：`core.schema-render` 演示页、Examples 导航与依赖图违背「生产不需要演示面时应可注销」及 playbook「标准模块按 Profile 装配」的产品包装意图。  
- 逻辑上属业务模块前的基架卫生，但不宜改写已 closed 的 VP-008 档案；用新长期 VP 承接，并在 W1 触及模块矩阵时按规则处理 `go` 消费有效性。

## 未选方案

| 方案 | 未选原因 |
|------|----------|
| 直接 reopen VP-008 | 可合法，但会重写已关门准入叙事；用户选择新长期程序 |
| 落入 VP-009 波次 | 意图是安全，不是符合性；避免台账混杂 |
| 落入业务 VP | 基架债不得倒逼领域交付 |
| 仅改文档不改代码 | 不闭合 as-built gap |

## 对 go 的影响（预告）

W1 预期变更 Profile 默认集 / 模块依赖 / Manifest baseline → 按 VP-008 规则，业务 VP 激活前须 freshness；若矩阵语义变化，旧 `go` 消费应暂挂直至 W1 回归证据落盘并由 `/govern` 留痕。
