---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-005-r4-fork-comparison
version: 0.1.0
---

# D-001 · 实验设计定档（2026-08-29）

## 决策

1. **演进集（用户裁决）**：`apps/api/v0.3.0`（`4f7cb0f1`）→ `apps/api/v0.4.0`（`00d97b5b`）真实演进：serve 面（`server` 包）新增 + CLI `serve` 子命令 + 模板薄封装 + config 面。（v0.3.0 的 `JoinIdentifiers` breaking 属上一演进段，不在本样本；迁移说明以 v0.3.0→v0.4.0 为准。）
2. **双模型口径**：
   - **fork 模型**：`fork-sim` 仓 = schema-ui-core 克隆（基线 = `apps/api/v0.3.0`），注入 2 个本地定制 commit（模拟业务 fork：① `apps/api/cmd/schema-ui/templates/main.go.tmpl` 加定制注释（与上游 v0.4.0 重写同文件 → 冲突点 1）② 根 `README.md` 定制（冲突点 2））→ `git merge apps/api/v0.4.0`（非 fast-forward 路径）→ 计数 `diff-filter=U` 冲突文件 + 人工按迁移说明解冲突（薄封装换装）+ `go build ./...` 验证 → 全程计时。
   - **包模型**：golden-field 新 worktree 重演 v0.3.0→v0.4.0：`git worktree` 检 v0.3.0 态（go.mod v0.3.0 + 旧冒烟 main）→ 按迁移说明 bump（go.mod → v0.4.0 · `cmd/server` 换薄封装 · web 六包 + `.npmrc` 钉 npmjs）→ `go mod tidy` + `go build` + serve 冒烟 → 计时；冲突计数 = 0（无 merge · 无 replace）。
3. **计时规则**：各阶段 `Measure-Command`/秒表；暖口径（本地 Go/pnpm 缓存预热）为主、冷口径注明；相对对比为结论依据，不主张跨机绝对秒数。
4. **契约迁移成本**：条目化——fork = 需人工 diff/改写点（main 形态、config 键、模板差异步骤数）；包 = changelog 迁移说明执行步骤数。
5. **输出**：`fork-comparison-report` 附件（耗时矩阵 + 冲突计数 + 迁移成本 + 定量结论）。

## 6. 执行偏差补记（v0.1.1 · 2026-08-29 · A-002 F-002 响应）

定制点 2 实跑采用 `QUICKSTART.md`（行尾追加）而非原定的根 `README.md`——复核确认演进集 v0.3.0→v0.4.0 **未改动** `README.md`（按原文不会产生内容冲突），故以演进集内确有改动的 `QUICKSTART.md` 作为冲突点 2（行尾追加场景 → Auto-merge）。冲突计数 1 仍由点 1（`main.go.tmpl` 整文件重写）驱动，结论方向不受影响。历史 `%TEMP%\fork-sim` 残留已删除（F-003 fixed）。

## 未选方案

| 项 | 未选 | 理由 |
|----|------|------|
| 合成演进集 | 自己造配置键/迁移/依赖样本 | 用户裁决用真实演进集；合成与 VP-022 R4 已核销演练重叠 |
| 冷缓存严格复现 | 每步清 Go/pnpm 全缓存 | 时间成本高且非相对对比所需；口径注明即可 |
| 远端 CI 对照 | 双模型各上 runner | 本波为方法学实证（本地相对计时）；hosted 已登记 R7 |