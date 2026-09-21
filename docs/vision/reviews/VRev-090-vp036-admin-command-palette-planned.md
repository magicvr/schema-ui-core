---
id: VRev-090-vp036-admin-command-palette-planned
doc_type: vision-review
title: VP-036 Admin 全局检索与 Command Palette · 计划阶段意图审视
source: self
scope: VP-036-admin-command-palette · planned
verdict: pass
open_required: 0
status: recorded
date: 2026-09-10
auditor: /vision
created: 2026-09-10
updated: 2026-09-10
parent: null
version: 0.1.0
---

# VRev-090 · VP-036 计划阶段意图审视

## 审视范围

新立 `VP-036-admin-command-palette`（Admin 功能分支 · 全局检索与 Command Palette）的计划阶段意图审视：

- Charter `@0.4.0` 对齐
- 结构选型（新 VP vs VP-010 波次）
- 首波范围与退出判据可判定性
- P-005 信息需求与激活门禁
- 与 VP-034、VP-009、VP-010、RT-X01/RT-X02 及业务域分支的边界

## 审视结论

**verdict: `pass`**（0 required，1 recommended）。

VP-036 的意图落在现行 Charter 成功边界 #3（产品化 Admin 体验）与 #5（模块可组合、避免 Shell 中央注册路径）内。它是向好演进的用户产品能力，不是 VP-010 的符合性整改；新 VP + 新 delivery 工作区的结构选择成立。首波将全局检索收窄为已注册页面、导航项和声明式动作，不预制实体全文索引或专用搜索基础设施。

本 Review 仅确认 `planned` 意图，不授权激活、开工作区或进入实现。

## Charter 对齐

| 项 | 检查结果 |
|----|---------|
| `vision_ref` = `schema-ui-core-admin-foundation@0.4.0` | ✅ 精确匹配唯一 active Charter |
| 落在成功边界内 | ✅ Admin 产品化体验与模块组合能力直接覆盖 |
| 不改变 Charter 目的/边界/非目标 | ✅ 不改 Charter，不新增业务域，不改变单主线模块策略 |
| 与现有路线图一致 | ✅ 承接 Admin 功能分支“体验增强”下一拍；不解除架构 gated 行 |

## 结构选型

| 问题 | 判断 | 依据 |
|------|------|------|
| 改 Charter 目的/边界？ | 否 | 是现有 Admin 体验方向内的有界增量 |
| 同愿景新纲领波次？ | 是 | 路线图已登记的体验增强下一拍 |
| 是否属于 VP-010 普通波次？ | 否 | 不是修复既有设计意图与 as-built 的偏差，而是新增发现能力 |
| 是否需要独立 Goal 树？ | 是 | 需独立冻结检索契约、权限/Profile 矩阵与跨模块回归 |
| 结论 | **新 VP + 新 delivery 工作区** | 计划阶段 0 区；激活后交 `/govern` scaffold |

## 首波边界与退出判据

- 可核验分母是当前已注册的页面、导航项和声明式动作；实体记录全文搜索不进入首波退出分母。
- `SearchableItem` / provider 聚合、稳定排序/去重、Profile/权限过滤、直接 URL 与导航分组联动均可形成证据矩阵。
- 快捷键、键盘焦点、ARIA、i18n、theme、加载/空态/错误态可由浏览器与自动化回归验证。
- Saved Views、批量结果中心、未保存保护、Toast 全局重做与新业务域明确排除，避免把多个可独立交付块塞入一个 VP。
- `RT-X01` / `RT-X02`、Redis、MQ、多实例仍保持 gated；若后续实体级搜索产生真实规模或查询需求，须另行 `/vision` 复核。

## P-005 信息就绪

| 信息项 | 级别 | 状态 | 门禁 |
|--------|------|------|------|
| I-036-001：页面/导航/动作精确分母与 Profile 覆盖 | required | collecting | 阻断 R1 范围冻结与 R2 聚合 |
| I-036-002：权限、Profile、直接 URL 与动作守卫语义 | required | collecting | 阻断 R1 方案冻结与 R3 实施 |
| I-036-003：字段、排序/去重、键盘/ARIA 与焦点语义 | required | collecting | 阻断 R1 方案冻结与 R3 实施 |
| I-036-004：实体搜索与 `RT-X01` / `RT-X02` 是否进入首波 | required | **verified (user decision)** | 首波不承诺实体全文搜索；保持 gated |
| I-036-005：最近搜索/固定项/持久化偏好 | non-blocking | deferred | 不影响首波，可在后续 UX VP 复核 |
| I-036-006：激活前 Admin 类 freshness | required | open | 阻断激活，不阻断 planned 登记 |

无 required 信息项阻断 `planned` 立项；激活前必须完成 I-036-006，R1 冻结前必须完成 I-036-001～003。

## Findings

### V-F123（recommended · 非阻断）

在激活或 R1 方案冻结前，建议把 `SearchableItem` 分母、结果排序/去重规则与权限/Profile 过滤写成一张可机器核对的矩阵；并明确“页面/导航/声明式动作检索”与“业务实体全文搜索”的边界。这样可以避免“全局搜索”在执行阶段无意扩大为需要 `RT-X01` / `RT-X02` 的第二个基础设施项目。

状态：`open · recommended`。不阻断 `planned` 登记；应由激活包或 R1 决策记录承接。

## 激活门禁（后续 `/vision`）

1. Admin 类 freshness review：核对协议 pin、依赖锁、迁移台账、Profile 默认集、Manifest/provenance 与区间变更；不满足时不得消费 VP-008 `go`。
2. 激活就绪 self Review：确认 I-036-006、slug 与工作区绑定边界。
3. 用户确认 delivery workspace slug 后，交 `/govern` 创建 workspace + Root；本 Review 不创建 Goal 五件套。

## 声明

- 本 Review source = `self`，不冒充 independent。
- 不改变 Charter、其它 VP、工作区或 Goal status/progress。
- open required = 0；V-F123 为 recommended，不阻断 `planned` 登记或后续激活准备。
- 下一步：若继续推进，使用 `/vision` 激活 VP-036；激活后交 `/govern` scaffold 工作区。
