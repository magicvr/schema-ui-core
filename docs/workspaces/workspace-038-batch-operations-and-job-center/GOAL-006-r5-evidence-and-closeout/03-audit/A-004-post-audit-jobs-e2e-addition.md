---
id: A-004-post-audit-jobs-e2e-addition
doc: audit-entry
parent: GOAL-006-r5-evidence-and-closeout
status: recorded
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# A-004 · 审计后补测（jobs 结果中心 e2e）与 F-002 闭合

- **source**：orchestrator（`/govern`，执行用户 2026-09-19 指令）
- **scope**：R5 C3 之后的**增量证据**——用户裁决「先补 jobs e2e 再关门」，故在 `A-003` 之后追加端到端覆盖并闭合 `A-001`/`A-002` 的 F-002。

## 1. 用户指令与范围

用户 2026-09-19 对「是否现在确认 VP-038 关门」的选择：**「先补 jobs e2e 再关门」**。据此本条目把 `A-001` F-002 / `A-002` F-002（浏览器回归未驱动 jobs 结果中心、默认 e2e 跑 mvp 不挂载 `admin.jobs`）由「登记为 bounded residual」改为**实现并闭合**。

## 2. 产物与证据

| 项 | 内容 |
|----|------|
| 新增用例 | `apps/web/e2e/jobs-result-center.spec.ts` |
| profile 守卫 | `test.skip(appProfile !== "admin", …)` —— `admin.jobs` 仅在 admin 默认集；mvp 下**显式跳过**（不伪装成已覆盖） |
| 端到端路径 | 登录（含首次强制改密）→ 侧边栏进 users 页 → 勾选真实行 → `导出所选` → **观察真实进度**（`[data-jobs-batch-export-progress]`）→ 终态出现下载按钮 → 取回 CSV（断言文件名 = 服务端 `users-selection.csv`）→ 侧边栏进 jobs 页（结果中心）→ 定位 `jobs.batch-export` 行且状态为**本地化终态** `Succeeded` → 打开行操作溢出菜单 → `Download result` **再次下载同一文件** |
| 单跑 | `APP_PROFILE=admin npx playwright test jobs-result-center.spec.ts` → **1 passed (20.8s)** |
| 双 profile 全量 | 默认（mvp）：**16 passed / 5 skipped / 0 failed**（第 5 个 skip 即本 spec 的守卫）；`APP_PROFILE=admin`：**17 passed / 4 skipped / 0 failed**（本 spec 实跑通过）。两条命令均 exit 0 |
| 发现并修正的实现细节 | 初版用 `page.goto("/users")` 会因整页重载丢失内存会话而落到 `Session expired` 宿主失败面；改为与既有 spec 一致的**侧边栏导航**（client-side routing），并在用例内注释原因 |

## 3. 对既有意见的闭合

| finding | 原状态（A-003） | 现状态 | 依据 |
|---------|-----------------|--------|------|
| `A-001` F-002 / `A-002` F-002（缺 jobs 端到端覆盖） | bounded residual（roadmap 登记） | **`fixed`** | 本节 §2；roadmap 对应行已由「登记（未做）」改为「2026-09-19 `fixed`」 |
| `A-001` F-001 / `A-002` F-001（e2e fresh-seed 顺序契约） | bounded residual | **bounded residual（保持）** | 用户指令只要求补 jobs e2e；该条为挂具契约观察，已在 roadmap 登记并保留触发条件 |

## 4. 诚实边界（审计覆盖声明）

- `A-002`（independent）的复审**早于**本次新增用例，故独立腿**未**审查 `jobs-result-center.spec.ts` 本身。
- 本条把风险与结论分开陈述：新增内容是**测试**（无产品代码改动），其证据力由**执行结果**给出（双 profile 全量 + 单跑），而非由文档断言；覆盖方向是**增强**而非削弱（原先被判「间接」的浏览器侧证据现在有了端到端一层）。
- 本条目**不**声称 independent 腿已复核该用例；若后续需要，可在关门前的复审中一并覆盖。
- 仍未覆盖（保持登记）：重试/取消的浏览器端到端路径（当前由交互级 24 例 + HTTP 契约测试覆盖）。

## 5. 门禁判定

- 相关意见（`A-001`/`A-002`/`A-003`/本条）：全部已响应；**开放 required = 0**。
- 判据 7 的剩余条件仍为**用户书面确认**（`I-038-019`）；本条目不修改任何 `status`/`progress`，不改 VP-038 状态。
