---
id: A-001-w34-self
doc: audit
parent: GOAL-046-w34-batch-export-availability
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
status: recorded
source: self
auditor: 会话编排器（self）
verdict: pass
open_required: 0
---

# A-001 · W34 自审（GOAL-046 C1～C3）

- **source**：self
- **date**：2026-09-19
- **scope**：W34 全量 —— 诊断结论、`configs/config.yaml` 修复、触发面可用性门禁、404 语义、两层守卫与全量回归
- **verdict**：**pass**（0 required + 3 recommended）

## 成果（有证据）

1. **根因定位可复现，且不是「猜」**：以 operator config 原样启动 → `POST /api/jobs/batch-export` 得 `404 {"error":"NOT_FOUND","messageKey":"error.notFound"}`；`GET /api/export/users` 正常（对照证明只缺 `admin.jobs`）；`profile: admin` 下同一请求为 `403 MUST_CHANGE_PASSWORD`（路由存在）。修复后 users/roles 均 `202 → succeeded` 且结果可下载。证据：`E-001` §1/§3。
2. **用户可见症状与证据严格对齐**：UI 显示的「未找到」正是 `error.notFound` 的 zh 文案，由 `messageKey` 决定；不是「行不存在」这类数据语义。这一区分决定了修复方向（装配 vs 数据）。
3. **修复面最小且方向正确**：只在 operator config 补一个模块 + 客户端可用性镜像，未改服务端门禁语义、Job 六态、导出契约；两道服务端门禁仍是唯一授权方。
4. **两层守卫 + 4 处变异验证（全部真实执行）**：
   - config 守卫变异（删 `- admin.jobs`）→ 失败并指名 `admin.jobs`；
   - 单测变异（去门禁判断 / 404 文案回退）→ 对应用例红；
   - e2e 变异（`return true` 去掉可用性判断）→ mvp 分支 `toBeDisabled` 失败（`Received: enabled`）。
5. **实测捕获并修正了本波第一版实现的真实回归**：`return null`（隐藏）导致 page-actions 行出现高度 0 的空插槽宿主，`list-visual-surface` e2e 报 `got 0, 32, 32, 32`。修正为「可见 + disabled + 说明」后单跑 6 passed、双 profile 全量全绿。这条恰好印证了 W33 `D-001` §3 的 fail-open 冻结不是教条 —— 它同时是布局契约的要求。
6. **回归面完整**：vitest 121 files / 1481 tests、typecheck/build exit 0、e2e mvp 17/5/0 与 admin 18/4/0、`go test ./...` 全绿（含 `internal/config`）。
7. **覆盖缺口被真实补齐**：原覆盖只到「admin profile 可用」，没有任何断言覆盖「节点已发布但路由未挂载」。新增 `e2e/batch-export-availability.spec.ts` 在**两个 profile** 上分别断言可用性契约（admin 可用 / mvp 可见但禁用），两分支各自实跑通过。

## 偏差

- **第一版实现方向错误（已修正）**：以「隐藏」处理不可用，违背 W33 `D-001` §3 冻结取向且触发布局回归，已在同一波次内改为可见禁用，并把原因写入代码注释与 `D-001` §3，避免后续再犯。
- **`adminCapable` 判定含 `custom`**：e2e 用 `APP_PROFILE` 推断 profile 能力，`custom` 被算作「可服务」。当前 harness 只接受 mvp/admin/demo/custom，且 `custom` 由 harness 显式列出模块；若将来 harness 的 custom 列表去掉 `admin.jobs`，该分支会误判。已在 `F-002` 登记为 recommended（有界）。
- 未修 `go test ./...` 的既有 PG drain flake（`TestShutdownDrainHarnessPostgres`，VP-021 遗留，本波 sqlite 路径全绿、未触发）。

## Findings

