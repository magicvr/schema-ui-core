---
doc_type: goal-attachment
id: r3-industry-comparison
parent: GOAL-004-r3-industry-comparison
status: draft
created: 2026-09-10
updated: 2026-09-10
version: 0.1.0
---

# R3 有界业界对照表（C2）

本表是 VP-035 方向级判据 3 的证据主体：四类冻死参照**每类至少一条**对照行，每行四格齐备（业界常见做法 → 本仓现状 → 分类 → 不推翻项），分类另附「影响的路线图行」与「复审触发」。

- **业界侧输入**：[industry-conventions-draft.md](industry-conventions-draft.md)（13 行素材，19 个来源经 `web_fetch` 实测 200；见 [industry-conventions-sources.txt](industry-conventions-sources.txt)）
- **本仓侧证据**：R2 [矩阵 v0.2.0](../../GOAL-003-r2-as-built-matrix/attachments/as-built-matrix.md) + 本轮亲自复核的精确锚点（下表内标 `复核` 者为本次逐一读过源码行）
- **分类词表**：仅 `现在修` / `仍 gated` / `接受残余` / `明确不做`（R3 D-001 冻结项 5）；「分类」列**只允许**四值本身，去向与路线图行写在同格方括号「〔去向〕」内，不得与分类值混写。
- **边界**：本表不改任何生产代码、端口语义、Profile 默认集；不消耗 trigger-gated 行；不写 `docs/vision/roadmap.md` 正文

## A · 四类参照对照（13 行）

