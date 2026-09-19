---
id: A-002-w34-independent
doc: audit
parent: GOAL-046-w34-batch-export-availability
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
status: recorded
source: independent
auditor: grok-4.6（grok build · high · `/audit`）
verdict: conditional
open_required: 2
audit_type: close-out
---

# A-002 · W34 独立交叉审计（GOAL-046 C1～C3）

- **source**：independent
- **auditor**：grok-4.6（grok build · 思考强度 high · `/audit`）
- **date**：2026-09-19
- **类型** / **scope**：close-out · W34 全量 —— 根因诊断、修复面是否越界、fail-open 取舍、两层守卫是否名实相符、残余登记
- **verdict**：**conditional**（2 required / med + 3 recommended；无 high required；无到期未闭合 required 信息项）
- **工作区**：`workspace-010-design-implementation-conformance`（`root_goal` = `GOAL-001-design-implementation-conformance`；`canonical_scope` 匹配；`shared_materials_catalog: none`；`primary_plan` = `VP-010-design-implementation-conformance`）

## 范围与区间

对照 `D-001`、`E-001`、self `A-001` 与 checkpoint `e1013c8b`（parent `784e9152`）。只审本区 GOAL-046 与其授权代码路径；跨区仅以 Q2 路径提及 `[workspace-038]` / VP-038，未读取其他工作区台账正文。

P-005：`I-046-001`/`002` required 已在 C1 前 verified；`I-046-003` non-blocking 已 verified；无 `deferred` required，无 accepted-residual。共享资料引用：无。

## 成果（有证据）

1. **根因成立，且更简单的替代解释不成立。** `git show 784e9152:apps/api/configs/config.yaml` 有 `profile: custom` 与 `- admin.data-transfer`，**没有** `admin.jobs`。`kernel/profile.go` 的 `ProfileAdmin` 含 `admin.jobs`，`ProfileMVP`/`ProfileDemo` 都不含。`modules/jobs/provider.go` 仅在 `submitter != nil` 时声明 `POST /api/jobs/batch-export`；生产装配 `composition.go` 在 `plan.HasModule("admin.jobs")` 时始终注入 submitter。未挂载时 `route_envelope.go` 走 `NOT_FOUND` → catalog `error.notFound` → zh-CN「未找到」（`apps/web/src/i18n/messages/zh-CN.json`）。对照：`data.export` 由 `admin.data-transfer` 贡献且修复前列表已有该项，故不是「缺导出权限」；JSON 信封也不是代理 HTML 404；空选择在组件里根本不会发请求。用户可见的「未找到」与未挂载路由对齐。
2. **修复面最小，未越 D-001 §6 禁止路径。** `git diff --stat 784e9152..e1013c8b` 仅 6 文件：`apps/api/configs/config.yaml`、`apps/api/internal/config/operator_config_coverage_test.go`、`apps/web/src/components/jobs-batch-export.tsx`、`apps/web/src/components/jobs-batch-export.test.tsx`、`apps/web/src/i18n/messages/{en-US,zh-CN}.json`。未改 `apps/api/internal/handler/jobs.go`、`apps/api/modules/jobs/export.go`、`apps/web/src/protocol/upstream/**`、`docs/schemas/**`。服务端两道门禁（`jobs.write` 然后 `data.export`）仍是唯一授权方；客户端是镜像。
3. **权限贡献与「两权限齐备 ≈ 路由存在且被授权」在生产装配上大致成立。** 非测试 Go 代码中 `"jobs.write"` 只由 `admin.jobs` 贡献，`"data.export"` 只由 `admin.data-transfer` 贡献。`admin.jobs` 在 plan 中时 submitter 非 nil，batch-export 路由会挂载。残留面见 F-004（客户端硬编码）与 fail-open 未知上下文。
4. **「可见 + disabled + 说明」的结果可接受，高度契约当前为绿。** 本轮 `list-visual-surface.spec.ts` **4 passed**（约 1.1 min）。W33 `D-001` §3 原文冻结的是**插槽布局** fail-open（schema 笔误不得让入口静默消失），不是「权限不足必须显示禁用」；真正卡住「return null」的是插槽宿主仍注册时留下高度 0 的空 flex 子元素。当前实现不再隐藏，高度契约通过。第一版「隐藏」回归叙述机械上可信，已在同一 commit 内修正，**不升级为开放 finding**。
5. **配置守卫对「删掉 / 注释掉 `- admin.jobs`」名实相符。** 本轮注释掉该行 → `TestOperatorConfigCoversAdminPreset` 失败并指名 `missing 1 admin-preset module(s): admin.jobs`；`git checkout -- configs/config.yaml` 后复跑 PASS。把 `profile: custom` 改成 `profile: admin` → **SKIP**（`operator config no longer uses an explicit custom module list`），随后已还原。
6. **本轮复跑通过的回归（不等于 C3 全量 e2e 数字已独立重跑）。** vitest **121 files / 1481 tests** 全绿；定向 vitest 17 passed；`go test ./internal/config/ ./modules/jobs/ -count=1` 全绿；可用性 e2e mvp **1 passed**、admin **1 passed**（工作区未跟踪的 spec，见 F-002）。

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| C1 诊断：404 来自未挂载路由，非数据缺失 | **达成** | 修复前 YAML 无 `admin.jobs`；envelope + profile 集合 + 贡献源 |
| C2 实施：config 补模块 + 触发面镜像 + 404 文案 | **达成** | `e1013c8b` 六文件；i18n 键 `schema.jobs.batchExport.unavailable` |
| C3 守卫名实相符 | **部分** | config 变异有效；XOR 单测空洞（F-001）；e2e 契约未入仓（F-002） |
| 修复未改服务端门禁 / Job 六态 / pinned / VP-038 status | **达成** | `git diff --name-only 784e9152..e1013c8b` |
| 信息门禁 I-046-001～003 | **达成** | 均 verified；无到期 required |
| C4 投影 / 残余登记 | **部分出现、内容超前** | 本意见**未写** roadmap/workspace/goal-tree。核验中后期工作区出现 W34 行，但写成「当日关门」且把未入仓 e2e 当作已交付守卫（F-005） |

