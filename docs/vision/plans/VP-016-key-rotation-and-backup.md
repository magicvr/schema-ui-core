---
doc_type: vision-plan
id: VP-016-key-rotation-and-backup
title: 密钥轮换与备份恢复合同（JWT + 轮换后恢复）
status: active
vision_ref: schema-ui-core-admin-foundation@0.2.0
lead_workspace: workspace-016-key-rotation-and-backup
created: 2026-08-22
updated: 2026-08-22
version: 0.2.0
parent: null
---

# VP-016 · 密钥轮换与备份恢复合同（JWT + 轮换后恢复）

## 状态与门闩（2026-08-22 · 已激活）

| 项 | 值 |
|----|-----|
| status | **`active`**（2026-08-22 用户确认：VRev-035 通过后激活并开区；V-F067/V-F068 → `fixed`） |
| **lead_workspace** | **`workspace-016-key-rotation-and-backup`**（Root `GOAL-001-key-rotation-and-backup`；唯一 delivery；**不**重开 workspace-015） |
| **Vision required** | **已满足**：VRev-035 `pass`，open required = 0；`V-F067`/`V-F068` recommended 已闭合 |
| **激活门闩** | 已满足（self Review + 架构类 freshness PASS + 用户书面「通过后开区」） |
| **组合位置** | 架构分支 A5；前提 = VP-013 A1 与 VP-014 A2、VP-015 A4 均已有界 `closed`；roadmap **RT-K01** 已 delivered，**RT-K03** 本 VP 交付；**RT-P05** 方言级 dump 已由 VP-013 交付，本 VP 只补轮换后恢复语义 |
| **完整 ≠ 架构清单无限扩张** | 本 VP 只承接 A5。A3 多实例/Redis/队列/优雅停机、KMS/HSM、静止加密、TLS、PITR、Admin 密钥页、业务域不进退出分母 |

## 意图

在 VP-003 单主线模块化内核、已交付的 YAML + env 密钥注入 fail-closed（RT-K01）与 VP-013 已交付的方言级备份（SQLite `VACUUM INTO`；PostgreSQL `pg_dump`/`pg_restore`）之上，把「密钥来自 env、无轮换合同」收成**可核对的内核轮换与恢复合同**：

1. **JWT 密钥轮换（RT-K03）**：生产可声明 **current + previous** 签名密钥；新签发只用 current；校验 current，失败再试 previous（重叠窗）。未配置 previous 时行为与今日单密钥相同。不是 KMS、不是 HSM、不是 Admin 密钥管理产品页。
2. **轮换后恢复语义（RT-P05 剩余面）**：**不重做** VP-013 的 dump 实现。本波证明：在既有 SQLite 快照与 PG 逻辑备份之上，完成密钥轮换后仍能从备份启动，且鉴权合同可核对（新密钥签发成立；重叠窗内旧 access 默认可验）。立即失效仅作 I-016-005 有界残余。
3. **内嵌默认**：本地双进程与 Compose **不**要求双密钥、不要求外部备份代理。缺 previous 密钥时进程仍能开发与快测。不得把「没 previous / 没 pg_dump 就不能启动」做成 mvp/dev 默认。
4. **生产向验收**：显式配置 current + previous 后，可核对一轮换路径与一轮换后恢复路径。

本 VP 属**架构分支**，不承载 Admin 功能页或业务域。不重开 VP-012 / VP-013。

## 配置面

密钥轮换由配置选择，**不是**改 Profile、也不是改模块矩阵：

- **缺省**：继续单一 `auth.jwt_secret` / `AUTH_JWT_SECRET`；无 previous 键；轮换不是 mvp/dev 启动硬依赖。
- **生产 / 本 VP 验收**：显式 current + previous（具体键名由 lead Root 方案冻结）。凭证走 YAML + env 插值、密钥 fail-closed，不把 secret 写入仓库。
- **生效方式**：本波默认 **进程重启后生效**。热加载 / 无重启轮换不进退出分母。
- 未配置 previous 时不得 fail-closed 挡住 mvp/dev。

