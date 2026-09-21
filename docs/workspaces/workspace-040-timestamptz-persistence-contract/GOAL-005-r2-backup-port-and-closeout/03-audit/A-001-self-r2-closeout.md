---
id: A-001-self-r2-closeout
doc: audit-entry
status: active
parent: GOAL-005-r2-backup-port-and-closeout
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
source: self
verdict: conditional
---

# A-001（self）· R2 阶段关门自审（M4 检查点 C）

- **source**: self（编排器自审）
- **日期**: 2026-09-20
- **scope**: R2 整体关门判据（Root `D-016` §4 M4）+ 本子目标 A/B/C；实现 commit `5204ee1f`（M4 A）与 `2a4269d6`（M4 B）；M1/M2 = `GOAL-003`（`done · 4/4`），M3 = `GOAL-004`（`done · 3/3`）。
- **verdict**: `conditional`（**等 independent 关门审计**；本自审不自证闭合）

## R2 判据逐项（可核对）

| 判据 | 证据 | 结果 |
|------|------|------|
| M1 共享 codec + 单测 | `internal/temporal`（11 单测全绿；固定 6 位、负毫秒 floor、999 ms 无进位、公元 9999、sentinel/NULL） | ✅ |
| M2 15 个 descriptor 落码 + 真实 checksum | `modules/*/migration/vp040_temporal.go`；冻结表 15 行；`attachments/r2-v73-v87-generated-statements-v0.1.md`；生成器 16/16 字节级复现 | ✅ |
| M3 仓储/谓词改造 + 双方言回归 + 金额列断言 | `GOAL-004` 全绿；`postgres_test.go` 时间列 `timestamptz(6)`/精度 6、金额列 `bigint`、leftover 21 名；真实 PG 路径通过 | ✅ |
| M4 `D-021` residual 三项完成并复审 | residual ①②③ 经 independent 复审判定可 `fixed`（`GOAL-003/03-audit/A-002`），闭合记录 `GOAL-002/03-audit/A-048` | ✅ |
| M4 Backup Port 类型表面与 provider | `kernel/backup.go`（仅 `CreateRecoveryPoint`）；`internal/backup`（SQLite/PG provider、90 列分母校验、C3 错误分类、restore harness） | ✅ |
| M4 C3 §4.2 调用点（用户裁决「全部接线」） | `store.recovery`/`migrate`/`postgres`：批前 A、批次内 C、批次后 B、§4.3 显式动作与有界重试、PG 对称形状校验 | ✅ |
| R2 self + independent 关门审计 | self = 本文件；**independent 未运行（本轮未发起）** | ❌ pending |

## 成果亮点（可核对）

- **两端真实工件链**：SQLite（`VACUUM INTO` → 新文件 restore → 90 列 + 账本指纹 + 样本）与 PG（容器 `pg_dump -F c` → 新库 `pg_restore` → `information_schema` 90 列 `timestamptz(6)` + 账本指纹）均实测通过。
- **反向断言（C3 §5.1/§6）**：A 类（预转换）与 C 类（批次中途混合形状）工件被以 `TimeContractMismatch` / `TemporalColumnSetIncomplete` 拒绝，且**显式断言不属于** `ArtifactNotFound` / `ArtifactUnreadable` / `ToolFailure`；缺文件单独归 `ArtifactNotFound`。
- **§4.3 可检测性**：删 marker 后重开 → 显式动作 + **恰好一次**补创；B 失败 → 启动不阻断、`RecoveryNote` 记录、**不写 marker**、`HasPoint=false`。
- **真实 PG 锚点端到端**：v86→v87 升级产生 A=1 / C=1 / B=1（真实 dump+restore+校验）并写入 DB 注释 marker；重开 noop；清 marker 后重开补创成功。

## 偏差与开放项（自评，交 independent）

1. **marker 载体**：SQLite 用 sidecar 文件、PG 用 `COMMENT ON DATABASE`（JSON+base64url）。理由：避免新增迁移（catalog 88 会改变指纹/丢账本表集合与 R1 冻结面）。是否符合「机械可核对身份」（C3 §3.1）由复审判定。
2. **C 类 artifact 的命名与 A 分离**：A=`<db>.batch-rollback-*`、C=`<db>.pre-v*`；命名是辅助，真正区分靠**实测形状**（C3 §3.1 判定规则）——本批未实现「把 A/C 喂给校验必须失败」之外的命名级防伪。
3. **PG 逐迁移 dump 的成本**：`snapshotBeforePendingPG` 每条 pending 迁移一次 `pg_dump`（升级 15 条 ≈ 数十秒）；fresh 库跳过。是否为 acceptable 由复审判定。
4. **跨版本矩阵**：实测组合 = server 15.4 + 容器客户端 15.19（同主版本）；`I-041-004`（15/16/17）仍为 R3 前复核项。
5. **`verifyIntegrityPG` 的边界**：当 ledger head ≥ 87 且形状未转换时报错；head < 87（部分升级）时**不**报错（与 §4.3 的「可续跑」一致）——该阈值硬编码 87，复审可判是否应改为「取 catalog 中最后一个 conversion descriptor 的版本」。
6. **`internal/backup` 的 harness 测试在包内**（`package backup`）并导入 `internal/store`；因此 store **不得**导入 backup（已用窄接口与 `internal/temporalcontract` 规避）。若未来有人让 store 依赖 backup，会导致测试包环。
7. **过程卫生**：本轮有一次目录级 `gofmt -w internal/store/` 造成 4 个无关测试文件的纯格式改动，已回滚（提交只含显式编辑路径）。规则：本仓存在既有 gofmt 漂移 + CRLF，禁止目录级 gofmt。
8. **未跑**：R2 关门所需的 independent 审计（本轮未发起）；`GOAL-005` 检查点 C 因此未完成。

## 自审边界（诚实声明）

- 本自审以可执行证据为主（真实 PG、真实 pg_dump/restore、双方言断言），未逐行复核全部 diff；M3 的 4 路并行子代理改造仍以聚合测试 + 抽样为验收依据（见 `GOAL-004/03-audit/A-001`）。
- 本文件**不**闭合任何 required finding，**不**推导 R2 `done`；R2 关门需 independent 关门审计 + 用户确认（Root 成功标准判据 6）。
