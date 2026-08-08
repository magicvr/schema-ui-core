---
id: GOAL-001-full-protocol-contract-v2-7-0
doc: execution-entry
record_id: E-002
status: recorded
parent: null
created: 2026-08-08
updated: 2026-08-08
version: 0.1.0
---

# E-002 · S0 差距盘点（I-001 闭合）

## 2026-08-08 · S0 差距盘点

### 已发生事实

1. 建立 S0 基线：`go test ./...`（apps/api）全包通过；`npm run test`（apps/web，vitest）**24 文件 / 495 测试全绿**；playwright 1.62.1 可用。
2. 逐项比对 `I-PROTO-001 v0.1.3` vs `protocol-inventory-v2.7.0.md` + `docs/schemas/`（6 结构契约）+ `apps/web/src/protocol/upstream/`（15 套 vendor fixture + provenance SHA pin）+ registry（24 type）+ 代码现状（`apps/web` renderer/conformance、`apps/api` handler/modules）。
3. 上游 pin `ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b` **本次可达**：拉取核对 `uploads/cases.json`（13 case）、`permissions-inheritance/cases.json`（17 case）、`docs/02-reaction-expression.md`、`docs/decisions/0022-table-selection-and-batch-request.md`、`docs/decisions/0012-upload-execution.md`、`docs/07-actions-contract.md`，作为 S2/S3 实现权威输入。
4. 差集结论（详见附件 `attachments/I-S0-001-gap-analysis-v0-1-3-to-full.md`）：
   - **无缺口域 7**：D-NODE、D-DATA、D-PERM（17/17）、D-APP（37+16）、D-VER（44+9）、D-VAL、D-COMP 已纳入 18 type 主路径。
   - **保真债 4 项**：D-EXPR 整引擎（reactions 16/16 当前排除）、D-ACT 批量执行（batchRequest 11 case）、D-TABLE 多选 UI（selection 状态机在 conformance 层已有）、D-COMP statCard/chart 渲染。
   - **未纳入 type 6 项**：statCard、chart、inputNumber、datePicker、dateRangePicker、upload。
   - **整域新增 1 项**：D-UPLOAD（前端 + 后端 + 13 fixture + 范例）。
   - **排除面转 include**：D-UPLOAD；reactions 16 case；batchRequest 11 case；6 个 registry type。**差集全部可纳入，无 exclude / 范围收缩需求**（I-002 无收缩 → N/A）。
5. 已记录 fixture 逐套执行现状（17 套：16 行为套件 + scenarios support-only）与 Go 后端缺口（批量端点、上传端点）。

### 证据

| 主张 | 路径 / 命令 |
|------|-------------|
| 差集分析全文 | `attachments/I-S0-001-gap-analysis-v0-1-3-to-full.md` |
| Go 基线全绿 | `cd apps/api && go test ./...`（exit 0，全包 ok） |
| Web 基线全绿 | `cd apps/web && npm run test`（495 passed / 24 files） |
| vendor 套件与 pin | `apps/web/src/protocol/upstream/provenance.json`（commit `ca9e5fe…`，15 artifacts SHA） |
| 上游拉取记录 | 本 E-002 §3 所列 6 个上游文件（pin commit 固定） |
| I-001 状态 | `00-meta.md` 信息表 → closed（证据链本 E-002 + 附件） |
