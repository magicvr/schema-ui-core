---
id: GOAL-024-w16-user-perspective-improvements
doc: audit-entry
record_id: A-005
source: independent
scope: GOAL-024 全目标完成情况 / W16-F01～W16-F10 落地核实 / S5 关门主张
verdict: fail
status: recorded
parent: GOAL-001-design-implementation-conformance
created: 2026-08-18
updated: 2026-08-18
version: 0.1.1
---

# A-005 · 独立审计 · GOAL-024 完成情况与 W16-F01～W16-F10 落地核实（2026-08-18）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型** / **scope**：close-out + execution-facts · GOAL-024 全目标完成情况、W16-F01～W16-F10 对照 D-001/D-002 与批 A/B/C 冻结方案的落地核验
- **verdict**：**fail**

## 范围与区间

- **工作区**：`workspace-010-design-implementation-conformance` · Root `GOAL-001-design-implementation-conformance` · canonical `docs/workspaces/workspace-010-design-implementation-conformance/` · `shared_materials_catalog: none`（无共享资料引用可核）
- **covered**：
  - GOAL-024 目标定义、D-001 台账、D-002 技术方案、D-003 分批、I-001、E-001/E-002、A-001～A-004
  - 下级整改子目标 GOAL-025 / GOAL-026 / GOAL-027 的 meta、冻结方案、实施事实与关门意见
  - W16-F01～W16-F10 对照代码、schema、迁移与定向测试
- **excluded**：
  - 未读取或比较其他工作区上下文；D-001 排除项（workspace-011 已规划的组织/公告/邮件短信/登录日志）仅作本区边界记录
  - 未重跑全量 Go/Web/e2e；本轮独立重跑了 W16 定向 Go 套件
  - 未做浏览器端到端点验

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| 工作区绑定合格 | `workspace.md`：id / root_goal / canonical_scope / `plan_refs`+`primary_plan`=VP-010；资料目录 `none` |
| S1 台账与排除项已冻结 | [D-001](../01-decision/D-001-w16-improvement-inventory.md)、[E-001](../02-execution/E-001-goal-establishment.md) |
| S2 技术方案与 I-001 关闭 | [D-002](../01-decision/D-002-w16-technical-design.md)；`00-meta` I-001 = verified（Renderer custom-node + schema 扩展） |
| S3 分批规划与渐进添加 | [D-003](../01-decision/D-003-w16-batch-subgoals.md)；子目标文件夹 GOAL-025/026/027 均存在 |
| 批 A 主体实现（F01/F07/F08） | migration 0038 + `auth.go` 门禁 + `force-password-change.tsx`；`POST /api/account/sessions/revoke-others` + `account-session-toolbar.tsx`（含 `reloadList`）；登录页 `data-captcha-refresh` + MFA 复制/下载 |
| 批 A 历史 required 已在子目标闭合 | GOAL-025 A-001 F-001/F-002 由 A-002 标 `fixed`；`account_self.go` 拒绝同密码；toolbar 调用 `crud?.reloadList()` |
| 批 B 后端契约（F02 URL 字段 / F03 API / F04 列格式） | `filelibrary.go` `downloadUrl`；`GET /api/import/{resource}/template` + `fieldErrors`；`wallet.json` `format: "currency"` + `formatCents` |
| 批 C 主体实现（F05 API / F06 下拉 / F09 Badge / F10 页脚） | `POST /api/scheduled-tasks/cron/preview`；`monitoring-auto-refresh.tsx` 0/5/10/30s；`dict_entries.badge_style` + `badgeStyleField`；settings/`/api/branding`/`App.tsx` footer |
| 本轮独立重跑 W16 定向 Go | `apps/api`：`go test ./internal/handler/ -run TestForcedPasswordChangeGateAndReissue\|TestRevokeOthers…\|TestImportTemplateUsers\|TestImportFieldErrors\|TestFileLibraryRowsExposeDownloadUrl\|TestCronPreviewEndpoint\|TestDictEntryBadgeStylePersists\|TestSettingsFooterFieldsPersist` → **ok**（10.555s） |
| 规划阶段自审 | A-001/A-002 `source: self` · pass（当时尚未实施，范围匹配） |
| 信息项 | 父目标 I-001 verified；子目标 I 项均标 closed。无到期未关闭 required 信息项 |

## 对照成功标准

