---
id: E-002
goal: GOAL-001-admin-functional-modules
title: R3 第四批（钱包/账务）批末统一验证尝试（环境受阻，未完成）
date: 2026-08-17
status: recorded
parent: null
created: 2026-08-17
updated: 2026-08-17
version: 1.0.0
---

# E-002 · R3 第四批批末统一验证尝试（2026-08-17）

## 背景

GOAL-019/020/021/022（S-14 钱包/账务批）已于 2026-08-16 同日关门，四处关门审计均约定「波次级 e2e 双 profile + V-007/V-008 容器冒烟留批末统一验证，批末必须补跑，失败则回流」。本条目记录 2026-08-17 的补跑尝试。

## 事实

| 项 | 结果 |
|----|------|
| `go test ./...`（apps/api） | 除 `cmd/server` 外全部包 **ok**；`TestServerProcessRestartPersistsUsers` 失败于其内部 `go build` 子服务 → `go: failed to trim cache: …\go-build\trim.txt: Access is denied`（本地沙箱拒绝 Go 构建缓存写入） |
| `go build` / `go vet`（apps/api） | 未通过（同上 cache Access denied） |
| `npm run test`（vitest） | 未通过（esbuild 启动原生服务 `spawn EPERM`，本地沙箱禁止捕获子进程 stdio） |
| `npm run build`（web） | 未通过（prebuild `generate-claim.mjs` 执行 `git rev-parse HEAD` → `spawnSync git EPERM`） |
| V-006 e2e 双 profile | 未执行（依赖上述 go run / vite dev 子进程，同受沙箱边界） |
| V-007 / V-008 容器冒烟 | 未执行（Docker engine 连接被沙箱拒绝：`dockerDesktopLinuxEngine` npipe permission denied） |

## 结论

- **波次级门禁未解除**：本次未能产生任何可采信的批末验证证据，**不构成 PASS**，不回流、不关闭 GOAL-019/020/021/022 的波次级验证义务。
- 上述失败均为**本地执行环境限制**（沙箱对构建缓存 / 子进程 stdio / Docker 管道的拒绝），非项目回归。
- 待补跑清单：V-001～V-008（含 e2e 双 profile 与 V-007 exit 8 / V-008 exit 0）；跑通后另记 E-003 并回填本行。
