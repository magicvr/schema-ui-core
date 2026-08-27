---
id: D-002
doc: decision-entry
goal: GOAL-001-account-email-identity
status: accepted
created: 2026-08-24
updated: 2026-08-24
version: 1.0.0
---

# D-002 · 冻结本 Root 直至 VP-017 再次关门

## 背景

用户书面（2026-08-24）：否决 VP-017 / workspace-017 组合层关门并升级 017 渠道分母；**VP-018 及对应工作区暂时冻结，在 017 重新关门之前不允许推进。**

本 Root 仍为 `0/4`，R1 未开工，尚无 `users` email 列变更。冻结不会回退开区 scaffold（D-001 / E-001 保留）。

## 决定

1. Root `status: active → blocked`。`progress` 保持 0/4。
2. **禁止**进入 R1 合同冻结、禁止改 `users` DDL、禁止创建 R1 子目标、禁止把 capture `Last()` 或 SMTP 锁成本 VP 验收权威。
3. 解冻条件（须同时）：VP-017 按**现行渠道分母**再次 `closed`；用户或 `/govern` 明确解冻。解冻后运输验收对齐当时默认渠道（预期 mock 出站记录），而不是历史 CaptureSink 专用叙事。
4. 本回合不改应用代码。I-001～I-006 状态不变（仍 collecting/registered）；冻结不等于信息项 verified。

## 为什么

- 017 运输面正在从 SMTP 专用升级到渠道模型；018 若继续会把错误的默认 sink 写进身份验收。
- 用户明确冻结，不是取消 VP-018。

## 未选方案

- 取消 VP-018 / Root `cancelled`：身份意图仍有效，只是排队。
- 继续 R1：违反用户书面冻结。
- 在 018 内重做邮件渠道：越界，运输归 017。

## 后续

- 等待 VP-017 再关门后由 `/govern` 解冻并重开 R1。
