---
id: A-001-w1-plan-freeze-design
doc: audit-entry
goal: GOAL-002-w1-examples-optional-module
source: self
status: open
created: 2026-08-11
updated: 2026-08-11
version: 1.0.0
---

# A-001 · W1 方案冻结设计审计（自审 · 实施前）

## 头字段

- **source**：self
- **auditor**：Claude Code（编排器自审）
- **类型 / scope**：design-plan / W1 方案冻结 D-002 的完备性与可行性（实施前）
- **verdict**：**conditional**（方案方向正确、信息门禁已闭合；但 homePageRef 推导机制存在一个需在拆分前解决的设计冲突）

## 范围与区间

审计 `GOAL-002` 的 `D-002` 方案冻结：模块 id、homePageRef 策略、Profile 默认、拆分目标形态、测试分母、go 消费影响。不审实施（尚未开始）、不审最终成功标准达成度（0/6）。

## 成果（有证据）

- 三项方案参数已由用户 2026-08-11 书面确认并落盘：`dev.examples`、homePageRef=首个启用的 admin 功能页、mvp/admin 默认关闭（`01-decision/D-002-w1-plan-freeze.md`；I-001～I-003 → verified）。
- 缺口盘点 G1–G4 有代码证据：`kernel/profile.go` profileDefaults（mvp/admin 均含 `core.schema-render`）、`composition.go:179`（无条件 `schemarendermodule.New()`）、`manifest/app-manifest.json`（8 范例 pageId + Examples 导航 + `homePageRef: overview`）、`kernel/module.go` BuiltinModules（`admin.*`/`core.manifest-route` DependsOn `core.schema-render`）。G5 已确认：Web Shell 运行时无硬编码业务路由（`apps/web/src/app/App.tsx` 动态读 `manifest.app.homePageRef`）。
- 本自审进一步核实两处实现事实：
  - `manifest.Aggregate`（`manifest.go:120-125`）对 `app` 块按 **canonical 全等**判冲突 —— 任何含不同 `homePageRef` 的 fragment 都会触发「app identity conflict」。
  - web consumer `app-manifest.ts:479-486` 要求 **pages 非空必须声明 homePageRef** 且不得指向未启用页（`MANIFEST_HOME_PAGE_UNKNOWN`）。

## 对照成功标准（design 层面）

| 标准 | 状态 | 证据 |
|------|------|------|
| S1 能力拆分 | 设计满足 | dev.examples 持有 8 pageId + Examples nav；core.schema-render 保留 schema/validation |
| S2 可选模块形态 | 设计满足 | 命名冻结 `dev.examples`；fragment + schema 归属在方案中 |
| S3 依赖剪枝 | 设计满足（细节待拆分步） | D-002 点 4：admin.*/manifest-route 不 DependsOn dev.examples |
| S4 Profile 默认 | 设计满足 | 默认关闭已冻结 |
| S5 产品面与 home | **设计缺口** | homePageRef 推导/覆写机制与现 manifest 聚合模型冲突（见 F-001） |
| S6 回归与 go 接口 | 部分 | 测试分母已识别未枚举完整（F-003）；go 暂挂路径已预告（F-005） |

## Findings

### F-001 · homePageRef 推导机制与 manifest 聚合冲突模型不兼容（med · required）

**描述**：D-002 点 2 提出「`dev.examples` 启用时可经 fragment 覆写回 `overview`」。但现行 `manifest.Aggregate` 将 `app` 块视为首个 fragment（`core.manifest-route` baseline）的 **canonical 全等**对象，任何 `dev.examples` fragment 携带不同 `homePageRef` 都会触发 `app identity conflict`（`manifest.go:123-124`）。因此**无法**用 fragment 覆写 homePageRef。同时 web consumer 强制 pages 非空时必须声明 `homePageRef`（`app-manifest.ts:479-486`），且拆分后 `users` 等 admin 页仍在 pages 中，homePageRef 不能缺省。

**证据**：`manifest.go:120-125`（app canonical 冲突）、`app-manifest.ts:479-486`（homePageRef 必填）、`manifest/app-manifest.json:10`（`homePageRef: overview` 现为静态 baseline）。

