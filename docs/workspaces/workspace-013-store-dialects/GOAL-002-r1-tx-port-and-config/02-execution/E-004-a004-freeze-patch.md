---
id: E-004
doc: execution-entry
goal: GOAL-002-r1-tx-port-and-config
status: recorded
parent: GOAL-001-store-dialects
created: 2026-08-20
updated: 2026-08-20
version: 1.0.0
---

# E-004 · A-004 响应：冻结合同 v1.2 落盘（2026-08-20）

## 已发生事实

- 用户本轮 `/govern` 指令「响应 GOAL-002 A-004」。D-003 采纳：A-004 全部 finding 走 `fixed`（含 F-003～F-005 recommended）。
- 冻结合同附件升至 **v1.2.0**（时间列按现行秒/毫秒沿用；R2 postgres `Open` 不 apply catalog）。
- A-005 记录闭合证据。GOAL-002 `status` 仍为 `done`；检查点仍 2/2。
- **未修改** `apps/api` 运行时。现行 jobs / operation_log 毫秒列、`INSERT OR IGNORE`、`authsession` 的 `sqlite_master`/`PRAGMA` 均仍在代码中，属 R3/R4 按新规则改写，不是本条目失败。
- Git checkpoint `cd46ce9`（仅 `docs/workspaces/workspace-013-store-dialects/`）。

## 证据

| 主张 | 路径 |
|------|------|
| 合同 v1.2 | `attachments/r1-tx-port-and-config-freeze.md` |
| 取舍 | `01-decision/D-003-a004-freeze-patch.md` |
| 响应意见 | `03-audit/A-005-a004-response.md` |
| 独立审原文（未改 verdict） | `03-audit/A-004-r1-freeze-v1-1-independent.md` |
| checkpoint | `cd46ce9` |

## 计划（非事实）

下一拍属 Root：按 v1.2 立项 R2 前先冻结驱动（Root I-002）。可选 `/audit` 复审 A-005 关闭证据。
