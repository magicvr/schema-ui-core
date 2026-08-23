---
id: GOAL-035-w24-e2e-dual-dialect-matrix
title: W24 · 浏览器 e2e 双数据库方言矩阵（收尾层双方言各测一次）
status: done
created: 2026-08-23
updated: 2026-08-23
parent: GOAL-001-design-implementation-conformance
version: 0.2.1
progress: 4/4
---

# GOAL-035 · W24 · 浏览器 E2E 双数据库方言矩阵

## 概述

承接 GOAL-034 关门后**用户书面复审**（2026-08-23）：「强制 `DB_DIALECT=sqlite` 属绕过测试；收尾层 e2e 应把两个方言都测一次（目前只测了一个方言）」。本波把 N-001 修复从「对抗 `.env` 改道」升级为**方言契约 + 双方言矩阵**：

1. 挂具显式声明方言（默认 sqlite；postgres 显式 opt-in），`.env` 因 process-env 优先在结构上无法再静默改道；
2. postgres 模式经 `apps/api/cmd/e2e-pgset` 自动 create → run → drop 专用 scratch 库（pgtest / CI api-postgres 既有模式）；
3. `globalSetup` 启动后校验 store 契约，违反即 fail-fast 诊断（而非 3 个 spec 后神秘 401）；
4. CI `browser-e2e` 扩为 `profile × dialect` 矩阵。

**实验事实（先于立项）**：专用 scratch Postgres（全新、带种子）跑全量浏览器 e2e **9/9 绿**（2026-08-23，`e2e-w23-pg-experiment.log`）——产品 pg 方言无缺陷，缺口全在挂具层。「是否能在 pg 方言全绿」已由实验关闭（I-001）。

## 成功标准

1. C1 冻结：D-001 设计（方言契约 / provisioning 复用 pgtest-CI 模式 / fail-fast 诊断语言）；
2. C2 实施：`playwright.config.ts` + `e2e/global-setup.ts` + `global-teardown.ts` + `cmd/e2e-pgset` + `test:e2e:postgres` 脚本 + README + CI 矩阵；
3. C3 回归：sqlite 全量绿 **且** postgres 全量绿（自动 provisioning 建/删）+ vitest + tsc/build + go test；
4. C4 关门：A-001 self 审计 + goal-tree/workspace.md 终态同步。

## 路线图（P-001 · 分母 = 4）

```text
S1 冻结   → C1 D-001 设计
S2 实施   → C2 挂具/工具/CI 落地
S3 回归   → C3 双方言全量绿
S4 关门   → C4 审计 + 台账同步
```

## 信息需求登记（P-005）

| 编号 | 问题 | 级别 | 影响门禁 | 状态 | 证据/结论 |
|------|------|------|----------|------|-----------|
| I-001 | 浏览器 e2e 在专用 Postgres（全新+种子）上能否全绿？还是存在真方言缺陷？ | required | C3 | **closed** | 2026-08-23 实验 9/9 绿（专用 scratch 库 `schema_ui_e2e_w23_242e2651`），产品方言无缺陷，缺口在挂具；证据摘要已固化于 `attachments/I-001-evidence.md`（A-002 F-002） |

## 边界与审计声明

- 不改产品代码语义（仅挂具/CI/工具/文档面；`cmd/e2e-pgset` 为 dev 工具）。
- 关门审计模式 `self`。