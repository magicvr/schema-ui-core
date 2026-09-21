---
id: E-016-npm-consumer-remediation
parent_goal: GOAL-001-nav-group-collapsible
doc: execution-entry
status: recorded
date: 2026-09-21
created: 2026-09-21
updated: 2026-09-21
version: 0.1.0
---

# E-016 · v0.7.0 合并后 npm 消费回归与修复

## 已发生事实

1. PR [#16](https://github.com/magicvr/schema-ui-core/pull/16) 已于 2026-09-21 合并到 `main`，merge commit 为 `d083e3374ca5a96b6778a72500d318202384575c`；合并后 main CI run `35573602384` 对该 SHA 的 9 个 job 全部 `success`。
2. 按 A-016 / A-017 采纳的顺序，先构建并实际发布了初始四个 npm patch：`@magicvr/schema-ui-protocol@0.2.16`、`lib@0.1.15`、`renderer@0.3.14`、`ui@0.1.12`。它们已在 registry 可见。
3. 六包隔离消费者验证发现上述新版本存在 Node ESM 消费问题：protocol 子路径扩展名 / JSON import attributes、renderer 内部包别名、UI/lib 子路径扩展名；另有 protocol schemas 未进入 tarball。旧六包版本组合的导入检查通过。此发现发生在发布后，因此这些版本仍留在 npm 历史中；本记录不表示已撤包或 deprecate。
4. 在 `codex/npm-esm-consumer-fix` 分支准备了 patch 修复：新增 ESM import normalizer 与单测；package alias rewrite 在完成后自动运行 `finalize-lib-dist.mjs`，确保 schema/host-support JSON 资源随包输出；增加 tarball 隔离消费者 smoke 并接入 CI；将新候选钉为 protocol `0.2.17`、lib `0.1.16`、renderer `0.3.15`、ui `0.1.13`，同步 CLI、QUICKSTART 与 peer dependencies。
5. 本地证据：package metadata/ESM 测试 8/8 通过；`apps/api` 下 `go test ./cmd/schema-ui` 通过；六包构建、alias rewrite + 自动 finalization 通过；在 OS 临时目录里对六个真实 npm tarball 执行 `npm install`，并安装 React peers 与包依赖后，6 个根入口与 protocol app-manifest / host-support / runtime schema validator / node schema JSON、lib utils、UI lib utils / button 共 13 项导入通过；npm publish dry-run 对四个新 patch 指向发布、对现有 shell/theme 版本 skip。
6. 修复分支尚未创建 PR；四个修复版 npm patch 尚未发布，`apps/api/v0.7.0` tag 与 GitHub Release 尚未创建。下一门禁是独立审计、PR 检查全绿、合并后 main CI 全绿，再发布并验证四个修复包。

## 范围与事实边界

- 本条记录 PR #16 合并后的实际发布与故障发现，不回写或抹去 E-014/E-015 发生时的候选预审事实。
- 六平台 CLI zip、SHA256SUMS、tag 和 GitHub Release 均未生成。
- 当前工作只改发布包生成 / 验证路径与新 patch 版本，不改变 Root / workspace / VP 状态。
