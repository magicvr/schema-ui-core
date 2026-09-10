---
doc_type: goal-attachment
id: roadmap-restatement-draft
parent: GOAL-005-r4-roadmap-draft-and-close
status: draft
created: 2026-09-10
updated: 2026-09-10
version: 0.1.0
---

# 下一版总路线图重述草案（交 `/vision` editorial）

本草案是 VP-035 判据 4 的交付物：把 R2 as-built 对照（17 面 + W1）与 R3 有界业界对照/缺口分类的结果，重述为下一版 `docs/vision/roadmap.md` 的草案。**它不是已冻结权威**：`docs/vision/roadmap.md` 的正文改动须经 `/vision` editorial（VRev 记录）后，与本文同一事务或紧随冻结执行。

- 输入证据：R2 [as-built 矩阵 v0.3.0](../../GOAL-003-r2-as-built-matrix/attachments/as-built-matrix.md)、R3 [业界对照表](../../GOAL-004-r3-industry-comparison/attachments/industry-comparison.md)、[缺口分类表 v0.2.0](../../GOAL-004-r3-industry-comparison/attachments/r3-gap-classification.md)、[I-035-003 判定](../../GOAL-004-r3-industry-comparison/attachments/r3-i035-003-determination.md)
- 边界：不改 Charter；不新增/消耗 trigger-gated 行；不实现任何「现在修」代码项；不重开已 closed VP

## 1 · 现状锚点（修正版）

**现行路线图第 104 行的锚点已过期**（G-002）。修正后的现状锚点应为：

> 现状锚点：**单进程** + **SQLite 文件库（小连接池，默认 4；内存库 1）** / **PostgreSQL 双方言** + 本地盘上传（S3 兼容适配器已交付） + 进程内 Job（六态） + 内存限流（原子窗口） + 进程内事件总线 + 可选 Prometheus 指标与 OTLP traces（缺省关闭） + JWT current/previous 轮换合同。Compose 已声明非目标含 TLS 终止与多实例（`compose.yaml`）。

与旧锚点的差异：`MaxOpenConns=1` → 小连接池（`apps/api/internal/store/store.go:29` `sqlitePoolDefault = 4`，`:104`–`113`）；补充 A2/A4/A5/A6/A7 已交付事实。

## 2 · A 序列（架构）现状

| 项 | 内容 | 现行状态（本草案口径） | 备注 |
|----|------|------------------------|------|
| A0 | 有界清单登记 | done | 本清单 + Store 双方言决策（RT-P03）已冻结 |
| A1 | 内核持久化端口 + PG 实现 + 台账对写 | **delivered** | VP-013 `closed` v0.3.0（2026-08-21）；residual = 无产品搬运器 |
| A2 | 对象存储适配器（本地盘默认） | **delivered** | VP-014 `closed` v0.3.0；residual = 无产品搬运器 |
| A3 | 多实例前置（就绪探针扩依赖；PG 锁/`SKIP LOCKED` vs Redis vs 外部队列评估） | **仍 trigger-gated** | 触发 = 多实例部署或 C 端业务域接入（RT-Q03/Q05 注）；优雅停机/排空已从 A3 拆出为 A7 |
| A4 | 指标 + OpenTelemetry | **delivered** | VP-015 `closed` v0.3.0；residual = otlp-sink 不解析、Store/对象/Job 指标不进分母 |
| A5 | 密钥轮换 / 备份恢复合同 | **delivered** | VP-016 `closed` v0.3.0；residual = 立即失效未选（**未获书面接受**）、mfa-wrap 窄口径（用户 2026-09-10 接受） |
| A6 | 出站邮件：发送端口 + 可切换渠道 | **delivered** | VP-017 `closed` v0.5.0（现行渠道分母） |
| A7 | 优雅停机 / 连接排空合同 | **delivered** | VP-021 `closed` v0.3.0；residual = 进程级 harness 以 linux CI 核销 |

**结论**：A 序列除 A3（多实例前置）外已全部交付。下一版路线图应把 A0–A7 从「建议顺序」改写为「已交付序列 + 唯一未触发项 A3」，并把新出现的候选（见 §5）单列，不再沿用「A0–A7 未冻结草案」的表述。

## 3 · 已交付 vs 下一拍（三分支）

### 架构分支

- **已交付**：A1/A2/A4/A5/A6/A7（见 §2）；RT-P03/P04-池化/RT-D02/RT-K03/RT-M03/RT-Q01/RT-S01/RT-O01-O02 等行均已 delivered。
- **下一拍（候选，未立项）**：RES-T03-tz（DB `timestamptz` 持久化合同；本轮证据：全仓 `*.sql` 中 `timestamptz` 命中 0、时间列仍 INTEGER）。
- **仍 gated（不消耗 trigger）**：A3、RT-Q02/Q03/Q05/Q06/Q07、RT-P06、RT-X01（搜索）、RT-S05（扫描执行器）、KMS/HSM。
- **明确不做**：RT-P07（文档库）、RT-P08（多租户物理隔离）、RT-N01–N05（ORM/先上 MQ/微服务拆分/GraphQL 网关/多云 K8s）。

### Admin 功能分支

