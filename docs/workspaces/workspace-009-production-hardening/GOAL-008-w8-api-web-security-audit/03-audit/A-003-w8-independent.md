---
id: A-003-w8-independent
doc: audit-entry
goal: GOAL-008-w8-api-web-security-audit
title: W8 F-001/F-002 close-out independent 复核
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
date: 2026-08-20
scope: A-001 F-001/F-002 required 修复实施与回归（close-out / implementation review）
verdict: pass
status: recorded
parent: GOAL-001-production-hardening
created: 2026-08-20
updated: 2026-08-20
version: 0.1.0
---

# A-003 · W8 F-001/F-002 close-out independent 复核

## 条目头

| 字段 | 值 |
|------|-----|
| **source** | independent |
| **auditor** | grok-build (grok-4.6 · reasoning high) |
| **类型** | close-out / code review |
| **scope** | A-001 F-001/F-002 required 修复实施与回归（close-out / implementation review） |
| **verdict** | **pass** |
| **工作区** | `workspace-009-production-hardening`（Root `GOAL-001-production-hardening`；canonical `docs/workspaces/workspace-009-production-hardening/`；`plan_refs`/`primary_plan` = `VP-009-production-hardening`；`shared_materials_catalog: none`） |

## 范围与区间

- **覆盖**：A-001 两条 required（F-001 分页整数溢出/切片 panic-DoS；F-002 生产 CSP 阻止 inline 首屏主题脚本）在现行代码中是否 genuine fixed；回归是否可重复核对。
- **方法**：工作区绑定核对 → 五件套通读 → 源码逐点抽验（`pagination.Bounds/Offset` 及全部调用点；`index.html` / `public/theme-init.js` / 生产 `dist` / nginx CSP）→ 本会话重跑 `cd apps/api && go test ./...`、`cd apps/web && npm test`、`npm run build`。未做动态 exploit / 真实浏览器 CSP 代理。
- **不覆盖**：不改 `status` / `progress` / 检查点 / 方案正文 / goal-tree。不把 A-001 F-003/F-004 升格为本波 required。不自行恢复 VP-008 `go` 宣称，不把本意见当作已关门。
- **排除**：F-003 localStorage refresh token（I-003 non-blocking）；F-004 development fallback 条件风险。

## P-005 / 工作区核对

| 核对项 | 结论 |
|--------|------|
| 工作区绑定 | `workspace.md`：`id=workspace-009-production-hardening`；Root `GOAL-001-production-hardening`；canonical 与 `goal-tree.md` 一致；`vision_role: delivery`；`plan_refs` + `primary_plan` = `VP-009-production-hardening`。Charter `schema-ui-core-admin-foundation@0.2.0`；VP-009 `vision_ref` 精确匹配。共享资料目录 `none`，本 scope 未把资料引用当关闭证据。 |
| I-001（finding 清单） | **verified**；本条不重开。A-001 已给出 F-001/F-002 required 与 F-003/F-004 非阻断处置。 |
| I-002（required 是否暂挂 VP-008 `go`） | **verified / closed by D-002**（整单采纳 F-001/F-002；闭合前暂挂 go 宣称）。本条确认代码闭合条件已满足；**恢复对外 go 宣称仍须 `/govern` 书面复核**，不得由本意见直接改宣称。 |
| I-003（localStorage refresh token 是否升级范围） | **open / non-blocking**；不阻断本 close-out scope。维持 A-001 F-003 原处置。 |
| 到期 required 信息项 | 无到期未关闭项阻断本 scope。I-002 绑的是宣称门禁，已由 D-002 合法关闭（go 暂挂）；不是本波代码闭合缺口。 |

## 成果（有证据）

### F-001 · 分页溢出/切片 panic → **fixed**

