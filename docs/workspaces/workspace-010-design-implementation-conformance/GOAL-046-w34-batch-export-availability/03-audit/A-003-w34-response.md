---
id: A-003-w34-response
doc: audit
parent: GOAL-046-w34-batch-export-availability
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
status: recorded
source: self
auditor: 会话编排器（`/govern` 响应）
verdict: pass
open_required: 0
---

# A-003 · W34 审计响应（A-001 self + A-002 independent）

- **source**：self（编排器响应；A-002 原文不改写，见 `03-audit/A-002-w34-independent.md`）
- **date**：2026-09-19
- **scope**：响应 A-002 的 2 required + 3 recommended，并修正 self `A-001` 中被独立审计推翻的两处结论
- **结论**：**F-001 / F-002 均 `fixed`；开放 required = 0**；`F-003`/`F-004` recommended 保留（与 self `A-001` 的 F-001/F-002 同向，已写明关闭要求）；`F-005` 反映的投影超前问题已按独立意见收回并改正。

## 响应表

| finding | 级别 | 响应 | 证据 |
|---------|------|------|------|
| `F-001`（XOR 单测名实不符、未选行） | required | **fixed** | 原例标题写「enabled」而断言 `disabled`，且**未选中任何行** —— 空选择本身就会 disabled，故「只检查 `jobs.write`」的错误实现同样能通过（独立判断正确）。已替换为**两个对称的判别性用例**：`keeps the trigger unavailable when only jobs.write is granted` 与 `... only data.export is granted`，两者都**先选中一行**再断言 `disabled` + `data-jobs-batch-export-unavailable="true"`。**变异验证（双向）**：把判据改成只查 `jobs.write` → 前者失败；只查 `data.export` → 后者失败；恢复后 10 passed |
| `F-002`（e2e 未入仓 + `custom` 推断与 harness 不一致） | required | **fixed** | ① 该 spec 已纳入版本控制（见 `E-001` §8 的 checkpoint）；② 独立指出 `custom` 不是「将来」而是**现在时**——`playwright.config.ts` 的 `customE2EModules` 当时确实不含 `admin.jobs`（与 operator config 同一种漂移）。已把 `admin.jobs` 补进该列表，使「custom = admin preset + channel.telegram」名副其实；③ 断言 `adminCapable` 含 `custom` 因此成立，并**实跑 `APP_PROFILE=custom` 复验**：1 passed；④ 复跑双 profile 全量 e2e，数字已反映该 always-on spec：mvp **18 passed / 5 skipped / 0 failed**、admin **19 passed / 4 skipped / 0 failed**（原 17/5 与 18/4 为入仓前数字，E-001 已更正） |
| `F-003`（config 守卫形态变化时 Skip 而非失败） | recommended | **保持 open（不阻断）** | 与 self `A-001` F-001 同向，独立变异确认（改 `profile: admin` → SKIP 且套件 PASS）。关闭要求已写明：config 形态变化时同步守卫，或对新形态显式断言 |
| `F-004`（客户端镜像漂移面） | recommended | **保持 open（不阻断）** | 与 self `A-001` F-002 同向；服务端仍是唯一授权方，404/403 均可诊断。关闭要求：新增服务端门禁时同步组件判据 |
| `F-005`（投影抢先写「当日关门」） | recommended | **fixed（已收回措辞）** | 独立意见成立：`roadmap.md`/`workspace.md`/`goal-tree.md` 的 W34 文本当时写成「立项并当日关门」，而独立 verdict 为 `conditional`、开放 required = 2，且把**未入仓**的 e2e 当成已锁死的守卫。已改为「立项，随 A-002 响应闭合」并把 e2e 描述为已入仓 + 已复跑；`goal-tree.md` 行与本目标 `00-meta` 的 `status`/`progress` 在响应完成后才转 `done · 4/4` |

## 对 self `A-001` 的更正（append-only，原 verdict 与 finding 原文不改写）

1. **`A-001` 的 `F-003` 关闭过早**：当时把「新增 e2e 契约」记为 `fixed`，但该文件尚未入仓，`E-001` §4 引用的是入仓前的全量数字。独立审计（`A-002` F-002）指出后已修正：文件入仓、harness 对齐、双 profile 复跑，关闭依据现在是可核对的。
2. **`A-001` 把「门禁缺一」用例算作有效守卫，属高估**：该用例不具备判别性（见 `F-001`），已由两个对称用例替换并以双向变异证明。
3. **`A-001` 对 W33 `D-001` §3 的引用属过度延伸**：独立意见指出该冻结针对**插槽布局**回落，不是「权限不足必须可见」。本响应接受该区分：真正卡住「`return null`」的是**空插槽宿主的高度契约**（`list-visual-surface` 实测 `got 0, 32, 32, 32`）；「可见 + disabled + 说明」这一结果仍然成立，但依据是「入口不撒谎 + 布局契约」，不是「W33 冻结要求必须显示」。`D-001` §3 的措辞已在 `02-execution/E-001` §5 与本响应中如实区分（不在原 `D-001` 上改写结论，仅在响应中澄清）。
4. **`E-001` §4 的数字已过期**（入仓前快照），已更新为含 always-on spec 的复跑结果。

## 响应后状态

- **开放 required = 0**（`F-001`/`F-002` 均 `fixed`）；recommended = 2（`F-003`/`F-004`，均不阻断，关闭要求已写明）。
- 复验（响应期间实跑）：vitest **121 files / 1482 tests**；typecheck exit 0；build exit 0；`go build`/`go vet` exit 0；`go test ./internal/config/ ./modules/jobs/ ./internal/handler/ -count=1` 全绿；e2e mvp **18/5/0**、admin **19/4/0**、`APP_PROFILE=custom` 可用性 spec **1 passed**。
- 未改：服务端门禁语义、Job 六态、导出契约、VP-038 `status`、workspace-038 台账、pinned 工件。

## 声明

本响应只追加响应与证据，不改写 `A-001`/`A-002` 的 verdict 与 finding 原文；不修改 Charter / VP / Goal 状态以外的治理事实。
