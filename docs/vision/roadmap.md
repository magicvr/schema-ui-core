---
doc_type: vision-roadmap
title: 愿景组合编排
status: active
created: 2026-07-31
updated: 2026-09-20
parent: null
version: 0.96.0
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
| 16 | [VP-016-key-rotation-and-backup](plans/VP-016-key-rotation-and-backup.md) | 架构 A5：JWT current+previous 轮换合同 + 既有备份上的轮换后恢复；单密钥仍为内嵌默认。 | VP-015 A4 已 closed；RT-K01 delivered；VP-013 方言级 dump 已交付；与 VP-009/010 正交；不进 A3 / KMS / PITR / 热加载 / Admin 页 / 业务域 | **closed**（2026-08-22 有界关门 · 架构 A5；lead: workspace-016-key-rotation-and-backup；Root done 5/5；VRev-036 `pass`；residual：I-016-005 立即失效未选 + `admin.mfa` wrapping 残余**已收窄**（现行代码支持 previous 解密 + 成功 TOTP 后惰性重包，W11 F-004；残余 = 仅凭恢复码路径不重包 / 无启动批量重包 / 无主动轮换重包，用户 2026-09-10 接受）） |
| 17 | [VP-017-outbound-mail](plans/VP-017-outbound-mail.md) | 架构 A6 升级：内核发送端口 + 可切换渠道（默认 mock 站内出站记录 + 生产 Resend；SMTP 适配器保留不删）+ 管理设置/试发。 | 用户 2026-08-24 否决同日 SMTP 专用有界关门（实施史不回退）。不进账号 email / 邀请 / 恢复状态机 / 模板 / 用户站内通知 / SMS / A3 | **closed**（2026-08-24 按**现行渠道分母**再关门 · v0.5.0；lead: workspace-017-outbound-mail；Root `done` 8/8；R5～R8 = GOAL-006～009 done；live 投递实跑 PASS；VRev-042） |
| 18 | [VP-018-account-email-identity](plans/VP-018-account-email-identity.md) | Admin 功能：账号邮箱身份（`users` email 可空 + 绑定/校验状态机 + 换绑）；消费 VP-017 `MailSender`。 | 硬前置 = VP-017 **再次** `closed`（现行渠道分母）。不进 IAM 恢复 / 邀请 / 密码策略 / SMS / 模板 / A3 | **closed**（2026-08-24 解冻当日关门 · v1.0.0；lead: workspace-018-account-email-identity；Root `done` 4/4；VRev-040 pass；A-002 independent 归零） |
| 19 | [VP-019-iam-recovery](plans/VP-019-iam-recovery.md) | Admin 功能 · IAM：密码策略 / 邀请入职 / 自助恢复状态机（忘密全链消费 VP-018 已校验邮箱 + VP-017 `MailSender`）。 | 硬前置 = VP-018 邮箱身份 + VP-017 运输（均已 `closed`）。不进 SMS / 模板中心 / 多邮箱 / 组织权限 / OIDC / 业务域 / A3；不改 Profile 默认集 | **closed**（2026-08-26 交付后关门 v0.3.0 · 用户书面确认；实现 2026-08-25 同日全链交付，Root done 4/4；关后 A-001 independent `pass` + A-002 recommended ×2 闭合；lead: workspace-019-iam-recovery） |
| 20 | [VP-020-timezone-number-currency-formatting](plans/VP-020-timezone-number-currency-formatting.md) | Admin 功能 · 时区 / 数字 / 货币**格式语义**：会话/用户级时区 + locale 数字/货币展示与输入合同（消费 VP-007 locale 运行时 / VP-005 设计系统）。 | 继承 VP-007 locale / VP-005 设计系统 / VP-011 用户角色边界；Admin 类 freshness（VP-019 时效 `66f5fd1f`）；DB 时区持久化合同仍归架构 RT-T03 | **closed**（v0.3.0 · 2026-08-27 关门 · 用户书面确认；lead workspace-020 结项 Root done 4/4；关门审计 A-001 self + A-002 grok independent 双 pass；VRev-045 self `pass`） |
| 21 | [VP-021-graceful-shutdown-and-connection-drain](plans/VP-021-graceful-shutdown-and-connection-drain.md) | 架构 RT-D02：优雅停机 / 连接排空**合同**（停机顺序、HTTP drain、运行中 Job 语义、双方言 Store 排空）；单进程 + Compose 基线。 | 继承 VP-012 Job 六态、VP-013 双方言 Store；与 VP-009/010 正交；A3 余项仍 trigger-gated | **closed**（v0.3.0 · 2026-08-27 关门 · 用户指令授权；lead workspace-021 结项 Root done 3/3；关闭双审 A-001 self + A-002 grok independent 0 开放；VRev-047 self `pass`；RT-D02 → **delivered**；residual = 进程级 harness / compose stop 以 linux CI 核销） |

