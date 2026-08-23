---
id: D-003-w10-scope-reconciliation
goal: GOAL-010-w10-api-web-security-audit
status: accepted
created: 2026-08-21
updated: 2026-08-21
parent: GOAL-010-w10-api-web-security-audit
version: 0.1.0
---

# D-003 · W10 修复范围调和：4 条 recommended 经源码核实作废（2026-08-21）

## 决定（编排器调和 · 依代码证据；用户可书面驳回）

S3 实施前逐条核实发现，D-002 采纳的 7 条中 **4 条 recommended 不成立**（审计误报），实际可实施范围为 **3 条**：

| F-ID | D-002 采纳时 | 核实结论 | 证据 |
|------|--------------|----------|------|
| F-001 | required HIGH | **fixed** | env.example + config_test.go 凭据全部替换为占位符/假值 |
| F-002 | recommended | **fixed** | 新增 `lib/fetch-timeout.ts` + auth-client/load-page/form-controls 默认路径接线 |
| F-003 | recommended | **作废（不成立）** | render.tsx:350-384 预览窗口仅写入静态模板（无任何不可信内容插值；iframe `sandbox=""` 禁脚本）；且 `noopener` 特性会使 `window.open` 返回 null，功能本身依赖该引用写入内容 |
| F-004 | recommended | **作废（误报）** | accounts.go:337-359 `RevokeRefreshToken` 已是防护式 UPDATE（`WHERE id=? AND revoked_at IS NULL` + RowsAffected 检查 + ErrAlreadyRevoked）；审计只读了 auth.go 调用层未核实仓库层 |
| F-005 | recommended | **作废（误报）** | render.tsx:418 `/^[._-]+|[._-]+$/g` 已剥离前导点；注释明确 "leading/trailing separators are trimmed" |
| F-006 | recommended | **作废（误报）** | service_credentials.go:151-155 `normalizedCredentialScopes` 去重后强制 1..64 上限 |
| F-007 | recommended | **fixed（升级为真实缺口）** | form-controls.tsx:78 字符类缺 `\\` 排除——WHATWG URL 解析将特殊 scheme 中 `\` 规范化为 `/`，`/\host` → 协议相对 `//host` 逃逸同源；11 处同型正则中仅此一处缺失 |

## 调和依据

- P-003 审计意见响应义务：发现意见与事实不符时如实记录，不作假修复。
- 对齐 W9 先例（其 D-002 调和表将原文 F-003 作废、补号 F-025）。
- **go 宣称暂挂不变**（D-002 §2）：F-001 为 HIGH required 且已修复，但 S4 复核未完成前继续暂挂。

## 用户裁决点

本调和改变 D-002 采纳清单的构成（7→3 实施 + 4 作废）。若用户驳回任一作废结论，按 P-004 以书面指令恢复该条并重开实施。