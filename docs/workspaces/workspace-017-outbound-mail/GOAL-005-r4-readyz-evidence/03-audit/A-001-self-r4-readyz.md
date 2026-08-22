---
id: GOAL-005-r4-readyz-evidence
doc: audit-entry
record_id: A-001
source: self
goal: GOAL-005-r4-readyz-evidence
status: recorded
parent: GOAL-001-outbound-mail
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
---

## A-001 · R4 显式路径证据与 readyz 扩依赖自审

- **source**: `self`
- **日期**: 2026-08-22
- **scope**: readyz 探测语义、显式路径证据面、关门叙事留痕、边界遵守
- **verdict**: **pass**（无开放 required）

### 核对结果

| # | 核对面 | 结论 | 证据 |
|---|--------|------|------|
| 1 | VP 判据 4 | 未配置 SMTP：probe=nil，readyz 语义零变化（双 profile lifecycle 测试继续全绿）；显式配置：ESMTP Ping 进入 readyz 变参（与 objectStore HeadBucket 同机制） | composition_mail_test / lifecycle 全绿 |
| 2 | 探针-投递同路径 | Ping 与 Send 共用 `tlsConfig()`（隐式 TLS/校验恒开/ServerName=host），无漂移面；明文端点必败测试在案 | smtp.go / smtp_test.go |
| 3 | 显式投递可核对 | 离线：TLS fake 断言 envelope/DATA/AUTH；live：`MAIL_SMTP_TEST_*` gated 测试镜像 s3_live 先例，operator 可复跑 | smtp_live_test.go |
| 4 | I-005/I-006 关门叙事 | HTML/MIME 不进分母（合同无该字段）；重启生效（启动时 Load→构造单例）——均有代码事实支撑并落 README | kernel/mail.go / README |
| 5 | 边界遵守 | 无第二拨号路径、无热加载、无消费方注册；fx.Provide 移除后无重复实例 | git diff 核对 |

### Findings

| F-ID | 级别 | 内容 | 处置 | 状态 |
|------|------|------|------|------|
| N-001 | note | live 测试本轮未实跑（无真实凭据）；离线 TLS harness + 协议断言已构成"与生产合同等价的 harness"证据，live 面作为 operator 可复跑补充 | 按 VP 判据 3 原文（live **或** 等价 harness）属**分母外 note**，非 residual——独立审计 A-002 已复核认可此定性 | **closed（note）** |

> N-001 按残余信息处理：不影响 R4 门禁——VP 判据 3 允许"live 或与生产合同等价的 harness"，本波以等价 harness 交付。若编排器判定需用户书面接受，升级至 P-004。

### 结论

R4 门禁达成且证据可回指；无开放 required finding。C1/C2/C3 完成，本目标可关门（done 3/3）。Root 四阶段全部完成，进入关门审计。
