---
id: GOAL-004-r3-binding-flow
doc: audit
status: active
parent: GOAL-001-account-email-identity
created: 2026-08-24
updated: 2026-08-24
version: 0.2.0
---

# 审计 · GOAL-004（R3 绑定/校验流）

> 本文件是稳定索引。正式意见写在 `03-audit/A-NNN-*.md`。

## 信息就绪核对（按 scope）

| 核对项 | 状态 | 备注 |
|--------|------|------|
| 影响本 scope 的 I-00N | I-005 **verified**（GOAL-004 D-001：TTL 10 分钟 / 冷却 60 秒）；I-006 **verified**（允许代填→pending）。Root `00-meta` 仍 collecting（台账滞后，见 A-001 F-003） | 最晚阶段 = R3 方案/实施；用户书面裁决 2026-08-24 |
| 到期 required 是否已 verified / residual | I-005 有书面裁决；落地见 A-001（Resend 冷却可核对；Bind 再签发见 F-002 recommended） | I-006 正向 HTTP 面未接通（A-001 F-001 required） |
| 资料引用是否固定且用户确认 | 不适用 | `shared_materials_catalog: none` |

## 意见台账索引

| A-ID | 日期 | source | scope | verdict | 开放 required | 文件 |
|------|------|--------|-------|---------|---------------|------|
| A-001 | 2026-08-24 | independent | R3 关门（0055 · 绑定/校验流 · 错误契约 · 最小页面；checkpoint `0ae17f09`） | conditional | 1（F-001） | [A-001-independent-r3-binding-flow-closeout.md](03-audit/A-001-independent-r3-binding-flow-closeout.md) |

## 结论状态

A-001 independent **`conditional`**（开放 required = 1 · F-001）。意见已落盘；**不**改 GOAL-004 `status` / 检查点 / `progress`，**不**改 goal-tree。响应与关门由 `/govern` 处理。

**编排器响应（/govern · 2026-08-24 · E-003）**：
- **F-001 → fixed**（commit `bd1cdff9`）：`email` 入 users 资源 `RawStringFields`，`""` 清空可达、非字符串 400；新增 HTTP 全链路测试 TestUsersPatchEmailPrefillFlows。authsession/handler/composition 复跑全绿。
- **F-002 → fixed**：同址 pending 重绑套用 60 秒冷却（不同地址换绑立即派发）；服务测试覆盖。
- **F-003 → fixed**：goal-tree 补登 GOAL-004；执行索引补 E-002/E-003；Root `00-meta` 信息表 I-005/I-006 → verified（权威表与镜像同步）。
- **F-004 → 口径对齐**：D-001 v1.1.0 §5 澄清 + 新增 §6。
- **N-1 / N-2 → 维持归属**（N-1 归 R4 证据面说明；配对不变量仓储层落点随本关门留痕）。

开放 required 归零。**GOAL-004 关门：done · 4/4。**
