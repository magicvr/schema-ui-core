# Grok Build · 独立交叉审计提示词（GOAL-002 R1 冻结）

> 执行方式：在仓库根 `schema-ui-core` 以 grok build（模型 grok-4.6 · 思考强度 high）调用 `/audit` 技能（`.grok/skills/audit/SKILL.md`），出具 `source: independent` 的**R1 冻结复审意见**。本文件由编排器预置，供独立会话只读核验。

## 任务

- **目标**：`GOAL-002-r1-denominator-and-contract-freeze`
- **工作区**：`workspace-038-batch-operations-and-job-center`（canonical `docs/workspaces/workspace-038-batch-operations-and-job-center/`；root_goal `GOAL-001-batch-operations-and-job-center`）
- **scope**：R1 冻结（C1～C3）全量复审——① 两项冻结矩阵是否**可机器核对**且与代码事实一致；② `I-038-001`～`003` 的 `verified` 关闭是否合法（有证据 + 用户裁决，而非把假设写成已验证）；③ 用户 P-004 裁决（C2=方案 B / C3=首波仅新建批量导出所选 / C1=管理作用域 + 新 `jobs.read`）是否被**如实**落盘，未选方案是否记录；④ 派生结论与用户裁决是否被诚实区分（D-001 §3.3）；⑤ 未定项 O-1～O-3 是否确属未定、未冒充已裁决；⑥ 边界：是否越界改动 `apps/**` 或任何 pinned 协议工件。
- **只读约束**：不修改 `00-meta` status/progress、不改 goal-tree、不改方案正文；只可追加审计意见（`03-audit/A-NNN` + 索引）或输出文本供代贴。

## 必读文件（按序）

1. `docs/workspaces/workspace-038-batch-operations-and-job-center/workspace.md` 与 `goal-tree.md`
2. `GOAL-002-r1-denominator-and-contract-freeze/00-meta.md`（4 检查点 / 审计模式 cross / 信息就绪）
3. `GOAL-002-…/01-decision.md` 与 `01-decision/D-001-r1-contract-and-denominator-freeze.md`（核心：§1 C2、§2 C1、§3 C3、§4 未选方案、§5 移交清单）
4. `GOAL-002-…/02-execution/E-001-r1-establishment-and-recon.md`、`E-002-c1-c3-freeze.md`
5. `GOAL-002-…/attachments/r1-job-kind-scope-matrix.md`、`r1-first-wave-denominator-matrix.md`
6. `GOAL-002-…/03-audit.md` 与 `03-audit/A-001-r1-freeze-self.md`
7. 证据基线：`GOAL-001-batch-operations-and-job-center/02-execution/E-002-r1-recon.md`；`GOAL-001-…/attachments/R1-recon-I-038-00{1,2,3}-*.md`
8. `docs/vision/plans/VP-038-batch-operations-and-job-center.md`（判据 1～7、非目标、首波范围）

## 重点核验清单（请自行读代码复验，不要只信文档）

- **矩阵 vs 代码**：`r1-job-kind-scope-matrix.md` 的每一行是否与 `apps/api/internal/jobs/**`、`apps/api/modules/wallet/jobs.go`、`apps/api/internal/handler/wallet.go`、`apps/api/modules/jobs/migration/migration.go` 一致。特别核验：
  - Job 生产种类分母是否真的 = 1（`grep -rn "Register\|RegisterWithTerminalHook" apps/api`）
  - `GetForActor` 是否确为 `id + kind + actor_id` 三重限定；是否真的**没有**任何跨 actor 读路径
  - 三个索引定义与「跨 actor `ORDER BY updated_at DESC` 不被覆盖」的结论是否成立
  - `idx_jobs_*` 在 sqlite 与 postgres 两侧是否一致
- **运行时门控**：`composition.go` 的 `jobRuntime.enabled` 是否确实只在 `plan.HasModule("admin.wallet")` 下置 true；`Start()`/`Stop()` 的 no-op 守卫是否存在。矩阵 §6 R-1 的结论是否成立。
- **冻结断言**：`internal/store/migrate_test.go` 的 `core.jobs/async_jobs` checksum 与 `modules/jobs/migration/migration_test.go` 的 `Version == 42` 是否如文档所述被冻结。
- **首波分母口径**：生产页面 schema 是否真的**零**批量动作（`batchMapping` / `requiresSelection` / `$selection.keys`）。分母修正为「4 个资源 / 6 条路由」是否准确（数一数 `resources.go` 的挂载条件与各 provider 声明的路由）。
- **同步 batch-delete 形状**：`resources.go` 的路由注册行、权限门、4 KiB 上限、`200 {"deleted": n}`、`BatchDeleter` 实现者（是否真的只有 users/roles 原子）是否与矩阵一致。
- **Breaking 标记**：users/roles `batch-delete` 与 `scheduled-tasks POST /{id}/run` 改异步是否确为 BREAKING（查前端依赖与测试断言）。
- **排除清单准确性**：X-1～X-11 的位置与「当前同步/后台」状态是否属实；是否有**遗漏**的长操作或批量操作未进清单。
- **协议 pin 零影响**：D-001 §1.1 的「零改动 pinned 工件」是否成立（方案 B 下是否需要碰 `docs/schemas/**` 或 `apps/web/src/protocol/upstream/**`）。
- **诚实性**：D-001 §3.3 把「同步 batch-delete 保持」标为派生结论而非用户裁决——该标注是否恰当（对照本轮 P-004 提问的实际范围）。
- **边界**：`git log`/`git show --stat` 复核两次 checkpoint 是否只含 `docs/workspaces/workspace-038-…/` 路径，`apps/**` 零改动。

## 输出要求

按 `.grok/skills/audit/SKILL.md` 与 `skills/prompts/05-independent-audit.md` 结构：verdict（pass | conditional | fail，附尺度）、Findings（F-00N；required | recommended；严重度；evidence 路径）、必改项汇总、与 self A-001 的异同、结论与给编排器/用户的下一步。若可写入：追加 `03-audit/A-002-*.md`（source: independent）并更新 `03-audit.md` 索引（不改 status/progress/goal-tree）。
