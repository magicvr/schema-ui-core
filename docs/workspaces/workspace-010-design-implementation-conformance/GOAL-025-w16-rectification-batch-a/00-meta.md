---
id: GOAL-025-w16-rectification-batch-a
title: W16 整改批 A · 安全与认证基线（F01 首次改密 / F07 一键下线其他 / F08 验证码刷新与 MFA 备份）
status: active
parent: GOAL-024-w16-user-perspective-improvements
created: 2026-08-17
updated: 2026-08-17
version: 0.1.0
progress: 2/4
---

# GOAL-025 · W16 整改批 A（安全与认证基线）

[GOAL-024](../GOAL-024-w16-user-perspective-improvements/00-meta.md)（W16）的**下级整改子目标（批 A）**：承接 D-003 分批规划，本子目标实施 **W16-F01 首次登录强制修改初始密码**、**W16-F07 个人中心一键下线其他设备**、**W16-F08 登录验证码主动刷新 + MFA 密钥复制与恢复码下载**。

## 当前边界

- **范围（本子目标实施）**：
  - W16-F01：`users.must_change_password` 字段、登录响应标记、强制改密拦截、改密后清标记。
  - W16-F07：`POST /api/account/sessions/revoke-others` 端点 + 个人中心“下线其他设备”按钮。
  - W16-F08：登录页验证码“换一题”、MFA 密钥复制与恢复码 txt 下载。
- **非范围**：批 B（F02/F03/F04）、批 C（F05/F06/F09/F10）不在此实施；不改 Profile 默认集 / 模块矩阵 / Manifest 装配语义；安全相关变更需独立审计与回归。

## 成功标准与路线图（P-001）

- [x] **S1 · 方案冻结**：F01/F07/F08 设计（端点/schema/存储契约/前端交互）+ 信息项登记（D-001）。
- [x] **S2 · 实施**：F01/F07/F08 代码与 schema/前端接线。
- [ ] **S3 · 测试与回归**：Go 全量 + Web vitest/tsc + 相关 e2e；涉登录/会话安全门禁时执行 independent 审计。
- [ ] **S4 · 自审与关门**：审计 + 台账同步 + goal-tree/workspace 同步。

progress: 由四个等权检查点派生（S1～S4）；当前 **2/4**。

## 审计策略

| 阶段 / 项 | 默认模式 | 说明 |
|-----------|----------|------|
| S1 冻结 | self | 方案落盘 + 证据核对 |
| F01/F07/F08 实施 | independent | 登录/会话安全门禁，按项目独立审计执行路径 |
| S4 关门 | self | 常规关门自审（需先闭合 independent findings） |

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|
| I-001 | required | F01 强制改密门禁的端点白名单（哪些端点允许 must_change_password=1 用户访问） | S1 方案冻结 | S1 | 梳理登录/改密/资料端点与中间件 | **closed** | D-001 §2.2：登录/刷新/登出/改密/资料/me + 验证码/MFA 必要端点，其余 403 |
| I-002 | non-blocking | F07 revoke-others 是否 bump token_version 导致当前 access token 也失效 | S1 方案冻结 | S1 | 会话/令牌模型核对 | **closed** | D-001 §3.1：bump + 吊销全部旧 refresh + 为当前设备重签新令牌，当前设备不中断 |

## 父目标

- [GOAL-024-w16-user-perspective-improvements](../GOAL-024-w16-user-perspective-improvements/00-meta.md)

## 台账布局

- `01-decision/`：D-NNN；`02-execution/`：E-NNN；`03-audit/`：A-NNN。
- 跨区引用用 Q2 路径（workspace-protocol §2.6）。
