---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-003-r2-public-release-channel
version: 0.1.0
---

# A-002 · GOAL-003 关门独立审计（source: independent · 2026-08-29）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型** / **scope**：close-out · GOAL-003（C1–C4 · E-002 发布/消费实证 · 凭据卫生 · 幂等 · `@schema-ui` 边界登记）
- **verdict**：**fail**
- **工作区**：`workspace-024-distribution-formalization`（Root `GOAL-001-distribution-formalization` · `canonical_scope` 本区 · `shared_materials_catalog: none`）

## 范围与区间

核对 `00-meta` 成功标准 C1–C4 是否被 E-002 / 公开 registry / 下游仓制品 **可重复**支持；并核凭据卫生、发布脚本幂等、D-001 §6 `@schema-ui` 候选边界。不改 `status` / `progress` / 方案正文 / goal-tree。

下游仓 `github.com/magicvr/golden-field` 仅作为 workspace.md 写明的实验消费实证对象读取，不引入其他工作区目标状态。

## 成果（有证据）

| 项 | 独立核对 | 证据路径 |
|----|----------|----------|
| C1 六包已在 npmjs 公开 | `npm view` 六包可见；版本 = lib/theme/ui/shell `0.1.0` · protocol/renderer `0.2.0`；无认证 GET tarball HTTP 200 | registry.npmjs.org `@magicvr/schema-ui-{lib,protocol,theme,ui,renderer,shell}`（lib 创建 `2026-08-29T08:53:53Z`） |
| C2 origin tag + 公共 proxy | 本地与 origin tag `apps/api/v0.4.0` = `00d97b5b64145dbf590465c05d314f18384dbe0f`；`proxy.golang.org` `@v/v0.4.0.info` HTTP 200；`GOPROXY=https://proxy.golang.org` `go list -m …@v0.4.0` 成功；tag 内含 `schema-ui serve` | `git ls-remote origin refs/tags/apps/api/v0.4.0`；golden-field `go.mod` require `v0.4.0`、无 `replace`；`go.sum` 含 v0.4.0 哈希 |
| 探针语义（当前 node_modules **以及** 审计员从 npmjs 新装） | 三探针全绿：protocol `2.9` / render html `1573B` / token `brand=2 ⊆ index=5` | golden-field `web/probe.mjs` · `probe-render.mjs` · `token-check.mjs`；隔离 `npm install`（空 userconfig + 独立 cache + `--registry https://registry.npmjs.org`）后复跑同样全绿 |
| 凭据卫生 | 仓库根 `.env` 被 `.gitignore` 排除且 **未入库**；HEAD 跟踪文件无 token 值；脚本读 `.env` `npm_token` → `os.tmpdir()` 临时 `.npmrc` → `finally rmSync`；stdio 不打印 token | `.gitignore` L17–20；`git ls-files --error-unmatch .env` 失败（未跟踪）；`scripts/publish-npmjs-packages.mjs` |
| 幂等 | 发布前 `npm view versions` 命中则 skip；403 `cannot publish over the previously published versions` 兜底 skip | 同脚本 L79–107 |
| `@schema-ui` 边界 | npmjs 上 `@schema-ui/schema-ui-lib` **E404**（未误发）；D-001 §6 / E-002 残余 1 / A-001 R-001 已登记为 org 就绪后正式化候选 | `npm view @schema-ui/schema-ui-lib` 404；D-001 §6 |
| I-024-001（最晚 R2） | 用户书面裁决存在（真实发布 + token 注入点 + §6 改用 `@magicvr` 先行） | D-001；不阻本 scope 的信息门禁本身 |
| I-024-002 / I-024-003 | 最晚 R3 / R4，不阻 R2 | `00-meta` 信息表 |
| 共享资料 | `none`，无引用被当成证据 | `workspace.md` |

## 对照成功标准

