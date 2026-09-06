---
id: GOAL-041-w29-api-web-protocol-conformance
doc: decision
status: active
parent: GOAL-001-design-implementation-conformance
created: 2026-09-06
updated: 2026-09-06
version: 0.2.0
---

# 决策记录 · GOAL-041

## 信息需求与阶段门禁

> 本文件是稳定索引。长决策写入 `01-decision/D-NNN-<slug>.md`。任何候选在证据不足时保持 `open` / `collecting`，不得被写成已验证的协议偏差或协议缺口。

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | `schema-ui-docs@v2.9.0` / `81aa1d8` 与本仓 provenance、schemas、registry、fixtures、claim/test 入口是否同源一致 | S1 分母冻结、S2 分类 | S1 | 对照上游正式工件、digest 与本仓 pinned 副本；区分 v2.7/v2.8 兼容基线与 v2.9 生产消费入口 | **verified**（答案已取得；canonical bytes 对应一致，冲突已登记但尚未分类） | — | E-002；`attachments/S1-protocol-denominator-v2.9.*`；canonical bytes 对应一致，C-001～C-004 待 S2 定性 |
| I-002 | required | API/Web 页面与控件完整分母是什么 | S1、S5 | S1 | 从 runtime Manifest `pages[]`、`schemaUrl`、API providers/fragments、Web route/renderer/custom registry 与测试入口生成目录 | **verified** | — | E-002；`attachments/S1-api-web-page-control-catalog.md` / `S1-api-web-page-control-inventory.json` |
| I-003 | required | 每个候选究竟是 implementation-gap、upstream-protocol-gap、custom-extension-candidate、explicitly-out 或 excluded | S2 方案冻结 | S2 | 逐项给出上游 schema/registry/ADR/fixture 与本仓 code/test/file:line 证据 | **verified**（C-001～C-014 全部给出唯一处置类别；upstream gap = 0） | — | `attachments/S2-candidate-classification.md`；D-002 §1/§2 |
| I-004 | required | 若存在 upstream-protocol-gap，上游是否已 accepted/merged，并有正式版本、commit、schema、registry、fixtures 与迁移/兼容说明 | 受影响项 S3→S4 实施 | S3 | 形成 `attachments/upstream-protocol-augmentation-report.md`；跟踪上游落地并固定本仓 provenance | **不适用**（有证据：S2 分类 upstream-protocol-gap = 0，不建空报告，D-001 §3） | 若后续确认缺口则恢复 | D-002 §2 统计 |
| I-005 | required | custom 候选是否确属上游允许/不负责范围，并具备 namespace、capability、schema/validator、failure、compatibility、fixtures 与退出/迁移条件 | 受影响项 S3→S4 实施 | S3 | 提交逐项方案与取舍，按 P-004 取得用户书面裁决 | **verified**（用户 P-004 书面裁决：本仓合法 custom + 保留现有键 + 新键规范；边界规范齐备） | — | D-003；`attachments/custom-extension-boundary.md` |
| I-006 | required | 页面/控件是否真实经过 Manifest→Schema→Renderer/API 契约链并在失败路径 fail-closed | S5 验收 | S5 | 运行 validator、正反 fixtures、代表性页面/E2E 和 API/Web 定向回归 | **verified**（35/35 D-VAL+Load+Render、5 组合 HTTP Manifest 快照、失败路径、全量回归绿） | — | E-006；`attachments/S5-coverage-matrix.md` + `S5-manifest-snapshots/` |
| I-007 | required | 变更是否影响 Profile 默认集、模块矩阵、Manifest 装配语义或共同门禁解释，从而影响 VP-008 `go` 消费有效性 | S5、S6 | S5 | 对照现行 `go` 候选身份与消费有效性规则，记录暂挂/恢复或无影响结论 | **verified（无影响，不暂挂）** | — | E-006（变更面核对） |
| I-008 | required | `cross` 审计的 self 与 independent provider 是否均有可核对输出，且 required findings 已合法闭合 | S2、S6 | S2 / S6 | 按阶段写入本目标 `03-audit/A-NNN-*.md`；provider 失败不得降级冒充 | open（S2 腿：A-001 self `conditional`；grok build independent 待合并） | provider 按项目级决策固定 | 03-audit 台账 |
| I-009 | non-blocking | 非核心页面、实验控件与未来业务域候选的后续覆盖顺序 | 后续波次 | S6 | 在分类目录中标记 deferred/excluded 与复核触发 | open | 责任人=维护者；复核=本波 S6 或新页面/控件进入生产 Manifest 时 | 待目录 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-09-06 | 协议分流、上游增补报告与协议先行整改门禁 | accepted | `01-decision/D-001-protocol-first-classification-and-remediation-gate.md` |
| D-002 | 2026-09-06 | S2 方案冻结 — C-001～C-014 证据分类与 S4/S5 契约 | accepted | `01-decision/D-002-s2-classification-and-scheme-freeze.md` |
| D-003 | 2026-09-06 | S3 决策 — C-009 custom 边界固定（用户 P-004 书面裁决） | accepted | `01-decision/D-003-s3-custom-boundary-fix.md` |

> legacy inline 的 `## D-NNN` 记录仍可保留并被读取；新记录从目录写入。编号在本目标内单调不复用。
