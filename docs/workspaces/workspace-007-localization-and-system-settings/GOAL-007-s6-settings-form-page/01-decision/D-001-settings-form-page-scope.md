---
id: D-001
doc: decision
title: S6 · 设置页表单/详情页改造范围与 recordSource 路径、审计模式
status: accepted
parent: GOAL-007-s6-settings-form-page
created: 2026-08-09
updated: 2026-08-09
version: 0.1.0
---

# D-001 · S6 设置页表单/详情页改造（2026-08-09）

## 决策

将 `settings` 页从「列表 + 弹窗编辑」重构为四类分组就地编辑的表单/详情页，用户选定方案 A。

1. **范围**：只改设置页形态 + 实现其所需的 renderer 预填能力；不改后端 API 语义（PATCH/POST 端点、字段、错误码不变）；不扩 Profile 可见性；不新增协议字段。`settings` 仍是 F-V029 分母固定面，测试与证据同步更新。
2. **预填路径 = 实现既有协议表面 `form.props.recordSource`**（ADR-0021，since 2.1）：
   - 盘点（I-001）：`component-registry.json` form 节点已定义 `recordSource`（`{method:GET, url, path?, query?, responseMapping:{field→响应路径}}`，要求 capability `form.record.load`、`responseMapping` 非空、search 模式禁止）；`constructRequest({kind:"recordSource",...})` 已实现 URL 构造与校验；渲染器 `parseRenderNode` 未放行、`FormView` 初值只认 `modalRow`。
   - 结论：**不新增字段**，把 `recordSource` 解析放行 + GET 预填 + `responseMapping` 映射 + capability 门禁 + reload 重预填 + 只读权限门禁接线进 renderer。
3. **settings schema 结构**：说明 + 4 个 recordSource 预填内联表单（每类一个，`submitAction` 复用现有 update* PATCH，节点级 `permissions.edit`+`permissionCascade.keys:["edit"]`）+ Restore defaults 改 `actionButton`（`actionId:"resetSettings"`、`permissionIntent:"edit"`、confirm）。删除表格、recordView、4 个 `open*` modal action、行内 Edit。meta 增 `form.record.load`，移除 `table.sort`/`actions.row.request`。**全文本复用现有 catalog 键**，不新增键。
4. **审计模式**：常规阶段 **self**（04）兜底；**C4 关门 = `independent`**（settings 为 F-V029 分母固定面 + 协议表面实现，风险中高；沿用 S5 惯例）。

## 为什么

- 设置页是**单例全局配置**：1×9 表格是 CRUD 列表机制错配，读/写断开、弹窗字段为空（工具条 `invokeAction(trigger, null)`）、空保存有清空风险。VP-007 §交付范围已写明「按 General / Branding / Localization / Appearance 四类组织」——四类就地编辑表单正是该产品结构的标准形态（Linear / Vercel / GitHub settings 同型）。
- 选 recordSource 而非新增字段：协议层已冻结此能力并有 conformance cases；实现它符合本仓库「实现既有冻结协议表面」的纪律，且避免自定义 `valueSource` 造成私有语义。
- 逐类独立表单（各带 Save）：映射现有 4 个 PATCH action，避免单一大表单误把未改动区段字段随空值提交。
- Root 需暂时回退关门：新增 S6 子目标使 `done` 不再成立；历史关门记录不重写，重新开根决策见 GOAL-001 `D-003`。

## 未选方案

| 方案 | 未选原因 |
|------|----------|
| 自定义 `form.props.valueSource` 新字段 | 协议层已有 `recordSource`（ADR-0021）并冻结验证；新增字段会造成私有语义与 conformance 分叉 |
| 保持表格，仅给 modal 预填（快速止血） | 形态仍是「列表 + 弹窗」，未解决读/写断开与分散编辑问题；用户选定方案 A |
| 单一大表单全字段一次保存 | 现有 PATCH 是分字段 bodyMapping；全字段提交会把未动区段字段当空值发送，覆盖风险高 |
| 用 section+text 当分组标题 | form 节点 registry 已有 `title`/`titleKey` props（同类未接线）；实现 form title 更贴合协议 |
