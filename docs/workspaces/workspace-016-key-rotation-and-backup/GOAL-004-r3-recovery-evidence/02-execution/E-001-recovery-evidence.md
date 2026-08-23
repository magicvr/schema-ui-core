---
id: E-001
doc: execution-entry
goal: GOAL-004-r3-recovery-evidence
status: recorded
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
---

# E-001 · R3 恢复循环自动化与双方言实跑证据

## 事实（2026-08-22）

新增 `internal/composition/post_rotation_recovery_test.go`，按 GOAL-004 D-001 剧本（T0–T5）实现两条可重复证据循环：

### SQLite（无条件必跑）

`TestSQLitePostRotationRecovery` — **PASS（0.55s）**：

1. T0：`testsupport.OpenStore` 真实 catalog 库 + `a.Login` 签出 `access_old[K1]` + `refresh_old`；
2. T1：`VACUUM INTO 'backup.sqlite'`（VP-013 SQLite 合同）；
3. T3/T4：备份文件直接作为 `db.path`，`NewApp(current=K2, cfg.AuthJWTSecretPrevious=K1).Start(ctx)` 全模块 Start+Ready 门禁绿；
4. T5：A1 旧 token 中间件 200 ✓；A2 新签发仅 K2 可验 ✓；A3 `Refresh(refresh_old)` 成功且新 access 为 K2 ✓。

### PostgreSQL（pgtest 门控 + 工具门控）

`TestPostgresPostRotationRecovery` — **PASS（7.8s，`R16_PG_DUMP_CONTAINER=r2-pg-probe`）**：

1. T0：`NewApp(K1)` 真实启动（catalog 48 迁移 + bootstrap）→ 停止 → `store.Open(nil catalog)` + `a.Login` 签出 K1 对；
2. T1：`docker exec r2-pg-probe pg_dump --no-password -h 192.168.31.213 -U sa -F c r16r3a`（17 客户端 → 15.4 服务端，VP-013 记录在案组合）；
3. T3：`CREATE DATABASE r16r3b` + `pg_restore`（stdin 喂入 custom format）。restore 报 **1 条 ignored error**（`SET transaction_timeout` 未识别）——GOAL-006 D-002（workspace-013 A-001 F-004）明文允许该告警类；测试不盲信 exit code，追加 `assertRestoredLedger`：live/restored 两库 `schema_migrations` count + `md5(string_agg(version||':'||checksum ORDER BY version))` 指纹一致才放行；
4. T4：`NewApp(K2, previous=K1)` 从恢复库完整启动，门禁全绿；
5. T5：A1/A2/A3 与 SQLite 同构断言全过（PG 侧另含 W16-F01 must_change_password 归零镜像步骤）。

## 证据判定

两方言「轮换后从既有备份启动 + 鉴权合同」成立：恢复库完整启动 ✓、旧 access 重叠窗可验 ✓、新签发只用 current ✓、opaque refresh 跨轮换+恢复连续 ✓。未新增任何备份实现（仅消费既有合同）。

## 对照检查点

- 检查点 1（I-004 决策）：done（D-001；Root I-004 verified）。
- 检查点 2（自动化循环）：done（本条）。
- 检查点 3（双方言实跑证据落盘）：done（本条即实跑记录；composition/auth/config 三包 `-count=1` 回归 ok，见提交 `1b8e9b0` 前验证）。
- 检查点 4（self 审计 + goal-tree）：进行中。
