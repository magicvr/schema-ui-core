---
doc_type: goal-attachment
id: r3-gap-classification
parent: GOAL-004-r3-industry-comparison
status: draft
created: 2026-09-10
updated: 2026-09-10
version: 0.1.0
---

# R3 缺口分类表（C3）

覆盖 R1 入册的 12 条 closed-VP residual（`RES-*`）+ R2 的 4 条缺口候选（`G-001`～`G-004`），逐条给出现状、证据、分类、影响的路线图行与复审触发。

- **分类词表**（R3 D-001 冻结项 5）：仅 `现在修` / `仍 gated` / `接受残余` / `明确不做`
- **裁决约束**（R3 D-001 用户裁决 A/B）：**「现在修」一律另立**，本 VP 只登记；residual 处置继承原 VP 留痕，**新增或扩大残余须用户逐条书面接受**
- **路线图正文**不在本轮改（R4 交 `/vision` editorial）

## 1 · R1 入册 12 条 residual

| ID | 现状（本轮复核） | 证据（精确锚点） | 分类 | 影响的路线图行 | 复审触发 |
|----|------------------|------------------|------|----------------|----------|
| RES-013-migrator | 无产品级 SQLite→PG 搬运器；`apps/api/cmd` 只有 `dbgdump`/`e2e-pgset`/`otlp-sink`/`schema-ui`/`server`，无迁移工具 | `apps/api/cmd/` 目录清单（本轮枚举） | **明确不做**（非本 VP 目标，也非路线图空白）〔原 VP 记录为「接受残余/仍非目标」；维持非目标，不新接受残余〕 | 路线图无对应 A 序列行；如需产品化须另立 | 出现真实部署迁移需求，或 PG 成为唯一支持方言 |
| RES-014-migrator | 无产品级本地盘→对象存储搬运器（同上，无对应 CLI） | `apps/api/cmd/` 目录清单；ObjectStore 端口 `kernel/objectstore.go:108` 只提供 CRUD | **明确不做**（同 RES-013） | 同上 | 出现真实对象存储搬迁需求 |
| RES-015-otlp-sink | `cmd/otlp-sink` 是**显式声明的极简 sink**：任意路径接收 POST、计数并逐行打印，文档头写明「NOT a collector: no parsing, no semantics, no storage. Do not use in production.」 | `apps/api/cmd/otlp-sink/main.go:1`–`8`（包注释含上述原文） | **接受残余**（显式 endpoint 才导出；sink 仅用于证据与排障）〔VP-015 F-003 已留痕〕 | A 序列无可观测栈选型行；由 R4 草案按需登记 | 需要真实 collector/后端存储时另立 |
| RES-015-metrics | 指标分母为 HTTP 请求计数/时延 + 启用模块 gauge + build info；**Store/对象/Job 指标仍未纳入** | `apps/api/internal/obs/observer.go:62`–`64`（`httpRequestsTotal`/`httpRequestDuration`/`kernelModulesEnabled`）、`:73`（`buildInfo`） | **仍 gated**（未触发；不被业界对照迫使——行 2.3/2.4 均判「明确不做」，见 [industry-comparison.md](industry-comparison.md)） | 路线图 A 序列可观测行（R4 草案登记为候选下一拍） | 需要容量/性能归因，或出现 Store/对象/Job 生产事故归因需求 |
| RES-016-revoke | 刷新令牌**已可撤销**（旋转 + 守卫式 UPDATE + `ErrTokenRevoked`）；但访问令牌（JWT）无立即失效通道 | `apps/api/internal/auth/auth.go:75`（`RevokeRefreshToken`）、`:85`（用户级全撤销）、`:304`–`305`（已撤销即拒） | **接受残余**（访问令牌短 TTL + 刷新撤销为当前合同）〔VP-016 I-016-005 已留痕〕 | 安全相关行归 VP-009 程序；架构主线无改动 | 出现「必须即时踢下线」合规要求，或访问 TTL 被拉长 |
| RES-016-mfa-wrap | **R1 记载已过期**：现行代码已实现「previous 解密 + 成功 TOTP 后惰性重包」；仍存窄口径 = 恢复码路径不重包、无启动批量重包、无主动轮换重包 | `apps/api/modules/mfa/service.go:57`–`60`（`prevKey`）、`:72`（`NewService(..., previousSecret)`）、`:151`+`:165`（Verify 路径解密后重包）、`:246`+`:259`（Confirm 路径）、`:335`+`:353`（Disable/RotateRecovery 路径）、`:329`–`333`（**恢复码分支提前返回，不重包**）、`:450`/`:466`（`decryptSecret`/`maybeRewrap`）；接线 `apps/api/internal/composition/composition.go:395`（传入 current+previous）；来源注释 W11 F-004 | **`接受残余`（用户 2026-09-10 再裁决 · 修正后窄口径）**：范围 = 「仅凭恢复码完成的路径不重包 + 无启动批量重包 + 无主动轮换重包」（见 §6） | RT-K03 行；`docs/vision/roadmap.md:223`、`roadmap.md:16`、`charter.md:73`、`workspaces.md:30,58`、VP-016 正文 `:115` 仍写旧表述 → **R4 文档卫生须同步修正**（归 G-005） | KMS/密钥波次；或出现「恢复码路径导致 MFA 不可解」的证据 |
| RES-021-harness | 进程级停机 harness 仍 `!windows`；compose stop 以 linux CI 核销 | R2 [矩阵行 11](../../GOAL-003-r2-as-built-matrix/attachments/as-built-matrix.md)；`internal/composition/shutdown_drain_test.go`（无 `!windows` 的 HTTP/Job 子测） | **接受残余**〔VP-021 V-F083 已留痕〕 | RT-D02 行保持 delivered | CI 失败，或 Windows 上出现停机回归 |
| RES-026-redis | 无 Redis 依赖（`go.mod` 直接依赖未含），接缝与触发条件已冻结；本期**不消耗** trigger | `apps/api/go.mod:5`–`22`；`docs/architecture/cache-redis-seam-and-track.md:62`–`79`（§2.6 接缝声明） | **仍 gated** | RT-Q03 / RT-Q05 保持 gated | 多实例部署或跨实例一致视图需求出现 |
| RES-028-broker | 无外部队列依赖；进程内 EventBus 为唯一运输 | `apps/api/go.mod:5`–`22`；`apps/api/kernel/eventbus.go:43` | **仍 gated** | RT-Q05、A3 | 跨进程可靠消息需求出现 |
| RES-030-keyfile | bot master key 文件与 DB 同目录（运维面限定）〔原 VP 已 `accepted-residual`〕 | 原 VP-030 R-009 留痕（本轮未重开） | **接受残余**（继承原接受，不重新裁决） | RT-M03 行保持 delivered | KMS 波次启动 |
| RES-T03-tz | 迁移与写入仍未使用 `timestamptz`（`apps/api` 内 `*.sql` 中 `timestamptz` 出现 0 次） | 本轮全仓 `*.sql` 检索：`timestamptz` 计数 = 0 | **现在修·但另立**（本 VP 不改 schema；建议作为下一拍候选登记） | 路线图 A 序列候选行（R4 草案建议） | R4 草案 editorial 后由 `/vision` 决定是否立项 |
| RES-P04-pool | **池化已交付**：SQLite 文件库默认小连接池（4），内存库仍 1；文档现状锚点仍写 `MaxOpenConns=1`（G-002） | `apps/api/internal/store/store.go:29`（`sqlitePoolDefault = 4`）、`:104`–`114`（内存库单连接分支）；文档锚点见 [r4-doc-hygiene-anchors.md](r4-doc-hygiene-anchors.md) G-002 | **明确不做**（读写分离/replica 仍 gated）+ **文档卫生**（归 G-002 / R4） | RT-P04 保持 trigger-gated，但**现状锚点须在 R4 修正** | 多实例或 PG 之后 |

