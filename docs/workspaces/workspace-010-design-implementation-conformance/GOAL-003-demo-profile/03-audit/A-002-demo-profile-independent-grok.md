---
id: A-002-demo-profile-independent-grok
doc: audit-entry
goal: GOAL-003-demo-profile
source: independent
status: closed
created: 2026-08-11
updated: 2026-08-11
version: 1.0.0
---

# A-002 · W2 demo Profile 独立波次审计（grok-build · 关门准备）

## 头字段

- **source**：independent
- **auditor**：grok-build@grok-4.5（reasoning-effort high；只读，未改任何文件/status）
- **类型 / scope**：execution-facts / W2 `demo` Profile 实施产物 vs S1–S6；实现正确性；VP-008 `go` 判定
- **verdict**：**conditional**（实施主体 S1–S4/S6 成立；**S5 的 QUICKSTART 缺口 → 1 条 required**，闭合后放行）

## 成果（有证据）

| 项 | 结论 | 证据 |
|----|------|------|
| S1 编译 Profile | 达成 | `profile.go` `ProfileDemo` + `profileDefaults[demo]`（= mvp 集 + dev.examples，无多无漏）；`ResolveProfile("demo",nil)` 成功；Source=profile.default；Precedence 全局一致 |
| S2 产品面 | 达成（主路径充分） | `deriveHomePageRef` dev.examples→overview；`TestDemoProfileManifest`（home=overview、users/roles/范例在集、无 settings/activity、schema 200）；e2e shell + schema-crud demo 分支 |
| S3 卫生保持 | 达成 | `TestDemoProfileIsNonProduction` + `TestDemoProfileManifest` 末尾；W1 `TestManifestHomePageRefDerivation` 仍绿；custom 语义未改 |
| S4 回归/烟测 | 达成（台账 + 本轮 API 子集） | E-001 记 go 23 包/web 746/三 Profile e2e；白名单放开 demo 未见泄漏到 mvp/admin |
| S5 文档 | **未完全达成** | 三处 README OK；**`QUICKSTART.md` 仍排除 demo**（F-001 required） |
| S6 go 判定 | 达成（有条件同意） | 「生产候选矩阵未变、demo 非生产向 → 不暂挂」在技术口径成立；留痕投影可加强（F-003） |

## Findings

### F-001 · S5 声称完成但 `QUICKSTART.md` 仍排除 `demo`（med · required）
**描述**：`00-meta` S5 写「README/QUICKSTART 标注 demo」并勾选完成；但 `QUICKSTART.md:33` 仍写「只接受 mvp、admin、custom」，L26/L45 示例无 demo。开发者按 QUICKSTART 会认为 `demo` 非法，与代码事实冲突 → S5 名不副实，阻断无条件关门。
**证据**：`GOAL-003/00-meta.md` S5；`QUICKSTART.md:26-33,45`；对比 `apps/api/README.md:50,65-69`、`README.md:83-97`。
**处置**：**fixed**（本波更新 QUICKSTART 接受值 + 非生产说明 + 示例）。

### F-002 · `apps/api/.env.example` 与部分架构叙述仍三 Profile（low · recommended）
**描述**：`.env.example:16` 注释 `mvp | admin | custom` 无 demo；`docs/architecture/module-architecture.md:72`、`VP-010` Profile 行等仍三元叙述。运行时不受影响，加剧文档—代码分叉。
**证据**：`.env.example:16`；`module-architecture.md:72`；`VP-010` Profile 行。
**处置**：**fixed**（`.env.example`）；architecture/VP 随符合性波次或 editorial 回贴（记录为可接受范围）。

### F-003 · `go`「不暂挂」判定合理，但对照 VP-008 字面触发留痕偏薄（low · recommended）
**描述**：E-001 以「生产默认未变 + demo 非生产」判定不暂挂，技术上成立；VP-008 字面含「Profile 默认集」变更触发。Root 波次表应补「W2 无影响/不暂挂」指针，避免 freshness 漏读。
**证据**：`E-001` §go；`VP-008` §go；`workspace.md`。
**处置**：**fixed**（Root 波次台账 W2 行已补「无影响、不触发暂挂；生产矩阵以 W1 恢复 digest 为准」）。

### F-004 · `TestDemoProfileManifest` / Source 断言充分度有界（low · recommended）
**描述**：(1) demo manifest 只抽检 3 个范例页，未断言 8 页全量；(2) `TestDemoProfileIsNonProduction` 断言 Source 未断言 Precedence；(3) 双 Profile 生命周期测试仍仅 mvp/admin。
**证据**：`composition_test.go`；`kernel_test.go`。
**处置**：**fixed**（`TestDemoProfileManifest` 补 8 范例 pageId 全量断言）；Precedence 全局一致、demo 生命周期由 e2e 间接覆盖，记录为可接受。

### F-005 · demo 下 localization e2e 全 skip（low · recommended）
**描述**：`localization.spec.ts` 分别 skip 非 admin/非 mvp；demo 运行 = 2 skip。演示面由 shell + schema-crud 覆盖。
**证据**：`localization.spec.ts:16,70`；E-001 demo 2p/2s。
**处置**：**accepted-residual**（I-002 烟测意图为演示面，shell+schema-crud 已覆盖；demo 专属本地化非本波范围）。

## 必改项汇总（required）

1. **F-001**：修正 `QUICKSTART.md` 对 `APP_PROFILE`/`demo` 的说明。**（已 fixed）**

## 与 self 审计（A-001）的异同

- self A-001 = **pass**（仅三 README 即认为 S5 达成）；本独立审 = **conditional**，指出 `QUICKSTART.md` 在 S5 标题内且错误 → required。
- S1–S4、S6 技术判断两审一致；S6「不暂挂」一致。
- 独立审新增 F-002（.env.example）、F-003（go 留痕投影）、F-004（断言充分度）。

## 结论 + 建议给用户 / 编排器的下一步

W2 `demo` Profile **实现正确、卫生未回归、e2e 分支合理**；`go` **不必暂挂**的判定成立。关门前须闭合 **F-001（required，已 fixed）**；F-002～F-005 建议处置或书面 residual。建议编排器：代贴本意见 → 合并响应（F-001/F-002/F-004 fixed、F-003 fixed、F-005 residual）→ required=0 → **GOAL-003 关门**（Root/VP-010 保持 active）。

## 声明

本意见不修改 status/progress/方案正文/goal-tree；响应、finding 闭合与阶段推进归 `/govern`；落盘由编排器代贴并保留 `source: independent`。