| 24 | [VP-024-distribution-formalization](plans/VP-024-distribution-formalization.md) | 分发形态正式化：cli+包 对外服务化——serve 壳 / npmjs 公开发布 / compose CI 实跑 / fork 对照计时 / renderer external 化 / 纯原子拆分 / 迁移工具化 + 方法 B 置顶宣告（收口 VP-022/023 go 后残余 7 项）。 | 继承 VP-022/023（合并 go 后清单）+ VP-003/005/006（pin 2.9.0）/008；**组合层平台波，与三分支正交**；不改 Charter（fork 与包并存维持）；与 VP-009/010 正交 | **closed**（2026-08-29 关闭 · v0.3.0；VRev-053 independent `pass`；八判据核销（GOAL-008 closure-report）；方法 B 置顶；残余四项登记（GOAL-008 D-001：hosted CI 实触发 / shell 类型面 / GH Packages 保留 / C 类 fork 包化候选）；lead workspace-024-distribution-formalization · Root done 7/7） |
| 25 | [VP-025-config-export-diff-dryrun-import](plans/VP-025-config-export-diff-dryrun-import.md) | Admin 功能 · 配置包导出 / diff / dry-run / 导入（基架能力剩余 #3）：可移植配置包 + 键级差量 + 只读预检 + 安全导入；密钥 fail-closed 保持、热加载不进分母、不改 Profile 默认集/Manifest（VP-008 `go` 红线）。 | 继承 RT-K01 配置系统（YAML+env 插值）与 VP-023/024 CLI/包产线 + VP-003 模块边界 / VP-007 设置面；**Admin 功能分支非门控未立项项（roadmap 点名）**；与 VP-009/010 正交 | **closed**（2026-08-30 v0.3.0 关闭 · 用户书面确认；六判据全满足（r4-evidence-matrix）· 关门双审 A-001 self `pass` + A-002 grok build independent（F-001～F-008 全 fixed · 开放 required=0）· VRev-055 `pass`；lead workspace-025 · Root `done` 4/4） |
| 23 | [VP-023-productionization-cli-package](plans/VP-023-productionization-cli-package.md) | 包消费产线化：发布运营（Go tag/proxy go get + npm registry）+ CLI（create/add/upgrade，对标 dotnet new + NuGet）+ 六包细化与 d.ts 自动化 + PG external/运维 + golden-field 从零上线与 fork→包迁移指南。 | 继承 VP-022（dist-lib/pack/双 golden/冻结面 v1.2.0/go 后清单）+ VP-003/005/006（pin 2.9.0）/008；**组合层平台波，与三分支正交**；不改 Charter（fork 与包并存维持）；与 VP-009/010 正交 | **closed**（2026-08-29 v0.3.0 关闭 · 用户 P-004 裁决 breaking 实演 v0.3.0 真实执行；六条判据达成；grok 独立双审 F-001～F-008 全闭合；lead workspace-023-productionization-cli-package · Root `done` 5/5；go 后清单 = serve 壳 / 六包 external 化 / 纯原子拆分 / fork 对照计时 / 迁移工具化 / 包公开可见性 / compose CI 实跑，已并入 VP-024 收口） |
| 22 | [VP-022-distribution-package-pilot](plans/VP-022-distribution-package-pilot.md) | 分发形态试点：**构建期包消费**最小闭环（Go 库模块 + npm 包组 + 空下游仓零冲突升级演练；对标 dotnet new + NuGet / Spring Boot starters）；不改 Charter、不弃 fork（保留为深度定制逃生舱）。 | 继承 VP-003 模块契约 / VP-004 playbook（仅评估）/ VP-005 主题覆盖 / VP-006 协议面 / VP-008 `go` 消费基线；**组合层平台波，与三分支正交**；与 VP-009/010 正交 | **closed**（2026-08-29 v0.4.0 关闭 · 用户 P-004 裁决：六条判据按有界口径满足；independent 双审闭合；Charter 0.3.0 strategic 随 GO 裁决落地（VR-050）；lead workspace-022 Root done 5/5） |
| 26 | [VP-026-cache-port](plans/VP-026-cache-port.md) | 架构 · **通用缓存端口**（H-002 同进程基座早期化 · 承接 RT-Q03）：Cache 端口（Get/Set/Delete + TTL）+ 绝对/滑动过期 + 可插拔策略接口 + **内存供应商（默认）** + **Redis 接缝声明（不实现）**；.NET IMemoryCache 式轻量分层。 | 继承 VP-003 模块契约 + Charter 0.4.0 成功边界 #6 / H-002；与 VP-027/028 按"触发条件独立 × 关门能力独立"分立；Redis 实现仍 trigger-gated；与 VP-009/010 正交 | **closed**（2026-09-01 v0.3.0 · **用户书面确认关门**：八条判据证据矩阵 verified · R1～R4 阶段 self + grok independent 双审闭合（开放 required=0）· VRev-061 `pass`；lead `workspace-026-cache-port` · Root `GOAL-001-cache-port` `done` 4/4；Redis 实现仍 trigger-gated（不消耗 RT-Q03）） |
| 27 | [VP-027-rate-limiter-port](plans/VP-027-rate-limiter-port.md) | 架构 · **通用限流器端口**（H-002 早期化 · 承接 RT-Q05）：RateLimiter 端口（Allow/Record/Reset/RetryAfter）+ 滑动窗口**内存供应商**（演进既有 loginRateLimiter）+ **7 处使用点完整迁移**（含 MFA verify 独立桶 / 邀请接受）+ **Redis 接缝声明（不实现）**；W12 D-002 窗口常量保持。 | 继承 VP-003 模块契约 + Charter 0.4.0 成功边界 #6 / H-002；与 VP-026/028 分立；Redis 实现仍 trigger-gated；与 VP-009/010 正交 | **closed**（2026-09-01 · v0.3.0 · **用户书面确认关门** · VRev-063 self `pass` · 判据 #1～#7 证据矩阵 7/7 · Root 双审 0 required；lead workspace-027-rate-limiter-port · Root done 4/4） |
| 28 | [VP-028-event-bus-port](plans/VP-028-event-bus-port.md) | 架构 · **进程内事件总线运输端口**（H-002 早期化 · 承接 RT-Q02 运输端口前置）：类型化 EventBus（Publish/Subscribe/Unsubscribe）+ 进程内 channel 实现 + **outbox/MQ 接缝声明（不实现）**；**不解除** Admin 功能分支 typed domain event 扩展接缝的 trigger-gated（应用契约仍归 Admin 功能）；EventBus ≠ Job 端口。 | 继承 VP-003 模块契约 + Charter 0.4.0 成功边界 #6 / H-002；与 VP-026/027 分立；outbox/broker 仍 trigger-gated；不重开 VP-012；与 VP-009/010 正交 | **closed**（2026-09-01 · v0.3.0 · Root `done` 4/4；lead `workspace-028-event-bus-port`；outbox/broker 仍 trigger-gated，不消耗 RT-Q02） |
| 29 | [VP-029-wallet-prepaid-instrument](plans/VP-029-wallet-prepaid-instrument.md) | Admin 功能 · **钱包预付资金凭证 + 外部主体接缝**：`(issuer, external_id) → subject_id`（不创建 `admin.users`）+ 卡密批次生成/导出/作废/核销入账（哈希存储、幂等 Redeem）。**R5**：Admin 已登录自助核销 HTTP + 「我的钱包」入口（入账 `owner_type=user`）。扩展 `admin.wallet`，**不是**支付业务域。 | 继承 VP-011 钱包账本；不重开 VP-011；与 VP-030/031 同批；硬前置于 030 身份与 031 扣款主体 | **closed**（2026-09-02 · v0.5.0 · 用户指令授权 · VRev-069 self `pass` · 十条判据全量 verified · Root done 5/5 · GOAL-005 独立审与 Root 关门自审 pass；lead `workspace-029-wallet-prepaid-instrument` 结项） |
| 30 | [VP-030-telegram-channel-runtime](plans/VP-030-telegram-channel-runtime.md) | 架构 · **C 端 Telegram 通道运行时**（对标 VP-017）：webhook + Update 分发端口 + SendMessage 文本 + `issuer=telegram` 主体映射 + Admin bot 设置。**不是**业务域、**不是**付费命令实现。 | 硬前置 = VP-029 主体接缝（已交付）；消费 VP-027 限流；激活前评估 C 端桶已落盘（进程内够用、不需要 Redis）；与 VP-009/010 正交 | **closed**（2026-09-05 · v0.3.0 · VRev-076 self `pass` · lead `workspace-030-telegram-channel-runtime` · Root done；R-009 bounded residual） |
| 31 | [VP-031-digital-offer-entitlement](plans/VP-031-digital-offer-entitlement.md) | 业务域 · **数字 Offer + 薄购买凭证 + 权益**（本仓首个业务域 VP）。服务视为可售 Offer，不是电商类目/SKU/税/库存/物流订单。 | VP-029 硬前置已满足；H-002 同进程书面确认；RT-Q03/Q05 均评估为本波不需要 Redis；R1 required 信息未冻结前不得进入 R2 | **closed**（2026-09-05 · v0.3.4 · 第 3 次关门：A-012 `conditional` 2 required → D-004 fixed ×2 + A-013 closed ×2 → A-014 independent `pass` 0 required → F-001 前置加固（A-015）· lead `workspace-031-digital-offer-entitlement` · Root `done` 4/4；可作为后继 VP 已验证前置） |
| 32 | [VP-032-rate-limiter-atomic-port](plans/VP-032-rate-limiter-atomic-port.md) | 架构 · **限流器端口原子化**（GOAL-001 A-008 R-007 residual 承接）：`kernel.RateLimiter` 新增原子 `AllowRecord` 与令牌化 `Reserve`/`Cancel`，迁移冻结 14 处使用点（4 立即消费 + 10 失败预算），消除 Allow→Record TOCTOU；`Allow`/`Record` 保留兼容；内存供应商实现；Redis 仍 RT-Q05 trigger-gated。 | 继承 VP-027 端口语义（**不重开** VP-027 关门事实）；VP-030 三桶限流直接受益；与 VP-009/010 正交 | **closed**（2026-09-04 · v0.3.0 · **用户书面确认关门** · VRev-074 self `pass` · 五判据全部 verified（E-004 矩阵）· Root A-001 self + A-002 grok independent 双 `pass` 0 required；lead workspace-032-rate-limiter-atomic-port 结项 · Root `GOAL-001-rate-limiter-atomic-port` `done` 3/3；失败预算口径承接 = GOAL-003 D-002 令牌化；Redis 实现仍 trigger-gated） |
| 33 | [VP-033-telegram-operator-console](plans/VP-033-telegram-operator-console.md) | Admin 功能 · **Telegram Bot 人工控制台**：连接状态（`getMe` / `setWebhook`）+ 入站模式开关（webhook \| 单实例 `getUpdates`）+ 业务占用位 + 未绑定人工 IM（代 bot 发言、无权限灰掉）。消费 VP-030 runtime，**不是**业务域。 | 硬前置 = VP-030 通道运行时已交付（不重开 030）；与 VP-031 占用位衔接；SSE/多实例 polling 仍 gated；与 VP-009/010 正交 | **closed**（2026-09-05 · v0.3.0 · VRev-077 self `pass` · lead `workspace-033-telegram-operator-console` · Root done 4/4） |
| 34 | [VP-034-nav-group-collapsible](plans/VP-034-nav-group-collapsible.md) | Admin 功能 · **导航分组折叠体验**：Admin Shell 左侧导航引入可折叠/展开分组（Group）；分组与模块解耦（模块可将导航注册到跨模块共享分组）；当前已注册 sidebar 导航全部纳入合理分组迁移与验证（`identity-access` / `content-data` / `operations` / `communications` / `commerce`）；直接 URL 进入时对应分组自动展开；playbook 更新注册规范。 | 继承 VP-003 模块架构 + VP-004 贡献 playbook + VP-005 设计系统；无硬前置；与 VP-009/010 正交 | **closed**（2026-09-09 · v0.4.0 · 用户书面确认 · VRev-085 self `pass` · 七条判据 verified · lead `workspace-034-nav-group-collapsible` · Root done 5/5；residual = Dashboard 现行 `workspace` 组 / GOAL-003） |
| 35 | [VP-035-foundation-architecture-health](plans/VP-035-foundation-architecture-health.md) | 架构 · **基架架构健康评估与路线图重述**：as-built 对照（内核/组合根/端口/Profile/文档）+ 有界业界对照（四类参照集，分类输入不是决策源头）+ 下一版总路线图草案交 `/vision` editorial。 | 基架交付波 VP-013～034 已收口；不替代 VP-009/010；不改 Charter；不消耗 Redis/MQ/多实例 trigger | **closed**（2026-09-10 · v0.3.0 · 六条方向级判据全部达成 · VRev-088 editorial + **VRev-089 self `pass`** · **A-019 grok build 4.6 high independent `pass`** · VR-075/VR-076 · lead `workspace-035-foundation-architecture-health` Root `done · 4/4`；激活记录 2026-09-09 v0.2.0） |

