---
id: workspace-009-production-hardening
title: 生产加固工作区（共享基架持续安全）
status: active
root_goal: GOAL-001-production-hardening
canonical_scope: docs/workspaces/workspace-009-production-hardening/
shared_materials_catalog: none
vision_role: delivery
plan_refs:
  - VP-009-production-hardening
primary_plan: VP-009-production-hardening
created: 2026-08-10
updated: 2026-08-26
version: 0.16.0
parent: null
---

# 工作区上下文 · 生产加固（持续安全）

本工作区是 [VP-009-production-hardening](../../vision/plans/VP-009-production-hardening.md)（`active` · **长期安全与健壮性程序**）的唯一 lead delivery workspace。

- **Root** 为长期程序容器（默认 `active`）。  
- **子目标** 为有界扫描/修复波次（可 `done`）。  
- 不因单波完成而关闭本区或 VP；不改变 Charter `primary_workspace`。

## 绑定

| 字段 | 当前值 | 说明 |
|------|--------|------|
| 工作区 ID | `workspace-009-production-hardening` | 与本区目标及资料引用的 `workspace_id` 一致 |
| Root Goal | `GOAL-001-production-hardening` | `parent: null`；长期容器 |
| canonical 范围 | `docs/workspaces/workspace-009-production-hardening/` | 本区唯一目标状态范围 |
| 共享资料目录 | `none` | 暂无固定共享资料 |
| 愿景角色 | `delivery` | VP-009 lead；不改变 Charter primary workspace |
| 规划对齐 | `primary_plan` = `VP-009-production-hardening` | 持续程序意图 |

## 愿景对齐

Charter：`schema-ui-core-admin-foundation@0.2.0`。  
VP-009 为共享基架持续安全程序；与 VP-008 `go` 消费有效性接口见该 VP。  
independent provider（沿用 workspace-008 D-002）：**grok build · grok-4.6 · high · `audit`**；波次 security 高影响默认 `cross`。

## 波次（实现层指针）

