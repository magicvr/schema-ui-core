---
title: 决策 · R5 · 工程化、fork 体验与集成关门
status: active
created: 2026-08-02
updated: 2026-08-03
parent: GOAL-001-production-admin-foundation
version: 0.1.7
---

# 决策 · GOAL-008

## D-001 · 承接 Root D-013 方案边界并立项 R5 子目标

- **日期**：2026-08-02
- **状态**：accepted
- **决定**：按 Root **D-013** 已冻结的 R5 方案边界建立本子目标，以 S1～S5 五个核心检查点承载工程化、fork 体验与集成关门；S6（最小操作日志）为可选加分项，是否实施由用户书面决定并在此记录。
- **依据**：Root `I-005` = verified（部署基线 A：文档双进程 + 可选 Docker Compose；建议计时口径：终点=登录+后台首页可交互、不含依赖下载、≥1 次独立复现；复现方法：文档步骤 + smoke 清单 + 独立复现记录，R5 落地 `scripts/smoke.sh`）；Root `I-006` = closed（方案甲：操作日志为 R5 可选加分 checkpoint）。
- **边界**：
  - 不创建完整生产运维 / CI-CD 部署流水线；**Docker Compose 是 R5 必须交付和验收的第二启动路径**（非可选加分项）；fork 使用者可选择本地双进程或 Compose 两条启动路径（F-001 澄清，见 D-002）。
  - 不扩大 `I-PROTO-001 v0.1.3`；不重开已关闭的 R1～R4 子目标。
  - S6 操作日志不阻断核心验收；是否实施另行决定并留痕。
- **影响**：登记三项实施前 required 信息项 `I-008-001`（环境/配置/容器契约）、`I-008-002`（计时复现协议与 smoke 判据）、`I-008-003`（operation_log 契约，仅当 S6 实施）；未关闭前不得实施受影响范围或验收对应检查点。Root R5 检查点**不**据此勾选。
- **后续**：先在 `GOAL-008` 内关闭 `I-008-001/002` 并记录实施决策，再进入 S1/S2 实施；`I-008-003` 随用户对 S6 的决定处理。

### 未选方案

- **在本目标内重新裁决已冻结的 R5 方案边界**：Root D-013 已由用户书面确认，子目标只承接，不重复裁决。
- **立项时即把操作日志设为 required**：与 VP-002「加分项非硬关门条件」及 D-013 方案甲冲突。
- **仅做文档不落地容器与 smoke**：会把 VP-002 #7（Docker 一键启动）与 #6（可复现 fork 验收）留在 R5 未交付。

## D-002 · 响应 A-001 F-001：明确 Compose 交付义务（非可选加分项）

- **日期**：2026-08-02
- **状态**：accepted
- **决定**：采纳 A-001 F-001 推荐口径——**Docker Compose 是 R5 必须交付和验收的第二启动路径**：对 fork 使用者而言可选择本地双进程或 Compose 两条启动路径；**完整生产拓扑 / CI-CD 部署流水线仍为非目标**。Compose 的交付义务对应成功标准 **S2**（核心检查点，计入 `0/5` 进度分母）；**不是**像 S6 那样的可选加分项。D-001 边界原「Docker 一键启动为可选加分路径」表述与本决定不一致，予以修订（见 D-001 边界修订注记）。
- **依据**：A-001（independent · conditional）F-001 required：D-001 把 Compose 写成「可选加分路径」与 S2 核心交付要求语义冲突，会污染 `I-008-001` 契约与 R5/Root 关门判据。VP-002 #7 把「Docker 一键启动」列入基础工程化成功标准。
- **影响**：`I-008-001` 的 Compose/容器契约按「必须交付的第二启动路径」冻结；S2 不可被当作可选而跳过；Root D-013「部署基线 A」与 VP-002 #7 对齐边界一致。`I-005` verified 状态不重开（本澄清只在子目标层消除交付义务歧义）。
- **未选方案**：把 Compose 正式降为加分项并同步 S2、进度分母、Root D-013 与 VP-002——用户未作此裁决，且会留下 VP #7 未交付；故维持核心交付路径。
- **后续**：S2 实施与验收时以本决定为准；`I-008-001` 冻结 Compose 契约时引用本决策。

