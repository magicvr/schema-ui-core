---
id: A-009
doc: audit-entry
record_id: A-009
source: self
auditor: /govern · grok-build (grok-4.6)
scope: 响应 A-008；闭合 F-001～F-004（recommended 一并 fixed）
verdict: pass
status: recorded
parent: GOAL-002-r1-tx-port-and-config
created: 2026-08-20
updated: 2026-08-20
version: 1.0.0
audit_type: response
---

# A-009 · 响应 A-008（2026-08-20）

- **source**：self（编排响应，**不是** independent）
- **auditor**：`/govern` · grok-build
- **类型** / **scope**：response · A-008 F-001～F-004
- **verdict**：pass（响应范围：recommended 已按 D-005 + 附件 v1.4 可核对修正）
- **工作区**：`workspace-013-store-dialects`

## 范围与区间

- **covered**：A-008 全部 findings；冻结合同 v1.3 → v1.4 文本；D-005 取舍；A-007 对 A-006 的宽度/`Open` 时序闭合仍有效。
- **excluded**：不改 A-008 / A-006 / A-004 / A-002 原文 verdict/findings；不审 R2 代码（仍无）；不改 GOAL-002 `status`/`progress`。
- **冲突裁决**：无 P-004.2。A-008 为 **pass**、开放 required = 0，无必改否决。recommended 写入合同 vs 登记 R2/R3：用户指令「响应」且未选 residual/overruled；本条与前三次响应同口径选 `fixed`。

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| 用户指令响应 A-008 | 本轮 `/govern` 输入「响应 GOAL-002 A-008」；落盘工作区 = `workspace-013-store-dialects` |
| 合同补丁 | `attachments/r1-tx-port-and-config-freeze.md` v1.4.0 |
| 取舍 | D-005 |
| 未改运行时 | 本回合仅 `docs/workspaces/workspace-013-store-dialects/` |

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| 1. Tx 端口合同：公共类型无 `*sql.Tx`；`Run` 为一事务 | **满足**（v1.4 补嵌套检测；时间宽度仍依 v1.3） | 附件 §2 嵌套检测；§3 时间存储 |
| 2. 配置键名与缺省/fail-closed；path 文件根闭合 | **满足**（v1.4 补扩展名谓词，拦住尚不存在的 `./data`） | 附件 §5 条 5 谓词 6 |
| 3. 未修改应用代码 | **满足** | E-006 |

方言 SQL（upsert/时间单位/时间宽度/`COLLATE NOCASE`/checksum 输入）写入合同后，成功标准 1–2 不被 A-008 recommended 卡住。**仍不等于** RT-P03 全端口（备份/迁移对写未冻）。R2 postgres 本拍证据仍是 Open+Ping，不是现行 `/readyz` 全绿。

## Findings

本条无新 required / recommended。

## 关闭证据表

| finding | 原建议 | 闭合 | 证据 |
|---------|--------|------|------|
| A-008 F-001 path 谓词放行尚不存在的 `./data` | recommended | **fixed** | 附件 §5 条 5 谓词 6（`Ext(Base(path))` 非空）；D-005 决定 1 |
| A-008 F-002 `COLLATE NOCASE` 未点名 | recommended | **fixed** | 附件 §3 点名清单 + 物理 SQL 列表；D-005 决定 2 |
| A-008 F-003 checksum 输入未钉 | recommended | **fixed** | 附件 §4 checksum 段（sqlite 历史 SQL + transform id；postgres 不进 digest）；D-005 决定 3 |
| A-008 F-004 嵌套 `Run` 检测 | recommended | **fixed** | 附件 §2「嵌套检测」；D-005 决定 4 |

## 仍开放项

- 本目标 **开放 required = 0**（A-008 原文即无 required；F-001～F-004 已 `fixed`）。
- Root I-002（PG 驱动）仍 open，阻断 **R2 方案冻结**，不阻断本条响应。
- R3/R4 必须按列把时间列写成 postgres `BIGINT`，点名改写 `COLLATE NOCASE`，checksum 不得改 sqlite 历史文本。现行代码未改。
- R2 本拍不得把现行 HTTP `/readyz` 模块门禁当作 postgres 探测打开的验收；`Run` 不得用 Store 级 `inRun` 检测嵌套。

## 必改项汇总

无未闭合 required。

## 结论 + 建议下一步

A-008 四条 recommended 均已 `fixed`。R2 可以 **对照 v1.4** 立项 Open/配置/Ping（仍受 Root I-002 约束），**不可以**对照 v1.3 省略扩展名谓词或 Store 级嵌套标志，**不可以**对照 v1.2 把 `INTEGER` 时间列抄进 postgres DDL。GOAL-002 保持 `done`，该状态不是冻结质量证明。建议下一拍：`/govern` 立项 R2 并先关 Root I-002。A-008 已为 pass，关闭证据复审可选。R2 方案/实施按项目默认补 independent（grok build）。

## 声明

本意见 `source: self`。不修改 A-008 独立意见正文，不把本条标为 independent。
