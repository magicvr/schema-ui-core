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
| A-001 | self | 2026-09-21 | GOAL-008 检查点 A（退出判据矩阵、判据 5 反向核验、残留清账）+ Root 关门就绪度 | conditional | 9 项成果可核对；开放 required = 0；两条限定（`L-1` 本环境绑定 / `L-2` 变更范围口径）已写入矩阵；`F-S-001`～`F-S-005` 提交 independent 复核；判据 6 未被越读，Root 关门仍待用户确认 | `03-audit/A-001-self-r3d-exit-matrix.md` |
| A-002 | independent（grok build · grok-4.6 · high） | 2026-09-21 | 关门审计：R3-D 退出矩阵六判据 + 跨目标开放 required + Root 关门就绪度（含独立复跑四个真实路径测试与判据 5 四类扫描） | **conditional**（开放 required = **0**） | 判据 1–5 **实质满足**；判据 6 部分满足（用户确认未发生）。**同意把 Root 交给用户确认关门**，**不同意**自行标 `done`。跨 GOAL-002～007 开放 required 逐目标核对 = 0。6 条 recommended：`F-I-001`（矩阵引用死链）、`F-I-002`（Root 信息表未同步）、`F-I-003`（goal-tree 说明段矛盾）、`F-I-004`（GOAL-004 overruled 载体偏弱）、`F-I-005`（GOAL-002 索引缺 A-048）、`F-I-006`（vision/workspace 投影过时）。复跑：四测试全 PASS 无 SKIP；判据 5 扫描与矩阵一致 | `03-audit/A-002-independent-r3d-root-closeout.md` |
| A-003 | self（编排器响应） | 2026-09-21 | 响应 A-002 全部 findings + 检查点 C 就绪声明 | **pass** | 6 条 recommended 全部闭合或登记（`F-I-001`～`F-I-005` fixed；`F-I-006` 登记为 C 动作）；检查点 B 判定完成；**Root 关门未发生**，三问确认包（`I-041-010` required + `I-041-011` + GOAL-004 overruled 确认）已提交用户 | `03-audit/A-003-response-to-a002-and-checkpoint-c-readiness.md` |
| A-004 | self（用户报告缺陷处置 + 用户裁决留痕） | 2026-09-21 | 用户报告 `.\dev.cmd start` 启动失败；`I-041-011` 与 GOAL-004 overruled 载体的用户裁决 | **pass** | 新 finding **`F-I-101`（required，用户要求先修）→ `fixed`**：`recoveryArtifactsDir` 相对路径致 `docker run -v` 失败（exit 125）→ 全分支经 `filepath.Abs` 绝对化 + 两条回归测试 + **一次性库上的端到端实测**（healthz/readyz 200、class-A 产物落盘 101,787 B）；`I-041-011` 按用户裁决只留退出矩阵/R3-C 附件；GOAL-004 `D-002` 按授权补写（结论不变）。Root 关门仍待用户确认 | `03-audit/A-004-self-dev-startup-finding-and-user-rulings.md` |
| A-005 | independent（grok build · grok-4.6 · high） | 2026-09-21 | 定向复审：`F-I-101`（required）是否合法闭合；同类相对路径隐患扫描；`L-1` 限定；用户裁决留痕；Root 关门就绪度 | **pass**（开放 required = **0**） | **同意 `F-I-101` 已 `fixed`**：独立复现旧函数相对挂载 exit 125、新函数挂载成功、**一次性库 `vp040_a005` 端到端**（healthz/readyz 200 + 产物 101,764 B）；13 处路径交给外部进程/容器的位置逐处判定，未发现第二处生产隐患；确认未触碰冻结物；**同意**再次把 Root 交给用户确认关门，**不同意**自行标 `done`。新 recommended：`F-I-102`（provider 不强制绝对 `WorkDir`；Abs 失败回退相对）、`F-I-103`（`00-meta` 的 `I-041-011` 未同步） | `03-audit/A-005-independent-fi101-closure.md` |
| A-006 | self（编排器响应） | 2026-09-21 | 响应 A-005 全部 findings + 重新提请 `I-041-010` | **pass** | `F-I-101` 维持 `fixed`；`F-I-102` **fixed**（`Create`/`Restore` 非绝对 `WorkDir` fail closed + `absoluteArtifactDir` 失败返回空 + 新增 3 例回归）；`F-I-103` **fixed**（`00-meta` 信息表同步 verified）；开放 required = 0，重新提请 Root 关门确认（并附 `.env` dev 库指向共享实例的环境建议） | `03-audit/A-006-response-to-a005-and-root-closeout-reask.md` |
| A-007 | self（编排器关门记录） | 2026-09-21 | 检查点 C：用户确认关门（判据 6）+ Root 关门记录 | **pass** | 用户 **2026-09-21** 书面确认关闭 Root（`01-decision/D-001`；Root `D-020`）：`GOAL-008` → `done · 3/3`、Root → `done · 3/3`、六条成功标准勾选、投影同步；跨目标开放 required = 0、未闭合 recommended = 0；dev 环境改指专用库 `schema_ui_dev`。**VP-040 波次关闭/Vision Review 属决策层，不在本关门范围** | 本条索引 + `02-execution/E-004` |

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
