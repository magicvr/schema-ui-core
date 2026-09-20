---
id: A-005-response-to-a004-and-checkpoint-c-closure
doc: audit-entry
status: active
parent: GOAL-006-r3-wire-formatter-and-unit-family-matrix
created: 2026-09-21
updated: 2026-09-21
version: 0.1.0
---

# A-005 · 编排器响应（A-004 `pass`）+ 检查点 C 关门记录

- **source**: `self`（编排器响应 / 关门记录，非独立意见）
- **日期**: 2026-09-21
- **scope**: 响应 `A-004`（independent · grok build · grok-4.6 · high，verdict `pass`，开放 required = 0）的全部内容，并在 required = 0 的前提下关闭 `GOAL-006` 检查点 C
- **verdict**: `pass`

## 1. 对 A-004 闭合判定的接受

A-004 逐条判定 `F-I-001` / `F-I-002` / `F-I-003` 均 **`fixed`**，并确认：开放 required = **0**；检查点 A 成立；`I-041-008` 与 `I-040-004` 维持 `verified`；可以进入检查点 C 关门；无 P-004 待裁项。编排器**接受**该判定，并逐项核对：

| A-002 finding | A-004 判定 | 编排器核对 |
|---------------|-----------|-----------|
| `F-I-001` PUT mail config | `fixed` | `mail_admin.go` PUT 与 GET 共用 `mailConfigWire`；回归存在且绿 |
| `F-I-002` MFA `enrolledAt` | `fixed` | 缺席 `null`、在场 `FormatWireTime`；Web 类型 `string \| null`；`user_mfa.created_at` 为 `NOT NULL` 且有打戳路径，无「已注册但零值」情形 |
| `F-I-003` 非零 offset | `fixed` | 按已冻结 `D-005` / GOAL-003 `D-001` §3 实现拒绝；一次性探针覆盖 `+08:00`/`-05:00`/`+0000`/`+00:60`/`+24:00`/逗号小数等拼写，无「codec 接受而 offset 探针读不到」的发散行；三个入站调用点未被误伤 |
| `F-I-004` 族身份 | `fixed` | provenance 测试由分母 unit + 生成描述符 `NonNull: false` + DDL 三处旁证 |
| `F-I-005` Web 自实现 | `fixed` | 4 时区「helper ↔ 生产 `formatDate`」对拍 |

## 2. 对两条新 recommended 的响应（均已 fixed）

### `F-I-101`（recommended/med）· **fixed**

A-004 指出 A-003 的主张过满：它实测「错误路径（`json.Marshal(now)`）与库值相等」的概率约 **10%/次**，两次循环的联合假通过约 1%，不是「可忽略」。编排器接受该实测，并**不再依赖巧合**：

- `internal/mail/runtime.go` 增加 `Switcher.SetClock`（与既有 `backup.Service.SetClock` 同形；生产仍为 `time.Now().UTC()`），`Update` 改用该时钟打戳。
- `TestMailConfigPutUpdatedAtIsCanonicalWire` 现在**冻结**时钟到 `2026-09-20T12:57:15.900000Z`，并断言：
  1. 前提可执行——`json.Marshal(frozen)` 确实**不是**规范形（旧实现会给出 `.9Z`）；
  2. 响应 `updated_at` 逐字节等于 `trailingZeroInstant`；
  3. 库中 `mail_config.updated_at` 也等于该值（写路径不丢微秒）；
  4. 响应与库值相等（保留原断言）。
- 于是该回归在**确定性**下失败于默认编解码器，不再有 ~1% 的假通过窗口。

### `F-I-102`（recommended/med）· **fixed（覆盖面扩大 + 措辞收回）**

A-004 用合成文件证明首版守卫存在六类绕过（A 无 tag 字段、B 大写键、C 裸名词键、D 值在下一行、E 变量名含 `Format`、F 运行时下标赋值），并指出「闭集 / 构建期失败」措辞过满。编排器**收回该措辞**（守卫文件头、`E-005`、`A-003` 的相应句子以后者为准）并扩大覆盖：

