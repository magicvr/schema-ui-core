---
id: A-005-w1-implementation-independent-grok
doc: audit-entry
goal: GOAL-002-w1-examples-optional-module
source: independent
status: closed
created: 2026-08-11
updated: 2026-08-11
version: 1.0.0
---

# A-005 · W1 实施波次独立审计（grok-build · 关门准备）

## 头字段

- **source**：independent
- **auditor**：grok-build@grok-4.5（reasoning-effort high；只读，未改任何文件/status）
- **类型 / scope**：execution-facts / W1 实施产物 vs D-003 §1–§6 / S1–S6；实现正确性；VP-008 `go` 恢复证据（**不重审**方案冻结，A-001/A-002/A-003 已闭环）
- **verdict**：**pass**（实施与 D-003 一致，S1–S6 有可核对证据；无阻断 bug；仅 recommended 补强项）

## 对照 D-003 §1–§6 与 S1–S6

| 项 | 结论 | 证据 |
|----|------|------|
| §1 home 机制 A | 一致 | fragment `app` 无 homePageRef；`StampHomePageRef`（`manifest.go:171-209`；`composition.go` 调用） |
| §2 推导表 | 实现正确 | `deriveHomePageRef`：dev.examples→overview；admin 序 users→roles→settings→activity；任意有页模块首页；无页 "" |
| §3 模块契约 | 一致 | dev.examples v2.0.0 / DependsOn schema-render+navigation / 8 页+examples fragment / 无 Provides / 无 Routes·Permissions·system-nav |
| §4 Profile 默认 | 一致 | `profileDefaults` 含 schema-render、不含 dev.examples（`kernel_test.go`） |
| §5 go 暂挂 | 触发已留痕；恢复证据技术面基本齐 | E-004 §go；digest `4a2b8cd…`（原文补强见 F-004） |
| §6 测试分母 | 主清单已改 | 见 F-001～F-003 |
| S1–S6 | 达成 | schema-render 0 页；dev/examples 8 页；mvp/admin home=users、无泄漏、启用恢复 overview；go/web 全绿 + 双 Profile e2e 3+3 |

## 成果（有证据）

1. G1–G4 纠偏落地：范例面迁出、schema-render 能力壳、默认无演示面。
2. 装配层 home（机制 A）：`Aggregate` 全等不被 home 破坏；stamp 注入/删除。
3. fragment app 全等：baseline + 4 admin + dev.examples + S2 probe 均为 `{appId,name,description}` 无 homePageRef。
4. 条件装配自洽：schema-render / dev.examples 均 `HasModule`；DependsOn 闭合由 `Registry.Resolve` fail-closed。
5. `core.manifest-route` 保留 DependsOn schema-render：与 D-003 §4/E-004 取舍一致，非演示依赖。
6. web consumer 约束：`app-manifest.ts:479-496` 非空 pages 必声明合法 homePageRef；默认装配满足。
7. 回归主路径绿：API 包本会话复跑通过；e2e 分母 `/overview`→`/users`。

## 特别核查

- **a) deriveHomePageRef**：四档与 D-003 §2 一致；dev.examples 依赖闭合（缺 schema-render/navigation → Resolve 失败）。「任意首页」取拓扑序首个有页模块（确定可复现；生产不触及）。**边分支缺单测**（F-002）。
- **b) StampHomePageRef**：不参与 Aggregate，不破坏 app 全等；map marshal 按键排序确定；pages/navigation Raw 保留；仅 5 顶层字段（未来扩字段会丢，F-005）。无独立单测（composition 集成覆盖）。
- **c) fragment 全等 + web 约束**：成立；空 pages + 省略 home 与 `PAGE_REF_WITH_EMPTY_PAGES` 规则兼容。
- **d) 条件装配**：mvp/admin 自洽；custom 加 dev.examples 测试路径通过。
- **e) 测试分母全局**：主清单已更新；e2e 已改；web 夹具仍为「范例启用」形态（非运行时，F-006）；默认 Manifest「无 Examples 导航组」API 断言偏弱（F-003）。
- **f) VP-008 go 恢复证据**：矩阵快照充分；digest 半充分（E-004 未写死 hash，F-004）；双 Profile 烟测台账充分；新增断言主路径充分；台账落点 E-004 + Root 波次指针。结论：**支持关门后由编排器恢复 go 消费**（范围=本波后矩阵；不免除业务 VP 激活前 freshness）。

