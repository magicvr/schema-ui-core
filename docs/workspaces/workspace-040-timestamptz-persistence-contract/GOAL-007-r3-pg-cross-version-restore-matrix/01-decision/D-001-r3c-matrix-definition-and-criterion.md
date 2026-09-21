---
id: D-001-r3c-matrix-definition-and-criterion
doc: decision-entry
status: accepted
parent: GOAL-007-r3-pg-cross-version-restore-matrix
created: 2026-09-21
updated: 2026-09-21
version: 0.1.0
---

# D-001 · R3-C 矩阵定义与 supported/unsupported 判定口径（检查点 A）

## 决定的来源

- Root `D-018` §2 第 4 项冻结了 R3-C 的目标：**PG 15/16/17（固定版本临时容器）上的 `pg_dump`/`pg_restore` 组合，逐组合记录 supported / unsupported**，并复核判据 4「升级后恢复」；§5 说明 `I-041-004` 为 **non-blocking（继承）**，由本目标收集并逐组合记录。
- Root `D-017` §3 约束②（用户裁决）限定证据定位：**容器结果只作 CI/reproducibility，不是生产就绪证据**；并限定破坏性动作只可作用于一次性/专用测试 database。
- 信息项 `I-041-009`（required，本目标新增）要求先冻结「supported / unsupported」判定口径与组合边界；本决策即该项的关闭依据。
- 先决事实：`attachments/r3c-pg-tool-compatibility-probe-v0.1.md`（2026-09-21 实测）。矩阵定义在**看到探测结果之后**冻结，但探测用平凡 schema、只回答工具行为，因此不构成「按结果反推口径」。

## 1. 实测版本事实（冻结 matrix 的取值）

| 轴 | 取值 | 说明 |
|----|------|------|
| server（容器） | `postgres:15-alpine` = **15.19**；`postgres:16` = **16.15**；`postgres:17-alpine` = **17.11** | 每个组合必须**重新记录**实测版本串（tag 可变） |
| client 工具 | `pg_dump`/`pg_restore` **15.19 / 16.15 / 17.11** | 与上面三个镜像同源；宿主无客户端二进制（`D-017` §3 约束③） |
| server（常驻） | **15.4**（`PG_TEST_*`） | major 15，与容器 15.19 同 major；**不**作为矩阵 server 轴取值，其路径由既有 C3 库内测试承担 |

## 2. 矩阵轴与组合枚举

矩阵按 **两个可分离的操作**记录，而不是笼统的「兼容性」：

**操作 P（dump）**：`server_major × client_major` → 3×3 = **9 格**
判据：`pg_dump -F c --no-owner` 是否成功产出归档（退出码 0 且归档存在）。

**操作 R（restore）**：`dumper_major × client_major × target_server_major` → 3×3×3 = **27 格**
判据：`pg_restore --exit-on-error --no-owner -d <新建空库>` 是否退出码 0。
（`--exit-on-error` 与本工作区 C3 路径 `internal/backup/provider_pg.go` 使用的旗标一致；口径必须以实际使用形态为准。）

**操作 U（升级后恢复有界核对）**：`dumper_major < target_server_major` 的对角线上方组合中，取**已知 supported** 的单元格，除退出码外还须校验恢复后的 canonical 形状（见 §4）。

## 3. supported / unsupported 判定口径（冻结）

| 判定 | 定义 |
|------|------|
| **supported** | 该格命令退出码 0，且（对 R）恢复后的 canonical 形状校验通过（§4） |
| **unsupported-toolgate** | 工具自身拒绝：`pg_dump` 报 `aborting because of server version mismatch`；或 `pg_restore` 报 `unsupported version (1.1x) in file header` |
| **unsupported-serverguc** | 工具接受但目标 server 拒绝：恢复会话设置 server 不认识的 GUC（实测为 PG 17 客户端的 `transaction_timeout` → server 15/16），`--exit-on-error` 下退出码非 0 |
| **unexpected-failure** | 不属上述三类却退出码非 0 → **必须**单独记录并作为 finding 处理，不得笼统归入「不支持」 |

