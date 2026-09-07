---
id: GOAL-040-w28-admin-passwd-convention
doc: audit-entry
record_id: A-001
status: closed
parent: GOAL-001-design-implementation-conformance
created: 2026-09-06
updated: 2026-09-06
version: 1.0.0
---

# A-001 · W28 关门自审（2026-09-06）

## A-001 · 关门自审（2026-09-06）
- **source**：self
- **auditor**：/govern 编排器
- **类型**：close-out
- **scope**：GOAL-040 W28 全波（ADMIN_PASSWD 约定 + AI 可发现性 + TEST_ADMIN 移除 + smoke 回退）
- **verdict**：pass

## 范围与区间

本审覆盖 2026-09-06 从 S1 决策冻结到 S4 关门全链：`apps/api` 代码面 TEST_ADMIN 移除、`ADMIN_PASSWD` 约定登记（.env.example / README / QUICKSTART / AGENTS.md）、smoke 回退、canonical 模板守卫更新。

## 成果（有证据）

1. **D-001 决策落盘**：`01-decision/D-001-admin-passwd-convention.md`（决定/理由/未选方案/影响/后续齐全）；I-001～I-003 verified（用户 2026-09-06 三项裁决）。
2. **TEST_ADMIN 移除**：config.go 字段+env 读取、bootstrap.go `EnsureTestAdmin`、composition.go upsert 分支、reconcile_test.go `TestEnsureTestAdmin`、.env.example 块全部移除。
3. **ADMIN_PASSWD 约定 + AI 可发现性**：`.env.example` 登记声明语义（API 不读取不重置）；`AGENTS.md` 新增「本地开发环境与 admin 凭据约定」小节（AI 助手每次会话加载，连现有库时从 .env/环境变量读取）；README/QUICKSTART 同步。
4. **smoke 回退**：`SMOKE_PASSWORD` 缺省回退 `ADMIN_PASSWD`；pre-release 显式传参不受影响。
5. **canonical 守卫**：`env_example_test.go` `declarationOnlyKeys` 白名单，注释说明合法例外。

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| `.env.example` 登记 ADMIN_PASSWD 声明约定 + 明确不重置 | ✅ | `.env.example` ADMIN_PASSWD 块；`env_example_test.go` 白名单注释 |
| AGENTS.md/QUICKSTART/README 让 AI 与自动化测试知道从 .env/环境变量取 | ✅ | AGENTS.md「本地开发环境」小节；QUICKSTART「连已有库的 admin 凭据」；README 配置键表 |
| TEST_ADMIN 全部移除 + 相关测试绿 | ✅ | `go test ./...` 0 FAIL；`rg` 代码面无残留 |
| smoke.sh 以 ADMIN_PASSWD 作为 SMOKE_PASSWORD 回退 | ✅ | smoke.sh L42/L126-130 头部与回退逻辑 |
| goal-tree 同步 + 关门无开放 required | ✅ | goal-tree 树/表已加 GOAL-040；本审 0 required |

## Findings

- **F-001 · 语义边界**（严重度：low · 建议：recommended · 状态：closed）
  描述：`ADMIN_PASSWD` 是声明而非 API 输入，若未来有人误以为可据此重置密码会误用。
  处置：`.env.example` 与 AGENTS.md 均明确「API 不读取、不据此重置密码」；`env_example_test.go` 白名单注释再加固。已闭合。

## 必改项汇总（required 列表）

无（0 required）。

## 结论 + 建议下一步

scope 内成功标准全部达成，无开放 required；用户已确认关门（2026-09-06）并授权验证后置 `status: done`（4/4）。

### 关门后验证证据（E-002 补充）

- 真实栈验证：API 连现有 postgres 库启动，`ADMIN_PASSWD` 登录 `admin` → HTTP 200 + token + `mustChangePassword=False`；`/api/accounts/me` roles=`admin,editor`。
- `SMOKE_PASSWORD` 未设 + `ADMIN_PASSWD` 导出 → `scripts/smoke.sh` SM-001~005 全 PASS（SM-006/007/008 预期 SKIP）。
- 后台进程已清理，端口 25080/25173 释放。Root/VP 保持 active 程序容器。
