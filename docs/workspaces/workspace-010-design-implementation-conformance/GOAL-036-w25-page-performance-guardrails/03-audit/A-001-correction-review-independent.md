---
title: A-001 · W25 修正情况与修改方式独立复审（independent）
source: independent
status: recorded
created: 2026-08-23
updated: 2026-08-23
parent: GOAL-036-w25-page-performance-guardrails
version: 0.1.0
auditor: grok-4.6（Grok Build /audit）
scope: 复审 GOAL-036 已实施修正（S1–S5 性能修复与防复发、I-001 e2e/孤儿清理、I-002 活栈测量）及修改方式是否合理；execution-facts + ad-hoc；非关门审计
verdict: conditional
---

# A-001 · W25 修正情况与修改方式独立复审（2026-08-23，independent）

## 范围与区间

用户指令：`/audit 审计工作区10目标36的修正情况，并审视其修改方式是否合理？`

工作区：`workspace-010-design-implementation-conformance`（Root `GOAL-001-design-implementation-conformance`；canonical 范围匹配；`shared_materials_catalog: none`）。

复审对象：GOAL-036（W25）截至 2026-08-23 的已落地修正，不是关门放行。

- 性能四因素修复（D-001 / E-001）与全盘/防复发（D-002 / E-002 / E-003）
- I-001 关闭声明（e2e 双 profile + 删用户孤儿链接，E-004）
- I-002 关闭声明（活栈计时，E-005）
- **修改方式**：池化/WAL、in-flight 合并、探活后写、schema 缓存、定向刷新、e2e 暴露缺陷的就地修补，是否与声明的不变量一致

本轮**未**复跑 Playwright e2e，**未**重做双栈计时；I-001/I-002 的「已跑过」主张按落盘证据 + 本轮可读代码/单测核对。SQLite 连接面不变量由本轮现场探针复验（探针文件用后删除，未入库）。

## 成果（有证据）

| 项 | 证据路径 | 核验方式 |
|----|----------|----------|
| 提交落地：产品改动在 `ba7d5c6`（父提交 `0878d7f` = W24 收盘） | `git show --stat ba7d5c6`（33 文件）；工作树相对该提交仅 I-002 文档未提交 | 本复审 `git log` / `git status` |
| Fix A：文件库池 4 + DSN `_busy_timeout=5000&_journal_mode=WAL&_synchronous=NORMAL`；`:memory:` 保单连接 | `apps/api/internal/store/store.go`；`store_wal_test.go` | 读码；本轮 `go test ./internal/store/ -run TestSQLiteDSNPragmas\|TestFileStoreWALPoolAndPragma\|TestMemoryStoreStaysSingleConnection\|TestFileStorePoolOverride` **ok** |
| postgres 路径未改 | `git show ba7d5c6 -- apps/api/internal/store/postgres.go` 空 | 提交足迹 |
| Fix C：`fetchList` 仅 in-flight 合并，无长期结果 memo；`reloadList` / transport 变更清空 in-flight | `apps/web/src/renderer/render.tsx`；`render.test.tsx` 三 statCard → 1 请求、statCard+chart → 1 请求 | 读码 |
| Fix B：wallet-ensure 探活后写；存在则零 POST、零 `reloadList` | `wallet-ensure.tsx` + `wallet-ensure.test.tsx` 4 例 | 读码 |
| Fix D：schema 文档 opt-in 缓存，失败不缓存，容量 64 | `load-page.ts` / `App.tsx` | 读码 |
| S5：`refreshList` + monitoring tick 只刷 `/api/system-monitoring/status`；注册校验镜像 `main.tsx` 9 个副作用导入 | `monitoring-auto-refresh.tsx`；`custom-components.schema.test.ts`；`main.tsx` L14–22 | 读码对照 |
| Playbook §6（1.1.0） | `docs/architecture/module-contribution-playbook.md` | 读码 |
| I-001 显式清理 `user_roles` / `user_mfa`（单删+批删）+ 单测 | `users_repository.go`；`TestDeleteUserCleansRoleAndMfaLinks` 本轮 **ok**；`TestDeleteUsersBatchCleansRoleAndMfaLinks` 本轮 **ok** | 读码 + `go test ./internal/modules/authsession/` |
| I-001 关闭材料：admin 首轮 8/9 失败点、探针状态机、修复后 9/9×2 叙述 | E-004 + `attachments/I-001-evidence.md` | 读材料（未复跑 e2e） |
| I-002 关闭材料：基线 `0878d7f` vs 当前，请求数 −47%～−86%，RTT150ms −1.4s；原始 JSON 五行+五行 | E-005 + `attachments/I-002-evidence.md`；`git rev-parse HEAD^` = `0878d7f` | 提交关系核对 + 算术核对（未重测） |

