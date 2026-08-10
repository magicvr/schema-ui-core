---
doc_type: vision-plan
id: VP-008-admin-module-readiness-and-foundation-convergence
title: Admin 业务模块准入与基架收敛
status: active
vision_ref: schema-ui-core-admin-foundation@0.2.0
lead_workspace: workspace-008-admin-module-readiness
created: 2026-08-10
updated: 2026-08-10
version: 0.12.0
parent: null
---

# VP-008 · Admin 业务模块准入与基架收敛

## 状态与门闩

| 项 | 值 |
|----|-----|
| status | **`active`**；2026-08-10 用户确认激活并建立单一 lead delivery workspace `workspace-008-admin-module-readiness`；仍未产生可消费 `go` |
| Vision required | 当前投影以 [Vision Review 台账](../reviews.md) 为唯一权威；本 VP 不独立维护计数。激活、宣称“方向已稳”或产生可消费 `go` 前，适用的 Vision required 必须合法闭合 |
| 实现入口 | lead delivery 工作区已由用户确认；物理 scaffold 与 Goal 推进由 `/govern` 承接 |
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
| 当前代码健康 | Go API、React Web、构建、静态检查、协议 fixture/conformance、双 Profile、关键 E2E、跨模块 UI 可访问性下限、fresh bootstrap、升级/reconcile 与失败语义形成可重复基线 |
| 模块接入能力 | 当前一方模块先按 `standard-admin`、`infra`、`core`、`other` 分级；`standard-admin` 按 Provider M1～M6 与 HTTP、Schema、Authorization、Navigation、Manifest、Persistence 六项核验，`infra`/`core` 按架构豁免表逐项标记 N/A 与理由；完成一次非领域化接入演练 |
| 前后台 UI 协议 | 将业务模块共性能力映射到 `schema-ui-docs@v2.7.0`；逐项分类为已覆盖、宿主实现缺口、上游协议缺口或明确非目标，并冻结前端宿主能力与自定义扩展边界 |
| 缺陷与缺漏治理 | 建立带严重度、影响门禁、责任边界、证据与关闭路径的台账；阻断项进入修复和回归闭环，非阻断项有延期理由与复核触发 |
| 准入证据 | 形成机器可复跑的验证入口、证据矩阵、风险声明、审计意见与用户最终 `go` / `no-go` 裁决 |

## 最小可枚举证据面

S0 必须冻结一份可版本化的准入分母，至少包含：

1. **代码与环境**：Git commit、Go/Node/package manager 版本、依赖锁文件、配置模板、数据库起始形态和所有验证命令。
2. **模块集合与适用矩阵**：当前编译候选 Provider、Profile 默认启用集、依赖闭包、模块 descriptor、模块分级标签（`standard-admin` / `infra` / `core` / `other`）、适用检查表（全六项 / 豁免项 / N/A+理由）、迁移与 system-data ownership、证据路径。`other` 只能作为带 owner、理由、复核触发的临时发现标签，S2 方案冻结前必须收敛为明确归属、明确不在模块契约内，或用户书面接受 residual；不得以 `other` 代替 N/A 或免检。
3. **运行形态**：至少覆盖 `mvp` / `admin` 两个真实 Profile；`custom` 的未知模块、缺依赖、冲突配置和失败语义进入内核测试分母。
4. **协议面**：`I-PROTO-FULL-001` 的 12/12 domain、24/24 registry type、16/16 behavior suite、320 total = **318 executed + 2 local adapter excluded** 的现行投影，以及本地主机对上游 case 的 adapter / exclude disposition。
5. **主流程与用例选取规则**：以 `scripts/smoke.sh` 的 SM-001～SM-005 作为 readiness、登录/身份、代表页的最低 smoke 下限；只有显式隔离的 `--disposable` 运行并通过 SM-006，才可声称种子可重复性。`mvp` / `admin` 每个实际可达 Runtime Manifest `pageId` 与 `schemaUrl` 都进入清单；每个声明 CRUD 的资源至少覆盖 list、detail、create/update、delete（若明确不支持则记录理由），每个可达写操作至少有成功、未授权/权限失败或校验失败路径。每条用例固定 `profile`、`pageId`/资源 id、`schemaUrl`、权限键、预期错误码和证据路径；未列入项只能以明确 residual 留痕，不得用单个代表页关闭 exit #2。
6. **模块接入演练**：以 test fixture、probe module 或等价有界方式验证新增标准模块无需修改 Renderer/Shell 中央业务注册，并能通过 Profile、权限、导航、Manifest、Schema、迁移和回归门禁。
7. **消费路径与升级边界**：至少核对 compose/容器启动与文档化 fork bootstrap 的可复跑入口；S0 至少固定一个受支持的升级来源与目标、升级前快照/恢复路径，并明确本 VP 不要求降级，除非用户另行扩 scope。fork/compose 必须指向文档化基线分支或 commit；超出该升级窗口、恢复边界或 fork 基线的兼容诉求记录为 `N/A`/residual 或回 `/vision` 扩 scope。若本 VP 明确不纳入消费路径，必须记录 `N/A` 理由、影响范围和重新纳入触发，不得用本地主线已运行替代消费路径证据。
8. **跨模块 UI 可访问性**：至少覆盖 Renderer/Shell 布局、导航/移动导航、schema-driven 表单/列表/详情/动作、模态与动态反馈、语言切换后的共享宿主行为；冻结键盘可达性、焦点管理、语义名称/角色/状态、错误/状态播报的可复跑断言或人工核对下限。未使用的宿主能力必须记录 N/A 理由、影响范围和重新纳入触发。

