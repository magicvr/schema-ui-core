你是本仓库的独立交叉审计员。请先加载本仓库 skill：.grok/skills/audit/SKILL.md（即 /audit 独立审计流程），按其要求执行一次 S4 关门前 independent 复核。只输出审计意见，不修改任何文件、不改任何 status/progress。

## 审计目标

工作区 workspace-009-production-hardening（canonical: docs/workspaces/workspace-009-production-hardening/）的波次目标 GOAL-011-w11-api-web-security-audit（W11）。

## 审计范围

1. A-001（source: independent，verdict: fail）的 6 条 required（F-001～F-006）在提交 72a5397 中的修复是否 genuine fixed：
   - F-001: apps/api/internal/modules/authsession/users_repository.go CreateUserManagement EXISTS 扫 bool（postgres 方言）
   - F-002: apps/api/internal/handler/resources.go 新增 TrashTxRecorder/TrashTxDeleter 同事务删除+回收站快照（dict types/entries、scheduled-tasks 实体），快照失败整体回滚
   - F-003: apps/api/internal/handler/auth.go MFA 分支限流 + modules/mfa/store/repository.go proof 懒清理与 fail_count<5 守卫
   - F-004: modules/mfa/service.go 双密钥（当前+previous）解密回退与成功后重封；composition.go 传入 AuthJWTSecretPrevious
   - F-005: modules/logincaptcha/store/repository.go 单语句守卫生成 DELETE（RowsAffected 判定胜者）
   - F-006: handler/wallet.go reconcile Decode 失败 400（io.EOF 空 body 保持全库哨兵）+ submit/cancel/retry 改 wallet.write
2. 修复是否改错既有逻辑或引入新缺陷（对照 A-001 原文的每条主张与既有测试语义）。
3. recommended 处置（02-execution/E-003-w11-recommended-disposition.md）：fixed 11 条（F-007/008/009/010/011/012/013/014/016/017/018）的代码证据；overruled 2 条（F-015 刷新家族吊销 vs 客户端跨标签重试协议 A-002 F-003；F-019 前端 custom action vs 服务端硬门禁）是否有据。
4. 回归主张复跑：apps/api `go vet ./...` 与 `go test ./...`（可只跑受影响包），apps/web `npx vitest run` / `npx tsc -b`（如可行）。

## 材料路径（先读）

- docs/workspaces/workspace-009-production-hardening/GOAL-011-w11-api-web-security-audit/00-meta.md
- .../01-decision/D-002-w11-scope-and-go-hold.md
- .../02-execution/E-002-w11-s3-implementation.md（含新增测试清单）
- .../02-execution/E-003-w11-recommended-disposition.md
- .../03-audit/A-001-w11-independent.md（原文 finding）
- .../03-audit/A-002-w11-self.md（self 审计）
- 工作区 workspace.md 与 docs/vision/alignment.md（绑定核对：Charter schema-ui-core-admin-foundation@0.2.0；primary_plan=VP-009-production-hardening）

## 输出要求

- 条目头：source=independent / auditor=grok-build（grok-4.6 · reasoning high）/ verdict / scope / 日期
- 逐条闭合判定表：F-001～F-006（fixed / 不成立 / 部分）＋证据；recommended 处置核对；回归复跑结果
- Findings：按 required / recommended / informational 分级
- 必改项汇总：开放 required 条数
- 明确声明：不改 status/progress；是否确认"代码闭合条件已满足"（最终关门仍归编排器/用户）