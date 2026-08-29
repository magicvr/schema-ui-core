---
id: GOAL-002-r1-serve-shell
title: R1 · serve 壳闭环（schema-ui serve · HTTP 壳 + config 装载 + assembly 服务器面 · RT-D02 接线）
status: done
parent: GOAL-001-distribution-formalization
created: 2026-08-29
updated: 2026-08-29
version: 0.3.0
progress: 5/5
---

# GOAL-002 · R1 · serve 壳闭环

## 概述

承接 Root R1 与 VP-024 判据 #1：把 VP-023 交付的「装配冒烟」骨架升级为**可运行的 HTTP 服务**——公开 serve 面（config 装载 + assembly 服务器面 + RT-D02 优雅停机接线），CLI `schema-ui serve` 与生成骨架共用同一实现，`create` 生成的骨架可直接 `serve` 启动。核销 go 后清单第 ① 项「`schema-ui serve` 壳」。

## 成功标准（可验证检查点）

- [x] C1：`schema-ui serve` 子命令存在；对 `create` 生成的骨架执行后可启动 HTTP 服务（/healthz、/readyz 可响应）——E2E-L1（CLI sqlite）+ E2E-L2（骨架薄封装 run，本地点缀）
- [x] C2：生成骨架 `cmd/server` 薄封装为同一 serve 面（`-dialect`/`-dsn`/`-config`/`-addr` 覆盖注入）——E2E-L2 编译+运行
- [x] C3：RT-D02 全序停机接线（`shutdown.starting` → drain → Store close → complete/错误）——`Run` ctx 干净排空单测 + config fail-closed；信号/退出码 = linux CI 登记（E-003 残余 1）
- [x] C4：双方言启动实证（SQLite E2E-L1/L2；PostgreSQL docker postgres:16 E2E-L3）
- [x] C5：config 装载 fail-closed（13 项单测：非法 shutdown_timeout / 裸 `${VAR}` / 方言配对 / 非 dev 密钥）；登录探针（E2E-L1/L2 login 200）

> 有界登记 ×2（E-003 残余节）：信号级 drain harness → R3 CI；registry 级骨架消费（无 replace）→ R2 发布核销。

## 方案与路线（P-001）

阶段串行：

| 阶段 | 内容 | 状态 |
|------|------|------|
| S1 | 设计冻结：serve 面构成（D-001 · 关键决策 I-001）｜公开面放置｜config 装载形态 | **已关门**（2026-08-29 · D-001 · 用户 P-004 裁决：方案 A + 薄封装） |
| S2 | 实现 serve 面（config 装载 + 组合装配 + 中央面接线 + RT-D02 生命周期） | **已关门**（2026-08-29 · E-003 证据） |
| S3 | CLI `serve` 子命令 + 模板 main 薄封装 + 骨架 config 示例 | **已关门**（2026-08-29 · E-003 证据） |
| S4 | 验收证据（双方言启动 / drain harness A·B / 探针）+ 自审关门 | **已关门**（2026-08-29 · A-001 self `pass` · 有界登记 ×2） |

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | serve 面构成范围（方案 A 标准下游组合闭环 / B 全量对等 / C 最小壳）与模板 main 形态（薄封装 vs 双入口） | S1 设计冻结 | S1 | 用户 P-004 裁决（附推荐） | **verified**（2026-08-29 用户裁决：方案 A + 薄封装单一形态） | — | D-001-r1-scope-and-wiring |
| I-002 | non-blocking | config 装载默认形态（内嵌默认 + 显式文件覆盖 + env） | S2 实现 | S2 | 参照内部 config.Load 语义镜像 | **verified**（2026-08-29 · D-001 §4 定案） | — | D-001-r1-scope-and-wiring |

## 父目标

- `GOAL-001-distribution-formalization`

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger。

## 备注

- 审计模式：S4 关门 default self（D-001 已定）；如出现 release/兼容语义变更再升 independent。
- 关联：original 挂账 = VP-023 go 后清单 ①；Root 大纲 R1；I-024-004（serve 接线方式）于本目标闭合。