---
id: GOAL-041-w29-api-web-protocol-conformance
title: W29 · schema-ui-docs v2.9.0 API/Web 页面控件符合性审视与上游协议闭环
status: active
parent: GOAL-001-design-implementation-conformance
created: 2026-09-06
updated: 2026-09-06
version: 0.2.0
progress: 1/6
---

# GOAL-041 · W29 · schema-ui-docs v2.9.0 API/Web 页面控件符合性审视与上游协议闭环

## 概述

以现行 Charter 固定的 `schema-ui-docs@v2.9.0`（pinned commit `81aa1d8`）为起始协议身份，对 `apps/api` 与 `apps/web` 的实际页面、控件、Manifest、Schema、registry、request/response/action 映射和运行时失败路径建立完整分母并逐项核验。

本目标不预设“手写 UI = 协议绕过”，而是把候选区分为：

1. **implementation-gap**：上游已有明确协议语义，本仓实现未遵循；应修正本仓实现。
2. **upstream-protocol-gap**：本仓存在合理需求，但上游没有足够的数据形状、状态/错误、安全边界、能力声明或 fixtures；先形成上游协议增补报告与提案。
3. **custom-extension-candidate**：上游明确允许或明确不负责，且本仓可以在允许范围内以 namespaced custom 控件/Host 边界承载；须另有可验证的兼容与 fail-closed 决策。
4. **explicitly-out / excluded**：明确属于 Host、身份/部署安全、业务域或本波边界外，不伪装成协议能力。

对确认为 `upstream-protocol-gap` 的候选，**上游协议 accepted/merged、正式版本与 commit 固定、本仓 provenance/schemas/registry/fixtures/兼容证据完成前，禁止实施或放行依赖该增补的本仓产品修改**。若存在需上游增补的条目，汇总报告固定为 `attachments/upstream-protocol-augmentation-report.md`。

## 边界

### 纳入

- `apps/api` 的页面 Schema 提供、Manifest/fragment 聚合、component/capability registry 映射、action/request/response/dataSource/permission 契约及相关失败路径。
- `apps/web` 的 manifest route → `schemaUrl` → page document → Renderer 链路，以及 form、table、recordView、actionButton、navigation、route binding、readonly、upload、reaction、custom 控件和错误/降级呈现。
- `schema-ui-docs@v2.9.0` 机器工件、本仓 provenance/schema/registry/fixtures、实际运行时 Manifest 与代表性页面流程之间的同源一致性。
- 确认后的本仓实现整改、上游协议增补等待门禁、合法 custom 边界与对应验证。

### 排除

- 不重写或重新开启已完成的 GOAL-004；仅引用其 v2.8 Host/App 边界、处置语义与历史证据。
- 不把 Login、branding、session、Shell、全局导航、错误边界等 Host 面一概判为违规；只核验其是否落在已固定 Host/App 契约或明确扩展边界内。
- 不承载纯安全漏洞、性能专项、业务域功能新增或身份平台/token 生命周期实现；命中安全主因时以限定引用移交 workspace-009。
- 不在上游协议未固定时于本仓自造稳定协议语义；不把静态 schema 通过、registry 存在或单个 happy path 当作运行时符合性证明。

## 成功标准 / 高层路线图

- [x] **S1 · v2.9 分母与候选目录**：已固定 v2.9 tag/commit 与 11/24/19/20 协议分母，建立 17 fragments / 35 页面 schema / 11 renderer nodes / 14 controls / 15 custom registrations 目录，并登记 C-001～C-014；证据见 E-002 与 `attachments/S1-*`。
- [ ] **S2 · 差异分类与方案冻结**：每项候选有协议/代码/测试证据，归入 implementation-gap、upstream-protocol-gap、custom-extension-candidate、explicitly-out 或 excluded；完成 self + independent 的方案级 cross 审视。
- [ ] **S3 · 上游协议或 custom 边界固定**：协议缺口取得 accepted/merged 的上游契约、正式版本/commit 与可消费机器工件；custom 候选取得用户书面裁决及 namespace/capability/schema/validator/failure/compatibility/fixture 边界。
- [ ] **S4 · API/Web 实现整改**：只按已固定的上游协议或合法 custom 契约修改实现；已有协议偏差全部有修复与防复发证据。
- [ ] **S5 · 运行时符合性验证**：validator、正反 fixtures、代表性页面与失败路径、API/Web 定向与全量回归可复跑；记录 VP-008 `go` 消费影响与暂挂/恢复结论。
- [ ] **S6 · 关门审计**：全部到期 required 信息项与 required findings 合法闭合；上游报告、custom 裁决、回归和 cross 关门意见完整；用户确认后才可 `done`。

`progress` 按上述六个等权检查点确定性派生；部分完成不计入。当前 S1 已完成，故为 `1/6`。

## 信息就绪

完整信息台账见 [01-decision.md](01-decision.md)。S1 已以 E-002 和 `attachments/S1-*` 回答 I-001/I-002；I-003 为 `collecting`，I-004～I-008 仍开放，I-009 为 non-blocking。S2 只允许证据补齐、逐项分类和方案 cross 审视，不越过上游协议/custom 决策门禁实施受影响产品代码。

## 审计模式

- 模式：`cross`。
- 最低要求：S2 方案冻结前形成 self 意见 + 项目默认 independent provider 的 `/audit` 意见；S6 关门前对运行时符合性、上游门禁、custom 边界与失败路径复审。
- provider：按项目级 `docs/architecture/independent-audit-execution.md` 使用本地 `grok build`（grok 4.6，思考强度 high）。若 provider 不可用或没有可核对输出，不得冒充 independent，相关门禁保持未满足。