- **已交付**：VP-011 四档、VP-012 横切契约、VP-017/018/019（邮件→邮箱身份→IAM 恢复）、VP-020（时区/数字/货币语义）、VP-025（配置导出/diff/dry-run/导入）、VP-029/030/031/033/034 等。
- **下一拍（非门控、未立项）**：体验增强（全局搜索 / Command Palette、Saved Views、批量结果中心、未保存保护、统一 Toast）。
- **仍 gated**：组织/部门/岗位与数据权限 `org`（触发 = 多组织 fork 或真实需求）；typed domain event、Notification Transport、OIDC/SSO/SCIM、Approval Gate、通用 Entitlement、SSE/WebSocket（VP-033 注记不解除）、外部连接器产品面、自定义 metadata、文件预览；文件扫描/隔离**策略**。

### 业务域分支

- **已立项**：VP-029（钱包预付费，Admin 资金通道）、VP-031（数字 Offer + 本域权益，`closed` v0.3.4）。
- **下一拍**：**无新触发则不要预开第二域**；候选 1～9 保持登记，一域一 VP，不得合并立项。
- **激活前置**：`/vision` 复核 + VP-008 `go` 消费有效性 + freshness review。

## 4 · residual 总账（R3 分类冻结，18 条）

| 分类 | 条数 | 条目 | 去向 |
|------|------|------|------|
| 现在修 | 6 | RES-T03-tz；G-001、G-002、G-003、G-005；G-006 | G-006 已执行（矩阵 v0.3.0）；G-001/G-003 走 R4 文档卫生；G-002/G-005 随本草案 `/vision` editorial；RES-T03-tz 登记为下一拍候选（本 VP 不实现） |
| 仍 gated | 3 | RES-015-metrics、RES-026-redis、RES-028-broker | 保持 RT 行状态，不消耗 trigger |
| 接受残余 | 4 | RES-015-otlp-sink、RES-021-harness、RES-030-keyfile（继承原 VP 留痕）、RES-016-mfa-wrap（用户 2026-09-10 按修正窄口径接受） | 复审触发见分类表 |
| 明确不做 | 5 | RES-013-migrator、RES-014-migrator、RES-016-revoke、RES-P04-pool、G-004 | 非目标或已声明的边界 |

明细与证据锚点见 [r3-gap-classification.md](../../GOAL-004-r3-industry-comparison/attachments/r3-gap-classification.md)（v0.2.0，经 independent A-005 `pass`）。

## 5 · 建议的路线图正文改动清单（供 `/vision` editorial 逐项裁决）

| # | 目标文件 | 锚点 | 现行文本 | 建议改为 | 依据 |
|---|----------|------|----------|----------|------|
| 1 | `docs/vision/roadmap.md` | `:104` | 现状锚点含 `MaxOpenConns=1` | §1 修正版锚点 | G-002 / R2 矩阵行 04 |
| 2 | `docs/vision/roadmap.md` | `:138` | RT-P04 现状 = `MaxOpenConns=1` | 「小连接池（默认 4）/内存库 1；读写分离与 replica 仍 trigger-gated」 | G-002；**不得**扩写为读写分离已交付 |
| 3 | `docs/vision/roadmap.md` | `:209` | RT-D02 行内「进程生命周期有，无明确 drain 合同」+ delivered | 删除自相矛盾表述，保留 delivered 与 VP-021 证据链 | G-002 |
| 4 | `docs/vision/roadmap.md` | `:290`–`304` | 「架构分支建议顺序（草案，未冻结）」A0–A7 | 改为「已交付序列 A1/A2/A4/A5/A6/A7 + 唯一未触发 A3」并单列新候选（RES-T03-tz） | §2/§3 |
| 5 | `docs/vision/roadmap.md` | `:16,223` | `admin.mfa` wrapping「不随 JWT previous 重包」 | 现行代码已惰性重包；残余改为窄口径（恢复码路径不重包等，用户 2026-09-10 接受） | G-005 |
| 6 | `docs/vision/charter.md` | `:73` | VP-016 residual 点名 mfa-wrap 旧表述 | 同步窄口径 | G-005 |
| 7 | `docs/vision/workspaces.md` | `:30,58` | 同上旧表述 + 本区 1/4 投影 | 同步窄口径；本区投影更新为 3/4（R3 完成后） | G-005 + 投影一致性 |
| 8 | `docs/vision/plans/VP-016-key-rotation-and-backup.md` | `:115` | 关门记录旧表述 | **保留历史时点原文**，加一行注记指向现行状态 | G-005；历史不回改 |
| 9 | `docs/vision/roadmap.md` | `:353` | Admin 功能最近一拍 = VP-034 | 更新为 VP-035 状态与后续 | 投影一致性 |
| 10 | `docs/vision/roadmap.md` | `:392` | 当前组合焦点 = VP-035 active | 按 VP-035 关门后的实际状态改写 | 投影一致性 |

## 6 · 本草案不做的事

- 不改 Charter 目的/成功边界/非目标/`vision_id@version`。
- 不清除任何 `trigger-gated` 行，不把「没有实现」写成缺口。
- 不把 `docs/architecture/**` 的文档卫生（G-001/G-003）混入本草案的 `docs/vision/**` 事务——两者分文件、分提交。
- 不把本草案写成已冻结权威；`/vision` editorial 与用户确认是前置。
