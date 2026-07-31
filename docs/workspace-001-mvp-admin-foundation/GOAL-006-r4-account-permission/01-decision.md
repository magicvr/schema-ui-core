---
id: GOAL-006-r4-account-permission
doc: decision
status: active
parent: GOAL-001-mvp-admin-foundation
created: 2026-07-31
updated: 2026-07-31
version: 0.1.0
---

# 决策记录 · GOAL-006

## 信息需求与阶段门禁

P-005 信息台账维护在 [00-meta.md](00-meta.md)。本目标的 `I-006-001` 在 R4 **方案冻结前**须验证；父目标 `I-PROTO-002` 在 R4 **实施前**须合法闭合。本目标不修改 Root 门禁状态。

## D-001 · 立项 R4 核心账号与权限子目标

**日期**：2026-07-31
**状态**：accepted

**决定**：

在 `GOAL-001-mvp-admin-foundation` 下创建 `GOAL-006-r4-account-permission`，将 R4 范围限定为账号权限最小 API 设计、`D-PERM` 映射冻结与前后端鉴权链路实现；Root 保持 `active`，纲领进度仍为 `3/6`。

**为什么**：

- Root 路线图把 R4 定义为「核心账号与权限」，且 `I-PROTO-002` 要求「设计最小 API + 对照 permissions-inheritance fixtures」。
- R3（GOAL-005）已完成 Admin 外壳/导航，其默认空 navigation context 需要真实身份/权限来源衔接——正是 R4 的产品边界。
- 协议清单将账号/权限映射为 `D-PERM` 核心能力，并有配套 behavioral fixture 可作验证来源。

**未选方案**：

- 在 R3 目标中补入鉴权：会越过 R3 已关门边界并绕过 `I-PROTO-002` 实施门禁。
- 吞并 R5 Renderer 全量与范例页：会提前越过 `I-PROTO-003` 验收/关门门禁。
- 以硬编码角色表代替 `D-PERM` 契约映射：无法证明与固定协议语义一致。

**影响**：

本目标进入 `active` 规划阶段；R4 方案冻结前验证 `I-006-001`，实施前闭合 `I-PROTO-002`。当前不修改 `apps/*`，不改变父目标或协议门禁状态。

## D-002 · 采用「信息就绪 → 方案冻结 → 实施 → 验证关门」的 R4 路线

日期：2026-07-31
**状态**：accepted

**决定**：

先处理 `I-006-001` 并闭合父目标 `I-PROTO-002`，再冻结账号权限最小 API、`D-PERM` 映射与前后端集成边界；在此之前不把待确认取舍写成实现事实。实现完成后必须补结构/行为/运行时证据和阶段自审，才可讨论 `done`。

**为什么**：

- 沿袭 R3（GOAL-005 D-002/D-005）已建立的信息门禁纪律：开放 required 信息项是方案冻结与受影响实施门禁。
- 账号/权限涉及安全边界，未验证的映射或伪造的「已验证」都会在验收时产生误导。

**未选方案**：

- 先写鉴权代码再补契约映射：会把未知项伪装成已决定行为。
- 仅以构建成功或登录页可打开作为 R4 关门证据：无法覆盖 `D-PERM`/权限继承契约。

**影响**：

`I-006-001` 与 `I-PROTO-002` 是 R4 方案冻结/实施门禁；必要时须按 P-004 由用户裁决 fixed、accepted-residual 或 user-overruled，不能静默放行。

## D-003 · 采用固定上游资料与明确的验证证据边界

日期：2026-07-31
**状态**：accepted

**决定**：

R4 规划以 `protocol-inventory-v2.7.0.md` 登记的 source commit `ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b` 为资料定位锚点，优先对照 `D-PERM` 与 `permissions-inheritance` 等固定 fixture；未来验收须记录 schema/fixture 或等价可核对证据。`conformance/reference-js`、`reference-python` 和 `runner` 仅可作为参考，不能单独证明兼容。

**为什么**：

协议清单明确区分 structural contract、behavioral fixture 和 excluded reference/runner；沿用 R3（GOAL-005 D-003）已确立的证据边界，避免把「找到路径」或「调用参考实现」误写成已验证。

**影响**：

该决定固定证据方向，不代表当前已有本地 schema、fixture、运行时或 conformance 结果。
