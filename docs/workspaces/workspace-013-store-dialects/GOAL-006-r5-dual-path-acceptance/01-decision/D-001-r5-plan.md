---
id: D-001
doc: decision-entry
goal: GOAL-006-r5-dual-path-acceptance
status: accepted
created: 2026-08-20
updated: 2026-08-20
version: 1.0.0
---

# D-001 · R5 方案：升级策略、备份合同、共事务与关门

## 触发

- Root R1–R4（GOAL-002..GOAL-005）done；portgres 完整启动已 live 证明。R5 收敛 Root I-001/I-004 并做生产向验收后关门 VP-013（退出判据 1–6）。

## 决定

1. **U0 基线固化**：sqlite `go test ./...` 0 FAIL + `TestFullCatalogPostgresBootstrapIntegration` + `TestCompositionPostgresStartup` 作为 R5 基线证据。
2. **U1 · 升级策略（I-001）**：结论导向——in-place int8 迁移 vs dump/restore vs fresh bootstrap。预期：**fresh bootstrap（重建 PG 库 + 双方言 catalog apply）为推荐**；存量数据经**导出/导入（dump/restore 或镜像命令）**迁移；若证据不足则以书面 residual 收口（范围 + 复审触发）。落盘 Decision + 抽样原型证据。
3. **U2 · 备份/恢复合同（I-004）**：替代 SQLite `VACUUM INTO` 的 PG 生产路径 = **pg_dump/pg_restore（custom format）或 `pg_basebackup`**（按部署规模）；落盘最小合同（何时用什么、恢复步骤、验证方式）并在 live PG 上做一次 dump→restore 可执行验证。
4. **U3 · 跨模块共事务**：在 PG 上验证一个跨模块事务（如 service-credential-use 审计 + operationlog 同事务回滚/提交）一致；`readyz` 模块门禁全绿（`TestCompositionPostgresStartup` 已含入门禁路径）。
5. **U4 · 关门**：VP-013 退出判据 1–6 逐条核对；self + **independent**（grok build，production 门禁）；无 open required 后 GOAL-006 `done` → **Root 5/5 关门**（workspace-013 结项）。

## 为什么

- 升级与备份是 VP-013「退出判据 2/4」的点名项；R5 是唯一承接处。
- PG 已能完整启动（R4），生产向验收（迁移、共事务、备份）只剩证据化工作。
- fresh bootstrap 是 V1 已默认路径（R3 全量 PG boot 已证明），fs 升级需额外原型；以证据定结论而非硬编码。

## 未选方案

- 把 I-001/I-004 余留为 deferred：R5 是它们的最晚需要阶段，不做=门禁未满足、不能关门。
- 用 SQLite `VACUUM INTO` 当作 PG 备份：无等价物；不选。
- 在 R5 重开 R1–R4：已完成；不选。

## 影响范围

- 文档（策略/合同/证据）+ 辅助脚本/测试（备份验证）；主体运行时代码 R1–R4 已冻结，R5 不改公共契约。`SCHEMA_UI_R2_PG_DSN` 门控测试为验证载体。

## 后续

U0→U4；U1/U2 先落盘决策与 required 信息项，U4 前 self+independent 关门。
