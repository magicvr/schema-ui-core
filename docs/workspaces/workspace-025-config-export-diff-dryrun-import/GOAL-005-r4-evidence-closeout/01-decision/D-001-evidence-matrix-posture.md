---
id: D-001
title: C1 证据矩阵口径（lead · 引用 D-004 / E-004 / E-005）
date: 2026-08-30
status: accepted
---

# D-001 · 证据矩阵口径（2026-08-30）

1. **矩阵结构**：VP-025 六条方向级退出判据逐行；每行 = 判据、阶段证据链接（合同/实现/测试/冒烟）、状态（verified）。
2. **证据源**：R1（GOAL-002 D-001/D-002 + A-001）、R2（GOAL-003 D-001 + E-002 + A-001）、R3（GOAL-004 D-001/D-002 + E-002/E-003 + A-001）、红线核账（Root E-004 预检 + 本目标全量核对）、测试（configpkg_test.go 18 用例）、回归（`go test ./...` 49 包）、CLI 冒烟（export/diff/dry-run/import 退出码与往返）。
3. **越界核账**：`git diff --name-only` 覆盖开区（cf68c7ce）→ 当前 HEAD 全提交面；红线域 = `apps/api/internal/store`（迁移台账）/ `kernel/profile.go`（Profile 装配）/ `apps/web/src/protocol/upstream`（provenance）/ 任何 migration 面。
4. **关闭条件**：六行全部 verified + 核账零触碰 → C1 关门（GOAL-005 1/3）。