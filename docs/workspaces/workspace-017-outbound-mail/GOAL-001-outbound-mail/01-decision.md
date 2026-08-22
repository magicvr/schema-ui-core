---
id: GOAL-001-outbound-mail
doc: decision
status: active
parent: null
created: 2026-08-22
updated: 2026-08-22
version: 0.1.0
---

# 决策记录 · GOAL-001

## 信息需求与阶段门禁

> 状态以 `00-meta.md` 信息表为准（本表为镜像，须保持同号同状态）。

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 默认 sink 形态与取报文 | R1 方案 | R1 冻结 | R1 决策 | **verified**（D-002） | — | capture sink 容量 1 + slog 双写；`internal/mail.CaptureSink.Last()` 取报文 |
| I-002 | required | Send To 基数 | R1 方案 | R1 冻结 | R1 决策 | **verified**（D-002） | — | 单收件人 `To string`；多收件人将来加法演进 |
| I-003 | required | STARTTLS vs 隐式 TLS | R2 方案 | R2 接入前 | R2 决策 | **verified**（D-003） | — | 唯一路径 = 隐式 TLS 465 |
| I-004 | required | 配置键名与凭证注入 | R2 方案 | R2 接入前 | R2 决策 | **verified**（D-003） | — | `mail.smtp.*` / `MAIL_SMTP_*`，四键必填规则 |
| I-005 | non-blocking | HTML/MIME 是否进分母 | 关门叙事 | R4 | R4 或不进分母 | collecting | — | VP I-017-005 |
| I-006 | non-blocking | 重启生效 / 热加载不进本波 | 关门叙事 | R4 | 已随 VP 冻结 | **registered**（V-F071） | — | VP-017 §配置面 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-22 | 开区 scaffold 与 A6 纲领路线图 | accepted | [D-001-workspace-root-establishment.md](01-decision/D-001-workspace-root-establishment.md) |
| D-002 | 2026-08-22 | R1 发送合同冻结（关闭 I-001 / I-002） | accepted | [D-002-r1-send-contract-freeze.md](01-decision/D-002-r1-send-contract-freeze.md) |
| D-003 | 2026-08-22 | R2 拨号路径与配置键冻结（关闭 I-003 / I-004） | accepted | [D-003-r2-dial-and-config-freeze.md](01-decision/D-003-r2-dial-and-config-freeze.md) |