## Findings

### F-001 · schema 404/200 断言仅存在于未提交工作区（med · recommended）
**描述**：`TestManifestHomePageRefDerivation` 的 `/api/schema/overview` 404/200 断言在已合入 commit `4a2b8cd` 中**缺失**，仅在工作区未提交改动中。行为正确，但落地回归分母与 commit 不一致。
**证据**：`git show 4a2b8cd:…/composition_test.go` vs 工作区 diff。
**处置**：随本波关门提交补齐（fixed）。

### F-002 · home 推导边分支缺自动化用例（med · recommended）
**描述**：推导表第 2/3/4 档（roles-only、无 admin 任意页、无页省略）无单测；custom profile 回归靠代码阅读。
**证据**：`composition.go:256-274`；测试仅覆盖 mvp/admin→users 与 +dev.examples→overview。
**处置**：新增 `TestDeriveHomePageRefBranches` 表测（fixed）。

### F-003 · 禁用路径「无 Examples 导航组」API 断言不完整（low · recommended）
**描述**：只抽查 overview/data-table/form-controls 页级泄漏；无 sidebar label=`Examples` 组断言。
**证据**：`composition_test.go`；`shell.spec.ts` 以链接 count=0 部分补偿。
**处置**：补充 Examples 导航组断言（fixed）。

### F-004 · go 恢复 digest 未写死台账（low · recommended）
**描述**：E-004 未写死 full hash，需从 git 推断。
**证据**：E-004 §go；`git rev-parse HEAD`。
**处置**：E-004 已补写 `4a2b8cd…`（fixed）。

### F-005 · StampHomePageRef 固定 5 字段 envelope（low · recommended）
**描述**：协议扩顶层字段会静默丢字段（当前协议安全）。
**证据**：`manifest.go:181-209`。
**处置**：accepted-residual（协议 envelope 固定，低风险可逆）。

### F-006 · web 夹具仍为「范例启用 + home=overview」形态（low · recommended）
**描述**：`test-fixtures/app-manifest.{mvp,admin}.json` 与默认 Profile 现实不一致，易误导维护者（非运行时回归失败源）。
**证据**：夹具文件；多处测试引用。
**处置**：accepted-residual（Renderer 夹具，非 Profile 契约；后续可改名 dogfood 语义）。

### F-007 · I-004 i18n 范例 key 非死 key（low · recommended · 信息项）
**描述**：dogfood 路径仍被 dev.examples fragment 引用，非死 key；可按「保留」闭合 deferred。
**证据**：dev/examples fragment titleKeys；00-meta I-004。
**处置**：I-004 以「保留」闭合（fixed/resolved）。

## 必改项汇总（required）

- **无。**

## 与既有意见的异同

与 self **A-004（pass）** 同向、无冲突。独立视角强调：(1) schema-404 断言仅在未提交工作区（F-001）；(2) 推导表边分支缺测（F-002）；(3) go digest 未写死台账（F-004）。均不构成放行阻断。

## 结论 + 建议给用户 / 编排器的下一步

W1 实施与 D-003 §1–§6 对齐，S1–S6 有代码与测试证据；`deriveHomePageRef`/`StampHomePageRef`/条件装配/fragment 全等未见阻断 bug。**无 required**，verdict = **pass**。建议编排器：代贴本意见 → 合并响应（F-001～F-007 按 fixed/residual 处置）→ required=0 → **W1 关门**（GOAL-002 done、goal-tree、Root 波次档案、VP-008 `go` 消费恢复留痕）；勿在未合并 independent 意见时仅凭 A-004 关门（cross 已定）。

## 声明

本意见不修改 status/progress/方案正文/goal-tree；响应、finding 闭合与阶段推进归 `/govern`；落盘由编排器代贴并保留 `source: independent`。
