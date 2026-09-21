---
id: A-015-release-v0.7.0-preflight-self
goal_id: GOAL-001-nav-group-collapsible
doc: audit-entry
source: self
auditor: current-session
date: 2026-09-21
scope: PR #16 apps/api/v0.7.0 release-candidate preflight, npm package versions, package-name mapping, hosted CI gate
verdict: pass
parent: null
version: 0.1.0
---

# A-015 · apps/api/v0.7.0 发布候选预检 self 审计

## 范围与基线

- 当前工作区：`workspace-034-nav-group-collapsible`；焦点目标：`GOAL-001-nav-group-collapsible`。Root 与 VP-034 已关闭；本意见沿用 A-013/A-014 既有模式，只追加后续发布预检，不改状态或 progress。
- 发布 PR：[#16](https://github.com/magicvr/schema-ui-core/pull/16)，head `971ebbce5256968e6e6c238c5c089473f5fddde9`，目标 `main`；尚未合并。
- 审计范围：Go/CLI 与六包版本钉、候选 npm 资产、包名转换、最终 PR Hosted CI；PostgreSQL provider 修复细节另由 workspace-040 `GOAL-005` 审计。

## 成果（有证据）

| 检查项 | 结果 | 证据 |
|---|---|---|
| 版本递进 | pass | `apps/api/cmd/schema-ui/main.go` 与 QUICKSTART 钉到 `apps/api/v0.7.0`；四个有源码变化的包升 patch，shell/theme 不变；npm registry 当前版本读取与选择一致 |
| 六包构建 | pass | `node scripts/build-lib-packages.mjs` → alias rewrite → pack 全部 exit 0，生成 6 个候选 tgz |
| 公开包名 | pass | `node --test scripts/npmjs-package-name.test.mjs` 4/4 通过；发布 dry-run 对新版本使用 `@magicvr/schema-ui-{lib,protocol,renderer,ui}`，对既有 shell/theme 正确跳过 |
| 初次 dry-run 暴露的问题 | 已修复 | 初次内部名 `@schema-ui/lib` 错映射为 `@magicvr/lib`；新增六包白名单映射，避免写入错误的 npm 包名；dry-run 未上传任何内容 |
| Hosted CI | pass | 最终 HEAD `971ebbce` 的 run `35569624987` 9/9 `success` |
| 外部发布门禁 | 未发生 | PR merge、npm 实际发布、Go tag 与 GitHub Release 尚未发生；本审计不把它们当作完成事实 |

## Findings

| finding | level | 状态 | 证据 / 响应 |
|---|---|---|---|
| F-S-001 · 内部 package short name 会生成错误 npm scope 路径 | required / high | **fixed** | `scripts/npmjs-package-name.mjs` 映射至固定六包命名；4 条回归测试通过；修复后的发布 dry-run 对四个新包名全部正确 |

## 结论

本预检 scope 判定 **pass**，F-S-001 已按 `fixed` 路径闭合，open required = 0。发布顺序仍是：最终 PR CI 全绿 → merge → 合并后 main CI 全绿 → 从合并 commit 构建/发布 npm 包并验证 → annotated Go tag → 六个 CLI zip + SHA256SUMS → 远端资产校验 → 发布 GitHub Release。上述外部步骤尚未发生。

## 声明

本意见只记录发布候选的已验证事实，不修改目标 `status` / `progress` 或 goal-tree。独立复核意见及后续响应由编排流程处理。
