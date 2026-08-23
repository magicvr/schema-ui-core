---
id: E-004
doc: execution-entry
goal: GOAL-001-key-rotation-and-backup
status: recorded
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
---

# E-004 · R3 阶段关门

## 事实（2026-08-22）

1. **I-004 关闭**：GOAL-004 D-001 verified —— 最小恢复剧本 T0–T5：备份点在轮换前（K1 运行中），恢复后以 K2+prev=K1 启动；断言 A1（旧 access 重叠窗可验）/ A2（新签发仅 current）/ A3（opaque refresh 跨轮换+恢复连续）。
2. **双方言实跑证据**（子目标 GOAL-004 E-001）：
   - SQLite：`VACUUM INTO` 备份 → 备份文件即库 → `NewApp(K2, prev=K1).Start` 门禁绿 → A1/A2/A3 ✓（`TestSQLitePostRotationRecovery` PASS）。
   - PG：`pg_dump -F c`（17 客户端 → 15.4 服务端，workspace-013 记录在案组合）→ `pg_restore` 至新库（1 条 ignored GUC 告警，按 GOAL-006 D-002 允许；ledger count+checksum 指纹核对一致）→ 恢复库完整启动 → A1/A2/A3 ✓（`TestPostgresPostRotationRecovery` PASS）。
3. **回归**：composition/auth/config 三包 `-count=1` 全 ok（exit 0）。零产品代码新增。
4. **自审**：GOAL-004 A-001（self · close-out）verdict pass，0 required。
5. **状态**：GOAL-004 `done` 4/4；Root 路线图 R3 → 完成；progress 3/5。

## 下一步（计划）

R4（GOAL-005，待开）：默认单密钥仍可用——整合既有证据（缺省 previous 的 config 不变式、compose 可选透传默认空、dev 启动路径），补齐缺口并关门。审计模式 self。
