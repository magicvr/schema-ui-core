---
id: A-001
doc: audit-entry
record_id: A-001
source: independent
auditor: grok-4.5（独立 CLI 会话，effort high）
audit_type: close-out
scope: S5 关门范围 — C1 证据矩阵（F-V029 分母）+ C2 真实入口验证 + 对 Root GOAL-001 进入 S5 关门放行的充分性（S0–S4 已 done 是否仍支撑 exit 1–6）
verdict: conditional
status: recorded
parent: GOAL-006-s5-evidence-and-closeout
created: 2026-08-09
updated: 2026-08-09
version: 0.1.0
---

# A-001 · S5 关门独立交叉审计（close-out）

## 范围与区间

| 项 | 值 |
|----|----|
| 被审目标 | `GOAL-006-s5-evidence-and-closeout` |
| 工作区 | `workspace-007-localization-and-system-settings`（`root_goal` = `GOAL-001-localization-and-system-settings`；`canonical_scope` 已校验） |
| audit_type | `close-out` |
| scope | C1 证据矩阵（F-V029 分母）+ C2 真实入口验证 + Root 进入 S5 关门放行充分性（S0–S4 是否仍支撑 VP-007 exit 1–6） |
| 日期 | 2026-08-09 |
| auditor | grok-4.5（独立 CLI 会话，effort high） |
| source | **independent** |

**只读材料（本区 + 代码 + VP-007 退出判据）**

- `workspace.md`、`goal-tree.md`
- GOAL-006：`00-meta`、`D-001`、`E-001`、`03-audit` 索引
- Root GOAL-001：`00-meta` 路线图、`attachments/F-V029-coverage-table-s0-freeze.md`、`attachments/S5-evidence-matrix.md`
- VP-007 方向级退出判据 1–6（只读；**未**写 vision reviews）
- 代码：`apps/web/e2e/localization.spec.ts`、i18n/startup-config/error 相关 vitest、`apps/api/internal/handler` branding/settings/localize/error 测试、`kernel/profile.go`
- S1–S4 目标 meta/自审（仅核对 done 主张与 exit 支撑链，不改其状态）

**工作区绑定核对**

| 项 | 结果 |
|----|------|
| `root_goal` / canonical | 与 goal-tree、Root 文件夹一致 |
| `plan_refs` / `primary_plan` | `VP-007-localization-and-system-settings` |
| `shared_materials_catalog` | `none`（无共享资料引用风险） |
| 跨工作区上下文 | **未读取、未比较** |

**P-005 信息门禁（本 scope）**

| ID | 状态 | 审计结论 |
|----|------|----------|
| GOAL-006 `I-001` | closed | 关门输入登记合理；**不等于** C1–C4 证据已全部可独立复核 |
| Root `I-L10N-001`～`005` | verified（D-002） | 与 S0 台账一致；无到期未决 required 信息项阻断本 scope 审阅 |

## 成果（有证据）

### 治理与阶段链

1. **工作区/树一致**：S0–S4 子目标均为 `done`；GOAL-006 `active` `progress: 2/4`（C1/C2 自标完成，C3/C4 未完成）与 `goal-tree.md` 叙述一致。
2. **S5 流程口径冻结**：`D-001` 明确 C1→C2→C3→C4、矩阵复用 F-V029、S5 关门 = independent、禁止代表性收缩分母——与 Root D-002 §4 / VP-007 exit 6 对齐。
3. **C1 矩阵落盘形态正确**：`GOAL-001/.../attachments/S5-evidence-matrix.md` 存在；行分母 8 格（zh-CN/en-US × mvp/admin × 匿名/已认证）；列覆盖固定 UI、12 pageId、M1–M4、权限正反、缺失翻译、配置刷新、错误回退；mvp 上 `settings`/`activity` 与 U7 的 N/A 理由与 `kernel/profile.go`（mvp 无 `admin.settings`/`admin.activity`）一致。
4. **分母结构证据（可复跑）**：独立会话复跑 web 子集 **9 files / 79 tests passed**（`locale`/`catalog`/`format`/`runtime`/`schema-keys.structural`/`ui-bilingual`/`locale-switcher`/`startup-config`/`error-localization`），与 E-001「79/79」声称一致。
5. **API 相关测试可复跑**：`go test ./internal/handler/ -run Branding|Settings|Localize|ErrorContract|Auth|Login` **ok**；`go test ./internal/kernel/ -run Profile|MVP` **ok**。
6. **C2 浏览器真实入口（部分可核对）**：
   - 源码 `apps/web/e2e/localization.spec.ts` 覆盖 admin：M1 语种切换 + `lang=zh-CN`、登录失败 envelope 协商（`UNAUTHORIZED` + zh message + `messageKey`）、登录后 overview 中文标题、M3 settings 保存投影站点标题、四类工具条可见、零 `pageerror`。
   - 产物 `apps/web/test-results/s5-settings-zh.png` 与 `.last-run.json`（`status: passed`）存在。
