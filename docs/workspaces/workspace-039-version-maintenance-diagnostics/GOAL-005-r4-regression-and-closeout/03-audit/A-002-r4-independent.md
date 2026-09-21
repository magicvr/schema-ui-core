---
id: A-002-r4-independent
doc: audit-entry
parent: GOAL-005-r4-regression-and-closeout
status: recorded
source: independent
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
auditor: grok-build (grok-4.6 · reasoning high)
---

# A-002 · R4 回归与关门独立交叉审计

| 字段 | 值 |
|------|-----|
| source | `independent` |
| auditor | grok-build（grok-4.6 · reasoning high · `/audit`） |
| date | 2026-09-19 |
| 类型 | `close-out` |
| scope | `[workspace-039-version-maintenance-diagnostics] GOAL-005-r4-regression-and-closeout` · VP-039 判据 1–7；R4 C1–C3 证据；C4 independent 半边；Root/VP 用户书面关门门禁 |
| verdict | **pass**（close proposal only; user gate remains open） |
| 开放 required findings | **0**（1 recommended） |

## 范围与区间

### 工作区绑定（已校验）

| 项 | 观察 | 结论 |
|----|------|------|
| `workspace.md` | `id=workspace-039-version-maintenance-diagnostics`；`root_goal=GOAL-001-version-maintenance-diagnostics`；`canonical_scope` 仅本区；`vision_role=delivery`；`primary_plan`/`plan_refs`=`VP-039-version-maintenance-diagnostics` | 绑定合格 |
| 共享资料 | `shared_materials_catalog: none` | 无非法固定引用；未把资料当事实 |
| 目标位置 | 本区扁平 `GOAL-005-r4-regression-and-closeout/`；`parent=GOAL-001-version-maintenance-diagnostics` | 在 canonical 内 |
| Root / VP | Root `00-meta` **`active · 3/4`**；VP-039 **`active` v0.2.0**；Charter `schema-ui-core-admin-foundation@0.4.0` **active** | 用户书面确认仍是开放门禁 |
| 跨区 | 残余挂具引用 `docs/vision/roadmap.md`（愿景台账，非他区目标状态）；未读取其他工作区五件套当作本区事实 | 合格 |

本意见**不**修改 `status` / 检查点 / 派生 `progress` / 方案正文 / `goal-tree`。独立意见 `pass` **不等于** Root/VP 已关闭。

### 本轮核验问题（用户指定；本重试不重跑全量套件）

1. 退出矩阵是否诚实区分 pass、既有 bounded harness residual、用户 close gate。
2. mvp full 17/5/1 + isolated 1 pass + admin 7/1/0 是否足以作为可运行范围证据；不把既有 fresh-seed 排序 residual 升格为 VP-039 产品 required。
3. Go / Vitest / typecheck 是否有证据路径。
4. R1–R3 cross chains 是否 required=0。
5. 是否存在越界改动。
6. Root/VP 仍 active 且用户书面确认仍是开放门禁。

上一次 independent 因重跑大套件超时且未落盘。本条按用户书面指令：**只读证据与代码边界**，不重跑 `go test ./...`、全量 Vitest、全量 Playwright。

## 成果（有证据）

### 证据路径（未复跑）

| 面 | 台账主张 | 本条独立核对（不复跑） |
|----|----------|------------------------|
| API 全量 | `apps/api`: `go test ./... -count=1` **PASS** | 证据路径 = `02-execution/E-002-r4-regression.md` + `attachments/r4-regression-matrix.md`，随 `a77feab6` 入仓。本条未复跑，不把「未复跑」写成产品失败 |
| Web Vitest | direct `vitest run` **123 files / 1489 tests PASS** | `git ls-files apps/web` 中 `*.test.*`/`*.spec.*` = 136；减去 13 个 `e2e/*.spec.ts` = **123**，与台账文件数一致。1489 用例数以 E-002 为路径，本条未重数 |
| TypeScript | direct `tsc -b --force` + `tsc -p e2e/tsconfig.json` **PASS** | 同一证据路径；`pnpm run typecheck`/`test` 的 Windows EPERM 被明确标为包装器问题而非产品失败，诚实 |
| R2/R3 定向 | RuntimeBanner / VersionChip / AuthGate / catalog | 代码仍在；R2 `5929e641` / R3 `6bb44c24` 已由各自 A-002 independent `pass` 审过 |

