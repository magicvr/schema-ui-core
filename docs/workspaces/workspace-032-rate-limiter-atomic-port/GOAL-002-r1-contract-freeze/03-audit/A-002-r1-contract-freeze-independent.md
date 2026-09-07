---
id: A-002-r1-contract-freeze-independent
doc: audit-entry
record_id: A-002
source: independent
scope: GOAL-002 R1 合同冻结关门（C1/C2 验收 + C3 关门证据）
verdict: conditional
status: recorded
parent: GOAL-002-r1-contract-freeze
created: 2026-09-03
updated: 2026-09-03
version: 0.1.0
auditor: grok-build (grok-4.6 · reasoning high)
audit_type: close-out
open_required: 1
---

# A-002 · GOAL-002 R1 合同冻结独立交叉审计（2026-09-03）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型** / **scope**：close-out（GOAL-002 R1 合同冻结 · C1 继承冻结 / C2 合同+端口落地 / C3 关门证据）
- **verdict**：**conditional**
- **open required**：1（F-001）

本意见不修改 `status` / `progress`；响应由 `/govern` 处理。

## 范围与区间

| 项 | 值 |
|----|-----|
| 工作区 | `workspace-032-rate-limiter-atomic-port`（Root `GOAL-001-rate-limiter-atomic-port`；canonical `docs/workspaces/workspace-032-rate-limiter-atomic-port/`；`shared_materials_catalog: none`） |
| 目标 | `GOAL-002-r1-contract-freeze`（parent: `GOAL-001-rate-limiter-atomic-port`） |
| 对齐 | `primary_plan` = `VP-032-rate-limiter-atomic-port`（active v0.2.0）；不审愿景层、不写 `docs/vision/reviews.md` |
| 覆盖 | D-001 / D-002；E-001 / E-002；A-001 self；`kernel.RateLimiter.AllowRecord`；Memory 单锁实现；合同级测试；14 处生产调用点未迁；I-032-001/002 门禁；红线（Redis / Profile / VP-027 重开） |
| 不覆盖 | R2 十四处迁移与 handler 行为等价；R3 证据矩阵 / 越界核账；Root 纲领关门；其它工作区正文（D-002 对 workspace-027 的引用仅核 Q2 路径形状，未读取该区） |
| 审计模式备注 | Root D-001：R1 阶段关门 default **self**，independent 留 R3。本次为用户 `/audit workspace-032 GOAL-002` + A-001 书面请求的额外交叉意见，不把 R1 升格为 R3 实证门禁。 |

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| 工作区绑定合格 | `workspace.md`：id / root_goal / canonical_scope / `plan_refs`+`primary_plan` 与 GOAL-002 `parent`、`plan_refs` 一致；资料目录 `none`，无无效共享资料引用 |
| C1 继承激活冻结，无新 P-004 | `01-decision/D-001-inherit-activation-freeze.md`：I-032-001/002 沿用 VRev-073；签名 `AllowRecord(key string, now time.Time) bool`；14 处全迁归 R2；Clear 无原子变体 |
| C2 合同冻结 | `01-decision/D-002-allowrecord-port-contract.md` v0.1.0：单锁顺序等价、剪枝/容量表、并发预算 `min(N,max)`、立即消费 vs 失败预算口径、14 处分母、R1/R2 切分、红线 |
| 内核接口加法且兼容方法保留 | `apps/api/kernel/ratelimit.go`：接口增 `AllowRecord`；`Allow`/`Record`/`Clear`/`RetryAfterSeconds` 仍在；godoc 标明新调用点 SHOULD 用 `AllowRecord`，Retry-After 仅在 Allow **或** AllowRecord false 后调用 |
| Memory 单锁实现 | `apps/api/internal/ratelimit/memory.go`：`AllowRecord` 在同一 `mu` 内 `allowLocked` 然后 `recordLocked`；false 不走 Record；`var _ kernel.RateLimiter = (*Memory)(nil)` |
| 合同级测试存在且本轮复跑绿 | `apps/api/kernel/ratelimit_test.go` stub 实现 `AllowRecord`；`memory_test.go`：`TestMemoryAllowRecordSequentialEquivalence`、`TestMemoryAllowRecordDenyDoesNotGrow`、`TestMemoryAllowRecordConcurrentBudget`；`TestMemoryConcurrent` 混入 `AllowRecord`。本轮（`apps/api`，`-count=1`）：`go test ./kernel/ ./internal/ratelimit/` 绿（kernel 0.562s / ratelimit 0.480s）；`go test -count=1 -race ./internal/ratelimit/` 绿（1.673s）；`go build ./...` 通过 |
| 生产 14 处未迁 | 当前 HEAD `98edb03e` 的 `--name-only` 不含 handler / telegram webhook / `go.mod` / Profile。工作区树内仍见 Allow→Record：`auth.go`、`captcha.go`、`account_self.go`、`recovery.go`、`mfa.go`（verify + `guardMFAStepUp` enroll/disable/recovery-rotate）、`invites.go`、`wallet_self.go`、`webhook.go` 三桶 |
| 未实现 Redis / 未改依赖锁 | `apps/api/go.mod` 无 redis 客户端；`98edb03e` 未改 `go.mod` |
| I-032-001/002 不阻断 C3 | 目标 `00-meta` + VP-032：两项 **verified**，最晚阶段 = 方案冻结 / C1；无 `deferred`、无到期开放 required、无残余风险待接受 |
| D-002 跨区引用形状合格 | `../../../workspace-027-rate-limiter-port/GOAL-002-r1-contract-freeze/01-decision/D-002-rate-limiter-port-contract.md`（文档 Q2）；本意见未读取该工作区正文 |

