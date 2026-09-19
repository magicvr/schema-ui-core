---
id: A-001-r5-closeout-self
doc: audit-entry
parent: GOAL-006-r5-evidence-and-closeout
status: recorded
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# A-001 · R5 关门自审（GOAL-006 C1～C2）

## A-001 · R5 C1～C2 自审（2026-09-19）

- **source**：self
- **auditor**：`/govern` 编排器（DeepSeek Harness 会话）
- **类型** / **scope**：`closeout` · GOAL-006 C1～C2（退出矩阵 / 浏览器回归）
- **verdict**：**pass**（0 required；2 recommended）
- **审计模式**：R5 为**关门审计**，按 `00-meta` §说明取 **`cross`**（self 本条 + independent 由 grok build 执行，见 C3）。

### 范围与区间

- 被审对象：`VP-038` §方向级退出判据 1～7 的证据充分性 + e2e 回归可信度。
- 审计区间：R1 `GOAL-002` → R5 本条目（HEAD `未定`；checkpoint 由 C4 登记）。
- 材料：`GOAL-002`～`GOAL-005` 的五件套；`GOAL-006/02-execution/E-001-r5-exit-matrix.md`；`apps/web/e2e/**`；`[workspace-010] GOAL-044`（跨区残余交付）。

### 成果（有证据）

| # | 主张 | 复核 |
|---|------|------|
| 1 | 判据 1～6 均有**可核对**证据（矩阵 / 测试 / 空 diff / 用户裁决），未用叙事代替 | `E-001` §1 逐条指向具体文件与测试名；本条目复读 `GOAL-002 D-001`、`GOAL-003 A-002`、`GOAL-004 A-002`、`GOAL-005 A-003` 与代码位置，无「声称但无证据」的条目 |
| 2 | 判据 7 的两项前置（退出矩阵、浏览器回归）已完成；独立意见与用户确认属 C3/C4，未提前宣称 | `E-001` §1 判据 7 标注「进行中」 |
| 3 | e2e 首次运行的失败被**定位到既有挂具缺陷**而非实现回退，且有两向证据（失败快照 + 隔离复验） | `E-001` §2 |
| 4 | 修复只让既有假设成真（重命名），未改断言、未放宽覆盖 | `git mv` + 注释；重跑 **16 passed / 4 skipped / 0 failed** |
| 5 | e2e 覆盖缺口（无 jobs 结果中心端到端用例）被**显式登记**而非省略 | `E-001` §3；本条目 F-002 |

### Findings

#### F-001 · e2e 顺序假设依赖「没有更早的文件」——同类缺陷可能再次出现

- 严重度：med
- 建议：**recommended**
- 描述：本次修复（`00-` 前缀）让 fresh-seed 用例真的先跑，但机制仍是**文件名排序**这一隐式约定。任何后续波次新增 `00-*` 之前的文件（或改动排序规则）都会再次打破它，而失败现象是「登录 401」这种与真实缺陷难以区分的形式。更稳的形态是把 fresh-seed 用例放进独立的 Playwright project/依赖链，或让挂具为它单独提供一个 fresh 库。
- 证据：`E-001` §2（根因与两向证据）；`apps/web/e2e/00-force-password-change.spec.ts` 注释。
- 状态：**open（recommended）** —— 处置建议：登记到 roadmap「未决项统一登记」（e2e 挂具顺序契约），由后续符合性波次决定是否改为 project 依赖。

#### F-002 · e2e 分母不含 jobs 结果中心端到端用例

- 严重度：med
- 建议：**recommended**
- 描述：判据 4（结果中心体验）与判据 5（权限/Profile 安全）的**浏览器侧**证据目前来自渲染/交互级测试（真实 `jobs.json` + 生产渲染链 + 真实路由断言）与 HTTP 契约测试，而不是真实浏览器端到端路径（提交批量导出 → 观察进度 → 下载）。这不是缺陷断言，但意味着「浏览器/自动化回归」对本 VP 新增面的覆盖是**间接**的。
- 证据：`apps/web/e2e/**` 11 个 spec 清单（`E-001` §3）；`GOAL-005` 的交互级测试 24 例。
- 状态：**open（recommended）** —— 处置建议：由用户/审计裁量为「R5 内补一条 e2e 用例」或「登记为后续波次项」。**不静默扩范围**。

### 必改项汇总（required）

**无。**

### 结论 + 建议下一步

判据 1～6 的证据链完整且可核对；e2e 恢复全绿并暴露/修复了一处既有挂具缺陷；覆盖缺口已诚实登记。R5 的 self 腿 **pass**。

**建议下一步**：C3 调用本地 grok build（grok 4.6 · high · `/audit`）做**关门 independent 审计**（任务书 `attachments/grok-prompt-r5-independent.md`），重点核验：① 判据 1～6 的证据是否真的充分、有无「用文档代替证据」；② e2e 失败根因判定是否成立（是否掩盖了真实回退）；③ 是否有未登记的范围外改动或 pinned 工件触碰；④ R4 三条残余回填 `fixed` 是否名实相符。

**本条不修改** `status`、检查点或派生 `progress`。