| 标准 | 状态 | 独立证据 | 缺口 |
|------|------|----------|------|
| C1 npmjs 真实发布六包 `@magicvr/schema-ui-*` | **达成** | 公开 `npm view` + 无认证 tarball 200 + `publishConfig.access=public` | 标题/概述仍写 `@schema-ui`（见 F-003） |
| C2 `apps/api/v0.4.0` origin tag + 公共 proxy `go get`（含 serve） | **达成** | origin tag 哈希一致；proxy.info 200；`go list -m`；tag 内 `cmdServe` / `schema-ui serve` | — |
| C3 golden-field 升级 + **npmjs 无凭据**安装 + 三探针 | **未达成（声明过宽）** | go 侧达标；`package.json` 已改版本号且无 `file:`；`.npmrc` 仅注释、无 GH 映射；探针在现有安装与 npmjs 新装上均绿 | **committed `web/pnpm-lock.yaml` 仍钉 GH Packages tarball**，integrity **不等于** npmjs dist.integrity；无认证 GET 该 tarball **HTTP 401**。E-002「lockfile 落盘」与 golden-field commit `c379a5c` **未改** `pnpm-lock.yaml`（只改了 `.npmrc`/`package.json`）矛盾。见 F-001 |
| C4 发布流程成文 + scope 迁移注记 | **基本达成** | 脚本入库（token 注入 / `--access public` / 幂等）；D-001 §6 为 C4 点名的 changelog 注记 | 脚本默认 `PUBLISH_SCOPE=@schema-ui` 与实发 `@magicvr` 不一致（F-002）；`02-execution.md` 索引未挂 E-002（F-003） |

## Findings

### F-001 · golden-field lockfile 仍指向 GitHub Packages，C3「npmjs 无凭据消费」不可从制品复现

- 严重度：high
- 建议：**required**
- 状态：open
- 描述：C3 与 E-002 主张 web 六包改为 npmjs 公开消费，且「`pnpm install`（`NPM_CONFIG_USERCONFIG=空`）+ lockfile 落盘」。独立核对：
  1. `golden-field/web/pnpm-lock.yaml` 六包 `resolution.tarball` 均为 `https://npm.pkg.github.com/download/@magicvr/schema-ui-…`。
  2. lockfile integrity 与 npmjs `dist.integrity` **六包全部不同**（例：lib lockfile `sha512-tndx…zMw==` vs npmjs `sha512-fNhA…6KQ==`）。
  3. 无认证 GET lockfile tarball URL → **HTTP 401**；无认证 GET npmjs tarball → **HTTP 200**（lib 约 52 054 B）。
  4. golden-field `c379a5c`（信息写「npmjs 公开发布 · 无凭据 pnpm install · lockfile 落盘」）实际 diff **不含** `web/pnpm-lock.yaml`。
  5. E-002 记录安装耗时 864ms，与「清依赖后走本地 store 缓存、未重解析 registry」相符，不能证明公开 npmjs 拉取。
  6. 审计员用 `package.json` 副本 + 空 userconfig + 隔离 npm cache + `registry.npmjs.org` 安装成功（9 packages / ~10s），lock `resolved` 全部为 `https://registry.npmjs.org/@magicvr/…tgz`，三探针全绿——说明 **C1 包可无凭据消费，但 golden-field 未把该解析写入 lockfile**。
- 证据：`C:\Users\magicvr\Documents\Code\golden-field\web\pnpm-lock.yaml` L38–57；npmjs `npm view @magicvr/schema-ui-* dist.integrity dist.tarball`；curl GH 401 / npmjs 200；`git show c379a5c --stat`（仅 `web/.npmrc`、`web/package.json`）。
- 关闭条件：在 golden-field 用空 userconfig、**不复用** GH store 缓存（或删 lock 后重装）生成新 `pnpm-lock.yaml`，使六包 `resolved`/`tarball` 指向 `registry.npmjs.org` 且 integrity 与当前 npmjs dist 一致并提交；再附一次可重复的无凭据安装日志。在此之前不得把 C3 视为已核销，不得将 GOAL-003 标 `done`。

### F-002 · 发布脚本默认 scope 仍为 `@schema-ui`，与实发 `@magicvr` 不一致

