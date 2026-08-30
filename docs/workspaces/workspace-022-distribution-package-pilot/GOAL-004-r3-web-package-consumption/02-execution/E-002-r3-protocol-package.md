---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-004-r3-web-package-consumption
version: 0.1.0
---

# E-002 · protocol 首包与消费闭环（2026-08-29）

## 打包链路（B 路径 · Vite lib）

- `apps/web/vite.lib.config.ts`：lib mode 单产物构建（entry `src/protocol/index.ts`；alias `@`/`@schemas`；无 react/tailwind 插件）。
- `apps/web/tsconfig.protocol.json`：declaration-only 编译（rootDir=src · DOM.Iterable · node/vite client types · exclude tests/upstream）。
- 产物（`.gitignore` 排除，可重建）：`apps/web/dist-lib/@schema-ui/protocol/` = `index.js`（306 kB · 144 modules · ajv/schema JSON 全部 bundle 自包含）+ `protocol/*.d.ts`（index/app-manifest/load-page/conformance 14）+ `lib/fetch-timeout.d.ts`；附 package.json（`exports` 指向 import/types）。
- 冗余 note：`theme-init.js` + `conformance-claim*.json` 为 `public/` 静态资源被 Vite 拷贝入产物（无害；正式发布流水线时过滤）。

## golden-web 消费验证

- `attachments/golden-web/`：独立 npm 项目（无源码依赖），`package.json` 仅 `@schema-ui/protocol: file:` 指向产物。
- `pnpm install` ✓（file: 链接解析）→ `node probe.mjs`：

```
golden-web protocol probe PASS · 2.9
```

- 断言：`APP_MANIFEST_PROTOCOL_VERSION=2.9`、支持窗含 2.8、`isValidExpression("$context.user.role == \"admin\"")`、`stripPathQuery`、`resolveSchemaUrl(base,schemaUrl,params)` —— 均为**包内功能行为**验证（非主仓代码路径）。

## 结论

判据 #2 前置句成立：npm 包可发布形态 + 下游仅包依赖消费 **PASS**。剩余 = S3 渲染闭环（renderer/shell/ui 包化与页面渲染 + Token 覆盖）——留下一轮（renderer 拆包需处理 React/Tailwind/CSS 产物面，属 B 路径内的增量）。

## 关联

- I-002：包边界设计 v0.1 已生效（六包 + peer 矩阵）；protocol 首包实证 → 设计高置信，renderer/shell 包化沿用同一链路。
- I-007（协议 pin 漂移）：登记 + 用户裁决「留 /vision 后续裁决」——本目标沿用代码事实（支持窗 2.7–2.9），不自行推翻。