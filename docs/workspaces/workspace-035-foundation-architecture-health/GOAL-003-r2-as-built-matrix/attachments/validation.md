---
status: active
created: 2026-09-09
updated: 2026-09-09
parent: GOAL-001-foundation-architecture-health
version: 0.1.0
---

# R2 本次验证

2026-09-09；Windows；代码基线 ebe6013c；生产源码未改。GOCACHE 使用 TEMP/schema-ui-w035-gocache。

| 命令（在对应 app 执行） | 结果 | 原始输出 |
|---|---|---|
| `go test ./kernel ./internal/cache ./internal/ratelimit ./internal/eventbus ./internal/objectstore ./internal/mail ./internal/obs ./internal/auth ./internal/jobs ./internal/store ./internal/manifest ./internal/composition ./assembly -count=1 -timeout=180s` | exit 0；12 包通过；assembly 无测试 | [Go](validation-go.txt) |
| `go test ./internal/channel/telegram ./internal/config ./internal/server -count=1 -timeout=180s` | exit 0；3 包通过 | [补充](validation-extra.txt) |
| `npm test -- src/app/navigation.test.ts src/protocol/app-manifest.test.ts src/protocol/load-page.test.ts src/app/App.integration.test.tsx` | exit 0；4 files / 48 tests passed | [Web](validation-web.txt) |

## 证据限制

本次不是 -race、全仓测试、浏览器 E2E 或生产 smoke。未配置/启动本任务专属 PG、S3、Telegram、Resend 外部服务；不把受环境跳过或 fake 服务用例计为 live 供应商证据。Go 简短输出仅证明包命令通过，不能据此统计所有子用例是否执行。PG drain/recovery 在源码显式以 PG_TEST_* 门控；Linux 进程信号/Compose stop 未重跑。

这些是本次有界架构对照的证据界限，不是本轮接受新的生产残余。