| 核对项 | 路径 | 结论 |
|--------|------|------|
| 共享安全边界 | `apps/api/internal/pagination/pagination.go` | `Bounds` 在 `page > lastPage` 时返回 `(total,total)`，不执行未保护的 `(page-1)*pageSize`。`lastPage = (total-1)/pageSize+1` 用 `int64` 除法。非正 `page`/`pageSize` 或 `total<=0` 返回 `(0,0)`。`Offset` 等于 `Bounds` 的 start。 |
| 原 exploit 用例已锁 | `apps/api/internal/pagination/pagination_test.go` | `huge page returns empty safely`：`page=92233720368547760, pageSize=100, total=25` → `(25,25)`，与 A-001 触发条件同形。另含 `MaxInt` page、巨大 pageSize、空 total、超越末页。 |
| A-001 点名的内存切片调用点 | `handler/datapermission.go:76`；`handler/account_self.go:342`；`handler/filelibrary.go:63` | 三处均改为 `pagination.Bounds(...)` 后再 `slice[start:end]`。`intParam`（`resources.go:899`）仍只拒绝 `<1`，**不再需要入口上限即可避免溢出切片**——超大 page 走空页语义。 |
| SQL `OFFSET` 调用点 | 见下表 14 处 | 全部 `pagination.Offset(page\|filter.Page, pageSize\|filter.PageSize, total)`；无遗漏裸 `(page-1)*pageSize` / `(filter.Page-1)*filter.PageSize`。 |
| 残余乘法检索 | `apps/api` 全仓 grep `(page - 1) * pageSize` 与 `(filter.Page - 1) * filter.PageSize` | 仅 `pagination.go` 自身在 `page <= lastPage` 保证下执行；调用点已清空。 |

**SQL OFFSET 抽验清单（本会话打开或与相邻 `pagination.Offset` 同行核对）：**

| 文件 | 用法 |
|------|------|
| `modules/authsession/users_repository.go:31` | `pagination.Offset(filter.Page, filter.PageSize, total)` |
| `modules/authsession/roles_repository.go:26` | 同上 |
| `modules/authsession/notifications_repository.go:163` | 同上 |
| `modules/authsession/service_credentials.go:84` | `pagination.Offset(page, pageSize, total)` |
| `modules/datadictionary/store/repository.go:111,273` | types / entries |
| `modules/wallet/store/repository.go:173,610,819` | accounts / ledger / reconciliation |
| `modules/operationlog/repository.go:230` | `ListOperationsFiltered` |
| `modules/recyclebin/store/repository.go:122` | recycle list |
| `modules/scheduledtasks/store/repository.go:113,267,318` | tasks / task runs / all runs |

`jobs/repository.go` 的 `LIMIT` 无 `OFFSET`，非分页溢出面。`mfa/totp.go` 的 `offset` 为 HMAC 动态截断，无关。

**原缺陷**：已认证用户提交极大正整数 `page`，`(page-1)*pageSize` 在 64 位溢出为负，随后切片或 SQL OFFSET 触发 panic/DoS。  
**现行为**：超大 page 映射为空页边界；内存切片与 SQL OFFSET 共用同一函数。A-001 触发条件不再成立。

### F-002 · CSP 阻止 inline 主题 bootstrap → **fixed**

| 核对项 | 路径 | 结论 |
|--------|------|------|
| 源 `index.html` 无 inline FOUC | `apps/web/index.html:14` | 仅 `<script src="/theme-init.js"></script>`；正文不再含 `localStorage.getItem("theme")`。第二个 script 为 Vite 入口（dev `/src/main.tsx`，prod hashed bundle）。 |
| 外部脚本内容 | `apps/web/public/theme-init.js` | 原 FOUC 逻辑（theme / prefers-color-scheme / `dark` class / `color-scheme`），同源静态文件。 |
| 生产 CSP | `apps/web/nginx.conf:29` | `script-src 'self'`，无 nonce/hash；外部 `/theme-init.js` 可由 `'self'` 放行。本修复未改 CSP。 |
| 测试锁 | `apps/web/src/theme/theme-init.test.ts` | 断言引用外部脚本、源 `index.html` 不含 inline `localStorage.getItem("theme")`、`public/theme-init.js` 存在且含原逻辑。本会话 3/3 通过。 |
| 生产构建产物 | `apps/web/dist/index.html`、`apps/web/dist/theme-init.js` | `npm run build` 后 `dist/index.html` 仍为 `<script src="/theme-init.js"></script>` + hashed module；`public/theme-init.js` 已复制到 `dist/theme-init.js`。无 inline bootstrap。 |

**原缺陷**：inline FOUC 脚本被生产 `script-src 'self'` 拦截，首屏主题可能错误或闪烁。  
**现行为**：bootstrap 为同源外部脚本；CSP 无需 hash/nonce。代码级不一致已消除。