| 36 | [VP-036-admin-command-palette](plans/VP-036-admin-command-palette.md) | Admin 功能 · **全局检索与 Command Palette**：以 `SearchableItem` / provider 接缝提供权限安全的页面、导航项与声明式动作检索；快捷键、键盘操作、直接跳转与导航分组联动。实体全文搜索不进首波。 | 继承 VP-034 导航分组、VP-003/004 模块贡献与 VP-005/007 体验基线；与 VP-009/010 正交；`RT-X01`/`RT-X02`、Redis/MQ/多实例仍 gated | **closed**（2026-09-14 · v0.3.0；VRev-093 self `pass`；七条方向级退出判据 verified；lead: `workspace-036-admin-command-palette`；Root `GOAL-001-admin-command-palette` done 4/4；open required = 0） |
| 37 | [VP-037-admin-workflow-continuity](plans/VP-037-admin-workflow-continuity.md) | Admin 功能 · **工作流连续性与安全反馈**：Saved Views（用户级列表视图保存与恢复）+ 未保存变更保护 + 统一 Toast/错误恢复 + 通用列表页视觉/筛选体验收敛；不承载实体全文检索、批量结果中心或业务域。 | 继承 VP-034 导航分组、VP-036 发现入口、VP-005/007 体验基线与 VP-012 横切契约；激活前须 Admin 类 freshness / VP-008 `go` 消费有效性；与 VP-009/010 正交；R6 不改变 shell、查询/重置合同或 Saved View 存储格式 | **closed**（2026-09-18 · v1.7.0 · 用户书面确认 · VRev-096 self `pass`；lead: `workspace-037-admin-workflow-continuity`；Root `GOAL-001-admin-workflow-continuity` **done · 6/6**；R1～R6 全部 done；非纲领整改子目标 GOAL-008/009/010/011 均 done · 4/4；open required = 0；关门后残余已统一收口/登记，见下节「未决项统一登记」） |
| 38 | [VP-038-batch-operations-and-job-center](plans/VP-038-batch-operations-and-job-center.md) | Admin 功能 · **批量操作与异步结果中心**：把 VP-012 显式排除的「通用 Job 管理页」与 roadmap「体验增强」未立项的「批量结果中心」收口为有界产品能力——已注册 Job 种类的列表/详情/进度/结果读取，**至少一条**真实批量操作以异步 Job 承接（202 + jobId），结果中心体验与权限/Profile 过滤；不重开 VP-012，不改同步 `batch-delete` 已交付语义。 | 继承 VP-012 Job 六态运行时与 `wallet.reconcile` 先例、VP-011 批量/导出导入面、VP-037 统一反馈与列表基线、VP-005/007 体验基线；**`I-038-004` 用户 2026-09-19 裁决 = 新建 `admin.jobs` 进 admin 默认集（Profile 内容扩展，不改装配语义，不暂挂 `go`）** + Admin 类 freshness PASS；与 VP-009/010 正交；Redis/MQ/多实例/搜索引擎仍 gated | **closed**（2026-09-19 激活并**同日关门** · v1.0.0 · **用户书面确认**；计划 self = VRev-098 `pass` · 激活 self = VRev-099 `pass` · 补做关门 Vision Review = VRev-100 self `conditional` → 报告内响应 `fixed` 后开放 required = 0；lead: `workspace-038-batch-operations-and-job-center` · Root `GOAL-001-batch-operations-and-job-center` **`done · 5/5`**；freshness `0c29c08` → `7e5ce891` PASS；关门后残余登记于「未决项统一登记」） |
| 39 | [VP-039-version-maintenance-diagnostics](plans/VP-039-version-maintenance-diagnostics.md) | Admin 功能 · **版本更新、维护提示与诊断报告**：承接 VP-012「运行时管理 UI 可后置」——维护/降级/只读持久横幅 + 版本身份/升级入口 + 已交付探活/就绪/状态字段的轻量诊断摘要。不重开 VP-012/015/025。 | 继承 VP-012 四模式写门禁与 Host 投影、VP-007 locale、VP-037 统一反馈、已有 `admin.system-monitoring` 状态行；**`I-039-004` 接受默认候选（不新模块、不改默认集，不暂挂 `go`）**；与 VP-009/010 正交；不并入 C1 | **closed**（2026-09-19 · v0.3.0；VRev-103 self `pass`；lead `workspace-039-version-maintenance-diagnostics`；Root `GOAL-001-version-maintenance-diagnostics` `done · 4/4`；用户书面确认；existing fresh-seed harness residual 保持登记） |
| 40 | [VP-040-timestamptz-persistence-contract](plans/VP-040-timestamptz-persistence-contract.md) | 架构 · **C1 DB 时间列 timestamptz 持久化合同**（`RES-T03-tz` / `RT-T03`）：PG 生产权威时间类型 + SQLite 合同平等物理类型 + 纳入分母的时间列双方言迁移；不引入 ORM/第三库。 | 继承 VP-013 双方言端口与 checksum 台账、VP-020 展示/输入时区合同、VP-035 评估；VP-039 波次已 `closed`；与 VP-009/010 正交；不消耗 A3/Redis/MQ | **active**（2026-09-20 · v0.2.0 · lead: workspace-040-timestamptz-persistence-contract；激活 self = VRev-104 `pass`；`I-040-001` collecting，R1 默认候选已登记） |

