---
doc_type: vision-review
id: VRev-035
status: active
source: self
created: 2026-08-22
updated: 2026-08-22
version: 0.2.0
parent: null
---

# VRev-035 · VP-016 意图完备 / 可行性 / 激活就绪（2026-08-22）

| 字段 | 值 |
|------|-----|
| source | self |
| auditor | `/vision` · 会话编排（grok-4.6） |
| scope | `VP-016-key-rotation-and-backup`（审视时 `planned` v0.1.0）意图完备、Charter 对齐、退出分母、与 roadmap A5 / RT-K03 / RT-P05 一致性、激活与开区就绪、VP-008 `go` 消费前新鲜度（架构类） |
| audit_type | vision-plan（意图 / 激活就绪） |
| verdict | pass |
| 建议 class | editorial |
| open required | 0 |

## 范围与结论

只读核对：`docs/architecture/principles.md` P-006、`docs/vision/alignment.md`、Charter `@0.2.0`、[VP-016-key-rotation-and-backup](../plans/VP-016-key-rotation-and-backup.md) v0.1.0、roadmap v0.35.0 A5 / RT-K01–K04 / RT-P05、VR-037、`apps/api/internal/auth`（HMAC-SHA256 单密钥签发/校验；服务凭证为 SHA-256 哈希的 opaque token）、`apps/api/internal/store/migrate.go`（SQLite `VACUUM INTO`）、VP-013 I-004（PG `pg_dump`/`pg_restore` 已交付）。未把 `planned` 读成已交付；本报告落盘时尚未改 VP status（激活与开区是用户本轮「通过后」的后续原子动作）。

**总判：pass（0 open required）。** 单愿景与 `vision_ref` 精确匹配；新 VP 承接架构 A5 的结构选型合法；退出分母与用户书面确认同构（JWT current+previous 轮换 + 既有备份上的轮换后恢复；不重做 dump）；方向足以激活并开新 delivery 工作区。两条 recommended 约束 Root 纲领/信息项、退出 1 措辞与 `go` 新鲜度留痕，不阻断激活。

## 核对事实

| 核对项 | 结论 | 证据 |
|--------|------|------|
| 单愿景 | **pass** | 唯一 active Charter `schema-ui-core-admin-foundation@0.2.0` |
| VP→Charter 机读 | **pass** | `vision_ref` 精确匹配 |
| 语义对齐 | **pass** | 可 fork 的 Go 后端内核能力（签名密钥轮换 + 恢复语义）；不把业务领域当成功条件；不改 Charter 非目标 |
| 最小完备 | **pass** | 意图、配置面、首波冻结、非目标、退出 1–6、邻接 VP、I-016-001～005、工作区表、短史均在 |
| 结构选型 | **pass** | 同愿景新纲领波次 → 新 VP；不重开 VP-012/013/014/015；不改 Charter；新 delivery 区（惯例 slug 见 D-001） |
| 与 A5 / RT-K03 / RT-P05 | **pass** | 轮换合同 = RT-K03；备份面明确消费 VP-013 dump 而非第二套实现；PITR / KMS 不进分母 |
| 退出分母有界 | **pass** | 明确排除热加载、`/readyz` 再扩、Admin 密钥页、DB/S3/种子密码轮换 |
| 配置面 | **pass** | 缺省单密钥；previous 为显式配置；未配置不 fail-closed 挡住 mvp/dev；重启生效 |
| 可行性 | **pass（工作量中、边界清）** | `SignAccessToken` / `ParseAccessToken` 今日单密钥 HMAC；接入点已知。服务凭证 `NewServiceCredentialToken` 为 SHA-256 哈希，**不**与 JWT secret 共用。SQLite `VACUUM INTO` 与 PG dump 已存在。工作是补 dual-key 校验与轮换后恢复剧本，不是换认证叙事 |
| 开放 VRev required | **pass** | 本报告前 open required = 0 |
| 过早交付主张 | **无** | `planned`、0 区；未主张驱动已写 |

## VP-008 `go` 消费前新鲜度（架构类）

VP-008 正文强制 freshness 的对象是**后续业务 VP**。VP-016 是架构分支，且自行把门闩写成激活前须复核。本表按该自设门闩做轻量复核，**不**把本 VP 误读成业务域解锁。

| 项 | 结论 |
|----|------|
| 原 `go` 候选 | `ed99e88`（2026-08-10，clean）；解锁 scope = 标准业务模块框架能力，不是密钥轮换 |
| 现行 HEAD | `57098c3`（`docs(governance): 修复文档、模板与审计台账中的无效相对链接及占位符`） |
| VP-009 | W1–W4 与 W6 done；W5 扫描 0 中高危未开子目标；无新的共享基架暂挂宣称 |
| VP-010 | W1–W13 done；`go` 无新暂挂 |
| Vision open required | 0 |
| F-007 residual | 上传授权深度仍 **deferred**（owner=VP-008 lead）。本 VP 不得借轮换面扩张授权 scope |
| 本 VP 是否改 Profile / 模块矩阵 / Manifest / 协议 pin | **意图否**。纯 auth 配置合同；若实施时证据显示改变，按消费有效性暂挂 |
| 复核结果 | **PASS（架构激活）**。不消费业务解锁 scope；不暂挂 `go` |

