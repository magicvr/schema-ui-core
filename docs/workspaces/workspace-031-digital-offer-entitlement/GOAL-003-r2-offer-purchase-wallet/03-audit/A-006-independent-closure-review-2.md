---
doc_type: goal-audit
id: A-006-independent-closure-review-2
parent: GOAL-003-r2-offer-purchase-wallet
date: 2026-09-05
status: closed
version: 1.0.0
source: independent
auditor: codex (gpt-5.6-sol, medium)
audit_type: finding-closure
verdict: conditional
open_required: 0
---

# A-006 · A-004 finding-closure 独立复审（第 2 轮，2026-09-05）

- **source**：independent
- **auditor**：codex (gpt-5.6-sol, medium)
- **类型**：finding-closure
- **scope**：仅复审 A-005 对 A-004 三项缺口的关闭：双库 purchase 矩阵、F-003 search Q 证据、F-004 版本/进度投影
- **verdict**：**conditional**
- **open required**：**0**（仍有 1 项 med recommended 证据缺口；另有 1 项 low recommended 台账异常）

## 范围与区间

本次仅核对 A-004 所列关闭缺口及 A-005 的对应 `fixed` 声明，不重审 GOAL-003 全量实施。核对区间为 2026-09-05 仓库现行内容：A-004、A-005、`purchase_test.go`、GOAL-003 `00-meta.md`、workspace-031 `goal-tree.md` 与 GOAL-003 `03-audit.md`。

动态验证在 `apps/api/` 执行：

1. `go test -count=1 -v -run 'TestPurchasePostgresAcceptance' ./modules/digitaloffer/service/`：**PASS**。本轮 `PG_TEST_*` 可用；输出显示 `single_transaction`、`terminal_failure_paths_leave_zero_residue`、`retry_exhaustion_leaves_zero_residue`、`search_handles_multibyte_and_oversized_Q` 等矩阵子测试均实际挂在 `TestPurchasePostgresAcceptance/...` 下并通过。
2. `go test -count=1 ./modules/digitaloffer/...`：**PASS**（service 包通过，其余无测试文件）。

## 关闭证据核对表

| A-004 缺口 | A-005 声明 | 复审核对 | 结论 |
|------------|------------|----------|------|
| A-002 F-002 未闭合：`runPurchaseMatrix` 硬编码 SQLite，双库矩阵声明不实 | 改为 env 工厂注入；SQLite 与真 PG 各跑一遍；PG 共享 store、per-scenario salt | `runPurchaseMatrix` 已为 `func runPurchaseMatrix(t *testing.T, newEnv func(t *testing.T) *testEnv)`，各矩阵场景均从 `newEnv(t)` 取环境（`apps/api/modules/digitaloffer/service/purchase_test.go:198-204,270-272,406-407,546-547`）。SQLite runner 注入 `newTestEnvOn(t, mustSQLite(t), false)`（同文件 `585-590`）；PG runner 打开 `kernel.DialectPostgres` store，并用捕获的共享 `st` + `pgSeq` salt 创建每场景 env（同文件 `636-680`）。本轮真 PG 动态输出确认完整矩阵子测试位于 `TestPurchasePostgresAcceptance` 下且全部通过。 | **关闭成立。** A-004 主判定已修复，A-002 F-002 可视为关闭。 |
| A-002 F-003 证据缺口：emoji Q 与 100→101 rune 边界断言缺失 | 已补 60 emoji Q、单 emoji 匹配、100→101 rune 边界断言、多字节匹配 | 已有 60 个 `🎫` 的 Q 并断言无错误（同文件 `554-561`），已有 CJK 与单 emoji 匹配结果断言（`563-568`），也构造并核验 100-rune 输入（`570-574`）。但 101-rune 调用仅断言 `err == nil`（`576-580`），没有观察截断后的 Q、返回结果差异，或以可区分数据证明第 101 rune 确实被丢弃；因此 A-005 所称“100→101 rune 边界断言”仍不能直接证明截断语义。另，`🎫` 单个字符是一个 rune；`strings.Repeat("🎫", 60)` 是 60 rune，注释“Emoji are multi-rune”不准确，但该输入仍能覆盖旧的按字节截断可能破坏 UTF-8 的风险。 | **部分关闭。** 多字节/emoji 输入与匹配已补；100→101 rune 截断证据仍有 med 缺口。 |
| A-002 F-004 投影缺口：GOAL-003 分母引用与 goal-tree ASCII 进度滞后 | `00-meta.md` 四处统一 v1.2.0；ASCII 树改为 3/4 | `00-meta.md` 的 `serves_summary`、概述、对齐/实施分母、信息就绪均引用 D-002 v1.2.0（`docs/workspaces/workspace-031-digital-offer-entitlement/GOAL-003-r2-offer-purchase-wallet/00-meta.md:13,20,22,44`）；其 frontmatter/progress 说明均为 3/4（同文件 `9,33`）。ASCII 树 GOAL-003 为 `active · 3/4`，状态表也为 3/4（`docs/workspaces/workspace-031-digital-offer-entitlement/goal-tree.md:10,20`）。 | **关闭成立。** |

