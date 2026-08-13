---
id: I-011-001
goal: GOAL-002-r1-bounded-research
title: 候选池与三档分档清单（R1 调研产出 · v1.1.0 响应 A-002）
date: 2026-08-14
status: verified
parent: GOAL-002-r1-bounded-research
created: 2026-08-14
updated: 2026-08-14
version: 1.1.0
---

# I-011-001 · 候选池与三档分档清单

> R1 有界调研产出（GOAL-002-r1-bounded-research）。判据见 D-001：一等公民=业界普遍+基架未覆盖；常用=高频非普遍；增补=低频按需。**v1.1.0（2026-08-14）响应 A-002（grok build independent）**：F-05/F-06 降为常用（S-13/S-14）；S-12 修正（不复用 0006）；C-01 修正 + 账号启停并入 F-03；增补补齐 B-10/B-11；`7 协议对照口径修正。变更记录见 `8。

## 1. 已覆盖（基架已交付，不重复立项）

| # | 能力 | 基架证据 |
|---|------|----------|
| C-01 | 账号管理（users CRUD、角色分配、管理员改密；**不含产品态启停**——启停见 F-03） | admin.users 模块；代表页 users |
| C-02 | 角色与权限（roles CRUD、RBAC、权限继承、assign） | admin.roles 模块；permissions-inheritance 证据 |
| C-03 | Schema 驱动 CRUD 全链路（list/detail/create/update/delete/search-form-table/batch/upload/reactions/长内容呈现） | renderer + schema-render + 代表页 |
| C-04 | 系统设置（General/Branding/Localization/Appearance 四类） | admin.settings（VP-007） |
| C-05 | 多语种运行时（zh-CN/en-US/auto） | VP-007 |
| C-06 | 设计系统与主题（design tokens、light/dark、shadcn 风格） | VP-005 |
| C-07 | 操作日志（activity 页面 + operationlog 模块 + 设置；含 auth.login/logout/refresh 事件） | admin.activity / core.operationlog |
| C-08 | 导航 / Manifest / Profile（mvp/admin/demo + 协议聚合） | manifest-route / navigation-capability |
| C-09 | 文件上传（表单控件级 + 授权/配额/所有权加固） | form-with-upload + VP-009 |
| C-10 | 用户改密（users CRUD 内、吊销 access token） | handler/users.go（VP-009 W4） |
| C-11 | 登录安全基线（失败自动锁定 423/限流/改密吊销） | auth 模块 + VP-009 |

## 2. 候选池来源（业界样本）

| 来源 | 锚点能力 |
|------|----------|
| 通用 admin 面板能力综述（Appwrite 博客） | 认证、RBAC、审计日志、仪表盘、通知、设置、文件管理、监控 |
| 电商 admin 参考实现（React+TS e-commerce admin） | 商品管理、订单跟踪、用户管理、分析、RBAC |
| Go admin 框架（simple-admin-core） | RBAC、API 管理、字典、日志 |
| 企业后台模板惯例（vue-element-admin / Ant Design Pro 生态） | 仪表盘、用户/角色、字典、日志、个人中心、监控、定时任务 |
| 多商户市场平台 admin 面板（eshop-plus 类） | 订单、钱包/余额、类目、通知、商品、营销 |
| 用户点名（2026-08-14 会话） | 订单、钱包为典型代表；类目、通知入候选池 |
| 协议清单（protocol-inventory-v2.7.0 `2.5 信息性场景） | grid-dashboard；上游样例 user-profile-* / order-list-*（R2 方案协议对照用） |

## 3. 一等公民（R2 第一批次 · 业界普遍 + 基架未覆盖）

| # | 能力 | 判定理由 | 备注 |
|---|------|----------|------|
| F-01 | 仪表盘 / 控制台（生产 Profile home） | overview 现仅 demo（dev.examples）；mvp/admin 生产面无 home dashboard；业界几乎必备 | 以标准模块（admin.dashboard）落地并进入 mvp/admin 默认启用集（Profile **内容**扩展，非装配语义变更，R2 方案写清）；overview.json 仅 section+text，dashboard 能力需新建（F-008） |
| F-02 | 数据导入 / 导出（CSV/Excel，schema 驱动） | 通用 CRUD 后台必备；基架无导出能力 | 共享能力模块；导出权限键 + 操作审计；**协议对照须在 R2 方案独立做**（不沿用 S3 9/0 外推，见 `8 必办-1） |
| F-03 | 个人中心与账户安全 + 账号启停（自助改密、会话列表/吊销、个人资料；管理员启用/停用/手动解锁） | 改密 API 已有（管理员视角）；自助 UI + 会话管理缺失；**产品态启停缺失**（C-11 仅失败自动锁定） | 复用改密/吊销逻辑；新增会话与启停端点；A-002 F-003 处置 |
| F-04 | 通知中心（站内通知、已读/未读、通知设置） | 业界普遍（站内信/通知铃铛）；基架无 | 与业务领域「通知」合并；模块边界（系统/业务通知、与公告/模板切分）在方案冻结（A-001 F-001 / A-002 F-005） |

## 4. 常用（R3 第二批次 · 高频非普遍）