## 组合门闩（用户 2026-08-08）

1. **协议优先于视觉**：在 VP-006 未 `closed` 前，**不得**将 VP-005 作为 `primary_plan` 推进实现，不得启动视觉优化波次。  
2. **MVP 子集不是终态成功条件**：`I-PROTO-001 v0.1.3` 是历史 MVP 冻结切片；整份 v2.7.0 契约由 VP-006 收口。  
3. 已关闭 VP-001～004 的历史证据与 status **不重写**。

> **协议 pin 注记（2026-08-14 · VR-020 · editorial）**：协议覆盖权威由 `v2.7.0` 升至 `v2.8.0`（additive 超集；`I-PROTO-FULL-001` v1.0.1 仍为 v2.7.0 历史分母、已被 v2.8.0 覆盖）。身份权威见 `apps/web/src/protocol/upstream/provenance-v2.8.json`。  
> **协议 pin 注记（2026-08-29 · VR-050 · strategic）**：协议来源升至 **`v2.9.0`**（pinned `81aa1d8`；支持窗 2.7–2.9 additive 超集——代码实现追认，I-007 闭合；Charter `@0.3.0`）。

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
> 现状锚点（2026-09-10 经 VP-035 R2/R3 复核修正）：**单进程** + **SQLite 文件库（小连接池，默认 4；内存库 1）/ PostgreSQL 双方言** + 本地盘上传（S3 兼容适配器已交付） + 进程内 Job（六态） + 内存限流（原子窗口） + 进程内事件总线 + 可选 Prometheus 指标与 OTLP traces（缺省关闭） + JWT current/previous 轮换合同。Compose 已声明非目标含 TLS 终止与多实例（`compose.yaml`）。
>
> 修正依据：[VP-035](plans/VP-035-foundation-architecture-health.md) as-built 矩阵（`apps/api/internal/store/store.go:29` `sqlitePoolDefault = 4`、`:104`–`113` 内存库 1）与 A2/A4/A5/A6/A7 交付事实。

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
| **planned** | 已立项 VP 承接（`plans/VP-0NN` 已存在）；交付后转 delivered |
| **trigger-gated** | 必须先有部署或产品触发（多实例、领域事件、全局搜索等） |
| **default-non-goal** | 列入以免遗忘；默认不做，除非具名 fork 需要 |

### 1. 持久化与数据生命周期

| id | 项 | 现状 | 状态 | 备注 |
|----|----|------|------|------|
| RT-P01 | SQLite 文件库 + 全局迁移台账 + 升级前快照 | `modernc.org/sqlite`；`VACUUM INTO`；模块 Persistence 贡献 | **delivered** | 内嵌默认；合同上与 PG 平等，不得残缺 |
| RT-P02 | PostgreSQL 方言实现 | VP-013（A1）已交付 | **delivered** | 生产权威实现；迁移方言 + 备份合同随 A1 收口（`pg_dump`/`pg_restore`、共事务、`readyz`） |
| RT-P03 | Store 双方言端口（无 ORM） | 内核持久化端口 + SQLite/PG 双方言实现（VP-013 交付） | **delivered** | A1 已 closed（2026-08-21）；全局 checksum 台账双 apply、公共面无 `*sql.Tx` |
| RT-P04 | 连接池 / 读写分离 / replica | 小连接池（SQLite 文件库默认 4、内存库 1；`sqlitePoolDefault`）；**读写分离 / replica 未实现** | **trigger-gated** | 池化本身已交付；多实例或 PG 之后才有意义。**不得**把已有池化扩写成读写分离 / replica 已交付 |
| RT-P05 | 备份 / 恢复 / PITR | SQLite `VACUUM INTO`；PG 逻辑备份 `pg_dump`/`pg_restore`（VP-013 I-004） | **registered**（方言级 dump = **delivered**；轮换后恢复 = **delivered**（VP-016）；PITR 仍 gated） | A1 已交付方言级恢复；A5 已补密钥轮换后的恢复语义，不重做 dump |
| RT-P06 | 加密静止数据 / 表级密钥 | 无 | **trigger-gated** | 合规触发；密钥见 RT-K\* |
| RT-P07 | 文档库（MongoDB 等） | 无 | **default-non-goal** | 第二数据模型；无具名 fork 需求则不做 |
| RT-P08 | 多租户物理隔离（独立库/schema） | 无；Charter 未要求多租户 | **default-non-goal** | 应用层 org 上下文在 Admin 功能分支 |

### 2. 缓存、锁、队列、调度

