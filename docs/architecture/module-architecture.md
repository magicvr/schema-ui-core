---
doc_type: architecture-decision
title: 单主线模块化 Admin 架构
status: active
created: 2026-08-04
updated: 2026-08-06
parent: null
version: 1.0.3
vision_ref: schema-ui-core-admin-foundation@0.2.0
serves: VP-003-modular-admin-architecture
---

# 单主线模块化 Admin 架构

## 决策

`schema-ui-core` 的目标架构是**单一代码主线上的模块化单体**。同一份可执行产物静态包含受支持的一方模块候选集，由启动时 Profile 和 `modules.enabled` 选择启用集合；MVP、完整 Admin 与后续 fork 起点是同一架构的不同配置，不再维护两套长期演进代码线。

本文件固化 [VP-003](../vision/plans/VP-003-modular-admin-architecture.md) 的终态架构边界。

> **历史评议输入**：`MODULE-ARCHITECTURE-DRAFT.md` 已在定稿时从当前工作树移除。其只读、可复核来源固定为 Git `72017c86313c75edfe04c71ec7266767509388bb:MODULE-ARCHITECTURE-DRAFT.md`（blob `e6473129ac52f7ae67284e356e3c4ddd47a217e6`）；下文的“历史评议输入”均指该版本。它不是现行架构权威。

## 1. 内核、组合根与 DI

### 1.1 边界

- **薄内核**只提供稳定的基础契约：配置、日志、数据库、HTTP 生命周期、认证/授权接口、模块描述与能力注册协议。内核包不得导入业务模块。
- **组合根**静态导入全部受支持的一方模块，解析 Profile 与 `modules.enabled`，校验依赖图，组装能力并统一管理启动/停止。
- **模块**不得通过修改中央业务注册表接入；它们只实现框架无关的模块契约和所需能力接口。

### 1.2 DI 选型

