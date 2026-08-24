---
id: D-002
doc: decision-entry
goal: GOAL-006-channel-provider-contract
status: accepted
created: 2026-08-24
updated: 2026-08-24
version: 1.0.0
---

# D-002 · R5 渠道供应商合同冻结（关闭 I-011 / 预冻 I-010）

## 背景

Root D-006 升级分母后，R5 必须在写 Resend 适配器或设置页之前把渠道合同一次冻清。四个决策点已于 2026-08-24 由用户书面裁决（本会话问答留痕）：

1. **I-011 mock 持久化 = DB 表 + 迁移**（非进程内扩容）。
2. **解析规则 = 显式键 `mail.channel`**。
3. **mock 保留策略 = 有界**，默认最近 500 条，保留条数可由管理员经 mock 渠道配置指定。
4. **I-010 一并预冻**（Resend 键名与 fail-closed）。

## 决定（可实施条款 · C2 依据）

### §1 渠道标识与注册形状

1. 第一期具名渠道 id 为小写枚举字符串：**`mock`**、**`resend`**、**`smtp`**。
2. 渠道是 `kernel.MailSender` 的**适配器实现**，注册/解析只发生在 composition root（`internal/composition`）；适配器代码住 `internal/mail`。模块与 handler 公共面**只见 `kernel.MailSender` / `MailMessage`**，不得出现 mock/resend/smtp 客户端类型（维持 R1 合同，R3 sweep 结论继续有效）。
3. 渠道集封闭：新增供应商须走新的决策，不在本期分母。
4. **SMTP 渠道保留**：R2 已落地的 SMTP 适配器、键名（`mail.smtp.*` / `MAIL_SMTP_*`）与隐式 TLS 拨号路径**原样保留**为可选渠道；它不再是唯一生产权威。不删除、不改写、不回退其实施史（重申 Root D-006 条款 2）。

### §2 解析规则（`mail.channel`）

1. 新增 YAML/env 配置键 **`mail.channel`**（env `MAIL_CHANNEL`），取值 `mock|resend|smtp|""`；空串 = 未显式选择。
2. 解析算法（composition 唯一实现点）：
   - **显式选择**：`mail.channel=X` → 解析 X 的适配器；X 所需配置不完整 → 启动 fail-closed（沿用 SMTP 防御性拒绝先例）。显式选 `mock` 时生产渠道块被忽略（允许，记 info 日志）。
   - **未显式选择（向后兼容推导）**：统计完整配置的生产渠道数：
     - 仅 Resend 块完整 → `resend`；仅 SMTP 块完整 → `smtp`；
     - 两者都完整 → **歧义，启动 fail-closed**，逼操作员显式设 `mail.channel`；
     - 都不完整/未配置 → `mock`（现行为保持：现有仅配 `mail.smtp.*` 的部署不受影响）。
3. 「块完整」判定沿用 SMTP 现行语义：出现任一该渠道键即视为「已触碰该块」，触碰即要求全部必填键齐备，缺项按 fail-closed 报错并点名缺失键（同 `validateMail` 现行文案风格）。
4. 本条只冻结**文件/环境层**的初始解析。管理面热切换（运行时改选渠道）归 R7 / Root I-009，其落库细节不得在本条下假装已解决。

### §3 mock 渠道合同（I-011 关闭）

1. **持久化 = Store 表**（双方言 SQLite/Postgres），经编译迁移目录新增迁移建表（R6 实施，表级设计归 R6；本条冻语义）。记录至少承载：收件人、主题、纯文本正文、创建时间；重启后仍可检视，多实例共享同一份。
2. **取信面 = 独立管理 API**（落实 I-012「独立 API」条款）：`GET /api/mail/outbox`（新→旧分页列表）+ `GET /api/mail/outbox/{id}`（含完整正文详情），管理员鉴权；不塞进 `/api/settings/default`，不做 PATCH。响应包络沿用现行 handler 约定（R6 细化）。
3. **保留策略 = 有界淘汰**：默认保留**最近 500 条**，插入超限淘汰最旧；上限值属 mock 渠道配置项，管理员可在设置「邮件」tab 调整（该调整面 R7 交付；R6 先以默认 500 实现，键位预留）。
4. **语义边界**：mock 是管理员出站记录（联调/验收工具），**不是**用户站内通知，不是 Notification Transport，无投递失败重试语义（写入成功即成功）。
5. **CaptureSink 处置**：`internal/mail.CaptureSink` 及其 `Last()` 测试面原样保留（单元测试继续用）；R6 起 composition 默认解析改为 mock 发布器，CaptureSink 不再是装配默认 sink。此为接线变更，不是对 R1～R4 实施史的回退或改写。

### §4 resend 渠道键名（I-010 预冻）

1. 配置键：**`mail.resend.api-key`**（env `MAIL_RESEND_API_KEY`；SECRET——仅允许 env / `configs/.env` 注入，禁止 YAML 明文，机制同 `MAIL_SMTP_PASSWORD`）、**`mail.resend.from`**（裸地址，校验规则同 SMTP From）。
2. fail-closed 规则：触碰任一 `mail.resend.*` 键即要求块内必填键齐备；缺项报错点名缺失键，prod 校验与 composition 防御性双层拒绝，机制镜像 SMTP 现行实现。
3. HTTP 调用细节（endpoint、错误映射、超时）与 live/harness 证据判据归 R6/R8，不在本条。
4. `readyz` 探针扩展仍按 R4 现状（仅 SMTP 显式路径），生产探针统一归 R8。

## 为什么

- DB 表是唯一同时满足「管理员可检视」「重启不丢」「多实例一致」的选项，且 Store 迁移目录已是成熟机制，边际成本低。
- 显式 `mail.channel` 让多渠道并存有确定解，且空值时完全保持现行为，存量部署零迁移负担；两生产渠道并存时的 fail-closed 歧义拒绝比静默优先级更安全。
- 有界默认 500 条防无限增长；上限可配置满足联调需要更久历史的场景。
- I-010 早关使 R6 开设即可开工，避免下一阶段再停一轮等裁决。

## 未选方案

- **进程内环形缓冲**：重启丢、多实例各自为政、管理端只能看本进程——用户裁决否决。
- **不加选择键纯存在性推导**：多渠道全配时无法表达操作员意图，只能靠删除配置「二选一」，运维体验差。
- **全量保留不设限**：无限增长风险交给未来清理，得不偿失。
- **mock 复用 CaptureSink 扩容**：容量语义与持久化需求不符，且会把测试替身变成产品面。
- **热切换细则并入本条**：I-009 未闭，写了就是假装解决；归 R7。

## 影响

- Root I-011 → **verified**（本决策）；Root I-010 → **verified**（提前于最晚阶段 R6 接入前关闭）。
- R6 可直接实施：mock 表 + 迁移 + outbox API + resend 适配器 + `mail.channel` 解析 + 默认接线切换。
- 不改应用代码（R5 边界不变）；不改 Profile/模块矩阵；VP-018 冻结不动。
