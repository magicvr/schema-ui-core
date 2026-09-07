---
doc_type: goal-audit
id: A-002-independent-closeout-audit
parent: GOAL-005-r4-evidence-closeout
date: 2026-09-05
status: open
source: independent
auditor: codex (gpt-5.6-sol, medium)
audit_type: close-out
scope: Root GOAL-001 关门审计；对照成功标准 1～8、VP-031 方向级退出判据、GOAL-005 E-001、子目标审计闭环与信息台账
verdict: conditional
open_required: 2
version: 1.0.0
---

# A-002 · Root 独立关门审计（close-out）

## A-002 · 数字 Offer 与权益 Root 独立关门审计（2026-09-05）

- **source**：independent
- **auditor**：codex (gpt-5.6-sol, medium)
- **类型 / scope**：close-out · 在 GOAL-005 E-001 证据矩阵之上，核对 Root GOAL-001 成功标准 1～8、VP-031 方向级退出判据 1～8、GOAL-002/003/004 审计闭环、I-031-001～005 与边界核账
- **verdict**：**conditional**
- **open required**：**2**

### 范围与区间

- 审计对象是 `workspace-031-digital-offer-entitlement` 的 Root 关门门禁；工作区绑定为 `root_goal = GOAL-001-digital-offer-entitlement`、`canonical_scope = docs/workspaces/workspace-031-digital-offer-entitlement/`、`primary_plan = VP-031-digital-offer-entitlement`，且共享资料目录为 `none`（`workspace.md:1-15,24-31`）。
- 成功分母采用 Root 成功标准 1～8（`GOAL-001-digital-offer-entitlement/00-meta.md:21-31`）与 VP-031 方向级退出判据 1～8（`docs/vision/plans/VP-031-digital-offer-entitlement.md:84-93`）；合同分母采用 GOAL-002 D-002 v1.2.0，信息裁决采用 D-001（`GOAL-002-r1-contract-freeze/01-decision/D-001-info-adjudication.md:10-27,36-39`）。
- 静态抽查覆盖 `apps/api/modules/digitaloffer/`、`apps/api/internal/handler/digitaloffer.go` 及指定测试；动态核验在 `apps/api/` Go module 内执行，定向测试通过，且 `go build ./...`、`go test ./...` 退出码为 0。仓库根目录不是 Go module，根目录执行相同命令会以 `directory prefix . does not contain main module` 失败，见 F-002。

### 成果（有证据）

1. **Offer/Admin/C 端主路径真实存在。** Admin Offer、purchase、entitlement 路由分别执行 `digitaloffer.read`、`digitaloffer.offer.manage`、`digitaloffer.entitlement.void` 权限门禁，公开目录只列 on-sale Offer 并实施 IP 限流（`apps/api/internal/handler/digitaloffer.go:45-72,74-154,157-267`）；Create/Update 与审计记录位于同一事务（`apps/api/modules/digitaloffer/service/service.go:155-188,201-247`）；handler 测试覆盖 Admin 401、创建/变更、公开投影、审计失败回滚与 429（`apps/api/internal/handler/digitaloffer_test.go:85-171,215-229`）。
2. **购买事务、幂等与并发证据成立。** Purchase 在一个 caller-owned transaction 内完成 purchase voucher、wallet freeze、deduct_frozen 与 entitlement；失败由事务整体回滚，属于 D-002 冻结的更强 fail-closed 形态（`apps/api/modules/digitaloffer/service/service.go:285-350,356-469`）。测试覆盖余额不足零残留、同 request 幂等、跨 Offer 冲突、并发同 request 收敛、余额竞争、重试耗尽与瞬态恢复（`apps/api/modules/digitaloffer/service/purchase_test.go:247-498`）。
3. **权益核验与消耗证据成立。** Check 聚合有效/过期/耗尽/作废语义；Consume 使用新事务重读、oldest-first 候选与条件 UPDATE 线性化（`apps/api/modules/digitaloffer/service/service.go:595-717`；`apps/api/modules/digitaloffer/store/store.go:638-685`）。测试覆盖 duration 过期、count 耗尽/不足零扣减、多行顺序与 void/consume 并发（`apps/api/modules/digitaloffer/service/check_consume_test.go:95-274`）。
4. **subject-only 与 Telegram 接缝成立。** Purchase 先做 `SubjectExistsInTx`，purchase/entitlement schema 与 store 均以 `subject_id` 为主体（`apps/api/modules/digitaloffer/service/service.go:374-395`；`apps/api/modules/digitaloffer/migration/migration.go:41-59,103-121`）。`RegisterTelegram` 注册 `price/buy/entitlements`，测试覆盖 Disabled 注入、真实 dispatcher 回复、身份门控、幂等与查询桶（`apps/api/modules/digitaloffer/service/service.go:728-760`；`apps/api/modules/digitaloffer/service/check_consume_test.go:291-370`）。
5. **三个子目标不存在开放 required。** GOAL-002 A-006 为 independent `pass`、`open_required: 0`，并确认 A-004 required 已按 fixed 闭合（`GOAL-002-r1-contract-freeze/03-audit/A-006-independent-closure-review-2.md:9-24,32-60`）；GOAL-003 A-008 为 independent `pass`、`open_required: 0`（`GOAL-003-r2-offer-purchase-wallet/03-audit/A-008-independent-closure-review-3.md:8-22,32-45`）；GOAL-004 A-002 为 independent `pass`、`open required 0`，A-003 已处置其两项 recommended（`GOAL-004-r3-entitlement-validation-telegram/03-audit/A-002-independent-implementation-audit.md:14-24,68-79`；`A-003-self-response-a002.md:14-27`）。
6. **I-031-001～005 无到期 required 残留。** D-001 记录三个 required 的用户书面裁决与两个 non-blocking 默认冻结，并明确回写为 verified（`GOAL-002-r1-contract-freeze/01-decision/D-001-info-adjudication.md:14-27,36-39`）；Root 与 VP-031 当前表均显示 I-031-001～005 为 verified（`GOAL-001-digital-offer-entitlement/00-meta.md:48-54`；`VP-031-digital-offer-entitlement.md:97-105`）。
7. **边界核账的业务事实可信。** 默认 Profile 未包含 `biz.digital-offer`（`apps/api/kernel/profile.go:25-114`）；迁移 0070 的 SQLite/PostgreSQL DDL 都只创建 `digital_offers`、`digital_purchases`、`digital_entitlements` 三表（`apps/api/modules/digitaloffer/migration/migration.go:21-75,83-137,153-163`）；store 未引用 `admin.users`，只操作本域三表与 `subject_id`（`apps/api/modules/digitaloffer/store/store.go:353-419,452-480,638-685`）。Charter 最近一次提交为 `1694dea74889ee0ca29949b3892bba4eea3e0829`（2026-09-01），早于 2026-09-05 的 R1～R3 记录；当前 Charter 仍为 0.4.0（`docs/vision/charter.md:5-10`；`git log -1 --follow -- docs/vision/charter.md`）。

