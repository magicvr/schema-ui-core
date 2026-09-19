---
id: A-002-r1-freeze-independent
doc: audit-entry
parent: GOAL-002-r1-denominator-and-contract-freeze
status: recorded
source: independent
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
auditor: grok-build (grok-4.6 · reasoning high)
---

# A-002 · R1 冻结独立交叉审计（C1～C3）

| 字段 | 值 |
|------|-----|
| source | `independent` |
| auditor | grok-build（grok-4.6 · reasoning high · `/audit`） |
| date | 2026-09-19 |
| 类型 | `design-plan` |
| scope | `[workspace-039-version-maintenance-diagnostics] GOAL-002-r1-denominator-and-contract-freeze` · R1 冻结 C1～C3 全量复审（D-001 + 两矩阵 + 侦察三份 + 对照现行代码） |
| verdict | **pass** |
| 开放 required findings | **0**（3 recommended） |

## 范围与区间

### 工作区绑定（已校验）

| 项 | 观察 | 结论 |
|----|------|------|
| `workspace.md` | `id=workspace-039-version-maintenance-diagnostics`；`root_goal=GOAL-001-version-maintenance-diagnostics`；`canonical_scope` 仅本区；`vision_role=delivery`；`primary_plan`/`plan_refs`=`VP-039-version-maintenance-diagnostics` | 绑定合格 |
| 共享资料 | `shared_materials_catalog: none` | 无非法固定引用；未把资料当事实 |
| 目标位置 | 本区扁平 `GOAL-002-r1-denominator-and-contract-freeze/`；`parent=GOAL-001-version-maintenance-diagnostics` | 在 canonical 内 |
| 跨区 | 本意见未把其他工作区目标状态当作本区事实 | 合格 |

`workspace.md` 正文仍写 GOAL-002 `0/4`、C1～C3「等 P-004」，与本目标 `00-meta`/`goal-tree` 的 `3/4` 及已冻结事实不同步。这是相邻台账漂移，**不**否定 D-001 冻结本身（见 F-003）。

### 本轮核验问题（用户指定）

1. 用户裁决 B/A/A 是否如实落盘；派生口径是否标成派生而非用户原话。
2. 矩阵 vs **现行**代码：maintenance Host 文档现为 `maintenance`——冻结是否诚实写成「将改」。
3. 改 Host **生产者**是否真能让 Shell 加载（`degraded` → `READY_DEGRADED`）。
4. 是否越界要求改 pinned upstream fixtures。
5. 横幅全登录 vs 版本仅 `monitoring.read` 是否自洽（mvp 无监控模块）。
6. 是否未改 `apps/**`。

本意见**不**修改 `status` / 检查点 / 派生 `progress` / 方案正文 / `goal-tree`。

## 成果（有证据）

| 工件 | 路径 | 独立核对 |
|------|------|----------|
| 用户三问 | `01-decision.md` P-004 表；`D-001` 开篇；`E-002` | maintenance 呈现 **B**；升级入口 **A**；Shell 可见性 **A**。未选方案已记。 |
| 派生 vs 原话 | `D-001` §1「派生口径（非第二轮产品分叉）」 | 生产者折叠为 Host `degraded`、不改消费者终态、`/me.runtimeMode` additive 均标为派生。用户原话停在「改 Host **生产者**投影使 Shell 仍加载」。 |
| 四模式矩阵 | `attachments/r1-mode-projection-matrix.md` | 写门禁/错误码与现行 `operational.go` 一致；Host 生产者列标明「本波将改」且 maintenance 现行为 `maintenance`。 |
| 诊断/版本矩阵 | `attachments/r1-diagnostic-field-matrix.md` | status 九字段与 `MonitoringStatusRow` 一致；Shell 横幅 vs 版本入口权限拆分与裁决 A 一致。 |
| 侦察 | `r1-recon-i-039-00{1,2,3}-*.md` | 只读事实；maintenance 终态阻断 Shell 的推论与代码一致。未被当成冻结决策。 |
| 信息项 | GOAL-002 `00-meta` / `01-decision` | `I-039-001`～`003` `verified` 锚到 D-001 + 矩阵。关闭证据可重复核对。 |
| 非实现 | `E-001`/`E-002`；git `e341c9c3` | 本目标未改 `apps/**`。 |

### 代码对照（现行，冻结尚未实施）

