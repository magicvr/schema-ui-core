---
id: A-008-root-a004-a005-closure-independent
goal: GOAL-001-shared-cross-module-contracts
doc: audit-entry
record_id: A-008
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
scope: finding-closure + close-out residual; A-004/A-005 F-001～F-009 remediations; E-009/A-007; affected API/Web verification; Root close-gate evidence; full API handler SQLite VACUUM timeout
audit_type: finding-closure
verdict: pass
status: recorded
parent: GOAL-001-shared-cross-module-contracts
created: 2026-08-19
updated: 2026-08-19
version: 0.1.0
reviews:
  - A-004
  - A-005
  - A-006
  - A-007
---

# A-008 · 独立复审 · A-004/A-005 finding 闭合与 Root 关门证据（2026-08-19）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high；项目级路径见 `docs/architecture/independent-audit-execution.md`）
- **类型**：finding-closure / Root residual close-gate
- **scope**：`workspace-012-shared-cross-module-contracts` Root `GOAL-001-shared-cross-module-contracts`。复核 A-004 F-001～F-007、A-005 F-008/F-009 的修复是否可重复核对；核对 E-009/A-007 主张；独立复跑受影响 API/Web 验证；判断 full API handler `TestNotificationPruneKeepsUnread` SQLite `VACUUM` 超时是否阻断 Root 关门。
- **verdict**：**pass**
- **required findings**：0
- **日期**：2026-08-19

## 范围与区间

- **工作区**：`workspace-012-shared-cross-module-contracts`（`workspace.md`：`id` 与路径一致；`root_goal` = `GOAL-001-shared-cross-module-contracts`；`canonical_scope` = `docs/workspaces/workspace-012-shared-cross-module-contracts/`；`shared_materials_catalog: none`；`vision_role: delivery`；`plan_refs` / `primary_plan` = `VP-012-shared-cross-module-contracts`）。
- **covered**：Root `00-meta` I-001/I-002；E-007～E-009；A-001～A-007；A-004/A-005 全部 finding 的现行代码与定向测试；本轮独立复跑的 API 受影响包、handler 定向切片、孤立 prune 测试、Web 结构化 i18n 与 Web 全量。
- **excluded**：不改 `status` / `progress` / `goal-tree` / `workspace` / 方案正文 / 业务代码；不关闭 VP-012；不读取或比较其他工作区治理上下文；不把派生 `progress=100` 当作闭合证据；不把未重跑的串行 `go test ./...` 写成全量通过。
- **共享资料**：目录为 `none`；无固定引用，不得当作事实或 finding 关闭依据。
- **auditor 立场**：只出意见。本条确认 A-007 的 `fixed` 主张是否可独立核对，不代替 `/govern` 响应。

## 本轮独立复验

在 `apps/api`：

| 命令 | 结果 |
|------|------|
| `go test ./internal/requestid ./internal/auth ./internal/jobs ./internal/modules/operationlog ./internal/docscheck -count=1 -timeout 180s` | **ok**：requestid 0.930s；auth 16.547s；jobs 2.096s；operationlog 5.942s；docscheck 0.688s |
| `go test ./internal/handler -run 'TestOperationalGate\|TestWriteLocalizedError\|TestAuthMiddlewareLocalized\|TestNotificationPruneKeepsUnread' -count=1 -timeout 180s -v` | **ok** 10.557s。含 `TestNotificationPruneKeepsUnread` **PASS 7.28s**；全部 OperationalGate / WriteLocalizedError / AuthMiddlewareLocalized 子测试通过 |

在 `apps/web`：

| 命令 | 结果 |
|------|------|
| `npm test -- --run src/i18n/schema-keys.structural.test.ts` | **4/4** 通过（552ms） |
| `npm test -- --run` | **72/72 files，1069/1069 tests** 通过（8.17s） |

本会话未重跑串行 `go test ./... -p 1 -timeout 300s`。该命令不是本条闭合 F-001～F-009 的必要证据；VACUUM 专项见下节。

## 工作区、信息门禁与 Root 关门边界

