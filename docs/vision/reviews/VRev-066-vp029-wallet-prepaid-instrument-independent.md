---
doc_type: vision-review
id: VRev-066
status: active
source: independent
created: 2026-09-02
updated: 2026-09-02
version: 0.1.0
parent: null
---

# VRev-066 · VP-029 钱包预付资金凭证 · 独立激活审视（2026-09-02）

| 字段 | 值 |
|------|-----|
| source | independent |
| auditor | grok-build（grok-4.6 · `/vision-audit`） |
| scope | `VP-029-wallet-prepaid-instrument` 意图完备 / Charter `@0.4.0` 对齐 / 退出判据可核验性 / P-005 / V-F110（资金哈希 · 并发双花 · 主体与 Admin 用户隔离）/ Admin 类 freshness（`29727510` → `b5c39dfb`） |
| audit_type | vision-plan（意图 / 激活就绪） |
| verdict | pass |
| 建议 class | editorial |
| open required | 0 |

## 范围与结论

只读核对（本报告写入前未改 Charter / VP / Goal status）：P-006 / `alignment.md`、Charter `@0.4.0`、[VP-029](../plans/VP-029-wallet-prepaid-instrument.md)（`planned` v0.1.0）、[VP-030](../plans/VP-030-telegram-channel-runtime.md)、[VP-031](../plans/VP-031-digital-offer-entitlement.md)、[VRev-065](VRev-065-c-end-paid-services-planned-self.md)、roadmap / revisions VR-059 / reviews.md（open required = 0）、workspaces.md，以及代码：

| 材料 | 核验用途 |
|------|----------|
| `apps/api/modules/wallet/migration/migration.go` | `owner_type` CHECK、账本 `entry_type`、三余额恒等 |
| `apps/api/internal/composition/composition.go` L567–572 | `ownerExists` = `UserByID`（W13 F-012） |
| `apps/api/internal/handler/wallet.go` L55–164 | 开户/adjust 拒绝非 `users` 行 |
| `apps/api/modules/wallet/store/repository.go` `GetOrCreateUserAccount` | 固定 `owner_type=user`；UNIQUE 冲突走新事务 |
| `apps/api/modules/wallet/store/concurrent_test.go` | 并发 get-or-create 已有回归 |
| `apps/api/modules/authsession/recovery.go` / `service_credentials.go` | 既有「哈希存、明文不落库」先例 |
| `apps/api/kernel/profile.go` | `admin.wallet` 已在 admin 默认集；本 VP 意图不改默认集 |
| workspace-009 Root `vp008_go_status` | `active`（2026-09-01 W16 恢复） |

对照 [VRev-065](VRev-065-c-end-paid-services-planned-self.md)（self · `pass`）只作参考，**不盲从**。本审视同时核销 V-F110 点名的三项（资金哈希、并发双花、主体隔离），并做 Admin 类 freshness。

**总判：pass（0 open required）。** 单愿景与 `vision_ref` 精确匹配；结构选型（Admin 功能增量、不重开 VP-011、不是支付域）成立；退出 1–7 方向可判定；Admin 类 freshness **PASS**，不暂挂 `go`。本 `pass` 允许**用户确认后激活 VP-029 并开 `workspace-029-wallet-prepaid-instrument`**。它不构成「卡密已可售 / 核销已上线 / 主体接缝已交付」的任何宣称。

用户本轮书面：「审视 vp-029 是否合理，没问题的话激活 vp-029，然后交编排器开设工作区」。

## 逐条审视

### 1. 与 Charter 0.4.0 对齐

**pass。** 未触发 P-006 §6.3 冲突谓词。

