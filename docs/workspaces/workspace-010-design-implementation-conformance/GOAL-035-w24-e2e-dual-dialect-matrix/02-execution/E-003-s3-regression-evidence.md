---
title: E-003 · S3 回归事实与连带整改（配置双载发现）
status: recorded
created: 2026-08-23
updated: 2026-08-23
parent: GOAL-035-w24-e2e-dual-dialect-matrix
version: 0.1.0
---

# E-003 · S3 回归事实（2026-08-23）

## 全量回归矩阵（最终代码态）

| 项 | 结果 | 说明 |
|----|------|------|
| e2e sqlite（默认契约，`npm run test:e2e`） | **9/9 passed**（1.7m） | 全新临时 SQLite；globalSetup 校验 DB 文件出现 |
| e2e postgres（`npm run test:e2e:postgres`） | **9/9 passed**（1.8m） | 自动 create → verify → run → drop，teardown 输出 `dropped schema_ui_e2e_*`；事后 `e2e-pgset list` = **0 遗留** |
| vitest | 76 文件 / **1088/1088** | 与 W23 基线一致 |
| go test ./... / go build / go vet | 全绿 | 含新 `cmd/e2e-pgset` |
| tsc -b --noEmit + vite build | exit 0 | conformance claim buildId 随构建刷新 |

## 连带发现与整改

**F-1（配置被求值两次 → 双份 scratch 库）→ fixed**
- 现象：首次 pg 全量跑后 `list` 发现 2 个 `schema_ui_e2e_*` 库，teardown 只处理 1 个。
- 根因：Playwright 真实运行会**两次求值 config**（主进程 + worker 配置序列化，两个独立进程；`--list` 只求值一次，插桩证实：pid/ppid 链 64008→40012，第二次为第一次的子进程）。config 顶层 provisioning 副作用因此执行两次；且 teardown 的 drop 在 API 仍持连接时执行会静默失败（无 FORCE）。
- 修复：
  1. `playwright.config.ts` 用 `E2E_PG_NAME` 环境守卫——第二次求值（子进程）继承首次非空值即**复用**同一 scratch 库，不再二次 create（插桩确认 `inh` 继承路径生效）；
  2. `cmd/e2e-pgset drop` 增加 `WITH (FORCE)`（PG≥13），失败回退普通 DROP——API 活跃连接不再阻塞清理；
  3. `global-teardown` 改 `stdio: inherit`——drop 失败可见，配合 `list`/`drop` 可手动清理。
- 验证：修复后两次 pg 运行均 1 create → 1 drop → **遗留 0**。

## 契约不变量复核（防回退）

- API webServer env 恒含显式 `DB_DIALECT` → `.env` 无法改道（process env 优先）；
- `globalSetup` 校验：sqlite = DB 文件出现；postgres = `verify`（schema_migrations 就位）；违反整体 fail-fast 并给出 W23 指向性诊断；
- CI `browser-e2e` = `profile × dialect` 矩阵，双方言每条 PR 强制执行。

## 结论

- C3 达成：双方言全量绿 + 密闭回归全绿 + provisioning 生命周期（建/验/删）闭环。
- 无产品语义/协议变化；关门审计模式 `self`（按 00-meta 边界声明）。