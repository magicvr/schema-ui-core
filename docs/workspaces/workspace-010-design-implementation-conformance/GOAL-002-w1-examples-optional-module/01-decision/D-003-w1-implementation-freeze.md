---
id: D-003-w1-implementation-freeze
doc: decision-entry
goal: GOAL-002-w1-examples-optional-module
status: accepted
created: 2026-08-11
updated: 2026-08-11
version: 1.0.0
---

# D-003 · W1 实施冻结附录（home 机制 / 模块契约 / go 暂挂 / 测试分母）

## 决定

1. **homePageRef 机制（R1 · 方案 A，用户确认）**：装配层统一推导并注入。
   - 各 fragment（baseline + admin.* + `dev.examples`）的 `app` 块**移除 `homePageRef`**，保持 `manifest.Aggregate` 对 `app` 的 canonical 全等通过（消除「派生 home」与全等冲突的硬冲突）。
   - 装配层（manifest 聚合/组合根）在发布前按启用集计算 `homePageRef` 并统一写入发布版 `app` 块；web consumer 对非空 pages 的必填校验（`app-manifest.ts:479-486`）保持满足。

2. **home 推导算法表（R2 · 用户确认）**：

   | 条件（按序判定） | `homePageRef` |
   |------|------|
   | `dev.examples` 启用 | `overview` |
   | 否则，启用集含 admin 功能页 | 首个启用的 admin 功能页（声明序 `users → roles → settings → activity`） |
   | 否则，启用集含任意页 | 首个启用的任意页 |
   | 启用集无页贡献 | 省略 `homePageRef` |

   web 侧校验：非空 pages 必须声明 `homePageRef` 且指向已启用页（`MANIFEST_HOME_PAGE_UNKNOWN`）。

3. **`dev.examples` 模块契约（R3）**：
   - id `dev.examples` · v2.0.0 · KernelAPIRange `>=2.0 <3.0`。
   - **DependsOn**：`core.schema-render`（schema 能力）、`core.navigation-capability`（导航）。
   - **Provides**：无新核心能力（不提供 CapabilitySchema 等；依赖 schema-render）。
   - **Contributions**：Pages = 8 范例 pageId（overview/data-table/admin-list-batch/data-display/search-form-table/form-controls/form-with-reactions/form-with-upload）；Navigation = top `overview` + sidebar Examples 组；Fragments = examples。
   - **schema 归属**：8 个范例 schema 文档自 `apps/api/internal/modules/schemarender/schema/` 迁至 `apps/api/internal/modules/dev/examples/schema/`（包 `examples`，目录与 id 命名空间 `dev.` 对齐）；`core.schema-render` 保留 CapabilitySchema/CapabilityValidation，**Pages 贡献清空**。
   - **六面豁免说明**：横切演示模块，不提供业务数据面（HTTP API / Authorization / Persistence 业务面）；范例复用真实 API（如 `/api/users`），不建 demo 后端；仍须满足模块贡献契约（DependsOn 闭合、无贡献冲突、fragment 无 secret）。
   - **组合根**：`schemarendermodule.New()` 与 `dev.examples` provider 均按 `plan.HasModule` 条件装配。

4. **Profile 默认（grok F-005 并入）**：`mvp`/`admin` 默认集**继续包含** `core.schema-render`（能力壳，admin 能力闭包依赖），**仅移除 `dev.examples`**；`kernel/profile.go` profileDefaults 更新；`kernel/kernel_test.go` mvp 默认列表断言同步。

5. **VP-008 `go` 暂挂（R4）**：
   - **触发时机** = 首个改变 Profile 默认集 / 模块矩阵 / Manifest 装配语义的**代码合入**（非方案冻结时）。
   - **恢复证据最低字段**：候选身份快照（profile/module 矩阵）、digest、双 Profile 烟测（mvp/admin）、新增断言（禁用时无 Examples 组 / 无 8 范例 pageId / schema 404、`homePageRef` 正确、启用时恢复）。
   - **台账落点**：GOAL-002 `02-execution` / `03-audit` + Root 台账指针。

6. **测试分母勾选清单（grok F-006 并入）**：`schemarender/provider_test.go:23`（8 页 → `dev.examples`）、`handler/schema_test.go`（overview 端点分母）、`composition/s2_access_drill_test.go:97-101,140`（probe app 块 / home）、`kernel/kernel_test.go:21-28`（mvp 默认列表）、`manifest/manifest_test.go`（ForModules 页投影 + app 冲突 case）、web e2e `shell.spec.ts:20,44-45` / `localization.spec.ts:43,91,99` / `schema-crud.spec.ts:16`（默认 Manifest 必含 overview、home→`/overview`）。拆分步先核对此清单再动工。

7. **i18n**：I-004（范例 key 清理）deferred non-blocking 保持；验收时清理或 residual 留痕。

## 为什么

- cross 审计完成：self **A-001** + independent **A-002**（grok-build@grok-4.5）均 `conditional`、findings 收敛无冲突；required 必改项 R1–R4 按用户裁决走 **fixed**（本附录落盘）。
- 用户 2026-08-11 确认：home 机制 **A（组合根统一 stamp）**；`dev.examples` 启用时 **overview 优先**；无 admin 功能页时**回退任意首个页**。
- 机制 A 不动 `Aggregate` 冲突语义，是 go 影响面最小的路径（grok 推荐）。

## 未选方案

| 方案 | 未选原因 |
|------|----------|
| 机制 B：扩展 `Aggregate` 允许 home 合并 | 更大的 Manifest 装配语义变更，go 消费影响面与测试面显著增大 |
| fragment 覆写 home | 与 `Aggregate` app canonical 全等冲突模型不相容（审计证实，`manifest.go:116-125`） |
| 默认移除整个 `core.schema-render` | 破坏 `StandardAdminCapabilities` 的 schema 能力闭包，admin 解析失败 |
| 无 admin 页时 fail closed 报错 | 降低 custom profile 组合自由度；回退任意首页 + 无页省略已满足 web 契约 |

## 影响与后续

- 信息项：I-001～I-003 已 verified；I-004 deferred；无新 open required。
- 审计闭合：R1–R4（required）+ F-005/F-006/F-007/F-008（recommended）→ **fixed**，见 A-003 响应记录。
- 下一步（roadmap 阶段 2）：**拆分与迁移实施**——先核对测试分母勾选清单，再动 kernel / composition / manifest / schema 归属 / web 代表路径。
- 审计模式：实施为 high-impact（模块矩阵 / Manifest 装配语义）→ `cross` 已定；实施后按风险再跑 self + independent。