分母中的 `excluded` / `N/A` 必须记录理由、影响范围和重新纳入触发；不得把“历史已关门”或“代表性测试通过”替代当前分母证据。

## 阻断与严重度量尺（S0 冻结）

本量尺是本 VP 的方向级门禁。S0 必须将量尺版本、适用 scope、证据分母和 audit scope 一并冻结；S1 只能按已冻结量尺登记和应用，不得重写量尺。S2 及以后不得在没有用户书面裁决的情况下扩大、收缩或改写 `required` 定义；新发现只可按既有量尺分类，若证据触发 blocker/major 条件则按量尺进入 required。

| 等级 | 默认状态 | 充分条件与处理 |
|------|----------|----------------|
| `blocker` | `required` | 阻止任一标准 Admin 模块启动、构建、依赖闭包、Manifest/Schema 发布或 fail-closed；破坏认证/授权、数据隔离、迁移完整性、协议兼容边界；使冻结证据不可复现；或发现影响所有未来标准业务模块的全局 protocol-gap。必须修复、用户书面接受 residual，或维持 `no-go`。 |
| `major` | 依 gate 判定 | 影响一个或多个跨模块必需能力或退出判据。若影响方案冻结、实施、验收或本 VP 冻结的生产化基架边界（容器、迁移、认证/授权、数据隔离、失败语义或消费路径），进入 `required`；若可在冻结 scope 内隔离且不影响准入，标为 non-blocking 并记录延期字段。这里的“生产边界”不等同完整生产认证、性能 SLO、灾备、威胁建模或全部运维控制。 |
| `minor` | `non-blocking` | 局部质量、体验、文档或可维护性问题，不影响冻结分母、必需流程、协议边界或后续业务 VP 解锁。进入延期/复核台账。 |
| `info` | `non-blocking` | 观察项、建议或证据增强项，不改变当前门禁；保留来源与复核触发。 |

**领域边界**：仅影响某个尚未建立的订单、钱包、类目、通知或其他具体领域的风险，默认不进入本 VP 的 `required`；只有用户书面扩展本 VP scope，或证据证明其实际影响跨模块基架/共同门禁时，才可按上述量尺升级。性能 SLO、完整威胁建模、Skills 发布矩阵等继续留在 Non-goals，除非用户另行扩 scope。

**台账映射**：信息未知使用 `I-READINESS-*`；实现缺陷与整改事实落在 lead 工作区 Goal ledger；VP 级阻断与意见响应留在 Vision Review。相同事实只保留一个 canonical finding id，其他台账通过链接引用，不得用三套状态互相覆盖。

## 证据基线有效性（S0 冻结）

每份 S0 分母和后续 S4 回归必须绑定不可变基线：候选 Git commit、构建 artifact digest（如有）、依赖锁文件 digest、Go/Node/package manager 与基础镜像版本、Profile 默认集、数据库迁移台账版本，以及 `schema-ui-docs` pin、adapter/exclude disposition。证据矩阵必须记录这些字段，并在 S5 确认结果对应用户裁决的候选 commit/artifact；不一致时不得沿用历史 pass。

**来源身份默认规则**：候选 commit 默认必须对应 clean checkout（`git status --porcelain` 为空，且没有未跟踪、staged 或工作树修改）。若验证有意使用未提交输入，必须在证据矩阵中记录：输入用途与 scope、owned-path manifest 及其 digest、不可变 patch/diff digest、未跟踪/生成文件清单及 digest、容器/外部输入清单及 digest，以及明确的纳入/排除理由；不得只记录 HEAD commit。未被 manifest 绑定的 dirty、未跟踪、生成或外部输入使该候选不可消费，必须 fail closed。

