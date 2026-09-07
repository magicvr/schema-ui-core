---
doc_type: goal-audit
id: A-002-r4-f001-response-self
parent: GOAL-001-telegram-operator-console
date: 2026-09-05
source: self
auditor: Codex govern
audit_type: finding-remediation-response
scope: 响应 A-001 F-001；polling 单实例 UI 警示与 Root R4 边界复核
verdict: pass
open_required: 0
version: 0.1.0
---

# A-002 · Root R4 F-001 修复响应（2026-09-05）

## 意见汇总

| 意见 | source | verdict | open required | 当前处理 |
|------|--------|---------|---------------|----------|
| A-001-r4-root-close-self-audit | self | fail | 1 | F-001 要求补齐 polling 单实例/丢 Update 的 UI 警示 |
| 本条 | self | pass | 0 | F-001 已以代码、双语文案和测试 fixed；无 residual/overrule |

## F-001 fixed 证据

- `apps/web/src/components/telegram-admin-tab.tsx:678-684` 在 `status.mode === polling`
  时渲染 `role="alert"` 与 `data-telegram-polling-warning`，文案明确单实例、多副本
  丢 Update 风险和非 HA 生产定位。
- `apps/web/src/i18n/messages/en-US.json:998-1001` 与
  `apps/web/src/i18n/messages/zh-CN.json:998-1001` 提供对应双语文案。
- `apps/web/src/components/telegram-admin-tab.test.tsx:118-150` 核对 polling 警示出现、
  webhook 状态下警示不出现；catalog 与 bilingual tests 通过。
- `E-028-r4-polling-warning-remediation.md` 记录了本次修复与完整 Web 验证；此前
  `da9d955e` 修复的 `form-controls.tsx:946-947` 构建错误仍保持已修复，未留下基线错误。

## 门禁结论

F-001 满足 `fixed` 合法闭合路径，`open_required: 0`。本响应不接受残余、不作
`user-overruled`，也不直接关闭 Root；Root 仍需由指定 provider 完成独立 R4 close-out。
