---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-003-r2-public-release-channel
version: 0.1.0
---

# E-002 · S2/S3 发布与无凭据消费实证（2026-08-29）

## 交付物

- `scripts/publish-npmjs-packages.mjs`：npmjs 公开发布脚本（读仓库根 `.env` 的 `npm_token` → 临时 `.npmrc`，token 不进日志/不入库；`--access public` + `publishConfig.access=public`；已发布版本幂等跳过（`npm view versions` 预检 + 403「cannot publish over」兜底）；PUBLISH_SCOPE 参数化）。
- `apps/api/v0.4.0` origin tag（serves 面 = server 包 + CLI serve + 新模板）。
- golden-field 升级：`go.mod` v0.3.0→v0.4.0（公共 proxy 免凭据下载 · go.sum 哈希）；`cmd/server` 换 thin wrapper（serve 面）；web 六包 → npmjs 公开消费形态（移除 GH Packages 映射）。

## 发布实证（S2 · C1/C2）

| 项 | 证据 | 结果 |
|----|------|------|
| npmjs 公开六包 | `@magicvr/schema-ui-{lib@0.1.0, protocol@0.2.0, renderer@0.2.0, shell@0.1.0, theme@0.1.0, ui@0.1.0}` 真实上传（lib/protocol 首轮落档，其余第二轮）；`npm view` 六包全部可见 | ✅ |
| 2FA 门禁处置 | 首次 E403（需 bypass-2FA token）→ 用户换 granular token → 通过；随后 E402（scoped 默认私有）→ `--access public` 修复 | ✅（知识项落脚本注释） |
| Go v0.4.0 发布 | origin tag 推送（`00d97b5b`）· golden-field `go mod tidy` 公共 proxy 下载成功 · go.sum 含 v0.4.0 哈希 | ✅ |

## 无凭据消费实证（S3 · C3）

| 项 | 证据 | 结果 |
|----|------|------|
| Go 消费 | golden-field `go mod tidy`（默认 GOPROXY，无认证）→ `go build ./cmd/server` OK → thin wrapper serve（sqlite）healthz 200 · login 200 | ✅ |
| Web 消费 | 清旧依赖后 `pnpm install`（`NPM_CONFIG_USERCONFIG=空` + 项目无 registry 映射）→ 六包 + react/react-dom 安装成功（864ms · lockfile 落盘）· `pnpm ls` 六包齐 | ✅ |
| 探针 | `probe.mjs`（protocol · 2.9）PASS · `probe-render.mjs`（html 1573B）PASS · `token-check.mjs`（brand=2 ⊆ index=5）PASS | ✅ |
| 无 token 断言 | 安装全程无 GH Packages 认证；用户级映射被空 userconfig 屏蔽 | ✅ |

## 发布流程成文（C4）

- 脚本 + 凭据注入点（`.env` `npm_token`，`*.env` 已在 .gitignore）；
- scope 迁移 changelog 注记：D-001 §6（`@magicvr` npmjs 公开先行 · `@schema-ui` = org 就绪后正式化候选；消费迁移 = 移除 GH Packages registry 映射）；
- 知识项：scoped 包默认私有 → `--access public`；404 传播延迟 → 403 幂等兜底；2FA → granular bypass token。

## 判据映射

| 判据 | 证据 | 结论 |
|------|------|------|
| C1（npmjs 真实发布六包） | 发布日志 + npm view 可见性 | ✅ |
| C2（Go v0.4.0 tag + proxy go get） | tag 推送 + go.mod/go.sum + tidy/build | ✅ |
| C3（golden-field 无凭据消费 + 探针） | pnpm install(空 userconfig) + 三探针 | ✅ |
| C4（流程成文 + scope 迁移注记） | 脚本/README 注记/D-001 §6 | ✅ |

## 残余与登记

1. `@schema-ui` org scope：正式化候选（D-001 §6；触发 = org 创建 + 用户指令，届时新包名发布 + 消费方迁移）。
2. 既有 GH Packages `@magicvr/schema-ui-*`（私有）：保留不删（历史消费面）；新消费一律指向 npmjs 公开版。