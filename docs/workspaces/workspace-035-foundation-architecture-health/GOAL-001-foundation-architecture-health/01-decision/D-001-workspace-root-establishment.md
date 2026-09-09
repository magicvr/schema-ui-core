---
doc_type: goal-decision
id: D-001-workspace-root-establishment
parent: GOAL-001-foundation-architecture-health
date: 2026-09-09
status: accepted
version: 0.1.0
---

# D-001 · 工作区 / Root 建立与开区决策

## 上下文

用户 2026-09-09 指令「OK 激活 VP-035，slug 用建议名」。激活门禁已满足：VRev-087 self `pass`（0 required）+ 架构类轻量 freshness PASS（`f2044cf3` → `5c341ec7`，不暂挂 `go`）。I-035-002 业界参照集已在计划阶段 verified。

## 决策

| 项 | 决定 |
|----|------|
| 工作区 | `workspace-035-foundation-architecture-health`（canonical `docs/workspaces/workspace-035-foundation-architecture-health/`） |
| Root | `GOAL-001-foundation-architecture-health`（`parent: null`；primary_plan = `VP-035-foundation-architecture-health`） |
| 愿景角色 | `delivery`（不改变 Charter primary workspace） |
| 纲领阶段 | R1 分母冻结 → R2 as-built 矩阵 → R3 业界对照+分类 → R4 路线图草案与关门（串行） |
| 审计模式 | 阶段关门 default **self**；R4 关门与路线图 editorial 冻结前建议 **independent**（V-F122 · `/vision-audit`） |
| 红线 | 不实现 Redis/MQ/K8s/ORM；不消耗 trigger-gated 行；不改 Profile 默认集；不重开已 closed VP；I-035-003 若对照要动 Charter 则停住 |

## 继承的激活冻结

| ID | 冻结结论 |
|----|----------|
| I-035-002 | 四类参照集：模块化单体+组合根；基础设施端口（内存默认、外部 gated）；Schema 驱动 Admin；同进程基座。每行四格。不得推翻 Charter 非目标。 |

I-035-001 / 004 / 005 留待 R1 用户冻结，不在开区时静默裁定。

## freshness 三字段

| 字段 | 值 |
|------|-----|
| consumer_vp | `VP-035-foundation-architecture-health`（vision_ref `schema-ui-core-admin-foundation@0.4.0`） |
| last_freshness_review_at | 2026-09-09（`f2044cf3` → `5c341ec7` · 架构类轻量 PASS · pin / 锁 / 迁移 / Profile 默认集零变更；区间 = VP-034 已审结目 + 发布包装） |
| next_freshness_review_trigger | 本 VP 关门前最终验证；或首个后续业务域 VP 激活（H-002 发现机制） |

## 未选方案

- 不把本意图做成 workspace-010 波次（用户已否决；010 不拥有路线图重写）。
- 不在开区时先改 `roadmap.md` 现状锚点（那是 R4 草案 + `/vision` editorial）。
- 不在开区时开始代码优化（I-035-005 未冻结）。
