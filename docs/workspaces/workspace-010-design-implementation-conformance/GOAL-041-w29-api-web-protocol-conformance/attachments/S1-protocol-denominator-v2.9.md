---
title: S1 · schema-ui-docs v2.9.0 协议分母
status: active
created: 2026-09-06
updated: 2026-09-06
parent: GOAL-041-w29-api-web-protocol-conformance
version: 0.1.0
---

# S1 · schema-ui-docs v2.9.0 协议分母

> 冻结 S1 对照对象，不在本附件判定 implementation gap、upstream protocol gap 或 custom 方案。完整机读记录见 `S1-protocol-denominator-v2.9.json`。

## 1. 固定身份

- 上游：`https://github.com/magicvr/schema-ui-docs`；本机只读 clone：`C:\Users\magicvr\Documents\Code\schema-ui-docs`。
- annotated tag：`v2.9.0`，tag object `463d563a243705380b4edb20c7c155f7fe65a02d`。
- tag 解引用 commit：`81aa1d8954717f4ebdcc695eed6fafaeafcebe8d`（`81aa1d8`）。
- 本仓被审快照：`ce66abab0897b644c695b8365f8324598bfe2bc8`。
- 核验入口：`git -C ../schema-ui-docs cat-file -p v2.9.0` 与 `git -C ../schema-ui-docs rev-parse v2.9.0^{commit}`。

## 2. 机器分母

- schema/registry 工件：**11**。
- component registry：**24** 个组件。
- capability registry：协议版本 `2.9`，**19** 个 capability。
- 上游 fixture：**20 suites / 450 cases**。
- 本仓 vendor：**19 suites / 437 cases**；未 vendor `scenarios`（13 cases）。
- `provenance-v2.9.json`：**30 artifacts**。

### 2.1 schema / registry provenance

| provenance 路径 | 上游 tag 路径 | 本仓路径 | digest | 结果 |
|---|---|---|---|---|
| `docs/schemas/action.schema.json` | `docs/schemas/action.schema.json` | `docs/schemas/action.schema.json` | `2f929e1fa4706a66a6e14bb4c58b2cd8a7c8c7b01f8e624c9c56adfb440af5c0` | byte match |
| `docs/schemas/app-manifest.schema.json` | `docs/schemas/app-manifest.schema.json` | `docs/schemas/app-manifest.schema.json` | `34a3354e245dbf3900744b5797edeb1ca5f2ac19872ac908d781274d47d68c55` | byte match |
| `docs/schemas/capability-registry.json` | `docs/schemas/capability-registry.json` | `docs/schemas/capability-registry.json` | `d7fa4043383be337fbc235d726f73c00ae0ac8ab614cca76474246b180585315` | byte match |
| `docs/schemas/component-registry.json` | `docs/schemas/component-registry.json` | `docs/schemas/component-registry.json` | `0714e1339b175a0822398d6105467aeb2dcba663903c141fee11d7911b2d21a9` | byte match |
| `docs/schemas/fixture-suite.schema.json` | `conformance/schemas/fixture-suite.schema.json` | `docs/schemas/fixture-suite.schema.json` | `a93f500bcb94f93e7a289af700227a407688f8205218fe5681269956b91fd1a3` | byte match; path alias differs |
| `docs/schemas/host-bootstrap.schema.json` | `docs/schemas/host-bootstrap.schema.json` | `docs/schemas/host-bootstrap.schema.json` | `e861430275c24ca5108e90ae9a15a8a6feb72f7871983ccb76fadf5384d3e55e` | byte match |
| `docs/schemas/host-conformance-claim.schema.json` | `docs/schemas/host-conformance-claim.schema.json` | `docs/schemas/host-conformance-claim.schema.json` | `da6389f874787e377a80ea3f1c4364605941220b856091adade561db6877a74e` | byte match |
| `docs/schemas/host-failure.schema.json` | `docs/schemas/host-failure.schema.json` | `docs/schemas/host-failure.schema.json` | `1923c5710173e601fb3052c97a8ab0d1ccf88d00fb24b8b660f7979e051531be` | byte match |
| `docs/schemas/node.schema.json` | `docs/schemas/node.schema.json` | `docs/schemas/node.schema.json` | `deade0b3146415856584ab9dbaa197db077ceb4b3f1e5512c6be449e65392630` | byte match |
| `docs/schemas/page.schema.json` | `docs/schemas/page.schema.json` | `docs/schemas/page.schema.json` | `213d9df40194221f4461fe2246601e53b040d7fab64d5b29a9021780669b0923` | byte match |
| `docs/schemas/reaction.schema.json` | `docs/schemas/reaction.schema.json` | `docs/schemas/reaction.schema.json` | `c35f1beb6b4a717114aeae5b7bc5ba523b507177658e44bc60752e70f8a415f4` | byte match |

