---
id: A-003
goal: GOAL-011-r3-s11-login-captcha
source: independent
date: 2026-08-14
scope: S5 关门 · 安全门禁（admin.login-captcha vs D-002 冻结方案）
verdict: fail
auditor: grok-build
audit_type: close-out
status: recorded
parent: GOAL-011-r3-s11-login-captcha
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# A-003 · independent 安全审计（S-11 实现 · 首轮）

## 结论

**verdict: fail**。5 个 required findings（F-001 过期未强制、F-002 Web 登录客户端未接线、F-003 settings bool/string 不匹配、F-004 删除失败 fail-open、F-005 真服务 HTTP 测试缺失）+ 3 个 recommended（F-006 Required fail-open、F-007 树中 WIP 破坏、F-008 文档不一致）。原始意见：`attachments/grok-audit-s11-closeout.txt`。

## Findings

- F-001（required, high）：Verify 未检查 expires_at（GetChallenge 只取 answer_hash）；惰性清理只在 CreateChallenge。
- F-002（required, high）：LoginPage/auth-client 未按 D-002 §5 预检与提交 captchaId/captchaAnswer；INVALID_CAPTCHA 未映射。
- F-003（required, high）：schema select 提交 "true"/"false" 字符串，PATCH 只接受 JSON bool → 400。
- F-004（required, med）：DeleteChallenge 错误被忽略（`_ =`），删除失败返回 nil 视为成功。
- F-005（required, med）：HTTP 测试用 fake；真 Service 缺 generate→login 200、403、captcha 失败不计锁定。
- F-006（recommended）：Required() 配置读取失败 fail-open。
- F-007（recommended）：审计时 S-12 WIP 未闭合（0025 缺失、recyclestore 引用失败）。
- F-008（recommended）：00-meta I-003 引用不存在的 D-003；01-decision.md I-001 状态 open 与 00-meta closed 冲突。

## 响应（required 全部修复）

- F-001/F-004：store 改为原子 `ConsumeChallenge`（同事务校验+删除+过期强制，expires_at/now 参与判定），Verify 映射全部失败为 ErrInvalidCaptcha；store/challenge 单测补过期拒绝与原子消费。
- F-002：auth-client.login 增可选 captcha（captchaId/captchaAnswer）；AuthContext/LoginPage 预检 GET /api/auth/captcha + 挑战输入 + INVALID_CAPTCHA 映射；LoginPage 单测 +3。
- F-003：settings GET 返回表单字符串 "true"/"false"，PATCH 接受 bool 或字符串（parseBoolValue，notifications 模式）。
- F-005：provider_test 真 Service：解题→登录 200；editor 403；25 次 captcha 失败后关闭门禁登录 200（不计锁定/限流）。
- F-006：Required() 配置读失败视为开启（fail-closed）。
- F-008：00-meta I-003 → D-002 §6 / A-002；01-decision.md I-001 → closed。
- F-007：S-12 完工后全量回归（见 A-004 复审）。

## 闭合

F-001~F-006、F-008 已 fixed（A-004 复审 + 全量回归确认）；F-007 随 S-12 收尾验证。