| 主张 | 现行代码 | 与冻结关系 |
|------|----------|------------|
| Host 生产者 `runtime.mode=maintenance` → 文档 `availability.mode=maintenance` | `apps/api/internal/handler/bootstrap.go` `bootstrapAvailability`：`case "maintenance": return {Mode: "maintenance"}`；`bootstrap_test.go` 期望 `maintenance`→`maintenance` | 矩阵「现行为 `maintenance` / 本波改 `degraded`」**诚实**。不是把将改写成已然。 |
| `read-only` 已折叠为 Host `degraded` | 同函数 `case "degraded", "read-only"` | 冻结保持该折叠。与侦察一致。 |
| 写门禁不改 | `operational.go`：maintenance→503/`SERVICE_MAINTENANCE`；degraded→`SERVICE_DEGRADED`；read-only→`SERVICE_READ_ONLY`；登录/恢复/邀请白名单 | 冻结「写门禁合同不改」与现行一致。 |
| Host **消费者** `maintenance` 终态 | `apps/web/src/host/bootstrap.ts` availability-gate：`mode === "maintenance"` → `MAINTENANCE`；`mode === "degraded"` 走完后 → `READY_DEGRADED`。非法 enum（含 `read-only`）在 `isValidBootstrapDocument` fail-closed。 | 冻结「不改消费者终态机 / 不新增 Host enum / 不发 `read-only`」与协议实现一致。 |
| 改生产者即可让 Shell 加载 | `apps/web/src/host/boot.ts`：`preLoad.result` 仅 `READY` 或 `READY_DEGRADED` 才继续拉 manifest；其余进 `HostFailureScreen`。 | **机械成立**：生产文档改发 `degraded` → `READY_DEGRADED` → Shell（含登录页）加载。匿名不再走 MAINTENANCE 终态，这是裁决 **B**（相对未选 **C** 混合匿名终态）的后果，不是隐瞒。 |
| 精确模式不在公开 bootstrap | `GET /api/accounts/me` 现行 `account.Session{User, Features}`，无 `runtimeMode` | 冻结把 `/me.runtimeMode` 标为 R2 additive（T-2），不是现行字段。诚实。 |
| status `availabilityMode` 原样 `runtime.mode` | `composition.go`：`systemmonitoring.New(..., string(cfg.RuntimeMode))` | 与 Host 折叠分离。A-001 F-001 成立。 |
| 版本身份 | `pkg/version`：`Version`/`Commit`/`BuiltAt`；`/healthz` 与 status 用 Version+Commit，无 BuiltAt | 与 C1 权威表一致。`BuiltAt` 不进首波。 |
| 诊断分母 | `MonitoringStatusRow`：`status, availabilityMode, ready, version, commit, uptimeSeconds, moduleCount, modules, dbSizeBytes` | 与冻结矩阵逐字段一致。 |
| mvp/demo 无监控模块 | `kernel/profile.go`：`admin.system-monitoring` 仅 `ProfileAdmin`；`monitoring.read` 只由该模块贡献 | 裁决 A 自洽：mvp/demo 可有横幅（`/me`），无版本 chip、无监控页。 |
| pinned fixtures | `apps/web/src/protocol/upstream/host-bootstrap.cases.json` 含 `maintenance-terminal`（消费者期望 `MAINTENANCE`） | 冻结明确**不改** `docs/schemas/**` 与 `protocol/upstream/**`。T-5：消费者仍实现终态；本仓生产文档不再发出该值。未越界。 |
| `apps/**` | `e341c9c3` 仅 `docs/workspaces/workspace-039-...`；工作树 `apps/` 无本目标 diff | 未改。 |

升级入口：仓库根 `QUICKSTART.md` 存在 `schema-ui upgrade` 节（方法 B 默认主路径）。C1 冻的是入口**类型**（该节 / 默认 GitHub blob `main/QUICKSTART.md`），精确公开 URL 允许 R2 钉死。可接受。

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| C1 版本身份与升级入口；`I-039-001` 可关闭 | **满足（合同）** | D-001 §2；诊断矩阵；`pkg/version` / healthz / status 同源；QUICKSTART 升级节存在；Shell 版本仅 `monitoring.read` |
| C2 四模式投影；`I-039-002` 可关闭 | **满足（合同）** | D-001 §1 区分用户 B vs 派生；矩阵 vs 现行代码诚实；生产者→`degraded`→`READY_DEGRADED` 可核对；写门禁不变 |
| C3 诊断分母 + 可见性；`I-039-003` 可关闭 | **满足（合同）** | 九字段分母；复用既有页；横幅全登录 / 版本 `monitoring.read`；mvp 无模块仍自洽 |
| 未实施、未改 pinned、未改默认集 | **满足** | 非目标与 D-001 §4；`apps/**` 未动 |
| 影响本 scope 的 required 信息项 | **无开放** | `I-039-001`～`003` verified 有裁决+矩阵；`I-039-006` deferred non-blocking 不进本门禁 |
| 本目标实现代码 | **正确未做** | 冻结阶段；R2/R3 移交 T-1～T-5 |

勾选 C1～C3 作为**方案冻结**检查点（非实施完成）与 `00-meta`「本目标不实现代码」一致。C4（self+independent、开放 required=0、Root R1 可投影）不在本意见勾选范围内。

## Findings

无 required。

### F-001 · R2 文案源必须是 runtime 名而非 Host 名

