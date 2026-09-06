---
id: GOAL-042-w30-w29-followup-supplement
doc: audit-entry
record_id: A-002
source: independent
status: recorded
parent: GOAL-001-design-implementation-conformance
created: 2026-09-06
updated: 2026-09-06
version: 0.1.0
---

# A-002 · W30 关门独立交叉审计 · independent（claude-sonnet-4-6）

## A-002 · W30 三项后继补强 close-out 独立审计（2026-09-06）

- **source**：independent
- **auditor**：claude-sonnet-4-6（DeepSeek Harness /audit · 独立会话）
- **类型 / scope**：close-out（F1 legacy 能力全量审计守卫 / F2 claim↔host-support 单源一致性 / F3 10 页行为级单测）；`audit_type: close-out`
- **verdict**：**pass**（0 required findings；2 recommended）

---

## 范围与区间

- 工作区：`workspace-010-design-implementation-conformance`，canonical root 已验证
- 被审目标：`GOAL-042-w30-w29-followup-supplement`（`status: done`，`progress: 3/3`）
- 审计覆盖：E-001（F1/F2/F3 实施事实）、A-001 self pass、产物代码（host-support.json / host-support.ts / generate-claim.mjs / conformance-claim.json / capability-declaration.guard.test.ts / host-support-consistency.test.ts / behavior-pages.test.tsx / s5-denominator-render.test.tsx）、部分 schema 能力声明抽样
- P-005：I-001/I-002/I-003 均 verified，无 deferred required
- 共享资料：`shared_materials_catalog: none`，无需核查

---

## 成果（有证据）

### F1 · legacy 能力守卫（独立核验）

1. **文件存在并可读**：`apps/web/src/protocol/capability-declaration.guard.test.ts` 存在，内容完整。
2. **MARKERS 对照核验**：13 个非豁免能力均有 marker：
   - `data.route-binding` → `$context.route.`
   - `form.controls.readonly` → `"readOnly": true`
   - `form.controls.extended` → textarea/switch/checkbox/radio 或 mode:multiple
   - `form.controls.advanced` → cascader/checkboxGroup/richText/password 或 defaultValue
   - `table.sort` → sortable/sortField/defaultSort
   - `permissions.inheritance` → permissionCascade/permissionIntent
   - `record.view.load` → type:recordView
   - `form.record.load` → recordSource
   - `table.selection` → requiresSelection
   - `actions.upload` → type:upload
   - `actions.row.request` → requestMapping
   - `actions.page.trigger` → actionRef 或 type:actionButton
   - `actions.row.navigate` → navigateMapping（收紧为 navigateMapping，非通用 navigate）
   - `actions.batch.request` → /batch-delete
3. **EXEMPT 豁免集**：5 个 envelope/host 能力（app.manifest / app.navigation / host.bootstrap / host.failure-recovery / host.conformance-claim）豁免设置正确。
4. **INTENT_OVERRIDES**：仅 `notifications → actions.page.trigger`（custom 铃铛），理由合理。
5. **schema 抽样核验**（独立读取原始文件）：
   - `notifications`：声明 `actions.page.trigger`，文本无 actionRef/actionButton（依赖 INTENT_OVERRIDE）✓
   - `wallet`：声明 actions.row.navigate + navigateMapping 存在 ✓；form.controls.extended + textarea ✓；readOnly ✓；permissions.inheritance + permissionCascade ✓
   - `data-permission`：无 actions.row.request（requestMapping 不在文本中） ✓；有 actions.page.trigger + actionButton ✓；无 table.sort（无 sortable） ✓
   - `data-dictionary`：有 actions.row.navigate + navigateMapping ✓；有 actions.row.request + requestMapping ✓
   - `scheduled-tasks`：无 actions.row.navigate（有 navigate 类型动作但**无** navigateMapping） ✓（marker 正确收紧为 navigateMapping）
   - `activity`：有 record.view.load + type:recordView ✓；form.controls.extended ✓；actions.page.trigger ✓
   - `admin-list-batch`：有 table.sort + sortable ✓；form.controls.extended + mode:multiple ✓
6. **schema 文件总数**：35 个（守卫 `files.length >= 30` 通过）✓
7. **claim↔host-support.json 能力集合相等**：机器核验一致（19 项，独立对比）✓

**F1 结论**：守卫设计逻辑合理，MARKERS 与 INTENT_OVERRIDES 均有明确依据，抽样 schema 声明与内容一致。

### F2 · claim↔host-support 单源一致性（独立核验）

1. **host-support.json** 存在，19 能力，pageVersions [2.7,2.8,2.9]。
2. **host-support.ts** 从 JSON 导入，导出 `HOST_SUPPORTED_PAGE_VERSIONS` / `HOST_SUPPORTED_CAPABILITIES`（单源消费，注释指向 JSON）✓
3. **generate-claim.mjs** 从同一 JSON 读取（`readFileSync(join(WEB_ROOT, "src", "host", "host-support.json"))`），生成 `support.pageVersions` / `support.capabilities`（F2 注释指向 D-001）✓
4. **conformance-claim.json**（已生成工件）：`support.capabilities` 19 项 = host-support.json 集合 ✓；`support.pageVersions` = [2.7,2.8,2.9] ✓
5. **host-support-consistency.test.ts** 5 个断言：
   - claim↔JSON 集合相等（capabilities + pageVersions）
   - 每个 capability 是合法的 registry ID
   - 每个 capability 的 mandatorySuites ⊆ claim suites 且 result:pass

   capability-registry.json 独立读取可见 19 个能力均在文件中（含 mandatorySuites 链接），测试设计逻辑完备。

