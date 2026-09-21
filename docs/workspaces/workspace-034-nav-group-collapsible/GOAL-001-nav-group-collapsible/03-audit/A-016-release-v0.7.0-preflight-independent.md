---
id: A-016-release-v0.7.0-preflight-independent
goal_id: GOAL-001-nav-group-collapsible
doc: audit-entry
source: independent
auditor: Grok Build grok-4.6 (high)
type: execution-facts
scope: PR #16 apps/api/v0.7.0 release-candidate and external release-gate preparation (CLI/QUICKSTART pins, six npm packages, package-name mapping, dry-run, Hosted CI, E-014/A-015; PR head 33b07b68)
date: 2026-09-21
verdict: pass
created: 2026-09-21
updated: 2026-09-21
parent: null
version: 0.1.0
---

# A-016 · apps/api/v0.7.0 发布候选预检 independent 审计（2026-09-21）

- **source**：independent
- **auditor**：Grok Build grok-4.6 (high)
- **类型** / **scope**：execution-facts · PR #16 `apps/api/v0.7.0` 发布候选与外部发布门禁准备（版本钉、六包关系、`npmjs-package-name` 归一化、dry-run、Hosted CI、E-014/A-015 事实记录；对照 A-015 F-S-001 `fixed` 声明）
- **verdict**：pass

## 范围与区间

工作区：`workspace-034-nav-group-collapsible`（`root_goal: GOAL-001-nav-group-collapsible`，`canonical_scope: docs/workspaces/workspace-034-nav-group-collapsible/`，`shared_materials_catalog: none`，`primary_plan: VP-034-nav-group-collapsible`）。Charter `schema-ui-core-admin-foundation@0.4.0` `active`。本意见只审本区 Root 的发布候选预检，未读取其他工作区目标正文。

审计对象：