## Findings

### F-001 · 「门禁缺一」单测与空选择禁用不可区分

- 严重度：med
- 建议：**required**
- 状态：open
- 描述：`jobs-batch-export.test.tsx` 中 `keeps the trigger enabled when only one of the two gates is granted` **标题与断言相反**（断言 `disabled === true`），且**未选行**。按钮在 `count === 0` 时本来就会 disabled，因此该例在「只检查 `jobs.write`、忽略 `data.export`」的错误实现下仍会绿。真正锁住「两门禁缺一」需要：选中至少一行 + 断言 `data-jobs-batch-export-unavailable="true"`（并建议补 `data.export` 有、`jobs.write` 无的对称例）。「两门禁都缺」那一例（检查 unavailable 标记与说明文案）仍然有效，不能替代 XOR。
- 证据：`apps/web/src/components/jobs-batch-export.test.tsx`（availability describe 内「only one of the two gates」例）；对照同文件空选择例与 `canBatchExport` 的 AND。
- 关闭要求：补选行 + unavailable 断言（建议标题改为 disabled / unavailable）；或删除该例并在「两门禁都缺」之外另写名实相符的 XOR 例。

### F-002 · 双 profile 可用性 e2e 未进入 checkpoint，self F-003 关闭过早

- 严重度：med
- 建议：**required**
- 状态：open
- 描述：self `A-001` 以新增 `apps/web/e2e/batch-export-availability.spec.ts` 将 F-003 标为 **fixed**。独立核验：该文件对 `e1013c8b` 为 **untracked**（`git ls-tree` / `git log` 均无；`git status` 为 `??`），也不在 `D-001` §6 授权路径里。`E-001` §4 引用的全量 e2e 数字（mvp 17/5/0、admin 18/4/0）与 W33 关门数字相同；若该 always-on spec 已计入全量套件，pass 数应各 +1。本轮在工作区跑该 spec：mvp 1 passed / admin 1 passed，说明**文件内容本身在默认 mvp/admin 上能工作**，但未入仓的守卫对 CI/克隆无效，不能作为 F-003 的关闭证据。
- 另：`adminCapable = appProfile === "admin" \|\| appProfile === "custom"`，而 `apps/web/playwright.config.ts` 的 `customE2EModules` **当前就不含** `admin.jobs`（有 `admin.data-transfer`）。self 写成「若将来 custom 列表去掉 admin.jobs」——独立判断这不是未来时，是**现在时**。默认 CI 只跑 mvp/admin，故尚未在主矩阵上假绿；一旦有人跑 `APP_PROFILE=custom`，该 spec 会按「可服务」断言 enabled，与 harness 模块集冲突。
- 证据：`git status --short`；`git diff --name-status 784e9152..e1013c8b`；`apps/web/e2e/batch-export-availability.spec.ts:21-22`；`apps/web/playwright.config.ts` `customE2EModules`；`E-001` §4 vs W33 回归数字。
- 关闭要求：① 将该 spec 纳入版本控制，或撤回 self F-003 的 `fixed`；② 修正 `adminCapable`（custom 与 harness 模块集对齐：要么 custom 列表补 `admin.jobs`，要么 custom 走「可见但禁用」分支）；③ 若入仓，复跑双 profile 全量 e2e，数字须反映该 spec。

