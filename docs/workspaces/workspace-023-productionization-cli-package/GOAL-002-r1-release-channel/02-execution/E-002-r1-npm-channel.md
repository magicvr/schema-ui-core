---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-002-r1-release-channel
version: 0.1.0
---

# E-002 · S2 npm 通道实证（2026-08-29）

## 发布（GitHub Packages · 凭据 = 用户本地 `.env` `github_token`）

- `scripts/publish-npm-packages.mjs`（含临时 .npmrc 认证注入；Windows `rm`/ESM require 缺陷已修）：dry-run 认证验证 → 真实发布：
  - `@magicvr/schema-ui-protocol@0.2.0`（73 KB tgz）
  - `@magicvr/schema-ui-renderer@0.1.0`（99 KB tgz）
- 包名 scope = `@magicvr`（GitHub Packages owner 匹配约束）；来源 = `apps/web/dist-lib/artifacts`（pack 一键）。

## 消费（golden-field/web · 真实 registry 语义）

- `package.json` 依赖改版本号（无 file:）；`.npmrc` 指向 `@magicvr:registry=https://npm.pkg.github.com`；认证 = **用户级 ~/.npmrc**（pnpm 安全策略拒绝项目级 env 展开——知识项）；
- `pnpm install` → **从 registry 下载**：

```
@magicvr/schema-ui-protocol@0.2.0
@magicvr/schema-ui-renderer@0.1.0
lockfile: tarball: https://npm.pkg.github.com/download/@magicvr/schema-ui-protocol/0.2.0/8b8565e… + integrity sha512-Y0L15… 
```

- 三探针全绿（protocol/render/token）。

## 知识项（产线流程追加）

1. pnpm：项目级 .npmrc 的 `${NODE_AUTH_TOKEN}` **不展开**（防泄露）；凭据必须 user-level `pnpm config set --location=user`。
2. GH Packages 包名 scope 必须 = owner（`@magicvr`），registry = `npm.pkg.github.com`。
3. 发布脚本幂等（同名同版本重发会 403，需 bump）。

## 状态

- I-023-002 **verified**；占位依赖清零：golden-field **无 file:、无 tarball 路径、锁文件指向 registry**（web 侧完成；Go 侧已自 S1 起无 replace）。
- S3（registry 升级演练）主体 = 绑定下次真实发布（R2 CLI 里程碑或六包发布时执行 bump→安装→回归），被本 S2 的「无占位安装」覆盖前置条件。