---
id: D-001
goal: GOAL-007-w7-api-web-security-audit
title: 开 W7：落盘独立审计、修复范围待确认
date: 2026-08-19
status: accepted
parent: GOAL-001-production-hardening
created: 2026-08-19
updated: 2026-08-19
version: 0.1.0
---

# D-001 · 开 W7：落盘独立审计、修复范围待确认

### 触发

用户 2026-08-19 `/govern`：在 workspace-009 新建子目标，落盘本会话对 `apps/api` + `apps/web` 的独立审计意见（用户明确：独立审计、不加载 skills）。

W6（GOAL-006）已关门；W5 为 0 中高危扫描未开子目标。本波有 **2 条 high required**，按程序约定必须开子目标承接。

### 决定

1. **开 W7 子目标** `GOAL-007-w7-api-web-security-audit`，`parent` = `GOAL-001-production-hardening`。Root / VP-009 保持 `active`，不因本波立项或完成而关门。
2. **本回合只落盘独立意见**（A-001，`source: independent`）。**不**改业务代码、**不**闭合 finding、**不**把沉默当 residual/overruled。
3. **finding 清单权威** = A-001。required 项（F-001～F-012）阻断本波实施放行与关门，直至三路径闭合。S2 起须用户确认修复范围（可整单采纳、可单条 residual/overruled）。
4. **审计模式**：本波含 security/auth/IDOR 高影响，后续实施与关门默认 `cross`（self + independent）。本条 A-001 已满足 independent 落盘；实施后仍须 self，关门前再跑 independent。
5. **VP-008 `go`**：High findings 是否暂挂消费有效性 **不自动裁决**（I-002）；S2 时询问。本回合不宣称 go 失效或仍有效。
6. **不重开** 已接受残余：localStorage refresh、匿名 schema/manifest、Compose 无 TLS、bcrypt cost、data-permission v1 未接线、development JWT 脚枪、单会话吊销不 bump `token_version`。

### 为什么

- VP-009 / Root 程序语义：扫描发现中高危 → 有界波次子目标；波次可关门，程序容器不关。
- P-003：独立意见必须进被审目标 `03-audit/A-NNN` + 索引，聊天不算放行依据。
- P-004：修复范围、residual、go 暂挂均须用户书面确认。

### 未选方案

- 把意见只写在聊天或 `raw/`：违反 P-003。
- 并入已关门的 GOAL-006：编号不复用新含义；W6 范围是低危扫描修补。
- 本回合直接开工修复：用户指令是落盘意见，不是实施。
- 因 High 自动暂挂 VP-008 `go`：无用户书面裁决。
