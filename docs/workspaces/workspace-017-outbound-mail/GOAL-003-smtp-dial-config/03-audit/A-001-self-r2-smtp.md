---
id: GOAL-003-smtp-dial-config
doc: audit-entry
record_id: A-001
source: self
goal: GOAL-003-smtp-dial-config
status: recorded
parent: GOAL-001-outbound-mail
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
---

## A-001 · R2 SMTP 接入与配置面阶段自审

- **source**: `self`
- **日期**: 2026-08-22
- **scope**: R2 拨号路径冻结、SMTP 适配器安全姿态、配置面 fail-closed、边界遵守
- **verdict**: **pass**（无开放 required）

### 核对结果

| # | 核对面 | 结论 | 证据 |
|---|--------|------|------|
| 1 | VP 对齐 | 单一可核对拨号路径 ✓；显式主机/端口/凭证/From ✓；凭证 YAML+env fail-closed 不入库 ✓；未配置不挡 mvp/dev ✓（validateMail untouched 放行） | VP-017 §SMTP 实现 / §配置面 |
| 2 | 安全姿态 | TLS 自首字节 + MinVersion1.2 + 校验恒开（rootCAs 只换锚集，无 InsecureSkipVerify 路径）；PlainAuth 仅 over TLS；Subject 控制字符拒发；错误只报键名不回显 secret | `smtp.go` / `config_mail_test.go` 泄密用例 |
| 3 | fail-closed 验证 | 部分块（host-only / port-only / 缺 from）启动即拒；非法 env 端口 Load 即拒；display-name From 拒绝——均有测试 | `config_mail_test.go` |
| 4 | 可核对性 | loopback TLS harness 断言 AUTH 身份、envelope、DATA 全文；明文端点必败证明无降级路径 | `smtp_test.go` 实跑输出 |
| 5 | 边界遵守 | composition / handler / readyz 零改动；capture sink 未提前实现；单方言未扩 | git diff 核对 |

### Findings

| F-ID | 级别 | 内容 | 处置 | 状态 |
|------|------|------|------|------|
| F-001 | minor | D-001 初稿未写明适配器级传输守卫（Subject 控制字符拒发 / 正文 CRLF 规范化），实现先于文本 | 当场补记 D-001 §1「传输守卫」，与代码行为一致 | **fixed** |

### 结论

R2 门禁（I-003/I-004）已由决策关闭并有测试证据；安全姿态符合"可核对唯一路径"冻结；无开放 required finding。C1/C2/C3 达成，本目标可关门（done 3/3）。R3 将落地 capture sink + composition 接线 + 公共面 sweep。
