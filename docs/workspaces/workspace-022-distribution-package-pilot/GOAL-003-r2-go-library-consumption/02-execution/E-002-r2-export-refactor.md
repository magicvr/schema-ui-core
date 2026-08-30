---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-003-r2-go-library-consumption
version: 0.1.0
---

# E-002 · 外移重构与装配首验事实（2026-08-29）

用户裁决方案 A 后执行（go 1.26.0）：

1. **git mv**：`apps/api/internal/kernel` → `apps/api/kernel`；`apps/api/internal/modules` → `apps/api/modules`（历史保留）。C 层 20 包留在 `internal/`。
2. **import 改写**：223 个 .go 文件批量替换 `apps/api/internal/modules` → `apps/api/modules`、`apps/api/internal/kernel` → `apps/api/kernel`；`docs/architecture/module-contribution-playbook.md` 同步（15 处；`internal/composition` 引用保留）。残留引用复查 = 0。
3. **编译回归**：`go build ./...` exit 0（全包编译通过）。
4. **测试回归**：`go test ./...` 全量——唯一失败 = `internal/composition TestShutdownDrainHarnessPostgres`（PG drain，全量并发下 EOF）；**单独复跑 PASS（3.20s）**；`internal/composition` 全包复跑 **ok（22.55s）** → 判定为全量并发时序敏感，非重构缺陷（A-001 F-003）。
5. **装配闭环首验（S3 前奏）**：golden-consumer（`golden.local/consumer` · replace 本仓）升级为最小真实组合根——`kernel.ResolveProfile("admin")` + `kernel.NewRegistry(kernel.BuiltinModules())` + `dashboard.New()`；`go run .` 输出：
   `kernel=2.0.0 profile=admin dashboard=admin.dashboard@2.0.0 range=>=2.0 <3.0 modules=23`（exit 0）。
   → **A/B 层外移后首次被外部模块成功消费**（internal 阻断解除的实证）。
6. **冻结面路径勘误**：freeze-face v1.0.0 → v1.0.1（`internal/kernel`/`internal/modules` 路径同步；契约内容不变，editorial）。

## 影响

- S2（外移重构）完成 ✅；S3 起点：最小装配首验 PASS，剩余 = 迁移台账 apply（双方言）+ users 装配（F-001 实测）+ B 层符号回填（F-002）+ 功能基线等价矩阵。