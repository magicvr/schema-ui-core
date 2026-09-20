---
id: D-019-r3a-i041-008-ruling-and-r3cd-slugs
doc: decision-entry
status: accepted
parent: null
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# D-019 · `I-041-008` 收口裁决与 R3-C/D 子目标 slug 预确认

## 决定的来源

- **用户 2026-09-20 P-004 裁决（本会话）**，两问：
  1. `I-041-008`（公共 wire 输出改为固定 6 位小数的破坏性 / 是否需兼容期）的收口方式；
  2. R3-C/D 子目标编号与 slug（AGENTS §11：新目标 slug 须经用户确认）。
- 依据：`D-018` §5（`I-041-008` required，最晚需要阶段 R3-A）；`D-003`/`VR-091`（固定 6 位输出）；inventory §Web consumer 与 §C2/R3 要求 3。

## 1. `I-041-008` 裁决：**无破坏性、不需兼容期** → `verified`

用户选择 **方案 A**：直接判定「仓内无破坏性依赖，不需兼容期」，信息项以 `verified` 关闭，不设过渡期、不记 residual。

### 判定所依据的证据（全部可复核）

| # | 命题 | 证据 |
|--:|------|------|
| 1 | 前端两个时间消费者都容忍变宽小数 | `apps/web/src/lib/datetime.ts:12-13` `(?:\.\d+)?`；`apps/web/src/host/failure.ts:93` `(?:\.\d+)?Z` |
| 2 | 前端不存在按小数位数截取/长度断言 | `apps/web/src` 全量 grep `split(".")`/`substring(`/`slice(0,1x)`/`length === 2x` → 全是协议版本号或无关路径，无时间宽度依赖 |
| 3 | 仓内 52 处 `.NNNZ` 字面量都是测试**输入**而非输出断言 | 按文件计数（`renderer/schema-crud.test.tsx` 7、`app/representative-pages.integration.test.tsx` 7、`host/return-intent.test.ts` 7 等）；它们被 Web 解析器当作**输入**读入，API 输出形状变化不影响 |
| 4 | 没有对外发布的「3 位小数」合同制品 | `apps/web/dist-lib/**` 为 `.gitignore:113` 忽略的构建产物（`git ls-files` = 0）；`docs/architecture`、`docs/vision`、README 无 3 位小数承诺；历史 3 位小数记载只存在于**已关闭**的旧工作区（workspace-002/010/020），已被 VP-040 取代 |
| 5 | 权威合同本身就是固定 6 位 | `VR-091`（用户 P-004 追加裁决）：「公共 API 时间输出统一为 6 位微秒 RFC3339 UTC `Z`」；`D-003` 同 |
| 6 | 输出精度是**加宽**而非收窄，合法 RFC3339 解析器一律接受 | `D-005` 输入侧已冻结 0/3/6/9 位兼容；本仓 `TestParseWireTimeCompatibility` 即该矩阵 |

### 结论

- 仓内**破坏性 = 0**：无 3 位小数依赖、无固定宽度断言、无发布制品承诺。
- 用户裁决**不设兼容期**、不记 residual；`I-041-008` → `verified`，检查点 R3-A 的最后一项门禁解除。
- 本裁决**不**改变 `D-003`/`VR-091`（输出仍为固定 6 位），也**不**改变 `D-005`（输入兼容矩阵）。

## 2. R3-C / R3-D 子目标 slug 预确认（AGENTS §11）

用户确认沿用提案：

| 阶段 | 编号与 slug | 对应 `D-018` §2 |
|------|-------------|-----------------|
| R3-C | `GOAL-007-r3-pg-cross-version-restore-matrix` | §2 第 4 项：PG 15/16/17 固定版本容器的 `pg_dump`/`pg_restore` 组合矩阵（逐组合 supported/unsupported）+ 判据 4「升级后恢复」有界核对；关闭或 residual `I-041-004` |
| R3-D | `GOAL-008-r3-exit-matrix-and-root-closeout` | §2 第 5 项：六条判据退出矩阵 + self & independent 关门审计 + **用户确认关门**（判据 6） |

- slug 经用户确认后**仍按 `D-018` §6 渐进立项**：R3-C/R3-D 在 `GOAL-006`（R3-A/B）关门后才创建，不预先批量建目标。
- 编号不嵌入工作区号；`parent` 为 `GOAL-001-timestamptz-persistence-contract`。

## 未选方案

- **`I-041-008` 方案 B（记 accepted-residual）**：用户未选。仓内证据为闭集（第 3、4 条），残余仅剩「仓外第三方消费者不可枚举」，用户判定该风险不适用本项目对外承诺范围。
- **`I-041-008` 方案 C（设兼容期/过渡说明）**：用户未选。固定 6 位输出本身已由 `VR-091` 冻结，设过渡期会让「一个合同两种形状」长期并存。
- **R3-C/D 改用其它编号或 slug**：用户未选，沿用提案。
- **R3-C/D 现在就立项**：违反 `D-018` §6 的渐进策略与 P-001「按阶段创建」，未采用。

## 影响与边界

- 本决策关闭 `I-041-008`（required → verified），**放行**检查点 R3-A 的收口判定；R3-B 收口见 `GOAL-006/02-execution/E-004`。
- 本决策**不**关闭 `I-041-004`（R3-C 收集并逐组合记录），**不**推导 Root `done`，**不**替代 R3-D 的用户确认关门（判据 6）。
- Root `progress` 仍按纲领路线图（R1/R2/R3）派生，不因本决策变化。
