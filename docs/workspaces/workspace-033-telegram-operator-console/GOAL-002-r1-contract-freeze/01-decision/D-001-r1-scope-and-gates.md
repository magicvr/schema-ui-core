---
doc_type: goal-decision
id: D-001-r1-scope-and-gates
parent: GOAL-002-r1-contract-freeze
date: 2026-09-04
status: superseded
version: 0.1.0
---

# D-001 · R1 范围与方案门禁

## 背景

Root 已按 VP-033 建立 R1 → R2 → R3 → R4 的串行路线。R1 横跨连接模式、轮询生命周期、Dispatcher 占用位、heartbeat、发言权反馈、公网 URL 与 Fake Bot API 验收，符合 P-001 的独立阶段范围。

## 已确认的范围

- 本目标只承载 R1 合同冻结与 R2 入口证据，不重开 VP-030，不修改默认 Profile，不实现付费命令或多实例 polling。
- `I-033-001`～`I-033-008` 的产品语义沿用 VP-033 与 Root 已有用户书面冻结；`I-033-009/010` 继续作为 non-blocking 未知。
- R1 阶段若进入方案冻结，必须先处理本记录列出的 `I-033-011`～`I-033-013` required 信息；本记录不把探子建议冒充用户决定。

## 待用户裁决的实现选择

| 门禁 | 推荐方案 | 未选/备选方向 | 影响 |
|------|----------|---------------|------|
| `I-033-011` 公网 URL 与 mode 的配置边界 | 复用现有 `auth.public_base_url` 作为显式 webhook origin；仅把非敏感 `mode` 持久化到 `telegram_config`，token/secret 继续加密 | 新增 Telegram 专属 `webhook_public_base_url` 配置/持久化面；或全部使用全局静态配置 | 决定迁移、重启与 webhook URL 的单一来源；现有 `auth.public_base_url` 已有 URL 校验，但不能默认推断为用户选择 |
| `I-033-012` 默认 | 新安装/旧行缺省为 `polling`；生产文档与部署配置推荐显式 `webhook`，不运行时猜测环境 | 根据 `runtime.mode` 自动推断，或所有环境默认 webhook | 决定首次连接、迁移默认值和开发/生产可预测性 |
| `I-033-013` 生命周期 | 在 Telegram 包内新增 connection manager + polling loop，暴露 `Start() error` / `Stop(context.Context) error`，由 composition `OnStop` 统一 drain | 直接由 HTTP handler 管理 goroutine，或复用 scheduled-tasks 而不增加 Telegram owner | 决定 SIGTERM 可验证性、并发互斥和停机是否等待 loop 退出 |

## 当前状态

本记录的待决策状态已由用户在 2026-09-04 书面裁决，并由 [D-002-r1-contract-freeze](D-002-r1-contract-freeze.md) 取代。原条目保留作为方案门的历史记录；R1 合同、验证矩阵和信息项投影以 D-002 为准。
