---
id: GOAL-001-full-protocol-contract-v2-7-0
doc: decision-entry
record_id: D-002
status: accepted
parent: null
created: 2026-08-08
updated: 2026-08-08
version: 0.1.0
---

## D-002 · 整份契约覆盖表 I-PROTO-FULL-001 冻结（S1）

### 触发

- S0 差距盘点闭合（E-002；I-001 = closed）：差集全部可纳入，**无 exclude / 范围收缩需求**（I-002 = N/A）。
- VP-006 exit 1 要求整份契约覆盖表以**新文件 + 新版本号**落盘并逐项给出 disposition；禁止就地改写 `I-PROTO-001 v0.1.3`。
- S1 为协议兼容高影响门禁 → 独立审计（用户指定 provider：**grok build**、模型 **grok 4.5**、思考强度 **high**）。

### 决定

1. 落盘覆盖表 **`attachments/I-PROTO-FULL-001-coverage-v2-7-0.md`**（`v1.0.0`），作为整份 v2.7.0 契约覆盖的**现行权威**；历史 `I-PROTO-001 v0.1.3` 保持只读未改。
2. **Disposition = 全量 include**：12/12 能力域、24/24 registry type、16/16 行为 fixture 套件（合计 320 case）全部 `include`；`include-partial` = 0；`exclude` = 0；残余清单为空。差集转 include 计数：5 域升格/新增（D-COMP/D-ACT/D-TABLE/D-FORM/D-UPLOAD）+ 6 registry type（statCard、chart、inputNumber、datePicker、dateRangePicker、upload）+ 40 fixture case（reactions 16、batchRequest 11、uploads 13）+ 2 后端端点族（批量、上传）。
3. 关闭信息项 **I-PROTO-FULL-001**（证据 = 本决策 + 覆盖表）；**I-002 = N/A**（无收缩，无用户 residual 需求；默认 include 纪律生效）。
4. 采纳实施批次 **B1–B6**（I-004）：B1 控件与展示 type；B2 $deps reaction 引擎；B3 批量；B4 上传；B5 保真/fail-closed；B6 范例+conformance 登记。
5. 维持的边界（非收缩，上游/产品既有边界）见覆盖表 §4.3：业务领域模块、完整 IAM 产品、跨页全选/部分成功、多租户市场、多版本并行矩阵、scenarios support-only、reference-* 非权威。
6. 保真度口径：契约语义可验证；**不**要求 VP-005 级视觉产品化（VP-005 实施保持冻结）。
7. 任何「已完整支持 v2.7.0」声明必须以本表 + S2–S5 实现/关门证据背书。

### 为什么

- 用户 2026-08-08 书面目标：**必须支持整份契约**；默认 include 是 VP-006 exit 1 纪律的默认路径，无需为「不做」开 residual。
- S0 证据显示缺口全部可验证实现（无技术不可行项），全量 include 是最诚实的范围定义。
- 新文件 + 新决策满足 VRev-012 F-V022（覆盖表权威落点），不破坏历史冻结证据。

### 未选方案

| 方案 | 未选原因 |
|------|----------|
| 就地升版改写 I-PROTO-001 v0.1.3 | 破坏历史 MVP 冻结证据；F-V022 禁止 |
| 部分域 include-partial 表达「暂缓」 | VP-006 exit 1 禁止用 partial 伪装整份契约；无用户 residual |
| 整域 exclude + residual | 差集全部可纳入，无必要；增加无依据的范围收缩 |
