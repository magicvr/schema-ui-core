---
doc_type: vision-plan
id: VP-008-admin-module-readiness-and-foundation-convergence
title: Admin 业务模块准入与基架收敛
status: planned
vision_ref: schema-ui-core-admin-foundation@0.2.0
lead_workspace:
created: 2026-08-10
updated: 2026-08-10
version: 0.1.0
parent: null
---

# VP-008 · Admin 业务模块准入与基架收敛

## 状态与门闩

| 项 | 值 |
|----|-----|
| status | **`planned`**；尚未激活或绑定工作区 |
| Vision required | 当前仓库级 open required = **0**；本 VP 激活前建议独立 Vision Review |
| 实现入口 | 激活与 lead delivery 工作区由后续 `/vision` 用户确认；物理 scaffold 与 Goal 推进交 `/govern` |
| 业务模块门闩 | 本 VP 未形成 `go` 结论前，不启动订单、钱包、类目、通知等正式业务模块 VP 的实现 |

## 意图

在 [VP-003](VP-003-modular-admin-architecture.md) 的单主线模块架构、[VP-004](VP-004-module-contribution-readiness.md) 的模块贡献操作契约、[VP-005](VP-005-design-system-and-ui-experience.md) 的产品化界面、[VP-006](VP-006-full-protocol-contract-v2-7-0.md) 的完整协议兼容面和 [VP-007](VP-007-localization-and-system-settings.md) 的 locale/settings 基础之上，在正式业务模块开发前完成一次面向当前代码主线的全基架准入与收敛波次。

本波次不把历史 VP 的 `closed` 直接等同于当前主线仍然无缺陷，也不以重新审判历史关门为目的。它以可重复运行的现状证据核对文档主张与实际代码，梳理代码缺陷、功能缺漏、治理漂移和前后台 UI 协议需求；对阻断业务模块开发的事项完成有界整改，并最终给出可审计的 `go` / `no-go` 结论。

### 用户确认（2026-08-10）

- 采用“全基架准入”范围，而不是仅围绕首个业务模块做局部扫描。
- 新建独立 VP，不重开 VP-001～007，不修改 Charter 的目的、成功边界或非目标。
- 本波次不是仅交付差距报告；阻断项须在本 VP 内修复、由用户合法接受残余，或明确维持 `no-go`。
- 前后台 UI 协议增补必须先区分“现有协议已覆盖但宿主未实现”与“上游协议确实缺失”；不得用未经治理的私有 Schema 语义赶业务进度。
- 本轮仅落盘 `planned` VP；激活、audit provider、lead delivery 工作区和 Root 名称留给后续门禁确认。

### 准入范围

| 区域 | 方向级范围 |
|------|------------|
| 治理一致性 | Charter、VP、roadmap、reviews、workspaces 与绑定区投影一致；历史状态不被当前摘要误写；开放 required 可发现 |
| 当前代码健康 | Go API、React Web、构建、静态检查、协议 fixture/conformance、双 Profile、关键 E2E、fresh bootstrap、升级/reconcile 与失败语义形成可重复基线 |
| 模块接入能力 | 当前一方标准 Admin 模块按 Provider M1～M6 与 HTTP、Schema、Authorization、Navigation、Manifest、Persistence 六项核验；完成一次非领域化接入演练 |
| 前后台 UI 协议 | 将业务模块共性能力映射到 `schema-ui-docs@v2.7.0`；逐项分类为已覆盖、宿主实现缺口、上游协议缺口或明确非目标 |
| 缺陷与缺漏治理 | 建立带严重度、影响门禁、责任边界、证据与关闭路径的台账；阻断项进入修复和回归闭环，非阻断项有延期理由与复核触发 |
| 准入证据 | 形成机器可复跑的验证入口、证据矩阵、风险声明、审计意见与用户最终 `go` / `no-go` 裁决 |

## 最小可枚举证据面

S0 必须冻结一份可版本化的准入分母，至少包含：

1. **代码与环境**：Git commit、Go/Node/package manager 版本、依赖锁文件、配置模板、数据库起始形态和所有验证命令。
2. **模块集合**：当前编译候选 Provider、Profile 默认启用集、依赖闭包、模块 descriptor、六项核心贡献、迁移与 system-data ownership。
3. **运行形态**：至少覆盖 `mvp` / `admin` 两个真实 Profile；`custom` 的未知模块、缺依赖、冲突配置和失败语义进入内核测试分母。
4. **协议面**：`I-PROTO-FULL-001` 的 12/12 domain、24/24 registry type、16/16 behavior suite、320/320 case 现行投影，以及本地主机对上游 case 的 adapter / exclude disposition。
5. **主流程**：启动与 readiness、登录/刷新/登出、权限正反例、Manifest/导航/Schema 页面、标准 CRUD、Settings、fresh bootstrap、既有库升级/reconcile、失败与恢复边界。
6. **模块接入演练**：以 test fixture、probe module 或等价有界方式验证新增标准模块无需修改 Renderer/Shell 中央业务注册，并能通过 Profile、权限、导航、Manifest、Schema、迁移和回归门禁。

