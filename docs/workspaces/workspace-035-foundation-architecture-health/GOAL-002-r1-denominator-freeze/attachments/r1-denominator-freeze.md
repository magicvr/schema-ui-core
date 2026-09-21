---
doc_type: goal-attachment
id: r1-denominator-freeze
parent: GOAL-002-r1-denominator-freeze
status: recorded
created: 2026-09-09
updated: 2026-09-09
version: 0.1.0
---

# R1 · 对照分母冻结

本附件是 I-035-001 / I-035-004 / I-035-005 的冻结正文。R2 对照矩阵必须按本表的包含行取证，不得把排除行写成缺口。

消费候选：HEAD `5c341ec7`（VRev-087 freshness 候选）。本冻结不消耗任何 trigger-gated 行。

## 1. I-035-001 · 包含 / 排除

### 1.1 包含（R2 必须出对照行）

| 面 | 路径 / 对象 | 对照什么 |
|----|-------------|---------|
| 模块契约 | `apps/api/kernel/module.go`、`provider.go`、`contribution.go`、`persistence.go` | 薄内核、Provider 六项、fail-closed、无业务 import |
| Profile | `apps/api/kernel/profile.go`；`apps/api/internal/config/config.default.yaml` | mvp/admin 默认集；compiled candidate 不得静默进默认 |
| 组合根 | `apps/api/internal/composition/`；`apps/api/internal/server/`；`apps/api/assembly/assembly.go` | 静态汇集、依赖图、生命周期、供应商类型不进公共面 |
| Store | `kernel/store.go` + `internal/store/` | 双方言端口；公共面无 `*sql.Tx` / 驱动类型 |
| ObjectStore | `kernel/objectstore.go` + `internal/objectstore/` | 本地盘默认 + S3 兼容；无 `os.File` 进公共面 |
| Cache | `kernel/cache.go` + `internal/cache/` + `docs/architecture/cache-redis-seam-and-track.md` | 内存默认；Redis 仅接缝声明 |
| RateLimiter | `kernel/ratelimit.go` + `internal/ratelimit/` | 内存默认；AllowRecord/Reserve/Cancel；Redis 仍 gated |
| EventBus | `kernel/eventbus.go` + `internal/eventbus/` | 进程内运输；outbox/MQ 仅接缝 |
| Mail | `kernel/mail.go` + `internal/mail/` | 发送端口；mock/Resend/SMTP 适配器不进公共面 |
| Observability | `internal/obs/`；composition metrics/tracing 接线 | Prometheus 类 scrape + OTLP；缺省 no-op |
| Shutdown | `kernel/lifecycle.go`；`internal/composition/shutdown_drain_test.go`；`internal/jobs/shutdown_reclaim_test.go` | 停机顺序 / HTTP drain / Job 语义 |
| Job 运行时 | `internal/jobs/` | RT-Q01 进程内六态；≠ EventBus |
| Telegram 通道端口 | `kernel/telegram.go` + `internal/channel/telegram/` | 内核端口 vs 适配器；无 SDK 类型进公共面 |
| 密钥轮换面 | JWT current+previous 接线（composition / config / auth）；不重做 dump | RT-K03 已交付合同是否仍成立 |
| Manifest 聚合 | `internal/manifest/`；composition NavGroup 归一化 | 模块 fragment 所有权；group 可选；无中央业务注册表 |
| 文档权威 | `docs/architecture/module-architecture.md`、`module-contribution-playbook.md`、`cache-redis-seam-and-track.md`、`overview.md`、`directory-layout.md`、`monorepo-layout.md` | as-designed vs as-built；overview 现时节过期记为文档缺口 |
| 泄漏抽检 | 模块 `provider.go` 与 handler 对 `*sql.Tx`、S3/Redis/Telegram SDK、供应商包的 import | 抽检即可，不把每个业务 handler 做成分母行 |

### 1.2 排除（R2 不得当成本 VP 缺口）

