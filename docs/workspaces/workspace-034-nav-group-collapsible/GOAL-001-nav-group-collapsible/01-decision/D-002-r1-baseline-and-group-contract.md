---
id: D-002-r1-baseline-and-group-contract
doc: decision-entry
parent_goal: GOAL-001-nav-group-collapsible
source: user-confirmed + /govern
status: accepted
date: 2026-09-07
created: 2026-09-07
updated: 2026-09-07
version: 0.1.0
---

# D-002 · R1 分组基线与契约优先方案

## 决策

用户确认：

1. 冻结 VP-034 当前提出的产品分组基线，作为 R1 的实现基线：
   - 分组顺序：`identity-access` → `content-data` → `operations` → `communications` → `commerce`。
   - 组内顺序采用 VP-034 §现有注册导航的表格顺序。
   - `menu_dashboard` 保持 sidebar 顶层主入口单例，但仍在迁移与回归分母内。
   - `dev.examples` 的既有 `Examples` 组保留，demo 回归不把它改写为 Admin 五组之一。
   - `menu_settings`、`menu_account`、`menu_wallet_self` 与通知面继续保留原 top/user slot 语义。
2. R2 采用**契约优先**实现：结构化 `NavigationContribution` 是分组语义的权威输入；Manifest 聚合负责把跨模块分组归一化为协议 `NavGroup`；未声明 `group` 的导航继续平铺。

## 方案边界（R2 实现前的冻结提案）

### 结构化贡献

在 `apps/api/kernel` 为导航贡献增加可选的分组元数据：

- `Group`：可选的具名分组描述；未设置表示保持当前平铺语义。
- 分组描述包含稳定 `Key`、组顺序、`Label`/`LabelKey` 至少一个，以及可选语义图标。
- 同一 `Key` 的跨模块声明必须使用一致的组元数据；不一致时聚合 fail closed。
- 分组 key 使用小写英文短横线命名，并作为跨模块共享契约标识；不以模块 id 自动加前缀。

### Manifest 归一化

- 现有模块 Manifest 仍可贡献单个导航链接，不要求模块维护中心注册表。
- composition 将已验证的结构化导航贡献传给 Manifest 聚合器。
- 聚合器仅对 sidebar 普通链接按 `NodeID → Group.Key` 归一化；top/user slot 不受影响。
- 聚合器按冻结的组顺序输出分组，组内保持现行 NodeID 排序；未分组条目保持兼容平铺。
- Manifest 中已有协议组（如 `dev.examples`）继续保留；不与结构化 Admin 分组隐式合并。

### Shell 行为

- Web 投影继续消费协议 `NavGroup`，不在 Shell 中按模块 id 做业务判断。
- R3 再实现分组折叠/展开、键盘交互、直接 URL 自动展开与状态保持；本决策不提前宣称这些行为已实现。

## 未选方案

- **Manifest 归一化优先**：只修改 Manifest、保留结构化贡献现状，会造成系统数据/Manifest 两套分组语义，未采纳。
- **组合根集中映射**：在 composition 维护 NodeID → group 中央业务映射，违反模块解耦与无中央注册分支的约束，未采纳。

## 风险与后续审计

- 这是跨 API kernel、Manifest 聚合与 Web Shell 的兼容性契约变更，实施前按 `cross` 模式执行 self + 本地 grok build（grok-4.6 · high）独立审计。
- I-034-002 的产品决定已由本记录闭合；真实 Manifest/Profile/UI 行为仍须在 R1/R3/R4 以代码与测试证据核对。
- I-034-003（折叠状态存储）继续作为 R3 前的 non-blocking 方案项，不影响本 R1/R2 契约决定。

## 证据

- 用户本轮对 `r1-group-baseline` 与 `group-contract-architecture` 的明确选择。
- `docs/vision/plans/VP-034-nav-group-collapsible.md` v0.3.0。
- `apps/api/kernel/contribution.go`、`apps/api/internal/manifest/manifest.go`、`apps/web/src/app/navigation.ts`、`apps/web/src/app/App.tsx` 的 R1 代码盘点。