| id | 项 | 现状 | 状态 | 备注 |
|----|----|------|------|------|
| RT-Q01 | 进程内 Job 六态 | VP-012 R4 | **delivered** | 外部队列当时显式推迟 |
| RT-Q02 | 外部消息队列 / Job broker | 无 | **trigger-gated** | 触发：多实例、跨机长任务、或领域事件要 fan-out。优先评估 PG `SKIP LOCKED`。**运输端口前置 = [VP-028](plans/VP-028-event-bus-port.md) `closed`**（2026-09-01 关门 · v0.3.0 · lead `workspace-028-event-bus-port` Root done 4/4；进程内 EventBus 运输端口 + outbox/MQ 接缝声明已交付；broker 运输仍 gated；不解除 Admin typed domain event gated——应用契约归 Admin 功能分支；不消耗 RT-Q02 trigger） |
| RT-Q03 | 缓存（Redis 等） | 无 | **trigger-gated** | 用途须先钉死：共享限流 / 分布式锁 / 热配置 / 查询缓存。禁止「先上 Redis 再找场景」。**触发条件（H-002 · VR-052）**：多实例部署 **或** C 端业务域模块正式接入同进程。业务域 VP 激活即视为触发条件成立，架构分支须在该 VP 开区前完成评估并在路线图中登记位置（可选「不需要」结论，但评估本身不可跳过）。**承接 = [VP-026-cache-port](plans/VP-026-cache-port.md) `closed`**（端口 + 内存默认 + 双策略 + Redis 接缝声明已交付；Redis 实现仍 gated——VP-026 不消耗 trigger）。**VP-031 激活评估（VRev-080 · 2026-09-05）**：同进程首波的 Offer/权益读取以权威存储为正确性来源，现有内存 Cache 仅按性能证据选择性使用；结论 = 本波不需要 Redis。多实例或跨实例共享失效需求触发复审 |
| RT-Q04 | 分布式锁 / leader election | 无 | **trigger-gated** | 定时任务、单飞 Job；PG advisory lock 可推迟 Redis |
| RT-Q05 | 登录/API 限流跨实例 | 进程内滑动窗口 | **trigger-gated** | 单实例够用；**触发条件（H-002 · VR-052）**：多实例部署 **或** C 端业务域模块接入且 C 端限流需求不可共用进程内 limiter。业务域 VP 激活即视为触发条件成立，须评估进程内 limiter 是否满足 C 端场景并登记路线图位置。**承接 = [VP-027-rate-limiter-port](plans/VP-027-rate-limiter-port.md) `closed`**（2026-09-01 关门 · v0.3.0 · lead `workspace-027-rate-limiter-port` Root done 4/4；端口 + 内存默认 + 7 处使用点迁移 + 接缝声明已交付；Redis 实现仍 gated）。**C 端 ingress 评估（VP-030 · VRev-070 §6，2026-09-03）**：webhook/`chat_id`/`telegram_user_id` 桶可被进程内 limiter 覆盖，结论 = 不需要 Redis；VP-030 已于 2026-09-05 经 VRev-076 `pass` 关门，不消耗本行 trigger。**端口原子化（VP-032 `closed` v0.3.0 · 2026-09-04 · VRev-074 self `pass`）**：进程内 Allow/Record TOCTOU 由 `AllowRecord` + 令牌化 `Reserve`/`Cancel` 收口；仍不消耗本行 Redis trigger。**VP-031 激活复核（VRev-080 · 2026-09-05）**：VP-030 ingress 桶继续覆盖通道入口；购买/权益业务端点采用独立进程内请求计数桶，结论 = 本波不需要 Redis。R1 冻结 key/阈值/拒绝语义，不使用 key-wide `Clear` 抹除历史 |
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
| RT-D02 | 优雅停机 / 连接排空 | 停机顺序 / HTTP drain / Job 语义 / 双方言 Store 排空均已交付 | **delivered**（VP-021 `closed` v0.3.0 · 2026-08-27） | 2026-08-26 立项 → 2026-08-27 交付；单进程基线；与 Job 租约相关部分仍随 A3 |
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
| RT-K03 | JWT/数据密钥轮换合同 | VP-016 已交付 current+previous；签发只用 current；校验允许重叠窗；缺 previous = 今日单密钥 | **delivered** | 有界 residual：**立即失效未选**（I-016-005，无用户书面残余接受，按未选设计后果记录）；`admin.mfa` wrapping **已支持** previous 解密 + 成功 TOTP 后惰性重包（W11 F-004），残余收窄为「仅凭恢复码完成的路径不重包 + 无启动批量重包 + 无主动轮换重包」（用户 2026-09-10 接受）。KMS/HSM 仍 gated |
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
| RT-T03 | 时区在持久化层的合同 | locale 在 Admin 功能面；DB `timestamptz` 尚未交付 | **active** | Admin 面「时区/数字/货币」= [VP-020](plans/VP-020-timezone-number-currency-formatting.md)（`closed` v0.3.0 · 2026-08-27）；DB 合同 = [VP-040-timestamptz-persistence-contract](plans/VP-040-timestamptz-persistence-contract.md) `active`（workspace-040 Root `active · 0/3`；`I-040-001` collecting，R1 最终冻结未完成） |

### 9. 出站消息

| id | 项 | 现状 | 状态 | 备注 |
|----|----|------|------|------|
| RT-M01 | 出站邮件发送端口 + 可切换渠道 | VP-017 v0.5.0 按现行分母 closed：mock 站内记录 + Resend（live PASS）+ 设置热切换/试发 | **delivered** | 2026-08-24 再关门放行（A-003/A-004 pass）；R1～R4 实施史保留 |
| RT-M02 | SMS / 其它推送运输 | 无 | **trigger-gated** | 用户 2026-08-22：审核麻烦，有真实需求再做 |
| RT-M03 | C 端聊天通道（Telegram Bot 运行时） | VP-030 + VP-033 | **delivered** | [VP-030-telegram-channel-runtime](plans/VP-030-telegram-channel-runtime.md) `closed` v0.3.0（2026-09-05 · VRev-076 · workspace-030 Root done）；[VP-033](plans/VP-033-telegram-operator-console.md) `closed` v0.3.0（2026-09-05 · VRev-077 · workspace-033 Root done）。入站 webhook + 出站 SendMessage + 分发端口 + 单实例有界 `getUpdates` 与人工控制台均已交付；**不是**业务域。Mini App / Stars、多实例 / HA 长轮询仍 gated |

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

### 架构分支序列（2026-09-10 经 VP-035 R4 重述；已交付序列 + 唯一未触发项）

```text
A0  本清单（已登记）；Store 双方言决策已冻结（RT-P03）                      [done]
A1  内核持久化端口 + PostgreSQL 实现 + 现有台账对写/翻译；                  [delivered · VP-013 closed v0.3.0]
    SQLite 保留为 dev/mvp/快测默认；生产 CI 以 PG 为权威
A2  对象存储适配器；本地盘保留为默认                                        [delivered · VP-014 closed v0.3.0]
A3  仅当需要多实例：就绪探针扩依赖、再评估 PG 锁/SKIP LOCKED                 [trigger-gated · 唯一未触发项]
    vs Redis vs 外部队列（优雅停机/排空已拆为 A7）
A4  指标 + OpenTelemetry                                                    [delivered · VP-015 closed v0.3.0]
A5  密钥轮换 / 备份恢复合同                                                  [delivered · VP-016 closed v0.3.0]
A6  出站邮件：内核发送端口 + 可切换渠道（mock 默认 + Resend 生产）；         [delivered · VP-017 closed v0.5.0]
    SMTP 适配器保留；SMS 不进
A7  优雅停机 / 连接排空合同（RT-D02 → VP-021 closed v0.3.0，2026-08-27；      [delivered]
    单进程基线先行，不与 A3 绑定）
```

**C1（2026-09-10 VP-035 R3 登记 → 2026-09-19 立项 → 2026-09-20 激活）**：

```text
C1  DB 时间列 timestamptz 持久化合同（RES-T03-tz）
    — 现状：`apps/api` 时间列仍 INTEGER；`timestamptz` 命中 0
    — 承接 = [VP-040-timestamptz-persistence-contract](plans/VP-040-timestamptz-persistence-contract.md) `active`
      （lead = workspace-040-timestamptz-persistence-contract；R1 默认候选已登记，最终合同仍由 Root R1 冻结）
```

**刻意后置**：MongoDB、ORM、Redis、消息队列、搜索引擎、K8s、SMS。它们是部署或产品触发的后果，或已否决的技术选型。

架构分支最近一拍（实现中）：**[VP-040-timestamptz-persistence-contract](plans/VP-040-timestamptz-persistence-contract.md) `active` v0.2.0**（2026-09-20 激活 · VRev-104 self `pass` · lead `workspace-040-timestamptz-persistence-contract`；Root `GOAL-001-timestamptz-persistence-contract` 初始 `active · 0/3`）。前一拍 **VP-035-foundation-architecture-health `closed` v0.3.0**；A3 余项仍 trigger-gated。

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
- 消费链：**A6 出站邮件（VP-017 已按现行分母再 `closed` v0.5.0）→ 账号邮箱身份（VP-018 已 `closed` v1.0.0）→ IAM（密码策略 / 邀请 / 自助恢复状态机 = [VP-019-iam-recovery](plans/VP-019-iam-recovery.md) 已 `closed` v0.3.0，2026-08-26 用户确认关门）**。

**基架能力剩余**

