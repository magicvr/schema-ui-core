---
id: D-001
doc: decision-entry
goal: GOAL-001-object-storage
status: accepted
created: 2026-08-21
updated: 2026-08-21
version: 1.0.0
---

# D-001 · 开区 scaffold 与 A2 纲领路线图

## 背景

用户确认：对 VP-014 做 self Review；通过后激活并开区；slug = `workspace-014-object-storage`。VRev-031 self `pass`（0 required）。A2 退出分母已由 VP-014 v0.1.0 / VR-031 冻结。

Root slug 取 `GOAL-001-object-storage`，与 VP / 工作区后缀对齐（同 workspace-013 惯例）。用户确认的是工作区 slug；Root 未另给名。

## 决策

1. 激活 `VP-014-object-storage`（v0.2.0 `planned → active`）。
2. lead 工作区 slug = `workspace-014-object-storage`；Root = `GOAL-001-object-storage`。
3. 纲领路线图 R1～R5：端口冻结 → S3 接入 → 三类落盘收口 → 公共面去本地路径 → 双路径证据。
4. 配置面：缺省仍为本地盘（`filepath.Dir(db.path)` 下 avatars / brand-assets / uploads）；S3 兼容为显式配置的生产/验收路径。不改 Compose 默认依赖。
5. 开区审计模式 **none**（可逆文档 scaffold）。R1 端口方案冻结起按内核门禁走 **self**；S3 / 生产路径实施按 **independent**（项目默认 grok build · grok-4.6 · `/audit`）。
6. 本回合**不**创建 R1 子目标、**不**改 `apps/api` 代码。

## 架构类 freshness（V-F062）

VP-008 强制 freshness 的对象是后续**业务** VP。本 VP 是架构 A2，按自身激活门闩做轻量复核：

| 项 | 值 |
|----|-----|
| 原 `go` 候选 | `ed99e88`（2026-08-10） |
| 现行 HEAD | `83526a4` |
| VP-009 / VP-010 | 无开放中高危；最近 W10 已恢复 `go` |
| Vision open required | 0 |
| F-007 | 上传授权深度仍 deferred；本 VP 不扩张授权 scope |
| 是否消费业务解锁 | **否**。不改 Profile / 模块矩阵 / Manifest 为意图 |
| 结果 | **PASS**。不暂挂 `go` |

`consumer_vp` = VP-014；`last_freshness_review_at` = 2026-08-21；`next_freshness_review_trigger` = 若实施改变共同门禁 / Profile / 模块矩阵则重做。

## 为什么

- 新纲领波次独立工作区，避免写入已 closed 的 workspace-013。
- slug 与 VP id 对齐，便于组合索引。
- I-001～I-005 满足 V-F061；freshness 表满足 V-F062。

## 未选方案

- 继续 `planned` 只写 VP：用户已要求激活并开区。
- 重开 workspace-013：VP-013 已关门且默认不接新区。
- 一开区就改 handler / composition：R1 方案未冻结。
- 默认改为必须 MinIO/S3：与 A2 内嵌默认冲突。
- 把本 VP 当业务 VP 走 VP-011 式完整 freshness 矩阵：解锁 scope 不匹配。
