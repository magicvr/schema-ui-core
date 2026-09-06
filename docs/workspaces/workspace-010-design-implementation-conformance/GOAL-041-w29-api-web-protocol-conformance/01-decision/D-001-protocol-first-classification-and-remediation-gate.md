---
id: GOAL-041-w29-api-web-protocol-conformance
doc: decision-entry
record_id: D-001
status: accepted
parent: GOAL-001-design-implementation-conformance
created: 2026-09-06
updated: 2026-09-06
version: 0.1.0
---

# D-001 · 协议分流、上游增补报告与协议先行整改门禁

## 触发

用户于 2026-09-06 要求在 workspace-010 开启新子目标，审视当前 API/Web 页面与控件是否遵循上游 `schema-ui-docs`；必须区分“上游已有协议但本仓未遵循”和“上游缺乏合理协议”，并在后者中进一步判断应推动上游增补还是在上游允许范围内使用本仓 custom 控件。若存在需上游增补项，须在本目标附件汇总报告，并把上游协议增补作为对应本仓修改的前置门禁。

立项前发现 VP-010 / workspace-010 / Root 对 Charter `@0.4.0` 的投影漂移；用户书面选择“对齐后建目标”。该 re-align 已由 VR-071 / VRev-081 记录后，本目标才建立为 W29。

## 决定

### 1. 协议分母

- 起始权威身份：现行 Charter 声明的 `schema-ui-docs@v2.9.0`、pinned commit `81aa1d8`。
- S1 必须重新核验上游正式工件、本仓 `provenance-v2.9.json`、schemas、capability/component registry、fixtures、runtime Manifest 与测试入口；Charter 声明本身不替代字节、schema、行为或运行时证据。
- v2.7/v2.8 保留为 additive 兼容历史与 GOAL-004 证据，不得被误写成当前全部页面控件的完成证明。

### 2. 候选分流

每个页面/控件候选必须进入有证据的唯一处置类别：

| 类别 | 判定与动作 |
|------|------------|
| `implementation-gap` | 上游已有可消费协议语义，本仓实现偏离；按上游修本仓，并补正反测试与运行时证据。 |
| `upstream-protocol-gap` | 合理需求存在，但上游缺少足够契约；先写上游增补报告/提案，等待正式协议身份固定，再实施依赖该语义的本仓修改。 |
| `custom-extension-candidate` | 上游明确允许扩展或明确不负责；只有 namespace/capability/schema/validator/failure/compatibility/fixtures 完整且用户书面裁决后，才可成为本仓合法 custom。 |
| `explicitly-out` | 明确属于 Host、身份/部署安全、业务域或其他责任层；记录归属与边界，不伪装成页面协议能力。 |
| `excluded` | 本波不覆盖；必须写理由、责任人或复核触发，不能用排除隐藏 required 缺口。 |

候选证据不足时保持 `open` / `collecting`，不得提前归类。

### 3. 上游协议增补报告与停止线

- 若确认存在一个或多个 `upstream-protocol-gap`，创建并维护 `attachments/upstream-protocol-augmentation-report.md`；至少包含候选 ID、用户/产品需求、当前实现证据、上游缺口证据、建议 shape/state/error/security/capability/fixture、兼容与迁移影响、优先级及上游跟踪引用。
- 对依赖该增补的本仓范围设置停止线：上游 accepted/merged、正式 tag/version/commit 可核对、本仓 provenance/schemas/registry/fixtures/兼容证据固定前，**禁止进入 S4 实施或宣称放行**。
- 上游没有增补项时，不创建空报告；以分类矩阵与 I-004 的“不适用”证据收口。

### 4. custom 门禁与用户裁决

`custom-extension-candidate` 不是默认逃生口。S3 前必须逐项向用户展示“推动上游增补 / 本仓 custom / explicitly-out”互斥方案与风险；未经用户书面裁决，不得选择 custom。合法 custom 至少固定 namespace、能力声明、schema/validator、失败语义、兼容策略、fixtures 与未来迁移/退出触发。

### 5. 审计与 go 影响

- 审计模式固定为 `cross`：S2 方案冻结与 S6 关门均要求 self + independent；independent provider 按项目级决策使用本地 `grok build`。provider 不可用时门禁保持未满足。
- 若整改改变 Profile 默认集、模块矩阵、Manifest 装配语义或共同门禁解释，进入 VP-008 `go` 消费暂挂/恢复路径；不凭“测试绿”静默判无影响。

## 未选方案

- 把所有手写 Host UI、Login、branding、session、Shell、全局导航或错误边界一概判为协议绕过。
- 仅扫描静态 JSON/schema，不核验 API/Web 运行时与失败路径。
- 在上游协议未固定时先改本仓，再以既成实现倒逼协议。
- 用无 namespace、无 validator、无兼容/迁移证据的 `custom` 节点规避上游协议。
- 重开或改写 GOAL-004 的历史结论；W29 只承接当前 v2.9 与现行实现的新增核验。

## 后续

S1 分母与目录 → S2 分类/方案 cross 审视 → S3 上游协议或 custom 边界固定 → S4 实现整改 → S5 运行时验证与 go 影响判定 → S6 cross 关门审计。当前只完成立项与治理门禁登记，未形成任何页面/控件符合性 verdict。
