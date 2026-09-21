---
id: A-003-response-to-a002-and-checkpoint-c-closure
doc: audit-entry
status: active
parent: GOAL-007-r3-pg-cross-version-restore-matrix
created: 2026-09-21
updated: 2026-09-21
version: 0.1.0
---

# A-003 · 编排器响应（A-002 `conditional`，开放 required = 0）+ 检查点 C 关门记录

- **source**: `self`（编排器响应 / 关门记录）
- **日期**: 2026-09-21
- **scope**: 响应 `A-002`（independent · grok build · grok-4.6 · high）的**全部** findings，并在开放 required = 0 的前提下关闭检查点 C
- **verdict**: `pass`

## 1. 对 A-002 判定的接受

A-002 独立复跑矩阵（`PASS 104.81s`）并取得与附件一致的结果（版本串、9/54 格、18 supported 形状 ok、24 toolgate、12 serverguc、0 unexpected），据此同意检查点 A/B 的**技术判据**、同意 `I-041-004` 维持 `verified`、确认无 P-004 待裁项、并明确**开放 required = 0**。编排器接受该判定。

其「有条件同意进入 C」的两个条件已在本轮完成：台账同步（`F-I-004`）与 `F-I-002` 的口径-实现一致性。三条 recommended 全部 `fixed`（`D-003`），第 4 条为编排器台账动作。

## 2. 逐条响应

| finding | 级别 | 处置 | 证据 |
|---------|------|------|------|
| `F-I-001` `does not exist` 标签过宽 | recommended | **fixed**：`setup-failure` 收窄为「目标 database 不存在（FATAL + database + does not exist）」；连接/超时/认证/磁盘/权限/对象-扩展-角色不存在全部 `unexpected-failure`；新增常驻 `TestMatrixClassifyOracle`（15 例，含审计的 12 条 oracle 输入与两条半边子串负例） | `pg_cross_version_matrix_test.go`；`D-003` §3 |
| `F-I-002` `D-001` §4.4 探针被替换未记录 | recommended | **fixed（字面实现，非补记替换）**：源库播种 `mail_config(id)`/`telegram_config(id)` 单例（`updated_at` NULL），形状校验新增「两列 NOT NULL 行数必须为 0」；保留 `users.locked_until` 的双向探针作为额外强度 | `D-003` §1 |
| `F-I-003` 抽样校验测不到未抽样对象缺失 | recommended | **fixed**：形状校验新增**冻结分母 44 张表**的存在性检查；残余边界（分母之外的对象/索引/扩展、全量 TOC 比对）书面记录在 `D-003` §2 | `D-003` §2 |
| `F-I-004` 台账未与检查点 B 对齐 | recommended | **fixed**：`GOAL-007/00-meta.md` 信息表 `I-041-004` → `verified`；`goal-tree.md` 树/纲领叙述/说明/状态表全部对齐到 `2/3`（C 完成前**不**写 `3/3`）；Root `D-018` §5 的 `I-041-004` 随本目标关门一并更新 | 见 §4 与 `goal-tree.md` |

### 修正后的复跑（证明三处修正未改变实测结论）

`VP040_PG_MATRIX=1 go test -count=1 -run TestPGCrossVersionRestoreMatrix -v -timeout 45m ./internal/backup/` → **PASS（105.07s）**：

- dump 9 格（6 supported）/ restore 54 格（**18 supported 且形状校验全通过**（现已含字面 sentinel 探针与 44 表存在性）/ 24 toolgate / 12 serverguc / **0 unexpected**）；
- 三个源库仍 `migrations at head 87, 87 ledger rows`，规范形样本逐字节一致。

证据已追加到 `attachments/r3c-pg-cross-version-matrix-v0.1.md` §「修正后复跑」。

## 3. A-002 对 `A-001` 提交项的结论（编排器接受）

- `F-S-003`（分类鲁棒性）：**部分成立并降为已知限制** → 已由 `F-I-001` 的收窄 + oracle 测试实质改善；「本轮 0 unexpected 不能证明标签永远正确」被接受为常量级事实。
- `F-S-004`（形状校验）：**「空库/仅 ledger/缺数据行」挡得住；「只恢复抽样对象」挡不住** → 已由 44 表存在性检查缩小；残余边界书面记录（`D-003` §2）。
- `F-S-005`（`D-001` §2 的 27 格 vs 实测 54 格）：同意为 note，`D-002` §2 已交叉引用；**不**回改 `D-001` 正文（保持「先冻结、后更正」的可追溯链）。
- `F-S-006`（常驻 15.4 未入 server 轴）：同意非缺口（同 major；`D-001` §7 已记录该取舍）。
- `N-001`～`N-005`：A-002 全部同意为非问题，编排器不再另行处置。

## 4. 检查点 C 关门记录

| 判据 | 结论 | 证据 |
|------|------|------|
| self 审计落盘 | 达成 | `A-001`（conditional，12 项成果；2 条偏差已 fixed） |
| independent 审计落盘 | 达成 | `A-002`（conditional，**开放 required = 0**，4 条 recommended） |
| required 合法闭合 | 达成 | 无 required；4 条 recommended 全部 `fixed`（`D-003` + 台账同步） |
| 关门前验证 | 达成 | 门控矩阵修正后复跑 **PASS（105.07s）**；默认基线 `go test -count=1 ./...` 64/64 包 ok；门控测试未设变量时 **SKIP** |
| 台账一致性 | 达成 | `00-meta` 信息表、`01-decision` 索引（D-001～D-003）、`02-execution`（E-001/E-003）、`03-audit`（A-001～A-003）、`goal-tree` 树+表+叙述全部对齐 |

**据此按既有用户裁决（子目标关门属非关键决策，可经交叉审计后静默执行）静默关闭 `GOAL-007`**：`status: done`、`progress: 3/3`。

**关门边界（不得读作 Root 关门）**：
- Root `GOAL-001` 仍为 `active · 2/3`；R3 的第三项 **R3-D（`GOAL-008-r3-exit-matrix-and-root-closeout`）尚未立项**，其范围含六条判据退出矩阵与**用户确认关门**（判据 6）。
- 容器矩阵的证据定位仍是 **CI/reproducibility**（`D-017` §3 约束②），不是生产就绪证据。
- **产品级「支持哪些组合」的书面承诺**（例如「17 client 不可恢复到 15/16」「dump 要求 client ≥ server」是否写入发布说明/退出矩阵）按 A-002 的建议**留给 R3-D 作 P-004 用户裁决**，本目标不预先裁定。

## 5. 移交

- R3-C 已交付：矩阵定义（`D-001`）、前置探测附件、实测记录（9+54 格）、规则更正（`D-002`）、三条 recommended 修正（`D-003`）、`I-041-004` 与 `I-041-009` 均 `verified`。
- R3-D 需要：Root 六条判据的逐条证据矩阵（判据 1–5 的证据分别来自 `GOAL-002`/`GOAL-003`–`005`/`GOAL-006`/本目标）、self + independent 关门审计、以及**用户确认关门**。
