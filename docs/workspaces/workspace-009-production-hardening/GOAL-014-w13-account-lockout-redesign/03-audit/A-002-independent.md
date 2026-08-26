---
status: active
created: 2026-08-26
updated: 2026-08-26
parent: GOAL-014-w13-account-lockout-redesign
version: 0.1.0
---

# A-002 · GOAL-014 S5 关门前 independent 复核（grok build）

> 本条目由编排器自独立审计输出转录落盘（审计工具按约束"只输出、不修改文件"；全文原件见 [附件](../attachments/audit-A-002-grok-output.txt)）。编排器响应见 [A-003](A-003-a002-response.md)。

| 字段 | 值 |
|------|-----|
| **source** | independent |
| **auditor** | grok-build（grok-4.6 · reasoning high · `/audit`） |
| **类型** | finding-closure / close-out（S5 关门前复核腿） |
| **scope** | D-002 冻结方案实施真实性（迁移 0061/auth 分层校验/dummy-burn/来源锁/全局天花板/失败零吊销/OnLockOpened/UnlockUser/handler loginClientIP）；回归锁 ×3 与旧契约测试改写；A-001 三条备注；上游 GOAL-013 F-007 代码面闭合 |
| **verdict** | **pass** |
| **日期** | 2026-08-26 |
| **被审 HEAD** | `cf5675f1`（实现基线 `26655b55` + 并发首插兜底） |

## 核心结论（摘要）

1. **D-002 逐要素 genuine 落地**：迁移 0061（含 catalog pin/checksum/identity 表清单）、Login 分层顺序与 dummy-burn 时序、来源锁仅拒施害来源、全局天花板 100+24h 滑动重启、失败路径零吊销（全仓生产调用点 = 0）、OnLockOpened 仅全局、UnlockUser 同事务清来源行、handler 传真实客户端身份——全部源码级核对通过。
2. **F-007 武器化形状在代码面消除**："知用户名 → 5 败锁全账号 15min 并踢全部设备"不复存在：来源锁只拒施害来源，全局熔断阈值 ×20 且不再吊销刷新令牌。
3. **回归三锁真实覆盖缺陷形状**；旧契约测试改写"按 D-002 改写契约而非弱化防护"——W7 F-009 防枚举信封、通知键、disable 不重复断言全部保留。
4. **独立复跑**（四包抽查）：go vet 0；auth/authsession/handler/store 全 ok；具名回归锁复跑 ok。
5. **A-001 三条备注均成立、非新缺陷**；残余面缩小为"全局 100 次熔断后 15min 账号级登录拒绝（设计意图）+ 全局锁窗内 Refresh rotate-before-checks（既有契约）"。

## Findings（开放 required = 0；均 recommended）

- **R-F001**（低）：台账未跟上 HEAD `cf5675f1` 并发兜底提交；建议补记 + 可选并发首插回归。
- **R-F002**（低）：全局锁窗内 Refresh 先轮换后校验（既有契约）——接受为 D-002 范围内残余或另开决策。
- **R-F003**（低）：TestAccountSourceLockKeepsSessions 注释语义不准；缺"来源锁不触发 OnLockOpened""来源锁不波及 Refresh"负向/正向断言。

## 结论

**verdict = pass。F-007 代码面可记 genuine fixed。** 建议 `/govern` 响应 recommended ×3 后进入 S6 用户书面关门（两目标关门顺序按 GOAL-013 D-003 问用户）。

完整逐项核对表、字段表原文见附件：[attachments/audit-A-002-grok-output.txt](../attachments/audit-A-002-grok-output.txt)。