## 2 · R2 缺口候选 G-001～G-004

| ID | 现状（本轮复核） | 证据（精确锚点） | 分类 | 影响的路线图行 | 复审触发 |
|----|------------------|------------------|------|----------------|----------|
| G-001 | `docs/architecture/overview.md` 的「当前阶段（现时）」过期：Charter 版本、产品树（`web/` FastAPI）、组合编排与工作区清单停在早期 | `docs/architecture/overview.md:64,71,72-73,74-79`（详见 [r4-doc-hygiene-anchors.md](r4-doc-hygiene-anchors.md) G-001） | **现在修（文档）** → 按裁决 A 落 **R4 文档卫生**（判据 4） | 架构文档权威面；不改任何架构结论 | R4 执行后由用户 editorial 确认 |
| G-002 | roadmap 现状锚点与 RT-P04/RT-D02 表述滞后或自相矛盾（`MaxOpenConns=1`；「无明确 drain 合同」与 delivered 并存） | `docs/vision/roadmap.md:104,138,209` | **现在修（文档）** → R4 文档卫生 + `/vision` editorial（`docs/vision/**` 须与草案同一事务或紧随冻结） | RT-P04 现状锚点、RT-D02 行内表述 | 同上；不得把「已池化」扩写成读写分离已交付 |
| G-003 | `cache-redis-seam-and-track.md` §2.6 只描述旧三元组（`Allow`/`Record`/`Clear`），未覆盖现行 `AllowRecord`/`Reserve`/`Cancel` 原子面 | 文档 `:68,73`；端口 `apps/api/kernel/ratelimit.go:40,44,52,65,71,76,79` | **现在修（文档）** → R4 文档卫生（架构文档，不属 `docs/vision/**`） | 接缝文档须能指导触发后替换；**不冻结** Redis 实现细节 | 触发 RT-Q03/Q05 立项前必须已修正 |
| G-004 | `assembly.NewAuthenticator` 为公开工厂但返回 internal 指针且只接 current 单密钥（对比组合根注入 current+previous） | `apps/api/assembly/assembly.go:26,37`；组合根双密钥 `apps/api/internal/composition/composition.go:254`–`259`；单密钥形态 `apps/api/internal/auth/auth.go:151` | **明确不做（B+ 边界内）+ 登记为 R4 草案的分发面候选**（用户 2026-09-10 裁决） | 路线图 B+ 公共面行（R4 草案登记候选） | 出现第一个包消费者需要 previous 密钥，或 B+ 面成为分发契约 |
| G-005 | R1 记载的 `RES-016-mfa-wrap` 表述过期：多处公开文档仍写「`admin.mfa` wrapping 不随 JWT previous 重包」，与现行代码（W11 F-004 惰性重包）不符 | 文档：`docs/vision/charter.md:73`、`docs/vision/roadmap.md:16,223`、`docs/vision/workspaces.md:30,58`、`docs/vision/plans/VP-016-key-rotation-and-backup.md:115`；代码：`modules/mfa/service.go:57-60,72,151,165,246,259,335,353` | **现在修（文档）** → R4 文档卫生 + `/vision` editorial（`docs/vision/**` 须与草案同一事务或紧随冻结） | RT-K03 行、VP-016 关门记录表述 | R4 执行后由用户 editorial 确认 |
| G-006 | R2 矩阵 v0.2.0 仍有 31 处锚点指向紧邻注释/空行/邻行（主张本身经 A-002 复核属实），其中 1 处为**主题性错锚**：行 10 引 `internal/config/config.go:493`，该行实为 `CacheMaxEntries`，主张（metrics/traces 默认关闭）对应 `:495`/`:498` | 清单见 [r3-asbuilt-anchors.md](r3-asbuilt-anchors.md) §D（本轮只读取证逐条复核）；主题性错锚：`apps/api/internal/config/config.go:493` vs `:495`/`:498`（本轮亲自复核） | **现在修（文档）** → R2 矩阵 v0.3.0 锚点校正（P-002 允许的只读断言/文档校正，不改任何主张或分类） | 判据 1 的证据精度；不影响路线图行 | 校正提交后由 R4 复核；若发现任一主张实质变化则须回到 A-00N 响应流程 |

