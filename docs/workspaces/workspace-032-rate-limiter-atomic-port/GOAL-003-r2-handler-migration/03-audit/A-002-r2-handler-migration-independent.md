---
id: A-002-r2-handler-migration-independent
doc: audit-entry
record_id: A-002
source: independent
scope: GOAL-003 R2 生产使用点迁移与 handler 回归（C1/C2/C3 全目标关门）
verdict: fail
status: recorded
parent: GOAL-003-r2-handler-migration
created: 2026-09-04
updated: 2026-09-04
version: 0.1.0
auditor: Codex (GPT-5 · independent /audit)
audit_type: close-out
open_required: 2
---

# A-002 · GOAL-003 R2 生产迁移独立交叉审计（2026-09-04）

- **source**：independent
- **auditor**：Codex（GPT-5 · independent `/audit`）
- **类型 / scope**：close-out（GOAL-003 R2 生产使用点迁移与 handler 回归；C1/C2/C3 全目标关门）
- **verdict**：**fail**
- **open required**：**2**（F-001、F-002）

本意见不修改 `status` / `progress` / 检查点 / `goal-tree`；响应与状态变更由 `/govern` 处理。

## 范围与区间

| 项 | 值 |
|----|----|
| 工作区 | `workspace-032-rate-limiter-atomic-port`；Root `GOAL-001-rate-limiter-atomic-port`；canonical `docs/workspaces/workspace-032-rate-limiter-atomic-port/` |
| 焦点目标 | `GOAL-003-r2-handler-migration`（parent: `GOAL-001-rate-limiter-atomic-port`；`active` · `2/3`） |
| 规划对齐 | `primary_plan` = `VP-032-rate-limiter-atomic-port`（`active` v0.2.0）；`vision_ref` = `schema-ui-core-admin-foundation@0.4.0` |
| 实施区间 | commit `5c3437768a55` → `b08798d4d4af`（当前 `dev` HEAD） |
| 覆盖 | D-001；E-001/E-002；A-001 self；14 处冻结生产使用点；新增并发测试；handler/telegram 与全仓 Go 回归；I-032-001/002；Redis/Profile/kernel 红线 |
| 不覆盖 | VP-032 愿景层关门；R3 最终证据矩阵；其他工作区正文 |
| 共享资料 | `shared_materials_catalog: none`；本意见未使用共享资料作为证据 |

## 成果（有证据）

1. **迁移分母已落地**：生产代码共 12 个直接 `.AllowRecord(` 表达式；其中 `guardMFAStepUp` 被 enroll / disable / recovery-rotate 三个语义使用点共用，因而对应冻结 **14/14** 使用点。`apps/api/internal/handler` 与 `apps/api/internal/channel/telegram` 生产码中已无 RateLimiter `.Allow(` / `.Record(`；仅 `resources.go` 保留无关的 `Trash.Record`。
2. **提交边界未越界**：`b08798d4` 只修改 10 个显式 code/test path（8 个生产文件 + `auth_test.go` + `webhook_test.go`），未改 `kernel` / Redis / Profile / Manifest / `go.mod`。
3. **现有测试绿灯可复现**：
   - `go test -count=1 ./internal/handler/... ./internal/channel/telegram/...` PASS（handler 29.119s；telegram 1.587s）。
   - `go test -count=1 -race -run 'Test(LoginRateLimit|PasswordChange|MFA|Recovery|InviteAccept|WalletSelf|Captcha)' ./internal/handler` PASS（70.443s）。
   - `go test -count=1 -race ./internal/channel/telegram/...` PASS（6.375s）。
   - `go test -count=1 ./...` PASS（含 `internal/docscheck`）。
4. **红线静态核账通过**：`Allow` / `Record` 兼容接口仍在；本 commit 未重开 VP-027，未引入 Redis，未改 Profile 默认集或其他内核端口。

## 对照成功标准

