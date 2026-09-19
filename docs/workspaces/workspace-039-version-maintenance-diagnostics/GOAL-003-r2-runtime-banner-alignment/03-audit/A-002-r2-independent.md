---
id: A-002-r2-independent
doc: audit-entry
parent: GOAL-003-r2-runtime-banner-alignment
status: recorded
source: independent
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
auditor: grok-build (grok-4.6 · reasoning high)
---

# A-002 · R2 实施独立交叉审计（T-1 / T-2 / T-3）

| 字段 | 值 |
|------|-----|
| source | `independent` |
| auditor | grok-build（grok-4.6 · reasoning high · `/audit`） |
| date | 2026-09-19 |
| 类型 | `execution-facts` |
| scope | `[workspace-039-version-maintenance-diagnostics] GOAL-003-r2-runtime-banner-alignment` · T-1 bootstrap 生产者、T-2 `/me`+`fetchMe`、T-3 Shell 横幅；对照 `[GOAL-002] D-001` 与两矩阵 |
| verdict | **pass** |
| 开放 required findings | **0**（1 recommended） |

## 范围与区间

### 工作区绑定（已校验）

| 项 | 观察 | 结论 |
|----|------|------|
| `workspace.md` | `id=workspace-039-version-maintenance-diagnostics`；`root_goal=GOAL-001-version-maintenance-diagnostics`；`canonical_scope` 仅本区；`vision_role=delivery`；`primary_plan`/`plan_refs`=`VP-039-version-maintenance-diagnostics` | 绑定合格 |
| 共享资料 | `shared_materials_catalog: none` | 无非法固定引用；未把资料当事实 |
| 目标位置 | 本区扁平 `GOAL-003-r2-runtime-banner-alignment/`；`parent=GOAL-001-version-maintenance-diagnostics` | 在 canonical 内 |
| 跨区 | 本意见用 Q2 引用本区 `[GOAL-002] D-001`；未把其他工作区目标状态当作本区事实 | 合格 |

`workspace.md` 纲领表仍写 R2「未开始」，与 GOAL-003 已实施事实不同步。这是相邻台账漂移，**不**否定 T-1/T-2/T-3 代码（见 F-001）。

本意见**不**修改 `status` / 检查点 / 派生 `progress` / 方案正文 / `goal-tree`。

### 本轮核验问题（用户指定）

1. maintenance Host **文档**现为 `degraded`；Host **消费者**仍对 `availability.mode=maintenance` 终态。
2. `/me` 精确 `runtimeMode`；`fetchMe` 不丢字段。
3. 横幅三分文案；`normal` 不渲染；登录页无 `App` 故无横幅。
4. T-4 版本 chip 未混入。
5. 未改 `apps/web/src/protocol/upstream/**`、`operational.go` 写门禁。

## 成果（有证据）

实施提交：`5929e641`（working tree clean）。对照 `[GOAL-002] D-001` §1 / §5 T-1～T-3 与 `attachments/r1-mode-projection-matrix.md`。

### T-1 · Host 生产者折叠

| 主张 | 独立核对 |
|------|----------|
| `maintenance`/`degraded`/`read-only` → Host `degraded` | `apps/api/internal/handler/bootstrap.go` `bootstrapAvailability`：三模式同一 `case`，返回 `{Mode: "degraded"}`；`normal`/空 → `normal`；非法值 fail-closed |
| 单测期望已改 | `bootstrap_test.go` `TestRegisterBootstrapProjectsRuntimeAvailability`：`maintenance` want `degraded`（不再 want `maintenance`） |
| Admin composition 接线 | `composition.go:807` 已有 `RegisterBootstrapWithAvailability(..., string(cfg.RuntimeMode))`（本提交未改此行；折叠发生在生产者函数） |
| `/me` 接线同步传入过程模式 | `composition.go:443` `RegisterWithMFAProbes(..., string(cfg.RuntimeMode), ...)` |
| 下游 serve 面 | `apps/api/server/serve.go` bootstrap 与 `/me` 固定 `"normal"`。该面无 `RuntimeMode`、无 operational gate，不是 Admin 产品进程（`cmd/server` → `composition.NewApp`）。E-002 已记事实，不构成本 scope 缺陷 |