变更按以下路径处理：

| 变更类别 | 处理 |
|----------|------|
| 仅影响已冻结条目的代码/配置、工具链、依赖锁、迁移、Profile、容器或 fork 基线 | 标记受影响分母项，重跑对应命令/用例并更新证据 digest；未受影响项可复用，但必须保留关联关系 |
| 改变分母范围、模块适用规则、协议 pin/adapter/exclude disposition 或风险语义 | 回流 S0 更新分母并由用户按 P-004 裁决；在回流完成前不得进入 S5 `go` |
| S5 后候选 commit/artifact 与裁决指向不一致 | 重跑全部受影响项；无法重跑或影响范围不明时保持 `no-go` |

新增 `I-READINESS-007` 登记基线字段、变更日志、受影响项重跑结果与 S5 一致性检查。该规则只约束本 VP 已声明的基架准入分母，不要求永久支持所有历史版本。

## `go` 消费有效性（S5 后冻结）

`go` 只适用于 S5 证据矩阵所指向的候选 commit/artifact、已绑定的 patch/owned-path manifest/输入 digest、已声明 Profile/模块集合、协议 pin/disposition、升级与 fork/compose 基线及其明确解锁 scope。后续业务 VP 只能消费该候选身份和 scope 内的 `go`，不得把一次历史 `go` 泛化为所有后续主线。

`go` 不是无期限凭证。本 VP 采用“消费前新鲜度复核”规则，不以一次 S5 结论覆盖所有未来业务 VP：每个后续业务 VP 在激活前，必须针对拟消费的候选身份与解锁 scope 完成一次轻量 freshness review，并把结果写入该 VP 的 S5/消费决策记录。最低复核字段为：候选 commit/artifact 与 patch/manifest/input digest 身份；冻结命令与关键证据是否仍可执行；外部输入、证书/镜像/包源和验证环境是否可用；以及最新 Goal/Vision open finding 与 accepted residual 投影。复核必须确认消费候选和 scope 与原 `go` 一致；不一致即按下列失效规则处理。

若 freshness review 失败、关键证据不可获得、外部输入/环境不可用，或 review 发现候选身份、scope、共同门禁语义已变化，`go` 立即暂挂且不得启动该业务 VP 实现；应回流 VP-008 重验证，或由用户按 P-004 对范围/残余作出书面裁决。S5 裁决至少记录 `go_issued_at`、`last_freshness_review_at`、`next_freshness_review_trigger`（每个后续业务 VP 激活前）、`consumer_vp`、复核结果、证据 digest 与暂停/回流路径。

下列变化触发 `go` 消费失效或暂挂：源代码/配置或有意 patch 改变；依赖锁、Go/Node/package manager、基础镜像或 artifact 改变；迁移台账、Profile 默认集、模块适用矩阵、容器/fork 基线改变；`schema-ui-docs` pin、adapter/exclude disposition 或协议兼容语义改变；认证/授权、数据隔离、fail-closed、可访问性或其他已冻结共同门禁语义改变；Charter/VP scope 或退出判据改变。触发后，在受影响项完成重验证并由 `/vision`/`/govern` 留痕前，后续业务实现门闩保持关闭。

| 触发类型 | 处理与解锁条件 |
|----------|----------------|
| 仅影响已冻结证据条目的实现/依赖/环境/输入变更 | 标记受影响项，按原命令/用例重跑并更新 source identity 与 digest；无影响判断也必须留痕。全部受影响项通过后，原 `go` 才可继续被消费。 |
| 改变分母范围、协议/风险语义或解锁 scope | 回流 S0，按 P-004 重新裁决；形成新的候选基线和完整 S5 证据前，原 `go` 不得解锁。 |
| 候选来源身份无法复核、输入未绑定、影响范围不明，或消费前 freshness review 未通过/无法完成 | `go` 立即暂挂/视为不可消费；无法完成重验证时保持 `no-go`。 |

扩展 `I-READINESS-009` 登记 `go` 的适用候选、解锁 scope、失效触发、消费前 freshness review 最低字段、受影响项重验证和重新启用条件；该规则不改变当前 `planned` 状态，只冻结未来业务 VP 的消费门禁。

## 可访问性准入边界（S0 冻结）

