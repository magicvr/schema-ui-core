---
id: GOAL-041-w29-api-web-protocol-conformance
doc: decision-entry
record_id: D-002
status: accepted
parent: GOAL-001-design-implementation-conformance
created: 2026-09-06
updated: 2026-09-06
version: 0.1.0
---

# D-002 · S2 方案冻结 — C-001～C-014 证据分类与 S4/S5 契约

## 触发

用户指令推进 GOAL-041 S2：逐项完成 C-001～C-014 的证据分类，并进行 self + independent 的 cross 方案审视（审计模式已在 00-meta 固定为 `cross`；independent provider 按 `docs/architecture/independent-audit-execution.md` 为本地 grok build · grok-4.6 · high）。

## 决定

### 1. 分类口径（在本目标内扩展 D-001 §2）

D-001 §2 的五个处置类别用于**存在偏差**的候选。S2 证据显示部分候选「有证据证明符合、无偏差」，为此在本目标内增补第 6 类结论：

| 类别 | 含义 |
|------|------|
| `no-gap（已核验符合）` | 上游协议 + 本仓实现 + 运行时证据一致，无处置动作；它不是放行信号，而是带证据的分类结论。 |

其余五类语义与 D-001 完全一致。每个候选必须且只能落入六类之一；证据不足保持 `open`/`collecting`。

### 2. C-001～C-014 处置汇总（完整证据见 `attachments/S2-candidate-classification.md`）

| ID | 处置类别 | S 后续 |
|----|----------|--------|
| C-001 | implementation-gap | S4：修正 provenance-v2.9.json 的 fixture-suite 路径与 note 口径 + stage3 守卫断言 |
| C-002 | implementation-gap | S4：`generate-claim.mjs` `pinnedUpstream.artifactVersion` 2.8.0 → 2.9.0 + 重生成 claim 三件套 + 防回归断言 |
| C-003 | implementation-gap | S4：台账卫生（provenance.json 标注兼容基线 / 清理 provenance-v2.8.json 死数据 / README 版本协商描述更新） |
| C-004 | implementation-gap（主类；模型层 no-gap 为上下文） | 版本协商模型合法（2.7 envelope + 2.9 页面 decoupled；渲染层按页门禁）；**S2 新发现（A-001/A-002 F-001）**：生产 `RenderPage` 未接线页面级能力协商（fixture-only），claim/HOST_SUPPORT 能力声明（7 项）未覆盖 35 页能力并集的 **11 项缺口**（actions.batch.request / actions.page.trigger / actions.row.navigate / actions.row.request / actions.upload / form.controls.advanced / form.controls.extended / form.record.load / permissions.inheritance / table.selection / table.sort）→ S4：接线页面级版本+能力门禁 **与** 扩展 claim/HOST_SUPPORT 至能力全集（两项同步做，否则现网页面 fail-closed）+ 一致性守卫 |
| C-005 | no-gap | digitaloffer 两页「声明未使用」为上游允许的保守声明；recommended 卫生项并入 C-009/S3 处置 |
| C-006 | implementation-gap | S4：dogfood ⊆ 模块联合守卫测试（pageId/route/schemaUrl 一致）或改生成/改名 |
| C-007 | excluded | 定义 S5 验收分母契约（§3）；S5 按 I-006 执行 |
| C-008 | excluded | 定义 S5 验收矩阵（§3）；S5 执行并留 HTTP Manifest 快照 |
| C-009 | custom-extension-candidate | S3：逐键固定 namespace/capability/schema/validator/failure/fixtures + 用户 P-004 书面裁决 |
| C-010 | implementation-gap | S4：渲染器未知 custom 呈现对齐上游（明显占位 + console.error 或 fail-closed UNKNOWN_COMPONENT_TYPE） |
| C-011 | no-gap | 白名单 + fail-closed 拒绝符合上游 custom action 契约；namespace 归 C-009 |
| C-012 | no-gap | 内页/Host 铃铛入口全部登记可达，无断链 |
| C-013 | explicitly-out | API-only 模块与 Host 层责任面；UI 经 C-009/C-011 面记账 |
| C-014 | excluded | 方法论规则：本波分类矩阵为现行权威；GOAL-004 仅历史边界语义 |

**upstream-protocol-gap = 0** → 按 D-001 §3 不创建空的上游增补报告；I-004 以「不适用」证据收口。**custom 候选 = C-009**（S3 触发 P-004）。

### 3. S5 验收契约（承接 C-007 / C-008；S2 只定义，S5 执行）

