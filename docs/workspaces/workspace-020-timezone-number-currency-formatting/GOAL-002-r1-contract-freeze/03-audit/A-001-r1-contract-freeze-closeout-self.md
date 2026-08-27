---
id: GOAL-002-r1-contract-freeze
doc: audit-entry
record_id: A-001
status: recorded
parent: GOAL-001-timezone-number-currency-formatting
created: 2026-08-26
updated: 2026-08-26
version: 0.1.0
---

## A-001 · R1 合同冻结关门审计（2026-08-26）

- **source**：self
- **auditor**：govern 编排器（本会话）
- **类型 / scope**：close-out · GOAL-002（R1 合同冻结）全量：`D-001` 合同正文 vs 用户裁决（Root `D-002`）vs 代码基现状 vs 五件套 / goal-tree / workspace 一致性
- **verdict**：**pass**

### 范围与区间

2026-08-26 立项至本审计时点。覆盖：① 合同正文 §0～§6 完整性；② 与用户裁决（I-001/I-002/I-005 accepted）一致性；③ 对 R2/R3 的可消费性（无开放 required 信息项）；④ 越界守卫（DDL / 迁移 / Profile 默认集 / `docs/contracts/`）；⑤ 台账与树一致性。共享资料引用：本区 `shared_materials_catalog: none`，无引用项需核对。

### 成果（有证据）

1. 合同正文落盘：`01-decision/D-001-r1-contract-freeze.md`（§0 效力 / §1 范围与落点 / §2 时区来源 L1～L4 / §3 数字货币落点 + API 机器合同 / §4 设置归属与字段 / §5 越界 / §6 消费指引）。
2. 用户裁决一致性：三个裁决项（会话级 auto + 站点兜底 + localStorage 覆盖；前端落点；Localization tab + 头部 locale 通道）逐条入文（§2 / §3 / §4）；裁决证据 = Root `01-decision/D-002`。
3. 代码基现状核对（C2）：`siteTimezone`/`defaultLocale` 已有（`apps/api/internal/handler/settings.go`）；金额 int64 JSON（`wallet.go`）；时间 RFC3339 毫秒 UTC（`rfc3339.go`）；前端 `Intl.*` + `schema-ui:locale` 单通道（`apps/web/src/i18n/`）；Localization tab 快测断言（`startup-config.test.tsx`）。
4. 台账一致性：五件套 + 三个 ledger 目录齐全；goal-tree 树/表含 GOAL-002（active 2/3）；workspace.md 纲领阶段 R1 → 已立项；Root meta/decision/execution 同步（`0a8e6f90`）。
5. 越界守卫：本波改动仅 `docs/workspaces/workspace-020-*/**`；未改 DDL / 迁移台账 / Profile 默认集 / 模块矩阵 / Manifest / `docs/contracts/`（git 提交范围可核对）。

### 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| 1. 合同正文落盘并可核对（时区 / 数字货币落点与机器合同 / 设置字段 / 内嵌默认 zh-CN+auto） | ✅ | `D-001` §2/§3/§4；现状证据见上「成果 3」 |
| 2. 可直接驱动 R2/R3（无歧义、无开放 required） | ✅ | `D-001` §6 消费指引；I-001/I-002/I-005 均 verified（Root D-002） |
| 3. 无越界（DDL/迁移/Profile 默认集/汇率/计费/RT-T03/`docs/contracts/`） | ✅ | 提交范围 `0a8e6f90`；§5 越界声明 |
| 4. 关门自审落盘且 open required = 0 | ✅（本条即审计；required = 0） | 本 A-001；必改项汇总为空 |

### Findings

- **F-001 · VP-020「设置面：数字默认」字段名的冻结解释需留痕确认**（low · recommended · open）
  - 描述：VP-020 首波冻结写「Settings Localization 下时区 / 数字 / 货币默认（归属与字段名由 lead Root R1 冻结，I-020-005）」。合同选择**无独立数字格式字段**（数字格式跟随有效 locale，§4.1），与 VP-020 意图「locale 驱动格式」一致，属 R1 的正当冻结解释；建议在关门决策中书面留痕（本条目即留痕），避免后续歧义。
  - 证据：`D-001` §4.1「数字格式（无独立字段）」；VP-020 首波冻结表。
- **F-002 · 内嵌默认货币映射表扩展边界待 R3 立项确认**（low · recommended · open）
  - 描述：§4.3 无配置时 `zh-CN → CNY`、`en-US → USD`；合同已声明 R3 可扩展且缺省不抛错。映射表的扩展边界（如何种 locale/货币组合必支持）建议在 R3（GOAL-004）立项时显式确认，避免实现期歧义。
  - 证据：`D-001` §4.3。

### 必改项汇总（required 列表）

无（0 条）。

### 结论 + 建议下一步

R1 合同冻结完成且证据可核对；scope 内无 required 必改项，无到期 required 信息项；verdict **pass**。F-001/F-002 为 recommended，不阻断关门，随 R3 立项/关门继续跟踪。建议：用户审阅 `D-001` → 确认 GOAL-002 `status: done`（R1 关门）→ 立项 GOAL-003（R2 时区语义）。如用户要求，可在关门确认前追加本地 grok build（grok-4.6 · high）`source: independent` 复核（本审计为 low-risk 文档合同，self 已足够）。