| # | 类别 | 业界常见做法 | 本仓现状（精确锚点） | 分类 | 不推翻项（本行不得用于） |
|---|------|--------------|----------------------|------|--------------------------|
| 1.1 | 模块化单体 + 组合根 | 模块化单体把「逻辑模块边界」与「单一部署单元」分开：模块只经 API 对外、内部包禁止被跨模块引用，整体仍编译部署为一个单元（Spring Modulith `verify()`；Microsoft 单体架构） | 薄内核 + 编译期候选集 + 显式依赖校验，单部署单元：`apps/api/kernel/module.go:11`（`KernelAPIVersion = "2.0.0"`）、`:265`（`Resolve` 入口，拒绝空集/重复/未知）、`:291`–`298`（依赖逐项 fail closed：未编译 `:293`–`295`、未启用 `:296`–`298`）；测试 `apps/api/kernel/kernel_test.go:89`（`TestRegistryRejectsUnknownMissingCycleConflictAndCapability`）〔复核〕 | **明确不做**（保持现行形态） | 不得据本行认定模块划分已达标，也不得作为拆服务或引入运行时插件市场的依据 |
| 1.2 | 模块化单体 + 组合根 | 依赖装配集中在组合根完成：入口把接口与实现接线，其余代码只声明依赖（Microsoft「composition root」；.NET Generic Host；Fx modules） | 唯一 Fx 装配在组合根：`apps/api/internal/composition/composition.go:111`（`NewApp` 为 Fx 组合根注释）、`:146`（`newCache` 进入 Fx 图）、`:312`–`314`（`cachePort kernel.Cache` / `eventBusPort kernel.EventBus` / `rateLimiters kernel.RateLimiterProvider` 由容器注入）、`:677`–`678`（`RegisterContributions` 在装配末端统一发布）；`assembly/assembly.go:26,37` 为公开工厂（B+ 边界，另见 G-004）〔复核〕 | **明确不做**（保持现状）→ 路线图行：无 | 不得据此要求替换 Fx、改成运行时反射发现或外部插件装载 |
| 1.3 | 模块化单体 + 组合根 | 装配图在启动阶段显式校验，缺依赖/冲突/作用域不匹配即失败退出（fx `ValidateApp`；Generic Host development 校验） | 校验先于发布：`kernel/module.go:265`–`300` 逐项返回 `CodeModule*` 错误；`apps/api/kernel/provider.go:70`（`RegisterContributions`）在 `:678` 调用点整体发布；测试 `kernel_test.go:89`、`:138`（`TestRegistryValidatesKernelAPIRanges`）〔复核〕 | **明确不做**（已具备，无需新增校验框架）→ 路线图行：无 | 不得反推本仓校验已完备到覆盖所有装配错误，也不得据此改动内核契约或错误语义 |
| 2.1 | 基础设施端口 | 外部依赖以端口接口表达，内存/回环适配器是常规实现之一，用于脱离真实基础设施运行与测试（Cockburn 端口与适配器 2005 原文） | 端口与内存默认实现成对存在：Store `kernel/store.go:30`（`Store`）/`internal/store/store.go:141`；ObjectStore `kernel/objectstore.go:108`/`internal/objectstore/local.go:119`；Cache `kernel/cache.go:35`/`internal/cache/memory.go:64`；RateLimiter `kernel/ratelimit.go:36`/`internal/ratelimit/memory.go:83`；EventBus `kernel/eventbus.go:43`/`internal/eventbus/memory.go:72`；Mail `kernel/mail.go:73`/`internal/mail/runtime.go:363`〔复核〕 | **明确不做**（形态即常规落法）→ 路线图行：本形态支撑 RT-\* 行的 delivered/gated 判定 | 不得据此认定每个端口都必须有内存实现；也不得把「内存实现是惯例」当作阻止未来真实适配器另立评估的理由 |
| 2.2 | 基础设施端口 | 外部后端按「可挂载资源」处理：只有配置句柄不同，替换本地/第三方实现不改代码（12-Factor IV） | 端口只暴露标准库/自有类型，供应商类型留在 `internal/`：`kernel/objectstore.go:109`–`114`（`[]byte`/`io` 语义）、`kernel/ratelimit.go:86`–`90`（`RateLimiterProvider` 工厂）；`apps/api/go.mod:5`–`22` 无 Redis/MQ/K8s/ORM 依赖〔复核〕 | **仍 gated**（Redis/MQ 触发条件未满足）→ 路线图行：RT-Q03/RT-Q05、A3 | 不得据此认为现在就该接入 Redis/MQ/K8s；也不得把「可挂载资源」当作外部适配器已实现或已验证的证据 |
| 2.3 | 基础设施端口 | 私有进程内缓存优先；需要跨实例一致视图时才引入共享缓存，且共享缓存不得作为权威存储，不可用时应能回退（Microsoft Caching guidance） | 内存缓存为默认且有界：`internal/cache/memory.go:64`（构造函数）、`:227` 区段（容量与拷贝语义）；组合根按配置预算构造单一实例：`composition.go:854`（`func newCache(cfg *config.Config) (kernel.Cache, error)`）、`:862`（`return cache.NewMemory(budget)`）〔复核〕 | **仍 gated**（共享缓存/多实例未触发）→ 路线图行：RT-Q03；`docs/architecture/cache-redis-seam-and-track.md` | 不得据此要求立刻实现 Redis 适配器；也不得把进程内缓存的多实例不一致语义直接判为缺陷 |
| 2.4 | 基础设施端口 | Go 惯例：接口定义在**消费方**包、不为 mock 在实现方定义、不在出现真实用例前预先定义（Go Code Review Comments · Interfaces） | 端口定义在内核且**已被消费方使用**：`composition.go:312`–`314` 注入 `kernel.Cache`/`EventBus`/`RateLimiterProvider`；`kernel/ratelimit.go:82`–`90` 的 provider 工厂注释写明「内存与未来 Redis 同合同、消费方零改动」〔复核〕 | **明确不做**（既有接缝不判违规）→ 路线图行：RT-Q03/Q05 维持 gated | 不得用本惯例主张删除既有内核端口、把已冻结接缝判为「提前抽象」违规，或据此重开端口公开语义 |
| 3.1 | Schema 驱动 Admin | 用声明式 schema（如 JSON Schema）描述结构与约束，由通用渲染器生成表单/页面，schema 兼作校验与文档，外观差异交给额外 UI 描述（json-schema.org；react-jsonschema-form） | 页面协议经 schemaUrl 加载并由宿主渲染：`apps/web/src/protocol/load-page.ts:81`（schemaUrl 解析）；`apps/web/src/main.tsx:76,79`（bootHost + manifestLoader）〔复核；自定义渲染组件注册与 schema 页面并存〕 | **明确不做**〔形态已成立；未实现部分不写成缺口〕 | 不得扩大解释为必须实现通用低代码/表单设计器，也不得要求把所有手写页面改造成 schema 渲染 |
| 3.2 | Schema 驱动 Admin | 页面与导航由各模块/插件以扩展点贡献、宿主自动发现并聚合，应用侧不做中心手写注册表（Backstage extension blueprints：`Nav items are auto-discovered from page extensions`；Django admin `autodiscover()`） | 贡献 → 投影 → 归一化链路：`apps/api/internal/manifest/manifest.go:43`（`NavigationPresentationsFromContributions` 投影）、`:213`（`NormalizeSidebarGroups`）、`:457`（`Aggregate`）；`composition.go:678`（聚合已启用贡献）；Web 侧只按 manifest 投影：`apps/web/src/app/navigation.ts:224`（`projectNavigation(manifest, ...)`）、`apps/web/src/app/App.tsx:1048`–`1050`（useMemo 调用点）〔复核〕 | **明确不做**（无中央业务导航表；形态即贡献式）→ 路线图行：无 | 不得据此要求引入第三方/运行时插件市场或动态装载；也不得把「应用侧保留自定义渲染组件注册」判为偏离惯例 |
| 3.3 | Schema 驱动 Admin | 权限由功能贡献方声明（可授权资源与动作），授权策略由部署/集成方集中实现，核心不内置具体授权逻辑（Backstage Permissions Overview：plugin authors 声明、integrators 配置 policy） | 权限**确实**随贡献声明、由宿主集中发布与校验：`apps/api/kernel/contribution.go:49`（`PermissionContribution` 含 `Permission`/`Resource`/`Action`/`PolicyID`）；注册面 `apps/api/kernel/provider.go:31`（`Authorization(...)`），`:242` 落实；贡献集持有 `provider.go:43`（`Permissions []PermissionContribution`）；导航节点可挂权限并**要求该权限已被声明**：`contribution.go:87`（`NavigationContribution.Permission`）、`provider.go:316`–`321`（权限索引）、`:349`（`if _, ok := permissions[n.Permission]; !ok` → 拒绝）；模块侧示例 `apps/api/modules/users/provider.go:107`–`111`（`users.enable`/`users.disable` 带 `PolicyID: authsessiondata.PolicyAdmin`）〔复核〕 | **明确不做**〔已具备：贡献声明 + 宿主集中判定，与业界形态一致；本 VP 不扩权限模型〕 | 不得用本行主张在本 VP 内新做权限模型、RBAC 或授权服务；也不得把「贡献声明权限」扩写成细粒度 ABAC 或列级权限已交付 |
| 4.1 | 同进程基座 | 默认从单体起步，只有复杂度（规模、团队、独立演进/伸缩）超出单体管理能力时才拆服务，分布式成本即 microservice premium（Fowler · MonolithFirst / MicroservicePremium） | 单进程基座：`apps/api/kernel/lifecycle.go:18`（`Start` 按 plan 拓扑启动）、`:68`（`stopModules` 逆序）；`composition.go:1151`–`1170`（停机顺序与错误合并）；`composition.go:1149`–`1171` 为单进程生命周期钩子〔复核〕 | **明确不做**（保持单进程默认）→ 路线图行：架构主线与 H-002 同进程基座 | 不得据此宣布永不拆分，也不得当作拒绝未来任何服务化评估的永久结论 |
| 4.2 | 同进程基座 | 同一部署单元内运行多个功能模块，共享宿主提供的基础服务并整体部署与伸缩（Microsoft 单体容器化；Backstage backend：模块与所属 plugin 同实例、共享 services） | 多模块共享单例端口与生命周期：`composition.go:111` 区段（Fx `Provide` 装配）、`:312`–`314`（共享端口注入）、`:430`（EventBus 单例注释：`kernel.EventBus` 由 Fx 持有进程生命周期）；模块清单见 `kernel/profile.go:25`（mvp/admin/demo 默认集）〔复核〕 | **明确不做**（形态即常规基座）→ 路线图行：RT-M03、`biz.*` 后续波次 | 不得把「同进程」绝对化为不可变更的架构承诺，也不得据此否决未来按模块边界做拆分评估 |
| 4.3 | 同进程基座 | 模块间主要交互走进程内事件发布/订阅；只有需要跨进程或对接外部系统时才外部化，且外部化需显式标注并引入 broker 依赖（Spring Modulith events） | 进程内 EventBus 为默认运输：`kernel/eventbus.go:43`（类型化 topic 契约）、`internal/eventbus/memory.go:72`（注册校验）、`:254`（订阅者运行/panic 隔离）；Job 六态为独立职责：`internal/jobs/model.go:19`〔复核〕 | **仍 gated**（outbox/MQ 未触发）→ 路线图行：RT-Q05、RES-028-broker | 不得据此把进程内 EventBus 扩写成事务持久化/可靠消息能力，也不得用本行否决未来按触发条件引入 broker 的评估 |

