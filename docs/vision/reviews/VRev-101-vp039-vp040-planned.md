---
id: VRev-101-vp039-vp040-planned
doc_type: vision-review
title: VP-039 版本/维护/诊断 + VP-040 timestamptz 合同 · 计划阶段意图审视
source: self
scope: VP-039-version-maintenance-diagnostics · planned；VP-040-timestamptz-persistence-contract · planned（停放）
verdict: pass
open_required: 0
status: recorded
date: 2026-09-19
auditor: /vision
created: 2026-09-19
updated: 2026-09-19
parent: null
version: 0.1.0
---

# VRev-101 · VP-039 / VP-040 计划阶段意图审视

## 审视范围

用户 2026-09-19 P-004 裁决选项 1 之后的计划阶段审视：

- 新建 `VP-039-version-maintenance-diagnostics`（Admin 功能 · 体验增强收口）`planned`
- 新建 `VP-040-timestamptz-persistence-contract`（架构 · C1）`planned`，**停放至 VP-039 波次之后**
- Charter `@0.4.0` 对齐
- 结构选型（两 VP 串行 vs 合成一 VP vs VP-010 波次）
- 退出判据可判定性与 P-005
- 本轮只读事实核对

本 Review **不授权激活、开工作区或进入实现**。

## 本轮只读事实核对

| 主张 | 证据 | 结论 |
|------|------|------|
| VP-012 把 maintenance UI 写成后置 | [VP-012](../plans/VP-012-shared-cross-module-contracts.md) 首波冻结表：「四模式、统一写门禁、Host/status 投影」已交付；**不进本 VP** = 「运行时管理 UI（原文「UI 可后置」）」 | ✅ 属实 |
| 四模式写门禁已存在 | `apps/api/internal/handler/operational.go`：`SERVICE_MAINTENANCE` / `SERVICE_DEGRADED` / `SERVICE_READ_ONLY`；登录/恢复/邀请白名单 | ✅ 属实 |
| Host bootstrap 已投影 runtime mode | `apps/api/internal/handler/bootstrap.go`：maintenance → Host `maintenance`；degraded **与** read-only → Host `degraded`（read-only 精确区分在 status，不新增协议 mode） | ✅ 属实 |
| 前端已能把维护码分类为 feedback，但无持久横幅 | `apps/web/src/renderer/feedback-policy.ts` 识别 `SERVICE_MAINTENANCE`；`apps/web/src` 无 `runtime.mode` / 持久维护横幅实现（本轮检索 0 命中） | ✅ 属实 |
| 版本身份已在系统监控暴露 | `apps/api/internal/handler/systemmonitoring.go` status 行含 `version` / `commit` / `uptimeSeconds`；Web 测试绑定这些字段 | ✅ 属实 |
| roadmap 体验增强仅剩本项未立项 | [roadmap.md](../roadmap.md)：「剩余未立项 = 版本与维护提示」；未决项登记「版本与维护提示」= 未立项 | ✅ 属实（本事务将改为 VP-039 `planned`） |
| C1 时间列仍为 INTEGER | VP-035 R3：`timestamptz` 命中 0；`RT-T03` `registered`；例：`modules/authsession/migration/migration.go` `created_at INTEGER NOT NULL` | ✅ 属实 |
| 两项不是 VP-010 符合性缺口 | VP-012 显式后置 UI；`RT-T03` 从未进入已交付分母。与 VP-036/037/038 结构先例一致 | ✅ 属实 |

## 审视结论

**verdict: `pass`**（0 required；2 recommended）。

VP-039 落在 Charter 成功边界 #3（产品化 Admin 体验）与 #5（模块可组合）。它只补 VP-012 后置的产品面，基础设施前置已交付，不消耗 gated trigger。

VP-040 落在成功边界 #6（基础设施端口）与 Store 双方言合同（VR-027 / VP-013）。它是真实 schema 合同，不是「改一列」。用户选择**另立并停放**成立：两项失败模式不同，退出判据不可混用。

结构选择成立：

- 不是 Charter strategic
- 不是 VP-010 波次
- 不是合成混分支 VP
- 各为新 VP + 将来各自的 delivery 工作区；现阶段 0 区

## Charter 对齐

| 项 | 检查结果 |
|----|---------|
| 两 VP `vision_ref` = `schema-ui-core-admin-foundation@0.4.0` | ✅ |
| 不改变目的 / 成功边界 / 非目标 | ✅ editorial 组合投影 only |
| 不消耗 A3 / Redis / MQ / 搜索引擎 / 文件扫描 trigger | ✅ 两 VP 均写明 |
| VP-040 激活硬门禁 | ✅ 不得在 VP-039 仍为 `planned`/`active` 时激活，除非用户书面改序 |