## 首波冻结（退出分母 = 架构 A5）

| 能力 | 本 VP 交付 | 不进本 VP |
|------|------------|-----------|
| JWT 轮换合同 | current + previous；签发只用 current；校验允许重叠窗；缺 previous = 今日单密钥 | KMS / HSM（RT-K02）；JWT 算法升级（仍 HMAC）；OIDC/SSO 密钥 |
| 密钥集合 | 默认应用签名密钥 = `AUTH_JWT_SECRET`（及配置面等价 previous） | `DB_PASSWORD`、S3 access/secret、`ADMIN_INITIAL_PASSWORD`、对象存储凭据轮换 |
| 轮换后恢复 | 在既有 `VACUUM INTO` 与 `pg_dump`/`pg_restore` 上核对：轮换后从备份启动 + 鉴权 | 第二套 dump 实现；PITR / `pg_basebackup`；产品级备份代理；SQLite→PG 搬运器 |
| 默认路径 | 单密钥仍能开发与快测；重启生效 | 强制 Compose 常驻双密钥或备份 sidecar；无重启热加载 |
| 既有注入 | 保留 RT-K01 YAML + env fail-closed | 重开 VP-002 认证方案；改 refresh 为第二套 JWT |

未在本轮确认的 **`/readyz` 再扩外部依赖**（RT-O02 注记）**不进本波**。PG `readyz` 与显式 S3 `readyz` 已由 VP-013 / VP-014 交付。

## 非目标

- 多实例 / Redis / 外部队列 / 分布式锁 / 进程分离 / 优雅停机（A3：RT-Q\* / RT-D02 / RT-D03 / RT-D04）
- 外部 Secret Provider / KMS / HSM（RT-K02）、静止数据加密（RT-P06）、TLS 终止（RT-K04 / RT-D05）
- PITR、WAL 归档、备份即服务、把 `pg_dump` 再实现一遍
- 产品级 SQLite→PG 搬运器（VP-013 residual）、本地盘→对象存储搬运器（VP-014 residual）
- Admin 密钥管理页、监控页、全局搜索；业务域页面
- 改变 Charter 边界；重开 VP-012 / VP-013 / VP-014 / VP-015；替代 VP-009 / VP-010
- 把 VP-015 residual（otlp-sink 不解析、Store/对象/Job 指标）并进本波

## 与相邻 VP 的边界

| VP | 关系 |
|----|------|
| **VP-003** | 遵守薄内核。轮换是内核 auth 配置合同，不是模块私建平行认证 |
| **VP-002** | 已 closed 的短 JWT + opaque refresh 方案不重开；本 VP 只补签名密钥轮换，不改 Bearer / refresh 哈希模型 |
| **VP-012** | 已 closed 的 API Token / 审计 envelope 不重开；本 VP 不把服务凭证哈希密钥并进分母，除非 R1 书面证明它今天与 JWT 共用同一枚 secret |
| **VP-013** | 已 closed 的 Store 双方言与方言级 `pg_dump`/`pg_restore` 保持；本 VP 消费该合同做轮换后恢复，不改持久化端口 |
| **VP-014 / VP-015** | 已 closed；不重开对象存储或可观测；S3 凭据轮换不进本波 |
| **VP-008 `go`** | 若实现改变 Profile 默认集 / 模块矩阵 / Manifest 装配 / 共同门禁，按消费有效性做 freshness review。纯密钥配置面若证据显示未改上述语义，不自动暂挂 `go`。激活前仍须架构类 freshness |
| **VP-009 / VP-010** | 安全 finding 与设计符合性 gap 仍归持续程序；本 VP 不扩扫描范围，也不做 KMS |
| **Admin 功能 / 业务域** | 只消费轮换合同；不得在本 VP 加密钥管理页或领域表 |

## 方向级退出判据

