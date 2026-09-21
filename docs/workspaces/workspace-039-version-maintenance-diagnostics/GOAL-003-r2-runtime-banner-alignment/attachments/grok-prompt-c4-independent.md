# Grok Build · 独立交叉审计（GOAL-003 R2）

仓库根 `/audit`。grok-4.6 · reasoning high。`source: independent`。不改 status/progress/goal-tree。

## 目标

`[workspace-039-version-maintenance-diagnostics] GOAL-003-r2-runtime-banner-alignment`

scope：T-1 bootstrap 生产者、T-2 `/me`+`fetchMe`、T-3 Shell 横幅。对照 GOAL-002 D-001。

## 必读

1. GOAL-003 `00-meta`、`D-001-r2-implementation-scope`、`E-002-r2-implementation`、`A-001-r2-self`
2. GOAL-002 `D-001-r1-contract-and-denominator-freeze` 与两矩阵
3. 代码：`handler/bootstrap.go`、`bootstrap_test.go`、`account.go`、`session.go`、`health.go` `RegisterWithMFAProbes`、`composition.go` 接线、`auth-client.ts` `fetchMe`/`parseRuntimeMode`、`runtime-banner.tsx`、`App.tsx`、`AuthGate.tsx`
4. 确认未改 `apps/web/src/protocol/upstream/**`、`operational.go` 写门禁

## 核验

- maintenance Host 文档现为 `degraded`；消费者仍对 `mode=maintenance` 终态
- `/me` 精确 runtimeMode；fetchMe 不丢字段
- 横幅三分文案；normal 不渲染；登录页无 App 故无横幅
- T-4 版本 chip 未混入

## 输出

verdict + Findings。可写 `03-audit/A-002-*.md` 并更新索引。
