---
id: E-008
doc: execution-entry
goal: GOAL-001-iam-recovery
status: recorded
created: 2026-08-25
updated: 2026-08-25
version: 1.0.0
---

# E-008 · 关后维护：邀请管理 UX 修正实录（Root done 之后）

> 本条目记录 Root（`done 4/4`，2026-08-25）关闭后，用户就邀请管理提出的连续 UX 修正。属维护类工作，未新建目标；全部按用户文字指示逐条落实并在此留痕。

## 时间线与交付（git 均在本工作区 repo）

| # | 用户指示 | 交付 | commit |
|---|----------|------|--------|
| 1 | 角色不要填文本，改下拉多选 | 角色输入重构为 `RoleMultiSelect`（trigger 徽章 + 勾选弹出层 + 外部点击关闭；`/api/roles` 目录加载；全不选禁用创建；表格角色列显示角色名）+ 3 组组件测试 | `6e3c97a5` |
| 2 | 邀请管理作为用户管理内页 + 入口按键 ／ 请沿仓库父子页先例评估 | 迁入子页 `/users-invites`（`navigate` 入口 + `users.invite` 权限显隐；第二 PageContribution + manifest fragment）；**顺修潜伏缺陷**：`createUser.bodyMapping` 缺 `roles` 映射（新建用户表单的角色选择此前未送达后端） | `b9ae434f` |
| 3 | 列表与操作分开，记录可筛选 | 拆「发起邀请/邀请记录」双卡；后端 `ListInvites` 增 `status` 过滤（pending/consumed/revoked/expired，SQL WHERE + 域测试）；前端筛选下拉联动 `?status=` | `930f0513` |
| 4 | 列表没有分页？ | 服务端分页接入：`page/pageSize` 参数 + `‹ x / y ›` 导航 + 10/20/50 条数切换；筛选/条数切换回第 1 页；创建成功回第 1 页刷新 | `cab06dbb` |
| 5 | 面包屑缺上级页导航；列表用「正常页面」样式（引擎自带筛选分页，不要自创样式） | 面包屑 `BREADCRUMB_PAGE_PARENTS` 补 `users-invites → users`；页面全面 schema 化（search 表单 status 筛选 + schema-table 分页 + 行级 revoke confirm）；操作卡收口 `invite-issue-card`（旧自绘列表退役） | `3d39b3e0` |
| 6 | 已撤销记录操作列仍显示撤销，不符合预期 | 后端列表响应派生 `revokable` 标志（仅 pending 真）；行级 revoke 加 `disabledWhen{field:revokable,equals:false}`；+四态 handler 单测 | `f58ad8f7` |
| 7 | 重发集成到操作列（免手动粘贴 id）；选中行应显示明细 | 中间迭代：行级 custom handler（`invites.resend`）剪贴板披露（`ceddcb2b` + recordView 明细抽屉）→ 用户提示应以协议为准查证 → 核实 `docs/schemas/action.schema.json`（`ModalAction.content` 任意 Node）后**回退剪贴板方案**，改协议模态：行级 `type:"modal"` + content custom `invite-resend-dialog`；渲染器为 modal content 注入 `modalRow`；弹窗内重发→内联披露新链接→复制→完成（关闭+刷新）；`recordView` 明细（id/角色/邮箱/状态/邀请人/过期/创建） | `ceddcb2b` → `43ce4f84` |

## 治理要点

- **判断修正**：#7 中间采用「custom handler + 剪贴板」绕过 schema 披露限制，后按用户提示回溯协议：`action.schema.json` 明确 `modal` 动作 + content 任意 Node + `closeModal` 行为——协议本就支持弹窗交互，剪贴板方案属过度绕行，已回退并协议化。经验：披露响应体能力的判断应以协议为准，而非渲染器实现现状。
- 未改动任何目标 status / progress（Root 保持 `done 4/4`）；本条目为维护事实记录。
- 验证：每轮 Go/web 相关测试全绿；全量 vitest 现为 **1113/1113**（含 W25 组件注册守卫与 D-VAL schema 校验扩展至 `users-invites`）；tsc 干净。