## 对照成功标准

| 标准 | 核验 | 结论 |
|------|------|------|
| C1 诊断 | D-001 四因素与代码落点一致（重复 GET / 挂载即写 / schema 重取 / 单连接+fsync） | **达成** |
| C2 钱包页实施 | 后端池/WAL + 前端合并/探活/schema 缓存均在 `ba7d5c6` | **达成，但 C2 后端改法引入 F-001** |
| C3 全盘扫描 | D-002 26 页台账；同类重复请求由全局合并覆盖，无需逐页改——机制方向正确 | **达成（扫描本身）** |
| C4 防复发 | 三支柱有文件；**漏钉 `foreign_keys` per-conn**（F-002） | **部分** |
| C5 自动回归 | E-003 自称 vitest 1096 / tsc 0 / store 绿；本轮复跑 store WAL 组 + authsession 两例删除测试均 ok。全量 vitest/e2e 未复跑 | **部分（定向复跑绿；全量未复验）** |
| C6 验证 | I-001 作为「e2e 是否全绿」信息项：材料完整、单测绿，本轮未复跑浏览器。I-002 作为 non-blocking 计时：JSON 自洽、基线提交属实。自审 A 条目此前为空（本条为独立审，不替代 self） | **I-001/I-002 关闭材料可核对；不构成关门** |

## 修改方式总评

**前端三条（合并 / 探活后写 / schema 缓存）方向合理，取舍留痕清楚。** 拒绝长期 URL 结果缓存以避免 T-02 reset 命中过期数据、`fetchList` 显式传入 transport、schema 缓存 opt-in 且失败不写，都与既有语义对齐。`DISPLAY_LIST_QUERY` 让探活与 statCard 共用同一 URL，是对症合并而不是改 RPC 形态。monitoring tick 改为 `refreshList("/status")`、事件表不随轮询刷新，有 D-002 书面范围。大表 `COUNT(*)` 出局正交，成立。

**后端「文件库小池 + WAL + NORMAL + busy_timeout」方向也对，但落地不完整——这是本波修改方式的主要缺陷。** SQLite 的 `PRAGMA foreign_keys` 是**连接级**状态。历史 `MaxOpenConns=1` 时，`migrate.go` `assertForeignKeysOn()` 在唯一连接上执行一次即覆盖全库。改为池 4 后仍只在 `open()`→`migrate()` 时对**第一条**连接 `Exec PRAGMA foreign_keys=ON`；DSN 三元组（busy/WAL/synchronous）进了每条连接，**`foreign_keys` 没有**。

本轮现场探针（`database/sql` 同时 `Conn()` 持有 4 条连接后读 `PRAGMA foreign_keys`，随后在最后一条连接上 `INSERT user + user_roles` → `DELETE users`，**不**走 I-001 的显式 `DELETE FROM user_roles`）：

```
held conn 0 foreign_keys=1
held conn 1 foreign_keys=0
held conn 2 foreign_keys=0
held conn 3 foreign_keys=0
user_roles leftover after DELETE users = 1   // CASCADE 未开火
```

DDL 写明 `user_roles.user_id REFERENCES users(id) ON DELETE CASCADE`（`authsession/migration/migration.go`）。W24 收盘 sqlite e2e 9/9（含 schema-crud 删用户后再删角色）；W25 把池改为 4 后同一用例在 admin 稳定失败。时间线与探针一致：**I-001 的「孤儿」首先是本波连接面改动关掉了 CASCADE 安全网的症状**，其次才是 `DeleteUser` SQL 从未显式清链接（字面为真，作为「W25 新缺陷根因」不完整）。

