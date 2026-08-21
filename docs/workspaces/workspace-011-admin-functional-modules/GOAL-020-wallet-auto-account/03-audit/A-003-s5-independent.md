---
id: A-003
goal: GOAL-020-wallet-auto-account
title: S5 关门独立审计 · 自动开户（方案+实现合并审）
date: 2026-08-16
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
scope: S5 关门（方案+实现合并审）
audit_type: close-out
verdict: conditional
status: recorded
parent: GOAL-020-wallet-auto-account
created: 2026-08-16
updated: 2026-08-16
version: 1.0.0
---

# A-003 · S5 关门独立审计（independent）

## 范围与区间（L0）

- **工作区**：`workspace-011-admin-functional-modules`；`root_goal` = `GOAL-001-admin-functional-modules`；`canonical_scope` 匹配；`shared_materials_catalog: none`（无资料引用可当事实）。
- **目标**：`GOAL-020-wallet-auto-account`（parent = Root；区内短 id 合法）。
- **scope**：S5 关门（close-out）——小目标合并审视：D-001 方案（触发面 / 幂等 / 审计形态 / 手动边界）+ S2/S3 实现一并独立审。data 门禁（自动开户 = 资金账户创建）。
- **保证等级**：L0（入口分离）。本意见不是第三方鉴证。
- **P-005**：`00-meta` / `01-decision` 中 I-001、I-002（required，最晚 S1）与 I-003（non-blocking，最晚 S1）均为 **verified**，无到期未闭合 required 信息项阻断关门。`03-audit.md` 索引仍写 I-001~I-003 open，与权威信息表不一致（见 F-005）。
- **上游基座（同区）**：GOAL-019 D-002 §1 UNIQUE(owner_type, owner_id, currency) / apply 表 / §3 既有端点。本目标声称复用 0031、无迁移、无权限/导航/Profile/协议变更。
- **本轮未做**：未复跑 web 1004/1004；未做真实多连接并发压测。go 侧本轮复跑 `apps/api` 的 `./internal/modules/wallet/store/` 与 `./internal/handler/`：**通过**（2026-08-16）。

## 成果（有证据）

| 主张 | 证据 | 核对 |
|------|------|------|
| 方案冻结：惰性 get-or-create + user 手动禁止 | `01-decision/D-001-auto-account-plan.md` §1–§3 | 达成（路径字符串见 F-004） |
| store SELECT→INSERT→UNIQUE 冲突重读代码存在 | `store/repository.go` 232–279（SELECT 237–244；INSERT 252–256；冲突重读 257–268） | 账户行幂等成立；`created` 见 F-001 |
| `created` 不靠余额启发式 | 标志来自「本次事务是否走到 INSERT」（249 行），非 `balance_*` | 顺序路径正确；冲突路径见 F-001 |
| by-owner 读 = `wallet.read`；调账 = `wallet.adjust` | `handler/wallet.go` 74–75、99–100 | 达成 |
| 自动开户审计 `wallet.account-create` + `"auto":true` | `handler/wallet.go` 90–93、137–140；`wallet_auto_test.go` TestWalletByOwnerAutoCreate 56–80（顺序二次读 = 1 次） | 顺序读路径达成；冲突/调账失败见 F-001/F-002 |
| 手动 user 创建 409 `WALLET_USER_AUTO_ONLY`；business 保留 | `handler/wallet.go` 163–167；`wallet_auto_test.go` 107–126 | 达成（system 手工路径代码允许，测试未覆盖） |
| 错误码双语 + 冻结集 + web catalog | `errorcatalog.go` 150（en/zh + `error.walletUserAutoOnly`）；`error_contract_test.go` 60；`apps/web/src/i18n/messages/{en-US,zh-CN}.json` 584 | 达成 |
| 路由改为 `/api/wallet/by-owner/`，避免与 `.../accounts/{id}/entries` 重叠 | `handler/wallet.go` 74、99；`provider.go` 156；`kernel/profile.go` 198；E-002 留痕 | 达成 |
| 复用 0031，无新迁移 | `migration/migration.go` 16–34（UNIQUE 仍在 0031）；`provider.go` 172；未见 0033 | 达成 |
| 无新权限/导航/Profile 默认集/协议 pin | `provider.go` 165 仍 `wallet.read/write/adjust`；`composition_test.go` admin 30/14；D-002 | 达成 |
| 前端新建表单去掉 user；展示键保留 | `schema/wallet.json` 75–78 仅 business/system；表列 `ownerType` 198 行无枚举过滤；`schema.wallet.owner.user` 仍在 en/zh 534 | 达成（既有 user 行可展示） |
| 顺序 get-or-create + 旧生命周期测试改走自动开户 | `store/repository_test.go` TestGetOrCreateUserAccount 280–318；`wallet_test.go` 136–138、258 | 顺序路径本轮测试通过 |
| S4 不暂挂 | D-002 | 与代码一致 |

## 对照成功标准（S1~S5）

