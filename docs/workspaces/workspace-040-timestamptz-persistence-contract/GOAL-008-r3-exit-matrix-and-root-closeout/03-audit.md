---
id: GOAL-008-r3-exit-matrix-and-root-closeout
doc: audit
status: active
parent: null
created: 2026-09-21
updated: 2026-09-21
version: 0.1.0
---

# 审计台账 · GOAL-008-r3-exit-matrix-and-root-closeout（R3-D）

> 本文件是唯一正式审计台账索引：`self` 与 `independent` **共用** `A-NNN` 序列。
> 每条意见正文在 `03-audit/A-NNN-<slug>.md`；本文件登记条目头（`source`/日期/scope/`verdict`）。
> 独立审计默认只写意见，不修改 `status`/`progress`/方案正文；响应归编排器。

## 意见索引

| A-ID | source | 日期 | scope | verdict | 摘要 | 文件 |
|------|--------|------|-------|---------|------|------|
| — | — | — | — | — | 尚无意见（本目标刚立项） | — |

## 待复审事项（编排器登记，供独立审计取证）

| # | 事项 | 证据位置 | 说明 |
|--:|------|----------|------|
| 1 | 退出矩阵每条判据的证据是否**指向真实产物**，而不是复述 `progress` 或「测试通过」 | 退出矩阵附件 + 各目标台账 | 逐条打开被引用的文件/测试确认 |
| 2 | 判据 1 的主张（合同冻结）是否与 `GOAL-002` 的实际冻结物一致，含 NULL/sentinel/截断语义 | `GOAL-002` `D-002`/`D-003`、inventory | 反例优先：有无被冻结但未落码的语义 |
| 3 | 判据 2（分母迁移 + checksum + 非时间 INTEGER 排除）是否有可核对台账 | `GOAL-003`/`004` 台账、`internal/store/migrate_test.go` 冻结 15 个 checksum | 90 列分母与 v73–v87 逐一对应 |
| 4 | 判据 3（写入读回一致 + VP-020 不漂移 + 含 PG 路径）是否有双方言与真实 PG 证据 | `GOAL-004`/`005`/`006`/`007` | 至少一条真实 PG 路径非 skip |
| 5 | 判据 4（升级后恢复）是否覆盖 SQLite 与 PG 两条路径，且 residual 已合法闭合（不得重记 `F-I-005` 为 open） | `GOAL-005` `A-004`/`A-005`、`GOAL-002` `A-048`、`GOAL-007` 附件 | 容器证据定位必须写明是 CI/reproducibility |
| 6 | **判据 5 的反向核验**：无 ORM/第三数据库/Redis/MQ/多实例；无 Admin 维护提示/业务域混入；驱动类型未泄漏进公共面 | 依赖清单 + 公共契约扫描 | 请自行复跑扫描，不采信叙述 |
| 7 | 残留/例外清账是否完整：已关闭项不得留在「open」，开放项不得被静默忽略 | 各目标 `03-audit` 索引 | 特别是 R1 `D-021` 的 `F-I-005` |
| 8 | 是否存在未合法闭合的 required / 必改 findings（跨全部子目标） | 各目标 `03-audit` 索引 | 三路径闭合：fixed / accepted-residual / user-overruled |
| 9 | **判据 6 的门禁是否被尊重**：Root `done` 是否确实等待用户书面确认 | `00-meta` 信息表 `I-041-010` | 无确认即 `done` = 违规放行 |
| 10 | 台账与投影一致性：goal-tree 树/表/叙述、Root `00-meta`、`docs/vision` 投影是否同步且不越读 | `goal-tree.md`、Root `00-meta.md`、`docs/vision/roadmap.md` | `docs/vision` 不是关门权威，但不得矛盾 |
