---
id: GOAL-001-wallet-prepaid-instrument
doc: audit
status: done
parent: null
created: 2026-09-02
updated: 2026-09-02
version: 0.1.5
---

# 审计 · GOAL-001-wallet-prepaid-instrument

## 信息就绪核对（按 scope）

| 核对项 | 状态 | 备注 |
|--------|------|------|
| 影响本 scope 的 I-00N | I-029-001 by-owner 门禁与冻结合同一致 | A-007：OwnerExists 只查 user 表（F-001 independently **fixed**）；主体校验在 `CreateAccount(subject)`/`Redeem` |
| 到期 required 是否已 verified / residual | **是** | A-002 三条 required 维持 closed（A-004）；A-005 F-001 → A-007 independent **fixed**；F-002～F-005 recommended → A-008 全部 `fixed`（含 0065 批次注册表 / PG e2e / expiresAt 范围 / 声明化导出） |
| 资料引用（若有）是否固定且用户确认 | 无 | `shared_materials_catalog: none` |

## 意见台账索引

| A-ID | 日期 | source | scope | verdict | 开放 required | 文件 |
|------|------|--------|-------|---------|---------------|------|
| A-001 | 2026-09-02 | independent | Root 钱包预付凭证与主体接缝全量 | **pass**（历史） | — | `03-audit/A-001-root-closeout-independent.md` |
| A-002 | 2026-09-02 | independent | Root 方案设计与代码实现（深度交叉审查） | **conditional** | 原文 3 required；闭合见 A-003/A-004 | `03-audit/A-002-root-design-and-code-independent.md` |
| A-003 | 2026-09-02 | self | A-002 独立审计意见合并响应与闭合 | **pass** | 3 required 主张 `fixed`（A-004 核验） | `03-audit/A-003-root-a002-closure-response.md` |
| A-004 | 2026-09-02 | independent | A-002 F-001/F-002/F-003 关闭证据复审 | **pass** | **0** required（2 recommended 不阻断） | `03-audit/A-004-a002-finding-closure-independent.md` |
| A-005 | 2026-09-02 | independent | Root 完成情况 · 方案设计与代码实现（不以治理文档为关门证据） | **conditional** | 原文 1 required（F-001）；闭合见 A-006 | `03-audit/A-005-root-design-code-closeout-independent.md` |
| A-006 | 2026-09-02 | self | A-005 合并响应与 F-001 闭合核验 | **pass** | 主张 F-001 `fixed`（A-007 核验） | `03-audit/A-006-a005-closure-response.md` |
| A-007 | 2026-09-02 | independent | A-005 F-001 关闭证据复审 | **pass** | **0** required（F-002～F-005 recommended 处置见 A-008） | `03-audit/A-007-a005-f001-closure-independent.md` |
| A-008 | 2026-09-02 | self | A-007 登记 + A-005 F-002～F-005 recommended 闭合核验 | **pass** | **0** required；recommended 全 `fixed` | `03-audit/A-008-a007-closure-response.md` |

## 结论状态

A-005 F-001（required）经 A-006 self → A-007 independent `pass` 合法闭合；F-002～F-005（recommended）由 **A-008** 全部按 `fixed` 处置：F-002 导出改为协议声明（downloadCsv）；F-003 PG Redeem/并发 e2e 落地并实证（docker postgres:15 PASS）；F-004 0065 `voucher_batches` 批次注册表（重复 batchId → 409）；F-005 expiresAt Unix 秒范围 fail-closed + 标签提示。**open required = 0 · recommended = 0**；Root 维持 `done`（4/4）。本索引不修改目标 `status`/`progress`。

> **E-007 修订注记（2026-09-02 · 用户反馈修复）**：A-008 F-002 的**实现载体**原为 `generateBatch.onSuccess.downloadCsv`，违反 pinned `OutcomeBehavior`（additionalProperties:false）→ 页面文档 D-VAL 失败（页面 Schema 错误）。已按 upstream pin 约束修订：声明改放生成表单节点的**业务 props**（node schema 允许），renderer/类型同步，onSuccess 还原纯 `reload`；语义不变（声明驱动 CSV 导出）。同时补齐 voucher 导航图标（card→Ticket）与 DefaultNavigationOrder 排序（menu_wallet 下一位）。详情见 `02-execution/E-007-vouchers-page-fixes.md`；11 件 pinned docs/schemas 未改动。
