---
title: 决策 · R4 · Schema 驱动 CRUD 与 SQLite 持久化闭环
status: active
created: 2026-08-02
updated: 2026-08-02
parent: GOAL-001-production-admin-foundation
version: 0.2.0
---

# 决策 · GOAL-007

## D-001 · 用一个端到端目标实施 Root D-010

- **日期**：2026-08-02
- **状态**：accepted
- **决定**：
  1. 以单一 R4 子目标承载 records 精确契约、SQLite 持久化、CRUD API、Schema 读写交互、权限/错误状态及重启回归；六个成功标准按依赖顺序推进。
  2. 延续 Root D-010：代表实体固定为 `records`；生产默认迁入 SQLite；重启保持为 required 验收；错误响应继续使用 HTTP status + 稳定 `code` + message 的统一 envelope。
  3. 立项只确定路线与门禁，不在缺少证据时枚举新的 error `code`、DDL、并发策略或 Schema action 映射。
  4. `I-007-001`、`I-007-002`、`I-007-003`、`I-007-004` 均为 required；每项必须在表中所列首个受影响实施或验收动作前由证据关闭并记录后续决策。
- **理由**：API、持久化、Schema action 与重启证据共同构成一个可验证的业务生命周期；拆成多个并列目标会形成无法独立验收的中间态。把未知项显式登记为 required，可在保留端到端交付边界的同时防止方案被代码隐式冻结。
- **实施门禁**：立项时四项均为 `open`。D-002/D-003 已关闭 `I-007-001`/`I-007-002` 并完成 S1/S2 契约冻结；`I-007-003`/`I-007-004` 仍 open，分别阻断首个 Schema 写交互代码与 S6 验收。

### 未选方案

- **按 API / DB / Web / 测试拆成四个并列目标**：依赖紧密且成功边界不可独立成立，会增加跨目标门禁与中间态。
- **沿用进程内 records 并只补 Schema 页面**：无法满足 D-010 的 SQLite 与重启保持 required 边界。
- **立项时先猜精确 error code / DDL / action 形状**：会把尚未收集和验证的信息伪装为决定，违反 P-005。

## D-002 · 冻结 records 精确 API 与错误契约（S1）

- **日期**：2026-08-02
- **状态**：accepted
- **决定**：
  1. 对外实体保持五字段 `id`/`name`/`status`/`owner`/`updatedAt`；不新增 `createdAt`。`id` 与 `updatedAt` 仅服务端管理；可编辑字段仅 `name`/`status`/`owner`（trim 后非空）；`status` 不做枚举白名单。
  2. 继承既有 list/detail/PATCH/DELETE 路径、查询参数、list envelope、权限键与错误 envelope `{"error","message"}`；稳定 code 全表见附件，含已有 `UNAUTHENTICATED`/`FORBIDDEN`/`INVALID_SORT_*`/`INVALID_PAGE*`/`RECORD_NOT_FOUND`/`INVALID_PATCH_*`。
  3. 新增 `POST /api/records`（`records.write`）：body 必填 `name`/`status`/`owner`；成功 **201** + 完整 record；`id` = `rec-` + 16 位小写 hex（`crypto/rand`）；失败 code 冻结为 `INVALID_CREATE_BODY` / `INVALID_CREATE_FIELD`（400），稀有内部失败 `INTERNAL`（500）。禁止把 create 错误复用 `INVALID_PATCH_*`，也不引入 R4 的 `RECORD_CONFLICT`/409。
  4. DELETE 保持 **204** 空体；并发语义为 last-write-wins（无乐观锁/version）；PATCH 忽略未知 JSON 键；body 上限 4 KiB。
  5. 正反矩阵 T-API-01～13 为 S3 实施与回归的最低 API 断言集。
- **理由**：对照现 handler/测试与 Root I-004 M-R4-01～06/08 可完整继承读改删与权限基线；唯一结构性缺口是 create，其 code 与 PATCH 分离可避免前端/测试歧义。不在契约层发明唯一约束或枚举，以免超出当前产品语义。
- **信息门禁**：`I-007-001` → `verified`；证据 [I-007-001-api-error-contract.md](attachments/I-007-001-api-error-contract.md)。本决策完成 **S1** 契约冻结，并放行 S3 中受 API/错误契约约束的代码变更；不构成 S3 已实现，也不关闭 `I-007-003`/`I-007-004`。

### 未选方案

- **create 复用 `INVALID_PATCH_BODY/FIELD`**：混淆操作语义，矩阵与前端映射更难维护。
- **引入 status 枚举或 name 唯一 / 409 CONFLICT**：当前实现与演示数据无此约束，扩大范围且无用户要求。
- **客户端指定 id 或乐观锁**：增加冲突面；与现 PATCH 形状不一致。
- **只收集不冻结**：无法勾选 S1，S3 仍会被 `I-007-001` 阻断。

## D-003 · 冻结 records SQLite DDL、0003 迁移、seed 与 repository（S2）

- **日期**：2026-08-02
- **状态**：accepted
- **决定**：
  1. 在既有迁移链追加 **`0003` / `records_persist` / transformID `0003:records-persist:v1`**；只创建 `records` 表与 `name`/`updated_at`/`owner` 索引，**不**改写 0001/0002 的 SQL 或 checksum 输入。
  2. DDL：`id TEXT PK`，`name`/`status`/`owner` TEXT NOT NULL（trim 非空 CHECK），`updated_at INTEGER` Unix 秒；无 FK、无 name 唯一、无 soft-delete。API RFC3339 ↔ DB Unix 秒在 repository 映射。
  3. 迁移 up 只建空表；业务种子走 **`seedRecords`**：在 `seedAdmin=true` 路径于 `seedRBAC` 之后执行；**仅当表行数为 0** 时插入与现 `staticRecords()` 对齐的 8 行（`rec-1`…`rec-8`）；非空则整段跳过，避免撤销用户删除或覆盖变更。
  4. 生产默认唯一数据源为 SQLite repository；废除进程内切片作为生产路径。写并发 = SQLite 单写者 + last-write-wins。Open 签名保持不变。
  5. 非空文件库在应用 pending（含 0003）前必须有一致性快照 + `integrity_check`；checksum 漂移与迁移事务失败 fail closed。T-DB-01～09 为 S3/S6 最低持久化断言。
- **理由**：复用 R3 runner/ledger/seed 模式可将 records 持久化纳入同一启动路径与恢复口径；空表才 seed 才能同时满足「新库有演示数据」与「删除/更新重启保持」。
- **信息门禁**：`I-007-002` → `verified`；证据 [I-007-002-sqlite-migration-plan.md](attachments/I-007-002-sqlite-migration-plan.md)。本决策完成 **S2** 结构冻结，并放行 S3 持久化代码变更；不构成 repository 已实现或 S6 重启证据，也不关闭 `I-007-003`/`I-007-004`。

### 未选方案

- **迁移内插入 8 条种子**：把演示数据绑死在 checksum 上，后续改种子即漂移；与 R3「迁移 vs seed」分层不一致。
- **按 id 永续 ensure 种子行**：用户 DELETE 后重启会插回，破坏删除持久化与 D-010 重启保持。
- **保留 handler 进程切片作为生产回落**：无法证明 SQLite 默认路径，且双源状态不一致。
- **改写 0001/0002 塞入 records**：破坏已部署库 checksum 与 R3 证据链。
