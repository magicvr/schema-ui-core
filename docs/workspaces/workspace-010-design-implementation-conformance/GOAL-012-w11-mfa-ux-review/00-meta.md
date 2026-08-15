---
id: GOAL-012-w11-mfa-ux-review
title: W11 · 个人中心 MFA 缺陷修复与全局 UX 审视整改
status: done
parent: GOAL-001-design-implementation-conformance
created: 2026-08-15
updated: 2026-08-15
version: 0.4.0
progress: 5/5
---

# GOAL-012 · W11 · 个人中心 MFA 缺陷修复与全局 UX 审视整改

VP-010 / workspace-010 的**第十一波**（用户 2026-08-15 立项）：落盘 2026-08-15 以真实用户视角对已实现 api/web 功能页的**全局 UX 审视改进项（U-01～U-14）**，并纳入用户补充报告的 **3 项个人中心 MFA 缺陷（M-01～M-03）**。问题全量清单与初步代码定位见 [01-decision.md](01-decision.md)。

## 当前边界

- 范围：个人中心（/account）MFA 绑定/解绑/确认缺陷修复（M-01～M-03）；全局 UX 审视改进（U-01～U-14，按优先级分批实施，S1 裁决分批范围）。
- 非范围：不改认证令牌协议本体（除非 I-001 裁决需要）；不做 VP-005 全量视觉基线重做；业务领域功能（订单/钱包/类目等）不在本波。

## 成功标准与路线图（P-001）

- [x] **S1 · 范围与优先级确认**：D-001 落盘（分批顺序 + I-001 修复方向 + I-003 二维码方案）；M-01～M-03 根因经代码级验证（E-001→E-002）
- [x] **S2 · MFA 缺陷修复**：M-01 二维码；M-02 停用不误登出/成功提示；M-03 绑定输错码报错重填（E-002，回归全绿）
- [x] **S3 · UX P0 实施**：U-01 用户角色多选 + U-02 角色权限动态化（E-003，optionsSource 本地扩展 + RBAC 目录端点）
- [x] **S4 · UX P1 实施**：U-03 Toast、U-04 表格搜索/筛选、U-05 行操作收纳、U-06 分页增强、U-07 空状态（E-004）
- [x] **S5 · 验证与关门**：Go 全量 + Web 全量回归绿（GO_ALL_OK / 1002/1002 / tsc 0）；审计 A-001 self pass + A-002 independent（grok）conditional→resolved + A-003 closeout self pass；goal-tree / workspace 同步；**2026-08-15 关门（5/5）**

progress: 由五个等权检查点派生（S1～S5）；当前 **5/5**（2026-08-15 关门）。

## 审计策略

本目标含**认证/MFA 语义改动**（401 映射、会话吊销用户提示）→ 按 P-003 默认 **self**；S2 涉及认证语义，关门审计建议 **cross**（independent provider 惯例 grok，执行时确认）。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 |
|----|------|-----------------|----------|--------------|-----------------|------|
| I-001 | required | M-02/M-03 根因确认与修复方案：无效动态码的 HTTP 语义（401 vs 400）与 authFetch 处理策略；解绑成功后会话吊销的用户提示方式 | 方案 | S2 | 复现验证 + 回归测试（E-002）；D-001/D-002 裁决 | **closed**（2026-08-15 实施并验证） |
| I-002 | required | UX P0 交互方案：角色多选数据源（/api/roles 动态选项）与协议表达；权限动态化（权限元数据接口形态与 schema 分组） | 方案 | S3 | D-002 裁决：optionsSource 本地扩展 + /api/permissions、/api/menu-items 目录端点；分组矩阵留 P2 | **closed**（2026-08-15） |
| I-003 | non-blocking | 二维码渲染实现方式：新增依赖 vs 自绘 canvas（离线可用、体积） | 方案 | S2 | D-001 裁决：qrcode-generator（MIT 零依赖，SVG 渲染） | **closed**（2026-08-15） |
| I-004 | required | UX P1 范围确认：Toast 方案；表格搜索/筛选是否需要扩展协议能力（schema 节点能力） | 方案 | S4 | D-003 裁决（Toast 本地 UI；搜索复用 search-form 模式；select 筛选留 P2） | **closed**（2026-08-15） |

## 父目标

- [GOAL-001-design-implementation-conformance](../GOAL-001-design-implementation-conformance/00-meta.md)

## 台账布局

- `01-decision/`：D-NNN 决策（S1 起）；`02-execution/`：E-NNN 事实记录；`03-audit/`：A-NNN 审计意见。
- 跨区引用（workspace-011 的 GOAL-016/017/018 等）用 Q2 路径，见 workspace-protocol §2.6。
