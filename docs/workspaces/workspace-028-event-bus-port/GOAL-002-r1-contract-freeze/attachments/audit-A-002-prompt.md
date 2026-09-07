/audit

你是独立交叉审计员（source: independent）。本轮**禁止落盘、禁止改任何文件**；只把正式意见以 Markdown 打到 stdout，供编排器写入 `03-audit/A-002-*.md`。

- 工作区：`workspace-028-event-bus-port`
- 被审目标：`GOAL-002-r1-contract-freeze`
- audit_type：close-out
- auditor：grok-build（grok-4.6 · reasoning high）
- 日期：2026-09-01
- 对照 self：`03-audit/A-001-contract-freeze-closeout-self.md`（verdict pass）
- 冻结分母：`01-decision/D-002-event-bus-port-contract.md` v0.1.0
- 裁决证据：`01-decision/D-001-info-adjudication.md`
- 实现：`apps/api/kernel/eventbus.go`、`apps/api/kernel/eventbus_test.go`
- **排除**：R2 进程内 channel 实现；R3 接缝/I-028-004 权属；R4 证据矩阵。

请：
1. 完整阅读仓库根 `AGENTS.md` 相关节、`skills/prompts/05-independent-audit.md`、工作区 `workspace.md` / `goal-tree.md`、被审目标五件套与 D-001/D-002、A-001、kernel 两文件。
2. 独立复跑（工作目录 `apps/api`）：`go vet ./kernel/...`；`go test ./kernel/... -count=1`；`go test ./kernel/ -count=1 -v -run Event`；`go build -o NUL ./kernel/`；`gofmt -l kernel/eventbus.go kernel/eventbus_test.go`。
3. 核：P-004 三项裁决是否落入 D-002；D-002 ↔ 实现逐节一致；快测覆盖合同 §10；I-028-004 升 required 且未伪装关闭；红线（go.mod 无 broker、不改 Profile/Manifest、不实现 outbox）。
4. 越界核账：相对 HEAD，代码变更应 ⊆ `apps/api/kernel/eventbus.go` + `eventbus_test.go` + workspace-028 文档 + VP-028 信息表。
5. 输出完整 A-002 报告（frontmatter + findings F-00N required|recommended|informational + verdict + open_required）。不要修改 status/progress。
