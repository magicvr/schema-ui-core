---
id: GOAL-014-w13-settings-tabs-and-topbar
title: W13 · 设置页 Tabs 化与顶栏/搜索交互打磨（设置页功能单元 Tabs / 移动端品牌条 / 搜索框组贴合 / 明暗-语种按键对调）
status: done
parent: GOAL-001-design-implementation-conformance
created: 2026-08-16
updated: 2026-08-16
version: 0.2.0
progress: 4/4
---

# GOAL-014 · W13 · 设置页 Tabs 化与顶栏/搜索交互打磨

VP-010 / workspace-010 的**第十三波**（用户 2026-08-16 点名立项）：四项产品面交互整改——① 设置页改为按功能单元 Tabs 切换（对齐个人中心页）；② 移动端布局下 logo 与网站标题单独占顶部分一条（避免挤占功能区）；③ 列表页筛选项的【文本框+搜索按键】贴合成一个语义组件且任何页面宽度下恒在同一行；④ 顶栏功能区亮/暗切换与语种切换按键对调（亮/暗在左、语种在右）。设计冻结见 01-decision/D-001-w13-freeze.md。

## 当前边界

- **范围（本波实施）**：T-01 设置页功能单元 Tabs；T-02 移动端品牌条；T-03 搜索框组【文本框+搜索键】贴合恒同行；T-04 顶栏亮暗/语种按键对调。
- **非范围**：不改协议 schema / 页面文档结构约束；不改 Go API；不改 Profile 默认集 / 模块矩阵 / Manifest 装配语义（go 判定：无影响、不暂挂）；不改设置项字段与动作契约；不改个人中心页既有 Tabs 行为。

## 成功标准与路线图（P-001）

- [x] **S1 · 设计冻结**：四项设计决策落盘（D-001；信息项登记；用户指令即范围确认）
- [x] **S2 · 实施**：T-01 设置页 schema + i18n；T-02 App 壳层品牌条；T-03 表单控件搜索框组贴合；T-04 顶栏按键顺序（E-002）
- [x] **S3 · 测试与回归**：相关单测/e2e 更新 + vitest 全量 + tsc（E-003）
- [x] **S4 · 自审与关门**：A-001 self 审计 + 台账同步 + goal-tree 同步（E-004）

progress: 由四个等权检查点派生（S1～S4）；当前 **4/4**（2026-08-16 关门；A-001 self pass；回归 vitest 1029/1029 + tsc 0 + Go 全量 0 FAIL + e2e admin/mvp 8/8）。

## 审计策略

| 阶段 / 项 | 默认模式 | 说明 |
|-----------|----------|------|
| S1 冻结 | none | 只读落盘 + 用户指令已明确 |
| T-01～T-04 | self | 呈现/交互层整改，可逆；无 security/data/migration/production/release/compatibility 门禁语义变化 |
| S4 关门 | self | 常规关门自审 |

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | non-blocking | 设置页功能单元切分与重置按钮归属 | S2 T-01 | S2 | 用户指令 + as-built 对照 | **verified** | — | D-001 §T-01：沿用现有五个表单单元；重置按钮留在 Tabs 外（任何 tab 下可达） |
| I-002 | non-blocking | 移动端断点（品牌条生效范围） | S2 T-02 | S2 | 既有壳层断点约定 | **verified** | — | D-001 §T-02：<lg（与抽屉导航/汉堡同断点） |

无开放 required 信息项。

## 父目标

- [GOAL-001-design-implementation-conformance](../GOAL-001-design-implementation-conformance/00-meta.md)

## 台账布局

- `01-decision/`：D-NNN；`02-execution/`：E-NNN；`03-audit/`：A-NNN。
- 跨区引用用 Q2 路径（workspace-protocol §2.6）。
