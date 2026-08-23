# I-002 证据 · 活栈计时复核（真实浏览器 × Go API + Vite dev）

## 方法

- **两栈对比**：基线 = HEAD 提交 `ba7d5c6` 的父提交 `0878d7f`（W24 收盘，修复前：SQLite `SetMaxOpenConns(1)` / wallet-ensure 无条件 POST+整页重拉 / 无 schema 缓存 / 无请求合并）经 `git worktree` 检出；当前 = 本工作树（修复后，含全部 W25 修复）。两栈均以 admin profile、独立临时 SQLite、`go run ./cmd/server` + Vite dev（`node vite.js`，WEB_PORT 25173）启动。
- **浏览器驱动**：Playwright Chromium（`locale: en-US`），流程：登录（含强制改密）→ `/my-wallet` 首访 → 整页 reload → RTT 150ms 模拟下 reload（CDP `Network.emulateNetworkConditions`）→ SPA 页内导航（`pushState`+`popstate`，dashboard→/my-wallet）→ SPA 二次回访。
- **指标**：页面相关 API 请求计数（`GET /api/wallet/me`、`GET /api/wallet/me/entries`、`GET /api/schema/my-wallet`、`POST /api/wallet/me`）与**呈现耗时**（导航起点 → 3 张 statCard 骨架全部消失，含开通/刷新波）。
- 原始 JSON 行（measure.js 输出，两栈各 5 轮）：见本文下方。

## 结果（页面相关请求数）

| 轮次 | 基线 | 当前 | 减幅 |
|------|------|------|------|
| 首访（钱包自动开通） | 9+3+2+1 = **15** | 2+2+2+1 = **7** | −53% |
| 整页刷新（已开通） | 9+3+2+1 = **15** | 3+3+2+0 = **8** | −47% |
| RTT 150ms 刷新 | 15 | 7 | −53% |
| SPA 导航进入 | 9+3+1+1 = **14** | 2+2+1+0 = **5** | −64% |
| **SPA 二次回访（已缓存）** | 9+3+1+1 = **14** | **1+1+0+0 = 2** | **−86%** |

- schema 文档缓存直接验证：SPA 二次回访轮 `GET /api/schema/my-wallet` 基线 1 → 当前 **0**。
- 基线每轮 `GET /api/wallet/me` = 9（3 statCard 独立请求 × 进页 + POST 触发的整页重拉波），当前合并后 = 1–3。

## 呈现耗时（readyMs）

| 轮次 | 基线 | 当前 | 变化 |
|------|------|------|------|
| 首访 | 518 ms | 408 ms | −110 ms（−21%） |
| 整页刷新 | 399 ms | 330 ms | −69 ms（−17%） |
| RTT 150ms 刷新 | 5 680 ms | 4 249 ms | **−1 431 ms（−25%）** |
| SPA 导航进入 | 135 ms | 71 ms | −64 ms（−47%） |
| SPA 二次回访 | 113 ms | **55 ms** | **−51%** |

慢网络（150 ms RTT）下收益 ≈ 请求数差 × RTT，绝对值约 1.4 s/次页面进入，为**可感提升**；本机 TCP 下减 17%–51%。

## 原始输出（measure.js）

```
[BASELINE] {"label":"first-open","readyMs":518,"GET /api/wallet/me":9,"GET /api/wallet/me/entries":3,"GET /api/schema/my-wallet":2,"POST /api/wallet/me":1}
[BASELINE] {"label":"refresh","readyMs":399,"GET /api/wallet/me":9,"GET /api/wallet/me/entries":3,"GET /api/schema/my-wallet":2,"POST /api/wallet/me":1}
[BASELINE] {"label":"refresh-rtt150","readyMs":5680,"GET /api/wallet/me":9,"GET /api/wallet/me/entries":3,"GET /api/schema/my-wallet":2,"POST /api/wallet/me":1}
[BASELINE] {"label":"spa-nav","readyMs":135,"GET /api/wallet/me":9,"GET /api/wallet/me/entries":3,"GET /api/schema/my-wallet":1,"POST /api/wallet/me":1}
[BASELINE] {"label":"spa-nav-cached","readyMs":113,"GET /api/wallet/me":9,"GET /api/wallet/me/entries":3,"GET /api/schema/my-wallet":1,"POST /api/wallet/me":1}
[CURRENT] {"label":"first-open","readyMs":408,"GET /api/wallet/me":2,"GET /api/wallet/me/entries":2,"GET /api/schema/my-wallet":2,"POST /api/wallet/me":1}
[CURRENT] {"label":"refresh","readyMs":330,"GET /api/wallet/me":3,"GET /api/wallet/me/entries":3,"GET /api/schema/my-wallet":2,"POST /api/wallet/me":0}
[CURRENT] {"label":"refresh-rtt150","readyMs":4249,"GET /api/wallet/me":2,"GET /api/wallet/me/entries":2,"GET /api/schema/my-wallet":2,"POST /api/wallet/me":0}
[CURRENT] {"label":"spa-nav","readyMs":71,"GET /api/wallet/me":2,"GET /api/wallet/me/entries":2,"GET /api/schema/my-wallet":1,"POST /api/wallet/me":0}
[CURRENT] {"label":"spa-nav-cached","readyMs":55,"GET /api/wallet/me":1,"GET /api/wallet/me/entries":1,"GET /api/schema/my-wallet":0,"POST /api/wallet/me":0}
```

## 说明与边界

- 计数含 React StrictMode 开发双效应（首访/刷新轮 schema ×2、基线 me ×9 含 POST 后整页重拉波）——两栈条件一致，不改变对比方向；SPA 导航轮为单效应，最干净。
- 整页 reload 清空前端进程内存，schema 缓存仅在**同会话 SPA 导航**间生效（轮 5 已验证）。
- 环境：Windows 本机（localhost TCP）；RTT 模拟轮放大请求数差异，用于"可感提升"判断。