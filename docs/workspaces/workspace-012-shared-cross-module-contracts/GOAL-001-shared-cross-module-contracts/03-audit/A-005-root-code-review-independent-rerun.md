---
id: A-005-root-code-review-independent-rerun
goal: GOAL-001-shared-cross-module-contracts
doc: audit-entry
record_id: A-005
source: independent
auditor: Codex（独立交叉审计，/audit 入口）
scope: GOAL-001 R1 through R6 actual code completion, reproducible validation, bugs and security
audit_type: close-out + ad-hoc
verdict: fail
status: recorded
parent: GOAL-001-shared-cross-module-contracts
created: 2026-08-19
updated: 2026-08-19
version: 0.1.0
---

# A-005 · 独立代码审查复核：完成门禁、缺陷与安全

- **source**：independent
- **auditor**：Codex（独立交叉审计，/audit 入口）
- **类型**：close-out + ad-hoc
- **scope**：GOAL-001 下 R1～R6 的实际代码完成度、可重复验证、bug 与安全漏洞
- **verdict**：fail
- **日期**：2026-08-19

## 范围与方法

本复核以代码、测试命令输出和可定位的实现为依据；工作区治理文件只用于确定 R1～R6 的审计范围，不作为成功证据。检查了 API 的 request-id、审计 detail、并发/幂等、Job、运行态门控、service credential 及其消费路径，检查了 Web schema/i18n 消费路径，并执行 API 全量测试、Web 全量测试和生产构建。

## 成果（有证据）

- API 在 `apps/api` 模块目录执行 `go test ./...` 通过；R1、R3、R4、R5、R6 均存在真实实现和定向测试，R2 的 auth/settings/users mutation 均调用 `operationlog.NewDetail`。
- Web `npm run build` 通过，说明 TypeScript/Vite 构建链可完成；这不等同于运行时/结构化契约测试全部通过。
- Service credential 认证按 SHA-256 hash 查找，拒绝 revoked/expired 凭据，scope 受既有 permission key 和创建者权限上限约束；创建响应才返回 raw secret，持久层与审计投影不保存 raw secret。
- 上传/下载路径有 ID 形状校验、owner 检查、内容嗅探、危险类型拒绝、`nosniff`、固定下载名和 HTTP 超时配置。本轮静态检查未确认认证绕过、路径穿越、SQL 注入或凭据明文泄露。

## Findings

### F-008 · Web 全量结构化 i18n 测试失败（required，medium）

- `apps/api/internal/modules/systemmonitoring/schema/system-monitoring.json:58-66` 使用 `schema.systemMonitoring.statCard.availability`。
- `apps/web/src/i18n/messages/en-US.json:287-294` 和 `apps/web/src/i18n/messages/zh-CN.json:287-294` 均未定义该 key。
- `apps/web/src/i18n/schema-keys.structural.test.ts:122-132` 会对该缺失分别报 en-US/zh-CN 错误。复现命令：在 `apps/web` 执行 `npm test -- --run`，结果为 **72 个测试文件中 1 个失败、1069 个测试中 1 个失败**，失败项为该结构化 key 完整性断言。

这直接违反 R5 成功标准中的“定向与全量测试通过”，也使 Root 的全量验证/关门主张不能从当前代码得到确认。修复前不能把 R5 或 Root 视为已完成；应补齐两种语言的 availability 文案（或改用已存在的 key），再重跑完整 Web 测试。

### F-009 · Job heartbeat 的取消查询错误会遗留 running Job（recommended，medium）

`apps/api/internal/jobs/runner.go:287-291` 在 `IsCancelRequested` 返回错误时直接 `return`，没有调用 `reporter.cancelExecution()`，也没有失败/取消当前 lease。handler goroutine 可能继续运行，而 Job 会保持 `running` 直到 lease 过期并被恢复扫描处理，造成延迟完成和重复恢复窗口。应在该错误路径取消 handler，并明确 lease 的失败/回收策略；新增数据库错误注入测试覆盖此行为。

## 安全审计结论

本轮静态审计未确认 critical/high 安全漏洞。已核对的保护包括 service credential hash-only 生命周期和 scope ceiling、上传内容/路径边界、文件 owner 访问控制、响应 `nosniff`/CSP、生产 JWT/dev-session 配置门槛、前端 return-intent 校验和 HTTP 超时。未执行动态渗透、模糊测试、依赖 CVE 扫描或生产反向代理/TLS 审查，因此不能将本结论解释为不存在安全漏洞的证明。

## 必改项汇总

- **F-008 required**：修复 system-monitoring availability i18n key，并使 Web 全量测试通过。
- **F-009 recommended**：修复 Job runner heartbeat 错误路径的 handler/lease 清理，补充回归测试。

## 结论与下一步

R1～R6 的主要 API 代码确实存在，且 API 测试与 Web 构建通过；但当前 Web 全量测试有可重复失败，因此“工作区 12 已确实完成”的关门结论不成立，verdict 为 **fail**。请使用 `/govern` 响应 F-008（建议 `fixed` 后复审）；F-009 可选择 `fixed` 或按范围和复审触发条件书面接受 residual。

## 声明

本意见不修改目标 `status`、`progress` 或 `goal-tree`；响应与状态推进由 `/govern` 处理。
