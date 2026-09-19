---
id: E-001-r2-establishment
doc: execution-entry
parent: GOAL-003-r2-generic-job-read-surface
status: recorded
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# E-001 · R2 立项与模块接线基线侦察

## 事实（2026-09-19）

### 1. 立项

按 P-001，Root `GOAL-001` 的纲领路线图已就位且 R1 已关门，故在 R2 阶段创建子目标 `GOAL-003-r2-generic-job-read-surface`（五件套 + 三个 ledger 目录齐全），承接 Root 的 R2 检查点与 R1 `D-001` §5 的 T-1～T-8。

`GOAL-003` 的 4 个显式检查点：

| 检查点 | 内容 | 承接 |
|--------|------|------|
| C1 | 模块与权限接线 | T-1、T-2 |
| C2 | 查询与索引 | T-3、T-4 / O-3（`I-038-007`） |
| C3 | 读面 API 与作用域 | T-5（读面部分）、`I-038-008` |
| C4 | R2 审计与投影 | `I-038-009` 方案冻结 |

审计模式在实施前按 P-002 判定为 **`cross`**（跨 actor 管理读面 + 可能的迁移索引 = 权限/数据边界 + 迁移高影响门禁）。

### 2. 信息项登记

| ID | 级别 | 问题 | 最晚阶段 |
|----|------|------|---------|
| `I-038-007` | required | 管理列表索引决策（O-3）：既有三索引是否足够；是否新增迁移版本；新增则两处冻结断言如何同步 | C2 前 |
| `I-038-008` | required | 结果 URL 泛化口径（现行 `walletJobToMap` 硬编码 wallet 权限门下的路径） | C3 前 |
| `I-038-009` | required | 前端触发机制与 capability 声明口径（O-1 / O-2） | R2 方案冻结前 |
| `I-038-010` | non-blocking | 导航分组与 i18n 键位 | R4 前 |

### 3. 模块接线基线（编排器只读核对）

以 `admin.scheduled-tasks` 为最近参考，R2 的 `admin.jobs` 模块需要以下接线面（**均已在代码中核实存在**）：

| # | 面 | 参考位置 |
|---|----|---------|
| 1 | `ModuleID` 常量 | `modules/scheduledtasks/provider.go:22`（同值另见 `schema/schema.go:9`、`migration/migration.go:13`） |
| 2 | `Descriptor()`：ID / Version / KernelAPIRange / DependsOn / `Requires: kernel.StandardAdminCapabilities()` / `Contributions{Routes,Pages,Navigation,Permissions,Fragments}` | `provider.go:46-69` |
| 3 | `Register()`：`reg.HTTP` / `reg.Schema` / `reg.Authorization` / `reg.Navigation` / `reg.Manifest` | `provider.go:75-122` |
| 4 | 路由贡献 + `a.Middleware` | `internal/handler/scheduledtasks.go:468-498`；wallet 的 `add()` 工厂 `handler/wallet.go:103-110` |
| 5 | 权限贡献（`PolicyID: authsessiondata.PolicyAdmin`、`SystemDataVersion`） | `provider.go:94-101` |
| 6 | 页面 schema 贡献 + `//go:embed` 的 schema 包 | `provider.go:81-93`；`schema/schema.go`（`PageIDs()` / `SchemaDocuments()`） |
| 7 | 导航贡献（`Group: &kernel.NavigationGroup{...}`、`Visibility`、`Permission`） | `provider.go:102-114` |
| 8 | Manifest fragment（`ProtocolVersion "2.7"`、`RequiredCapabilities`、`pages[]`、`navigation.sidebar[]` 带 `visibleWhen`） | `provider.go:115-121`；`manifest/manifest.go`；`manifest/fragment.json` |
| 9 | 进 admin 默认集：`profileDefaults[ProfileAdmin]` + `BuiltinModules` 描述符 | `kernel/profile.go:46-93`、`:193`（wallet 在 `:202`） |
| 10 | 组合根按 plan 装配 `if plan.HasModule("admin.jobs")` | `internal/composition/composition.go:555-557`（wallet 先例 `:582-600`） |

**会因新模块而失败的计数型断言**（R1 侦察已登记，R2 须同步）：

| 断言 | 位置 | 现值 |
|------|------|------|
| admin profile 权限/导航计数 | `internal/composition/composition_test.go:529` | `wantPermissions: 34, wantNavigation: 18` |
| provider 贡献键与实际 `reg.*` 逐键一致 | `kernel/provider.go:178`（`stringSetEqual`） | — |

### 4. R2 必须显式处理的已知缺口

**R-1（来自 R1 矩阵 §6）**：`jobRuntime.enabled` 仅在 `plan.HasModule("admin.wallet")` 时置 true（`composition.go:588`），`Start()`/`Stop()` 在 `enabled=false` 时为 no-op（`:185-197`）。⇒ 含 `admin.jobs` 但不含 `admin.wallet` 的 Profile 下 Job runner 不会启动。admin 默认集含 wallet 故当前无矛盾，但 R2 须显式处理（改为按 `admin.jobs` 或 `admin.wallet` 任一存在即启用）。

### 5. 本轮**未**做的事（边界）

- **未**改动 `apps/**`：本轮仅治理文档写入。
- **未**冻结任何 R2 方案项（索引 / 结果 URL / 前端触发 / capability 口径）：待侦察证据齐备后落盘，涉及方案选型的按 P-004 询问用户。
- **未**执行审计（C4 待 C1～C3 完成后进行）。
- **未**创建 R3～R5 子目标。

### 6. Git checkpoint

见提交记录（R2 立项，覆盖路径 `docs/workspaces/workspace-038-batch-operations-and-job-center/`）。