- PR [#16](https://github.com/magicvr/schema-ui-core/pull/16) 当前 head `33b07b68d419d5fc28b202aa1469428fb2a9283c`（`dev` → `main`）
- 版本钉提交 `971ebbce5256968e6e6c238c5c089473f5fddde9`（`chore(release): pin v0.7.0 package face`）
- E-014、A-015 对预检事实与 F-S-001 `fixed` 的声明
- CLI / QUICKSTART、六包版本与依赖、`scripts/npmjs-package-name.mjs` 及其测试、发布 dry-run、Hosted CI
- 真实可核对的 npm registry 与 GitHub tag/Release 状态

本意见未修改 `status` / `progress` / goal-tree / 决策或执行正文 / 代码。本地工作树相对该 SHA 另有未提交改动，不纳入本条证据。

## 成果（有证据）

| 检查项 | 结果 | 证据 |
|---|---|---|
| 工作区绑定 | 合格 | `workspace.md` Root/canonical/`plan_refs`+`primary_plan` 一致；资料目录 `none` |
| P-005 信息项 | 不影响本门禁 | I-034-001～005 均 `verified`，最晚阶段为 R1–R5，已关闭；本预检未引入新的到期 required 信息项 |
| PR 状态 | 未合并 | PR #16 `state: OPEN`，`isDraft: false`，`mergeable: MERGEABLE`，`mergeStateStatus: CLEAN`，`headRefOid: 33b07b68…` |
| 候选 SHA 关系 | 可核对 | `971ebbce` 为版本钉；`33b07b68` = 其上 `docs(audit): record v0.7.0 preflight evidence`。A-015/E-014 记录的「最终候选 HEAD」是钉版本提交；当前 PR head 已前移到文档提交 |
| Go/CLI 版本钉 | 一致 | `apps/api/cmd/schema-ui/main.go`：`apiVersion=v0.7.0`，protocol `0.2.16` / lib `0.1.15` / renderer `0.3.14` / ui `0.1.12` / shell `0.1.6` / theme `0.1.4` |
| QUICKSTART | 一致 | 方法 B 注记同一 `apps/api/v0.7.0` 包面；create 模板使用 `@magicvr/schema-ui-*` |
| 六包源码面 | 一致 | `scripts/rewrite-lib-aliases.mjs` 的 `versions` 与 renderer peers 钉到上述新版本；shell 保持已发布 `0.1.6` 且 peer `@magicvr/schema-ui-protocol: ^0.2.12`（`^0.2.12` 覆盖 `0.2.16`） |
| 包名映射修复 | 可重复核对 | 旧路径把 `@schema-ui/lib` 的 leaf `lib` 拼成 `@magicvr/lib`。现 `npmjsPackageName()` 白名单六包并输出 `${scope}/schema-ui-${name}`。独立跑 `node --test scripts/npmjs-package-name.test.mjs` **4/4 pass** |
| 全六包名映射 | 可重复核对 | 内部名与 tarball 名 12 组输入均映射到 `@magicvr/schema-ui-{lib,protocol,renderer,ui,shell,theme}` |
| 错误名未入库 | 可核对 | npmjs 上 `@magicvr/lib` 与 `@schema-ui/lib` 均为 404 |
| 本地 dry-run | dry-run 成功，**不是发布** | 本会话 `node scripts/publish-npmjs-packages.mjs --dry-run` exit 0：skip shell@0.1.6、theme@0.1.4；dry-run 输出 `+ @magicvr/schema-ui-lib@0.1.15` / protocol@0.2.16 / renderer@0.3.14 / ui@0.1.12。产物目录 `apps/web/dist-lib/` 被 gitignore，未当作已发布资产 |
| Hosted CI（当前 PR head） | 9/9 success | run [35570762993](https://github.com/magicvr/schema-ui-core/actions/runs/35570762993) `headSha=33b07b68…`：web / api / api-postgres / 四组 browser E2E / 两组 container smoke 均为 `success`。web job 含步骤 `Test npmjs package name mapping` = success |
| Hosted CI（钉版本提交） | 9/9 success（父 SHA） | A-015 引用的 run [35569624987](https://github.com/magicvr/schema-ui-core/actions/runs/35569624987) `headSha=971ebbce…`，结论 `success`。合并门禁以当前 PR head 的 35570762993 为准 |
| A-015 F-S-001 `fixed` | 闭合证据充分 | 映射模块、测试、publish 脚本改写、CI 步骤与独立 dry-run 公开名均指向 `@magicvr/schema-ui-*`；本条确认该声明可重复核对 |
| npm 外部门禁 | **尚未发生** | registry 当前 latest：protocol `0.2.15`、lib `0.1.14`、renderer `0.3.13`、ui `0.1.11`；shell `0.1.6` 与 theme `0.1.4` 已存在。`0.2.16` / `0.1.15` / `0.3.14` / `0.1.12` **未发布** |
| Git tag / Release | **尚未发生** | 远端无 `apps/api/v0.7.0` tag；本地亦无。GitHub 最新 Release 仍为 [apps/api/v0.6.0](https://github.com/magicvr/schema-ui-core/releases/tag/apps/api/v0.6.0)（2026-09-08）。`apps/api/v0.6.1` tag 存在，不构成本次 v0.7.0 发布 |

## 对照成功标准

### A. PR 合并前门禁

| 标准 | 状态 | 证据 |
|------|------|------|
| 版本钉 CLI / QUICKSTART / 六包一致 | 满足 | `main.go`、QUICKSTART、`rewrite-lib-aliases.mjs` 同源数字 |
| A-015 包名 finding 已可核对闭合 | 满足 | 见上表 F-S-001 复核 |
| PR head Hosted CI 全绿 | 满足 | 35570762993 · 9/9 `success` · SHA = `33b07b68` |
| PR 可合并且未提前宣称已合并 | 满足（未合并） | OPEN / MERGEABLE / CLEAN |

**本门禁判定：可以进入合并。** `pass` 只覆盖合并前预检，不覆盖后续外部动作。

### B. 合并后 tag / npm / GitHub Release 门禁

| 标准 | 状态 | 证据 |
|------|------|------|
| PR 已合并到 `main` | 未发生 | PR #16 仍 OPEN |
| 合并后 `main` CI 全绿 | 未发生 | 无 merge commit、无 main push run |
| 四包新版本已在 npmjs 可见 | 未发生 | latest 仍为 0.2.15 / 0.1.14 / 0.3.13 / 0.1.11 |
| annotated tag `apps/api/v0.7.0` | 未发生 | `git ls-remote` 无该 tag |
| 六平台 CLI zip + SHA256SUMS + GitHub Release | 未发生 | 最新正式 Release 仍为 v0.6.0 |

**本门禁判定：尚未满足。** 本地 dry-run 的 `published … (dry-run)` 日志不是 registry 写入。CLI 已把 create 骨架钉到尚未发布的四包版本，因此 tag 必须发生在这四包 registry 可见之后。

## Findings

### F-001 · 合并后须先让四包在 npmjs 可见，再打 `apps/api/v0.7.0` tag

- 严重度：med
- 建议：recommended
- 状态：open
- 影响门禁：合并后 npm 发布校验、annotated tag、GitHub Release（**不阻断 PR 合并**）
- 描述：`schema-ui create` 与 QUICKSTART 已声明 protocol `0.2.16`、lib `0.1.15`、renderer `0.3.14`、ui `0.1.12`。这四个版本在 npmjs 上仍不存在。若先推 tag / 先发 GitHub Release，从 `apps/api/v0.7.0` 安装的 CLI 会生成无法 `pnpm install` 的骨架。A-015 已写出顺序「merge → main CI → 构建/发布 npm 并验证 → annotated tag → CLI zip → Release」；本条把它标为合并后必须遵守的开放建议，而不是已完成事实。
- 证据：`apps/api/cmd/schema-ui/main.go` 常量；QUICKSTART 包面注记；`npm view` latest 与版本数组；dry-run 对四包走 publish 分支、对 shell/theme 走 skip。

无其他 required / recommended finding。未发生的 merge / npm / tag / Release 记为外部门禁事实，不记为失败 finding，也不当作已完成。

## 必改项汇总

- 无。PR 合并前门禁无 open required。

## 开放非阻断建议

- **F-001（recommended / med）**：合并后先发布并核对 `@magicvr/schema-ui-{protocol,lib,renderer,ui}` 的新 patch 在 `https://registry.npmjs.org` 可见，再创建 annotated `apps/api/v0.7.0` 与 GitHub Release。shell@0.1.6 与 theme@0.1.4 保持幂等跳过。

## 与既有意见的异同

- 与 A-015 self 同向：版本钉一致、包名映射已修、Hosted CI 绿、外部发布未发生、F-S-001 可作为 `fixed`。
- 与 A-015 的 SHA/CI 差：A-015 记录最终候选为 `971ebbce` + run `35569624987`。当前 PR head 是 `33b07b68` + run `35570762993`。文档提交未改版本钉或映射代码；合并应使用当前 PR head，并以 35570762993 为合并前 CI 证据。
- 与 A-014 independent（v0.6.0 预检）同模式：`pass` 不等于 tag/npm/Release 已完成；未执行的外部门禁不写成成果。
- A-015 把初次 dry-run 暴露的 `@magicvr/lib` 错映射标为 F-S-001 required/high 并声明 `fixed`。本条独立复核该闭合：旧 `split("/").pop()` 路径已替换，测试与 CI 步骤覆盖回归，错误名未出现在 npmjs。闭合证据充分。

## 结论 + 建议给编排器/用户的下一步

本 scope 判定 **pass**。PR 合并前门禁可放行；合并后 npm / tag / GitHub Release 门禁全部未执行。不要把 dry-run 或 A-015 `pass` 写成发布完成。

建议 `/govern`：汇总 A-015 + A-016；在 PR #16 head 仍为 `33b07b68` 且 35570762993 全绿的前提下合并；合并后等 main CI 全绿，发布四包并 `npm view` 复核，再 annotated tag 与 Release；响应 F-001。

## 声明

本意见不修改 status/progress；响应由 /govern 处理。保证等级为 P-003 **L0**（入口分离），不是第三方鉴证。