口径由**实测行为**分类（`attachments/r3c-pg-tool-compatibility-probe-v0.1.md` §1–§3），不是按预期填表；探测给出的预期规则为：

- P：`client_major ≥ server_major`；
- R：`dumper_major ≤ client_major ≤ target_server_major`。

矩阵本体若出现与上述规则不符的格，按 `unexpected-failure` 处理并逐格记明。

## 4. 恢复后的形状校验（操作 U 与所有 supported 格）

仅凭退出码不足以判定 supported：恢复出来的库必须仍是 VP-040 转换后的合同形状。每个 supported 格至少校验：

1. 迁移台账到 head：`schema_migrations` 最大版本 = **87**，行数与源库一致；
2. 列形状：抽样列的类型为 `timestamptz(6)`（`information_schema.columns.datetime_precision = 6`）；
3. 值形状：抽样瞬时逐字节等于源库值（例如 `.900000` 微秒尾零必须保留，不得变成 `.9`）；
4. sentinel 族：`mail_config.updated_at` / `telegram_config.updated_at` 保持 NULL（legacy `0` 不复活）。

## 5. 驱动与可复现入口（冻结）

- **机制**：`docker network` + `--network-alias`（server 与 client 容器同网）、命名卷交换归档、`pg_isready` 就绪门；每个组合一条可复制的命令（写入矩阵记录的证据列）。
- **schema 来源**：矩阵使用**本仓真实的 VP-040 迁移链**（不是平凡 schema）——server 容器起好后用本仓 store 走完整迁移至 v87，再播种代表性行。这样一次运行同时回答「迁移链在 16/17 上是否成立」。
- **驱动归属**：`apps/api/internal/backup/` 下的 Go 测试，**环境变量门控**（沿用 `VP040_GENERATE=1` 先例），默认不跑（需要 docker 与网络）。矩阵记录由驱动输出，落到 `attachments/r3c-pg-cross-version-matrix-v0.1.md`。
- **失败即停止**：驱动遇到 `unexpected-failure` 必须整体失败（fail closed），不得把它写成 unsupported 继续。
- **资源清理**：驱动自建容器/网络/卷，结束时删除；常驻 15.4 实例只允许创建/删除一次性测试 database。

## 6. `I-041-009` 关闭

`I-041-009`（required）→ **verified**：判定口径与组合边界已按本决策冻结（§2–§3），且先于矩阵本体落盘。

## 7. 本决策**不**裁定的事项

- **不**把任何 unsupported 组合升级为「产品不支持承诺」的裁决——那属产品/愿景层，若矩阵出现需要在发布日期/支持矩阵层面作出承诺的组合，按 P-004 单独提请用户裁决并留痕。
- **不**把容器结果当作生产就绪证据（`D-017` §3 约束②维持）。
- **不**关闭 `I-041-004`：该项由矩阵本体（检查点 B）的逐组合记录关闭。

## 未选方案

- **只跑「同版本」舒适组合**：无法回答 `D-018` §2 第 4 项的「逐组合 supported/unsupported」，未采用。
- **用平凡 schema 代替真实迁移链**：省时间但不能回答「迁移链在 16/17 上是否成立」，且与判据 4「升级后恢复」的证据强度不匹配，未采用。
- **把矩阵做成常驻 CI job**：`D-017` §3 约束②已把容器定位为 CI/reproducibility，但把三容器矩阵塞进默认测试会让本仓 `go test ./...` 依赖 docker 与外部网络；采用环境变量门控，默认跳过。
- **把常驻 15.4 实例作为 server 轴的第 4 取值**：其 major 与容器 15.19 相同，工具层面规则一致，但会引入对共享实例的额外负载与不可复现因素；其应用层路径已由既有 C3 库内测试覆盖，未采用。
