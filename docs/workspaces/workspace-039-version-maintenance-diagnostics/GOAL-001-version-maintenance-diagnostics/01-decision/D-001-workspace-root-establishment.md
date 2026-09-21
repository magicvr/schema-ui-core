---
id: D-001-workspace-root-establishment
doc: decision-entry
parent: GOAL-001-version-maintenance-diagnostics
status: accepted
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# D-001 · 工作区与 Root 建立：VP-039 激活落盘、`I-039-004` 裁决与 freshness 记录

## 背景

用户 2026-09-19 指令「`/vision` 走流程激活 VP-039，再交 `/govern` 开区」，承接同日已确认的 [VP-039-version-maintenance-diagnostics](../../../../vision/plans/VP-039-version-maintenance-diagnostics.md)（`planned` v0.1.0，计划 self = `VRev-101` `pass`）。

`/vision` 已完成激活事务：VP-039 `planned → active` v0.2.0，激活 self = `VRev-102` `pass`（0 required）。本决策记录 `/govern` 侧的建区落地与两项激活门禁的最终结论。

## 决策

### 1 · 用户 P-004 裁决：`I-039-004` 承载面

**裁决 = 接受默认候选**：

- **Shell 持久横幅**（maintenance / degraded / read-only）
- **复用既有 `admin.system-monitoring`** 状态行作为版本/诊断字段来源
- **不新建模块**；**不改** Profile 默认集
- `ResolveProfile` / Manifest 聚合 / 协议 pin / 共同门禁 **零改动**

**结论**：**不暂挂 VP-008 `go`**。

**约束**：若实施期改为新模块进默认集或改 Host availability 枚举 / `ResolveProfile`，须按 freshness / `go` 规则暂停与复核。

**未选**：新模块进 admin 默认集（会触发与 VP-038 `admin.jobs` 同类的内容扩展裁定，本波不需要）；把横幅只做成一次性 Toast（与「持久可感知」意图不符）。

### 2 · Admin 类 freshness（`I-039-005`）

区间 = `7e5ce891`（VP-038 激活基线 · `VRev-099` PASS）→ HEAD `6197e802`。

| 域 | 结果 |
|----|------|
| 协议 pin（`v2.9.0` / `81aa1d8`） | **PASS**（零变更） |
| 依赖锁 | **PASS**（`go.mod` / `package.json` 零 diff） |
| 迁移台账 | **PASS（已审结）**：jobs v72 等属 VP-038 |
| Profile 默认集与装配 | **PASS（已审结）**：`admin.jobs` 追加已由 VRev-099 裁定为内容扩展；本 VP 不再追加 |
| provenance | **PASS**（零变更） |

区间提交均可追溯至 VP-038 实现/关门与 workspace-010 W32–W34。`go.mod` 无 redis / kafka / amqp / rabbit / gorm / ent 新增。

**结果**：**freshness PASS，不暂挂 `go`**。

### 3 · 工作区与 Root 建立

| 项 | 值 |
|----|----|
| workspace_id | `workspace-039-version-maintenance-diagnostics` |
| canonical 范围 | `docs/workspaces/workspace-039-version-maintenance-diagnostics/` |
| vision_role | `delivery` |
| `plan_refs` / `primary_plan` | `VP-039-version-maintenance-diagnostics` |
| Root | `GOAL-001-version-maintenance-diagnostics`（`parent: null`） |
| Root 初始状态 | `active · 0/4`（纲领 R1→R4） |

slug 由用户「走流程」沿用 VP-013～038 惯例，非静默占位。

## 影响

- R1 前仍须关闭 `I-039-001`～`003`。
- 本轮**未**改动 `apps/**`。
- 不改变 Charter `primary_workspace`；不激活 VP-040。
