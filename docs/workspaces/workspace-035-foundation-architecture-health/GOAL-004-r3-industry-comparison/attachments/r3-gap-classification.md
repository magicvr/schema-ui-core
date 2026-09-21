---
doc_type: goal-attachment
id: r3-gap-classification
parent: GOAL-004-r3-industry-comparison
status: draft
created: 2026-09-10
updated: 2026-09-10
version: 0.2.0
---

# R3 缺口分类表（C3）

覆盖 R1 入册的 12 条 closed-VP residual（`RES-*`）+ R2 的 4 条缺口候选（`G-001`～`G-004`）+ 本轮新增发现（`G-005`、`G-006`），共 **18 个唯一条目**，逐条给出现状、证据、分类、去向与复审触发。

- **分类词表**（R3 D-001 冻结项 5）：分类列**只允许** `现在修` / `仍 gated` / `接受残余` / `明确不做` 四值。执行状态、去向与范围限定一律写在「去向 / 部署」与「复审触发」列，**不得**写入分类列（v0.2.0 修正：v0.1.0 曾在分类列写「现在修（文档）」「已执行（文档）」等复合值，违反冻结词表）。
- **裁决约束**（R3 D-001 用户裁决 A/B）：**「现在修」一律另立**，本 VP 只登记；residual 处置继承原 VP 留痕，**新增或扩大残余须用户逐条书面接受**。
- **路线图正文**不在本轮改（R4 交 `/vision` editorial）。

## 1 · R1 入册 12 条 residual

