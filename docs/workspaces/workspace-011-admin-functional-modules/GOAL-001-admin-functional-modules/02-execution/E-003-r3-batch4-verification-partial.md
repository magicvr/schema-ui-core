---
id: E-003
goal: GOAL-001-admin-functional-modules
title: R3 第四批批末统一验证（V-001～V-006 绿；V-007/V-008 环境受阻）
date: 2026-08-17
status: recorded
parent: null
created: 2026-08-17
updated: 2026-08-17
version: 1.0.0
---

# E-003 · R3 第四批批末统一验证（2026-08-17）

## 结果

| 项 | 命令 | 结果 |
|----|------|------|
| V-001 | `cd apps/api && go build ./...` | ✅ exit 0 |
| V-002 | `cd apps/api && go test ./...` | ✅ 全绿（exit 0） |
| V-003 | `cd apps/api && go vet ./...` | ✅ exit 0 |
| V-004 | `cd apps/web && npm run test` | ✅ vitest **1038/1038** |
| V-005 | `cd apps/web && npm run build` | ✅ `tsc -b` + `vite build`（1851 modules transformed） |
| V-006 | `APP_PROFILE=mvp / admin && npm run test:e2e` | ✅ mvp **8 passed** / admin **8 passed**（各 1 skipped，为 profile 对应变体） |
| V-007 | `bash scripts/smoke.sh` | ⛔ 未执行（Docker engine 连接被本地沙箱拒绝） |
| V-008 | `bash scripts/smoke.sh --disposable` | ⛔ 未执行（同上） |

## 执行环境说明

- `GOCACHE` 指向工作区内 `.gocache`（本地沙箱禁止写 `AppData` 默认 Go 缓存）；`GIT_COMMIT` 传完整 40 位 HEAD（`816ff5f7afe3b30017854877c4e1ca3b1a7894a9`），使 claim 生成跳过 `git` 子进程。
- V-007/V-008 依赖 Docker engine（npipe `dockerDesktopLinuxEngine`），本地沙箱在 workspace-write 模式拒绝连接；批准策略为 `never`，无法升级重试。

## 结论

- **V-001～V-006 全绿**：钱包批（GOAL-019/020/021/022）代码 / 回归 / e2e 双 profile 层面无回归。
- **V-007/V-008 容器冒烟仍未执行** → 波次级门禁**部分解除**（e2e 双 profile 已补跑；容器冒烟待本机正常环境补跑后回填，另记 E-004）。
- 本条目**不构成完整 PASS**，不回流、不关闭容器冒烟义务。
