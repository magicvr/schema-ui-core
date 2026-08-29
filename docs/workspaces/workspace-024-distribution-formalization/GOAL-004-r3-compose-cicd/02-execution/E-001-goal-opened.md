---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-004-r3-compose-cicd
version: 0.1.0
---

# E-001 · 目标建立（2026-08-29）

1. **立项**：承接 Root 纲领 R3（compose/CI 实跑 · VP-024 判据 #3 · go 后清单 ③ compose CI 实跑 + R1 残余 1 信号 harness + workspace-023 F-001/F-007 残留核销）；goal-tree 同步。
2. **现状核验**：
   - 主仓 `compose.yaml`：api（secrets fail-closed · readyz healthcheck · `stop_grace_period: 15s`）+ web（nginx 反代）——**从未本环境实跑**（GOAL-005 F-001）；
   - golden-field `consumer-regression.yml`：仍含 GH Packages token 步骤（npmjs 公开后已过时）+ `go run ./cmd/server` 旧冒烟形态（thin wrapper 下会挂起）；
   - golden-field 四个 web 探针（probe/render/six/token-check）均为**静态校验**（无 HTTP 依赖）；
   - `apps/api/Dockerfile` 基底待 S1 确认（harness 工具面）。
3. **信息门禁**：I-024-002（CI 环境）为 R3 前置 → 本目标核销。