---
id: GOAL-001-wallet-prepaid-instrument
doc: execution-entry
record_id: E-007
status: recorded
parent: null
created: 2026-09-02
updated: 2026-09-02
version: 0.1.0
---

# E-007 · 预付凭证页面修复：Schema 校验、导航图标与排序（2026-09-02）

## 用户报告

1. 预付凭证页面显示「页面 Schema 错误」；
2. 给预付凭证导航菜单合适的图标；
3. 注册排序移到「钱包」下一位。

## 根因与修复事实

### 1. 页面 Schema 错误 → **fixed**（声明载体修正，A-008 F-002 实现细节修订）
- **根因**：A-008（F-002）把 CSV 导出声明 `downloadCsv` 放进 `generateBatch.onSuccess`。但 `onSuccess` 受 **upstream pin** `docs/schemas/action.schema.json` 的 `OutcomeBehavior` 严格结构约束（`additionalProperties: false`，provenance-v2.9 锁定 11 件 docs/schemas 全量 sha256），页面文档在运行时 D-VAL 失败 → `PAGE_SCHEMA_INVALID` → 壳层显示「页面 Schema 错误」。**禁止改 pinned schema**（workspace-011 GOAL-004 E-003 先例：download 类扩展不得走行为扩展，必须用协议允许的扩展点）。
- **修正**：声明改放**生成模态表单节点的业务 props**（`openGenerate.content.props.downloadCsv`）——node schema 的 `props` 是上游明确放开的「业务级配置参数」区（仅禁 CSS 名），文档重新通过 D-VAL：
  - `modules/wallet/schema/wallet-vouchers.json`：`generateBatch.onSuccess` 还原为 `{"behavior":"reload"}`；`openGenerate` 表单 props 增 `downloadCsv`（columns/fileName 同前）；
  - `apps/web/src/renderer/render.tsx` `submitForm` 改为读取**表单节点 props** 的声明（不再读 action.onSuccess）；语义不变：声明存在 + 列含 `code` + 响应 `items[].code` 才同手势导出；
  - `apps/web/src/renderer/render.types.ts`：`RenderFormNode.props` 增 `downloadCsv?`，`parseRenderNode` form 投影白名单纳入该业务键（此前 form props 只投影已知键，未知键会被丢弃——这也是仅改 JSON 不够的原因）。
- **回归守卫**（新增 `apps/web/src/protocol/conformance/wallet-vouchers-downloadcsv-schema.test.ts`）：对**真实模块文档**跑 D-VAL 断言通过；反向断言把 `downloadCsv` 塞回 `action.onSuccess` 必须失败（防回归到 pin 违规载体）。

### 2. 导航图标 → **fixed**（`card` → Ticket）
- 根因：wallet manifest fragment 的 voucher 导航项语义图标名为 `"card"`（`modules/wallet/manifest/fragment.json` sidebar），但壳层 `apps/web/src/app/App.tsx` 的 `iconRegistry` 没有 `card` → `iconFor` 返回 null，导航项无图标。
- 修正：`iconRegistry` 增 `card: Ticket`（lucide Ticket 图标，语义贴合卡密凭证），import 同步。

### 3. 排序移到「钱包」下一位 → **fixed**
- 根因：`menu_wallet_vouchers` 不在 `kernel.DefaultNavigationOrder` → 合并/发布排序回退 legacy Order，落到 Data permission 之后。
- 修正：`kernel/provider.go` `DefaultNavigationOrder` 在 `menu_wallet` 之后插入 `menu_wallet_vouchers`（注释 VP-029 R3）；同步快照 `navigation_order_test.go` 与 composition 冻结侧栏列表（Prepaid vouchers 移到 Wallet 之后）。

## 验证（exit 0）

```text
apps/api: go test ./kernel ./internal/composition ./internal/manifest ./internal/handler ./modules/wallet/... → 全 ok
apps/web: npx vitest run src → 91 files / 1191 tests PASS
  含：wallet-vouchers 文档 D-VAL 守卫（真实模块 JSON · onSuccess 反向断言）、
      renderer CSV 声明化用例（表单 props 载体 · 44/44）
```

未触碰 11 件 pinned docs/schemas（sha256 与 provenance-v2.9 保持一致）；JSON 均可解析。

## 产物

- `apps/api/modules/wallet/schema/wallet-vouchers.json`（声明载体表单 props）
- `apps/web/src/renderer/render.tsx` / `render.types.ts` / `render.test.tsx`
- `apps/web/src/app/App.tsx`（iconRegistry card→Ticket）
- `apps/web/src/protocol/conformance/wallet-vouchers-downloadcsv-schema.test.ts`（新增守卫）
- `apps/api/kernel/provider.go` + `kernel/navigation_order_test.go` + `internal/composition/composition_test.go`（排序）