## D-003 · 冻结 I-008-001 环境/配置与容器部署契约

- **日期**：2026-08-02
- **状态**：accepted
- **决定**：采纳 [I-008-001-engineering-contract.md](attachments/I-008-001-engineering-contract.md)（v1.0.0）作为 S1/S2 的**环境/配置与容器部署契约**；`I-008-001` → **`verified`**，解除 S1（环境/配置基线）与 S2（容器一键启动）的方案冻结/立项目门禁。契约冻结内容：env 键全集与 dev/prod 行为（§1）、health/启动验证（§2）、Compose 服务/镜像/DB volume/探针（§3）、SPA fallback 与 `/api` 反代（§4）、依赖/超时/失败行为（§5）、CI 入口（§6）与验收清单 C-001～C-007（§7）。
- **依据**：Root D-013 部署基线 A（Compose 必须交付、fork 用户可选）；A-001 R-002 最低清单（production env/secrets、DB volume、SPA fallback 与 `/api` 反代、API/Web readiness、服务依赖/超时/失败行为、CI 入口）；GOAL-008 A-003（self · pass）确认立项与方案边界后可进入本方案冻结门禁。
- **边界**：
  - 契约冻结的是**形状与判据**，不替代 S1/S2 实现（Dockerfile/compose.yaml/README/nginx.conf 属实施，由实施留痕）。
  - 完整生产运维 / CI-CD 部署流水线、TLS、多实例仍为非目标。
  - `I-008-002`（计时复现协议 + smoke 判据）保持 `open required`，继续阻断 S3/S4。
- **影响**：`I-008-001` → verified；S1 环境/配置基线可进入实施准备；S2 容器一键启动可进入实施准备（均仍需在实施时留痕并对照验收清单）。Root R5 检查点**不**据此勾选；GOAL-008 `active / 0/5` 不变。
- **后续**：实施 S1（env 清单 + health/启动验证说明 + dev/prod 区分文档）；再实施 S2（Dockerfile × 2 + compose.yaml + nginx 反代 + CI smoke 入口）；S2 完成后建议一次实施向审计。

### 未选方案

- **在实施期临场决定容器/反代细节**：违反「先冻结再编码」；契约 C-001～C-007 为 S1/S2 验收判据。
- **把 15 分钟计时与 smoke 退出码纳入本契约**：属 `I-008-002` scope，避免契约混叠。

## D-004 · 冻结 I-008-002 计时复现与 smoke 验收协议

- **日期**：2026-08-03
- **状态**：accepted
- **决定**：采纳 [I-008-002-fork-reproduction-protocol.md](attachments/I-008-002-fork-reproduction-protocol.md)（v0.1.0，冻结参照 revision `5e27019482eb8d0695c402b784860233bbc90c39`）作为 S3/S4 的计时复现与 smoke 验收协议；`I-008-002` → **`verified`**。协议固定两条复现路径、依赖缓存排除与计时起止、四个 ≤15 分钟终点、独立复现记录字段、SM-001～SM-006、`scripts/smoke.sh` 最小输入、退出码和 disposable seed-reset 安全边界。
- **依据**：Root D-013 已由用户书面采纳终点=登录成功+后台首页可交互（列表加载）、不含依赖下载、≥1 次独立复现及「文档步骤 + smoke 清单 + 独立复现记录」方法；A-001 R-002 要求在关闭本项前固化工具/平台、缓存、计时、失败/重试、证据字段和机器可判定退出码。本轮实测了冻结时的工具版本与 repository ref；没有把这些静态事实写成 S3/S4 验收。
- **边界**：
  - 本决定只解除 S3/S4 的方案冻结/实施前**信息**门禁；不实施 QUICKSTART/README、`scripts/smoke.sh`，不产生独立复现、CI run 或实际 smoke 通过结论。
  - 当前 Windows 记录主机没有可用的 WSL `bash`，协议因此要求 S4 在 Linux CI 或已验证可用的 Git Bash/WSL 中执行；该环境探测不是脚本失败或通过证据。
  - `I-008-003` 保持 open（仅在用户决定实施 S6 时阻断）；S1/S2、GOAL-008 `active / 2/5`、Root R5 `4/5` 不变。
