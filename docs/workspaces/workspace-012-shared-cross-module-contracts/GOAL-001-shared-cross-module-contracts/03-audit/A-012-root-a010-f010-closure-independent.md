---
id: A-012-root-a010-f010-closure-independent
goal: GOAL-001-shared-cross-module-contracts
doc: audit-entry
record_id: A-012
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
scope: finding-closure；A-010 F-010 关闭复审；R6 使用审计 fail-closed 契约与 Root I-002 受影响门禁
audit_type: finding-closure
verdict: pass
status: recorded
parent: GOAL-001-shared-cross-module-contracts
created: 2026-08-19
updated: 2026-08-19
version: 0.1.0
reviews:
  - A-010
  - A-011
---

# A-012 · 独立复审 · A-010 F-010 闭合（2026-08-19）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high；项目级路径见 `docs/architecture/independent-audit-execution.md`）
- **类型**：finding-closure
- **scope**：`workspace-012-shared-cross-module-contracts` Root `GOAL-001-shared-cross-module-contracts`。只复核 A-010 F-010（R6 使用审计失败仍放行请求）是否已按 `fixed` 路径可重复核对；核对 R6 D-004 / E-008、Root A-011 / E-011 主张，以及受影响的 R6 S3 与 Root I-002。
- **verdict**：**pass**
- **required findings**：0
- **日期**：2026-08-19

## 范围与区间

- **工作区**：`workspace-012-shared-cross-module-contracts`（`workspace.md`：`id` 与路径一致；`root_goal` = `GOAL-001-shared-cross-module-contracts`；`canonical_scope` = `docs/workspaces/workspace-012-shared-cross-module-contracts/`；`shared_materials_catalog: none`；`vision_role: delivery`；`plan_refs` / `primary_plan` = `VP-012-shared-cross-module-contracts`）。
- **covered**：A-010 F-010 原文；A-011/E-011 的 `fixed` 主张；R6 D-004/E-008；现行 `authenticateServiceCredential`、生产 composition 事务 recorder、`MarkServiceCredentialUsedWithAudit`、定向故障注入测试；本轮独立复跑的 auth / authsession / composition / handler `TestServiceCredential*`。
- **excluded**：不改 `status` / `progress` / `goal-tree` / `workspace` / 方案正文 / 业务代码；不重开整个 Root close-out；不重审 A-004～A-009；不读取或比较其他工作区；不把派生 `progress=100` 当作闭合证据；本会话未重跑 API `go test ./...` 与 Web 全量。
- **共享资料**：目录为 `none`；无固定引用，不得当作事实或 finding 关闭依据。
- **auditor 立场**：只出意见。本条确认 A-011 对 F-010 的 `fixed` 主张是否可独立核对，不代替 `/govern` 再响应。

## 本轮独立复验

在 `apps/api`（2026-08-19 本会话）：

| 命令 | 结果 |
|------|------|
| `go test ./internal/auth ./internal/modules/authsession ./internal/composition -count=1 -timeout 180s` | **ok**：auth 25.694s；authsession 33.147s；composition 23.192s |
| `go test ./internal/handler -run TestServiceCredential -count=1 -timeout 180s` | **ok** 8.507s |

本会话未重跑串行 `go test ./...` 或 Web `npm test`。F-010 的必要证据是认证链 fail-closed 与事务回滚，不依赖全量回归复述。

## 工作区、信息门禁与受影响门禁

| 检查项 | 结论 | 证据 |
|--------|------|------|
| 工作区绑定 | 通过 | `workspace.md` `id` / `root_goal` / `canonical_scope` 与 Root `parent: null`、路径一致 |
| `plan_refs` / `primary_plan` | 通过 | workspace 与 Root 均挂 `VP-012-shared-cross-module-contracts` |
| 共享资料 | 无引用 | `shared_materials_catalog: none` |
| I-001 | 维持 verified | non-blocking；不构成本条阻断 |
| I-002 | 维持 verified（本 scope） | required；最晚阶段 = Root 关门。A-010 曾以 F-010 否定「R1～R6 全部合法闭合」的无条件主张。本条独立核对后，该阻断已消除。无 `deferred` required，无用户书面 residual/overrule |
| R6 成功标准 3 | 本条范围内满足 | 创建/吊销事务审计既有；使用审计现为 fail-closed 且与 `last_used_at` 同事务；secret/hash/header 不进入 `ServiceCredentialUse` 或 503 错误体 |

## 成果（有证据）

A-010 F-010 原文指出：`auth.go` 以 `_ =` 丢弃 `MarkServiceCredentialUsed` 与 `serviceCredentialUseRecorder` 错误，并无论成败都调用 `next.ServeHTTP`。现行代码中这两处丢弃已不存在。