1. JWT 轮换合同已落地：可配置 current + previous；新签发只用 current；重叠窗内 previous 可验 access。立即失效仅可作为 I-016-005 有界残余，须用户书面接受。
2. 未配置 previous 时，本地/Compose 默认仍能开发与快测；轮换不是 mvp/dev 启动硬依赖。
3. 轮换后恢复：在既有 SQLite `VACUUM INTO` **与** PG `pg_dump`/`pg_restore` 路径上，密钥轮换后从备份启动且鉴权可核对（允许分路径取证，但两方言都须有证据）。
4. 生产向验收以显式双密钥配置为准：一轮换路径 **与** 一轮换后恢复路径都须有可核对证据。
5. 未进入 A3 / KMS / PITR / Admin 功能 / 业务域；未改 Charter；未假装交付热加载或第二套 dump。
6. 开放 required finding = 0（或已合法闭合）。

详细纲领阶段由 lead Root `GOAL-001-key-rotation-and-backup`（P-001）书写：R1 轮换合同冻结 → R2 JWT 双密钥实现 → R3 轮换后恢复证据 → R4 默认单密钥仍可用 → R5 显式双密钥轮换路径 **与** 恢复路径证据。

## 信息需求（P-005）

允许带未知立项。下列不影响「本 VP 意图已冻结」，但必须在对应阶段前关闭或经用户接受残余。

| id | 要回答的问题 | 级别 | 影响门禁 | 最晚阶段 | 状态 |
|----|--------------|------|----------|----------|------|
| I-016-001 | current / previous 配置键名、生产 fail-closed 规则、secret 长度/熵是否沿用现行 `ValidateProd`。禁止把 secret 写入仓库或日志。 | required | 方案冻结 / 实施 | R1 合同冻结 | collecting |
| I-016-002 | 本波密钥集合是否仅 `AUTH_JWT_SECRET`。若发现第二枚应用签名密钥与 JWT 共用，须书面纳入或出局。DB / S3 / 种子密码默认出局。 | required | 方案冻结 | R1 合同冻结 | collecting |
| I-016-003 | 重叠窗语义：旧 access 在 previous 下可验多久；是否使用 JWT `kid`；refresh 是否受签名密钥轮换影响（opaque refresh 预期不受）。 | required | 方案冻结 / 实施 | R2 接入前 | collecting |
| I-016-004 | 轮换后恢复的最小剧本：备份点相对轮换点的先后、两方言各自证据命令、鉴权断言（login / 旧 access / 新 access）。不重做 dump 实现。 | required | 方案冻结 | R3 接入前 | collecting |
| I-016-005 | 重叠窗内旧 access 立即失效是否可被用户选为有界残余。默认建议：previous 可验，避免 15m access 全断。 | non-blocking | 退出 1 措辞 | R2 | collecting |

## 工作区绑定

| workspace_id | root_goal | role | joined | notes |
|--------------|-----------|------|--------|-------|
| workspace-016-key-rotation-and-backup | GOAL-001-key-rotation-and-backup | lead | 2026-08-22 | 2026-08-22 用户确认激活并开区；惯例 slug（D-001 留痕）；不重开 workspace-015 |

## 关门记录

（仅 `closed` / `abandoned` 时填写。）

| date | outcome | summary | evidence_links | residuals |
|------|---------|---------|----------------|-----------|
| — | — | — | — | — |

## 规划修订短史

| date | change |
|------|--------|
| 2026-08-22 | 初创 `planned`：用户确认按路线图主路径新建本 VP 承接架构 A5；退出分母 = JWT current+previous 轮换 + 既有备份上的轮换后恢复；不重做 dump；A3 / KMS / PITR / 热加载 / `/readyz` 再扩 / Admin 页 / 业务域不进分母。未激活、未开区 |
| 2026-08-22 | VRev-035 self `pass`（0 required）；用户确认激活并开区。v0.2.0 `planned → active`；lead = `workspace-016-key-rotation-and-backup`；退出 1 editorial 收口为 previous 默认可验；Root 承接 P-001 与 I-00N（V-F067）及架构类 freshness（V-F068） |
