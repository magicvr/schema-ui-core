---
id: A-001-self-r3c-matrix
doc: audit-entry
status: active
parent: GOAL-007-r3-pg-cross-version-restore-matrix
created: 2026-09-21
updated: 2026-09-21
version: 0.1.0
---

# A-001 · self 审计：R3-C 检查点 A/B（矩阵定义、实测记录、口径更正）

- **source**: `self`
- **日期**: 2026-09-21
- **scope**: `GOAL-007` 检查点 A（`D-001` 矩阵定义与判定口径）与 B（`E-003` 实测记录、`D-002` 规则更正、`I-041-004` 收口）；边界依据 Root `D-018` §2 第 4 项/§5、`D-017` §3 约束②；不含 R3-D。
- **verdict**: `conditional`（**开放 required = 0**；2 条 recommended 提交 independent 复核，1 条已在本轮内部修正）

## 成果（逐项可核对）

| # | 判据 | 结论 | 证据 |
|--:|------|------|------|
| 1 | 判定口径与组合边界**先于**矩阵本体冻结 | **达成** | `D-001`（2026-09-21，先于 `E-003`）；`I-041-009` → verified |
| 2 | 组合枚举覆盖 server×client 关键组合 | **达成** | dump 9 格（6 supported / 3 toolgate）；restore 54 格（6 归档 × 3 client × 3 target） |
| 3 | unsupported **逐格**记录原因与退出码 | **达成** | `attachments/r3c-pg-cross-version-matrix-v0.1.md` 两张表逐行含 exit 与首条诊断 |
| 4 | 升级后恢复有界核对（跨版本） | **达成** | 18 个 supported 格含 15→16、15→17、16→17 升级恢复，逐格执行 `D-001` §4 四项校验 |
| 5 | 形状校验足以区分「恢复成功」与「恢复成空/降级」 | **达成（自审已构造反证）** | 空库会使 `schema_migrations` 查询失败 → `matrixReadFacts` 直接 `t.Fatalf`；微秒尾零逐字节断言可捕捉「降为 3 位」 |
| 6 | 驱动 fail closed | **达成** | `matrixClassify` 的 `unexpected-failure` / `setup-failure` 分支直接 `t.Fatalf`；本次运行该计数为 0 |
| 7 | 默认基线不依赖 docker | **达成** | 未设 `VP040_PG_MATRIX` 时 `--- SKIP`；本仓 `go test ./...` 仍不依赖 docker |
| 8 | 破坏性动作限定一次性/专用测试 database | **达成** | 容器内新建库 + 宿主 `t.TempDir()`；常驻 15.4 未参与（`E-003` 清理节） |
| 9 | 证据定位未被越读为生产就绪 | **达成** | `D-002` §4 与附件「边界」节明示 CI/reproducibility；不作产品级支持承诺 |
| 10 | `I-041-004` 收口 | **达成** | `D-002` §3：以逐格记录满足信息需求 → **verified**（非 residual） |
| 11 | 附带结论（迁移链在 16/17 至 v87 成立）有可核对证据 | **达成** | 三个源库分别在 15.19/16.15/17.11 上用真实 catalog 建到 head 87；`t.Logf` 记录 head 与行数 |
| 12 | v1–v87 canonical SQL/checksum 不可变 | **未触碰** | 本轮只新增 `internal/backup` 的一个测试文件；无 migration/descriptor 改动 |

## 偏差（本轮内部发现并已修正）

| ID | 级别 | 偏差 | 处理 |
|----|------|------|------|
| `F-S-001` | recommended | 首版驱动草案存在缺陷（`matrixAdminCache` 从未填充、DSN 以字符串替换推导容器地址、残留调试代码、以字符串匹配报告统计），在写入磁盘前整体重写为结构化结果 + 显式 DSN 构造 | **fixed**（`git` 仅记录重写后的版本） |
| `F-S-002` | recommended | `D-001` §3 把探测观察**外推**为 `dumper ≤ client ≤ target_server`；矩阵实测给出反例（16 client → 15 server 为 supported） | **fixed**：新增 `D-002` 更正规则为 `client ≥ server`、`client ≥ dumper_client`、`client == 17 ⇒ target == 17`；`D-001` 保持原文以便追溯；分类与枚举未被推翻 |

## 提交 independent 复核的事项

| ID | 级别 | 事项 | 独立复核点 |
|----|------|------|-----------|
| `F-S-003` | recommended | `matrixClassify` 的三类 + `unexpected-failure` 是否会把未知失败误判为 unsupported（文本匹配的鲁棒性） | 构造未知错误文本（例如网络中断、磁盘满、权限拒绝）确认落入 `unexpected-failure`；确认 `setup-failure` 不会被读成 restore 结果 |
| `F-S-004` | recommended | 形状校验四项是否可能被「部分恢复」骗过（例如只恢复了 ledger 未恢复数据行） | 反例优先：制造缺表/缺行的恢复结果，确认四项中至少一项失败 |
| `F-S-005` | note | `D-001` §2 写「27 格」而实测为 54 格（6 归档 × 3 × 3）；`D-002` §2 已注明 | 是否需要把 `D-001` §2 的枚举式改为「归档数 × 3 × 3」的表述（当前以 `D-002` 交叉引用处理） |
| `F-S-006` | note | 驱动只在容器 server 上验证，未把常驻 15.4 作为 server 轴（`D-001` §7 未选方案）；其应用层路径由既有 C3 库内测试覆盖 | 该取舍是否导致判据 4 证据缺口 |

## 已核对为「非问题」的项（供审计反证）

| ID | 项 | 依据 |
|----|----|------|
| `N-001` | 归档由哪个 server 生成是否构成额外目标约束 | 实测：`src16_by_16` → 16 client → **15 server** 为 supported 且形状通过，故无额外约束（`D-002` §1 第 5 条） |
| `N-002` | `pg_restore -l` 可读性是否等价于可恢复 | 两者在本次矩阵中一致（toolgate 格在 `-l` 阶段即失败），但矩阵只以实际恢复为判据，未把 `-l` 当证据 |
| `N-003` | 镜像 tag 可变导致结论漂移 | `D-001` §1 要求每次重新记录版本串；附件记录了本次 15.19/16.15/17.11 |
| `N-004` | `--exit-on-error` 会放大差异（不带该旗标时 `transaction_timeout` 只是 per-statement 失败） | `D-001` §3 已冻结「以 C3 实际使用的旗标为准」；C3 生产路径 `provider_pg.go` 使用 `--exit-on-error` |
| `N-005` | 矩阵是否修改了常驻实例 | 未使用常驻 15.4；容器与宿主临时目录在 `t.Cleanup`/`t.TempDir()` 回收 |

## 自审可复跑证据

- `cd apps/api && go build ./...` → exit 0；`go vet ./internal/backup/` → exit 0。
- 门控矩阵：`VP040_PG_MATRIX=1 go test -count=1 -run TestPGCrossVersionRestoreMatrix -v -timeout 45m ./internal/backup/` → **PASS（105.30s）**，`unexpected-failure` = 0，形状校验 18 次全通过。
- 默认隔离：不设门控同命令 → **SKIP**。
- 全仓 `go test -count=1 ./...` 见 `02-execution.md`。

## 自审结论

检查点 A/B 的判据均有可执行证据，开放 required = **0**；两条偏差已修正（其中 `F-S-002` 是**规则陈述**的更正，不影响判定分类）；`F-S-003`/`F-S-004` 提交 grok independent 复核。检查点 C 的成立以独立意见与其响应为准。
