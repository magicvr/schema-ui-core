---
id: E-004
goal: GOAL-006-w5-recordview-declared-fields
title: S4 · HEAD 全量回归（V-001～V-006 绿；V-007/V-008 基础设施受阻）
date: 2026-08-14
status: recorded
parent: GOAL-006-w5-recordview-declared-fields
created: 2026-08-14
updated: 2026-08-14
version: 0.1.0
---

# E-004 · S4 · HEAD 全量回归

## 候选身份

- **HEAD `c8ae108`**（feat(vision): 新增 VP-011 并将协议 pin 升至 v2.8.0；此前 working tree clean）。
- 运行候选即当前主线；`go` 消费候选身份由本次重验证更新（VP-011 §消费决策记录）。

## 冻结矩阵回归结果（2026-08-14 实测）

| 命令 | 结果 |
|------|------|
| V-001 `go build ./...`（apps/api） | ✅ exit 0 |
| V-002 `go test ./...`（apps/api） | ✅ 全包 ok（cmd/server、auth、handler、composition、modules 等） |
| V-003 `go vet ./...`（apps/api） | ✅ exit 0 |
| V-004 `npm test`（apps/web） | ✅ **48 文件 / 889 测试** pass |
| V-005 `npm run build`（apps/web） | ✅ tsc + vite 成功；claim 重生成 `buildId=git:c8ae108d1026127d66dd3198a4401ffc9baf4355`，digest `sha256:7f8d8c7b…` |
| V-006 e2e（playwright，mvp / admin） | ✅ **8 pass + 1 skip（双 Profile）**；E2E_MVP_EXIT=0、E2E_ADMIN_EXIT=0 |
| V-007 smoke（`bash scripts/smoke.sh`） | ⚠️ **本地不可复跑**：`docker compose up --build` 的 web 镜像构建失败（容器内 `git rev-parse` 失败）→ 见 F-1 |
| V-008 disposable smoke（`--disposable`） | ⚠️ **本地不可复跑**：依赖 F-1 的隔离 compose 镜像；脚本本身 SM-001 前置即失败 → 见 F-1 |

## F-1 · 容器 smoke 复现性破损（post-go，W3 引入）

- **事实**：`apps/web/scripts/generate-claim.mjs` 自 VP-010 W3（commit `5e4c384`，2026-08-13，晚于 2026-08-10 的 `go` 签发）起，`buildId` 强制 `git rev-parse HEAD`（无 env 回退）；而 `.dockerignore`（自 VP-002 `719e331` 起）排除 `.git`。→ **compose/容器 web 镜像构建必然失败**（实测：`fatal: not a git repository`，web build 7/7 CANCELED）。CI `container-smoke`（r6-basic-matrix，`docker compose -p … build`）同路径，现行代码下同样会失败。
- **影响**：V-007 / V-008 冻结证据在 HEAD 不可复现；VP-008 §`go` 消费有效性「冻结命令与关键证据可执行性」复核受影响。W3 关门时仅本地 vitest/build 重验，未复验容器构建路径。
- **处置建议**：VP-010 W6 波次修复（generate-claim 支持 `GIT_COMMIT` env 回退 + Dockerfile ARG + compose/CI 接线），或用户按 P-004 裁决 residual / 暂挂。本 finding 不阻断 GOAL-006 自身范围（recordView 波次代码全部绿），但**阻断 `go` 消费的干净重验证**。
