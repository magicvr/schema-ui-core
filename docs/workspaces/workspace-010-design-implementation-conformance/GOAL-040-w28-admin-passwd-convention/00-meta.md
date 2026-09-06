---
id: GOAL-040-w28-admin-passwd-convention
title: W28 · 现有库 admin 凭据约定（ADMIN_PASSWD）与测试账户机制退役
status: done
parent: GOAL-001-design-implementation-conformance
created: 2026-09-06
updated: 2026-09-06
version: 0.3.0
progress: 4/4
---

# GOAL-040 · W28 · 现有库 admin 凭据约定（ADMIN_PASSWD）与测试账户机制退役

## 概述

管理员手动测试后，登录强制改密（`must_change_password`）使 `admin` 的现密码不可知：`ADMIN_INITIAL_PASSWORD` 只对零用户库的 fresh bootstrap 生效，自动化测试与 AI 助手连**现有库**时因不知道 admin 现密码而撞墙。本波次落盘并实施两项决策：

1. **ADMIN_PASSWD 约定**：维护者在本地 gitignored 的 `apps/api/configs/.env` 设置 `ADMIN_PASSWD`，声明当前库 `admin` 的现密码；自动化测试与 AI 助手从环境变量 / `.env` 读取后登录。`ADMIN_PASSWD` 是**纯声明约定**——API 不读取、不据此重置密码。
2. **AI 助手可发现性**：让 AI 助手知道可以从 `.env`（或环境变量）拿到该密码。
3. **TEST_ADMIN 机制退役**：删除 `TEST_ADMIN_USERNAME/TEST_ADMIN_PASSWORD` 测试账户机制（E-011 引入，boot-time 每次启动 upsert 一个可用 admin，形如后门且不可发现；AI 助手与自动化测试通常不知道如何用它）。

## 成功标准

- [x] `apps/api/configs/.env.example`（唯一 canonical 模板）登记 `ADMIN_PASSWD` 约定，语义为「声明当前库 admin 现密码，供自动化测试与 AI 助手读取」，并明确「API 不读取、不重置密码」。
- [x] 仓库根 `AGENTS.md` / `QUICKSTART.md` / `apps/api/README.md` 让 AI 助手与自动化测试知道：连现有库时从 `.env` / 环境变量的 `ADMIN_PASSWD` 读取 admin 密码。
- [x] `TEST_ADMIN_USERNAME/TEST_ADMIN_PASSWORD` 相关代码全部移除（config / bootstrap / composition / 测试 / .env.example），`go test` 相关包全绿。
- [x] `scripts/smoke.sh` 支持以 `ADMIN_PASSWD` 作为 `SMOKE_PASSWORD` 的回退来源（连现有库路径不再撞墙）。

## 派生进度展示

`progress: 4/4` 由上方 4 个显式成功检查点等权计算；检查点变化后重算并同步 `goal-tree.md`。progress 不放行阶段、不关闭 finding、不推导 done。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | `ADMIN_PASSWD` 以哪份 gitignored `.env` 为权威？ | 方案冻结 | S1 | 用户裁决 | verified | — | 2026-09-06 用户确认：`apps/api/configs/.env`（与 E-011 既有测试凭据约定一致；canonical 模板 `.env.example` 登记） |
| I-002 | required | 是否本轮包含 TEST_ADMIN 移除实施？ | 方案冻结 | S1 | 用户裁决 | verified | — | 2026-09-06 用户确认：包含 |
| I-003 | required | `ADMIN_PASSWD` 是否被 API 读取并用于重置 admin 密码？ | 方案冻结 / 实施 | S1 | 设计决策 | verified | — | D-001：**不读取、不重置**（纯声明约定；避免重蹈 TEST_ADMIN 后门形状，不与 must_change_password 语义冲突） |

## 父目标

- `GOAL-001-design-implementation-conformance`（持续符合性程序 Root）

## 台账布局

三个可追加台账使用同名平铺目录：`01-decision/`、`02-execution/`、`03-audit/`；索引文件保留 frontmatter、摘要和条目索引；独立记录使用 `D-NNN-*`、`E-NNN-*`、`A-NNN-*` 文件。

## 备注

历史记录：TEST_ADMIN 机制由 [workspace-030 E-011](../../workspace-030-telegram-channel-runtime/GOAL-001-telegram-channel-runtime/02-execution/E-011-operator-config-and-test-admin.md)（2026-09-03）引入；本波次将其退役，历史五件套不改写。
