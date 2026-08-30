---
title: A-004 · W15 审计意见响应与闭合记录（S6）
status: active
created: 2026-08-30
updated: 2026-08-30
parent: null
version: 0.1.0
---

# A-004 · W15 审计意见响应与闭合记录（S6）

日期：2026-08-30 · 编排器合并响应（P-003）：响应 A-002（self pass）+ A-003（independent pass）+ A-001 分母；全部相关意见 no-conflict、同向。

## 响应摘要

| 来源 | verdict | 本响应 |
|------|---------|--------|
| A-001 independent intake | conditional（6 required + 1 recommended） | 分母全部响应，见下 |
| A-002 self | pass | 采纳；两条 note（S-001/S-002）= 有据记录，无需动作 |
| A-003 grok-build independent | pass（0 required；recommended F-008/F-009 + notes N-001～N-003） | F-008/F-009 立即实施；notes 全部留痕/修复，见下 |

## A-001 分母闭合（P-003 三路径：fixed）

| Finding | 级别 | 闭合路径 | 证据 |
|---------|------|----------|------|
| F-001 | P1 required | **fixed** | E-001；代码默认/内嵌 YAML/create 模板回环 + 空 env fail-closed；负例 `TestLoadConfigRequiresExplicitAppEnv`/`TestLoadConfigDefaults`；A-003 §F-001 逐条表 |
| F-002 | P1 required | **fixed** | E-001；`ValidateJWTSecretStrength` 单一来源 + server 复用；短/纯字母/纯数字负例；A-003 §F-002 |
| F-003 | P1 required | **fixed** | E-001；`ValidateSeedPassword`（8–72 非空）+ `resolveSeedHash`/`bootstrapAdmin` 非 dev fail-closed，dev 回退保留；启动负例 ×2；A-003 §F-003 |
| F-004 | P2 required | **fixed** | E-001；`requireActiveSecondFactor` CAS `AdvanceLastUsedStep`；重放拒绝测试；A-003 §F-004 |
| F-005 | P2 required | **fixed** | E-002；`invite-accept.tsx` 读 token 即 `replaceState` 清理；本波新增回归锁 F-008 |
| F-006 | P1 required | **fixed** | E-002；13 suite fixture 根统一 + guard 测试 + README；vitest 1183/1183（基线 76 失败）；A-003 §F-006 |
| F-007 | P3 recommended（用户裁决 = fixed） | **fixed** | E-003；LocalStore `0700/0600` + 权限测试（非 Windows）；A-003 §F-007 |

**闭合后开放 required = 0。**

## A-003 新发现响应

- **F-008（recommended）→ fixed**：新增 `apps/web/src/components/invite-accept.test.tsx`（jsdom ×3：token 清理且保留其它 query / 无 token 不触碰 history / 二次挂载幂等）。回归：`vitest src/components/invite-accept.test.tsx` 3/3 pass；全量 vitest 复跑见 E-004。
- **F-009（recommended）→ fixed**：`server/config_test.go` 中 `TestLoadConfigInvalidShutdownTimeout`（`0s`/`-1s`）与 `TestLoadConfigDialectPairing` 的自定义 YAML 显式声明 `app.env: development`，确保用例真正咬到目标分支而非空 APP_ENV 门禁（F-001 后假绿修正）。回归：`go test ./server/ -count=1` pass。
- **N-001（note）**：F-007 权限断言在 Windows 按设计 SKIP；已记录于 A-003/A-002；Linux/darwin CI 复跑为后续选项，不构成实现缺口。
- **N-002（note）**：`server/config.go` `AppEnv` 注释已同步现行语义（空串 = 拒绝，不再是 development 缺省）。`03-audit.md` 信息表在本波后续已更新为 verified/done 态。
- **N-003（note）**：`TestShutdownDrainHarnessPostgres`（VP-021 harness，非 W15 范围）在重型并行负载下曾 flake；隔离与第二次全量均绿。留痕，不升格；建议后续波次错峰运行全量 API 与 Web 测试。

## 门禁状态

- 开放 required findings：**0**（全部 fixed，证据链可核对）。
- 到期 required 信息项：无（I-001～I-005 均 verified，A-003 复核）。
- 意见冲突：无（self 与 independent 同向 pass）。
- 关门剩余动作：S6 审计腿完成；**用户书面授权后**执行 D-003（`status: done` + goal-tree 同步）。