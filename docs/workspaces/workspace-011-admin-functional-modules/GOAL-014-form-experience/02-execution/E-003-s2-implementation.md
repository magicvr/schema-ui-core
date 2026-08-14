---
id: E-003
goal: GOAL-014-form-experience
date: 2026-08-14
status: recorded
parent: GOAL-014-form-experience
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# E-003 · S2 实现完成（字段级错误 + 前端校验 + 布局）

## 事实

- 2026-08-14：S2 实现完成（D-002 §2/§3/§4）。

## 服务端（fieldErrors）

- errorcatalog：`Body` 对字段校验码（INVALID_CREATE_FIELD/INVALID_PATCH_FIELD）拼接具体字段原因（其余码保持纯本地化文本，不泄露诊断）；新增 `FieldError{field,reason}` 与 `BodyWithFields`（可选 fieldErrors 键，向后兼容）。
- handler/localize.go：新增 `writeLocalizedFieldError`（cataloged 走 BodyWithFields，uncataloged 也附 fieldErrors）。
- resources.go：create/patch 字段校验失败路径改走 writeLocalizedFieldError（fe.field/fe.reason 结构化输出）。

## 前端

- resource.ts：`FieldError` 接口 + `ResourceApiError.fieldErrors` + readEnvelope 解析 fieldErrors。
- form-controls.ts：`FormControlField` 扩展 required/pattern/minLength/maxLength；`validateFieldValues` 纯函数校验器（REQUIRED/PATTERN/MIN_LENGTH/MAX_LENGTH/MIN_VALUE/MAX_VALUE；布尔字段跳过 required）。
- form-controls.tsx：`FormControls` 加 `fieldErrors`/`columns` props；**默认单列**（移除硬编码两列），columns>1 响应式网格（上限 4，移动端单列）；FieldControl 支持 error 内联展示。
- render.tsx：handleSubmit 提交前 validateFieldValues（失败阻止提交 + 字段内联）；服务端 fieldErrors 回显到字段；表单 columns 透传。

## Schema 示范

- data-dictionary.json：create/edit 表单 key/name 加 required（name 加 min/maxLength 1..64）。
- dictionary-entries.json：dictKey/entryKey/label 加 required。

## 测试

- form-controls.test.ts +4（validateFieldValues：required/布尔跳过/pattern/长度/数值/非必填空值跳过）。
- 全量 web vitest 911/911 绿；API 回归见提交时状态。

## 遗留

- S3 验证记录（含真实 HTTP fieldErrors 冒烟）、S4 go 判定、S5 grok 审计。