### 回归（本会话执行，非转述 A-002）

| 命令 | 结果 |
|------|------|
| `cd apps/api && go test ./...` | **exit 0**；含 `internal/pagination` |
| `cd apps/web && npm test` | **73 files / 1072 tests passed**（含 `theme-init.test.ts` 3 tests） |
| `cd apps/web && npm run build` | **exit 0**（tsc -b + vite build；chunk size 警告与本 scope 无关） |

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| S1 独立意见落盘 | 已达成（前序） | A-001 |
| S2 用户确认 required 范围与 go 影响 | 已达成（前序） | D-002：整单采纳 F-001/F-002；闭合前暂挂 VP-008 go |
| S3 按确认范围实施并回归 | 已达成 | E-002 + 本条源码抽验 + 本会话 API/Web/build 全绿 |
| S4 self/independent 复核确认 required 合法闭合 | **本条判定代码闭合条件已满足** | A-002 self pass + 本条 independent pass。本意见不改 status；关门与 go 恢复由 `/govern` 响应后执行 |

## Findings

无开放 required。A-001 F-001 / F-002 可核对为 **fixed**。

| finding（A-001） | 本条判定 | 证据 |
|------------------|----------|------|
| F-001 | **fixed** | `pagination.Bounds/Offset` + 3 处内存切片 + 14 处 SQL OFFSET + 极大 page 单测 + `go test ./...` |
| F-002 | **fixed** | `index.html` 外部 `/theme-init.js`；`public/` 与 `dist/` 均有该文件；`script-src 'self'` 可放行；`theme-init.test.ts` + `npm test` / `npm run build` |

### 本条 recommended（不阻断闭合）

| ID | 严重度 | 建议 | 状态 | 说明 |
|----|--------|------|------|------|
| F-001 | low | recommended | open | `apps/web/src/main.tsx:37` 仍留「synchronous **inline** script」旧注释，下一行已改为 external。无运行时影响；建议 `/govern` 顺手删旧行。 |
| F-002 | low | recommended | open | 闭合证据是源码 + 构建产物 + 文件锁测试，**不是**带真实 CSP 响应头的浏览器回归。A-001 曾建议浏览器层检查；代码级 F-002 已成立，该项不升格 required。 |

A-001 F-003 / F-004 维持原 recommended / conditional，本条不重开、不升格。

## 必改项汇总

- **开放 required = 0**
- 无未合法闭合的 required finding。
- 无到期且影响本 scope 的 required 信息项。

## 与 A-001 / A-002 的异同

- **与 A-001**：同向认定 F-001/F-002 为当时 required。本条独立复核现行代码后判定两条 **genuine fixed**，不沿用 A-002 结论作为唯一依据。F-003/F-004 维持 A-001 非阻断处置。
- **与 A-002**：同向 pass、同向 F-001/F-002 fixed。差异：本条 `source: independent`；本会话重跑全量 API/Web/build 并打开全部 OFFSET/切片调用点与 `dist/` 产物；A-002 未记录 `dist/index.html` / `dist/theme-init.js` 构建抽验。本条另记 2 条 low recommended（旧注释、无浏览器 CSP 测试），不构成与 A-002 的冲突。
- **不冲突**：无需 P-004 裁决。

## 结论 + 建议给 `/govern` 的下一步

**verdict: pass。** A-001 F-001、F-002 关闭证据充分、可重复核对。

建议编排器：

1. 响应本条：将 A-001 F-001/F-002 记为 `fixed`（P-003 合法闭合）。
2. 勾选 S4；按检查点重算派生 progress；将 GOAL-008 按用户授权推进关门（本意见不改 status）。
3. 按 D-002「闭合后恢复宣称前应复核」写 **D-003**（或等价决策）恢复 VP-008 `go` 宣称——**本意见不恢复 go**。
4. 同步工作区指针：`workspace.md` 波次表与 `goal-tree.md` 仍写 GOAL-008 `draft` /「待用户裁决修复范围」，与现行 `00-meta` `active` + D-002/E-002 滞后；由 `/govern` 在响应时更正。
5. I-003 保持 open non-blocking；本波 recommended F-001/F-002 可顺手清理或登记后续维护，不阻断关门。

## 声明

本意见不修改 status/progress；响应由 /govern 处理。
