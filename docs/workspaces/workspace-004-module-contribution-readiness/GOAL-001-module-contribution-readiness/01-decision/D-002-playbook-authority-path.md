---
id: D-002-playbook-authority-path
goal_id: GOAL-001-module-contribution-readiness
status: accepted
created: 2026-08-06
updated: 2026-08-06
version: 0.1.0
parent: null
---

# D-002 · 冻结 playbook 权威路径（I-001）

## 决定了什么

1. **权威操作入口**定为新建：  
   `docs/architecture/module-contribution-playbook.md`  
   （产品模块贡献方法论 / 操作 playbook；must + DO NOT + Core-vs-module 同一入口）
2. **架构终态边界**仍为 `docs/architecture/module-architecture.md`；在其 **§9** 链出 playbook，不把 playbook 正文整段并入决策文。
3. **发现路径**：`docs/architecture/overview.md` 与根 `QUICKSTART.md` §5 必须链到 playbook（S3）。
4. **不**修订 `principles.md` / workspace-protocol；**不**默认交付脚手架；**不**默认改 AGENTS/Skills。

## 为什么

- S1 盘点（[s1-gap-inventory.md](../attachments/s1-gap-inventory.md)）显示架构终态与代码范例已在，缺的是可执行操作清单与发现路径。
- 单一权威入口便于验证与 AI 发现（VP-004 F-V016/F-V017 路径）。
- 新建 authoring 文可避免 VP-003 终态决策文与「如何接模块」操作文混读。

## 未选方案

| 方案 | 未选原因 |
|------|----------|
| 仅扩展 module-architecture 长文而不拆 playbook | 终态边界与操作清单混读；S2 验证难定位 |
| 只放在 QUICKSTART | 权威性不足；与 architecture 索引脱节 |
| 写入 principles | 越界治理 MUST（VP-004 Non-goals） |

## I-001

| 字段 | 值 |
|------|-----|
| 状态 | **verified** |
| 证据 | 本决策 + `docs/architecture/module-contribution-playbook.md` 落盘 + module-architecture §9 链出 |