本 VP 只冻结跨模块共享宿主的可访问性下限，不把完整 WCAG、性能 SLO、威胁建模或领域专有体验纳入默认 required。S0/S3 至少建立以下矩阵：

| 共享宿主 | 最低可判定项 | 证据形式 |
|----------|--------------|----------|
| Renderer/Shell 布局、导航与移动导航 | 键盘可达全部交互控件；焦点可见、顺序稳定；移动导航打开/关闭后的焦点去向可预期 | 可复跑键盘断言 + 必要的人工核对 |
| schema-driven 表单、列表、详情与动作 | 控件有可计算名称、角色和状态；校验/禁用/加载状态可观察；错误不会只依赖颜色 | 可复跑语义/状态断言 + 代表页面人工核对 |
| 模态与动态反馈 | 模态焦点进入/约束/恢复成立；成功、错误、异步状态有可感知反馈 | 可复跑焦点/状态断言 + 人工核对 |
| 语言切换与共享文案 | 切换后名称、错误、状态和焦点仍可感知；不存在仅依赖语言或颜色的隐含状态 | 双 Profile/locale 可复跑断言 + 人工核对 |

每个未使用宿主能力的 `N/A` 必须记录理由、影响范围和重新纳入触发；暂缓项必须记录 owner、复核日期/触发条件。失败按既有量尺分类：跨模块键盘不可达、焦点丢失或共享状态/错误不可感知且影响标准模块时为 `blocker/required` 候选；可隔离的宿主缺口按 `major` gate 判定；局部体验或证据增强项按 `minor/info`，不得静默扩大 scope。S5 证据矩阵除通用最小列外，必须能指回宿主、断言/人工核对、Profile、结果和 N/A/residual 理由。

## 方向级退出判据

在同时满足下列方向、且均有工作区 Q2 证据时，本 VP **可以**提议 `closed`：

1. **治理与事实基线一致**
   愿景层当前投影、工作区绑定和历史关门叙事一致；本 VP 的代码/环境/模块/协议/流程分母已冻结，文档主张能够指向当前主线证据，不存在会误导后续业务 VP 的状态漂移。

2. **当前主线健康可重复验证**
   API、Web、构建、静态检查、协议 fixture/conformance、双 Profile、跨模块 UI 可访问性下限和关键 E2E 具备可复跑入口；fresh bootstrap、升级/reconcile、权限与 fail-closed 路径有正反证据。失败项不得以“历史曾通过”关闭。

3. **标准模块接入路径经现网验证**
   当前模块先按 `standard-admin`、`infra`、`core`、`other` 形成名册；仅 `standard-admin` 对 Provider M1～M6 和核心六项形成全量矩阵，已声明的 `infra`/`core` 豁免项逐项 N/A 并附架构理由，禁止将 N/A 记为 blocker；`other` 必须在 S2 前收敛或有用户书面 residual。至少一次非领域化接入演练证明新增模块不需要修改 Renderer/Shell 中央业务注册，不私建平行认证、授权或数据库路径，并能被 Profile、Manifest、Schema、权限、导航和全局迁移台账正确处理。probe/fixture 默认只进入 test-only 候选集，不得进入 `mvp` / `admin` 默认启用集。

4. **前后台 UI 协议决策边界冻结**
   面向后续业务模块共性能力的需求矩阵已逐项映射到 `schema-ui-docs@v2.7.0`，并形成前端宿主能力矩阵，至少记录 component/action/reaction/page 能力、已实现/宿主缺口/明确非目标、证据路径及对应 Profile。业务模块默认只使用协议驱动 UI；协议外自定义 React 组件或路由不在本 VP 放行范围，除非另有用户书面扩 scope、Manifest 声明、测试和不修改 Renderer/Shell 中央业务注册的准入条件。前端宿主矩阵还必须记录共享 UI 可访问性下限及其证据。现有协议已覆盖但宿主缺失的项进入实现缺口；协议确实缺失的项只能走上游提案或用户确认的版本化兼容决策。影响全部未来标准业务模块的全局 protocol-gap 默认阻断 `go`；除非用户书面接受有范围、有期限和复审触发的 residual，或兼容决策已经落盘，否则受影响范围保持 `no-go`，不得以私有 Schema 语义放行。

5. **阻断缺陷完成合法闭环**
   所有按本 VP 已冻结量尺判为 `required` 的缺陷、信息项与 findings 均为 `fixed`、用户书面 `accepted-residual` 或 `user-overruled`；非阻断项记录延期理由、责任人和复核日期/触发条件。只登记不整改不能形成 `go`。

