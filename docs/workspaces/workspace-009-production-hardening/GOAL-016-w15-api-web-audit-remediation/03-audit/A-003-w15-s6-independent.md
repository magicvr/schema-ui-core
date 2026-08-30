---
status: active
created: 2026-08-30
updated: 2026-08-30
parent: GOAL-016-w15-api-web-audit-remediation
version: 0.1.0
---

# A-003 · W15 S6 关门前 independent 复核（grok build）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high · `/audit`）
- **日期**：2026-08-30
- **类型**：close-out / finding-closure / execution-facts
- **scope**：A-001 分母 required F-001～F-006 + recommended F-007 是否 genuine-fixed；对照源码与测试，不轻信台账；独立复跑 API/Web 回归。实现 checkpoint = `609cd6d6`。
- **verdict**：**pass**
- **工作区**：`workspace-009-production-hardening`（Root `GOAL-001-production-hardening`；canonical `docs/workspaces/workspace-009-production-hardening/`；`vision_role: delivery`；`primary_plan` = `VP-009-production-hardening`；`shared_materials_catalog: none`）
- **被审 HEAD**：`e20391a4`（A-002 文档提交）。`609cd6d6..HEAD` 仅治理文档（11 文件、无实现 diff）；实现切片仍为 `609cd6d6`。

## 范围与区间

- **覆盖**：D-002 冻结方案相对 A-001 分母的代码闭合；E-001～E-003 实施主张；A-002 self 逐条核验是否可重复；I-001～I-005 对 S6 的信息门禁。
- **方法**：只读工作区 + GOAL-016 五件套 → 点名路径源码/测试抽验（含调用点）→ 本会话独立复跑（命令与结果见下）。未做动态 exploit、未起 compose、未改密。
- **不覆盖**：不改 `status` / `progress` / 检查点 / 方案正文 / goal-tree。不自行闭合 finding、不恢复 VP-008 `go`、不把本意见当作已关门。
- **排除**：R-001（refresh token `localStorage` residual）按 D-001 不重开。A-001 的 nginx referrer 建议按 D-002 不在本波。

## P-005 / 工作区核对

| 核对项 | 结论 |
|--------|------|
| 工作区绑定 | `workspace.md`：`id=workspace-009-production-hardening`；Root `GOAL-001-production-hardening`；canonical 与 goal-tree 一致；`vision_role: delivery`；`primary_plan` = `VP-009-production-hardening`。Charter `schema-ui-core-admin-foundation@0.3.0`。共享资料目录 `none`，本 scope 未把资料引用当关闭证据。未读取其他工作区作为关闭依据。 |
| I-001 provider | **verified**（D-002 §2）：grok build · grok-4.6 · high · `/audit`。本条即该腿。 |
| I-002 serve 边界 | **verified**（D-002）：公共 serve 是受支持下游入口；默认回环 + 显式才对外。 |
| I-003 bootstrap 策略 | **verified**（D-002）：冻结 8–72 字节非空；fresh 库与 0057 播种默认等价。 |
| I-004 fixture 根 | **verified**（D-002）：canonical = `apps/api/modules`。 |
| I-005 LocalStore 威胁模型 | **verified**（D-002 用户书面 **fixed** 0700/0600）。 |
| 到期 required 信息项 | 无。S6 无开放 required 信息门禁。 |
| 共享资料 | 无 |

索引文件 `03-audit.md` 顶部信息表在本条写入前仍写 I-001～I-005 `open` / 目标 `draft`，落后于 D-002/00-meta（见 note N-002）。不以滞后索引否定已 verified 的信息项。

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| 6 required + 1 recommended 均有可核对代码改动，与 D-002 对齐 | 见逐条闭合判定；无偷换范围 |
| A-002 回归主张可重复核对 | 本会话独立复跑（见下） |
| 实现切片未被后续文档提交改写 | `git diff --stat 609cd6d6..HEAD` 仅 `docs/workspaces/workspace-009-production-hardening/**` |
| 未引入阻断性新缺陷 | 源码抽验未见把 fail-closed 改回 fail-open；未见 MFA 登录 CAS 回退；未见 fixture 根回退 |