| 主张 | 独立结论 | 证据 |
|------|----------|------|
| 生产路径不再 best-effort 放行 | **成立** | `apps/api/internal/composition/composition.go:242-258` 只安装 `SetServiceCredentialUseTransactionalRecorder`。`apps/api/internal/auth/auth.go:620-633`：事务 recorder 存在时走 `MarkServiceCredentialUsedWithAudit`；类型断言失败或函数返回错误均写 503 `STORAGE_UNAVAILABLE` 并 `return`，之后才 `next.ServeHTTP`。 |
| 使用审计与 `last_used_at` 原子提交 | **成立** | `apps/api/internal/modules/authsession/service_credentials.go:180-195`：同一 `withTx` 内先 UPDATE `last_used_at`，再调用 audit；audit 错误回滚。`RecordOperationTx` 使用调用方 `tx`（`apps/api/internal/modules/operationlog/repository.go:151-157`），`SetOperationLogError` 可注入失败。 |
| 兼容旧 seam 亦 fail-closed | **成立** | `auth.go:635-645`：非事务 recorder 或 `MarkServiceCredentialUsed` 任一错误同样 503 并 return。此路径不是生产 composition 热路径。 |
| 策略已书面冻结为 required 持久化 | **成立** | 区内 `GOAL-007-r6-api-token-service-credential/01-decision/D-004-r6-use-audit-fail-closed.md`：用户确认 `fixed`，拒绝 residual/overrule；未采用 D-003 best-effort。 |
| 故障注入证明 handler 不执行且元数据不落盘 | **成立** | `auth_test.go:68-95` 事务 recorder 失败 → 503、`called=false`、`LastUsedAt==nil`。`authsession/service_credentials_test.go:99-119` 仓库层回滚。`handler/service_credentials_test.go:186-195` 真实 HTTP 组装 + `SetOperationLogError` → 503 `STORAGE_UNAVAILABLE`、`LastUsedAt==nil`。本轮上述包测试均通过。 |
| 成功路径仍写使用审计且不泄露 secret | **成立** | `handler/service_credentials_test.go:87-105`：存在 `service-credentials.use`；detail 含 `scopeCount`；detail 不含 raw secret。`ServiceCredentialUse` 字段只有 credentialId/name/scopeCount/method/path/correlation/at。 |

## 对照成功标准（本 scope）

| 标准 | 状态 | 证据 |
|------|------|------|
| A-010 F-010 关闭声明可重复核对 | 已达成 | 原丢弃路径已删除；生产事务 fail-closed；本轮定向测试通过 |
| R6 S3「使用审计可靠存在」 | 已达成（本条） | 失败返回 503 且不进入下游；成功写入 `service-credentials.use` |
| Root I-002 受 F-010 阻断 | 已解除 | 见上表；I-002 在本 finding 范围内维持 verified |

## Findings

本条无新的 required / recommended finding。

### A-010 F-010 闭合判定

| 字段 | 值 |
|------|-----|
| 原 finding | A-010 F-010 · R6 使用审计失败时仍放行请求 |
| level | required / medium |
| A-011 主张 | fixed |
| 本轮独立结论 | **同意 fixed** |
| closure | P-003 `fixed`：可核对修正 + D-004/E-008/E-011/A-011 留痕 + 本条独立复验 |
| 观察（非 finding） | `TestServiceCredentialMiddlewareFailsClosedWhenUseMetadataFails` 覆盖的是旧 `MarkServiceCredentialUsed` seam，不是生产 `MarkServiceCredentialUsedWithAudit` 的 UPDATE 注入。生产 UPDATE 失败与审计失败共用同一 `err != nil → 503` 分支，且仓库层/HTTP 组装已覆盖审计失败回滚。此测试缝差异不否定 F-010 闭合。 |

## 必改项汇总

无。开放 required = 0。

## 与既有意见的异同

- **A-010**（independent / fail）：原文与 F-010 保持不变。当时证据针对 `_ =` 丢弃错误；该代码事实已被本轮源码核对否定。
- **A-011**（self / pass）：主张 F-010 `fixed`、用户选 `fixed`、生产事务 fail-closed。本条同意其闭合结论；不把 self 记录当作独立证据本身，而以现行代码与本轮复跑测试为据。
- 不与其他 Root 意见冲突。本条不重开 A-004～A-009。

## 结论 + 建议给编排器/用户的下一步

A-010 F-010 的关闭证据充分、可重复核对。R6 使用审计失败不再放行下游请求；生产路径下审计与 `last_used_at` 同事务提交或回滚。本 scope verdict 为 **pass**。

建议 `/govern` 记录已接收本独立复审；无需再为 F-010 选择 residual/overrule。本条不修改 Root `status`/`progress`。

## 声明

本意见为 `source: independent`，只追加审计意见，不修改目标 `status`、`progress`、`goal-tree` 或既有审计原文。响应由 `/govern` 处理。
