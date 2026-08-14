---
id: D-001
goal: GOAL-006-w6-scan-findings-remediation
title: W6 修复范围与技术取舍
date: 2026-08-15
status: accepted
parent: GOAL-001-production-hardening
created: 2026-08-15
updated: 2026-08-15
version: 0.1.0
---

# D-001 · W6 修复范围与技术取舍

## 背景

2026-08-15 本会话对 `apps/api` + `apps/web` 逐文件审视，确认以下低危缺陷（本波承接；中高危已在 W1–W5 处理完毕，W5 扫描 0 中高危）：

| # | 位置 | 事实 | 定级 |
|---|------|------|------|
| F1 | `apps/api/internal/modules/scheduledtasks/scheduler.go` tick()（`cron.go` Next()） | 每 30s tick 对每个启用任务调用 `fields.Next(now)`；未到期时该函数以分钟步进扫描至多 5 年（最多约 262 万次迭代）。tick 仅需判断“当前分钟槽是否匹配”，应改用 `fields.Matches(slot)` | 低危 · 性能（CPU 空转） |
| F2 | `apps/api/internal/modules/recyclebin/service.go` Restore/isConflict + `handler/recyclebin.go` writeRecycleError | 还原孤儿 dict entry（父 `dict_types` 已删除）时 `CreateEntry` 返回 `ErrDictKeyNotFound`，`isConflict` 未识别 → 500 INTERNAL | 低危 · 错误映射 |
| F3 | `apps/web/src/app/branding.ts` isSafeBrandingUrl | 仅接受同源路径与 http(s)；管理员配置 Base64 内联小图标（`data:image/png;base64,…`）会被置空。注：API 侧 normalizeLogoURL 亦需核对 | 低危 · 可选增强 |

## 决策

1. **F1 采纳修复**：`tick` 内以 `slot := now.Truncate(time.Minute)`；先 `if !fields.Matches(slot) { continue }`，仅当匹配时才沿用原执行路径（无需 Next 预判；lastRun 分钟槽去重已保证单次执行）。保留 `Next()` 供未来“下次运行时间”展示用，不删除。附回归测试：不匹配槽不触发 Execute、匹配槽触发一次。
2. **F2 采纳修复**：`isConflict` 增加 `store.ErrDictKeyNotFound` 识别，映射为 `DomainError{Status: 409, Code: "DICT_KEY_NOT_FOUND", Message: "parent dict type does not exist"}`（409 与既有 RECYCLE_RESTORE_CONFLICT 语义一致；snapshot 保留可重试）。附回归测试：孤儿 entry 还原返回 409 而非 500。
3. **F3 不采纳（user-overruled）**：API `normalizeLogoURL`（settings/repository.go:276）与 errorcatalog `INVALID_LOGO_URL` 均明确拒绝 data: URI（"必须为空、同源路径或 http(s) URL"）；web `startup-config.test.tsx:212` 已断言 `data:image/svg+xml` 为 false。拒绝 data: 是 web/API 一致的有意安全收紧（防 SVG 脚本载荷 / XSS 混淆面），保持现状，本项不再实现。
4. 本波不改协议 pin、Profile 默认集、模块矩阵与 Manifest 装配（延续 W5 go 判定）。

## 未选方案

- F1 引入“下次运行时间”缓存/持久化：超出本波范围，无产品需求支撑。
- F2 自动跳过孤儿快照或删除快照：数据安全优先，保留快照 + 明确错误，用户可先恢复父类型再重试。