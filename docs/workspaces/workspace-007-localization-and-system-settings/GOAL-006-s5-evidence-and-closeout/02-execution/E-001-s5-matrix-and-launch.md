---
id: E-001
doc: execution
title: S5 · C1 证据矩阵 + C2 真实入口启动验证
status: recorded
parent: GOAL-006-s5-evidence-and-closeout
created: 2026-08-09
updated: 2026-08-09
version: 0.2.0
---

# E-001 · S5 C1 证据矩阵 + C2 真实入口（2026-08-09）

## 已发生事实

### C1 · 证据矩阵

- 落盘 Root `attachments/S5-evidence-matrix.md`，复用 F-V029 同一分母（行 = zh-CN/en-US × mvp/admin × 匿名/已认证；列 = 固定 UI / 12 pageId 并集 / M1～M4 / 权限正反 / 缺失翻译 / 配置刷新 / 错误回退）。
- 非 N/A 单元格均绑定可核对证据路径（vitest / Go 测试名、e2e spec、scratch 捕获）；N/A 仅 Profile 不可达并注明模块边界（`admin.settings`/`admin.activity` ∉ mvp）。

### C2 · 真实入口（已执行）

1. **API 构建 + 双启动**
   - `go build -o {SCRATCH}/s5-launch/server.exe ./cmd/server`（apps/api）exit 0。
   - 两次独立启动周期（独立 SQLite，`APP_PROFILE=admin`，`HTTP_ADDR=:18080`），各 `GET /api/branding`：
     - body 完全一致：
       `{"siteTitle":"Schema UI Core","logoUrl":"","logoUrlLight":"","logoUrlDark":"","faviconUrl":"","defaultLocale":"auto","supportedLocales":["zh-CN","en-US"],"siteTimezone":"auto","defaultTheme":"auto"}`
     - 字段健全性：siteTitle 非空；supportedLocales 含 zh-CN 与 en-US；defaultLocale/defaultTheme/siteTimezone 存在。
   - 捕获：`{SCRATCH}/s5-launch/run1.json`、`run2.json`、`compare.log`、`go-build.log`。

2. **Web 构建**
   - `npm run build`（apps/web）exit 0（tsc -b + vite build，~11s）。
   - 捕获：`{SCRATCH}/s5-launch/web-build.log`。

3. **Playwright e2e（admin）**
   - `npx playwright test e2e/localization.spec.ts`，`APP_PROFILE=admin`，`WEB_PORT=9999`。
   - **1 passed**（8.2s）：zh 切换、`document.documentElement.lang=zh-CN`、登录失败 envelope 协商、settings 保存投影、零 pageerror。
   - 截图：`apps/web/test-results/s5-settings-zh.png`（副本 `{SCRATCH}/s5-launch/s5-settings-zh.png`）。
   - 日志：`{SCRATCH}/s5-launch/e2e-localization.log`。

4. **矩阵引用测试刷新**
   - vitest 子集 79/79 全绿（i18n + startup-config + error-localization + locale-switcher）→ `{SCRATCH}/s5-tests/web-i18n-related.log`。
   - `go test ./internal/handler/ -run Branding|Settings|Localize|ErrorContract|Auth` 全绿 → `{SCRATCH}/s5-tests/api-handler-related.log`。

## 证据

| 主张 | 路径 / 命令 |
|------|-------------|
| S5 矩阵 | `GOAL-001-.../attachments/S5-evidence-matrix.md` |
| F-V029 分母 | `GOAL-001-.../attachments/F-V029-coverage-table-s0-freeze.md` |
| API dual-run | `{SCRATCH}/s5-launch/run{1,2}.json` + `compare.log` |
| Web build | `{SCRATCH}/s5-launch/web-build.log` |
| e2e | `apps/web/e2e/localization.spec.ts` + e2e log + screenshot |
| unit refresh | `{SCRATCH}/s5-tests/*.log` |

## 检查点

- [x] C1 完成
- [x] C2 完成