### Playwright 计数可静态核对（不复跑）

`e2e/*.spec.ts` 中 `test(` 共 **23** 条。mvp 默认 Profile 下 `test.skip`：

- `telegram-operator-layout.spec.ts` ×3（`APP_PROFILE=custom`）
- `jobs-result-center.spec.ts` ×1（admin-only）
- `localization.spec.ts` ×1（admin-only settings）

= **5 skipped**。23 − 5 skip − 1 fail = **17 passed**。与 E-002「17/5/1」算术一致。

admin 切片命令（shell + console-health + host-failure + localization）：1+1+4+2 = 8 条；其中 localization 的 mvp-only 边界在 `APP_PROFILE=admin` 下 skip → **7 passed / 1 skipped / 0 failed**。与 E-002 一致。

唯一失败位点：`list-visual-surface.spec.ts:82` 是视觉契约 `test(` 起始行，首个 `await` 为 `openRolesList` → `signInAsAdmin`（`sign-in.ts`）。失败形态是 helper fallback「Sign in 按钮启用超时」，不是搜索框/布局断言。同文件另 3 条（:221 + roles/users pageSize）在全量中计入 passed，支持「顺序挂具」而非「VP-039 产品断言失败」。隔离 1/1 以 E-002 为路径，本条未复跑。

该机制已登记于 `docs/vision/roadmap.md`（fresh-seed 顺序契约 bounded residual）。**不升格为 VP-039 产品 required。**

### 代码边界（HEAD `a77feab6`）

VP-039 产品区间 `6197e802..HEAD` 仅两笔功能提交：`5929e641`（R2）与 `6bb44c24`（R3）。R4 `b6e6eb00..a77feab6` 产品树仅 `version-chip.test.tsx` 一行 matcher 改写（`toHaveBeenCalledWith(expect.stringContaining(...))`），属回归断言风格，不削弱门禁。

| 红线 | 独立核对 |
|------|----------|
| `apps/api/kernel/profile.go` 默认集 | `git diff 6197e802 HEAD --` **空**。`ProfileMVP`/`ProfileDemo` 仍无 `admin.system-monitoring`；该模块仅 `ProfileAdmin` |
| `apps/web/src/protocol/upstream/**` | 同区间 log/diff **空** |
| `handler/operational.go` 写门禁 | 同区间 diff **空**。`SERVICE_MAINTENANCE` / `SERVICE_READ_ONLY` / `SERVICE_DEGRADED` + recovery allowlist 仍在 |
| Host 生产者折叠 | `bootstrap.go` `bootstrapAvailability`：`maintenance`/`degraded`/`read-only` → Host `degraded`；精确模式走 `/me.runtimeMode` |
| `/me.runtimeMode` | `account.go` `meHandler`：空串 → `normal`，否则原样 |
| Shell 横幅 | `runtime-banner.tsx` 仅三非 `normal` 模式渲染；`AuthGate` 未登录走 `LoginPage`，不挂 `App`；登录后传 `session?.runtimeMode` |
| VersionChip | 无 `monitoring.read` 不 fetch；只读 `items[0].version`；QUICKSTART blob；诊断按钮随 `system-monitoring` 页存在 |
| Redis / MQ / timestamptz / 新模块 | VP-039 产品 diff 19 个 api/web 文件，无上述引入 |

### R1–R3 cross

| 目标 | 索引 | 开放 required |
|------|------|----------------|
| GOAL-002 | A-001/A-002 `pass` + A-003 响应 | **0**；`done · 4/4` |
| GOAL-003 | A-001/A-002 `pass` + A-003 响应 | **0**；`done · 4/4` |
| GOAL-004 | A-001/A-002 `pass` + A-003 响应 | **0**；`done · 4/4` |

