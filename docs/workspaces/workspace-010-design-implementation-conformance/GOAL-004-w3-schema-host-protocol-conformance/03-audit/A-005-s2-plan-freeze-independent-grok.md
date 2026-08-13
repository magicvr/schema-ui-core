---
id: A-005
goal_id: GOAL-004-w3-schema-host-protocol-conformance
title: 独立审计 · S2 方案出口 cross 审视（grok build · grok 4.6 · xhigh）
source: independent
scope: S2 方案出口四门禁（附件 §6）+ I-001/I-002/I-005/I-006 证据链 + S4 工作清单覆盖
verdict: conditional
provider: grok build（model grok-4.6，reasoning-effort xhigh）
created: 2026-08-13
updated: 2026-08-13
parent: GOAL-004-w3-schema-host-protocol-conformance
version: 0.1.0
---

# A-005 · 独立审计（source: independent · grok build）

> 原文由独立审计会话（grok build，grok 4.6，xhigh）产出，经编排器代贴落盘并保留
> `source: independent`。本轮只读，未修改任何文件。

## 1. verdict: conditional（BLOCKING_COUNT=4）

核心协议事实成立：§1c 与 ADR-0034 D10 **95/95、处置取值 0 差异**；ADR-0035/0036/0037 均
`accepted`；v2.8.0 正式制品（tag `521cff8`）host 三 schema / registry / 三 fixture suite
字节级可核对。**不能无条件放行 S2 冻结**：I-001 仍 `collecting`、S2 scope 的 self 未落盘、
S4 工作清单对 adopt-now 消费项覆盖不全，且工作区有一份未入索引、自称 `accepted` 的 D-002
预写了未经用户书面接受的 `accepted-residual`。

## 2. 逐门禁判定表

| 门禁 | 判定 | 证据 |
|------|------|------|
| (a) 每一行均有 adopt/reserve/out 处置和理由 | **满足** | §1c 95 行均有处置 + 压缩理由；§3 91 能力候选 + §2.2 IMP-001～004 = 95；与 D10 ID 集合一致；§6 四框未勾选，未把 H0 列冒充 S2 出口 |
| (b) P0 项有 schema/状态机/能力/错误/安全/fixtures 可核对提案 | **满足** | P0 adopt-now 34 项由 accepted 0035/0036/0037 + v2.8.0 正式制品覆盖；P0 reserve 8、out 5 未登记 capability；本仓 registry 仅新增三个 host capability，无空能力 |
| (c) IMP-001～004 已裁定为实现偏离或合法 Host extension | **满足（实质）** | 与 ADR-0034 D6 一致：IMP-001=out、IMP-002=adopt、IMP-003=adopt、IMP-004=reserve（overlay ADR） |
| (d) 完成 cross 方案审视，required findings 合法闭合 | **不满足** | A-001～A-003 scope=H0 标签同步、A-004=I-003 身份；缺 S2 方案 scope 的 self；I-001 collecting；independent 本轮 4 条 required 未闭合 |

## 3. Findings

### F-1 · P1（required）· I-001 / I-002 台账未达 S2 冻结可放行状态
- 位置：`00-meta.md` 信息表 + `01-decision.md`
- 证据：I-001 在 00-meta / 01-decision 均为 `collecting`；00-meta I-002 `verified` 但 01-decision
  I-002 仍 `collecting` 且证据写 ADR-0034 D10 `proposed`
- 建议：以同一证据（§1c 95/95 + D10 处置 + accepted ADR）将 I-001 置 `verified`；01-decision
  I-002 与 00-meta 对齐为 `verified`，删除 `proposed` 措辞

### F-2 · P1（required）· 缺 S2 方案 scope 的 self，cross 未完成
- 位置：`03-audit.md` 台账 / A-001～A-004
- 证据：A-001 scope=H0 标签同步，A-002/A-003 同 H0，A-004=I-003 身份；00-meta 审计策略要求
  S2 方案冻结至少一条 self + 一条 independent；03-audit 结论也写「S2 方案冻结仍需后续 cross 审计」
- 建议：补一条 S2 方案 scope 的 self（source: self），与本 independent 组成 cross

### F-3 · P1（required）· 未入索引的 D-002 自称 accepted，且预写未经用户书面接受的 residual
- 位置：`01-decision/D-002-s2-plan-freeze-disposition-and-s4-worklist.md`（git 未索引；01-decision.md 索引只有 D-001）
- 证据：frontmatter `status: accepted`；§3.1 把「account-locked 生产源缺位」写为
  `accepted-residual`，但无用户书面接受（P-003/P-004）
- 建议：D-002 登记入 01-decision 索引；`accepted-residual` 改为「拟议 residual，待用户 P-004
  书面决策（S6 关门时点）」；D-002 的 S2 冻结结论在 cross findings 闭合前不视为已放行