| 成功标准 | 结论 | 证据 / 缺口 |
|----------|------|-------------|
| 1. 14 处分母全覆盖 | **已达成** | 12 个直接 `AllowRecord` 表达式 + `guardMFAStepUp` 三个语义调用点 = 14；无 RateLimiter Allow→Record 剩余配对 |
| 2. 行为等价 | **未达成** | 失败预算在结果已知前用 `AllowRecord` 占槽，但 `Clear` 删除该 key **全部**历史，无法只回滚本次占槽；已产生可证实的计数与清空语义回归（F-001） |
| 3. handler / channel 回归全绿 | **字面达成，但不能证明等价** | 现有套件与 `-race` 均绿；关键混合历史转换没有测试（F-002） |
| 4. 红线保持 | **已达成** | commit path 清单与静态搜索 |

## 信息就绪核对（P-005）

| ID | 级别 | 当前投影 | 本次核对 |
|----|------|----------|----------|
| I-032-001 | required | `verified` | 签名与 bool 返回语义已落地，不是本次缺口 |
| I-032-002 | required | `verified` | “14 处全迁”已实施；但“Clear 无需原子变体 / 失败预算行为等价”被 F-001 的实施证据否定。`/audit` 不改信息项状态；须由 `/govern` 回流重审该决策，必要时登记新信息项 |

到期且影响本 scope 的新证据冲突：**1 类**（I-032-002 的 Clear/行为等价结论）。未有用户书面 residual 或 overrule 记录，因此按 P-004/P-005 保持 fail closed。

## Findings

### F-001 · `AllowRecord` + key-wide `Clear` 破坏失败预算语义

| 字段 | 值 |
|------|----|
| 严重度 | **high** |
| 建议 | **required** |
| status | **open** |
| 影响门禁 | GOAL-003 C2 行为等价；C3 关门；VP-032 退出判据 #2/#5 |

**证据：**

1. `apps/api/internal/ratelimit/memory.go:153-161` 的 `Clear` 直接 `delete(l.attempts, key)`，删除该 key 的全部历史；它不是“取消本次 AllowRecord 占槽”。R1 D-002 却在 `D-002-allowrecord-port-contract.md:56,85-87,146` 拒绝原子 Clear / CompareAndClear，并将所有失败预算假定为旧路径都是“成功 Clear”。
2. `apps/api/internal/handler/recovery.go:93-117`：入口 `AllowRecord` 占槽后，`ErrRecoveryNotAvailable` 分支不返回，继续落到 `Clear(key)`。父提交中该分支是 `Record` 且没有后续 Clear；现实现会使 no-path/unknown-account 请求不再累计，从而无法到达 20 次后的 429。
3. `apps/api/internal/handler/auth.go:109-124`：已有密码失败历史时，一次无效 CAPTCHA 先占槽再 `Clear(limiterKey)`，会同时清掉之前的密码失败。旧路径的 CAPTCHA 失败不修改此 bucket；新路径允许客户端用无效 CAPTCHA 重置登录失败预算。
4. 同类未冻结语义变更还出现在 `recovery.go:146-238`（非猜测分支/成功清全历史，部分内部错误反而留下占槽）、`mfa.go:143-167`（新增成功 Clear）、`invites.go:327-358`（新增成功 Clear）。E-002:23-29 把这些新 Clear 写成等价实施，但提交前源码中对应成功路径并不都 Clear。

**结论：** `AllowRecord` 只能原子表达“每一次放行都立即消费”。对“只计指定失败”的路径，在结果已知前乐观占槽，之后又只有 key-wide `Clear`，无法同时保留历史失败并回滚当次非计数占槽。A-001:49 的“行为等价已达成”与代码事实冲突。

**闭合要求：**

- 由 `/govern` 回流重审 I-032-002 / D-002 的失败预算协议；建议优先采用能区分当次占槽的原子 reservation + commit/cancel（或等价的 tokenized rollback），而不是 key-wide Clear。
- 对 10 个失败预算语义点逐一冻结“哪些结果计数 / 哪些不计数 / 成功是否清历史”，并修复实现。若要把旧的 failure-only 预算改成 all-attempt 预算或允许成功清历史，这是产品/安全语义变更，须依 P-004 由用户书面裁决，不能以“实现更保守”默认带过。

### F-002 · 现有测试分母未覆盖历史保留与 no-path 累计

