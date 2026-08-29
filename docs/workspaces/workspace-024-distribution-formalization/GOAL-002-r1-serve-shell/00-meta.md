---
id: GOAL-002-r1-serve-shell
title: R1 · serve 壳闭环（schema-ui serve · HTTP 壳 + config 装载 + assembly 服务器面 · RT-D02 接线）
status: active
parent: GOAL-001-distribution-formalization
created: 2026-08-29
updated: 2026-08-29
version: 0.1.0
progress: 0/5
---

# GOAL-002 · R1 · serve 壳闭环

## 概述

承接 Root R1 与 VP-024 判据 #1：把 VP-023 交付的「装配冒烟」骨架升级为**可运行的 HTTP 服务**——公开 serve 面（config 装载 + assembly 服务器面 + RT-D02 优雅停机接线），CLI `schema-ui serve` 与生成骨架共用同一实现，`create` 生成的骨架可直接 `serve` 启动。核销 go 后清单第 ① 项「`schema-ui serve` 壳」。

## 成功标准（可验证检查点）

- [ ] C1：`schema-ui serve` 子命令存在；对 `create` 生成的骨架执行后可启动 HTTP 服务（/healthz、/readyz 可响应）
- [ ] C2：生成骨架 `cmd/server` 薄封装为同一 serve 面（`-dialect`/`-dsn` 保留为覆盖注入；`-config` 支持外部配置）
- [ ] C3：RT-D02 全序停机接线：SIGTERM → `shutdown.starting` → drain（`http.shutdown_timeout` 预算，默认 10s）→ Store close → `shutdown.complete`/exit 0；预算耗尽 → exit 1
- [ ] C4：双方言启动实证（SQLite 必做；PostgreSQL 经 docker 复现或登记环境限制并给等价口径）
- [ ] C5：config 装载 fail-closed（非法 `http.shutdown_timeout ≤ 0` 或裸 `${VAR}` 未设置 → 拒绝启动）；探针（token-check 等）经 serve 登录路径 PASS

## 方案与路线（P-001）

阶段串行：

| 阶段 | 内容 | 状态 |
|------|------|------|
| S1 | 设计冻结：serve 面构成（D-001 · 关键决策 I-001）｜公开面放置｜config 装载形态 | 待用户裁决 |
| S2 | 实现 serve 面（config 装载 + 组合装配 + 中央面接线 + RT-D02 生命周期） | 依赖 S1 |
| S3 | CLI `serve` 子命令 + 模板 main 薄封装 + 骨架 config 示例 | 依赖 S2 |
| S4 | 验收证据（双方言启动 / drain harness A·B / 探针）+ 自审关门 | 依赖 S2/S3 |

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | serve 面构成范围（方案 A 标准下游组合闭环 / B 全量对等 / C 最小壳）与模板 main 形态（薄封装 vs 双入口） | S1 设计冻结 | S1 | 用户 P-004 裁决（附推荐） | open | S1 前必须闭合 | 待确认 |
| I-002 | non-blocking | config 装载默认形态（内嵌默认 + 显式文件覆盖 + env） | S2 实现 | S2 | 参照内部 config.Load 语义镜像 | open | — | 待确认 |

## 父目标

- `GOAL-001-distribution-formalization`

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger。

## 备注

- 审计模式：S4 关门 default self（D-001 已定）；如出现 release/兼容语义变更再升 independent。
- 关联：original 挂账 = VP-023 go 后清单 ①；Root 大纲 R1；I-024-004（serve 接线方式）于本目标闭合。