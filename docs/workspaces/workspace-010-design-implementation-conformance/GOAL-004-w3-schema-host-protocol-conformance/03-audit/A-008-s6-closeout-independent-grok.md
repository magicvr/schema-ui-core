---
id: A-008
goal_id: GOAL-004-w3-schema-host-protocol-conformance
title: 独立审计 · S6 关门审计（grok build · grok 4.6 · xhigh）
source: independent
scope: GOAL-004 全量关门检查（信息项/findings/检查点/residual/go 影响/claim 一致性）
verdict: conditional
provider: grok build（model grok-4.6，reasoning-effort xhigh）
created: 2026-08-13
updated: 2026-08-13
parent: GOAL-004-w3-schema-host-protocol-conformance
version: 0.1.0
---

# A-008 · 独立审计（source: independent · grok build）

> 原文由独立审计会话（grok build，grok 4.6，xhigh）产出，经编排器代贴落盘并保留
> `source: independent`。本轮只读，未修改任何文件。意见不修改 `status`/`progress`。

## 1. verdict: conditional（BLOCKING_COUNT=2）

S1–S5 主体证据成立，claim 三处哈希一致，go 不暂挂可复核。不能无条件关门：account-locked
仍无用户书面 P-004（F-1）；S4-1 生产路径未捕获 `location.search`，与「恢复 path?query」的
完成声明不一致（F-2）。

## 2. 逐项关门检查表（摘要）

| 检查项 | 判定 |
|--------|------|
| I-001/I-002/I-003/I-005/I-006 | 满足（verified；95/95、accepted ADR、正式 pin `521cff8`/`4fae4605…`、migration/registry 弃用机制） |
| I-004 provider | 满足（provider=`grok build` 属实；证据列三处漂移 → F-3） |
| I-007 | non-blocking，不阻断 |
| A-001～A-006 required findings | 满足（三路径闭合；P2 acknowledged 非 required） |
| S1/S2/S3 | 满足 |
| S4 | **不满足**（S4-1 生产捕获丢 query → F-2；其余 S4-2～S4-7 有代码与测试） |
| S5 | 满足（文档证据；vitest 871、tsc 0、Go ok、fixtures 96+41+16、Playwright 7+1skip；本轮未复跑） |
| S6 | 部分满足（A-007 self + 本 independent；开放 required 未闭合前不勾） |
| account-locked 拟议 residual | **不满足（关门前置未决）**：仓库内无用户书面 accept/overrule（F-1） |
| claim 三处一致性 | 满足（content `4fae4605…`/fixture `7aacf133…`/buildId `git:5e4c384…` 逐字一致；HEAD≠buildId 属「先生成再提交」已记录模式；report residuals `[]` 正确） |
| go 影响 | 满足 · **不暂挂**（S4 未改 Profile 默认集/模块矩阵/Manifest 装配语义；2.8 pin 属 S3 additive MINOR） |
| 台账一致性 | 部分满足（I-004 三处不同；A-007 预写 A-008；结论段未吸收 A-007 → F-3/F-4） |

## 3. Findings

### F-1 · P1（required）· account-locked 拟议 residual 无用户书面裁决
- 位置：`01-decision/D-002` §3.1；`02-execution/E-007`；`03-audit/A-007`
- 证据：D-002 写明「S6 须经用户 P-004 书面接受或驳回，本决定不预写接受」；全仓无用户书面
  `accepted-residual`/`user-overruled`。映射层存在（`bootstrap.ts` locked→`ACCOUNT_LOCKED`；
  `failure.ts` 禁 reauth、只允许 home/support；fixtures pin）；`bootstrapAuthFor`/`AuthContext`/
  `auth-client` 从不产出 `locked`。AUTH-008 为 `reserve-extension`，与「本波不做锁定产品」一致，
  但不能代替书面接受。
- 建议：用户书面三选一并落盘：`accepted-residual`（范围=本波不新增账号锁定安全特性；复审触发=
  认证迭代引入锁定状态时）／实现生产源／`user-overruled`。在此之前不得 `status: done`。

