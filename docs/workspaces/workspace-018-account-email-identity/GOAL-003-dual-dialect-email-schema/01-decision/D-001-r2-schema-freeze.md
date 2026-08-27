---
id: D-001
doc: decision-entry
goal: GOAL-003-dual-dialect-email-schema
status: accepted
created: 2026-08-24
updated: 2026-08-24
version: 1.0.0
---

# D-001 · R2 schema 设计冻结（迁移 0054）

## 背景

R1 合同（GOAL-002 D-001 §1/§2/§3/§6）已冻结身份面语义。R2 把它落成物理 schema。勘察结论：全仓 catalog 头 = v53（operationlog mail events）；authsession 已参与全局 checksum 台账；新增迁移须同步五处黄金断言（workspace-010 GOAL-033 先例）；users 的全部 INSERT 均显式列名，无非测试 `SELECT *`。

## 决定（条款）

1. **版本**：`Version: 54`，`Name: "account_email_identity"`，transform tag `"0054:account-email-identity:v1"`；归属 `core.auth-session` 迁移贡献（Descriptors 追加，不改任何既有 DDL 字符串 → 既有 checksum 不变）。
2. **DDL（三条，双方言同一文本 = 可移植，`ApplyPostgres: nil`，先例同 0011/0038）**：
   - `ALTER TABLE users ADD COLUMN email TEXT`
   - `ALTER TABLE users ADD COLUMN email_status TEXT CHECK (email_status IN ('pending','verified'))`
   - `CREATE UNIQUE INDEX idx_users_email_lower ON users(lower(email))`
3. **状态语义**：`email IS NULL ⇒ unbound`（无视 email_status）；email 非空时 email_status ∈ {pending, verified}。存量行 ALTER 后为 (NULL, NULL) = unbound ✓。
4. **唯一性承载**：表达式唯一索引 `lower(email)`——PG 与 SQLite 均支持表达式索引；NULL 在两方言唯一索引中互异 → 多账号无邮箱共存。SQLite lower() 仅 ASCII 折叠的差异由应用层归一补偿（GOAL-002 A-001 F-2 移交项，R3 仓储实现时落）。
5. **黄金断言同步**：identity.go head 53→54；identity_test lockedHeadExtraTables[54]={}（无新对象，注明理由）；migrate_test 两处计数与链尾；operations_test / restart_test 尾断言。
6. **不做**：不加 token 表（I-005/R3）；不改 fingerprintR2 收紧口径（遗留无 ledger 库保持 fail-closed）；不写仓储层邮箱语义。

## 为什么

- 单列 + 状态列直接映射合同三态，最小面；CHECK 内联枚举符合仓内风格（roles.system 等）。
- 可移植 DDL 免去成对 PG body，降低双方言漂移面；表达式索引是「lower 唯一」在两方言的标准解。

## 未选方案

- `email_verified INTEGER 0/1`：丢 pending 态，需再加列才能区分待校验，否决。
- CITEXT（PG）/ COLLATE NOCASE（sqlite）：仅列级大小写折叠，无法做函数索引的唯一槽且方言分叉，否决。
- 部分 WHERE email IS NOT NULL 索引：与全表达式索引等价但多一条款，从简。

## 后续

- 实现迁移 + 黄金断言 + 专项测试 → go build/test 绿 → git checkpoint → independent 审计（grok build）→ 关门。
