---
id: D-003-w18-a002-response
doc: decision-entry
goal: GOAL-019-w18-api-web-security-hardening
date: 2026-09-22
status: accepted
parent: GOAL-001-production-hardening
version: 0.1.0
---

# D-003 · 响应 A-002 的修复路径

## 决定

采纳 A-002 的 F-001～F-003 作为本波 required 响应范围，并采用以下路径：

1. `serve` 当前明确未装配 MFA verifier，因此不在本波扩大为完整 MFA 服务重构；在打开 store 后、listener/runtime 启动前查询 `user_mfa.status = 'active'`。查询失败或命中 active enrollment 均 fail closed，并关闭 store。
2. Web 目标 URL 解析移除当前 realm 的 `instanceof Request/URL` 依赖，改为读取结构化 `url` / `href`；不可安全解析的对象返回不可认证目标，不再字符串化为当前 origin。
3. 修正 goal-tree、执行索引和审计信息表，保持 S6/I-004 open，直到修复后再次通过干净上下文独立审计。

## 取舍

完整装配 store-backed MFA verifier 会改变当前 serve 的模块边界、配置面和认证端点契约，超出本波“修复已确认旁路”的最小安全闭环；fail-closed 会拒绝承载已有 active MFA enrollment 的 serve 实例，但不会改变无 MFA enrollment 数据库的既有启动路径。该残余兼容性代价通过明确错误和启动前检测暴露，而不是静默降级。

## 放行边界

本决定不关闭 A-002 的 required findings；修复事实记录在 E-004，self 响应与复审前状态记录在 A-003，最终闭合必须由新的 clean-context Reviewer 独立确认。