分母中的 `excluded` / `N/A` 必须记录理由、影响范围和重新纳入触发；不得把“历史已关门”或“代表性测试通过”替代当前分母证据。

## 方向级退出判据

在同时满足下列方向、且均有工作区 Q2 证据时，本 VP **可以**提议 `closed`：

1. **治理与事实基线一致**
   愿景层当前投影、工作区绑定和历史关门叙事一致；本 VP 的代码/环境/模块/协议/流程分母已冻结，文档主张能够指向当前主线证据，不存在会误导后续业务 VP 的状态漂移。

2. **当前主线健康可重复验证**
   API、Web、构建、静态检查、协议 fixture/conformance、双 Profile 和关键 E2E 具备可复跑入口；fresh bootstrap、升级/reconcile、权限与 fail-closed 路径有正反证据。失败项不得以“历史曾通过”关闭。

3. **标准模块接入路径经现网验证**
   当前标准 Admin 模块对 Provider M1～M6 和核心六项形成逐项矩阵；至少一次非领域化接入演练证明新增模块不需要修改 Renderer/Shell 中央业务注册，不私建平行认证、授权或数据库路径，并能被 Profile、Manifest、Schema、权限、导航和全局迁移台账正确处理。

4. **前后台 UI 协议决策边界冻结**
   面向后续业务模块共性能力的需求矩阵已逐项映射到 `schema-ui-docs@v2.7.0`。现有协议已覆盖但宿主缺失的项进入实现缺口；协议确实缺失的项只能走上游提案或用户确认的版本化兼容决策。受未决协议项影响的业务范围保持 `no-go`，不得以私有 Schema 语义放行。

5. **阻断缺陷完成合法闭环**
   所有影响业务模块方案冻结、实施、验收或生产边界的 required 缺陷、信息项与 findings 均为 `fixed`、用户书面 `accepted-residual` 或 `user-overruled`；非阻断项记录延期理由、责任人和复核日期/触发条件。只登记不整改不能形成 `go`。

6. **准入结论可审计且可复用**
   lead 工作区 Root 完成约定范围；Goal/Vision open required = 0；完成覆盖兼容性、数据/迁移、生产与治理一致性的 cross 审计；用户基于证据矩阵书面确认 `go` 或 `no-go`。只有 `go` 才允许把后续业务 VP 从规划推进到实现。

## 信息门禁

| id | 问题 / 所需信息 | 级别 | 最晚阶段 | 验证 / 决策动作 | 初始状态 |
|----|-----------------|------|----------|-----------------|----------|
| `I-READINESS-001` | 当前主线的可复跑验证分母、环境版本、Profile、数据库起始形态与关键流程究竟是什么？ | required | S0 结束前 | 从 CI、README、脚本、package/go 配置与真实运行入口抽取命令矩阵，执行首轮基线并冻结证据路径 | open |
| `I-READINESS-002` | 当前一方模块是否仍逐项满足 Provider M1～M6、核心六项、依赖闭包、Profile 和全局迁移台账契约？ | required | S2 方案冻结前 | 对 compiled providers 与真实模块做清单核验；用接入演练和冲突/失败测试验证文档未漂移 | open |
| `I-READINESS-003` | 上游 fixture/conformance 的本地 adapter、explicit exclude 与现行 `I-PROTO-FULL-001` 主张是否一致，哪些会影响未来业务模块？ | required | S3 协议判断前 | 对 12/24/16/320 分母和本地 disposition 逐项复核，特别核对 error envelope 与多轮 `$deps` reactions | open |
| `I-READINESS-004` | 在尚未选择首个领域模块时，哪些跨模块共性能力足以构成全基架准入分母？ | required | S0 结束前 | 从订单、钱包、类目、通知候选抽取共性模式，只冻结列表/详情/写操作/状态流转/权限/审计/迁移/反馈等框架能力，不预设领域模型 | open |
| `I-READINESS-005` | cross 审计使用哪个 independent provider，覆盖哪些 compatibility/data/migration/production scope？ | required | Root S0 实施前 | 激活 VP 时由用户指定会话可用 provider；记录 self + independent 的 scope 与最低证据要求 | open |

