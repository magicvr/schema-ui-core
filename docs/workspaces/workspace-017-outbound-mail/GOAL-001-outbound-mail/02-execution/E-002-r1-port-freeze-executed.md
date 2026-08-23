---
id: GOAL-001-outbound-mail
doc: execution-entry
record_id: E-002
goal: GOAL-001-outbound-mail
status: recorded
parent: null
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
---

# E-002 · R1 执行：合同冻结 + kernel 端口代码（2026-08-22）

## 已发生事实

- 开设子目标 `GOAL-002-port-contract-freeze`（五件套齐全，parent = 本 Root），承接 R1 治理上下文。
- 子目标 D-001 冻结发送合同；Root D-002 关闭 I-001/I-002（均 → verified），两处信息表已同号同步。
- 端口代码落地：`apps/api/internal/kernel/mail.go`（`MailMessage` / `Validate` / `MailSender`）+ `mail_test.go`。
- 实跑证据（go1.26.0 windows/amd64）：`go test ./internal/kernel/ -run Mail -v` 全 PASS（3 测试 / 9 子用例）；`go build ./...` OK；`go test ./internal/kernel/` 全绿；`go vet ./internal/kernel/` 干净。
- Root 路线图 R1 → 已完成；progress 1/4。R1 冻结前未改 composition / handler / readyz（符合"冻结前不改 apps/api 行为面"约束——本次仅新增 kernel 合同文件，无行为变更）。

## 证据

| 主张 | 路径 |
|------|------|
| 子目标合同正文 | `../GOAL-002-port-contract-freeze/01-decision/D-001-r1-contract-freeze.md` |
| 本 Root 决策 | [D-002-r1-send-contract-freeze.md](../01-decision/D-002-r1-send-contract-freeze.md) |
| 端口代码 | `apps/api/internal/kernel/mail.go` |
| 子目标执行条目 | `../GOAL-002-port-contract-freeze/02-execution/E-001-port-code-and-tests.md` |

## 未做

- 未实现任何适配器（SMTP 归 R2、capture sink 落地归 R3）；未动配置键与 `readyz`；未改 web 端。
