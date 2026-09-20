---
id: D-004-r2-migration-ownership-user-decision
doc: decision-entry
status: accepted
parent: null
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# D-004 · R2 migration 归属与 catalog 形态裁决

## 用户裁决

R2 采用**按模块追加 migration**：每个拥有时间表的模块追加自己的 conversion migration version/`Apply`/`ApplyPostgres`；公共 codec/测试契约可放在共享 internal package；保持模块表所有权与 compiled catalog append-only。

## 强制不变量

- 不修改 v1–v72 的 canonical SQL、checksum、identity 或历史 Apply 语义。
- 新 conversion migration 从 v73 之后追加，并在 compiled catalog 中保持版本连续、checksum 唯一、两方言成对 Apply。
- 模块只迁移自己拥有的表；Store runner 只负责执行/ledger/snapshot，不集中接管模块表 DDL。
- R2 设计必须列出 `CHECK`、partial index、NULL/default、WHERE/ORDER/predicate 的重建顺序；不能只改列类型。
- 公共 wire 统一 6 位 RFC3339 的 formatter/parser 影响单列为 C2/C4/R3 交付，不泄漏 pgx/driver 类型。

## 未选方案

- 单一 platform conversion migration：会跨越模块表所有权，除非后续用户重新裁决，不采用。
- 修改历史 DDL 让 fresh bootstrap 直接生成新形状：会破坏已 apply catalog 的 checksum fail-closed 合同，不采用。
