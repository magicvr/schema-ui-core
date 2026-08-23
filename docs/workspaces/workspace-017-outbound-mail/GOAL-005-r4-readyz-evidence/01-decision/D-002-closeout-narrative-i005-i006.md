---
id: GOAL-005-r4-readyz-evidence
doc: decision-entry
record_id: D-002
status: accepted
parent: GOAL-001-outbound-mail
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
---

## D-002 · 关门叙事留痕：HTML/MIME 不进分母；重启生效（关闭 I-005 / I-006）

### 触发

Root 信息表 I-005（HTML/MIME 是否作为可选体，最晚 R4）、I-006（生效方式）须在关门叙事前留痕关闭。

### 决定

1. **I-005 → verified：HTML/MIME 不进本波退出分母**。已交付合同只有 `TextBody` 纯文本（R1 冻结、R2/R3 实现一致）；VP 建议方向即"纯文本进分母，HTML 不进"。未来如需 HTML，属加法演进（新增字段/方法 + 方案审计），不破坏既有合同。
2. **I-006 → registered→closed（叙事投影）**：生效方式 = 进程重启后读取配置（`Load()` 在启动时一次性解析 mail.smtp 并构造 sender 单例）；无热加载。与 VP §配置面冻结决策同构，V-F071 台账投影就此闭合。

证据：`kernel/mail.go` 合同面（无 HTML 字段）；`config.Load()` 启动时序；README「出站邮件」节明示两者。

### 未选方案

- HTML 可选体进分母：超出 VP 首波冻结的退出分母，且引入 MIME 编码/注入面。未选。
