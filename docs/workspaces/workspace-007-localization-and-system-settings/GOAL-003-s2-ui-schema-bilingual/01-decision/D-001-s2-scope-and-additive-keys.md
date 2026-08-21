---
id: D-001
doc: decision
title: S2 · 实施范围、key 解析约定与 additive registry 扩展登记
status: accepted
parent: GOAL-003-s2-ui-schema-bilingual
created: 2026-08-09
updated: 2026-08-09
version: 0.1.0
---

# D-001 · S2 实施范围与 key 解析约定

## 触发

Root D-002（I-L10N-001 前端 key 解析）已冻结；S2 按 F-V029 分母实施前冻结本阶段解析约定与 additive 扩展。

## 决定

1. **解析约定**（渲染时统一）：任何用户可见文本 prop 的解析顺序为
   `*Key 存在 → t(*Key)`（catalog 当前语种 → en-US → key 本身，缺失可观察）→ 否则字面 `label`/`title`/`text`（协议规范文本）。字面文本为 en-US 规范基线；需要双语化的文本一律加 `*Key`。
2. **Manifest 层（C2）**：前端 `navigation.ts`/`App.tsx` 解析 `labelKey`/`titleKey`；API 侧 manifest 数据（`app-manifest.json` + 4 个模块 fragment）为 pageId 与导航项补 `titleKey`/`labelKey`（协议已声明字段，additive）。catalog 键空间 `manifest.*`。
3. **Schema 层（C3）**：12 个 page schema 文档为用户可见文本补 `*Key`。**关键边界（S2 审计核对修正）**：`docs/schemas/component-registry.json` 是 **schema-ui-docs@2.7.0 上游 pin 制品**（`stage3-fixtures.test.ts` I-PROTO-004 对 sha256 校验），**禁止改写**；其已声明的 `labelKey`/`titleKey`/`contentKey`/`options.labelKey` 即上游 key 字段，直接使用。以下四个缺口字段（上游 registry 未声明）作为**本地页面文档约定**（非 registry 扩展、非上游语义主张；上游 `node.schema.json` 的 `props` 为开放对象，文档级合法）登记，Renderer 解析并遵循冻结回退链：
   - `form.props.submitLabelKey`（提交按钮）
   - `table.props.actions[].confirmKey` / `table.props.toolbar[].confirmKey`（确认文案）
   - `text.props.textKey`（文本节点）
   - 表单字段 `placeholderKey`（占位符）
4. **固定 UI（C1）**：登录页/Shell/通用反馈/通用组件/Manifest 失败面的全部用户可见文案迁入 catalog（键空间 `login.*`/`shell.*`/`feedback.*`/`manifestFailure.*`）；en-US 值 = 现状英文文本（既有测试断言不变）。
5. **M4（C4）**：缺失 key 流程用 schema 夹具制造：`labelKey` 指向不存在键 → 回退字面 label、事件可观察、主流程可完成。
6. **证据矩阵**：F-V029 表 U 行与页面行证据路径随测试落盘回填（`{SCRATCH}/unit-s2-web.log`、测试名）。

## 未选方案

- 不按 locale 重写 schema 文档（服务端 overlay，已排除）。
- 不给所有字面文本强制补 key（ID/数值等双语同形文本保留字面，减少无意义 churn）。
- 不引入 `t()` 之外的插值/复数机制（参数插值已支持）。

## 影响

- 实施范围 = 前端 7 组件 + 12 schema 数据 + 5 manifest 数据 + registry 4 处 additive + catalog 双向键集；测试全绿 + 新增双语断言。
