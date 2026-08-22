---
id: D-001
doc: decision-entry
goal: GOAL-001-key-rotation-and-backup
status: accepted
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
---

# D-001 · 开区 scaffold 与 A5 纲领路线图

## 背景

用户确认：对 VP-016 做意图审视；没有问题的话交 `/govern` 开区。未另写工作区 slug。VRev-035 self `pass`（0 required）。A5 退出分母已由 VP-016 v0.1.0 / VR-037 冻结，激活时 editorial 收口退出 1。

工作区 slug 按 VP-013 / VP-014 / VP-015 惯例取 `workspace-016-key-rotation-and-backup`（`workspace-NNN` + VP slug）。Root slug 取 `GOAL-001-key-rotation-and-backup`，与 VP / 工作区后缀对齐。用户确认的是「开区」动作；slug 为惯例推导并在此留痕，不是另一次口头点名。

## 决策

1. 激活 `VP-016-key-rotation-and-backup`（v0.2.0 `planned → active`）。
2. lead 工作区 slug = `workspace-016-key-rotation-and-backup`；Root = `GOAL-001-key-rotation-and-backup`。
3. 纲领路线图 R1～R5：轮换合同冻结 → JWT 双密钥实现 → 轮换后恢复证据 → 默认单密钥仍可用 → 显式双密钥轮换路径 **与** 恢复路径证据。
4. 配置面：缺省单 `AUTH_JWT_SECRET`；previous 为显式配置的生产/验收路径。不改 Compose 默认依赖。重启生效，热加载不进本波。
5. 开区审计模式 **none**（可逆文档 scaffold）。R1 合同冻结起按内核门禁走 **self**；密钥轮换生产路径实施按 **independent**（项目默认 grok build · grok-4.6 · `/audit`）。
6. 本回合**不**创建 R1 子目标、**不**改 `apps/api` 代码。

## 架构类 freshness（V-F068）

VP-008 强制 freshness 的对象是后续**业务** VP。本 VP 是架构 A5，按自身激活门闩做轻量复核：

| 项 | 值 |
|----|-----|
| 原 `go` 候选 | `ed99e88`（2026-08-10） |
| 现行 HEAD | `57098c3` |
| VP-009 / VP-010 | 无开放中高危暂挂宣称 |
| Vision open required | 0 |
| F-007 | 上传授权深度仍 deferred；本 VP 不扩张授权 scope |
| 是否消费业务解锁 | **否**。不改 Profile / 模块矩阵 / Manifest 为意图 |
| 结果 | **PASS**。不暂挂 `go` |

`consumer_vp` = VP-016；`last_freshness_review_at` = 2026-08-22；`next_freshness_review_trigger` = 若实施改变共同门禁 / Profile / 模块矩阵则重做。

## 为什么

- 新纲领波次独立工作区，避免写入已 closed 的 workspace-015。
- slug 与 VP id 对齐，便于组合索引。
- I-001～I-005 满足 V-F067；freshness 表满足 V-F068。
- 开区时已见：`ParseAccessToken` 单密钥；服务凭证 SHA-256 opaque。后者只作 I-002 输入，R1 才书面出局。

## 未选方案

- 继续 `planned` 只写 VP：用户已要求通过后开区。
- 重开 workspace-015：VP-015 已关门且默认不接新区。
- 一开区就改 JWT 解析：R1 合同未冻结。
- 默认改为必须有 previous：与 A5 内嵌默认冲突。
- 把本 VP 当业务 VP 走 VP-011 式完整 freshness 矩阵：解锁 scope 不匹配。
- 等待用户另写 slug 再开区：用户本轮已指令「开区」；惯例与 VP-013/014/015 同构，记录推导以免静默。
