---
id: E-003-s3-regression-and-go-verdict
doc: execution-entry
goal: GOAL-038-w26-email-display-and-mail-pages
date: 2026-08-26
author: govern orchestrator（S3 回归）
---

# E-003 · S3 回归结果与 go 消费判定

## 回归命令与结果

| 套件 | 命令 | 结果 |
|------|------|------|
| Go 全量 | `go test ./...`（apps/api） | **全绿，0 FAIL**（含 store 迁移头快照 v60 同步后重跑；`internal/store` ok 30.6s） |
| Web 单测 | `npx vitest run`（apps/web） | **81 文件 / 1116 测试全过** |
| 类型检查 | `npx tsc --noEmit` | **0 错误** |
| 生产构建 | `npm run build` | **成功**（chunk >500kB 警告为既有现象，非本次引入） |

迁移快照维护（W26 新增 0060 后的强制同步，均有测试锁定）：`store/identity.go completeFingerprintCatalogHead 59→60` + `lockedHeadExtraTables[60]={}`（纯加列无新对象）；`migrate_test.go / operations_test.go / restart_test.go` applied 头 v59→v60 `mail_outbox_channels`；`TestCompiledMigrationCatalogOwnership` want 增 `{"core.persistence","mail_outbox_channels","6f9d3771…"}`。

## go 消费判定（VP-010 接口）

- **改动面**：admin.settings 模块 additive 产品面扩展（+2 页面 / +2 导航节点 / 设置页移除 tab-mail）；users 资源读面 additive 字段（email/emailStatus/emailStatusStyle）；mail_outbox 表加列（默认值即存量语义）；GET /api/mail/outbox 列表项 additive 字段。
- **未动**：协议 pin（2.7）、Profile 默认集与模块矩阵语义（模块集合不变，仅贡献清单内 additive）、权限键集合（零新增）、既有路由契约（仅响应体 additive 字段）、破坏性数据迁移（无）。
- **判定**：**无影响、不暂挂**。additive 页面/字段对既有消费方向后兼容；出站记录列表携带正文属 D-001 §2.1 显式契约修订（有界：retention ≤500、pageSize ≤200）。

## 遗留观察（non-blocking）

- ~~e2e（playwright 双方言矩阵）不在本波 S3 冻结范围~~ → **2026-08-26 用户指示补跑，双方言全绿**（见下节）。

## E-003 补充 · e2e 双方言矩阵补跑（2026-08-26，用户指示）

| 方言 | 命令 | 结果 |
|------|------|------|
| sqlite | `npm run test:e2e`（apps/web，临时隔离库） | **9 passed / 1 skipped**（skip = admin 专属用例在 mvp profile 的预期跳过），exit 0，约 1.9 分钟 |
| postgres | `npm run test:e2e:postgres`（`cmd/e2e-pgset` scratch 库 create→run→drop；凭据取自 gitignored `apps/api/configs/.env` DB_* 键） | **9 passed / 1 skipped**（同上预期跳过），exit 0；teardown 输出确认 scratch 库 `schema_ui_e2e_mt9hakmzcm1ji1` 已 drop |

- 跑前探针：`e2e-pgset create/verify/drop` 生命周期验证通过（verify=1 为新库无 schema_migrations 的预期 fail-fast）。
- 两轮串行执行（共用 WEB_PORT 25173 / API 25080，挂具 `reuseExistingServer: false` 禁并行复用）。
- 结论：W26 改动在双数据库方言下的浏览器验收面均绿，GOAL-038 关门证据链补全（不改变 A-001 pass 与 done 4/4 状态，仅扩充回归证据）。
