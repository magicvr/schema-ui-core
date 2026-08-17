---
id: GOAL-024-w16-user-perspective-improvements
title: W16 · 真实用户视角未计划改进项台账与规划（W16-F01～W16-F10）+ 整改承接
status: active
parent: GOAL-001-design-implementation-conformance
created: 2026-08-17
updated: 2026-08-17
version: 0.2.0
progress: 6/8
---

# GOAL-024 · W16 · 真实用户视角未计划改进项台账与规划

VP-010 / workspace-010 的**第十六波**（用户 2026-08-17 点名立项）：承接真实用户视角审视中**尚未在 workspace-011 规划**的 10 项核心体验改进项，建立问题台账、技术方案设计与分批实施规划，并通过下级整改子目标分批实施。

## 背景与边界

- **排除项（已在 workspace-011 规划）**：
  - `B-10 组织 / 部门 / 岗位`（部门树与组织架构层级）
  - `S-05 公告管理` / `B-09 消息模板`（管理员系统广播/公告下发）
  - `B-09 邮件/短信配置`（SMTP/SMS 基础配置）
  - `B-11 登录日志独立视图`
- **纳入项（W16-F01～W16-F10）**：
  - W16-F01：首次登录强制修改初始密码（安全合规）
  - W16-F02：文件库在线直接预览（图片/PDF）与一键复制链接
  - W16-F03：数据导入弹窗支持下载 CSV 模板与逐行错误定位反馈
  - W16-F04：钱包余额金额“分转元”标准格式化展示与调账警示
  - W16-F05：定时任务 Cron 表达式自然语言解析与未来运行时间预览
  - W16-F06：系统监控页面支持定时自动轮询刷新（5s/10s/30s）
  - W16-F07：个人中心活跃会话支持“一键下线其他所有设备”
  - W16-F08：登录算术验证码主动刷新 + MFA 密钥复制与恢复码 txt 下载
  - W16-F09：数据字典项支持 Badge 颜色/标签风格配置
  - W16-F10：系统设置支持页脚版权文字与备案号

## 成功标准与路线图（P-001）

- [x] **S1 · 改进项台账与范围冻结**：建立 W16-F01～W16-F10 详细台账（D-001）与各模块改动范围，确认排除项与纳入项边界。
- [x] **S2 · 技术方案与接口设计**：明确各改进项的 Schema 契约、API 路由、前端交互与存储迁移方案（D-002）；关闭 I-001。
- [x] **S3 · 实施分批与子目标规划**：按风险与模块关联度完成分批规划（D-003：批 A/B/C）；先创建批 A 子目标 [GOAL-025-w16-rectification-batch-a](../GOAL-025-w16-rectification-batch-a/00-meta.md)，批 B/C 渐进添加。
- [x] **S4 · 阶段审计与就绪确认**：完成台账与规划的自审（A-002 pass），确认门禁满足，进入批 A 实施。
- [x] **R1 · 整改批 A 完成**——安全与认证基线（W16-F01 / W16-F07 / W16-F08），子目标 [GOAL-025](../GOAL-025-w16-rectification-batch-a/00-meta.md)（done · 4/4）。
- [x] **R2 · 整改批 B 完成**——核心资产与数据交互（W16-F02 / W16-F03 / W16-F04），子目标 [GOAL-026](../GOAL-026-w16-rectification-batch-b/00-meta.md)（done · 4/4）。
- [ ] **R3 · 整改批 C 完成**——系统运维与通用外观（W16-F05 / W16-F06 / W16-F09 / W16-F10），子目标渐进添加。
- [ ] **S5 · 关门**：全部整改子目标 done + 终审 + goal-tree/workspace 同步为 done。

> 当前进度：6/8（S1～S4 + R1 + R2 完成；R3/S5 未完成，GOAL-024 保持 active）。

## 审计策略

| 阶段 | 模式 | 说明 |
|------|------|------|
| S1～S3 规划与台账 | self | 台账与接口方案审查（A-001/A-002） |
| S4 阶段审计 | self / independent | 实施前台账与分批规划审核；安全相关批 A 实施时按子目标升级 independent |
| R1～R3 整改批 | 按子目标 | 每个子目标独立方案/实施/回归/审计关门 |

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|
| I-001 | required | W16-F01~F10 各项在现有前端 Renderer 体系中的兼容实现路径（是否需要新增 custom-component 或仅扩展原生 schema） | S2 方案设计 | S2 | 查阅 renderer 与 form-controls 机制 | **verified** | D-002 §2：custom-node 注册表 + schema 表格/表单扩展两条路径；`custom-components.ts`、`schema-table.tsx`、`form-controls.ts(x)` 已核验 |

## 父目标

- [GOAL-001-design-implementation-conformance](../GOAL-001-design-implementation-conformance/00-meta.md)

## 子目标（整改承接，渐进添加）

| 批 | 子目标 | status | 范围 |
|----|--------|--------|------|
| A | [GOAL-025-w16-rectification-batch-a](../GOAL-025-w16-rectification-batch-a/00-meta.md) | active · 0/4 | W16-F01 首次改密 / W16-F07 一键下线其他 / W16-F08 验证码刷新 + MFA 备份 |
| B | [GOAL-026-w16-rectification-batch-b](../GOAL-026-w16-rectification-batch-b/00-meta.md) | done · 4/4 | W16-F02 文件预览复制 / W16-F03 导入模板与错误定位 / W16-F04 金额格式化与调账警示 |
| C | 待渐进添加（预计 GOAL-027） | — | W16-F05 Cron 预览 / W16-F06 监控自动刷新 / W16-F09 字典 Badge / W16-F10 页脚版权 |

## 台账布局

- `01-decision/`：D-001 等决策与台账；
- `02-execution/`：E-001 等执行事实；
- `03-audit/`：A-001 等审计报告；
- `attachments/`：长文附件（可选）。