### 对照 Root 成功标准

| 标准 1～8 | 状态 | 证据 |
|-----------|------|------|
| 1 · Offer CRUD + Admin 页面/权限/审计 + C 端列表 | **达成** | handler 路由与公开投影：`apps/api/internal/handler/digitaloffer.go:45-267,302-319`；事务审计：`service/service.go:155-247`；测试：`internal/handler/digitaloffer_test.go:85-171,215-229`。 |
| 2 · 余额不足、freeze→deduct_frozen、fail-closed/unfreeze 语义、并发 | **达成** | 购买事务：`service/service.go:285-469`；失败/幂等/并发/重试矩阵：`service/purchase_test.go:247-498,599-712`；D-002 以单事务回滚取代可见补偿窗口。 |
| 3 · 有效/过期/耗尽可测，服务前统一核验 | **达成** | Check/Consume：`service/service.go:595-717`；条件更新：`store/store.go:638-685`；测试：`service/check_consume_test.go:95-274`。 |
| 4 · 只挂 subject_id，不创建 admin.users | **达成** | `service/service.go:374-395`；`migration/migration.go:41-59,103-121`；`store/store.go:353-480`。 |
| 5 · Telegram enabled 注册命令；disabled 不依赖 Bot API | **达成** | `RegisterTelegram` 与 Disabled/Reply/RateLimit 测试：`service/service.go:728-760`；`service/check_consume_test.go:291-370`。 |
| 6 · freshness 与 RT-Q03/Q05 留痕 | **达成** | Root freshness 三字段：`GOAL-001/00-meta.md:56-62`；VP-031 激活门禁与评估：`VP-031-digital-offer-entitlement.md:11-18,75-93`。 |
| 7 · 红线/Profile/通用接缝/Charter 边界 | **达成（业务事实）** | Profile、迁移、store、Charter 历史核账见上；D-002 红线：`GOAL-002-r1-contract-freeze/01-decision/D-002-digital-offer-contract.md:252-261`。 |
| 8 · 关门前 open required = 0 | **未达成** | 子目标历史 required 均已关闭，但本 A-002 新增 F-001、F-002 两项 med required；关闭前必须按 P-003 合法闭合。 |

### Findings

#### F-001 · close-out 权威投影与 ledger 索引未同步（med · required）

**事实与证据**：