## A2 · 本轮取证中修正/新增的核对事实

本节记录 C2 期间**由证据推翻的先前表述**，以及 R2 矩阵之外新核对的锚点。不修改 R2 原始记录（A-001/A-002/A-003 与矩阵 v0.2.0 保持原样）。

| # | 项 | 先前表述 | 复核结论 | 影响 |
|---|----|----------|----------|------|
| 1 | 行 3.3 权限面 | 本轮初稿曾按 R2 矩阵行 15/W1 推断「未见权限随贡献声明的对应面」 | **推断不成立**：`kernel/contribution.go:49` 的 `PermissionContribution` + `provider.go:31/242` 注册面 + `provider.go:349` 的「导航权限必须已被声明」校验，构成完整的「贡献声明 + 宿主集中发布」链 | 行 3.3 分类由「登记为显式边界」改为**明确不做（已具备）**；此项属本轮自查修正，不涉及 R2 结论（R2 矩阵行 15 只声明「不参与 sidebar 归一化」，未主张权限面缺席） |
| 2 | 注册面数量 | R2 矩阵行 01 记「Registrar 六项不能等同于六个 `Register` 调用」 | 复核确认：`provider.go:19`–`35` 接口含 `Authorization`/`Navigation`/`Manifest`/`Configuration` 等分面方法，`:242`/`:256`/`:274`/`:288` 为各自落地 | 与 R2 结论一致，无需修改 |
| 3 | G-003 端锚点 | R2/本轮初稿写 `kernel/ratelimit.go:36` 起 | 精确锚点：`:40 Allow`、`:44 Record`、`:52 AllowRecord`、`:65 Reserve`、`:71 Cancel`、`:76 RetryAfterSeconds`、`:79 Clear` | 已同步到 [r4-doc-hygiene-anchors.md](r4-doc-hygiene-anchors.md)（G-003 表） |
| 4 | gated 依赖缺席 | R2 矩阵行 06 记「go.mod 未引入 Redis」 | 复核 `apps/api/go.mod:5`–`22`：直接依赖含 fx/pgx/sqlite/aws-sdk-s3/prometheus/otel，**无** Redis、broker、K8s、ORM | 行 2.2 的「仍 gated」分类有可核对证据 |