1. **页面分母**：全部 35 页（含 8 个 v2.9 页）逐页执行 validator（D-VAL）→ `loadPageDocument` → `RenderPage` 链路 + 代表性 API handler 触达。**D-VAL 基线经 A-002 F-002 纠正后已修复**：`all-module-schemas-dval.test.ts` 递归遍历全部 `**/schema/*.json`，实测 **35/35 绿（38 tests）**（此前一层 walker 漏扫 telegram×2 + dev/examples×8）。渲染级 representative 10/35 加 v2.9 定向测试，S5 补齐至分母并记录覆盖矩阵。
2. **验收矩阵**：默认三 profile（MVP 6 / Admin 22 / Demo 14）× 显式 custom 模块组合（Telegram、Digital Offer 等），每个组合保留实际 HTTP Manifest 快照（当前无快照，S5 补）。
3. 该契约是 I-006 的执行口径；S5 前不视为已验证。

### 4. 审计（cross）

- S2 方案冻结：self [A-001-s2-classification-self](../03-audit/A-001-s2-classification-self.md) + grok build independent [A-002-s2-classification-independent](../03-audit/A-002-s2-classification-independent.md)（grok-4.6 · high · `/audit`）。
- required findings 按 P-003 三路径合法闭合前，不进入 S3 或 S4 实施；实施须按本决定工作清单，不得扩大范围。

### 4b. cross 意见合并响应（A-001 + A-002 → A-003）

| Finding | source | 处置 | 证据 |
|---------|--------|------|------|
| F-001 生产页面级能力协商缺位 + claim/HOST_SUPPORT 覆盖不足（required · med） | A-001 + A-002（同向无冲突） | **closed（accepted-residual · 用户 P-004 书面裁决 2026-09-06）**：S4 必做（接线页面级版本+能力门禁 与 扩展 claim/HOST_SUPPORT 至能力全集，两项同步实施）；复审触发 = 进入 S4 实施该清单项时以完成证据闭合，最迟 S6 关门前 | D-002 §2 C-004 行（11 项能力并集缺口）；A-003 |
| F-002 D-VAL「35/35」名不副实（required · med） | A-002 | **fixed**：walker 递归化（`all-module-schemas-dval.test.ts`），实测 35/35 绿（38 tests）；D-002 §3 基线同步纠正 | 本决定 §3；`02-execution/E-003` |
| F-003 C-009 W25 守卫未覆盖嵌套模块（recommended · low） | A-002 | **fixed**：`custom-components.schema.test.ts` walker 递归化 + 补 `telegram-admin-tab` import，实测绿 | `E-003` |
| F-004 C-004 双类别 vs 唯一处置（recommended · low） | A-002 | **fixed**：C-004 主类改标 implementation-gap，模型层 no-gap 降为上下文说明 | §2 C-004 行 / 分类矩阵 |
| F-005 台账计数/去向不一致（recommended · low） | A-002 | **fixed**：E-003 计数改 ×6；C-005 卫生项去向统一为 S3（随 C-009） | `E-003` / 分类矩阵 C-005 行 |

- 无 P-004 冲突项（A-001/A-002 同向）。F-001 按用户书面 accepted-residual 闭合（S4 承接 + 复审触发）；F-002～F-005 已 fixed。**S2 required findings 全闭合，S2 可标记完成（progress 2/6）**。

### 5. 方法论规则（承接 C-014）

GOAL-004 的 95/95 处置与旧 fixture 绿灯**不得**作为 W29 现行符合性证明；其 adopt-now / reserve-extension / explicitly-out 语义仅作为历史边界参考。现行分类以上述矩阵与 S5/S6 证据为准。

## 未选方案

- 把 C-004 的 2.7 envelope 判为违规并强推 Manifest 升 2.9（无 2.9 专属 manifest 字段在用；2.7 是最低兼容合同，decoupled 合法）。
- 把 C-005 的「声明未使用」判为 implementation-gap 或上游缺口（上游为单向约束，保守声明合法）。
- 为 C-007/C-008 在 S2 直接实施全量页面测试与 HTTP 快照（属 S5 I-006 验收范围，避免在方案冻结前开展验收）。
- 因 zero upstream gap 而创建空的上游增补报告（D-001 §3 禁止）。

## 影响

- **I-003**（逐项分类）：`collecting` → **verified**（`attachments/S2-candidate-classification.md` 14 项证据齐全）。
- **I-004**（upstream gap）：**不适用**（有证据：upstream-protocol-gap = 0，见分类矩阵统计）。
- **I-008**（cross 审计）：S2 腿 = A-001 self + grok independent；闭合见 03-audit 台账。
- **I-005**（custom 门禁）：C-009 触发 S3，本决定不替代用户裁决。
- S4 工作清单即 §2 表中的 implementation-gap（C-001/002/003/004-残余/006/010）与 recommended 卫生项；S2 期间新发现的 C-004 残余项（生产页面级能力门禁 + claim 能力覆盖）已并入清单。
