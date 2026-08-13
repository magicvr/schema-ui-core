---
id: E-003
goal: GOAL-006-w5-recordview-declared-fields
title: S3 · dev/文档卫生（commit 5c309ff / c420e5d）
date: 2026-08-13
status: recorded
parent: GOAL-006-w5-recordview-declared-fields
created: 2026-08-14
updated: 2026-08-14
version: 0.1.0
---

# E-003 · S3 · dev/文档卫生

## 事实

- commit `5c309ff`（2026-08-13，fix(dev)）：dev 脚本**等待 API ready 后启动 Web**，stop 按 **PID 精确停止**（修复按端口误停/过早启动问题）。
- commit `c420e5d`（2026-08-13，docs(quickstart)）：QUICKSTART 修正 `dev.cmd` 调用前缀与排版。

同批未归档提交，归入本波 S3；不改变任何生产行为。

## 证据

- `git show --stat 5c309ff` / `c420e5d`。
