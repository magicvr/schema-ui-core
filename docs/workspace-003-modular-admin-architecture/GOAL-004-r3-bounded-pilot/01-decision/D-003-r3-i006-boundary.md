---
id: D-003-r3-i006-boundary
doc: decision
goal: GOAL-004-r3-bounded-pilot
date: 2026-08-05
status: accepted-for-audit
---

# D-003 · I-006 静态入口、兼容和回滚边界

## 决定

将 API 聚合 Manifest 作为唯一生产运行时入口；Web `public` Manifest、API
嵌入式基线和页面 JSON 在 R3 实施期间只能作为源码/开发/测试输入，不能形成
生产静态兜底。Settings/Activity 的后端路由和 Schema 归属必须随 Profile 选择，
而 operationlog 的持久化和核心能力必须始终存在。

## 当前盘点结果

| 类别 | 当前入口 | C1 边界 |
|------|----------|---------|
| Web Manifest | `apps/web/public/.well-known/schema-ui/app-manifest.json` | R3 C4 前可作为开发/测试 fixture；不得由最终生产镜像提供 |
| API Manifest 基线 | `apps/api/internal/manifest/app-manifest.json`、`manifest.Default` | 保留为构建源码基线；生产响应必须来自 `ForModules` 聚合结果 |
| API Manifest 路由 | `handler.RegisterManifest`、`/.well-known/schema-ui/app-manifest.json` | 保留；由同一 `kernel.Plan` 决定页面和导航 |
| API Schema | 当前 `handler/fixtures/schema/*.json` | Settings/Activity 在 C2 迁入模块包；其他页面不在 R3 中擅自迁移 |
| Vite 代理 | Manifest 精确路径和 `/api` 代理 | 保留为开发期同源链路，不把失败降级到静态文件 |
| Nginx/Docker | 精确代理；构建后删除 `dist` Manifest | 保留并增加可核对的最终产物检查 |
| 中心业务注册 | `handler.Register` 固定挂载 Settings/Activity/Operations | C2 移除试点模块的固定硬编码，由组合根按已解析计划应用模块贡献 |
| Host/Shell 事件 | `BRANDING_CHANGED_EVENT` 在 Host 触发、App/LoginPage 消费 | C2 泛化为模块事件贡献；在 V3 前不得把当前实现误称为通用契约 |

## 兼容窗口与告警

兼容窗口以 R3 C4 D 门为边界：在 C4 通过前，开发/测试可以保留 Web
`public` fixture，便于协议测试和无 API 的静态检查；从最终生产镜像中删除该
文件的规则持续有效。开发链路若命中静态 fixture，必须产生明确告警并记录命中
路径；生产 Nginx/API 链路不得静默回退。C1 当前只冻结该要求，告警实现和同构建
验证仍是 C2/C3 的待取证项。

## 移除触发

满足以下全部条件后才允许进入 R6 的旧路径移除：

1. C2 的四类病灶均有代码证据；
2. C3 在同一 Web 构建上通过 V-1～V-4、冲突 fail-closed、Settings 事件和数据
   保留检查；
3. C4 self 与 Grok independent 审计均无开放 required finding；
4. R6 另行记录生产镜像、容器代理、双 Profile、升级/恢复和完整回归证据。

任一条件未满足时，只能保留边界内的开发/测试输入，不得声称移除完成，也不得
冻结 R4 的全量迁移方案。

## 回滚触发与恢复

以下任一情况触发停止推进并回滚到上一份 API/Web 构建和上一份 Profile 配置：

- 启用模块的 Manifest、Schema 或业务路由返回错误/缺失；
- disabled Profile 仍能访问 Settings/Activity 路由或菜单；
- 同一 Web 构建对不同 Profile 产生错误页面集合；
- Settings 事件不刷新 Host/Shell，或 generic Renderer 出现试点专用分支；
- operationlog 不再可写/可读，或禁用模块导致系统数据被删除或用户字段改变。

恢复时只回退应用构建和启用配置，保留 SQLite 文件、迁移表、operation_log 和
用户字段；恢复后必须依次核验 `readyz`、Manifest、启用/禁用路由、Schema、
operationlog 读写和数据计数/关键字段。C1 的演练证据仍待 C2/C3 取证。

## 非结论

本记录冻结的是 R3 实施边界，不宣称四类病灶已切除、兼容告警已实现、回滚已
演练，也不把本记录当作 R3 D 门或 VP 退出证据。
