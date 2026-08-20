---
id: E-006
doc: execution-entry
goal: GOAL-002-r1-tx-port-and-config
status: recorded
parent: GOAL-001-store-dialects
created: 2026-08-20
updated: 2026-08-20
version: 1.0.0
---

# E-006 · A-008 响应：冻结合同 v1.4 落盘（2026-08-20）

## 已发生事实

- 用户本轮 `/govern` 指令「响应 GOAL-002 A-008」。D-005 采纳：A-008 全部 finding 走 `fixed`（四条均为 recommended；与 A-002 / A-004 / A-006 响应同口径）。
- 冻结合同附件升至 **v1.4.0**（path 扩展名谓词；`COLLATE NOCASE` 点名；checksum 输入 = sqlite 历史 SQL；嵌套 `Run` 按回调检测）。
- A-009 记录闭合证据。GOAL-002 `status` 仍为 `done`；检查点仍 2/2。
- **未修改** `apps/api` 运行时。现行 `INSERT OR IGNORE`、`COLLATE NOCASE`、`authsession` 的 `sqlite_master`/`PRAGMA`、钱包 `INTEGER` 余额均仍在代码中，属 R3/R4 按规则改写，不是本条目失败。

## 证据

| 主张 | 路径 |
|------|------|
| 合同 v1.4 | `attachments/r1-tx-port-and-config-freeze.md` |
| 取舍 | `01-decision/D-005-a008-freeze-patch.md` |
| 响应意见 | `03-audit/A-009-a008-response.md` |
| 独立审原文（未改 verdict） | `03-audit/A-008-r1-freeze-v1-3-independent.md` |

## 计划（非事实）

下一拍属 Root：按 v1.4 立项 R2 前先冻结驱动（Root I-002）。A-008 为 pass，关闭证据复审可选。
