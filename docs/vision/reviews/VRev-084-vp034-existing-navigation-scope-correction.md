---
id: VRev-084
doc_type: vision-review
title: VP-034 既有导航纳入范围修正审视
source: self
scope: VP-034-nav-group-collapsible · v0.3.0 scope correction
verdict: pass
date: 2026-09-07
auditor: /vision
created: 2026-09-07
updated: 2026-09-07
parent: null
version: 0.1.0
---

# VRev-084 · VP-034 既有导航纳入范围修正审视

## 触发与范围

用户对已激活的 VP-034 作出明确范围修正：当前已经注册的导航不能因属于既有模块而排除在本 VP 外，应按产品语义完成合理分组，然后交 `/govern` 开设工作区。

本次审视覆盖：

- VP-034 v0.3.0 scope correction 是否仍在 Charter 边界内
- 当前已注册 sidebar 导航是否已登记为明确的执行覆盖范围
- 分组是否保持与模块解耦、支持跨模块共组
- 既有 top/user slot 是否被错误纳入或静默丢失
- 是否影响 VP-008 `go`、Charter 或 VP-010 的对齐门禁

## 修正事实

VP-034 由 v0.2.0 更新为 v0.3.0，保留 `status: active`、`vision_ref: schema-ui-core-admin-foundation@0.4.0` 和 lead workspace 绑定。

修正后的有效范围：

1. 当前已注册 sidebar 节点必须全部纳入迁移与验证；每个节点要么归入合理分组，要么以有明确 UX 理由的顶层单例例外留痕。
2. 初始分组基线为：
   - `identity-access`：`menu_users`、`menu_roles`、`menu_data_permission`
   - `content-data`：`menu_files`、`menu_dictionary`
   - `operations`：`menu_activity`、`menu_monitoring`、`menu_scheduled_tasks`、`menu_recycle_bin`
   - `communications`：`menu_mail`、`menu_mail_outbox`、`menu_telegram`
   - `commerce`：`menu_wallet`、`menu_wallet_vouchers`、`menu_digitaloffer_offers`、`menu_digitaloffer_entitlements`、`menu_digitaloffer_purchases`
   - `menu_dashboard`：保留为左侧顶层主入口单例，不是排除项。
3. `dev.examples` 既有 `Examples` 组保留并纳入回归。
4. `menu_settings`、`menu_account`、`menu_wallet_self` 和通知面继续保留原 top/user slot 语义；它们不搬到 sidebar，但会被纳入 slot 回归，避免被分组逻辑丢失。
5. 未来模块仍可选择不声明 `group`；本次修正针对当前已注册 sidebar，不把“已有模块”当作排除理由。

当前代码盘点依据：

- `apps/api/kernel/profile.go` 的编译候选、Profile 默认集与 Navigation keys
- 各模块 `apps/api/modules/**/provider.go` 的 NavigationContribution
- 各模块 `apps/api/modules/**/manifest/fragment.json` 的 slot 与 navigation 声明
- 当前候选 HEAD：`f2044cf3`

## Charter / 组合对齐

| 项 | 结论 |
|----|------|
| Charter 版本 | `schema-ui-core-admin-foundation@0.4.0` |
| Charter 成功边界 | #3 产品化 Admin 体验；#5 模块贡献导航且 Shell 不维护模块中央注册，均直接覆盖 |
| Charter 修改 | **无**；不改变目的、边界或非目标 |
| VP-010 关系 | **正交**；本次是 Admin 功能增强，不是 as-designed/as-built 偏差整改，不创建 W31 |
| VP-008 `go` | **不受影响**；仅改变 Admin 功能 VP 的交付分母，不改 Profile 默认集、协议 pin、迁移或 go 适用性 |
| Vision Review required | **0** |

## 可判定性与门禁

新增/修正后的退出判据仍可验证，并补强为既有导航覆盖：

- Shell group 渲染与折叠/展开可访问性
- 可选 `group` 注册字段与跨模块共组聚合
- 当前已注册 sidebar 节点的分组迁移、顶层单例例外和 demo 既有组
- 直接 URL 进入时自动展开所属分组
- top/user slot 兼容与 profile/custom 回归
- playbook 对 group key、组内顺序、共组、未分组例外和迁移规范的说明

新增 P-005 信息项已写入 VP-034：I-034-001～005。当前没有阻断工作区 scaffold 的 Vision required；I-034-002（精确组名/顺序）和 I-034-004/005 在实现阶段按 Root R1～R4 的目标门禁关闭，不能提前写成已验证事实。

## 结论

**verdict: `pass`**（0 required；本轮无新 required finding）。

VP-034 v0.3.0 scope correction 合法、与现行 Charter/VP-008 对齐、不会触发 strategic re-align 或暂停 `go`。可以按用户指令交 `/govern` 开设 `workspace-034-nav-group-collapsible` 并建立 Root 五件套。

VRev-082 的 V-F121 recommended 继续保留为 open recommended；本次已将其要求（group key 命名空间）纳入 VP-034 退出判据与 P-005/R1 工作范围，不阻断开区。

## 历史记录声明

VRev-082 与 VRev-083 保留其原始审视时间点、原始 scope 和 verdict，不改写历史；本 VRev-084 是对现行 VP-034 交付范围的后续修正记录。
