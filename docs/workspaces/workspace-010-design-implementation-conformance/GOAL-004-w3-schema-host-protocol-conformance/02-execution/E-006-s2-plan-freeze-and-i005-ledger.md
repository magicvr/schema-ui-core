---
id: E-006
goal_id: GOAL-004-w3-schema-host-protocol-conformance
title: S2 方案冻结 + I-005 台账收尾 + cross 审视闭环
status: recorded
created: 2026-08-13
updated: 2026-08-13
parent: GOAL-004-w3-schema-host-protocol-conformance
version: 0.1.0
---

# E-006 · S2 方案冻结 + I-005 台账收尾 + cross 审视闭环

## 已完成事实

### 1. I-005 台账收尾（verified）

- **迁移矩阵 / 兼容规则**：上游 `docs/migrations/2.7-to-2.8.md`（tag `521cff8`）为 additive
  MINOR：双轨兼容矩阵、零强制迁移动作、capability 门控（仅声明 `host.*` capability 的 Host
  才受新增契约约束）。
- **弃用机制**：vendored `docs/schemas/capability-registry.json` 的 `deprecatedSince` /
  `removedIn` 字段由 `apps/web/src/host/claim.ts` 的 `validateRegistry` 消费（`removedIn`
  严格晚于 `deprecatedSince`、依赖闭包、DAG 无环）。
- **正反 fixtures**：version-negotiation（2.8 向量）随 stage3 套件执行；host 三 suite 96
  fixtures + app-manifest 41 零排除；app-manifest 解析对 2.7/2.8 严格协商（`PROTOCOL_VERSION_TOO_LOW`
  / `MISSING_REQUIRED_CAPABILITY`）与 return intent 敏感 key 拒绝均已有 fixtures 覆盖。
- I-005 → `verified`（00-meta / 01-decision / 03-audit 三处同步）。

### 2. S2 方案冻结（D-002）

- 附件 §1c 95/95 处置与上游 ADR-0034 D10 机械比对 0 差异；adopt-now 行 shape/state/security/
  fixtures 由 accepted ADR-0035/0036/0037 + v2.8.0 正式制品覆盖（A-005 独立复核）。
- 附件 §6 四复选框勾选；(a)(b)(c) 实质满足，(d) 由 A-005（independent）+ A-006（self）组成
  cross 审视并闭合 required findings。
- I-001 → `verified`；I-002 三处一致为 `verified`。

### 3. cross 审视与 findings 闭合

- **A-005**（independent · grok build · grok 4.6 · xhigh）：`conditional`，BLOCKING_COUNT=4，
  F-1～F-4（P1 required）+ F-5～F-8（P2）。机械比对 §1c vs D10 = 95/95 / 0 差异。
- **A-006**（self）：`pass`。F-1～F-4 整改闭合（台账 verified、补 self、D-002 入索引 +
  residual 改拟议、S4 清单补 IMP-002 与 residual 纠错）。

### 4. 台账一致性修正（A-005 F-5/F-6/F-7 响应）

- host fixtures 计数统一为 **96**（E-004、00-meta I-005、D-002）；
- 附件 §6 H0 注记 `proposed` → `accepted`；§2.2 IMP-004 → 已裁定 reserve；
- 00-meta S3 正文 tag `593f625` → `521cff8`；progress 正文 `1/6` → `2/6`（S2 勾选后 3/6）。

## 阻塞 / 风险

无阻断。A-005 的 F-5～F-8 为 P2（已 fixed / acknowledged）；O-001/O-002/O-003（vendored 行尾
与 app-navigation 2.7 pin）为 recommended 观察项，转 S6 后续。

## 关联信息项

- I-001 / I-002 / I-005 → `verified`（见 §1/§2）。
- I-003 / I-004 / I-006 → `verified`（维持）。

## 下一步（计划）

S4 工作清单（D-002 §2）：S4-1 return intent 登录流接入、S4-2 reauth-required 映射、
S4-3 hostOwnedPaths、S4-4 `$deps` residual 纠错、S4-7 IMP-002 证据固化；随后 S5 验证、S6 关门。

## progress

S2 检查点完成 → **3/6**（S1 + S2 + S3）。