6. **准入结论可审计且可复用**
   lead 工作区 Root 完成约定范围；Goal/Vision open required = 0；完成覆盖兼容性、数据/迁移、生产化基架边界、跨模块 UI 可访问性与治理一致性的 cross 审计；用户基于证据矩阵书面确认一种明确决策形状，且证据基线、来源身份与裁决候选 commit/artifact 一致。只有在适用候选和消费有效性规则均满足时的 `go` 才允许把后续业务 VP 从规划推进到实现；`conditional-go` / `partial-go` 不得作为关闭状态或业务解锁凭证。

## 准入决策形状

若本 VP 绑定多于一个工作区，S5 `go` 与 VP 关门提案显式采用 [对齐契约](../alignment.md) 的多工作区责任规则：`lead_workspace` 必填，且只由 lead 发起可消费的 `go` / 关门提案；规范化决策记录必须指向 lead Root 的证据矩阵，并通过 Q2 路径聚合其余工作区的 support 证据。矩阵须逐区列明纳入的 exit 证据、Goal open finding、accepted residual 及未纳入项与理由，同时列出仓库级 Vision open required 投影。只有该矩阵完整、所有适用 required 均已合法闭合，并由用户书面确认后，`go` 才可消费；任一关键 support 证据不可获得、影响范围不明或 required 仍开放时均保持 `no-go`。局部 Goal 通过或单一工作区证据不得解释为整个 VP 已准入。单工作区时仍须用户书面确认；`lead_workspace` 可按对齐契约省略或等于该区。

### Residual 手递与关闭后所有者

任何后续业务 VP 在激活时消费本 VP 的 `go`，必须在其消费决策记录中引用适用 residual id 列表、影响的共同门禁、当前 owner、到期/触发条件和复审入口；不得只写“继承 VP-008 residual”。若 VP-008 已 `closed` 后 shared-foundation residual 触发，相关 `go` 消费立即暂停，由 `/vision` 选择 reopen VP-008 或建立新的准入 VP，并把决定、影响 scope 与复核证据链接回消费记录。

S5 的用户裁决必须落盘以下最小字段：`decision`、日期、证据矩阵链接、Goal/Vision finding 闭合状态、accepted residual（如有）、受影响/解锁 scope、适用候选 commit/artifact、source identity（clean 或 patch/manifest/input digest）、`go_issued_at`、`last_freshness_review_at`、`next_freshness_review_trigger`、go 生效/失效触发，以及对 roadmap 业务门闩的生效语句。证据矩阵至少包含：`exit_id`、分母项 id、命令/手续、结果、Q2 证据路径、residual/N/A 理由。

| decision | 关门与解锁语义 |
|----------|----------------|
| `go` | 仅当所有 required 项均已合法闭合、证据矩阵完整、来源身份可复核且适用候选/解锁 scope 明确时成立。允许携带已书面接受、明确不影响解锁 scope 的 residual；residual、期限、责任人与复审触发必须被后续业务 VP 消费。每个后续业务 VP 激活前还必须完成并记录 freshness review；触发 `go` 失效规则或 freshness review 未通过后，后续业务实现门闩自动暂挂，直到受影响项完成重验证并形成可消费的新鲜证据。VP 才可提议 `closed`，并解锁后续业务 VP 的实现。 |
| `conditional-go` / `partial-go` | 仅表示过程中的有界判断或实验，不是 VP status，不得提议 `closed`，不得解锁任何领域业务 VP，也不得按订单/钱包/类目/通知分别切开准入。VP 保持 `active`（若尚未激活则保持 `planned`），直到形成 `go` 或 `no-go`。 |
| `no-go` | 不解锁后续业务 VP；若意图仍继续，VP 保持 `active` 进入整改循环；若用户决定放弃，走明确的 `abandoned` 路径。`no-go` 不得作为成功关门或隐式 residual。`abandoned` 同样不解锁后续业务 VP 实现，除非另有新的准入 VP 形成 `go`。 |

probe/fixture 的生命周期必须在 S2 记录 owner、稳定 id/version、加入的 test-only 候选集、退出条件和清理/保留结果；S5 前必须移除，或以明确 test-only 资产保留，绝不进入 `mvp` / `admin` 默认启用集或生产 Manifest。

## 愿景层与实现层落点（激活前冻结）