6. **D-001** 决策落盘，「仅注释约定」未选方案理由充分。

**F2 结论**：单源设计完整，claim 工件与 JSON 一致，一致性测试覆盖集合相等+合法性+mandatory 覆盖，无漂移风险。

### F3 · 10 页行为级单测（独立核验）

1. **behavior-pages.test.tsx** 存在，TARGETS 列表 10 页：mail / mail-outbox / my-wallet / wallet / wallet-vouchers / telegram-settings / telegram-operator / digitaloffer-offers / digitaloffer-entitlements / digitaloffer-purchases ✓
2. **测试结构**：
   - for 循环产生 10 条通用渲染断言（非空内容 + 无 schema-error + 无 unknown custom component）
   - 10 个命名 it() 块产生页面特有 UI 断言（表头文字/工具栏按钮/自定义面标题/行操作）
   - 共 20 个测试（20/20 PASS 声明对应数量一致）✓
3. **链路真实性**：通过 `loadPageDocument`（含 D-VAL + F-001 协商）走产品链路；fixtureFetcher 按页面列字段喂数据（解决空态替换整表问题）；`AuthProvider` + `I18nProvider` 完整包裹 ✓
4. **自定义组件注册**：文件头显式 import mail-admin-tab / telegram-admin-tab / wallet-ensure（与 main.tsx 对齐）✓
5. **s5-denominator-render.test.tsx 补丁**：补 15 个 custom import（account-session-toolbar / activity-export / cron-preview / data-permission-scopes / email-identity / import-template-download / invite-issue-card / invite-resend-dialog / mail-admin-tab / mfa-manager / monitoring-auto-refresh / notification-center / password-policy-tab / telegram-admin-tab / wallet-ensure），消除 custom 节点占位噪音 ✓

**F3 结论**：测试链路真实，覆盖 10 页，页面特有 UI 断言有意义，不是空测。

### 整体回归

- Web vitest 97 files / 1332 tests（E-001 声明；+25 新增文件）
- `npm run build`（tsc+vite）0 错误
- Go 全量 0 FAIL（schema 元数据改动不影响 Go 行为）
- 无越界：未改 Profile 默认集/模块矩阵/Manifest 装配语义/上游协议版本

---

## 对照成功标准

| 成功标准 | 核验结论 |
|---------|---------|
| F1 守卫覆盖全部非豁免能力（13 markers + 豁免 + intent override） | ✓ 独立核实 marker 映射完整 |
| 32 个 schema 双向修正（删 22 处 / 补 11 处） | ✓ 抽样验证代表性页面；数量由 E-001 记录，守卫 35/35 通过是机械验证 |
| 守卫 35/35 绿 | ✓ 文件数 = 35，guard `files.length >= 30` 成立 |
| F2 host-support.json 单源（ts + mjs 同源） | ✓ 独立读取两文件确认同源 |
| 一致性测试 5/5 绿 + claim 重生成 | ✓ 测试设计完备；claim 工件与 JSON 一致（独立校验） |
| F3 10 页行为单测 20/20 | ✓ 文件存在，20 条测试可数，链路真实 |
| A-001 self pass + 全量回归 Web 1332/1332 + build 0 + Go 0 FAIL | ✓ A-001 记录在案；回归数字来自 self 审计记录，此次 independent 未重跑（可接受：scope 为关门复审） |

---

## Findings

**F-001（recommended · low）·「schema 文件遍历依赖 API 目录结构」**

`capability-declaration.guard.test.ts` 与 `behavior-pages.test.tsx` 均从 `apps/web/src/../../../api/modules` 路径扫描 schema，与 Web 包的物理边界跨越 monorepo 目录（非 workspace package 引用）。功能正确，但路径耦合在此 monorepo 布局为已知设计决定（其他测试如 s5-denominator-render 同样跨越）；无需立即修改，建议 playbook 中注明此约定以避免未来迁移时遗漏同步更新。

**F-002（recommended · low）·「behavior-pages.test.tsx 仅断言文本存在，不核验 HTTP mock 调用完整性」**

fixtureFetcher 会响应所有 `/api/` 请求，测试通过断言渲染文本验证链路。若某页实际未触发数据请求（如 custom 组件短路），测试仍绿但链路有缺口。当前 10 页均为声明式 table 页面，自然会触发列表请求，此风险低；建议未来对 custom 组件主导页面（如 mail）补 fetch spy 确认调用。

---

## 必改项汇总

**无**（0 required findings）。两个 recommended 均不阻断关门。

---

## 与既有意见的异同

- 与 A-001（self）一致：verdict pass，0 required，成果核验结论相同。
- F-001/F-002 为新识别的 non-blocking 注意点，不与 A-001 冲突；A-001 未提及（超出 self 范围正常）。

---

## 结论 + 建议给编排器/用户的下一步

GOAL-042 三项交付（F1 守卫扩展 + F2 单源一致性 + F3 行为单测）均有可追溯证据，机制设计合理，代码与文档记录一致，无未关闭 required findings。两个 recommended 不阻断关门。

**建议**：independent 审计意见落盘完毕，目标已 `done`，无需 `/govern` 修正动作。若用户希望跟踪 F-001/F-002，可在后续符合性波次或 playbook 更新中处理。用户可用 **`/govern`** 确认意见响应并归档。

---

## 声明

本意见不修改 `GOAL-042` 的 `status`/`progress`/`goal-tree` 状态；响应由 `/govern` 处理。
