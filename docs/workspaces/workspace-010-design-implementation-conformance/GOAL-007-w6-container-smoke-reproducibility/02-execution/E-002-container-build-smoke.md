---
id: E-002
goal: GOAL-007-w6-container-smoke-reproducibility
title: S3 · 容器构建 + smoke 复跑（V-007/V-008 全绿）
date: 2026-08-14
status: recorded
parent: GOAL-007-w6-container-smoke-reproducibility
created: 2026-08-14
updated: 2026-08-14
version: 0.1.0
---

# E-002 · S3 · 容器构建 + smoke 复跑

## 候选

HEAD `c8ae108` + W6 修复（未提交，待 checkpoint）；隔离 compose project `ci-freshness-20260814`（APP_PROFILE=mvp；GIT_COMMIT=`c8ae108d…`）。

## 实测结果（2026-08-14）

| 项 | 结果 |
|----|------|
| claim 双路径自检（git 回退 / GIT_COMMIT env） | ✅ LOCAL_CLAIM_EXIT=0、ENV_CLAIM_EXIT=0 |
| nginx -t（修复后配置） | ✅ syntax ok / test successful |
| `docker compose build`（web 镜像含 claim + nginx 修复） | ✅ BUILD_EXIT=0 |
| `docker compose up -d` + 双端点 ready | ✅ BOTH_READY（25080 readyz + 25081 /） |
| **V-007** `bash scripts/smoke.sh`（SMOKE_EXPECTED_PROFILE=mvp） | ✅ SM-001~005 + SM-007 全 PASS；SM-006 skip（非 disposable）；**exit 8**（部分绿 = 预期语义） |
| **V-008** `bash scripts/smoke.sh --disposable` | ✅ SM-001~005 + SM-007 + **SM-006 PASS**（种子可重复性 + 重启持久化）；**exit 0**（完整绿） |

## F-1 闭合（三部分）

| id | 根因 | 修复 | 状态 |
|----|------|------|------|
| F-1a | claim 脚本 git 强制（W3 `5e4c384`，post-go） | `GIT_COMMIT` env 回退 + Dockerfile ARG + compose args + CI 接线 | fixed（构建通过） |
| F-1b | nginx `upstream` 嵌在 `server {}` 内（VP-009 W3 `7dbc3b5`，post-go） | 移到 http 层；nginx -t 通过；容器正常服务 | fixed（web 可达、登录代理 PASS） |
| F-1c | smoke.sh SM-007 仍按 S5 页面集断言（W1 拆出 dev.examples 后 mvp=users+roles、admin=+settings+activity） | SM-007 改为按 profile 取必需页集；mvp/admin/demo 分派 | fixed（SM-007 PASS） |

## 残余观察（非阻断）

- `apps/web/src/test-fixtures/app-manifest.{mvp,admin}.json` 仍为 W1 前页面集（10/12 页），与运行时投影（2/4 页）不一致——供单测使用、非 smoke 断言源；建议后续波次核对是否更新。
