---
id: D-001-workspace-s0-bindings
doc: decision-entry
goal: GOAL-001-admin-module-readiness
status: accepted
created: 2026-08-10
updated: 2026-08-10
version: 0.1.0
---

# D-001 · 开区绑定、Root 命名与 independent provider

## 用户裁决

用户于 2026-08-10 确认：

- 采用**单工作区**模式；不创建 VP-008 的 support workspace。
- lead workspace 为 `workspace-008-admin-module-readiness`，角色为 `delivery`。
- Root 标题为 `Admin 业务模块准入与基架收敛`，slug 为 `admin-module-readiness`，完整 ID 为 `GOAL-001-admin-module-readiness`，`parent: null`。
- `I-READINESS-005` 的 independent provider 由本次编排选定为 **GitHub Copilot · `/audit`**，后续采用 `cross` 模式（self + independent）。

## provider 选择依据与边界

VRev-024 与 VRev-026 的既有 independent 报告均记录 `GitHub Copilot · /vision-audit`。本次沿用同一 provider 名称以保持审计链可追溯，但把执行入口明确为 Goal 层的 `/audit`；历史 VRev 不能充当本 Root 的独立执行证据。

独立审计覆盖 compatibility、data、migration、production/release 以及跨工作区/跨层治理语义；self 与 independent 必须分别记录 scope、证据路径和 findings。provider 不可用或没有可核对输出时，相关门禁保持未满足，不由编排器冒充 independent。

## 未选方案

- **多工作区**：当前没有已确认的 support 范围，先保持单 lead；新增 support 需要重新记录绑定和证据聚合规则。
- **将 provider 记录为“当前会话”**：无法提供独立性，违反 P-004；故固定为后续指定的 GitHub Copilot `/audit` 会话。
- **在开区时直接关闭 I-READINESS-005**：provider 选择只是输入，审计执行与可核对输出仍待 S0/S5。