## 不构成 fail / 不新开 required 的诚实边界

1. 退出 1 写「重叠窗内 previous 可验 **或** 书面冻结为立即失效」。冻结表与意图 §1 已把 previous 可验写成默认交付。属措辞，不是方向空洞；激活时 editorial 收口。I-016-005 仍可作有界残余。
2. VP 建议的 Root R5 写成「激活后 freshness / 回归」，会与 VP 激活门闩混淆。激活后 Root 纲领须把 R5 收成双路径证据（默认单密钥 + 显式轮换/恢复），freshness 留在激活记录。
3. I-016-001～004 仍 `collecting`。键名、重叠窗时长、是否 `kid`、恢复剧本先后不是激活阻断；最晚 R1/R2/R3。现行代码显示服务凭证不与 JWT secret 共用，但是 R1 书面出局，不是本 Review 已 verified。
4. 本 `pass` 允许激活与开区，**不是** R1 轮换合同已冻结，也不是可以开始无设计地改 `ParseAccessToken`。
5. 无独立 Vision Review 不是 alignment 强制项。若用户要求交叉审视，另走 `/vision-audit`。
6. 用户本轮确认开区但未另写 slug。不把惯例 slug 写成用户已口头点名；开区记录须留痕推导。

## Findings

### V-F067 · recommended · Root 须写出纲领阶段，收口退出 1，并登记 I-00N

- level: `recommended`
- status: `open`
- severity: medium
- impact: 若不开区就改 JWT 解析，或把重叠窗、kid、恢复剧本、密钥集合混成未登记未知，A5 会在 R1/R2 才爆。退出 1 的「或」会让只交「换 env 重启、旧 token 全断」的人误读成可关门。
- finding: |
  VP 已建议 R1→R5 与 I-016-001～005，但 `planned` 无工作区。激活后 lead Root **必须**：
  1. 写出串行纲领（轮换合同冻结 → JWT 双密钥实现 → 轮换后恢复证据 → 默认单密钥仍可用 → 显式双密钥轮换路径 **与** 恢复路径证据）。R5 不得写成「激活后 freshness」。
  2. 把 I-016-001/002 登记为 required（最晚 R1）；I-016-003 required（最晚 R2）；I-016-004 required（最晚 R3）；I-016-005 保持 non-blocking。
  3. 激活时把退出 1 收口为：previous **默认可验**；立即失效只走 I-016-005 有界残余。
- evidence: VP-016 退出 1–4；I-016-001～005；冻结表「校验允许重叠窗」；`ParseAccessToken` 单密钥。
- closure: |
  Root `00-meta` 含 P-001 阶段表 + 上述 I-00N；VP 退出 1 作 editorial 修订。不要求本 Review 落盘时已经有键名/重叠窗答案。
- 建议 class: `editorial`

### V-F068 · recommended · 激活记录须留下架构类 freshness 结论，避免被读成已消费业务 `go`

- level: `recommended`
- status: `open`
- severity: low
- impact: VP-016 把门闩写成「激活前 freshness review」。若激活记录只写「开区」而不点名：本 VP 非业务域、不改 Profile 意图、F-007 不升格、现行无 `go` 暂挂，后续读者会把架构 A5 误读成已走 VP-011 那种业务消费。
- finding: 激活时在 VP 短史或 lead Root D-001 写入上表复核结论与候选/HEAD 指针。
- evidence: VP-008 §`go` 消费有效性（业务 VP）；VP-016 激活门闩；GOAL-007 D-001 原候选 `ed99e88`。
- close requirement: D-001 或 VP 激活短史含 freshness 表；不要求重开 VP-008。
- 建议 class: `editorial`

### 声明

本意见不直接修改 Charter / VP / Goal status。required finding 的响应由 `/vision` 追加在本报告中；原 verdict 与 finding 原文不得改写。

### 门禁含义

- Vision Review **open required = 0**。
- **允许**：用户确认后激活 VP-016、开新 delivery 工作区、按 V-F067/V-F068 写 Root 纲领与 freshness。
- **禁止**：把本 `pass` 写成 R1 已冻结或 JWT 解析已改；重开 workspace-015；把 VP-016 当业务 VP 消费 `go` 解锁 scope。

### 响应（2026-08-22 · `/vision` 激活与 `/govern` 开区）

不回溯改写原 verdict `pass` 与 finding 正文。

| finding | 闭合 | 证据 |
|---------|------|------|
| V-F067 | **fixed** | Root `GOAL-001-key-rotation-and-backup` 纲领 R1～R5（R5 = 双路径证据，不是「激活后 freshness」）；I-001～I-004 required collecting（对应 I-016-001～004）；I-005 non-blocking（I-016-005）。VP-016 v0.2.0 退出 1 改为 previous **默认可验**；立即失效只走 I-016-005。答案仍 open，登记已满足 close requirement。 |
| V-F068 | **fixed** | D-001 架构类 freshness：原 `go` `ed99e88`；HEAD `57098c3`；非业务解锁；不暂挂 `go`；F-007 不升格。 |

用户书面「对 VP-016 做意图审视；没有问题的话交 `/govern` 开区」已执行：VP `active`；lead `workspace-016-key-rotation-and-backup`（惯例 slug，D-001 留痕）；Root scaffold。本 scope **0 open required、0 open recommended**。
