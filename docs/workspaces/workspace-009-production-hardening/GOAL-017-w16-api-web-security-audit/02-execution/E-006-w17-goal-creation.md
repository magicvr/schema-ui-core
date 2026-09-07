---
id: E-006-w17-goal-creation
title: W17 子目标创建（承接 F-003 残余）
date: 2026-09-01
status: recorded
---

# E-006 · W17 子目标创建（承接 F-003 残余）

## 背景

W16 (GOAL-017) 已完成并关门，但 F-003 (refresh token localStorage XSS 风险) 作为 accepted-residual 延期到后续波次。按治理规则（P-005 信息就绪 + 残余风险复审触发），需在合适时机承载修复工作。

## 执行事实

**日期**: 2026-09-01

**动作**: 创建 GOAL-018 (W17 · Refresh Token httpOnly Cookie 双模式架构)

**范围**:
1. **五件套文件**:
   - `00-meta.md`: 目标定义、成功标准（S1–S5 共 5 检查点）、信息需求（I-001～I-004）
   - `01-decision.md`: 信息就绪台账（4 项）、决策索引（待 S1 填写）
   - `02-execution.md`: 执行索引（待实施后填写）
   - `03-audit.md`: 审计索引（预期 S4 自审 + S5 独立审计）

2. **目标属性**:
   - `id`: GOAL-018-w17-refresh-token-httponly
   - `parent`: GOAL-001-production-hardening
   - `status`: draft（等待 I-001/I-002/I-003 required 信息就绪后进入 active）
   - `progress`: 0/5（5 个阶段检查点）

3. **技术方案（初步）**:
   - **优先模式**: httpOnly cookie（`Secure`, `SameSite=Strict`, `HttpOnly`, `Path=/api/auth`）
   - **回退模式**: 保留 `X-Refresh-Token` header 支持（非浏览器环境）
   - **API 改造**: 3 个端点（login/refresh/logout）
   - **Web 改造**: cookie 优先 + localStorage 回退检测

4. **信息需求**（阻断 S1 方案冻结）:
   - I-001 (required): Cookie 属性配置
   - I-002 (required): 非浏览器环境兼容性策略
   - I-003 (required): Token 轮换时的 cookie 更新策略
   - I-004 (non-blocking): 开发环境 cookie 配置

5. **目标树更新**:
   - 树结构：GOAL-018 添加到 GOAL-001 子节点
   - 状态表格：新增 GOAL-018 行
   - W17 叙事段落：关联 W16 残余移交

## 产物路径

- `docs/workspaces/workspace-009-production-hardening/GOAL-018-w17-refresh-token-httponly/00-meta.md`
- `docs/workspaces/workspace-009-production-hardening/GOAL-018-w17-refresh-token-httponly/01-decision.md`
- `docs/workspaces/workspace-009-production-hardening/GOAL-018-w17-refresh-token-httponly/02-execution.md`
- `docs/workspaces/workspace-009-production-hardening/GOAL-018-w17-refresh-token-httponly/03-audit.md`
- `docs/workspaces/workspace-009-production-hardening/goal-tree.md`（已同步）

## 关联

- **父目标**: GOAL-001-production-hardening (Root)
- **来源**: GOAL-017 (W16) F-003 accepted-residual
- **引用**: GOAL-017 A-001 F-003, D-002

## 后续步骤

1. **S1 方案冻结**（需用户确认）:
   - 解决 I-001/I-002/I-003 required 信息需求
   - 详细设计 cookie 属性配置
   - 明确向后兼容性策略
   - 完成测试计划

2. **S2–S4 实施**: API 端 → Web 端 → 集成验证

3. **S5 审计与关门**: Self + Independent（按 P-003 建议 `independent` 或 `cross`，安全改造）

## 验证

- [x] 五件套文件已创建并填充真实 id/title/parent
- [x] `goal-tree.md` 已同步（树 + 表格 + 叙事）
- [x] W16 叙事中已添加残余移交到 GOAL-018
- [x] 目标处于 draft 状态，信息门禁明确