本 VP 只定义方向级准入契约：意图、适用范围、非目标、严重度量尺、方向级退出判据、`go`/`no-go` 消费规则，以及 S0～S5 的高层阶段关系。VP 中保留的 `I-READINESS-*` 字段、证据矩阵列、命令/手续类别、生命周期字段和检查表，只是实现层必须填充的最小结构模板，不是已执行事实、运行结果或 Goal 进度。

VP-008 激活并建立 lead workspace 后，每轮 S0～S5 的具体分母实例、信息项状态、候选 commit/artifact、命令与基线、代码路径、缺陷与整改、证据矩阵、self/independent 审计意见、Goal finding、延期/residual 和 progress，必须落在该 workspace 的 Root/Goal 五件套及其 ledger；不得继续把这些执行事实追加到本 VP。若实现发现改变 Charter/VP 目的、共同门禁或解锁范围的战略偏差，先回 `/vision`；其余实施与证据由 `/govern` 承接。

`go` 的回归治理按影响范围分层：仅影响单一业务模块且不改变共享基架准入、`go` 适用性或冻结风险语义的问题，由该业务 VP 在自己的 Root/Goal 台账中整改与审计；影响共享基架、`go` 适用性或共同风险语义的问题，立即暂停 `go`，由 `/vision` 决定重开 VP-008 或建立新的准入 VP，随后由 `/govern` 承接实现与证据。业务 VP 不得以自身 Goal 记录替代 Vision Review，也不得自行修改 VP 状态或方向级门闩。

## 信息门禁

| id | 问题 / 所需信息 | 级别 | 最晚阶段 | 验证 / 决策动作 | 初始状态 |
|----|-----------------|------|----------|-----------------|----------|
| `I-READINESS-001` | 当前主线的可复跑验证分母、环境版本、Profile、数据库起始形态与关键流程究竟是什么？ | required | S0 结束前 | 从 CI、README、脚本、package/go 配置与真实运行入口抽取命令矩阵，执行首轮基线并冻结证据路径 | open |
| `I-READINESS-002` | 当前一方模块按分级标签与适用检查表是否满足对应 Provider M1～M6、核心贡献、依赖闭包、Profile 和全局迁移台账契约？ | required | S2 方案冻结前 | 对 compiled providers 与真实模块建立 `standard-admin` / `infra` / `core` / `other` 名册；按全六项或架构豁免逐项核验，N/A 必须有理由与证据；用接入演练和冲突/失败测试验证文档未漂移 | open |
| `I-READINESS-003` | 上游 fixture/conformance 的本地 adapter、explicit exclude 与现行 `I-PROTO-FULL-001` 主张是否一致，哪些会影响未来业务模块？ | required | S3 协议判断前 | 已由 workspace-008 S0/S3 与 workspace-005 v1.0.1 / D-003 / E-007 逐项复核；两项 exclusion 为 error-envelope adapter 差异，不影响域级协议范围 | **verified** |
| `I-READINESS-004` | 在尚未选择首个领域模块时，哪些跨模块共性能力足以构成全基架准入分母？ | required | S0 结束前 | 从订单、钱包、类目、通知候选抽取共性模式，只冻结列表/详情/写操作/状态流转/权限/审计/迁移/反馈等框架能力，不预设领域模型；领域特有项默认不进本 VP required | open |
| `I-READINESS-005` | cross 审计使用哪个 independent provider，覆盖哪些 compatibility/data/migration/production scope？ | required | Root S0 实施前 | 激活 VP 时由用户指定会话可用 provider；记录 self + independent 的 scope 与最低证据要求 | open |
| `I-READINESS-006` | 阻断/严重度量尺、台账映射与 S1 只应用规则是否已经冻结？ | required | S0 结束前 | 记录本节量尺版本、适用 scope、证据分母、audit scope 与用户确认；S1 起只按量尺分类，不重写定义 | open |
| `I-READINESS-007` | S0/S4 证据基线是否绑定候选 commit/artifact/lockfile/环境/pin，且变更触发、受影响项重跑与 S5 一致性检查是否已留痕？ | required | S0 结束前 | 冻结本节基线字段与变更分类；记录每次受影响项重跑，范围/协议/风险语义变化回流 S0 并经用户裁决 | open |
| `I-READINESS-008` | 跨模块 UI 可访问性下限、断言/人工核对、N/A/延期触发与失败严重度映射是否已经冻结？ | required | S0 结束前 | 为共享 Renderer/Shell、导航、表单/列表/详情/动作、模态/动态反馈和语言切换固定键盘、焦点、语义名称/角色/状态及错误/状态播报证据；记录未使用宿主的 N/A 理由和重新纳入触发 | open |
| `I-READINESS-009` | `go` 的适用候选、解锁 scope、失效触发、消费前 freshness review、受影响项重验证与重新启用条件是否已经冻结？ | required | S0 结束前 | 记录 S5 候选 commit/artifact、clean 或 patch/manifest/input digest、消费边界；为每个后续业务 VP 激活前记录候选身份、关键证据可执行性、外部输入/环境可用性、finding/residual 投影与复核结果；留痕触发后的暂挂/回流和重新验证证据 | open |

