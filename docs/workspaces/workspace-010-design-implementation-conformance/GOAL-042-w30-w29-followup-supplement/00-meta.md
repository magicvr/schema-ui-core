---
id: GOAL-042-w30-w29-followup-supplement
title: W30 · W29 后继补充 — legacy 能力保守声明审计、claim↔host-support 一致性、10 页行为级单测
status: done
parent: GOAL-001-design-implementation-conformance
created: 2026-09-06
updated: 2026-09-06
version: 0.2.0
progress: 3/3
---

# GOAL-042 · W30 · W29 后继补充

## 概述

承接 [GOAL-041-w29-api-web-protocol-conformance](../GOAL-041-w29-api-web-protocol-conformance/00-meta.md)（2026-09-06 关门，`done` 6/6）留下的三个 non-blocking 后续项（用户书面指令 2026-09-06：「随便一起做掉，治理层面作为本目标后继补充」）：

1. **F1 · legacy 能力保守声明全量审计**：对全部 35 页的 `meta.requiredCapabilities` 与文档内实际使用做全量能力审计（不止 v2.9 能力对），清理「声明未使用」的保守声明（上游单向约束下合法，但为卫生删除或登记意图），并把 `capability-declaration.guard.test.ts` 从「仅 v2.9 能力对」扩展为「全部非豁免能力」强制。
2. **F2 · claim↔host-support 机械一致性**：把 host 支持集（页面版本 + 能力）收敛为单一 JSON 真值（`apps/web/src/host/host-support.json`），`host-support.ts` 与 `generate-claim.mjs` 同源消费，并新增一致性测试（claim 与 host-support、mandatory suites 覆盖）。
3. **F3 · 10 页行为级单测**：对 S5 分母中仅「全分母渲染 + e2e 兜底」的 10 页补行为级单测——`mail`、`mail-outbox`、`my-wallet`、`wallet`、`wallet-vouchers`、`telegram-settings`、`telegram-operator`、`digitaloffer-offers`、`digitaloffer-entitlements`、`digitaloffer-purchases`。

## 边界

- 仅处理 W29 遗留的审计/一致性/测试补强；不改变协议语义、不新增 capability、不改 Profile/模块矩阵/Manifest 装配。
- F1 删除声明属于 schema 元数据卫生（与 C-005 先例一致）；若某页声明对应「custom 组件内使用」等无法文本标记的能力，保留并登记意图。
- 沿用 W29 契约：上游单向约束（用时须声明）、envelope 能力（app.manifest/app.navigation）与 host 级能力（host.*）豁免。

## 成功标准 / 路线图

- [x] **F1 · legacy 能力全量审计**：守卫覆盖全部非豁免能力（13 markers + 豁免 + intent override）；32 个 schema 双向修正（删 22 处未使用声明、补 11 处欠声明）；守卫 35/35 绿。
- [x] **F2 · claim↔host-support 一致性**：`host-support.json` 单源落地（host-support.ts + generate-claim.mjs 同源）；一致性测试 5/5 绿（集合相等 + 能力 ID 合法 + mandatory suites 覆盖）；claim 重生成。
- [x] **F3 · 10 页行为级单测**：`behavior-pages.test.tsx` 20/20 绿（真实链路 + 页面特有 UI 断言）。
- [x] 自审 + 关门：A-001 self **pass**（0 required）；全量回归 Web vitest 1332/1332 + build 0 + Go 0 FAIL；**子目标关门经审计执行（2026-09-06）→ `status: done` / 3/3**。

`progress` 按上述三个等权检查点派生；F1/F2/F3 全部完成，故为 `3/3`；目标 `status: done`。

## 信息就绪

| ID | 级别 | 所需信息 | 影响门禁 | 状态 | 证据/结论 |
|----|------|----------|----------|------|-----------|
| I-001 | required | F1 审计后未使用声明的处置口径（删除 vs 保留+意图登记） | F1 完成 | **verified**（用户 2026-09-06 书面：一起做掉，按卫生清理） | 删除 22 处 genuinely-unused；notifications actions.page.trigger 保留并登记意图（custom 铃铛） |
| I-002 | required | claim↔host-support 单源形态 | F2 完成 | verified（D-001） | host-support.json 单源 + 一致性测试 5/5 |
| I-003 | non-blocking | 10 页行为断言的基线形状 | F3 | verified | 由各页 schema 列字段派生；20/20 绿 |

## 审计模式

- 模式：`self`（测试/卫生/审计补强，无 security/data/migration/production 门禁；沿用 W29 已固定的协议契约）。
- 关门：A-001 self 关门审计 **pass**（0 required）→ **子目标关门经审计执行（2026-09-06）**。若实施中触达跨边界语义则升级询问（本波未触达）。
