---
id: E-038-a044-fixes-and-index-hygiene
doc_type: goal-execution-entry
status: recorded
date: 2026-09-20
parent: GOAL-002-r1-contract-and-denominator-freeze
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-038 · A-044 响应：三处残留修正与索引卫生

## 事实

1. **独立审计 A-044 完成并落盘**（`source: independent`，`auditor: grok-build (grok-4.6 · reasoning high)`）：verdict `conditional`，**开放 required 仍为 2**（F-I-002、F-I-005），并明确判定 **R1 不具备关门条件**。
2. **A-044 的独立 PG 复现扩宽了上一轮结论**：它在 `postgres:15.19` / `16.15` / `17.11` **三个版本**上复现原毫秒式的误差，并找到**误差起点约 8×10¹³ ms**（部分 remainder 偏差达 **±16 µs**），而上一轮我只报了单点（公元 9999 年，+8 µs）。整数拆分式在点名样本上**零误差**；秒族精确；`timestamptz(6)` → `timestamp with time zone` / 精度 6。临时容器 `w040-a044-pg16/15/17`（端口 15441/15442/15443）已由审计方销毁，`compose.yaml` 未改。
3. **A-044 独立复算表数 = 44**，与我的推导一致。
4. **本轮修正（A-045 响应）**：
   - **C3 边界 §5 断言 4**：删去「转换是 fail-closed，B 内不应再有负值」这一**假命题**；改为「B 的 round-trip 不含负值样本，理由是样本按业务意义选取，**不是** B 内不可能有负值」，并把负值断言明确落到 `TestLegacyArtifactMustFail`（A/C）与 `m0`（**仅 voucher `#72`/`#73`** 触发回滚）。
   - **Root `D-015` + child `D-012` 标弃用**：两处均加 **⚠️ 原式已弃用** 块（含受影响量级、三版本复现、`::numeric` 无效、「不得再实现或引用」、指向实测记录）。
   - **mechanism §2 残留「一律不进表达式」**：改写为「不进 `CASE` 分支 + **政策按列分档**」，保留分档表。
   - **child `01-decision.md` 索引补齐**：由 D-012 补至 **D-021**；显式说明 **`D-013`～`D-016` 在本 child 内不存在**（历史未使用、允许空洞、不复用）。信息需求表四条 I-040 状态由过时的「待…」改为**当前真实状态**。
   - **`goal-tree.md` 内部不一致**（本轮自查发现，非审计点名）：纲领路线图行原写 `R1 … 0/4` 而状态表为 `1/4`；已统一为 `1/4`，并把 `primary_plan` 版本由 v0.2.0 更正为 **v0.2.3**，补入 GOAL-002 当前状态说明。
5. **F-I-005 的闭合路径需用户书面裁决**（A-044 判定其**既非 `fixed` 也非完整 `accepted-residual`**）：`D-021` 有范围句但未覆盖全部延期子项、未写复审触发。已在 A-045 §3 列出 A（记 `accepted-residual`）与 B（保持 open）两条路径，**未自行选择**。

## 证据

- A-044 全文：`03-audit/A-044-r1-independent-fi002-fi005-pg-ms-and-d021.md`。
- A-045 响应：`03-audit/A-045-r1-self-response-to-a044.md`。
- 变更载体：`attachments/r1-c3-backup-recovery-boundary-v1.0-fc.md` §5；`GOAL-001/01-decision/D-015-negative-instant-truncation.md`；`01-decision/D-012-v73-allocation-negative-truncation.md`；`attachments/r1-c2-sqlite-rebuild-mechanism-v1.0-fc.md` §2；`01-decision.md`；`goal-tree.md`。
- 本轮 `apps/**` **未修改**。

## 状态评估

- **开放 required = 2**；**本条未闭合任何 required**。
- **R1 关门条件：A-044 判定不具备**（F-I-002 待复审、F-I-005 待用户 residual 裁决）。
- 下一步：取 F-I-005 的用户书面裁决；随后把本轮修正与 `D-021` 一并送 `/audit` 复审。
