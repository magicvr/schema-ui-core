---
id: D-002-w9-finding-inventory
doc: decision-entry
goal: GOAL-009-w9-api-web-security-audit
record_id: D-002
status: accepted
parent: GOAL-001-production-hardening
created: 2026-08-21
updated: 2026-08-21
version: 0.1.0
---

# D-002 · W9 finding 清单调和（A-002 F-001 响应）

### 触发

用户 2026-08-21 `/govern` 书面指令：响应 [A-002](../03-audit/A-002-w9-a001-reasonableness.md)，先调和 A-001 finding 清单（作废空洞 F-003、给全文 P2-7 编号、修正 22/12/11 计数），再裁决 I-002。本条只完成**前半**：冻结消费用清单。I-002（本波修哪些 required、是否暂挂 VP-008 go）**仍 open**，不得本条静默裁。

### 决定

1. **S2/S3 消费权威**为本表，而不是 A-001 原文的「F-001～F-012」或 `03-audit.md` 曾写的「22」。A-001 原文保留为独立意见历史；消费指针见 A-001 文首勘误。
2. **F-003 作废**：A-001 从未定义该条（无标题、无正文）。编号**不复用、不改挂**全文 P2-7 或其他缺陷（AGENTS：禁止把已取消编号赋予新含义）。
3. **全文 P2-7 编号为 F-025**（cron DOM/DOW 用 AND，非 POSIX「两者受限取 OR」）。严重度 med、建议 required（与 A-001 全文 P2 其余 9 条同级进入清单；本波是否实施仍由 I-002 裁）。
4. **计数**：required **12** = 2 high（F-001、F-002）+ 10 med（F-004～F-012 与 F-025，其中无 F-003）。推荐 F-013～F-023（11）；info F-024（1）；作废 F-003（1）。索引「22」为抄写错误，作废。
5. **I-001** 以本表为 verified 证据。不勾选 S2、不改 I-002。

### 调和后 required 清单（12 · 均 open · 本波是否修待 I-002）

| 消费 ID | 来源 | 严重度 | A-002 独立判定 | 摘要 |
|---------|------|--------|----------------|------|
| F-001 | A-001 | high | 成立，保留 high required | 钱包 `isUniqueViolation` 仅 SQLite 文案；PG 幂等/去重契约失效；无资金错账 |
| F-002 | A-001 | high | 成立，保留 high required | 生产 nginx 未代理 host-bootstrap → SPA 200 HTML → boot `HOST_PROTOCOL_REJECTED`。容器仍可 healthy |
| F-004 | A-001 | med | 代码成立；「助长暴力破解」过述 | `RecordLoginFailure` 非原子；有登录限流/captcha |
| F-005 | A-001 | med | 成立 | TOTP 跨事务 check-then-act；proof 非原子消费 |
| F-006 | A-001 | med | 成立 | 恢复码整表回写丢失更新/双花 |
| F-007 | A-001 | med | 代码成立；属可用性 | `internal/jobs/runner.go` handler goroutine 与 scheduler 无 recover |
| F-008 | A-001 | med | 部分成立 | 声明 `permissionIntent` 却漏 `key` 时 UI 门禁跳过；未标记默认 allow 是协议 |
| F-009 | A-001 | med | 运行时成立 | cascade 缺 source 时 `return true`；L2 仅测试调用 |
| F-010 | A-001 | med | 潜伏成立 | `delete()` Get 出错跳过归属预检；无生产 Scoper |
| F-011 | A-001 | med | 成立 | 凭据重名匹配 SQLite 文案，PG 约束名不命中 → 500 |
| F-012 | A-001（全文 P2-1；A-001 编号错位） | med | 成立 | scheduledtasks `LIKE OR … AND` 未加括号，绕过 enabled/status 过滤 |
| F-025 | A-001 全文 P2-7（原无 F-ID） | med | 成立 | `cron.go:99-107` DOM 与 DOW 用 AND；`0 0 1 * 1` 只在 1 号逢周一触发 |

**F-025 定义（补编号，非新审计主张）**：`apps/api/internal/modules/scheduledtasks/store/cron.go` `Matches` 要求 DOM **与** DOW 均命中。Vixie/POSIX 在两者均非 `*` 时取 OR。后果是少调度，不是多跑或越权。

### 不作废、不改号

- F-012 保持 scheduledtasks WHERE（不把 F-012 改回 F-003）。
- F-013～F-024 保持 A-001 原分级（recommended / info）。
- A-001 代码 findings 的闭合状态仍为 **open**。本条只调和身份与计数。

### 为什么

- A-002 F-001（required）阻断把自相矛盾清单当 S2 输入；用户书面要求先调和。
- 作废 F-003 而不把 P2-7 塞进 F-003，避免「空号改挂新含义」。
- 给 P2-7 新号 F-025，使全文 10 条 P2 与 10 条 med required 一一对应（2 P1 + 10 P2 = 12）。

### 未选方案

- **改写 A-001 正文、把 F-012 改号为 F-003**：破坏独立意见历史；用户要求作废空洞 F-003，不是重编号已有条。
- **把 P2-7 挂到 F-003**：赋予空号新含义，违规。
- **本条同时裁 I-002 整单采纳/最少集/go 暂挂**：用户说「再裁决」，但未写明范围；P-004 禁止静默自动裁。I-002 另条决策。
- **把 I-001 改回 collecting**：清单已可枚举且与全文 P1/P2 对齐；状态保持 verified，证据改为「A-001 + A-002 + 本表」。

### 影响

- I-001 证据更新；S2 仍等 I-002。
- A-002 F-001 → `fixed`（见 A-003）。A-002 F-002～F-004 recommended 转入 I-002 选项说明，不阻断清单调和。
- 不实施代码、不暂挂/恢复 VP-008 go。

### 后续

用户书面选择 I-002 后另写 D-003（范围 + go 宣称），再勾选 S2、开工 S3。