- **影响**：`I-008-002` → verified；S3/S4 可开始实施并应严格对照 v0.1.0。首次 S3 独立复现仍须满足 `<=900` 秒和四个终点，首次 S4 实现必须输出稳定检查项及退出码；完成这些实施事实后才可勾选相应成功标准。
- **后续**：实施 S3（fork 文档 + 独立复现记录）和 S4（`scripts/smoke.sh` + 本地/CI 运行证据），随后以阶段审计核对协议、实现与证据的一致性。

### 未选方案

- **把 Root D-013 的高层建议直接当成 S3/S4 已完成**：其只固定方向，未冻结本目标所需的工具基线、失败规则、记录字段与退出码，也不构成实施或验收事实。
- **在没有可用 bash 的当前 Windows 主机上把不存在的脚本记为已验证**：当前只记录环境事实；脚本实施后必须在实际可用的 POSIX shell/CI 上执行。
- **把独立复现、脚本实现与协议冻结合并为同一次结论**：会把未发生的 S3/S4 工作伪装成已完成，违反信息门禁与事实记录边界。

## D-005 · 响应 A-008：修订 I-008-002 协议至 v0.1.1（F-004 fixed）

- **日期**：2026-08-03
- **状态**：accepted
- **决定**：采纳 A-008（independent · design-plan · conditional）`conditional`。**F-004（required/medium）→ `fixed`**——[I-008-002 协议](attachments/I-008-002-fork-reproduction-protocol.md) 补丁至 **v0.1.1**：S4 验收**强制** ≥1 次 disposable/隔离运行且 `SM-006=PASS`；非 disposable 默认路径的 exit 0 不得单独作为 S4「种子可重复」关闭证据；保持「拒绝对普通开发库 destructive reset」安全边界（不安全 destructive → exit 2）。**同时吸收 R-008-001～004**（recommended，同版补丁一并处理）。
- **依据**：A-008 F-004 required——协议 v0.1.0 将 SM-006 标为「仅 disposable 模式」，未强制时可在从未验证种子幂等/不重复播种的情况下 SM-001～SM-005 全绿并 exit 0，与 S4 成功标准字面「种子可重复」及 Root D-013 smoke 清单漂移。R-008-001～004 为 recommended：默认 URL 未按路径区分（R-008-001）、「后台首页」操作化未对齐 manifest（R-008-002）、SM-005「SPA root」判据偏虚（R-008-003）、退出码 6 混用不安全 destructive 与种子断言失败（R-008-004）。
- **边界**：
  - 本决定只修订协议判据，**不**实施 S3/S4、不勾选检查点、不产生独立复现或脚本运行证据。
  - `I-008-002` 维持 `verified`（权威版本 v0.1.1）；不重开 S1/S2，不改 `I-008-001`。
  - `I-008-003` 保持 open（仅当 S6 实施时阻断）。
- **影响**：S4 实施与勾选必须以 v0.1.1 判据为准；CI `container-smoke` 或等价 job 需覆盖至少一次 disposable `SM-006=PASS`；QUICKSTART 必须覆盖 `list-edit-lifecycle` 路由（终点 4）。
- **后续**：实施 S3（fork 文档 + ≥1 次独立复现记录）与 S4（`scripts/smoke.sh` + 本地/CI 全绿 + disposable SM-006 证据）；S4 完成后建议实施向审计（self 或 `/audit`）。

