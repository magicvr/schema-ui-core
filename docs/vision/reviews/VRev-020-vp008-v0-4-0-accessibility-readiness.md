---
doc_type: vision-review
id: VRev-020
status: active
source: independent
created: 2026-08-10
updated: 2026-08-10
version: 0.2.0
parent: null
---

# VRev-020 · VP-008 v0.4.0 独立意图复审 · 清晰度 / 准入边界 / 未考虑项（2026-08-10）

| 字段 | 值 |
|------|-----|
| source | independent |
| auditor | GitHub Copilot · `/vision-audit` |
| scope | `VP-008-admin-module-readiness-and-foundation-convergence`（`planned` v0.4.0）；用户关注：意图是否足够清晰、是否存在问题、还有哪些尚未考虑到的问题 |
| audit_type | vision-plan |
| prior_review | [VRev-017](VRev-017-vp008-intent-clarity-readiness-gates.md)、[VRev-018](VRev-018-vp008-v0-2-0-post-closure-intent-reaudit.md)、[VRev-019](VRev-019-vp008-v0-3-0-evidence-validity-review.md)（其 required findings 均已记录为 `fixed`） |
| verdict | conditional |
| 建议 class | editorial |

## 范围与结论

只读核对 [P-005 / P-006](../../architecture/principles.md)、[愿景对齐契约](../alignment.md)、[Charter](../charter.md) `schema-ui-core-admin-foundation@0.2.0`、[组合编排](../roadmap.md)、[VP-008 v0.4.0](../plans/VP-008-admin-module-readiness-and-foundation-convergence.md)、既有 VP-008 VRev 响应，以及当前代码与 fork/Compose 入口。

**结论**：VP-008 的方向和边界已足够清晰，可作为正式业务模块开发前的独立准入波次：它继承单主线、模块贡献、UI、完整协议和 locale/settings 基础；明确不重开历史 VP；将 `go` 作为后续业务 VP 的唯一解锁条件；并冻结了阶段、严重度量尺、证据基线与变更后重验证规则。其 `planned`、0 workspace 状态与组合编排一致，未见 Charter 对齐链或状态投影冲突。

但「全基架准入」仍未定义可访问性（a11y）在 S0 分母中的最小、可判定证据边界。现有代码包含语义角色、ARIA 属性和焦点样式，历史交付也有局部 a11y 契约；但这些是实施迹象，不能替代 VP-008 对冻结分母、失败处置和 S5 `go` 门闩的明确规定。由于业务模块会复用 Renderer、Shell、表单、模态和移动导航等跨模块 UI 宿主，缺少此判据可能让将来每个模块都可通过既有功能/协议测试，却共同继承键盘操作或辅助技术可用性缺陷。

本审视未运行构建、E2E 或辅助技术验证；以上运行时结论均为**证据不足**，不得解读为当前产品存在或不存在具体可访问性缺陷。

## Findings

### F-V045 · S0 准入分母缺少跨模块 UI 可访问性判据

| 字段 | 值 |
|------|-----|
| level | required |
| status | open |
| severity | major |
| scope | VP-008 S0 证据分母、S3 共性能力判断、S4 回归与 S5 `go` 证据矩阵 |
| evidence | VP-008 已要求冻结代码/环境/模块/协议/流程/用例和证据基线，但未枚举键盘可达性、焦点管理、语义/名称/状态、动态反馈或语言切换后的辅助技术行为；当前前端存在 `aria-*`、`role`、`focus-visible` 等局部实现与测试，尚非 VP-008 的可复跑准入分母。 |
| impact gate | 未关闭前，不得以「全基架」或 S5 `go` 声称跨模块 UI 宿主已完成可判定准入；该 finding 不要求在本 VP 实施领域模块，也不把完整 WCAG、性能 SLO 或全量威胁建模静默扩入 scope。 |
| close requirement | `/vision` 需在不改变 Charter 的前提下，明确本 VP a11y 范围及其处理方式：至少固定覆盖的共享 UI 宿主与可复跑断言/人工核对下限、N/A/延期的理由和重新纳入触发、失败时按既有严重度量尺的分类与 S5 矩阵列。若用户决定此风险不属于 VP-008，须按 P-004 记录范围、理由和业务 VP 前的复审触发，作为 `accepted-residual` 或 `user-overruled`；否则修订 VP 后重新独立审视。 |

## 其他已识别但非阻断问题

- **运行证据尚未形成**：VP-008 仍未激活、未绑定 lead workspace；当前无法证明冻结基线、模块接入演练、升级恢复、双 Profile 和 `go` 审计是否可复跑。这是计划状态的正常未知，不单独升级 finding；激活后由 `/govern` 按 S0～S5 收集证据。
- **安全与运维边界已被有意识排除**：VP-008 把完整威胁建模、性能 SLO 和完整运维控制列为 Non-goals；当前 Compose 对生产 secrets 采用 fail-closed。若后续 S0 发现跨模块认证、授权、数据隔离、迁移完整性或证据可复现性缺陷，既有 `blocker` 量尺已经要求进入 required，不可用本 Review 的范围说明降级。
- **业务领域风险保持隔离**：订单、钱包、类目、通知的领域模型和专有风险不应被预先塞入 VP-008；只有有证据表明其影响共同基架门闩时，才按既有量尺升级或由用户扩 scope。

## 声明

本独立意见不直接修改 Charter、VP、Goal status、progress 或既有 finding。`F-V045` 的响应、闭合或风险接受必须由 `/vision` 追加到本报告；实现工作仅在 VP 激活并建立工作区后交 `/govern`。原 verdict 与 finding 原文不得改写。

## `/vision` 响应（2026-08-10）

### 决策

- 用户确认按建议采纳本报告的 `conditional` verdict 与 `editorial` 建议 class；保留本报告原 verdict、原 finding 与原始结论，不改写历史审计文本。
- F-V045 采用 `fixed` 路径：VP-008 v0.5.0 已增加“可访问性准入边界（S0 冻结）”，限定跨模块共享 Renderer/Shell、导航/移动导航、schema-driven 表单/列表/详情/动作、模态/动态反馈与语言切换宿主；冻结键盘、焦点、语义名称/角色/状态、错误/状态反馈的可复跑断言或人工核对下限、N/A/延期触发、严重度映射与 S5 矩阵回指。
- VP-008 继续保持 `planned`、0 workspace；本响应不宣称方向已稳、不激活、不创建工作区。F-V045 已闭合，但实际可访问性断言、人工核对和失败证据仍须由后续 `/govern` 在 S0/S3/S4/S5 产生；本响应不把完整 WCAG、性能 SLO、威胁建模或领域实现引入默认 scope。

### Finding 响应台账

| finding | 原 level | 响应状态 | 响应摘要 | 证据 |
|---------|----------|----------|----------|------|
| F-V045 | required | **fixed** | 增加共享 UI 宿主可访问性范围、键盘/焦点/语义/动态反馈/语言切换最低断言或人工核对、N/A/延期触发、失败严重度和 S5 证据矩阵回指 | [VP-008 v0.5.0](../plans/VP-008-admin-module-readiness-and-foundation-convergence.md)「可访问性准入边界（S0 冻结）」、`I-READINESS-008` |

### 当前门禁

本响应将本报告的 open required 投影降为 **0**；原始 verdict `conditional` 继续保留。VP-008 仍是 `planned`、0 workspace；后续 `/govern` 进入实现前，必须按 v0.5.0 的可访问性分母和 S0/S3/S4/S5 证据规则冻结并验证共享宿主下限。任何新发现只能按已冻结量尺分类，不能借 S1 扫描重定义退出范围。
