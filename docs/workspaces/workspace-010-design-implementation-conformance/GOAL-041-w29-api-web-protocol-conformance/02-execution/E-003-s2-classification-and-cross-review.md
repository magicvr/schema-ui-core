---
id: GOAL-041-w29-api-web-protocol-conformance
doc: execution-entry
record_id: E-003
status: recorded
parent: GOAL-001-design-implementation-conformance
created: 2026-09-06
updated: 2026-09-06
version: 0.1.0
---

# E-003 · S2 证据分类与方案冻结（C-001～C-014 + cross 审视）

## 已执行事实

1. **身份复核**：上游 `v2.9.0` 解引用 `81aa1d8954717f4ebdcc695eed6fafaeafcebe8d`（本机只读 clone 实查一致）；本仓 HEAD `1427f9b7`（唯一提交 = S1 文档提交，产品代码快照与 S1 的 `ce66abab` 相同，工作树干净）。
2. **逐项证据收集**：对照上游 schema/registry/fixtures/ADR 与本仓 code/test/file:line，为 C-001～C-014 全部补齐证据并给出唯一处置类别（见 `attachments/S2-candidate-classification.md`）。关键核验点：
   - C-001：上游 `conformance/schemas/fixture-suite.schema.json` 与本仓 `docs/schemas/fixture-suite.schema.json` LF 归一化 digest 均为 `a93f500b…`（重算），provenance 登记路径与 note 计数不符。
   - C-002：`generate-claim.mjs` 中 `report.pinnedUpstream.artifactVersion = "2.8.0"` vs claim `2.9.0`；生成的 `public/protocol/conformance-local-report.json` 现为 2.8.0；无测试断言 2.8.0。
   - C-003：`upstream-fixtures.test.ts` 固定 v2.7 provenance.json（含 2.9 线 digest 混标）；`provenance-v2.8.json` 无消费者；README 版本协商描述过时。
   - C-004：17 fragment 全部 2.7 envelope；Web 适配器支持 2.7/2.8/2.9；8 页声明 2.9；渲染层按页门禁（form-controls.types.ts / permissions.ts / `gateDataRouteBinding`）。
   - C-005：wallet-entries / dictionary-entries / wallet 声明且实际使用（`$context.route.params.*`、`readOnly:true`）；digitaloffer 两页声明未使用（上游单向约束，保守声明合法）。
   - C-006：mvp-dogfood 13 页 / admin-dogfood 26 页（2.7）与真实 profile 投影不同，无防漂移守卫。
   - C-007：D-VAL 全量 35/35 绿；渲染级 representative 10/35 + v2.9 定向测试。
   - C-008：真实 `ForModulesWithFragments` 投影 MVP 6 / Admin 22 / Demo 14；Telegram/Digital Offer 仅显式 custom。
   - C-009：15 注册键全有注册（W25 守卫）；`runtime-schema-validate.ts` 本地 `component` 扩展保持上游 schema 字节不变；无统一 namespace。
   - C-010：未知标准 node fail-closed `RENDER_UNKNOWN_NODE_TYPE`；未知 custom 内联文字 fallback（弱于上游 §1.1/§1.3）。
   - C-011：5 handler 全在白名单 + 非白名单 `CUSTOM_HANDLER_NOT_FOUND` fail-closed（符合 07-actions-contract §6）。
   - C-012：内页全部登记且经 navigate action / Host 铃铛可达；重复/缺失/orphan = 0。
   - C-013：admin.data-transfer / admin.mfa / admin.login-captcha 均为无页面模块（profile.go 描述符核实）。
   - C-014：GOAL-004 95/95 处置为历史边界，不以旧绿灯替代现行核验。