> 重要口径：上游 tag 中 fixture suite schema 的真实路径是 `conformance/schemas/fixture-suite.schema.json`；本仓 provenance 把它登记为 `docs/schemas/fixture-suite.schema.json`。字节 digest 一致，但来源路径不一致，已登记 C-001。

### 2.2 component registry（24）

`grid`, `section`, `tabs`, `statCard`, `chart`, `text`, `recordView`, `table`, `actionButton`, `form`, `input`, `textarea`, `switch`, `checkbox`, `radio`, `inputNumber`, `datePicker`, `dateRangePicker`, `upload`, `select`, `cascader`, `checkboxGroup`, `richText`, `password`

### 2.3 capability registry（19）

`actions.upload`, `actions.row.request`, `actions.page.trigger`, `actions.row.navigate`, `form.record.load`, `table.selection`, `actions.batch.request`, `permissions.inheritance`, `record.view.load`, `app.manifest`, `app.navigation`, `table.sort`, `form.controls.extended`, `form.controls.advanced`, `host.bootstrap`, `host.failure-recovery`, `host.conformance-claim`, `data.route-binding`, `form.controls.readonly`

### 2.4 fixture suites

| suite | 上游 cases | 本仓 vendor | 本仓 cases | 结果 |
|---|---:|---|---:|---|
| `actions` | 11 | yes | 11 | match |
| `app-manifest` | 41 | yes | 41 | match |
| `app-navigation` | 16 | yes | 16 | match |
| `component-format` | 5 | yes | 5 | match |
| `host-bootstrap` | 23 | yes | 23 | match |
| `host-conformance-claim` | 30 | yes | 30 | match |
| `host-failure` | 43 | yes | 43 | match |
| `permissions-inheritance` | 17 | yes | 17 | match |
| `query-serialization` | 16 | yes | 16 | match |
| `reactions` | 16 | yes | 16 | match |
| `request-construction` | 81 | yes | 81 | match |
| `request-lifecycle` | 4 | yes | 4 | match |
| `response-mapping` | 23 | yes | 23 | match |
| `runtime-defaults` | 9 | yes | 9 | match |
| `scenarios` | 13 | no | — | not pinned |
| `search-table` | 11 | yes | 11 | match |
| `static-data` | 9 | yes | 9 | match |
| `table-sort` | 14 | yes | 14 | match |
| `uploads` | 13 | yes | 13 | match |
| `version-negotiation` | 55 | yes | 55 | match |

## 3. v2.9 增量

- `data.route-binding`（ADR-0039）：DataRef `params` 对 `$context.route.query.*` / `$context.route.params.*` 的完整单值绑定；要求页面 2.9 与 capability。
- `form.controls.readonly`（ADR-0040）：字段不可编辑但继续参与 values 与提交投影；不等同于权限或 reaction `disabled`。
- v2.7/v2.8 仅作为兼容历史；GOAL-004 的 95 个 Host/App 候选与 v2.8 处置不替代当前页面、控件和运行时核验。

## 4. 当前本仓消费入口

| 入口 | 已核事实 |
|---|---|
| `provenance-v2.9.json` | 30 artifacts；除来源路径口径问题外，本地文件与所对应上游 tag 工件 digest 一致 |
| `stage3-fixtures.test.ts` | 固定 v2.9.0/full commit，逐项校验 provenance digest |
| `app-manifest.ts` | 支持 2.7/2.8/2.9；当前适配器常量为 2.9 |
| `generate-claim.mjs` / public claim | 最终 claim 为 2.9.0；支持 page 2.7/2.9、Manifest 2.7/2.8/2.9 |
| legacy provenance | v2.7/v2.8 仍存在；兼容回归与 current pin 的权威边界待 S2 定性 |

## 5. S1 冲突登记

1. C-001：fixture-suite schema 来源路径与 provenance 路径不一致；note 同时写“docs/schemas 11 件 + 20 suites（scenarios 未 vendor）”，实际为上游机器工件 10+1、vendor fixtures 19。
2. C-002：claim generator 的 `report.pinnedUpstream.artifactVersion` 为 2.8.0，最终 claim 为 2.9.0。
3. C-003：旧 v2.7 provenance/test pin 与 v2.9 stage3 pin 并存，须区分 compatibility 与 current evidence。
4. C-004：API baseline/fragments 的 Manifest protocol 仍为 2.7，Web 适配器/claim 当前支持 2.9，且 35 页中有 8 页声明 2.9。

以上均为 `collecting`，S1 不提前分类。