| ID | 现状（本轮复核） | 证据（精确锚点） | 分类 | 去向 / 部署 | 影响的路线图行 | 复审触发 |
|----|------------------|------------------|------|-------------|----------------|----------|
| RES-013-migrator | 无产品级 SQLite→PG 搬运器；`apps/api/cmd` 只有 `dbgdump`/`e2e-pgset`/`otlp-sink`/`schema-ui`/`server` | `apps/api/cmd/` 目录清单（本轮枚举） | **明确不做** | 维持非目标；产品化须另立 | 无对应 A 序列行 | 出现真实部署迁移需求，或 PG 成为唯一支持方言 |
| RES-014-migrator | 无产品级本地盘→对象存储搬运器（同上，无对应 CLI） | `apps/api/cmd/` 目录清单；ObjectStore 端口 `kernel/objectstore.go:108` 只提供 CRUD | **明确不做** | 维持非目标；产品化须另立 | 同上 | 出现真实对象存储搬迁需求 |
| RES-015-otlp-sink | `cmd/otlp-sink` 是**显式声明的极简 sink**：任意路径接收 POST、计数并逐行打印；包注释写明「NOT a collector: no parsing, no semantics, no storage. Do not use in production.」 | `apps/api/cmd/otlp-sink/main.go:1`–`8`（含上述原文） | **接受残余** | 已在 VP-015 关门事务内以 recommended「F-003 文档化残余」闭合（[A-002](../../../workspace-015-observability/GOAL-001-observability/03-audit/A-002-independent-root-closeout.md)；`03-audit.md:33`：开放 required = 0） | A 序列无可观测栈选型行 | 需要真实 collector/后端存储时另立 |
| RES-015-metrics | 指标分母为 HTTP 请求计数/时延 + 启用模块 gauge + build info；**Store/对象/Job 指标仍未纳入** | `apps/api/internal/obs/observer.go:62`–`64`（`httpRequestsTotal`/`httpRequestDuration`/`kernelModulesEnabled`）、`:73`（`buildInfo`） | **仍 gated** | 不消耗 trigger；未受业界对照迫使（[industry-comparison.md](industry-comparison.md) 行 2.3/2.4 均判明确不做） | A 序列可观测行（R4 草案登记为候选下一拍） | 需要容量/性能归因，或出现 Store/对象/Job 生产事故归因需求 |
| RES-016-revoke | 刷新令牌**已可撤销**（旋转 + 守卫式 UPDATE + `ErrTokenRevoked`）；但访问令牌（JWT）**未实现**立即失效 | `apps/api/internal/auth/auth.go:75`（`RevokeRefreshToken`）、`:85`（用户级全撤销）、`:304`–`305`（已撤销即拒） | **明确不做** | ⚠ **不再记为「接受残余」**（v0.2.0 修正）：原记录 [Root I-005](../../../workspace-016-key-rotation-and-backup/GOAL-001-key-rotation-and-backup/00-meta.md) 现为 `collecting`，且注明「用户书面残余时才改变退出 1」——**不存在合法的用户书面残余接受**。当前形态是「立即失效未选」的设计后果，按未实现能力登记；若用户希望改为残余风险，须按裁决 B 逐条书面接受 | 安全相关行归 VP-009 程序；架构主线无改动 | 出现「必须即时踢下线」合规要求，或访问 TTL 被拉长 |
| RES-016-mfa-wrap | **R1 记载已过期**：现行代码已实现「previous 解密 + 成功 TOTP 后惰性重包」；仍存窄口径 = 恢复码路径不重包、无启动批量重包、无主动轮换重包 | `apps/api/modules/mfa/service.go:57`–`60`、`:72`、`:151`+`:165`、`:246`+`:259`、`:335`+`:353`、`:329`–`333`（恢复码分支提前返回）、`:450`/`:466`；接线 `composition.go:395`；来源注释 W11 F-004 | **接受残余** | 用户 2026-09-10 就**修正后窄口径**书面再裁决接受（见 §6）；旧表述的文档同步归 G-005 | RT-K03 行 | KMS/密钥波次；或出现「恢复码路径导致 MFA 不可解」的证据 |
| RES-021-harness | 进程级停机 harness 仍 `!windows`；compose stop 以 linux CI 核销 | R2 [矩阵行 11](../../GOAL-003-r2-as-built-matrix/attachments/as-built-matrix.md)；`internal/composition/shutdown_drain_test.go` | **接受残余** | 由 [VRev-047](../../../vision/reviews/VRev-047-vp021-closeout.md) `V-F083`（recommended · 低）在 VP-021 关门事务内登记的有界残余，范围与复审触发明确；闭合 = CI 证据采集或用户书面 accepted-residual | RT-D02 行保持 delivered | 该 CI 流程失败，或 Windows 出现停机回归 |
| RES-026-redis | 无 Redis 依赖，接缝与触发条件已冻结 | `apps/api/go.mod:5`–`22`；`docs/architecture/cache-redis-seam-and-track.md:62`–`79` | **仍 gated** | 本期不消耗 trigger | RT-Q03 / RT-Q05 保持 gated | 多实例部署或跨实例一致视图需求出现 |
| RES-028-broker | 无外部队列依赖；进程内 EventBus 为唯一运输 | `apps/api/go.mod:5`–`22`；`apps/api/kernel/eventbus.go:87` | **仍 gated** | 本期不消耗 trigger | RT-Q05、A3 | 跨进程可靠消息需求出现 |
| RES-030-keyfile | bot master key 文件与 DB 同目录 | [E-010](../../../workspace-030-telegram-channel-runtime/GOAL-001-telegram-channel-runtime/02-execution/E-010-a008-response-and-r5.md):21「R-009（默认密钥文件与 DB 同目录）→ accepted-residual（**用户书面**）」 | **接受残余** | 继承原 VP 用户书面接受（裁决 B），不重新裁决 | RT-M03 行保持 delivered | KMS 波次启动 |
| RES-T03-tz | 迁移与写入仍未使用 `timestamptz`（`apps/api` 内 `*.sql` 中 `timestamptz` 出现 0 次；时间列仍 INTEGER，如 `modules/authsession/migration/migration.go:41`） | 本轮全仓检索：`timestamptz` 命中 = 0；`docs/vision/roadmap.md:252` 仍 registered | **现在修** | 本 VP 红线禁止改 schema → 按裁决 A 另立；在 R4 路线图草案中登记为**下一拍候选**，由 `/vision` editorial 决定是否立项 | A 序列候选行（R4 草案建议） | R4 editorial 后由 `/vision` 决定立项或明确不做 |
| RES-P04-pool | 池化已交付（SQLite 文件库默认 4、内存库 1）；**读写分离 / replica 未实现**；文档现状锚点仍写 `MaxOpenConns=1` | `apps/api/internal/store/store.go:29`（`sqlitePoolDefault = 4`）、`:104`–`113`（内存 1 / 文件池）；文档锚点见 [r4-doc-hygiene-anchors.md](r4-doc-hygiene-anchors.md) G-002 | **明确不做** | 仅对「读写分离 / replica」而言；池化本身**已交付**。文档锚点过期另归 G-002（文档卫生） | RT-P04 保持 trigger-gated，但现状锚点须在 R4 修正 | 多实例或 PG 之后 |

## 2 · R2 缺口候选与本轮新增发现 G-001～G-006

