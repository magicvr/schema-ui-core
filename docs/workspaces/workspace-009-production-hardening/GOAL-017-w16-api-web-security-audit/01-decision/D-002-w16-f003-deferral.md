---
id: D-002
title: "F-003 延期处置决策"
status: pending
created: 2026-08-30
updated: 2026-08-30
parent: GOAL-017-w16-api-web-security-audit
version: 1.0.0
---

# D-002 · F-003 延期处置决策

## 决策点

**F-003 (M-1, P2 recommended)**: Refresh token localStorage → httpOnly cookie

经 E-003 评估，F-003 需要实质性双端改造（API login/refresh/logout + Web 客户端逻辑），建议延期到后续波次或独立子目标完整实施。

## 背景

### 现状
- Refresh token 当前存储在 `localStorage`
- 存在 XSS 攻击下的 token 泄露风险（若 XSS 防护失效）
- 该权衡在 GOAL-005 D-002 中已知并记录

### 提议方案（E-003）
双模式架构：
1. Web SPA 默认使用 httpOnly cookie（防 XSS）
2. API 同时支持 cookie 和 `X-Refresh-Token` header（向后兼容）
3. 保留 localStorage 代码路径供客户端模式使用

### 实施工作量
- **API 端**（3 个端点改造）:
  - `/api/auth/login`: 成功后设置 httpOnly cookie + 返回 tokens
  - `/api/auth/refresh`: 优先读 cookie，无 cookie 时读 header
  - `/api/auth/logout`: 清除 cookie
- **Web 端**（客户端逻辑改写）:
  - 保留 localStorage 操作但默认不使用
  - 添加 cookie 模式检测与回退逻辑
- **测试需求**:
  - 单元测试：cookie 优先/header 回退逻辑
  - 集成测试：login → refresh → logout 完整流程
  - 回归测试：确保不破坏现有 localStorage 模式
- **文档**: 双模式架构说明与选择建议

### 当前波次状态
- ✅ Required (P1): F-001/F-002 已修复
- ✅ Recommended (P2): F-004/F-005 已由先前工作区处理
- 🔄 Recommended (P2): F-003 待处置
- ✅ Informational (P3): F-006/F-007/F-008/F-009 已评估

## 选项分析

### 选项 A: 延期到后续波次（建议）

**优点**:
- ✅ 允许 W16 快速关门（P1 required 已全部修复）
- ✅ 可在后续波次或独立子目标中完整设计、实施和测试
- ✅ 不阻塞其他安全修复工作的推进
- ✅ 保持波次聚焦与可验证性

**缺点**:
- ⚠️ 延长 localStorage 风险窗口（但当前 XSS 防护已有多层防御）

**所需决策**:
- 用户书面接受 **accepted-residual**:
  - 范围：F-003 localStorage 风险
  - 理由：需实质性双端改造，延期专注完整实施
  - 复审触发：W17+ 或独立子目标开设时
  - 责任人：安全加固程序（Root Goal）

**后续行动**:
- 在 Root Goal E-008 或后续波次登记 F-003 为待修复项
- 开设独立子目标或在 W17+ 中实施

---

### 选项 B: 当前波次实施

**优点**:
- ✅ 尽早消除 localStorage 风险

**缺点**:
- ⚠️ 需要 API+Web 双端完整改造（估计 1-2 天工作量）
- ⚠️ 增加复杂度（双模式架构需要完整测试）
- ⚠️ 可能引入新的回归风险
- ⚠️ 阻塞当前波次关门（P1 已完成，延长验证周期）

**风险**:
- 双模式实施不充分可能引入新漏洞
- 回归测试不完整可能破坏现有客户端

---

### 选项 C: 拒绝修复（不建议）

**理由为何不可行**:
- ❌ F-003 为 recommended 级别，且有清晰的安全改进价值
- ❌ localStorage 风险在 XSS 防护失效时真实存在
- ❌ 业界最佳实践是使用 httpOnly cookie

---

## 建议

**推荐选项 A**（延期到后续波次）：

1. **当前波次**:
   - 用户书面接受 F-003 为 **accepted-residual**（范围与复审条件明确）
   - 完成 W16 关门（P1 required 已修复，P2/P3 已处置）
   - 恢复 VP-008 go 宣称

2. **后续波次**:
   - 在 W17+ 或独立子目标中完整实施 F-003
   - 双模式架构完整设计、开发、测试
   - 向后兼容性完整验证

**理由**:
- 符合波次聚焦原则（P-001）：小目标可直接执行，大改造需路线图
- 允许 F-003 获得足够的设计、实施和测试时间
- 不阻塞其他安全工作推进
- 残余风险已知、有界、可追踪

---

## 决策记录

**待用户确认**:

请选择以下之一：

- [ ] **选项 A**（建议）: 接受 F-003 为 accepted-residual，延期到后续波次
  - 需书面确认：接受范围、复审触发条件
  - 落盘位置：本文件 § 用户裁决
  
- [ ] **选项 B**: 当前波次实施 F-003
  - 需确认：愿意接受额外实施周期（1-2 天）
  - 需补充：双模式详细设计与测试计划
  
- [ ] **其他决策**: 请明确指示

---

## 用户裁决

**决策日期**: 2026-08-30  
**选择**: **选项 A** — 延期到后续波次  
**理由**: 
- F-003 需要实质性双端改造（API + Web），估计 1-2 天工作量
- 当前波次 P1 required (F-001/F-002) 已全部修复
- P2 recommended (F-004/F-005) 已由先前工作区处理
- 回归测试全部通过（go vet/test, vitest, tsc）
- 延期允许 W16 快速关门，F-003 可在后续波次完整设计、开发、测试

**后续行动**:
1. 更新 A-001，将 F-003 标记为 **accepted-residual**
2. 登记残余风险：
   - **范围**: F-003 refresh token localStorage 风险（XSS 防护失效时可能泄露）
   - **复审触发**: W17+ 或独立子目标开设时实施双模式架构
   - **责任人**: 持续安全程序（GOAL-001-production-hardening）
3. 在 Root Goal 或 W17 中登记 F-003 为待修复项
4. 继续 W16 关门流程（S4 自审 → S6 独立审计）

---

## 影响

- **门禁**: I-004 (non-blocking) 需用户裁决后才可进入 S4 自审
- **下一步**: 
  - 若选项 A → 更新 A-001，F-003 标记为 accepted-residual，进入 S4
  - 若选项 B → 补充详细设计，开设 E-005 实施 F-003
