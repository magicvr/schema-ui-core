---
id: D-001
doc: decision-entry
goal: GOAL-001-store-dialects
status: accepted
created: 2026-08-20
updated: 2026-08-20
version: 1.0.0
---

# D-001 · 开区 scaffold 与 A1 纲领路线图

## 背景

用户确认：对 VP-013 做意图审视后激活并开工作区；slug 授权由编排选定。VRev-029 self `pass`（0 required）。RT-P03 / VR-027 已冻结双方言决策。

## 决策

1. 激活 `VP-013-store-dialects`（v0.2.0 `planned → active`）。
2. lead 工作区 slug = `workspace-013-store-dialects`；Root = `GOAL-001-store-dialects`。
3. 纲领路线图 R1～R5：端口冻结 → PG 接入 → 台账对写 → 公共面收口 → 双路径证据。
4. 配置面：缺省仍为 SQLite `db.path`；PostgreSQL 为显式 DSN 的生产/验收路径。不改 Compose 默认依赖。
5. 开区审计模式 **none**（可逆文档 scaffold）。R1 端口方案冻结起按内核/数据门禁走 **self**，迁移与生产路径实施按 **independent**（项目默认 grok build）。
6. 本回合**不**创建 R1 子目标、**不**改 `apps/api` 代码。

## 为什么

- 新纲领波次独立工作区，避免写入已 closed 的 workspace-012。
- slug 与 VP id 对齐，便于组合索引。
- I-001～I-004 满足 V-F058；配置面满足 V-F059。

## 未选方案

- 继续 `planned` 只写 VP：用户已要求激活并开区。
- 重开 workspace-012：VP-012 已关门且默认不接新区。
- 一开区就改 store 包：R1 方案未冻结。
- 默认改为必须 PostgreSQL：与 RT-P03 内嵌默认冲突。
