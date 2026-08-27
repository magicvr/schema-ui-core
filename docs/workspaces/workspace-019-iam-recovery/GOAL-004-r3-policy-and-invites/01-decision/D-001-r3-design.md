---
id: D-001
doc: decision-entry
goal: GOAL-004-r3-policy-and-invites
status: accepted
created: 2026-08-25
updated: 2026-08-25
version: 1.0.0
---

# D-001 · R3 密码策略 + 邀请入职方案冻结

## 背景

R1 合同（GOAL-002 D-001 §2/§3）与 I-003/I-004/I-005/I-007（verified）冻结产品语义。2026-08-25 用户裁决新增两点：①**受邀账号初始角色 = 管理员发布邀请时指定**（激活后按邀请角色建号）；②Web 新建用户表单应支持直接选角色（后端 `POST /api/users` 已收 roles，表单缺字段 → 列入 C4）。技术细节按惯例留痕：纯链接邀请允许不带邮箱。链接形态（A-001 F-003 回写定稿）：激活页路径为相对路由 `/invite/accept?token=…`；管理 API 响应与邀请邮件中的链接由请求 Host（含 X-Forwarded-Proto 推导 scheme）拼接为**绝对链接**——比纯相对路径更实用，邮件场景必须可点击。

## 方案条款（冻结）

### §1 迁移 0057/0058（owner：core.auth-session）

- **0057 password_policy**（可移植，ApplyPostgres nil）：单例行 `id INTEGER PK CHECK(id=1)`、`min_length` DEFAULT 8、`min_categories` DEFAULT 0、`history_depth` DEFAULT 0。
- **0057 user_password_history**（PG 成对 BIGINT 时间列）：id TEXT PK、user_id FK CASCADE、password_hash、created_at。
- **0058 user_invites**（PG 成对）：id TEXT PK（公开令牌 id）、token_hash UNIQUE、roles JSON NOT NULL、invited_by FK users、email 可空、expires_at、consumed_at 可空、revoked_at 可空、last_sent_at、created_at。

### §2 密码策略域与四口强制

- `authsession.ValidateNewPassword(userID, plain) error`：读策略行 → 基线 8–72 字节非空白 → min_categories>0 时字符类别计数（小写/大写/数字/其他）不足即拒 → history_depth>0 时对最近 N 条历史 bcrypt 比较命中即拒。错误哨兵映射 INVALID_PASSWORD（复用）。
- 历史捕获：UpdateUser 带 PasswordHash patch 时同事务把旧 hash 推入历史并按 depth 裁剪（捕获为 best-effort，失败不阻断设密）；创建账号无历史写入。校验语义注明（A-001 F-004 回写）：**轮换前的当前密码不在历史中**——它在成功轮换的同一事务内才入史，因此「改回当前密码」需先完成一次轮换后才被拦；这是既定口径而非缺陷。
- 四口接线（设密前调用 ValidateNewPassword）：users Create / users Update(password) / account_self changePassword / recovery complete。
- 配置面：admin.settings 模块新增 GET/PATCH `/api/settings/password-policy`；**GET 走 `settings.read`、PATCH 走 `settings.write`**（A-001 F-004 回写，与邮件 tab 读/写分权同形）；PATCH 范围校验 minLength∈[8,72]、minCategories∈[0,4]、historyDepth∈[0,10]；UI 形状为 admin.settings tab 扩展（C4）。

### §3 邀请入职（角色由邀请指定 · 用户裁决 2026-08-25）

- 管理 API（admin.users 模块，权限 `users.invite` PolicyAdmin）：
  - POST `/api/users/invites` {email?, expiresInDays?, roles[]}：roles ≥1 且全部存在（fail-closed）；默认有效期 7 天（I-005）；返回 invite + token（仅此一次）+ 相对链接；带 email 时经 MailSender 发邀请信（Host 推导基址），发送失败报 EMAIL_SEND_FAILED。
  - GET `/api/users/invites`；DELETE `/api/users/invites/{id}`（撤销即时失效）；POST `/api/users/invites/{id}/resend`（撤旧发新新 token，60s 冷却）。
- 激活 API（公开面中央注册）：POST `/api/auth/invite/accept` {token, username, name, password}：token sha256 定位 → live 校验（未过期/未消费/未撤销）→ 角色 re-validate（有被删角色 → INVITE_ROLE_GONE，需管理员重发）→ 用户名唯一冲突 USERNAME_TAKEN fail-closed → 密码过 ValidateNewPassword → 同事务建号（角色=邀请指定，MustChangePassword=false）+ 消费邀请 → 不签发会话（D-001 §4 投影）。
- 错误码新增：INVALID_INVITE_BODY / INVITE_INVALID（未知/过期/已用/已撤销统一）/ INVITE_ROLE_GONE；复用 USERNAME_TAKEN / INVALID_PASSWORD / EMAIL_SEND_FAILED。

### §4 Web 面（C4）

设置页密码策略字段组；用户页邀请管理面板；公开 `/invite/accept` 激活页；新建用户表单补 roles checkboxGroup（复用 edit-user-roles-form 的 checkboxGroup/optionsSource 形状）；i18n zh/en。

### §5 边界

渐进生效不扫存量；无 SMS/模板中心/多邮箱；不动 Profile 默认集；重开相邻 VP 契约禁止。

## 未选方案

- 受邀账号固定 viewer / 无角色：与用户裁决（邀请时指定角色）不符，未选。
- 邀请表挂 admin.users 迁移 owner：IAM 能力归 core.auth-session 扩展（VP-003 边界），未选。
- 激活时直接签发会话：偏离「回登录页」投影口径，未选。

## 后续

C2 迁移 → C3 后端实施与测试 → C4 Web → C5 independent + self 关门。
