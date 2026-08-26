---
doc_type: vision-roadmap
title: 愿景组合编排
status: active
created: 2026-07-31
updated: 2026-08-26
parent: null
version: 0.44.0
---

# 组合编排 · Schema UI Core Admin 基架

本文件索引已落盘的 VP 与用户确认的后续方向；它不是 Goal 路线图，也不汇总 progress%。

## 已落盘意图

| 顺序 | VP | 意图 | 前置 | 状态 |
|------|----|------|------|------|
| 1 | [VP-001-mvp-admin-foundation](plans/VP-001-mvp-admin-foundation.md) | 初始化 React + Go Admin MVP，覆盖固定协议来源、核心账号权限与协议范例验证。 | 无 | **closed**（2026-08-01；lead: workspace-001-mvp-admin-foundation） |
| 2 | [VP-002-production-admin-foundation](plans/VP-002-production-admin-foundation.md) | 在 I-PROTO-001 冻结子集之上，交付可直接 fork 使用的生产级 Schema 驱动 Admin 基架。 | 继承 VP-001 协议验证基线 | **closed**（2026-08-04；lead: workspace-002-production-admin-foundation） |
| 3 | [VP-003-modular-admin-architecture](plans/VP-003-modular-admin-architecture.md) | 单主线模块化单体：薄内核、模块契约、Fx、Profile、后端聚合 Manifest。 | 继承 VP-002；strategic re-align 见 VRev-006 | **closed**（2026-08-06；lead: workspace-003-modular-admin-architecture） |
| 4 | [VP-004-module-contribution-readiness](plans/VP-004-module-contribution-readiness.md) | 一方模块贡献 playbook 与 Core vs 模块归属方法论。 | 继承 VP-003 | **closed**（2026-08-06；lead: workspace-004-module-contribution-readiness） |
| 5 | [VP-006-full-protocol-contract-v2-7-0](plans/VP-006-full-protocol-contract-v2-7-0.md) | **`schema-ui-docs@v2.7.0` 整份契约**可验证兼容：覆盖表升版、Renderer/后端实现、范例与验证；纠正「长期停留在 MVP 子集」的组合焦点。 | 继承 VP-003/004；以 inventory + 上游 pin 为权威；`I-PROTO-001 v0.1.3` 仅作升版起点 | **closed**（2026-08-08 用户书面确认；lead: workspace-005-full-protocol-contract-v2-7-0；`I-PROTO-FULL-001` v1.0.1 = 12/12 include、318 executed + 2 local adapter excluded） |
| 6 | [VP-005-design-system-and-ui-experience](plans/VP-005-design-system-and-ui-experience.md) | Design Token、shadcn/ui 风格、Renderer/Shell 视觉与状态工效产品化。 | 继承 VP-003/004 + **VP-006 已 closed 的整份协议面**；VRev-011 `F-V018`/`F-V019`/`F-V020` → **fixed**（v0.3.0） | **closed**（2026-08-09 用户书面确认；v0.5.0；lead: `workspace-006-design-system-and-ui-experience`；Root `GOAL-001-design-system-and-ui-experience` `done 5/5`） |
| 7 | [VP-007-localization-and-system-settings](plans/VP-007-localization-and-system-settings.md) | 建立 `zh-CN` / `en-US` 多语种运行时与 `auto` 解析，并把既有 Settings 产品化为 General / Branding / Localization / Appearance 四类系统设置。 | 继承 VP-003/004 模块边界、VP-005 设计系统与 VP-006 完整协议面；不改变双 Profile 的 Settings 可见性边界 | **closed**（2026-08-09 用户书面确认；lead: `workspace-007-localization-and-system-settings`，Root done 6/6） |
| 8 | [VP-008-admin-module-readiness-and-foundation-convergence](plans/VP-008-admin-module-readiness-and-foundation-convergence.md) | 在正式业务模块开发前，对当前代码主线执行全基架准入：现状扫描、代码/功能/治理缺口、UI 协议判断、阻断整改与 `go`/`no-go`。 | 继承 VP-003/004 模块架构与贡献契约、VP-005 设计系统、VP-006 完整协议面、VP-007 locale/settings；不重开历史 VP | **closed**（2026-08-10 用户书面确认；候选 `ed99e88` clean，S0–S5 完成、open required = 0、`go` 签发；lead: workspace-008-admin-module-readiness，Root `GOAL-001-admin-module-readiness` done 6/6） |
| 9 | [VP-009-production-hardening](plans/VP-009-production-hardening.md) | 生产加固：**共享基架持续安全与健壮性程序**（周期扫描、波次修复、与 VP-008 `go` 消费有效性接口）；具体 finding 由工作区波次子目标承接。 | 继承 VP-003/004/005/006/007 + **VP-008 `go` 消费有效性**；共享基架缺陷可暂挂/恢复 `go` | **active**（2026-08-10 语义纠正为长期程序；曾误 `closed` 已撤销；lead: workspace-009-production-hardening；Root **active** 程序容器；波次 W1–W4 与 W6–W12 均 done，W5 扫描 0 中高危未开子目标；W12 = 2026-08-26 跨区限流评估收官） |
| 10 | [VP-010-design-implementation-conformance](plans/VP-010-design-implementation-conformance.md) | 设计意图与实现符合性：**共享基架持续对齐程序**（周期对照 as-designed / as-built、conformance gap 分流、波次整改、与 VP-008 `go` 消费有效性接口、与 VP-009 正交）；具体 gap 由工作区波次子目标承接。 | 继承 VP-003/004/005/006/007/008 + **VP-008 `go` 消费有效性**；与 **VP-009** 正交（安全 vs 符合性） | **active**（2026-08-11 用户确认类 VP-009 长期程序；lead: workspace-010-design-implementation-conformance；Root **active** 程序容器；波次 W1–W13 均 done，`go` 均无新暂挂） |
| 11 | [VP-011-admin-functional-modules](plans/VP-011-admin-functional-modules.md) | 标准 Admin 功能模块（通用模块 + 常用业务领域）分档交付：有界调研 → 三档分档 → 分波实现；一等公民 / 常用 / 增补。 | 继承 VP-008 `go` 消费有效性（freshness review **PASS**，候选 `f14ab9d`）+ VP-009/010 无开放阻断；VP-001～008 已固化协议/架构/设计/locale 基线 | **closed**（2026-08-18 有界关门；lead: workspace-011-admin-functional-modules；Root done；四档能力地图上提至本 roadmap） |
| 12 | [VP-012-shared-cross-module-contracts](plans/VP-012-shared-cross-module-contracts.md) | 共享横切契约与平台基架：correlation、审计模型、并发/幂等、异步 Job、maintenance 门控、API Token；不承载业务领域。 | 继承 VP-011 的 R5 四档能力地图；与 VP-009/VP-010 正交分流；不改变 Charter 边界 | **closed**（2026-08-19 完整关门 · 首波；lead: workspace-012-shared-cross-module-contracts；Root done 6/6；后续 session/effective actor、保留/归档、其余 writer envelope 移交本文件 Admin 功能分支） |
| 13 | [VP-013-store-dialects](plans/VP-013-store-dialects.md) | 架构 A1：内核持久化端口 + PostgreSQL 实现 + 现有迁移台账对写；SQLite 保留为内嵌默认；无 ORM。 | RT-P03 已冻结（VR-027）；继承 VP-003 模块化内核与全局台账；与 VP-009/010 正交；不进 A2+ 与 Admin/业务域 | **closed**（2026-08-21 有界关门 · 架构 A1；lead: workspace-013-store-dialects；Root done 5/5；residual：无产品 SQLite→PG 搬运器，见 GOAL-006 D-002） |
| 14 | [VP-014-object-storage](plans/VP-014-object-storage.md) | 架构 A2：内核对象存储端口 + S3 兼容实现；本地盘保留为内嵌默认。 | VP-013 A1 已 closed；RT-S01 delivered；与 VP-009/010 正交；不进签名 URL / 分片 / 扫描 / CDN / 搬运器，不进 A3+ 与 Admin/业务域 | **closed**（2026-08-21 有界关门 · 架构 A2；lead: workspace-014-object-storage；Root done 5/5；VRev-032 `pass`；residual：无产品本地盘→对象存储搬运器，见 I-014-004） |
| 15 | [VP-015-observability](plans/VP-015-observability.md) | 架构 A4：Prometheus 类指标导出 + OpenTelemetry traces；无收集器仍为内嵌默认。 | VP-014 A2 已 closed；RT-O01/O02 delivered；与 VP-009/010 正交；不进 A3 / A5 / Sentry / 剖析 / Admin 页 / 业务域 | **closed**（2026-08-22 有界关门 · 架构 A4；lead: workspace-015-observability；Root done 5/5；VRev-034 `pass`；residual：otlp-sink 不解析 + Store/对象/Job 指标不进分母） |
| 16 | [VP-016-key-rotation-and-backup](plans/VP-016-key-rotation-and-backup.md) | 架构 A5：JWT current+previous 轮换合同 + 既有备份上的轮换后恢复；单密钥仍为内嵌默认。 | VP-015 A4 已 closed；RT-K01 delivered；VP-013 方言级 dump 已交付；与 VP-009/010 正交；不进 A3 / KMS / PITR / 热加载 / Admin 页 / 业务域 | **closed**（2026-08-22 有界关门 · 架构 A5；lead: workspace-016-key-rotation-and-backup；Root done 5/5；VRev-036 `pass`；residual：I-016-005 立即失效未选 + `admin.mfa` wrapping 不随 JWT previous 重包） |
| 17 | [VP-017-outbound-mail](plans/VP-017-outbound-mail.md) | 架构 A6 升级：内核发送端口 + 可切换渠道（默认 mock 站内出站记录 + 生产 Resend；SMTP 适配器保留不删）+ 管理设置/试发。 | 用户 2026-08-24 否决同日 SMTP 专用有界关门（实施史不回退）。不进账号 email / 邀请 / 恢复状态机 / 模板 / 用户站内通知 / SMS / A3 | **closed**（2026-08-24 按**现行渠道分母**再关门 · v0.5.0；lead: workspace-017-outbound-mail；Root `done` 8/8；R5～R8 = GOAL-006～009 done；live 投递实跑 PASS；VRev-042） |
| 18 | [VP-018-account-email-identity](plans/VP-018-account-email-identity.md) | Admin 功能：账号邮箱身份（`users` email 可空 + 绑定/校验状态机 + 换绑）；消费 VP-017 `MailSender`。 | 硬前置 = VP-017 **再次** `closed`（现行渠道分母）。不进 IAM 恢复 / 邀请 / 密码策略 / SMS / 模板 / A3 | **closed**（2026-08-24 解冻当日关门 · v1.0.0；lead: workspace-018-account-email-identity；Root `done` 4/4；VRev-040 pass；A-002 independent 归零） |
| 19 | [VP-019-iam-recovery](plans/VP-019-iam-recovery.md) | Admin 功能 · IAM：密码策略 / 邀请入职 / 自助恢复状态机（忘密全链消费 VP-018 已校验邮箱 + VP-017 `MailSender`）。 | 硬前置 = VP-018 邮箱身份 + VP-017 运输（均已 `closed`）。不进 SMS / 模板中心 / 多邮箱 / 组织权限 / OIDC / 业务域 / A3；不改 Profile 默认集 | **active**（2026-08-25 激活；VRev-043 independent `pass`，required = 0；Admin 类 freshness PASS `092bf37` → `66f5fd1f`；lead: workspace-019-iam-recovery） |

