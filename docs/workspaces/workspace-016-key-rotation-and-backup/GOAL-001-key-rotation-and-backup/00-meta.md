---
id: GOAL-001-key-rotation-and-backup
title: 密钥轮换与备份恢复合同（JWT + 轮换后恢复）
status: active
parent: null
created: 2026-08-22
updated: 2026-08-22
version: 0.4.0
progress: 3/5
plan_refs:
  - VP-016-key-rotation-and-backup
primary_plan: VP-016-key-rotation-and-backup
serves_summary: 交付架构 A5：JWT current+previous 轮换合同 + 既有备份上的轮换后恢复；单密钥仍为 dev/mvp/快测默认。不承载 KMS / PITR / A3 / Admin 密钥页或业务域。
---

# GOAL-001 · 密钥轮换与备份恢复合同（JWT + 轮换后恢复）

## 概述

本 Root 承载 [VP-016-key-rotation-and-backup](../../../vision/plans/VP-016-key-rotation-and-backup.md)（**`active`**）的实现：在已有 YAML + env 密钥注入 fail-closed 与 VP-013 方言级 dump 之上，补齐 JWT 双密钥轮换合同，并核对轮换后从既有备份启动仍能鉴权。

**边界**：不强制本地默认改成必须有 previous 密钥或外部备份代理；不承接 KMS、PITR、热加载、Admin 密钥页或业务域表。不重做 `pg_dump`/`VACUUM INTO`。安全 finding → VP-009；符合性 gap → VP-010。

## 纲领路线图（P-001）

| 阶段 | 内容 | 先后 | 状态 |
|------|------|------|------|
| R1 | **轮换合同与配置面冻结**：current/previous 键名、生产 fail-closed、熵规则（I-001）；本波密钥集合是否仅 JWT（I-002）；缺省单密钥。 | 起点 | **完成**（D-002 + GOAL-002 A-001 pass；配置面落地） |
| R2 | **JWT 双密钥实现**：签发只用 current；校验 current 再 previous；重叠窗 / `kid` / refresh 不受签名密钥影响（I-003）；重启生效。 | 依赖 R1 | **完成**（GOAL-003 done · self A-001 + independent A-002 双 pass） |
| R3 | **轮换后恢复证据**：在既有 SQLite `VACUUM INTO` 与 PG `pg_dump`/`pg_restore` 上核对轮换后启动 + 鉴权（I-004）。不重做 dump。 | 依赖 R2 | **完成**（GOAL-004 done · A-001 self pass；双方言循环全绿） |
| R4 | **默认单密钥仍可用**：未配置 previous 时本地/Compose 仍能开发与快测；轮换不是启动硬依赖。 | 依赖 R2 | 未开始 |
| R5 | **双路径证据**：显式双密钥下，一轮换路径 **与** 一轮换后恢复路径都有可核对证据。 | 依赖 R3/R4 | 未开始 |

`progress` = 已完成阶段数 / 5。当前 **2/5**（R1、R2 完成）。

## 成功标准（方向级）

1. JWT 轮换合同落地：可配置 current + previous；新签发只用 current；重叠窗内 previous 可验 access。
2. 未配置 previous 时本地/Compose 默认仍能开发与快测。
3. 轮换后恢复：两方言既有备份路径上，轮换后从备份启动且鉴权可核对。
4. 显式双密钥配置下，一轮换路径 **与** 一轮换后恢复路径都有可核对证据。
5. 未进入 A3 / KMS / PITR / Admin 功能 / 业务域；未改 Charter；未假装交付热加载或第二套 dump。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | current / previous 配置键名、生产 fail-closed、熵规则沿用 `ValidateProd`（≥32 字符且同时含字母与数字）；secret 不入库、不进日志 | R1 方案冻结 / 实施 | R1 合同冻结 | R1 决策 | **verified**（D-002） | — | D-002 §1；`config.go:202/496/953-964`、`main.go:74-85` |
| I-002 | required | 本波密钥集合仅 `AUTH_JWT_SECRET`（+等价 previous 键）。服务凭证为 SHA-256 opaque hash，不与 JWT secret 共用；D-002 书面出局 | R1 方案冻结 | R1 合同冻结 | R1 决策 | **verified**（D-002） | — | D-002 §2；`auth.NewServiceCredentialToken` = CSPRNG + SHA-256 |
| I-003 | required | 重叠窗 = previous 配置存续期（退役 ≥access_ttl 后移除并重启）；不用 JWT `kid`；refresh 为 opaque SHA-256，不受签名密钥轮换影响 | R2 方案冻结 / 实施 | R2 接入前 | R2 决策 | **verified**（GOAL-003 D-001） | — | GOAL-003 D-001；`auth.NewOpaqueToken`/`HashToken`/`RefreshTokenByHash` |
| I-004 | required | 轮换后恢复最小剧本：备份点在轮换前（K1 运行中），恢复后以 K2+prev=K1 启动并断言 A1 旧 access 可验 / A2 新签发仅 current / A3 refresh opaque 连续；SQLite `VACUUM INTO`、PG `pg_dump -F c`→`pg_restore`，不重做 dump | R3 方案冻结 | R3 接入前 | R3 决策 | **verified**（GOAL-004 D-001） | — | GOAL-004 D-001；VP-013 备份合同（GOAL-006 D-002）沿用 |
| I-005 | non-blocking | 重叠窗内旧 access 立即失效是否接受为有界残余。默认：previous 可验 | 退出 1 措辞 | R2 | 用户书面残余时才改变退出 1 | collecting | — | 对应 VP I-016-005 |

## 父目标

- null（Root；Charter `schema-ui-core-admin-foundation@0.2.0` / VP-016）

## 台账布局

新目标为三个可追加台账创建同名平铺目录：`01-decision/`、`02-execution/`、`03-audit/`。索引文件保留 frontmatter、摘要和条目索引；独立记录使用 `D-NNN-*`、`E-NNN-*`、`A-NNN-*` 文件。
