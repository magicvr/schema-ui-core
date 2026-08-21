---
id: GOAL-024-w16-user-perspective-improvements
doc: audit-entry
record_id: A-007
source: independent
scope: A-005 F-001～F-005 关闭复审（finding-closure）
verdict: pass
status: recorded
parent: GOAL-001-design-implementation-conformance
created: 2026-08-18
updated: 2026-08-18
version: 0.1.0
---

# A-007 · 独立审计 · A-005 finding 关闭复审（2026-08-18）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型** / **scope**：finding-closure · 复审 A-005 required F-001/F-002 与 recommended F-003/F-004/F-005；对照 A-006 关闭声明与 `fbe7c40` 代码
- **verdict**：**pass**

## 范围与区间

- **工作区**：`workspace-010-design-implementation-conformance` · Root `GOAL-001-design-implementation-conformance` · canonical 范围匹配 · `shared_materials_catalog: none`
- **covered**：A-005 五条 findings 的关闭证据；A-006 self 响应；GOAL-026 A-002；E-003；相关前端实现与定向测试
- **excluded**：未重审批 A/C 未改动代码；未做浏览器端到端点验；未重跑全量 Go/Web
- **信息项**：父目标 I-001 仍为 verified；无到期 required 信息门禁

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| P-004 冲突已留痕 | A-006：用户书面采纳 A-005，先修正再关门 |
| A-005 F-001 预览鉴权通道已改 | `render.tsx` `library.preview`/`library.copyLink` 经 `fetcher` 拉 blob 再 `createObjectURL`；不再 `window.open` 裸下载路径 |
| A-005 F-001 有可重复测试 | `download-behavior.test.tsx` `library.preview fetches a blob and opens its object URL`；本轮 `npx vitest run src/renderer/download-behavior.test.tsx` → **7/7 passed** |
| A-005 F-002 模板入口存在 | `import-template-download.tsx` 用 `authFetch` 拉 `GET /api/import/users/template` 并触发下载；`users.json` 挂 `import-template-block`；`main.tsx` 注册；i18n `importTemplate.download` |
| A-005 F-002 200 `fieldErrors` 会留在模态并列出 | `runRequest` 解析成功体 `fieldErrors`（含 `rowNumber`）；`submitForm` 在有 fieldErrors 时不关模态；`FormView` `data-import-error-rows` 渲染 `#行号 / 字段 / 原因` |
| A-005 F-003 调账二次确认 | `FormView` 对 `adjust-wallet-form` 的负数或 `abs>100000` 调用 `window.confirm`；`wallet.json` 该 form id 存在；i18n `schema.wallet.adjustConfirm` |
| A-005 F-005 台账收口 | `00-meta` 子目标表 GOAL-025 已 `done · 4/4`；`01/02/03` 索引 `status: done`；`E-003` 已补 R1～R3/S5 |

## 对照 A-005 关闭声明

| A-005 finding | A-006 声明 | 本轮判定 | 说明 |
|---------------|------------|----------|------|
| F-001 required · F02 预览/复制 | fixed | **可闭合（fixed）** | 原 401 + `attachment` 路径已不走裸 `window.open`。复制对象是 `blob:` URL，见本条 F-001 recommended |
| F-002 required · F03 模板 + 行错误 | fixed | **可闭合（fixed）** | 行错误列表已落地。模板入口在用户页而非导入模态内，见本条 F-002 recommended |
| F-003 recommended · F04 二次确认 | fixed | **可闭合（fixed）** | 调账表单已有 confirm |
| F-004 recommended · F05 字段防抖 + 中文人话 | fixed | **不能按 fixed 接受** | 仅给独立控件加了 400ms 防抖；仍不读取任务表单 `cron` 字段；`describeCron` 仍是英文 stub |
| F-005 recommended · 台账残留 | fixed | **可闭合（fixed）** | 与 A-004 recommended 文档残留一并收口 |

## Findings

### F-001 · 预览依赖异步 `window.open`；复制链接是短命 `blob:` URL

| 字段 | 值 |
|------|-----|
| level | recommended |
| severity | med |
| status | fixed |
| closure | A-010 / GOAL-029 A-001（W18） |
| evidence | `apps/web/src/renderer/render.tsx`：`await fetcher` 之后 `window.open(objectUrl)`；`library.copyLink` 写入 `blob:`；两端均未 `revokeObjectURL`。 |

A-005 的 required 鉴权缺口已补。残余：异步 `window.open` 可能被弹窗拦截（本轮未浏览器点验）；`blob:` 不能作为可分享/可粘贴到外部的文件 URL；大文件 object URL 未释放。冻结方案中的 Lightbox 仍未做。不阻断 A-005 F-001 required 闭合。

### F-002 · 模板入口不在导入模态；行错误列表无定向测试

| 字段 | 值 |
|------|-----|
| level | recommended |
| severity | low |
| status | fixed |
| closure | A-010 / GOAL-029 A-001（W18） |
| evidence | `users.json` `import-template-block` 在页面 `body`，`openImport` 模态仍只有 `file` 字段；仓库内无 `data-import-error-rows` / 导入 200 `fieldErrors` 的 vitest。`import-template-download.tsx` 在 `!response.ok` 时静默返回。 |

用户已能在用户页下载模板，导入失败行可在模态内看到。与 D-001/GOAL-026 D-001「模态内链接」不完全一致，但不构成再开 required。

### F-003 · A-005 F-004 仍开放，A-006 标 fixed 过满

| 字段 | 值 |
|------|-----|
| level | recommended |
| severity | med |
| status | fixed |
| closure | A-009 / GOAL-028 A-001（W17） |
| evidence | `cron-preview.tsx` 独立输入 + 400ms 防抖；`scheduled-tasks.json` 仍是页面块，不绑创建/编辑表单 `cron`；`apps/api/internal/handler/scheduledtasks.go` `describeCron` 仍返回 `"every minute"` / `"every hour at minute N"` / `"cron schedule (5-field)"`。 |

A-006 把 A-005 F-004 整条标 fixed 不成立。防抖是增量，冻结方案的「字段下方即时中文语义」仍未交付。

## 必改项汇总

- **开放 required：0**
- A-005 F-001 / F-002 关闭证据可重复核对，允许按 `fixed` 闭合。
- 本条 F-001～F-003 均为 recommended，不阻断维持关门。

## 与既有意见的异同

| 意见 | 关系 |
|------|------|
| A-005 independent fail | 两条 required 的代码缺口本轮已补；不同意继续用 A-005 阻断关门。 |
| A-006 self pass | **部分同意**。required 闭合成立；F-003/F-005 recommended 闭合成立；**不同意** A-005 F-004 已 fixed。 |
| A-004 independent pass | 冲突已由用户采纳 A-005 后修正；不再作为完成证据。 |
| GOAL-026 A-002 | 与父目标 A-006 同步声明 F02/F03 fixed，与本轮 required 判定一致。 |

## 结论 + 建议给编排器/用户的下一步

A-005 的两条 required 可以闭合。GOAL-024 维持 `done` 在必改门禁上不再被本轮意见阻断。A-006 对 A-005 F-004 的 `fixed` 应改回 open，或由用户按 residual 接受「独立预览控件 + 英文短描述」。

建议 `/govern`：把 A-005 F-001/F-002 记为 `fixed`（引用本条）；A-005 F-004 / 本条 F-003 保持 recommended open，或书面 residual。

## 声明

本意见为独立交叉审计（`source: independent`），不修改目标 `status` / `progress` / goal-tree。响应由 `/govern` 处理。