| 类 | 处置 |
|----|------|
| B 大写键（`"EnrolledAt"`） | **已覆盖**：键正则改为 any-case `*at` |
| C 裸名词键（`"created"`） | **记录为已知缺口**：同词也用于计数与已投影字符串，纳入后产生假阳性（实测 `"updated": len(...)`、`"deleted": deleted`、`"created": row.Created`），故不纳入 |
| D 值在下一行 | **已覆盖**：键后为空时把后续行并入同一语句判值 |
| E 变量名含 `Format` | **已覆盖**：必须匹配真实 formatter **调用**（`FormatWireTime(`/`FormatWireTimePtr(`/`temporal.FormatWire(`/`temporal.MustFormat(`） |
| F 运行时下标赋值 | **已覆盖**：新增 `["...At"] = <value>` 赋值形态扫描 |
| A 无 tag 导出字段随结构体上线 | **记录为已知缺口**：判定「该结构体是否真的进 `writeJSON`」需要类型信息，正则扫描无法在不淹没于内部记录类型假阳性的前提下做到 |

- **决定性验证（本轮自测）**：临时探针文件同时放入 B/D/E/F 四类 → 守卫**失败并逐条点名**（`zz_guard_probe2.go:6/11/17/22`），随后删除（`Test-Path` = False）；`A-002` 那一类的合成文件此前已验证两条守卫各自失败。
- **诚实边界**：两条守卫是**已知绕过类的棘轮**，不是「时间到达公共 wire 已闭集」的证明；A、C 两类为书面缺口，A-004 已确认当前仓内**无**这两类的现行 HTTP 漏网；若将来出现现行漏网，按 findings 升级处理。

## 3. 检查点 C 关门记录

| 判据 | 结论 | 证据 |
|------|------|------|
| self 审计落盘 | 达成 | `A-001`（conditional，后由 A-002 复核并升级 `F-S-004`） |
| independent 审计落盘 | 达成 | `A-002`（`fail`，3 条 required）→ `A-004`（`pass`，开放 required = 0） |
| required 合法闭合 | 达成 | `F-I-001`/`002`/`003` 全部 `fixed`（A-004 判定 + 编排器核对）；无 `accepted-residual` / `user-overruled` 需求 |
| recommended 处置 | 达成 | `F-I-004`/`F-I-005`/`F-I-101`/`F-I-102` 全部 fixed 或书面记录缺口 |
| 关门前验证 | 达成 | `go build ./...` exit 0；`go test -count=1 ./...` 64/64 包 ok；`apps/web` `vitest run` 124 文件 / 1508 用例 ok |

**据此按既有用户裁决（子目标关门属非关键决策，可经交叉审计后静默执行）静默关闭 `GOAL-006`**：`status: done`、`progress: 3/3`。

**关门边界（不得读作 Root 关门）**：
- Root `GOAL-001` 仍为 `active · 2/3`；R3-C（`GOAL-007-r3-pg-cross-version-restore-matrix`）与 R3-D（`GOAL-008-r3-exit-matrix-and-root-closeout`）尚未立项。
- `I-041-004`（PG 15/16/17 跨版本 `pg_restore` 兼容矩阵）仍 **open**，归 R3-C。
- 本目标只承载 R3-A/B；Root 判据 3/4/6 的收口分别在 R3-B（已完成）、R3-C、R3-D。

## 4. 遗留与移交

- **R3-C 待立项**：slug `GOAL-007-r3-pg-cross-version-restore-matrix`（用户 2026-09-20 预确认，Root `D-019` §2）；范围 `D-018` §2 第 4 项。
- **R3-D 待立项**：slug `GOAL-008-r3-exit-matrix-and-root-closeout`；范围 `D-018` §2 第 5 项，含**用户确认关门**（判据 6）。
- 本目标无未闭合 required；`F-I-102` 的两类书面缺口在 R3-C/D 期间无需回访，除非出现现行漏网。
