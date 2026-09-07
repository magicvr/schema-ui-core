---
title: S3 · custom 扩展边界规范（C-009 裁决产物）
status: active
created: 2026-09-06
updated: 2026-09-06
parent: GOAL-041-w29-api-web-protocol-conformance
version: 0.1.0
---

# S3 · custom 扩展边界规范

> 本文件固定 C-009 的合法 custom 边界，是 D-001 §4 custom 门禁的落地产物。用户 P-004 书面裁决（2026-09-06）：**① 处置路径 = 本仓合法 custom；② namespace 策略 = 保留现有 15 键 + 所有权登记 + 新键规范（不批量重命名）；③ C-005 子项 = 删除 digitaloffer 两页未使用能力声明（S4 实施）**。裁决留痕见 [D-003-s3-custom-boundary-fix](../01-decision/D-003-s3-custom-boundary-fix.md)。

## 1. 协议依据（上游）

- `08-renderer-spec.md §1.1`（schema-ui-docs@v2.9.0）：核心组件注册表只包含 `component-registry.json` 声明的类型；宿主可通过私有 API（如 `Renderer.register()`）安装 **Host Extension**；扩展**不属于** v2.0 页面协议的核心身份，扩展页面不得以核心协议身份跨 Renderer 互操作；扩展必须由宿主显式安装，未安装的扩展 type 在标准入口以 `UNKNOWN_COMPONENT_TYPE` 拒绝，不得静默降级为核心组件。
- 上游通用组件注册表共 24 个组件（grid/section/tabs/statCard/chart/text/recordView/table/actionButton/form/input/…/password），均为通用 UI 语义，不覆盖业务域控件。
- `07-actions-contract.md §6`：`type: custom + handler` 仅允许前端白名单预注册处理函数（动作层扩展，与节点层 custom 并列的独立扩展面）。

## 2. 本仓 custom 扩展面（15 键 · 所有权登记）

| key | 使用形态（次数） | 所属页面 | 归属模块 | 注册证据 |
|-----|------------------|----------|----------|----------|
| `account-session-toolbar` | body custom ×1 | account | admin.account | `apps/web/src/components/account-session-toolbar.tsx:84` |
| `activity-export` | body custom ×1 | activity | admin.activity | `apps/web/src/components/activity-export.tsx:84` |
| `cron-preview` | afterComponent ×2 | scheduled-tasks / task-runs | admin.scheduled-tasks | `apps/web/src/components/cron-preview.tsx:134` |
| `data-permission-scopes` | body custom ×1 | data-permission | admin.data-permission | `apps/web/src/components/data-permission-scopes.tsx:282` |
| `email-identity` | body custom ×1 | account | admin.account | `apps/web/src/components/email-identity.tsx:195` |
| `import-template-download` | afterComponent ×1 | users | admin.users | `apps/web/src/components/import-template-download.tsx:70` |
| `invite-issue-card` | body custom ×1 | users-invites | admin.users | `apps/web/src/components/invite-issue-card.tsx:248` |
| `invite-resend-dialog` | action content custom ×1 | users-invites | admin.users | `apps/web/src/components/invite-resend-dialog.tsx:133` |
| `mail-admin-tab` | body custom ×1 | mail | admin.settings | `apps/web/src/components/mail-admin-tab.tsx:350` |
| `mfa-manager` | body custom ×1 | account | admin.account | `apps/web/src/components/mfa-manager.tsx:379` |
| `monitoring-auto-refresh` | body custom ×1 | system-monitoring | admin.system-monitoring | `apps/web/src/components/monitoring-auto-refresh.tsx:62` |
| `notification-center` | body custom ×1 | notifications | admin.notifications | `apps/web/src/components/notification-center.tsx:250` |
| `password-policy-tab` | body custom ×1 | settings | admin.settings | `apps/web/src/components/password-policy-tab.tsx:137` |
| `telegram-admin-tab` | body custom ×2 | telegram-settings / telegram-operator | channel.telegram | `apps/web/src/components/telegram-admin-tab.tsx:1089` |
| `wallet-ensure` | body custom ×1 | my-wallet | admin.wallet | `apps/web/src/components/wallet-ensure.tsx:112` |

使用形态合计：body custom 13 处（12 键，telegram-admin-tab 双页）、action content 1 处、afterComponent 3 处（2 键）。全部 15 键均有生产注册（`main.tsx` side-effect import）与守卫测试覆盖（`custom-components.schema.test.ts`，S2 已递归化）。

## 3. 边界决定（用户 P-004 书面裁决后固定）

### 3.1 处置路径

**本仓合法 custom**。依据：上游 §1.1 明确允许且不负责 Host Extension；15 键均为本应用业务域控件（钱包/电报/邮件/MFA/通知/权限/监控等），不属于上游通用组件语义；推动上游增补会把业务域控件灌入通用协议注册表（与协议定位冲突），explicitly-out 与「schema 内声明 custom 节点」的事实矛盾。裁决见 D-003。

