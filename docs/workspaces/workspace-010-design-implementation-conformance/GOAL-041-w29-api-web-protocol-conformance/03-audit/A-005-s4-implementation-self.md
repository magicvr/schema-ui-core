---
id: GOAL-041-w29-api-web-protocol-conformance
doc: audit-entry
record_id: A-005
source: self
status: recorded
parent: GOAL-001-design-implementation-conformance
created: 2026-09-06
updated: 2026-09-06
version: 0.1.0
---

# A-005 · S4 实现整改 · self 审计

## A-005 · S4 实现整改（2026-09-06）
- **source**：self
- **auditor**：schema-ui-core 编排器（DeepSeek Harness /govern）
- **类型 / scope**：execution-facts（S4 阶段）；C-001/002/003/004(含 F-001)/006/010 + C-005 子项
- **verdict**：pass

## 范围与区间

按 meta，cross 最低要求落在 S2/S6；S4 为实施阶段，本条为 self 复核实施事实与门禁（independent 对运行时/自定义边界的复审在 S6）。核验对象：E-005 所列实现是否如实、回归是否全绿、F-001 的 accepted-residual 是否以 S4 完成证据闭合、是否有越界改动。

## 成果（有证据）

1. **逐项实现核验**（对 E-005 事实抽查代码）：
   - C-001：provenance-v2.9.json 路径 = `conformance/schemas/fixture-suite.schema.json` + note 计数修正；stage3 守卫断言存在且通过。
   - C-002：`generate-claim.mjs` `pinnedUpstream.artifactVersion = "2.9.0"`；重生成 claim 三件套（buildId `git:13670154…`）C0/C1 校验通过（claim-artifact.test 5 tests）。
   - C-003：provenance.json note 前缀基线标注；README 版本协商节；bootstrap.ts 注释；provenance-v2.8.json 保留为历史（vision/Charter 引用，未删除——正确处置）。
   - C-004/F-001：`host/host-support.ts`（19 能力 + 3 页面版本单源）；`load-page.ts` 页面级协商 fail-closed（2 新错误码）；claim 19 能力/12 suites；负例测试 2 条；**全部 35 页过门禁**（D-VAL + representative 绿，无 fail-closed 误伤）。
   - C-005：digitaloffer 两页未使用声明删除；capability-declaration 守卫仅锁 v2.9 能力对（legacy 保守声明按上游单向约束合法，non-blocking 后续）。
   - C-006：dogfood 守卫上线即捕获真实漂移（dictionary-entries route），fixture 修正 + sha 重算。
   - C-010：未知 custom 明显占位 + console.error + 用例；representative 页补 custom import。
2. **回归全绿**：Web vitest **94 files / 1271 tests PASS**；`tsc -b` + `vite build` 0；Go `go test ./...` 0 FAIL；Go 定向（manifest/composition/digitaloffer）PASS。
3. **无越界改动**：git diff 范围 = GOAL-041 治理文档 + 上述 17 个代码/测试/证据文件 + digitaloffer schema ×2；未改 Profile 默认集/模块矩阵/Manifest 装配语义/共同门禁解释。
4. **F-001 accepted-residual → fixed 闭合**：页面级能力门禁已接线（load-page）+ 支持集/claim 已扩展到 19 能力（host-support + boot + claim）+ 一致性由 host-support↔claim 同步注释与 claim-artifact C1 校验保证；S2 复审触发满足，F-001 可按 `fixed` 闭合。

## 对照成功标准（S4）

| 标准 | 状态 | 证据 |
|------|------|------|
| 只按已固定上游协议或合法 custom 契约修改实现 | 达成 | D-002/D-003 为唯一契约；改动清单见 E-005 |
| 已有协议偏差全部有修复与防复发证据 | 达成 | C-001~C-010 + F-001 修复 + 守卫/负例（stage3/load-page/capability-declaration/dogfood/render） |
| 全量回归可复跑 | 达成 | vitest 1271 / tsc+build 0 / go 全绿 |

## Findings

无 required / recommended findings。

- 说明 1（non-blocking）：claim `support.capabilities` 与 host-support 的同步靠注释约定 + claim C1 校验（注册表未知能力会 fail），未做生成器与常量间的机械一致性测试——建议后续波次加（触发 = 新增 capability 时）。
- 说明 2（non-blocking）：legacy 能力保守声明全量审计为后续项（E-005 附注）。

## 必改项汇总

无。

## 结论 + 建议下一步

S4 通过（0 开放 required）：实现与防复发证据齐备、回归全绿、F-001 已按 S4 完成证据闭合（fixed）。建议下一步：S5 运行时符合性验证（I-006：HTTP Manifest 快照 × profile 矩阵、覆盖矩阵、go 影响 I-007 定稿），随后 S6 cross 关门审计（self + grok independent + 用户确认）。