### T-2 · `/me.runtimeMode` + `fetchMe`

| 主张 | 独立核对 |
|------|----------|
| 结构体 additive | `account/session.go` `RuntimeMode string \`json:"runtimeMode"\``（无 `omitempty`） |
| 处理器原样字符串；空 = `normal` | `handler/account.go` `meHandler`：空串规范化为 `normal`，否则原样写入 |
| 注册链 | `health.go` `RegisterWithMFAProbes` 接受 `runtimeMode` → `accountsHandler`；`RegisterWithMFA` 默认 `"normal"`（测试/简路径） |
| `/me` 测试 | `account_test.go`：`RegisterWithReadiness` 默认 → `normal`；`RegisterWithMFAProbes(..., "maintenance")` → `"maintenance"` |
| `fetchMe` 不丢字段 | `auth-client.ts`：`AuthSession.runtimeMode` 必填；`fetchMe` 读 `body.runtimeMode` 经 `parseRuntimeMode`；缺省/非法回落 `normal` |
| 会话入口都走 `fetchMe` | `login` / MFA verify / `restoreSession` 成功路径均 `return await fetchMe()` |
| 客户端测试 | 缺字段 → `normal`（既有 permissions 用例）；显式 `"maintenance"` 投影用例 |

### T-3 · Shell 横幅

| 主张 | 独立核对 |
|------|----------|
| 只读 `/me.runtimeMode` | `AuthGate` 传 `session?.runtimeMode`；`App` 注释禁止 Host `availability.mode`；`RuntimeBanner` 只吃该 prop |
| 三分文案且不塌成「降级」 | `runtime-banner.tsx` 三 key；zh-CN / en-US 三条文案互不相同；maintenance 提到登录/恢复仍可用，与写门禁白名单一致 |
| `normal` / 缺省不渲染 | `runtimeMode !== maintenance|degraded|read-only` → `null`；单测锁定 |
| 登录页无横幅 | `AuthGate` `unauthenticated` → `LoginPage`，不挂 `App`。`auth-gate.wiring.test.tsx`：无会话时 `captured.appProps === null` 且渲染登录表单 |
| 浅色/深色 | `text-amber-950 dark:text-amber-100` + `bg-amber-500/10` |
| 位置 | `App` 壳层顶部、sticky topbar 之前，`role="status"` |

### 红线与非目标

| 主张 | 独立核对 |
|------|----------|
| Host **消费者**仍终态 `maintenance` | `host/bootstrap.ts:275`：`mode === "maintenance"` → `MAINTENANCE`；`mode === "degraded"` → `READY_DEGRADED`（:317–318）。`boot.ts:195` 仅 `READY`/`READY_DEGRADED` 继续拉 manifest |
| pinned fixtures 未改 | `5929e641` 文件列表不含 `apps/web/src/protocol/upstream/**`。`host-bootstrap.cases.json` `maintenance-terminal` 仍期望 `result: MAINTENANCE` |
| 写门禁未改 | `5929e641` 不含 `operational.go`。`operational_test.go` 仍：maintenance→503/`SERVICE_MAINTENANCE`；degraded→`SERVICE_DEGRADED`；read-only→`SERVICE_READ_ONLY`；登录/恢复/邀请白名单 |
| T-4 版本 chip 未混入 | 提交与 `apps/web/src` 无 version chip / QUICKSTART 入口；横幅不展示 `Version`/`Commit` |
| T-5 HostFailure maintenance 路径 | 本提交未改 `HostFailureScreen` / `evaluateBootstrap` 终态机 |

### 本轮复跑（独立，非转述 E-002）