## 组合门闩（用户 2026-08-08）

1. **协议优先于视觉**：在 VP-006 未 `closed` 前，**不得**将 VP-005 作为 `primary_plan` 推进实现，不得启动视觉优化波次。  
2. **MVP 子集不是终态成功条件**：`I-PROTO-001 v0.1.3` 是历史 MVP 冻结切片；整份 v2.7.0 契约由 VP-006 收口。  
3. 已关闭 VP-001～004 的历史证据与 status **不重写**。

> **协议 pin 注记（2026-08-14 · VR-020 · editorial）**：协议覆盖权威由 `v2.7.0` 升至 `v2.8.0`（additive 超集；`I-PROTO-FULL-001` v1.0.1 仍为 v2.7.0 历史分母、已被 v2.8.0 覆盖）。身份权威见 `apps/web/src/protocol/upstream/provenance-v2.8.json`。

---

## 三分支后续方向（2026-08-20）

用户确认组合层后续方向拆成三条可并行轨道：**架构**、**Admin 功能**、**业务域**。  
2026-08-20 上午曾登记为「架构 / 产品」二分（VR-025）；同日改为三分（VR-026）。仍共用 VP-003 单主线，**不是**已退役的 [dual-track-contract.md](dual-track-contract.md) 两套代码线。