## 3 · 需用户裁决的项（P-004）

| # | 项 | 选项 | 本条建议 | 理由 |
|---|----|------|----------|------|
| 1 | **G-004** | `明确不做` / `接受残余` / `现在修（另立）` | **`明确不做`（B+ 边界内）+ R4 草案登记为分发面候选** | ✅ 用户 2026-09-10 已裁决 |
| 2 | **RES-016-mfa-wrap** | `接受残余`（窄口径） / `现在修（另立）` | **`接受残余`** | ✅ 用户 2026-09-10 已裁决（按修正后窄口径；原文见 §6） |
| 3 | **RES-013 / RES-014 搬运器** | `明确不做` / `接受残余` | **`明确不做`（维持非目标）** | ✅ 用户 2026-09-10 已裁决 |
| 4 | **RES-T03-tz** | 登记为下一拍照候选 / `明确不做` | **登记为下一拍候选（由 `/vision` 在 R4 editorial 决定）** | ✅ 用户 2026-09-10 已裁决 |

未列入裁决的项：RES-015-otlp-sink、RES-016-revoke、RES-021-harness、RES-030-keyfile 为**继承原 VP 已有留痕**的接受残余（裁决 B：不重新裁决）；RES-015-metrics、RES-026-redis、RES-028-broker 为**仍 gated**（不消耗 trigger）；RES-P04-pool 为**明确不做 + 文档卫生**；G-001～G-003、G-005、G-006 为**现在修（文档）→ R4**。

## 4 · 分类统计