| 字段 | 值 |
|------|----|
| 严重度 | med |
| 建议 | **required** |
| status | **open** |
| 影响门禁 | GOAL-003 C2 回归证据；C3 关门 |

**证据：**

- 新增 handler 测试只有 `auth_test.go:362` 的“纯失败并发不穿透”和 `auth_test.go:402` 的“成功清空”；Telegram 新测试只覆盖立即消费 IP 桶。
- `recovery_test.go:80-100` 对 `ErrRecoveryNotAvailable` 只发送一次请求，没有验证其应累计至 429。
- `recovery_test.go:195-210` 的 non-guess 测试从空 bucket 开始，只能证明每轮后为空，不能发现 `Clear` 同时删掉进入该分支前已累积的失败历史。
- 因此，现有套件全绿与 A-001:49 的“失败预算净状态等价”之间没有充分证据链。

**闭合要求：**在 F-001 的语义方案冻结后，补充至少下列回归：

1. 同一 IP|account 的 `ErrRecoveryNotAvailable` 连续请求能精确在预算后返回 429。
2. 已有登录密码失败时，无效 CAPTCHA 不得删除既有失败历史。
3. 对 recovery complete / MFA verify / invite accept 至少覆盖“先有历史失败，再走成功或非计数分支”的混合序列，断言结果与经 P-004 冻结的语义一致。
4. 复跑 handler / telegram 全量、上述安全用例 `-race` 及 `go test -count=1 ./...`。

### F-003 · Root / workspace 的 R2 路线图投影滞后

| 字段 | 值 |
|------|----|
| 严重度 | low |
| 建议 | recommended |
| status | open |
| 影响门禁 | 不单独阻断 GOAL-003 技术修复；在关门响应时应同步 |

`GOAL-001-rate-limiter-atomic-port/00-meta.md:37` 仍写“待迁 14 处调用点”，`workspace.md:49` 仍只写“GOAL-003 已立项”；而 GOAL-003 与 `goal-tree.md` 已投影 C1/C2 完成。建议 `/govern` 响应时修正 Root/workspace 路线图文字，保持 P-001 路线图可追踪。

## 必改项汇总

1. **F-001（high / required）**：失败预算的 `AllowRecord` + key-wide `Clear` 破坏旧计数与历史保留语义；已造成 recovery start no-path 不累计和无效 CAPTCHA 可清登录失败的安全回归。
2. **F-002（med / required）**：补齐能防止 F-001 回归的混合历史/no-path 测试，并重建“行为等价”证据链。

开放必改项数：**2**。在 F-001/F-002 未按 `fixed` / `accepted-residual` / `user-overruled` 合法闭合前，不得关门 GOAL-003，也不得以 A-001 `pass` 单独放行 R2。

## 与既有意见的异同

| 项 | A-001 self | A-002 independent |
|----|------------|-------------------|
| 14 处静态迁移 | pass | **同意** |
| commit 边界 / 红线 | pass | **同意** |
| 现有测试是否全绿 | pass | **同意，已独立复跑** |
| 失败预算行为等价 | 已达成 | **不同意**（F-001） |
| 关门 | pass / 0 required | **fail / 2 required** |

A-001 `pass` 与本意见 `fail` 覆盖同一 close-out scope，符合 P-004.2 的 verdict 冲突条件。编排器须展示冲突、给出建议并留下用户裁决/修复记录；未响应前 fail closed。本意见建议修复 F-001/F-002，不建议 residual 或 overrule。

## 结论与下一步

GOAL-003 的“14 处已迁移”、“现有套件全绿”和“未越界”都可复核；但 close-out 的关键主张“失败预算行为等价”不实，且包含可用于绕过限流的路径，因此 verdict 为 **fail**。

建议给编排器的下一句：

> `/govern workspace-032 GOAL-003 响应 A-001/A-002：按 fixed 路径处理 F-001/F-002，先回流失败预算语义与 I-032-002，修复并补测后再复审；当前不关门。`

## 声明

本意见 `source: independent`，为 L0 入口分离级交叉意见，不等同于外部法定鉴证。本意见不修改目标 `status` / 检查点 / 派生 `progress` / 方案正文 / `goal-tree`；响应、修正、finding 闭合与状态推进由 `/govern` 处理。
