---
id: GOAL-001-outbound-mail
doc: decision
status: active
parent: null
created: 2026-08-22
updated: 2026-08-24
version: 0.2.0
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
| I-005 | non-blocking | HTML/MIME 是否进分母 | 关门叙事 | R4 | R4 或不进分母 | **verified**（D-005） | — | 纯文本进分母；HTML 不进（合同无该字段） |
| I-006 | non-blocking | 重启生效 / 热加载不进本波 | 关门叙事 | R4 | 已随 VP 冻结 | **verified**（D-005 · V-F071 闭合） | — | 重启生效；启动时构造单例 |
| I-007 | required | 第一期渠道集 mock+Resend；SMTP 保留 | 现行分母 / R5 | R5 | D-006 | **verified**（D-006） | — | 用户采纳讨论方案 |
| I-008 | required | mock = 管理员出站记录 | R5 / R6 | R5 | D-006 | **verified**（D-006） | — | 非用户通知 |
| I-009 | required | 热切换密钥/失败语义/多实例 | R7 | R7 实施前 | R7 决策 | collecting | — | I-017-009 |
| I-010 | required | Resend 配置键 | R6 | R6 接入前 | R6 决策 | collecting | — | I-017-010 |
| I-011 | required | mock 持久化形态 | R5 | R5 | GOAL-006 | collecting | — | I-017-011 |
| I-012 | required | 设置「邮件」tab 形状 | R7 | R5 可冻形状 | D-006 | **verified**（D-006） | — | 独立 API |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-22 | 开区 scaffold 与 A6 纲领路线图 | accepted | [D-001-workspace-root-establishment.md](01-decision/D-001-workspace-root-establishment.md) |
| D-002 | 2026-08-22 | R1 发送合同冻结（关闭 I-001 / I-002） | accepted | [D-002-r1-send-contract-freeze.md](01-decision/D-002-r1-send-contract-freeze.md) |
| D-003 | 2026-08-22 | R2 拨号路径与配置键冻结（关闭 I-003 / I-004） | accepted | [D-003-r2-dial-and-config-freeze.md](01-decision/D-003-r2-dial-and-config-freeze.md) |
| D-004 | 2026-08-22 | R3 默认 sink 接线与公共面 sweep 规则 | accepted | [D-004-r3-default-sink-wiring.md](01-decision/D-004-r3-default-sink-wiring.md) |
| D-005 | 2026-08-22 | R4 readyz 探测、显式路径证据与关门叙事（关闭 I-005 / I-006） | accepted（R4 技术事实仍成立；组合层关门结论由 D-006 supersede） | [D-005-r4-readyz-and-closeout.md](01-decision/D-005-r4-readyz-and-closeout.md) |
| D-006 | 2026-08-24 | 否决 Root/VP 关门并升级渠道分母（关闭 I-007 / I-008 / I-012） | accepted | [D-006-revoke-closeout-channel-upgrade.md](01-decision/D-006-revoke-closeout-channel-upgrade.md) |
