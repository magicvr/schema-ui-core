---
id: E-036-a042-negative-policy-scoping
doc_type: goal-execution-entry
status: recorded
date: 2026-09-20
parent: GOAL-002-r1-contract-and-denominator-freeze
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-036 · A-042 响应：负值政策按列分档与表数更正为 44

## 事实

1. **独立审计 A-042 完成并落盘**（`source: independent`，`auditor: grok-build (grok-4.6 · reasoning high)`）：verdict `conditional`，**开放 required 由 3 降至 2**（F-I-002、F-I-005）。
   - **closed**：**F-I-004**（A-040 §G 七项均已对位）、**F-I-027**（三源已收口）。
   - **F-I-002 仍 open**，剩余 = 冻结包把「负值一律 `m0` fail closed」写成**全局规则**，与 Root `D-012`（仅 voucher）／Root `D-015`（负 epoch 合法 instant）**及已绿测试矛盾**。
   - **F-I-025 仍 open**：「20 张带时间列的表」≠ 独立计数 **44**。
   - 无新 required / recommended。
   - A-042 独立跑 `go test ./internal/w040contracttest/ -v -count=1`（在 `apps/api`）→ 三测试全绿；独立核算各边界期望值一致；确认现行测试**未**把 `D-012` 泛化到普通列。
2. **本轮修正（A-043 响应）**：
   - **负值政策按列分档**，修四份冻结载体：`r1-c2-predicate-exact-sql-v1.0-fc.md`（§1 说明块 + §6 `m0` 判定）、`r1-c2-sqlite-rebuild-mechanism-v1.0-fc.md`（§2）、`r1-c2-per-table-rebuild-ddl-v1.0-fc.md`（§0）、`r1-c2-per-column-conversion-contract-v1.0-fc.md`（§3.5）。统一为：**仅** `#72`/`#73`（voucher 两列）负值 fail closed（**Root** `D-012`）；**其余全部时间列正常转换**（**Root** `D-015`：负 epoch 是合法 instant）。另在已 superseded 的 runbook 正文就地标注「`negative-invalid` 已移除」。
   - **表数 20 → 44**：以「倒数第二段 = 表名」对 inventory v0.3 的 90 行机械去重，实测得 **44 张不同的表**（分母仍 90 列）；rebuild 附件 §0 已改并写明推导方法。
3. **根因（如实记录）**：这是**编排器自身的传播缺口**——第 7 轮的可执行测试**已经**暴露出「负值政策被过度泛化」（我当时确实修了测试代码），但**没有**把同一更正传播到冻结包正文，导致「合同正文与已绿测试互相矛盾」。A-042 正是从此处抓出的。
4. **本轮 `apps/**` 未修改**（纯文档收口）。

## 证据

- A-042 全文：`03-audit/A-042-r1-independent-e035-c3-g7-and-fi002-tests.md`。
- A-043 响应：`03-audit/A-043-r1-self-response-to-a042.md`。
- 表数复核：对 `attachments/r1-time-column-inventory-v0.3.md` 机械去重 → **44**（清单见 A-043 §3）。
- 附件变更：上述四份 C2 载体 + superseded runbook 正文标注。

## 状态评估

- **开放 required = 2**（F-I-002、F-I-005）；**本条未闭合任何 required**（F-I-004 / F-I-027 的 closed 由 A-042 判定）。
- 下一步需用户裁决两项（已在 A-043 §5 列出）：
  1. **PG 侧从未经验验证**（本机无 PG 客户端工具、5432 无监听、compose 无 PG 服务；但 Docker 与 `postgres:16`/`15-alpine`/`17-alpine` 镜像在本地，可临时起容器）；
  2. **F-I-005 的关门口径**（哈希须对真实迁移语句求值 → 结构性依赖 R2）。
