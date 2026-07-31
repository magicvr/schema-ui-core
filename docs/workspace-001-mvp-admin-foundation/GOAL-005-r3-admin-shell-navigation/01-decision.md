---
id: GOAL-005-r3-admin-shell-navigation
doc: decision
status: active
parent: GOAL-001-mvp-admin-foundation
created: 2026-07-31
updated: 2026-07-31
version: 0.1.1
---

# 决策记录 · GOAL-005

## 信息需求与阶段门禁

P-005 信息台账维护在 [00-meta.md](00-meta.md)。本目标的 `I-005-001` 至 `I-005-005` 均为 required/open；在对应最晚阶段前未 verified，且没有用户书面接受的有界 residual 时，不得冻结方案或进入受影响的实施门禁。

父目标的 `I-PROTO-002` 与 `I-PROTO-003` 分别属于 R4/R5 的后续门禁。本目标不修改其状态，也不以 `I-PROTO-001=verified` 替代它们。

## D-001 · 立项 R3 Admin 外壳与导航子目标

**日期**：2026-07-31
**状态**：accepted

**决定**：

在 `GOAL-001-mvp-admin-foundation` 下创建 `GOAL-005-r3-admin-shell-navigation`，将 R3 范围限定为 App manifest 装载、Admin shell、导航入口和路由语义规划/实施；Root 仍保持 `active`，纲领进度仍为 `2/6`。

**为什么**：

- Root 路线图已把 R3 定义为“Admin 外壳与导航”，R2 的 D-009/A-006 已记录冻结后进入 R3 规划。
- R1 已明确不包含 Admin 导航壳、多业务路由和协议 Renderer 全量；R3 需要有独立的范围、信息项和审计台账。
- 协议清单将 `D-APP` 的 React 主责定义为装载与导航壳，能够为本目标提供范围锚点，但不能替代本地实现决策。

**未选方案**：

- 在 R1 目标中直接补入 Admin shell 或业务路由：会越过既有 R1/R3 边界。
- 把账号权限、Renderer 全量或业务范例一起纳入本目标：会提前吞并 R4/R5，且绕过其 required 信息门禁。
- 以硬编码业务路由树代替 manifest/导航契约：无法证明与 `S-09`、`app-navigation` 语义一致。

**影响**：

本目标进入 `active` 规划阶段；当前不改 `apps/web`，不改变父目标或协议门禁状态。

## D-002 · 采用“信息就绪 → 方案冻结 → 实施 → 验证关门”的 R3 路线

日期：2026-07-31
**状态**：accepted

**决定**：

先处理 `I-005-001` 至 `I-005-005`，再冻结 manifest 最小子集、路由映射、默认/fallback/active-route 语义和 shell 产品边界；在此之前不把待确认取舍写成实现事实。实现完成后必须补结构/行为/运行时证据和阶段自审，才可讨论 `done`。

**为什么**：

当前 Web 仍是 R1 单页占位，`main.tsx` 没有 router，仓库也没有本地 manifest loader、navigation 或 shell 实现。现有协议资料只固定了上游路径与验证方向，不能据此推断本地路由行为。

**未选方案**：

- 先写一个临时 shell，再事后补 manifest 与路由语义：会把未知项伪装成已决定行为。
- 仅以构建成功或页面可打开作为 R3 关门证据：无法覆盖 manifest/navigation 契约。

**影响**：

开放 required 信息项是 R3 方案冻结和受影响实施门禁；必要时须按 P-004 由用户裁决 fixed、accepted-residual 或 user-overruled，不能静默放行。

## D-003 · 采用固定上游资料与明确的验证证据边界

日期：2026-07-31

状态：accepted

**决定**：

R3 规划以 `protocol-inventory-v2.7.0.md` 登记的 source commit `ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b` 为资料定位锚点，优先对照 `S-09`、`app-manifest` 和 `app-navigation`；未来验收须记录 schema/fixture 或等价可核对证据。`conformance/reference-js`、`reference-python` 和 `runner` 仅可作为参考，不能单独证明兼容。

**为什么**：

协议清单明确区分 structural contract、behavioral fixture 和 excluded reference/runner。将资料版本和证据类别写入计划，可以避免把“找到路径”或“调用参考实现”误写成 R3 已验证。

**影响**：

该决定固定证据方向，不代表当前已有本地 schema、fixture、运行时或 conformance 结果。

## D-004 · 响应 A-001，统一方案冻结门禁并显式登记父级依赖

**日期**：2026-07-31
**状态**：accepted

**决定**：

采用 A-001 对 F-001 的推荐修正：将 `I-005-003` 和 `I-005-004` 的最晚需要阶段统一为「方案冻结前」，使信息表与 D-002 的“先处理 `I-005-001` 至 `I-005-005`，再冻结”保持一致。采纳 F-002：`I-005-001` 显式依赖 Root `I-PROTO-004` 的 vendor vs pin 工程策略；关闭该项时必须记录所选校验方式和失败边界，但不把 Root 的 non-blocking 项静默升格为本目标独立 required。

**为什么**：

- A-001 已证明原信息表与 D-002 对方案冻结的要求存在矛盾；采用推荐修正可使编排器只有一个可核对的冻结门槛。
- `app-manifest` 的本地校验方式会影响 `I-005-001` 的证据边界；显式关联 Root `I-PROTO-004` 可避免把“已知上游路径”误写成“本地已可校验”。

**影响**：

这只闭合 A-001 的 F-001 文档一致性 finding，不验证或关闭任何 `I-005-*`。R3 仍不得方案冻结、实施或关门；用户在未来拟用独立审计推进方案冻结前，仍须按 P-004.1 决定是否进行同 scope 自审。