| 分支 | 承接原四档 | 管什么 | 不管什么 |
|------|------------|--------|----------|
| **架构分支** | 运行时平台（四档未覆盖的缺口） | fork 部署依赖的存储方言、缓存、队列、对象存储、可观测、多实例、密钥与备份 | 不新建 Admin 页面或领域模块；不重开 VP-012；不替代 VP-009/VP-010 |
| **Admin 功能分支** | 原 Tier A 剩余 + Tier B + Tier C | 通用 Admin 能力、扩展接缝、体验增强 | 不引入第二套持久化栈；不把订单/库存/CMS 等塞进 Admin VP |
| **业务域分支** | 原 Tier D | 成立后的真实业务领域（Catalog、订单、支付、库存、CMS…） | 不私建 Redis/MQ/对象存储；不把 IAM/SSO/搜索当业务域 |

**并行规则**

1. 三分支可同时各有一个 active 交付 VP，互不作为对方的硬前置。
2. Store / 迁移方言 / 连接池 / 进程模型 / 外部中间件 → **架构**。Admin 与业务域只消费稳定内核接口。
3. 领域事件 / 通知 / SSO 的**应用契约** → **Admin 功能**；outbox、broker、SMTP、对象存储等**运输实现** → **架构**。禁止两套平行队列。
4. 订单、支付、库存等实体与流程 → **业务域**。可消费 Admin 接缝（审批、通知、事件），但不得在业务 VP 里重做基架。
5. 共享基架安全缺陷 → VP-009；设计/实现 gap → VP-010。二者是持续程序，不是第四分支。
6. 登记 ≠ 立项。具体 VP 仍须 `/vision` 冻结退出分母后再交 `/govern` 开区。
7. 横切契约增量**不默认重开 VP-012**。

**原四档 → 三分支（只读映射）**

