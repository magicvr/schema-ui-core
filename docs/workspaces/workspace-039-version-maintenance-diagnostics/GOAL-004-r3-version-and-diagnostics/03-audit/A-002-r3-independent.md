---
id: A-002-r3-independent
doc: audit-entry
parent: GOAL-004-r3-version-and-diagnostics
status: recorded
source: independent
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
auditor: grok-build (grok-4.6 · reasoning high)
---

# A-002 · R3 实施独立交叉审计（T-4）

| 字段 | 值 |
|------|-----|
| source | `independent` |
| auditor | grok-build（grok-4.6 · reasoning high · `/audit`） |
| date | 2026-09-19 |
| 类型 | `execution-facts` |
| scope | `[workspace-039-version-maintenance-diagnostics] GOAL-004-r3-version-and-diagnostics` · T-4 版本 chip + QUICKSTART + 诊断入口；对照 `[GOAL-002] D-001` §2/§3/§5 T-4 与本目标 D-001 |
| verdict | **pass** |
| 开放 required findings | **0**（2 recommended） |

## 范围与区间

### 工作区绑定（已校验）

| 项 | 观察 | 结论 |
|----|------|------|
| `workspace.md` | `id=workspace-039-version-maintenance-diagnostics`；`root_goal=GOAL-001-version-maintenance-diagnostics`；`canonical_scope` 仅本区；`vision_role=delivery`；`primary_plan`/`plan_refs`=`VP-039-version-maintenance-diagnostics` | 绑定合格 |
| 共享资料 | `shared_materials_catalog: none` | 无非法固定引用；未把资料当事实 |
| 目标位置 | 本区扁平 `GOAL-004-r3-version-and-diagnostics/`；`parent=GOAL-001-version-maintenance-diagnostics` | 在 canonical 内 |
| 跨区 | 本意见用 Q2 引用本区 `[GOAL-002] D-001`；未读取其他工作区目标状态 | 合格 |

`workspace.md` 纲领表仍写 R3「未开始」，与 GOAL-004 已实施事实不同步。这是相邻台账漂移，**不**否定 T-4 代码（见 F-001）。

本意见**不**修改 `status` / 检查点 / 派生 `progress` / 方案正文 / `goal-tree`。

### 本轮核验问题（用户指定）

1. 无 `monitoring.read` 不 fetch。
2. 只展示 `version`，不展示 `commit`。
3. QUICKSTART URL 为冻结 GitHub blob。
4. 诊断入口仅在有 `system-monitoring` 页面时出现。
5. mvp 默认无 chip。
6. 未改 `kernel/profile.go` 默认集、未新建模块、未改 `protocol/upstream/**`。

## 成果（有证据）

实施提交：`6bb44c24`（working tree clean）。对照 `[GOAL-002] D-001` §2 / §3 / §5 T-4 与本目标 `D-001-r3-implementation-scope`。

提交文件仅：`version-chip.tsx` / `version-chip.test.tsx` / `App.tsx` / 双 locale 三键 / GOAL-004 五件套与 `goal-tree.md`。不含 `apps/api/kernel/profile.go`、`apps/api/modules/**`、`apps/web/src/protocol/upstream/**`。

### 权限门与 fetch 短路（C1）

| 主张 | 独立核对 |
|------|----------|
| 无 `monitoring.read` 不请求 status | `version-chip.tsx` `useEffect`：`if (!canReadMonitoring) { setVersion(null); return; }`，在调用 `fetcher` 之前返回 |
| 无权限不渲染 | 同文件：`if (!canReadMonitoring \|\| version === null) return null` |
| `hasMonitoringRead` | `Array.isArray(permissions) && permissions.includes("monitoring.read")`；`undefined` / 非数组 → false |
| 生产接线 | `App.tsx`：`canReadMonitoring={hasMonitoringRead(currentUser?.permissions)}`；`AuthGate` 把 `user`（`/me` → `parseAuthUser` 保留 `permissions`）传入 `currentUser` |
| 失败隐藏、不打断 Shell | 非 `ok` / 空 version / `catch` → `setVersion(null)` → 不渲染 |
| 登录页无 chip | `AuthGate` `unauthenticated` → `LoginPage`，不挂 `App`。`auth-gate.wiring.test.tsx`：无会话时 `captured.appProps === null` |
| 单测 | `version-chip.test.tsx`：`canReadMonitoring: false` → `fetcher` 未被调用且无 `[data-version-chip]` |

fetcher 走 `resourceFetcher`（AuthGate 的 Bearer / refresh 包装）。匿名 fallback 仅在未注入 fetcher 时发生；生产路径已注入。无权限时根本不发请求，服务端 `requirePermission(..., "monitoring.read")` 仍是第二道门。

