---
id: A-002-response
doc: audit-response
goal: GOAL-012-w11-mfa-ux-review
source: self (orchestrator response)
date: 2026-08-15
status: closed
parent: GOAL-001-design-implementation-conformance
created: 2026-08-15
updated: 2026-08-15
version: 0.1.0
---

# A-002 响应 · findings 闭合记录（P-003）

响应对象：A-002-s2-s4-independent（grok build · grok-4.6 · reasoning high，2026-08-15，verdict: conditional）。

## F-001（required）→ fixed

- 主张：I-004（Toast / 搜索是否扩协议）最晚阶段 S4 仍 open，无 D-003。
- 响应：审计运行期间 D-003 已落盘并闭合 I-004（Toast 本地 UI；搜索复用既有 search-form 模式；select 筛选因后端无 filters 解析留 P2）。证据：01-decision/D-003-s4-scope-confirmation.md；00-meta / 01-decision.md 信息表 I-004 = closed。
- **fixed**（可核对：D-003 文件 + 索引行）。

## F-002（recommended）→ fixed

- 主张：本地 optionsSource 字符串形态与上游 registry 对象形态（{url, labelField, valueField}，component-registry.json since 0.2）同名不同形。
- 响应：渲染层与 schema 全部改为上游对象形态——form-controls.ts 类型、render.ts 解析白名单、useDynamicOptions（url 校验 + params 拼接 + valueField/labelField 映射）、users.json（/api/roles + params.pageSize=100 + key/name）、roles.json ×4（/api/permissions key/key、/api/menu-items id/label）。测试更新并通过。
- **fixed**。

## F-003（recommended）→ fixed

- 主张：回收站 ListItems 支持 q 但未加搜索表单。
- 响应：recycle-bin.json 补 search 表单（targetTable recycle-bin-table）；body 由裸 table 改为 section 包裹。
- **fixed**。

## F-004（recommended）→ fixed

- 主张：Toast 仍占文档流。
- 响应：FeedbackRegion 改为 fixed 右上浮动（right-4 top-16 z-50 max-w-sm shadow-lg），不占布局。
- **fixed**。

## F-005（recommended）→ fixed

- 主张：QR 静区声明不实（编码器矩阵无静区）。
- 响应：qr-code.tsx 增加 4 模块白色静区（viewBox 扩展 + 模块偏移）。
- **fixed**。

## F-006（recommended）→ fixed

- 主张：fakeMFAService.RotateRecovery 不校验码，缺错码测试。
- 响应：fake 校验 123456/RECOVERY；TestMFASelfService 新增 rotate 错码 400 断言。
- **fixed**。

## F-007（recommended）→ fixed

- 主张：目录端点缺 403 用例。
- 响应：rbac_catalog_test.go 新增无 roles.read 用户访问两端点 → 403 断言。
- **fixed**。

## 结论

全部 findings 已按 fixed 路径闭合（可核对修正）；A-002 verdict 由 conditional 解除。认证语义部分（401/400 分轨）grok 已核对通过，无需重开 independent。