| # | 能力 | 判定理由 |
|---|------|----------|
| S-01 | 数据字典（枚举/字典管理） | 企业后台高频（中文生态尤甚） |
| S-02 | 文件 / 附件库（统一文件管理、引用、清理；复用 upload 基建） | 高频；基架仅控件级上传 |
| S-03 | 系统监控与错误日志查看（health/指标/错误日志 UI） | 运维型后台高频 |
| S-04 | 定时任务管理（cron 后台） | 平台型后台高频 |
| S-05 | 公告管理 | 门户/多用户后台高频 |
| S-06 | API 令牌管理 | 平台/服务型后台高频 |
| S-07 | 类目管理（用户候选池） | 目录型业务必备 |
| S-08 | 商品管理（业界惯例补充入池） | 电商 fork 目标普遍 |
| S-09 | 数据权限（行级/数据范围） | 企业后台高频；扩展 RBAC |
| S-10 | MFA / 2FA | 安全敏感后台高频且增长中 |
| S-11 | 登录验证码 | 防爆破补充（基架已有锁定/限流） |
| S-12 | 回收站 / 软删除管理 | 高频；**需新持久化（tombstone/回收表）+ 管理 UI**——不复用 0006 records_retire（那是演示实体退场 DROP，非产品基建；A-002 F-002 处置） |
| S-13 | 订单管理（用户点名典型） | **A-002 F-001 降档**：领域模块，非「几乎所有 Admin」；按需在 R3 或 fork 时启用；若 R2 先行须声明最小实体/桩（F-008） |
| S-14 | 钱包 / 账务（余额、流水、对账） | **A-002 F-001 降档**：同上；余额变动审计 + 迁移基建 |

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
| B-10 | 组织 / 部门 / 岗位 | 企业型后台按需（A-002 F-006 显式补入） |
| B-11 | 登录日志独立视图 | 审计诉求强化时（事件源已在 operationlog，缺独立产品面；A-002 F-006 显式补入） |

## 6. 明确不入池（非目标）

- 多租户 / 白标 SaaS 控制台：Charter 非目标。
- 运行时插件市场 / 远程模块：Charter 非目标。
- 协议语义扩展（新 capability）：走上游提案或 /vision 兼容决策。

## 7. 交付依赖核对（I-002，v1.1.0 修正）

- 一等公民 4 项（F-01～F-04）均可按标准 Admin 六项以新模块落地。
- **协议对照口径（A-002 F-004 修正）**：S3-protocol-judgment 的 9 covered/0 protocol-gap 是 VP-008 准入波次对当时共性能力的分类，**不得外推**为 F-01/F-02 的协议判定；R2 方案须对 dashboard / 导出做独立协议对照（protocol-inventory `2.5 grid-dashboard 信息性场景、上游 user-profile-*/order-list-* 样例、node.schema.json export 扩展动作键），并沿用「呈现自由 + fail-open」处置时须留痕。
- F-04 通知中心模块边界（系统/业务通知 vs S-05 公告 / B-09 模板）在方案冻结。
- F-01 的 admin.dashboard 进入 mvp/admin 默认启用集属 Profile 内容扩展：用既有模块贡献机制 + 更新 adminFunctionalOrder，不改装配语义（R2 方案写清，防误触「不改变 Profile 默认集」字面门禁）。
- 领域模块（S-13/S-14）保持领域问题留领域台账；共享基架问题回流 VP-009/VP-010。

## 8. 审计响应（A-002 grok build independent，2026-08-14 用户裁决 A）

| finding | 级别 | 闭合 | 证据 |
|---------|------|------|------|
| F-001 档位漂移（订单/钱包） | required high | **fixed（用户裁决 A：降档）** | S-13/S-14 移入 `4 常用；R2 = F-01～F-04 |
| F-002 S-12 复用 0006 | required med | **fixed** | `4 S-12 修正：需新持久化 + 管理 UI |
| F-003 C-01 过标 + 启停缺口 | required med | **fixed** | `1 C-01 修正；启停并入 F-03 |
| F-004 `7 9/0 外推 | recommended med | **fixed（口径修正）** | `7 协议对照口径改写 |
| F-005 通知边界 | recommended med | 登记 R2 方案必办 | `7 + 必办-3 |
| F-006 组织/登录日志缺位 | recommended low | **fixed（显式补入）** | `5 B-10/B-11 |
| F-007 信息项/指针未同步 | recommended low | **fixed** | Root 00-meta I-001 verified、workspace.md/附件路径同步 |
| F-008 home 装配/依赖声明 | recommended low | 登记 R2 方案必办 | `7 + 必办-4 |

### R2 方案必办清单（立项时逐项核对）

1. F-01/F-02 独立协议对照（grid-dashboard / export 扩展键 / user-profile-* / order-list-* 样例），呈现自由 + fail-open 处置留痕（A-002 F-004）。
2. F-04 通知中心模块边界冻结：系统通知 / 业务通知 / 公告（S-05）/ 消息模板（B-09）切分（A-001 F-001 / A-002 F-005）。
3. F-01 admin.dashboard 进入 mvp/admin 默认集 = Profile 内容扩展：贡献机制 + adminFunctionalOrder 更新 + 装配语义不变声明（A-002 F-008）。
4. F-05 相关：若订单先行须声明最小实体或桩；与 S-07/S-08 依赖显式（A-002 F-008）。
5. F-03 账号启停端点与权限键（users.enable/disable）、会话列表/吊销端点（A-002 F-003 落地）。