- 严重度：med
- 建议：recommended
- 状态：open
- 描述：`scripts/publish-npmjs-packages.mjs` `PUBLISH_SCOPE` 默认 `"@schema-ui"`，文件头仍写「改名为 `@schema-ui/schema-ui-*`」。C1/D-001 §6 实发为 `@magicvr`。未设环境变量再跑脚本会打向未创建的 `@schema-ui` org（当前 404），既不能幂等跳过已发布的 `@magicvr` 包，也与成文流程不符。
- 证据：`scripts/publish-npmjs-packages.mjs` L12–13、L25；对比 C1 / D-001 §6。
- 建议修复：默认改为 `@magicvr`，或在脚本/README 把 `PUBLISH_SCOPE=@magicvr` 写成强制前置；头注释与 D-001 §6 对齐。

### F-003 · 台账与索引未随 §6 裁决对齐

- 严重度：med
- 建议：recommended
- 状态：open
- 描述：
  - `00-meta` title / 概述仍写 npmjs `@schema-ui/schema-ui-*` 与 golden-field 消费 `@schema-ui` 六包，与已勾选的 C1（`@magicvr`）及 D-001 §6 冲突。
  - I-024-001「证据/结论」仍写「`@schema-ui` 公开 scope」，未记 §6 改用 `@magicvr` 先行（授权门禁本身成立，但结论句过期）。
  - D-001 决策条 1 仍写公开 scope = `@schema-ui`，变更只在 §6。
  - `02-execution.md` 索引只有 E-001，**未登记 E-002**（文件已存在）。
- 证据：`00-meta.md` L3/L16/L38；`01-decision.md` 索引；`02-execution.md` 索引；D-001 L13–16 vs L30–35。
- 建议：`/govern` 响应时改措辞与索引，不把 progress 当完成证明。

### F-004 · CLI `create` 模板仍写入 GH Packages registry 映射

- 严重度：med
- 建议：recommended
- 状态：open
- 描述：D-001 §6 消费迁移 = 移除 `web/.npmrc` 的 GH Packages 映射后 `pnpm add @magicvr/schema-ui-*` 即 npmjs 公开版。但 `apps/api/cmd/schema-ui/templates/web/npmrc.tmpl` 仍写 `@magicvr:registry=https://npm.pkg.github.com`。`schema-ui create` 新骨架默认仍走私有 registry，与「公开发布通道」消费面不一致。属 C4/通道完整性质而非 C1 发布失败；可在本目标补模板，或明确登记到 R5/R7。
- 证据：`apps/api/cmd/schema-ui/templates/web/npmrc.tmpl`。

### F-005 · `@schema-ui` org 正式化仍为候选（同意 A-001 R-001）

- 严重度：low
- 建议：recommended
- 状态：open（登记；触发 = org 创建 + 用户指令）
- 描述：独立确认 `@schema-ui/schema-ui-lib` 未出现在 npmjs。候选边界已在 D-001 §6 / E-002 残余 1 登记，不阻 C1 在 `@magicvr` 下成立。迁移方案（新包名 + 消费方清单）仍待 org 就绪。
- 证据：npm view 404；D-001 §6；A-001 R-001。

### F-006 · GH Packages 私有同名包退役留 R7（同意 A-001 R-002）

- 严重度：low
- 建议：recommended
- 状态：open（R2 边界外）
- 描述：D-001 明确 GH Packages `@magicvr/schema-ui-*` 保留不删。F-001 说明两 registry 同名同版本 **integrity 不同**，lockfile 不改写就会继续命中私有面。退役/文档双轨说明属 R7，但 F-001 闭合前新克隆 golden-field 无法无凭据安装。
- 证据：D-001 未选方案表；E-002 残余 2。

## 必改项汇总

1. **F-001（required / high）**：golden-field `web/pnpm-lock.yaml` 必须重解析并提交为 npmjs 公开 tarball（integrity 与 `registry.npmjs.org` 当前 dist 一致），并留下可重复的无凭据安装证据。未闭合前 **禁止** GOAL-003 `done`，禁止把 C3 当已核销。

无其他 required。F-002～F-006 为 recommended。

## 与既有意见的异同（A-001 self）