| Charter 条款 | 独立判定 |
|--------------|----------|
| 非目标「不建设特定业务领域的终端产品」；钱包列为后续 VP **候选能力** | 本 VP 扩展已交付的 `admin.wallet` 资金通道，不把 Telegram 产品写成成功条件 |
| 非目标「不预制 C 端 API 的业务逻辑」 | 交付通道无关主体接缝 + 模块级 `Redeem`；Telegram/Offer 分属 030/031。可选 C 端 HTTP 核销默认倾向「不做」 |
| 成功边界 #6 / H-002 | Admin 功能，不消耗 RT-Q03/Q05 trigger；C 端流量评估挂 VP-030。H-002 再确认是 VP-031 业务域门禁，不要求本 VP 重裁同进程 |
| `vision_ref` | `schema-ui-core-admin-foundation@0.4.0` 精确匹配 |

将卡密入金放在 Admin 功能、而不是路线图业务域「支付/结算」，与 roadmap 已有注记及用户 2026-09-02 书面切分一致。凭证 ≠ 支付网关。

### 2. 切分、硬前置与组合位置

**pass。** 与 VP-030（身份消费者）/ VP-031（资金原语消费者）的硬前置方向一致；本 VP 不得 import Telegram。三分支并行规则允许 Admin 交付 VP 与架构程序并存。VP-028 计划文件已 `closed` v0.3.0，架构交付位空闲（roadmap 第 28 行仍写 `active`，见 V-F113）。

不重开 VP-011：账本原语（`adjust` / `freeze` / `unfreeze` / `deduct_frozen`、三余额、不可变流水、幂等键）保持只读基线。凭证核销是新入金通道，不是第二套余额。

### 3. 退出判据可判定性

**pass。**

| 判据 | 独立判定 |
|------|----------|
| 1 主体接缝 | **可判定**：幂等 get-or-create、未登记不能开户、不写 `admin.users`、查询不依赖钱包已启（V-F109 已写入）。落点仍是 I-029-001 |
| 2 凭证生命周期 | **可判定**：哈希存储、一次性明文、导出/作废/过期拒绝、明文不进审计原文均可测。算法与码熵未冻，属 R1（V-F111），不使判据本身不可核验 |
| 3 核销原子且幂等 | **可判定**：「成功则账本与状态一致 / 重复不双记 / 并发双花 fail-closed」有明确负向证据形态。UNIQUE + 同事务 CAS 的具体合同属 R1 |
| 4 账本不变式 | **可判定**：既有钱包测试 + 新入金类型纳入 apply 表。现行 `entry_type` CHECK 已含 `adjust/freeze/unfreeze/deduct_frozen`（0033）；新类型 vs `adjust+ref_type=voucher` 由 I-029-002 冻 |
| 5 Admin 可操作 | **可判定**：协议页面 + 权限键 + 操作审计。权限键由 I-029-003 冻 |
| 6 边界 | **可判定**：未改 Charter / 默认 Profile / 支付网关 / Telegram 依赖 / 重开 VP-011 |
| 7 审计闭合 | **过程可判定** |

非目标充分：电商三件套、Stars/支付网关、完整 C 端 IAM、把 Bot 用户写入 `admin.users`、Redis/多实例、解禁 typed domain event，均已点名。

### 4. V-F110 三项（本审视主责）

#### 4.1 资金哈希

方向级已够用：只存哈希、明文仅生成时一次性出示、不进审计原文。本仓已有两类可抄先例，**不得混用**：

| 先例 | 形态 | 对本 VP |
|------|------|---------|
| 恢复码（`recovery.go`） | 6 位、sha256 at rest、常时比较 | **不可**当卡密默认：熵太低 |
| 服务凭证（`service_credentials.go`） | 高熵 token、`token_hash` 落库、原文永不持久化 | **可**作默认候选 |

bcrypt 适合口令、不适合批次生成/核销热路径。算法 / pepper / 字母表与长度 / 常时比较必须在 R1 冻结，见 V-F111。不升 required：判据 2 已要求「哈希存储 + 明文不落库」，缺的是实现合同而非方向。

#### 4.2 并发双花

方向级已够用。现行钱包已具备可复用零件：`UNIQUE(account_id, idempotency_key)`、账户 `version` CAS、并发 get-or-create 回归（`concurrent_test.go`；PG 失败 INSERT 必须换新事务，W9 F-001）。核销还缺一张**凭证行**的唯一键与状态机：`UNIQUE(code_hash)` + 同事务「未用→已核销 AND 入金」；并发失败者 fail-closed；同一 code 对同一主体重复 Redeem 不双记。属 R1，见 V-F111。