## 对照成功标准

| 标准 | 本 scope | 状态 | 证据 |
|------|----------|------|------|
| S1 立项与 A-001 落盘 | 前序 | 达成 | A-001 + D-001 |
| S2 方案冻结 + I-001～I-005 | 前序 | 达成 | D-002 |
| S3 API 修正 + 回归 | 前序 + 本条复跑 | 达成 | E-001；`go vet` 0；W15 包与全量 `go test` |
| S4 Web 修正 + 回归 | 前序 + 本条复跑 | 达成 | E-002；`tsc -b` 0；vitest 89/1183；`vite build` 0 |
| S5 F-007 + 全量验证 | 前序 + 本条复跑 | 达成（Windows 上 POSIX 权限测试按设计 skip，见 N-001） | E-003；`local.go` 0700/0600 |
| S6 self + independent 复核 | **本条判定代码闭合条件已满足** | 达成（意见层） | A-002 self pass + 本条 independent pass。合法闭合与关门由 `/govern` 响应 |

## 独立复跑（本会话真实执行）

工作目录：仓库根。日期：2026-08-30。主机：Windows。

| 命令 | 结果 |
|------|------|
| `Set-Location apps/api; go vet ./...` | **exit 0**（0 输出） |
| `Set-Location apps/api; go test ./...`（首次，走 cache） | 全包 `ok` / `?`，**exit 0**（cache hit，不作独立执行证据） |
| `Set-Location apps/api; go test ./... -count=1`（第一次强制执行） | **exit 1**：`TestShutdownDrainHarnessPostgres`（`internal/composition`）`postgres in-flight request failed during drain: EOF`。该测试属 VP-021 停机 harness，**不在 W15 分母**；当时 web vitest 并行。见 N-003。 |
| `go test ./internal/composition -count=1 -run TestShutdownDrainHarnessPostgres` | **ok** 6.411s（隔离复跑通过） |
| `Set-Location apps/api; go test ./... -count=1`（第二次，无 web 并发） | **exit 0**，全包 `ok`/`?`，0 FAIL。`internal/composition` 33.145s ok |
| W15 定向：`./server` `./cmd/server` `./internal/config` `./modules/authsession` `./modules/mfa` `-count=1` | 全部 **ok** |
| `go test ./internal/objectstore -count=1 -v -run TestLocalPutTightenedUnixPermissions` | **SKIP**：`POSIX file modes are not enforced on Windows` |
| `Set-Location apps/web; npx tsc -b` | **exit 0** |
| `Set-Location apps/web; .\node_modules\.bin\vitest.cmd run` | **89 passed (89) / 1183 passed (1183)**，Duration 11.96s |
| `Set-Location apps/web; .\node_modules\.bin\vite.cmd build` | **exit 0**；既存 chunk >500 kB 警告（`index-BzLQch-0.js` 795.25 kB） |

结论：W15 分母相关包与 Web 全量回归可重复为绿。全量 API 无缓存套件在无并发干扰下为绿。第一次 `-count=1` 的 composition PG drain 失败不构成 W15 实现回退。

## 逐条闭合判定（A-001 分母）

### F-001 · 默认暴露面收紧 → **genuine-fixed**

