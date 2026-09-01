# I-026-004 评估 · mail cachedAdapter 是否迁移到通用 Cache 端口（2026-09-01）

> 信息项（non-blocking）：「既有 mail runtime `cachedAdapter` 是否迁移到端口（评估，不强制；其版本戳失效语义可能不匹配通用 TTL）」。**用户确认（2026-09-01）：不迁移，评估留痕**。判据 #2 评估面闭合。

## 1. 被评估对象（事实）

`apps/api/internal/mail/runtime.go`：

```go
type Switcher struct {
    ...
    mu    sync.Mutex
    cached *cachedAdapter   // 每次 Send 路径上的发送器缓存
}
type cachedAdapter struct {
    updatedAt int64         // mail_config.updated_at（UnixMilli）
    sender    kernel.MailSender
}
```

- `Send` 每次先 `LoadRuntime()`（**每次 Send 都读 `mail_config WHERE id = 1`**），随后 `currentSender()`：`cached.updatedAt == cfg.UpdatedAt` → 复用缓存发送器；否则 `buildAdapter(*cfg)` 重建并覆盖缓存。
- 失效源 = **DB 行版本戳**（admin 在设置页保存渠道配置时 bump `updated_at`）；切换语义 = **事件驱动、零延迟热切**（保存即生效，下一条 Send 立即用新渠道）。
- `buildAdapter` 为纯对象构造（`NewOutboxSink` / `NewResend` / `NewSMTP`），无昂贵资源。

## 2. 与通用 Cache 端口的语义对照

| 维度 | mail cachedAdapter | kernel.Cache 端口（R1 冻结） |
|------|--------------------|------------------------------|
| 失效模型 | 版本戳（`updated_at` 相等性） | TTL（绝对 / 滑动过期） |
| 失效延迟 | 零（事件驱动） | TTL 窗（时间驱动） |
| 读取原语 | 按版本条件命中 | Get/Set/Delete（无版本条件读取） |
| 价值面 | 缓存 sender 对象（构建廉价） | 通用值缓存 |

## 3. 候选迁移方案与判定

| 方案 | 描述 | 判定 |
|------|------|------|
| A. TTL 近似 | `kernel.Cache.Set("mail:runtime", sender, AbsoluteExpiry{30s})`；命中即用，过期重建 | **否决**：渠道切换将延迟 ≤ 30s，违反 VP-017「热切」语义（保存即生效）；需连带修改 mail 合同 |
| B. 版本戳作 key | key = `"mail:runtime:"+updatedAt`，值 = sender | **否决**：Get 前仍须 DB 读拿 `updated_at` 判版本（`LoadRuntime` 未省掉），缓存零收益 + key 膨胀 |
| C. 部分迁移（快照） | 仅缓存 `RuntimeConfig` 快照（含 updated_at），替代每次 Send 的 DB 读 | **否决**：同上——端口无版本条件读取原语，判断「快照是否最新」仍需 DB 读；收益为零 |
| D. 不迁移（保留现状） | `cachedAdapter` 留在 mail 内部（已并发安全、合同冻结于 VP-017） | **采纳**（用户确认）：语义精确匹配 mail 需求；无行为漂移；维护成本不变 |

## 4. 结论

通用 Cache 端口的设计目标（TTL 驱动的供应商无关值缓存）与 mail 配置缓存的失效模型（版本戳驱动的零延迟热切）**不匹配**；端口当前无「版本条件读取」原语（不为此扩契约——无消费者收益）。mail 保留自有 `cachedAdapter`，`runtime.go` 零改动。若未来 mail 需要通用缓存（如渠道发送频率控制等），届时另行评估；本评估不构成任何 mail 行为承诺变更（VP-017 合同不变）。

## 5. 证据链

- 被审代码：`apps/api/internal/mail/runtime.go`（L116～L229：Switcher / cachedAdapter / currentSender）。
- 端口合同：workspace-026 `GOAL-002/01-decision/D-002-cache-port-contract.md`（§1～§8）。
- 裁决：workspace-026 `GOAL-004/01-decision/D-001-r3-adjudication.md`（2026-09-01 用户确认不迁移）。