7. **Profile 边界代码事实**：`apps/api/internal/kernel/profile.go` 中 mvp 模块集不含 settings/activity；admin 包含——支撑 exit 4 与矩阵 N/A 语义（非 e2e 层）。
8. **S0–S4 阶段自审链**：GOAL-002～005 均有 `source: self` 且 `verdict: pass`、开放 required = 0 的落盘意见；Root S0 有 independent A-001 → A-002 响应闭合记录。信息门禁 `I-L10N-001`～`005` 均 `verified`。

### 对照成功标准 / VP-007 exit 1–6（本轮独立判断）

| 判据 | 主张来源 | 独立可核对强度 | 结论摘要 |
|------|----------|----------------|----------|
| **exit 1** 语种解析/用户控制/`lang`/格式化 | S1 单元 + e2e zh 切换 | **中–高** | `locale.test.ts`/`runtime`/`format.test.ts` + e2e `lang` 可支撑方向；**缺** en-US 浏览器端完整路径与 mvp 真实入口 |
| **exit 2** 可维护翻译面 + 分母无硬编码英文路径 | S2 structural + ui-bilingual + F-V029/S5 矩阵 | **中** | 12 schema 全文本 `*Key` + 双语 catalog 对称（`schema-keys.structural.test.ts`）为强结构证据；**运行时双语渲染** primarily 固定 UI + users；roles 与多枚 example page 仍主要依赖 structural，F-V029 原标注「渲染证据 S5」未充分兑现 |
| **exit 3** 四类设置产品面 | S3 + startup-config + e2e M3 | **高** | Go settings/branding + web startup-config + e2e 保存投影可支撑 |
| **exit 4** 双 Profile + 公开启动边界 | profile 编译默认 + branding 公开读 | **中–高** | 模块集与 handler 测试支撑；**缺** mvp Profile 真实入口/浏览器反证 settings 不可达（单元/composition 侧有相关事实，e2e 仅 admin） |
| **exit 5** 错误码 + 前端保底 + 路径 (a) | S4 + e2e Accept-Language | **高** | handler localize/error_contract + error-localization + e2e 登录失败协商可支撑；I-L10N-004 路径 (a) 实施证据链成立 |
| **exit 6** 同分异分母矩阵 + required=0 + 用户确认 | S5 矩阵 + 本审计 + C4 | **不足以致无条件放行** | 矩阵**形态**合规；多格证据可复跑；但 C2 双启动体一致与 unit/build 日志 **耐久路径缺失**；本 A-001 产生开放 required；用户书面关门（C4）尚未发生 |

## Findings

### F-001 · required · high · C2 双启动 / scratch 捕获不可独立复核

- **主张**：E-001 / S5 矩阵声称 API 真实构建后 ≥2 次启动，`GET /api/branding` 响应体完全一致，并捕获 `{SCRATCH}/s5-launch/run1.json`、`run2.json`、`compare.log`、`go-build.log`、`web-build.log`、`e2e-localization.log` 与 `{SCRATCH}/s5-tests/*`。
- **独立核对**：
  - 在仓库内 `docs/workspaces/workspace-007-...`、`apps/api/bin`、`apps/web/test-results` 及常见 scratch 目录 **未找到** `run1.json` / `run2.json` / `compare.log` / `web-build.log` / `e2e-localization.log` / `s5-tests` 日志。
  - GOAL-006 `attachments/` 为空（仅 README）。
  - `apps/api/bin/out.log` / `out2.log` 时间戳为 **2026-07-31**，内容为 server starting，**不能**作为 2026-08-09 S5 branding 体一致证据。
  - 可核对缓解：`apps/api/bin/schema-ui-core-api.exe` 修改时间为 2026-08-09；e2e 截图与 spec 存在；本会话复跑 79 vitest + handler 测试全绿。