### 未选方案

- **仅把 SM-006 当可选、仍允许非 disposable exit 0 作为 S4 种子关闭证据**：与 A-008 F-004 要求及 S4/D-013 字面「种子可重复」冲突。
- **暂缓协议修订而直接实施 S4**：会把未冻结的验收判据带入实施，违反「先冻结再编码」。
- **将 `I-008-002` 短暂回退 collecting**：协议主体与计时/复现字段在 v0.1.0 已可核对且未被证伪；补丁强化 S4 强制力即可维持 verified（A-008 亦未要求回退）。

## D-006 · 实施 S3 fork 文档与 S4 smoke 脚本

- **日期**：2026-08-03
- **状态**：accepted
- **决定**：按 [I-008-002 协议 v0.1.1](attachments/I-008-002-fork-reproduction-protocol.md) 实施 **S3**（fork 文档 + ≥1 次独立复现记录）与 **S4**（`scripts/smoke.sh` + 本地/CI smoke 全绿，含 ≥1 次 disposable `SM-006=PASS` 证据）。
  - **S3 交付**：根 `QUICKSTART.md` fork 上手段（前置/双启动路径/四终点/命令行与完整 smoke 用法/接业务指引）；独立复现记录 [R5-S3-REPRO-001](attachments/R5-S3-REPRO-001.md)（compose 路径，四终点全 PASS，`34.5s ≤ 900s`）。
  - **S4 交付**：`scripts/smoke.sh` 实现 SM-001～SM-006 + 退出码 `0/2/3/4/5/6/70` + `--disposable` 安全开关（无该开关不得执行种子 reset；不输出 token/password/secret）；本地 Git Bash + Docker disposable 运行 **SM-006=PASS**（log `r5-smoke-disposable-local.txt`）；CI `.github/workflows/r6-basic-matrix.yml` `container-smoke` 接入 `bash scripts/smoke.sh --disposable`（GitHub Actions 每次 runner 为隔离 project/volume，满足「CI 默认 disposable」）。
- **依据**：I-008-002 verified（v0.1.1）已解除 S3/S4 信息门禁；Root D-013 终点与 smoke 清单；A-008 F-004（S4 强制 disposable SM-006）；GOAL-008 成功标准 S3/S4。Git Bash（本机 `C:\Program Files\Git\bin\bash.exe`）可用，补足 D-004 记录的 WSL bash 不可用约束。
- **边界**：
  - 完整生产运维 / CI-CD 部署流水线仍非目标；S4 不替代 S3 的浏览器可交互验证（SM-005 仅 HTTP + SPA root，浏览器终点由 S3 独立复现承担）。
  - 复现记录为 `same-operator-clean-session`（协议 §3.3 允许同操作者，须声明独立性）；后续可邀请独立 `/audit` 或用户侧复现交叉核验。
  - `I-008-003`（S6 操作日志）保持 open；S5 仍待实施。
- **影响**：S3/S4 检查点勾选（`2/5 → 4/5`）；S5（阶段审计与 Root R5 勾选/关门条件评估）为最后一个核心检查点；Root R5 检查点仍待 GOAL-008 完成证据后勾选。
- **后续**：S3/S4 实施完成后做一次实施向审计（self 或 `/audit`），再实施 S5。

### 未选方案

- **只做文档、不落地脚本与独立复现**：会留下 S4「种子可重复」无机器判定证据，违反 I-008-002 §5 与 S4 成功标准字面。
- **在不可用的 WSL bash 上伪造脚本运行结果**：D-004 已记录 WSL bash 失败；改用已验证可用的 Git Bash，不把未运行当已通过。
- **把 S3 复现记录写成非独立（复用既有服务/DB）**：协议 §3.3 要求隔离 shell/checkout/DB；本轮先 `down -v` 清卷、全新 `up -d`，并声明 `same-operator-clean-session`。

## D-007 · 响应 A-011：F-005～F-009 fixed 路径整改（协议 v0.1.2 + 重写 smoke 守卫）

