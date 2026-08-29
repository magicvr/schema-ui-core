---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-006-r5-release-and-gono-go
version: 0.1.0
---

# A-002 · R5 关门独立审计（source: independent · 2026-08-29）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型**：close-out
- **scope**：GOAL-006 S1–S4（发布可复现 / golden tarball 回归 / go/no-go / 关门就绪）+ VP-022 退出判据 #4/#5/#6，并核对 #1–#3 作为触发框架「#1–#5 全部达成」的证据链
- **verdict**：**conditional**
- **audit_type**：close-out
- **工作区**：`workspace-022-distribution-package-pilot` · Root `GOAL-001-distribution-package-pilot` · `shared_materials_catalog: none`

## 范围与区间

独立核对（不改 status/progress/方案正文）：

1. VP-022 六条方向级退出判据的证据链（GOAL-002～006 的 03-audit / 02-execution、golden-consumer / golden-web、go/no-go 报告）。
2. GOAL-006 本目标 S1–S3 事实与 S4 关门就绪（含 I-003 / I-007 最晚阶段 = R5）。
3. Charter `@0.3.0` strategic（VR-050）与试点证据 / VP-022 正文的一致性。
4. 本工作区范围与共享资料引用（catalog = none → 无资料引用可当关闭证据）。

本条是 R5/发布门禁的独立意见；Root 关门就绪另见 `GOAL-001` `A-001`。

## 成果（有证据）

独立本轮复跑（2026-08-29）：

| 项 | 结果 | 证据 |
|----|------|------|
| golden-consumer 装配冒烟 | `go run` exit 0：`kernel=2.0.0 profile=admin dialect=sqlite fresh=true contrib{routes=10 pages=2 perms=3 nav=1 frag=1} users_module=admin.users` | `GOAL-003-r2-go-library-consumption/attachments/golden-consumer/`（仍经 `replace` → `apps/api`） |
| golden-web 三探针 | protocol PASS（`2.9`）· render PASS（1573 B）· token override PASS（brand=2 ⊆ index=5） | `GOAL-004-r3-web-package-consumption/attachments/golden-web/{probe,probe-render,token-check}.mjs` |
| V2 能力 `normalizePageID` | `normalizePageID(' Roles ') → 'roles'`；`normalizePageID('  Users ') → 'users'` | 本轮 `node` 直调 `@schema-ui/protocol`（探针脚本本身未固化该断言） |
| npm 一键 pack | `node scripts/pack-npm-packages.mjs` 复跑产出 `schema-ui-protocol-0.2.0.tgz` · `schema-ui-renderer-0.1.0.tgz` | `scripts/pack-npm-packages.mjs`；`apps/web/dist-lib/artifacts/` |
| Go tag | 本地存在 `v0.0.2` = `c353fb58`（R4 关门提交）；`origin` 无此 tag | `git tag` / `git ls-remote --tags origin v0.0.2` 空 |
| go/no-go + 用户 GO | 报告存在；D-001 记录 GO + I-003 定案 + I-007 pin 2.9.0 | `attachments/gono-go-report-v1.md`；`01-decision/D-001-r5-release-and-gono-go.md` |
| Charter 0.3.0 | 成功边界 #1 已追加构建期包消费；非目标已澄清；pin `v2.9.0`（`81aa1d8`）；22 VP `vision_ref` → `@0.3.0` | `docs/vision/charter.md`；`docs/vision/revisions.md` VR-050；`docs/vision/reviews/VRev-050-charter-0-3-0-strategic.md`（self） |

既有 self 关门（不替代本条）：

- GOAL-002 A-001 self `conditional` → 用户 D-002 关门；F-001/F-002 在 GOAL-003 A-002 回填 `fixed`。
- GOAL-003 A-002 self `pass`（判据 #1 满足声明；诚实写明 replace；F-005 PG residual）。
- GOAL-004 A-001 self `pass`（判据 #2；F-006 d.ts recommended）。
- GOAL-005 A-001 self `pass`（判据 #3；样本 = A 层/协议 additive + 迁移 0063；配置键/依赖「R5 补」）。
- GOAL-006 A-001 self `pass`（S1–S3；0 required；S4 = 本独立审）。

## 对照成功标准