- **判定**：双启动 body 一致与 web production build 的 **耐久、可重复核对证据不足**。不得仅凭叙述将 C2「API dual-run 体一致」视为已关闭的关门证据。
- **建议闭合路径**：
  - `fixed`：将 dual-run 的两次 branding JSON + 比对结果 + build 日志（或等价可重复脚本输出）落入 GOAL-006 `attachments/`（或 Root attachments 并在 E 条更新**仓库内**路径），并保证路径可从 E-001/矩阵点击核对；**或**在受控环境下重跑 dual-run 并落盘。
  - 若仅接受「可重跑命令」而无产物：须用户书面 `accepted-residual`（范围=耐久产物缺失、缓解=命令可重跑、责任人、复审触发）。

### F-002 · required · med · 分母 pageId 渲染层证据未齐（F-V029「S5 渲染」承诺）

- **主张**：S5 矩阵将 12 pageId 非 N/A 单元格标 ✓；F-V029 对 overview/roles 等写明「渲染证据 S5 矩阵」。
- **独立核对**：
  - **强渲染/双语**：固定 UI（`ui-bilingual`）、users 列表/写表单/确认（`ui-bilingual` + `schema-crud`）、overview 中文标题（e2e admin）、settings 四类（`startup-config` + e2e）。
  - **主要为 structural**：`data-display` / `data-table` / `search-form-table` / `form-controls` / `form-with-reactions` / `form-with-upload` / `roles` / `activity` 等在 S5 矩阵证据列指向 `schema-keys.structural.test.ts`（key 完备 + 双语 catalog 存在），**无**对应 zh-CN/en-US 页面渲染断言。
  - `admin-list-batch` 另引 `representative-pages.integration.test.tsx`（路径实际为 `apps/web/src/app/representative-pages.integration.test.tsx`，文件存在；**未**证明其以双 locale 解析 catalog 文案）。
  - F-V029 对 **roles** 明确「渲染证据 S5 矩阵」——当前矩阵仍无 roles 双语渲染路径。
- **判定**：exit 2 / exit 6 允许测试名作证据，但 **不得把「catalog 键齐全」静默等同于「运行时不存在硬编码英文完成路径」的渲染证明**。在 S5 关门 scope 下，对 F-V029 已预告 S5 补渲染的格子，证据等级不足。
- **建议闭合路径**：
  - `fixed`：至少为 **roles**（F-V029 明文）及矩阵中仅 structural 的 pageId 补可核对渲染证据（vitest 双语渲染或 e2e 抽样），并回填矩阵单元格路径；**或**
  - 用户书面 `accepted-residual`：明确接受「结构完备 + 代表页渲染（users/overview/settings）替代全并集渲染」的残余范围、缓解与复审触发（**禁止**无裁决的静默收缩）。

### F-003 · required · med · 双 Profile 真实入口矩阵不对称（mvp 无真实入口验证）

- **主张**：exit 6 行分母含 `mvp`；E-001 C2 以 `APP_PROFILE=admin` 跑 playwright；矩阵 mvp 行大量 ✓ 复用与 admin 相同的前端单元测试。
- **独立核对**：
  - e2e **仅 admin**；无 mvp 启动、无 mvp 下 settings 不可达的浏览器反证、无 mvp 匿名/已认证真实入口捕获。
  - 共享运行时 + `schema-keys`/`ui-bilingual` 对「同一前端 build」有合理性，但 **不等于** 已用真实入口验证 mvp 行。
  - Profile 编译边界有 Go 证据（支撑 exit 4 的模块集半边），不能单独关闭 exit 6 的 mvp 行「真实入口」期望（D-001 C2 / 计划验证步骤 4 语境）。
- **判定**：在宣称「双 Profile 验证矩阵与真实入口」已完成并支撑 Root 关门时，**mvp 真实入口侧证据不足**。
- **建议闭合路径**：
  - `fixed`：mvp Profile 至少一次真实入口（API branding 可读 + web 构建产物加载；settings 路由/导航不可达的可核对断言）并写入矩阵/E 条；**或**
  - 用户书面 `accepted-residual`：接受「mvp 行以共享单元 + profile 模块集测试代替浏览器真实入口」的范围与复审条件。