**影响**：不解决则 S5 无法在装配层实现确定性 homePageRef。

**建议实现路径**（待拆分步落地）：把 homePageRef 推导上移到**装配层**（composition 或 manifest 聚合函数签名），按启用集计算——`dev.examples` 启用 → `overview`；否则 → 启用集中的首个 admin 功能页（按冻结声明序）；无任何 admin 功能页 → 首个启用页，仍无 → 报配置错（满足 web 必填约束）。baseline `app-manifest.json` 不再静态携带示例专用 `homePageRef: overview`。

### F-002 · `core.manifest-route` 对 `core.schema-render` 的 DependsOn 去留未定（med · recommended）

**描述**：D-002 只冻结了「admin.* 与 core.manifest-route 不 DependsOn `dev.examples`」，未定 `core.manifest-route` 是否仍依赖 `core.schema-render`。现行 `module.go:96` 声明 `core.manifest-route` DependsOn `["core.server-registration", "core.schema-render"]`。若保留，custom profile 启用 manifest-route 而无 schema-render 会解析失败；若移除，需确认 manifest 装配不消费 schema 能力。

**证据**：`kernel/module.go:96`。

**影响**：影响 custom profile 组合自由度与模块矩阵语义（go 影响面）。

**建议**：拆分步核验 manifest 装配是否真正依赖 schema-render 能力；倾向移除以扩大 custom 组合自由度，但需一并更新受影响的 `BuiltinModules` 描述与 profile 测试分母。

### F-003 · 测试分母未枚举完整（med · recommended）

**描述**：已识别 3 处既有断言假设「范例页在默认集 / homePageRef=overview」，拆分后需更新或重定向：
- `modules/schemarender/provider_test.go:23`：断言 `core.schema-render` provider 发布 8 个范例 pageId → 页面归属将移到 `dev.examples`。
- `handler/schema_test.go:31,45`：断言 `/api/schema/overview` 可解析、meta.pageId=overview → 需按 dev.examples 是否在测试启用集调整。
- `composition/s2_access_drill_test.go:100`：断言 Manifest `homePageRef: overview` → 需按新推导规则更新。

**证据**：上述三处 `:行`。

**影响**：分母漏改会导致回归误判或拆分后测试失败。

**建议**：拆分步开始时先完整枚举 composition / manifest / profile / web 代表路径的测试分母，形成显式清单后一并更新，避免拆分到一半再修。

### F-004 · web i18n 死 key 清理为 deferred 非阻断（low · recommended）

**描述**：`en-US.json` / `zh-CN.json` 含 `manifest.title.overview`、`manifest.nav.examples`、`schema.overview.*` 等范例专用 key；`dev.examples` 默认关闭后为死 key。I-004 已标 deferred non-blocking（owner=本波，复核=S6）。

**证据**：`apps/web/src/i18n/messages/en-US.json:7,19,21,173` 等。

**建议**：验收时清理或按 I-004 residual 留痕；不阻断功能。

### F-005 · VP-008 go 消费暂挂留痕时机（low · recommended）

**描述**：D-002 点 6 已「预告」变更 Profile 默认集/模块矩阵/Manifest 语义将触发 VP-008 §go 消费暂挂/重验证，但尚未正式留痕暂挂。矩阵语义实际变化发生在拆分实施落地时，而非方案冻结时。

**证据**：`D-002-w1-plan-freeze.md` 点 6；`VP-010` §与 VP-008 的接口。

**建议**：拆分实施落盘（组合根/profile 变更提交）时，同步在 `03-audit` 或 Root 台账正式记录 `go` 消费暂挂 + 恢复证据要求，避免漂移窗口。

## 必改项汇总（required）

- **F-001**（med required）：拆分步前确定 homePageRef 装配层推导机制，替换「fragment 覆写」假设；满足 web 必填约束。

## 结论 + 建议下一步

方案冻结方向正确、参数已钉死、信息门禁闭合；可进入拆分实施，但 **F-001 须在编码前落定实现路径**（装配层推导 homePageRef），F-002/F-003 同步在拆分步处理。本自审为 `self`；独立交叉审计（auditor: grok-build@grok-4.5）将另行落盘 A-002，编排器随后合并响应。
