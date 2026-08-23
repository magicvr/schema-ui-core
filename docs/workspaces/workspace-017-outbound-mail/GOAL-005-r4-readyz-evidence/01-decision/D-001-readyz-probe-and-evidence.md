---
id: GOAL-005-r4-readyz-evidence
doc: decision-entry
record_id: D-001
status: accepted
parent: GOAL-001-outbound-mail
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
---

## D-001 · R4 readyz 探测语义与显式路径证据面

### 触发

Root R4 门禁：仅显式配置 SMTP 后 `readyz` 才扩该依赖（未配置不得 not-ready）；显式配置后可核对至少一封投递。R2 已冻结唯一拨号路径。

### 决定

1. **Ping probe**：`(*mail.SMTP).Ping(ctx)` —— 复用 `Send` 的冻结拨号形态（隐式 TLS、ServerName=host、MinVersion TLS1.2、校验恒开），完成最小 ESMTP 往返（220 问候 + QUIT）。不进 kernel 端口，仅 composition 消费。
2. **接线语义**：`newMailSender` 改为 objectStore 同构三元返回 `(sender, probe, error)`——capture 缺省时 probe 为 nil（readyz 语义零变化）；显式配置时 probe = `sender.Ping`，经 `RegisterWithMFAProbes` 变参进入 readyz（复用 VP-014 GOAL-003 机制：共享 readyz 截止、任一失败整体 unavailable）。装配从 fx.Provide 收敛为 NewApp 内直接构造（R3 的容器注入无实际消费方，且 probe 必须同步到达 handler 注册点——与 objectStore 完全同构）。
3. **证据面**：
   - 离线：loopback TLS fake 断言 Ping 全绿 / 明文端点必败；composition 断言缺省 probe=nil、显式 probe 非 nil。
   - live：`internal/mail/smtp_live_test.go` 镜像 `s3_live_test.go` 先例——`MAIL_SMTP_TEST_*` 六个 env 全设才跑真实 465 端点投递一封并核对，否则干净 skip（plain `go test ./...` 保持离线）。

### 为什么

- probe 与 Send 共用一条拨号路径：readyz 绿 = 投递路径的 TLS/证书/可达性全部可核对，不存在"探针绿但发送走另一条路"的漂移面。
- 直接构造收敛：消除无人消费的第二实例，保持与 objectStore 单例装配一致。

### 未选方案

- **probe 进 kernel.MailSender 接口**：污染生产端口（VP-014 HeadBucket 同判例）。未选。
- **AUTH 后再 QUIT 的重探测**：探测只需可达性+TLS+横幅，认证正确性由投递测试覆盖。未选。
