---
title: E-004 · 关门验证修复——shell.spec 匿名 schema 探测契约对齐 F-010
status: active
created: 2026-08-26
updated: 2026-08-26
parent: null
version: 1.0.0
---

# E-004 · 关门验证修复——shell.spec 匿名 schema 探测契约对齐 F-010

## 起因（关门验证发现）

R-001 并入后的**全量 e2e 套件**首轮出现 1 例失败：`shell.spec.ts` 「login gates the shell…」在匿名探测断言处 `Expected: 404 / Received: 401`（`/api/schema/settings|activity`，mvp 档）。

## 定性：既有陈旧契约，非本波回归

- 本波改动范围仅 web 装配/测试/文档，未触碰任何 API 代码；
- 失败语义 = W13 F-010（`b7954235`）「/api/schema 挂认证」的真实后效：schema 路由**全局注册 + 鉴权中间件前置**，匿名请求在任何 profile 下都先吃 401，模块存在性判定在其后；
- 该断言仍是 F-010 之前的旧契约（admin 档期望匿名 200、mvp 档期望 404）；
- 根目录 `e2e-baseline.log` 显示上一次全量 e2e 尝试因 postgres 方言配置 **API 启动即失败**（LIFECYCLE_START_FAILED），即 F-010 落地后该套件从未完整跑通过——陈旧断言因此一直未被暴露。

## 修复

`apps/web/e2e/shell.spec.ts` 两行期望更新并加注：

```ts
expect(settingsSchema.status()).toBe(401);
expect(activitySchema.status()).toBe(401);
```

数据路由对照断言（`/api/settings`、`/api/operations` 的 `isAdminProfile ? 401 : 404`）保持不变——它们按模块注册，行为与 F-010 无冲突且当前通过。spec 内注释指回本条 E-004。

## 验证

修复后全量 e2e 套件复跑结果见 [D-002](../../01-decision/D-002-w14-closeout.md) 关门记录。