### 版本权威与不展示 commit（C2）

| 主张 | 独立核对 |
|------|----------|
| 读 `GET /api/system-monitoring/status` 信封 `items[0].version` | `STATUS_URL`；`body.items?.[0]?.version`；仅 `typeof === "string" && !== ""` 才展示 |
| 同源 | `handler/systemmonitoring.go` `Version: version.Version`（`pkg/version`）；本提交未改该 handler |
| 不展示 `commit` / `BuiltAt` / 模块列表 | 渲染只输出 `version` 文本 + 升级链 + 可选诊断按钮。JSON 里的 `commit`/`modules` 不进 DOM |
| 单测 | mock 含 `commit: "deadbeef"` 与 `modules`；`textContent` 为 `1.2.3` 且不含 `deadbeef` |

监控页保持既有 version/commit（R1 合同「监控页保持现有」）；Shell 不复制该面。符合「不展示」而非「不传输」：有 `monitoring.read` 的用户本就可以打开监控页看到 commit。

### 升级链接（C3 半边）

| 主张 | 独立核对 |
|------|----------|
| 冻结 URL | `QUICKSTART_UPGRADE_URL = "https://github.com/magicvr/schema-ui-core/blob/main/QUICKSTART.md"`，与本目标 D-001 逐字一致 |
| R1 默认 | `[GOAL-002] D-001` §2：默认 GitHub blob `main/QUICKSTART.md`；R3 未另钉 SHA（本目标 01-decision：「无新 P-004。QUICKSTART URL 按 R1 默认冻结」） |
| 目标文档存在升级节 | 仓库根 `QUICKSTART.md` 含 `schema-ui upgrade` 节（方法 B 主路径） |
| 单测 | 断言 `a[href="${QUICKSTART_UPGRADE_URL}"]` 存在 |

### 诊断入口（C3 半边）

| 主张 | 独立核对 |
|------|----------|
| 有 `system-monitoring` 页才给按钮 | `App.tsx`：`manifest.pages.some((page) => page.pageId === "system-monitoring") ? () => onNavigate("/system-monitoring") : undefined` |
| 组件侧缺省不渲染按钮 | `onOpenDiagnostics !== undefined` 才输出诊断 `<button>` |
| 不新建页 | `6bb44c24` 无新 page / schema / 模块；导航硬编码既有路由 `/system-monitoring` |
| 自动化 | **无**。`version-chip.test.tsx` 三例均不传 `onOpenDiagnostics`；`App.integration.test.tsx` 未断言该接线（见 F-002） |

代码路径可静态核对为正确；缺的是回归锁，不是行为错误。

### mvp / demo 默认无 chip

| 主张 | 独立核对 |
|------|----------|
| mvp 默认集不含 `admin.system-monitoring` | `kernel/profile.go` `ProfileMVP`：users/roles/account/dashboard/notifications + core；**无** system-monitoring。本提交未改该文件（上次提交 `ea6e4006`，jobs R4） |
| demo 同 mvp + `dev.examples` | `ProfileDemo` 同样不含 system-monitoring |
| `monitoring.read` 仅 admin 模块贡献 | `BuiltinModules` / `composition_test.go`：「S-03 … `monitoring.read` … to admin only」；mvp `wantPermissions: 11` 不含该键 |
| 产品面 | 无该权限 → `hasMonitoringRead` false → 不 fetch、不渲染。即便误授权限字符串，mvp 无 status 路由（404）与无页面（无诊断按钮），chip 仍隐藏 |

派生机制是权限 + 页面存在，不是 profile 名硬编码。与 D-001「mvp/demo 默认无此权限 → 无 chip」一致。

### 红线

| 主张 | 独立核对 |
|------|----------|
| 未改 `ResolveProfile` / 默认集 | `git diff 5929e641 6bb44c24 -- apps/api/kernel/profile.go` 空 |
| 未新建模块 | 同区间 `apps/api/modules` 空；无新 `admin.*` 目录 |
| 未改 pinned `protocol/upstream/**` | 同区间该路径空 |
| 不新建诊断页 | 无新 schema/page；复用既有 `admin.system-monitoring` |

### 本轮复跑（独立，非转述 E-002）

`apps/web` vitest（直接 `vitest run`，绕过 pnpm install hook）：

- `version-chip.test.tsx` 3
- `catalog.test.ts` 12（含 zh-CN/en-US 键集合相等；`versionChip.*` 双 locale 已齐）
- `App.integration.test.tsx` 16
- `auth-gate.wiring.test.tsx` 2
- `ui-bilingual.test.tsx` 12

**5 files / 45 passed**。

