# Grok Build · 独立交叉审计提示词（GOAL-041 S6 关门）

> 执行方式：在仓库根 `schema-ui-core` 以 grok build（模型 grok-4.6 · 思考强度 high）调用 `/audit` 技能（`.grok/skills/audit/SKILL.md`），出具 `source: independent` 的**关门复审意见**。本文件由编排器预置，供独立会话只读核验。

## 任务

- **目标**：`GOAL-041-w29-api-web-protocol-conformance`
- **工作区**：`workspace-010-design-implementation-conformance`（canonical `docs/workspaces/workspace-010-design-implementation-conformance/`；root_goal `GOAL-001-design-implementation-conformance`）
- **scope**：S6 关门 close-out 全目标复审——① S2 分类与 D-002 冻结是否如实；② S3 C-009 custom 边界（D-003 + `attachments/custom-extension-boundary.md`）是否符合 D-001 §4 要素；③ S4 实现整改是否真实（重点核验：页面级版本+能力门禁接线、claim/HOST_SUPPORT 19 能力、digitaloffer 声明删除、dogfood 守卫、C-010 呈现）；④ S5 运行时验证是否如实（35/35 分母、5 组合 HTTP 快照、失败路径、go 影响判定）；⑤ 审计台账 A-001～A-007 required findings 是否全部按 P-003 合法闭合（F-001 的 accepted-residual→S4 fixed 链）；⑥ 信息项 I-001～I-009 是否与门禁一致。
- **只读约束**：不修改 `00-meta` status/progress、不改 goal-tree、不改方案正文；只可追加审计意见（03-audit/A-NNN + 索引）或输出文本供代贴。

## 必读文件（按序）

1. `docs/workspaces/workspace-010-design-implementation-conformance/workspace.md`
2. `docs/workspaces/workspace-010-design-implementation-conformance/goal-tree.md`
3. `GOAL-041-w29-api-web-protocol-conformance/00-meta.md`（路线图 6 检查点 / 审计模式 cross / 信息就绪）
4. `GOAL-041-w29-api-web-protocol-conformance/01-decision.md`（I-001～I-009）与 `01-decision/D-001/002/003`
5. `GOAL-041-w29-api-web-protocol-conformance/02-execution/E-002～E-006`
6. `GOAL-041-w29-api-web-protocol-conformance/03-audit.md` 与 `03-audit/A-001～A-007`
7. `GOAL-041-w29-api-web-protocol-conformance/attachments/S1-*`、`S2-candidate-classification.md`、`custom-extension-boundary.md`、`S5-coverage-matrix.md`、`S5-manifest-snapshots/*`

## 重点核验清单

- **页面级能力门禁（F-001 fixed）**：`apps/web/src/protocol/load-page.ts` 是否在 D-VAL 后接线 `UNSUPPORTED_PROTOCOL_VERSION` / `MISSING_REQUIRED_CAPABILITY`；`apps/web/src/host/host-support.ts` 19 能力全集与 `boot.ts`/claim 是否一致；35 页是否全部能过门禁（`denominator-render.test.tsx`）。
- **claim 一致性**：`apps/web/public/protocol/conformance-claim.json` 19 能力 / 12 suites / artifactVersion 2.9.0；`claim-artifact.test.ts` C0/C1 通过。
- **provenance / 台账**：`provenance-v2.9.json` fixture-suite 路径 `conformance/schemas/…` + stage3 守卫；`provenance.json` 基线标注；`provenance-v2.8.json` 保留为历史。
- **custom 边界**：`attachments/custom-extension-boundary.md` 是否覆盖 D-001 §4 全部要素（namespace/capability/schema-validator/failure/compatibility/fixtures/退出迁移）+ 用户书面裁决（D-003）。
- **S5 分母**：`denominator-render.test.tsx` 35 页；`s5_manifest_snapshot_test.go` 5 组合（mvp 6/admin 22/demo 14/admin+digitaloffer 25/admin+telegram 24）；快照文件存在且 page 数与断言一致。
- **失败路径**：load-page 负例、C-010 占位、unknown handler fail-closed、D-VAL、manifest 错误码。
- **go 影响**：I-007 无影响不暂挂——变更面是否确实不含 Profile 默认集/模块矩阵/Manifest 装配语义。
- **required 闭合**：F-001（accepted-residual 用户书面 → S4 fixed）、F-002～F-005（fixed）是否有可核对证据；无开放 required。
- **信息项**：I-001～I-009 状态与门禁阶段是否一致。

## 输出要求

按 `.grok/skills/audit/SKILL.md` 与 `skills/prompts/05-independent-audit.md` 结构：verdict（pass | conditional | fail，附尺度）、Findings（F-00N；required | recommended；严重度；evidence 路径）、必改项汇总、与 self A-007 的异同、结论与给编排器/用户的下一步。若可写入：追加 `03-audit/A-008-*.md`（source: independent）并更新索引（不改 status/progress/goal-tree）。
