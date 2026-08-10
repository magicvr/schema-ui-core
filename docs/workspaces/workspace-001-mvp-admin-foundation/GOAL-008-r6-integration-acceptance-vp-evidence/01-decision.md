---
id: GOAL-008-r6-integration-acceptance-vp-evidence
doc: decision
status: done
parent: GOAL-001-mvp-admin-foundation
created: 2026-08-01
updated: 2026-08-01
version: 0.3.0
---

# 决策记录 · GOAL-008

## 信息需求与阶段门禁

权威表见 [00-meta.md](00-meta.md)「信息就绪与未知项」。`I-008-001`～`I-008-005` 均为 required：允许规划如何收集信息，但在阶段 1 结束前不得把验收合同写成已冻结，也不得进入受影响的 R6 验收执行或 VP 关门提案。

## D-001 · 立项 R6 集成验收与 VP 证据子目标

**日期**：2026-08-01
**状态**：accepted

**决定**：

1. 在 `GOAL-001-mvp-admin-foundation` 下创建 `GOAL-008-r6-integration-acceptance-vp-evidence`，状态为 `active`、处于规划期。
2. R6 只承接集成验收与 VP 工作区证据：对照 VP-001 三条退出判据建立验收合同、执行集成验证、汇编可复核证据并完成 Goal close-out 审计。
3. Root 路线图 R6 标记为「规划中」，Root `progress` 保持 `5/6`；立项不等于 R6 完成、Root `done` 或 VP 可关门。
4. 当前不修改 `apps/*`、VP status 或既有 R1-R5 关门事实。

**为什么**：

- 用户明确调用 `/govern 规划 R6 — 集成验收与 VP 证据`。
- Root 已将 R6 定义为唯一未完成纲领检查点；VP-001 要求退出证据落在挂接工作区目标记录，而当前没有 R6 专属目标或集中证据索引。
- 范围横跨运行态、协议回归、账号权限集成、证据格式与 VP 映射，满足 P-001 的分阶段条件，应先立项并写高层路线图与 P-005 信息门禁。

**未选方案**：

- **只在 Root 继续追加零散证据**：会把 R6 多阶段执行与 Root 长期记录混在一起，难以形成可关门的有界证据路径。
- **直接修改 VP 为 closed**：无 R6 工作区证据与关门审计，违反 alignment 的区证据和用户确认门禁。
- **把 R5 的 395 项测试直接当作 R6 完成**：R5 scope 不含干净启动、浏览器/API 集成、机器可读证据包或 VP 三判据汇编。

**影响**：

- 新目标与 `goal-tree.md` 登记为 `active`；Root R6 进入规划中。
- 新增 `I-008-001`～`I-008-005` required 信息项；它们阻断阶段 1 冻结及相应执行/关门主张。
- Root / VP status、Root `progress` 与应用代码均不变。

## D-002 · 提出“四阶段 + 机器可读证据索引”的 R6 计划

**日期**：2026-08-01
**状态**：accepted（阶段 1 计划审视 A-002 pass 后，用户 `/govern` 授权冻结）

**决定**：

提出以下 R6 执行结构，详见 [attachments/R6-acceptance-plan.md](attachments/R6-acceptance-plan.md) v0.2.0：

1. 验收合同与证据计划冻结；
2. 集成验收执行；
3. VP 证据汇编与缺口整改；
4. R6 关门审计与 VP 提案输入。

证据包拟采用机器可读 index 统领运行结果、revision/environment identity、命令与退出码、排除/residual 和 SHA-256；具体 schema、最低环境矩阵与账号权限 oracle 仍由 `I-008-002`～`I-008-005` 决定。

**为什么**：

- VP 退出判据覆盖三个不同主张，单一“测试通过”结论无法表达各自主张、边界与残余风险。
- 现有 Web/API 测试、pinned fixture/SHA 与 R5 登记表是强输入，但默认输出主要为文本，且仓库未发现 CI workflow、浏览器 E2E 或统一 evidence writer。
- 先冻结证据合同再执行，可避免验收后补写口径，保证失败、排除与环境缺口在收集时即被记录。

**未选方案**：

- **先跑所有现有命令，再决定证据格式**：可能丢失 revision、环境、退出码或排除边界，无法稳定复核。
- **要求完整发布/部署才允许 R6**：VP-001 当前要求可运行、可 fork 与集成可验证；发布/部署是否构成最低门禁仍需 `I-008-005` 决定，不在规划草案中擅自升级。
- **仅生成一份 Markdown 总结**：缺少机器可解析结果与文件摘要，难以证明证据集合完整且未漂移。

**影响**：

- **本决策已冻结**（2026-08-01，A-002 pass + 用户授权）：阶段 1 冻结候选计划正式成为 R6 执行计划；阶段 2 可开始。
- `I-008-001`～`I-008-005` 均 `verified`（见 [03-audit.md](03-audit.md) A-002 / F-008-001 响应）；CI 首跑 green。
- 阶段 2 执行须按 evidence index 持久化结果；Root `progress: 5/6`、VP-001 `active` 不变。

