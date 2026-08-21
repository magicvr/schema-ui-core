---
id: E-003
goal: GOAL-020-wallet-auto-account
title: 用户反馈修复批（流水页 404/命名/面包屑/状态弹窗）+ A-003~A-005 闭合
date: 2026-08-16
status: recorded
parent: GOAL-020-wallet-auto-account
created: 2026-08-16
updated: 2026-08-16
version: 1.0.0
---

# E-003 · 用户反馈修复批 + 审计闭合（2026-08-16）

## 事实

- **A-003（grok independent · conditional）响应**：F-001（并发冲突 created 误报）→ fixed（冲突重读分支 isNew=false + TestGetOrCreateUserAccountConcurrent 8 goroutine 并发测试）；F-002（调账失败开户不审计）→ fixed（account-create 审计前置于 Mutate + TestWalletByOwnerAdjustFailureStillAuditsCreate）；F-003（id 随机后缀）→ fixed；F-004（D-001 路径勘误 v1.1.0）→ fixed；F-005（403 循环/错误码表/台账卫生）→ fixed。
- **A-004（grok finding-closure）**：F-002/F-003/F-004 closed；F-001 点名并发测试缺失。
- **A-005（grok finding-closure）verdict pass（0 required）**：并发测试真实存在且通过；代码 isNew=false 静态核对；1 条 recommended 残留（MaxOpenConns(1) 串行化使 UNIQUE 冲突重读运行时不可达——代码可核对，登记为残余）。
- **用户反馈批（2026-08-16）**：
  1. 流水页 404「钱包账户不存在」→ 根因：行导航 navigate 无 navigateMapping 时 URL 模板原样使用（/wallet-entries/{id} 字面 → accountId={id}）；修复：wallet.json 表格 entries 动作条目加 navigateMapping path {id: $row.id}；新增 wallet-navigate 渲染测试断言导航目标 /wallet-entries/acct-1。
  2. 菜单/页面命名「钱包/账务」→「钱包」（zh-CN manifest.title/nav.wallet）。
  3. 流水页为钱包内页（面包屑）→ App.tsx BREADCRUMB_PAGE_PARENTS 加 wallet-entries→wallet。
  4. 状态弹窗 FORM_VERSION_TOO_LOW/FORM_CAPABILITY_REQUIRED → wallet.json protocolVersion 2.7→2.9 + requiredCapabilities 加 form.controls.readonly（readOnly 字段门禁）。
  5. 顺带核查字典页同款导航：data-dictionary 表格 entries 动作条目本就带 navigateMapping（path dictKey + query dictTypeName），无同款 bug；此前误加的定义级 navigateMapping 已还原（actions 定义不允许该字段）。
- **验证**：go 全量全绿；web 全量 **1005/1005**（+1 wallet-navigate 测试）；schema D-VAL 全过。
