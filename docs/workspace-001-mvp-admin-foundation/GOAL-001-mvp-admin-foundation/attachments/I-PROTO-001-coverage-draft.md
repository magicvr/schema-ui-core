---
title: I-PROTO-001 · MVP 协议覆盖纳入/排除表（草案）
status: draft
doc_type: info-collect-draft
created: 2026-07-31
updated: 2026-07-31
parent: GOAL-001-mvp-admin-foundation
version: 0.1.0
related_info: I-PROTO-001
related_decision: D-005
source_inventory: docs/vision/protocol-inventory-v2.7.0.md
inventory_pin: ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b
freeze_status: not-frozen
---

# I-PROTO-001 · MVP 协议覆盖纳入/排除表（草案）

> **性质**：R2 信息收集草案（`I-PROTO-001` → `collecting`）。  
> **不是**：已冻结的 MVP 覆盖子集；**不得**据此主张“支持全部协议功能”或放行 R4/R5 实施范围。  
> **权威输入**：[protocol-inventory-v2.7.0.md](../../../vision/protocol-inventory-v2.7.0.md) §3（`mvp_candidate` 仅提示）。  
> **闭合条件**：用户书面确认本表（或修订版）→ Root `01-decision` 冻结决策 → `I-PROTO-001` → `verified`。

## 0. 冻结原则（草案约定）

1. **Charter / VP 边界优先**：核心账号与权限必纳；每一纳入域须可落前后端路径 + 范例 + 验证（R5 再登记 `I-PROTO-003`）。
2. **清单全量 ≠ 覆盖全量**：inventory 是能力全集；本表只选 VP-001 MVP 子集。
3. **`mvp_candidate` 可推翻**：yes/partial/optional 只作初判；本表为决策输入。
4. **业务非目标**：订单/钱包/类目/通知等内容域不纳入；上游样例仅可作结构模板并改写。
5. **partial 必须写清边界**：写入「纳入子面 / 明确排除子面」，禁止裸 `partial` 充当冻结。
6. **结构契约随纳入域带入**：被纳入域所依赖的 `docs/schemas/*.json` 自动进入验证基线，不另开“只纳 schema、不纳语义”的空洞行。

## 1. 能力域（domain_id）草案

| domain_id | 草案 disposition | 纳入边界（草案） | 明确排除 / 延后（草案） | 主要验证入口 | 备注 |
|-----------|------------------|------------------|-------------------------|--------------|------|
| D-NODE | **include** | Node/Page 树解析与合法 page 产出/消费 | 无（基座） | structural：`node`/`page` schema | 全站基座 |
| D-EXPR | **include** | 表单/列表常用 reaction 与 visible-when 级表达式 | 冷门比较语义边角、与未纳控件绑定的复杂表达式 | fixtures：`reactions` | 支撑权限显隐与表单联动 |
| D-COMP | **include-partial** | MVP 组件子集：布局/基础表单控件/表格/按钮/反馈等 Admin 最小集 | 完整 registry；扩展/进阶控件全集（见 D-FORM） | fixtures：`component-format`（子集断言）+ registry 声明 | 须在冻结时附「组件 type 白名单」或链到后续决策 |
| D-DATA | **include** | 列表/详情 datasource、response mapping、query 序列化、static-data | 与未纳域强绑定的专用 API 形状 | fixtures：`request-construction`、`response-mapping`、`query-serialization`、`static-data` | Go 列表/详情 API 主责 |
| D-ACT | **include** | page/row 级 action 触发与 request lifecycle；编辑/删除等通用后端 request | 复杂 batch 多选语义可与 D-TABLE 对齐后裁剪；业务订单动作 | fixtures：`actions`、`request-lifecycle` | 与 D-PERM intent 协同 |
| D-PERM | **include** | 账号会话最小闭环 + 权限继承 / intent → UI 显隐禁用 + Go 鉴权模型 | 完整 IAM 产品（SSO 联邦、细粒度审计后台等）除非 R4 决策扩写 | fixtures：`permissions-inheritance`；场景 `permission-inheritance` | **核心必纳**（Charter） |
| D-APP | **include** | App manifest 装载、导航壳、路由入口 | 多租户应用市场、动态远程 manifest 编排 | fixtures：`app-manifest`、`app-navigation` | 主要在 R3 落地外壳 |
| D-TABLE | **include-partial** | 排序声明、搜索表、基础列表交互 | **延后**：多选批量（ADR 0022）可作为 R5+ 或完整实现线；不阻塞 R2 基线 | fixtures：`table-sort`、`search-table` | 列表 Admin 刚需；batch 另标 |
| D-FORM | **include-partial** | 基础表单控件 + 校验展示 + 编辑回填（ADR 0021 所需最小面） | **排除 MVP**：2.6/2.7 扩展进阶控件（cascader / checkbox-group / rich-text / password 专用面 / form-default-value 全套等，ADR 0028–0033） | structural + 基础 scenarios；不强制 `form-controls-advanced` 全过 | 与 D-COMP 白名单对齐 |
| D-UPLOAD | **exclude** | — | 上传 UI/端点/fixtures `uploads` 整域 | — | inventory `optional`；完整实现线候选 |
| D-VER | **include** | `supportedCapabilities` / 版本协商 / runtime defaults（baseURL 等） | 多协议大版本并行矩阵 | fixtures：`version-negotiation`、`runtime-defaults` | 兼容边界声明必需 |
| D-VAL | **include** | 构建时/加载时结构校验（schemas）；可选服务端校验不强制 R2 | 自研替代上游 schema 语义 | `docs/06-validation.md` + 6 schemas | 支撑 I-PROTO-004 工程策略 |

