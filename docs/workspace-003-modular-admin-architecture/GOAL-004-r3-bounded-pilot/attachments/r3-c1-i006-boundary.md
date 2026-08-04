---
title: R3 C1 · I-006 入口与回滚边界证据
date: 2026-08-05
status: working-evidence
---

# R3 C1 · I-006 入口与回滚边界证据

## 证据范围

本附件记录源码盘点和决策边界，不把静态扫描当作生产运行、CI、容器、发布或
R3 通过证据。文件行号以当前工作树为准；实现修改后必须重新取证。

## Manifest 与生产链路

- Web 静态 Manifest：`apps/web/public/.well-known/schema-ui/app-manifest.json:1-64`，
  被协议测试和代表性页面测试读取。
- API Manifest 源：`apps/api/internal/manifest/app-manifest.json:1-143`；
  `apps/api/internal/manifest/manifest.go:12-13` 通过 `go:embed` 纳入。
- `manifest.Default`/`ForModules`：`apps/api/internal/manifest/manifest.go:34-42`；
  其中 `ForModules` 根据已解析模块过滤页面/导航并聚合。
- 运行时路由：`apps/api/internal/handler/manifest.go:16-42` 注册公共 Manifest，
  `apps/api/internal/composition/composition.go:87-100` 用 `plan.IDs()` 产生聚合。
- 开发代理：`apps/web/vite.config.ts:14-29`；容器代理：
  `apps/web/nginx.conf:10-32`；Docker 最终层：`apps/web/Dockerfile:14-18`。
  最终 `dist` 中的静态 Manifest 会被删除，Nginx 对精确路径转发 API。

## Schema 与路由入口

- 当前全局 Schema embed：`apps/api/internal/handler/schema.go:9-16`；运行时读取
  `fixtures/schema/*.json` 并在 `:24-68` 提供 `/api/schema/{pageId}`。
- `activity.json`、`settings.json` 仍在
  `apps/api/internal/handler/fixtures/schema/`，这是 C2 需要迁入模块包的试点
  所有权病灶。
- 中心注册：`apps/api/internal/handler/health.go:27-36` 固定调用
  `registerOperations`、`settingsHandler` 和 `schemasHandler`，未使用 Profile
  决定是否挂载后两个业务面。
- Profile：`apps/api/internal/kernel/profile.go:24-47,91-107` 表示 MVP 只启用
  `core.operationlog`，Admin 另启用 `admin.settings` 和 `admin.activity`。

## Host/Shell 边界

- Host 触发：`apps/web/src/main.tsx:29-52` 在 Settings PATCH 成功后调用
  `notifyBrandingChanged`。
- 事件定义：`apps/web/src/app/branding.ts:13-19`。
- 消费：`apps/web/src/app/App.tsx:391-412` 和
  `apps/web/src/app/LoginPage.tsx:29-51` 重新读取 branding。
- 通用 Renderer 明确不承载该职责：`apps/web/src/renderer/render.tsx:266-268`。
- 当前没有该事件的独立自动化测试；因此 C2/C3 必须补充事件贡献和运行证据。

## 当前缺口

1. C1 之前没有固定的兼容截止/触发记录；D-003 现已规定以 R3 C4 D 门为开发
   fixture 窗口边界，但告警行为尚未实现。
2. 当前没有 disabled Profile 的真实路由清单、重启后菜单/Manifest 对照。
3. 当前没有回滚演练、数据计数/关键字段保留核对或恢复后 `readyz`/Manifest/
   operationlog 证据。

## C1 结论

入口盘点已完成，边界已形成，但 I-006 仍为 `collecting`，直到独立审计和后续
实现/演练证据闭合 A-001 的 required findings。
