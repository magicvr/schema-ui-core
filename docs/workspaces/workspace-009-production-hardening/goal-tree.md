---
title: 目标树 · workspace-009-production-hardening
status: active
created: 2026-08-10
updated: 2026-08-17
parent: null
version: 0.10.1
workspace_id: workspace-009-production-hardening
---

# 目标树 · 生产加固（共享基架持续安全与健壮性）

> 工作区：`workspace-009-production-hardening`
> canonical：`docs/workspaces/workspace-009-production-hardening/`
> Root：`GOAL-001-production-hardening`，**长期程序容器 · active**
> primary_plan：`VP-009-production-hardening`，**active**
## 树
```text
GOAL-001-production-hardening [active]  · 持续安全程序
├── GOAL-002-audit-findings-remediation [done] (16/16)    · W1
├── GOAL-003-upload-ownership-hardening [done] (4/4)      · W2
├── GOAL-004-w3-security-audit-remediation [done] (8/8)   · W3
├── GOAL-005-w4-security-audit-remediation [done] (8/8)   · W4
└── GOAL-006-w6-scan-findings-remediation [done] (4/4)    · W6
```

**W6（2026-08-15）**：承接本会话对 api/web 的代码审视——scheduler 未到期任务 5 年分钟空扫描改 O(1) Matches 快速路径（每日一次诊断保留）；回收站还原孤儿字典项 500 退化改 409 DICT_KEY_NOT_FOUND（快照保留可重试）；branding data:image 内联评估后**不采纳**（API normalizeLogoURL 与 errorcatalog 均拒绝，web 测试锁定，保持一致收紧）。`go test ./...` 全绿，self 审计 A-001 pass，开放 required = 0；**2026-08-17 补记用户授权关门（D-002）+ close-out self 审计 A-002 pass**，`status: done` 维持。
**W5（2026-08-14 扫描）**：全量审计 **0 中高危**（L-001～L-006 低危就地修补，见 Root [E-002](GOAL-001-production-hardening/02-execution/E-002-w5-scan-zero-midhigh.md)）；按程序约定未开子目标。**go 判定：无影响、不暂挂**（安全加固与已冻结 fail-closed 语义一致；未改 Profile 默认集 / 模块矩阵 / Manifest 装配 / 协议 pin）。
Root **保持 active**。W1–W4 为已关门波次档案；W4 承接 2026-08-11 新一批 api/web 全量审计修复（限流驱逐、上传权限门+配额、改密吊销 access token、前端异常捕获、URL 校验、启动加固、文案）。

## 状态表

| id | title | parent | status | progress | updated |
|----|-------|--------|--------|----------|---------|
| GOAL-001-production-hardening | 生产加固（共享基架持续安全与健壮性） | null | active | —（程序容器，不用 n/n→done） | 2026-08-10 |
| GOAL-002-audit-findings-remediation | 审查发现修复（W1） | GOAL-001-production-hardening | done | 16/16 | 2026-08-10 |
| GOAL-003-upload-ownership-hardening | 上传所有权与下载鉴权加固（W2） | GOAL-001-production-hardening | done | 4/4 | 2026-08-10 |
| GOAL-004-w3-security-audit-remediation | W3 安全审计发现修复（api/web） | GOAL-001-production-hardening | done | 8/8 | 2026-08-11 |
| GOAL-005-w4-security-audit-remediation | W4 安全审计发现修复（api/web） | GOAL-001-production-hardening | done | 8/8 | 2026-08-11 |
| GOAL-006-w6-scan-findings-remediation | W6 扫描审计发现修复（api/web） | GOAL-001-production-hardening | done | 4/4 | 2026-08-15 |
| — | W5 scan（0 中高危；低危就地修补，未开子目标） | GOAL-001-production-hardening | — | — | 2026-08-14 |

## 维护说明

- Root 是长期能力容器；`status: done` 仅在程序废弃或 `primary_plan` 迁移且用户确认时使用。
- 波次 progress 只写在子目标；不得用波次完成数推导 Root done。
- 层级唯一来源是目标 `00-meta.md` 的 `parent`。