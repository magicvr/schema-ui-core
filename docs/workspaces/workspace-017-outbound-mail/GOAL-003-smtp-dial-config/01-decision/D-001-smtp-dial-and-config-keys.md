---
id: GOAL-003-smtp-dial-config
doc: decision-entry
record_id: D-001
status: accepted
parent: GOAL-001-outbound-mail
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
---

## D-001 · R2 拨号路径与配置键冻结（关闭 I-003 / I-004）

### 触发

Root R2 门禁：I-003（STARTTLS 587 vs 隐式 TLS 465，只钉一种可核对路径）、I-004（配置键名与凭证注入）须在 SMTP 接入实施前由决策关闭。VP-017 §SMTP 实现："显式主机/端口/凭证/From；STARTTLS 或隐式 TLS 由 lead Root 方案冻结为一种可核对路径"。

### 决定

**1. I-003 → verified：唯一拨号路径 = 隐式 TLS（465），证书校验强制开启**

- 实现形态：`tls.Dial`（TLS 自第一个字节起）→ `smtp.NewClient`；`tls.Config{ServerName: host, MinVersion: TLS1.2}`；禁止任何 `InsecureSkipVerify` 出口。
- AUTH：仅 `smtp.PlainAuth` over TLS（stdlib 在明文上拒发凭据，本路径全程 TLS）。
- **传输守卫（适配器级，非产品策略）**：`Subject` 含 CR/LF 或控制字符 → 拒发（header 注入面）；`TextBody` 换行规范化为 RFC 5322 CRLF（线格式编码，内容字节不变）。端口层仍只校验 To——守卫发生在 SMTP 适配器内并在代码注释留痕。
- 理由：隐式 TLS 从连接建立即加密，无 STARTTLS 的降级/剥离歧义，是"可核对路径"的最强形态（RFC 8314 对 submission 的推荐方向）；stdlib 原生支持，零新依赖。

**2. I-004 → verified：配置键名与注入规则**

```yaml
mail:
  smtp:
    host:      # MAIL_SMTP_HOST
    port:      # MAIL_SMTP_PORT；缺省 465
    username:  # MAIL_SMTP_USERNAME
    password:  # MAIL_SMTP_PASSWORD —— secret，只走 env / configs/.env 插值
    from:      # MAIL_SMTP_FROM；默认发件人（bare addr-spec）
```

- **fail-closed 配对规则**：五个键全空 = 未配置（合法，走默认 sink）；任意一个非空则 `host`/`username`/`password`/`from` 全部必填（port 缺省 465），缺失即拒绝启动，错误只报键名不回显值。`port` 若设置必须为 1–65535；`from` 必须为 bare RFC 5322 addr-spec。
- 校验挂 `ValidateProd` 链（每个环境含 development 生效，同 db/objects 先例）；secret 不入库的契约与 DB_PASSWORD / STORAGE_OBJECTS_S3_* 同文。

### 为什么

- 单路径钉死满足 VP"一种可核对路径"；R4 显式证据 harness 只需覆盖这一条拨号形态。
- 配置面完全复用既有三层机制（embedded default → YAML → env 插值 + CONFIG_ENV_FILE secret 层），operator 心智与排错路径不变。

### 未选方案

- **STARTTLS 587**：明文起步 + 协商升级，可核对性弱于隐式 TLS；部分企业中继只开 587 是其相对优势，但本波以可核对性优先。未选（将来如需第二路径属加法演进，须先过方案审计）。
- **双端口同时支持**：违反 VP 单方言冻结。未选。
- **无认证 SMTP（username/password 可选）**：VP 明示凭证为显式配置组成部分；放行无认证会制造 open-relay 误配面。未选。