显式补删 `user_roles`/`user_mfa` **仍值得保留**（防御纵深；`user_mfa` 的 FK **没有** `ON DELETE CASCADE`，FK 开启时本就会挡住删用户）。但它不能代替把 `foreign_keys` 放进 DSN。postgres 本来就强制 FK，所以「postgres 零改动」成立，也解释了为何生产方言没暴露此回归。

## Findings

### F-001 · SQLite 池化未把 `foreign_keys` 迁入每条连接（CASCADE 在 3/4 连接上失效）

- 严重度：**high**
- 建议：**required**
- 关联：C2 后端改法；I-001 关闭叙事；P-005 I-001 作为 e2e 门禁仍可视为「测绿」，但**完整性回归未关**
- 描述：`sqliteDSNParams` 含 busy/WAL/synchronous，不含 foreign_keys。`assertForeignKeysOn` 仍是对 `s.db` 的一次性 `Exec`+`QueryRow`，池化后只保证 migrate 用过的那条连接。本轮探针：4 连接中 3 条 `foreign_keys=0`，`ON DELETE CASCADE` 不执行，删用户留下 `user_roles`。同一机制还会让其它依赖 SQLite FK 的行为（`notifications` 的 `ON DELETE CASCADE`、RBAC RESTRICT、`TestForeignKeyEnabled` 所覆盖的 `refresh_tokens` 插入校验）在**非 migrate 连接**上失效。HTTP 活栈必然打开多连接，与 e2e 首轮 admin 失败吻合。
- 证据：`store.go` `sqliteDSNParams`；`migrate.go:31-34,253-264`；本轮探针日志（上节）；E-004「W24 后同用例在 W25 首轮红」；`user_roles` DDL CASCADE。
- 状态：open
- 处置：将 `_foreign_keys=on`（或与现有 mattn-compat 下划线参数等价的写法）并入 `sqliteDSNParams`，使**每条**池连接开启 FK；保留 I-001 显式清理作为纵深。补回归见 F-002。

### F-002 · S5 连接面防复发测试未钉死 per-conn `foreign_keys`，既有 FK 测试在池化后假绿

- 严重度：**med**
- 建议：**required**
- 描述：`store_wal_test.go` 断言池大小、`journal_mode=wal`、`busy_timeout=5000`、`synchronous=1`、内存库单连接——**不**断言每条连接 `foreign_keys=ON`，也不在多连接持有下测 CASCADE。`TestForeignKeyEnabled` / `TestRBACConstraintsAndIndexes` 对 `st.db` 顺序 `Exec`，复用 migrate 那条已开 FK 的连接；本轮两者仍 **ok**，因此不能发现 F-001。S5「回退 MaxOpenConns=1 立即红」成立，但「池化不得破坏既有 FK 不变量」没有栅栏。
- 证据：`store_wal_test.go`；`migrate_test.go` `TestForeignKeyEnabled`（单次 `st.db.Exec`）；本轮 `go test` 该组通过 vs 探针 3/4 连接 FK=0。
- 状态：open
- 处置：在 store 测试中同时持有 `sqlitePoolDefault` 条 `Conn`，断言每条 `PRAGMA foreign_keys=1`，并在非第一条连接上复现「删用户 → `user_roles` 为 0」（CASCADE）以及「插入缺用户的 `refresh_tokens` 失败」。

### F-003 · `refreshList` 不丢 in-flight，与 `reloadList` 不对称

- 严重度：**med**
- 建议：**recommended**
- 描述：`reloadList` 有意清空 `listInFlight`，避免慢请求期间的重载加入「变更前」那一发。`refreshList` 只 bump per-URL token，`useDisplayData` 重跑 effect 时 `fetchList` 仍会 join 仍在飞的同一 URL Promise——重叠的定向刷新拿不到新数据。monitoring 最短 5s，通常不会踩中；作为「修改方式」与 `reloadList` 的注释契约不一致，慢网络/连点时漏刷。
- 证据：`render.tsx` `reloadList` vs `refreshList`/`fetchList`；`render.test.tsx` 只覆盖「刷新后目标 URL +1、其它不变」，不覆盖 in-flight 重叠。
- 状态：open

### F-004 · I-002 测量不可从仓库复现；台账卫生滞后