#### 4.3 主体与 Admin 用户隔离

方向级已够用，且与现行代码的张力**正是本 VP 要关的门**，不是计划缺陷：

| 现行事实 | VP 要求 |
|----------|---------|
| `ownerExists` = `authRepository.UserByID`（composition.go L569–572） | 承认已登记主体，**不**要求 `users` 行 |
| 开户/adjust 对未知 owner 回 `USER_NOT_FOUND`（wallet.go L129–134, L162–164） | 未登记主体不能开户；已登记主体可以 |
| `owner_type CHECK IN ('user','business','system')`；`GetOrCreateUserAccount` 写死 `OwnerUser` | 须冻结是否新增 `'subject'`（或等价）——写入 I-029-001，不得在 VP 层假装已选 |
| W13 F-012「禁止孤儿账本」 | 孤儿相对**主体登记表**，不再相对 `admin.users` |

V-F109（主体查询不依赖 `admin.wallet` 已启）仍约束 I-029-001 的落点选择：把主体表塞进仅随钱包模块加载的 persistence，会再次打断 030/031。见 V-F112。

### 5. P-005

**pass**（带 recommended 补项）。I-029-001～005 字段完整，required 三项均最晚 R1，未伪装已决。缺口：

1. **I-029-001** 未点名 `owner_type` 枚举与 `OwnerExistsFunc` 从 `UserByID` 迁到主体登记（V-F112）。
2. **缺 I-029-006 required（最晚 R1）**：哈希算法 / 码熵 / 常时比较 / 核销 UNIQUE+同事务（V-F111）。
3. **I-029-005**：若选 C 端 HTTP 自助核销，限流评估变成**本 VP** 义务，不得推到未激活的 VP-030。默认「模块 API 必做、HTTP 可选」合理。

不得把主体落点、entry_type、权限键、哈希算法在 VP 层假装已选。

### 6. Admin 类 freshness（VP-008 `go` 消费有效性）

VP-008 强制 freshness 的对象是**后续业务 VP**。VP-029 是 Admin 功能（扩展 `admin.wallet`），不是 Tier D。按激活门闩做 Admin 类复核：

| 项 | 结论 |
|----|------|
| 原 `go` 候选 | `ed99e88`（2026-08-10，clean）；解锁 scope = 标准业务模块框架能力，不是钱包卡密本身 |
| VP-008 go 现行 | **active**（workspace-009 Root `vp008_go_status: active`；最近恢复 = 2026-09-01 W16 S6） |
| 本轮基线 | `29727510`（VP-027 关门 / VP-028 激活基线 · VRev-064） |
| 现行 HEAD | `b5c39dfb`（`docs(vision): 新增 VP-029~031 愿景计划及 VRev-065 自审记录`） |
| 工作树 | HEAD 干净 |
| 五域 `29727510`→`b5c39dfb` | **零变更**：`kernel/profile.go` / `apps/web/src/protocol/upstream` / `go.mod`+`go.sum` / 根 `package.json` / `pnpm-lock.yaml` |
| 区间代码 | VP-028 EventBus 已审结目交付 + VP-009 W16/W17（JWT/CORS/httpOnly refresh cookie）+ 上传 accept 校验。**未改**钱包模块 / Profile 默认集 / 协议 pin |
| 共同门禁 | 认证/授权/fail-closed 语义无新的暂挂条件 |
| F-007 residual | 上传授权深度仍 deferred（owner = VP-008 lead）。本 VP 不得借卡密面扩张上传授权 |
| 本 VP 意图是否改 Profile / 模块矩阵 / Manifest / 协议 pin | **意图否**（钱包增量 + 主体持久化贡献；不进 `mvp`/`admin` 默认集） |
| 数据库意图 | **意图会变**（主体表 + 凭证表 + 可能的 `owner_type`/`entry_type` CHECK）。属本波分母，不是激活前已存在的失效 |
| **结果** | **PASS（Admin 激活）**。不消费 Tier D 解锁；不暂挂 `go` |