| 核对项 | 路径 | 结论 |
|--------|------|------|
| 代码默认 | `apps/api/server/config.go:109` | `HTTPAddr: "127.0.0.1:25080"` |
| 空 `APP_ENV` | `server/config.go:263-265` `validate()` 首条 | 空串 fail-closed（`refusing to guess`）；不再把空 env 当 development |
| 内嵌默认 YAML | `server/config.default.yaml:11` | `addr: "127.0.0.1:25080"`；`app.env: ${APP_ENV:-development}` 仍显式 pin development |
| create 骨架 | `cmd/schema-ui/templates/config.yaml.tmpl:12` | 同步回环 |
| `-addr` 显式放开 | `cmd/schema-ui/main.go:181-201` | 空 `-addr` 不覆盖；非空在 LoadConfig/validate **之后**写入，符合 D-002「显式才对外」 |
| dev 回退 | `serve.go` `resolveSecret` / `bootstrapAdmin` | 缺密钥仍回退 `dev-only-insecure-jwt-secret-change-me`；缺密码仍回退 `admin`。仅 development 路径可达（非 dev 已由 validate 强制密钥/密码） |
| 负例 | `config_test.go` `TestLoadConfigDefaults` / `TestLoadConfigRequiresExplicitAppEnv` | 默认回环断言；省略 `app.env` 的自定义 YAML 拒绝；显式 development 可装载 |

主仓 `internal/config` 仍默认 `:25080`（`cmd/server`），**不在 A-001/D-002 分母**（A-001 已区分主服务 `ValidateProd`）。compose 不发布宿主端口仍是 W7 边界。不把该残余升格为本波 required。

结构体注释 `AppEnv string // "" = development（缺省）` 与现行 `validate()` 不一致（见 N-002）。

### F-002 · 非 development JWT secret 强度 → **genuine-fixed**

| 核对项 | 路径 | 结论 |
|--------|------|------|
| 单一来源 | `internal/config/config.go:1517-1533` | 导出 `ValidateJWTSecretStrength`：`len ≥ 32` + `containsLettersAndDigits`；错误只带 key 名、不带密钥值 |
| `ValidateProd` 复用 | `config.go:1115-1127` | 当前密钥与 previous 密钥均调用同一函数 |
| 公共 serve | `server/config.go:303-306` | 非 development 调用同一函数 |
| 无第二份规则 | grep `minJWTSecretLen` | 仅 `internal/config` 定义 |
| 负例（server） | `TestLoadConfigJWTSecretStrengthNonDev` | short / all-letter / all-digit 启动失败；`0123456789abcdef`×2（32 混合）通过 |
| 负例（主仓） | `internal/config/config_test.go` ValidateProd 短/纯字母/纯数字 | 既有生产负例仍在 |

### F-003 · bootstrap 种子策略 → **genuine-fixed**

| 核对项 | 路径 | 结论 |
|--------|------|------|
| 策略函数 | `modules/authsession/password_policy.go:100-105` | `ValidateSeedPassword`：8–72 字节、`TrimSpace` 非空；复用 `policyMinLengthFloor/Ceiling` |
| 主服务入口 | `cmd/server/main.go:100-117` | 非 development 先策略再 bcrypt；空种子直接 fail；development 保留 `admin` 回退 |
| 公共 serve | `server/serve.go:272-290` | 非 development 先策略再 bcrypt；空种子先填 `admin` 再被 5 字节策略拒绝（validate 已强制非空，属纵深） |
| 策略负例 | `password_policy_test.go` `TestValidateSeedPassword` | `""` / 空白 / 7 / 73 拒绝；8 / 72 / 普通合规通过 |
| 启动负例 | `cmd/server/main_test.go` `TestResolveSeedHashPolicy`；`serve_test.go` `TestRunRejectsWeakSeedPasswordNonDev` | production 缺省/weak/7/73 失败；compliant 通过；dev 回退保持；4 字节种子在监听前失败（错误含 `bootstrap seed password`） |

与 A-002 S-001 一致：config 层只强制非空，策略在 bootstrap 层咬合——D-002 §3 写明的分层，不是漏洞。HTTP 改密仍走 `ValidateNewPassword`（可配置策略行），不重叠。

### F-004 · MFA step-up TOTP CAS → **genuine-fixed**