| 原档 | 现归属 |
|------|--------|
| Tier A 应用能力剩余 | Admin 功能 |
| Tier B 扩展接缝 | Admin 功能（运输面归架构） |
| Tier C 体验增强 | Admin 功能 |
| Tier D 真实业务领域 | 业务域 |
| 未进四档的运行时平台 | 架构 |

---

## 架构分支

> 性质：P-005 有界清单。整波退出分母尚未立项冻结。  
> **已冻结（2026-08-20 用户确认）**：Store 双方言决策，见下节。  
> 现状锚点：单进程 + SQLite（`MaxOpenConns=1`）+ 本地盘上传 + 进程内 Job + 内存限流。Compose 已声明非目标含 TLS 终止与多实例（`compose.yaml`）。

### 已冻结：Store 双方言（RT-P03）

用户确认（2026-08-20）。未建 VP、未写驱动。

| 项 | 决定 |
|----|------|
| ORM | **不引入**（GORM / ent / AutoMigrate 均否）。不自研查询构造器或 session。 |
| 方言集合 | 只支持 **PostgreSQL** 与 **SQLite**。禁止第三库、禁止「支持所有数据库」。 |
| 内核 Store | **持久化端口**（连接/事务/占位符/upsert/时间类型/迁移 runner/备份与就绪）。**不是**业务仓库。禁止 `*sql.Tx`、驱动类型进入 handler 与模块公共契约。 |
| 业务对接 | Handler / 其他模块只打**本模块 Repository**。模块拥有表与 Persistence 贡献；逻辑 schema **一份**，物理 SQL 可以按方言成对。 |
| 生产权威 | **PostgreSQL**：生产 fork 推荐与架构验收（升级、备份、共事务、CI）以 PG 为准。 |
| 内嵌默认 | **SQLite**：dev / mvp / 快测 / 当前 `db.path` 与 Compose 卷继续默认。不因「生产首要」强制本地必须有 PG。 |
| 合同平等 | 两实现走同一迁移台账与同一表结构语义。SQLite **不是**可残缺的缩水库；新迁移须两方言都能 apply + checksum。 |
| 迁移 | 仍为全局、不可变、带 checksum 的台账。禁止用 ORM 推 schema。PG 备份合同替换 `VACUUM INTO`，不删 SQLite 快照路径。 |

### 判定

| 状态 | 含义 |
|------|------|
| **delivered** | 主线已有，不再单独立项 |
| **registered** | 已收集；触发后经 `/vision` 立项 |
| **trigger-gated** | 必须先有部署或产品触发（多实例、领域事件、全局搜索等） |
| **default-non-goal** | 列入以免遗忘；默认不做，除非具名 fork 需要 |

### 1. 持久化与数据生命周期

| id | 项 | 现状 | 状态 | 备注 |
|----|----|------|------|------|
| RT-P01 | SQLite 文件库 + 全局迁移台账 + 升级前快照 | `modernc.org/sqlite`；`VACUUM INTO`；模块 Persistence 贡献 | **delivered** | 内嵌默认；合同上与 PG 平等，不得残缺 |
| RT-P02 | PostgreSQL 方言实现 | 无 | **registered** | 生产权威实现；硬问题是迁移方言 + 备份合同 |
| RT-P03 | Store 双方言端口（无 ORM） | Store 即 SQLite 平台 | **registered**（决策已冻结，实现未做） | 见上节；A1 前置 |
| RT-P04 | 连接池 / 读写分离 / replica | `MaxOpenConns=1` | **trigger-gated** | 多实例或 PG 之后才有意义 |
| RT-P05 | 备份 / 恢复 / PITR | SQLite `VACUUM INTO`；PG 逻辑备份 `pg_dump`/`pg_restore`（VP-013 I-004） | **registered**（方言级 dump = **delivered**；轮换后恢复 = **delivered**（VP-016）；PITR 仍 gated） | A1 已交付方言级恢复；A5 已补密钥轮换后的恢复语义，不重做 dump |
| RT-P06 | 加密静止数据 / 表级密钥 | 无 | **trigger-gated** | 合规触发；密钥见 RT-K\* |
| RT-P07 | 文档库（MongoDB 等） | 无 | **default-non-goal** | 第二数据模型；无具名 fork 需求则不做 |
| RT-P08 | 多租户物理隔离（独立库/schema） | 无；Charter 未要求多租户 | **default-non-goal** | 应用层 org 上下文在 Admin 功能分支 |

### 2. 缓存、锁、队列、调度

