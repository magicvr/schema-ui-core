---
doc_type: vision-review
id: VRev-001
status: active
source: self
created: 2026-07-31
updated: 2026-07-31
version: 0.1.0
parent: null
---

# VRev-001 · Charter 初建审视

- source: `self`
- date: `2026-07-31`
- scope: `schema-ui-core-admin-foundation@0.1.0` 与 `VP-001-mvp-admin-foundation`
- verdict: `conditional`
- suggested_class: `no-change`

### 结论

Charter、首个 VP 和组合编排与用户确认的方向一致；`schema-ui-docs` 的 `v2.7.0` 标签、固定提交和 protocol manifest 已完成外部核验。当前仓库尚无 React/Go 实现，也尚未将完整协议清单固定到本地，因此不得宣称 MVP 已具备协议覆盖能力。

### Findings

#### F-V001 · 固定协议的完整实施清单尚未落盘

- level: `required`
- status: `fixed`
- closed_at: `2026-07-31`
- closed_by: `/vision` · V6 响应（用户指令：提取清单 + 映射；闭合前不冻结覆盖）
- impact: VP-001 的协议覆盖范围与实施计划冻结；任何“支持全部协议功能”的实现主张。
- finding: 当前只确认了外部 tag、commit 和 manifest。实施前必须从 `ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b` 提取语义规范、`docs/schemas/*.json` 与 conformance 范围，并映射到 React Renderer、Go 数据/动作接口、范例页面与验证路径。
- closure: 在未来工作区的受控决策或信息登记中固化该清单及证据链接，并经 `/govern` 在受影响的实施门禁前核验。
- resolution: |
  已落盘 [protocol-inventory-v2.7.0.md](../protocol-inventory-v2.7.0.md)：从 pinned commit 提取 semanticSpecs / structuralContracts（6 schemas）/ behavioralContracts（17 fixture suites）/ ADR 索引 / 信息性场景，并映射 React·Go·范例·验证与 `mvp_candidate` 提示。
  **明确未做**：不冻结 VP-001 协议覆盖子集；不宣称已实现协议兼容。
  实施收集与门禁核验仍交 **`/govern`**（建议 I-PROTO-001…004 有界信息项）。
- evidence_links:
  - `docs/vision/protocol-inventory-v2.7.0.md`
  - external manifest @ `ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b`

#### F-V002 · 分发契约的 canonical 目录缺失

- level: `required`
- status: `fixed`
- closed_at: `2026-07-31`
- closed_by: `/vision` · V6 响应（路径：恢复 `docs/contracts/`）
- impact: 将当前仓库表述为完整治理安装，以及任何非引导性的工作区推进。
- finding: 本仓库分发 Skills 消费适配器，但 [alignment.md](../alignment.md) 的 Minimal Complete Install 要求此情形存在 `docs/contracts/`；当前目录不存在，只有 `skills/contracts/`。
- closure: 从可核验上游恢复或生成 canonical `docs/contracts/`，或经明确决策调整分发模型及相关规则，使规则、文档入口和实际路径一致。
- resolution: |
  从现有 stage 镜像 `skills/contracts/` **整树复制**恢复 canonical `docs/contracts/`（14 文件，SHA-256 逐字节一致）。
  分发模型未改：`docs/contracts/` = canonical；`skills/contracts/` = 镜像。
  闭合后允许按 checklist 将消费契约标为 present；**完整安装**仍以全部 MUST 与冷启动顺序为准。
  用户确认：在本 finding 合法闭合**前**不得非引导开区、不得宣称完整安装（本响应完成闭合后，开区仍须 `/govern` 且挂 VP）。
- evidence_links:
  - `docs/contracts/`（与 `skills/contracts/` 逐字节一致）
  - `docs/vision/consumer-checklist.md`

#### F-V003 · 双线分支的维护契约尚未定义

- level: `recommended`
- status: `open`
- impact: 后续建立双线分支时的 fork 预期与兼容性沟通。
- finding: 用户已确认维护 MVP 基架线与完整 Admin 实现线，但尚未定义命名、版本/协议兼容策略、回合并方向和变更发布方式。
- closure: 在对应后续 VP 建立前记录分支与兼容策略。

### 响应

| date | actor | summary |
|------|-------|---------|
| 2026-07-31 | `/vision` | 响应 VRev-001 + VRev-002：`F-V002` → `fixed`（恢复 docs/contracts）；`F-V001` → `fixed`（协议清单+映射落盘，**不**冻结覆盖）；`F-V004`/`F-V005` → `fixed`（recommended）；`F-V003` 仍 open。 |

门禁更新：required 已合法闭合 → 不再以 F-V001/F-V002 阻断完整安装声明中的契约行，或阻断「可进入开区引导」。开区执行仍归 **`/govern`**（slug 用户确认、挂 `primary_plan`）。**仍禁止**在覆盖子集未冻结前主张“支持全部协议功能”。

---

> **迁移说明（2026-08-07）**：本报告自 legacy inline `docs/vision/reviews.md` 原样拆出，编号与历史结论未改；相对链接已按 `reviews/` 目录深度调整。