| 核对项 | 路径 | 结论 |
|--------|------|------|
| 公共 helper | `modules/mfa/service.go:321-354` `requireActiveSecondFactor` | `ValidateTotp` 成功取匹配 step → `AdvanceLastUsedStep`；`advanced=false` → `ErrMFAInvalid`；`maybeRewrap` 仅 CAS 赢者 |
| 调用方 | `Disable` `:269`；`RotateRecovery` `:277`；`VerifySecondFactor` `:369` | 三处共用同一 helper，无旁路 |
| CAS SQL | `modules/mfa/store/repository.go:198-219` | `UPDATE … WHERE last_used_step < ?`；`affected == 1` 才算赢 |
| 重放测试 | `service_test.go` `TestServiceStepUpTotpReplayRejected` | 同码二次 `RotateRecovery` → `ErrMFAInvalid`；`VerifySecondFactor` 重放拒绝；下一 step 新码仍可用 |
| 仓库并发守卫 | `repository_test.go` `TestAdvanceLastUsedStepGuardedCAS` | 同 step 第二次 `advanced=false`（登录路径同语义） |

`Disable` 未单独写重放用例，但无独立实现。Confirm 入学仍用 `SetLastUsedStep`（W13 F-004 匹配步进），不在本条分母。

### F-005 · 邀请 token URL 清理 → **genuine-fixed**

| 核对项 | 路径 | 结论 |
|--------|------|------|
| 读取后 scrub | `apps/web/src/components/invite-accept.tsx:45-58` | `useState` 初始化：`URLSearchParams.get("token")` 后立即 `searchParams.delete("token")` + `history.replaceState`；其它 query 保留 |
| 提交路径 | 同文件 `:79-80` / `:98` | token 只存在 React state；成功后 `window.location.href = "/"` 不变 |
| StrictMode | `main.tsx` 包 `<StrictMode>`；React 19 对 lazy init **第二次调用结果丢弃** | 第一次读 token 并 scrub；第二次 URL 已无 token（幂等）；state 保留第一次返回值。A-002「第二次挂载 query 已无 token → 幂等」对 URL 副作用成立 |

无组件级 `replaceState` 回归测试（见 F-008）。D-002 明确不改 nginx referrer；页面本身也无 `no-referrer` meta——属分母外残余，不重开 required。

### F-006 · canonical fixture 根 → **genuine-fixed**

| 核对项 | 路径 | 结论 |
|--------|------|------|
| 退役路径 | grep `api/internal/modules` under `apps/web` | **仅** `fixture-root.guard.test.ts` 自身描述/扫描字符串（guard 豁免 `SELF`） |
| 13 suite 切换 | `all-module-schemas-dval`、`custom-components.schema`、`error-localization`、`row-action-bindings`、`ui-bilingual`、`schema-keys.structural`（fragment ×10）、`s5-denominator-render`、`representative-pages`、`schema-dictionary-entries`、`schema-crud`、`wallet-navigate`、`startup-config`、`representative-pages.integration` | 全部 `../../../api/modules` 或等价 `apps/api/modules`；`load-page.test.ts` 仅注释 |
| guard | `src/protocol/fixture-root.guard.test.ts` | canonical 根存在 + `src` 无退役路径引用。本会话 vitest 该文件 **2 tests passed** |
| README | `apps/web/README.md:69` | `apps/api/modules/` + `dev/examples/schema`；无 `internal/modules` |
| 回归 | vitest **1183/1183**（基线 A-001：76 failed / 1157） | 质量门禁解除 |

### F-007 · LocalStore 0700/0600 → **genuine-fixed**（recommended；用户裁决 = fixed）

| 核对项 | 路径 | 结论 |
|--------|------|------|
| `Put` | `internal/objectstore/local.go:119-151` | `MkdirAll(..., 0o700)`；tmp `Chmod 0o600` 后 rename；sidecar `WriteFile(..., 0o600)` |
| 既有文件 | 同文件注释 | 不强制改写（D-002：无迁移 churn）。旧 0644 对象直到下次 Put 才收紧——书面范围内残余 |
| 测试 | `local_test.go` `TestLocalPutTightenedUnixPermissions` | 非 Windows 断言 dir 0700、body/sidecar 0600。本会话 Windows：**SKIP**（见 N-001） |