| 阶段 / 检查点 | 定义 | 本轮判定 | 证据 |
|---------------|------|----------|------|
| S1 台账与范围冻结 | W16-F01～F10 + 排除项 | **完成** | D-001 |
| S2 技术方案；关闭 I-001 | 契约/路由/前端路径 | **完成** | D-002 |
| S3 分批与子目标规划 | 批 A/B/C 渐进添加 | **完成** | D-003 + 三个子目标文件夹 |
| S4 规划就绪自审 | 实施前审核 | **完成**（规划 scope） | A-002 |
| R1 批 A | F01/F07/F08 可核对落地 | **基本完成** | GOAL-025 done；代码与定向测试支持 |
| R2 批 B | F02/F03/F04 按冻结方案落地 | **未完成** | 见 F-001、F-002；F04 部分见 F-003 |
| R3 批 C | F05/F06/F09/F10 按冻结方案落地 | **大部分完成** | F06/F09/F10 可核对；F05 能力在但偏离冻结交互，见 F-004 |
| S5 关门 | 10 项改进对用户成立 + 终审 | **主张不成立** | A-003/A-004 宣称 10 项完成；本轮核验 F02/F03 用户问题仍在 |

## Findings

### F-001 · W16-F02「在线直接预览」按当前实现不可用

| 字段 | 值 |
|------|-----|
| level | required |
| severity | high |
| status | fixed |
| closure | A-007 复审可闭合 + A-008 / D-004（2026-08-18） |
| evidence | `apps/web/src/renderer/render.tsx` `library.preview` 仅 `window.open(url)`；`apps/api/internal/auth/auth.go` Middleware 只认 `Authorization: Bearer`；`apps/api/internal/handler/filelibrary.go` 下载端点强制 `Content-Disposition: attachment`。GOAL-026 D-001 要求图片 Lightbox + PDF/文本新标签预览。 |

D-001 问题是「无法在浏览器内直接预览，只能下载」。现实现：

1. 新标签不会带 access token → 生产路径下该 URL 会 401，而不是打开文件。
2. 即便鉴权过去，响应也是 `attachment`，浏览器会下载而不是预览。
3. 未实现冻结方案中的 `file-preview-lightbox`。
4. `library.copyLink` 复制的是相对 API 路径 `/api/library/files/{id}/download`，同样无法作为可打开的完整访问 URL。

对照：同文件里的 `library.download` 才走 `authFetch` + blob。预览没有复用这条已工作的鉴权通道。A-004 把 `download-behavior.test.tsx` 记成 F02 预览证据，该文件测的是 `library.download`，不是 preview。

### F-002 · W16-F03 前端「模板下载 + 逐行错误表」未落地

| 字段 | 值 |
|------|-----|
| level | required |
| severity | high |
| status | fixed |
| closure | A-007 复审可闭合 + A-008 / D-004（2026-08-18） |
| evidence | `apps/api/internal/modules/users/schema/users.json` `openImport` 表单只有 `file` 上传字段，无模板链接；`submitImport.onSuccess.behavior = reload`；`apps/web/src/renderer/render.tsx` `submitForm` 在 `result.ok` 时只 toast + `reloadList` + 关模态。导入接口在部分行失败时仍 **HTTP 200** 返回 `fieldErrors`。GOAL-026 D-001 §3 明确要求模态内「下载 CSV 模板」和 `rowNumber/field/reason` 表格/列表。 |

后端契约存在且本轮 Go 测试通过（`TestImportTemplateUsers`、`TestImportFieldErrors`）。用户原问题是「不知表头」和「无法定位第几行哪个字段」。当前导入模态仍只有上传控件；成功响应里的 `fieldErrors` 不会被渲染。GOAL-014 的 form `fieldErrors` echo 只作用于失败响应且映射到同名表单字段，无法展示导入行号明细。

### F-003 · W16-F04 调账仅有警示文案，无二次确认

| 字段 | 值 |
|------|-----|
| level | recommended |
| severity | med |
| status | fixed |
| closure | A-006 实施 + A-007 同意闭合 |
| evidence | `apps/web/src/renderer/form-controls.tsx` 对 `amountDelta` 仅渲染 `schema.wallet.adjustWarning`；form-controls 无 confirm。GOAL-024 D-002 / GOAL-026 D-001 均写「高亮警示与二次确认」。货币列 `format: "currency"` / `formatCents` 已落地。 |

### F-004 · W16-F05 Cron 预览未按冻结方案接到表单字段，描述也不是中文人话

