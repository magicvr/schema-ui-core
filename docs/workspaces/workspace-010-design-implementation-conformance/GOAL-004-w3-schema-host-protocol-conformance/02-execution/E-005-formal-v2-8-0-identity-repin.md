---
id: E-005
goal_id: GOAL-004-w3-schema-host-protocol-conformance
title: v2.8.0 正式身份纠偏 — 重 pin 521cff8 / content 4fae4605 并重生成 claim
status: recorded
created: 2026-08-13
updated: 2026-08-13
parent: GOAL-004-w3-schema-host-protocol-conformance
version: 0.1.0
---

# E-005 · v2.8.0 正式身份纠偏 — 重 pin 521cff8 / content 4fae4605 并重生成 claim

## 背景

上游 `schema-ui-docs` 发布 **v2.8.1**（tag @ `ceb099f`），其审计 0080（V379–V385）裁定：正式
`v2.8.0` 身份为 **tag `v2.8.0` @ `521cff8`**、content `sha256:4fae4605…`、artifact
`sha256:6cdbffcc…`、fixture 树 `7aacf133…`（不变）；本仓此前 pin 的 `593f625` / content
`40690917…` 被裁定为 **H4 预备 commit，不是正式 pin**（回归自 0068/V309）。v2.8.1 为 PATCH，
机器契约未变（schemas、versioned fixtures 与正式 v2.8.0 字节一致；`protocolVersion` 仍为
`"2.8"`，2.7+2.8 协商不受影响）。

## 已完成事实

### 1. 字节级核验（本仓 vendored 工件 vs 正式 tag）

在本地上游仓核对（`git diff 593f625 521cff8 -- docs/schemas conformance` 仅
`conformance/README.md` 2 行措辞，未 vendor 路径）；本仓 vendored 工件逐项 sha256 比对
`521cff8` 结果：

| 工件 | 结果 |
|------|------|
| `docs/schemas/host-bootstrap.schema.json` 等 4 个 host/2.8 schema + `capability-registry.json` | 与 `521cff8` 字节一致 |
| `apps/web/src/protocol/upstream/{host-bootstrap,host-failure,host-conformance-claim,app-manifest}.cases.json`（LF 归一） | 与 `521cff8` 字节一致 |
| fixture 树 digest `7aacf133…` | 三方（593f625 / 521cff8 / 本仓）一致 |

**结论：无需重 vendor，无需改生产代码。** 观察项（不属本轮）：`app-navigation.cases.json`
为 2.7.0 历史 pin（provenance.json `11b01170…`）；`docs/schemas/app-manifest.schema.json`
本地 CRLF 与上游 LF 为行尾差异（LF 归一 sha 与 provenance `34a3354e…` 一致，仓库无
`.gitattributes` 归一）。

### 2. 代码侧 pin 修正（commit `fd641c6`）

- `apps/web/scripts/generate-claim.mjs`：`UPSTREAM_SOURCE_COMMIT` = `521cff8`、
  `UPSTREAM_PROTOCOL_CONTENT_SHA256` = `4fae46058d01bb62d8ff5a17b35f57021a417302c9d8b932916e17ab8acf3c30`、
  fixture `7aacf133…` 不变；
- `apps/web/src/protocol/upstream/provenance-v2.8.json`：sourceCommit + note（正式身份与
  纠偏注记，artifactDigest `6cdbffcc…` 一并记录）；
- `apps/web/src/protocol/upstream/provenance.json`：重 pin note 的 tag 引用；
- `apps/web/src/protocol/app-manifest.ts`：`APP_MANIFEST_SOURCE` → `521cff8`；
- `apps/web/src/protocol/upstream-fixtures.test.ts` / `upstream-host-fixtures.test.ts`：
  注释与期望值同步。

### 3. Claim 重生成（绑定修正 commit）

`node apps/web/scripts/generate-claim.mjs` →
`buildId=git:fd641c6e9dada57fe335ebadd56c33079303b13f`，claim canonical digest
`sha256:bbea0ccfe38f58c753f58f97bd480be44dd385fade337593e1c3d29987eb9620`。
`conformance-claim.json` 的 `protocolArtifact.contentSha256` 现为正式
`4fae4605…`；`conformance.fixtureSha256` 仍为 `7aacf133…`；report/evidence/host 三处
buildId 逐字一致。

### 4. 门禁重跑（2026-08-13）

- `claim-artifact.test.ts` 5/5、`upstream-fixtures.test.ts` 59/59、
  `upstream-host-fixtures.test.ts` **99/99 零排除**；
- apps/web 全量 vitest：**46 文件 862 通过**；`tsc --noEmit` 无错误。
- Go API 侧本轮无改动（pin 身份仅在 web 消费链路与 claim 生成脚本）。

## 阻塞 / 风险

无。上游 v2.8.1 对机器契约零变更，2.7+2.8 协商与 production Host 模块均不受影响。

## 关联信息项

- **I-003**：维持 `verified`；证据更新为正式身份（tag `521cff8` / content `4fae4605…` /
  fixture `7aacf133…`，以上游审计 0080 V379 为权威）。
- **I-005**：上游已交付件不变；本仓 S3 台账录入收尾仍待完成（本轮不涉）。

## 下一步（计划）

1. I-005 本仓台账录入收尾（S3 台账剩余项）；
2. S2 出口门禁：I-001/I-002 逐项闭环 + S2 方案级 cross 审视（I-004 provider 已指定）；
3. S4 残余迭代（E-004 residual 2–6：页面协议 2.7 mandatory residual、return intent 端到端、
   304/ETag 复用、session adapter 终态、hostOwnedPaths）。

## progress

维持 **2/6**。S3 检查点已于 E-004 完成；本轮为同一检查点的证据身份纠偏，不新增检查点。