### F-4 · P1（required）· S4 工作清单未覆盖全部 adopt-now 且本仓需实现的项目
- 位置：`00-meta.md` S4 检查点（一句话）、`02-execution/E-004` 已登记 residual、D-002 §1.4/§2
- 证据：
  - IMP-002（adopt-now，D6=实现偏离）：`provider.go:90` `Label: "Users"` 与
    `fragment.json:25-30` `label`/`labelKey` 并存；Web 消费已按 labelKey 命中投影，但 provider
    `NavigationContribution.Label` 未在 S4 清单给出动作；
  - E-004 residual #2（multi-round `$deps` 未实现）与 `stage3-fixtures.test.ts:351-356`
    reactions 零排除矛盾（过时）；`generate-claim.mjs` residuals 仍写「子集未闭环」；
  - E-004 residual #3/#5/#6（return intent 接入、session adapter、hostOwnedPaths）未进 S4 清单
- 建议：IMP-002 显式写入 S4 清单；按 stage3 事实纠错 `$deps` residual 并更新 claim 文本；
  residual #3/#5/#6 补进 D-002 §2

### F-5 · P2（recommended，不计 BLOCKING）· I-005 台账未收尾而 00-meta 已勾选 S3
- 位置：`00-meta.md` S3 检查点 + I-005、附件 §6 S3 框
- 证据：I-005 `collecting`（最晚 S3 固定前）；00-meta S3 已 `[x]` 且正文写 tag `593f625`
  （E-005/A-004 已纠为 H4 预备身份）；附件 §6 S3 框仍 `[ ]` 且写「I-003/I-005 均为 open」
- 建议：S3 正文改为 `521cff8`；完成 I-005 台账录入（维持 S3 勾选并注明「S3 已固定、I-005 台账后置闭合」）

### F-6 · P2 · 附件 §6 / §2.2 / progress 说明陈旧
- 位置：附件 `I-HOST-APP-001` §6、§2.2；`00-meta.md` progress 说明
- 证据：§6 写 ADR-0034 D10 `proposed`、S2/S3 未通过、I-003/I-005 open；§2.2 IMP-004「未裁定」；
  progress 字段 `2/6` 但正文写 `progress: 1/6`
- 建议：按 accept + v2.8.0 更新 §6 H0 注记与 §2.2 裁决句；progress 正文对齐 `2/6`

### F-7 · P2 · host fixtures 计数应为 96 而非 99
- 位置：`00-meta.md` I-005、`E-004`、`D-002`
- 证据：host-bootstrap 23 + host-failure 43 + host-conformance-claim 30 = **96**（app-manifest 41）
- 建议：统一为 96（99 实为 96 fixture + 3 条 sha256 pin 断言的测试数）

### F-8 · P2 · `app-manifest.schema.json` 行尾（O-001 观察项）
- 位置：`docs/schemas/app-manifest.schema.json`（A-004 O-001）
- 证据：原始 sha256 ≠ provenance pin；去 CR 后 = pin `34a3354e…`；同 directory 其它 schema 为 LF
- 建议：统一 LF 或加 `.gitattributes`（非 S2 阻断）

## 4. 机械比对结果（§1c vs ADR-0034 D10）

| 项 | 值 |
|----|----|
| §1c ID 数 | 95 |
| D10 ID 数 | 95 |
| §1c 缺失 | 0 |
| D10 缺失 | 0 |
| 处置取值差异 | 0 |
| 双方分布 | adopt-now=43 / reserve-extension=46 / explicitly-out=6 |
| 来源 | §3 91 + IMP-001～004 = 95 |

## 5. 结论

**按现状不放行 S2 冻结**。放行前需闭合 F-1～F-4（required）并经 `/govern` 落盘；F-5～F-8 建议
一并处理。核心协议事实（95/95、accepted ADR、正式制品 pin）无异议，无 P0。

**BLOCKING_COUNT=4**

## 编排器响应（2026-08-13）

| Finding | 处置 | 说明 |
|---------|------|------|
| F-1 (P1 required) | fixed | 00-meta + 01-decision 的 I-001/I-002 均置 `verified`（证据 = §1c 95/95 + accepted ADR-0034～0037 + 本 A-005 机械比对）；删除 01-decision I-002 的 `proposed` 措辞 |
| F-2 (P1 required) | fixed | 补 S2 方案 scope 的 self：A-006（source: self，pass）；与本 A-005 组成 cross |
| F-3 (P1 required) | fixed | D-002 登记入 01-decision 索引；§3 residual 改为「拟议，S6 用户 P-004 书面决策」，不预写 `accepted-residual` |
| F-4 (P1 required) | fixed | D-002 §2 补 S4-7（IMP-002 证据固化）；S4-4 按 stage3 事实纠错 `$deps` residual（引擎已 e18edce 实现，零排除）并更新 `generate-claim.mjs` 文本；residual #3/#5/#6 补入 S4-1/2/3 |
| F-5 (P2) | fixed | I-005 → `verified`（E-006 录入消费证据）；00-meta S3 正文 tag 改 `521cff8`；附件 §6 S3 四框勾选 |
| F-6 (P2) | fixed | 附件 §6 H0 注记改 `accepted`；§2.2 IMP-004 改「已裁定 reserve」；progress 正文对齐 `2/6` |
| F-7 (P2) | fixed | host fixtures 计数统一为 96（E-004、00-meta I-005、D-002） |
| F-8 (P2) | acknowledged | 行尾观察项并入 A-006 O-001/O-003（recommended，非阻断） |

**复核（A-006，self，pass）：** F-1～F-4 required 整改经自审验证闭合；S2 方案出口四门禁满足。