- GOAL-004 `03-audit.md` 只索引 A-002，但目录中实际存在 A-001、A-002、A-003；A-003 是对 A-002 两项 recommended 的正式响应，因此“索引 + ledger 共同构成正式台账”的投影不完整（`GOAL-004-r3-entitlement-validation-telegram/03-audit.md:11-16`；`03-audit/A-001-self-implementation-review.md`；`03-audit/A-003-self-response-a002.md:10-27`）。
- GOAL-005 的 `02-execution.md` 仍只有占位行，未索引已经存在且 `status: done` 的 E-001；`03-audit.md` 在本条写入前也只有占位行、未索引已存在的 A-001（`GOAL-005-r4-evidence-closeout/02-execution.md:11-19`；`GOAL-005-r4-evidence-closeout/03-audit.md:11-19`；`E-001-evidence-matrix.md:1-10`；`A-001-self-closeout-review.md:10-29`）。
- Root `03-audit.md` 仍把 I-031-001～005 记为 open，并称尚未进入 R1；这与 D-001、Root `00-meta.md` 和 VP-031 的 verified 状态冲突（`GOAL-001-digital-offer-entitlement/03-audit.md:16-24`；`GOAL-001-digital-offer-entitlement/00-meta.md:48-54`；`D-001-info-adjudication.md:36-39`）。
- canonical 投影还有陈旧重复：`workspace.md` 同时记录 R3 已关门/R4 进行中与 R3/R4 待开始；Root `00-meta.md` 同样重复；`goal-tree.md` ASCII 树写 Root `active · 1/4`，而状态表写 `3/4`（`workspace.md:33-42`；`GOAL-001-digital-offer-entitlement/00-meta.md:33-45`；`goal-tree.md:7-20`）。

**影响**：业务实现证据虽成立，但 close-out 的权威状态/进度/信息与正式 ledger 投影不唯一，不能在此状态下声称成功标准 8 已可重复核对。

**必须修正**：由 `/govern` 同步缺失索引，清除陈旧重复投影，统一 Root 信息台账及 workspace/Root/goal-tree/GOAL-005 的阶段事实；修正后以 finding-closure 复审核对。不得仅以当前聊天说明代替落盘。

#### F-002 · E-001 的“全仓 go build/test 绿”命令声明不可按字面复现（med · required）

**事实与证据**：E-001 写“全仓 `go build ./...` + `go test ./...` 回归绿”（`GOAL-005-r4-evidence-closeout/02-execution/E-001-evidence-matrix.md:15-20`）。仓库根目录没有 `go.mod`，从根目录执行两条命令均失败，错误为 `pattern ./...: directory prefix . does not contain main module or its selected dependencies`；在真实 Go module `apps/api/` 内，两条命令本轮退出码均为 0，定向数字 Offer 测试也通过。

**影响**：这是证据矩阵中的可复现命令边界错误，不推翻数字 Offer 实现与测试通过事实，但会使关门证据按记录命令重放失败。

**必须修正**：将 E-001 明确改为“在 `apps/api/` module 执行 `go build ./...` 与 `go test ./...` 均通过”，记录 cwd、可写 GOCACHE 与 PostgreSQL 测试是否运行/跳过；随后复审该证据声明。

### 必改项汇总

| Finding | 级别 | 状态 | 关门影响 |
|---------|------|------|----------|
| F-001 · close-out 投影与 ledger 索引未同步 | med required | open | 阻断 Root/GOAL-005 关门；需统一正式台账与状态投影后复审 |
| F-002 · E-001 全仓 Go 命令边界不可复现 | med required | open | 阻断证据矩阵无缺口声明；需修正 cwd/命令证据后复审 |

### 与既有意见的异同

- **同意 GOAL-005 A-001 的实现层判断**：Root 标准 1～7 的业务代码、测试、子目标审计与边界主张有充分证据（`GOAL-005-r4-evidence-closeout/03-audit/A-001-self-closeout-review.md:18-29`）。
- **不同意 A-001 的无条件 pass/0 findings**：A-001 未发现 GOAL-004/GOAL-005 索引缺口、Root/工作区陈旧投影及 E-001 命令 cwd 不可复现，因此本审计将 close-out verdict 收紧为 `conditional`。
- **不推翻三个子目标的最终 independent 结论**：GOAL-002 A-006、GOAL-003 A-008、GOAL-004 A-002 在各自限定 scope 内均为 pass 且 open required 0；本次 findings 针对 Root/R4 关门证据与 canonical 投影，不追溯改写原独立意见。

### 结论 + 建议下一步

- **Root 当前不可关门，不得将 `GOAL-001-digital-offer-entitlement` 标为 `done`。** 标准 1～7 的业务/边界事实成立，I-031-001～005 已 verified，三个子目标的历史 required 均已合法闭合；但 F-001、F-002 两项 med required 尚开放，因此成功标准 8 与 P-003/P-005 关门门禁未满足。
- 建议下一步：用 `/govern` 响应 A-002，按 fixed 路径同步 projection/ledger 并修正 E-001 命令边界；完成后发起限定 scope 的 independent finding-closure 复审。只有复审确认 `open required = 0` 后，编排器才可推进 GOAL-005/Root 状态与 goal-tree/VP 投影。

### 声明

本意见仅追加 `source: independent` 的 A-002 并更新 GOAL-005 `03-audit.md` 索引；**不修改** Root/子目标 `00-meta.md` 的 `status/progress`、`goal-tree.md`、01/02 台账、D-001/D-002 合同正文、VP-031 或业务代码。finding 响应、合法闭合与状态推进由 `/govern` 处理。
