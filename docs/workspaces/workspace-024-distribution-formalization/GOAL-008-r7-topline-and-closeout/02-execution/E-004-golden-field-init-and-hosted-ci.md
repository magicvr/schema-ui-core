---
status: active
created: 2026-08-30
updated: 2026-08-30
parent: GOAL-008-r7-topline-and-closeout
version: 0.1.0
---

# E-004 · golden-field 线上初始化与 hosted CI 首跑闭环（2026-08-30）

- **背景**：残余 1（hosted CI 实触发 · consumer-regression）复审触发兑现——用户 P-004 授权（2026-08-30：「采用初始化并推送」「推送后立即尝试触发」）。
- **初始化**：本地 12 commits（VP-023/024 实证史）**首次推送** origin `main`（HEAD `8631d53`：README 更新至冻结面 v1.4.0 终态 · 移除跟踪二进制 `server.exe` + `.gitignore` 补全二进制/日志 · `web/pnpm-workspace.yaml` 终值 version-set 补交）。远端从空仓变为**公开消费实证的持久可复现宿主**（closure-report 判据 #2/#3/#5 证据链可克隆；A-002 F-005「origin 头部不可见」根因消除）。
- **hosted 首跑三连**（全部 `workflow_dispatch` · ubuntu-latest）：
  1. `33286154992` ❌ —— `pnpm/action-setup@v4` 默认读**仓库根** `package.json` 的 `packageManager`；`web/package.json` 中的声明不被识别 → 修复 commit `ba052e7`（显式 `version: 11.11.0`；web 声明保留供本地一致性）。
  2. `33286191334` ❌ —— Go `@latest` 升级/build ✓ · pnpm 免凭据 install ✓ · **四探针全绿** · `shutdown.complete` 断言 ✓，但步骤仍 exit 1（收尾段 trap 内 kill 无 `|| true`、wait 返回值受 `set -e` 影响）→ 修复 commit `8ef02e9`（trap/kill/wait 全 `|| true` + 显式 `exit 0` + 诊断输出）。
  3. `33286302663` ✅ **PASS**（1m9s）：apps/api @latest + 六包免凭据消费 · serve 冒烟（healthz/readyz 200）· 四探针全绿（protocol 2.9 / render 1573 B / six-package / token）· SIGTERM → `shutdown.complete`（RT-D02 出口）断言通过。
- **结论**：**残余 1 → 核销**（hosted 实触发完成）；A-002 F-005 → **fixed**。宿主证据：https://github.com/magicvr/golden-field/actions/runs/33286302663
- **副产品认知**：hosted 与本地执行面存在差异（action-setup 版本源解析、bash 收尾语义）——正是「登记而非 hosted acceptance」口径（D-001 / closure-report §3-1）的价值实证；两处修复均回写 workflow 注释留痕。
- **golden-field commits**：`8631d53`（初始化）· `ba052e7` / `8ef02e9`（hosted 修复）· `52b7220`（README 闭环记录）。
- **主仓 git checkpoint**：`95f5d78b`（GOAL-008 台账 · 残余 1 核销与 F-005 闭合 · E-004 首版）。