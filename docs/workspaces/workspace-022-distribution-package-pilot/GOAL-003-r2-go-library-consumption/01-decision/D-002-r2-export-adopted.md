---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-003-r2-go-library-consumption
version: 0.1.0
---

# D-002 · internal 外移方案采纳（用户裁决）

用户 2026-08-29 书面裁决：**方案 A · 目录提升**（`internal/kernel` → `apps/api/kernel`；`internal/modules` → `apps/api/modules`；C 层保持 internal；单 go.mod 不变，G1 粗粒度保持）。

执行要点（D-001 决策点响应）：
1. 外移范围 = 冻结面清单 §1/§2 严格对应（A 层 kernel 全量 + B 层 modules 全量）；C 层（composition/auth/handler/server/store 方言等 20 包）不动。
2. F-001 联动（`users.New(*auth.Authenticator)` 的 C 层泄漏）：eve 按计划留 S3 验证矩阵实测（不预先裁定）；本轮 S3 首验选零依赖模块（dashboard）通过 → 泄漏仍待 users 装配实测。
3. 主仓等价性闸门：`go build ./...` exit 0 + `go test ./...` 全量（仅 1 个 PG drain harness 全量并发时序敏感失败，单跑与全包复跑均 PASS，判定非重构缺陷，见 A-001 F-003）+ `internal/composition` 全包复跑 PASS。