| 标准 | 状态 | 证据 | 缺口 |
|------|------|------|------|
| S1 发布流水线 | **部分** | pack 脚本独立复跑双 tgz；Go 本地 tag `v0.0.2` 存在 | 无 CI 接入；tag 指向 R4 提交、**不含** pack 脚本与 GOAL-006；tag 未推 origin |
| S2 golden tarball 回归 | **达成（有界）** | `package.json` 指向 `file:…/*.tgz`；三探针 + V2 能力本轮全绿 | V2 断言未写入探针脚本；安装语义 = 本地 tarball，非 registry |
| S3 go/no-go + 用户裁决 | **报告与裁决存在** | 报告 + D-001 GO + VR-050 已执行 | 报告将 #1–#5 标 ✅，与判据字面证据不完全同构（见 F-001/F-002/F-003/F-005） |
| S4 关门 | **未就绪** | 本条 | 开放 required；I-003/I-007 登记未闭合；台账索引空；goal-tree 无本目标 |
| 判据 #4 冻结面 | **方向满足、定稿有瑕疵** | `GOAL-002/attachments/freeze-face-v0.1.0.md` 正文 v1.2.0；semver + changelog 模板 | 文件名仍 v0.1.0；§4 R3/R5 勾选未改；§6 控制字符 |
| 判据 #5 发布可复现 | **部分** | npm pack 可复现 + golden-web tarball 绿 | 「脚本/**CI** 一键 Go tag **+** npm 包组」未齐；Go 消费仍 replace |
| 判据 #6 报告 | **产物在、触发框架过宽** | 报告六表 + Charter 草案 + 用户 GO | fork 对照耗时为推理；「全部达成」声明过宽 |

## 信息门禁（P-005）

| ID | 级别 | 最晚阶段 | 登记状态 | 本轮判定 |
|----|------|----------|----------|----------|
| I-003 | required | R5 | GOAL-006 meta = **collecting**；D-001 已定案（本地 tarball + 本地 tag + proxy 文档；registry 上传 = go 后） | **未在信息登记闭合**。决策文件存在，但 meta 仍 collecting，Root I-003 仍 open。到期 required 未 verified / 无 `accepted-residual` 书面范围 |
| I-007 | required | R5 | GOAL-006 meta = 「S3 交用户裁决」；D-001 声称 pin 2.9.0 已执行；Charter 0.3.0 已改 | **事实层已发生，登记未闭合**。GOAL-004/005 仍 collecting/registered |
| F-005（GOAL-003） | recommended residual | R4/R5 复审 | PG external 未测 | R5 E-001 **未**执行该复审 |
| 共享资料 | — | — | catalog = none | 无引用；未把跨区资料当证据 |

## Findings

### F-001 · 判据 #1 未以 `go get …@<tag>` 实证（消费面仍是 replace）

- **严重度**：high
- **建议**：required
- **状态**：open
- **描述**：VP-022 判据 #1 要求空下游仓**仅** `go get github.com/magicvr/schema-ui-core/apps/api@<tag>`。golden-consumer `go.mod` 明确 `replace … => ../../../../../../apps/api`，注释写「R5 落成前不依赖真实 tag」。R5 宣称 Go tag `v0.0.2`，但未去掉 replace、未演示 tag 消费。本轮：去掉 replace + `GOPROXY=off` → `module lookup disabled`；`git ls-remote --tags origin v0.0.2` 为空（仅本地 tag）。装配冒烟经 replace **可以**复现（见上表），故不是「不能装配」，而是**发布面消费路径名不副实**。GOAL-003 D-001 已否决「一直用 replace」作为 released 消费。Charter 0.3.0 成功边界 #1 现以 `go get` 为路径名，与可复现证据不一致。
- **证据**：`GOAL-003/attachments/golden-consumer/go.mod`；GOAL-003 A-002 核对点 #1；本轮 GOPROXY=off 失败；本地 tag `c353fb58`（R4 提交，不含 R5 pack 脚本）。

### F-002 · 判据 #3 字面三类演进未齐；R5 未补测；R4 决策文件缺失

- **严重度**：high
- **建议**：required
- **状态**：open
- **描述**：VP-022 判据 #3：「至少含**配置键变更 + 新增迁移 + 依赖更新**」。实做 = kernel `JoinKeys` + protocol `normalizePageID` + 迁移 0063。changelog **明文**「无配置键或依赖变更……留给 R5 发布回归补测」。GOAL-005 A-001 亦写「配置键/依赖样本 = … R5 补」。R5 E-001 / go/no-go **没有**该补测。GOAL-005 `01-decision.md` 索引指向 `01-decision/D-001-r4-drill-samples.md`，但 **`01-decision/` 目录与该文件不存在**——无法核对用户是否书面接受以 A 层/协议 additive 替换 VP 字面三类。go/no-go 仍将判据 #3 标 ✅。冲突 0 / 无 merge / 探针绿是真的；VP 字面分母不齐也是真的。
- **证据**：`docs/vision/plans/VP-022-distribution-package-pilot.md` §方向级退出判据 #3；`GOAL-005/attachments/changelog-upgrade-drill-v2.md`；`GOAL-005/03-audit/A-001-r4-s1-s4-closeout.md` 核对点 #1；`GOAL-005/01-decision.md`（悬空链接）；GOAL-006 E-001 / gono-go-report-v1 §1。

### F-003 · 判据 #5 「脚本/CI 一键 Go tag + npm 包组」未齐