- **日期**：2026-08-03
- **状态**：accepted
- **决定**：采纳 A-011（independent · execution-facts · fail）`fail` 与用户 P-004 §3.2 裁决「**fixed 路径整改 F-005～F-009**」（未选 residual/overruled）。整改内容：
  - **F-005（required/high）→ fixed**：重新生成 clean-ref 独立复现记录 [R5-S3-REPRO-002](attachments/R5-S3-REPRO-002.md)——计时起点为 `docker compose up -d --build`，项目 build/配置/迁移/种子/登录/页面加载均计入计时（仅依赖下载与镜像层获取按协议 §3.1 排除），`≤900s`。
  - **F-006（required/medium）→ fixed**：R5-S3-REPRO-002 按协议 §4 全字段落盘——clean ref（干净 worktree，无未提交 diff）、checks 逐项含 smoke 输出、result 含命令/日志输出路径与失败/重试记录。
  - **F-007（required/medium）→ fixed**：[I-008-002 协议](attachments/I-008-002-fork-reproduction-protocol.md) 补丁 **v0.1.2**——新增部分绿退出码 `8`；非 disposable 路径**不得**以 `0` 退出（SM-006 未运行时输出非完整绿摘要 + exit 8），CI 不得把部分绿当 S4 完整绿。
  - **F-008（required/high）→ fixed**：协议 v0.1.2 §5.1 冻结 **disposable 隔离守卫**——强制 `SMOKE_ISOLATION_ID` + `SMOKE_DISPOSABLE_CONFIRM=yes`；脚本机器校验运行 project 与 `<id>_db-data` 卷绑定（`docker compose -p <id> ps -q api` + `docker inspect` 挂载检查），不满足 → exit 2；**删除 `SMOKE_RESTART_CMD` 任意 `eval`**，重启由脚本以校验过的隔离 project 直接执行；CI 改用显式隔离 project `-p ci-container-smoke` 并把隔离身份写入 job/step 环境与日志。`scripts/smoke.sh` 同步重写（含 R-011：重启后 readiness 重判 → exit 3）。
  - **F-009（required/medium）→ fixed**：对修复后的当前 revision 触发并保留 `container-smoke` 成功 run 证据（run URL、head SHA、SM-001～006 输出、teardown 结果）。
- **依据**：A-011 F-005～F-009（全部 required）证据与要求；I-008-002 协议 v0.1.1 §3/§4/§5；A-010（self · pass）主体事实（浏览器/本地 disposable）经 A-011 点验成立但证据边界不足。P-004 §3.2：同 scope verdict 冲突由用户书面裁决走 fixed。
- **边界**：
  - 本决定只整改 S3/S4 实施与验收证据，不重开 S1/S2、不关闭 `I-008-001`、不勾选 S5、不改 Root R5。
  - `I-008-002` 维持 `verified`，权威版本升至 **v0.1.2**（D-007）。
  - S3/S4 检查点保持勾选（`4/5`），其验收有效性以 F-005～F-009 关闭证据为条件。
- **影响**：`scripts/smoke.sh` 退出码语义变更（非 disposable 不再 exit 0）；CI `container-smoke` 使用隔离 project；QUICKSTART smoke 用法更新；R5-S3-REPRO-001 被 R5-S3-REPRO-002 取代为 S3 达标证据（历史保留）。
- **后续**：完成复现记录与 CI run 证据后落盘 A-011 响应节；建议对 F-005～F-009 关闭证据做一次 `/audit` finding-closure 复审，再进入 S5。

### 未选方案

- **F-005～F-009 走 accepted-residual / user-overruled**：用户书面裁决选择 fixed；计时/安全/CI 证据缺口属可修复项，无残余价值。
- **仅修订协议不重写脚本**：退出码与隔离守卫不落实现，机器可判定性无法成立。
- **保留任意 `SMOKE_RESTART_CMD` 仅加文档警告**：无法机器校验隔离，违反协议 fail-closed 安全边界。

