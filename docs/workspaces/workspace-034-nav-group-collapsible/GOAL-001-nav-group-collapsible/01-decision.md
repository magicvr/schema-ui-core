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
| I-034-001 | required | 当前已注册 sidebar/top/user 导航的清单、slot 与 Profile 覆盖 | R1 分母 / R4 回归 | R1 | 代码盘点 + 运行时 manifest 矩阵 | verified（静态 R1 分母） | R4 运行时 Manifest/Profile harness 核验 | [R1 矩阵](attachments/r1-navigation-profile-slot-matrix.md)；E-002 |
| I-034-002 | required | 初始五组的标题、组内顺序、Dashboard 顶层单例是否符合使用语义 | R1 方案冻结 | R1 | 用户确认基线；实现后的 Manifest/Shell 行为仍需验证 | verified（决策） | R1 代码/回归证据继续核对 | D-002；VP-034 初始基线 |
| I-034-003 | non-blocking | 折叠状态的会话/浏览器持久化选择 | R3 | R3 | 用户确认 sessionStorage；实现阶段补交互测试 | verified（用户决策） | R3 以测试核对读写与容错 | D-005 |
| I-034-004 | required | 直接 URL → 所属分组自动展开覆盖矩阵（含已登记内页/动态路径） | R3/R4 | R3 | route projection matrix + e2e；覆盖父级页面关系 | verified（R3/R4 矩阵） | R5 关门前再核对 | D-005；E-008；`apps/web/src/app/nav-groups-r4.test.ts` |
| I-034-005 | required | optional/custom/demo profile 的跨模块共组聚合 | R2/R4 | R4 | profile manifest harness | verified（R4 runtime matrix） | R5 关门前再核对 | E-008；`apps/api/internal/composition/nav_group_r4_test.go` |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-09-07 | VP-034 scope correction 与 R1-R5 实施路线 | accepted | `01-decision/D-001-vp034-scope-and-roadmap.md` |
| D-002 | 2026-09-07 | R1 分组基线与契约优先方案 | accepted | `01-decision/D-002-r1-baseline-and-group-contract.md` |
| D-003 | 2026-09-07 | 响应 A-002：R1 证据闭合与 R2 开放门禁 | accepted | `01-decision/D-003-a002-response-and-open-gates.md` |
| D-004 | 2026-09-07 | R2 组序、混排与冲突契约 | accepted | `01-decision/D-004-r2-group-order-and-conflict-contract.md` |
| D-005 | 2026-09-07 | R3 折叠状态与直接 URL 范围 | accepted | `01-decision/D-005-r3-collapse-and-deep-link-scope.md` |

## 当前方案边界

- 当前已注册 sidebar 导航不因属于既有模块而排除；必须在 R4 完成迁移与验证。
- 分组按产品语义跨模块组织，不按模块一对一创建组。
- Dashboard 是顶层主入口单例；top/user slot 不搬到 sidebar，但纳入兼容回归。
- R1 已按用户确认冻结五组 key、组顺序、组内顺序与 Dashboard 例外；Manifest/Shell 的真实行为仍必须以代码与测试证据核对。
- R2 采用契约优先：结构化 `NavigationContribution` 提供可选 group 元数据，Manifest 聚合归一化跨模块共组；未声明 group 继续平铺。
- 同一 group key 的跨模块元数据不一致时必须 fail closed；不得由 composition 维护 NodeID → group 的中央业务映射。
- A-002 independent 的 F-002～F-004 已由 D-004/A-004 合法 fixed；R2 代码已实施并通过 A-005 self + A-006 grok independent，A-007 已响应 recommended。
- R3 按 D-005 采用 sessionStorage；直接 URL 自动展开覆盖 sidebar 直接 route 以及已登记内页/动态路径，实际行为仍待 R3 测试证据。
