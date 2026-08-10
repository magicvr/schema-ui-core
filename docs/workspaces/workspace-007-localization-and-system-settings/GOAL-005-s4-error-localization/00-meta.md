---
id: GOAL-005-s4-error-localization
title: S4 · 后端用户可见反馈本地化（稳定错误码 + 有界服务端协商）
status: done
parent: GOAL-001-localization-and-system-settings
created: 2026-08-09
updated: 2026-08-09
version: 0.2.0
progress: 5/5
---

# GOAL-005 · S4 · 后端用户可见反馈本地化

## 概述

承接 Root [GOAL-001-localization-and-system-settings](../GOAL-001-localization-and-system-settings/00-meta.md) 的 **S4 阶段**：按 I-L10N-004 用户书面选定的 **exit 5 路径 (a) 有界服务端 locale 协商**实施——错误码保持稳定可机读；错误 envelope 兼容扩展 `{error, message, messageKey?, params?}`；已编目错误（认证/验证/设置/资源/上传）按 `Accept-Language` 返回对应语种 `message` 并声明 `Content-Language`；未编目/`INTERNAL` 错误保持英文通用文案（安全回退、不泄露诊断）；前端按码/key/参数以当前语种呈现用户可见反馈（不可降级保底）。

**方案依据**：Root D-002 §I-L10N-004（路径 a）+ 附录 A（错误码枚举钉死）。本目标只实施与验证，不重新决策。

**范围纪律**：不翻译开发日志/内部诊断文本；不以翻译文案替代机读错误码；不泄露堆栈/内部细节；`INTERNAL` 永不本地化。

## 成功标准（可验收 · 等权检查点 · 共 5 项）

- [x] **C1**：错误码契约测试：全仓 `writeError`/`writeLocalizedError`/DomainError 码集合与 D-002 附录 A 枚举一致（31 字面量 + 8 域码族），可回归防漂移。
- [x] **C2**：服务端协商：已编目码按 `Accept-Language`（zh-CN/en-US）返回对应语种 `message` + `messageKey`/`params` + `Content-Language` 头；未编目/`INTERNAL` 返回英文通用文案且无 `messageKey`；失败回退（无支持语种 → en-US）。
- [x] **C3**：认证/验证/设置错误端到端：登录失败（INVALID_CREDENTIALS 经协商）、设置校验（INVALID_TIMEZONE 等）、资源验证错误双语可验证（Go 测试驱动真实 handler）。
- [x] **C4**：前端保底：前端请求携带当前语种（Accept-Language）；`readResourceApiError` 解析 `messageKey`/`params`；反馈呈现优先 catalog（按码/key/参数）→ 服务端 message 回退；未编目错误安全呈现且不泄露诊断。
- [x] **C5**：验证：Go 全绿（新增协商/契约测试）；vitest 全绿（新增前端保底测试）；`npm run build` 通过；证据捕获 `{SCRATCH}`；F-V029 错误回退列证据回填。

## 派生进度展示

`progress: 5/5` 由上方 5 个等权检查点派生；仅为展示，不放行阶段、不关闭 finding。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 错误码全集、编目范围与协商/回退语义是否齐备 | C1–C5 | 实施前 | 读 Root D-002 §I-L10N-004 + 附录 A + 全仓 writeError 盘点 | **closed** | — | Root D-002 冻结（2026-08-09 用户裁决路径 a） |

## 父目标

- [GOAL-001-localization-and-system-settings](../GOAL-001-localization-and-system-settings/00-meta.md)（Root；本目标为 S4 阶段子目标）

## 台账布局

本目标使用 ledger 目录：`01-decision/`、`02-execution/`、`03-audit/`。索引文件保留 frontmatter 与条目表；独立记录使用 `D-NNN-*`、`E-NNN-*`、`A-NNN-*`。
