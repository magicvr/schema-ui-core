---
id: E-004
doc: execution-entry
goal: GOAL-004-r3-policy-and-invites
status: recorded
created: 2026-08-25
updated: 2026-08-25
version: 1.0.0
---

# E-004 · C4 满：Web 面落地

2026-08-25 完成：

- 设置页新增「密码策略」tab（password-policy-tab 自定义组件，GET/PATCH /api/settings/password-policy；settings.json 挂 custom 节点）。
- 用户页顶部新增「邀请管理」面板（user-invites-panel：生成（角色随邀请/邮箱可选/天数）+ 列表 + 重发/撤销；token/link 仅响应时一次性披露）。
- 公开 /invite/accept?token=… 激活页（invite-accept.tsx；未认证路由分支接入 main.tsx AuthGate；成功不签发会话回登录）。
- 新建用户表单补 roles checkboxGroup（复用 edit-user-roles-form 的 optionsSource 形状；后端已支持）——用户诉求②。
- i18n zh/en 各 ~30 键（invite.*）；W25 守卫 side-effect imports 补齐；两组件自注册。
- web tsc 干净；vitest 1105/1105 绿。
