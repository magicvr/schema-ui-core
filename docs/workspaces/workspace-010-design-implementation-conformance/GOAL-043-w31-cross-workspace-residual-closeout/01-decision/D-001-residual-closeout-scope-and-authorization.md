---
id: D-001-residual-closeout-scope-and-authorization
doc: decision
status: accepted
goal_id: GOAL-043-w31-cross-workspace-residual-closeout
created: 2026-09-18
updated: 2026-09-18
parent: GOAL-001-design-implementation-conformance
version: 1.0.0
---

# D-001 · 残余收口范围、跨工作区授权与 `tsc` 余项口径

## 用户指令（原文）

> 现在处理的原则是，能现在处理的直接处理掉（可以在工作区10添加一个子目标承载治理上下文）。现在暂时不需要处理的，我们需要确保他们在路线图中被正确的统一登记，以便日后知道是什么情况，而非散落在治理文档各处（可能以后被忘了）。

## 决定

### 1. 承载位置 = workspace-010 的波次子目标

VP-010 是**持续符合性程序**（`active`，Root 为长期程序容器），其实现层以有界波次子目标承接整改——与本次「跨工作区残余统一收口」的性质一致，且避免重开已关门的 workspace-037/VP-037。故开设 `GOAL-043-w31-cross-workspace-residual-closeout`。

### 2. 跨工作区授权与可写范围（用户指令显式授权）

| 可写 | 说明 |
|------|------|
| `apps/web/e2e/list-visual-surface.spec.ts` | 增暗色计算背景断言、把分页契约断言扩展到第二个页面（workspace-037 GOAL-009 交付物） |
| `apps/web/src/app/HostFailureScreen.tsx` | 仅把文案映射表与通用键导出，供对照测试读取（无行为变更） |
| `apps/web/src/host/resource-feedback-parity.test.ts`（新增） | Host 终态与普通 resource 反馈的直接对照测试 |
| `docs/vision/roadmap.md` | 新增「未决项统一登记」节（与 `/vision` 协同） |
| 本目标台账 + workspace-037 相关 finding 的**闭合回填注记** | 按 P-003：只加闭合注记与证据指针，不改其 `status`/`progress`/历史正文 |

**禁止**：重开 VP-037/workspace-037；解除任何 gated 能力；逐条考古 269 行叙述式 `tsc` 记录；改 `apps/api`；重写历史记录正文。

### 3. `tsc` 文档余项的口径（`GOAL-008 A-002 F-002`）

- **可执行面**（`package.json` / CI workflow / `scripts/`）已由 `GOAL-008`/`GOAL-010` 的守卫与 CI 门禁锁定为 0 处非检查型调用——这是"证据可信"的充要条件，不再依赖历史文档形态。
- **文档侧**：扫描得 354 行提到 `tsc` 而命令形态无法从文本唯一确定（其中 45 行指 `npm run build`（= `tsc -b && vite build`，形态确定）、40 行为字面 `tsc --noEmit`（本区已更正或已加勘误注记）、269 行为叙述式如「tsc clean / tsc 0 / tsc 未受影响」）。
- **处置**：收口为 **bounded residual**（不逐条考古），触发条件 = **某条历史记录被再次当作类型检查证据引用时**，按 `GOAL-008 D-001` 的正确口径复核并在引用处注明；数量、口径与触发条件登记进路线图，避免散落遗忘。

### 4. `V-F124` 的处置

其要求（首波页面分母、状态字段、Profile/权限覆盖、Saved View 持久化边界、dirty-state/反馈类型的机器可核对矩阵）**实质已由 R1 交付**（`r1-denominator-matrix.json` 24/58、`r1-form-matrix.json`、`r1-state-feedback-matrix.md` + `D-003`/`D-004`/`D-005`，`I-037-001`～`004` verified）。故不重复劳动，而是**提请 `/vision` 以 Vision Review 复核并闭合**（愿景层判断不由 `/govern` 代作）。

## 未选方案

- **重开 workspace-037 承载**：会把已关门的 Root/VP 重新打开，与本波"收口"意图相反；不采纳。
- **逐条考古全部 354 行 `tsc` 记录**：文本无法唯一确定形态，收益与成本不成比例，且可执行面已被守卫锁死；不采纳。
- **把 gated 能力也一并推进**：用户明确"暂时不需要处理的"只登记不推进；不采纳。
- **由 `/govern` 直接判 `V-F124` fixed**：属愿景层 finding，须由 `/vision` 的审视作出；不采纳。
- **只登记不修复可处理项**：与用户"能现在处理的直接处理掉"相悖；不采纳。

## 边界与不变量

- 不改变任何目标的 `status`/`progress`（除本目标自身）；workspace-037 的门状态与历史结论不变。
- 登记节只描述**现状与触发条件**，不得把 deferred/recommended 写成已验证或已承诺。
- 本目标关门不影响 Root `GOAL-001-design-implementation-conformance` 的 active 程序容器状态。