1. 密码策略、邀请、账号恢复状态机（自助恢复 **硬前置** = VP-017 + 账号邮箱；邀请仍可用管理员出示链接）  
2. 组织 / 部门 / 岗位，以及数据权限 `org` 扩展 —— **trigger-gated**（2026-08-29 用户书面：自用优先下调权；触发 = 多组织/多团队 fork 消费或真实多组织管理需求出现，再经 `/vision` 立项；应用层 org 上下文仍归 Admin 分支，见架构 RT-P08 注）  
3. 配置包导出、diff、dry-run、导入 —— 已由 **[VP-025-config-export-diff-dryrun-import](plans/VP-025-config-export-diff-dryrun-import.md)** 承接并于 2026-08-30 **`closed`**（v0.3.0 · 用户书面确认 · roadmap 明文点名「其后非门控未立项 = 基架能力剩余 #3」）  
4. 文件扫描 / 隔离**策略**（执行器见架构 RT-S05）  
5. 时间、时区、数字、货币格式语义（持久化时区合同见架构 RT-T03）＝ **[VP-020-timezone-number-currency-formatting](plans/VP-020-timezone-number-currency-formatting.md) `closed` v0.3.0**（2026-08-26 激活 · lead `workspace-020`；2026-08-27 关门）  
6. impersonation / effective actor 产品化——仅当再次出现  
7. **账号邮箱身份**（`users` email 字段 + 校验状态机）——由 **VP-018** 承接并于 2026-08-24 **已 `closed`**（v1.0.0）；自助恢复的身份前置；运输面见架构 RT-M01（delivered）

**扩展接缝**（全部 trigger-gated；运输实现归架构）

typed domain event、Notification Transport、OIDC/SSO/SCIM、Approval Gate、Entitlement、多组织 context、SSE/WebSocket、外部连接器/Secret 的产品面、自定义 metadata/tags、文件预览。

> **SSE 注记（2026-09-03）**：VP-033 人工控制台首波用 Admin 短轮询 + 控制台 heartbeat，**不**解除本行 SSE/WebSocket trigger-gated。

> **Entitlement 注记（2026-09-02）**：VP-031 只交付**数字 Offer 本域**权益表与校验，**不**解除本行通用 Entitlement / Approval Gate 接缝的 trigger-gated。

**体验增强**

| 全局搜索 / Command Palette、Saved Views、批量结果中心、未保存保护、统一 Toast/错误恢复、版本与维护提示。全局搜索若需要专用引擎，拉动架构 RT-X01。**收口进度（2026-09-19）**：Command Palette = VP-036 `closed`；Saved Views / 未保存保护 / 统一反馈 = VP-037 `closed`；批量结果中心 = VP-038 `closed` v1.0.0；**版本与维护提示 = [VP-039](plans/VP-039-version-maintenance-diagnostics.md) `closed` v0.3.0**。体验增强清单无剩余未立项项。

**当前一拍（已交付并关门）为 [VP-038-batch-operations-and-job-center](plans/VP-038-batch-operations-and-job-center.md)**（2026-09-19 激活并**同日关门** · **`closed` v1.0.0** · 用户书面确认 · Root `GOAL-001-batch-operations-and-job-center` `done · 5/5` · 判据 1～7 达成 · R5 cross 审计 self `A-001` + grok-build independent `A-002` 均 `pass` · 补做关门 Vision Review = [VRev-100](reviews/VRev-100-vp038-batch-operations-and-job-center-closeout.md)）：把 VP-012 显式排除的「通用 Job 管理页」与本节未立项的「批量结果中心」收口为有界产品能力——已注册 Job 种类的可见性与结果读取 + **至少一条**真实批量操作以异步 Job 承接 + 结果中心体验；实体全文检索、Saved Views 重做、新业务域、组织/权限域与 Redis/MQ/多实例不进入本波。**计划 self = VRev-098 `pass` · 激活 self = VRev-099 `pass`（0 required）**；lead `workspace-038-batch-operations-and-job-center`（Root `GOAL-001-batch-operations-and-job-center`），2026-09-19 用户确认激活并由 `/govern` 开区。`I-038-004` 用户 P-004 裁决 = 新建 `admin.jobs` 进 admin 默认集（**Profile 内容扩展，不改装配语义，不暂挂 `go`**）；Admin 类 freshness `0c29c08` → `7e5ce891` **PASS**。

当前一拍（已交付并关门）为 **[VP-037-admin-workflow-continuity](plans/VP-037-admin-workflow-continuity.md)**：把已注册列表页的视图复用、离开前安全保护、成功/失败反馈与列表视觉/筛选体验收敛成一组可验证的工作流连续性能力；实体全文搜索、批量结果中心、组织/权限域与新业务域不进入本波。delivery workspace = `workspace-037-admin-workflow-continuity`；R1～R6 全部完成，Root `done · 6/6`，**2026-09-18 `closed` v1.7.0**（用户书面确认 · VRev-096 self `pass`）。

Admin 功能上一拍：**[VP-019-iam-recovery](plans/VP-019-iam-recovery.md)（IAM：密码策略 / 邀请入职 / 自助恢复状态机）——2026-08-25 激活并同日全链交付，2026-08-26 `closed` v0.3.0（用户书面确认；Root done 4/4；关后 A-001/A-002 pass）**；硬前置 = VP-018 已校验邮箱（已 `closed` v1.0.0）+ VP-017 运输（已按现行分母再 `closed` v0.5.0）。不要把恢复状态机打进 VP-018。再下一截（已交付并关门）：**[VP-020-timezone-number-currency-formatting](plans/VP-020-timezone-number-currency-formatting.md) `closed` v0.3.0**（2026-08-26 激活并开区 · 2026-08-27 关门 · 时区/数字/货币格式语义，基架能力剩余 #5 交付完成；lead `workspace-020-timezone-number-currency-formatting` 结项；关门审计双腿 pass）；其后非门控未立项 = 配置包导出/diff/dry-run/导入（基架能力剩余 #3 · **已由 [VP-025](plans/VP-025-config-export-diff-dryrun-import.md) 交付并 `closed`**）与体验增强（全局搜索 / Command Palette 等）；组织/部门/岗位 + 数据权限 `org`（#2）已于 2026-08-29 按用户指示降权为 **trigger-gated**（见「基架能力剩余」）。

Admin 功能最近一拍：**[VP-039-version-maintenance-diagnostics](plans/VP-039-version-maintenance-diagnostics.md) `closed` v0.3.0**（2026-09-19 激活并关门；用户书面确认；VRev-103 self `pass`；workspace-039 Root `done · 4/4`；R1–R4 全部 done；Goal cross required=0；existing fresh-seed harness bounded residual 保持登记）。其前一拍 = **[VP-038-batch-operations-and-job-center](plans/VP-038-batch-operations-and-job-center.md) `closed` v1.0.0**。当前 Admin 功能分支无 active 交付 VP；文件扫描/隔离策略与组织/部门/岗位仍未立项或 trigger-gated。

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

已有的钱包模块是 VP-011 交付的 Admin 演示/能力面，**不**等于本分支「支付/结算」业务域已成立。VP-029 的卡密入金仍属 Admin 功能（资金通道），不把本表第 3 项提前成立。

**已立项（2026-09-02 · 真实触发 = 下游 Telegram 付费服务 · 用户确认同进程）**

| VP | 收窄后的域 | 明确不做 |
|----|------------|----------|
| [VP-031-digital-offer-entitlement](plans/VP-031-digital-offer-entitlement.md) `closed` v0.3.4 | 数字 Offer + 薄购买凭证 + 本域权益 | 类目树、SKU/税/库存、物流订单、支付网关、通用 Entitlement 框架 |

业务域下一拍：无新触发则不要预开第二域。不要把候选 1～9 打进同一个 VP。

---

**当前组合焦点**：**active 交付 VP = [VP-040-timestamptz-persistence-contract](plans/VP-040-timestamptz-persistence-contract.md)**（架构 C1 · v0.2.0 · lead `workspace-040-timestamptz-persistence-contract` · Root `active · 0/3`；VRev-104 `pass`）。最近关门 = **[VP-039-version-maintenance-diagnostics](plans/VP-039-version-maintenance-diagnostics.md) `closed` v0.3.0**（Root `done · 4/4`；VRev-103 `pass`；用户书面确认）。持续程序 VP-009 / VP-010。组织/部门/岗位、实体全文检索、文件扫描与 A3/Redis/MQ 仍 gated。

