# 独立关门审计任务（workspace-017 · GOAL-001）

你以独立审计者身份工作。请先加载本仓库的 `/audit` 技能（skills 目录内 `audit`）并按其规则执行；若技能不可用，按以下等效规则执行。

## 审计对象

- 工作区：`docs/workspaces/workspace-017-outbound-mail/`（Root：`GOAL-001-outbound-mail`，progress 4/4，拟关门）
- 被审面：Root 关门（R1～R4 全阶段 + 方向级退出判据），含四个子目标：
  - GOAL-002-port-contract-freeze（R1 合同冻结 + kernel 端口）
  - GOAL-003-smtp-dial-config（R2 SMTP 隐式 TLS 465 唯一路径 + mail.smtp 配置面）
  - GOAL-004-default-sink-surface-sweep（R3 capture sink + composition 接线 + 公共面 sweep）
  - GOAL-005-r4-readyz-evidence（R4 readyz 扩依赖 + 显式路径证据 + I-005/I-006 叙事）

## 你必须核对的材料

1. Root 五件套（00-meta 路线图与信息表、01-decision D-001～D-005、02-execution E-001～E-005、03-audit 台账）。
2. 四个子目标的五件套与 A-001 self 审计。
3. 愿景层对齐：`docs/vision/plans/VP-017-outbound-mail.md`（退出判据、非目标、信息需求 I-017-001～006）。
4. 代码事实（可读仓库）：`apps/api/internal/kernel/mail.go`、`apps/api/internal/mail/*`、`apps/api/internal/config/config.go` 中 mail 相关段、`apps/api/internal/composition/composition.go` 中 newMailSender/newMux 接线、`apps/api/README.md` 出站邮件节。
5. 已知环境事项：`go test ./...` 在 `internal/store` 有两个 live-Postgres 集成测试失败（TestOpenPostgresProbeIntegration / TestPostgresMigrateRunnerIntegration）；`git diff be23164..HEAD -- apps/api/internal/store` 为空——该包未被本工作区改动，属共享 probe 库遗留状态的环境问题。请独立判断其是否应阻塞本次关门。

## 输出要求

只输出**一个** Markdown 审计意见条目（将被编排器代贴入 Root `03-audit/A-NNN`，source 标为 independent）。格式：

```markdown
## A-00N · Root 关门独立审计（<一句话 scope>）

- **source**: independent
- **日期**: <今天>
- **scope**: <覆盖面>
- **verdict**: pass | conditional | fail

### 核对面与结论
<表格或列表：每项核对点 → 结论 → 证据路径>

### Findings
<表：F-ID | 级别(required/recommended/minor) | 内容 | 建议 | 是否可 fixed 闭合；无则写"无">

### 结论
<是否同意 Root done；条件或保留意见>
```

约束：只出意见，不修改任何文件；不重复 self 审计已闭合的 minor 除非你认为闭合不成立；发现 required 问题必须给出可执行的修正建议。
