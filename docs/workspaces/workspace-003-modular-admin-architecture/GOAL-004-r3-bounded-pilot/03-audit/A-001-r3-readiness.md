---
id: A-001-r3-readiness
doc: audit
goal: GOAL-004-r3-bounded-pilot
source: self
date: 2026-08-05
scope: R3 建立、C1 I-006 入口盘点、兼容/告警/移除/回滚边界
verdict: conditional
---

# A-001 · R3 C1 自审

## 结论

`conditional`。源码入口、生产代理边界、Profile 模块集合和 Host/Shell 事件
职责已经形成可追踪事实；但当前代码仍固定挂载试点路由，且没有开发期静态
fixture 告警、Profile 路由禁用、回滚演练或同构建 V-1～V-4 证据。

## 已核对事实

- `apps/api/internal/composition/composition.go` 按 `kernel.Plan` 调用
  `manifest.ForModules` 和 `handler.RegisterManifest`。
- `apps/web/Dockerfile` 在构建后删除 `dist/.well-known/schema-ui/app-manifest.json`；
  `apps/web/nginx.conf` 对同一路径使用精确 API 代理；Vite 同时代理 Manifest 和
  `/api`。
- `apps/api/internal/manifest/manifest.go` 将嵌入基线标注为源码 fixture，
  `ForModules` 负责发布聚合。
- `apps/api/internal/handler/health.go` 的 `Register` 仍固定注册
  `registerOperations`、`settingsHandler` 和 `schemasHandler`。
- `BRANDING_CHANGED_EVENT` 在 Host 触发、App/LoginPage 消费，通用 Renderer
  没有试点专用分支；但缺少独立自动化测试。

## Findings

### F-C1-001 · I-006 的边界尚未通过独立核验

- level: `required`
- status: `open`
- impact: C1 关闭和 R3 方案冻结
- finding: D-003 已列出保留/移除边界，但尚未由 Grok independent 审计确认，
  也未在同一构建上核对模块禁用与生产无静态兜底。
- closure: 需要独立审计意见，并补上可复核的实现/验证证据；否则保持 `open`。

### F-C1-002 · 告警和回滚尚无演练证据

- level: `required`
- status: `open`
- impact: I-006-02、I-006-03 和后续 R6 移除
- finding: D-003 只规定了开发告警、回滚触发和数据保留要求；当前没有告警
  行为测试、模块禁用演练、数据保留核对或恢复后端点证据。
- closure: 在 C2/C3 补齐可运行的告警/禁用/回滚演练，并由 self 与 independent
  审计核验后，才可使用 `fixed` 关闭。

## 阶段结论

C1 暂不通过。继续范围只能是 I-006 证据补齐和有界实现准备；不得把本条意见
解释为 R3 失败，也不得跳过 C1 直接进入 R4。
