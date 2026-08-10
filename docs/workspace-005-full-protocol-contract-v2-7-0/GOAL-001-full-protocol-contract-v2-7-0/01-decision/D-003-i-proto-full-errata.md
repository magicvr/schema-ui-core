---
id: GOAL-001-full-protocol-contract-v2-7-0
doc: decision-entry
record_id: D-003
status: accepted
parent: null
created: 2026-08-10
updated: 2026-08-10
version: 0.1.0
---

## D-003 · I-PROTO-FULL-001 执行分母勘误

### 触发

- workspace-008 S0/S3 复核确认：16 个行为套件的总分母仍为 320，但本地 conformance 为 **318 执行通过 + 2 明确排除**。
- `I-PROTO-FULL-001` v1.0.0 将 app-manifest 写成 `37/37`、将合计写成 `320/320 全绿`，与 `upstream-fixtures.test.ts` 的真实 adapter disposition 不一致。

### 决定

1. 将 `I-PROTO-FULL-001` 升为 **v1.0.1 勘误版**，保留 12/12 能力域、24/24 registry type、16/16 行为套件 `include`。
2. 执行基线改为 **320 total = 318 executed + 2 local adapter excluded**；app-manifest 为 **35/37 executed + 2 excluded**。
3. 两项排除固定为 `m1-missing-app-manifest-capability` 与 `m1-navigation-without-capability`。原因：上游 fixture 期望 `CAPABILITY_REQUIRED`，R3 hand-written host validator 返回 `MISSING_REQUIRED_CAPABILITY`；该 error-envelope 差异位于冻结 R3 子集之外。
4. 这两项是 fixture adapter execution exclusion，不是能力域、registry 或 suite 范围收缩；因此 I-002 仍为 N/A，VP-006 / Root 既有 `closed` / `done` 状态不重开。
5. 复审触发：下一次协议 pin/disposition 变更，或上述任一错误包络发生变化。

### 历史记录纪律

- D-002、E-003、E-005、A-001、A-002 及 VRev 历史原文不改写；它们在 2026-08-08 时记录的 `320/320` / `0 exclude` 主张由本决策和 v1.0.1 统一勘误。
- 当前消费者必须引用 v1.0.1 的 `318+2` 口径，不得继续把 `320/320` 解释为 320 个 case 全部执行通过。