| id | 项 | 现状 | 状态 | 备注 |
|----|----|------|------|------|
| RT-Q01 | 进程内 Job 六态 | VP-012 R4 | **delivered** | 外部队列当时显式推迟 |
| RT-Q02 | 外部消息队列 / Job broker | 无 | **trigger-gated** | 触发：多实例、跨机长任务、或领域事件要 fan-out。优先评估 PG `SKIP LOCKED` |
| RT-Q03 | 缓存（Redis 等） | 无 | **trigger-gated** | 用途须先钉死：共享限流 / 分布式锁 / 热配置 / 查询缓存。禁止「先上 Redis 再找场景」 |
| RT-Q04 | 分布式锁 / leader election | 无 | **trigger-gated** | 定时任务、单飞 Job；PG advisory lock 可推迟 Redis |
| RT-Q05 | 登录/API 限流跨实例 | 进程内滑动窗口 | **trigger-gated** | 单实例够用；多实例才需要共享存储 |
| RT-Q06 | 事务 outbox / inbox | 无 | **trigger-gated** | Admin 功能「领域事件」契约的运输前置；先 DB outbox，后可选 broker |
| RT-Q07 | 分布式 cron | 定时任务模块，单进程 | **trigger-gated** | 与 RT-Q04 绑定 |

### 3. 对象与文件存储

| id | 项 | 现状 | 状态 | 备注 |
|----|----|------|------|------|
| RT-S01 | 本地盘（头像/品牌/上传/库文件） | `./data` | **delivered** | 内嵌默认；合同上与 S3 兼容实现平等，不得残缺 |
| RT-S02 | 对象存储适配器（S3 兼容） | VP-014 已交付内核端口 + S3 兼容 + 本地盘默认 | **delivered** | 有界 residual：无产品搬运器（I-014-004） |
| RT-S03 | 签名 URL / TTL / 直传 | 无 | **trigger-gated** | 依赖 RT-S02；**不进** VP-014 退出分母 |
| RT-S04 | 分片/断点上传 | 无 | **trigger-gated** | Admin 接缝的平台面；**不进** VP-014 |
| RT-S05 | 恶意内容扫描执行面 | 无 | **trigger-gated** | Admin 定策略，架构接扫描器；**不进** VP-014 |
| RT-S06 | CDN / 公共资源分发 | 无 | **trigger-gated** | 品牌资源公网分发时；**不进** VP-014 |

### 已冻结：对象存储 A2 退出分母（VP-014）

用户确认（2026-08-21）。VP-014 已有界 `closed`；lead `workspace-014-object-storage`（Root `done 5/5`）。

| 项 | 决定 |
|----|------|
| 对象存储方言 | 只支持 **S3 兼容** 与 **本地盘**。禁止 Azure Blob / GCS native 作为第三方言。 |
| 内核 | **对象存储端口**（put / get / delete / exists + 命名空间隔离）。**不是**业务文件管理器。禁止本地路径 / `os.File` 进入 handler 与模块公共契约。 |
| 生产权威 | **S3 兼容**：生产 fork 推荐与本 VP 验收。 |
| 内嵌默认 | **本地盘**：现有 avatars / brand-assets / uploads 继续默认。不因「生产首要」强制本地必须有 MinIO/S3。 |
| 合同平等 | 两实现走同一端口语义。本地盘 **不是** 可残缺的缩水实现。 |
| 读面 | 继续经 API 代理；签名 URL / 直传不进本波。 |
| 存量 | **不提供**产品级本地盘→对象存储搬运器。 |

### 4. 可观测性

| id | 项 | 现状 | 状态 | 备注 |
|----|----|------|------|------|
| RT-O01 | 结构化日志 + correlation / request-id | VP-012 R1 | **delivered** | |
| RT-O02 | `/healthz` `/readyz` | 探活 + 迁移/模块图就绪 | **delivered** | 外部依赖（PG/Redis/S3）接入后须扩展 ready |
| RT-O03 | 指标（Prometheus 等） | VP-015 已交付专用 scrape listener + `suc_*` 系列（含 `module_id`） | **delivered** | 有界 residual：Store/对象/Job 指标不进分母（I-015-003） |
| RT-O04 | 分布式 tracing（OpenTelemetry） | VP-015 已交付 OTLP/HTTP SERVER span + `correlation.request_id` | **delivered** | 缺省 no-op；显式 endpoint 才导出。in-repo sink 不解析 |
| RT-O05 | 剖析 / 连续性能剖析 | 无 | **trigger-gated** | **不进** VP-015 |
| RT-O06 | 错误汇聚（Sentry 类） | 无 | **trigger-gated** | 可后置于 RT-O04；**不进** VP-015 |

### 已冻结：可观测 A4 退出分母（VP-015）

用户确认（2026-08-21）。VP-015 已于 2026-08-22 有界 `closed`；lead `workspace-015-observability`（Root `done 5/5`）。

| 项 | 决定 |
|----|------|
| 指标 | Prometheus 类 pull 面；系列至少携带 `module_id`。不是 Grafana 产品或 Admin 仪表盘。 |
| Tracing | OpenTelemetry / OTLP 导出；与现有 request-id / correlation 可关联。不是替换结构化日志。 |
| 内嵌默认 | 无 Prometheus / collector / Jaeger 仍能开发与快测。不得把收集器做成 mvp/dev 启动硬依赖。 |
| 生产验收 | 显式配置后可核对 scrape **与** 至少一条 trace 导出（可分路径）。 |
| 不进本波 | A3 多实例/Redis/队列；A5 密钥轮换；Sentry（RT-O06）；连续剖析（RT-O05）；Admin 监控页；业务域。 |

