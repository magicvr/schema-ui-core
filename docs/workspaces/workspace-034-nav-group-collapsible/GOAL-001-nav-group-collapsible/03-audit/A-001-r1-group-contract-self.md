---
id: A-001-r1-group-contract-self
doc: audit-entry
parent_goal: GOAL-001-nav-group-collapsible
source: self
auditor: /govern（schema-ui-core 编排器）
type: design-plan
scope: R1 分组基线与 R2 契约优先方案
date: 2026-09-07
verdict: pass
created: 2026-09-07
updated: 2026-09-07
version: 0.1.0
---

# A-001 · R1 分组基线与 R2 契约优先方案自审

## 范围与区间

本自审覆盖用户确认后的 R1 产品基线、跨模块分组契约提案，以及从现有代码到 R2 实施的边界。尚未把任何代码实现或 R2/R3 检查点写成已完成事实。

## 成果（有证据）

1. 用户确认了五组顺序、组内顺序、Dashboard 顶层单例、Examples 保留和 top/user slot 边界，已记录于 [D-002](../01-decision/D-002-r1-baseline-and-group-contract.md)。
2. 现有 Web 协议已经接受 `NavGroup`，但 `apps/web/src/app/App.tsx` 只静态渲染组；现有 API `NavigationContribution` 与 Manifest 聚合尚未提供跨模块 group 归一化。
3. D-002 选择以结构化导航贡献作为分组语义输入，并保留无 group 的平铺兼容路径；不引入 composition 中央业务映射。
4. 相关 Vision required = 0；Goal 当前无历史 A 条目或开放 required finding；I-034-002 已有用户决策证据，但实现/回归证据仍待后续阶段。

## 对照成功标准

| 项 | 自审结论 |
|---|---|
| R1 产品基线可判定 | pass：D-002 已冻结；行为证据仍待实现验证 |
| 模块可声明可选 group | pass：D-002 明确结构化契约方向；代码尚未实施 |
| 跨模块共组不依赖中央业务注册 | pass：D-002 明确由 Manifest 聚合归一化 |
| 未声明 group 向后兼容 | pass：D-002 保留平铺语义 |
| 不误改 top/user slot | pass：D-002 明确只归一化 sidebar 普通链接 |

## Findings

本轮无 required 或 recommended finding。该结论只表示方案在进入独立审计前没有发现已知阻断，不表示 R2 代码或 R3/R4 行为已经完成。

## 审计模式

本 scope 跨 API kernel、Manifest 聚合与 Web Shell，并涉及模块兼容契约，按 `cross` 处理：本条 self 之后调用项目指定的本地 grok build（grok-4.6 · high）执行 independent design-plan audit。independent 意见不改目标状态；响应由 `/govern` 合并。

## 结论与下一步

**verdict: `pass`**（可交独立审计；在 independent 结果合并前不宣称 R2 冻结或放行）。

下一步：由本地 grok build 以 `/audit` 对 D-002、现有 kernel/Manifest/Web 边界和 I-034-001/002 进行独立复核；之后由 `/govern` 响应全部意见，再进入 R2 实施。

## 声明

本意见为 `source: self`，不修改目标 status/progress，不关闭任何 required 信息项，不替代 independent 意见。