| 分类 | 条数 | 条目 |
|------|------|------|
| 现在修（文档） | 6 | G-001、G-002、G-003、G-005、G-006（+ G-002 含 roadmap editorial） |
| 现在修（另立） | 1 | RES-T03-tz（登记为下一拍候选） |
| 仍 gated | 3 | RES-015-metrics、RES-026-redis、RES-028-broker |
| 接受残余 | 5 | RES-015-otlp-sink、RES-016-revoke、RES-021-harness、RES-030-keyfile（继承原留痕）+ RES-016-mfa-wrap（用户 2026-09-10 按修正后窄口径再裁决） |
| 明确不做 | 4 | RES-013-migrator、RES-014-migrator、RES-P04-pool、G-004（用户 2026-09-10 裁决） |

合计 19 条 = R1 入册 12 + R2 候选 4 + 本轮新发现 3（G-005 过期表述、G-006 锚点漂移、RES-016-mfa-wrap 记载修正归入原条目）。每类均有证据锚点与复审触发，无「感觉该优化」条目。

## 5 · 边界

- 本表不授权任何代码/端口/Profile 改动；分类为「现在修（文档）」的条目只在 R4 的文档卫生范围内执行（G-006 为 R2 矩阵自身的锚点校正，属只读断言/文档修正，不改任何主张或分类）。
- 本表不关闭任何 finding，不改 `docs/vision/roadmap.md`；`docs/vision/**` 的改动须与 R4 草案同一事务或紧随 `/vision` 冻结。
- 全部待裁决项已由用户于 2026-09-10 裁决（RES-016-mfa-wrap 经修正前提后于同日再裁决），无遗留裁决项。
- G-006（R2 矩阵锚点校正）在 C4 之前或 R4 执行均可；执行时只改锚点文字，不改主张、分类或 verdict。

## 6 · 中途发现与再裁决留痕（RES-016-mfa-wrap）

本轮只读取证子代理报告 `RES-016-mfa-wrap` 与 R1 记载**不一致**（一致 11 / 不一致 1 / 无法核对 0），编排器随后亲自复核源码确认：

| 项 | R1 记载 | 现行代码 | 结论 |
|----|---------|----------|------|
| 语义 | 「`admin.mfa` wrapping **不随** JWT previous 重包」 | previous 秘钥已参与解密，且成功 TOTP 后**惰性重包**（W11 F-004）：`modules/mfa/service.go:57-60,72,151,165,246,259,335,353`；接线 `composition.go:395` | R1 记载**过期** |
| 仍存的窄残余 | — | 仅凭恢复码完成的路径提前返回、不重包（`:329`–`333`）；无启动批量重包；无主动轮换重包 | 残余范围**显著小于** R1 记载 |

处置：

1. 用户原裁决（「接受残余」）所依据的**前提**与现行事实不符，故本表未直接沿用原措辞，而是提出**修正后的窄口径**重问；用户于 2026-09-10 **再次裁决 `接受残余`**。范围 = 上表「仍存的窄残余」（仅凭恢复码完成的路径不重包 + 无启动批量重包 + 无主动轮换重包），复审触发 = KMS/密钥波次，或出现「恢复码路径导致 MFA 不可解」的证据。
2. 旧表述仍留在多处公开文档（`docs/vision/charter.md:73`、`docs/vision/roadmap.md:16,223`、`docs/vision/workspaces.md:30,58`、`docs/vision/plans/VP-016-key-rotation-and-backup.md:115`），登记为 **G-005**，归 R4 文档卫生 + `/vision` editorial。
3. 残留原文（VP-016 关门时点）**不回改**：VP-016 正文与 VRev-036 保留其历史时点表述，仅由 R4 在现行文档面加注/修正。

## 7 · R2 矩阵锚点漂移（G-006）

只读取证复核发现 R2 矩阵 v0.2.0 尚有 31 处锚点指向紧邻注释/空行/邻行，另有 7 处「不算错行但可再精确」。**主张本身未受影响**（A-002 已独立复核 17 面 + W1 的主张真实性），唯一主题性错锚为行 10 的 `internal/config/config.go:493`（实为 `CacheMaxEntries`；主张对应 `:495`/`:498`，本轮亲自复核）。

处置：登记为 **G-006（现在修·文档）**，由 R4 或在 C4 前以矩阵 v0.3.0 一次性校正；校正只改锚点文字，不改任何主张、分类或 verdict。完整清单见 [r3-asbuilt-anchors.md](r3-asbuilt-anchors.md) §D 与 [r3-anchor-commands.txt](r3-anchor-commands.txt)。