### 5. 进程、部署、网络

| id | 项 | 现状 | 状态 | 备注 |
|----|----|------|------|------|
| RT-D01 | 本地双进程 + Compose 一键 | VP-002 | **delivered** | 单 API 容器 + SQLite 卷 |
| RT-D02 | 优雅停机 / 连接排空 | 进程生命周期有，无明确 drain 合同 | **registered** | 多实例与 Job 租约相关 |
| RT-D03 | API 与 worker 进程分离 | Job 跑在 API 进程内 | **trigger-gated** | 长任务与 HTTP 隔离时 |
| RT-D04 | 多实例 / 水平扩展 | Compose 非目标 | **trigger-gated** | 拉动 RT-P02/P04、RT-Q\*、RT-S02、RT-Q05 |
| RT-D05 | TLS 终止 / 证书 | 无；API 不直接暴露 | **trigger-gated** | fork 生产反向代理可外置 |
| RT-D06 | Kubernetes / Helm / Operator | 无 | **default-non-goal** | 不把 K8s 当基架成功条件 |
| RT-D07 | 服务网格 / mTLS 网格 | 无 | **default-non-goal** | |
| RT-D08 | 反向代理、CORS、trusted proxies | Nginx + 可配 CIDR | **delivered** | |

### 6. 密钥、配置注入、密码学平台

| id | 项 | 现状 | 状态 | 备注 |
|----|----|------|------|------|
| RT-K01 | YAML + env 插值；密钥 fail-closed | VP-002/W7 | **delivered** | |
| RT-K02 | 外部 Secret Provider / KMS / HSM | 无 | **trigger-gated** | VP-012 API Token 推迟 HSM |
| RT-K03 | JWT/数据密钥轮换合同 | VP-016 已交付 current+previous；签发只用 current；校验允许重叠窗；缺 previous = 今日单密钥 | **delivered** | 有界 residual：立即失效未选（I-016-005）；`admin.mfa` wrapping 不随 JWT previous 重包。KMS/HSM 仍 gated |
| RT-K04 | 传输中加密（TLS） | 依赖外置代理 | **trigger-gated** | 同 RT-D05 |

### 已冻结：密钥轮换与恢复 A5 退出分母（VP-016）

用户确认（2026-08-22）。VP-016 已于 2026-08-22 有界 `closed`；lead `workspace-016-key-rotation-and-backup`（Root `GOAL-001-key-rotation-and-backup` `done 5/5`）。

| 项 | 决定 |
|----|------|
| 密钥集合 | 默认应用签名密钥 = `AUTH_JWT_SECRET`。禁止把 DB / S3 / 种子密码当成本波轮换对象。 |
| 轮换合同 | 可声明 **current + previous**；新签发只用 current；校验允许重叠窗。缺 previous = 今日单密钥。 |
| 生效 | 进程重启后生效。热加载不进本波。 |
| 备份面 | **不重做** VP-013 方言级 dump。本波只核对：轮换后从既有 SQLite `VACUUM INTO` 与 PG `pg_dump`/`pg_restore` 启动 + 鉴权。 |
| 内嵌默认 | 无 previous、无外部备份代理仍能开发与快测。 |
| 不进本波 | A3 多实例/Redis/队列/优雅停机；KMS/HSM（RT-K02）；静止加密（RT-P06）；TLS；PITR；`/readyz` 再扩；Admin 密钥页；业务域。 |

### 7. 搜索引擎与专用索引

| id | 项 | 现状 | 状态 | 备注 |
|----|----|------|------|------|
| RT-X01 | 专用搜索引擎（Meilisearch / OpenSearch 等） | 无 | **trigger-gated** | 由 Admin 功能「全局搜索」拉动；未立项搜索 UX 前不上引擎 |
| RT-X02 | DB 全文检索 | SQLite FTS 未作为平台能力 | **trigger-gated** | 可作搜索的过渡实现 |

### 8. 时间、标识、时钟

| id | 项 | 现状 | 状态 | 备注 |
|----|----|------|------|------|
| RT-T01 | 请求级 ID | VP-012 | **delivered** | |
| RT-T02 | 分布式 ID / 时间权威 | 应用侧 UUID/随机 token | **trigger-gated** | 多实例写入热点时再评 |
| RT-T03 | 时区在持久化层的合同 | locale 在 Admin 功能面 | **registered** | 与 Admin 功能「时区/数字/货币」接缝；DB `timestamptz` 属架构 |

### 9. 出站消息

| id | 项 | 现状 | 状态 | 备注 |
|----|----|------|------|------|
| RT-M01 | 出站邮件发送端口 + 可切换渠道 | VP-017 v0.5.0 按现行分母 closed：mock 站内记录 + Resend（live PASS）+ 设置热切换/试发 | **delivered** | 2026-08-24 再关门放行（A-003/A-004 pass）；R1～R4 实施史保留 |
| RT-M02 | SMS / 其它推送运输 | 无 | **trigger-gated** | 用户 2026-08-22：审核麻烦，有真实需求再做 |

