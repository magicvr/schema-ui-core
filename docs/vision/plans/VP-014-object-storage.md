---
doc_type: vision-plan
id: VP-014-object-storage
title: 对象存储适配器（S3 兼容 + 本地盘内嵌）
status: active
vision_ref: schema-ui-core-admin-foundation@0.2.0
lead_workspace: workspace-014-object-storage
created: 2026-08-21
updated: 2026-08-21
version: 0.2.0
parent: null
---

# VP-014 · 对象存储适配器（S3 兼容 + 本地盘内嵌）

## 状态与门闩（2026-08-21 · active）

| 项 | 值 |
|----|-----|
| status | **`active`**（2026-08-21 用户书面确认：VRev-031 `pass` 后激活并开区） |
| **lead_workspace** | **`workspace-014-object-storage`**（Root `GOAL-001-object-storage`；唯一 delivery） |
| **Vision required** | **已满足**：VRev-031 `pass`，open required = 0；V-F061/V-F062 recommended 由激活 + Root scaffold 闭合 |
| **激活门闩（现行）** | 已激活；实现证据在 lead 区。改变 Profile / 模块矩阵 / Manifest / 共同门禁时按 VP-008 `go` 消费有效性暂挂 |
| **组合位置** | 架构分支 A2；前提 = VP-013 有界 `closed`（A1）+ roadmap **RT-S01** 已 delivered、**RT-S02** 本 VP 冻结退出分母 |
| **完整 ≠ 架构清单无限扩张** | 本 VP 只承接 A2。签名 URL / 分片 / 扫描 / CDN / 产品搬运器、A3 多实例/Redis/队列、A4 可观测、A5 密钥轮换不进退出分母 |

## 意图

在 VP-003 单主线模块化内核与 VP-013 已交付的持久化端口之上，把现行 **本地盘-only 文件落盘**收成**内核对象存储端口**，并交付 **S3 兼容**与 **本地盘**两个实现：

1. **内核端口**：put / get / delete / exists；用命名空间（或等价隔离）分开头像、品牌资源、通用上传。不是业务文件管理器、不是 CMS。Handler 与模块公共契约不得把本地路径或 `os.File` 当作存储合同。
2. **S3 兼容实现**：生产 fork 推荐与本 VP 验收权威（配置接入、`readyz` 扩依赖、三类既有落盘的读写删除证据）。
3. **本地盘保留**为 dev / mvp / 快测 / 当前 Compose 的默认路径；合同上与 S3 兼容实现 **平等**，不得残缺。
4. **现有三类落盘**：`avatars`、`brand-assets`、`uploads`（含文件库与 data-transfer 导入所用的共享上传目录）改走同一端口；默认仍由 `filepath.Dir(db.path)` 派生本地根，与现行 `config.yaml` 注释一致。

读路径可继续走既有 API（如 `GET /api/account/avatars/{id}`）。本 VP **不**把签名 URL / 直传改成产品面。

本 VP 属**架构分支**，不承载 Admin 功能页或业务域。

## 配置面

对象存储由配置选择，**不是**改 Profile、也不是改模块矩阵：

- **缺省**：继续本地盘；本地双进程与 Compose 卷默认不变。没有 MinIO / S3 仍能开发与快测。不得把「没对象存储就不能启动」做成 mvp/dev 默认。
- **生产 / 本 VP 验收**：显式 S3 兼容端点、桶、凭证（具体键名由 lead Root 方案冻结）。凭证走 YAML + env 插值、密钥 fail-closed，不把 secret 写入仓库。
- 未配置对象存储时，就绪探针保持现有语义；**仅当**显式配置了 S3 兼容后端时，`readyz` 才扩该依赖。

## 首波冻结（退出分母 = 架构 A2）

| 能力 | 本 VP 交付 | 不进本 VP |
|------|------------|-----------|
| 内核对象存储端口 | put / get / delete / exists；命名空间隔离三类落盘；公共面无本地路径 / `os.File` | 业务文件管理器、媒体库产品、跨模块上帝「万能盘」 |
| S3 兼容实现 | 端点/桶/凭证、put/get/delete 证据、配置后 `readyz` 扩依赖 | Azure Blob / GCS native 作为第三方言；强制本地/Compose 必须有 MinIO/S3 |
| 本地盘实现 | 现有 `avatars` / `brand-assets` / `uploads` 路径继续为默认；同一端口语义 | 把本地盘做成缩水/落后实现 |
| 既有读面 | 继续经 API 代理读写（与现行头像/品牌/上传路由同构） | 签名 URL、TTL、直传、CDN |
| 存量 | 新写入走端口；不提供产品级本地盘→对象存储搬运器 | 自动迁移既有文件；跨后端 in-place 搬运 |

## 非目标