`I-039-001`～`005` **verified**；`I-039-006` deferred **non-blocking**。无到期且影响本 scope 的 required 信息项。

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| 判据 1 分母与契约冻结 | **达成** | GOAL-002 `done · 4/4`；cross required=0 |
| 判据 2 维护提示可感知 | **达成** | GOAL-003 `done · 4/4`；横幅三分 + Host 折叠 + 写门禁未改 |
| 判据 3 版本提示可核对 | **达成** | GOAL-004 `done · 4/4`；VersionChip 权限门 + QUICKSTART |
| 判据 4 诊断摘要只读 | **达成** | 复用既有 `admin.system-monitoring`；无新诊断页 |
| 判据 5 体验与权限 | **达成（可运行证据）** | Vitest 123 文件可核对；mvp 17/5/1 算术与 skip 清单可核对；admin 7/1/0 切片清单可核对；demo/custom 走既有 Profile/Manifest，与 D-001 一致 |
| 判据 6 范围保持 | **达成** | profile / upstream / operational / Redis/MQ/tz 红线未破 |
| 判据 7 证据与审计 + 用户确认 | **独立半边满足；用户确认仍开放** | 本条 + A-001；Root/VP 仍 `active`。E-003 将该行标「待完成」诚实 |
| R4 C1 回归矩阵 | **满足** | `attachments/r4-regression-matrix.md` |
| R4 C2 自动化验证 | **满足（台账 + 静态核对）** | E-002；本条未复跑全量 |
| R4 C3 边界与残余 | **满足** | 残余分类诚实；未升格为 VP-039 required |
| R4 C4 cross | **independent 半边本条；勾选归 /govern** | 开放 required=0；不得把本意见写成已改 status |
| P-004 Root/VP 关门 | **仍开放** | 须用户书面确认 |

退出矩阵（E-003）诚实三分：1–6 达成、浏览器 residual 单列、判据 7 待用户确认。未把 17/5/1 写成全绿，也未把 residual 写成 VP-039 产品缺陷。

## Findings

无 required。

### F-001 · `workspace.md` R4 行仍写 GOAL-005 `0/4`

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：GOAL-005 `00-meta` 与 `goal-tree` 为 `active · 3/4`（C1–C3 已勾）。`workspace.md` 纲领表 R4 仍写 `GOAL-005` `0/4`。不否定回归证据，也不构成产品越界。同类：GOAL-003/004 A-002 台账漂移。`/govern` 响应时可刷新工作区说明。
- 证据：`docs/workspaces/workspace-039-version-maintenance-diagnostics/workspace.md` 纲领表 R4 行 vs `GOAL-005` `00-meta.md` `progress: 3/4`
- 不阻断提请用户确认关闭。

## 必改项汇总

无。开放 required = 0。

## 与既有意见的异同

| 来源 | 关系 |
|------|------|
| GOAL-005 A-001 self `pass` | **同意** 0 required、17/5/1 不作产品失败、Root/VP 用户门禁仍开。self 未写 `workspace.md` 0/4 漂移 → 本条 F-001（recommended）。self 转述套件结果；本条用 git 区间、spec 库存与红线 diff 交叉，**按指令不复跑** |
| GOAL-002/003/004 A-002 | 各 `pass`、A-003 后 required=0。R1–R3 链未重新打开 |
| A-001「无新产品残余」 | **同意**。fresh-seed 仍挂 `docs/vision/roadmap.md`，本 VP 不新开 required |

无冲突。无 P-004 冲突项（用户确认关门是预设开放门禁，不是意见冲突）。

## 结论 + 建议给编排器/用户的下一步

**verdict = pass。** R4 退出矩阵诚实；可运行范围证据充分（Go/Vitest/typecheck 有台账路径；Playwright 17/5/1 与 admin 7/1/0 可被 spec 库存核对）；既有 fresh-seed residual 不升格为 VP-039 产品 required；R1–R3 cross required=0；产品红线未破；Root 与 VP-039 仍 `active`。

建议 `/govern`：

1. 响应 A-001 + A-002（本条）。F-001 可修 `workspace.md` 说明，非门禁。
2. 向用户提请书面确认：「确认 VP-039 与 workspace-039 Root 关门」。
3. 用户确认前保持 GOAL-005 / Root / VP **active**；不要把本意见写成已改 `status`/`progress`。

## 声明

本意见不修改 status/progress；响应由 `/govern` 处理。