### F-003 · 配置守卫在形态变化时 Skip 而非失败

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：与 self `A-001` F-001 同向。手写扫描只认 `profile: custom` + `list:` 下 `- <id>`。本轮把 `profile: custom` 改为 `profile: admin` → **SKIP 且套件 PASS**（已还原）。内联数组 `list: [a, b]` 因扫不到 `- ` 项会 `Fatalf`（偏 fail-closed）。注释掉模块行会失败（本轮已证实）。跳过是显式的，但「不变量失效却不报警」。
- 证据：`operator_config_coverage_test.go:44-46`；本轮 preset 变异输出 `operator config no longer uses an explicit custom module list`。
- 关闭要求：形态变化时同步守卫，或对非 custom 形态显式失败/断言新契约。

### F-004 · 客户端可用性镜像与服务端门禁清单存在漂移面

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：与 self `A-001` F-002 同向。`canBatchExport` 硬编码 `jobs.write` + `data.export`。服务端再加第三道门禁时，入口可能仍显示可用、提交被拒；404 文案已能诊断未挂载，403 仍回显 `body.message`。服务端仍是唯一授权方。上下文完全没有 `permissions` 时 fail-open（视为可用）会让 mvp 裸挂具保持可点，属有意取舍。
- 证据：`apps/web/src/components/jobs-batch-export.tsx` `canBatchExport`；`apps/api/internal/handler/jobs.go` 提交路由两道 `requirePermission`。
- 关闭要求：新增服务端门禁时同步组件判据（PR 检查项或 `D-001` §3 增补）。

### F-005 · C4 投影把未闭合的独立审计写成「当日关门」

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：本意见开读时 `roadmap.md`「最近更新」止于 W33 / VRev-100，未决项无 W34 行。核验过程中（非本审计员写入）工作区出现：`roadmap.md` 新增 `[workspace-038]`「导出所选」404 行并标 `fixed`，且点名未入仓的 `batch-export-availability.spec.ts`；`workspace.md` 追加 W34 段并写「立项并当日关门」+ self F-003 本波 `fixed`。登记缺陷本身是 C4 该做的；问题是投影**抢先关门**，并把 F-002 指出的未入仓 e2e 写成已交付守卫。独立 verdict 为 conditional，开放 required = 2，与「当日关门」冲突。
- 证据：开读时的 roadmap「最近更新」；后来 `git diff` 显示 `docs/vision/roadmap.md`、`workspace.md`、`goal-tree.md` 被他方修改。本意见未改这三份文件。
- 关闭要求：由 `/govern` 在响应 F-001/F-002 之后再改投影：未决项可保留缺陷事实，但不得在开放 required 清零前写 `done` / 「当日关门」；e2e 入仓前不要把它写成已锁死的守卫。

## 必改项汇总

1. **F-001**：XOR 可用性用例必须在选中行后断言 `unavailable`，不能只断言 `disabled`。
2. **F-002**：可用性 e2e 要么入仓并修正 `adminCapable`/`customE2EModules` 对齐，要么撤回 self F-003 的 `fixed`。未入仓不得当作覆盖缺口已永久锁死。

开放 required = **2**。recommended 不阻断，但 F-005 属于 C4 清单。

## 与既有意见的异同（self A-001）

