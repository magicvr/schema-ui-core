---
id: E-002
doc: execution-entry
goal: GOAL-004-r3-binding-flow
status: recorded
parent: GOAL-001-account-email-identity
created: 2026-08-24
updated: 2026-08-24
version: 1.0.0
---

# E-002 · R3 实现与验证（2026-08-24）

## 已发生事实

- **迁移 0055 `email_verification_challenges`**（authsession/migration/migration.go）：成对方言 DDL（时间列 INTEGER|BIGINT）；checksum `1556bda2…b754` 入冻结目录；黄金断言全套随动（head 55、lockedHead[55]、completeLostLedgerTables / postV1CatalogTables / identity_test 夹具、migrate_test 链尾×2、operations/restart 尾断言）。
- **服务层**（authsession/email_identity.go）：Bind（占槽冲突 fail-closed `EMAIL_TAKEN`；同址已 verified 幂等不重发；发信失败整体回滚）/ Verify（常量时间比较；失败计数与过期清理走独立事务——runner 对错误回滚的语义经测试驱动修正）/ Resend（60 秒冷却）；换绑=覆写释放旧址（合同 §5）。
- **HTTP 面**：POST /api/account/email/bind|verify|resend（admin.account 自助，无权限键）；7 个新错误码入 Catalog 并登记错误契约冻结清单；kernel/profile.go 路由清单同步。
- **I-006 代填**：UserPatch.Email（users PATCH）→ 置 pending / 清空回 unbound；存储层冲突检查；mapUserStoreError 映射；profile GET 增读 email/emailStatus（EmailIdentityState）。
- **最小页面**：email-identity.tsx 自注册组件（mfa-manager 先例），account.json profile tab 挂 custom 节点；i18n zh/en 各 12 键。

## 验证矩阵

| 面 | 结果 |
|----|------|
| go build ./... | 干净 |
| store / authsession / kernel / composition / handler 全量 | 全绿（含新增 8 个服务流测试 + 黄金断言） |
| PostgreSQL 17 集成 `-run Postgres` | ok（全 catalog 含 0055 实跑） |
| Web vitest（组件 3 用例 + schema-keys 结构） | 7 passed |
| `npm run build`（tsc -b && vite build） | 成功 |

## 未做

- independent 审计待执行（R3 关门门禁）；Root 收口待审计后。