| 字段 | 值 |
|------|-----|
| level | recommended |
| severity | med |
| status | fixed |
| closure | A-009 / GOAL-028 A-001（W17 字段绑定 + 中文 describeCron） |
| evidence | `cron-preview.tsx` 是独立输入+提交，无防抖、不读取任务表单 `cron` 字段；`scheduled-tasks.json` 把它挂在页面块而非 Cron 字段下方。`describeCron` 只返回 `"every minute"` / `"every hour at minute N"` / `"cron schedule (5-field)"`。D-001 要求中文人话；GOAL-027 D-001 要求字段下方防抖预览。API 与未来 3 次 `nextRuns` 本身存在且测试通过。 |

### F-005 · 父目标台账未与关门状态对齐

| 字段 | 值 |
|------|-----|
| level | recommended |
| severity | low |
| status | fixed |
| closure | A-006 台账收口 + A-007 同意闭合 |
| evidence | `00-meta.md` 子目标表仍写 GOAL-025 `active · 0/4`；`01-decision.md` / `02-execution.md` / `03-audit.md` frontmatter `status: active`。`02-execution.md` 只有 E-001/E-002（规划），没有 R1～R3/S5 实施完成事实条目。goal-tree / workspace 波次表已写 done 8/8。 |

此条不单独构成 fail；说明 S5 台账收口不完整。A-004 F-001/F-002 已点出部分残留，本条合并并补上执行索引缺口。

## 必改项汇总

1. **F-001（required / high）**：按鉴权通道实现真正的在线预览（复用已鉴权 fetch + blob/Lightbox 或 inline 预览端点）；复制链接给出当前会话可用的完整 URL，而不是裸相对下载路径。
2. **F-002（required / high）**：导入模态增加模板下载入口；导入 200 响应中的 `fieldErrors` 必须以行号/字段/原因展示，不能只 reload。

未闭合 required = **2**（审计当时）。**编排闭合（2026-08-18 · A-008）**：F-001 / F-002 → `fixed`；F-003 / F-005 → `fixed`；F-004 保持 recommended `open`。现开放 required = **0**。

## 与既有意见的异同

| 意见 | 关系 |
|------|------|
| GOAL-024 A-001 / A-002 self | 规划阶段 pass，本轮同意（当时明确排除实施）。 |
| GOAL-024 A-003 self close-out pass | **不同意**。A-003 只引用子目标 `done` 与「全量回归」，未按 D-001 用户问题核对 F02/F03 前端。子目标 done ≠ 10 项用户问题已解。 |
| GOAL-024 A-004 independent pass（gemini-3.7-flash-high） | **不同意其完成结论**。A-004 把 F02 预览记成已落地、把 `download-behavior.test.tsx` 当作预览证据、把 F03 前端记成 `form-controls`/`vitest 覆盖`。本轮按代码路径复核不成立。A-004 的 recommended 文档残留（00-meta 子表、索引 status）本轮并入 F-005。 |
| GOAL-025 A-001 independent conditional | 其 required F-001/F-002 代码侧已补；recommended F-003/F-005/F-006/F-007 仍开放，不升格为本轮父目标 required。 |
| GOAL-026 / GOAL-027 A-001 self pass | 批 B/C 关门自审同样未核对 F02 预览鉴权与 F03 导入模态，证据链偏「有测试即完成」。 |

## 结论 + 建议给编排器/用户的下一步

规划与分批治理（S1～S4）成立；批 A 与批 C 的大部分项有可核对实现。但 **W16 对真实用户承诺的 10 项里，F02 预览与 F03 导入反馈在产品面上仍未解决**，与 `status: done` / A-003·A-004 `pass` 的关门主张冲突。

建议 `/govern`：

1. 将本意见纳入 GOAL-024（及受影响的 GOAL-026）响应，**不要**在 F-001/F-002 未按三路径闭合前维持或重申关门。
2. 优先在 GOAL-026 范围修正 F02/F03；F-003/F-004 可同批或另记 residual。
3. 响应时补父目标执行事实（R1～R3/S5）并同步 `00-meta` 子目标表。
4. A-004 与本意见在关门结论上冲突：按 P-004 展示后由用户裁决采信哪一份；本审计建议以本条代码核验为准。

## 声明

本意见为独立交叉审计（`source: independent`），不修改目标 `status` / `progress` / goal-tree。响应、finding 闭合与是否重新开门由 `/govern` 处理。
