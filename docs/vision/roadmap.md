---
doc_type: vision-roadmap
title: 总路线图 · 愿景组合编排
status: active
created: 2026-07-31
updated: 2026-09-22
parent: null
version: 0.98.1
---

# 总路线图 · Schema UI Core Admin 基架

本文件回答：现在处于哪里、已经交付什么、后续方向是什么、哪些事项仍有条件或残余。它索引 VP 与用户确认的方向；Goal 状态和执行证据以各工作区为准，不在这里汇总 progress%。

导航：[当前状态](#当前状态) · [VP 总览](#已落盘意图) · [编排规则](#三分支后续方向2026-08-20) · [架构](#架构分支) · [Admin](#admin-功能分支) · [业务域](#业务域分支) · [未决登记](#未决项统一登记残余--悬置--trigger-gated)

详细范围、版本、lead、审计编号、RT 能力清单与历史闭合记录移至 [roadmap-reference.md](roadmap-reference.md)。原有关键信息保留在主表或参考页，避免把执行历史反复复制成“当前状态”。

## 当前状态

| 项 | 现状 |
|----|------|
| 愿景 | [Charter](charter.md) `schema-ui-core-admin-foundation@0.4.0`，active；单主线、薄内核、模块与启动时 Profile；fork 与构建期包消费并存 |
| 协议基线 | `schema-ui-docs@v2.9.0`，pin `81aa1d8`；支持窗 v2.7–v2.9（additive 超集） |
| 交付 VP | **当前无 active 交付 VP**；VP-001～040 中除两项持续程序外均已 closed |
| 持续程序 | **VP-009 安全加固、VP-010 设计 / 实现符合性保持 active**；波次完成不等于程序关闭。波次现状见各区 [安全目标树](../workspaces/workspace-009-production-hardening/goal-tree.md)、[符合性目标树](../workspaces/workspace-010-design-implementation-conformance/goal-tree.md) |
| 最近交付 | 架构 VP-040：2026-09-21 closed v0.3.0，六条退出判据满足、VRev-105 pass；Admin VP-039：2026-09-19 closed v0.3.0；业务域 VP-031：2026-09-05 closed v0.3.4 |
| 下一步边界 | 候选与未决项按下文触发条件复核，再经 `/vision` 立项；当前登记不等于实施承诺，也不预开第二业务域 |

## 已落盘意图

按能力分组，**不表示新的执行顺序**。完整 40 项按编号排列的范围、依赖与关门依据见[VP 完整索引](roadmap-reference.md#vp-完整索引)；历史执行中 VP-006 先于 VP-005。

| 类别 | VP | 能力 | 状态 |
|------|----|------|------|
| 基线与准入 | [VP-001](plans/VP-001-mvp-admin-foundation.md) | Admin MVP 基线 | **closed** |
| 基线与准入 | [VP-002](plans/VP-002-production-admin-foundation.md) | 生产级 Admin 基架 | **closed** |
| 基线与准入 | [VP-003](plans/VP-003-modular-admin-architecture.md) | 模块化单体、薄内核与 Profile | **closed** |
| 基线与准入 | [VP-004](plans/VP-004-module-contribution-readiness.md) | 一方模块贡献 playbook | **closed** |
| 基线与准入 | [VP-005](plans/VP-005-design-system-and-ui-experience.md) | 设计系统与 UI 体验 | **closed** |
| 基线与准入 | [VP-006](plans/VP-006-full-protocol-contract-v2-7-0.md) | 完整 v2.7.0 协议兼容 | **closed** |
| 基线与准入 | [VP-007](plans/VP-007-localization-and-system-settings.md) | 多语种与系统设置 | **closed** |
| 基线与准入 | [VP-008](plans/VP-008-admin-module-readiness-and-foundation-convergence.md) | 全基架准入与 go | **closed** |
| 持续程序 | [VP-009](plans/VP-009-production-hardening.md) | 持续安全与健壮性程序 | **active**（长期程序） |
| 持续程序 | [VP-010](plans/VP-010-design-implementation-conformance.md) | 持续设计 / 实现符合性程序 | **active**（长期程序） |
| Admin 与共享契约 | [VP-011](plans/VP-011-admin-functional-modules.md) | 标准 Admin 功能模块 | **closed** |
| Admin 与共享契约 | [VP-012](plans/VP-012-shared-cross-module-contracts.md) | 共享横切契约首波 | **closed** |
| Admin 与共享契约 | [VP-018](plans/VP-018-account-email-identity.md) | 账号邮箱身份 | **closed** |
| Admin 与共享契约 | [VP-019](plans/VP-019-iam-recovery.md) | IAM 密码策略、邀请与自助恢复 | **closed** |
| Admin 与共享契约 | [VP-020](plans/VP-020-timezone-number-currency-formatting.md) | 时区、数字与货币格式 | **closed** |
| Admin 与共享契约 | [VP-025](plans/VP-025-config-export-diff-dryrun-import.md) | 配置导出 / diff / dry-run / 导入 | **closed** |
| Admin 与共享契约 | [VP-029](plans/VP-029-wallet-prepaid-instrument.md) | 钱包预付凭证与外部主体接缝 | **closed** |
| Admin 与共享契约 | [VP-033](plans/VP-033-telegram-operator-console.md) | Telegram 人工控制台 | **closed** |
| Admin 与共享契约 | [VP-034](plans/VP-034-nav-group-collapsible.md) | 导航分组折叠 | **closed** |
| Admin 与共享契约 | [VP-036](plans/VP-036-admin-command-palette.md) | 页面 / 导航 / 动作检索与 Command Palette | **closed** |
| Admin 与共享契约 | [VP-037](plans/VP-037-admin-workflow-continuity.md) | Saved Views、未保存保护与统一反馈 | **closed** |
| Admin 与共享契约 | [VP-038](plans/VP-038-batch-operations-and-job-center.md) | 批量操作与异步结果中心 | **closed** |
| Admin 与共享契约 | [VP-039](plans/VP-039-version-maintenance-diagnostics.md) | 版本、维护提示与诊断报告 | **closed** |
| 架构 | [VP-013](plans/VP-013-store-dialects.md) | Store 双方言（SQLite / PostgreSQL） | **closed** |
| 架构 | [VP-014](plans/VP-014-object-storage.md) | 对象存储（本地盘 / S3） | **closed** |
| 架构 | [VP-015](plans/VP-015-observability.md) | 指标与 OpenTelemetry | **closed** |
| 架构 | [VP-016](plans/VP-016-key-rotation-and-backup.md) | 密钥轮换与备份恢复 | **closed** |
| 架构 | [VP-017](plans/VP-017-outbound-mail.md) | 出站邮件与可切换渠道 | **closed** |
| 架构 | [VP-021](plans/VP-021-graceful-shutdown-and-connection-drain.md) | 优雅停机与连接排空 | **closed** |
| 架构 | [VP-026](plans/VP-026-cache-port.md) | 内存 Cache 端口与 Redis 接缝 | **closed** |
| 架构 | [VP-027](plans/VP-027-rate-limiter-port.md) | 内存限流端口与使用点迁移 | **closed** |
| 架构 | [VP-028](plans/VP-028-event-bus-port.md) | 进程内 EventBus 与 outbox/MQ 接缝 | **closed** |
| 架构 | [VP-030](plans/VP-030-telegram-channel-runtime.md) | Telegram 通道运行时 | **closed** |
| 架构 | [VP-032](plans/VP-032-rate-limiter-atomic-port.md) | 限流端口原子化 | **closed** |
| 架构 | [VP-035](plans/VP-035-foundation-architecture-health.md) | 基架架构健康评估 | **closed** |
| 架构 | [VP-040](plans/VP-040-timestamptz-persistence-contract.md) | 双方言绝对时刻持久化与 Backup SPI | **closed** |
| 组合层分发平台波 | [VP-022](plans/VP-022-distribution-package-pilot.md) | 构建期包消费试点 | **closed** |
| 组合层分发平台波 | [VP-023](plans/VP-023-productionization-cli-package.md) | CLI 与包消费产线化 | **closed** |
| 组合层分发平台波 | [VP-024](plans/VP-024-distribution-formalization.md) | 分发形态正式化 | **closed** |
| 业务域 | [VP-031](plans/VP-031-digital-offer-entitlement.md) | 数字 Offer、薄购买凭证与本域权益 | **closed** |

## 组合门闩（用户 2026-08-08）

1. **协议优先于视觉**：在 VP-006 未 `closed` 前，**不得**将 VP-005 作为 `primary_plan` 推进实现，不得启动视觉优化波次。  
2. **MVP 子集不是终态成功条件**：`I-PROTO-001 v0.1.3` 是历史 MVP 冻结切片；整份 v2.7.0 契约由 VP-006 收口。  
3. 已关闭 VP-001～004 的历史证据与 status **不重写**。

> **协议 pin 注记（2026-08-14 · VR-020 · editorial）**：协议覆盖权威由 `v2.7.0` 升至 `v2.8.0`（additive 超集；`I-PROTO-FULL-001` v1.0.1 仍为 v2.7.0 历史分母、已被 v2.8.0 覆盖）。身份权威见 `apps/web/src/protocol/upstream/provenance-v2.8.json`。  
> **协议 pin 注记（2026-08-29 · VR-050 · strategic）**：协议来源升至 **`v2.9.0`**（pinned `81aa1d8`；支持窗 2.7–2.9 additive 超集——代码实现追认，I-007 闭合；Charter `@0.3.0`）。

上述“协议先于视觉”门闩已随 VP-006 / VP-005 关门满足，保留作为历史顺序约束。**VP-008 的 go 不是永久凭证**：后继 VP 激活前按其消费有效性规则核对拟消费候选与 scope，完成 freshness review；共享基架缺陷可暂挂 go，须重验证后恢复。

## 三分支后续方向（2026-08-20）

用户确认组合层后续方向拆成三条可并行轨道：**架构**、**Admin 功能**、**业务域**。  
2026-08-20 上午曾登记为「架构 / 产品」二分（VR-025）；同日改为三分（VR-026）。仍共用 VP-003 单主线，**不是**已退役的 [dual-track-contract.md](dual-track-contract.md) 两套代码线。

| 分支 | 承接原四档 | 管什么 | 不管什么 |
|------|------------|--------|----------|
| **架构分支** | 运行时平台（四档未覆盖的缺口） | fork 部署依赖的存储方言、缓存、队列、对象存储、可观测、多实例、密钥与备份 | 不新建 Admin 页面或领域模块；不重开 VP-012；不替代 VP-009/VP-010 |
| **Admin 功能分支** | Tier A 应用能力剩余 + Tier B 扩展接缝 + Tier C 体验增强 | 通用 Admin 能力、扩展接缝、体验增强 | 不引入第二套持久化栈；不把订单/库存/CMS 等塞进 Admin VP |
| **业务域分支** | 原 Tier D | 成立后的真实业务领域（Catalog、订单、支付、库存、CMS…） | 不私建 Redis/MQ/对象存储；不把 IAM/SSO/搜索当业务域 |

**并行规则**

1. 三分支可同时各有一个 active 交付 VP；分支之间不设一概而论的硬前置，具体能力依赖仍须满足（如 VP-029 → VP-030/031）。
2. Store / 迁移方言 / 连接池 / 进程模型 / 外部中间件 → **架构**。Admin 与业务域只消费稳定内核接口。
3. 领域事件 / 通知 / SSO 的**应用契约** → **Admin 功能**；outbox、broker、SMTP、对象存储等**运输实现** → **架构**。禁止两套平行队列。
4. 订单、支付、库存等实体与流程 → **业务域**。可消费 Admin 接缝（审批、通知、事件），但不得在业务 VP 里重做基架。
5. 共享基架安全缺陷 → VP-009；设计/实现 gap → VP-010。二者是持续程序，不是第四分支。
6. 登记 ≠ 立项。具体 VP 仍须 `/vision` 冻结退出分母后再交 `/govern` 开区。
7. 横切契约增量**不默认重开 VP-012**。

## 架构分支

当前基线：单 API 进程；SQLite 文件库（默认池 4，内存库 1）/ PostgreSQL；本地盘 / S3；进程内 Job 六态、Cache、原子限流与 EventBus；可选 Prometheus / OTLP（默认关闭）；JWT current/previous；邮件 mock / Resend（SMTP 保留）。Compose 不承诺 TLS 终止与多实例。

| 序列 | 能力 | 状态 / 承接 |
|------|------|-------------|
| A0 | 架构能力登记、Store 双方言决策 | 已完成；RT-P03 冻结 |
| A1 | 内核持久化端口、PG 与 SQLite、checksum 迁移及方言备份 | delivered · VP-013 |
| A2 | 对象存储端口、本地盘与 S3 | delivered · VP-014 |
| A3 | 多实例、依赖就绪扩展、锁 / 队列方案评估 | **trigger-gated**；A0–A7 序列中唯一未交付的门控项 |
| A4 | 指标与 OpenTelemetry | delivered · VP-015 |
| A5 | 密钥轮换与轮换后恢复 | delivered · VP-016 |
| A6 | 邮件端口、mock / Resend、配置热切换与试发 | delivered · VP-017；SMS 后置 |
| A7 | 单进程优雅停机与连接排空 | delivered · VP-021；已从 A3 拆出 |
| C1 / RT-T03 | PG `timestamptz(6)`、SQLite fixed-6 UTC RFC3339 `TEXT`；统一微秒 wire 与 Backup SPI/Service | delivered · VP-040；90 列 / 44 表，双方言原地迁移 |

**H-002 触发评估**：C 端业务域接入同进程时，必须在开区前评估 RT-Q03 缓存与 RT-Q05 限流，允许有证据地得出“不需要 Redis”。VP-030/031 已作该评估；多实例或跨实例共享失效需求触发复核。VP-026/027/028/032 已交付内存端口与接缝，**不等于 Redis / broker / outbox 已实现**；EventBus 也不等于 Job 或 typed domain event 应用契约。

**冻结边界**：PG 是生产验收权威，SQLite 保留内嵌默认且合同平等；S3 与本地盘同理。无 ORM、无第三数据库；模块拥有 Repository / Persistence，驱动与事务类型不进入公共契约。MongoDB、ORM、微服务 / 运行时插件、默认 CQRS / 事件溯源、K8s / 服务网格不属于默认成功条件。

完整 [RT 清单与状态判定](roadmap-reference.md#架构分支)覆盖数据生命周期、缓存 / 锁 / 队列、文件、可观测、部署、密钥、搜索、时间和消息；[冻结合同](roadmap-reference.md#架构已冻结合同)保留每波具体边界与非目标。签名 URL、分片、扫描执行器、CDN、剖析、Sentry、KMS/HSM、PITR、TLS 等仍按各自触发条件处理，不能由“A3 唯一”推断为已经交付。

## Admin 功能分支

承接 Tier A 剩余、Tier B 接缝与 Tier C 体验。VP-011 标准模块与 VP-012 横切契约首波已 closed；保留 / 归档与 session envelope 已在 workspace-012 增量完成。effective actor 仍为当前 actor；impersonation 须重新出现产品需求。

### 已交付

| 能力 | 承接与边界 |
|------|------------|
| 邮箱身份 → 密码策略 / 邀请 / 自助恢复 | VP-018 → VP-019；运输依赖 VP-017。自助恢复证明为持有已校验邮箱；管理员重置仍走 `must_change_password` 特权路径；邀请可由管理员出示链接；SMS 后置 |
| 时区、数字与货币格式 | VP-020；DB 持久化合同另归 VP-040 |
| 配置导出 / diff / dry-run / 导入 | VP-025；密钥 fail-closed，热加载不在该波范围 |
| 钱包预付凭证与外部主体 | VP-029；不等于支付 / 结算业务域 |
| Telegram 人工控制台 | VP-033 消费 VP-030 runtime；短轮询 + heartbeat，不解除 SSE/WebSocket 门控 |
| 导航与体验增强 | VP-034 导航分组；VP-036 Command Palette；VP-037 Saved Views / 未保存保护 / 统一反馈；VP-038 批量作业结果中心；VP-039 版本 / 维护 / 诊断 |

**体验增强边界**：上述波次均已关门；VP-036 检索的是已注册页面、导航与声明式动作，**实体全文检索仍未实现**。VP-038 新增 `admin.jobs` 进入 admin 默认集是用户确认的 Profile 内容扩展，不改变装配语义、不暂挂 go；不改同步 `batch-delete` 语义。VP-039 不新建模块、不改默认集。

### 仍开放（非立即实施）

- **基架能力剩余**：原 #2 组织 / 部门 / 岗位与 `org` 数据权限（2026-08-29 用户降权为 trigger-gated）；原 #4 文件扫描 / 隔离策略（执行器归 RT-S05）；原 #6 impersonation。原 #1、#3、#5、#7 已交付，保留编号映射于[Admin 参考](roadmap-reference.md#admin-能力与交付依据)。
- **扩展接缝**（全部 trigger-gated）：typed domain event、Notification Transport、OIDC/SSO/SCIM、Approval Gate、通用 Entitlement、多组织 context、SSE/WebSocket、外部连接器 / Secret 产品面、自定义 metadata/tags、文件预览。VP-031 的本域权益不解除通用 Entitlement / Approval Gate 门控；运输实现归架构。
- **体验增强剩余**：实体全文检索；需要专用引擎时拉动 RT-X01，DB 全文检索见 RT-X02。协作视图等范围决策见下方 C 类登记。

## 业务域分支

承接原四档 **Tier D**。Charter 非目标仍然成立：本项目不把特定业务领域的终端产品写成愿景成功条件；领域只在业务成立后独立立项。

| 约束 | 说明 |
|------|------|
| 激活前 | 须 `/vision` 复核；消费 VP-008 `go`（候选 `ed99e88`）并对拟消费候选做 freshness review |
| 默认承载 | VP-003 架构 + VP-004 playbook + VP-006 协议面 + VP-005 设计系统 + VP-007 locale/settings |
| 问题分流 | 领域问题留在该业务 VP 台账；共享基架安全 → VP-009；平台缺口 → 架构分支；通用 Admin 缺口 → Admin 功能分支 |
| 禁止 | 用业务模块倒逼恢复长期双线、跳过协议覆盖、私增协议语义、私建 Redis/MQ/第二数据库 |

### 候选域（成立后再立项；一域一 VP）

1. Catalog、商品、SKU、价格、税  
2. 库存、仓库、预留、调拨  
3. 订单、支付、退款、退货  
4. 物流、包裹、履约  
5. 营销、优惠券、促销规则  
6. 订阅、计费、发票、用量  
7. 工单、客服、CRM  
8. CMS、内容发布、知识库  
9. 支付网关、ERP、物流、CRM 连接器（领域侧；通用连接器接缝在 Admin 功能）

已有的钱包模块是 VP-011 交付的 Admin 演示/能力面，**不**等于本分支「支付/结算」业务域已成立。VP-029 的卡密入金仍属 Admin 功能（资金通道），不把本表第 3 项提前成立。

**已立项（2026-09-02 · 真实触发 = 下游 Telegram 付费服务 · 用户确认同进程）**

| VP | 收窄后的域 | 明确不做 |
|----|------------|----------|
| [VP-031-digital-offer-entitlement](plans/VP-031-digital-offer-entitlement.md) `closed` v0.3.4 | 数字 Offer + 薄购买凭证 + 本域权益 | 类目树、SKU/税/库存、物流订单、支付网关、通用 Entitlement 框架 |

业务域下一拍：无新触发则不要预开第二域。不要把候选 1～9 打进同一个 VP。

## 未决项统一登记（残余 / 悬置 / trigger-gated）

本节登记已交付范围之外的残余、悬置与触发条件，不是实施承诺。新增或闭合事项须同步此处并保留原工作区决策 / 审计证据；不得把 deferred / recommended 写成 verified。**已 fixed 或已交付的条目移至[闭合留档](roadmap-reference.md#已闭合事项留档)**，不再混入开放清单。

### 一、有界残余与证据补齐（B 类）

| 编号 | 内容 | 现状 | 触发条件 | 责任人 | 证据 |
|------|------|------|----------|--------|------|
| [GOAL-008 A-002 F-002](../workspaces/workspace-037-admin-workflow-continuity/GOAL-008-typecheck-evidence-convention/03-audit/) | 全仓 `tsc` 简写未逐条裁定（354 行形态不可唯一确定） | **bounded residual**：可执行面由守卫 + CI 门禁锁死；文档侧不逐条考古 | **历史记录被再次当作类型检查证据引用时**，按 `GOAL-008 D-001` 口径复核并注明 | `/govern`（引用时） | [GOAL-043 E-003](../workspaces/workspace-010-design-implementation-conformance/GOAL-043-w31-cross-workspace-residual-closeout/02-execution/E-003-tsc-evidence-residue-closeout.md)；[GOAL-008 D-001](../workspaces/workspace-037-admin-workflow-continuity/GOAL-008-typecheck-evidence-convention/01-decision/D-001-typecheck-convention-and-scope.md) |
| 本地扩展清单 | 本地扩展登记（`valueLabels`/`badgeStyleField`/`truncate`/`width`/`minWidth` 等列级与节点级本地扩展缺少**单一清单**） | **登记（部分补齐）**：[GOAL-045](../workspaces/workspace-010-design-implementation-conformance/GOAL-045-w33-list-actions-slot-and-roles-trigger/00-meta.md) 已把「列表页 actions 左侧插槽」加入本清单并记录其取值与约定；其余扩展仍分散在各波次决策里 | 出现「这是 pinned 还是本地扩展」的实际争议，或后续协议波次需要一次性核对时 | workspace-010（后续符合性波次） | [GOAL-044 A-001 F-003](../workspaces/workspace-010-design-implementation-conformance/GOAL-044-w32-r4-residual-seams/03-audit/)；[GOAL-045 D-001 §1/§4 与 A-002](../workspaces/workspace-010-design-implementation-conformance/GOAL-045-w33-list-actions-slot-and-roles-trigger/00-meta.md) |
| [GOAL-006 A-001/A-002 F-001](../workspaces/workspace-038-batch-operations-and-job-center/GOAL-006-r5-evidence-and-closeout/03-audit/) | e2e 挂具的 **fresh-seed 顺序契约**仍依赖文件名排序（`00-` 前缀）：全量套件共用一块 scratch 库，任何更靠前的新 spec 若使用 sign-in helper（其会自动完成强制改密）都会再次消费该前提，失败形态是难与真实缺陷区分的「登录 401」 | **bounded residual（2026-09-19 登记）**：本次已由重命名修复并两向复验（隔离通过 / 全量 16 passed / 4 skipped / 0 failed），但机制仍是隐式约定 | 后续波次新增 e2e 文件，或排序/挂具策略变更时 | workspace-010（后续符合性波次；e2e 挂具属 W23/W24/W25 一系） | [GOAL-006 E-001 §2、A-001/A-002 F-001](../workspaces/workspace-038-batch-operations-and-job-center/GOAL-006-r5-evidence-and-closeout/00-meta.md)；`apps/web/e2e/00-force-password-change.spec.ts` 注释 |

本地扩展的已知字段、层次、取值与回退语义见[本地扩展清单](roadmap-reference.md#本地扩展清单)。jobs 结果中心浏览器主路径已补齐；重试 / 取消的浏览器端到端覆盖随实现扩展时补，原闭合记录仍保留在参考页。

**早期波次的有界范围与残余索引**（保留各波次原定性，不能据此推断已验证或用户已接受）：

| 来源 | 保留的范围 / 残余 | 后续依据 |
|------|--------------------|----------|
| VP-013 / VP-014 | 无产品级 SQLite→PG 搬运器、无本地盘→对象存储搬运器 | [VP-013 GOAL-006 D-002](../workspaces/workspace-013-store-dialects/GOAL-006-r5-dual-path-acceptance/01-decision/) / [VP-014 I-014-004](plans/VP-014-object-storage.md)；见完整 VP 与冻结合同 |
| VP-015 | in-repo OTLP sink 不解析；Store / 对象 / Job 指标未进首波分母 | I-015-003 / VP-015；后续范围另立项 |
| VP-016 | 立即失效未选（I-016-005，**无书面 residual 接受**）；MFA 已支持 previous 解密与成功 TOTP 后惰性重包，剩恢复码路径 / 启动批量 / 主动轮换不重包 | W11 F-004；MFA 收窄残余于 2026-09-10 获用户接受 |
| VP-021 | 进程级 harness / compose stop 由 Linux CI 核销 | 原波次证据，不宣称本地已核销 |
| VP-024 | hosted CI 实触发、shell 类型面、GH Packages 保留、C 类 fork 包化候选 | [GOAL-008 D-001 四项登记](../workspaces/workspace-024-distribution-formalization/GOAL-008-r7-topline-and-closeout/01-decision/D-001-residuals-and-topline.md) |
| VP-030 / VP-034 | R-009 bounded residual；Dashboard 现行 `workspace` 组残余 | [VP-030](plans/VP-030-telegram-channel-runtime.md)；[VP-034 GOAL-003](../workspaces/workspace-034-nav-group-collapsible/GOAL-003-sidebar-engine-navigation/00-meta.md) |

早期条目的完整限定与历史版本见[VP 完整索引](roadmap-reference.md#vp-完整索引)，不在本次编辑中新增接受、关闭或实施决定。

### 二、悬置的范围决策（C 类：非缺陷，等需求再定）

| 编号 | 内容 | 现状 | 触发条件 | 责任人 | 证据 |
|------|------|------|----------|--------|------|
| `I-037-005` | 跨用户共享视图、最近使用/收藏、协作权限是否进入后续波次 | 首波明确只做**个人级**工作流；`deferred · non-blocking` | **真实协作需求出现** | `/vision` | [VP-037 首波范围](plans/VP-037-admin-workflow-continuity.md)；[Root 信息表](../workspaces/workspace-037-admin-workflow-continuity/GOAL-001-admin-workflow-continuity/00-meta.md) |
| `I-038-006` | 历史作业保留 / 清理、归档与容量上限 | `deferred · non-blocking`；首波只做可见性与结果读取，未新建数据生命周期程序 | 作业表容量、保留期或合规出现真实需求 | `/vision` | [VP-038 信息表](plans/VP-038-batch-operations-and-job-center.md)；已有延期，非新增实施承诺 |

### 三、未推进 / trigger-gated 能力（A 类）

| 能力 | 现状 | 触发条件 | 责任人 / 下一步 | 出处 |
|------|------|----------|-----------------|------|
| 实体全文检索（`RT-X01` 专用引擎 / `RT-X02` DB 全文检索） | 架构触发项，**未实现**；VP-036 首波只做已注册页面/导航/声明式动作检索；VP-038 只做作业/批量结果的可见性，**不含**实体检索 | 真实**实体级**搜索需求 + 规模证据 | `/vision` 立新 VP；引擎路线拉动架构 `RT-X01` | roadmap「体验增强」；VP-036 边界表；VP-038 边界表 |
| 组织·部门·岗位 + 数据权限 `org` | 基架能力剩余 #2；**2026-08-29 用户书面降权**为 trigger-gated | 多组织/多团队 fork 消费，或真实多组织管理需求 | `/vision`（应用层 org 上下文归 Admin 分支） | roadmap「基架能力剩余」#2 |
| 新业务域 | **未立项**；无新触发不预开第二域 | 真实业务需求 | `/vision` | roadmap「业务域」 |
| 多实例及 Redis / MQ 等外部基础设施 | 架构 **A3** / `RT-Q03`，**A0–A7 序列唯一未交付的门控项**；第二持久化栈仍受 Charter / 双方言合同限制 | 多实例或 C 端域接入触发评估；VP-030/031 已评估本波不需要 Redis，后续按 RT-Q03/Q05 复核 | 架构分支 + `/vision` | roadmap「基架能力剩余」；[Redis 接缝与轨道](../architecture/cache-redis-seam-and-track.md) |
| 文件扫描 / 隔离**策略**（执行器见 `RT-S05`） | 基架能力剩余 #4，**未立项** | 真实需求 | `/vision` | roadmap「基架能力剩余」#4 |
| 扩展接缝：typed domain event、Notification Transport、OIDC/SSO/SCIM、Approval Gate、Entitlement、多组织 context、SSE/WebSocket、外部连接器/Secret 产品面、自定义 metadata/tags、文件预览 | 全部 **trigger-gated**（SSE 注记：VP-033 用短轮询，**不**解除本行；Entitlement 注记：VP-031 只交付数字 Offer 本域，**不**解除本行） | 各自真实需求 | `/vision` | roadmap「扩展接缝」 |

此外，impersonation / effective actor 产品化仅在需求重新出现时复核；架构各 RT 项的独立触发条件以[完整能力清单](roadmap-reference.md#架构分支)为准。

## 单主线模块化策略

未来 fork 起点统一由同一代码主线、模块候选集与启动时 Profile 表达，权威见 [module-architecture.md](../architecture/module-architecture.md) 和 VP-003。原 [dual-track-contract.md](dual-track-contract.md) 已转为历史记录。