## Findings（本条新发现）

开放 **required = 0**。下列不阻断 S6 代码闭合。

### F-008 · recommended · 低 · 邀请页缺少 `replaceState` 回归锁

F-005 实现正确，但 `apps/web` 无针对 `InviteAcceptPage` 初始化器的 history/search 断言。后续重构容易把 scrub 挪到提交成功之后或删掉。建议：jsdom 测「挂载后 `location.search` 无 `token` 且 `history.replaceState` 被调用一次」。

### F-009 · recommended · 低 · 若干 server 配置负例在 F-001 后假绿

`TestLoadConfigInvalidShutdownTimeout` 的 `"0s"`/`"-1s"` 与 `TestLoadConfigDialectPairing` 的自定义 YAML **省略 `app.env`**。`validate()` 现在先拒绝空 `APP_ENV`，测试仍 `err != nil` 通过，但不再咬到超时/方言分支。`"abc"` 超时仍在 ParseDuration 阶段失败（真负例）。产品 fail-closed 未破。建议：这些 YAML 显式写 `app.env: development`（或断言错误字符串含目标原因）。

### N-001 · note · F-007 POSIX 权限测试在本审计主机 skip

Windows 不强制 POSIX mode。本条以源码 + skip 路径为证据，**未在本机实证 0700/0600**。Linux/darwin CI 或后续 Unix 复跑才能锁权限位。不构成实现缺口。

### N-002 · note · 注释与审计索引滞后

1. `server/config.go:39` 仍写 `AppEnv "" = development（缺省）`，与 `validate()` 空 env fail-closed 矛盾。
2. 本条写入前 `03-audit.md` 信息表仍写 I-001～I-005 `open`、目标 `draft`。不否定 D-002。

### N-003 · note · `TestShutdownDrainHarnessPostgres` 在负载下曾 flake

本会话第一次 `go test ./... -count=1`（与 vitest 并行）该用例 EOF 失败；隔离复跑与第二次全量均绿。测试来自 VP-021（`1295de83`），W15 未改 `internal/composition`。不升格 required；全量套件与重型 Web 测试并行时不稳定。

## 必改项汇总

**无。** A-001 required F-001～F-006 与 recommended F-007 均 genuine-fixed。本条无新 required。

正式 `fixed` 标记、S6 检查点勾选、用户书面关门仍须 `/govern` 按 P-003 三路径响应（A-004）后才能改 status。

## 与既有意见的异同

| 来源 | 关系 |
|------|------|
| A-001 independent intake | 分母全部落地；本条不重开 R-001、不把主仓 `:25080` 并入 F-001 |
| A-002 self | 逐条结论一致。采纳 S-001（种子分层）与 S-002（guard 扫描成本）为非缺陷。本条补 F-008/F-009 测试质量建议与 N-001～N-003，不与 self 冲突 |
| D-002 | 冻结方案与 as-built 一致；F-007=fixed 有用户书面裁决 |

## 结论 + 建议给编排器/用户的下一步

**verdict = pass。** 代码闭合条件已满足；开放必改 = 0；S6 independent 腿已完成。

建议 `/govern`：

1. 响应本条：将 A-001 F-001～F-007 标为 `fixed`（证据 = 本条逐条表 + E-001～E-003 + 本会话回归）。
2. 对 F-008/F-009：本波可 `accepted-residual`（关门后补测试）或立即修；非关门阻断。
3. 勾选 S6 检查点、同步 goal-tree；**用户书面授权后再 `status: done`**。
4. 不在本意见中恢复或改写 VP-008 `go`。

## 声明

本意见不修改 status/progress/goal-tree；响应由 `/govern` 处理。