| 面 | 排除 | 理由 |
|----|------|------|
| Admin 体验增强 | Command Palette、Saved Views、统一 Toast、未保存保护 | 产品项；路线图可建议下一拍，不在本 VP 实现 |
| 业务域产品 | `biz.digital-offer` 业务规则、钱包账本产品面 | 已 closed VP-031/029；只抽检端口泄漏 |
| VP-009 / VP-010 波次史 | 安全 finding 台账、UX 符合性波次 | 正交持续程序 |
| 前端产品页 | `apps/web` 除 Manifest/NavGroup 消费与协议 host 之外 | 不是内核/端口健康 |
| 未实现的 gated 供应商 | Redis 客户端、MQ broker、搜索引擎、K8s | 评估只确认仍为接缝/缺席，不把「没实现」写成缺陷 |
| 分发包装 | npm/CLI 发布残余（VP-024） | 除非证明破坏内核端口公共面 |
| VP-034 Dashboard 组 | `workspace` 组 vs 历史顶层单例 | Admin 展示残余，点名原 VP，不进架构端口分母 |

Web 只抽：Shell 是否仍只消费 Manifest NavGroup、有无中央业务导航注册。

## 2. I-035-004 · closed VP residual 入册

分类只决定**是否进入本 VP 登记册供 R2/R3 对照**。进入 ≠ 现在修。

| ID | 来源 | 内容 | 入册？ | R2/R3 预期分类（草案，R3 才冻结） |
|----|------|------|--------|----------------------------------|
| RES-013-migrator | VP-013 D-002 | 无产品 SQLite→PG 搬运器 | **入册** | 接受残余 / 仍非目标，除非产品触发 |
| RES-014-migrator | VP-014 I-014-004 | 无产品本地盘→对象存储搬运器 | **入册** | 同上 |
| RES-015-otlp-sink | VP-015 F-003 | in-repo otlp-sink 不解析 | **入册** | 接受残余（显式 endpoint 才导出） |
| RES-015-metrics | VP-015 I-015-003 | Store/对象/Job 指标不进分母 | **入册** | 仍 gated 或另立，R3 对照后定 |
| RES-016-revoke | VP-016 I-016-005 | JWT 立即失效未选 | **入册** | 接受残余 |
| RES-016-mfa-wrap | VP-016 | `admin.mfa` wrapping 不随 JWT previous 重包 | **入册** | 接受残余或另立，R3 定 |
| RES-021-harness | VP-021 V-F083 | 进程级 harness `!windows`；compose stop 以 linux CI 核销 | **入册** | 接受残余；复审 = CI 失败 |
| RES-026-redis | VP-026/027/032 | Redis 实现仍 RT-Q03/Q05 gated | **入册** | 仍 gated；不消耗 trigger |
| RES-028-broker | VP-028 | outbox/MQ 仍 gated | **入册** | 仍 gated |
| RES-030-keyfile | VP-030 R-009 | bot master key 文件与 DB 同目录 | **入册** | 已 accepted-residual；复审 = KMS 波 |
| RES-T03-tz | roadmap RT-T03 | DB `timestamptz` 未做 | **入册** | registered；R3 定是否建议下一拍 |
| RES-P04-pool | RT-P04 | `MaxOpenConns` 现状 vs 文档锚点 | **入册** | 仍 gated；R2 核实现状锚点是否过期 |
| RES-024-dist | VP-024 D-001 四项 | hosted CI / shell 类型面 / GH Packages / C 类 fork 包化 | **不入册** | 分发面，非内核端口 |
| RES-034-dashboard | VP-034 GOAL-003 | Dashboard 现行 `workspace` 组 | **不入册** | 保留原 VP 点名 |
| RES-017-smtp | VP-017 历史 SMTP 分母 | 已被渠道模型取代 | **不入册** | 历史，不再当缺口 |

## 3. I-035-005 · 「现在修」vs 另立

| 规则 | 决定 |
|------|------|
| 默认 | **另立**：R2/R3 分类为「现在修」的代码/端口变更，登记后新开 VP 或交 VP-009/010 波次，**不在本 VP 内改端口实现** |
| 本 VP 内允许 | ① 文档卫生（`overview.md` 现时节、roadmap 现状锚点/A 序列）——属判据 4，在 R4 做；② 对照矩阵/测试注释/抽检用的只读断言，不改变公开端口语义 |
| 本 VP 内禁止 | 实现 Redis/MQ/搜索/多实例；改 Store/Cache/RateLimiter/EventBus/Mail/ObjectStore/Telegram 公开接口；改 Profile 默认集；重开已 closed VP；把 Admin 体验增强当整改 |

R1 不授权任何生产代码整改。I-035-003（对照是否迫使改 Charter）仍 collecting，最晚 R3。