3. **方案冻结**：`01-decision/D-002-s2-classification-and-scheme-freeze.md`——分类口径（含 no-gap 第 6 类）、处置汇总、S5 验收契约、I-004 不适用证据、C-014 方法论规则。
4. **定向验证**（HEAD 复跑，与 S1 证据核对）：
   - Web：`npm test -- --run`（stage3-fixtures / upstream-fixtures / upstream-host-fixtures / app-manifest / representative-pages / all-module-schemas-dval / custom-components.schema）→ **7 files / 487 tests PASS**（2026-09-06）。
   - Go：`go test -count=1 ./internal/manifest ./internal/composition`——manifest 通过；composition 中 `TestShutdownDrainHarnessPostgres` 首次运行遇 `postgres in-flight request failed during drain: EOF`（该测试由 `PG_TEST_*` 门控、需真实 Postgres，属已知 VP-021 drain harness flake，workspace-009 W15 有先例）；**隔离复跑 `go test -count=1 ./internal/manifest` 与 `go test -count=1 -run TestShutdownDrain ./internal/composition` 均 PASS**（ok 0.384s / 7.509s），确认首次失败为 drain flake 而非代码回归。产品代码本波零改动，与 S1 快照一致。
5. **cross 审视**：self [A-001-s2-classification-self](../03-audit/A-001-s2-classification-self.md)（conditional；F-001）→ grok build independent [A-002-s2-classification-independent](../03-audit/A-002-s2-classification-independent.md)（grok-4.6 · high · `/audit`；conditional；F-001 加强 + F-002～F-005）→ 合并响应见 [D-002 §4b](../01-decision/D-002-s2-classification-and-scheme-freeze.md) 与 [A-003-s2-a002-response](../03-audit/A-003-s2-a002-response.md)。
6. **A-002 意见落实（证据卫生修复，属 S2 证据补齐范围）**：
   - **F-002（required）**：`all-module-schemas-dval.test.ts` walker 由一层 `modules/<top>/schema` 改为递归 `**/schema/*.json`（原漏扫 `channel/telegram/schema` 2 页 + `dev/examples/schema` 8 页 = 实扫 25/35）；修复后实测 **35/35 绿（38 tests）**。
   - **F-003（recommended）**：`custom-components.schema.test.ts` walker 同步递归化 + 补 `telegram-admin-tab` side-effect import（原守卫漏扫 telegram 两页 custom 引用）；实测绿。
   - **F-004 / F-005（recommended）**：C-004 主类改标 implementation-gap（模型 no-gap 为上下文）；E-003 计数 ×5 → ×6；C-005 卫生项去向统一为 S3（随 C-009）。

## 产物

- `attachments/S2-candidate-classification.md`（分类矩阵，14 项逐项证据 + 处置）
- `01-decision/D-002-s2-classification-and-scheme-freeze.md`（方案冻结 + cross 合并响应）
- `03-audit/A-001-s2-classification-self.md`（self）· `03-audit/A-002-s2-classification-independent.md`（grok build independent）· `03-audit/A-003-s2-a002-response.md`（合并响应）
- 证据卫生修复：`apps/web/src/protocol/all-module-schemas-dval.test.ts`、`apps/web/src/renderer/custom-components.schema.test.ts`（walker 递归化）

## 阶段结论

S2 分类与方案冻结完成：upstream-protocol-gap = 0（I-004 不适用）；implementation-gap ×6（C-001/002/003/004/006/010，S4 整改清单）；custom-extension-candidate ×1（C-009，S3 P-004）；explicitly-out ×1（C-013）；excluded ×2（C-007/008 S5 契约 + C-014 方法论规则）；no-gap ×3（C-005/011/012）。I-003 → verified。cross 合并后：F-002～F-005 已闭合（fixed），**F-001 按用户 P-004 书面裁决 accepted-residual（S4 承接 + 复审触发）**——S2 required 全闭合，S2 完成（progress 2/6）。

## Git checkpoint

- **checkpoint `2854b898`**（2026-09-06，S2 方案冻结后）：scope = GOAL-041 五件套 + ledger + attachments + goal-tree/workspace 同步 + 两处守卫测试（`all-module-schemas-dval.test.ts` / `custom-components.schema.test.ts`）；验证 = Web 受影响子集 10 files / 511 tests PASS（含 D-VAL 35/35 与 custom 守卫）+ Go `./internal/manifest` PASS + `TestShutdownDrain*` 隔离 PASS。仅暂存显式 owned paths，无 `git add -A`。