### 3.2 namespace 约定

1. **现有 15 键保留不变**（兼容性：schema/registry/fixtures/e2e 零破坏）；本表即所有权登记。
2. **新键规范**：`<module-scope>-<feature>`（如 `wallet-ensure` 已合规范；`channel.telegram` 模块键建议形如 `telegram-<feature>`）。登记时同步更新本表与 playbook §6.2。
3. **防碰撞规则**：custom 键必须全局唯一，且**不得与上游 `component-registry.json` 任何核心类型重名**；唯一性由注册 Map 与守卫测试保证，核心类型碰撞由注册评审（playbook §6.2）把关。

### 3.3 capability

custom 节点**不要求新增协议 capability**：扩展面是 Host-owned 呈现能力，不进入 claim `support.capabilities`；页面 `meta.requiredCapabilities` 仍只声明协议特性门禁（F-001 的 claim 覆盖问题与此无关，S4 另行处理）。custom 节点的动作层 handler（C-011）走 `actions.custom` 契约，亦无新增 capability。

### 3.4 schema / validator

- 结构校验：`apps/web/src/protocol/conformance/runtime-schema-validate.ts` 在 vendored `node.schema.json` 上叠加本地 `component` 属性（**上游 schema 字节不变**，provenance digest 守卫已核）；`custom` 节点必须携带 `component: string`。
- D-VAL：全部 35 页（含 13 处 custom 使用）过 `validatePageDocument`，S2 后实测 35/35 绿。
- 新 custom 键必须过注册守卫（`custom-components.schema.test.ts` 递归扫描全部模块 schema 的 `{type:"custom",component}` 引用）。

### 3.5 failure 语义（与 C-010 联动）

- 未知**标准 node**：`RENDER_UNKNOWN_NODE_TYPE` fail-closed（现有行为，符合上游 §1.1 拒绝语义）。
- 未知 **custom component**：S4 按 C-010 实施对齐——渲染明显占位（含 component 键）+ `console.error`（含 node id），或标准入口 fail-closed `UNKNOWN_COMPONENT_TYPE`；当前内联文字 fallback 为过渡行为。
- 未知 **custom handler**（动作层）：`CUSTOM_HANDLER_NOT_FOUND` fail-closed（已符合 07-actions-contract §6）。

### 3.6 compatibility

- 扩展页不得以核心协议身份跨 Renderer 互操作（上游 §1.1）：含 custom 节点的页面 schema 是**本 Host 特定**的，跨 Renderer 消费不保证。
- 扩展不改变上游 schema 字节与 24 组件注册表；`custom` 类型 + `component` 键是本仓 Host 的私有扩展入口。

### 3.7 fixtures

- 结构：D-VAL 35/35（含全部 custom 使用页）。
- 注册：`custom-components.schema.test.ts` 守卫（递归 + 全部键 import）。
- 行为：`representative-pages.test.tsx` 已渲染含 custom 的页面（users-invites 等 10 页）；unknown-custom 失败用例（C-010）在 S4 补齐；v2.9 页定向测试（dictionary-entries/wallet-entries）不涉 custom。
- S5：按 D-002 §3 验收矩阵对含 custom 页面补真实链路触达。

### 3.8 退出 / 迁移触发

| 触发 | 动作 |
|------|------|
| 上游未来引入等价核心组件 | 评估迁移：替换 custom 节点为核心组件 + 删除注册 + 更新本表与守卫；需 new custom 裁决或 S6 复审 |
| custom 键不再被任何 schema 引用 | 移除注册 + 组件文件（若有） + 更新本表；守卫测试自动转红直至移除 |
| 新增页面 / 新 custom 键 | 按 §3.2 规范命名 + 注册 + 本表登记 + 守卫覆盖 |
| 上游协议版本线更新（如 v2.10+） | 复核 §1.1 扩展模型是否变化；变化时按 S6 复审路径处理 |

## 4. C-005 子项（用户裁决）

digitaloffer 两页「声明未使用」的能力声明**删除**（S4 实施）：`digitaloffer-entitlements.json` 移除 `meta.requiredCapabilities` 中的 `data.route-binding`；`digitaloffer-offers.json` 移除 `form.controls.readonly`。同时补一条守卫断言（页面声明的能力 ⊆ 实际使用或文档化意图），防止再次出现声明-使用漂移。删除动作属于 D-002 §4 门禁下 S4 实施范围（不改变协议语义，仅让声明如实）。

## 5. 与本目标其他项的关系

- C-010（未知 custom 呈现）与 C-011（动作层白名单）已在此边界内定位；C-010 的 S4 实施受 §3.5 约束。
- F-001（生产页面级能力门禁 + claim/HOST_SUPPORT 覆盖）**不属于** custom 边界，S4 另行实施。
- 上游协议增补报告：不创建（upstream-protocol-gap = 0，I-004 不适用）。
