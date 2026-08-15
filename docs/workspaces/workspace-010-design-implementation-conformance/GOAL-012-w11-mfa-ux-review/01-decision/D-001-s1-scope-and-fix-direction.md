---
id: D-001-s1-scope-and-fix-direction
doc: decision-entry
goal: GOAL-012-w11-mfa-ux-review
date: 2026-08-15
status: accepted
parent: GOAL-001-design-implementation-conformance
created: 2026-08-15
updated: 2026-08-15
version: 0.1.0
---

# D-001 · S1 范围与修复方向裁决

## 背景

用户 2026-08-15 启动目标轮次指令：「推进目标，直到顺利关门」；会话审批策略为 never（免逐项审批）。按 P-004，把本目标的裁决点以书面形式记录并接受（用户轮次指令即书面授权）。

## 裁决

1. **分批范围与顺序**（对应 00-meta 路线图）：按文档既定顺序推进——
   - S2 先修 M-01～M-03（用户实测报告的 P0 缺陷）；
   - S3 实施 U-01/U-02（P0）；
   - S4 实施 U-03～U-07（P1）；
   - U-08～U-14（P2）不设硬性完成要求，S5 关门时按剩余资源选择性实施并留痕（不完成不阻断关门）。
2. **I-001 修复方向（accepted）**：
   - 自服务端点（/api/mfa/confirm、/api/mfa/disable、/api/mfa/recovery/rotate）的 ErrMFAInvalid 从 **401 改为 400**（业务校验失败语义），前端 authFetch 不再误判为会话丢失；
   - 登录二步验证 /api/auth/mfa/verify 保持 401（由 mfaVerify 直接解析，不经 authFetch，行为不变）；
   - 解绑成功后服务端吊销全部会话（A-004 F-002 安全语义不变），前端改为：设置 sessionStorage["mfa.disabledNotice"] 标记 → 本地 logout()（AuthContext 置 unauthenticated）→ 登录页消费标记显示「MFA 已解绑，请重新登录」横幅。
3. **I-003 二维码方案（accepted）**：引入 npm 依赖 **qrcode-generator**（MIT、零依赖、纯 JS、~10KB gzip），在 React 中以 SVG 渲染模块矩阵（不依赖 canvas，jsdom 可测、离线可用）。废弃「自绘 canvas」选项（工作量/正确性风险高）。

## 未选方案

- 错误码保持 401 + 前端特判 MFA 端点：脆弱（前端需维护端点清单，与 authFetch 通用契约冲突）。
- 解绑成功后服务端不吊销当前会话：削弱安全语义，违背 A-004 F-002 parity。
- 自绘 QR 编码器：Reed-Solomon/掩码实现复杂、正确性风险高，测试成本大。
