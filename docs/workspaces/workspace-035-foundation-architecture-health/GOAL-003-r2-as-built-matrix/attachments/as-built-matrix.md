---
title: R2 as-built 对照矩阵
status: active
created: 2026-09-09
updated: 2026-09-09
parent: GOAL-003-r2-as-built-matrix
version: 0.1.0
---

# R2 as-built 对照矩阵

基线 `ebe6013c`；与 R1 消费候选 `5c341ec7` 的差异为 VP-035 治理文档。按 R1 包含表 17 面 + Web 附加抽检逐行取证。路径均相对仓库根；行号指基线代码。本文是有界架构对照，不是全仓安全、生产服务或协议全量认证。测试执行结果另见 validation.md；列出测试名不等于实跑通过。

## 分母对照

| # | 面 / 设计预期 | 当前代码事实与证据 | 代表验证 / 差异 |
|---|---|---|---|
| 01 | 模块契约；薄内核、六项、fail closed | `apps/api/kernel/module.go:11,265` 固定 API 2.0.0 / Resolve 校验依赖及冲突；`kernel/provider.go:19,27,70` Provider + Registrar，校验后整体发布；`kernel/persistence.go:47` CollectPersistence 独立全局入口 | `TestDualProfileContractMatrix`、`TestCollectPersistenceVersionGap`。Persistence 不在 Registrar：六项不能等同于六个 Register 调用；users 的空迁移归属 core.auth-session，`modules/users/provider.go:65` |
| 02 | Profile 默认集 / compiled 不自动启用 | `apps/api/kernel/profile.go:25` mvp 11、admin 23、demo 12 项；custom 须显式集；`internal/composition/composition.go:83` ResolvePlan 先显式再 Profile，依赖不补齐；`internal/config/config.default.yaml` 提供 app 配置入口 | `TestBuiltinProfilesResolveDeterministically`、`TestS2AccessDrill_ProbeAbsentFromDefaultProfiles`；Telegram/digital-offer 不在默认集 |
| 03 | 静态组合、公共面边界 | `internal/composition/composition.go:127,509,677` Fx 在组合根，按 plan 静态装配 Provider；`internal/server/server.go:17` HTTP 超时及中间件；`assembly/assembly.go:26` OpenStore 返回 kernel.Store | `TestAppStartsAndStopsDualProfiles`。B+ `assembly.NewAuthenticator` 在 `assembly.go:37` 返回 internal auth 指针且仅 current 单密钥，不能宣称所有公开工厂都已返回 kernel 接口；登记 G-004 待分类，不据此改代码 |
| 04 | Store 双方言 / 无驱动句柄 | `kernel/store.go:30,43` Store/Tx/Rows 接口；`internal/store/store.go:141` Run callback、rollback/repanic、嵌套拒绝；`internal/composition/composition.go:198` 全候选 catalog 先于 open；`internal/store/postgres.go` PG 实现 | `TestKernelStoreRunCommitAndRollback`、`TestMigrateMidApplyFailureNoResidueThenReopen`。`kernel/store.go:24` ErrNoRows 是标准库 sentinel 别名，非 *sql.Tx 或供应商句柄泄漏。SQLite 文件池默认 4（store.go:29,104），内存池 1；G-002 文档锚点过期 |
| 05 | ObjectStore 本地/S3、公共面无 os.File | `kernel/objectstore.go:30` namespace/object ID 验证；接口使用标准 io 类型；`internal/objectstore/local.go:119` 文件/元数据写入；`s3.go:47,117` 供应商适配及 not-found 映射 | `TestLocalPutRollsBackWhenMetaWriteFails`、`TestS3TransportErrorsPropagate`；本地与 fake S3 证据不等于 live S3 验收 |
| 06 | Cache 内存默认 / Redis 接缝 | `kernel/cache.go:35,117` namespace/key/expiry 约束；`internal/cache/memory.go:64,227` 有界内存 + 拷贝保持空 slice hit；composition.go:853 构造内存 | `TestMemoryFailClosedValidation`、`TestMemoryConcurrentBudgetBound`；go.mod 未引入 Redis；接缝仍未实现 |
| 07 | RateLimiter 内存、原子预留/取消 | `kernel/ratelimit.go:35` 同时保留 Allow/Record/Clear 和 AllowRecord/Reserve/Cancel；`internal/ratelimit/memory.go:83` 窗口/容量，Cancel 仅取消本次 token | `TestMemoryReserveCancelsOnlyItsOwnSlot`、`TestMemoryReserveConcurrentBudget`；G-003：Redis 接缝 §2.6 尚只描述旧 Allow/Record/Clear，不足以指导现行替换 |
| 08 | EventBus 进程内运输 | `kernel/eventbus.go:43` 类型化 topic 契约；`internal/eventbus/memory.go:72,254` 校验、背压、panic 隔离；composition.go:871 单例内存 | `TestPublishBlocksOnFullBuffer`、`TestStopDrainsAndRejects`、`TestPayloadIsolation`；不提供事务持久化、跨进程可靠消息，不据缺席判缺陷 |
| 09 | Mail 端口/渠道隔离 | `kernel/mail.go:33,73` 消息/发送接口；`internal/mail/runtime.go:105,363` 渠道切换、失败保留旧渠道；composition.go:1024 配置 mock retention；`assembly.go:43` 默认 OutboxSink | `TestSwitcherHotSwitchSemantics`、`TestSMTPSendFailsClosedWithoutAuthAdvertisement`；mock/Resend/SMTP 与邮件产品分离；CaptureSink 是历史实现，不能当现行默认 |
| 10 | Observability 缺省无外部依赖 | `internal/config/config.go:493` metrics/traces 默认关闭；`internal/obs/tracing.go:48` noop/显式 exporter；`server.go:37` 专用 listener；composition.go:262,289 接线 | `TestTracingDisabledIsNoOp`、`TestTracingExporterDeliversToOTLPSink`、`TestServerEnabledBindFailureFailsClosed`。仅当前指标分母；不扩大到 Store/对象/Job；exporter 初始化失败退 noop 属可观测策略，不应笼统宣称一切失败都阻止启动 |
| 11 | Shutdown 顺序 / Job 中断重跑 | `kernel/lifecycle.go:18,67` 拓扑启动/逆序停止及失败清理；composition.go:1151 HTTP drain → retention/metrics → jobs → eventbus → runtime → store → tracing，错误合并 | `TestShutdownDrainHarness`、`TestShutdownInterruptLeaseReclaim`；HTTP/Job 单测可在 Windows 执行，不能把所有 drain 测试统称 !windows；进程/Compose Linux 证据本轮未重跑 |
| 12 | Job 六态、区别 EventBus | `internal/jobs/model.go:19` queued/running/succeeded/failed/cancelled/expired；runner.go:125,185,214 Stop 取消并等待 worker，扫描恢复/过期 | `internal/jobs/shutdown_reclaim_test.go`；持久任务状态/lease 与内存 EventBus 是不同职责；不等于外部 broker |
| 13 | Telegram 端口无 SDK | `kernel/telegram.go:45,76` 本仓消息/分发契约；`internal/channel/telegram/disabled.go:9` sender 返回 disabled，dispatcher 注册 no-op；runtime.go:90 参数校验 | `TestTelegramFxShutdownDrainsPollingReceiver`；composition.go:601,637 启用计划与运行时共享。disabled 注册 no-op 是显式禁用形态，不将其自动判为故障 |
| 14 | JWT current/previous | composition.go:255 注入 current + previous；`internal/auth/auth.go:424` current 验证失败再试 previous，保留过期与方法校验 | `TestDualKeyRotationOverlapWindow`、`TestNewAuthenticatorWiresPreviousSecret`；B+ factory 仅 current 的限制另见 G-004；未重做 dump |
| 15 | Manifest 聚合与 NavGroup | `internal/manifest/manifest.go:40,213,457` 结构化贡献投影/分组/聚合；composition.go:677,771 聚合已启用贡献；group 可选、冲突拒绝、top/user 不参与 sidebar 归一化 | `TestAggregateRejectsProtocolAppPageAndNavigationConflicts`、`TestNormalizeSidebarGroupsRejectsMetadataConflictAndNonSidebarNode`、`TestR4NavigationGroupProfileMatrix` |
| 16 | 六份 architecture 文档 | module-architecture §1–7 主线与当前实现一致；playbook §1/1.3/2 覆盖 compiled-global 与 group；cache seam §2.6 落后于 RateLimiter；overview 描述旧 FastAPI 治理 UI；directory-layout 指向 apps；monorepo-layout 为初建阶段记录 | G-001/G-002/G-003：现时描述滞后。module-architecture `vision_ref @0.2.0` 是旧架构记录，不用于覆盖当前 VP→Charter @0.4.0 链；monorepo 初建章节保留历史语境 |
| 17 | Provider/handler 供应商泄漏抽检 | 搜索 `apps/api/modules` 的 *sql.Tx / redis.Client / s3.Client / pgx / tgbot / aws.：*sql.Tx 命中为测试/历史注释，未发现目标供应商客户端生产签名；digitaloffer/provider.go:25,68 注入 kernel.TelegramDispatcher/Sender | 有界文本抽检，不是所有依赖的类型证明；标准 []byte、io、JSON 负载不属于供应商类型泄漏；assembly B+ 单独列 G-004，不用抽检无命中掩盖 |
| W1 | Web Shell/Host 只消费 Manifest | `apps/web/src/main.tsx:76` bootHost + manifest loader；`app/App.tsx:1049` projectNavigation；`app/navigation.ts:224` 以 manifest 输入投影；`protocol/load-page.ts:81` schemaUrl 加载 | 本轮 Web 4 files / 48 tests passed（navigation/app-manifest/load-page/App.integration）；未发现中央业务导航表。自定义渲染组件注册不等于中央业务导航注册 |

