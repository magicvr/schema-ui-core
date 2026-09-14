---
doc_type: vision-plan
id: VP-036-admin-command-palette
title: Admin 全局检索与 Command Palette
status: planned
vision_ref: schema-ui-core-admin-foundation@0.4.0
lead_workspace:
created: 2026-09-10
updated: 2026-09-10
version: 0.1.0
parent: null
---

# VP-036 · Admin 全局检索与 Command Palette

## 状态与激活门禁

| 项 | 值 |
|-----|-----|
| status | **`planned`**（2026-09-10 · 0 区） |
| 组合位置 | **Admin 功能分支 · 体验增强**；承接 VP-034 导航分组后的跨模块发现能力 |
| Vision Review | 计划阶段 [VRev-090](../reviews/VRev-090-vp036-admin-command-palette-planned.md) self `pass`；本条不是激活许可 |
| 激活前置 | Admin 类 freshness review + 激活就绪 self Review；之后由 `/govern` scaffold delivery 工作区 |
| 基础设施边界 | 首波不消耗 `RT-X01` 搜索引擎、Redis、MQ 或多实例 trigger |

## 意图

在现有 Schema 驱动 Admin Shell、导航分组、模块贡献模型、权限与 Profile 体系之上，建立一个**键盘优先、跨模块的 Admin 发现入口**。本 VP 的首波将“全局检索”收窄为对已注册的页面、导航项和声明式动作进行检索，并提供统一的 `SearchableItem` / provider 接缝；结果必须遵守当前 Profile、权限和既有路由守卫。

该 VP 是面向用户的 Admin 体验增强，不是 VP-010 的符合性整改，也不承载新的业务域。实体记录全文检索仅作为后续可能的独立范围，不进入本 VP 首波退出分母。

## 首波范围与边界

| 范围 | 本 VP 首波 | 不在本 VP |
|------|-----------|-----------|
| 发现对象 | 已注册模块、页面、导航项、声明式 UI 动作 | 业务实体全文索引、任意数据库表扫描 |
| 入口 | 全局快捷键、Command Palette、键盘导航、直接跳转 | 远程命令执行、绕过既有权限/路由守卫 |
| 聚合 | 跨模块 `SearchableItem` / provider 聚合、稳定排序与去重 | 为每个模块增加 Shell 中央业务分支 |
| 可见性 | 按当前 Profile 与权限过滤；直接 URL 行为与现有守卫一致 | 新的权限绕过、组织权限模型、SSO/多租户 |
| 体验 | 中英文、浅色/深色、加载/空态/错误态、键盘可访问 | Saved Views、批量结果中心、未保存保护、Toast 全局重做 |
| 基础设施 | 先用现有注册信息和既有查询能力；保留未来接缝 | 专用搜索引擎、Redis、MQ、跨进程索引、多实例 |

## 与相邻 VP / 路线图的边界

| VP / 方向 | 关系 |
|-----------|------|
| **VP-034** | 复用已交付的导航分组、直接 URL 与激活态语义；不重开 VP-034，不把 Dashboard residual 扩大为本 VP 的历史修订 |
| **VP-010** | 本 VP 是向好演进的产品能力，不是 as-designed / as-built 偏差整改；若实现发现既有协议符合性问题，转 VP-010 |
| **VP-009** | 共享基架安全问题仍归持续生产加固；Command Palette 不降低既有 API/路由鉴权 |
| **VP-008 `go`** | 默认 additive Admin 产品面；若改变 Profile 默认集、模块矩阵或 Manifest 装配语义，必须暂停并按 freshness / `go` 规则复核 |
| **RT-X01 / RT-X02** | 专用搜索引擎与 DB 全文检索保持 `trigger-gated`；只有实体级搜索需求和规模证据成立后再经 `/vision` 决定 |
| **业务域分支** | 不新增 Catalog、订单、支付、库存、CMS 等业务域；没有新触发时不预开第二个业务域 |

## 方向级退出判据

在同时满足下列方向时，本 VP **可以**有界或完整关门（证据必须在工作区目标内）：

1. **分母与契约**：当前纳入范围的页面、导航项和声明式动作有可核对清单；`SearchableItem` / provider 契约、字段语义、版本化与排除项已冻结。
2. **跨模块聚合**：来自现有及可选模块的条目可稳定聚合、去重和排序；不要求 Shell 为每个模块维护中央业务注册分支。
3. **权限与 Profile 安全**：结果按当前 Profile 与权限过滤；直接 URL、既有路由守卫和动作权限不被 Command Palette 绕过；mvp/admin/demo/custom 覆盖有矩阵证据。
4. **Command Palette 体验**：快捷键、打开/关闭、键盘移动/确认/退出、焦点管理、搜索输入、加载/空态/错误态可用，并满足既有中英文与浅色/深色产品约定。
5. **导航联动与回归**：从 Palette 进入已分组页面时，当前导航分组正确展开；既有导航、top/user slot、模块贡献和直接 URL 回归通过。
6. **基础设施与范围保持**：未实现或解除 `RT-X01`、Redis、MQ、多实例等 gated 能力；未把 Saved Views、批量中心、未保存保护、Toast 全局重做或第二业务域混入本 VP。
7. **证据与审计**：退出矩阵、浏览器/自动化回归和必要的独立意见已落盘，开放 required finding = 0，并经用户确认关门。

