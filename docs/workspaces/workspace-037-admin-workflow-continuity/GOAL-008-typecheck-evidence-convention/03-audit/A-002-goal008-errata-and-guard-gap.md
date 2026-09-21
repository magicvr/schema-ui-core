---
id: A-002-goal008-errata-and-guard-gap
doc: audit-entry
status: recorded
goal_id: GOAL-008-typecheck-evidence-convention
created: 2026-09-18
updated: 2026-09-18
parent: GOAL-001-admin-workflow-continuity
version: 1.0.0
---

# A-002 · 跨区勘误执行复核与守卫缺口（self）

- **source**：self
- **日期**：2026-09-18
- **scope**：`I-008-004` 授权后的跨工作区勘误执行（`E-006`）＋ 防复发守卫对「非 `-b` 有效形态」的覆盖度
- **verdict**：`conditional`

## 成果（有证据）

| 核对项 | 结论 | 证据 |
|--------|------|------|
| 用户授权是否落盘 | 是 | Root `D-015`；`00-meta.md` `I-008-004` → `verified` |
| 勘误是否只加注记、不改历史 | 是 | 11 处注记均含「原记录保留不改」，原命令/结论/verdict 未变；`git diff` 仅新增行（无删除型改写） |
| 跨区引用是否合规 | 是 | 注记统一使用 Q2 路径指向本目标目录，未嵌入工作区号进 goal id |
| 有效形态是否被误改 | 否 | workspace-002 两处 `-p e2e/tsconfig.json` 与 workspace-037 的 `-p tsconfig.app.json` 复核为有效，未改动 |
| 是否改动他人 `status`/`progress`/verdict | 否 | 仅新增注记；他人 `03-audit.md` 索引与 verdict 未变 |
| 空转判据是否有实测支撑 | 是 | `E-006` 注入错误实测表（4 条命令的退出码） |

## Findings

### F-001 · 防复发守卫接受 `-p <solution-style config>`，同类空转仍可复现

- 严重度：medium
- 建议：recommended
- 状态：**open**
- 关联：`D-001`（判据表述）、`E-003`（守卫实现）、`E-006`（新实测证据）
- 描述：`apps/web/src/typecheck-convention.guard.test.ts` 以 `CHECKING_FLAG = /(-b|--build|-p|--project)/` 判定「是否为检查型调用」，只要求出现 `-p` 令牌，**不校验被指向的配置自身是否选择源文件**。因此 `tsc --noEmit -p tsconfig.json`（根 solution-style，`files: []`）会被判为合规，而它实测 **exit 0、未检查任何文件**——正是本目标要防的失效模式，只是换了个 `-p` 外壳。守卫自身的 `selectsOwnSources()` 已具备判断能力，但只用于断言根配置与 e2e 配置的形态，未用于校验 `-p` 的目标。
- 证据：`E-006`「本次新增的第二层事实」实测表；守卫源码 `CHECKING_FLAG`（第 87 行）与 `selectsOwnSources`（第 63–76 行）的使用点。
- 影响面：当前仓库**没有**任何可执行面（`package.json` / CI workflow / `scripts/`）使用 `-p tsconfig.json`，唯一出现处在 workspace-011 的历史文档中（已加勘误）；因此本缺口目前不产生失实证据，属预防性缺口。
- 建议处置：在守卫中把「`-p` 目标配置必须 `selectsOwnSources`」纳入断言（并补 `tsc --noEmit -p tsconfig.json` 的变异用例）。修复需要改动已关门目标的交付物，应由用户裁决：新开整改子目标，或并入后续轮次。
- 闭合要求：改动守卫并经变异验证后 `fixed`；或用户书面接受为残余（须写明范围与复审触发）。

### F-002 · 全仓「tsc 简写」未被裁定

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：本次范围限定为字面 `tsc --noEmit` 与 `-p tsconfig.json` 条目；全仓另有 300+ 行以「tsc clean / tsc 0 / tsc 全绿」等简写描述验证，其中的命令形态无法从文本唯一确定。凡同条目另有 `npm run build`（`tsc -b && vite build`）或显式 `-b`/有效 `-p` 者实质检查成立，已在 `E-006` 注明。
- 建议处置：不逐条考古；由 `README` 约定 + 守卫 + CI 门禁保证**新增**记录不再出现无形态的 `tsc` 证据。若用户要求彻底清理，另开范围明确的目标。

## 结论

跨区勘误执行**可核对**，`I-008-004` 可转 `verified`。verdict 为 `conditional`：`F-001`（守卫缺口，recommended）与 `F-002`（简写未裁定，recommended）保持 open，均不阻断任何门禁，但 `F-001` 建议在下一轮加固。本意见不改 `status`/`progress`。