- `go test ./internal/handler ./internal/account ./internal/composition ./server -count=1` → ok（含 `operational_test` 与 composition R5 写门禁）
- vitest：`runtime-banner` / `auth-client` / `auth-gate.wiring` / `auth-context` → 4 files / 33 passed
- vitest `src/protocol`（含 `upstream-host-fixtures.test.ts` 99）→ 12 files / 558 passed

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| T-1 生产者 `maintenance`→Host `degraded` + 单测 | **满足** | `bootstrap.go` / `bootstrap_test.go`；composition 传入 `cfg.RuntimeMode` |
| T-2 `/me.runtimeMode` 原样 + `fetchMe` 投影 | **满足** | `account.go` / `session.go` / `auth-client.ts`；handler + vitest |
| T-3 三分文案、i18n、浅色深色、不挡登录页；只读 `/me` | **满足** | `runtime-banner.tsx` + 双 locale；`AuthGate`/`App` 接线；登录页 wiring 测试 |
| 消费者仍实现 `maintenance` 终态；本仓生产文档不再发出该值 | **满足** | `evaluateBootstrap` 未改；生产者不再 `case "maintenance": Mode maintenance` |
| 写门禁与 pinned upstream 不改 | **满足** | `5929e641` 文件列表；protocol fixtures 仍绿 |
| T-4 未混入 | **满足** | 无版本 chip 代码 |
| 空 runtime 配置视为 `normal` | **满足** | `bootstrapAvailability` 与 `meHandler` 均把 `""` 当 `normal`；config 默认 `RuntimeModeNormal` |
| 影响本 scope 的 required 信息项 | **无开放** | `I-039-002` 已在 GOAL-002 R1 `verified`；GOAL-003 无新 I-00N；`I-039-006` deferred non-blocking 不进本门禁 |

GOAL-003 检查点 C1–C4 **不由本意见勾选**。C4 的 independent 半边即本条；勾选与 `progress` 归 `/govern`。

## Findings

无 required。

### F-001 · 相邻台账仍写 R2 未开始（非实施缺陷）

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：T-1/T-2/T-3 代码与 GOAL-003 五件套（E-002、A-001）已记录实施事实。下列**相邻**表面仍停在立项前：`workspace.md` 纲领表 R2 =「未开始」；`goal-tree.md` 说明段仍写「R2 未立项」（树与状态表已有 GOAL-003 `active · 0/4`）。不构成「横幅未做」。`/govern` 响应 C4 时应刷新 workspace 上下文说明，避免下一波误读。
- 证据：`docs/workspaces/workspace-039-version-maintenance-diagnostics/workspace.md` 纲领表；同目录 `goal-tree.md`「说明」段 vs 树/状态表
- 与 A-001：self 未写台账漂移。不阻断 C4。

## 必改项汇总

无。开放 required = 0。

## 与既有意见的异同

| 来源 | 关系 |
|------|------|
| GOAL-003 A-001 self `pass` | **同意** T-1/T-2/T-3 与 D-001 一致、写门禁仍绿、pinned 未改、T-4 未做。本条补上消费者终态机、protocol fixtures 复跑、`fetchMe` 全会话入口、登录页 wiring 测试与提交级红线核对。 |
| GOAL-002 A-002 F-001（横幅禁止读 Host `availability.mode`） | **实施已落地**（本目标证据）：横幅只读 `session.runtimeMode`。本独立意见不关闭 GOAL-002 台账项（该目标已关门）。 |
| GOAL-002 A-002 F-002（`fetchMe` 会丢 additive 字段） | **实施已落地**：`fetchMe`/`AuthSession`/`parseRuntimeMode` 已投影。 |
| GOAL-002 A-002 F-003（相邻台账漂移） | 同类现象在本区 `workspace.md` 再现 → 本条 F-001。 |

无冲突。无 P-004 冲突项。

## 结论 + 建议给编排器/用户的下一步

**verdict = pass。** T-1/T-2/T-3 与 `[GOAL-002] D-001` 派生口径一致：生产 Host 文档对 maintenance 发 `degraded`（Shell 走 `READY_DEGRADED`），精确模式经 `/me.runtimeMode` 到已登录横幅；消费者终态机与写门禁、pinned fixtures 未改；版本 chip 未混入。

建议 `/govern`：

1. 响应 A-001 + A-002（本条 recommended F-001 可修台账说明，非门禁）。
2. 若接受实施证据，勾选 GOAL-003 C1–C3；C4 在 self+independent 均 `pass` 且开放 required=0 后勾选。
3. 不要把本意见写成已改 `status`/`progress`。

## 声明

本意见不修改 status/progress；响应由 `/govern` 处理。
