---
id: GOAL-005-r3-admin-shell-navigation
doc: execution
status: done
parent: GOAL-001-mvp-admin-foundation
created: 2026-07-31
updated: 2026-07-31
version: 0.4.0
---

# 执行记录 · GOAL-005

## 时间线

### 2026-07-31 · R3 目标立项与范围登记

- `/govern` 复核当前显式工作区、Charter/VP 对齐、Root 路线图、R2 D-009/A-006 以及 R1 Web 边界。
- 创建本目标五件套和 `attachments/` 目录，并将 `GOAL-005-r3-admin-shell-navigation` 挂到 `GOAL-001-mvp-admin-foundation`。
- 将 R3 范围记录为 App manifest 装载、Admin shell、导航入口和路由语义；明确排除 R4 权限、R5 Renderer/业务范例及完整协议支持。
- 登记 `I-005-001` 至 `I-005-005` 为 required/open，记录方案冻结前的验证动作；当前没有把任何未知写成 verified。
- 同步工作区 `goal-tree.md`，Root `status` 仍为 `active`、`progress` 仍为 `2/6`。
- 本次没有修改 `apps/web`，没有产生 manifest loader、router、navigation 或 shell 的实现证据；父目标 `I-PROTO-002` / `I-PROTO-003` 未改变。

### 2026-07-31 · `/govern` 响应 A-001

- 采用 D-004 修正 A-001 F-001：`I-005-003` 和 `I-005-004` 的最晚需要阶段均改为「方案冻结前」，与 D-002 一致。
- 采纳 A-001 F-002：在 `I-005-001` 中显式关联 Root `I-PROTO-004`；未改变其在 Root 中的 `open` / `non-blocking` 状态，也未将其伪装成已验证。
- 采纳 A-001 F-003：Root 路线图的 R3 文案改为「规划中」，仅反映本目标已进入规划阶段；Root `progress` 保持 `2/6`。
- 本次没有修改 `apps/web`，没有收集或验证任何 `I-005-*`，也没有放行方案冻结、实现或 `done`。本响应不是同 scope 自审；是否需要自审仍待用户按 P-004.1 决定。

### 2026-07-31 · R3 规划阶段同 scope 自审计

- 按用户明确请求执行 GOAL-005 同 scope 自审，形成 [A-003](03-audit.md)；A-002 保持为编排响应，不替代 self 审计。
- 主线程复核工作区、A-001/A-002 闭合证据、固定协议入口、`apps/web` 当前源码与 Git 状态；确认没有代码、测试或协议接入变更。
- 在 `apps/web` 执行 `npm run build` 通过；`npm test` 因没有 `test` script 失败。该结果只记录当前 R1 骨架构建事实，不构成 R3 验收。
- 未修改 `apps/web`、Root `status/progress`、GOAL-005 `status` 或 `goal-tree.md`；`I-005-001` 至 `I-005-005` 仍为 `required/open`。

上述三段记录的是规划阶段当时的工作树快照。以下条目记录其后的真实实施事实，不改写历史审计结论。

### 2026-07-31 · R3 方案冻结与工作树实施

- 记录决策 [D-005](01-decision.md)：冻结 R3 的 2.7 manifest 子集、pinned artifact/hash、三 slot navigation projection、D4a route/active/fallback 语义、参数 pageRef href 规则、shell 固定区域和 R4/R5 边界。
- `I-005-001` 至 `I-005-005` 已依据 D-005 和下列可核对产物标为 `verified`；Root `I-PROTO-004` 仍保持 `open` / `non-blocking`。
- 在 `apps/web/src/protocol/app-manifest.ts` 实现 exact `2.7` manifest validation、`loadAppManifest()`、D4a route matching、home/deep-link/fallback、schema/logo URL resolution 和受限 navigation expression evaluation。
- 在 `apps/web/src/app/navigation.ts` 实现 `top` / `sidebar` / `user` projection、group pruning、context filtering、active-route 和参数 pageRef 绑定；在 `apps/web/src/app/App.tsx` 实现 header、desktop sidebar、mobile navigation、main surface、History API 和 unknown-route fallback；`main.tsx` 在 loader 失败时渲染 `ManifestFailure`。
- 添加真实静态 manifest：`apps/web/public/.well-known/schema-ui/app-manifest.json`；添加固定来源记录：`apps/web/src/protocol/upstream/provenance.json`。

### 2026-07-31 · R3 验证事实

- `apps/web/src/protocol/upstream-fixtures.test.ts` 校验 pinned schema、两个 behavior fixture 的 SHA-256 与来源 commit；机器断言执行 35/37 个 app-manifest cases、16/16 个 app-navigation cases，并登记两条 error-envelope 排除理由。negotiation/decoupled cases 的适配器只用于 fixture 对照，不扩展生产 host API。
- `apps/web/src/protocol/app-manifest.test.ts` 通过 `loadAppManifest()` 验证 `public/.well-known/schema-ui/app-manifest.json` 的真实字节；`apps/web/src/app/App.integration.test.tsx` 覆盖 root→home、站内 History 导航、popstate、未知路由 fallback、参数链接/context 和 ManifestFailure surface。
- 在 `apps/web` 执行 `npm test`：4 个测试文件、73 个测试全部通过（13 manifest unit + 3 navigation unit + 53 pinned fixture/provenance + 4 shell integration）。
- 在 `apps/web` 执行 `npm run build`：`tsc -b && vite build` 成功。
- 截至本条记录，R3 实现是未提交工作树事实；本条不声称代码已进入 HEAD、已发布或已通过完整协议 conformance。最终命令输出与运行时检查在关门审计响应中记录。

### 2026-07-31 · dev server 运行时复核

- 复用当前 `apps/web` dev server `http://127.0.0.1:4173/` 进行 HTTP 检查：`/.well-known/schema-ui/app-manifest.json` 返回 `200`、`application/json`，`protocolVersion` 为 `2.7`，包含 4 个 pages；根入口返回 `200` 并提供应用 boot shell。
- `App.integration.test.tsx` 已对根路径到 home 的 replace、站内 History API 导航、popstate、未知路由 fallback、参数链接/context 和 ManifestFailure surface 作行为断言；本条 HTTP 结果与该集成证据共同构成 R3 运行时入口复核，不宣称完整生产发布。

## 目标关门事实

- 用户已按 P-004.1 明确选择执行与 A-004 同 scope 的实施阶段 self-audit；A-006 已追加到 `03-audit.md`，verdict 为 `pass`。
- A-004 F-003 已由 A-006 以 `fixed` 合法闭合；F-004～F-006 保持 recommended、非阻断跟进。
- GOAL-005 已标为 `done`，Root R3 检查点和 `goal-tree.md` 已同步；Root 仍为 `active`、progress `3/6`。

## 完成后边界

1. R4 前闭合 `I-PROTO-002`，R5 验收前闭合 `I-PROTO-003`。
2. 按 R4-R6 推进；开放 required 信息项到期前不得越过对应门禁。

## 进度评估

R3 已完成：D-005 冻结方案，五项 required 信息已 verified，工作树实现、73 项自动化测试、构建、固定 fixture 对照和 dev server 入口复核均已落盘；A-006 self-audit 通过且无开放 required finding。Root `progress` 已从 `2/6` 同步为 `3/6`，不代表 R4-R6 或完整协议支持完成。