### 现行：出站邮件 A6 退出分母（VP-017 · 2026-08-24 重开）

用户确认（2026-08-24）：否决同日 SMTP 专用有界关门；实施史不回退。lead `workspace-017-outbound-mail`（Root `active`）。

| 项 | 决定 |
|----|------|
| 运输方言 | **可切换渠道**。第一期：mock（默认）+ Resend（生产）。R2 SMTP 适配器**保留不删**，不再是唯一生产权威。禁止 SMS。 |
| 内核 | **发送端口**（to / subject / text body；默认 From 来自配置）。**不是**通知产品。禁止供应商客户端类型进入 handler 与模块公共契约。 |
| 生产权威 | **Resend**（显式配置）。SMTP 为已实施兼容渠道。 |
| 内嵌默认 | **mock 站内出站记录**（管理员可检视）。无生产渠道仍能开发与快测。历史 `CaptureSink` 实施保留，由 R5/R6 决定是否升级为持久化记录。 |
| 合同平等 | 各渠道走同一端口语义。mock **不是**可残缺的缩水实现（至少可在管理面取出报文）。 |
| 发送模型 | 同步 `Send`。事务 outbox / 外部邮件队列不进本波（RT-Q06 仍 gated）。 |
| 管理面 | 设置「邮件」tab：选渠道、填配置、热切换、试发（同一 `MailSender`）。 |
| 消费者 | 本波用测试/harness + 管理试发。账号 email、校验、邀请、自助恢复不进本波（VP-018 已 `closed` 交付账号邮箱身份；邀请 / 自助恢复 = VP-019）。 |

历史 SMTP 专用分母（2026-08-22 冻结、2026-08-24 曾用于有界关门）保留在 VP-017 历史关门节，**不再**作为组合层成功定义。

### 10. 明确不作为基架一等公民（除非具名触发）

| id | 项 | 理由 |
|----|----|------|
| RT-N01 | MongoDB / 泛「支持所有数据库」 | 第二模型与无限方言会拆内核；双方言集合已冻结为 PG + SQLite |
| RT-N06 | 引入 ORM（GORM / ent / AutoMigrate） | 与全局 checksum 台账、模块拥有 Persistence、薄内核边界冲突；查询层也不用 ORM 换方言 |
| RT-N02 | 先上 Kafka/Rabbit 再补事件模型 | 与 RT-Q06 顺序相反 |
| RT-N03 | 微服务拆分 / 运行时插件市场 | VP-003 非目标 |
| RT-N04 | GraphQL 网关、CQRS/事件溯源默认化 | 无 Charter 要求 |
| RT-N05 | 多云、服务网格、K8s Operator | 部署细节交给 fork |

### 架构分支建议顺序（草案，未冻结）

```text
A0  本清单（已登记）；Store 双方言决策已冻结（RT-P03）
A1  内核持久化端口 + PostgreSQL 实现 + 现有台账对写/翻译；
    SQLite 保留为 dev/mvp/快测默认；生产 CI 以 PG 为权威
A2  对象存储适配器；本地盘保留为默认
A3  仅当需要多实例：就绪探针扩依赖、优雅停机、
    再评估 PG 锁/SKIP LOCKED vs Redis vs 外部队列
A4  指标 + OpenTelemetry（可与 A1 部分并行）
A5  密钥轮换 / 备份恢复合同（随 A1 或紧随其后）
A6  出站邮件：内核发送端口 + 可切换渠道（mock 默认 + Resend 生产）；SMTP 适配器保留；SMS 不进
```

**刻意后置**：MongoDB、ORM、Redis、消息队列、搜索引擎、K8s、SMS。它们是部署或产品触发的后果，或已否决的技术选型。

架构分支当前拍：**[VP-017-outbound-mail](plans/VP-017-outbound-mail.md) `closed`**（2026-08-24 按现行渠道分母再关门 · v0.5.0；RT-M01 delivered）。A3 仍 trigger-gated（多实例才评估就绪探针扩依赖 / 优雅停机 / PG 锁 vs Redis vs 队列）。

---

## Admin 功能分支

承接原四档的 **A 剩余 + B + C**。历史证据：`workspace-011` 的 `I-011-002`。  
VP-011 已交付的标准 Admin 模块（用户/角色/设置/钱包演示面等）不重开。本分支只登记**尚未立项**的通用 Admin 能力。

横切**契约**首波已由 VP-012 `closed` 交付（correlation / 审计 / 并发幂等 / 进程内 Job / maintenance / API Token；保留/归档与 session envelope 已在 workspace-012 增量完成）。effective actor 冻结为当前 actor；不做 impersonation，除非再次出现产品触发。

### 仍开放（非立即实施）

**产品事实（2026-08-22 用户确认）**