这些未知不阻断 VP 以 `planned` 落盘。`I-READINESS-005`、`I-READINESS-006`、`I-READINESS-007`、`I-READINESS-008` 与 `I-READINESS-009` 在实现层/S0 启动前必须关闭；其余各项在最晚阶段前未关闭时阻断相应方案冻结、整改或准入裁决。任何范围收缩、有界实验或 residual 必须按 P-004 由用户书面裁决。

## 建议实现阶段（供后续 `/govern` 建立 Root 纲领路线图时参考）

| 阶段 | 目的 | 建议检查点 |
|------|------|------------|
| S0 | 准入分母与门禁冻结 | 修正/复核治理投影；关闭 I-READINESS-001/004/005/006/007/008/009；冻结代码、环境、模块、协议、流程、用例选取规则、严重度量尺、证据基线有效性、来源身份、UI 可访问性下限、go 消费有效性、消费前 freshness review 与 audit scope |
| S1 | 当前情况扫描 | 只按 S0 冻结的量尺、模块适用检查表和用例选取规则登记代码缺陷、功能缺漏、治理漂移、测试/文档偏差与风险；完成冻结分母中每条命令/用例/模块检查表的 pass/fail/N/A+理由登记且无未分类项，方可进入 S2/S3；不得重写量尺或把领域特有项默认升为 required |
| S2 | 模块契约与接入演练 | 关闭 I-READINESS-002；逐模块核验 M1～M6/核心六项；完成非领域化接入演练及依赖、权限、Profile、迁移正反测试 |
| S3 | UI 协议与共性能力判断 | 关闭 I-READINESS-003；冻结业务共性能力映射与 covered/host-gap/protocol-gap/non-goal 分类；需协议变更时回 `/vision`/上游决策 |
| S4 | 阻断整改与回归 | 按风险修复 required 缺陷，补齐功能/测试/文档/治理投影；非阻断项合法延期；重跑冻结分母并保留失败到通过证据，包含跨模块 UI 可访问性断言/人工核对。超出 S0 分母的新 blocker 必须回流 S0/用户裁决扩 scope，不得静默扩大 required 整改集；完成全部 required 闭环和冻结分母回归，方可进入 S5 |
| S5 | 准入审计与裁决 | 完成证据矩阵、self + independent cross 审计、finding 响应与用户 `go` / `no-go` 决策；仅 `go` 解锁业务模块实现；`abandoned` 不解锁 |

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
| `workspace-008-admin-module-readiness` | `GOAL-001-admin-module-readiness` | lead | 2026-08-10 | `active`，单区 lead delivery；workspace/Root 已由 `/govern` scaffold，当前无 support workspace、无可消费 `go`。 |

## 关门记录

（`closed` 提议已由 lead workspace `go` 支持；正式 status 变更走 `/vision`。）

| date | outcome | summary | evidence_links | residuals |
|------|---------|---------|----------------|-----------|
| 2026-08-10 | `go` 签发（可提议 closed） | 用户书面签发 `go`，候选 `ed99e88`（clean）；S0–S5 全部完成、open required = 0；解锁后续标准业务模块实现 | workspace-008 [GOAL-007 D-001](../workspace-008-admin-module-readiness/GOAL-007-s5-admission-audit-and-verdict/01-decision/D-001-s5-go-decision.md) + [S5-evidence-matrix](../workspace-008-admin-module-readiness/GOAL-007-s5-admission-audit-and-verdict/attachments/S5-evidence-matrix.md) | F-007（上传授权深度）deferred，owner=VP-008 lead；后续业务 VP 消费时须 freshness review |

## 规划修订短史

