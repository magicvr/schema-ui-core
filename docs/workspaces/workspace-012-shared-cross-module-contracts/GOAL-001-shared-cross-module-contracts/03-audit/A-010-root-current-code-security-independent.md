---
id: A-010-root-current-code-security-independent
goal: GOAL-001-shared-cross-module-contracts
doc: audit-entry
record_id: A-010
source: independent
auditor: Codex independent audit
scope: workspace-012 close-out; current R1-R6 implementation, regression evidence, bug and security review; focused recheck of R6 service-credential use auditing
verdict: fail
status: recorded
parent: null
created: 2026-08-19
updated: 2026-08-19
version: 0.1.0
---

# A-010 · 独立代码与安全审计 · workspace-012 当前实现（2026-08-19）

## 范围与方法

本意见审查 `docs/workspaces/workspace-012-shared-cross-module-contracts/` 的 Root `GOAL-001-shared-cross-module-contracts` 及其 R1-R6 代码交付。工作区治理文件只用于识别声称的范围、成功标准和历史审计，不作为代码成功依据。

已核对：

- API：`go test ./...`，全部包通过；其中 `internal/handler` 通过，耗时 333.086s。
- Web：`npm test -- --run`，72/72 test files、1069/1069 tests 通过。
- R1 request-id、R3 wallet CAS/idempotency、R4 Job lease/state machine、R5 operational gate、R6 service-credential authentication/management 的现行源码与局部测试。
- R6 创建/吊销事务审计与使用审计的不同错误处理路径。

## 成果（代码证据）

1. R1 request-id 仅接受受限字符集和 128 字符长度，并对无效输入生成新值，未发现基于该路径的 header injection 证据（`apps/api/internal/requestid/requestid.go:24-29,53-66`）。
2. R3/R4/R5 的关键状态、CAS/idempotency、lease/终态与已注册 mutation gate 均能在现行代码中定位；本轮全量 API/Web 回归未发现新的可直接复现失败。
3. R6 管理创建/吊销使用事务审计，审计失败会使领域变更失败或回滚（`apps/api/internal/handler/service_credentials.go:190-207,218-240`；`apps/api/internal/modules/authsession/service_credentials.go`）。

## Findings

### F-010 · required · medium · R6 使用审计失败时仍放行请求

- **主张/影响门禁**：GOAL-007 R6 成功标准 3 要求创建、使用、吊销审计含 actor/credential/correlation 且不泄露 secret；Root I-002 与 close-out 要求 R1-R6 合法闭合。该 finding 影响 R6 S3 与 Root close-out。
- **证据**：`apps/api/internal/auth/auth.go:594` 以 `_ = a.repository.MarkServiceCredentialUsed(credential.ID, now)` 丢弃使用时间持久化错误；`apps/api/internal/auth/auth.go:595-604` 对 `a.serviceCredentialUseRecorder(...)` 同样以 `_ =` 丢弃错误；`apps/api/internal/auth/auth.go:606` 无论上述失败与否都执行 `next.ServeHTTP(...)`。
- **链路**：生产组合根的 recorder 将使用事件写入 operation log（`apps/api/internal/composition/composition.go:241-256`）。operation-log repository 明确可返回持久化/事务错误（`apps/api/internal/modules/operationlog/repository.go:142-156`）。因此数据库不可用或审计写入失败时，有效 service credential 仍可访问业务 API，而对应 `service-credentials.use` 事件缺失；`last_used_at` 也可能缺失。
- **验证缺口**：现有 R6 测试验证成功使用审计和创建/吊销失败回滚，但没有对“使用审计失败后 handler 不得执行”或明确 best-effort 监控策略的回归测试。`SetOperationLogError` 可作为故障注入入口（`apps/api/internal/modules/operationlog/repository.go:152-157`）。
- **风险**：这是审计完整性与安全可追责性缺陷，不是凭据本身的权限绕过；在日志存储故障期间，机器凭据调用会形成不可见的成功操作，削弱事件调查、撤销后核验和合规审计。
- **建议**：由 `/govern` 选择并落盘一种明确策略：若“使用审计”是 required，则在 recorder/必要持久化失败时返回 5xx/503 且不调用 downstream，并以事务或 durable outbox 保证凭据使用与审计的一致性；若产品确实接受 best-effort，则必须由用户书面记录 `accepted-residual` 的范围、期限、监控和复审触发，并相应收窄成功标准。无上述闭合前，F-010 保持开放 required。

## 对照成功标准与结论

- **实现路径**：R1-R6 均能定位到源码/测试，且本轮 API/Web 全量回归通过。
- **真实消费与装配**：现行 wallet/jobs/operationlog/route composition 路径存在；本轮未发现 Profile、Manifest 或协议资产被测试直接击穿。
- **关门结论**：不能确认 workspace-012 已“确实完成”。F-010 直接否定 R6 “使用审计可靠存在”的无条件主张，因此本 scope verdict 为 **fail**，而不是 pass/conditional。

## 必改项汇总

| finding | 级别 | 状态 | 影响 |
|---|---|---|---|
| F-010 | required / medium | open | R6 S3、Root close-out |

## 声明

本意见为 `source: independent`，只追加审计意见，不修改目标 `status`、`progress`、`goal-tree` 或既有审计原文。请使用 `/govern` 响应 F-010，完成修正、用户书面 residual/overrule 或后续独立复审。