| 点 | A-001 self | A-002 independent |
|----|------------|-------------------|
| C1 / C2 / 凭据卫生 / 幂等 | ✅ | 同意（独立 registry / tag / 脚本 / gitignore 复核） |
| C3 | ✅（空 userconfig + lockfile 落盘 + 三探针） | **不同意**：探针与「package.json 可从 npmjs 安装」成立，但 **committed lockfile 仍是 GH Packages，无凭据不可复现** |
| C4 | ✅ | 基本同意；补 F-002/F-003/F-004 |
| `@schema-ui` 候选 / GH 私有包保留 | R-001 / R-002 recommended | 同意 → F-005 / F-006 |
| required | 0 | **F-001 ×1** |
| verdict | conditional（self 侧 pass，待独立审） | **fail** |

**冲突（P-004）**：self 认为 C3 已证且无必改；independent 认为 C3 关闭声明对 lockfile **名不副实**，F-001 为关门 required。编排器须展示冲突、给建议、等用户裁决；未闭合 F-001 不得放行关门。

独立审计建议：采纳 F-001（`fixed`：重装 lockfile）；不要 `accepted-residual` 把「锁文件仍 401」说成已无凭据公开消费。

## 信息门禁（P-005）

- I-024-001 required · 最晚 R2：用户裁决 + 真实发布已发生 → 门禁本身 **verified**。结论句仍写 `@schema-ui` 记入 F-003，不另开 required。
- I-024-002 / I-024-003：最晚 R3/R4，不进入本 close-out 阻断集。
- 无共享资料引用。

## 结论 + 建议给编排器/用户的下一步

C1（npmjs 六包公开）、C2（v0.4.0 tag + 公共 proxy + serve 面）、凭据不入库、脚本幂等、`@schema-ui` 未误发——独立可重复。C3 在 **golden-field 可复现制品**上未完成：lockfile 仍要 GH token（401），与「无凭据 npmjs 消费」和 E-002「lockfile 落盘」不符。按尺度属 **关键主张名不副实** → **fail**。

建议 `/govern`：

1. 先响应 F-001：在 golden-field 删除或重建 `web/pnpm-lock.yaml`，空 userconfig + 非 GH store 重装，确认 lock 指向 npmjs 后提交；把安装日志写入 E-00N。
2. 按需顺手 F-002（脚本默认 scope）、F-003（meta/索引措辞）、F-004（`npmrc.tmpl`）。
3. F-001 `fixed` 并经复核前 **不要** 把 GOAL-003 标 `done`、不要把 Root 计为 2/7。
4. F-005/F-006 保持登记（org 触发 / R7）。

## 声明

本意见不修改 status / progress / 检查点 / goal-tree；响应、finding 闭合与关门由 `/govern` 处理。

---

## 响应（2026-08-29 · /govern · source: self）

- **F-001（required）→ fixed**：根因 = pnpm 无视 `NPM_CONFIG_USERCONFIG`、沿用用户级 `@magicvr→GH Packages` 映射（VP-023 遗留）；修复 = golden-field 项目级 `.npmrc` 正向钉死 `@magicvr:registry=https://registry.npmjs.org` + 删 lock/node_modules 后全新空 store 重装 → **lockfile 六包 tarball 全为 `registry.npmjs.org`（GH 残留 0），integrity 与 npmjs dist 一致**；golden-field commit `fb957a9`；证据/根因详见 E-003。
- **F-002 → fixed**：发布脚本默认 `PUBLISH_SCOPE=@magicvr`，头注释对齐 D-001 §6。
- **F-003 → fixed**：00-meta title/概述/I-024-001 结论句、D-001 决策条 1 → `@magicvr` 先行（§6 保留历史）；02-execution 索引补 E-002/E-003。
- **F-004 → fixed（下一发布生效）**：`templates/web/npmrc.tmpl` 改 npmjs 正向钉死；v0.4.0 已发旧模板影响面注册到 R5 发布核销。
- **F-005 / F-006 → 保持登记**（同意 A-001 R-001/R-002：org 触发 / R7）。

全部 required 合法闭合（fixed ×1 + 无 residual）。C3 在可复现制品上成立（lockfile = npmjs 公开 tarball · 无凭据可装）。GOAL-003 可关门。