### F-004 · recommended · low · M2 浏览器路径未覆盖写表单

- F-V029 M2 要求有权限账号完成至少一次用户或角色写表单；e2e 仅登录后 overview。
- 单元层 `ui-bilingual` / `schema-crud` 覆盖 users 写路径，**不**将本条升为 required，但 Root 关门材料宜在矩阵/M2 证据列区分「单元写表单」与「浏览器写表单」。

### F-005 · recommended · low · 矩阵证据路径未写全限定文件路径

- 多处仅写测试文件名（如 `ui-bilingual.test.tsx`），依赖读者已知 `apps/web/src/...`；`representative-pages.integration.test.tsx` 实际在 `src/app/` 而非 `src/renderer/`。
- 建议回填稳定相对路径，降低复审歧义（非阻断，除非与 F-001/F-002 一并整改）。

## 必改项汇总

| ID | 级别 | 摘要 | 阻断 |
|----|------|------|------|
| **F-001** | required · high | C2 dual-run / build / unit 刷新的 scratch 捕获在仓库内不可复核；不得无条件采信 body 一致叙事 | **是** — 阻断无条件 S5/Root 关门放行 |
| **F-002** | required · med | 并集 pageId（尤 roles 及多枚 example 页）缺少 F-V029 预告的 S5 渲染层证据；structural ≠ 渲染证明 | **是** — 阻断无条件 exit 2/6 关闭 |
| **F-003** | required · med | mvp 无真实入口验证，矩阵 mvp 行与「双 Profile 真实入口」主张不对称 | **是** — 阻断无条件 exit 6 / C2 完整认定 |
| F-004 | recommended | M2 e2e 未写表单 | 否 |
| F-005 | recommended | 证据路径宜全限定 | 否 |

**开放 required findings = 3（F-001、F-002、F-003）**。按 P-003：存在未合法闭合 required 时，**不得**将 GOAL-006 / Root 标 `done`，**不得**无条件放行 S5 关门。

## 与既有意见的异同

| 既有 | 关系 |
|------|------|
| GOAL-006 此前无正式 A 条目 | 本 A-001 为序列首条 |
| S1–S4 self pass | **不否定**阶段内检查点在各自 scope 的自审结论；本意见 scope 是 **S5 关门与 exit 1–6 充分性**，标准更严 |
| Root S0 independent A-001（conditional → 已响应闭合） | 模式一致：independent 可出 conditional + required；闭合归 `/govern` + 用户裁决 |
| S2/S3/S4 自审将浏览器 e2e 留待 S5 | 本审计确认 e2e **部分兑现**（admin zh M1/M3/错误协商），但 **未**完整兑现双 Profile × 双 locale × 分母渲染 |

## 结论

**verdict: conditional**

- C1 矩阵**形态**与 F-V029 分母对齐，N/A 边界正确，大量单元格有可复跑单元/Go 证据；C2 **浏览器 admin 路径**有 spec + 截图 + last-run；S0–S4 对 exit 1/3/5 与 exit 4 模块半边仍有可核对支撑。
- **不可无条件放行** Root/S5 关门：F-001（C2 耐久 dual-run/build 证据缺失）、F-002（分母渲染层缺口）、F-003（mvp 真实入口缺失）为开放 required。exit 6 的「开放 required = 0 + 用户确认」尚未满足（本意见新增 required；C4 未执行）。

### 建议给编排器 / 用户的下一步

1. 用 **`/govern` 响应 A-001**：逐条处理 F-001～F-003（`fixed` 优先；residual/overruled 须用户书面留痕）。
2. 建议优先序：`F-001` 落盘 dual-run 产物 → `F-003` mvp 最小真实入口或 residual → `F-002` roles/关键 page 渲染补证或 residual。
3. required 全闭合后，再执行 GOAL-006 **C4**（用户书面关门确认 → Root done / VP-007 关门记录 / goal-tree 同步）。**不要**在本审计未响应前勾选 C3/C4 或改 status。

## 声明

本意见 `source: independent`，**不修改**任何目标 `status` / `progress` / 检查点勾选 / 方案正文 / `goal-tree`。  
响应、finding 闭合与阶段放行由用户通过 **`/govern`** 处理。
