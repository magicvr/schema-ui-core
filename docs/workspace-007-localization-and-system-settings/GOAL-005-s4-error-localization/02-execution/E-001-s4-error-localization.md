---
id: E-001
doc: execution
title: S4 · 稳定错误码 + 有界服务端 locale 协商 + 前端本地化保底
status: recorded
parent: GOAL-005-s4-error-localization
created: 2026-08-09
updated: 2026-08-09
version: 0.1.0
---

# E-001 · S4 后端用户可见反馈本地化（2026-08-09）

## 事实

- **共享编目包** `internal/errorcatalog`：43 个已编目码（D-002 附录 A 冻结集；`INTERNAL` 刻意缺席）；每码 {messageKey, En, Zh}；`NegotiateLocale`（Accept-Language 解析：zh-CN/zh → zh-CN，en 前缀 → en-US，无匹配 → en-US 回退）；`Body` 组装 `{error, message, messageKey}` envelope。
- **Go 写路径**：`handler/localize.go` `writeLocalizedError`（编目 → 本地化 message + Content-Language 头；未编目 → 英文通用 + 无 messageKey）；58 处 `writeError` 调用点转换为 `writeLocalizedError(w, r, ...)`（auth/account/schema/settings/resources/upload）；`requirePermission`/`writeEntityError`/`writeSettingsError` 改签名携带 `r`；`internal/auth` 中间件错误本地化（auth 包独立 writeLocalizedError 共享 errorcatalog）。
- **错误码契约测试** `error_contract_test.go`：扫描 handler + auth 源码字面量，断言 = 冻结集（34 字面量 + 10 域码），防漂移；编目覆盖/缺席（INTERNAL 不编目）断言。
- **协商测试** `localize_test.go`：zh/en/前缀/回退/Content-Language；INTERNAL 英文无 key 无 params；端到端：登录失败（UNAUTHORIZED zh）、设置校验（INVALID_TIMEZONE zh + messageKey）、中间件（UNAUTHENTICATED zh）。
- **Web（C4 保底）**：`getActiveLocale`/`setActiveLocale` 注册表（provider 同步）；`auth-client` 全部请求携带 `Accept-Language: <activeLocale>`；`readResourceApiError` 解析 `messageKey`/`params`（`ResourceApiError` 扩展）；`SchemaCrudFeedback`/`ActionResult` 透传；`FeedbackRegion`/表单错误呈现优先 catalog（按 key/params，当前语种 → en-US），未知 key 回退服务端 message（不渲染裸 key）；web catalog 镜像 errorcatalog 43 键（脚本从 Go 单一来源生成，双 catalog 键集对称）。
- **验证**：Go 全包通过（新增 12 项本地化/契约测试）；vitest **711/711**（39 文件，新增 5 项 S4 测试）；`npm run build` exit 0；证据捕获 `{SCRATCH}/unit-s4-web.log`、`{SCRATCH}/unit-s4-api.log`。

## 产物

| 路径 | 说明 |
|------|------|
| `apps/api/internal/errorcatalog/` | 共享编目 + 协商（C2） |
| `apps/api/internal/handler/{localize.go,localize_test.go,error_contract_test.go}` | 写路径 + 契约（C1/C2/C3） |
| `apps/api/internal/{auth,handler}/*.go` 58 处调用点 | 编目错误本地化（C3） |
| `apps/web/src/i18n/runtime.tsx`、`account/auth-client.ts` | Accept-Language 携带（C4） |
| `apps/web/src/renderer/{resource.ts,render.tsx}` | envelope 解析 + 呈现保底（C4） |
| `apps/web/src/renderer/error-localization.test.tsx` | C4 测试 5 项 |

## 里程碑 checkpoint

- commit：`<S4 checkpoint hash 待填>`（owned paths = 上述全部路径 + 本目标治理文档，显式 `git add` 无 `-A`）。