| ID | 现状（本轮复核） | 证据（精确锚点） | 分类 | 去向 / 部署 | 影响的路线图行 | 复审触发 |
|----|------------------|------------------|------|-------------|----------------|----------|
| G-001 | `docs/architecture/overview.md` 的「当前阶段（现时）」过期：Charter 版本、产品树（`web/` FastAPI）、组合编排与工作区清单停在早期 | `docs/architecture/overview.md:64,71,72-73,74-79`（详见 [r4-doc-hygiene-anchors.md](r4-doc-hygiene-anchors.md) G-001） | **现在修** | 文档卫生，按裁决 A 落 R4（判据 4） | 架构文档权威面；不改任何架构结论 | R4 执行后由用户 editorial 确认 |
| G-002 | roadmap 现状锚点与 RT-P04/RT-D02 表述滞后或自相矛盾（`MaxOpenConns=1`；「无明确 drain 合同」与 delivered 并存） | `docs/vision/roadmap.md:104,138,209` | **现在修** | 文档卫生落 R4；`docs/vision/**` 部分须与草案同一事务或紧随 `/vision` 冻结 | RT-P04 现状锚点、RT-D02 行内表述 | 同上；不得把「已池化」扩写成读写分离已交付 |
| G-003 | `cache-redis-seam-and-track.md` §2.6 只描述旧三元组（`Allow`/`Record`/`Clear`），未覆盖现行 `AllowRecord`/`Reserve`/`Cancel` 原子面（该文档全文这三个名字 0 命中） | 文档 `:68,73`；端口 `apps/api/kernel/ratelimit.go:40,44,52,65,71,76,79` | **现在修** | 文档卫生落 R4（架构文档，不属 `docs/vision/**`） | 接缝文档须能指导触发后替换；**不冻结** Redis 实现细节 | 触发 RT-Q03/Q05 立项前必须已修正 |
| G-004 | `assembly.NewAuthenticator` 为公开工厂但返回 internal 指针且只接 current 单密钥（组合根注入 current+previous） | `apps/api/assembly/assembly.go:26,37,38`；组合根 `apps/api/internal/composition/composition.go:254`–`259`；单密钥形态 `apps/api/internal/auth/auth.go:151` | **明确不做** | 用户 2026-09-10 裁决：B+ 边界内，无包消费者失败证据；登记为 R4 草案的分发面候选 | 路线图 B+ 公共面行（R4 草案登记候选） | 出现第一个包消费者需要 previous 密钥，或 B+ 面成为分发契约 |
| G-005 | R1 记载的 `RES-016-mfa-wrap` 表述过期：多处公开文档仍写「`admin.mfa` wrapping 不随 JWT previous 重包」，与现行代码（W11 F-004 惰性重包）不符 | 文档：`docs/vision/charter.md:73`、`docs/vision/roadmap.md:16,223`、`docs/vision/workspaces.md:30,58`、`docs/vision/plans/VP-016-key-rotation-and-backup.md:115`；代码：`modules/mfa/service.go:57-60,72,151,165,246,259,335,353` | **现在修** | 文档卫生落 R4 + `/vision` editorial；VP-016 历史原文不回改 | RT-K03 行、VP-016 关门记录表述 | R4 执行后由用户 editorial 确认 |
| G-006 | R2 矩阵 v0.2.0 有 31 处锚点指向紧邻注释/空行/邻行（主张本身经 R2 A-002 独立复核属实），含 1 处**主题性错锚**：行 10 引 `internal/config/config.go:493`，该行实为 `CacheMaxEntries`，主张对应 `:495`/`:498` | 清单见 [r3-asbuilt-anchors.md](r3-asbuilt-anchors.md) §D；主题性错锚 `apps/api/internal/config/config.go:493` vs `:495`/`:498`（本轮亲自复核） | **现在修** | **已执行**（矩阵 v0.3.0，提交 `625e2945`）：31 处漂移 + 主题性错锚全部校正，每条校正本会话读回源码确认；主张/分类/verdict 未变 | 判据 1 的证据精度；不影响路线图行 | 已完成；R4 复核若发现主张实质变化则回到 A-00N 响应流程 |

## 3 · 需用户裁决的项（P-004）

| # | 项 | 裁决状态 |
|---|----|----------|
| 1 | **G-004** | ✅ 用户 2026-09-10 裁决 `明确不做`（B+ 边界内）+ 登记为 R4 草案分发面候选 |
| 2 | **RES-016-mfa-wrap** | ✅ 用户 2026-09-10 **修正前提后再裁决** `接受残余`（窄口径，见 §6） |
| 3 | **RES-013 / RES-014 搬运器** | ✅ 用户 2026-09-10 裁决维持 `明确不做`（非目标） |
| 4 | **RES-T03-tz** | ✅ 用户 2026-09-10 裁决登记为下一拍照候选（由 `/vision` 在 R4 editorial 决定立项） |

未列入裁决的项及依据：RES-030-keyfile（原 VP 用户书面 accepted-residual）、RES-015-otlp-sink（VP-015 关门事务内 recommended 文档化残余、开放 required = 0）、RES-021-harness（VRev-047 V-F083 登记的有界残余，范围与复审触发明确）——按裁决 B 继承，不重新裁决；RES-015-metrics、RES-026-redis、RES-028-broker 为**仍 gated**（不消耗 trigger）；RES-P04-pool、RES-013/014、G-004 为**明确不做**；G-001/G-002/G-003/G-005/G-006 为**现在修（文档，落 R4 或已执行）**。

