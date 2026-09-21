---
id: D-002-fi002-user-overruled-record
doc: decision-entry
status: accepted
parent: GOAL-004-r2-repository-and-predicate-rewrites
created: 2026-09-21
updated: 2026-09-21
version: 0.1.0
---

# D-002 · `F-I-002` 的 `user-overruled` 正式载体（结论不变，仅补齐书面依据）

## 决定的来源

- 事项：`GOAL-004` 的 independent 审计 `A-002` 把 `mail.PublicView.UpdatedAt` 由值语义 `time.Time` 改为 `*time.Time`、NULL 投影为 JSON `null`，判为 **recommended** finding（`F-I-002`）。
- 用户 **2026-09-20** P-004 书面裁决：**接受** `*time.Time` + JSON `null` 属 R2 的必要后果。
- 该裁决此前只落在编排器响应 `GOAL-004/03-audit/A-003-response-to-independent-m3.md` 的自述里，**没有独立 `D-` 条目**；`GOAL-008` 的 independent 关门审计 `A-002`（`F-I-004`）指出该载体偏弱。
- 用户 **2026-09-21** 在 Root 关门确认包中**授权补一条 `D-` 条目**（结论不变）。

## 裁决内容（正式记录）

| 项 | 内容 |
|----|------|
| 被裁 finding | `F-I-002`（recommended）：`mail.PublicView.UpdatedAt` 由 `time.Time` 改 `*time.Time`，NULL → JSON `null`，改变 JSON 形状 |
| 用户选择 | **`user-overruled`**：接受该变化是 R2 的**必要后果**，不要求改回 |
| 理由（记录在案） | `mail_config.updated_at` 为 NULL 表示「从未配置」；此前读面会**伪造** `time.UnixMilli(0)`（1970）瞬时。改为 `null` 是诚实投影；Web 客户端**本来就只在字段是 string 时使用**（`mail-admin-tab.tsx`），已容忍 `null` |
| 范围 | 仅限该字段的**类型与 NULL 语义**；不扩展到其它模型的 wire 形状 |
| 边界 | 本裁决**不**豁免固定 6 位输出合同（Root `D-003`）：该字段的**非空值**仍必须走共享 formatter |

## 后续一致性（R3-A/B 已落实，无需再裁）

- R3-A 发现该字段的**非空值**由 `encoding/json` 默认编解码器输出变宽小数（违反 `D-003`），已改为处理器侧投影 `mailConfigResponse` + `mailConfigWire`（GET 与 PUT 共用），**模型类型与 NULL 语义保持不变**（`*time.Time` + JSON `null`）。
- 因此本 `user-overruled` 与 Root `D-003` 的固定 6 位合同**同时成立**：类型/null 语义按 R2 裁决，非空值格式按 R3 合同。

## 未选方案

- **改回值语义 `time.Time`**：会恢复 1970 伪造瞬时，与 R2 的诚实读面目标相反；用户已明确不接受。未采用。
- **把该字段退回 R3 与 wire formatter 同批**：R3-A/B 事实上已把非空值的格式化收进共享 formatter，但**模型类型保持** `*time.Time`；用户接受的是类型与 null 语义，故此选项与最终实现等价而非替代。未采用为单独路径。
- **不补 `D-` 条目、仅保留审计响应自述**：载体偏弱（`GOAL-008/A-002` `F-I-004`），用户已授权补写。未采用。

## 影响与边界

- 本决策**不**重开 `GOAL-004` 的关门结论，**不**改变任何代码，只补齐用户书面依据的落盘载体。
- `GOAL-004/03-audit/A-003` 的修订段继续保留（说明「等用户选择」已不再开放）。
