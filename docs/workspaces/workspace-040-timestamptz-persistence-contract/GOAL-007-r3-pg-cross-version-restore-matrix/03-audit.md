---
id: GOAL-007-r3-pg-cross-version-restore-matrix
doc: audit
status: active
parent: null
created: 2026-09-21
updated: 2026-09-21
version: 0.1.0
---

# 审计台账 · GOAL-007-r3-pg-cross-version-restore-matrix（R3-C）

> 本文件是唯一正式审计台账索引：`self` 与 `independent` **共用** `A-NNN` 序列。
> 每条意见正文在 `03-audit/A-NNN-<slug>.md`；本文件登记条目头（`source`/日期/scope/`verdict`）。
> 独立审计默认只写意见，不修改 `status`/`progress`/方案正文；响应归编排器。

## 意见索引

| A-ID | source | 日期 | scope | verdict | 摘要 | 文件 |
|------|--------|------|-------|---------|------|------|
| A-001 | self | 2026-09-21 | GOAL-007 检查点 A/B（矩阵定义、实测记录、口径更正、`I-041-004` 收口） | conditional | 12 项成果可核对；开放 required = 0；`F-S-001`/`F-S-002` 已 fixed（后者为 `D-001` §3 外推规则的实测更正 → `D-002`）；`F-S-003`～`F-S-006` 提交 independent 复核；`N-001`～`N-005` 已核对为非问题 | `03-audit/A-001-self-r3c-matrix.md` |

## 待复审事项（编排器登记，供独立审计取证）

| # | 事项 | 证据位置 | 说明 |
|--:|------|----------|------|
| 1 | 组合定义是否真的覆盖 server×client 的关键组合，而非只跑「同版本」的舒适组合 | `attachments/r3c-pg-cross-version-matrix-v0.1.md` | 实测为 9 dump 格 + 54 restore 格（6 归档 × 3 client × 3 target），无跳过格 |
| 2 | 「supported / unsupported」判定口径是否事先冻结、是否可复现（命令、版本串、退出码） | `D-001` §2–§3 + 记录 | 口径先于矩阵本体落盘；`D-002` 只更正外推句，未改分类 |
| 3 | unsupported 组合是否**逐条**记录原因，而不是被「全部支持」概括 | 矩阵记录表 | 36 个 unsupported 格逐行含退出码与首条诊断 |
| 4 | 升级后恢复有界核对是否真的跨版本，且有形状/校验证据 | `D-001` §4 + 驱动 | 18 个 supported 格（含 15→16/17 与 16→17 升级恢复）逐格执行 ledger head / `timestamptz(6)` / 微秒逐字节 / sentinel NULL 四项校验 |
| 5 | 破坏性动作是否严格限于一次性/专用测试 database | 驱动 + 环境声明 | 全部为容器内新建库与宿主 `t.TempDir()`；常驻 15.4 未参与 |
| 6 | 是否把 Docker 容器结果误当生产就绪证据（`D-017` §3 约束②） | 附件与 `D-002` §4 | 明示定位为 CI/reproducibility；不作产品级支持承诺 |
| 7 | `I-041-004` 的关闭/residual 是否有对应证据与范围说明 | `01-decision.md` 信息表 + `D-002` §3 | 以 verified 关闭（逐格记录满足信息需求），非 residual |
| 8 | **驱动自身的可信度**：`matrixClassify` 是否会把未知失败误判为 unsupported；`setup-failure` 分支是否会掩盖真实问题 | `pg_cross_version_matrix_test.go` | 请构造反例：未知错误文本是否可能落入三类之一；DB 创建失败是否可能被读成 restore 失败 |
| 9 | **形状校验是否足够**：四项校验能否被「恢复出一个空库」骗过 | 驱动 `matrixVerifyShape` | 空库会让 `schema_migrations` 查询失败或 head 为空 → 应报 problem；请独立确认 |
| 10 | **版本/口径绑定的诚实性**：结论是否被正确限定在本次实测的镜像与旗标形态 | 附件「边界」节 + `D-002` §4 | 禁止把结论外推到未测版本或不同旗标组合 |
| 11 | 源库数据是否真能区分「恢复成功」与「恢复成空表」 | 驱动播种 | 3 行固定 id 数据 + 精确值比较；请核对是否可能两库都为空却判通过 |
| 12 | 迁移链在 16/17 成立的附带结论是否有独立证据 | 驱动 `t.Logf` + 源库事实读取 | 三个 server 各自 `migrations at head 87` 并从库内读回播种值 |
