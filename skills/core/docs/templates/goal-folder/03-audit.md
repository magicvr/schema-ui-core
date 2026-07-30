---
id: GOAL-042-optimize-readme
doc: audit
status: active
parent: GOAL-040-docs-quality
created: 2026-03-01
updated: 2026-07-29
version: 0.3.0
---

# 审计 · GOAL-042

> **模板说明（复制后请改）**：本文件是目标的**唯一正式意见台账**（P-003）。  
> 每一条正式意见必须是可扫描的 **`A-00N` 编号节**，含 `source` / `scope` / `verdict`；禁止仅散文复盘、仅聊天或仅附件无索引。  
> 下方「示例 A-001」演示结构；真实目标应删除示例内容，从 `A-001` 起按事实追加。  
> 「可选 · 阶段性复盘散文」**不能**替代编号意见。

## 信息就绪核对（按 scope）

> 审视规划、实施或关门时，核对相关 I-00N 的级别、最晚需要阶段、延期复核、证据、决策或经用户接受的残余风险。未关闭的 required 信息项应作为 finding，不得被写成“已知”或“已完成”。

| 核对项 | 状态 | 备注 |
|--------|------|------|
| 影响本 scope 的 I-00N | 待填 | 编号与最晚阶段 |
| 到期 required 是否已 verified / residual | 待填 | 无则列为 finding |
| 资料引用（若有）是否固定且用户确认 | 待填 / 无 | 缺字段 fail closed |

---

## 意见台账索引（可选总表）

| A-ID | 日期 | source | scope | verdict | 开放 required |
|------|------|--------|-------|---------|---------------|
| A-001 | YYYY-MM-DD | self \| independent | … | pass \| conditional \| fail | N |

---

## A-001 · ＜意见标题＞（YYYY-MM-DD）

- **source**：`self` | `independent`（强制）
- **auditor**：（工具 / 模型 / 人；可知时填写）
- **类型**：`stage` | `close-out` | `response` | `ad-hoc` | …
- **scope**：审什么（目标整体 / 某阶段 / 某门禁 / 某 finding 关闭证据）
- **verdict**：`pass` | `conditional` | `fail`
- **长文**：（可选）`attachments/audit-A-001-….md` — 有长文时本节仍须保留摘要 + findings 索引

### 范围与区间

- 覆盖区间 / 不覆盖项：
- 与既有意见关系：（无冲突 / 冲突见 P-004）

### 成果（有证据）

| 主张 | 证据路径 / 决策号 / 执行条目 |
|------|------------------------------|
| … | … |

### 对照成功标准（scope 内适用时）

| 成功标准 | 状态 | 证据 |
|----------|------|------|
| … | 达成 / 部分 / 未开始 / 证据不足 | … |

### Findings

#### F-001 · ＜标题＞

| 字段 | 值 |
|------|-----|
| **级别** | `required` \| `recommended` |
| **严重度** | `low` \| `med` \| `high` |
| **影响门禁** | 方案冻结 / 实施 / 验收 / 关门 / 其他 |
| **状态** | `open` \| `fixed` \| `accepted-residual` \| `user-overruled` |
| **描述** | … |
| **证据** | 路径或引用 |
| **关闭要求** | 修正到何种可核对状态 |
| **闭合留痕** | （闭合后）决策/响应节链接；fixed 须可核对产物 |

> 编号在**本目标** `03-audit` 内按 finding 递增或按条目局部编号，响应时须能唯一引用（建议全局 `F-00N` 或 `A-00N/F-00N`）。所有正式意见必须是 A-00N 编号节，无纯散文复盘。

### 必改项汇总（required）

- [ ] F-00N · …（open / 已闭合路径）

### 建议项（recommended）

- …

### 结论与下一步

- 一句话结论：
- 建议编排器下一步：
- **声明**：独立意见默认**不**改本目标 `status` / `progress`；响应与状态变更走 `/govern` + 用户确认。

---

## A-00N · 响应 ＜被响应的 A-ID＞（YYYY-MM-DD）

> **response** 模式：编排器/自侧响应记录。`source` 一般为 `self`（响应记录）；**禁止**把编排响应标成 `independent`。

- **source**：`self`
- **类型**：`response`
- **scope**：响应 A-00N / F-00N…
- **verdict**：（对关闭证据的判定，可选）`pass` | `conditional` | `fail`

### 关闭证据表

| Finding / I-00N | 闭合路径 | 证据 |
|-----------------|----------|------|
| F-001 | fixed \| accepted-residual \| user-overruled | 路径 / D-00N |
| I-00N | verified \| accepted-residual | … |

### 仍开放项

- …

### 冲突裁决（若有）

- 指向 `01-decision` D-00N；未决不得放行对应门禁。

---

## 可选 · 阶段性复盘散文

> 可保留叙事性复盘，**不能**代替上方 `A-00N`。无编号节的散文**不**作为放行/关门依据。

### 成果 / 偏差 / 改进建议 / 非正式结论

- …
