---
id: D-001
doc: decision-entry
goal: GOAL-004-r3-recovery-evidence
status: accepted
created: 2026-08-22
updated: 2026-08-22
version: 1.1.0
---

# D-001 · R3 轮换后恢复最小剧本（关闭 I-004）与证据方案

> v1.1（2026-08-22）：按 Root A-002 F-003 勘误 PG 版本组合措辞；判定标准不变。

## I-004 裁定（verified）

### 剧本时序（备份点相对轮换点）

```text
T0  K1 单密钥运行；产生真实会话（login → access_old[K1] + refresh_old）
T1  备份点：对运行中的库取备份（两方言各自命令，见下）
T2  轮换配置：current=K2, previous=K1（进程重启生效——本波无热加载）
T3  恢复点：把 T1 备份恢复为可用数据库（SQLite=备份文件即库；PG=新库 pg_restore）
T4  以轮换后配置从恢复库启动应用 → 全模块 Start+Ready 门禁须绿
T5  鉴权断言（在恢复库 + 轮换后密钥集上）：
    A1 旧 access_old[K1] 经中间件验签通过（重叠窗内 previous 可验）
    A2 login 新签发 access 仅能被 K2 验证、不能被 K1 验证（签发只用 current）
    A3 refresh_old 在恢复库上仍可 Refresh 成功并签出 K2 新对（opaque 会话跨轮换+恢复连续）
```

### 两方言证据命令（不重做 dump，全部既有工具）

| 方言 | 备份（T1） | 恢复（T3） |
|------|-----------|-----------|
| SQLite | `VACUUM INTO '<backup>'`（VP-013 sqlite 合同；测试内经 `database/sql` 执行） | 备份文件本身即合法库，直接作为 `db.path` 启动 |
| PG | `pg_dump -F c -U <user> <db> -f <dump>`（官方客户端；允许 [workspace-013] GOAL-006 D-002 已记录的跨版本客户端 GUC 告警类——restore 出现 1 条 ignored error 属预期，功能不受影响） | `createdb <rest>` + `pg_restore -U <user> -d <rest> <dump>` |

> v1.1 勘误（Root A-002 F-003）：初版写「17 客户端 ↔ 17 服务端同主版本无 GUC 告警」，与实跑环境（17 客户端 → 15.4 服务端，1 条 ignored GUC error）不一致。合同以「官方工具 + ledger 指纹核对」为准，不依赖客户端/服务端同主版本。

### 断言判据

- T4 组合启动 = `composition.NewApp(...).Start(ctx)` 无错（等价 readyz 门禁路径全绿）。
- A1/A2/A3 全部成立才记「轮换后恢复合同成立」。

## 证据方案

1. **可重复自动化**（committed 测试，`internal/composition/post_rotation_recovery_test.go`）：
   - `TestSQLitePostRotationRecovery`：完整 T0–T5 循环，**无条件必跑**（文件库，无外部依赖）。PG 侧同构循环 `TestPostgresPostRotationRecovery` 按 `pgtest.DSN()` 门控；dump/restore 工具按「PATH 上有 `pg_dump`/`pg_restore`，或环境变量 `R16_PG_DUMP_CONTAINER` 指向带客户端工具的容器」门控，二者皆无则跳过（不伪造证据）。
   - 密钥值使用强随机形字符串（K1≠K2）；access TTL 15m 内完成全部断言。
2. **实跑记录**：E 条目登记双方言命令与输出摘要（含 PG round-trip 的 restore 库 boot 绿与三条鉴权断言结果）。

## 为什么

- 密钥不在库中 ⇒ 恢复风险不在"数据丢密钥"，而在"旧会话数据 × 新密钥配置"组合语义；剧本直击该组合。
- A3 是关键增强断言：refresh 为 SHA-256 opaque 存库，跨轮换与恢复都应连续——这是"备份可用性 × 轮换安全性"同时成立的直接证据。
- 自动化放在 composition 包：复用 `NewApp` 真实启动门禁与 `pgtest` 门控惯例，避免另起测试装配。

## 未选方案

- 在 Go 测试里用 `CREATE DATABASE ... TEMPLATE` 替代 pg_dump/restore：那是另一条复制路径，违反「不重做 dump」边界，不冒充合同证据。
- 只做脚本不做 committed 测试：workspace-015 先例表明可重复自动化是回归保障；脚本仅作 live 补充。
- 断言里包含"旧 access 立即失效"：I-005 non-blocking 默认已冻结为 previous 可验，不测残余分支。