- 忘密码要 **自助恢复** 与 **管理员重置** 两种。
- 自助恢复的证明依据 = 用户持有事先绑在账号上的 **邮箱**（验证码或链接只是投递形态）。没有已校验邮箱 + 出站投递，登录页「忘记密码」是空转。
- 出站通道 **先做邮件**；**SMS 后置**（审核成本高，有真实需求再做）。
- 管理员重置继续走既有特权路径（`must_change_password`），不冒充自助恢复。
- 消费链：**A6 出站邮件（VP-017 已按现行分母再 `closed` v0.5.0）→ 账号邮箱身份（VP-018 已 `closed` v1.0.0）→ IAM（密码策略 / 邀请 / 自助恢复状态机 = [VP-019-iam-recovery](plans/VP-019-iam-recovery.md) **`active`**，2026-08-25 激活）**。

**基架能力剩余**

1. 密码策略、邀请、账号恢复状态机（自助恢复 **硬前置** = VP-017 + 账号邮箱；邀请仍可用管理员出示链接）  
2. 组织 / 部门 / 岗位，以及数据权限 `org` 扩展  
3. 配置包导出、diff、dry-run、导入  
4. 文件扫描 / 隔离**策略**（执行器见架构 RT-S05）  
5. 时间、时区、数字、货币格式语义（持久化时区合同见架构 RT-T03）  
6. impersonation / effective actor 产品化——仅当再次出现  
7. **账号邮箱身份**（`users` email 字段 + 校验状态机）——由 **VP-018** 承接并于 2026-08-24 **已 `closed`**（v1.0.0）；自助恢复的身份前置；运输面见架构 RT-M01（delivered）

**扩展接缝**（全部 trigger-gated；运输实现归架构）

typed domain event、Notification Transport、OIDC/SSO/SCIM、Approval Gate、Entitlement、多组织 context、SSE/WebSocket、外部连接器/Secret 的产品面、自定义 metadata/tags、文件预览。

**体验增强**

全局搜索 / Command Palette、Saved Views、批量结果中心、未保存保护、统一 Toast/错误恢复、版本与维护提示。全局搜索若需要专用引擎，拉动架构 RT-X01。

Admin 功能下一拍：**[VP-019-iam-recovery](plans/VP-019-iam-recovery.md)（IAM：密码策略 / 邀请入职 / 自助恢复状态机）——2026-08-25 已激活（`active`）**；硬前置 = VP-018 已校验邮箱（已 `closed` v1.0.0）+ VP-017 运输（已按现行分母再 `closed` v0.5.0）均满足。不要把恢复状态机打进 VP-018。密码策略虽可并行规划，本波三者同为退出分母。再下一截（未立项）：组织/部门/岗位 + 数据权限 `org` 等，见「基架能力剩余」。

---

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

已有的钱包模块是 VP-011 交付的 Admin 演示/能力面，**不**等于本分支「支付/结算」业务域已成立。

业务域下一拍：仅当某个域有真实业务触发时 `/vision` 建该域 VP；不要把多个域打进同一个 VP，也不要在无触发时预先开区。

---

**当前组合焦点**：**[VP-019-iam-recovery](plans/VP-019-iam-recovery.md)**（Admin 功能 · IAM：密码策略 / 邀请入职 / 自助恢复状态机；2026-08-25 **`active`**，VRev-043 independent `pass`）。[VP-017-outbound-mail](plans/VP-017-outbound-mail.md) 已于 2026-08-24 按**现行渠道分母**再 `closed`（v0.5.0 · 架构 A6；RT-M01 delivered）；**[VP-018-account-email-identity](plans/VP-018-account-email-identity.md) 已于 2026-08-24 同日 `closed`**（v1.0.0 · 账号邮箱身份）。[VP-016-key-rotation-and-backup](plans/VP-016-key-rotation-and-backup.md) 已于 2026-08-22 有界 `closed`（架构 A5）。**[VP-015-observability](plans/VP-015-observability.md) 已于 2026-08-22 有界 `closed`**（架构 A4）。**[VP-014-object-storage](plans/VP-014-object-storage.md) 已于 2026-08-21 有界 `closed`**（架构 A2）。**[VP-013-store-dialects](plans/VP-013-store-dialects.md) 已于 2026-08-21 有界 `closed`**（架构 A1）。后续方向按 **架构** / **Admin 功能** / **业务域** 三分支并行登记。持续程序 = **VP-009 `active`** 与 **VP-010 `active`**。VP-001～008、VP-011～018 均为历史 `closed`（VP-017 为 2026-08-24 按现行分母再关门；VP-018 同日关门）。VP-008 `go` 消费有效性在无新的共享基架阻断时保持可消费。协议覆盖权威 `I-PROTO-FULL-001`（v2.7.0 历史分母，被 v2.8.0 覆盖）。

## 单主线模块化策略

未来 fork 起点统一由同一代码主线、模块候选集与启动时 Profile 表达，权威见 [module-architecture.md](../architecture/module-architecture.md) 和 VP-003。原 [dual-track-contract.md](dual-track-contract.md) 已转为历史记录。