说明：后三套**未覆盖** VersionChip 接线；全绿证明未回归既有壳层，不能单独当作 C3 诊断入口的锁。

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| C1 无 `monitoring.read` 不渲染、不请求 status | **满足** | `version-chip.tsx` 短路 + 单测；`App`/`AuthGate`/`parseAuthUser` 权限链 |
| C2 从 status 读 `version`（非 commit） | **满足** | 只取 `items[0].version`；单测锁定不含 commit |
| C3 QUICKSTART blob + 诊断入口仅有页时出现 | **满足（代码）**；诊断入口自动化不足 | URL 与 D-001 逐字一致；`App.tsx` pageId 门。测试缺口见 F-002 |
| mvp 默认无 chip | **满足** | profile 默认集 + 权限贡献面 + 客户端门禁 |
| 未改 Profile / 未新模块 / 未改 pinned 协议 | **满足** | `6bb44c24` 文件列表与路径 diff |
| 影响本 scope 的 required 信息项 | **无开放** | `I-039-001` / `I-039-003` 已在 GOAL-002 R1 `verified`（最晚阶段 R1）；GOAL-004 无新 I-00N；`I-039-006` deferred non-blocking |

GOAL-004 检查点 C1–C4 **不由本意见勾选**。C4 的 independent 半边即本条；勾选与 `progress` 归 `/govern`。

## Findings

无 required。

### F-001 · 相邻台账仍写 R3 未开始（非实施缺陷）

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：T-4 代码与 GOAL-004 五件套（E-002、A-001）已记录实施事实。`workspace.md` 纲领表 R3 仍为「未开始」。不构成「版本 chip 未做」。`/govern` 响应 C4 时应刷新工作区上下文说明，避免 R4 误读。
- 证据：`docs/workspaces/workspace-039-version-maintenance-diagnostics/workspace.md` 纲领表 R3 行
- 与 A-001：self 未写台账漂移。不阻断 C4。同类：GOAL-003 A-002 F-001。

### F-002 · 诊断入口无自动化锁

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：C3「有 `system-monitoring` 页才出现诊断入口」只存在于 `App.tsx` 四行接线。`version-chip.test.tsx` 不覆盖 `onOpenDiagnostics` 有/无；`App.integration.test.tsx` 不传 `currentUser.permissions`、不断言诊断按钮。行为经代码审阅成立，但以后改 manifest 门或误传回调时测试不会红。建议补：无回调不渲染按钮；有回调点击调用 `onNavigate("/system-monitoring")`；可选 mvp manifest 无该页。
- 证据：`apps/web/src/app/version-chip.test.tsx`（3 例均无 `onOpenDiagnostics`）；`apps/web/src/app/App.tsx:1347-1354`；`App.integration.test.tsx` 无 `currentUser` / `data-version-chip`
- 不阻断 C4：生产接线可静态核对。

## 必改项汇总

无。开放 required = 0。

## 与既有意见的异同

| 来源 | 关系 |
|------|------|
| GOAL-004 A-001 self `pass` | **同意**权限门与 fetch 短路、不带 commit、诊断复用既有页、红线未破。self 证据过薄（无路径级核对）。本条补上提交文件列表、`/me` 权限链、profile 默认集、QUICKSTART 原文、vitest 复跑与诊断测试缺口。 |
| `[GOAL-002] D-001` T-4 / §2 / §3 | **实施已落地**：Shell 仅 `monitoring.read` 见 `Version` + QUICKSTART blob；不展示 Commit/模块清单；不新建诊断页。 |
| GOAL-003 A-002（T-4 未混入） | R2 当时成立；本提交 `6bb44c24` 才引入 chip，范围正确。 |
| GOAL-003 A-002 F-001 台账漂移 | 同类现象在 `workspace.md` R3 行再现 → 本条 F-001。 |

无冲突。无 P-004 冲突项。

## 结论 + 建议给编排器/用户的下一步

**verdict = pass。** T-4 与 `[GOAL-002] D-001` 及本目标 D-001 一致：无 `monitoring.read` 不 fetch；Shell 只展示 `version`；升级链为冻结 GitHub blob；诊断入口随 `system-monitoring` 页面存在而出现；mvp/demo 默认无该权限故无 chip；Profile 默认集 / 新模块 / pinned upstream 未动。

建议 `/govern`：

1. 响应 A-001 + A-002（本条 recommended F-001/F-002 可修台账说明或补测，非门禁）。
2. 若接受实施证据，勾选 GOAL-004 C1–C3；C4 在 self+independent 均 `pass` 且开放 required=0 后勾选。
3. 不要把本意见写成已改 `status`/`progress`。

## 声明

本意见不修改 status/progress；响应由 `/govern` 处理。
