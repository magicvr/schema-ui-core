---
id: E-004
goal: GOAL-001-admin-functional-modules
title: R3 第四批批末统一验证完成（V-001～V-008 全绿）
date: 2026-08-17
status: recorded
parent: null
created: 2026-08-17
updated: 2026-08-17
version: 1.0.0
---

# E-004 · R3 第四批批末统一验证完成（2026-08-17）

## 结果

| 项 | 命令 | 结果 |
|----|------|------|
| V-001 | `cd apps/api && go build ./...` | ✅ exit 0 |
| V-002 | `cd apps/api && go test ./...` | ✅ 全绿（exit 0） |
| V-003 | `cd apps/api && go vet ./...` | ✅ exit 0 |
| V-004 | `cd apps/web && npm run test` | ✅ vitest **1038/1038** |
| V-005 | `cd apps/web && npm run build` | ✅ `tsc -b` + `vite build` |
| V-006 | e2e 双 profile | ✅ mvp **8 passed** / admin **8 passed**（各 1 skipped，profile 对应变体） |
| V-007 | `bash scripts/smoke.sh`（mvp，非 disposable） | ✅ **exit 8**（SM-001～005 + SM-007 PASS；SM-006 SKIP 为预期） |
| V-008 | `bash scripts/smoke.sh --disposable`（隔离 `ci-s5-batch4-20260817`） | ✅ **exit 0**（SM-001～007 + SM-006 PASS；种子可重复性 + 重启持久化） |

## 环境说明

- 本地执行使用 `GOCACHE` 指向工作区内目录；`GIT_COMMIT` 传完整 40 位 HEAD `816ff5f7afe3b30017854877c4e1ca3b1a7894a9` 生成 claim。
- V-007 用默认栈（`schema-ui-core`，端口 25080/25081）；V-008 用隔离 project（`ci-s5-batch4-20260817`，端口 25180/25181），完成后已 `docker compose down`（隔离 project 带 `-v` 清理卷）。

## 结论

- **R3 第四批（钱包/账务：GOAL-019/020/021/022）批末统一验证完成，波次级门禁解除。**
- V-001～V-008 全绿；本批未改变 Profile 默认集 / 模块矩阵 / Manifest 装配语义 / 协议 pin，`go` 判定无影响、不暂挂。
- E-002（初次受阻）与 E-003（V-001～V-006 部分绿）作为过程记录保留；本 E-004 为最终收口。