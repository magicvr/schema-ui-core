---
id: GOAL-001-object-storage
title: 对象存储适配器（S3 兼容 + 本地盘内嵌）
status: active
parent: null
created: 2026-08-21
updated: 2026-08-21
version: 0.2.0
progress: 1/5
plan_refs:
  - VP-014-object-storage
primary_plan: VP-014-object-storage
serves_summary: 交付架构 A2：内核对象存储端口 + S3 兼容实现；本地盘保留为 dev/mvp/快测默认。不承载签名 URL / 分片 / 扫描 / CDN / 产品搬运器，也不承载 Admin 功能或业务域。
---

# GOAL-001 · 对象存储适配器（S3 兼容 + 本地盘内嵌）

## 概述

本 Root 承载 [VP-014-object-storage](../../../vision/plans/VP-014-object-storage.md)（**`active`**）的实现：把现行本地盘-only 文件落盘收成内核对象存储端口，补齐 S3 兼容实现，并把现有三类落盘（avatars / brand-assets / uploads）改走同一端口。

**边界**：不强制本地默认改成必须有 MinIO/S3；不承接签名 URL / 分片 / 扫描 / CDN / 产品搬运器；不加 Admin 页面或业务域表。安全 finding → VP-009；符合性 gap → VP-010。

## 纲领路线图（P-001）

| 阶段 | 内容 | 先后 | 状态 |
|------|------|------|------|
| R1 | **端口与配置面冻结**：公共类型无本地路径 / `os.File`；命名空间隔离三类落盘；缺省本地盘；S3 为显式配置；裁定 List/GC 是否进端口（I-005）与桶模型（I-002）。承载子目标：GOAL-002-object-port-freeze | 起点 | **已完成**（2026-08-21：实现 + self/independent 审计闭合，见其 E/A 台账） |
| R2 | **S3 兼容接入**：驱动公约数（I-001）、凭证键名（I-003）、配置后 `readyz` 扩依赖。承载子目标：GOAL-003-object-s3-driver | 依赖 R1 | 进行中（D-004/D-005 已裁，方案见 GOAL-003 D-001） |
| R3 | **三类落盘收口**：avatars / brand-assets / uploads（含 file-library 与 data-transfer 共享上传目录）走同一端口 | 依赖 R2 | 未开始 |
| R4 | **公共面收口**：Handler / 模块公共契约去掉本地路径与 `os.File` | 依赖 R1；可与 R3 部分并行 | 未开始 |
| R5 | **双路径证据**：本地盘默认路径回归 + S3 兼容生产向验收（配置接入、读写删除、就绪探针） | 依赖 R3/R4 | 未开始 |

`progress` = 已完成阶段数 / 5。当前 `1/5`（R1 完成）。progress 不放行、不关门。

## 成功标准（方向级）

1. 内核对象存储端口落地；公共契约无本地路径 / `os.File`。
2. S3 兼容实现对三类落盘可核对 put / get / delete；显式配置时 `readyz` 扩依赖。
3. 本地盘仍为本地/Compose 默认；无对象存储仍能开发与快测。
4. 生产向验收以 S3 兼容为准（配置接入、读写删除、就绪探针之一可核对）。
5. 未引入第三方言；未改 Charter；未进 Admin 功能/业务域；未假装交付签名 URL / 分片 / 扫描 / CDN / 产品搬运器。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | S3 API 子集与驱动：MinIO / R2 / AWS 的最低公约数；禁止第三对象存储方言 | R2 方案冻结 / 实施 | R2 实施前 | R2 决策 + 编译/接入证据 | verified（2026-08-21） | — | D-004：aws-sdk-go-v2；Put/Get/Delete/Head/ListV2+HeadBucket；对应 VP I-014-001 |
| I-002 | required | 桶模型：单桶 + 前缀 vs 多桶；三类落盘的 key 隔离规则 | R1 方案冻结 | R1 端口冻结 | R1 决策 | verified（2026-08-21） | — | D-002：单桶 + 命名空间前缀；key 为 namespace/32hex；对应 VP I-014-002 |
| I-003 | required | 配置键名与凭证注入（YAML + env fail-closed；secret 不入库） | R2 方案冻结 | R2 实施前 | R2 决策 | verified（2026-08-21） | — | D-005：键名沿 GOAL-002 D-001；static credentials 显式构造；对应 VP I-014-003 |
| I-004 | non-blocking | 存量本地文件如何进入对象存储？ | R5 关门叙事 | R5 | 点名 residual | **recorded** | 用户已裁决不进退出分母 | 不提供产品搬运器；既有存量 = 继续本地或运维自备拷贝（对应 VP I-014-004） |
| I-005 | required | 启动 GC / 列举是否属于端口（List），还是本地适配器私有 + DB 引用集？ | R1 方案冻结 | R1 端口冻结 | R1 决策 | verified（2026-08-21） | — | D-003：List+Stat 进端口，GC 与配额留应用层经原语实现；uploads 无 DB 引用集，排除仅靠 DB 方案 |

## 父目标

- null（Root；Charter `schema-ui-core-admin-foundation@0.2.0` / VP-014）

## 台账布局

新目标为三个可追加台账创建同名平铺目录：`01-decision/`、`02-execution/`、`03-audit/`。索引文件保留 frontmatter、摘要和条目索引；独立记录使用 `D-NNN-*`、`E-NNN-*`、`A-NNN-*` 文件。