## 结构选型

| 问题 | 判断 | 依据 |
|------|------|------|
| 改 Charter 目的/边界？ | 否 | 现有 Admin 体验 + 已登记架构 C1 |
| 同愿景新纲领波次？ | 是 | 两项各一波次 |
| VP-010 普通波次？ | 否 | 后置产品能力 + 从未交付的架构合同 |
| 合成一个 VP？ | 否（用户已否） | 分支、风险、退出判据均不同 |
| 结论 | **VP-039 下一拍 planned；VP-040 planned 停放** | 均 0 区；激活分两轮 `/vision` |

## 退出判据与 P-005

**VP-039** 七条判据可判定。无 required 信息项阻断 `planned`。激活前必须关闭 `I-039-004`（承载面）与 `I-039-005`（freshness）。默认候选 = Shell 横幅 + 复用 `admin.system-monitoring`，若改新模块进默认集才可能触及 `go`。

**VP-040** 六条判据可判定，但 `I-040-001`（SQLite 物理类型）是硬未知，**必须在 R1 冻结**，建议激活前至少有默认候选。无 required 阻断 `planned` 登记。`I-040-005` 阻断激活。

## Findings

### V-F131（recommended · 非阻断 · VP-039）

`I-039-004` 是本 VP 唯一可能触及 VP-008 `go` 失效触发项的地方。默认候选（Shell 横幅 + 复用 system-monitoring、不改默认集）不暂挂 `go`。若激活包改为新模块进 admin 默认集，须像 VP-038 `I-038-004` 一样做 P-004，并独立核对是内容扩展还是装配语义变更。

状态：`open · recommended`。不阻断 `planned`；由激活包承接。

### 响应（2026-09-19 · `/vision` · VRev-102 激活事务）

`V-F131` → **fixed**。用户「走流程激活」接受默认候选：Shell 横幅 + 复用 `admin.system-monitoring`，不新建模块、不改 Profile 默认集 → 不暂挂 `go`。证据：[VRev-102](VRev-102-vp039-activation.md)；VP-039 `I-039-004` verified。原 verdict / finding 原文不改写。

### V-F132（recommended · 非阻断 · VP-040）

SQLite 无原生 `timestamptz`。激活后若把「PG 改 `timestamptz`、SQLite 继续 INTEGER 且无合同说明」写成已交付，将违反 VP-013 合同平等。R1 必须在 TEXT RFC3339 / INTEGER epoch+约定 / 其他方案中冻结一项，并写清零值与编解码。

状态：`open · recommended`。不阻断 `planned`；由 `I-040-001` 承接。**阻断 VP-040 的 R1 冻结**。

### 响应（2026-09-20 · `/vision` · VRev-104 激活事务）

`V-F132` → **fixed（激活前置子要求）**：本轮已登记 R1 默认候选：SQLite `INTEGER` / PostgreSQL `BIGINT`，暂按 UTC Unix seconds 的合同平等语义；证据见 [VRev-104](VRev-104-vp040-timestamptz-persistence-contract-activation.md) 与 `[workspace-040] GOAL-001` D-001。`I-040-001` 仍为 `collecting`，R1 仍须最终冻结物理类型、精度、NULL/零值、编解码与毫秒字段分母；原 finding 与原 verdict 不改写。

## 激活门禁（后续 `/vision`）

**VP-039（下一拍）**

1. 用户 P-004 裁决 `I-039-004`（或书面接受默认候选）
2. Admin 类 freshness
3. 激活就绪 self Review
4. 用户确认 workspace / Root slug → `/govern` scaffold

**VP-040（停放）**

1. VP-039 `closed`，或用户书面改序
2. `I-040-001` 至少有默认候选
3. 架构类 freshness + 激活就绪 Review + slug

## 声明

- source = `self`，不冒充 independent。
- 不改变 Charter 目的/边界/非目标；不改变其它 VP 或 Goal status/progress。
- open required = 0；`V-F131` 已于激活事务 `fixed`；`V-F132` 仍 recommended（VP-040 R1）。
- 下一步：若继续推进 VP-039，使用 `/vision` 完成 `I-039-004` + freshness 后激活；VP-040 保持 planned。

## 当前响应投影（2026-09-20）

`V-F132` 已按 VRev-104 **fixed（激活前置子要求）**；VP-040 已满足激活门禁并进入 `active`。该响应不把 `I-040-001` 写成 verified：R1 最终合同仍由 `[workspace-040] GOAL-001-timestamptz-persistence-contract` 的决策与证据承接。
