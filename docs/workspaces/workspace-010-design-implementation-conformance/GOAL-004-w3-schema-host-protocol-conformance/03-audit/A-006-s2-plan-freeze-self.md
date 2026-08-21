---
id: A-006
goal_id: GOAL-004-w3-schema-host-protocol-conformance
title: 自审 · S2 方案出口（source: self）
source: self
scope: S2 方案冻结四门禁（附件 §6）+ D-002 方案冻结 + A-005 findings 整改
verdict: pass
created: 2026-08-13
updated: 2026-08-13
parent: GOAL-004-w3-schema-host-protocol-conformance
version: 0.1.0
---

# A-006 · 自审（source: self）· S2 方案出口

## 范围与区间

对照附件 §6 的 S2 四门禁与 00-meta 审计策略（S2 方案冻结需 self + independent 各一条），
复核 D-002 方案冻结及 A-005（independent，conditional，BLOCKING_COUNT=4）required findings
的整改闭合。

## 成果（有证据）

| 门禁 | 判定 | 证据 |
|------|------|------|
| (a) 每行处置+理由 | 满足 | 附件 §1c 95 行；机械比对 D10 0 差异（A-005 §4） |
| (b) P0 提案可核对 | 满足 | P0 adopt-now 34 项由 accepted ADR + v2.8.0 制品覆盖；registry 仅三个 host capability |
| (c) IMP-001～004 裁定 | 满足 | D6 一致（out/adopt/adopt/reserve） |
| (d) cross 审视闭合 | 满足 | A-005（independent）+ 本 A-006（self）；A-005 F-1～F-4 已闭合（见下） |

## 对照 A-005 findings（required 整改）

| Finding | 处置 | 证据 |
|---------|------|------|
| F-1 I-001/I-002 台账未达 S2 冻结状态 | fixed | 00-meta + 01-decision I-001/I-002 → `verified`（证据 = §1c 95/95 + accepted ADR + A-005 机械比对）；删除 01-decision `proposed` 措辞 |
| F-2 缺 S2 scope self | fixed | 本 A-006（source: self，S2 方案 scope）落盘 |
| F-3 D-002 未入索引 + 预写 accepted-residual | fixed | D-002 登记入 01-decision 索引；residual §3 改为「拟议，S6 用户 P-004 决策」，不预写接受 |
| F-4 S4 清单覆盖不全 + `$deps` residual 过时 | fixed | D-002 §2 补 S4-7（IMP-002）；S4-4 纠错 `$deps`（stage3 零排除 + generate-claim 文本已更新）；residual #3/#5/#6 补入 S4-1/2/3 |

## Findings（本轮新增）

无新增 required。非阻断观察项（转 O-001/O-002/O-003，recommended）：
- **O-001**（承接 A-004）：vendored 工件行尾混用（`app-manifest.schema.json` CRLF），无
  `.gitattributes` 归一化；测试路径已有归一化兜底、门禁全过。建议后续统一 LF。
- **O-002**（承接 A-004）：`app-navigation.cases.json` 仍 2.7.0 历史 pin（`11b01170…`），
  上游 2.8 线有 protocolVersion 字段推进；是否升级属 I-005/S2 后续范围。
- **O-003**（A-005 F-8 归并）：`fixture-suite.schema.json` 等 vendored 工件 LF 一致性，
  与 O-001 同根，合并跟踪。

## 结论 + 建议下一步

S2 方案出口四门禁满足，A-005 的 4 条 required finding 已合法闭合；S2 冻结成立（D-002）。
建议下一步：执行 S4 工作清单（S4-1 return intent 接入、S4-2 reauth-required 映射、
S4-3 hostOwnedPaths、S4-4 纠错、S4-7 IMP-002 证据固化），随后 S5 验证与 S6 关门。
