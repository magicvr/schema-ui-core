---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-003-r2-public-release-channel
version: 0.1.0
---

# E-003 · A-002 响应与修复（2026-08-29 · /govern）

## F-001（required / high）· lockfile 仍钉 GH Packages —— **fixed**

**根因**（A-002 证据链确认后定位）：pnpm v11 **无视 `NPM_CONFIG_USERCONFIG`**，仍读取用户级 `~/.npmrc` 的 `@magicvr:registry=https://npm.pkg.github.com`（VP-023 时代 `pnpm config set --location=user` 遗留）。因此「删 lockfile + 空 userconfig」的重装实际仍按用户级映射解析 GH Packages metadata，并用用户级 GH token 下载成功——install 看似成功、lock 却写 GH tarball（401 无凭据不可复现）。

**修复**（golden-field commit `fb957a9`）：
1. 项目级 `web/.npmrc` **正向钉死** `@magicvr:registry=https://registry.npmjs.org`（项目级覆盖用户级；公开源，无 token）。
2. 删除 lock/node_modules 后以**全新空 store**（`--store-dir` 临时目录）重装：```GH tarball 残留 = 0```；lockfile 六包 `resolution.tarball` 全部 `registry.npmjs.org`；lib integrity `sha512-fNhAMTnBS+FEX…` 与 `npm view … dist.integrity` **一致**。
3. 三探针基线不受影响（包内容同 npmjs 版；grok 独立隔离安装已复跑全绿）。
4. **知识项**：pnpm 的 scope 映射优先级 = 项目级 > 用户级；NPM_CONFIG_USERCONFIG 仅影响 npm 客户端，pnpm 不适用——公开消费必须项目级正向钉 scope registry。

## F-002（recommended）· 脚本默认 scope —— **fixed**

`scripts/publish-npmjs-packages.mjs` 默认 `PUBLISH_SCOPE` 改 `@magicvr`；头注释与 D-001 §6 对齐（`@schema-ui` 保留为 org 就绪后候选覆写点）。

## F-003（recommended）· 台账措辞/索引 —— **fixed**

- `00-meta` title/概述/I-024-001 结论句 → `@magicvr` 先行（§6 记录历史）；
- D-001 决策条 1 增注「§6 变更」指针（不改历史原文，追加变更节）；
- `02-execution.md` 索引补 E-002/E-003。

## F-004（recommended）· CLI 模板 npmrc —— **fixed**（下一发布生效）

`apps/api/cmd/schema-ui/templates/web/npmrc.tmpl` 改为 npmjs 正向钉死映射（无 token、无 GH Packages）；**生效版本 = 下一次模块 tag**（模板随发布分发；v0.4.0 已发布的旧模板仅影响该版本骨架，注册到 R5 发布时核销）。

## F-005 / F-006（recommended）· 候选登记 —— **保持登记**（同意 A-001 R-001/R-002）

- `@schema-ui` org 正式化 = org 创建 + 用户指令触发；
- GH Packages 私有同名包退役评述 = R7 收口报告。

## 复核

F-001 关闭证据 = `golden-field/web/pnpm-lock.yaml`（npmjs tarball + 匹配 integrity）+ `web/.npmrc` 钉死行 + 空 store 重装日志；无凭据可复现（项目内无任何 token，用户级映射被项目级覆盖）。