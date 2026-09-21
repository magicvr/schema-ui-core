---
id: r3c-pg-tool-compatibility-probe-v0.1
doc_type: evidence-attachment
title: R3-C 前置探测：PG 客户端工具跨版本兼容行为（15.19 / 16.15 / 17.11）
status: recorded
created: 2026-09-21
updated: 2026-09-21
parent: GOAL-007-r3-pg-cross-version-restore-matrix
version: 0.1.0
---

# R3-C 前置探测：客户端工具跨版本兼容行为

本附件记录**先于**矩阵定义（`D-001`）所做的一次探测，目的是在冻结判定口径前先取得实测事实，而不是按预期反推口径。探测使用**平凡 schema**（单表 `probe(id int primary key, at timestamptz(6))` + 一行 `.900000Z`），因此它只回答「工具本身怎么行为」，不代表 VP-040 真实 schema 的结论——后者由矩阵本体（`r3c-pg-cross-version-matrix-v0.1.md`）承担。

## 环境与版本（实测，非沿用文档）

| 角色 | 版本 | 取得方式 |
|------|------|----------|
| server（容器 15） | **15.19** | `docker run --rm postgres:15-alpine postgres --version` |
| server（容器 16） | **16.15**（Debian 16.15-1.pgdg13+2） | `docker run --rm postgres:16 postgres --version` |
| server（容器 17） | **17.11** | `docker run --rm postgres:17-alpine postgres --version` |
| server（常驻，`PG_TEST_*`） | **15.4**（Debian 15.4-2.pgdg120+1） | `psql -tAc "select version()"` |
| client 工具 15 | **pg_dump / pg_restore 15.19** | `postgres:15-alpine` |
| client 工具 16 | **pg_dump / pg_restore 16.15** | `postgres:16` |
| client 工具 17 | **pg_dump / pg_restore 17.11** | `postgres:17-alpine` |

机制：`docker network create` + `--network-alias`，server 与 client 容器同网；产物经命名卷 `/vp040` 交换；就绪门用 `pg_isready`。宿主无 `psql`/`pg_dump`/`pg_restore`（沿用 `D-017` §3 约束③）。

## 归档格式版本（直接读头部字节）

命令：`head -c 7 /vp040/<archive> | od -An -tu1` → `PGDMP` + 版本三元组。

| 产物（由同 major 的 client 生成） | 头部字节 | 归档格式版本 |
|-----------------------------------|----------|--------------|
| `a15.dump`（server 15.19，client 15.19） | `80 71 68 77 80 1 14` | **1.14** |
| `a16.dump`（server 16.15，client 16.15） | `80 71 68 77 80 1 15` | **1.15** |
| `a17.dump`（server 17.11，client 17.11） | `80 71 68 77 80 1 16` | **1.16** |

## 1. `pg_dump`：client 与 server 的 major 关系

| server | client 15.19 | client 16.15 | client 17.11 |
|--------|--------------|--------------|--------------|
| 16.15 | **拒绝**：`pg_dump: error: aborting because of server version mismatch` / `detail: server version: 16.15 …; pg_dump version: 15.19`（exit 1） | 接受（exit 0） | 接受（exit 0） |

结论（实测）：**`pg_dump` 要求 client major ≥ server major**；低版本 client 直接拒绝，不是「尽力而为」。

## 2. `pg_restore` 读取归档：client 与归档格式版本的关系

命令：`pg_restore -l <archive>`（不连库，仅读归档）。

| 归档 | client 15.19 | client 16.15 | client 17.11 |
|------|--------------|--------------|--------------|
| `a15.dump`（1.14） | 可读 | 可读 | 可读 |
| `a16.dump`（1.15） | **拒绝** `unsupported version (1.15) in file header` | 可读 | 可读 |
| `a17.dump`（1.16） | **拒绝** `unsupported version (1.16) in file header` | **拒绝** `unsupported version (1.16) in file header` | 可读 |

结论（实测）：**`pg_restore` 要求 client major ≥ 生成归档的 client major**（归档格式版本随 dumper major 递增：1.14 / 1.15 / 1.16）。

## 3. `pg_restore` 写入目标 server：9×3 组合实测

每格新建一个空 database 后执行 `pg_restore --exit-on-error --no-owner`。

| 归档 | client | → server 15.19 | → server 16.15 | → server 17.11 |
|------|--------|----------------|----------------|----------------|
| `a15.dump` | 15.19 | 成功 | 成功 | 成功 |
| `a15.dump` | 16.15 | 成功 | 成功 | 成功 |
| `a15.dump` | 17.11 | **失败** `unrecognized configuration parameter "transaction_timeout"` | **失败** 同上 | 成功 |
| `a16.dump` | 15.19 | 失败（header 1.15） | 失败（header 1.15） | 失败（header 1.15） |
| `a16.dump` | 16.15 | 成功 | 成功 | 成功 |
| `a16.dump` | 17.11 | **失败** `transaction_timeout` | **失败** `transaction_timeout` | 成功 |
| `a17.dump` | 15.19 | 失败（header 1.16） | 失败（header 1.16） | 失败（header 1.16） |
| `a17.dump` | 16.15 | 失败（header 1.16） | 失败（header 1.16） | 失败（header 1.16） |
| `a17.dump` | 17.11 | **失败** `transaction_timeout` | **失败** `transaction_timeout` | 成功 |

结论（实测）：**`pg_restore` 要求 dumper major ≤ client major ≤ 目标 server major**。两个独立失败模式：

1. **归档格式门**（client < dumper）：`unsupported version (1.1x) in file header`。
2. **server 端 GUC 门**（client 17 → server 15/16）：17 的客户端工具在恢复会话中设置 `transaction_timeout`（PG 17 引入的 GUC），15/16 server 不识别该参数；`--exit-on-error` 下立即失败。注意这一门**与归档由哪个版本生成无关**：`a15.dump`（1.14）用 17 client 恢复到 s15/s16 同样失败。

## 4. 与本工作区既有实测的关系（避免重复劳动）

`D-015` 记载：毫秒族 PG 表达式曾于 **PG 15.19 / 16.15 / 17.11** 三版本独立复现并更正（A-044）。那是**表达式/DDL 层**的实测（`psql` 执行 SQL），**不**覆盖：完整迁移链在 16/17 上的应用、`pg_dump`/`pg_restore` 跨版本组合、恢复后的 canonical 形状校验。R3-C 只承担后者。

## 5. 探测清理

探测创建的资源（`vp040c-net` 网络、`vp040c-s15/s16/s17` 容器、`vp040c-vol` 卷）在矩阵驱动落地前删除，避免与驱动自建资源冲突；常驻 15.4 实例未被本探测修改（只读 `version()` 与容器内新建临时 DB）。

## 边界

- 本附件的结论**仅**由本次实测支持，且限定在列出的三个镜像版本与上述命令形态；镜像 tag 可变，矩阵本体必须**每个组合重新记录版本串**。
- `--exit-on-error` 会放大差异：不带该旗标时 `transaction_timeout` 只是 per-statement 失败。矩阵口径（`D-001`）以 C3 实际使用的旗标为准。