| 波次 | 子目标 | status |
|------|--------|--------|
| W1 | GOAL-002-audit-findings-remediation | done |
| W2 | GOAL-003-upload-ownership-hardening | done |
| W3 | GOAL-004-w3-security-audit-remediation | done |
| W4 | GOAL-005-w4-security-audit-remediation | done |
| W5 scan | （未开子目标） | 2026-08-14 全量审计 **0 中高危**；低危就地修补见 Root E-002；**go 判定：无影响、不暂挂**（安全加固与冻结 fail-closed 语义一致；未改 Profile/模块矩阵/Manifest 装配/协议 pin） |
| W6 | GOAL-006-w6-scan-findings-remediation | **done**（4/4 · 2026-08-15 完成；2026-08-17 补记用户授权关门 D-002 + close-out self 审计 A-002 pass） |
| W7 | GOAL-007-w7-api-web-security-audit | **done**（4/4 · 2026-08-19 独立审计 A-001 fail → D-002 整单采纳 → E-002/E-003 实施 → self A-002 + independent A-004 pass → **A-005 独立代码复核 pass + VP-008 go 宣称恢复（D-003）**） |
| W8 | GOAL-008-w8-api-web-security-audit | **done**（4/4 · 2026-08-20：A-001 fail → D-002 整单采纳 + go 暂挂 → E-002 修复 → A-002 self pass + A-003 independent pass → D-003 恢复 VP-008 go 宣称；真实浏览器/CSP 回归 E-004 通过；发版前冒烟流程集成 E-005） |
| W9 | GOAL-009-w9-api-web-security-audit | **done**（4/4 · 2026-08-21：A-001 fail + A-002 conditional → D-002 清单 required 12（F-003 作废，F-025=P2-7）→ D-003 整单采纳 + go 暂挂 → E-004 S3 实施 12/12 修复 + 回归全绿 + A-004 self pass → **A-005 grok-build（grok-4.6 · high）independent pass：12/12 genuine fixed** → A-006 闭合记录 fixed ×12、开放 required = 0 → E-005 三条 recommended 全部实施并锁定 → D-004 恢复 VP-008 go 宣称） |
| W10 | GOAL-010-w10-api-web-security-audit | **done**（4/4 · 2026-08-21：A-001 independent conditional（1 HIGH required）→ D-002 整单 7 条 + go 暂挂 → D-003 调和 4 误报作废 → E-002 修复 3 条 + 回归全绿 + A-002 self pass → **A-003 grok-build（grok-4.6 · high）independent pass：3/3 genuine fixed** → 用户书面闭合授权 → E-003 A-003 recommended ×3 全部修正 + 索引同步 → A-004 闭合记录 fixed ×6 + 作废 ×4、开放 required = 0 → **D-004 关门 + 恢复 VP-008 go 宣称**；残余移交：数据库密码轮换） |
| W11 | GOAL-011-w11-api-web-security-audit | **done**（4/4 · 2026-08-22：A-001 independent fail（6 required）→ D-002 整单采纳 + go 暂挂 → E-002 S3 实施 6/6 + E-003 recommended 处置（fixed 11 + overruled 2 有据）→ A-002 self pass → **A-003 grok-build（grok-4.6 · high）independent pass：6/6 genuine fixed + 真实 PG 复跑** → A-004 闭合记录（开放 required = 0；I-003 关闭）→ **D-004 关门 + 恢复 VP-008 go 宣称**；残余移交：密码轮换 + R-001/R-002 + F-009 lastRun） |
| W12 | GOAL-012-w12-multi-instance-rate-limiting | **done**（4/4 · 2026-08-26 评估型收官：承接跨区登记项 [workspace-019 E-009 §F-002](../../workspace-019-iam-recovery/GOAL-001-iam-recovery/02-execution/E-009-a001-finding-fixes.md)；D-002 三项用户裁决 = 维持单实例官方边界 / 载体预登记 Redis 方向 / 零码变更；self A-001 `pass` 关门；复审触发 =「多实例部署形态出现」） |
| W13 | GOAL-013-w13-api-web-security-audit（子目标 [GOAL-014](GOAL-014-w13-account-lockout-redesign/00-meta.md) 承载 F-007） | **done**（6/6 · 2026-08-26：A-001 conditional（required F-001～F-004 + P3/B 全量分母，D-001 范围=一次修完）→ S2～S5 全量实施（checkpoints `9da0084e`/`b7954235`/`e93f7228`；D-002 三裁决：F-007=fixed→GOAL-014、F-013=accepted-residual+Root E-008 硬门、F-020=HSTS+保留 img-src）→ A-002 self pass → **A-003 grok-build（grok-4.6 · high）independent pass：required ×4 genuine fixed** → A-004 响应 ×3 → GOAL-014 分层锁定模型落地（迁移 0061、来源锁 5/15min、全局天花板 100/24h 滑动、失败零吊销；真实 PG 复核；A-001/A-002 双 pass）→ **D-003/D-004 用户书面关门：两目标一并 done**；残余移交：I-001 TLS 拓扑 deferred、F-013 复审硬门（Root E-008）、GOAL-014 Refresh 残余（R-F002）） |
| W14 | [GOAL-015-w14-schema-auth-wiring-lock](GOAL-015-w14-schema-auth-wiring-lock/00-meta.md) | **done**（4/4 · 2026-08-26：用户报障全页「无法显示此页面」→ 定位 F-010 后生产入口缺 `schemaFetcher` 认证传输（匿名 401 · 测试装配≠生产装配）；hotfix 用户确认落地 + AuthGate 模块化 + 生产装配回归锁 ×2（vitest 1130/1130 · tsc 0 · 变异验证红→绿）；R-001 并入 = e2e Bearer 冒烟 fixed；附带对齐 shell.spec 匿名 schema 探测陈旧契约至 F-010；Playwright e2e 10 passed / exit 0；self A-001 pass → **D-002 用户书面关门 done**） |

## 固定共享资料引用

> `shared-materials/index.json` 只能提供候选路径与摘要。缺完整引用字段的行无效。

| reference_id | workspace_id | material_id | source | version | sha256 | purpose | local_record | status |
|--------------|--------------|-------------|--------|---------|--------|---------|--------------|--------|
| — | — | — | — | — | — | — | — | — |
