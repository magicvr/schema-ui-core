# Grok Build · 独立交叉审计提示词（GOAL-041 S2）

> 执行方式：在仓库根 `schema-ui-core` 以 grok build（模型 grok-4.6 · 思考强度 high）调用 `/audit` 技能（`.grok/skills/audit/SKILL.md`），按下方 scope 出具 `source: independent` 意见。本文件由编排器预置，供独立会话只读核验。

## 任务

对以下目标与 scope 做**独立交叉审计**（source: independent，audit_type: design-plan）：

- **目标**：`GOAL-041-w29-api-web-protocol-conformance`
- **工作区**：`workspace-010-design-implementation-conformance`（canonical root `docs/workspaces/workspace-010-design-implementation-conformance/`；root_goal `GOAL-001-design-implementation-conformance`）
- **scope**：S2 证据分类与方案冻结 —— C-001～C-014 的处置类别是否证据充分、唯一且如实；`01-decision/D-002-s2-classification-and-scheme-freeze.md` 冻结范围与 S4/S5 契约是否成立；I-003（分类 verified）/ I-004（upstream gap 不适用）/ I-008（cross 审计）状态是否一致；是否存在被漏判的候选或类别。
- **只读约束**：不修改 `00-meta` status/progress、不改 goal-tree、不改方案正文；只可追加审计意见（03-audit/A-NNN + 03-audit.md 索引）或输出意见文本供代贴。

## 必读文件（按序）

1. `docs/workspaces/workspace-010-design-implementation-conformance/workspace.md`
2. `docs/workspaces/workspace-010-design-implementation-conformance/goal-tree.md`
3. `GOAL-041-w29-api-web-protocol-conformance/00-meta.md`（路线图/审计模式：S2 要求 cross）
4. `GOAL-041-w29-api-web-protocol-conformance/01-decision.md`（I-001～I-009）与 `01-decision/D-001-*.md`、`01-decision/D-002-s2-classification-and-scheme-freeze.md`
5. `GOAL-041-w29-api-web-protocol-conformance/02-execution/E-002-*`、`E-003-*`
6. `GOAL-041-w29-api-web-protocol-conformance/03-audit.md` 与 `03-audit/A-001-s2-classification-self.md`（self 意见）
7. `GOAL-041-w29-api-web-protocol-conformance/attachments/S1-*` 与 **`attachments/S2-candidate-classification.md`（核心被审物）**

## 重点核验清单（建议逐项）

- **C-001**：上游 `conformance/schemas/fixture-suite.schema.json` 与本仓 `docs/schemas/fixture-suite.schema.json` LF 归一化 digest（期望 `a93f500b…`）；`provenance-v2.9.json` 登记路径与 note 计数口径。
- **C-002**：`apps/web/scripts/generate-claim.mjs` 中 `report.pinnedUpstream.artifactVersion = "2.8.0"` vs claim `2.9.0`；`apps/web/public/protocol/conformance-local-report.json` 现状。
- **C-003**：`apps/web/src/protocol/upstream-fixtures.test.ts`（v2.7 provenance.json 兼容回归 + 2.9 线 digest 混标）；`provenance-v2.8.json` 是否无消费者；`apps/web/src/protocol/README.md` 是否过时。
- **C-004**：17 fragment 是否全 2.7 envelope；Web 适配器 `app-manifest.ts` 支持 2.7/2.8/2.9；8 页 2.9；上游 decoupled 协商是否合法；**重点**：生产 `RenderPage`（`apps/web/src/renderer/render.tsx:3026`）是否接线页面级版本+能力协商（`version-negotiate.ts` 仅 fixture）；claim `support.capabilities` 与 `boot.ts HOST_SUPPORT` 7 项 vs 页面 `meta.requiredCapabilities` 全集。
- **C-005**：wallet-entries / dictionary-entries / wallet 的 `$context.route.params.*` 与 `readOnly` 实际使用；digitaloffer 两页「声明未使用」是否上游允许（单向约束）。
- **C-006**：mvp/admin-dogfood 页数与协议；与真实 profile 投影差异；无防漂移守卫。
- **C-007 / C-008**：D-VAL 35/35、representative 10/35、真实 `ForModulesWithFragments` 投影（MVP 6/Admin 22/Demo 14）；S5 验收契约是否如实记录。
- **C-009**：15 注册键全部有注册（`custom-components.schema.test.ts`）；`runtime-schema-validate.ts` 本地 `component` 扩展保持上游 schema 字节不变；key 无 namespace；是否应判 custom-extension-candidate。
- **C-010**：未知标准 node fail-closed vs 未知 custom 内联文字 fallback；对照上游 01-node-protocol §「明显占位」与 renderer-spec §1.1/§1.3。
- **C-011**：`render.tsx CUSTOM_HANDLER_URLS` 5 handler + `CUSTOM_HANDLER_NOT_FOUND` fail-closed vs 上游 07-actions-contract §6。
- **C-012**：内页（dictionary-entries/task-runs/wallet-entries/telegram-operator/users-invites）与 notifications（Host 铃铛）的可达性登记。
- **C-013**：`profile.go` 描述符确认 admin.data-transfer / admin.mfa / admin.login-captcha 无 Pages。
- **C-014**：GOAL-004 95/95 处置的边界语义是否被正确限定为历史。

## 输出要求

按 `.grok/skills/audit/SKILL.md` 与 `skills/prompts/05-independent-audit.md` 的结构输出：
- verdict：pass | conditional | fail（附尺度说明）
- Findings（F-00N；required | recommended；严重度；每条的 evidence 路径）
- 必改项汇总
- 与 self A-001 的异同（重点：F-001 生产页面级能力协商/claim 覆盖不足的判定是否成立）
- 结论与给编排器/用户的下一步

若可写入：在 `GOAL-041-w29-api-web-protocol-conformance/03-audit/` 追加 `A-002-*.md`（source: independent）并更新 `03-audit.md` 索引（不改 status/progress/goal-tree）。