## 对照成功标准

| 成功标准（GOAL-002 `00-meta`） | 状态 | 证据 |
|-------------------------------|------|------|
| 1. 合同冻结：签名、单锁语义、顺序等价、false 路径不登记 key，写入 D-002 且可测试 | **已达成** | D-002 §1–§3；Memory 实现与上述三项新测试 + 本轮复跑 |
| 2. 兼容：`Allow`/`Record`/`Clear`/`RetryAfterSeconds` 语义不变；Retry-After 仅在 Allow 或 AllowRecord false 后调用 | **已达成**（接口与既有测试） | 接口保留；既有 `TestMemoryAllowDoesNotRegisterKey` / 滑动窗口 / RetryAfter / 容量回落仍绿；godoc 已更新 |
| 3. 可编译：kernel stub 与 Memory 均实现新方法；`go test` kernel + ratelimit 绿 | **已达成** | stub + Memory；本轮 `go test` / `-race` / `go build ./...` |
| 4. 未越界：不迁 14 处；不实现 Redis；不改 Profile 默认集；不重开 VP-027 | **已达成** | `98edb03e` 文件清单；生产调用点仍为 Allow+Record；`go.mod` 无 redis |

## 信息就绪核对（P-005）

| ID | 级别 | 最晚阶段 | 状态 | 是否阻断本 scope |
|----|------|----------|------|------------------|
| I-032-001 | required | C1 / 方案冻结 | **verified**（VRev-073 → D-001 / D-002 §1） | 否 |
| I-032-002 | required | C1 / 方案冻结 | **verified**（VRev-073 → D-001 / D-002 §4/§5） | 否 |

到期且影响本 scope 的开放 required 信息项：**0**。无共享资料引用被当成事实或关闭证据。

## Findings

### F-001 · 执行台账 checkpoint SHA 与当前 HEAD 不一致

| 字段 | 值 |
|------|-----|
| 严重度 | med |
| 建议 | **required** |
| status | **open** |
| evidence | GOAL-002 `02-execution/E-002-contract-frozen.md` 写 `bdfe925f`；Root `GOAL-001/.../E-002-r1-contract-freeze.md` 同 SHA；`git merge-base --is-ancestor bdfe925f HEAD` 失败；HEAD = `98edb03e`（同 message 的 amend，CommitDate 晚 36s）；`git diff bdfe925f 98edb03e` 仅上述两份 E-002 各多一行 SHA 文本，代码树其余相同 |
| closure | 待 `/govern` 把 GOAL-002（及如需 Root 镜像）E-002 的 checkpoint 改为当前分支上的 `98edb03e`，或另做可核对 checkpoint 并只暂存 owned paths |