## 纲领路线图（实现层由 `/govern` 承接）

```text
R1 范围与信息冻结：页面/导航/动作分母、权限/Profile 语义、排序/去重、快捷键与实体搜索排除
  → R2 SearchableItem 契约与跨模块 provider 聚合
  → R3 Admin Shell Command Palette、键盘可访问、i18n/theme、直接路由与分组联动
  → R4 Profile×权限×路由回归、证据矩阵、边界复核与关门
```

纲领阶段通常串行；同一阶段内可并行的实现目标须在 R1 边界冻结后由 `/govern` 按证据创建。

## 信息需求（P-005）

| id | 要回答的问题 | 级别 | 影响门禁 | 最晚阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|--------------|------|----------|----------|------------------|------|-------------|-------------|
| I-036-001 | 当前所有可检索页面、导航项和声明式动作的精确分母是什么？需覆盖哪些 Profile / optional module？ | required | R1 范围冻结、R2 聚合、R4 回归 | R1 | 扫描模块 provider、Manifest fragment、现有导航/动作注册，建立 item→module→profile 矩阵 | collecting | R1 冻结前复核 | 首波只承诺注册页面/导航/声明式动作；实体记录不进分母 |
| I-036-002 | Palette 结果的权限、Profile、直接 URL 与动作守卫如何与现有语义保持一致？ | required | R1/R3/R4 安全与回归 | R1 | 对照现有权限过滤、路由守卫、Profile 装配和 action binding；建立允许/拒绝矩阵 | collecting | R1 方案冻结前复核 | 不新增权限绕过；具体实现待 R1 冻结 |
| I-036-003 | 结果字段、标签本地化、排序/去重、快捷键、结果上限和焦点/ARIA 语义是什么？ | required | R1/R3/R4 体验验收 | R1 | 用户体验方案 + 既有设计系统/locale 约定 + 浏览器可访问性验证 | collecting | R1 方案冻结前复核 | 需避免将 UX 偏好伪装成已冻结实现细节 |
| I-036-004 | 首波是否承诺实体级全局搜索，以及是否触发 `RT-X01` / `RT-X02`？ | required | VP 范围与基础设施门禁 | R1 | 用户书面确认的本 VP 边界；若后续新增实体搜索，另行 `/vision` 复核 | **verified (user decision)** | 实体搜索需求或规模证据出现时复核 | 首波不承诺实体全文搜索；`RT-X01` / `RT-X02` 保持 trigger-gated |
| I-036-005 | 最近搜索、固定项或持久化偏好是否需要进入首波？ | non-blocking | R3 体验扩展 | R3 | 实现前评估；若需要，另立体验 VP 或追加有界范围决策 | deferred | 不进入首波；下一 UX VP 规划时复核 | Saved Views / 最近项不作为本 VP 退出条件 |
| I-036-006 | 激活时当前代码候选是否仍满足 Admin 类 freshness 与 VP-008 `go` 消费有效性？ | required | 激活 | 激活前 | `/vision` 复核协议 pin、依赖锁、迁移、Profile 默认集、provenance 与区间变更 | open | 激活前复核 | 计划阶段不宣称已通过 freshness |

## 工作区绑定

| workspace_id | root_goal | role | joined | notes |
|--------------|-----------|------|--------|-------|
| — | — | — | — | `planned` 可为 0 区；激活后由用户确认 workspace slug，再交 `/govern` scaffold |

## 关门记录

（仅 `closed` / `abandoned` 时填写。）

| date | outcome | summary | evidence_links | residuals |
|------|---------|---------|----------------|-----------|
| — | — | — | — | — |

## 规划修订短史

| date | change |
|------|--------|
| 2026-09-10 | 初创 `planned`；用户确认下一拍采用 Admin 功能体验增强方向，首波冻结为权限安全的 Command Palette + 模块/导航/声明式动作检索；实体级全文搜索、RT-X01/RT-X02、Saved Views、批量中心、未保存保护与 Toast 全局重做不进首波；不改 Charter、不新开工作区、不消耗 gated trigger。 |

## 声明

本文件是已确认的 Vision Plan 意图，不是 Goal 五件套、实现事实或 progress 权威。激活、工作区创建、阶段执行和 Goal 审计分别由 `/vision` 与 `/govern` 按对齐契约承接。
