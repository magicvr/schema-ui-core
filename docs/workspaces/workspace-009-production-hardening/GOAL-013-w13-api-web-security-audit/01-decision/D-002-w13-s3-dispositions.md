---
status: active
created: 2026-08-26
updated: 2026-08-26
parent: GOAL-013-w13-api-web-security-audit
version: 0.1.0
---

# D-002 · W13 S3/S5 处置裁决（2026-08-26，用户书面）

**背景**：S2 完成后进入 S3（API P3 与健壮性批）。按 D-001，F-007 / F-013 属"三路径裁决留痕"例外项；实施中另发现 F-020 的 img-src 部分与受支持功能冲突。三项经结构化提问由用户书面裁决（2026-08-26）。

## 决策 1 · F-007（账号锁定可作定向 DoS）→ **fixed**（承载于子目标）

- **用户选择**：「新建一个当前目标的下级子目标，承载治理上下文，fixed。」
- **落位**：新建 `GOAL-014-w13-account-lockout-redesign`（parent = GOAL-013），承载锁定模型重设计的完整治理上下文（方案、实施、审计、关门）。
- **方向**：现模型（5 连败 → 全局账号锁 15min + 吊销全部刷新令牌）允许知道用户名的攻击者反复锁定任意账号；重设计为不依赖"全局账号锁"作为唯一爆破防线的模型（如按 IP+账号联合维度计数/指数退避），保留管理员可见性（OnLockOpened）。
- **对 GOAL-013 关门的影响**：F-007 为 recommended（非 required），不构成 A-001 required 门禁；GOAL-013 关门检查时该 finding 以「fixed · 实施承载于 GOAL-014」登记。GOAL-014 未完成前 GOAL-013 是否先行 `done` 由关门审计（S6）时再请用户确认，不在本决策内静默裁定。
- **未选方案**：accepted-residual / user-overruled —— 用户明确要求修复。

## 决策 2 · F-013（自助行级 scope TOCTOU，当前休眠）→ **accepted-residual**

- **用户选择**：「accepted-residual（推荐）」。
- **残余内容**：`handler/resources.go` update/delete 的 self-scope 检查为 Go 侧先读后写、UPDATE 无属主谓词，存在理论竞争窗口；当前生产未启用任何 self-scope 行级角色（休眠）。
- **范围**：仅限 scope=self 的资源写路径；scope 为空/admin 的路径无此问题。
- **复审触发（硬性）**：首个使用行级 self-scope 的生产角色上线之前，必须完成谓词化条件 UPDATE 改造并经独立审计；该触发同时登记于 Root 执行台账。

## 决策 3 · F-020 img-src 部分 → **保留 https:（只实施 HSTS）**

- **用户选择**：「保留 https:+只做 HSTS（推荐）」。
- **理由（审查证据）**：`settings/repository.normalizeLogoURL` 本就接受 http(s) 绝对 URL 作为品牌图地址——img-src https: 是受支持功能的必要条件，不是配置疏漏。
- **已实施**：nginx 增加 HSTS（RFC 6797：HTTP 响应上的头被浏览器忽略，TLS 接入后自动生效；I-001 TLS 终结拓扑仍属运维侧）；CSP 其余部分不动。
