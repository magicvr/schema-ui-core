---
id: D-001-workspace-root-scaffold
goal_id: GOAL-001-module-contribution-readiness
status: accepted
created: 2026-08-06
updated: 2026-08-06
version: 0.1.0
parent: null
---

# D-001 · 开区与 Root 设立；挂 VP-004

## 决定了什么

1. 在用户已确认 slug 的前提下，建立显式工作区 `workspace-004-module-contribution-readiness`（`vision_role: delivery`）。
2. 设立 Root `GOAL-001-module-contribution-readiness`（`parent: null`），`plan_refs` / `primary_plan` = `VP-004-module-contribution-readiness`。
3. Root 纲领采用 VP-004 建议的 **S1–S4**（盘点冻结 → playbook 落盘 → 发现路径与抽检 → 审计与 VP 关门提案）；检查点等权，建区不勾选。
4. 主交付边界继承 VP-004：产品模块贡献方法论/playbook 落 `docs/architecture/`；不默认脚手架；不修订 principles/治理 MUST；AI 发现路径默认 overview + QUICKSTART 充分。

## 为什么

- VP-004 已由 `/vision` 激活，lead 绑定本工作区；物理 scaffold 交 `/govern`（本决策）。
- Charter `@0.2.0` 未改；结构选型为独立 delivery 区，禁止吸收进 closed workspace-003。
- VRev-010 开放 required = 0；F-V016/F-V017 已 editorial fixed。

## 未选方案

| 方案 | 未选原因 |
|------|----------|
| 在 workspace-003 继续立项 | 违反 VP-004 绑定与 closed VP 不重开约定 |
| 开区即写 playbook 正文 | 缺 S1 权威路径冻结（I-001） |
| 默认纳入脚手架/AGENTS 改写 | 与 VP-004 退出分母与 F-V017 路径 a 冲突 |

## 影响

- 实现推进焦点：`[workspace-004-module-contribution-readiness] GOAL-001-module-contribution-readiness`
- 下一步：S1 盘点与 I-001 权威路径冻结（可先决策/子目标，不机械拆分）
