---
id: E-001
doc: execution
title: S2 · 固定 UI 双语化 + Manifest/Schema key 真解析 + M4 缺失 key 流程
status: recorded
parent: GOAL-003-s2-ui-schema-bilingual
created: 2026-08-09
updated: 2026-08-09
version: 0.1.0
---

# E-001 · S2 前端与 Schema 覆盖（2026-08-09）

## 事实

- **Catalog 扩展**：`messages/en-US.json` + `zh-CN.json` 从 5 键扩展到 ~200 键（manifest.*、login.*、shell.*、manifestFailure.*、feedback.*、schema.*）；两 catalog 键集一致（新增对称测试断言）；en-US 值 = 现状英文文本（既有断言零改动通过）。
- **固定 UI（C1）**：`LoginPage`（含错误码→key 映射）、`App`（Shell 顶栏/侧栏/用户导航/反馈面/路由回退/页面 Schema 错误/移动抽屉）、`ManifestFailure`、`confirm`（Cancel/Confirm）、`render.tsx`（成功反馈/Submitting/Search/Submit/RecordView/statCard/chart/Tabs/Text 面）、`schema-table`（选中 aria/空态/分页/表头）、`data-table`（Loading/空态/移动卡片 aria）全部迁入 catalog。
- **Manifest 真解析（C2）**：`navigation.ts` labelFor 经 catalog 解析 `labelKey`/`titleKey`（`page.titleKey → page.title → pageRef`），`App.tsx` pageTitle 用 `resolveTextProp(titleKey, title)`；API 侧 manifest 数据（`app-manifest.json` + users/roles/settings/activity 4 个 fragment）与 web 测试夹具补 `titleKey`/`labelKey`（协议已声明字段，additive）。
- **Schema 真解析（C3）**：12 个 page schema 文档补 134 处 `*Key`（labelKey/textKey/submitLabelKey/confirmKey/placeholderKey）；Renderer 解析链补齐：`render.ts` `parseRenderNode`/`gateRenderFormFields` 透传 `submitLabelKey`/`labelKey`/`placeholderKey`/`options.labelKey`（原实现会丢弃——已修复），`form-controls`/`schema-table`/`render.tsx` 用 `resolveTextProp` 解析（key → 字面 → 默认）。
- **关键边界修正（S2 审计核对）**：`docs/schemas/component-registry.json` 是上游 pin 制品（I-PROTO-004 sha256 校验），**未改写**（曾尝试 additive 扩展后 revert）；其已声明的 `labelKey`/`titleKey`/`contentKey`/`options.labelKey` 即上游 key 字段；四个缺口字段（`submitLabelKey`/`confirmKey`/`textKey`/`placeholderKey`）登记为**本地页面文档约定**（上游 `props` 开放、文档级合法），见 D-001 修订与 Root D-002 修正。
- **M4（C4）**：`ui-bilingual.test.tsx` 制造缺失 key（column labelKey 指向不存在键）→ `schema-ui:missing-translation` 事件（locale/key）+ 字面文本回退 + 主流程（列表/工具栏）可完成。
- **证据矩阵**：F-V029 表 U1～U7 与页面行证据路径回填（见附件表）。
- **验证**：vitest **686/686**（36 文件）全绿，新增 S2 测试 12（双语 chrome、manifest key 解析与回退、schema labelKey/submitLabelKey/confirmKey 双语、M4 缺失 key、catalog 键集对称）；`go test ./...` 全绿（schema/manifest 数据变更兼容）；`npm run build` exit 0；输出捕获 `{SCRATCH}/unit-s2-web.log`。

## 产物

| 路径 | 说明 |
|------|------|
| `src/i18n/messages/*.json`（~200 键） | 双语 catalog（C1/C2/C3 键空间） |
| `src/i18n/catalog.ts` | `resolveTextProp`（key→字面→默认）、`translate` literalFallback 参数 |
| `src/i18n/runtime.tsx` | `useTranslate` 容忍钩子（无 Provider 时 en-US 安全默认） |
| `src/app/*`、`src/components/*`、`src/renderer/*` | 固定 UI 与 schema 文本解析（C1/C3） |
| `apps/api/internal/**/schema/*.json`（12）+ `manifest/*.json`（5） | 数据侧 `*Key`（C2/C3） |
| `src/i18n/ui-bilingual.test.tsx` | C1–C4 双语集成测试（12 项） |

## 里程碑 checkpoint

- commit：`5ea6247`（2026-08-09，S2；owned paths = 上述全部路径 + 本目标治理文档，显式 `git add` 无 `-A`）。