| 检查项 | 结论 | 证据 |
|--------|------|------|
| 工作区绑定 | 通过 | `workspace.md` `id` / `root_goal` / `canonical_scope` 与 Root `parent: null`、路径一致 |
| `plan_refs` / `primary_plan` | 通过 | workspace 与 Root 均挂 `VP-012-shared-cross-module-contracts` |
| VP → Charter | 通过 | VP-012 `vision_ref` = `schema-ui-core-admin-foundation@0.2.0`；Charter `status: active`、版本 `0.2.0` |
| 共享资料 | 无引用 | `shared_materials_catalog: none`；本条未用共享资料关闭任何 finding |
| I-001 | 维持 verified | non-blocking；R1 开始前；消费方已由 R1～R6 真实路径回答 |
| I-002 | 维持 verified（本 scope） | required；最晚阶段 = Root 关门。A-002 已对六子目标最终链与四条方向成功标准独立 pass。本条补齐 A-005 之后唯一 required（F-008）与全部 recommended 修复的独立核对。无 `deferred` required，无用户书面 residual/overrule |
| VP-012 | 仍为 `active` | 本条只审 Root finding 闭合与关门证据，**不是** VP-012 关门 |

## 成果（有证据）

A-007/E-009 声称的实现修复，除「heartbeat 错误注入测试已补」外，均可在现行代码中定位。

### A-004 / A-005 finding 闭合表

| 原 finding | 级别 | A-007 主张 | 本轮独立结论 | 证据 |
|---|---|---|---|---|
| A-004 F-001 Claim 使用 `context.Background()` | recommended | fixed | **同意 fixed** | `apps/api/internal/jobs/runner.go:264`：`r.repo.Claim(ctx, ...)`；heartbeat 查询/续租亦用同一 `ctx`（L288、L299） |
| A-004 F-002 / A-005 F-009 heartbeat `IsCancelRequested` 错误直接 return | recommended | fixed | **实现同意 fixed**；A-005 所要的错误注入测试仍缺，见本条 F-001 | `runner.go:288-304`：查询或续租失败先 `reporter.cancelExecution()`，且仅当 `ctx.Err() == nil` 时 `abortLease`；`abortLease`（L317-325）用 background `Fail(..., "JOB_HANDLER_FAILED")` 并 `notifyTerminal`。Stop 路径仍跳过 abort，与 `TestRunnerStopLeavesJobRecoverable`（`runner_test.go:267-301`）一致 |
| A-004 F-003 request-id 回退可碰撞 | recommended | fixed | **同意 fixed** | `apps/api/internal/requestid/requestid.go:54-62`：CSPRNG 失败回退 `UnixNano` + `fallbackSequence.Add(1)`；`requestid_test.go:62-70` 注入失败 reader 后两个 ID 唯一且 `Valid` |
| A-004 F-004 `newID()` 回退可预测 | recommended | fixed | **同意 fixed** | `apps/api/internal/auth/auth.go:410-420`：同样原子序列；`NewServiceCredentialToken`（L392-398）仍无回退，CSPRNG 失败返回 error。`auth_test.go:41-49` 回退唯一性通过 |
| A-004 F-005 `writeLocalizedError` 重复 | recommended | fixed | **同意 fixed** | 共享实现 `apps/api/internal/errorcatalog/writer.go:13-33`；`handler/localize.go:17-19` 与 `auth/auth.go:425-427` 均委托。`writeLocalizedFieldError` 仍留在 handler（L25-41），属于 fieldErrors 扩展面，不是第二套基础包络 |
| A-004 F-006 allowlist 精确路径 | recommended | fixed | **同意 fixed** | `handler/operational.go:63-81` 显式 `operationalRecoveryPaths` registry，覆盖 login/refresh/logout/mfa-verify/改密及全部 MFA 自助恢复路径。这是 A-004 建议的「显式注册」路径，不是前缀匹配。`operational_test.go:103-120` 对全部 9 条 recovery path 断言 204 |
| A-004 F-007 `redactValue` 不支持类型直接失败 | recommended | fixed | **同意 fixed** | `operationlog/detail.go:154-167`：default 分支 JSON 归一化后再递归脱敏；不可序列化则 `fmt.Sprint`。`detail_test.go:89-116` 覆盖 `map[string]string`、`[]int`、带 `password` 的 struct |
| A-005 F-008 Web 结构化 i18n key 缺失 | **required** | fixed（A-006/E-008） | **同意 fixed** | schema `system-monitoring.json:58-66` 使用 `schema.systemMonitoring.statCard.availability`；`en-US.json:292` = `Availability`；`zh-CN.json:292` = `可用性`。本轮结构化 4/4、全量 72/72 · 1069/1069 通过 |

