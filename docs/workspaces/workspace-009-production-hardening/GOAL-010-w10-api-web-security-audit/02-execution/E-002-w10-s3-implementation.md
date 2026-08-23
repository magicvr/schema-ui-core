---
id: E-002-w10-s3-implementation
goal: GOAL-010-w10-api-web-security-audit
status: done
created: 2026-08-21
updated: 2026-08-21
parent: GOAL-010-w10-api-web-security-audit
version: 0.1.0
---

# E-002 · W10 S3 实施：3 条修复 + 回归全绿（2026-08-21）

## 事实

### F-001 · env.example 硬编码凭据清除（required HIGH · fixed）

1. `apps/api/configs/env.example`：`DB_PASSWORD`/`DB_DSN`/`PG_TEST_PASSWORD` 中的真实密码与内网 IP 全部替换为占位符（`replace-with-a-strong-password` / `your-db-host.example.internal` 等），并加"NEVER commit real credentials"注记；`sslmode` 示例从 `disable` 收紧为 `require`。
2. 全仓扫描（`Ss.110110|192.168.31.213`）追加发现 `apps/api/internal/config/config_test.go:827-847` 测试夹具含同一真实凭据（测试文件入库）——同波修复：替换为 `test-only-password`/`db.example.internal`/`appdb`/`appuser` 假值，断言同步更新。
3. 本波 A-001 及附件中的明文值加脱敏注记（`<leaked-password>`），不再于版本控制文件复现秘密。
4. 残余项（移交用户）：W9/W13 历史台账以证据形式记录过该密码（闭合记录不改写）；**该数据库环境的密码应视为已泄露并轮换**。

### F-002 · Web fetch 统一超时（recommended · fixed）

1. 新增 `apps/web/src/lib/fetch-timeout.ts`：`withTimeout(fetchImpl?, timeoutMs=30_000)` 返回 fetch 兼容包装——超时 abort、与调用方 signal 组合、已中止 signal 快速失败（不发起请求）、settle 后清 timer、**调用时**解析全局 fetch（兼容测试 stub）。
2. 接线（默认生产路径）：`auth-client.ts`（postJSON + authFetch 首次与 401 重试，覆盖 login/refresh/logout/mfa-verify 及全部 authFetch 消费方：main.tsx resourceFetcher、config-events、mfa-manager、wallet-ensure、force-password-change、cron-preview、account-session-toolbar、import-template-download、fetchMe）；`protocol/load-page.ts`（页面 schema 默认 fetcher；注入 fetcher 不变保测试确定性）；`renderer/form-controls.tsx`（动态选项源默认 fetcher）。
3. 新增 `src/lib/fetch-timeout.test.ts`（5 用例：透传/超时中止/signal 组合/已中止快速失败/默认 30s）。

### F-007 · 选项源 URL 反斜杠缺口（recommended · fixed）

1. `renderer/form-controls.tsx:78` 字符类 `[^\s\#]` → `[^\s\\\#]`：封堵 WHATWG URL 将 `\` 规范化为 `/` 导致 `/\host` 变协议相对 `//host` 逃逸同源的缺口；注释说明载荷理由。
2. 全仓核查 11 处同型正则：其余 10 处（app-manifest/branding/resource/request-construction/upload-orchestration/render/form-controls.ts）均已含 `\\` 排除，无需改动。
3. `form-controls.dynamic-options.test.tsx` 新增反斜杠回归用例（`/\evil.example`、`/api\evil` 等拒绝且不发起 fetch）。

### 误报调和（F-003/F-004/F-005/F-006 → 作废）

逐条源码核实不成立，证据见 [D-003](../01-decision/D-003-w10-scope-reconciliation.md)。未作任何"假修复"。

### 回归验证

- API：`go vet ./...` 0 告警；`go test ./...` 全绿（含改动包 config/handler/authsession 复跑 ok）。
- Web：`npx vitest run` **76 files / 1083 tests 全绿**（含新增 6 用例）；`npm run build`（tsc -b + vite build）通过。

## 产物

- `apps/api/configs/env.example`、`apps/api/internal/config/config_test.go`
- `apps/web/src/lib/fetch-timeout.ts`（+test）、`apps/web/src/account/auth-client.ts`、`apps/web/src/protocol/load-page.ts`、`apps/web/src/renderer/form-controls.tsx`（+dynamic-options.test 回归）

## 后续（计划，非事实）

- S4：independent 复核（工作区惯例 grok provider，本会话不可用——I-003 待用户裁决：追加 grok `/audit` 或接受本 self 审计关门）。
- go 宣称恢复：待 S4 后用户书面裁决（对齐 W9 D-004）。