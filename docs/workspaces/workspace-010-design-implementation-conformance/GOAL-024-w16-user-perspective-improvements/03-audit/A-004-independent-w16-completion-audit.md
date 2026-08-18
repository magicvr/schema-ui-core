---
id: A-004
goal: GOAL-024-w16-user-perspective-improvements
title: 独立审计 · GOAL-024 整体完成情况与 W16-F01～W16-F10 落地核实
source: independent
date: 2026-08-17
verdict: pass
scope: GOAL-024 整体完成情况、W16-F01～W16-F10 实施落地、子目标分批关门与回归证据
---

# A-004 · 独立审计 · GOAL-024 整体完成情况与 W16-F01～W16-F10 落地核实

- **source**: independent
- **auditor**: gemini-3.7-flash-high
- **type / scope**: close-out / execution-facts · GOAL-024 全目标完成情况与 W16-F01～W16-F10 落地核验
- **verdict**: **pass**

## 1. 范围与区间

- **covered**:
  - `GOAL-024` 目标定义、台账界定（D-001）、技术方案（D-002）、分批规划（D-003）、信息项（I-001）与关门自审（A-003）；
  - 下级整改子目标完成情况：批 A `GOAL-025`（F01/F07/F08）、批 B `GOAL-026`（F02/F03/F04）、批 C `GOAL-027`（F05/F06/F09/F10）；
  - 10 项核心体验改进点（W16-F01～W16-F10）的实际代码实现、数据库迁移、API 契约、前端组件与自动化测试覆盖；
  - Web 端全量 vitest 与 tsc 类型检查，Go 端单元测试与 W16 专门测试套件。
- **excluded**:
  - `workspace-011` 规划的外部功能项（B-10 组织架构、S-05 公告管理、B-09 消息模板、B-11 登录日志视图等，已由 D-001 明确排除）。

## 2. 成果与证据（逐项核验）

| 改进项 ID | 改进主题 | 涉及子目标 | 后端/存储证据 | 前端/UI 证据 | 测试与验证证据 |
|---|---|---|---|---|---|
| **W16-F01** | 首次登录强制修改初始密码 | GOAL-025 (批 A) | Migration 0038 (`users.must_change_password`)；`internal/handler/account_self.go`（拒绝新旧同密码，改密后重置标记并 bump session）；`internal/auth/auth.go` 门禁 | `apps/web/src/components/force-password-change.tsx`；`apps/web/src/app/LoginPage.tsx` 识别并拦截 | Go: `TestForcedPasswordChangeGateAndReissue`；Web: `force-password-change.test.tsx` (2 tests)、`LoginPage.test.tsx` (13 tests) |
| **W16-F02** | 文件库在线直接预览与一键复制直链 | GOAL-026 (批 B) | `apps/api/internal/handler/filelibrary.go` (`downloadUrl` 随文件行下发) | `apps/api/internal/modules/filelibrary/schema/file-library.json` (`preview`/`copyLink` actions)；`apps/web/src/renderer/render.tsx` | Go: `TestFileLibraryRowsExposeDownloadUrl`；Web: `download-behavior.test.tsx` |
| **W16-F03** | 数据导入 CSV 模板与逐行错误定位 | GOAL-026 (批 B) | `apps/api/internal/handler/import.go` (`GET /api/import/{resource}/template` 返回 CSV；导入响应返回 `fieldErrors: [{ rowNumber, field, reason }]`) | `apps/web/src/renderer/form-controls.tsx` (fieldErrors 定位高亮) | Go: `TestImportTemplateUsers`、`TestImportFieldErrors`；Web: vitest 覆盖 |
| **W16-F04** | 钱包金额分转元格式化与调账警示 | GOAL-026 (批 B) | `apps/api/internal/modules/wallet/schema/wallet.json` (`format: "currency"`) | `apps/web/src/renderer/schema-table.tsx` (`format: "currency"` 分转元处理)；`form-controls.tsx` (调账警示提示) | Web: `schema-table.test.tsx` (货币格式化用例) |
| **W16-F05** | Cron 自然语言解析与未来运行预览 | GOAL-027 (批 C) | `apps/api/internal/handler/scheduledtasks.go` (`POST /api/scheduled-tasks/cron/preview` 返回自然语言描述与接下来 3 次运行时间) | `apps/web/src/components/cron-preview.tsx` (输入防抖与预览渲染) | Go: `TestCronPreviewEndpoint` |
| **W16-F06** | 系统监控页面定时自动轮询刷新 | GOAL-027 (批 C) | 无需后端改动 (复用既有监控接口) | `apps/web/src/components/monitoring-auto-refresh.tsx` (关闭/5s/10s/30s 下拉与定时 refetch) | Web: vitest 1056/1056 全量 |
| **W16-F07** | 个人中心一键下线其他所有设备 | GOAL-025 (批 A) | `apps/api/internal/handler/account_self.go` (`POST /api/account/sessions/revoke-others`)；`internal/account/session.go` (`RevokeOtherSessions`) | `apps/web/src/components/account-session-toolbar.tsx` (下线其他按钮与会话列表即时刷新) | Go: `TestRevokeOthersReissuesTokensAndRevokesOtherSessions`；Web: `account-session-toolbar.test.tsx` |
| **W16-F08** | 验证码主动刷新 + MFA 密钥复制与恢复码 txt 下载 | GOAL-025 (批 A) | 复用既有验证码与 MFA 接口 | `apps/web/src/app/LoginPage.tsx` (点击验证码重新请求)；`apps/web/src/components/mfa-manager.tsx` (一键复制密钥与导出 txt 恢复码) | Web: `mfa-manager.test.tsx`、`LoginPage.test.tsx` |
| **W16-F09** | 数据字典项 Badge 颜色/标签风格配置 | GOAL-027 (批 C) | Migration 0039 (`dict_entries.badge_style`)；`internal/modules/datadictionary/` 契约与 CRUD 端点 | `apps/web/src/renderer/schema-table.tsx` (`badgeStyleField` 映射 Tag 变体) | Go: `TestDictEntryBadgeStylePersists`；Web: `schema-dictionary-entries.test.tsx` |
| **W16-F10** | 系统设置支持页脚版权与备案号 | GOAL-027 (批 C) | Migration 0040 (`site_settings.copyright_text`, `icp_number`)；`internal/modules/settings/` 契约与 `/api/branding` 端点 | `apps/web/src/app/branding.ts`、`apps/web/src/app/App.tsx` (Shell 页脚渲染) | Go: `TestSettingsFooterFieldsPersist`；Web: `startup-config.test.tsx` |

