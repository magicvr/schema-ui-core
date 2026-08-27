---
status: active
created: 2026-08-26
updated: 2026-08-26
parent: GOAL-013-w13-api-web-security-audit
version: 0.1.0
---

# D-001 · W13 落位与范围裁决（2026-08-26，用户书面）

**背景**：2026-08-26 会话完成 api/web 全量审查（方法、覆盖与全部发现见 A-001 及其附件）。发现 P1×1 + P2×3 + P3 加固/健壮性若干 + 既有残余若干。用户指令开子目标承载治理上下文并修复。

**决策 1 · 落位**：新建 `GOAL-013-w13-api-web-security-audit`，挂 workspace-009-production-hardening 的 Root `GOAL-001-production-hardening` 下，作为 W13 波次。
- 已选（用户，结构化提问）：「workspace-009 · W13 波次（推荐）」。
- 未选方案及理由：
  - 其他现有工作区（如 2026-08-25 关门的 workspace-019）：默认不接新区，重开须用户另行确认；且本批发现横跨 auth/MFA/上传/前端/部署，不属单一已关门域。
  - 新建工作区 + VP：本批为持续安全程序内常规波次规模，不构成新的独立意图（结构选型判定树 → 同 Root 子目标）。

**决策 2 · 范围**：本波一次性修复全部发现。
- 已选（用户，结构化提问）：「全部发现一次修完」——P1 + P2 必修；P3 加固（API/Web）与部署/运维项全量纳入分母。
- 未选方案及理由：
  - 「只修 P1+P2」：分母过窄，P3 仍需后续波次承接；
  - 「P1+P2 必修、P3 登记待裁」：用户明确选择扩大分母，避免多轮波次往返成本。

**范围内处置类例外**（非"直接修"，须按三路径裁决留痕）：
- F-007（账号锁定可作定向 DoS 的策略权衡）、F-013（自助 scope TOCTOU，当前休眠）→ 实施时 fixed / accepted-residual / overruled 三选一并记录；
- F-020 的 HSTS 部分受 I-001（TLS 终结拓扑未知，non-blocking）约束 → 按可行范围实施，拓扑项文档化；
- localStorage refresh token 残余维持 tokens.ts 既有 D-002 书面接受，本波不重开。

**审计模式**（P-004 登记）：security 高影响门禁 → 关门前置 **independent** 审计（项目级默认执行路径：grok build · grok-4.6 · reasoning high · `/audit`），S6 执行；provider 失败不得静默降级、不得由编排器冒充 independent。S3/S5 中任何 residual / overruled 处置须用户逐条书面确认后方可闭合。
