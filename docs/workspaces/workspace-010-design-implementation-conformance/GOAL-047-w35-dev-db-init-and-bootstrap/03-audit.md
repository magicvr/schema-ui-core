---
id: GOAL-047-w35-dev-db-init-and-bootstrap
doc: audit
status: active
parent: null
created: 2026-09-21
updated: 2026-09-21
version: 0.1.0
---

# 审计台账 · GOAL-047-w35-dev-db-init-and-bootstrap（W35）

> 本文件是唯一正式审计台账索引：`self` 与 `independent` **共用** `A-NNN` 序列。
> 每条意见正文在 `03-audit/A-NNN-<slug>.md`；本文件登记条目头（`source`/日期/scope/`verdict`）。
> 独立审计默认只写意见，不修改 `status`/`progress`/方案正文；响应归编排器。

## 意见索引

| A-ID | source | 日期 | scope | verdict | 摘要 | 文件 |
|------|--------|------|-------|---------|------|------|
| A-001 | self | 2026-09-21 | GOAL-047 检查点 A/B（自举、快捷方式、hint、文档、reset→init→start 实测） | conditional | 7 项成果可核对；开放 required = 0；`F-S-001`～`F-S-004` recommended 提交 independent 复核 | `03-audit/A-001-self-w35-db-init.md` |

## 待复审事项（编排器登记，供独立审计取证）

| # | 事项 | 证据位置 | 说明 |
|--:|------|----------|------|
| 1 | 自举修复是否真的让工具能在**空实例**上创建第一个库，且不改变业务连接语义 | `cmd/e2e-pgset` | 反例优先：维护连接回退是否可能被误用于业务路径；回退顺序（`postgres` → `template1`）是否可靠 |
| 2 | 初始化快捷方式是否**幂等**且**不破坏既有库** | 快捷方式实现 + 实测 | 库已存在时的行为；`drop` 是否仍是显式独立动作 |
| 3 | 快捷方式是否**一次建齐 dev 与 test** 库，且与 `.env`/`config.yaml` 的库名来源一致（不硬编码） | 实现 + 文档 | 库名应从配置读取，而非写死 |
| 4 | 启动错误可操作化是否**不泄漏秘密**、且不改变既有错误分类/语义 | `internal/store` 或 `internal/composition` | 只加提示，不改分类；DSN 口令必须保持脱敏 |
| 5 | 文档章节是否与真实命令一致（可复制即用），并覆盖多方言差异（SQLite 无需初始化）与失败排查 | README/QUICKSTART | 命令与实现必须逐字一致 |
| 6 | 是否有可复跑的「删库 → 初始化 → `dev.cmd start` 全绿」证据 | 执行台账 | 不接受只有叙述 |
| 7 | 是否越界改动（迁移链/checksum/codec/wire/备份 Port、自动建库语义） | diff 范围 | 越界即 required |