- 严重度：med
- 建议：recommended
- 状态：open
- 描述：本波 Host 文档对 `maintenance`/`degraded`/`read-only` **一律** `availability.mode=degraded`。若横幅或 Toast 读 Host 模式，三分文案会塌成同一个「降级」，直接违反 D-001 §1.5。system-monitoring `availabilityMode` 已是原样 `runtime.mode`（`composition.go` 传入 `cfg.RuntimeMode`）。`/me.runtimeMode` 合同也是原样字符串。R2 应用同一常量/类型，且 **Shell 横幅只读 `/me.runtimeMode`**。
- 证据：`D-001` §1.2–1.5；`bootstrap.go` 现行 vs 将改；`host/bootstrap.ts:275,317-318`；`systemmonitoring` + `composition.go:557`
- 与 A-001：同意其 F-001；本条补上「禁止用 Host `availability.mode` 当产品文案」的实施约束。不阻断 C4。

### F-002 · 现行 `fetchMe` 会丢掉 additive `runtimeMode`

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：`GET /api/accounts/me` 今日只返回 `user`+`features`。Web `fetchMe`（`apps/web/src/account/auth-client.ts`）把 JSON 收窄为 `{ user, features }`，**丢弃**其余字段。R2 T-2 不能只改 Go 结构体：必须扩展 `AuthSession`/`fetchMe` 才把 `runtimeMode` 送到 Shell。冻结把该字段标为派生 additive 仍然正确；这是移交缺口，不是合同谎言。
- 证据：`account.go` `meHandler`；`account/session.go`；`auth-client.ts` `AuthSession` 与 `fetchMe` 返回 `{ user, features: body.features ?? {} }`
- 与 A-001：self 未写客户端投影。不阻断 C4。

### F-003 · 相邻台账未与冻结同步（非合同缺陷）

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：冻结正文（D-001、两矩阵、GOAL-002 信息表、检查点 C1–C3）内部一致。下列**相邻**表面仍停在侦察后、P-004 前：`workspace.md` 纲领表写 GOAL-002 `0/4`「等 P-004」；GOAL-002 `00-meta` 备注仍写 `progress: 0/4`（frontmatter 已是 `3/4`）；Root `01-decision.md` 信息表 I-039-001～003 仍 `open`（同文件声明权威在 Root `00-meta`，后者已 verified）；VP-039 信息表 I-039-001～003 仍 `open`（愿景目录不是 goal-tree 权威）。不构成「裁决未落盘」。`/govern` 响应 C4 时应刷新 workspace 上下文与 Root 索引副本；VP 表由 `/vision` 同步。
- 证据：上述路径 vs `GOAL-002/00-meta.md`、`01-decision.md`、`goal-tree.md`
- 与 A-001：self 未记。不阻断本冻结门禁。

## 必改项汇总

无。开放 required = 0。

## 与既有意见的异同（A-001 self）

| 项 | A-001 (self) | 本意见 (independent) |
|----|----------------|----------------------|
| verdict | pass | **同意 pass** |
| 三问 B/A/A 落盘 | 一一对应 | 同意；并核到 D-001 / 决策索引 / E-002 |
| 派生口径标注 | 已区分 | 同意；生产者 `degraded`、不改消费者、`/me` 字段均非用户原话 |
| 矩阵 vs 现行 maintenance 文档 | 「将改」诚实 | **独立复现** `bootstrap.go` + `bootstrap_test.go` |
| 生产者能否让 Shell 加载 | 断言一致 | **独立复现** availability-gate + `boot.ts` 预加载终态 |
| pinned fixtures | 避免改 pinned | 同意；`host-bootstrap.cases.json` `maintenance-terminal` 应保留；R2 只改本仓生产者测试 |
| `apps/**` | 未改 | 同意（`e341c9c3` 仅 docs） |
| mvp 可见性 | 未展开 | 补证：`ProfileMVP`/`ProfileDemo` 无 `admin.system-monitoring` |
| findings | F-001 recommended（同名字面值） | 同意并收紧为 F-001；**新增** F-002（`fetchMe` 丢字段）、F-003（相邻台账漂移） |
| 冲突 | — | **无冲突**。无 P-004 意见冲突需用户裁 |

## 结论 + 建议给编排器/用户的下一步

独立意见：**pass**。C1～C3 作为信息/方案冻结可核对；用户裁决 B/A/A 如实落盘；派生未冒充原话；矩阵对现行代码诚实；改 Host 生产者让 Shell 加载的机制成立；未越界改 pinned upstream；横幅/版本权限拆分与 mvp 缺模块自洽；本目标未改 `apps/**`。

建议 `/govern`：

1. 响应 A-001 F-001 与本 A-002 F-001～F-003（recommended：可带入 R2 移交，不阻断勾选 C4）。
2. 开放 required=0 后，由用户确认是否勾选 C4 并将 Root R1 检查点投影（本意见不改检查点/`goal-tree`）。
3. 刷新 `workspace.md` 与 GOAL-002 `00-meta` 备注中的过期 `0/4` 表述（F-003）。
4. VP-039 信息表 I-039-001～003 仍 `open` → 另走 `/vision` 同步，勿在 Goal 台账里改 VP。

## 声明

本意见 `source: independent`，**不**修改目标 `status` / `progress` / 检查点 / `goal-tree` 状态列 / 方案正文 / VP status。响应由 `/govern` 处理。
