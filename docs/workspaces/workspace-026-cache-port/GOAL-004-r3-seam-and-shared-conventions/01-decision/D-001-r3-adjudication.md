---
doc_type: goal-decision
id: D-001-r3-adjudication
parent: GOAL-004-r3-seam-and-shared-conventions
date: 2026-09-01
status: accepted
version: 0.1.0
---

# D-001 · R3 裁决：I-026-004 不迁移 + F-002 fx 容器挂载（2026-09-01 用户裁决 ×2）

## 上下文

R3 两项前置裁决点（P-004 / P-005）：I-026-004（mail cachedAdapter 迁移评估）与 A-002 F-002 义务（组合根 kernel.Cache 单一实例的挂载方式）。编排器提交带建议的选项（基于 `internal/mail/runtime.go` 实读与组合根 Fx 结构），2026-09-01 用户裁决**全部采纳建议项**。

## 裁决记录

| 项 | 选项 | 裁决 |
|----|------|------|
| I-026-004 | ① 不迁移，评估留痕 ② 迁移到端口（TTL 近似） ③ 部分迁移（快照缓存） | **采纳①**：mail `cachedAdapter` 版本戳失效语义（`mail_config.updated_at` 驱动 · 零延迟热切）与通用 TTL 端口不匹配——TTL 近似会造成渠道切换最多延迟一个 TTL 窗（行为漂移，违反 VP-017 即时热切语义）；版本戳作 key 又仍需每次 DB 读判版本（零收益）。保留 mail 自有实现；评估留痕（attachments）；判据 #2 评估面闭合。 |
| F-002（挂载义务） | ① fx 容器持有 + newMux 注入点 ② 无消费者不持有（接入点登记） | **采纳①**：`fx.Provide(newCache)` 将单一实例注册进 Fx 容器（进程级长生命周期持有）；`newMux` 依赖注入使构造沿 mux→server→lifecycle 依赖链 eager（fail-closed 保持）；newMux 参数 = 首个消费者显式接入点。满足 F-002 字面义务且不推翻 R2「组合根单一实例」冻结。 |

## 未选方案

| 项 | 未选 | 理由 |
|----|------|------|
| TTL 近似迁移（②） | 未选 | 渠道切换延迟 ≤ TTL 窗 = VP-017「热切」语义漂移；mail 合同需连带修改，收益为零 |
| 版本戳作 key 的部分迁移（③） | 未选 | Get 前仍需 DB 读拿 updated_at（LoadRuntime 未省），双读/版本键技巧复杂度 > 收益 |
| 无消费者不持有（②） | 未选 | 与 F-002 字面义务有差距（实例将被 GC），需审计再认可；fx 方案成本低且诚实 |

## 影响

- C2 落盘分母 = 本合同：架构短文（`docs/architecture/cache-redis-seam-and-track.md`：§2 接缝声明 + §3 轨道约定）+ mail 评估附件 + 组合根 fx 改造（`fx.Provide(newCache)` + `newMux`/`newMuxWithExtraProviders` 注入 + 4 个测试调用点更新）。