- 签名 URL / TTL / 直传（RT-S03）、分片 / 断点上传（RT-S04）、恶意内容扫描执行器（RT-S05）、CDN（RT-S06）
- 产品级「本地盘 → 对象存储」搬运器（与 VP-013 D-002 同类：运维自备）
- Redis、外部队列、多实例（A3）、OpenTelemetry/指标导出（A4）、JWT 密钥轮换 / KMS / TLS 终止（A5 / RT-K\* / RT-D05）
- Azure Blob、GCS native、泛「支持所有对象存储」
- 改变 Charter 边界；Admin 功能分支（含文件扫描**策略**）或业务域页面
- 重开 VP-012 / VP-013；替代 VP-009 / VP-010

## 与相邻 VP 的边界

| VP | 关系 |
|----|------|
| **VP-003** | 遵守薄内核；对象存储是内核能力，模块只消费端口。不改模块化合同 |
| **VP-013** | 已 closed 的 Store 双方言保持；本 VP 不改数据库端口。文件根今天从 `db.path` 目录派生，本地适配器可继续用该约定，但公共契约改为对象端口 |
| **VP-008 `go`** | 若实现改变 Profile 默认集 / 模块矩阵 / Manifest 装配 / 共同门禁，按消费有效性做 freshness review。纯存储后端接入若证据显示未改上述语义，不自动暂挂 `go`。激活前仍须消费前 freshness review |
| **VP-009 / VP-010** | 上传安全 finding 与设计符合性 gap 仍归持续程序；本 VP 不扩扫描范围，也不做恶意扫描执行器 |
| **VP-012** | 已 closed 的应用契约不重开；本 VP 不新增 Job/审计模型 |
| **Admin 功能 / 业务域** | 只消费对象端口；不得在本 VP 加领域表、CMS 或 Admin 扫描策略页 |

## 方向级退出判据

1. 内核对象存储端口已落地；handler 与模块公共契约不再把本地路径 / `os.File` 当作存储合同。
2. S3 兼容实现对现有三类落盘（avatars / brand-assets / uploads）可核对 put / get / delete；显式配置时 `readyz` 扩该依赖。
3. 本地盘默认路径仍可用；两实现端口语义一致；没有对象存储仍能开发与快测。
4. 生产向验收以 S3 兼容为准（至少：配置接入、读写删除、就绪探针之一可核对）。
5. 未引入第二对象存储方言；未改 Charter；未进入 Admin 功能 / 业务域范围；签名 URL / 分片 / 扫描 / CDN / 产品搬运器均未假装交付。
6. 开放 required finding = 0（或已合法闭合）。

详细纲领阶段由激活后的 lead Root（P-001）书写。方向级建议顺序（非工作区事实）：R1 端口冻结 → R2 S3 兼容接入 → R3 三类落盘收口 → R4 公共面去本地路径 → R5 双路径证据。

## 信息需求（P-005）

允许带未知立项。下列不影响「本 VP 意图已冻结」，但必须在对应阶段前关闭或经用户接受残余。

| id | 要回答的问题 | 级别 | 影响门禁 | 最晚阶段 | 状态 |
|----|--------------|------|----------|----------|------|
| I-014-001 | S3 API 子集与驱动：MinIO / R2 / AWS 的最低公约数是什么？禁止第三对象存储方言。 | required | 方案冻结 / 实施 | R2 接入前 | open；开区后由 Root 登记并收集 |
| I-014-002 | 桶模型：单桶 + 前缀 vs 多桶；三类落盘的 key 隔离规则。 | required | 方案冻结 | R1 端口冻结 | open；开区后由 Root 登记并收集 |
| I-014-003 | 配置键名与凭证注入（YAML + env fail-closed；secret 不入库）。 | required | 方案冻结 | R2 接入前 | open；开区后由 Root 登记并收集 |
| I-014-004 | 存量本地文件如何进入对象存储？ | non-blocking | 关门叙事 | R5 | **已裁决不进退出分母**：不提供产品搬运器；既有存量 = 继续本地或运维自备拷贝。开区后点名 residual 即可 |

## 工作区绑定

| workspace_id | root_goal | role | joined | notes |
|--------------|-----------|------|--------|-------|
| workspace-014-object-storage | GOAL-001-object-storage | lead | 2026-08-21 | 2026-08-21 用户确认激活并开区；唯一 delivery；不重开 workspace-013 |

## 关门记录

（仅 `closed` / `abandoned` 时填写。）

| date | outcome | summary | evidence_links | residuals |
|------|---------|---------|----------------|-----------|
| — | — | — | — | — |

## 规划修订短史

| date | change |
|------|--------|
| 2026-08-21 | 初创 `planned`：用户确认新建本 VP 承接架构 A2；退出分母 = 内核对象存储端口 + S3 兼容实现 + 本地盘默认；签名 URL / 分片 / 扫描 / CDN / 搬运器不进分母。未激活、未开区 |
| 2026-08-21 | VRev-031 self `pass`（0 required）；用户确认激活并开区。v0.2.0 `planned → active`；lead = `workspace-014-object-storage`；Root 承接 P-001 与 I-00N（V-F061）及架构类 freshness（V-F062） |