| date | version | change |
|------|---------|--------|
| 2026-08-10 | `0.1.0` | 用户确认采用全基架准入结构并落盘：新建 planned VP-008；范围含现状扫描、代码/功能/治理缺口、UI 协议判断、阻断整改和 `go`/`no-go`；不激活、不建区、不重开历史 VP。 |
| 2026-08-10 | `0.2.0` | 响应 VRev-017：采纳 `conditional` / `editorial`，保留原 verdict 与 finding 原文；F-V032 → fixed，新增 S0 冻结的阻断/严重度量尺、S1 只应用不重写、领域特有项默认不进 required；同批固定 F-V033/034/035 的决策形状、E2E/CRUD 选取规则与 protocol-gap/probe 生命周期。保持 `planned`、0 区；不激活或开区。 |
| 2026-08-10 | `0.3.0` | 响应 VRev-018：采纳 `conditional` / `editorial`，保留原 verdict 与 finding 原文；F-V036 → fixed，补齐模块分级与适用检查表、前端宿主矩阵、S1/S4 完成界、S5 证据矩阵、`abandoned` 门闩与 fork/容器分母规则；F-V037～F-V040 同批 fixed。保持 `planned`、0 区；不激活或开区。 |
| 2026-08-10 | `0.4.0` | 响应 VRev-019：采纳 `conditional` / `editorial`，保留原 verdict 与 finding 原文；F-V041 → fixed，增加证据基线有效性、变更触发/重验证与候选提交一致性；F-V042～F-V044 同批 fixed。保持 `planned`、0 区；不激活或开区。 |
| 2026-08-10 | `0.5.0` | 响应 VRev-020：采纳 `conditional` / `editorial`，保留原 verdict 与 finding 原文；F-V045 → fixed，增加跨模块 UI 可访问性下限、断言/人工核对、N/A/延期触发与严重度映射，并接入 S0/S4/S5。保持 `planned`、0 区；不激活或开区。 |
| 2026-08-10 | `0.6.0` | 响应 VRev-021：采纳 `conditional` / `editorial`，保留原 verdict 与 finding 原文；F-V046 → fixed，增加 clean checkout 默认、未提交输入 patch/manifest/digest 绑定；F-V047 → fixed，规定 go 适用候选、失效触发与重验证规则。保持 `planned`、0 区；不激活或开区。 |
| 2026-08-10 | `0.7.0` | 响应 VRev-022：采纳 `conditional` / `editorial`，保留原 verdict 与 finding 原文；F-V048 → fixed，选择每个后续业务 VP 激活前的消费前 freshness review，冻结最低复核字段、S5 记录字段及失败后的暂挂/回流路径，并接入 roadmap 业务门闩。保持 `planned`、0 区；不激活或开区。 |
| 2026-08-10 | `0.8.0` | 响应 VRev-023：采纳 `conditional` / `editorial`，保留原 verdict 与 finding 原文；V-F049 → fixed，冻结愿景层方向契约与实现层 Root/Goal 台账边界；V-F050 recommended 同批 fixed，按领域局部问题与共享基架问题区分 `go` 回归治理所有者。保持 `planned`、0 区；不激活或开区。 |
| 2026-08-10 | `0.9.0` | 响应 VRev-024 / VRev-025：采纳 `conditional` / `editorial`，保留原 verdict 与 finding 原文；V-F051 → fixed，显式采用多工作区 lead 提案、support 证据聚合、用户书面确认及缺证据即 `no-go` 规则；V-F052 → fixed，删除 VP 内独立计数并以 `reviews.md` 为唯一 Vision required 投影。保持 `planned`、0 区；不激活或开区。 |
| 2026-08-10 | `0.10.0` | 响应 VRev-026：采纳 `pass` / `no-change`，保留原 verdict 与 finding 原文；V-F053 recommended → fixed，增加后续业务 VP residual 手递字段与 VP-008 closed 后 shared-foundation residual 的 `/vision` reopen/新准入 VP 所有者。用户随后确认激活 VP-008；状态改为 `active`，0 区进入 14 日空转宽限，未创建 workspace/Root/Goal 或产生 `go`。 |
| 2026-08-10 | `0.11.0` | 用户确认单工作区 lead `workspace-008-admin-module-readiness`、Root `GOAL-001-admin-module-readiness` 与 GitHub Copilot `/audit` independent provider；`/govern` 完成 workspace/Root scaffold。VP-008 保持 `active`，单区已绑定但仍未产生 `go`；S0 required 信息项与运行证据继续由工作区台账承接。 |
| 2026-08-10 | `0.12.0` | editorial 勘误投影：继承 workspace-005 `I-PROTO-FULL-001` v1.0.1 与 workspace-008 A-003 的现行分母，改为 320 total = 318 executed + 2 local adapter excluded；不改变 VP-008 意图、门闩或 `go` 状态。 |