## 缺口候选（R3 用户分类前不作整改授权）

| ID | 证据与影响 | 建议候选 / 门禁 |
|---|---|---|
| G-001 | overview 当前阶段仍 Charter @0.2.0 / VP-005 active，且产品树写 FastAPI web；实际 apps/api + apps/web | R4 文档卫生；不能用于判断产品尚未实现 |
| G-002 | roadmap 架构锚点/RT-P04 写 MaxOpenConns=1，RT-D02 同行称无明确 drain 合同却标 delivered；代码已池化且有 drain | R4 editorial 草案逐项修正；读写分离/replica 仍 gated，不能把已有池化扩写成这些已交付 |
| G-003 | cache-redis-seam-and-track.md §2.6 旧端口描述缺 AllowRecord/Reserve/Cancel | 建议现在修文档并交 owner 确认边界；或明确登记到触发前必修。不得替 owner 冻结新的 Redis 实现，也不接受语义缺口为残余 |
| G-004 | assembly.NewAuthenticator 返回 *internal/auth.Authenticator 且只有 current；B+ 注释明确类型推断和单密钥形态 | 当前是可核对的公共工厂限制，非 kernel 供应商句柄泄漏。是否构成优先整改、另立兼容性波次，留 R3/用户；没有包消费者失败证据，不宣称生产故障 |

R1 residual 总账仍保留 12 条入册/3 条排除；R2 不改变原 residual 接受范围，也不重开原 VP。R3 须逐条建议分类与复审条件，不能靠本阶段矩阵自动接受残余。