⚠ **RES-016-revoke 更正**（v0.2.0）：v0.1.0 曾把它记为「接受残余（继承原留痕）」。复核原 VP 记录后确认 [Root I-005](../../../workspace-016-key-rotation-and-backup/GOAL-001-key-rotation-and-backup/00-meta.md) 仍为 `collecting`、并注明「用户书面残余时才改变退出 1」——**不存在合法书面接受**，故改记 `明确不做`（保留用户日后按裁决 B 接受为残余的可能）。

## 4 · 分类统计（18 个唯一条目）

| 分类 | 条数 | 条目 |
|------|------|------|
| 现在修 | 6 | RES-T03-tz（→ 另立为下一拍候选）；G-001、G-002、G-003、G-005（→ R4 文档卫生）；G-006（→ 已执行，矩阵 v0.3.0） |
| 仍 gated | 3 | RES-015-metrics、RES-026-redis、RES-028-broker |
| 接受残余 | 4 | RES-015-otlp-sink、RES-016-mfa-wrap、RES-021-harness、RES-030-keyfile |
| 明确不做 | 5 | RES-013-migrator、RES-014-migrator、RES-016-revoke、RES-P04-pool、G-004 |
| **合计** | **18** | 12（R1 入册）+ 4（R2 候选）+ 2（本轮新增 G-005/G-006） |

合计 6 + 3 + 4 + 5 = **18**，与唯一条目数一致。v0.1.0 的「19 条」为计数错误（把 18 个唯一 ID 误记为 19），v0.2.0 更正。

## 5 · 边界

- 本表不授权任何代码/端口/Profile 改动；分类为「现在修」的条目只在 R4 的文档卫生范围内执行（G-006 已完成，属只读断言/文档修正，未改任何主张或分类）。
- 本表不关闭任何 finding，不改 `docs/vision/roadmap.md`；`docs/vision/**` 的改动须与 R4 草案同一事务或紧随 `/vision` 冻结。
- 全部待裁决项已由用户于 2026-09-10 裁决（RES-016-mfa-wrap 经修正前提后于同日再裁决），无遗留裁决项。

## 6 · 中途发现与再裁决留痕（RES-016-mfa-wrap）

本轮只读取证子代理报告 `RES-016-mfa-wrap` 与 R1 记载**不一致**（一致 11 / 不一致 1 / 无法核对 0），编排器随后亲自复核源码确认：

| 项 | R1 记载 | 现行代码 | 结论 |
|----|---------|----------|------|
| 语义 | 「`admin.mfa` wrapping **不随** JWT previous 重包」 | previous 秘钥已参与解密，且成功 TOTP 后**惰性重包**（W11 F-004）：`modules/mfa/service.go:57-60,72,151,165,246,259,335,353`；接线 `composition.go:395` | R1 记载**过期** |
| 仍存的窄残余 | — | 仅凭恢复码完成的路径提前返回、不重包（`:329`–`333`）；无启动批量重包；无主动轮换重包 | 残余范围**显著小于** R1 记载 |

处置：

1. 用户原裁决（「接受残余」）所依据的**前提**与现行事实不符，故本表未直接沿用原措辞，而是提出**修正后的窄口径**重问；用户于 2026-09-10 **再次裁决 `接受残余`**。范围 = 上表「仍存的窄残余」，复审触发 = KMS/密钥波次，或出现「恢复码路径导致 MFA 不可解」的证据。
2. 旧表述仍留在多处公开文档（见 G-005），归 R4 文档卫生 + `/vision` editorial。
3. 残留原文（VP-016 关门时点）**不回改**：VP-016 正文与 VRev-036 保留其历史时点表述，仅由 R4 在现行文档面加注/修正。

## 7 · R2 矩阵锚点漂移（G-006）

只读取证复核发现 R2 矩阵 v0.2.0 尚有 31 处锚点指向紧邻注释/空行/邻行，另有 7 处「不算错行但可再精确」。**主张本身未受影响**（R2 A-002 已独立复核 17 面 + W1 的主张真实性），唯一主题性错锚为行 10 的 `internal/config/config.go:493`（实为 `CacheMaxEntries`；主张对应 `:495`/`:498`，本轮亲自复核）。

处置：**已执行** —— 矩阵升 v0.3.0 一次性校正（提交 `625e2945`），只改锚点文字，未改任何主张、分类或 verdict。完整清单见 [r3-asbuilt-anchors.md](r3-asbuilt-anchors.md) §D 与 [r3-anchor-commands.txt](r3-anchor-commands.txt)。