这些未知不阻断 VP 以 `planned` 落盘。`I-READINESS-005` 在实现层启动前必须关闭；其余各项在最晚阶段前未关闭时阻断相应方案冻结、整改或准入裁决。任何范围收缩、有界实验或 residual 必须按 P-004 由用户书面裁决。

## 建议实现阶段（供后续 `/govern` 建立 Root 纲领路线图时参考）

| 阶段 | 目的 | 建议检查点 |
|------|------|------------|
| S0 | 准入分母与门禁冻结 | 修正/复核治理投影；关闭 I-READINESS-001/004/005；冻结代码、环境、模块、协议、流程与 audit scope |
| S1 | 当前情况扫描 | 执行完整基线；形成代码缺陷、功能缺漏、治理漂移、测试/文档偏差与风险台账，不把扫描结果写成已修复事实 |
| S2 | 模块契约与接入演练 | 关闭 I-READINESS-002；逐模块核验 M1～M6/核心六项；完成非领域化接入演练及依赖、权限、Profile、迁移正反测试 |
| S3 | UI 协议与共性能力判断 | 关闭 I-READINESS-003；冻结业务共性能力映射与 covered/host-gap/protocol-gap/non-goal 分类；需协议变更时回 `/vision`/上游决策 |
| S4 | 阻断整改与回归 | 按风险修复 required 缺陷，补齐功能/测试/文档/治理投影；非阻断项合法延期；重跑冻结分母并保留失败到通过证据 |
| S5 | 准入审计与裁决 | 完成证据矩阵、self + independent cross 审计、finding 响应与用户 `go` / `no-go` 决策；仅 `go` 解锁业务模块实现 |

阶段和子目标由实现层根据 P-001/P-005 进一步冻结；本表不是 Goal progress 或已实施事实。

## Non-goals（非目标）

- 不在本 VP 实现订单、钱包、类目、通知或其他领域产品模块，也不预先冻结其领域模型。
- 不重开 VP-001～007，不改写其历史关门证据、Goal status 或 progress。
- 不修改 Charter 的目的、方向级成功边界或非目标；本 VP 是现行愿景下的新准入波次。
- 不以全面重构、清零所有技术债或追求零 recommended issue 作为关门条件；只有影响准入的 required 项阻断。
- 不重新定义 `schema-ui-docs@v2.7.0`，不以本地私有 Schema 扩展替代上游提案或明确兼容决策。
- 不用扫描报告、测试数量或历史 pass 单独证明生产就绪；结论必须对照冻结分母与当前证据。
- 不为 VP 建 Goal 五件套，不在 `docs/vision/` 记录 progress%。

## 与前后 VP 的关系

| VP | 关系 |
|----|------|
| VP-002 | 消费生产级 Admin/fork 基线作为待复核历史主张；不重开其工作区。 |
| VP-003 / VP-004 | 继承单主线、Profile、模块契约与贡献 playbook；以现网矩阵和接入演练检查是否漂移。 |
| VP-005 | 继承设计系统、Renderer/Shell 产品化基线；只修复阻断业务模块准入的 UI/UX 缺陷。 |
| VP-006 | 继承 `schema-ui-docs@v2.7.0` 全量兼容主张与 `I-PROTO-FULL-001`；重新核对本地 adapter/exclude disposition 对未来模块的影响。 |
| VP-007 | 继承双语、Settings、错误反馈与双 Profile 基础；核对当前主线而不重开历史交付。 |
| 后续业务 VP | 必须消费本 VP 的 `go` 证据、开放 residual 与协议决策；不得把本 VP 的未决项当作已交付能力。 |

## 工作区绑定

| workspace_id | root_goal | role | joined | notes |
|--------------|-----------|------|--------|-------|
| — | — | lead | — | `planned`，0 区；激活、lead delivery 与 Root 命名待后续 `/vision` 用户确认，物理 scaffold 交 `/govern` |

## 关门记录

（仅 `closed` / `abandoned` 时填写。）

| date | outcome | summary | evidence_links | residuals |
|------|---------|---------|----------------|-----------|
| — | — | — | — | — |

## 规划修订短史

| date | version | change |
|------|---------|--------|
| 2026-08-10 | `0.1.0` | 用户确认采用全基架准入结构并落盘：新建 planned VP-008；范围含现状扫描、代码/功能/治理缺口、UI 协议判断、阻断整改和 `go`/`no-go`；不激活、不建区、不重开历史 VP。 |