P-002 要求实施事实可指回证据。`bdfe925f` 不在当前 `dev` 祖先链上，按 E-002 无法从 HEAD 历史复现该 checkpoint。A-001 已改引 `98edb03e`（正确），但执行台账仍写失效对象。这不否定 C2 代码/测试本身，但关门证据链不完整，**不得在未闭合本条时把 GOAL-002 标 `done`**。

### F-002 · AllowRecord true 路径的容量驱逐缺少专用测试

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | recommended |
| status | open |
| evidence | D-002 §2 表：AllowRecord true「缺席则可能 FIFO 驱逐」。实现走共享 `recordLocked`（`memory.go`），既有 `TestMemoryAllowDoesNotRegisterKey` / `TestMemorySlidingWindowSemantics` 只覆盖 `Record` 驱逐；无 AllowRecord 喷洒 distinct key 的对等断言 |
| closure | 非关门阻断。R2 前或合同测试补强时可加；接受残余则书面范围=「驱逐语义由 `recordLocked` 共享，回归仍走 Record 测试」 |

不构成 D-002 §9 R1 验收缺口（§9 列的是顺序等价 / 拒绝不登记 / 并发预算 / 既有回归，本轮均绿）。

## 必改项汇总

1. **F-001**（required / med）：更正 GOAL-002 E-002（建议同步 Root E-002）checkpoint 为当前分支祖先 `98edb03e`，或重新打可核对 checkpoint。未闭合前不得 `status: done`。

开放必改项数：**1**。

## 与既有意见的异同

| 项 | A-001 self（pass · 0 required） | 本意见 independent |
|----|-------------------------------|-------------------|
| 合同 / 接口 / Memory 单锁 / 合同测试 | 达成 | **同意**；本轮独立复跑 kernel + ratelimit + `-race` 仍绿 |
| 14 处未迁 / 无 Redis | 达成（A-001 引 `98edb03e --stat`） | **同意**（独立核 HEAD 文件清单与调用点） |
| I-032-001/002 | verified，不阻断 | **同意** |
| handler / telegram 套件 | A-001 记 `./internal/handler/...` 与 telegram 绿 | **未复跑**（超出 D-002 §9 R1 合同级验收；不作为本意见关闭或否决依据） |
| Git checkpoint | A-001 写 `98edb03e` | **部分不同意**：A-001 的 SHA 对应当前 HEAD；E-002 仍写非祖先 `bdfe925f` → **F-001** |
| 关门放行 | self `pass`，建议接着 `/audit` | **conditional**：技术验收可核对，执行台账 SHA 必改后再关门 |

无 verdict 对撞需要用户在「合同是否冻结」上二选一；冲突仅在「执行 SHA 是否已可重复核对」。建议按 F-001 **fixed**，不建议 residual/overruled。

## 结论 + 建议给编排器/用户的下一步

C1/C2 的合同、端口加法、Memory 单锁与合同级测试**名实相符**，红线未破，信息门禁已关闭。独立意见不能无条件放行 GOAL-002 关门，因为执行台账指向一个不在当前分支上的 checkpoint。

建议 `/govern`：

1. 响应 A-001 + A-002；将 F-001 **fixed**（改 E-002 SHA 或重打 checkpoint）。
2. F-002 可带入 R2 或书面 accepted-residual（低优先级）。
3. F-001 闭合后将 GOAL-002 C3 关门、`status: done`、Root progress → 1/3，再立项 R2（14 处迁移）。不要用 `progress: 2/3` 代替本条闭合。

## 声明

本意见 `source: independent`，不修改目标 `status` / 检查点 / 派生 `progress` / 方案正文 / goal-tree。响应与状态变更走 `/govern` 与用户裁决。
