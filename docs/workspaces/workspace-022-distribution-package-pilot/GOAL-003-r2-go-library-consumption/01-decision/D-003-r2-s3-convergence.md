---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-003-r2-go-library-consumption
version: 0.1.0
---

# D-003 · C 层泄漏收敛方案提案（关键决策 · 待用户裁决）

背景（E-003）：泄漏面 = **类型命名面**（非包加载）。模块可外部 import（users exit 0）；下游唯一命名需求 = `*auth.Authenticator`（`internal/auth`，17 处引用）；`handler`（18 处）是模块内部实现细节无需命名；`store` 经 `authsession.TxRunner`（B 层接口）间接。**Go 类型推断允许不命名消费**——`a := assembly.NewAuthenticator(cfg)` 后可直接传参。

## 方案

| 方案 | 内容 | 冻结面影响 | 成本 | 判定 |
|------|------|-----------|------|------|
| **β · 公开装配工厂（推荐）** | 新增 `apps/api/assembly` 公开包（同模块树可 import `internal/*`），导出最小工厂集（`NewTxRunner` / `NewAuthenticator` / `NewMailSender` 等，输入 = config/kernel/B 层接口）；下游组合根 = kernel + assembly + 模块自选，**零外移、零契约修订** | 新增 B+ 公开面（additive · minor 合规） | 一个包 + golden-consumer 装配验证 | ★ |
| **α · 有限外移** | `internal/auth` + `internal/store`（+依赖闭包）提升为公开路径 | 公开面显著扩大（auth/store 入 semver 约束），与「薄基架面」哲学有张力 | 中（含依赖闭包盘点） | 备选（go 后评估） |
| **β' · kernel 接口化** | 模块构造签名改 kernel 级接口（长期纯净） | 契约修订（breaking 流程） | 高（22 模块签名全过） | 长期演进候选，本期不做 |

## 推荐 β 的落地范围（本目标 S3）

1. `apps/api/assembly` v0.0.0：工厂最小集 = users 装配链所需（`NewTxRunner(cfg) authsession.TxRunner`、`NewAuthenticator(cfg) *auth.Authenticator`、`NewMailSender(cfg) kernel.MailSender`）；内部 import `internal/store`/`internal/auth`/`internal/mail`/`internal/config`。
2. golden-consumer：users 全链装配（authn + authsession.Repository + operationlog.Recorder + MailSender）→ `go build` + 启动冒烟（registry/contribution 注册）+ **双方言迁移 apply**（SQLite + PG 可选）。
3. 冻结面 v1.1.0：增列 B+ 层（assembly 导出签名）+ 回填 B 层符号盘点（F-002）。
4. 不预改 kernel/module 构造签名；α/β' 留 go/no-go 后评估。

## 未选方案

- 不采用「下游全量自建基架装配」（需外移全部运行底座，公开面爆炸）。
- 不采用纯 β'（本期成本高、无外部消费者、可后置）。