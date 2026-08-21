---
id: D-003-r6-c63-contribution-validation-lifecycle
doc: decision-entry
goal: GOAL-013-r6-old-path-removal
source: orchestrator
date: 2026-08-06
status: accepted
---

# D-003 · R6 C6.3 Contribution、校验与生命周期终态契约

## 决策

C6.3 按四个可独立验证的切片实施；全部完成并通过 cross 审计后，R6-I003 才可从
`collecting` 改为 `verified`。

### 1. Schema document bytes

- `PageContribution` 增加必填的原始 JSON document bytes；`ContributionSet.Pages` 同时
  承载页面元数据、owner 与发布字节，不建立第二份 document map 真相源。
- Registrar 在注册期验证 document 是 JSON object、可确定性重编码，且
  `meta.pageId == PageID`；写入临时 set 前复制 bytes，调用方后续修改不得影响聚合结果。
- `core.schema-render` 建立一方 provider，贡献 `overview`、`data-table`、
  `search-form-table`、`form-controls`、`form-with-reactions` 五个现有 core 页面；四个
  Admin provider 各自提交模块内已 embed 的 document bytes。
- composition 只把 finalized `set.Pages` 交给 Schema handler。删除中心
  `staticSchemaDocuments`、`schemaOwnerMap` 与 nil/test fallback；禁用模块没有 Page
  contribution，也就没有可发布 document。

### 2. Configuration runtime contribution

- 新增框架无关 `ConfigurationContribution`、`Registrar.Configuration` 与
  `ContributionSet.Configurations`；规范 identity 为 `Key == Namespace`，descriptor
  继续通过 `ConfigNamespaces` 做注册前声明。
- namespace 使用稳定的 ASCII 小写点分 identifier；defaults 必须是 JSON object、可
  确定性重编码；validator 必填，且 defaults 必须在注册期通过 validator。defaults bytes
  写入 set 前复制；namespace 全局冲突 fail closed，finalize 后按 namespace 排序。
- `admin.settings` 首个声明并注册 `settings.branding`：defaults 为现有
  `siteTitle`/`logoUrl` 默认值，validator 固定非空标题与空值、同源路径或 HTTP(S) logo
  语义。Settings PATCH 的 config-change header 从该 contribution 的 namespace 传入，
  handler 不再持有 `settings.branding` 私有常量。禁用 `admin.settings` 时聚合集中没有该
  namespace。
- 本切片只冻结和实现现有运行时配置贡献，不新增通用配置 CRUD、秘密配置发布或新的
  前端协议面。

### 3. PolicyID / Visibility validation

- kernel 只负责通用语法：单个 policy reference 必须符合
  `<segment>(.<segment>)*`；segment 为小写 ASCII 字母开头，后接小写字母/数字，允许
  由单个 `-` 分隔的非空子段。现有 `system.admin`、`system.admin-editor`、
  `system.admin-editor-viewer` 均合法。
- C6.3 的 `Visibility` v1 明确为**单个 policy reference**，不是尚未版本化的布尔表达式；
  与 `PolicyID` 使用同一语法校验。若未来引入表达式，须另做版本化 grammar/兼容决策。
- owner-specific 语义仍由 `core.auth-session` 的 `rolesForPolicy` allowlist 校验。格式非法
  在 Registrar 阶段返回 `MODULE_INVALID`；格式合法但当前 auth-session 不认识的 policy
  在 system-data reconcile / readiness 前 fail closed。kernel 不反向导入认证模块。

### 4. Lifecycle matrix

- `kernel.Runtime` 拥有模块 hook 清理语义：Start 失败只逆序清理此前成功 Start 的模块；
  Ready 失败逆序清理全部已 Start 模块；失败模块的原始 stable code/module 保留，cleanup
  error 追加到 detail。`Stop` 继续清理全部模块、返回首个错误并清空 started，重复 Stop
  为 no-op。
- composition 拥有进程资源：listener/store 在 Start 或 Ready 失败后关闭；只有全部
  Start+Ready 成功后才设置 readiness gate。kernel 已清理 Ready 失败的模块 hooks 后，
  composition 不再重复 Stop。
- 自动化矩阵以真实 `mvp`/`admin` resolved Plan 参数化，覆盖 Start+Ready+Stop success、
  Start failure reverse cleanup、Ready failure reverse cleanup、Stop error continuation；另
  运行两种 Profile 的 Fx app 成功启动/停止与端口占用失败，证明组合根资源路径。

## 理由

- Schema bytes 与 owner 放在同一个 finalized contribution，才能关闭 Root A-010 F-003b，
  并消除当前“owner 动态、bytes 静态”的双真相。
- Configuration 原先只有 descriptor metadata，无法证明 defaults/validator 或运行时
  namespace 归属；结构化 Registrar 与现有五类 surface 的双检模式一致。
- kernel 只验证可移植语法、auth-session 验证已知语义，可保持薄内核并在启动前拒绝
  未知策略。
- Ready cleanup 下沉到 Runtime，使框架无关生命周期契约自身成立；composition 只处理
  listener/store，职责与失败证据可分别测试。

## 未选方案

- **handler 继续静态 merge bytes、ContributionSet 只给 owner**：保留双轨，不能关闭
  F-003b。
- **Configuration 仅加常量或只填 `ConfigNamespaces`**：没有 defaults/validator/runtime
  聚合，仍不满足架构 §2.2。
- **kernel 直接导入 auth-session policy allowlist**：形成反向业务依赖，破坏薄内核。
- **在 C6.3 发明完整 Visibility 布尔表达式语言**：超出现有数据和冻结范围，增加未验证
  兼容面；本轮固定单 policy reference。
- **只靠 composition 在 Ready 失败后调用 Stop**：`Runtime` 独立使用时契约不完整，且
  难以对模块清理语义做稳定单测。

## 影响与后续

- R6-I003 已有可实施契约，但在四切片代码、测试与 cross 审计完成前保持 `collecting`。
- 依次提交 Schema bytes、Configuration+Policy、Lifecycle matrix；每个切片运行 API
  `go test ./...` 与 `go vet ./...`，重要节点单独 Git checkpoint。
- C6.3 不修改 `progress: 2/4`；只有 self + Grok independent 无开放 required finding
  并由 `/govern` 响应后，才勾选 C6.3。