> 2026-09-18 当前投影修订：VP-037 为 **`closed` v1.7.0**；Root 为 **`done · 6/6`**，R5 `GOAL-006` 为 `done · 4/4`，R6 `GOAL-007` 为 `done · 8/8`；关门依据 = 用户书面确认（`GOAL-006` `D-002`，前置条件为整改子目标 `GOAL-011` 修正两个分页/文案缺陷）+ `VRev-096` self `pass`。仍开放项见下节统一登记。

---

## 未决项统一登记（残余 / 悬置 / trigger-gated）

> **用途**：本区是"已交付范围之外的未决事项"的**统一登记处**——不是待办清单、不是承诺、也不代表已验证。目的只有一个：日后任何人对某个未实现或有界接受的能力有疑问时，能在这里一眼看到**它是什么、为什么不现在做、什么条件下做、谁负责、证据在哪**，而不必翻遍各工作区台账。
> **维护约定**：新增或闭合任何残余/悬置/触发项时**必须同步本节**（与 goal-tree、`03-audit` 台账同级要求）。登记只描述现状与触发条件，禁止把 deferred/recommended 写成已验证或已承诺。
> **最近更新**：2026-09-20（⑪ VP-040 `planned → active` v0.2.0 + workspace-040 / Root scaffold，VRev-104 `pass`，I-040-001 默认候选登记但 R1 仍 collecting；⑩ VP-039 `active → closed` v0.3.0 + workspace-039 Root `done · 4/4`，VRev-103 `pass`，用户书面确认；⑨ 激活开区；⑧ 立项。①–⑦ 见同日 VP-038 关门与残余收口记录）。

### 一、有界残余（B 类：实现已交付并验证，剩覆盖/文档加固）

| 编号 | 内容 | 现状 | 触发条件 | 责任人 | 证据 |
|------|------|------|----------|--------|------|
| `GOAL-009 A-001 F-001` | 列表视觉 e2e 未断言暗色下开关的**计算背景** | **2026-09-18 `fixed`**（暗色四条不变量 + 变异验证） | — | workspace-010 GOAL-043 | `GOAL-043 E-002` §1 |
| `GOAL-009 A-001 F-002` | 浏览器守卫仅覆盖单一页面 | **2026-09-18 `fixed`**（分页契约改为 roles + users 双页参数化） | — | workspace-010 GOAL-043 | `GOAL-043 E-002` §2 |
| `GOAL-005 A-002 F-002` / `R5-I-005` | Host 终态与普通 resource 反馈无直接对照断言 | **2026-09-18 `fixed`**（新增跨表对照测试；Host 文案表导出供读取） | — | workspace-010 GOAL-043 | `GOAL-043 E-002` §3 |
| `GOAL-008 A-002 F-002` | 全仓 `tsc` 简写未逐条裁定（354 行形态不可唯一确定） | **bounded residual**：可执行面由守卫 + CI 门禁锁死；文档侧不逐条考古 | **历史记录被再次当作类型检查证据引用时**，按 `GOAL-008 D-001` 口径复核并注明 | `/govern`（引用时） | `GOAL-043 E-003`；`GOAL-008 D-001` |
| `V-F124` | 首波页面/状态/权限/持久化矩阵（防执行期滑向共享视图/实体搜索/第二套基础设施） | **2026-09-18 `fixed`**：实质要求已由 R1 交付（`r1-denominator-matrix.json` 24/58、`r1-form-matrix.json`、`r1-state-feedback-matrix.md` + `D-003`～`D-005`），愿景层经 `VRev-097` 复核闭合 | — | `/vision` | `VRev-097`；`GOAL-002` R1 交付物 |
| `[workspace-038] GOAL-005 A-001 F-002/F-003/F-004` | R4 结果中心三条 low 级项：状态/错误文本未逐值本地化；自动刷新依赖 `reloadList()` 清空选择；无进行中作业时仍按档位轮询 | **2026-09-19 `fixed`**（用户裁决移交 → `[workspace-010]` `GOAL-044-w32-r4-residual-seams` 交付通用能力：列级 `valueLabels`、`refreshTable` 定向刷新 seam、`activeStatuses` 空闲判定；三条均经变异验证） | — | workspace-010 GOAL-044 | `GOAL-044` `D-001`/`E-001`/`A-001`；`[workspace-038] GOAL-005 A-003` §2 |
| 本地扩展登记（`valueLabels`/`badgeStyleField`/`truncate`/`width`/`minWidth` 等列级与节点级本地扩展缺少**单一清单**） | **登记（部分补齐）**：`[workspace-010] GOAL-045` 已把「列表页 actions 左侧插槽」加入本清单并记录其取值与约定；其余扩展仍分散在各波次决策里 | 出现「这是 pinned 还是本地扩展」的实际争议，或后续协议波次需要一次性核对时 | workspace-010（后续符合性波次） | `GOAL-044 A-001 F-003`；`GOAL-045 D-001` §1/§4；`GOAL-045 A-002` |
| `[workspace-038]` roles 页导出触发面（R3 冻结分母含 `users`/`roles`，但页面 UI 只接到 users） | **2026-09-19 `fixed`**：`roles-table` 增多选、页面增 `roles-batch-export` 节点（`resource: roles`，与 users 同一组件与同一插槽）；**后端零改动**（分母本就含 roles） | 已闭合；后续新资源接入导出时按同一形态扩展 | workspace-010 GOAL-045 | `[workspace-010] GOAL-045 E-001` §2；`apps/api/modules/roles/schema/roles.json` |
| `[workspace-038]` 「导出所选」在 operator config 下**必定 404**（`configs/config.yaml` 的 custom 内联列表自称镜像 admin preset，却漏了 VP-038 新加的 `admin.jobs` → 路由未挂载；而节点属于 `admin.users`/`admin.roles`，在所有 profile 照常渲染）；`mvp`/`demo` 为同类形态（既无 `jobs.write` 也无 `data.export`）；另发现 e2e harness 的 `customE2EModules` 同源漂移 | **2026-09-19 `fixed`**（用户报告 → 原样复现定位）：operator config 补 `admin.jobs`；触发面按路由真实的两道门禁渲染为**可见 + disabled + 说明**（隐藏会留下高度 0 的空插槽宿主、破坏 page-actions 高度契约）；提交命中 404 改报「未启用批量导出」而非裸 `NOT_FOUND`；新增 **config 超集守卫**（`internal/config` `TestOperatorConfigCoversAdminPreset`）+ **双 profile e2e 可用性契约**（`apps/web/e2e/batch-export-availability.spec.ts`，并同步修正 harness `customE2EModules`）；6 处变异验证（含双向判别性）。独立审计 `A-002`（grok build）`conditional` 指出「XOR 单测不具判别性 + e2e 未入仓」，经 `A-003` 响应两条 `fixed`，开放 required = 0 | — | workspace-010 GOAL-046 | `[workspace-010] GOAL-046 D-001/E-001/A-001/A-002/A-003` |

**本地扩展清单（截至 2026-09-19，供「pinned vs 本地扩展」判断）**：

