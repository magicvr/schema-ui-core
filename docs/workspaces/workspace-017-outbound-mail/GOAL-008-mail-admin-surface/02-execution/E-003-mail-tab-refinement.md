---
id: E-003
doc: execution-entry
goal: GOAL-008-mail-admin-surface
status: recorded
parent: GOAL-001-outbound-mail
created: 2026-08-24
updated: 2026-08-24
version: 1.0.0
---

# E-003 · 关门后 UX 精修：邮件 tab 按渠道聚焦（2026-08-24）

## 背景

用户在 VP-017 v0.5.0 closed 之后书面提出三项设置页 UI 优化，并要求纳入治理上下文。经核：三项均为 R7 管理面的**呈现层精修**，不改变任何已冻结合同语义（D-002 渠道集/解析、D-007 密钥写后不可读、试发走同一端口均保持），故以本执行条目记录于 R7 所属目标的台账，**不重开** Root/VP/子目标状态。

## 已发生事实

1. **按渠道选择性显示设置项**：新增自注册自定义组件 `mail-admin-tab`（沿 data-permission-scopes 先例）——渠道 select 驱动条件渲染：mock → 保留条数；resend → api-key/from；smtp → host/port/username/password/from。settings.json 的 declarative 表单无法表达跨节点联动（R5 D-EXPR reactions 冻结为 $context-only，无字段值触发），故采用既有 custom-component 扩展缝。
2. **mock 出站记录仅 mock 渠道时显示**：记录表移入组件内，随 `channel === "mock"` 条件渲染。
3. **试发自定义标题/正文**：`POST /api/mail/test-send` 接受可选 `subject`/`body`（空值回退默认）；组件提供 to/subject/body 输入；审计 detail 增记 subject。
4. settings.json：tab-mail 子节点收敛为单一 custom 节点；移除失效的 updateMail/mailTestSend actions。i18n 双语各补 13 键。

## 测试证据

- 新增 `mail-admin-tab.test.tsx`：①三渠道条件显隐断言（mock 表可见/resend·smtp 字段互斥）；②PUT 平铺键精确匹配 + 试发 subject/body 透传断言。
- handler 测试扩展：test-send 自定义 subject/body 持久化进 outbox 断言。
- 复跑确认：web vitest 全量 78 文件 / 1100 用例全绿（含 W25 schema↔registry 防复发守卫——新组件已登记）；tsc -b + vite build 通过；api handler/mail/config 包绿。

## 未做

- 未改公共端口合同、密钥读取面、渠道解析算法；Root/VP 关门状态不变。