- 严重度：**low**
- 建议：**recommended**
- 描述：E-005 写明测量脚本在 `%TEMP%\w25-metrics\`、不入库。附件有 JSON 五行，基线 commit 属实、算术自洽，作为 non-blocking 关闭材料**够用**，但他人无法复跑。另：`attachments/README.md` 仍写「当前为空」；`00-meta` 路线图 S6 仍把 I-002 写成未勾，与信息表 `closed` 不一致。I-002 文档本身尚未进入 `ba7d5c6`（工作树 `??`）。
- 证据：E-005「会话临时目录」；`attachments/README.md`；`00-meta.md` 路线图 L46 vs 信息表 I-002；`git status`。
- 状态：open

### F-005 · 批删回归测试名义含 MFA、断言只有 `user_roles`

- 严重度：**low**
- 建议：**recommended**
- 描述：`TestDeleteUserCleansRoleAndMfaLinks` 播种并断言 `user_mfa=0`。`TestDeleteUsersBatchCleansRoleAndMfaLinks` 不插入 `user_mfa`、不计数 MFA，批路径 MFA 清理无回归钉。
- 证据：`users_repository_test.go` 两例对照。
- 状态：open

## 必改项汇总

1. **F-001（high required）**：`sqliteDSNParams`（或等价 per-connection 钩子）必须对**每一条**文件库连接打开 `foreign_keys`；用多连接探针证明 CASCADE 恢复。保留 I-001 显式删链接。
2. **F-002（med required）**：把 per-conn `foreign_keys` + CASCADE/RESTRICT 行为写入 store 防复发测试，避免再出现「WAL 测试绿、活栈 FK 灭」。

无其余 required。F-003～F-005 为 recommended，不单独阻断，但 F-001 未修前**不得**把 C2 后端改法视为已安全，也**不得**将 I-001 理解为「sqlite 引用完整性已恢复」。

## 信息项（P-005）

| 项 | 级别 | 最晚阶段 | 本轮判定 |
|----|------|----------|----------|
| I-001 e2e 双 profile × sqlite 是否全绿 | required · C6 | 材料称修复后 9/9×2；本轮未复跑 Playwright。作为「测绿」可维持 closed，**前提是编排器承认 F-001 是并行完整性缺口，不是 I-001 已覆盖的根因** | closed 维持，但叙事需修正 |
| I-002 活栈可感提升 | non-blocking · C6 | JSON 与基线 commit 可核对；脚本未入库（F-004） | closed 维持 |

无到期未关的 required 信息项被本意见直接 reopen。F-001 是新的实施正确性 finding，不是新的 I-00N。

## 与既有意见的异同

此前目标 `03-audit` 无 A 条目（索引曾写「S6 前不新增 A」）。无 self/independent 历史可对照。与执行台账的差异：E-004 将 admin e2e 失败归因于「DeleteUser 从不清理 user_roles」；本意见同意 SQL 字面，并补上 **W25 池化关闭 CASCADE** 这一本波引入的机制原因。

## 结论 + 建议给编排器/用户的下一步

**verdict: conditional。**

性能修正（请求合并、探活后写、schema 缓存、WAL/池化方向、monitoring 定向刷新、Playbook）**真实存在且前端取舍合理**；I-002 显示的可感提升材料内部一致。I-001 把 e2e 重新打绿、并补了本应存在的显式清理，作为 e2e 门禁关闭**可以成立**。

**修改方式不合理的核心**是：把 SQLite 从单连接改成池时，没有把连接级 `foreign_keys` 随 DSN 带到每一条连接。这让本波自己的 e2e 红了一发，也让既有 FK 测试假绿。未修 F-001/F-002 前，不建议 self 关门，也不建议把「页面性能问题以后不会再出现」理解成「sqlite 完整性面无回归」。

建议 `/govern` 响应：

1. 先修 F-001（DSN 级 `foreign_keys`）并落地 F-002 多连接回归；再决定是否复跑一次 sqlite e2e 作为 F-001 关闭证据。
2. F-003～F-005 顺手修或书面 residual（I-002 脚本不入库可接受，但请改 README / 路线图勾选以免台账自相矛盾）。
3. 本意见不替代 S6 自审；F-001 闭合后再跑 self / 考虑关门。
4. 目标保持 `active` 5/6，由编排器处理 finding，审计不改 status。

## 声明

本意见不修改 status/progress/方案正文/goal-tree 状态列；响应与任何状态变更由 `/govern` 处理。
