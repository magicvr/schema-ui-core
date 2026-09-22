---
id: A-010-w18-post-close-independent
doc: audit-entry
goal: GOAL-019-w18-api-web-security-hardening
date: 2026-09-22
source: independent
auditor: Grok 4.7 / grok-build / independent
type: close-out
scope: W18 关门复核：API JWT fail-closed、下游 MFA 启动探针、Web 跨 realm 凭据边界、A-005 闭合与用户关门记录
verdict: pass
open_required: 0
parent: GOAL-001-production-hardening
version: 0.1.0
---

# A-010 · W18 关门后独立复核（2026-09-22）

- **source**：independent
- **auditor**：Grok 4.7 / grok-build
- **类型** / **scope**：close-out · 已关门 GOAL-019 的 required 修复、A-005 闭合证据与 D-004 用户关门记录
- **verdict**：pass
- **开放 required**：0

## 范围与区间

工作区 `workspace-009-production-hardening`（`workspace.md`：`root_goal` = `GOAL-001-production-hardening`，`canonical_scope` = 本区目录，`shared_materials_catalog: none`，`plan_refs` / `primary_plan` = `VP-009-production-hardening`）。被审目标在 canonical 范围内，`parent` = Root。本区无固定共享资料引用。

本意见只复核 W18 已冻结并宣称闭合的边界，以及关门记录是否与代码和台账一致。范围是：

- D-002 API MAJOR：`server.Run` 对直接传入的配置 fail-closed，非 development 不得使用公开固定 JWT
- A-002 F-001：下游 `serve` 未装配 MFA verifier 时，active enrollment 或探测失败必须拒绝启动
- D-002 / A-002 F-002 / A-004 F-001：`authFetch` 不得把 Bearer / refresh 附到跨源请求，包括跨 realm `Request` / `URL`
- A-005 F-001 / F-002 的闭合证据，以及 D-004 / A-009 的用户关门记录
- P-005：I-001～I-004 的最晚阶段、状态与证据

未把本波已接受的 development fallback、未跑生产 PostgreSQL、以及启动后才写入的共享库 enrollment 重新升级为新的 required。`workspace.md` 波次表停在 W15，状态权威以本区 `goal-tree.md` 为准；该表漂移不改变本目标的绑定字段。

## 成果（有证据）

### API JWT

`apps/api/server/serve.go` 的 `Run` 在打开 store、装配认证和监听之前调用 `cfg.validate()`。`apps/api/server/config.go` 的 `validate` 在 `AppEnv != "development"` 时拒绝空 `AUTH_JWT_SECRET`，并调用 `ValidateJWTSecretStrength`。`resolveSecret` 只在 `development` 返回 `dev-only-insecure-jwt-secret-change-me`，其他环境返回空字符串。`apps/api/cmd/schema-ui/main.go` 的 `cmdServe` 经 `server.Serve` 进入这条路径。主进程 `apps/api/cmd/server/main.go` 的 `resolveJWTSecret` 是另一条已有 fail-closed 路径，本波没有改写它。

本轮复跑：`apps/api` 下 `go test ./server/ ./modules/mfa/store/ -count=1`，两个包均为 `ok`。其中包含 `TestRunRejectsNonDevelopmentJWTSecretBeforeStartup` 与 `TestResolveSecretKeepsDevelopmentFallbackOnly`。

### 下游 MFA 启动探针

`Run` 在 listener 之前调用 `mfastore.NewRepository(st).HasActiveEnrollment()`。查询失败或存在 `user_mfa.status = 'active'` 时关闭 store 并返回错误。`compiled.PersistenceCatalog()` 包含 `mfamigration.Provider`。`Authenticator.Login` 只在 `a.mfa != nil` 时拦截第二因素；下游 `RegisterWithMFAProbes` 传入的 verifier 为 nil，因此启动探针是这条装配的控制点。空表与 pending 由 `TestHasActiveEnrollment` 锁定为非 active。

### Web 跨 realm 凭据

`apps/web/src/account/auth-client.ts` 的 `resolveTargetUrl` 对字符串使用 `new URL(input, origin)`，对对象只读 `url` / `href` 字符串，不可解析或读取抛错时返回 null。`isSameOrigin` 为假时不写 `Authorization` 与 `X-Refresh-Token`。

本轮复跑：

