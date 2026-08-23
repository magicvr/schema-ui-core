---
title: E-006 · S5 回归事实：e2e 补跑履约 + 基线实验定性既有回归；全量回归绿
status: recorded
created: 2026-08-23
updated: 2026-08-23
parent: GOAL-033-w22-residual-closeout
version: 0.1.0
---

# E-006 · S5 回归与 A1 补跑结论（2026-08-23）

## C17 · 全量回归证据

| 项 | 结果 |
|----|------|
| `go build ./...`（apps/api） | exit 0 |
| `go test ./...` 全量 | 全部包 ok；`internal/store` 首轮 5 失败 = 新增迁移 0049 的黄金断言未同步 → 修复后 **store 包 ok（45.044s）** |
| 黄金断言修复清单 | `identity.go` completeFingerprintCatalogHead 48→49；`identity_test.go` lockedHeadExtraTables[49]={}（data-only 头迁移先例，放宽 `len(extra)==0` 守卫并留注释）；`migrate_test.go` [1..48]→[1..49] + want 表补 0049 行 + catalog len 49；`operations_test.go` / `restart_test.go` 尾断言 48→49 |
| vitest 全量（apps/web） | **76 文件 / 1088 测试全绿**（覆盖 H3 改名 + A4 变更） |
| tsc -b --noEmit + npm run build | **TSC_EXIT=0；BUILD_EXIT=0（vite ✓ built in 13.40s；chunk>500kB 警告为既有非阻断）** |

## C2 · A1 e2e admin M3 补跑 —— 履约完成，闭合叙事调整

时间线：
1. 首次补跑（子代理）：环境阻塞**确证解除**（8011–8110 排除区间消失；API 25080 / Vite 25173 成功启动）；暴露 strict-mode 选择器冲突 → 本波修复 5 处 `exact:true`（见 E-005）。
2. 二次补跑（编排器）：选择器问题消除，但 M1 断言失败——登录成功后停留在 `/` 而非 `/dashboard`（`localization.spec.ts:62`，稳定复现、非负载抖动）。
3. **基线实验（决定性）**：`git stash` 全部 W22 apps/ 改动后在 HEAD 上重跑 → **同样失败**。结论：该 `/dashboard` 断言漂移是**先于 W22 存在的既有回归**（W14–W21 某波引入，疑似 home 推导/路由面），与本目标六项整改及原始端口 residual 均无因果关系。

处置（P-003 合法叙事）：
- W7·GOAL-007 F-002 residual 的实质 = 「端口排除区间导致浏览器证据不可产出」。现触发条件（区间解除）已兑现并完成补跑附日志（`attachments/e2e-admin-m3-rerun-final.log` + 历史 `e2e-admin-m3-rerun.log`）→ 该残余**按复审履行完毕关闭**，不标 fixed（补跑未产出门禁级通过证据）。
- 新发现登记为本目标 **N-001（non-blocking · 移交）**：admin 登录后未跳转 `/dashboard`（先于 W22 存在）；建议下一符合性波次以 home 推导链路为首要排查点；不阻断本目标关门（超出其范围且非其引入）。

## 进度影响

C2 ✓（履约+定性）、C16 ✓（vitest 确证）→ 累计 **15/18**；C17 待 tsc/build 追加记录后 ✓（16/18）；C18 关门审计待执行。