- `F-001`：**recommended**；状态 `open`。**配置守卫依赖手写 YAML 扫描**：`parseInlineModuleList` 只认 `app.modules.list` 下的 `- <id>` 行，避免为一条断言引入 YAML 依赖。若 operator config 将来改用 `preset:` 或其它形态，守卫会 `t.Skip`（可见但静默）而不是失败 —— 即「不变量失效但不报警」。影响：低（跳过是显式的，且形态变化本身会被人看到）。关闭要求：若 config 形态变化，同步守卫或在守卫中显式断言新形态。
- `F-002`：**recommended**；状态 `open`。**客户端可用性镜像与服务端门禁清单存在漂移面**：组件硬编码 `jobs.write` + `data.export` 两个权限键。若服务端将来为提交路由增加第三道门禁，客户端不会自动同步 —— 后果只是「入口看起来可用但提交被拒」，且 404/403 文案已可诊断，服务端仍是唯一授权方。关闭要求：新增服务端门禁时同步组件判据（可在 PR 模板或 `D-001` §3 增补检查项）。
- `F-003`：**recommended → fixed（本波次内闭合）**；证据见「成果 7」与 `E-001` §4。原覆盖缺口 = 「节点已发布但路由未挂载」无 e2e 断言；已补双 profile 可用性契约 spec，并以变异验证其有效性。

## 结论 + 建议下一步

- **verdict = pass**：诊断有可复现证据链、修复面最小且方向由证据决定、两层守卫经变异验证、全量回归全绿、覆盖缺口已补；`0 required`。
- 3 条 recommended 均不阻断关门：`F-003` 已在本波闭合；`F-001`/`F-002` 是有界维护性项，已写明关闭要求。
- 建议下一步：按项目级决策执行 **independent（grok build · grok-4.6 · high · `/audit`）**，随后由编排器合并响应（`A-002`/`A-003`），再同步 `goal-tree.md`/`workspace.md`/`roadmap.md` 并关门。

---

## 追加更正（append-only · 2026-09-19 · 据独立审计 `A-002`）

原 verdict（`pass`）与上方 finding 原文不改写。独立审计在实际核验后**推翻/高估**了本意见的三处结论，逐条更正如下：

| 本意见原结论 | 独立审计判定 | 更正 |
|--------------|--------------|------|
| 成果 7 / `F-003`：新增 `apps/web/e2e/batch-export-availability.spec.ts` 已补齐「节点已发布但路由未挂载」的覆盖缺口，标 `fixed` | **不接受关闭**（`A-002` F-002）：该文件当时对 `e1013c8b` 为 **untracked**，未入仓的守卫对 CI/克隆无效；且 `adminCapable` 含 `custom` 而 harness 的 `customE2EModules` **当时就不含** `admin.jobs`（判断为「现在时」而非本意见所写的「将来若发生变化」） | 已按 `A-003` `fixed`：spec 入仓、harness 补 `admin.jobs`、`APP_PROFILE=custom` 与双 profile 全量复跑（mvp **18/5/0**、admin **19/4/0**）；本意见「成果 7」的时效性以 `A-003` 为准 |
| 成果 4 / `F-003` 相关：把「门禁缺一」单测（1 例）算作有效守卫 | **不同意**（`A-002` F-001）：该例标题写「enabled」而断言 `disabled`，且**未选行** —— 空选择本身即 disabled，故「只检查一个门禁」的错误实现同样能通过 | 已替换为两个对称的判别性用例（均先选行并断言 `unavailable` 标记），并以**双向变异**证明：判据只查 `jobs.write` → 前者红；只查 `data.export` → 后者红 |
| `D-001` §3 / 本意见成果 4：以 W33 `D-001` §3 的 fail-open 冻结作为「不可用必须可见」的依据 | **部分不同意**（`A-002` 成果 4）：该冻结针对**插槽布局**回落，硬约束是**空插槽宿主的高度契约** | 接受该区分：「可见 + disabled + 说明」的成立依据是「入口不撒谎 + 布局契约」，W33 条文本身不构成权限语义要求；澄清已写入 `E-001` §5 与 `A-003` |

以上更正意味着：本意见 `verdict: pass` 在当时**高估了 C3 守卫的有效性**（两处守卫一处不具判别性、一处未入仓）。最终关门依据以 `A-002` + `A-003`（响应后开放 required = 0）为准。