- **严重度**：med
- **建议**：required
- **状态**：open
- **描述**：npm 侧 pack 脚本本轮可复现，golden-web tarball 探针绿——**这一半成立**。缺口：(1) `.github/workflows` **无** `pack-npm-packages`；(2) 未复用 `scripts/pre-release-smoke.sh` 思路接到一键入口；(3) `v0.0.2` 是 R4 关门 tag（`c353fb58`，2026-08-29 11:22），R5 提交 `96d7dca4` 在 tag **之后**，tag 树不含 `scripts/pack-npm-packages.mjs` 与 GOAL-006；(4) I-003 用户定案把 registry/proxy 发布列为 go 后动作，可作残余范围，但 GOAL-006 meta 仍 collecting，没有 `accepted-residual` 留痕。Windows 下 pack 使用 `shell: true`，本轮有 DEP0190 弃用警告（不阻断产物）。
- **证据**：本轮 `node scripts/pack-npm-packages.mjs`；`git cat-file -e v0.0.2:scripts/pack-npm-packages.mjs` 失败；`.github/workflows/r6-basic-matrix.yml` 无 pack；GOAL-006 `00-meta.md` I-003。

### F-004 · 本目标台账与 P-005 登记未达到关门可核对状态

- **严重度**：med
- **建议**：required
- **状态**：open
- **描述**：(1) `01-decision.md` / `02-execution.md` / `03-audit.md` 在本条写入前为 **0 字节**，而 `01-decision/D-001`、`02-execution/E-001`、`03-audit/A-001` 已存在——A-001 未入索引，违反 P-003「索引 + A 条目共同构成唯一正式台账」。(2) I-003/I-007 最晚阶段已到 R5，meta 未标 verified，也无 residual 书面范围。(3) `goal-tree.md` **没有**本目标（树止于 GOAL-005）；AGENTS §7：新建目标不更新 goal-tree = 任务未完成。本条只修复 03-audit 索引，不改 goal-tree / meta。
- **证据**：本目标三个索引文件；`goal-tree.md`；`00-meta.md` 信息表。

### F-005 · go/no-go「#1–#5 全部达成」过宽，且已驱动 Charter 0.3.0

- **严重度**：med
- **建议**：required
- **状态**：open
- **描述**：VP-022 触发框架「倾向推进」的前提是判据 #1–#5 **全部达成** + 冲突 0 + 升级耗时 ≤ fork 对照 + golden 稳定。报告全标 ✅ 后用户 GO，VR-050 已把包消费写入 Charter 成功边界 #1。独立核：冲突 0 与 golden 稳定（tarball/replace 语义下）可核对；#1/#3/#5 字面未齐（F-001～F-003）；fork 对照耗时为「分钟级 vs 小时级」推理，不是并排计时实验。Charter 0.3.0 **文本**与报告 §3 草案一致（本条不否定用户 GO 的发生）；独立审认定：**不能**把触发框架「全部达成」当作已核实事实。是否接受残余并维持 0.3.0，须 `/govern` + P-004 书面 residual / overruled，而不是本条改 Charter。
- **证据**：`attachments/gono-go-report-v1.md` §1–2；`docs/vision/charter.md` 成功边界 #1；`docs/vision/reviews/VRev-050-charter-0-3-0-strategic.md`（self；独立复核指向本 A-002）。

### F-006 · 冻结面 v1.2.0 定稿内不一致（文件名 / §4 勾选 / 控制字符）

- **严重度**：low
- **建议**：recommended
- **状态**：open
- **描述**：R5 声称 freeze-face v1.2.0 定稿。正文确有 §2c Web 面与 §6 发布注记，但文件名仍 `freeze-face-v0.1.0.md`；§4 仍未勾选「R3 Web 冻结面」「R5 发布通道」；§6 出现 `v`/`apps` 等控制字符损坏。不影响 A/B/B+ 分层已成文这一方向结论。
- **证据**：`GOAL-002-r1-contract-freeze/attachments/freeze-face-v0.1.0.md` 标题、§4、§6。

### F-007 · 挂账复审与对照实验未做完

- **严重度**：low
- **建议**：recommended
- **状态**：open
- **描述**：F-005（PG external 消费）复审触发 = R4 或 R5，R5 未跑。F-006 d.ts TS5056 仍为 go 后。判据 #2 的 shell/ui 未独立成包（D-002 粗粒度单包；该文件标题仍「待用户裁决」）。判据 #6 fork 对照无计时表。
- **证据**：GOAL-003 A-002 F-005；GOAL-006 E-001「残余」；GOAL-004 D-002 标题；gono-go-report-v1 §2 耗时行。

## 必改项汇总

