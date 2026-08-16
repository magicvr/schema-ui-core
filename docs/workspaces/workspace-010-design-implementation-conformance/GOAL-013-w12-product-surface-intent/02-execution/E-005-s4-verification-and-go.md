---
id: E-005-s4-verification-and-go
doc: execution-entry
goal: GOAL-013-w12-product-surface-intent
date: 2026-08-16
status: recorded
parent: GOAL-001-design-implementation-conformance
created: 2026-08-16
updated: 2026-08-16
version: 0.1.0
---

# E-005 · S4 验证与关门（回归 + go 判定）

## 事实

### 回归证据（2026-08-16）

| 面 | 命令 | 结果 |
|----|------|------|
| Go 全量 | `go test ./... -count=1`（apps/api） | **0 FAIL**（全部包 ok） |
| Web 全量 | `npx vitest run`（apps/web） | **1027/1027**（63 文件；基线 1022 + 新增 5 例） |
| 类型检查 | `tsc -b`（apps/web） | **0 错误** |
| D-VAL 结构回归 | 全模块 schema（all-module-schemas-dval.test.ts） | 绿（含改造后 12 页 + account tabs） |

### go 判定（T-06 部署契约变化）

- **变化面**：模块启用集权威从 env（`APP_PROFILE` / `APP_MODULES_ENABLED`）移到 `configs/config.yaml` 的 `app.profile` / `app.modules`；compose 增加 configs 只读挂载；dev.cmd 改用 CONFIG_FILE overlay；文档教学面改写。
- **不变面**：mvp/admin/demo 三档模块成员与 `kernel.profileDefaults` 完全一致（本波未改 profile.go 默认集）；Manifest 装配语义、模块矩阵、Web 构建页面集均不变；密钥等敏感项仍走 VAR 插值（I-006）。
- **判定**：**部署契约变化、默认集不变 → 不暂挂业务消费；不触发「改 Profile 默认集」门禁**。operator 需按新契约迁移（写 YAML 而非 env）；config_test「旧 env 选择器被忽略」用例保证迁移期不静默双轨。

### 审计

- A-001 self（S3/S4）：pass，无 required findings。
- A-002 independent（grok-4.6 · high）：**conditional**，required **F-001**（关门台账先于本独立意见落盘：预写 done / 预写 A-002 pass）。响应见 E-006（台账纠正 + F-003/F-004 修复 + Web 全量复跑）。
- 已核对 P-005 信息项：I-001～I-006 全部 verified（D-002～D-008 落盘）；I-005 non-blocking 已在 D-002 闭合。

### Git 提交

- `074d8ad`：S3-P0/P1（T-05/T-01/T-03/T-02）
- `f056c9d`：S3-P2（T-06）
- 关门文档提交：随本 E 条目。