### F-2 · P1（required）· S4-1 生产捕获丢弃当前 URL query，完成声明过满
- 位置：`apps/web/src/host/return-intent.ts` `captureReturnIntent`；调用点 `AuthContext.tsx`、
  `boot.ts`；对照 `return-intent.test.ts`、E-007 S4-1、ADR-0036 D6
- 证据：无参调用时 `options?.query ?? {}`，不解析 `window.location.search`；单测全部显式注入
  query，未覆盖生产调用点。`tab`/`view`/`page` 等 allowlist 键在 reauth 后会丢。空 query 协议
  合法（收窄），但 S4-1/E-007 写的是「恢复 path?query」，与生产不符。
- 建议：生产捕获从 `window.location.search` 取 query，再走 `validateReturnIntent`；补生产调用点
  测试。或登记 residual 并改措辞。fixed / accepted-residual / user-overruled 之前 S4 不得视为完成。

### F-3 · P2（recommended）· I-004 / A-007 预支 A-008
- 位置：`00-meta.md` I-004；`03-audit.md` I-004；`A-007` 检查表
- 证据：00-meta 写 A-001/A-002；03-audit 写 A-006/A-005；A-007 把尚未存在的 A-008 写成已落盘。
- 建议：I-004 统一为「provider=`grok build`；S2=A-005/A-006；S6=A-007+A-008」；本意见落盘后
  再写 A-008 已落盘。

### F-4 · P2（recommended）· `03-audit.md` 结论段未吸收 A-007
- 位置：`03-audit.md`「结论状态」末段
- 证据：索引进了 A-007，正文仍写「S6 关门仍需各自 scope 的后续 cross 审计」。
- 建议：补 A-007/A-008 摘要与开放前置。

## 4. residual 判定

**account-locked 是否必须用户决策：是，关门前置。** 建议用户接受 `accepted-residual`（范围=
本波不新增账号锁定安全特性；复审触发=认证迭代引入锁定状态时）。

- `$deps`/页面 2.7 mandatory：已关闭；304/ETag：不是关门 residual（200-only 合规）；return
  intent allowlist 收窄：设计约束已实现，不是缺口。
- 未登记新发现：生产 return-intent 不捕获 `location.search`（F-2）；O-001/O-002/O-003 仍为
  recommended，不阻断。

## 5. 结论

**不可关门**（conditional）。放行前必须：①用户对 account-locked 书面 P-004；②闭合 F-2；
③本 independent 落盘（A-008）并由编排器响应。F-3/F-4 建议同批改。go 不暂挂。progress 保持
5/6，S6 保持未勾，直到上述项闭合。

**BLOCKING_COUNT=2**

## 编排器响应（2026-08-13）

| Finding | 处置 | 说明 |
|---------|------|------|
| F-1 (P1 required) | **closed（已实现生产源）** | 用户 P-004 裁决「实现生产源」（2026-08-13）：E-008 完成锁策略 + 423 + Host 终态全链路；D-002 §3 residual #1 关闭 |
| F-2 (P1 required) | fixed | `captureReturnIntent` 无参调用改为解析 `window.location.search`（`parseLocationQuery`）再走 `validateReturnIntent`；补无参生产路径测试（`return-intent.test.ts`）；commit `…` |
| F-3 (P2) | fixed | I-004 三处统一为「provider=`grok build`；S2=A-005/A-006；S6=A-007/A-008」；A-007 检查表更新为 A-008 已落盘 |
| F-4 (P2) | fixed | 03-audit 结论段吸收 A-007（self，conditional）与本 A-008（independent，conditional，BLOCKING=2）摘要 |

**用户裁决（F-1，P-004，2026-08-13）：** **实现生产源**（用户选择）。E-008 已完成：
migration 12（`failed_login_count` + `locked_until`）、5 失败/15 分钟锁窗、锁开撤销 refresh
token、423 `ACCOUNT_LOCKED`（error catalog 双语）、`AuthContext.locked`、`adapterAuthFor` 单一
来源、`HOST_ACCOUNT_LOCKED` 终态（home/support only）。Go internal 全绿、web vitest 875、
tsc 0、Playwright 7+1skip。F-1 **closed**；BLOCKING_COUNT=2 → 0。
