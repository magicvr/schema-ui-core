---
id: D-001
doc: decision-entry
goal: GOAL-003-r2-self-recovery-flow
status: accepted
created: 2026-08-25
updated: 2026-08-25
version: 1.0.0
---

# D-001 · R2 自助恢复方案冻结（API / 错误码 / MFA 门 / 会话）

## 背景

R1 合同（Root D-002 + GOAL-002 D-001 §1）冻结了产品语义。本决策把语义落成可实施方案；无新 P-004 裁决点（全部为合同条款的机械投影，唯 N-3「完成是否直接签发会话」按默认建议执行：**完成后回登录页，不直接签发会话**——本条即 R2 定稿确认，依据：恢复链路不产出会话面最小，且与登录页既有 UX 无冲突）。

## 方案条款（冻结）

### §1 迁移 0056

`password_recovery_challenges`：与 0055 `email_verification_challenges` 同构（user_id PK + FK ON DELETE CASCADE、code_hash、expires_at、sent_at、attempt_count）；SQLite INTEGER / PostgreSQL BIGINT 时间列成对迁移体；checksum 入台账并同步黄金断言。

### §2 API 形状（公开面 · 预认证 · 中央注册）

- **`POST /api/auth/recovery/start`** `{account}` → 成功形 `202 {"status":"dispatched"}`。
  - 账号定位：`account` 字段 = 用户名（精确，同登录口径）或该账号已绑定且 `verified` 的邮箱（lower 折叠）。二者命中任一即可。
  - 命中且账号 enabled 且邮箱 verified → 冷却检查（60 s）→ 建/换挑战 → 经 MailSender 发信。
  - 未命中 / 无邮箱 / 未验证 / disabled → **同一成功形响应但不发信**（防枚举 fail-closed；W7 F-009 口径）。
  - body 非法 → 400 INVALID_RECOVERY_BODY；发送失败 → 补偿删挑战 + 502 EMAIL_SEND_FAILED（此时挑战存在性已暴露，属 VP-009 记录在案的残余，不在本波分母）。
- **`POST /api/auth/recovery/complete`** `{account, code, newPassword, secondFactorCode?, recoveryCode?}` → 204。
  - 校验顺序：body → 账号定位（未命中/无挑战/错码 → 400 RECOVERY_CODE_INVALID 统一码；过期 → RECOVERY_CODE_EXPIRED）→ MFA 门（如登记）→ 密码基线校验（现行 8–72 字节非空白，INVALID_PASSWORD）→ 设密。
  - 任一失败路径消耗挑战 attempt_count（≤5 次作废，含第二因子失败），封死「绕过邮箱码预算爆破 TOTP」。
- 两端点共用 loginRateLimiter 模型（IP|identifier 桶；complete 失败也 record）。

### §3 MFA 第二因子门

- admin.mfa 新增导出方法 `VerifySecondFactor(userID, code, recoveryCode, now)` = 既有 `requireActiveSecondFactor` 的薄封装（TOTP 窗口校验或一次性恢复码消费；nil-receiver fail-closed）。
- 组装根以 true-nil 接口（`handler.RecoverySecondFactor { Required; VerifySecondFactor }`）注入恢复 handler；模块关闭 → Required=false 完全跳过（合同 §1：未登记 MFA 不受影响）。

### §4 会话语义（I-008 投影落地）

设密走 `UpdateUser(PasswordHash)`：token_version+1、refresh 全撤销、must_change_password 清除——与自愿改密同一原子事务语义。完成后不直接签发会话（见背景）。

### §5 审计与通知

审计复用 `operationlog.EventAccountPasswordChange`（detail action=`password-recovery`，携带 username；不新增 event CHECK 迁移）；通知复用 `NotifyAccountEvent("account.password-changed")` 本地化面。

### §6 错误码（catalog 新增）

`INVALID_RECOVERY_BODY` / `RECOVERY_CODE_INVALID` / `RECOVERY_CODE_EXPIRED` / `RECOVERY_SECOND_FACTOR_REQUIRED`（MFA 账号 complete 缺第二因子字段）；复用 `INVALID_PASSWORD` / `MFA_INVALID` / `EMAIL_SEND_FAILED` / `RATE_LIMITED`。

## 为什么

- 与 w018 已验证的两阶段派发+补偿、挑战表形态、错误目录模式同构，实施面最小。
- 直接校验第二因子（无 login proof）避免复用 mfa_proofs 的「proof 即半张门票」语义泄漏到预认证恢复面；猜测预算由挑战 attempt 上限承担。

## 未选方案

- 复用 email_verification_challenges 表：生命周期语义耦合绑定流，未被选。
- verify/complete 两步 + 服务端 resetProof 表：多一张表多一种令牌，收益仅为 UX 分步，未被选（Web 端单表单分域收集即可）。
- 新增 operation_log 事件名：需 CHECK 重建迁移，detail action 已足够区分，未被选。

## 后续

- C2 迁移落地 → C3 域逻辑/API → C4 e2e + Web 面 → C5 independent（grok build · grok-4.6 · high）+ self 关门。
