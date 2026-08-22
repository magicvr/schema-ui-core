---
id: GOAL-002-port-contract-freeze
doc: decision-entry
record_id: D-001
status: accepted
parent: GOAL-001-outbound-mail
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
---

## D-001 · R1 发送合同冻结（sink 形态 / To 基数 / 公共面规则）

### 触发

Root R1 门禁：I-001（默认 sink 形态与测试取报文）、I-002（Send To 基数）必须在方案冻结前由证据关闭。VP-017 已给出方向级建议（capture/log sink 可取出最后一封；建议单收件人），VRev-037/038 均 pass 后激活。本决策按 VP 建议方向冻结为可实施合同。

### 决定

**1. 端口形态（`kernel.MailSender`）**

```go
type MailMessage struct {
    To       string // 唯一收件人，RFC 5322 addr-spec
    Subject  string // 原样传递，不做产品策略
    TextBody string // 纯文本正文；HTML/MIME 不进本波分母（I-005）
}

type MailSender interface {
    Send(ctx context.Context, msg MailMessage) error
}
```

- 同步 `Send`：无队列、无 outbox、无重试；失败以 error 返回，由调用方处理（VP 非目标已排除 outbox / 外部队列）。
- **From 不进 `MailMessage`**：默认发件人来自配置（VP §意图），由适配器在投递时补齐。调用方永远不需要知道发件头。
- 合同级校验只有一条：`To` 必须非空且为 **bare RFC 5322 addr-spec**（`net/mail.ParseAddress` 可解析，且解析结果与输入一致——display-name 形式如 `Alice <alice@example.com>` 一律拒绝）；`Subject`/`TextBody` 原样传递。业务字段策略（主题模板、正文长度等）留给消费者。

**2. I-001 关闭：默认 sink = 进程内 capture sink + 结构化日志双写**

- 未配置 SMTP 时进程正常启动，发送走 capture sink：容量 1 的内存环（只保留最后一封），同时打一行结构化日志（slog）。
- 测试取报文方式：`internal/mail.CaptureSink.Last() (kernel.MailMessage, bool)` 取出最后一封，`Reset()` 清空。取报文 API 属 dev/test 适配器面（`internal/mail` 包），**不进 kernel 公共合同**——生产消费者只见 `kernel.MailSender`。
- 落地归 R3 子目标；本决策只冻形态。

**3. I-002 关闭：单次 Send 单收件人**

- `To string` 非集合。理由：VP 明示"建议单收件人，降低转发面"；本波全部已知消费场景（校验邮件、恢复码）均为单收件人；多收件人将来可用加法方法演进，不破坏本合同。

**4. 公共面类型边界**

- handler 与模块 Provider 的公共契约只能引用 `kernel.MailSender` / `kernel.MailMessage`；禁止在任何公共签名出现 `net/smtp` 或具体 SMTP 客户端类型。SMTP 拨号细节封死在 `internal/mail` 适配器内（R2）。核对动作归 R3 公共面 sweep。

### 为什么

- 端口放 kernel、适配器放 internal/mail：完全复用 VP-014 objectstore 先例（`kernel.ObjectStore` + `internal/objectstore`），薄内核不引入新依赖（stdlib `net/smtp`/`net/mail` 足够，不加第三方邮件库）。
- From 由配置注入而非消息携带：防止调用方伪造发件人（From 伪造是 VP-009 持续程序的关注点，端口层直接消除该面）。
- capture 容量固定 1：满足"测试可取出最后一封"的最低合同；不做多封历史缓存，避免 dev 默认变成隐形邮件存档。

### 未选方案

- **只写结构化日志的 sink**：测试断言需解析日志文本，脆且慢；未选。
- **强制 Compose 常驻 Mailhog**：VP 明确不进分母；且给 mvp/dev 加硬依赖违反"没 SMTP 也能启动"。未选。
- **To 为小集合（[]string）**：转发面更大、去重/校验语义复杂化；单收件人不满足时再加法扩展。未选。
- **第三方邮件库（gomail 等）**：薄内核不新增重量级依赖；stdlib `net/smtp` + STARTTLS 已覆盖唯一钉死的拨号路径（R2 决策）。未选。
