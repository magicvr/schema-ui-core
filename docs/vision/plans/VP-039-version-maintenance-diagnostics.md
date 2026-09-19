---
doc_type: vision-plan
id: VP-039-version-maintenance-diagnostics
title: Admin 版本更新、维护提示与诊断报告
status: active
vision_ref: schema-ui-core-admin-foundation@0.4.0
lead_workspace: workspace-039-version-maintenance-diagnostics
created: 2026-09-19
updated: 2026-09-19
version: 0.2.0
parent: null
---

# VP-039 · Admin 版本更新、维护提示与诊断报告

## 状态、激活与关门门禁

| 项 | 值 |
|-----|-----|
| status | **`active`**（2026-09-19 · v0.2.0 · lead `workspace-039-version-maintenance-diagnostics`） |
| 组合位置 | **Admin 功能分支 · 体验增强**；体验清单收口（版本更新 / 维护提示 / 诊断报告） |
| 计划阶段 Vision Review | [VRev-101](../reviews/VRev-101-vp039-vp040-planned.md) self `pass`（0 required；`V-F131` 激活事务内 `fixed`） |
| 激活就绪 Vision Review | [VRev-102](../reviews/VRev-102-vp039-activation.md) self `pass`（0 required） |
| 激活门禁 | **已满足**：① `I-039-004` 接受默认候选（Shell 横幅 + 复用 `admin.system-monitoring`，不新模块、不改默认集，不暂挂 `go`）；② `I-039-005` Admin 类 freshness **PASS**（`7e5ce891` → `6197e802`）；③ 激活就绪 self Review；④ slug 按惯例确认 |
| 基础设施边界 | 首波不消耗 Redis、MQ、多实例、搜索引擎、`timestamptz` schema 迁移或文件扫描 trigger；不重开 VP-012 / VP-015 / VP-025 |
| 与 VP-040 | **正交**。C1 `timestamptz` 由 [VP-040](VP-040-timestamptz-persistence-contract.md) 另立；本 VP **不并入** schema 迁移。用户 2026-09-19 裁决：先本 VP，VP-040 本波之后再激活 |

## 用户已裁决（2026-09-19 · P-004）

| 项 | 裁决 |
|----|------|
| 组合层下一拍 | **本 VP**（体验增强收口）；不合成 C1；不塞 VP-010 |
| 架构候选 C1（`RES-T03-tz`） | **另立 VP-040**，保持 `planned`，**本波之后再激活** |
| workspace / Root slug | `workspace-039-version-maintenance-diagnostics` / `GOAL-001-version-maintenance-diagnostics`（走流程惯例） |
| `I-039-004` 承载面 | **默认候选**：Shell 持久横幅 + 复用 `admin.system-monitoring`；不新建模块；不改 Profile 默认集；**不暂挂 `go`** |

## 意图

VP-012 已交付 maintenance / degraded / read-only **四模式写门禁**、Host bootstrap 可用性投影与 `SERVICE_*` 错误码，当时明文把「运行时管理 UI」写成 **UI 可后置**。写拒绝已经 fail-closed，但 Admin Shell 没有持久、可感知的维护/降级/只读横幅；操作者主要在请求失败后才看见 Toast。

版本身份已在 `admin.system-monitoring` 状态行暴露（`version` / `commit` / `uptimeSeconds`），对 fork / 包消费者仍缺少**跨页面的版本与升级提示**（当前版本从哪来、已知升级说明入口在哪）。诊断面已有 `/healthz` `/readyz` 与系统监控页，但没有面向管理员的、与维护模式/版本绑定的轻量摘要。

本 VP 把「版本更新、维护提示和诊断报告」收口为有界 Admin 产品能力：让已交付的运行时模式、版本身份与探活/就绪**可被看见、可被理解、可被权限过滤**，而不是再建一套可观测或配置迁移产品。

本 VP 是面向用户的 Admin 体验增强，不是 VP-010 的符合性整改，不承载业务域，也不重开 VP-012。

## 首波范围与边界

