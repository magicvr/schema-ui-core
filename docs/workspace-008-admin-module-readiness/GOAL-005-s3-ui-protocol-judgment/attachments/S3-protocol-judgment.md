# S3 · UI 协议与共享能力判断 · 证据

> 本文件为 GOAL-005-s3-ui-protocol-judgment 的执行证据。候选基线：S2 提交 `39a0737`（clean）。协议 pin：`schema-ui-docs@v2.7.0` / commit `ca9e5fe…`。

## 1. I-READINESS-003 闭合 · fixture/conformance 与 I-PROTO-FULL-001 一致性

### 1.1 实测核对

| 主张 | I-PROTO-FULL-001 声明 | S3 实测 | 结论 |
|------|----------------------|---------|------|
| 能力域 | 12/12 | 12 域（D-NODE/D-EXPR/D-COMP/D-DATA/D-ACT/D-PERM/D-APP/D-TABLE/D-FORM/D-UPLOAD/D-VER/D-VAL） | ✅ 成立 |
| registry type | 24/24 | `docs/schemas/component-registry.json` components = 24 | ✅ 成立 |
| 行为套件 | 16/16 include | `apps/web/src/protocol/upstream/*.cases.json` = 16 套件 | ✅ 成立 |
| case 总量 | v1.0.0 claimed 320/320 全绿 | **320 total = 318 执行 + 2 排除**，执行项全部通过 | ✅ 已由 workspace-005 v1.0.1 勘误 |
| exclude | v1.0.0 claimed 0 | **2 local adapter execution exclusions**：`m1-missing-app-manifest-capability`、`m1-navigation-without-capability`（错误信封 `MISSING_REQUIRED_CAPABILITY` vs 上游 `CAPABILITY_REQUIRED`） | ✅ 已由 workspace-005 v1.0.1 勘误 |

### 1.2 F-001 调和结论

- **现行权威 disposition**：318 执行 + 2 排除，已由 S0 [D-003 §5](../GOAL-001-admin-module-readiness/01-decision/D-003-s0-denominator-freeze.md) 冻结，并经本 S3 实测复核（conformance 全绿）。
- **`I-PROTO-FULL-001` 文档漂移**已处置：workspace-005 `I-PROTO-FULL-001` v1.0.1 + D-003/E-007 正式记录 12/12 · 24/24 · 16/16 · 318+2，且明确两项 local adapter exclusion 不是域级收缩。
- **处置结果**：I-003 → **verified**；跨区 F-001 已由 workspace-008 A-003 以 `fixed` 路径闭合，不再保留 documented residual。
- I-003 → **verified**。

## 2. 共享能力映射（S0 D-003 §13 → schema-ui-docs@v2.7.0）

| 共享能力（框架级） | 协议面 | 分类 | 证据 |
|--------------------|--------|------|------|
| list/detail 读表面 | D-DATA（datasource/response mapping）、D-NODE | **covered** | data-table / search-form-table / record-view；`request-construction`/`response-mapping`/`static-data` conformance |
| 写操作 create/update/delete | D-ACT（actions） | **covered** | schema-crud、admin-list-edit-lifecycle；`actions`/`request-lifecycle` conformance |
| 状态流转 | D-ACT（request lifecycle）、D-EXPR | **covered** | `request-lifecycle` conformance、form-with-reactions |
| 权限 read/write/assign | D-PERM（permission inheritance/intent） | **covered** | permissions-inheritance conformance；D-PERM 引擎；users/roles/settings 权限键 |
| 操作审计 | Go 侧 host 能力（core.operationlog） | **covered（host）** | 协议无审计 UI 域；host 经 operationlog 记录写操作；无 protocol-gap |
| 迁移与 system-data reconcile | 基础设施（非 UI 协议） | **covered（host）** | store 迁移 0001-0010、system-data reconcile；协议不涉及 |
| 表单/校验/反馈 | D-FORM、D-VAL、D-EXPR | **covered** | 24 registry types 含表单控件；`06-validation` + schemas；错误反馈（errorcatalog） |
| 导航与 Manifest 发布 | D-APP（app manifest/navigation） | **covered** | core.manifest-route + navigation-capability；`app-manifest`/`app-navigation` conformance |
| 双语与设置 | D-APP（titleKey/labelKey 本地化）、host i18n + settings | **covered（host+协议键）** | manifest titleKey/labelKey；i18n en-US/zh-CN；admin.settings |

## 3. 前端宿主能力矩阵（冻结）

| 能力面 | 已实现 | 宿主缺口 | 非目标 | Profile | 证据路径 |
|--------|--------|----------|--------|---------|----------|
| component 渲染（24 registry types） | ✅ 全覆盖 | — | — | mvp/admin | renderer + component-format conformance |
| action 触发/生命周期 | ✅ 全覆盖 | — | — | mvp/admin | actions/request-lifecycle conformance |
| reaction 表达式 | ✅ 全覆盖 | — | — | mvp/admin | reactions conformance |
| page 渲染（schema-driven） | ✅ 12/10 页 | — | — | mvp/admin | representative-pages / s2 演练 |
| 权限显隐/禁用（D-PERM） | ✅ | — | — | mvp/admin | permissions-inheritance + D-PERM |
| manifest/导航 | ✅ | — | — | mvp/admin | app-manifest/navigation |
| 表单/校验 | ✅ | — | — | mvp/admin | form-controls + validation |
| 上传 | ✅ 端点 | **F-007 host-gap**：服务端仅认证无权限键 | — | mvp/admin | uploads conformance + upload.go |
| 可访问性（焦点管理） | 静态 aria 部分 | **F-002 host-gap（required→S4）**：模态/抽屉无焦点约束/恢复 | — | mvp/admin | S1-findings F-002 |
| 双语/设置 | ✅ | — | — | admin（设置页） | i18n + settings |
| 领域特有 UI（订单/钱包/类目/通知） | — | — | **non-goal** | — | VP-008 Non-goals |

## 4. 回流决策

- **protocol-gap**：无。框架级共性能力（S0 §13 全部 9 项）均由 `schema-ui-docs@v2.7.0` 覆盖或属 host 能力；**未发现需要上游协议变更或私有 Schema 扩展的全局缺口**。
- **host-gap 进入 S4**：F-002（a11y 焦点管理，required）、F-007（上传授权深度，S3 判为 host-gap 需 S2/S4 评估是否补权限键）。
- **不需回 `/vision`**：无协议变更需求；`I-PROTO-FULL-001` v1.0.1 勘误已由 workspace-005 D-003/E-007 与本区 A-003 完成，属于执行分母校正，不改变愿景范围。
- **全局 protocol-gap 阻断判断**：不触发（无 protocol-gap）；`go` 候选不受协议缺口阻断。

## 5. 汇总

| 分类 | 数量 | 项 |
|------|------|----|
| covered | 9 | S0 §13 全部框架共享能力 |
| host-gap | 2 | F-002（required→S4）、F-007（待评估） |
| protocol-gap | 0 | — |
| non-goal | 1 | 领域特有 UI（订单/钱包/类目/通知） |
