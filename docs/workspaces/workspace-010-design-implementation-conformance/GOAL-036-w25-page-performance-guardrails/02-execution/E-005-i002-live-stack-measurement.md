---
date: 2026-08-23
scope: GOAL-036 S6（I-002 活栈计时复核）
---

# E-005 · I-002 活栈计时复核（真实浏览器对比实测）

## 方法（事实）

- 基线 = `0878d7f`（W24 收盘，修复前）经 `git worktree` 检出；当前 = 本工作树（修复后）。两栈均 admin profile + 隔离临时 SQLite + `go run ./cmd/server` + Vite dev（25173），由 `run-stack.ps1` 驱动。
- Playwright Chromium（en-US）驱动：登录（含强制改密）→ `/my-wallet` 首访 → reload → RTT 150ms 模拟 reload（CDP）→ SPA 页内导航 → SPA 二次回访；统计页面相关 API 请求数与"3 张 statCard 呈现就绪"耗时。
- 完整原始数据：`attachments/I-002-evidence.md`。

## 结果摘要（页面相关请求数 / 呈现耗时）

| 轮次 | 基线请求 | 当前请求 | 基线耗时 | 当前耗时 |
|------|----------|----------|----------|----------|
| 首访（开通） | 15 | 7 | 518 ms | 408 ms |
| 整页刷新 | 15 | 8 | 399 ms | 330 ms |
| RTT 150ms 刷新 | 15 | 7 | 5 680 ms | **4 249 ms（−1.4 s）** |
| SPA 导航进入 | 14 | 5 | 135 ms | 71 ms |
| SPA 二次回访 | 14 | **2** | 113 ms | **55 ms** |

- schema 文档缓存直接实证：SPA 二次回访轮 `GET /api/schema/my-wallet` 基线 1 → 当前 **0**。
- 结论：请求数 −47%～−86%（视轮次），慢网络（150 ms RTT）下呈现耗时 −25%（绝对值 −1.4 s/次），本机 TCP 下 −17%～−51%——**可感提升成立**。

## 边界（如实记录）

- 整页 reload 清空前端模块内存，schema 缓存仅在**同会话 SPA 导航**间生效（轮 5 实证）；Web 构建产物（vite build）下 StrictMode 双效应不存在，请求数与耗时预期更低。
- 计数含 StrictMode 双效应（两栈一致）；本机 localhost TCP 不放大网络层差异，RTT 轮为慢网络代理指标。
- 测量脚本与驱动器位于会话临时目录（`%TEMP%\w25-metrics\`，不入库）；worktree 基线已清理。

## 结论

**I-002 closed**（non-blocking → 已提交证据）。C6 剩余：自审 A-001 与关门台账同步。