## D-003 · 记录阶段 1 本地能力基线与证据边界

**日期**：2026-08-01
**状态**：accepted

**决定**：

1. 将当前 revision `7d20acc7702bcc0e514f787c455bf9c93d5b832f` 的 Windows/amd64 能力盘点记录为阶段 1 的规划输入：Node `v22.17.0`、npm `10.9.2`、Go `1.26.0`、Vitest `3.2.7`；Web `15` 个测试文件 / `395` 项测试通过且 build 通过，API `go test ./...` 与 `go build ./...` 通过。
2. 将实际入口固定为候选命令与运行态：`apps/web` 下 `npm test` / `npm run build`，`apps/api` 下 `go test ./...` / `go build ./...`；API `GET /healthz`、Web `:5173` 到 API `:8080` 的 `/api` proxy 作为后续验收入口候选。
3. 新增 `attachments/evidence-index.schema.json` 与 `attachments/evidence-index.dry-run.json` 作为 `I-008-004` 的 draft 候选；它们必须保持 `planning`/`draft` 语义，不能当作正式 R6 evidence artifact。
4. 记录未发现 `.github/workflows`、专用浏览器 E2E runner 或 JSON/JUnit/evidence writer；这使 `I-008-005` 继续处于收集状态，不静默把当前 Windows 结果扩展为 Linux/CI 或浏览器证据。

**为什么**：

- 这些事实均由当前工作树的命令、配置和只读扫描直接核对，能够缩小 `I-008-002`、`I-008-004`、`I-008-005` 的未知范围。
- R4 已有的单端、fixture 和 HTTP 事实可以作为 `I-008-003` 的输入，但不能替代尚未冻结的 API→Web/Renderer/动作链正向/拒绝 oracle。

**未选方案**：

- **把本地命令全绿写成阶段 1 已冻结或 R6 已验收**：会绕过 revision/environment identity、artifact hash、跨层 oracle 和平台矩阵门禁。
- **把没有 CI/browser runner 写成已接受 residual**：当前没有用户书面范围、期限、缓解与复审触发，不满足 P-004.4。

**影响**：

- `I-008-001`～`I-008-005` 均进入“已收集但未闭合”的可追踪状态（`collecting`）；阶段 1 仍未冻结，阶段 2 仍关闭。
- Root `progress: 5/6`、GOAL-008 `active`、VP-001 `active` 与应用代码不变。
- 下一步是验证 draft schema/dry-run、冻结跨层账号权限 oracle 和最低环境矩阵候选，再进行同 scope 计划审视；若需要 residual/有界实验，先按 P-004.4 请求用户裁决。

## D-004 · 最低环境矩阵采用“本轮搭建最小 CI + 浏览器矩阵”（用户裁决 · I-008-005）

**日期**：2026-08-01
**状态**：accepted（用户 `/govern` 裁决）

**决定**：

1. `I-008-005` 的答案是**不**接受“Windows-only + 平台 residual”，而是**本轮搭建最小 CI + 浏览器矩阵**作为 R6 最低验收矩阵。
2. 具体机制：
   - **Windows/amd64 本地**：开发机命令（web/api test+build、双服务启动、Playwright）——本机已验证（阶段 1 dry-run + E2E 通过）。
   - **Linux/amd64 CI 等价**：新增 `.github/workflows/r6-basic-matrix.yml`（web / api / browser-e2e 三 job，Node 22 + Go 1.26 + Playwright Chromium）。
   - **浏览器 E2E**：Playwright，`apps/web/e2e/shell.spec.ts`，webServer 启动双服务，验证 shell 渲染 + `/api` proxy 账号上下文。
3. CI 实际首跑发生在推送到远端后；在跑绿前，不得把“已配置”写成“已跑绿”，Linux 等价证据按阶段 2 执行事实处理。

**为什么**：

- 用户明确选择“本轮搭建最小 CI + 浏览器矩阵”而非接受有界 residual（P-004.4 裁决点，用户决定）。
- VP-001 退出判据要求“可运行、可 fork、可验证”；缺少 Linux/CI 与浏览器证据会削弱 fork 与集成验证主张。
- Playwright 为最小浏览器 runner；workflow 覆盖干净安装（`npm ci`）与 Linux 等价。

**未选方案**：

- **接受有界 residual（Windows-only）**：用户未选；缺失 Linux/浏览器证据会让 VP-001 判据 1/3 的“可 fork、可验证”主张证据不足。
- **延后裁决**：用户未选；阶段 1 冻结继续受阻。

**影响**：