`consumer_vp` = VP-029；`last_freshness_review_at` = 2026-09-02；`next_freshness_review_trigger` = 若实施改变共同门禁 / Profile 默认集 / 模块矩阵 / Manifest 装配语义 / 协议 pin，或激活 VP-031（业务域 freshness + H-002）则重做。

### 7. 不构成 fail / 不新开 required 的诚实边界

- 主体表 / 凭证表 / `owner_type` CHECK **会变**——属本波分母，R1 冻结后实施。
- roadmap 第 28 行与 `workspaces.md` 的 workspace-028 行落后于 VP-028 `closed` / Root `done` 4/4。这是组合索引卫生，不否定 VP-029 意图（V-F113）。
- README / Charter「当前组合焦点」长段仍偏历史，不在本 scope 强制重写。

## Findings

### 必改（required）

无。

### 建议（recommended）

| id | level | 状态 | 一句话 |
|----|-------|------|--------|
| V-F111 | recommended | open | Root/VP R1 补 I-029-006：哈希算法与码熵、常时比较、核销 UNIQUE+同事务双花合同；I-029-005 若选 HTTP 则本 VP 做限流评估 |
| V-F112 | recommended | open | I-029-001 扩写：`owner_type` 枚举 + `OwnerExists` 迁到主体登记（禁止回退 `UserByID`）；落点不得绑死钱包已启 |
| V-F113 | recommended | open | 激活记录留下 Admin 类 freshness（`29727510`→`b5c39dfb`，不暂挂 `go`）；核销 V-F110；同步 VP-028 已 closed 的组合索引 |

#### V-F111 · recommended · 资金哈希与并发双花的 R1 合同

- level: `recommended`；status: open；severity: medium
- impact: 若 R1 不冻算法与唯一键，实施可能用恢复码那套 6 位 sha256（可枚举）或 bcrypt 批次（热路径不可用），或把核销做成「先改状态再入金」的两事务窗口。
- finding: lead Root 信息表增补 **I-029-006 required（最晚 R1）**：① 哈希算法（默认候选 = 高熵码的 SHA-256 或 HMAC-SHA256+pepper；**禁止**把 6 位恢复码或 bcrypt 当卡密默认）；② 码字母表与长度（熵下限）；③ 核销常时比较；④ `UNIQUE(code_hash)`（或等价）+ 同事务「未用→已核销 AND 账本入金」，并发失败者 fail-closed，重复 Redeem 不双记。I-029-005：若选 C 端 HTTP 自助核销，RT-Q05 精神的限流评估在本 VP 完成，不得推到未激活的 VP-030。
- evidence: VP-029 判据 2/3；`recovery.go` L4–5 vs `service_credentials.go` L17–18；`wallet/store/concurrent_test.go`；W9 F-001（PG 失败 INSERT 须换事务）。
- closure: VP P-005 表 + Root `00-meta` 含 I-029-006；不要求本 Review 落盘时已经有算法答案。
- 建议 class: `editorial`

#### V-F112 · recommended · 主体隔离落地合同

- level: `recommended`；status: open；severity: medium
- impact: 若只扩 `owner_id` 仍走 `UserByID` 门闩，C 端主体无法开户；若 `owner_type` 不扩 CHECK，主体行插不进；若主体 persistence 随 `admin.wallet` 启用才加载，V-F109 回潮。
- finding: I-029-001 在 R1 必须同时冻结：主体落点（薄模块 / `authsession` / 钱包表——公共契约通道无关）+ **`owner_type` 是否新增取值** + **`OwnerExistsFunc` 改为「已登记主体」，禁止回退 `UserByID`**。W13 F-012 重解释为相对主体登记表禁孤儿。查询/get-or-create 不依赖 `admin.wallet` 已启。
- evidence: `composition.go` L567–572；`wallet.go` L129–134；`migration.go` L24 `CHECK (owner_type IN ('user','business','system'))`；VRev-065 V-F109。
- closure: I-029-001 问题陈述扩写进 VP 与 Root 信息表。不要求本 Review 已选定落点。
- 建议 class: `editorial`

