---
id: GOAL-013-w12-product-surface-intent
title: W12 · 产品面交互意图对齐（顶栏菜单 / 列表搜索 / 个人中心 / 我的钱包 / 回收站时间 / 模块配置）
status: active
parent: GOAL-001-design-implementation-conformance
created: 2026-08-16
updated: 2026-08-16
version: 0.7.0
progress: 2/4
---

# GOAL-013 · W12 · 产品面交互意图对齐

VP-010 / workspace-010 的**第十二波**（用户 2026-08-16 `/govern` 点名立项）：把六条产品面意图对照 as-built，先对齐再按阶段实施。问题清单与现状证据见 [01-decision.md](01-decision.md) 与 [01-decision/D-001-intent-inventory.md](01-decision/D-001-intent-inventory.md)。

## 当前边界

- **范围（本波实施）**：T-01 顶栏用户下拉；T-02 列表搜索；T-03 个人中心 Tabs；T-05 回收站删除时间；T-06 模块启用只认 YAML。
- **移交**：T-04「我的钱包」不在本波实施，见 D-005 → [GOAL-022-my-wallet-self-service](../../workspace-011-admin-functional-modules/GOAL-022-my-wallet-self-service/00-meta.md)（Q2）。
- **非范围**：不改 Charter `@0.2.0`；不重开 VP-005 全量视觉基线；不把本波写成 VP-010 关门；不默认改 Profile 生产默认集（T-06 若触及则走 go 判定）。

## 成功标准与路线图（P-001）

- [x] **S1 · 意图盘点**：六条意图对照 as-built 落盘；信息项登记；用户确认范围与分批（D-001～D-008，2026-08-16）
- [x] **S2 · 方案冻结**：逐项设计决策（D-002～D-008）；required I-001～I-004、I-006 已闭合
- [ ] **S3 · 按冻结范围实施**：可分批（建议 P0 = T-05 + T-01；P1 = T-03 + T-02；P2 = T-06。T-04 已移交 workspace-011）
- [ ] **S4 · 验证与关门**：回归绿 + 自审（T-04/T-06 若改数据面或 Profile/模块语义则加 independent / go 判定）

progress: 由四个等权检查点派生（S1～S4）；当前 **2/4**。

## 审计策略

| 阶段 / 项 | 默认模式 | 说明 |
|-----------|----------|------|
| S1 盘点 | `none` | 只读落盘 |
| T-01 / T-03 / T-05 | `self` | 壳层与呈现，可逆 |
| T-02 | `self` | 复用已有 form 控件白名单；不扩协议则不必 independent |
| T-04 | `independent`（若新增自服务资金面）或 `self`（若只是只读包装既有账户 API） | S2 按范围裁定；不得静默降级 |
| T-06 | `self`；若改 Profile 默认集 / 模块矩阵 / Manifest 装配 → 加 go 判定，必要时 `independent` | 与 W7 / VP-008 `go` 接口 |

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 各列表页应暴露哪些搜索/筛选字段（字段、控件类型、后端是否已支持） | S2 冻结 T-02 | S2 | 用户选择题 | **verified** | — | D-003 矩阵 |
| I-002 | required | 个人中心选项卡信息架构（资料 / 安全 / 会话 / MFA 如何分组） | S2 冻结 T-03 | S2 | 用户选择题 | **verified** | — | D-004：资料 / 安全 / 会话 |
| I-003 | required | 「我的钱包」归属与范围：本区壳层+页面 vs 回流 workspace-011；只读 vs 自助流水/充值口径 | S2 冻结 T-04 | S2 | 用户书面 Other | **verified** | — | D-005：本波不做；[GOAL-022](../../workspace-011-admin-functional-modules/GOAL-022-my-wallet-self-service/00-meta.md) |
| I-004 | required | T-06 是「发现/文档/YAML 列表形态」还是「改变 Profile 默认启用集」 | S2 冻结 T-06 | S2 | 用户书面 Other | **verified** | — | D-007：启用集只认 YAML；不改三档成员 |
| I-006 | required | 「不再用 env/.env」是仅针对模块启用，还是全局废除（含 JWT 等密钥插值） | S2/S3 T-06 | S3 改 compose/加载器前 | 用户选择题 | **verified** | — | D-008：只取消模块启用 env |
| I-005 | non-blocking | 顶栏触发器形态（头像+姓名 vs 仅头像；移动端是否仍保留抽屉内用户链） | S2 T-01 | S2 | 用户选择题 | **verified** | — | D-002：头像+姓名；抽屉不重复用户链 |

## 父目标

- [GOAL-001-design-implementation-conformance](../GOAL-001-design-implementation-conformance/00-meta.md)

## 台账布局

- `01-decision/`：D-NNN；`02-execution/`：E-NNN；`03-audit/`：A-NNN。
- 跨区引用用 Q2 路径（workspace-protocol §2.6）。