## 3. 对照成功标准

| 阶段 / 检查点 | 定义 | 达成状态 | 验证与证据 |
|---|---|---|---|
| **S1 · 改进项台账与范围冻结** | W16-F01～F10 详细台账与排除项冻结 | **完成** | `01-decision/D-001-w16-improvement-inventory.md`、`02-execution/E-001-goal-establishment.md` |
| **S2 · 技术方案与接口设计** | 各项 Schema 契约、API 路由、前端交互设计；关闭 I-001 | **完成** | `01-decision/D-002-w16-technical-design.md`、`I-001` 状态已为 `verified` |
| **S3 · 实施分批与子目标规划** | 批 A/B/C 分批与渐进子目标规划 | **完成** | `01-decision/D-003-w16-batch-subgoals.md` |
| **S4 · 阶段审计与就绪确认** | 台账与分批方案自审通过 | **完成** | `03-audit/A-002-self-planning-review.md` (pass) |
| **R1 · 整改批 A 完成** | 安全与认证基线 (F01/F07/F08)，子目标 GOAL-025 (4/4) | **完成** | GOAL-025 done；含 independent 审计 A-001 及闭合响应 A-002 (F-001/F-002 fixed) |
| **R2 · 整改批 B 完成** | 核心资产与数据交互 (F02/F03/F04)，子目标 GOAL-026 (4/4) | **完成** | GOAL-026 done；A-001 关门自审 pass；测试全绿 |
| **R3 · 整改批 C 完成** | 系统运维与通用外观 (F05/F06/F09/F10)，子目标 GOAL-027 (4/4) | **完成** | GOAL-027 done；A-001 关门自审 pass；测试全绿 |
| **S5 · 关门** | 全部整改子目标 done + 全量回归 + 关门审计 | **完成** | `03-audit/A-003-closeout.md`；`goal-tree.md` / `workspace.md` 已同步 |

## 4. Findings

### F-001 · GOAL-024 00-meta.md 中的子目标表格 GOAL-025 状态残留未同步
- **level**: recommended
- **severity**: low
- **status**: open
- **evidence**: `GOAL-024/00-meta.md` 第 70 行表格中显示 `GOAL-025-w16-rectification-batch-a | active · 0/4`，而实际上 GOAL-025 在 S5 关门检查点及 `goal-tree.md` 中已正确标记为 `done · 4/4`。
- **remediation**: 在后续常规文档维护时将该单元格同步为 `done · 4/4`。

### F-002 · GOAL-024 索引文件 frontmatter status 仍为 active
- **level**: recommended
- **severity**: low
- **status**: open
- **evidence**: `01-decision.md`、`02-execution.md`、`03-audit.md` 索引文件的 frontmatter `status: active`，而目标主表 `00-meta.md` 已为 `done`。
- **remediation**: 可在编排器维护索引元数据时统一刷新。

## 5. 必改项汇总（Required Findings）

- **开放必改项数**：0
- 无阻断关门或影响质量门禁的 required 项。

## 6. 与既有意见的异同

- **GOAL-024 既有自审 (A-001, A-002, A-003)**：结论均为 pass。本次独立审计通过独立重跑自动化测试套件（Go W16 测试集、Web 69 文件 1056 个 vitest、tsc 类型检查）和查阅实际 git 变更，确认自审结论真实客观、证据确凿。
- **GOAL-025 历史独立审计 (A-001 grok)**：曾提出 2 项 required findings（F-001 刷新会话列表、F-002 拒绝同密码改密），核查确认已在 A-002 中通过代码修改与测试用例完全闭合（`fixed`）。

## 7. 结论与建议

- **审计结论**：**PASS**
- **审视评价**：`GOAL-024-w16-user-perspective-improvements` 及其所辖的三个整改子目标（`GOAL-025`、`GOAL-026`、`GOAL-027`）严格遵循 P-001～P-005 治理原则。10 项未计划改进项均具备完备的台账记录、技术方案设计、高质量代码落地与双端自动化测试。目标完成事实真实充分。
- **建议下一步**：
  1. 通过 `/govern` 留痕并响应本次独立审计意见；
  2. 顺带微调 F-001（`00-meta.md` 子目标表格文字残留）；
  3. 确认 workspace-010 当前状态，继续推进工作区内的后续任务或进行工作区总结。

## 8. 声明

本意见为独立交叉审计意见（`source: independent`），不修改目标 `status`/`progress`；响应与状态维护归由 `/govern` 处理。