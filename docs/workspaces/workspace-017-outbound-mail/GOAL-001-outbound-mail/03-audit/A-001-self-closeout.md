---
id: GOAL-001-outbound-mail
doc: audit-entry
record_id: A-001
source: self
goal: GOAL-001-outbound-mail
status: recorded
parent: null
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
---

## A-001 · Root 关门自审（R1～R4 全阶段 · 方向级判据核对）

- **source**: `self`
- **日期**: 2026-08-22
- **scope**: Root 全量——VP-017 六条方向级退出判据、治理结构完整性、证据可回指性
- **verdict**: **pass**（无开放 required）

### 退出判据逐条核对（VP-017 §方向级退出判据）

| # | 判据 | 结论 | 证据 |
|---|------|------|------|
| 1 | 内核发送端口落地；公共面不再把 SMTP 客户端类型当发送合同 | ✅ `kernel.MailSender`/`MailMessage` 唯一合同；rg 证明 `net/smtp` 仅在 internal/mail，handler/modules 零引用 | GOAL-002 D-001、GOAL-004 E-001 sweep |
| 2 | 未配置 SMTP 仍能开发与快测；capture/log sink 测试可取最后一封 | ✅ capture 缺省 + 双 profile lifecycle 启动不变量测试全绿；Send→Last 回读断言通过 | GOAL-004 E-001/A-001 |
| 3 | 显式 SMTP 可核对至少一封投递；配置不完整 fail-closed | ✅ loopback TLS harness 断言 envelope/DATA/AUTH 全会话（与生产合同等价）；live env-gated 测试就绪；validateMail 部分块拒收测试在案 | GOAL-003 E-001、GOAL-005 E-001 |
| 4 | 仅显式配置后 readyz 扩依赖；未配置不得 not-ready | ✅ probe=nil 缺省；显式时 ESMTP Ping 入变参（objectStore 同机制） | GOAL-005 D-001/E-001 |
| 5 | 无 SMS / 第二邮件方言；未改 Charter；未进账号 email/邀请/恢复/模板/业务域 | ✅ 单一拨号路径冻结；git 历史仅 mail/config/composition/README 面；users 表未动 | 各子目标边界节 + git log |
| 6 | 开放 required finding = 0（或合法闭合） | ✅ 四个子目标 A-001 均 self pass；F-minor 均已 fixed 复验 | goal-tree 状态表 |

### 治理结构核对（AGENTS §12）

编号单调（GOAL-002～005，无复用）；id=文件夹名；五件套+ledger 目录齐全；路线图 4/4 与 progress 一致；goal-tree 已同步；决策/执行/审计均落盘且证据路径有效；I-001～I-006 全部 verified。

### Findings

无 required。N-001（live 测试未实跑）已在 GOAL-005 A-001 按"与生产合同等价的 harness"判据留痕为残余，范围与复核触发明确。

### 结论

四阶段全部完成、六条判据全绿、无开放 required。建议进入独立关门审计（项目默认 provider = 本地 grok build `/audit`），independent 意见落盘并响应后本 Root 可 `done`。
