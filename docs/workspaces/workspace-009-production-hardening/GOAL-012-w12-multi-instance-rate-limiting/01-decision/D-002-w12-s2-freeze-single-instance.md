---
id: GOAL-012-w12-multi-instance-rate-limiting
doc: decision-entry
record_id: D-002
status: accepted
goal: GOAL-012-w12-multi-instance-rate-limiting
created: 2026-08-26
updated: 2026-08-26
version: 1.0.0
---

# D-002 · S2 方案冻结：维持单实例官方边界 · 载体预登记 Redis 方向 · 零码收官

> 用户裁决留痕：2026-08-26 会话内 `ask_user_question` 三问三答（I-001 / I-002 / 处置），本条目为完整书面记录。

## §1 · I-001 裁决（required · verified）：维持单实例官方边界

**决定**：产品部署拓扑意图维持「单实例为官方支持形态」，不升级多实例宣称。认证限流（login/recovery `loginRateLimiter`）按节点预算分摊 = **已文档化部署边界**，非缺陷、非任务残余；复审触发 =「多实例部署形态出现」（对齐 GOAL-033 A-002 R-A6-1 先例口径）。

**证据基础（四处既有书面边界 + 程序先例）**：

| 出处 | 表述 |
|------|------|
| `README.md` L86 | 「完整生产运维 / CI-CD 部署流水线、TLS、多实例为**非目标**」 |
| `compose.yaml` 头注 | "Non-goal: … TLS termination, multi-instance scaling" |
| workspace-002 I-008-001 工程契约（Root D-013 边界） | 多实例水平扩展 = 非目标 |
| vision roadmap RT-Q05 / RT-D04 | 限流跨实例与多实例均 **trigger-gated**（「单实例够用；多实例才需要共享存储」） |
| 波次先例 | W3 D-001 明确不做「多实例限流后端」；GOAL-002/W3 00-meta residual 接受 |

PG 方言（VP-013）为生产权威存储 ≠ 多实例宣称——该 VP 非目标清单明确「A3 多实例不进退出分母」。

## §2 · I-002 裁决（required · verified）：预登记载体方向 = Redis 等进程外依赖

**决定**：若未来触发多实例（A3），共享限流状态载体的预登记技术方向为 **Redis 等进程外专用存储**（用户裁决；编排器建议为内核 Store 新表，未被采纳——两案论据并录如下）。

| 方案 | 论据 | 论据 |
|------|------|------|
| **Redis 等进程外依赖（采纳）** | 限流为高频计数热路径，TTL 键 + 原子 INCR 是专用语义，避免主存储争用；业界标准、跨语言生态成熟 | 引入新基础设施依赖，偏离最小部署哲学（VP-015 曾为此选无收集器默认） |
| 内核 Store 新表（未采纳） | 双方言迁移机制成熟零新依赖；恢复码 CAS / 渠道表先例；对齐 RT-Q05 倾向 | pre-auth 热路径每次查 DB，SQLite `MaxOpenConns=1` 下争用敏感；滑动窗口在 SQL 侧表达较繁 |

**边界澄清**：预登记 ≠ 本波实施承诺，也 ≠ 触发时的最终冻结——A3 真正立项时仍须走正式方案冻结（届时连同连接池/分布式 ID/优雅停机等 RT-D04 联动项一并裁决）。roadmap RT-Q05 保持 trigger-gated 登记不变。

## §3 · I-003 裁决（non-blocking · closed）：单实例边界下语义保持现状

login/recovery 两桶独立键空间与预算（15 min / 20 次 / `IP|identifier`）、`Retry-After` 语义、窗口常量均维持现状不动；多实例下的统一语义问题随 A3 触发归入 §2 的届时冻结范围。

## §4 · S3 缩减与收官方式

- **S3 缩减为零代码变更**（评估型波次的合法收官形态）：单实例边界已在 README / compose / 代码注释（rate_limit.go L12–16）三处如实声明，追加文档属冗余；roadmap 登记已存在且与本波结论一致，无需改动。
- **审计模式确定（P-004）**：本波无任何代码/行为变更，security 高影响不成立 → S4 采用 **self 复核**（不强制 cross grok 腿）。
- 用户已选「零代码变更直接复核关门」= 书面关门授权；S4 self pass 后 GOAL-012 即 `done`。

## 未选方案

- **升级为支持多实例**：推翻四处非目标宣言 + RT-D04 全家桶联动，远超有界波次，应经 `/vision` 重开架构分支。未选。
- **单实例边界下限流器先行上 Store**：收益在可预见拓扑内不兑现，违反 trigger-gated 哲学。未选。
- **载体预登记为 Store 新表**：见 §2 对照，用户裁 Redis 方向。未选（论据保留供 A3 触发时复用）。
- **README 补部署注意句**：边界已三处声明，冗余。未选。