#### V-F113 · recommended · 激活 freshness 留痕 + 组合索引卫生

- level: `recommended`；status: open；severity: low
- impact: 若激活只写「开区」而不点名非 Tier D、五域零变更、go 未暂挂、基线 `29727510`、HEAD `b5c39dfb`，后续读者会把卡密误读成业务域解锁。roadmap 第 28 行仍标 VP-028 `active` 会让人以为架构交付位被占。
- finding: 激活时在 VP 短史或 lead Root D-001 写入上表 freshness；核销 V-F110；`roadmap.md` 第 28 行与 `workspaces.md` workspace-028 行对齐 VP-028 `closed` / Root `done` 4/4。
- evidence: VP-028 计划文件 `status: closed` v0.3.0；workspace-028 Root `status: done` `progress: 4/4`；roadmap L47 仍写 `active`；workspaces.md L43 仍写 `active · 0/4`。
- close requirement: D-001 或 VP 激活短史含 freshness 表；组合索引与 VP-028 实际 status 同拍。不要求重开 VP-008/028。
- 建议 class: `editorial`

## 门禁含义

- 本 scope **open required = 0**。
- **允许**：用户确认后激活 VP-029、开新 delivery 工作区、按 V-F111/112/113 写 Root 纲领 / I-00N / freshness / 索引。
- **禁止**：把本 `pass` 写成卡密已可售、核销已上线、或主体接缝已交付；重开 VP-011；把 VP-029 当 Tier D 消费订单/支付解锁；为卡密把新模块打进 `mvp`/`admin` 默认集；把 Bot 用户写入 `admin.users`。

## 结论

**VP-029 是否合理、是否可以激活并开区？可以。** 用户确认后：`/vision` 激活 VP、`/govern` 开 `workspace-029-wallet-prepaid-instrument`。开区后第一件事是 R1 合同冻结（I-029-001～003 + I-029-006），不要直接改钱包 DDL。

## 声明

本意见 `source: independent`，**不直接修改** Charter / VP / Goal status，不自行闭合任何 finding。required finding 的响应由 `/vision` 追加在正式 VRev 报告中；原 verdict 与 finding 原文不得改写。实施工作交 `/govern`。

## `/vision` 响应（2026-09-02 · V-F110/111/112/113 → fixed）

本 VRev 落盘与 VP-029 激活/开区同事务完成（用户本轮书面指令：「审视 vp-029 是否合理，没问题的话激活 vp-029，然后交编排器开设工作区」）：

- **V-F110 → `fixed`**：本独立审视即 VRev-065 所请求的激活前 independent Review；三项方向级可判定。细节合同下沉 V-F111/112。
- **V-F111 → `fixed`**：VP-029 v0.2.0 P-005 增补 **I-029-006 required（最晚 R1）**（哈希算法 / 码熵 / 常时比较 / UNIQUE+同事务双花）；I-029-005 增写「若选 HTTP 则本 VP 做限流评估」。投影 lead Root `GOAL-001-wallet-prepaid-instrument` 信息表。
- **V-F112 → `fixed`**：I-029-001 问题陈述扩写 `owner_type` 枚举、`OwnerExistsFunc` 禁止回退 `UserByID`、不依赖钱包已启、W13 F-012 相对主体登记表。投影 Root 信息表。
- **V-F113 → `fixed`**：VP-029 激活短史 + Root `D-001-workspace-root-establishment` 写入 Admin 类 freshness（`29727510` → `b5c39dfb`；PASS，不暂挂 `go`）；roadmap 第 28 行 / `workspaces.md` workspace-028 行对齐 VP-028 `closed` / Root `done` 4/4。

原 verdict（pass）与 finding 原文未改写；本响应为 append-only 补充。