| 点 | self A-001 | 本意见 |
|----|------------|--------|
| 根因 / 修复面最小 | 同意 | 同意；独立用 `784e9152` YAML 与 diff 复核 |
| fail-open 必须可见禁用 | 视为 W33 §3 冻结 | **部分不同意**：W33 §3 是布局回落；硬约束是空插槽宿主高度。结果（disabled + 说明）仍可接受，不要求改实现 |
| config 守卫 | F-001 recommended（Skip） | 同意，独立变异确认 Skip；编号本意见 F-003 |
| 客户端镜像漂移 | F-002 recommended | 同意，本意见 F-004 |
| e2e 覆盖缺口 F-003 | **fixed** | **不接受关闭**：文件 untracked，且 custom 推断与 harness 模块集不一致（F-002 required） |
| 门禁缺一用例 | 算作 +1 有效守卫 | **不同意**：未选行，名不副实（F-001 required） |
| 第一版隐藏回归 | 成果 / 已修正偏差 | 同意不升级为开放 finding |
| verdict | pass（0 required） | **conditional**（2 required） |

无 P-004 冲突项需用户在两条 required 之间二选一；两条都是「补证据/补守卫」，不是产品方向对立。

## 结论 + 建议给编排器/用户的下一步

**verdict = conditional。** 用户报告的「未找到」根因、config 补齐、客户端不再撒谎、服务端门禁未改——这些成立，且本轮命令复验支持。不能无条件关门的原因是 C3 声称的两处守卫（XOR 单测、双 profile e2e 契约）名实不符：前者测不到缺一门禁，后者不在 checkpoint 里。

**关门建议：** W34 **尚不具备关门条件**。请用 `/govern` 响应本意见：

1. 先闭合 F-001、F-002（fixed：补测 + 入仓/对齐；或对 F-002 选择「撤回 F-003 关闭、把 e2e 降为后续」——若选后者须用户书面 overruled/residual，因 self 已把覆盖缺口标 fixed）。
2. C4 投影（roadmap / workspace.md / goal-tree）已有他方草稿，须按 F-005 收回「当日关门」措辞，等 required 闭合后再同步；Root 保持 active。
3. 响应落盘后再考虑关门；**不要**在开放 required 未闭合时把 GOAL-046 标 `done`。
4. `00-meta.md` 概述已写「立项并当日关门（done · 4/4）」与 `status: active` / `progress: 0/4` 矛盾——属编排卫生，由 `/govern` 改，本意见未动 meta。

## 本轮实际执行的命令与观察

| 命令 | 观察 |
|------|------|
| `git show --stat e1013c8b` / `git diff --stat 784e9152..e1013c8b` | 6 文件，+285/−6；未触禁止路径 |
| `git show 784e9152:apps/api/configs/config.yaml` | 有 `admin.data-transfer`，无 `admin.jobs` |
| `git status` / `git ls-tree e1013c8b` | `batch-export-availability.spec.ts` **untracked** |
| `go test ./internal/config/ -run TestOperatorConfigCoversAdminPreset -count=1 -v` | PASS |
| 注释掉 `- admin.jobs` 后同测试 | FAIL：`missing 1 admin-preset module(s): admin.jobs`；**已 `git checkout` 还原**，复跑 PASS |
| `profile: custom` → `profile: admin` 后同测试 | SKIP 且 PASS；**已还原**，复跑 PASS |
| `go test ./internal/config/ ./modules/jobs/ -count=1` | 两包 ok |
| `npm test -- jobs-batch-export.test.tsx jobs-batch-export-roles.test.tsx list-actions-slot.test.tsx`（`apps/web`） | 3 files / 17 tests passed |
| `npm test`（`apps/web` 全量） | **121 files / 1481 tests** passed |
| `npx playwright test e2e/list-visual-surface.spec.ts --reporter=line` | **4 passed**（~1.1 min） |
| 同上 `e2e/batch-export-availability.spec.ts`（默认 mvp） | **1 passed**（16.0s） |
| `APP_PROFILE=admin` 再跑该 spec | **1 passed**（15.3s） |

未做：未再起「修复前 YAML」活进程 curl 404（机制由修复前 YAML + envelope + 路由挂载条件覆盖）；未把组件改回 `return null` 再跑高度契约；未跑 `go test ./...` 全仓。变异改动均已还原。

## 声明

本意见不修改 status/progress/goal-tree/workspace.md/roadmap.md/方案正文/实现代码。响应由 `/govern` 处理。
