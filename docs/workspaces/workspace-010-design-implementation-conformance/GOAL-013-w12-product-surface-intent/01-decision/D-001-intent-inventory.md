---
id: D-001-intent-inventory
doc: decision-entry
goal: GOAL-013-w12-product-surface-intent
date: 2026-08-16
status: accepted
parent: GOAL-001-design-implementation-conformance
created: 2026-08-16
updated: 2026-08-16
version: 0.1.0
---

# D-001 · 开波范围：六条产品面意图（设计未冻结）

## 裁决（范围）

用户 2026-08-16 书面点名在 workspace-010 开设意图对齐子目标，纳入下列六条（编号 T-01～T-06）。**本条只接受「这六条属于本波」**；交互细节、字段矩阵、工作区归属与 YAML 语义在 S2 冻结，不得把下列「用户方向 / 工作假设」写成已验证方案。

## T-01 · 顶栏用户区收入下拉选单

| 项 | 内容 |
|----|------|
| 用户方向 | 个人中心、设置、退出登录应收进下拉；横铺在移动端不可用。入口排序目标：个人中心 → 我的钱包 → 设置 → 退出登录 |
| as-built | `apps/web/src/app/App.tsx`：`projection.user` 在桌面为横向 `nav`（`hidden lg:flex`，**小于 lg 直接隐藏**）；移动端用户链进汉堡抽屉；退出登录是顶栏常驻描边按钮。注释仍写「W8 follow-up: 个人中心 / 设置 / 退出登录」。`admin.account` / `admin.settings` 各自向 `navigation.user` 贡献一项 |
| 工作假设 | 桌面与移动共用同一用户下拉（头像/姓名触发）；退出登录进菜单末项，不再单独横铺。移动抽屉是否仍重复用户链由 I-005 定 |
| 未冻结 | 触发器形态；键盘/焦点；是否保留抽屉内副本 |

## T-02 · 列表搜索按业务重设（渲染器能力核对）

| 项 | 内容 |
|----|------|
| 用户方向 | 用户、角色、钱包、操作日志、文件库、数据字典等列表的「语义不明文本搜索」应按业务重设；怀疑搜索项渲染器不完整 |
| as-built（schema） | 下列页的搜索表单都是单个 `id: "q"` + `labelKey: "feedback.search"`：`users`、`roles`、`activity`、`wallet`、`wallet-entries`、`file-library`、`data-dictionary`、`dictionary-entries`、`recycle-bin`、`scheduled-tasks`、`task-runs`。W11 U-04 / A-002 F-003 补的就是这层通用 `q`。`data-permission`、`notifications`、`system-monitoring` **有表无搜索表单** |
| as-built（渲染器） | 控件白名单已较完整（`form-controls.ts`：`input` / `select` / `inputNumber` / `datePicker` / `dateRangePicker` / `textarea` / `switch` / `checkbox` / `radio` / `cascader` / `checkboxGroup` / `password` / `upload`；W11 另有 `optionsSource`）。**缺口主要在各业务 schema 只用了 `input`+`q`，不是控件类型未实现。** search-mode 表单已绑定 `targetTable`。筛选 `select` 在 W11 D-003 被明确留 P2 |
| 工作假设 | 本波按页设计字段矩阵，优先用已有控件；缺后端 query 的字段先登记，不假装已支持 |
| 未冻结 | I-001 字段矩阵；是否把无搜索的三页纳入本波 |

## T-03 · 个人中心顶部选项卡

| 项 | 内容 |
|----|------|
| 用户方向 | 顶部选项卡切换，每单元单独显示；体感更好、语义更明确 |
| as-built | `account.json` 单页纵向平铺：资料表单、改密表单、`mfa-manager` custom、会话表。W11 U-11（P2，未做）已写「Tabs 分区（资料/安全/会话）」 |
| 工作假设 | 采纳 Tabs；候选分组：资料 \| 安全（密码+MFA）\| 会话。需确认协议/渲染器是否已有 tab 容器，或本波加壳层分区 |
| 未冻结 | I-002 分组；无 tab 节点时的实现路径 |

## T-04 · 「我的钱包」自服务页 + 顶栏入口

| 项 | 内容 |
|----|------|
| 用户方向 | 对标个人中心，做用户查看/管理自己钱包的页面；入口在个人中心之后 |
| as-built | `admin.wallet` 是**管理端账本列表**（`/wallet`、`/wallet-entries/{id}`），侧栏导航，不是「当前用户的钱包」。workspace-011 GOAL-019/020/021 已交付账本/开户/冻结扣款。无 `navigation.user` 的自服务页 |
| 边界风险 | VP-010 写明不承载钱包等业务模块实现。用户点名本区开目标 → 先问 I-003 |
| 建议（待用户裁） | **本区**定壳层入口顺序与页面 IA；**自服务 API/账本语义**优先回流 workspace-011（或本区只做只读包装既有 get-or-create 账户 API）。未裁决前不实施 T-04 |
| 未冻结 | 只读 vs 管理；充值/提现是否出现（建议本波不做支付通道） |

## T-05 · 回收站「删除时间」格式化

| 项 | 内容 |
|----|------|
| 用户方向 | 删除时间未格式化，不友好 |
| as-built | handler `recycleItemToMap` 输出 `deletedAt: item.DeletedAt.Unix()`（**秒级整数**）。`formatDisplayTime` **只接受 ISO-8601 字符串**，数字直接 `null`（`datetime.test.ts` 对 `12345` 断言 null）。表格因此显示原始 Unix 秒 |
| 工作假设 | 与其它列表对齐：API 改输出 ISO-8601，或 `formatDisplayTime` 识别秒/毫秒。优先 API 出 ISO，避免前端猜测量级 |
| 未冻结 | 选哪条路径（低风险，S2 可默认 API ISO） |

## T-06 · 启动启用集用 config.yaml 指定

| 项 | 内容 |
|----|------|
| 用户方向 | 启动启用哪些功能仍像命令行指定；希望用 config.yaml 方便配置 |
| as-built | **W7（GOAL-008）已支持** `app.profile` 与 `app.modules_enabled`（`config.default.yaml`）。加载序：YAML → env 覆盖（`APP_PROFILE` / `APP_MODULES_ENABLED`）。仓库**没有**检入的 `configs/config.yaml`（仅 embed 默认）。README / QUICKSTART / compose 仍以 export env 为第一教学面。`modules_enabled` 是逗号字符串，不是 YAML 列表 |
| 工作假设 | 默认按「配置面卫生」做：检入示例 `configs/config.yaml`、文档改 YAML 优先、`modules_enabled` 改为 YAML 序列。**不**静默改 mvp/admin 默认模块集 |
| 未冻结 | I-004：卫生 vs 改默认启用集 |

## 未选（本条不采用）

- 六条拆成六个子目标再开工：违反 P-001；先本波路线图，S2 后再按需要拆。
- 把 T-04 直接写成已在本区实施的功能模块：与 VP-010 边界冲突，须先裁 I-003。
- 把 T-06 写成「从零做配置体系」：W7 已关门，重复立项。

## 建议的 S3 分批（待确认）

| 批 | 项 | 理由 |
|----|----|------|
| P0 | T-05、T-01 | 缺陷小、意图清、不依赖 I-001～I-004 |
| P1 | T-03、T-02 | 体验面；T-02 依赖字段矩阵 |
| P2 | T-04、T-06 | 跨区 / 配置语义，先裁 I-003/I-004 |