| 范围 | 本 VP 首波 | 不在本 VP |
|------|-----------|-----------|
| 维护提示 | 当 `runtime.mode` 为 maintenance / degraded / read-only 时，已登录 Admin 显示持久、本地化、浅色/深色可用的横幅；写拒绝继续走既有门禁，失败反馈与横幅语义一致 | 重做四模式写门禁；热切换 `runtime.mode` 而不重启（热加载不进分母）；伪造 Host 新 availability 枚举 |
| 版本提示 | 授权管理员可看到当前版本身份（复用已有 `pkg/version` / 系统监控字段）；提供升级说明入口（文档/changelog 链接或站内只读说明，R1 冻结） | 自动升级、远程拉取发行说明、运行时插件市场、把 fork 与包消费做成两套版本面 |
| 诊断报告 | 面向管理员的轻量摘要：运行时模式、探活/就绪、版本/commit、已启用模块计数等**已交付字段**的只读聚合 | Grafana / Sentry / 连续剖析 / 新指标系列；重开 VP-015；把 Job/Store/对象存储指标扩进分母 |
| 承载面 | R1 前由 `I-039-004` 冻结：默认候选 = **Shell 横幅 + 复用 `admin.system-monitoring`**，不新增模块、不改 Profile 默认集 | 默认为新模块进默认集；改 `ResolveProfile` / Manifest 装配语义 |
| 体验 | 中英文、浅色/深色、键盘可访问、与 VP-037 统一反馈约定一致 | Command Palette / Saved Views / 批量结果中心重做；实体全文检索 |
| 基础设施 | 消费既有 config `runtime.mode`、Host bootstrap、operational gate、system-monitoring status | Redis / MQ / 多实例 / `timestamptz` / 文件扫描执行器 |

## 与相邻 VP / 路线图的边界

| VP / 方向 | 关系 |
|-----------|------|
| **VP-012** | **消费**已交付的四模式写门禁、Host/status 投影与 `SERVICE_*` 语义，**不重开** VP-012；本 VP 承接当时「UI 可后置」 |
| **VP-015** | 不重开可观测 VP；不把本 VP 做成 Grafana/Sentry；探活/就绪只读消费 |
| **VP-025** | 配置包导出/diff/dry-run/导入已 closed；本 VP 不重做配置迁移 |
| **VP-007** | 复用 locale / settings 体验基线；维护提示文案走既有 i18n，不新开 locale 运行时 |
| **VP-037** | 复用统一 Toast/错误恢复；横幅是持续状态，不是把 maintenance 只做成一次性 Toast |
| **VP-038** | 不重开作业中心；诊断摘要不是 Job 结果中心 |
| **VP-010** | 本 VP 是向好演进的产品能力，不是 as-designed / as-built 偏差整改；若实现发现既有协议符合性问题，转 VP-010 |
| **VP-009** | 共享基架安全问题仍归持续生产加固；横幅/诊断不得降低鉴权 |
| **VP-008 `go`** | 默认候选不改 Profile 默认集与装配语义 → **不暂挂 `go`**。若激活前 `I-039-004` 改为新模块进默认集或改装配语义，须按 freshness / `go` 规则暂停与复核 |
| **VP-040** | 架构 C1 时间列合同；**不并入**本 VP，本波不激活 |
| **业务域分支** | 不新增 Catalog、订单、支付等业务域 |

## 方向级退出判据

在同时满足下列方向时，本 VP **可以**有界或完整关门（证据必须在工作区目标内）：

1. **分母与契约冻结**：版本身份字段、四种 `runtime.mode` 的横幅/Host/写拒绝投影、诊断摘要字段与排除项形成可机器核对的矩阵；承载面（`I-039-004`）已冻结。
2. **维护提示可感知**：maintenance / degraded / read-only 下，已登录 Admin 有持久横幅；业务写继续被既有门禁拒绝；登录/恢复/邀请放行路径不因横幅而中断（与 VP-012/VP-019 白名单一致）。
3. **版本提示可核对**：授权角色能看到当前版本身份，并到达 R1 冻结的升级说明入口；未授权角色 fail-closed，不泄露内部 commit/模块清单（若分母排除）。
4. **诊断摘要只读**：摘要只聚合已交付探活/就绪/状态字段；缺省无 Prometheus collector 仍可用；不把 VP-015 residual 指标扩进分母。
5. **体验与权限**：中英文、浅色/深色、加载/空态/错误态与既有约定一致；mvp / admin / demo 覆盖有矩阵；不绕过路由守卫。
6. **范围保持**：未实现或解除 Redis / MQ / 多实例 / `timestamptz` / 文件扫描 / 实体检索；未重开 VP-012/015/025；未把 VP-040 混入本 VP。
7. **证据与审计**：退出矩阵、浏览器/自动化回归和必要的独立意见已落盘，开放 required finding = 0，并经用户确认关门。

## 纲领路线图（实现层由 `/govern` 承接）

```text
R1 范围与信息冻结：模式×横幅×错误码矩阵、版本身份与升级入口、诊断字段分母、承载面（I-039-004）
  → R2 维护横幅 + 与 operational gate / Host bootstrap 语义对齐
  → R3 版本提示 + 诊断摘要（复用或有界扩展 system-monitoring）
  → R4 Profile×权限×主题回归、证据矩阵、边界复核与关门
```

