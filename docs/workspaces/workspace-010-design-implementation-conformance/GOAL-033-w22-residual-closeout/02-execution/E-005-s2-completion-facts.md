---
title: E-005 · S2 完成事实：A2/A4/A5/A6/H3 落地与验证；e2e 选择器回归修复
status: recorded
created: 2026-08-23
updated: 2026-08-23
parent: GOAL-033-w22-residual-closeout
version: 0.1.0
---

# E-005 · S2/S4 完成事实（2026-08-23）

子代理多次在长测试阶段静默中断，其代码改动完整保留于工作树；本目标编排器直接接管验证并收尾。

## C3 · A2 种子 admin 迁移（W10·GOAL-025 F-003）—— 完成

- 实现：`apps/api/internal/modules/authsession/migration/migration.go` 新增迁移 **0049**（`migrateSeedAdminMustChangePassword`，+29 行）；测试 `apps/api/internal/store/migrate_0049_test.go`。
- 验证（本编排器复跑）：`go build ./...` 干净；`TestMigrate0049BackfillsSeedAdminMustChangePassword` PASS、`TestMigrate0049IdempotentOnFreshDatabase` PASS（store 包 ok 2.388s）。

## C5 · A4 密码可见切换（W6·F-VUI-011）—— 完成

- `LoginPage.tsx` 增加 Eye/EyeOff toggle（lucide-react 既有依赖，`[data-password-toggle]`，aria-label 接 i18n）；i18n 双语 key：`login.password.show|hide`（显示密码/Show password 等）；组件测试 +3（LoginPage.test.tsx 16/16 过、tsc -b --noEmit 零错）。
- 连带整改见下节「e2e 选择器回归」。

## C6 · A5 上传内容嗅探（W9·GOAL-002 N-002）—— 完成

- `upload.go`（+36 行）：入库嗅探 `<svg`/`<script` 标记与 `on*=` 事件处理器形态（含嗅探窗口边界），命中拒绝入库；下载侧安全头不变。
- 测试：`TestContainsActiveContentEventHandler`、`TestContainsActiveContentSniffWindowBound`、`TestUploadA5ContentSniffEndToEnd` 全 PASS；upload 包既有契约测试（attachment/owner-only/非 hex id/files.write/quota）12/12 PASS（handler 包 ok 6.867s）。

## C7 · A6 MFA verify 独立限流（W9·GOAL-011 R-001）—— 完成

- `mfa.go`（+43 行）：verify 端点独立限流桶；测试 `TestMFAVerifyRateLimit`、`TestMFAVerifyRateLimitPerIP`、`TestMFAVerifyRateLimitDoesNotBlockNormalFlow` 全 PASS。

## C16 · H3 dogfood 重命名（W10·W1 F-006）—— 完成（编排器接管）

- 子代理两轮中断留下半成品（新文件已建、引用未改）。接管完成：8 个测试文件 9 处引用 `app-manifest.(admin|mvp).json` → `app-manifest.(admin|mvp)-dogfood.json`；原文件删除；残留引用 0。内容逐字节不变（纯改名）。vitest 全量验证见 E-006。

## 连带整改 · e2e 密码定位符回归修复（A4 引入）

A1 补跑日志证实：新增切换按钮 aria-label「显示密码」使既有 `getByLabel("密码")` 子串匹配命中双元素触发 Playwright strict-mode（`localization.spec.ts:20`）。修复：`localization.spec.ts`（2 处）、`sign-in.ts`（2 处）、`force-password-change.spec.ts`（1 处）共 5 处改为 `{ exact: true }` 并注明缘由。该修复属本波自引入回归的必要纠正，随关门审计一并送审。

## 环境结论（A1 前置）

A1 日志（`attachments/e2e-admin-m3-rerun.log`）证实：8011–8110 排除区间消失、API(25080)/Vite(25173) 均成功启动 —— W7 原环境 residual 的阻塞**已解除**，仅剩选择器问题（已修复），补跑重验见 E-006。

## 进度

C3/C5/C6/C7/C16 完成 → 累计 **14/18**。剩余：C2（e2e 补跑重验）、C17（全量回归）、C18（关门审计）。