### 1.1 汇总计数（草案）

| disposition | domain_id |
|-------------|-----------|
| include | D-NODE, D-EXPR, D-DATA, D-ACT, D-PERM, D-APP, D-VER, D-VAL（8） |
| include-partial | D-COMP, D-TABLE, D-FORM（3） |
| exclude | D-UPLOAD（1） |

## 2. Fixture suite 草案映射

| suite id | 草案 disposition | 绑定 domain | 说明 |
|----------|------------------|-------------|------|
| actions | include | D-ACT | |
| app-manifest | include | D-APP | |
| app-navigation | include | D-APP | |
| component-format | include-partial | D-COMP | 仅对抗白名单 type / 声明子集 |
| permissions-inheritance | include | D-PERM | 核心 |
| query-serialization | include | D-DATA | |
| reactions | include | D-EXPR | |
| request-construction | include | D-DATA | |
| request-lifecycle | include | D-ACT | |
| response-mapping | include | D-DATA | |
| runtime-defaults | include | D-VER | |
| search-table | include | D-TABLE | |
| static-data | include | D-DATA | |
| table-sort | include | D-TABLE | |
| version-negotiation | include | D-VER | |
| uploads | **exclude** | D-UPLOAD | 与域排除一致 |
| scenarios | **support-only** | 多域 | 信息性场景；**不是**单独 conformance 门禁；可作范例候选源 |

> 上游 `conformance/reference-*` / `runner` 仍为 manifest **excluded** 参考实现，不得单独证明兼容。

## 3. 范例场景候选（草案 · 非 I-PROTO-003 闭合）

R5 前仅作规划提示；路径以 inventory §2.5 为准：

| 优先级 | 场景 / 样例 | 服务域 |
|--------|-------------|--------|
| P0 | permission-inheritance | D-PERM |
| P0 | Admin 外壳 + 导航（自建，对齐 app-manifest） | D-APP |
| P1 | search-form-table / data-table | D-DATA, D-TABLE |
| P1 | admin-list-edit-lifecycle 或 list-detail（去业务化改写） | D-ACT, D-DATA, D-FORM |
| P2 | form-with-reactions（基础） | D-EXPR, D-FORM |
| 不做（MVP） | form-with-upload；form-controls-advanced/extended 全量；order-* 业务样例原文 | D-UPLOAD / 非目标 |

## 4. 与后续信息项的接口

| 信息项 | 本草案之后 |
|--------|------------|
| I-PROTO-001 | 用户确认/修订本表 → 冻结决策 → `verified` |
| I-PROTO-002 | 仅对 **include 的 D-PERM** 设计最小 API 与映射（R4 前） |
| I-PROTO-003 | 仅为 **include / include-partial** 域登记范例路径与验证入口（R5 前） |
| I-PROTO-004 | vendor vs pin 不阻断本表；建议 R3 前定，以支持 schema/fixture 校验 |

## 5. 待用户确认的开放点（不影响「草案落盘」，影响「冻结」）

| # | 问题 | 建议默认（可改） |
|---|------|------------------|
| Q1 | D-TABLE 多选批量是否进入 MVP？ | **否**（exclude-from-MVP / 完整实现线） |
| Q2 | D-FORM 是否纳入任一 2.6/2.7 扩展控件？ | **否**（基础控件 only） |
| Q3 | D-UPLOAD 是否升格 include？ | **否**（保持 exclude） |
| Q4 | D-COMP 白名单粒度：冻结时写 type 列表，还是 R3 再附？ | 冻结时写**原则 + 初表**，R3 可增补不得偷偷扩域 |
| Q5 | `scenarios` suite 是否任何一条进自动化门禁？ | **否**；仅范例/手工路径 |

## 6. 变更规则

- 本文件 `status: draft` 时修改 = 信息收集，不自动 `verified`。  
- 冻结后变更覆盖子集 = 新决策（D-00N）+ 修订本表 version + 重估 I-PROTO-002/003。  
- 不得在 chat-only 确认下把本草案当放行证据。