## 专项：full API handler SQLite VACUUM 超时是否阻断

**结论：不是阻断。不得作为 required finding，也不得据此重开 Root 关门门禁。**

| 主张 | 本轮核对 |
|------|----------|
| E-009/A-007 未把串行 `go test ./... -p 1 -timeout 300s` 写成 full API pass | **属实**。A-007 L57 明确「未被计为 finding，也未被表述为 full API pass」 |
| 超时指向 `TestNotificationPruneKeepsUnread` 的 VACUUM 初始化 | 该测试本身（`notifications_test.go:234-266`）只做 500 条通知 prune，不直接调用 `VACUUM`。`VACUUM INTO` 仅出现在 `store/migrate.go:231-245` 的升级快照：文件库且已有业务行、且存在 pending `version >= 2` 时才执行。`newAuthTestEnv`（`testhelpers_test.go:86-92`）每次打开 `t.TempDir()` 新空库 |
| 该测试是否本身挂死 / 产品缺陷 | **否**。本轮孤立复跑 `TestNotificationPruneKeepsUnread` **7.28s PASS**，与同命令中的 OperationalGate / WriteLocalizedError 切片一起 10.557s 通过 |
| 是否由 F-001～F-009 修复引入 | **证据不足指向本次变更**。修复面是 jobs/requestid/auth/errorcatalog/operational/operationlog/i18n，不触及 notifications prune 或 migrate snapshot |
| 与 Root 成功标准 / I-002 的关系 | Root 四条方向成功标准不依赖「每一次串行 `./... -p 1 -timeout 300s` 全绿」。A-002 close-out 本身也只复跑定向包，未把未重跑全量写成不实。A-005 当时 `go test ./...` 通过；其后 300s 串行预算耗尽是套件时长/环境问题，不是当前可复现的测试失败 |

因此：I-002 不因该超时重新变成开放 required。后续若要收紧套件预算，应作为工程维护（加长 timeout、并行、或减少每测 bcrypt/migrate 成本），不是本 Root 的关门阻断。

## 对照 Root 方向成功标准（本 scope 残余）

| # | 标准 | 本轮 | 证据 |
|---|------|------|------|
| 1 | 每个契约有可验证实现路径 | 维持达成 | A-002 已核；本轮 requestid/auth/jobs/operationlog/handler 定向与 Web 全量仍绿 |
| 2 | 至少一个真实模块消费 | 维持达成 | 未撤销 wallet / operationlog / auth / Host / service-credential 消费面 |
| 3 | 不改变 Profile/矩阵/Manifest/协议 pin | 本条未发现回写 | 本轮修复未改 kernel profile / protocol pin；docscheck ok |
| 4 | Tier D 不进入 | 维持达成 | 修复均为横切基架，无新业务域 |

A-005 用以否决关门的唯一 required（F-008）现已可重复核对为 fixed。A-004 当时即为 conditional + 无 required。本条不重开 Root。

## Findings

### F-001 · Job heartbeat 数据库错误路径仍无注入回归测试

| 字段 | 值 |
|------|-----|
| level | recommended |
| 严重度 | low |
| status | open |
| evidence | A-005 F-009 原文要求「新增数据库错误注入测试覆盖此行为」。`apps/api/internal/jobs/runner_test.go` 现有用例覆盖 cancel、retry、heartbeat 防双 Claim、Stop 可恢复（L156-301），**没有任何**对 `IsCancelRequested` / `Heartbeat` 返回错误后调用 `abortLease` 的注入测试。`abortLease` 仅在 `runner.go:317-325` 出现 |