1. **F-001**：要么补上「无 replace 的 `go get @v0.0.2`（或等价 GOPROXY/私有 proxy）消费实证」，要么由用户书面 `accepted-residual`（范围：G1 本地 replace；复审触发：首次 origin tag / proxy 发布）。
2. **F-002**：要么补配置键 + 依赖更新样本并再跑零冲突演练，要么用户书面接受「本波样本 = A 层/协议 additive + 迁移 0063」为 VP #3 残余（并补齐缺失的 R4 D-001 文件或等价决策留痕）。
3. **F-003**：要么 CI/一键入口覆盖 Go tag + npm pack，要么用户书面接受「npm = 本地 tarball 脚本；Go = 本地 tag；registry/proxy = go 后」为 I-003 残余，并把 I-003 登记改为 verified/residual。
4. **F-004**：补齐 GOAL-006 三个索引、闭合 I-003/I-007 登记、把本目标写入 `goal-tree.md`（由 `/govern` 执行）。
5. **F-005**：对触发框架过宽声明做 P-004：维持 Charter 0.3.0 并接受证据有界，或补证据后再确认 GO 口径。

未闭合以上 required 前，**不得**将 GOAL-006 标 `done`，也不得把 VP-022 六条判据当作已无条件满足。

## 与既有意见的异同

| 条目 | 同 | 异 |
|------|----|----|
| GOAL-006 A-001 self `pass`（S1–S3） | npm pack、tarball 探针、报告与 GO 留痕 — 独立复跑同意这些**产物存在** | self 0 required，并写「判据 #5 方向满足」。独立审：方向有证据，但 #1/#3/#5 字面与台账门禁仍有 required |
| GOAL-003 A-002 self | 同意装配闭环与 replace 诚实披露 | R5 未兑现「真实 tag 消费」延期 |
| GOAL-005 A-001 self | 同意冲突 0、无 merge、探针绿 | 「R5 补」配置键/依赖未发生；报告仍标 #3 ✅ |
| VRev-050 self `pass` | 同意修订分类、22 VP `vision_ref`、pin 2.9 为代码追认 | 触发框架证据被 VRev 当作「判据 #1–5 全绿」；本条不背书该全绿声明 |

无与 self 相反的「探针失败」结论；冲突在**完成口径**（字面判据 vs 有界试点），不是「没做」。P-004：是否 residual 须用户书面，本条不代裁。

## 结论 + 建议给编排器/用户的下一步

**verdict = conditional。** R5 的 npm tarball 发布脚本与 golden-web 回归可独立复现；Go 装配冒烟可独立复现；go/no-go 报告与用户 GO、Charter 0.3.0 文本落地都是事实。不能无条件关门：`go get @tag` 未实证、判据 #3 字面分母不齐、CI/一键 tag+npm 未齐、P-005 登记与台账索引未闭合，且触发框架「全部达成」过宽。

建议 `/govern`：

1. 展示 F-001～F-005，请用户逐条选 `fixed` / `accepted-residual` / `user-overruled`。
2. 在 residual 或补证完成前保持 GOAL-006 `active`，不要标 `done`。
3. Root 关门另响应 `GOAL-001` A-001（同源缺口 + 工作区/愿景对齐卫生）。

## 声明

本意见不修改 status / progress / 方案正文 / goal-tree 状态列。响应由 `/govern` 处理。
---

## 响应（2026-08-29 · /govern · source: self）

用户 P-004 书面裁决（逐条）：

| finding | 路径 | 范围 / 复审触发 |
|---------|------|-----------------|
| F-001 | **accepted-residual** | G1 本地 replace 为试点消费实证；复审 = 首次 origin tag / Go proxy 发布前必补 go get @tag 实证（go 后首个动作） |
| F-002 | **accepted-residual** | 本波样本 = A 层/协议 additive + 迁移 0063 为判据 #3 有界口径；配置键/依赖样本 = go 后发布回归补测；R4 D-001 决策文件已补建（fixed 部分） |
| F-003 | **accepted-residual** | 本波口径 = npm tgz 一键脚本 + 本地 Go tag；CI 接入与 registry/proxy 发布 = go 后；I-003 登记 = verified + residual 范围 |
| F-004 | **fixed** | GOAL-006 三索引已建、goal-tree 已含本目标、I-003/I-007 登记闭合（见 meta） |
| F-005 | **accepted-residual（维持 Charter 0.3.0）** | 用户书面接受触发框架证据有界性：判据按有界试点口径满足；残余 = F-001~003 接受项 + go 后清单 |
| F-006 | **fixed** | freeze-face 更名 v1.2.0.md + §4 勾选 + §6 控制字符清理 |
| F-007 | **fixed（部分）** | F-005 PG external（无环境）与 F-006 d.ts 保持 go 后清单；fork 对照计时 = residual（go 后并排实验）；挂账表已梳理 |

全部 required 已按三路径合法闭合；GOAL-006 可关门。