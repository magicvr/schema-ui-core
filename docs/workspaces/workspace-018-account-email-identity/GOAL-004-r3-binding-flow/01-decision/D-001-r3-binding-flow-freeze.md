---
id: D-001
doc: decision-entry
goal: GOAL-004-r3-binding-flow
status: accepted
created: 2026-08-24
updated: 2026-08-24
version: 1.0.0
---

# D-001 · R3 绑定/校验流方案冻结

## 背景

R1 合同（GOAL-002 D-001）+ R2 schema（0054）就绪。2026-08-24 用户四项裁决关闭 I-005 / I-006 并确认分母：**验证码 TTL = 10 分钟；重发冷却 = 60 秒；允许管理员代填待校验；分母含最小 Admin 页面**。端口事实：`kernel.MailSender.Send(ctx, MailMessage{To,Subject,TextBody})` 同步单收件人、From 由适配器盖章；错误面 `errorcatalog.Catalog` + `writeError`。

## 决定（条款）

**§1 迁移 0055 `email_verification_challenges`**（双方言成对，非可移植——时间列方言差）：
`user_id TEXT PK REFERENCES users(id) ON DELETE CASCADE`；`code_hash TEXT NOT NULL`（sha256 hex）；`expires_at` / `sent_at` INTEGER|BIGINT NOT NULL（unix 秒）；`attempt_count INTEGER NOT NULL DEFAULT 0`。每用户至多一条活跃挑战（PK 即幂等覆写）。黄金断言全套随动：head 55、lockedHead[55]、completeLostLedgerTables / postV1CatalogTables 增补、冻结目录行。

**§2 服务语义**（authsession 模块内）：
- **Bind(userID, email)**：trim 后形状校验（含 @、≤254、无空白）；`lower(email)` 冲突查他号 → `EMAIL_TAKEN` fail-closed（pending/verified 均占槽）；同址同号已 verified → 幂等成功不重发；写入 users.email（原样）+ status='pending'；生成 6 位码（crypto/rand），存 sha256，经 MailSender 发信；发信失败 → 整体回滚并返回 `EMAIL_SEND_FAILED`。
- **Verify(userID, code)**：无挑战或过期 → 删陈旧 + `EMAIL_CODE_EXPIRED`；常量时间比较；不匹配 → attempt_count++，≥5 作废挑战 `EMAIL_CODE_INVALID`；匹配 → status='verified' + 删挑战。
- **Resend(userID)**：仅 pending 可用；距 sent_at < 60s → `EMAIL_RESEND_COOLDOWN`；否则新码覆写并发信。
- **换绑 = Bind 覆写**（合同 §5）：旧址槽随覆写释放，新址重新占槽待校验。
- 归一补偿（N-1）：比较与唯一性判定统一走 SQL `lower()`；应用层 trim + 大小写折叠后再比对哈希输入。

**§3 API 面（admin.account 自助 + admin.users 代填）**：
- `POST /api/account/email/bind {email}` / `POST /api/account/email/verify {code}` / `POST /api/account/email/resend {}`——自助身份面（登录即可，无权限键）；响应 `{status}`。
- 管理员代填：既有用户 PATCH 面新增可选 `email` 字段 → 置 pending（不得直接置 verified）；清空 email = 回 unbound。
- 新错误码入 Catalog：`EMAIL_TAKEN / EMAIL_INVALID / EMAIL_NOT_PENDING / EMAIL_CODE_INVALID / EMAIL_CODE_EXPIRED / EMAIL_RESEND_COOLDOWN / EMAIL_SEND_FAILED`。
- 操作日志沿用既有 users.* 事件枚举，不扩枚举。

**§4 最小页面**：账号页一处「邮箱绑定」卡（状态徽标 unbound/pending/verified + 输入框 + 发起/验证/重发按钮），消费上述三端点；不复制邮件设置 tab，不做模板中心。

**§5 尝试上限**：5 次失败作废挑战须重发（实现级防爆破，MFA 先例口径；非合同条款）。

## 为什么

- 挑战表独立于 users：换绑覆写不动挑战历史结构，PK= user_id 天然单活跃挑战；ON DELETE CASCADE 随账号清理。
- 成对 DDL 只差时间列类型，照 accountLock 先例；服务层全部走端口与 Catalog，零内核改动。

## 未选方案

- users 行内加 code 三列：污染身份主表且 CASCADE 缺失，否决。
- 复用 MFA/登录限流器做尝试上限：跨模块耦合，行内计数最简。
- HTML 邮件：VP-017 明确本波纯文本。

## 后续

实施 → 测试绿 → commit checkpoint → independent 审计（grok build · grok-4.6 · high）→ 关门。
