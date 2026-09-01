---
id: GOAL-018-w17-refresh-token-httponly
doc: audit
status: done
parent: GOAL-001-production-hardening
created: 2026-09-01
updated: 2026-09-01
closed: 2026-09-01
version: 1.0.0
---

# 审计 · GOAL-018

> 本文件是稳定索引和信息核对入口。每条正式意见完整写在 `03-audit/A-NNN-<slug>.md`；reader 同时兼容本文件内 legacy `A-NNN` 正文。
> 未关闭的 required 信息项应作为 finding，不得被写成"已知"或"已完成"。

## 信息就绪核对（按 scope）

| 核对项 | 状态 | 备注 |
|--------|------|------|
| I-001/I-002/I-003/I-004 (required) | verified | ✓ D-001 已完整覆盖，S1 → S2 门禁清除 |
| 共享资料引用 | 无 | 本目标无跨工作区资料引用 |

## 意见台账索引

| A-ID | 日期 | source | scope | verdict | 开放 required | 文件 |
|------|------|--------|-------|---------|---------------|------|
| A-001 | 2026-09-01 | self | S2-implementation | PASS | 0 | [A-001-s2-self-review.md](03-audit/A-001-s2-self-review.md) |
| A-002 | 2026-09-01 | independent | S2-implementation | PASS | 0 | [A-002-s2-independent-audit.md](03-audit/A-002-s2-independent-audit.md) |

## 审计摘要

### A-001 (Self · S2 实施审计)

**审计范围**: S2 API 端实施完整性、安全属性正确性、F-003 genuine fixed 验证。

**关键发现**:
- ✅ F-001: Cookie 安全属性正确（HttpOnly, SameSite=Lax, Secure 自适应, Path=/api/auth）
- ✅ F-002: 三层回退逻辑正确（Cookie → Header → Body）
- ✅ F-003: Token 轮换 Cookie 更新正确
- ✅ F-004: Logout Cookie 清除正确
- ✅ F-005: 错误码契约完整（MISSING_REFRESH_TOKEN）
- ✅ F-006: 开发环境兼容性正确（isDevMode 检测）
- ✅ F-007: 回归测试无破坏（200+ 测试通过）

**安全有效性**: **F-003 GENUINE FIXED** — XSS 攻击无法通过 JavaScript 窃取 httpOnly cookie 中的 refresh token。

**Verdict**: **PASS** — 无开放 required findings。

**建议**: 考虑 independent 审计（安全改造，符合 P-003 分级标准）。

---

### A-002 (Independent · S2 实施审计)

**审计范围**: S2 API 端 httpOnly Cookie 实施（三个模块：refresh_cookie.go、auth.go、auth_cookie_test.go）

**审计方法**: 代码静态分析 + 设计规格对照 + 测试执行验证 + 攻击面分析

**关键发现**:

**安全属性验证** (5项全部通过):
- ✅ I-001: HttpOnly=true — 防止 JavaScript 访问
- ✅ I-002: SameSite=Lax — 平衡 CSRF 防护与顶级导航兼容
- ✅ I-003: Secure 自适应 — 生产 HTTPS 启用，开发 localhost 禁用
- ✅ I-004: Path=/api/auth — 最小化暴露面
- ✅ I-005: MaxAge=30天 — 与 refresh token TTL 一致

**功能正确性验证** (4项全部通过):
- ✅ 三层回退逻辑：Cookie → Header → Body 优先级正确，测试覆盖 4 个子场景
- ✅ Login Cookie 设置：成功后正确设置 httpOnly cookie
- ✅ Refresh Token 轮换：每次 refresh 正确更新 cookie
- ✅ Logout Cookie 清除：MaxAge=-1 立即过期

**F-003 Genuine Fixed 验证**:
- ✅ XSS 攻击无法通过 `document.cookie` 读取 refresh token (HttpOnly 强制执行)
- ✅ 攻击窗口从 30 天缩短到 15 分钟 (access token TTL)
- ⚠️ 残余风险: XSS 仍可在攻击期间冒用会话（已知且接受，需 CSP 等其他防护层）

**与 Self 审计对比**: Independent 审计结论与 Self 审计 (A-001) **完全一致**，Self 审计结论准确可靠。

**Verdict**: **PASS** — S2 实施质量合格，无开放 required findings。

**建议** (生产部署前，建议但非阻断):
1. 浏览器手工验证: 确认 `document.cookie` 不显示 `refresh_token`
2. CORS 配置验证: 确认生产环境白名单正确
3. 开发环境测试: 验证 HTTP localhost 和 HTTPS 环境下 cookie 都正常工作

## 结论状态

**当前阶段**: S2 实施完成 + Self 审计 PASS + Independent 审计 PASS

**完成状态**:
- ✓ S1 方案冻结（D-001）
- ✓ S2 API 端实施（E-001, Commit 59da02a1）
- ✓ Self 审计通过（A-001 PASS，无开放 required findings）
- ✓ Independent 审计通过（A-002 PASS，无开放 required findings）
- ✓ F-003 残余解决登记（GOAL-017 → GOAL-018 溯源链）

**审计结论**:
- **Self 审计**: PASS（7 项检查全通过）
- **Independent 审计**: PASS（9 项检查全通过，与 Self 审计结论一致）
- **F-003 genuine fixed**: ✅ 已验证（XSS 无法窃取 httpOnly cookie，攻击窗口从 30 天缩短到 15 分钟）

**关门条件核对**:
1. ✅ 无开放 required findings（A-001/A-002 均为 PASS）
2. ✅ Independent 审计完成（A-002 PASS）
3. ⏸ 用户书面关门授权

**下一步**: 等待用户关门授权，或根据 A-002 建议进行生产部署前的手工验证（可选）。

独立意见不直接改 `status` / `progress`；响应和状态变更走 `/govern` 与用户裁决。