| 标准 | 状态 | 证据 / 缺口 |
|------|------|-------------|
| S1 方案冻结 | 达成（有文档漂移） | D-001 冻结语义与实现一致；冻结路径仍写 `/api/wallet/accounts/by-owner/...`（F-004） |
| S2 实现 | 部分 | 触发面/门禁/错误码/前端选项落地；并发 `created` 与调账失败审计不成立（F-001、F-002） |
| S3 验证 | 部分 | 顺序路径与门禁 401 有测试；本轮 store+handler go test 通过。D-001 §4「并发冲突重读」无测试执行；403 新路由未进 `TestWalletRoutesGates` 第二圈（F-005）。web 全量本轮未复跑（E-002 声称 1004/1004） |
| S4 go 影响判定 + 自审 | 达成 | D-002 + A-002；组合根权限/导航未增 |
| S5 关门 | 未放行 | 本条 independent **conditional**；存在未闭合 med required |

`progress`（meta 4/5、goal-tree 0/5）不是放行或闭合依据。

## Findings

### F-001 · 并发 UNIQUE 冲突重读后 `created` 仍为 true（会重复审计）

- **严重度**：med
- **建议**：required
- **影响门禁**：S5 关门 / data（账户创建审计幂等）
- **状态**：open
- **描述**：`GetOrCreateUserAccount` 在 SELECT 未命中后于插入前把 `isNew = true`（`repository.go` 249）。INSERT 撞上 `UNIQUE(owner_type, owner_id, currency)` 时会重读已有行并 `return nil`（257–268），**没有**把 `isNew` 置回 `false`。handler 用该 bool 决定是否写 `wallet.account-create`（`wallet.go` 90–93、137–140）。同 owner 并发 get-or-create（读+读，或读+调账）时，败方仍报告新建 → **重复 auto 审计**。账户行本身不会双插（约束兜底成立）。`created` 判断不依赖余额，但也不等于「本事务是否成功 INSERT」。
- **证据**：
  - `apps/api/internal/modules/wallet/store/repository.go` 232–279（尤其 249 vs 257–268）
  - `apps/api/internal/handler/wallet.go` 85–93、122–140
  - `repository_test.go` TestGetOrCreateUserAccount 只做顺序二次调用（第二次走 241–244 SELECT 命中），注释写「concurrent-style」但**未执行**冲突重读；`wallet_auto_test.go` 同样只断言顺序二次读。
  - D-001 §1/§4 与 A-002 主张「冲突重读 + created 准确」——代码路径存在，标志语义在冲突分支不成立。
- **建议修正**：冲突重读成功后 `isNew = false`（或仅在 INSERT 成功分支置 true）；补一条可复现的 UNIQUE 冲突重读测试，断言 `created=false` 且不写第二笔 `wallet.account-create`。

### F-002 · by-owner 调账在 Mutate 失败时已开户但不写 account-create

- **严重度**：med
- **建议**：required
- **影响门禁**：S5 关门 / data（资金账户创建必须留审计）
- **状态**：open
- **描述**：`POST /api/wallet/by-owner/{ownerId}/adjust` 先 `GetOrCreateUserAccount`，再 `Mutate`，**仅当 Mutate 成功**才按 `created` 记 `wallet.account-create`（`wallet.go` 122–140）。JSON 能解但业务非法时（缺 memo、`amountDelta==0`、透支等）账户已以零余额落库，本请求不审计；后续成功读/调账时 `created=false`，**该账户可能永远没有** `wallet.account-create`。这不是 GOAL-019「operationlog 与账本非同事务」残余，而是把创建审计绑在后继调账成功上。读端点（90–93）是开户后立即审计，写端点不一致。
- **证据**：
  - `apps/api/internal/handler/wallet.go` 121–136（失败直接 `writeWalletError` return）相对 137–142（成功后才双审计）
  - `provider.go` 109–111：`Memo==""` 在 service.Mutate 才拒
  - `wallet_auto_test.go` TestWalletByOwnerAdjustAutoCreate 只断言余额/流水，**未**断言 `account-create`+`adjust`
- **建议修正**（择一并测）：开户成功立即记 auto 审计（与读路径对齐）；或先校验调账载荷再 get-or-create；或 Mutate 失败且本次新建时补偿审计。测试覆盖「非法调账后账户已存在 + 审计条数」。

### F-003 · 自动开户 id 未复用 `newID` 随机后缀，PK 冲突会被误判成 owner UNIQUE

- **严重度**：med
- **建议**：recommended
- **影响门禁**：可靠性（非账户双花）
- **状态**：open
- **描述**：GOAL-019 `newID`（`provider.go` 42–47）是毫秒前缀 + 12 字节随机，专门防 PK 碰撞。GetOrCreate 使用 `fmt.Sprintf("%016x%d", now.UnixMilli(), now.UnixNano()%1e9)`（`repository.go` 251），同一 `now` 值（测试夹具、同纳秒时钟）给**不同 owner** 会生成同一 `id`。`isUniqueViolation` 无法区分 PRIMARY KEY 与 UNIQUE(owner)（691–693）。PK 冲突会走「当 owner 并发成功」重读，对另一 owner 重读落空 → 500。fail-closed，可换时钟重试，不绑错户。
- **证据**：`repository.go` 250–268、691–693；`provider.go` 38–47。
- **建议**：插入前走 `newID`；或仅在 owner UNIQUE 冲突时重读。