- `apps/web`：`npm run test -- src/account/auth-client.test.ts`，34 passed
- `apps/web`：`npx playwright test e2e/auth-client-cross-realm.spec.ts --project=chromium`，1 passed。该用例使用 iframe 的 `Request` / `URL` 构造器，并断言同源 session-list 带两项凭据、跨源 `Request` 与 `URL` 不带、同源 `/api/auth/login` 的 401 不触发 refresh

### 治理与信息项

| 项 | 核对 |
|----|------|
| I-001 | verified。非 development 缺密钥 / 弱密钥在 `Run` 启动前失败，证据见 `serve.go`、`config.go` 与上述 `go test` |
| I-002 | verified。跨 realm 凭据边界有实现、Vitest 与本轮 Chromium 复跑 |
| I-003 | verified。本轮复跑了区分性 API 包、`auth-client` 单测和指定 Chromium spec。A-007 记录的全量 `go test ./...`、Web 124/1516 与 `npm run typecheck` 本轮未整包重跑 |
| I-004 | verified。A-007 `pass` 且 open_required=0；A-008 将 A-005 F-001/F-002 记为 `fixed`；D-004 为用户书面 `ok，done`；本轮未发现新的 required |
| A-005 F-002 哈希 | `A-002` SHA256 `250E9356426AEF700887B88CE2356CA118A5AAFECE62EC4E723294B1A97A8139`，`A-004` SHA256 `02FC8825FFB8A0E5FBDAF73C3892A204269613EDE500ECB6F98199B6325593C7`，与 E-005 登记值一致 |
| 关门记录 | `status: done`，goal-tree 主树与状态表均为 done · 6/6，D-004 / A-009 / E-007 一致。Root 保持 active |

无到期未关闭的 required 信息项，无待用户接受的 residual。

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| 扫描范围、基线与证据已落盘，事实与推断分开 | 已达成 | E-002、D-002 |
| required 范围、非目标与信息门禁已冻结 | 已达成 | D-002；I-001～I-004 均有状态与证据 |
| 本波 required 安全问题已修复并有回归 | 已达成 | 上文三处代码与本轮复跑 |
| 验证与风险相称，未把既有失败写成新成果 | 已达成 | 区分性测试本轮通过；全量分母保持为 A-007 的历史记录 |
| self 与 independent 意见已落盘，required 按三路径闭合 | 已达成 | A-001～A-009；A-005 两条为 `fixed`；D-004 书面关门 |
| 目标与 goal-tree、执行台账、审计台账一致 | 部分 | 状态、进度、D-004 与 goal-tree 一致。路线图 S6 旁注见 F-001 |

## Findings

### F-001 · 路线图 S6 旁注仍写用户关门待定

- 严重度：low
- 建议：recommended
- 状态：open

`00-meta.md` 路线图第 6 点已经勾选，但括号内仍是「A-007 pass；用户关门裁决待定」。同文件 frontmatter 为 `status: done`，D-004、A-009、E-007 与 `goal-tree.md` 都已记录用户书面关门。旁注会让路线图读起来像关门尚未发生。

这不推翻 `done`，也不构成未闭合 required。编排器宜只改正这句旁注，使之指向 D-004。

## 必改项汇总

无。开放 required = 0。

## 与既有意见的异同

与 A-007 的代码结论一致：JWT 启动校验、active MFA 启动探针、跨 realm 凭据边界均成立，A-005 F-001/F-002 可以维持 `fixed`。A-007 当时目标仍是 `active`，S6 旁注写「用户关门裁决待定」与当时事实相符；D-004 之后该句成为过期叙述，故本轮新增 recommended F-001。

不改写 A-002 / A-004 / A-005 / A-007 的 source、verdict 或 open_required。

## 结论与建议

`verdict: pass`。W18 宣称闭合的三项安全边界与当前代码、区分性回归和用户关门记录一致。没有新的 required finding，不需要重新打开目标。

建议编排器用 `/govern` 响应本意见：记录 A-010，并修正 `00-meta.md` 路线图 S6 旁注。不要因本意见改 `status` 或派生 progress。

本轮没有复跑全量 `go test ./...`、Web 全量 Vitest、`tsc` 或完整 Playwright，也没有连接生产 PostgreSQL。这些缺口不形成开放 required。

## 声明

本意见为 `source: independent`，不修改目标 `status`、检查点、派生 `progress`、方案正文或 goal-tree。响应与旁注修正由 `/govern` 处理。
