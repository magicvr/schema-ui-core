---
id: A-039-r1-self-response-to-a038
doc_type: goal-audit-entry
source: self
auditor: /govern
date: 2026-09-20
scope: GOAL-002-r1-contract-and-denominator-freeze · response to A-038 / F-I-026 closed, F-I-002 three design items fixed, F-I-004 C3 boundary drafted
verdict: conditional
open_required: 3
status: recorded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-001-timestamptz-persistence-contract
version: 0.1.0
---

# A-039 · R1 self response to A-038

- **source**：self（编排器响应，**不**冒充 independent）
- **verdict**：conditional
- **开放 required**：**3**（F-I-002、F-I-004、F-I-005）——由 4 降至 3

## 1. 接受 A-038 的判定

| finding | A-038 verdict | 本响应动作 |
|---------|---------------|------------|
| **F-I-026**（PG 显式 DDL 缺席 v78） | **closed**（§4.1 v78 两列 NN 秒族骨架；未误入 F-5；v73–v87 的 **90 列全部被点名，无缺失**） | 接受闭合 |
| **F-I-002.1 第 1 项**（`dict_entries` 可粘贴 CREATE） | **fixed**（10 列、`badge_style` 末列、UNIQUE/FK 与 live cid 一致） | 接受 |
| **F-I-002.1 第 3 项**（`#72/#73` 单路径） | **fixed**（conversion contract 内无其它「双分支」USING） | 接受 |
| **三序更正**（cid / `sqlite_master` 文本 / Go 行号） | **正确、基本完整**；补：`Descriptors()` 的 Version 字段序与 cid 同序 | 接受并采纳该补充 |
| **A-034 接受的形式** | **保持**，未过度或不足展开 | 接受 |
| **F-I-002 整条** | **仍 open**：剩余为 ① 用例仍是 ID、非法/越界可执行测试未发生；② freeze-candidate ≠ 实施 | 见 §3 |
| 新 finding | **无**新 required / recommended | — |
| F-I-004 / F-I-005 | 维持 open | 见 §2 |
| F-I-025（recommended） | 收窄仍 open（附件「20 张时间列表」表述） | 见 §2.3 |

**A-038 的证据窗口处置正确**：它明确声明 `attachments/r1-c3-backup-recovery-boundary-v1.0-fc.md` 当时为**未跟踪文件、不在该 commit、不在其 scope**，不作为 F-I-004 证据。本响应照录——该文件在**本轮**才提交，**尚未被任何 independent 复审**。

## 2. 本轮动作

### 2.1 响应 A-038 §D 的补充

采纳「`Descriptors()` 的 `Version` 字段序与 live cid 同序」这一补充，作为列序交叉核对的**第三重**证据（cid / `sqlite_master` 折入文本 / `Descriptors()` 版本序三者同序；唯一会误导的是**模块文件物理行号序**）。

### 2.2 F-I-004 · C3 备份/回滚边界**首次落盘**

新增 `attachments/r1-c3-backup-recovery-boundary-v1.0-fc.md`（`status: freeze-candidate`），逐条对位 F-I-004 的五项关闭要求：

| # | 关闭要求 | 本文件落点 |
|--:|----------|-----------|
| 1 | 唯一区分 pre-conversion rollback artifact 与 post-conversion `CreateRecoveryPoint` artifact | §2 三类产物（A pre-conversion / B post-conversion / C batch boundary）与三条硬规则 |
| 2 | PG restore 不得再用预转换 `<artifact>` 当目标形状校验输入；给转换后 artifact 独立 token | §3 `<rollback-artifact>` vs `<recovery-artifact>` 双 token |
| 3 | `CreateRecoveryPoint` 的包路径与 before/after 调用点 | §4.1 包路径（kernel 仅接口 + `internal/backup` 编排）；§4.2 三个调用点（批次前生成 A/C、批次成功校验后调用、**禁止**在迁移事务内调用） |
| 4 | restore-to-new-db 可执行 harness | §5 harness 规格（两侧同构 + 6 项断言 + PG 版本兼容显式记录） |
| 5 | 旧 dump 不得通过新合同校验 | §6 `TestLegacyArtifactMustFail` 反向断言；C3 提交须同时具备 B 通过 + A/C 被拒 |

**现状事实（实测，已在 §1 列表）**：`apps/` 中 `CreateRecoveryPoint` / `BackupService` / `RecoveryPoint` **0 匹配**；`kernel.Store` 无 Backup 接口；`internal/temporal` 不存在；SQLite 现有 `snapshotBeforePending` 是 per-migration rollback 点。→ **C3 全部为设计，无任何实现**。

### 2.3 F-I-025（recommended）未闭合项

- 附件自称「20 张时间列表」的表述尚未改（A-036/A-038 两次点名）。下一轮随其他附件更新一并处理。

## 3. F-I-002 的结构性观察（**需用户注意**）

A-038 明确列出 F-I-002 剩余两项：

1. **「用例仍是 ID；非法/越界可执行测试未发生」**——这是**唯一**仍是 F-I-002 实质剩余的项；
2. 「freeze-candidate ≠ 实施」——属边界声明，不是可闭合项。

**观察**：第 1 项要求**可执行测试**。若在**一次性内存库**上按合同 DDL 建表并跑负值/越界/边界用例，既可满足「可执行测试」又不触碰生产 schema（不违反 C2 未冻结前的「不得实施 schema 变更」门禁）；若必须在真实迁移代码上跑，则属 R2 落码，**在 C2 冻结前无法闭合 F-I-002**——这会构成结构性死锁（F-I-002 挡 C2 冻结，C2 冻结又挡跑测试）。

→ **下一轮建议就此向用户确认口径**（P-004）：F-I-002 的「可执行测试」是允许在**一次性验证库**上跑（可在 R1 内闭合），还是必须绑定真实迁移（则 F-I-002 只能随 R2 闭合，需调整 R1 关门口径）。

## 4. 仍开放（不得放行）

**F-I-002 / F-I-004 / F-I-005**；C2/C3 未冻结，R2 未放行。**本响应未闭合任何 required**（F-I-026 的 closed 由 A-038 判定）。

## 5. 下一步

1. `/audit` 复审**本轮新落盘的 C3 边界文件**（F-I-004）——该文件从未被复审。
2. 就 §3 的 F-I-002 测试口径向用户取 P-004 裁决。
3. F-I-005 需 R2 落码时才能记录真实哈希，属结构性依赖。
