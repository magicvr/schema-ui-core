---
id: GOAL-001-modular-admin-architecture
doc: decision-entry
record_id: D-002
status: accepted
parent: null
created: 2026-08-04
updated: 2026-08-04
version: 0.1.0
---

## D-002 · 响应 A-002：根目标设计补强（映射、Profile 切分、协议继承）

### 触发

[A-002](../03-audit/A-002-root-goal-design-review.md)（independent，`verdict: conditional`）指出 F-001～F-003 为 R1 方案冻结前 required，并建议同批处理 F-004～F-006。用户通过 `/govern` 指令「响应 工作3 目标1 A-002」启动闭环。

### 决定

1. **采纳 A-002 全部 findings 的修复路径为 `fixed`**（F-001～F-006），不采用 residual 或 overruled。
2. **F-001**：成功边界拆成「阶段可验收结果」与「VP 关门必证七条」两层；增加 `R 阶段 → VP exit # → 证据类型` 映射表；明文 `progress=6/6` 不得推导 Root `done` 或 VP `closed`。
3. **F-002**：R1 仅允许 Profile **候选/依赖盘点**；`mvp`/`admin` **精确模块集合与配置覆盖顺序**以 I-004 在 **R2 方案冻结前** verified 为准；R5 做运维/配置收敛，不回写否定 R2 冻结集除非新决策。
4. **F-003 / I-007**：本 Root 实施**默认不扩大** VP-003 继承的 `I-PROTO-001 v0.1.3` 覆盖范围；扩大 domain、改 exclude 或升上游协议版本必须新决策并递增覆盖表版本。新增 required 信息项 **I-007**（最晚 R1 方案冻结前）：确认继承范围可读且与迁移模块清单一致。权威仍在 VP-003 与 Q2 覆盖表路径，本区不复制他区过程状态。
5. **F-004**：在 Root meta 预置各高影响阶段建议审计模式（R1 冻结 / R3 门闩 / R6 关门至少 `independent` 等）；编排时仍可按风险确认，不得静默降级。
6. **F-005**：R3 路线图行增加 VP R3 门闩稳定锚点，并硬约束「未满足 A+B+C+D 不得进入 R4」。
7. **F-006**：默认按 R 阶段建实施子目标；信息项仅当收集有独立范围/证据价值时升格；禁止为 I-001～I-007 机械双目标。

### 为什么

设计补强可在信息不完备时完成，且直接解除 A-002 对「根目标设计可治理性」的 required 门禁，降低 R1 方案冻结时的映射歧义、Profile 阶段误冻与协议范围静默扩大风险。同批消化 recommended 项可减少后续 P-004 成本。

### 未选方案

- **accepted-residual / user-overruled 任一 required finding**：不采用；修复成本低、可逆、且建议修复路径明确。
- **仅改文档不增 I-007**：不采用；P-005 台账需能阻断 R1/R4 范围静默漂移。
- **现在创建 R1 或信息收集子目标**：不采用；本决策只做设计补强，I-001～I-007 仍 open，是否拆子目标按 F-006 约定另议。

### 影响与后续

- 修正落盘：`00-meta.md` v0.2.0；响应审计 [A-003](../03-audit/A-003-a002-response.md)；执行事实 [E-002](../02-execution/E-002-a002-design-response.md)。
- **不**勾选 R1–R6；**不**将 I-001～I-007 标为 verified；**不**放行 R1 方案冻结。
- R1 方案冻结前仍须：I-001～I-003、I-007 verified（或合规 residual）；设计侧 F-001～F-003 已闭合。
- 建议可选：`/audit` 对 A-003 关闭证据做轻量复审（非强制）。