| 扩展 | 层 | 取值/形态 | 约定 |
|------|----|-----------|------|
| `badgeStyleField` | 表格列 | 行字段名（`success`/`warning`/`destructive`/`info`/其它=neutral） | 单元格渲染为彩色徽标；空值显示 `—` |
| `truncate` / `width` / `minWidth` | 表格列 | boolean / px 或 CSS 长度 | 单行截断 + 标题提示 / 列宽提示 |
| `valueLabels` | 表格列 | `{ "值": "i18n 键" }` | 展示层映射；键缺失或值未映射 → **回落原始值**（fail-open） |
| `slot` + `targetTable` | custom 节点 | `slot: "list-page-actions"` | 节点投放进目标表格列表控件行**左段**；未知 slot 值按未声明处理（**原地渲染，无提示**）；目标表不存在亦原地渲染（fail-open） |
| `component` | custom 节点 | 组件注册键 | 节点级键（非 `props`）；未注册 → 明显占位 + `console.error` |
| `format: "tag"` / `tagMap` | 表格列 | — | **pinned 定义，本仓库未实现**；与 `valueLabels` 不竞争（前者字面量映射，后者 i18n 键） |
| `[workspace-038] GOAL-006 A-001/A-002 F-001` | e2e 挂具的 **fresh-seed 顺序契约**仍依赖文件名排序（`00-` 前缀）：全量套件共用一块 scratch 库，任何更靠前的新 spec 若使用 sign-in helper（其会自动完成强制改密）都会再次消费该前提，失败形态是难与真实缺陷区分的「登录 401」 | **bounded residual（2026-09-19 登记）**：本次已由重命名修复并两向复验（隔离通过 / 全量 16 passed / 4 skipped / 0 failed），但机制仍是隐式约定 | 后续波次新增 e2e 文件，或排序/挂具策略变更时 | workspace-010（后续符合性波次；e2e 挂具属 W23/W24/W25 一系） | `[workspace-038] GOAL-006 E-001` §2；`A-001/A-002` F-001；`apps/web/e2e/00-force-password-change.spec.ts` 注释 |
| `[workspace-038] GOAL-006 A-001/A-002 F-002` | 浏览器回归**未驱动** jobs 结果中心：11 个 e2e spec 无「提交批量导出 → 观察进度 → 下载」路径；且 `npm run test:e2e` 默认 `APP_PROFILE=mvp`（**不含** `admin.jobs`），本 VP 新增面在默认 e2e 中模块级缺席 | **2026-09-19 `fixed`**（用户指令「先补 jobs e2e 再关门」）：新增 `apps/web/e2e/jobs-result-center.spec.ts`（admin profile 专用，mvp 下 `test.skip` 显式跳过）——真实浏览器路径：选择行 → 提交异步导出 → 观察真实进度 → 终态下载（文件名 = 服务端 `users-selection.csv`）→ 结果中心读同一作业（本地化 `Succeeded`）→ 行操作再次下载；双 profile 全量复跑：mvp 16 passed / 5 skipped / 0 failed，admin 17 passed / 4 skipped / 0 failed（均 exit 0） | 未覆盖部分（重试/取消的浏览器端到端路径）随实现扩展时补 | `/govern`（后续波次） | `[workspace-038] GOAL-006 E-001` §2/§3；`A-004`；`apps/web/e2e/jobs-result-center.spec.ts` |

### 二、悬置的范围决策（C 类：非缺陷，等需求再定）

| 编号 | 内容 | 现状 | 触发条件 | 责任人 | 证据 |
|------|------|------|----------|--------|------|
| `I-037-005` | 跨用户共享视图、最近使用/收藏、协作权限是否进入后续波次 | 首波明确只做**个人级**工作流；`deferred · non-blocking` | **真实协作需求出现** | `/vision` | VP-037 首波范围表；Root `GOAL-001-admin-workflow-continuity` 信息表 |

### 三、未推进 / trigger-gated 能力（A 类：从未进入任何 VP 的实现分母）

| 能力 | 现状 | 触发条件 | 责任人 / 下一步 | 出处 |
|------|------|----------|-----------------|------|
| 实体全文检索（`RT-X01` 专用引擎 / `RT-X02` DB 全文检索） | 架构触发项，**未实现**；VP-036 首波只做已注册页面/导航/声明式动作检索；VP-038 只做作业/批量结果的可见性，**不含**实体检索 | 真实**实体级**搜索需求 + 规模证据 | `/vision` 立新 VP；引擎路线拉动架构 `RT-X01` | roadmap「体验增强」；VP-036 边界表；VP-038 边界表 |
| 批量结果中心 | **已交付并关门**：[VP-038-batch-operations-and-job-center](plans/VP-038-batch-operations-and-job-center.md)（**`closed` v1.0.0** · 2026-09-19 用户书面确认 · lead `workspace-038-batch-operations-and-job-center`）；R1～R5 子目标全部 `done · 4/4`（Root `GOAL-001-batch-operations-and-job-center` **`done · 5/5`**）；同清单的 Saved Views / 未保存保护 / 统一反馈已由 VP-037 交付 | 已交付；关门后残余登记于本节 | `/vision`（后续波次） | VP-038；roadmap「体验增强」 |
| 组织·部门·岗位 + 数据权限 `org` | 基架能力剩余 #2；**2026-08-29 用户书面降权**为 trigger-gated | 多组织/多团队 fork 消费，或真实多组织管理需求 | `/vision`（应用层 org 上下文归 Admin 分支） | roadmap「基架能力剩余」#2 |
| 新业务域 | **未立项**；无新触发不预开第二域 | 真实业务需求 | `/vision` | roadmap「业务域」 |
| Redis / MQ / 多实例（+ 第二持久化栈） | 架构 **A3** / `RT-Q03`，架构骨架**唯一未触发项** | 多实例部署，或 C 端业务域模块正式接入同进程 | 架构分支 + `/vision` | roadmap「基架能力剩余」；`architecture/cache-redis-seam-and-track.md` |
| 文件扫描 / 隔离**策略**（执行器见 `RT-S05`） | 基架能力剩余 #4，**未立项** | 真实需求 | `/vision` | roadmap「基架能力剩余」#4 |
| 版本与维护提示 | **已关门**：[VP-039-version-maintenance-diagnostics](plans/VP-039-version-maintenance-diagnostics.md) `closed` v0.3.0 · lead `workspace-039-version-maintenance-diagnostics` · VRev-103 `pass` · Root `done · 4/4` | 已交付 | `/vision` 关门记录 | roadmap「体验增强」；VP-039 |
| DB `timestamptz` 持久化合同（C1 / `RT-T03` / `RES-T03-tz`） | **实现中**：[VP-040-timestamptz-persistence-contract](plans/VP-040-timestamptz-persistence-contract.md) `active` v0.2.0 · `workspace-040-timestamptz-persistence-contract` / Root `active · 0/3`；`I-040-001` collecting（默认候选已登记，最终 R1 合同未冻结） | R1 合同/分母冻结与后续迁移门禁 | `/govern`（R1 先冻结） | roadmap 架构 C1；VP-040；VRev-104 |
| 扩展接缝：typed domain event、Notification Transport、OIDC/SSO/SCIM、Approval Gate、Entitlement、多组织 context、SSE/WebSocket、外部连接器/Secret 产品面、自定义 metadata/tags、文件预览 | 全部 **trigger-gated**（SSE 注记：VP-033 用短轮询，**不**解除本行；Entitlement 注记：VP-031 只交付数字 Offer 本域，**不**解除本行） | 各自真实需求 | `/vision` | roadmap「扩展接缝」 |

## 单主线模块化策略

未来 fork 起点统一由同一代码主线、模块候选集与启动时 Profile 表达，权威见 [module-architecture.md](../architecture/module-architecture.md) 和 VP-003。原 [dual-track-contract.md](dual-track-contract.md) 已转为历史记录。
