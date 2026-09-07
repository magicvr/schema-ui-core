---
id: GOAL-001-nav-group-collapsible
doc: decision
status: active
parent: null
created: 2026-09-07
updated: 2026-09-07
version: 0.1.0
---

# 决策记录 · GOAL-001-nav-group-collapsible

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|------------------|------|-------------|-------------|
| I-034-001 | required | 当前已注册 sidebar/top/user 导航的清单、slot 与 Profile 覆盖 | R1 分母 / R4 回归 | R1 | 代码盘点 + 运行时 manifest 矩阵 | collecting（代码盘点完成） | R1 补核 | VP-034 §现有注册导航 |
| I-034-002 | required | 初始五组的标题、组内顺序、Dashboard 顶层单例是否符合使用语义 | R1 方案冻结 | R1 | Shell/UI 复核并留痕 | collecting | R1 结束前复核 | VP-034 初始基线 |
| I-034-003 | non-blocking | 折叠状态的会话/浏览器持久化选择 | R3 | R2 | 实现阶段决策与交互测试 | open | 不阻断 R1/R2 | 待确认 |
| I-034-004 | required | 直接 URL → 所属分组自动展开覆盖矩阵 | R3/R4 | R3 | route matrix + e2e | open | — | 待确认 |
| I-034-005 | required | optional/custom/demo profile 的跨模块共组聚合 | R2/R4 | R4 | profile manifest harness | open | — | 待确认 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-09-07 | VP-034 scope correction 与 R1-R5 实施路线 | accepted | `01-decision/D-001-vp034-scope-and-roadmap.md` |

## 当前方案边界

- 当前已注册 sidebar 导航不因属于既有模块而排除；必须在 R4 完成迁移与验证。
- 分组按产品语义跨模块组织，不按模块一对一创建组。
- Dashboard 是顶层主入口单例；top/user slot 不搬到 sidebar，但纳入兼容回归。
- R1 在实现前冻结最终标题、顺序和组 key；VP 中的五组是初始基线，不把未复核内容伪装成最终事实。
