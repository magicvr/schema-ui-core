---
id: D-001-vp034-scope-and-roadmap
doc: decision-entry
parent_goal: GOAL-001-nav-group-collapsible
source: user-confirmed + /vision
status: accepted
date: 2026-09-07
created: 2026-09-07
updated: 2026-09-07
version: 0.1.0
---

# D-001 · VP-034 scope correction 与 R1-R5 实施路线

## 决策

用户明确修正 VP-034：当前已经注册的导航不能因属于既有模块而排除在目标外，需要将其按产品语义合理分组。随后交 `/govern` 开设 delivery workspace。

VP-034 已更新为 v0.3.0，当前有效范围是：

- 当前已注册 sidebar 节点全部纳入分组迁移与验证；明确的 Dashboard 顶层单例可作为有意例外，但不能排除在回归外。
- 分组跨模块解耦，允许不同模块向同一组注册；不采用一个模块一个组。
- 初始基线为 `identity-access`、`content-data`、`operations`、`communications`、`commerce` 五组。
- `dev.examples` 的既有 Examples 组保留。
- Settings、Account、My wallet、通知等 top/user slot 保留原 slot，不搬到 sidebar，但必须做兼容回归。
- 直接 URL 进入已分组页面时所属组自动展开。

## 依据与边界

- VP：`docs/vision/plans/VP-034-nav-group-collapsible.md` v0.3.0
- Scope correction Vision Review：`docs/vision/reviews/VRev-084-vp034-existing-navigation-scope-correction.md`，self `pass`，0 required
- 当前代码盘点基线：HEAD `f2044cf3`
- 不修改 Charter、VP-008 `go` 语义或 VP-010 长期程序边界。

## 实施路线

R1 先冻结导航清单、Profile/slot 矩阵、分组标题与顺序；R2 冻结注册/聚合契约；R3 实现 Shell 交互与直接 URL 自动展开；R4 完成现有 sidebar 全量迁移并验证 optional/custom/demo；R5 完成证据、回归和 Goal 审计。

## 未决信息

R1 尚需在真实 UI 上复核初始组名、组内顺序和 Dashboard 单例呈现；该项保持 `I-034-002 = collecting`，不被本决策记录伪装为最终实现事实。