---
id: GOAL-014-w13-account-lockout-redesign
title: W13-F007 账号锁定模型重设计（fixed · 承载自 GOAL-013）
status: done
created: 2026-08-26
updated: 2026-08-26
parent: GOAL-013-w13-api-web-security-audit
version: 0.5.0
progress: 6/6
---

# GOAL-014 · W13-F007 账号锁定模型重设计

## 意图

承接 GOAL-013（[workspace-009] W13）A-001 F-007 的用户裁决（2026-08-26，D-001）：**fixed**——现账号锁定模型（连续 5 次密码错误 → 全局账号锁 15 分钟 + 吊销全部刷新令牌，`internal/auth/auth.go` LockThresholdFailures/LockWindow + RecordLoginFailure/RevokeAllRefreshTokensForUser）允许知道用户名的攻击者反复锁定任意账号并强制其全部设备下线（定向 DoS）。本目标承载锁定模型重设计的方案、实施与审计闭环。

审计意见来源：[GOAL-013 A-001](../GOAL-013-w13-api-web-security-audit/03-audit/A-001-w13-security-review-findings.md)；处置裁决：[01-decision/D-001](01-decision/D-001-w13-f007-fixed-adjudication.md)。

## 路线图（progress 来源：以下 6 个检查点等权）

- [x] **S1 立项与裁决落盘** —— 用户三路径选择 fixed（经子目标承载）；五件套建立
- [x] **S2 方案设计与冻结** —— 候选模型对比（A 纯 IP 锁 / B 全局退避 / C 分层），选定模型 C（IP 维度锁 + 高阈值全局熔断 + 移除失败触发吊销），写入 D-002 并冻结
- [x] **S3 实施** —— 迁移 0061 login_failures + users.last_login_failure_at；auth.Login 分层校验（来源锁 5/15min → 全局锁 100/24h 滑动）；移除失败触发会话吊销；UnlockUser 清来源行；handler 登录传真实客户端身份（checkpoint `26655b55`）
- [x] **S4 回归** —— go vet ./... 0 输出；go test ./... -count=1 全绿 46 包（含 store 目录头 pin 更新至 v61）；缺陷形状回归锁 ×3（来源隔离 / 全局制动+通知恰一次+自愈 / 失败不吊销令牌）
- [x] **S5 审计** —— self 审计 pass（A-001）→ independent 审计 pass（A-002 · grok build · grok-4.6 · reasoning high，开放 required=0；recommended ×3 已由 A-003 全部响应闭合）
- [x] **S6 关门** —— 用户书面关门确认已获得（2026-08-26，[D-003](01-decision/D-003-closeout.md)）：与 GOAL-013 一并 `status: done`

## 残余移交（留痕）

- 全局熔断后的 15min 账号级登录拒绝——**设计意图**（分布式猜测制动），非缺陷。
- 全局锁窗内 Refresh rotate-before-checks（仅锁窗内主动出示的那张令牌被消费）——A-003 R-F002 accepted-residual；复审触发：产品提出"全局锁窗内出示也不丢会话"时另开决策改 Refresh 校验顺序。

## 边界

- 不改变登录失败响应的防枚举时序语义（missing-user/locked/disabled 的 dummy-bcrypt 保持）。
- 不改变 ErrAccountLocked 冻结线契约（wire code 不变）。
- 本目标完成情况不影响 GOAL-013 其余 finding 的推进；两目标关门顺序由 S6 阶段用户裁决。