## D-008 · 响应 A-012：F-005 fixed——无项目编译缓存重做 S3 计时复现（R-012 handled）

- **日期**：2026-08-03
- **状态**：accepted
- **决定**：采纳 A-012（independent · finding-closure · fail）对 **F-005（required/high）** 的关闭证据要求，按 `fixed` 路径重做一次**禁用/隔离项目编译缓存**的 S3 计时复现：
  - 预 T0（排除项内）执行 `docker rmi schema-ui-core-api:local schema-ui-core-web:local` + `docker builder prune -a -f`，**整体禁用 BuildKit 结果缓存**；基础镜像保持本地（镜像拉取排除项）。
  - 新记录 [R5-S3-REPRO-003](attachments/R5-S3-REPRO-003.md)：clean worktree（`1961e5a`，detached clean）+ **预 T0 `git status --short` 空** + 四终点全 PASS + **`64.833s ≤ 900s`（单调原始读数 403981233142700→404046066135300）**；BuildKit 归档输出中 **`go build`（12.8s）与 `npm run build`（6.1s）均实际执行、全程仅 1 条平凡 `CACHED`**——不再存在 A-012 所指的「编译层 CACHED」。
  - **R-012（recommended/low）→ handled**：预 T0 `git status --short` 与运行后状态分别留存；截图写于 worktree 外（`apps/web/test-results/`，gitignored）再归档并记录 sha256；单调计时工具原始读数（node `process.hrtime.bigint()` ns）逐终点落盘。
- **依据**：A-012 F-005（required/high：`13.5s` 不能证明 build-included，须隔离/禁用编译缓存重做）与 R-012；I-008-002 协议 v0.1.2 §3.1（项目自身编译不得预先完成）/§3.2（`.env` 写入在计时内）/§3.3（失败留痕）。用户书面指示「按 fixed 重做一次禁用/隔离项目编译缓存的 S3 计时复现，再请求仅针对 F-005 的关闭复审」。
- **边界**：
  - 仅重做 S3 计时复现证据与关闭留痕；不重开 F-006～F-009（A-012 已确认其关闭证据可维持）；不改协议 v0.1.2 正文、不改 `scripts/`、不改产品代码。
  - `I-008-002` 维持 `verified`（v0.1.2）；S3/S4 检查点维持勾选（`4/5`）；**未**推进 S5、未勾选 Root R5、未关门。
- **影响**：`R5-S3-REPRO-002` 被取代为历史记录（A-012 判定其编译层 `CACHED` 证据不足）；`R5-S3-REPRO-003` 成为 F-005 的 S3 计时证据。run log 内含按 §3.3 留痕的两次失败尝试（attempt #1 驱动脚本 stderr 中断；attempt #2 PowerShell→curl `-d` 引号被吞、登录体未送达——均为 runner 工具链故障，非被测栈故障，修正后 attempt #3 通过）。
- **后续**：落盘 A-012 响应节（03-audit）并同步 goal-tree；请求 `/audit` 仅针对 **F-005** 的 finding-closure 关闭复审；复审 pass 前不得推进 S5、勾选 Root R5 或关门。

### 未选方案

- **F-005 走 accepted-residual / user-overruled**：A-012 指明证据缺口属可修复项，且用户书面指示走 fixed；无残余价值。
- **仅在 REPRO-002 上补注、不重跑**：无法提供「编译层非 CACHED」的可核对新证据，不满足 F-005 字面要求。
- **用 `docker compose up --build` 加临时 BuildKit 环境变量绕过**：`docker compose up` 不支持 `--no-cache` 直传，且依赖未经协议冻结的私有机制；采用 `docker rmi` + `docker builder prune -a` 的协议可陈述、可复核做法。
- **修改 Dockerfile 强制 `--no-cache`**：改变产品实现为证据服务，且影响 CI 构建行为，超出 fixed 最小范围。
