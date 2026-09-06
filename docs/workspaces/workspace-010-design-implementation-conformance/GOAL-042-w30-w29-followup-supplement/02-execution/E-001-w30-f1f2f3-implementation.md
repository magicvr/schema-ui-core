---
id: GOAL-042-w30-w29-followup-supplement
doc: execution-entry
record_id: E-001
status: recorded
parent: GOAL-001-design-implementation-conformance
created: 2026-09-06
updated: 2026-09-06
version: 0.1.0
---

# E-001 · W30 执行 — F1/F2/F3 三项后继补强

## 已执行事实（2026-09-06）

### F1 · legacy 能力保守声明全量审计（capability 声明↔使用双向对齐）

1. `apps/web/src/protocol/capability-declaration.guard.test.ts` 从「仅 v2.9 能力对」扩展为**全部非豁免能力**（13 个 marker + envelope/host 豁免 + `INTENT_OVERRIDES` 仅 notifications→actions.page.trigger，custom 铃铛消费）。
2. 双向审计结果并修正 32 个 schema 文件的 `meta.requiredCapabilities`：
   - **删除 22 处未使用保守声明**：permissions.inheritance ×10（account/telegram-settings/admin-list-batch/notifications/recycle-bin/task-runs/mail-outbox/system-monitoring/wallet-entries 等）、actions.row.request ×6（activity/data-permission/digitaloffer-offers/notifications/wallet/wallet-vouchers）、form.controls.extended ×9（form-with-upload/search-form-table/file-library/notifications/settings/users-invites/users/my-wallet/wallet-vouchers）、table.sort ×1（data-permission）、actions.row.navigate ×1（scheduled-tasks，页面级 navigate 非 row）。
   - **补 11 处欠声明**：record.view.load ×5（activity/mail-outbox/users/users-invites/roles）、actions.page.trigger ×4（activity/form-with-upload/recycle-bin/digitaloffer-entitlements）、actions.row.navigate ×1（data-dictionary，确有 navigateMapping）、table.sort ×1（admin-list-batch，确有 sortable）+ form.controls.extended ×1（admin-list-batch，确有 mode:multiple）。
   - 关键判据核对：navigateMapping 仅 data-dictionary/wallet 使用（row 级导航）；data-permission 表无 sortable；search-form-table 有 sortable（保留声明）；users/roles 有 recordView。
3. 守卫扩展后 **35/35 PASS**；`actions.row.navigate` marker 收紧为 `navigateMapping`（页面级 navigate 动作不计数）。

### F2 · claim↔host-support 机械一致性（单源 JSON + 一致性测试）

4. 新建 `apps/web/src/host/host-support.json`（supportedPageVersions [2.7,2.8,2.9] + supportedCapabilities 19 项）为**唯一权威**：
   - `host-support.ts` 从 JSON 导入并导出 `HOST_SUPPORTED_PAGE_VERSIONS` / `HOST_SUPPORTED_CAPABILITIES`（load-page/boot 消费不变）。
   - `generate-claim.mjs` 构建期读取同一 JSON 生成 `support.pageVersions` / `support.capabilities`（`ARTIFACT_VERSION 2.9.0`、12 suites 不变）。
5. 新增 `apps/web/src/host/host-support-consistency.test.ts`（5 断言）：claim↔JSON 能力集/页面版本集合相等、能力 ID 合法（⊆ capability-registry）、每个能力 mandatorySuites ⊆ claim suites 且 pass。**5/5 PASS** + claim-artifact C0/C1 **5/5 PASS**；重生成 claim 三件套。

### F3 · 10 页行为级单测

6. 新增 `apps/web/src/renderer/behavior-pages.test.tsx`：mail / mail-outbox / my-wallet / wallet / wallet-vouchers / telegram-settings / telegram-operator / digitaloffer-offers / digitaloffer-entitlements / digitaloffer-purchases——每页走真实链路（D-VAL → loadPageDocument 含 F-001 协商 → RenderPage），按页面自身列字段喂一条样本行（空态会替换整表，故喂数据使表头/行操作渲染），断言页面特有 UI（表头/工具栏/自定义面标题/行操作）。**20/20 PASS**。
7. 附带修复 `src/i18n/s5-denominator-render.test.tsx`：补 15 个 custom 组件 side-effect import（真实 App 渲染 custom 节点，消除 activity-export 等占位噪音）。

## 验证

- 定向：capability-declaration 35/35、host-support-consistency 5/5、claim-artifact 5/5、behavior-pages 20/20、denominator 36/36、D-VAL 38。
- 全量回归：见下方补录（Web vitest 全量、`npm run build`、Go 全量）。

## 产物

- `apps/web/src/host/host-support.json`（新，单源）
- `apps/web/src/host/host-support.ts` / `scripts/generate-claim.mjs`（JSON 同源）
- `apps/web/src/host/host-support-consistency.test.ts`（新）
- `apps/web/src/protocol/capability-declaration.guard.test.ts`（全量能力守卫）
- `apps/web/src/renderer/behavior-pages.test.tsx`（新）
- `apps/web/src/i18n/s5-denominator-render.test.tsx`（补 custom import）
- 32 个 API 模块 schema 的 `meta.requiredCapabilities` 修正
- 重生成 `apps/web/public/protocol/`（claim 三件套）

## 阶段结论

F1/F2/F3 全部完成并有可核对证据；守卫/一致性测试成为防复发机制。审计与关门见 03-audit。
