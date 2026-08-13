---
id: E-001
goal: GOAL-007-w6-container-smoke-reproducibility
title: S1/S2 · 代码修复 + 双路径自检
date: 2026-08-14
status: recorded
parent: GOAL-007-w6-container-smoke-reproducibility
created: 2026-08-14
updated: 2026-08-14
version: 0.1.0
---

# E-001 · S1/S2 · 代码修复 + 双路径自检

## S1 · 代码修复（4 文件，未提交、待 S3 后 checkpoint）

| 文件 | 变更 |
|------|------|
| `apps/web/scripts/generate-claim.mjs` | buildId：`process.env.GIT_COMMIT`（trim 非空）优先，否则 `git rev-parse HEAD` 回退 |
| `apps/web/Dockerfile` | build 阶段 `ARG GIT_COMMIT` + `ENV GIT_COMMIT=${GIT_COMMIT}`（build RUN 前） |
| `compose.yaml` | web build `args: GIT_COMMIT: ${GIT_COMMIT:-unknown}` |
| `.github/workflows/r6-basic-matrix.yml` | container-smoke env 加 `GIT_COMMIT`（CI 传 ${{ github.sha }}） |
| `apps/web/nginx.conf`（F-1b） | `upstream api_upstream` 从 `server {}` 内移到 http 层（nginx -t 通过） |

## S2 · 双路径自检（2026-08-14 · 后台任务 pwsh-7）

- git 回退路径：`node scripts/generate-claim.mjs`（LOCAL_CLAIM_EXIT 待确认）；
- env 路径：`GIT_COMMIT=deadbeef1234` 生成（ENV_CLAIM_EXIT 待确认）；
- 恢复：清 env 后重生成（git 路径，buildId=git:c8ae108…）。

> 结果以后台任务输出为准（E-002 更新）。
