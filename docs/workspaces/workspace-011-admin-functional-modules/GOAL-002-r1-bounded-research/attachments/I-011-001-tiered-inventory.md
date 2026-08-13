---
id: I-011-001
goal: GOAL-002-r1-bounded-research
title: 候选池与三档分档清单（R1 调研产出）
date: 2026-08-14
status: verified
parent: GOAL-002-r1-bounded-research
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# I-011-001 · 候选池与三档分档清单

> R1 有界调研产出（GOAL-002-r1-bounded-research）。判据见 D-001：一等公民=业界普遍+基架未覆盖；常用=高频非普遍；增补=低频按需。来源逐项登记（`2）。本清单是 R2/R3/R4 立项的依据，已回写 Root 纲领路线图。

## 1. 已覆盖（基架已交付，不重复立项）

| # | 能力 | 基架证据 |
|---|------|----------|
| C-01 | 账号管理（users CRUD、状态/锁定、权限键 users.read/write） | admin.users 模块；代表页 users |
| C-02 | 角色与权限（roles CRUD、RBAC、权限继承、assign） | admin.roles 模块；permissions-inheritance 证据 |
| C-03 | Schema 驱动 CRUD 全链路（list/detail/create/update/delete/search-form-table/batch/upload/reactions/长内容呈现） | renderer + schema-render + 代表页（data-table/search-form-table/form-controls/form-with-reactions/form-with-upload/admin-list-batch/data-display） |
| C-04 | 系统设置（General/Branding/Localization/Appearance 四类） | admin.settings（VP-007） |
| C-05 | 多语种运行时（zh-CN/en-US/auto） | VP-007 |
| C-06 | 设计系统与主题（design tokens、light/dark、shadcn 风格） | VP-005 |
| C-07 | 操作日志（activity 页面 + operationlog 模块 + 设置） | admin.activity / core.operationlog |
| C-08 | 导航 / Manifest / Profile（mvp/admin/demo + 协议聚合） | manifest-route / navigation-capability（VP-003/006/010） |
| C-09 | 文件上传（表单控件级 + 授权/配额/所有权加固） | form-with-upload + VP-009 W2/W4 |
| C-10 | 用户改密（users CRUD 内、吊销 access token） | handler/users.go（VP-009 W4 改密吊销） |
| C-11 | 登录安全基线（锁定 423/限流/改密吊销/账号锁定产品语义） | auth 模块 + VP-009 |

## 2. 候选池来源（业界样本）

| 来源 | 锚点能力 |
|------|----------|
| 通用 admin 面板能力综述（Appwrite 博客：build internal tools/admin panels） | 认证、RBAC、审计日志、仪表盘、通知、设置、文件管理、监控 |
| 电商 admin 参考实现（React+TS e-commerce admin） | 商品管理、订单跟踪、用户管理、分析、基于角色的认证 |
| Go admin 框架（simple-admin-core） | RBAC、API 管理、字典、日志 |
| 企业后台模板惯例（vue-element-admin / Ant Design Pro 生态；中文后台模板综述） | 仪表盘、用户/角色、字典、日志、个人中心、监控、定时任务 |
| 多商户市场平台 admin 面板（eshop-plus 类） | 订单、钱包/余额、类目、通知、商品、营销 |
| 用户点名（2026-08-14 会话） | 订单、钱包为典型代表；类目、通知入候选池 |

> 注：web_search 检索到的具体链接见附件记录（appwrite.io/blog/post/build-internal-tools-quickly；github.com/TahaNacibe/ecommerce-admin；pkg.go.dev/suyuan32/simple-admin-core；wrteam.in 多商户面板；youzan.com 电商一体化综述）。判档以「业界普遍性 × 基架缺口」为准，个别项标注不确定可回炉。

## 3. 一等公民（R2 第一批次 · 业界普遍 + 基架未覆盖）

| # | 能力 | 判定理由 | 备注 |
|---|------|----------|------|
| F-01 | 仪表盘 / 控制台（生产 Profile home） | overview 现仅存在于 demo（dev.examples）；mvp/admin 生产面无 home dashboard；业界几乎必备 | 复用 overview 页能力，以标准模块形态（admin.dashboard）落地 |
| F-02 | 数据导入 / 导出（CSV/Excel，schema 驱动） | 通用 CRUD 后台必备（用户/角色/订单等列表导出、批量导入）；基架无导出能力 | 共享能力模块，供所有资源复用；导出权限键 + 操作审计 |
| F-03 | 个人中心与账户安全（自助改密、会话列表/吊销、个人资料） | 改密 API 已有（users CRUD 内，管理员视角）；**自助个人中心 + 会话管理 UI 缺失**；业界必备 | 复用改密/吊销逻辑；新增会话列表/吊销端点 |
| F-04 | 通知中心（站内通知、已读/未读、通知设置） | 业界普遍（站内信/通知铃铛）；基架无 | 与业务领域「通知」合并为一个模块；非操作日志 |
| F-05 | 订单管理（用户点名典型） | 电商/平台 fork 目标普遍；用户点名代表 | 领域模块，含列表/详情/状态流转/审计（走标准六项） |
| F-06 | 钱包 / 账务（余额、流水、对账） | 用户点名典型；平台后台普遍 | 领域模块；余额变动审计 + 迁移基建 |

## 4. 常用（R3 第二批次 · 高频非普遍）

| # | 能力 | 判定理由 |
|---|------|----------|
| S-01 | 数据字典（枚举/字典管理） | 企业后台高频（中文生态尤甚）；非所有 admin 必备 |
| S-02 | 文件 / 附件库（统一文件管理、引用、清理；复用 upload 基建） | 高频；基架仅控件级上传 |
| S-03 | 系统监控与错误日志查看（health/指标/错误日志 UI） | 运维型后台高频 |
| S-04 | 定时任务管理（cron 后台） | 平台型后台高频 |
| S-05 | 公告管理 | 门户/多用户后台高频 |
| S-06 | API 令牌管理 | 平台/服务型后台高频 |
| S-07 | 类目管理（用户候选池） | 目录型业务必备；用户入池 |
| S-08 | 商品管理（业界惯例补充入池） | 电商 fork 目标普遍；非用户点名但同族 |
| S-09 | 数据权限（行级/数据范围） | 企业后台高频；扩展 RBAC |
| S-10 | MFA / 2FA | 安全敏感后台高频且增长中 |
| S-11 | 登录验证码 | 防爆破补充（基架已有锁定/限流） |
| S-12 | 回收站 / 软删除管理（复用 records retirement 基建） | 高频；基建已有（0006 records_retire） |

## 5. 增补（R4 backlog · 低频按需）

| # | 能力 | 触发条件建议 |
|---|------|--------------|
| B-01 | Webhook 管理 | 出现外部集成需求时 |
| B-02 | 报表 / 图表中心 | 领域数据可视化需求明确时 |
| B-03 | 优惠券 / 营销 | 电商 fork 需要时 |
| B-04 | 物流 / 履约 | 电商 fork 需要时 |
| B-05 | 订阅 / 套餐 | 多租户/计费场景（当前 charter 非目标，需先 /vision） |
| B-06 | 工单 / 客服 | 客服场景 |
| B-07 | 库存管理 | 电商 fork 需要时 |
| B-08 | 帮助 / 关于页 | 任意时间（低成本） |
| B-09 | 消息模板 / 邮件配置 | 通知中心落地后的自然延伸 |

## 6. 明确不入池（非目标）

- 多租户 / 白标 SaaS 控制台：Charter 非目标（不建设特定业务领域终端产品；fork 基架不内建多租户）。
- 运行时插件市场 / 远程模块：Charter 非目标。
- 协议语义扩展（新 capability）：走上游提案或 /vision 兼容决策，不经业务模块私增。

## 7. 交付依赖核对（I-002）

- 一等公民 6 项均可按标准 Admin 六项（http/schema/authorization/navigation/manifest/persistence）以新模块落地；F-01/F-02 需核对协议面对 dashboard/导出是否有既有定义（S3-protocol-judgment 已冻结 9 covered/0 protocol-gap——预计呈现自由或本地契约 + fail-open，与 W5 同模式）。
- F-04 通知中心涉及「通知」领域与通用能力合并，立项时需明确模块边界（通用通知 vs 业务通知）。
- 领域模块（F-05/F-06）保持领域问题留领域台账；共享基架问题回流 VP-009/VP-010。
