---
id: E-002-w8-implementation
doc: execution-entry
goal: GOAL-008-w8-api-web-security-audit
date: 2026-08-20
status: recorded
parent: GOAL-001-production-hardening
created: 2026-08-20
updated: 2026-08-20
version: 0.1.0
---

# E-002 · W8 F-001/F-002 required 修复实施与回归

## 事实（已发生）

- 用户目标轮次指令授权推进 GOAL-008 至闭门；D-002 确认整单采纳 A-001 F-001/F-002 并暂挂 VP-008 go 宣称。
- F-001 分页整数溢出/切片 panic 修复：新增 `apps/api/internal/pagination` 共享安全分页工具（`Bounds`/`Offset`），所有内存切片分页与 SQL `OFFSET` 计算改用该工具，避免 `(page-1)*pageSize` 溢出。
- F-002 CSP inline 首屏主题脚本：把 `index.html` 内联 FOUC bootstrap 迁移到 `apps/web/public/theme-init.js`，`index.html` 改为 `<script src="/theme-init.js"></script>`；生产 CSP `script-src 'self'` 即可允许，无需 inline hash/nonce。

## 文件级改动

| 文件 | 改动 |
|------|------|
| `apps/api/internal/pagination/pagination.go` | 新增：溢出安全的 `Bounds(page,pageSize,total)` 与 `Offset(...)` |
| `apps/api/internal/pagination/pagination_test.go` | 新增：极大 page/pageSize、空 total、超越末页、边界用例 |
| `apps/api/internal/handler/resources.go`（intParam 消费端经 `datapermission.go`/`account_self.go`/`filelibrary.go` 处理） | 通用分页入口未引入业务上限；改用共享 Bounds 安全切片 |
| `apps/api/internal/handler/datapermission.go` | 内存分页改 `pagination.Bounds` |
| `apps/api/internal/handler/account_self.go` | 内存分页改 `pagination.Bounds` |
| `apps/api/internal/handler/filelibrary.go` | 内存分页改 `pagination.Bounds` |
| `apps/api/internal/modules/{operationlog,authsession,datadictionary,wallet,recyclebin,scheduledtasks}/*/repository.go` | SQL `OFFSET` 改 `pagination.Offset(page,pageSize,total)` |
| `apps/web/index.html` | inline bootstrap 移除，改外部 `/theme-init.js` |
| `apps/web/public/theme-init.js` | 新增：外部 FOUC bootstrap（原 inline 逻辑） |
| `apps/web/src/theme/theme-init.test.ts` | 新增：锁定 index.html 无 inline script、引用外部脚本、public 文件存在 |
| `apps/web/src/theme/theme.ts`、`src/main.tsx` | 注释更新为外部脚本事实 |

## 回归证据

- `cd apps/api && go test ./...`：**全部通过**（exit 0，含新增 `pagination` 包测试）。
- `cd apps/web && npm test`：**1072/1072 全绿**（73 个测试文件；较上一轮 1069 增加 3 个 theme-init 断言）。
- `cd apps/web && npm run build`：**通过**（tsc -b + vite build exit 0；dist 产物生成，仅有 chunk size 警告）。
- 未发现生产代码新增 inline script；CSP 配置 `script-src 'self'` 无需修改。

## 备注

- F-003/F-004 按 A-001 原判定保持 non-required/conditional，不在本波 required 闭合范围。
- VP-008 go 宣称维持暂挂，待 independent 复核后按 D-003 恢复流程处理。