- 新增 `.github/workflows/r6-basic-matrix.yml`、`apps/web/playwright.config.ts`、`apps/web/e2e/shell.spec.ts`；`apps/web/package.json` 增加 `@playwright/test` 与 `test:e2e`；`apps/web/vite.config.ts` 固定 `host:127.0.0.1`（保证浏览器矩阵在 Windows/IPv6 与 CI 上确定性）。
- `I-008-005` → 由“待决定”转为“决定 + 机制已建”，状态进入可审视的冻结候选；Linux/CI 跑绿仍为阶段 2 执行事实。
- Root `progress: 5/6`、GOAL-008 `active`、VP-001 `active` 不变。

## D-005 · R6 验收矩阵与证据形状升级为冻结候选（I-008-001/003/004）

**日期**：2026-08-01
**状态**：accepted（A-002 计划审视 pass 后随 D-002 冻结）

**决定**：

1. [R6-acceptance-plan.md](attachments/R6-acceptance-plan.md) v0.2.0 将 VP 三条退出判据映射为验收矩阵（C-001～C-008），并把最低环境矩阵（§4c）、evidence schema 校验结论（§3）登记为冻结候选。
2. [account-permission-oracle.md](attachments/account-permission-oracle.md) v0.1.0 登记账号权限跨层正向/拒绝 oracle（I-008-003）。
3. `evidence-index.schema.json` 经 [validate-evidence-dry-run.mjs](attachments/validate-evidence-dry-run.mjs) 校验：可解析、5 个 artifact SHA-256 可重算；仍为 draft，正式 acceptance index 属阶段 2。
4. 本决策不把阶段 1 标为冻结，不放行阶段 2；D-002 仍为 proposed 直至计划审视通过。

**为什么**：

- 阶段 1 退出条件要求五项 required 信息有证据结论或合规 residual；本轮已把矩阵、环境、oracle 与 schema 形状落到候选，并实际执行本地 dry-run/E2E 作为证据。
- 生成候选后，由同 scope 计划审视（A-002）核对后再由用户确认冻结，避免“自写自冻结”。

**未选方案**：

- **直接冻结**：跳过计划审视；违反阶段 1 门禁（需审视通过 + 0 开放 required）。
- **维持全部 collecting**：不反映本轮已发生的本地能力、矩阵与 oracle 证据。

**影响**：

- 计划附件与 oracle 附件成为冻结候选；`I-008-001`/`I-008-003`/`I-008-004` 有证据路径可审视。
- 阶段 1 仍未冻结；阶段 2 仍关闭；Root `progress: 5/6`、VP-001 `active` 不变。

## D-006 · 响应 A-006 并授权 GOAL-008 关门（2026-08-01）

**日期**：2026-08-01
**状态**：accepted
**关联意见**：[03-audit.md](03-audit.md) A-006（independent · close-out · pass）
**关联决策**：Root [D-016](../GOAL-001-mvp-admin-foundation/01-decision.md)

**决定**：

1. 接受 A-006 `pass` 结论（与 A-005 self close-out `pass` 同向、无冲突）；本目标无未合法闭合的 required finding，`I-008-001`～`I-008-005` 均 verified。
2. **GOAL-008 → `done`**（用户 `/govern` 授权）：R6 四阶段完成、验收矩阵 C-001～C-008 全执行、`evidence-index.json`（mode: acceptance）7 artifact SHA-256 verified、VP 三判据 Q2 汇编落盘。
3. recommended 处置：`F-008-003`（文档标签同步冻结/完成语义）与 `F-008-004`（evidence builder 锚定 `apps/web` 依赖、可自 `attachments/` 重跑）→ **fixed**；`F-008-002`（CI Node 20 弃用 / go.sum 缓存 skip）→ tracked，不阻断。
4. Root 纲领 **R6** → 完成；`progress` 5/6 → **6/6**；Root status → `done`（Root D-016）。
5. VP-001 保持 `active`：关门提案由 `/vision` 读取工作区 Q2 证据并经用户确认；R6 `pass` 不等于 VP 已关门或完整协议支持。

**为什么**：

- A-005 与 A-006 双关门审计均 pass、开放 required=0，满足「相关意见无未合法闭合 required、相关信息项无未处理关门 required、至少有阶段/关门向审计、成功标准可核对」的关门条件。
- 用户本轮 `/govern` 明确指令：响应 A-006、授权 GOAL-008 → done、决定 Root R6 / progress 6/6 / status。

**未选方案**：

- **保持 active 待 VP 关门**：目标自身关门条件已满足；VP 关门是 `/vision` 的独立决策，不阻断本目标 `done`。
- **不改文档标签 / 不修 builder**：A-006 的 recommended 可修可不修；两者成本低且直接改善台账可读性与证据可重跑性，本轮一并处理。

**影响**：

- GOAL-008 `status` → `done`；goal-tree 同步；Root R6 完成、`progress: 6/6`、Root status → `done`。
- VP-001 仍 `active`；`/vision` 可读取本目标 [vp-evidence-assembly.md](attachments/vp-evidence-assembly.md) 与 evidence 目录作为关门提案输入。
