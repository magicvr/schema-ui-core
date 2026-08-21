---
id: E-006
goal: GOAL-015-dict-inner-page-breadcrumb
date: 2026-08-14
status: recorded
parent: GOAL-015-dict-inner-page-breadcrumb
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# E-006 · vendor 重 pin：上游 v2.9.0（81aa1d8）

## 事实

### 1. 上游核查

- 上游 HEAD = 81aa1d8（v2.9.0，PR #39 feat/2.9-route-binding-readonly，审计 0082 A-002 round 2 pass）。
- 相对本仓 pin（521cff8 / v2.8.0）：docs/schemas 4 件变更（capability-registry、component-registry、node.schema、page.schema——新增 `data.route-binding` 与 `form.controls.readonly` capability、params 路由绑定描述、表单字段 readOnly 声明）；20 个 conformance fixture 套件算法线升 2.9（fixture digest 89baddbc…）；request-construction +6（data-ref-route-* + readonly 投影）、version-negotiation +4、scenarios +1（admin-list-route-filter-lifecycle——即本目标的业务闭环样例）。

### 2. 重 pin 动作

- docs/schemas 11 件与 apps/web/src/protocol/upstream/ 19 个 cases.json 全量替换为 81aa1d8 字节（scenarios 未 vendor）。
- 新建 `provenance-v2.9.json`（30 项 artifact，LF 规范化 sha256）；legacy provenance.json 中 app-manifest/app-navigation cases 条目同步（注记指向 v2.9 pin）。
- 测试 pin 更新：stage3-fixtures → provenance-v2.9.json（SOURCE_COMMIT 81aa1d8 / 2.9.0）；upstream-fixtures APP_MANIFEST_FIXTURE/APP_NAVIGATION_FIXTURE SHA；upstream-host-fixtures 3 个 host fixture SHA；permissions-inheritance SHA；APP_MANIFEST_SOURCE → 81aa1d8。

### 3. 适配器修正（消费侧 2.9 支持）

- request-construction.ts `buildDataRef`：params 支持 `$context.route.query.*` / `$context.route.params.*` 整值绑定，缺失键/无 route 输入 → tombstone 删除（data-ref-route-* 5 用例）。
- app-manifest.ts：APP_MANIFEST_SUPPORTED_PROTOCOL_VERSIONS += 2.9；APP_MANIFEST_PROTOCOL_VERSION 2.9；returnIntentQueryKeys 版本门槛改为 >= 2.8（2.9 manifest 通过）；APP_MANIFEST_SOURCE 更新。
- upstream-fixtures.test.ts hostManifestValue 2.9 直通（此前 2.9 会被改写为 2.7）。

### 4. 验证

- protocol + host 套件 489/489 绿（含新 data-ref-route-*、readonly 投影、2.9 协商用例）。
