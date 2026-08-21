---
id: GOAL-002-s1-locale-core
doc: audit-entry
record_id: A-001
source: self
auditor: schema-ui-core 编排器（grok build）
audit_type: execution-facts
scope: S1 实施 · C1–C6 检查点证据
verdict: pass
status: recorded
parent: GOAL-002-s1-locale-core
created: 2026-08-09
updated: 2026-08-09
version: 0.1.0
---

# A-001 · S1 多语种核心自审（execution-facts）

## 范围与区间

- scope：GOAL-002（S1）全部实施产物与 C1～C6 检查点；不含页面文案双语化（S2）与 API（S3/S4）。
- 依据：Root D-002（I-L10N-001/002 冻结语义）+ 本目标 D-001（单元边界）。

## 成果（有证据）

| 检查点 | 证据 | 判定 |
|--------|------|------|
| C1 locale 解析/优先级 | `src/i18n/locale.test.ts`（12 断言：显式 > 系统默认 > 浏览器 > en-US；BCP47 归一化；auto 语义；登录前后同一函数） | ✓ |
| C2 资源装载 | `messages/en-US.json`/`zh-CN.json` 纯数据；`catalog.test.ts` 完整性断言；装载失败降级路径=回退链 | ✓ |
| C3 缺失 key 可观察+安全回退 | `catalog.test.ts`：`schema-ui:missing-translation` 事件（locale/key/path）、按 locale+key 去重、回退链（当前→en-US→key）、不渲染空 | ✓ |
| C4 用户切换 | `locale-switcher.test.tsx`（可达性=仅需 Provider、无权限 prop、标签双语）+ `runtime.test.tsx`（localStorage 单通道、auto 移除键、立即生效） | ✓ |
| C5 lang + 格式化 | `runtime.test.tsx`（`<html lang>` 挂载/切换跟随）+ `format.test.ts`（zh-CN/en-US、IANA 时区、无效时区/时间戳安全降级） | ✓ |
| C6 验证 | vitest 674/674 全绿（含新增 45）；`npm run build` exit 0；输出捕获 `{SCRATCH}/unit-s1-web.log` | ✓ |

## 对照成功标准

- C1～C6 全部有 shipped-函数级证据；无 mock 被测单元、无测试专用复制实现。
- 与冻结语义一致：localStorage 单通道（登出不清除——runtime 无任何登出钩子，天然满足）；优先级函数在 AuthGate 前后同一实例。

## Findings

- F-001（recommended）：`formatDate` 的 `dateStyle/timeStyle` 为 Intl 简写（S3 时区语义落地时确认与 `siteTimezone` 组合正确——非本阶段阻断）。
- 无 required findings。

## 必改项汇总

无。

## 结论

**verdict: pass** — S1 检查点全部有可重复证据；`I-L10N-001/002` 冻结语义已落地为 shipped 运行时。S2 可在本阶段之上实施页面双语化与 `titleKey`/`labelKey` 真解析。

## 声明

本意见不修改 status/progress；响应与放行由编排器处理。
