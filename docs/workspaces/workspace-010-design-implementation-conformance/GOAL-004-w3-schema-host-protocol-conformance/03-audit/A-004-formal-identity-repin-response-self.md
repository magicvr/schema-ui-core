---
id: A-004
goal_id: GOAL-004-w3-schema-host-protocol-conformance
title: 响应 · I-003 证据身份纠偏自审（上游审计 0080 V379 外部裁定）
source: self
scope: I-003 证据（v2.8.0 正式 pin 身份）+ E-005 纠偏动作与门禁重跑
verdict: pass
created: 2026-08-13
updated: 2026-08-13
parent: GOAL-004-w3-schema-host-protocol-conformance
version: 0.1.0
---

# A-004 · 响应自审（source: self）· I-003 证据身份纠偏

## 范围与区间

上游 v2.8.1 审计 0080 V379 裁定本仓 I-003 证据所绑定的 `593f625` / content `40690917…`
为 H4 预备身份而非正式 tag（正式 = `521cff8` / `4fae4605…` / artifact `6cdbffcc…` /
fixture `7aacf133…`）。本响应核对并修正 E-004 / 00-meta / 01-decision / 03-audit 中
I-003 的证据身份，记录于 E-005。不改变任何 status 与 progress。

## 成果（有证据）

| 项 | 证据 |
|----|------|
| 上游权威裁定 | `schema-ui-docs` v2.8.1 CHANGELOG V379 + `docs/audit/0080-2026-08-13-review.md` |
| vendored 工件与正式 tag 字节一致（无需重 vendor） | 逐项 sha256 比对（host 四 schema + registry + 四 fixture suite，LF 归一）；fixture digest `7aacf133…` 三方一致 |
| 代码侧 pin 修正 | commit `fd641c6`（6 文件，14+/11-） |
| claim 重生成 | buildId `git:fd641c6…`；contentSha256 `4fae4605…`；canonical digest `bbea0ccf…` |
| 门禁重跑 | claim-artifact 5/5；upstream-fixtures 59/59；upstream-host-fixtures 99/99 零排除；apps/web vitest 862/862；tsc 0 错误 |

## 对照成功标准（S3 检查点）

| 标准 | 状态 | 证据 |
|------|------|------|
| S3 新协议到手（上游增补合并、形成可固定版本、本仓固定引用） | 已达成（E-004），证据身份本轮纠偏为正式 tag | E-004 + E-005 §1–§4 |

## Findings

无新增 required。观察项（不阻断，recommended 级、留待后续）：
- **O-001 · recommended**：`docs/schemas/app-manifest.schema.json` 为 CRLF 而 provenance
  sha 为 LF 归一值；上游其他 vendored 工件亦存在行尾混用。仓库无 `.gitattributes`
  归一化。当前测试路径均有归一化兜底、门禁全过；建议后续统一 vendored 工件行尾或
  加 `.gitattributes`，降低字节级核验的跨平台脆弱性。
- **O-002 · recommended**：`app-navigation.cases.json` 仍为 2.7.0 历史 pin（provenance.json
  权威 `11b01170…`），上游 2.8 线的该 suite 有 protocolVersion 字段推进。是否随 2.8
  线升级属 I-005/S2 范围，本轮不动。

## 必改项汇总

无。

## 结论 + 建议下一步

I-003 证据现与上游权威一致：正式 v2.8.0 = tag @ `521cff8` / content `4fae4605…` /
fixture `7aacf133…`；停止线解除的表述基础恢复干净。S4 生产实现无需任何代码改动。
建议下一步：I-005 台账录入收尾 → S2 出口门禁（I-001/I-002 + cross 方案审视）。
