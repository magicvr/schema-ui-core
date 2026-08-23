---
id: A-002-w22-security-independent
doc: audit-entry
goal: GOAL-033-w22-residual-closeout
source: independent
date: 2026-08-23
scope: 安全面整改复核：A5 上传内容嗅探（upload.go）+ A6 MFA verify 独立限流（mfa.go）；连带覆盖关门叙事合法性
audit_type: finding-closure
verdict: pass
---

# A-002 · 安全面整改独立复核

## 范围与区间

`git diff` 全量核对 `apps/api/internal/handler/upload.go`（A5）与 `mfa.go`（A6）；对照 W9·GOAL-002 N-002 与 W9·GOAL-011 R-001 的原始残余条款；旁及 E-005/E-006 关门叙事。

## 成果（有证据）

- **A5**：8 KiB 头部嗅探窗口 + `<svg`/`<script`/`<?xml` 标记 + `(?i)\bon[a-z]+=` 事件处理器启发式；拒绝在 MIME 之上、下载侧 attachment/CSP sandbox/nosniff 不变——与 N-002 残余「安全边界=下载头」的既定模型一致且收紧。测试含端到端（`TestUploadA5ContentSniffEndToEnd`）。
- **A6**：独立桶 `newLoginRateLimiter(15m, 10, 1<<16)`，不复用登录桶；键=client IP（复用登录既有代理信任规则）；失败才计数、成功不占预算；429 + Retry-After + 本地化 RATE_LIMITED。3 条限流测试 PASS。

## Findings

| 编号 | 级别 | 描述 |
|------|------|------|
| R-A5-1 | recommended · 记录 | `on[a-z]+=` 对头 8 KiB 含独立 `on*=` 词形的良性文本存在理论误报面；拒绝属安全侧偏移，且有下载头兜底，可接受——建议后续若出现误报反馈再收窄为 `on(click|error|load|mouseover|focus)=` 枚举。 |
| R-A6-1 | recommended · 记录 | 限流器为进程内存态，重启清零（与登录限流器同姿态，parity 可接受）；多实例部署时桶不共享——当前单写者 SQLite 部署形态下无影响。 |

**required：0。**

## 结论 + 建议

verdict **pass**。两条 recommended 由编排器记录在案即可，无需代码改动。关门放行待用户书面确认（编排器下一步）。

### 声明

本意见不修改 status/progress；响应由 /govern 处理。
