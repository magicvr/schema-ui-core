---
id: A-001-demo-profile-wave
doc: audit-entry
goal: GOAL-003-demo-profile
source: self
status: closed
created: 2026-08-11
updated: 2026-08-11
version: 1.0.0
---

# A-001 · W2 demo Profile 实施波次审计（self · 关门准备）

> **闭合**：本意见无 required；recommended F-001（localization skip）经 A-003 以 residual 处置，F-002（go 判定依赖）经 A-003 §go 记录；A-002 的 required F-001（QUICKSTART）已 fixed。

## 头字段

- **source**：self
- **auditor**：Claude Code（编排器自审）
- **类型 / scope**：execution-facts / W2 `demo` Profile 实施产物 + 关门就绪
- **verdict**：**pass**（无 required；实施与用户确认方案一致，S1–S6 有证据）

## 范围与区间

审计 `E-001` 实施产物：`ProfileDemo` + `profileDefaults[demo]`、测试分母（kernel/composition/config）、playwright 白名单 + e2e 分支、README 文档、`go` 判定；对照成功标准 S1–S6。

## 成果（有证据）

- **S1 编译 Profile**：`kernel/profile.go` 新增 `ProfileDemo = "demo"` + `profileDefaults[demo]` = mvp 集 + `dev.examples`；`TestBuiltinProfilesResolveDeterministically` 断言 demo 模块列表与 Resolve 成功；`TestDemoProfileIsNonProduction` 断言 mvp/admin 不含 dev.examples、demo Source=profile.default。
- **S2 产品面**：`TestDemoProfileManifest` 断言 demo manifest `homePageRef=overview`、含 users/roles/overview/data-table/form-controls、无 settings/activity、`/api/schema/overview` 200；demo e2e（shell + schema-crud）浏览器验证 home→`/overview` + Examples 导航 + users CRUD。
- **S3 卫生保持**：`TestDemoProfileIsNonProduction` + `TestDemoProfileManifest` 末尾断言 mvp/admin 无 dev.examples；W1 的 `TestManifestHomePageRefDerivation` 保持绿。
- **S4 回归/烟测**：`go test ./...` 23 包全绿；web 746 测试；Playwright 三 Profile（mvp 3+1skip / admin 3+1skip / demo 2+2skip）全绿。
- **S5 文档**：apps/api/README、根 README、apps/web/README 补 `demo`（非生产向、用途、启用方式）。
- **S6 go 判定**：`E-001` §go 记录——mvp/admin 生产默认未变、demo 非生产向 → **`go` 保持有效、不触发暂挂**；业务 VP 以 demo 为候选时触发 freshness。

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| S1 编译 Profile | 达成 | kernel_test demo 断言 |
| S2 产品面 | 达成 | TestDemoProfileManifest + demo e2e |
| S3 卫生保持 | 达成 | TestDemoProfileIsNonProduction + W1 断言绿 |
| S4 回归/烟测 | 达成 | go/web 全绿 + 三 Profile e2e |
| S5 文档 | 达成 | 三处 README |
| S6 go 接口 | 达成 | E-001 §go 判定留痕 |

## Findings

### F-001 · demo 下 localization e2e 全 skip（low · recommended）
**描述**：`localization.spec.ts` 两测分别限 admin/mvp，demo 下 skip。demo 的演示面烟测由 shell + schema-crud 覆盖；本地化边界证据仍由 mvp/admin 承担。可接受，但若后续要验证「demo 下本地化」，需显式扩展。
**证据**：`localization.spec.ts` skip 条件；demo e2e 2 passed/2 skipped。
**建议**：非阻断；按需扩展。

### F-002 · go 判定依赖「demo 非生产」约定（low · recommended）
**描述**：S6 判定「不触发暂挂」基于 mvp/admin 生产默认未变 + demo 非生产向。若未来某业务 VP 把 demo 作为候选，或 demo 被意外纳入生产包装，判定失效。
**证据**：`E-001` §go；`TestDemoProfileIsNonProduction`。
**建议**：非阻断；已留痕触发条件（业务 VP 以 demo 为候选 → freshness），且 demo 未进入 compose/生产默认。

## 必改项汇总（required）

- **无。**

## 结论 + 建议下一步

实施与用户确认方案一致，S1–S6 达成，回归证据完整（go/web 全绿 + 三 Profile e2e），**无开放 required**。建议：由 grok-build@grok-4.5 做 independent 波次审计（A-002）交叉核对，合并响应后进入 W2 关门。
