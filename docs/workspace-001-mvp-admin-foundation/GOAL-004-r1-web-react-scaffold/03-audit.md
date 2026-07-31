---
id: GOAL-004-r1-web-react-scaffold
doc: audit
status: active
parent: GOAL-001-mvp-admin-foundation
created: 2026-07-31
updated: 2026-07-31
version: 0.3.0
---

# 审计 · GOAL-004

> 本文件是目标的唯一正式意见台账（P-003）。正式意见必须为可扫描的 `A-00N` 编号节。

## 信息就绪核对

| 核对项 | 状态 | 备注 |
|--------|------|------|
| 影响本 scope 的 I-00N | I-004-001/002 **verified**（D-002） | 分层策略 B 已锁 |
| 到期 required | 无 | 可冻结骨架目录并实施 |
| 资料引用 | 无 | 平行仓外部参考 |

## 意见台账索引

| A-ID | 日期 | source | scope | verdict | 开放 required |
|------|------|--------|-------|---------|---------------|
| A-001 | 2026-07-31 | independent | 目标定义 + 设计/计划（R1 React 骨架） | conditional | F-001 → 见 A-002 闭合 |
| A-002 | 2026-07-31 | self（编排响应） | 响应 A-001 · F-001 | pass | 0 open required |

---

## A-001 · R1 目标设计交叉审 · React Web 骨架（2026-07-31）

- **source**：independent
- **auditor**：Grok · `/audit`（独立交叉审计）
- **类型**：goal-definition / design-plan
- **scope**：GOAL-004 目标定义与 R1 设计合理性；对照 D-004、GOAL-002/003 与 R3 边界
- **verdict**：**conditional**
- **完整意见**：本节约

### 范围与区间

- 工作区：`workspace-001-mvp-admin-foundation`
- 设计/计划审；实施仍为 0（`02-execution.md`）
- 外壳完整度属 R3、协议范例属 R5 —— 用此对照是否越界设计

### 成果（有证据）

| 项 | 证据路径 |
|----|----------|
| 工具链与 D-004 一致（Vite/React/TS/npm + Tailwind/shadcn） | `00-meta.md`；Root D-004；本目标 D-001 |
| 明确「平行仓无 Tailwind，本目标需新增 UI 基线」 | `00-meta.md` 备注；避免 pure copy 偏差 |
| 业务 mock / 订单等排除清楚 | D-001.5；成功标准最后一条 |
| R1 主题最小占位 vs R3 外壳 — 边界基本清楚 | 成功标准 + 概述 |
| 与 003 对称的可运行骨架切分 | goal-tree R1 三子目标 |

### 对照成功标准（设计充分性）

| 标准 | 设计评价 |
|------|----------|
| package.json + lock + dev/build | 充分、可验证 |
| React 19 基线 | 与平行仓一致；已写死版本需在实施时核对生态兼容（非阻断设计） |
| Tailwind + shadcn「初始化**或等价**」 | 存在解释弹性（F-002） |
| 浅/深色最小占位 | 合理；与 R3 完整产品化可区分 |
| 无业务默认路由树 | 充分 |

### Findings

#### F-001 · required · med · I-004-002（host/protocol/renderer 是否 R1 预建）影响骨架形状却标 non-blocking

- **现象**：信息项承认「目录深度 / 骨架冻结」受 I-004-002 影响，级别仍为 non-blocking；meta 建议「预建空分层防后续大挪」，决策未锁定。
- **风险**：实施者若只建扁平 `src/app`，R3/R5 引入 protocol/renderer 时大挪；若预建过深空包，又增加 R1 噪声与「假完成」感。两种都可接受，**不可长期双悬**。
- **证据**：`00-meta.md` I-004-002；`01-decision.md` 信息表。
- **建议**：实施前用户选定并写入 D-00N：**(A)** R1 仅 `src/app`（+components/ui）；**(B)** R1 预建空 `host`/`protocol`/`renderer` 目录 + README 边界一句。选定后 I-004-002 → verified。是否升 required：建议 **required @ 骨架目录冻结前**（可与首次目录落盘同一门禁）。

#### F-002 · recommended · med · shadcn「或等价」削弱可验收性

- **现象**：成功标准允许「shadcn/ui 初始化**或等价**组件目录约定」。
- **风险**：「等价」可被解释为任意组件文件夹而无 shadcn 工具链，与 Charter/D-004「Tailwind + shadcn」方向摩擦。
- **建议**：改为：**必须**可指回 shadcn 初始化痕迹（`components.json` 或文档记载的 init 命令与 `components/ui` 约定）；禁止用无工具链的手写 div 冒充 shadcn 基线。I-004-001 预设（new-york 等）保持 non-blocking 可接受。

#### F-003 · recommended · low · 与 GOAL-002 目录所有权交界

- 同 GOAL-002 A-001 F-001 / GOAL-003 A-001 F-003：`apps/web` 首次创建权须在编排响应中统一。

#### F-004 · recommended · low · R1 UI 基线 vs R3 外壳的验收防漂移

- **现象**：R1 含主题切换占位；R3 为 Admin 外壳与导航。
- **风险**：R1 做成完整侧栏/多页壳会吞并 R3。
- **建议**：成功标准或验收附注：**R1 通过条件不含** App manifest 导航壳、多业务路由；单页/占位页 + 主题切换即可。

### 必改项汇总

1. **F-001**：在骨架目录冻结前锁定 I-004-002（扁平 vs 预建分层）并 verified。

### 与既有意见的异同

- 无历史 A-00N。跨目标一致：三分法合理；本目标特有风险是 **前端分层未决** 与 **shadcn 验收措辞**。

### 结论 + 建议给编排器/用户的下一步

**结论**：GOAL-004 将 Charter 的 Tailwind/shadcn 方向前移到 R1、同时把外壳留给 R3，**设计取向正确**；与平行仓差异处理诚实。因目录分层门禁未锁 → **conditional**。

**建议 `/govern`**：响应 F-001（选定分层策略）+ 可选收紧 F-002 措辞；与 002 交界、003 module path 一并处理后并行实施 003/004。

### 声明

本意见不修改 status/progress；响应由 `/govern` 处理。

---

## A-002 · 编排响应 · A-001 F-001 闭合（2026-07-31）

- **source**：self（编排响应）
- **auditor**：Grok · `/govern`
- **类型**：response
- **scope**：响应 A-001 required F-001（I-004-002 分层）；收紧 F-002 shadcn 措辞
- **verdict**：**pass**
- **P-004.1**：用户指令闭合后推进 R1 → 不另作自审

### 关闭证据

| Finding / I-00N | 状态 | 证据 |
|-----------------|------|------|
| F-001 I-004-002 分层 | **fixed** · 方案 **(B)** 预建空 host/protocol/renderer | `01-decision.md` D-002；`00-meta` |
| I-004-002 | **verified** | D-002 + 后续目录证据 |
| I-004-001 | **verified** | new-york 预设 |
| F-002 shadcn 措辞 | recommended → 成功标准要求 `components.json`/init 痕迹 | `00-meta` |
| F-003 目录交界 | 服从 GOAL-002 D-002 | D-002.8 |
| F-004 R1 vs R3 | 成功标准写明不含导航壳 | `00-meta` |

### 仍开放项

- 无开放 required finding。

### 结论

可实施 `apps/web` 骨架。