采用 [Uber Fx](https://github.com/uber-go/fx) 作为 Go 组合根的依赖装配与生命周期工具，但不让 `fx.In`、`fx.Out`、`fx.Option` 等类型进入模块公共契约。Fx 是实现细节，不是模块 API。

[Google Wire](https://github.com/google/wire) 不采用：其官方仓库已归档并声明不再维护；同时，纯编译期注入不能单独表达本项目启动时 Profile 选择。若未来替换 Fx，只允许影响组合根，不应迫使业务模块改写公共契约。

## 2. 模块契约

每个模块必须声明稳定且不可复用的 `id`、版本、内核 API 兼容范围，以及带兼容版本约束的显式依赖（`DependsOn`；无依赖可为空列表，但必须声明）。

### 2.1 一方标准 Admin 功能模块的核心贡献（必须）

下列六项对应固定历史评议输入 D-3 的**必须**贡献点。一方**标准 Admin 功能模块**（提供可装配业务能力的 users / roles / settings / activity 等当前语义能力）在迁入统一契约时**必须**实现，不得以「能力可选」为由永久缺省。历史演示实体 `records` 已由 `0006 records_retire` 退场，不属于当前标准模块集合；其迁移账本与历史 operation-log 事件只作为兼容证据保留：

| 能力 | 贡献内容 | 级别 |
|------|----------|------|
| HTTP | 路由、处理器、中间件挂点 | **必须** |
| Schema | 页面、资源描述、数据源与动作 | **必须** |
| Authorization | 权限键与后端授权规则 | **必须** |
| Navigation | 导航节点与可见性表达式 | **必须** |
| Manifest | App Manifest 片段与前端运行时配置 | **必须** |
| Persistence | 迁移与系统数据 reconciliation | **必须** |

横切基础设施模块（如始终启用的 `operationlog`）若不暴露标准管理 UI，可经显式架构说明豁免 Schema/Navigation/Manifest 中不适用项，但不得借此让标准 Admin 功能模块缺省上述六项。

### 2.2 按需能力（可选）

下列能力**按需**实现；未实现即表示本模块不贡献该能力。**「按需」仅适用于本节，不得覆盖 §2.1 核心六项：**

| 能力 | 贡献内容 | 级别 |
|------|----------|------|
| Configuration | 带命名空间、默认值和验证规则的配置 | 可选 |
| Lifecycle | 启动、就绪、停止钩子 | 可选（有副作用或外部依赖时强烈建议） |
| Observability | 健康、指标、日志与审计事件贡献 | 可选（有独立运行特征时强烈建议） |

`DependsOn`、配置定义与种子/系统数据（历史评议输入的可选三项）分别落在模块描述的依赖字段、Configuration 能力与 Persistence 的 reconcile/bootstrap 路径中；无配置或无种子时可为空，但依赖与冲突仍 fail closed。

注册与启动必须确定性、可重复并 fail closed：未知模块、重复 ID、缺失/循环依赖、未启用的依赖、重复路由、页面、导航键、权限键或配置命名空间，均阻止启动并给出可定位错误。依赖不会被静默自动启用。

## 3. Profile 与启用语义

- `mvp`、`admin`、`custom` 等 Profile 是版本化配置，不是分支名；Profile 展开后仍形成显式 `modules.enabled` 集合。
- 运行时只能在**已编译候选集**中选择模块，不支持 `.so`、远程下载、热插拔或运行中启停。
- Profile 的默认值、覆盖优先级和最终解析结果必须可观察；未知 ID、冲突配置或不满足依赖闭包时拒绝启动。
- Profile 控制能力暴露与 UI 组合，不承诺从二进制中物理移除代码。需要物理裁剪时，由 fork 或独立构建目标负责。

## 4. 数据生命周期

### 4.1 迁移

- 保留**一个全局、不可变、带 checksum 的迁移台账**和全局顺序；沿用现有事务、连续性校验、升级前快照与恢复约束。
- 所有已编译的一方模块迁移都参与全局排序并执行，**不以模块是否启用为条件**。禁用模块不回滚、不删除其表或数据。
- 已发布迁移不得改写；重复编号、checksum 漂移、缺口或失败迁移必须 fail closed。
- 数据库中存在当前二进制不认识的已应用迁移时必须 fail closed。模块退役后，其已发布迁移描述符、编号与 checksum 仍须以 tombstone 形式保留在全局目录；退役不自动生成 drop migration。

该规则避免 Profile 切换产生不可证明的数据库形态，也保留现有实例升级与恢复证据链。

### 4.2 系统数据与种子

- **fresh bootstrap** 仅负责新库首次建立所需的初始管理员、基础角色等启动数据。
- **system-data reconcile** 以版本化、幂等方式补齐权限键、系统菜单和模块拥有的固定数据，适用于既有实例升级。
- reconciliation 不得覆盖用户拥有的字段或删除用户数据；所有权、变更规则和冲突处理必须显式。

## 5. Manifest、Schema 与前端装配

- 后端在启动时聚合已启用模块的 Manifest、Schema、导航与运行时配置，先校验再发布。
- 聚合结果必须声明协议版本、模块 API 版本和前端 `requiredCapabilities`；前端 build 不满足要求时确定性拒绝，不得降级成局部可用的未知页面。
- 规范入口为同源公共端点 `/.well-known/schema-ui/app-manifest.json`。它必须在登录前可读取、内容确定、无秘密和无用户个性化信息；后端认证与授权仍是业务 API 的最终权限边界。
- 响应支持确定性 `ETag` 和 revalidation；Vite 与生产 Nginx 必须将该精确路径代理到 API。
- 生产环境不允许静态 Manifest 静默兜底。迁移期开发兼容路径必须有期限、显式告警和删除门禁。
- 同一前端 build 必须能消费不同 Profile 的聚合结果；标准模块的增减不得要求修改 Shell 的中央路由、导航或页面注册代码。

## 6. 横切能力边界

- `operationlog` 是始终启用的横切基础设施，负责关键写操作记录；`activity` 是可选的只读查询与 UI 能力。关闭 Activity UI 不得关闭审计记录。
- `settings` 作为普通模块贡献品牌与设置能力。现有 Shell 专用通知或硬编码回调应由通用配置贡献/变更事件能力替代。
- 认证、授权、数据库、错误协议和日志属于内核能力；模块只消费稳定接口，不能私建平行实现。

## 7. 生命周期、可观测性与验证

- 组合根先完成图校验和注册，再按拓扑顺序启动；停止顺序反向执行。启动失败必须清理已启动资源。
- `/healthz` 只表示进程存活；`/readyz` 必须反映模块图、迁移、系统数据 reconciliation 与必需依赖是否就绪。
- 日志与指标至少携带 `module_id`；模块启用集、版本、Profile、依赖图和 Manifest 摘要应可诊断但不得泄露秘密。
- 验证矩阵至少覆盖：单模块契约、依赖/冲突失败、迁移与升级、fresh/reconcile、`mvp`/`admin` 双 Profile、同一前端 build、认证授权、Activity 关闭仍记录操作日志、启动/停止与失败清理、容器和 fork 路径。

## 8. 迁移与退出边界

实施路线可以先用 `activity`（及 `operationlog` 拆分）与 `settings` 验证切口，但**试点不是终态，也不能替代 VP-003 的退出标准**。终态必须将现有一方 Admin 能力迁入统一契约，并删除中央模块清单、静态生产 Manifest、Shell 特例和其他已经被新架构替代的装配路径。

R3 有界试点的**通过门闩**（继承固定历史评议输入 §4.5 / §5，细节以 [VP-003](../vision/plans/VP-003-modular-admin-architecture.md) 路线图为准）至少包括：

1. 试点模块按 §2 完整迁入（含核心六项及声明的依赖/配置/种子路径）；
2. 切除四个旧架构病灶：中心化业务 Register、全局 Schema fixtures 对试点的占用、生产静态 Manifest 依赖、Host/Shell 对试点的硬编码特例；
3. 验证 V-1 启停、V-2 禁用不可见、V-3 同一前端 build 零 diff、V-4 Schema 驱动渲染；
4. 未通过则先加固 Kernel，**不得**以「模块已写出」放行全量存量迁移（R4）。

本决策不包含：运行时插件市场、第三方不受信任模块加载、跨进程微服务拆分、全量上游协议扩张、业务领域模块本身。协议范围以 [VP-003 的继承协议基线](../vision/plans/VP-003-modular-admin-architecture.md) 为准；扩大范围须先有新的决策、递增覆盖版本与验证。上述能力需要新的愿景或明确兼容决策。

## 9. 操作 playbook（一方模块贡献）

架构边界以本文 §1–§8 为权威。将边界**操作化**为「必须 / 禁止 / 归属判定」的可执行清单见：

→ **[module-contribution-playbook.md](module-contribution-playbook.md)**（VP-004；产品模块贡献方法论，非治理 principles 修订）