纲领阶段通常串行；同一阶段内可并行的实现目标须在 R1 边界冻结后由 `/govern` 按证据创建。

## 信息需求（P-005）

| id | 要回答的问题 | 级别 | 影响门禁 | 最晚阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|--------------|------|----------|----------|------------------|------|-------------|-------------|
| I-039-001 | 版本身份的权威字段与升级说明入口是什么？（`pkg/version`、系统监控行、包/CLI 版本是否同一分母） | required | R1 范围冻结、R3 版本提示 | R1 | 扫描 `pkg/version`、system-monitoring status、发布/changelog 入口；冻结「看得到什么 / 链到哪里」 | **verified** | — | workspace-039 GOAL-002 D-001 §2（QUICKSTART 链接；Shell 版本仅 monitoring.read） |
| I-039-002 | 四种 `runtime.mode` 在横幅、Host bootstrap、写门禁错误码上如何一一投影？read-only 继续映射 Host `degraded` 是否保持？ | required | R1/R2 | R1 | 对照 `bootstrap.go`、`operational.go`、error catalog、前端 feedback-policy | **verified** | — | GOAL-002 D-001 §1：用户裁决 B；生产者 maintenance→Host `degraded`；精确模式 `/me.runtimeMode` |
| I-039-003 | 诊断摘要的字段分母是什么？哪些已有 system-monitoring 字段直接复用，哪些不进首波？ | required | R1/R3 | R1 | 对照 `/healthz` `/readyz` 与 status 行；显式排除 Grafana/Sentry/VP-015 residual 指标 | **verified** | — | GOAL-002 D-001 §3；`r1-diagnostic-field-matrix.md` |
| I-039-004 | 承载面：Shell 横幅 vs 复用 `admin.system-monitoring` vs 新模块？是否进默认集？ | required | 激活、`go` 消费 | 激活前 | 默认候选 = Shell 横幅 + 复用 system-monitoring、不新增模块、不改默认集 | **verified** | 实施期若改新模块进默认集须复核 `go` | 2026-09-19 用户「走流程激活」接受默认候选；VRev-102 |
| I-039-005 | 激活时当前代码候选是否仍满足 Admin 类 freshness 与 VP-008 `go` 消费有效性？ | required | 激活与开区 | 激活前 | `/vision` 复核协议 pin、依赖锁、迁移台账、Profile 默认集、provenance 与区间变更 | **verified** | 下次涉及 Admin 类基线/默认集/协议身份的区间变更时复核 | `7e5ce891` → `6197e802` 五域 PASS（区间 = VP-038 已审结目 + W32–W34）；VRev-102 |
| I-039-006 | 是否需要管理员在运行中切换 `runtime.mode`（相对重启生效）？ | non-blocking | 不进首波 | — | 热切换 = 配置热加载类能力，与 VP-016/025 红线同类 | deferred | 真实运维需求出现时由 `/vision` 复核 | 首波不承诺热切换 |

## 工作区绑定

| workspace_id | root_goal | role | joined | notes |
|--------------|-----------|------|--------|-------|
| workspace-039-version-maintenance-diagnostics | GOAL-001-version-maintenance-diagnostics | delivery | 2026-09-19 | `/govern` scaffold；Root 初始 `active · 0/4`；不改变 Charter primary |

## 关门记录

（仅 `closed` / `abandoned` 时填写。）

| date | outcome | summary | evidence_links | residuals |
|------|---------|---------|----------------|-----------|
| — | — | — | — | — |

## 规划修订短史

| date | change |
|------|--------|
| 2026-09-19 | 初创 `planned` v0.1.0 · 0 区。用户确认选项 1：下一拍 = 本 VP；C1 另立 VP-040 且本波不激活；不塞 VP-010、不合成混分支 VP。计划阶段 self = [VRev-101](../reviews/VRev-101-vp039-vp040-planned.md)。 |
| 2026-09-19 | 用户指令走流程激活：`planned → active` v0.2.0。`I-039-004`/`I-039-005` verified；激活 self = [VRev-102](../reviews/VRev-102-vp039-activation.md) `pass`；lead `workspace-039-version-maintenance-diagnostics` 交 `/govern` 开区。VP-040 保持 planned 停放。 |

## 声明

本文件是已确认的 Vision Plan 意图，不是 Goal 五件套、实现事实或 progress 权威。激活、工作区创建、阶段执行和 Goal 审计分别由 `/vision` 与 `/govern` 按对齐契约承接。