实现本身已关闭原缺陷（见上表）。缺测试不把 A-004 F-002 / A-005 F-009 打回 open required，也不阻断 Root。建议 `/govern` 后续补一条错误注入测试，或以书面 residual 接受该验证缺口。

### F-002 · `finish()` 在取消查询失败时仍直接放弃终态写入

| 字段 | 值 |
|------|-----|
| level | recommended |
| 严重度 | low |
| status | open |
| evidence | `apps/api/internal/jobs/runner.go:327-335`：handler 已返回后，`IsCancelRequested(context.Background(), lease)` 若出错则 `return`，不调用 `abortLease` / `Fail` / `FinalizeCancel`。L357-359 的 Complete 失败回退路径同样如此 |

这与已修的 heartbeat 循环不是同一条路径。handler 已结束，租约仍会到期回收，故不升 required。建议与 F-001 一并决定是补对称 cleanup 还是接受「完成路径仍依赖 lease 过期」的残余。

## 必改项汇总

**无。** required = 0。

VACUUM 超时：**非阻断，非本条 finding。**

## 与既有意见的异同

| 点 | A-005 / A-007 | 本意见 |
|----|----------------|--------|
| F-008 required | A-005 fail；A-006/A-007 称 fixed | **同意 fixed**；本轮 Web 1069/1069 独立复跑 |
| F-001～F-007 实现 | A-007 称全部 fixed | **同意**；逐条 file:line 核对 |
| F-009 实现 | A-007 称与 F-002 一并 fixed | **同意实现已修** |
| F-009 错误注入测试 | A-007 写「jobs 测试通过」 | **不同意把现有 jobs 测试当作该路径的覆盖**；本条 F-001 recommended |
| 串行 full API VACUUM 超时 | A-007 不计 finding、不称 pass | **同意**；本轮孤立测试 7.28s 通过，明确 **非阻断** |
| Root 关门 | A-002 pass 后已 `done`；A-005 曾 fail；A-007 等待本独立复审 | 本条 **pass**，不改 status；不因 recommended 测试缺口或套件预算重开 Root |
| P-004 冲突 | A-004 conditional vs A-005 fail 已由 A-006 按 F-008 `fixed` 响应 | 本条与 A-007 无「一要一否」的 required 冲突 |

## 结论 + 建议给编排器/用户的下一步

A-004 F-001～F-007 与 A-005 F-008 的修复证据充分、可重复核对。A-005 F-009 的 **实现** 已修；所要的错误注入测试仍缺，仅构成 recommended。full API handler SQLite VACUUM 超时 **不是阻断**：孤立 `TestNotificationPruneKeepsUnread` 本轮 7.28s 通过，且 A-007 未把它伪装成全量通过。

开放 required = 0。无到期且影响本 scope 的 required 信息项。

建议 `/govern`：

1. 响应本条 A-008（`pass`，required=0）。可将 A-004 F-001～F-007、A-005 F-008、A-005 F-009（实现）维持 `fixed`。
2. 对本条 F-001/F-002：补测试/对称 cleanup，或书面 `accepted-residual`（范围 + 复审触发）。二者皆不阻断维持 Root `done`。
3. **不要**把串行 `go test ./... -p 1 -timeout 300s` 的历史超时写成产品回归，也 **不要**把本条写成 VP-012 已关闭。
4. 本意见不改 `status` / `progress` / 检查点 / 方案正文 / `goal-tree` / 业务代码。

## 声明

本意见 `source: independent`，不修改目标 `status` / `progress` / 检查点 / 方案正文 / `goal-tree.md` / `workspace.md` / 子目标五件套 / 业务代码 / `docs/vision`。响应、finding 闭合与状态变更由 `/govern` 处理。保证等级为框架默认 **L0**（入口分离），不得表述为第三方鉴证。