## B · 行级备注（R2 缺口的落点）

| 行 | 与 R2 缺口候选的关系 |
|----|---------------------|
| 2.3 / 2.4 | 端口形态本身判「明确不做」；相关**文档滞后**（`cache-redis-seam-and-track.md` §2.6 未覆盖 `AllowRecord`/`Reserve`/`Cancel`）是单独缺口 G-003，落 R4 文档卫生或另立，不改变本行分类 |
| 3.3 | 与 G-004 无直接关系；G-004（`assembly.NewAuthenticator` 公共工厂限制）另见 C3 分类表 |
| 全部 | 本表不把 R1 §1.2 排除项（Admin 体验增强、业务域产品面、gated 缺席、分发残余）写成缺口 |

## C · 判据 3 对照

| 判据 3 要求 | 本表证据 |
|-------------|----------|
| 四类参照每类至少一条对照行 | 类别 1：3 行（1.1–1.3）；类别 2：4 行（2.1–2.4）；类别 3：3 行（3.1–3.3）；类别 4：3 行（4.1–4.3） |
| 每行含四格（业界常见 → 本仓现状 → 分类 → 不推翻项） | 全部 13 行齐备；分类均取自四值词表 |
| 对照未导致静默改 Charter | 全部行未产生「必须改 Charter 非目标」的结论；裁决留 I-035-003 逐行判定（C4） |
| 不把业界实践当决策源 | 每行「不推翻项」逐行写明禁用方向；未使用个人博客/营销页，来源 19 条全部实测 |
