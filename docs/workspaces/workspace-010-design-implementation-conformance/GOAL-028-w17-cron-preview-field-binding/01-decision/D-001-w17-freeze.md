---
id: D-001
goal: GOAL-028-w17-cron-preview-field-binding
title: W17 方案冻结：Cron 字段绑定与中文 describeCron
status: accepted
created: 2026-08-18
updated: 2026-08-18
version: 0.1.0
parent: GOAL-001-design-implementation-conformance
---

# D-001 · W17 方案冻结

## 1. 触发

GOAL-024 A-007 确认 A-005 F-004 不能标 `fixed`：`cron-preview` 仍是页面独立输入，不读创建/编辑表单的 `cron`；`describeCron` 仍返回 `"every minute"` / `"every hour at minute N"` / `"cron schedule (5-field)"`。用户 2026-08-18 选择开后续波次补这两点，不重开 GOAL-024。

## 2. 决定

### 2.1 字段绑定

- 创建/编辑任务模态的 `cron` 输入增加**本地**字段属性 `afterComponent: "cron-preview"`。
- `gateRenderFormFields` 透传该属性；`FormControls` 在字段控件下方渲染已注册 custom component，并把当前字段值作为 `node.props.bindValue` 传入。
- `cron-preview` 在存在 `bindValue` 时进入绑定模式：不再渲染独立输入/提交，只对绑定值做 400ms 防抖预览。
- 从 `scheduled-tasks.json` 页面 `body` 移除 `cron-preview-block`。
- **不**新增协议 field `type`，不改 Host/App 契约。

### 2.2 describeCron

- 预览端点按 `Accept-Language` 协商（`errorcatalog.Negotiate`，与错误本地化同一套：`zh-CN` / `en-US`）。
- `zh-CN` 必须是中文人话，禁止再返回英文 stub。
- 至少覆盖：每分钟、每小时第 N 分钟、每天 HH:MM、每周X HH:MM、每月 D 日 HH:MM、每 N 分钟；其余回退为「5 段 Cron 计划」/ `5-field cron schedule`（不再用 `cron schedule (5-field)` 这种实现者口吻）。
- 响应形状仍是 `{ description, nextRuns }`。

## 3. 为什么

- GOAL-027 D-001 / GOAL-024 D-002 原方案就是「Cron 字段下方挂预览 + 中文描述」；W16 只交付了独立控件与英文 stub。
- `afterComponent` 把绑定留在 Host 表单层，避免为预览发明协议控件或把表单值塞进页面 context。

## 4. 未选方案

| 方案 | 未选理由 |
|------|----------|
| 在 Renderer 为 `cron` 新增官方 field type | 超出本波；会碰协议白名单 |
| 保留页面独立预览块 + 同步表单 | 双输入源，用户仍会对「另一个框」 |
| `description` 只出中文、忽略语种 | 与现有 Accept-Language 协商不一致 |
| 书面 residual 不改代码 | 用户已点名开波次补 |

## 5. 影响

- 前端：`form-controls.ts(x)`、`cron-preview.tsx`、`scheduled-tasks.json`、定向 vitest。
- 后端：`describeCron` + 预览端点测例。
- go：不改 Profile / 模块矩阵 / Manifest → 不暂挂。

## 6. 后续

S2 按本决策实施；S3 定向测试至少覆盖绑定预览与 `zh-CN` 描述。