## 新 Findings

### F-006-001 · 100→101 rune 边界测试未证明实际截断（recommended · med）

- **证据**：测试只确认 `long` 恰为 100 rune（`apps/api/modules/digitaloffer/service/purchase_test.go:570-574`）；对 `long + "a"` 和 101 个 ASCII rune 的调用仅检查无错误（同文件 `576-580`）。即使实现完全不截断，只要底层查询接受长字符串，这些断言仍会通过。
- **影响**：A-004 对 A-002 F-003 指出的“第 101 rune 被截去”证据缺口尚未完整关闭；不影响本轮已确认的双库矩阵主修复，但不满足 A-005 的完整 `fixed` 表述。
- **建议**：构造可区分的 100/101-rune 查询数据，使“截断到前 100 rune”与“不截断/错误截断”产生不同结果；或对可注入/可观察的规范化查询值做直接断言，并在 SQLite 与 PG 矩阵中保持同一场景。

### F-006-002 · 审计索引重复登记 A-004（recommended · low）

- **证据**：GOAL-003 `03-audit.md` 已在 `:18` 登记 A-004、`:19` 登记 A-005，但 `:20` 又重复登记同一 A-004 并链接到同一文件（`docs/workspaces/workspace-031-digital-offer-entitlement/GOAL-003-r2-offer-purchase-wallet/03-audit.md:18-20`）。
- **影响**：与审计编号共用递增序列及索引唯一性的台账要求不一致，可能造成自动投影或后续审计读取歧义；不改变 A-004 正文 verdict。
- **建议**：由 `/govern` 在响应本意见时去重索引，保留与 A-004 正文一致的一条权威记录。

## 结论 + 建议下一步

**结论：conditional。** A-004 的主 required 缺口（A-002 F-002）已真实关闭：矩阵改为工厂注入，本轮真 PostgreSQL 动态运行明确显示矩阵子测试在 PG acceptance 下全部通过；F-004 投影也已同步。F-003 的 emoji/多字节覆盖已补，但 100→101 rune 测试仍仅证明“不报错”，不能证明第 101 rune 被截去，因此保留 1 项 med recommended 证据缺口；另发现审计索引重复 A-004 的 low recommended 台账异常。

建议 `/govern`：补强 100→101 rune 的可区分结果断言，去重 `03-audit.md` 中重复 A-004 索引，再发起针对 F-006-001 的轻量 closure 复审。GOAL-003 在此之前不宜以 A-005 的“全部 fixed”作为无条件关门依据。

## 声明

本意见仅追加 `source: independent` 的 A-006 与审计索引；不修改 GOAL-003 的 `status`、`progress`、检查点、goal-tree、01/02 台账、D-002 合同正文或业务代码。finding 响应与状态推进由 `/govern` 处理。
