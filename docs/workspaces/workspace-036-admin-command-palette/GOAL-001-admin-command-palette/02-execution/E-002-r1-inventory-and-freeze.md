---
doc_type: goal-execution
id: E-002-r1-inventory-and-freeze
status: recorded
created: 2026-09-14
updated: 2026-09-14
parent: GOAL-001-admin-command-palette
version: 0.1.0
---

# E-002 · R1 分母盘点与冻结（2026-09-14）

## 事实

- 已复核工作区 `workspace.md`、Root 五件套、Charter `@0.4.0`、alignment、workspace protocol、VP-036 与 Vision Review；工作区绑定与 Vision open required 均满足，资料目录为 `none`，本轮没有共享资料引用。
- 已扫描 `apps/web/src/protocol/app-manifest.ts`、`apps/web/src/app/navigation.ts`、`apps/api/internal/manifest/manifest.go`、`apps/api/internal/composition/composition.go`、`apps/api/kernel/profile.go`、各模块 Manifest fragment 与 Schema。确认 pinned AppManifest 不承载 module/profile/provider/action 字段，不能把 provider 元数据静默写入公共 Manifest。
- 已形成 `attachments/r1-searchable-item-matrix.md`：记录 mvp/admin/demo/custom 页面、导航、顶层 action 定义与可收录直接触发器分母，以及可见导航/通知 Shell 入口规则。矩阵事实来源为代码与 Schema，不把已读的其他工作区台账当作本区事实。
- 用户 2026-09-14 书面确认 R1 推荐的分母、provider v1、排序/去重/上限、快捷键、ARIA/focus、Profile 证据策略；已写入 `01-decision/D-002-r1-scope-contract-ux-freeze.md`。
- 已识别实现前的安全缺口：当前页面级 `SchemaCrudProvider.invokeAction` 的 modal/navigate/custom 分支需统一执行权限复核，且 actionButton 的 node id 目标必须传给 executor；该缺口未被当作已完成，转入 R2/R3 修正与测试。

## 信息台账事实

- I-036-001：矩阵已回答页面/导航/动作分母与四 Profile 候选覆盖；状态可更新为 `verified`（用户决策 + 代码盘点）。
- I-036-002：现有 permission/Profile/route/action executor 语义与程序化调用安全要求已冻结；状态可更新为 `verified`（用户决策 + 现有实现证据），执行安全修正仍是后续成功标准证据。
- I-036-003：字段/关键词、排序/去重、12 条上限、快捷键、ARIA/focus、加载/空态/错误态与双语/主题口径已冻结；状态可更新为 `verified`（用户决策），浏览器可用性仍待 R3/R4 验证。
- I-036-004、I-036-006 保持 `verified`；I-036-005 保持 `deferred` / `non-blocking`，不进入首波。

## 阻塞 / 风险

- R1 方案冻结不再有开放 required 信息项；实现尚未开始，不能据此宣称 R2/R3/R4 完成。
- R2/R3 必须保留后端鉴权、既有路由与 Profile 语义，并把已识别的 programmatic gating 缺口纳入测试；若实现方式需要改变 action 类型、参数或 Manifest 协议，须回到用户裁决点。

## 下一步（计划）

- 实现 provider v1 的类型、聚合/冲突、匹配/排序/12 条裁剪与内置 Manifest provider，并补纯函数测试。