### F-004 · D-001 冻结路径与实现不一致（已由 E-002 改前缀）

- **严重度**：low
- **建议**：recommended
- **影响门禁**：方案台账可核对性（不阻断功能）
- **状态**：open
- **描述**：D-001 §1 仍写 `GET/POST /api/wallet/accounts/by-owner/{ownerId}[ /adjust]`。实现、Descriptor、`kernel/profile.go` BuiltinModules、测试均为 `/api/wallet/by-owner/...`。E-002 已说明 Go 1.22 ServeMux 与 `GET /api/wallet/accounts/{id}/entries` 重叠。语义未变，冻结方案未勘误。
- **证据**：D-001 §1；E-002「路由形态修正」；`handler/wallet.go` 74、99；`provider.go` 156；`kernel/profile.go` 198。

### F-005 · 验证与台账卫生缺口（403 新路由、wallet 码表、索引漂移）

- **严重度**：low
- **建议**：recommended
- **影响门禁**：S3 证据完整度 / 关门卫生（非功能错误）
- **状态**：open
- **描述**：
  1. `TestWalletRoutesGates` 401 含 by-owner（`wallet_test.go` 91–92）；403 循环（109–120）**未**含新路由。实现确为 `requirePermission`，属测试缺口。E-002/A-002「401/403 覆盖新路由」对 403 证据不足。
  2. `TestWalletErrorCodesCataloged`（`wallet_test.go` 318）未列入 `WALLET_USER_AUTO_ONLY`；冻结集 + `TestErrorCatalogCoversFrozenCodesExceptInternal` 仍覆盖该码。
  3. `03-audit.md` 信息就绪表仍写 I-001~I-003 open；权威表已 verified。
  4. `goal-tree.md` GOAL-020 `progress` = 0/5，`00-meta` = 4/5。progress 非门禁，但关门前宜同步。
- **证据**：上列路径；`03-audit.md` 13–20 行；`goal-tree.md` 40、69 行；`00-meta.md` `progress: 4/5`。

## 必改项汇总

| ID | 级别 | 一句话 | 放行条件 |
|----|------|--------|----------|
| F-001 | required / med | UNIQUE 冲突重读须 `created=false`，且不重复 `wallet.account-create` | 改代码 + 冲突重读测试可核对 |
| F-002 | required / med | by-owner 调账不得留下「已开户、无 account-create」 | 审计时点或校验顺序修正 + 失败路径测试 |

F-003～F-005 为 recommended，不单独阻断关门。

## 与既有意见异同（A-001 / A-002）

| 条目 | 异同 |
|------|------|
| A-001 self · 立项 pass | 同意。五件套、parent、I-00N 登记、S1–S5 路线图成立。本条不回溯立项。 |
| A-002 self · S2–S4 pass、无 finding | **部分不同意**。门禁键、手动 409、错误码三处一致、路由前缀修正、S4 不暂挂、顺序路径测试——同意。A-002 将 TestGetOrCreateUserAccount 记为「冲突重读 + created 准确 ✅」**过宽**：该测试未打 UNIQUE 冲突分支；冲突分支 `created` 仍为 true（F-001）。调账双审计与 403 新路由的测试覆盖亦不足（F-002、F-005）。 |
| P-004.2 | **不构成** verdict 相反（本条 `conditional`，非 `fail`）。A-002 未主张「F-001/F-002 主题不必改」。开放 required 由 /govern 汇总响应，不得因 A-002 pass 静默放行 S5。 |

## 结论 + 关门放行条件

**verdict: conditional。** 自动开户主路径（顺序 get-or-create、零余额、UNIQUE 兜底代码、读/调账门禁、手动 user 409、错误码契约、无迁移/无新权限、前端去掉 user 选项且不伤展示）有证据。data 门禁上 **账户创建审计的幂等与完备** 未达到 D-001 / 成功标准所写：并发冲突会重复审计（F-001），非法调账会开户而不审计（F-002）。

**不得无条件将本目标标 `done`。** 关门最低条件：

1. F-001、F-002 按 P-003 合法闭合（`fixed` 且可核对，或用户书面 `accepted-residual` / `user-overruled` 并写清范围与复审触发）。
2. 闭合留痕写在本目标 `03-audit` 响应节或 `02-execution`；口头不算。
3. 本意见不修改 `status` / `progress` / 方案正文 / goal-tree。响应与是否关门走 **/govern**。

建议编排器下一步：先修 F-001/F-002（冲突分支复位 `created` + 调账失败审计时点），补测试后响应本 A-003；F-003/F-004/F-005 可同批或列为残余。

## 声明

本意见 `source: independent`，保证等级 **L0**。不修改目标 `status`/`progress`，不关闭检查点，不改方案或代码。响应、finding 闭合与关门由用户通过 `/govern` 处理。
