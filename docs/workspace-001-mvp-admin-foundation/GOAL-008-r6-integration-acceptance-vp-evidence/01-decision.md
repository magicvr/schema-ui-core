---
id: GOAL-008-r6-integration-acceptance-vp-evidence
doc: decision
status: active
parent: GOAL-001-mvp-admin-foundation
created: 2026-08-01
updated: 2026-08-01
version: 0.1.1
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
**状态**：proposed（规划草案，尚未冻结）

**决定**：

提出以下 R6 执行结构，详见 [attachments/R6-acceptance-plan.md](attachments/R6-acceptance-plan.md) v0.1.0：

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

- 本决策不放行阶段 2；`I-008-001` 保持 `collecting`，其余 required 保持 `open`。
- 计划阶段下一步是收集/验证五个信息项，并对冻结候选做同 scope 计划审视。

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
