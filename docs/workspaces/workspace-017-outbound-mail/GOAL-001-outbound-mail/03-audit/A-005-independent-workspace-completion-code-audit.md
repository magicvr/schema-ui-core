---
id: A-005
doc: audit-entry
goal: GOAL-001-outbound-mail
status: recorded
parent: null
created: 2026-08-24
updated: 2026-08-24
version: 1.0.0
source: independent
audit_scope: 工作区 workspace-017-outbound-mail 完成情况复核（R1～R8 全阶段；execution-facts + 关门复核；用户指令「独立审视代码完成情况，不以治理文件为准」）
verdict: conditional
---

# A-005 · independent · 工作区完成情况代码级独立复核（2026-08-24）

> 审计提供方：ox-alpha（DeepSeek Harness `/audit` 会话）。方法：**只读扫描治理文件获取主张清单后，逐条下到代码第一手核对**；回归不在信任台账结论的前提下复用——`go build` / 全量 `go test ./...` / web `tsc -b && vite build` / 全量 vitest 均为本会话独立重跑。本条目不改任何 status/progress/goal-tree。

## 范围与区间

- 对象：`workspace-017-outbound-mail` 全部 9 个目标所声称的交付（R1～R8），重点是 R5～R8 现行分母的代码事实。
- 区间：工作树当前状态（HEAD `e83ff61`，working tree clean）。
- 方法说明：先从各目标 meta/E 条目提取**可证伪主张**（文件路径、行为语义、测试名、迁移号、错误码、i18n 键），再逐条到源码/测试/配置中核对；能重跑的回归全部独立重跑。

## 成果（有证据 · 第一手核实）

| 主张 | 独立核对结果 | 证据 |
|------|--------------|------|
| R1 端口唯一合同、公共面无供应商类型 | ✅ 属实 | `apps/api/internal/kernel/mail.go`（`MailSender` 单方法端口 + `MailMessage.Validate` 单收件人裸地址校验，From 由适配器盖章）；handler 消费 `OutboxReader`/`MailAdminService` 抽象，无客户端类型 |
| R6 配置层：`mail.channel` + `mail.resend.*` + 冻结解析算法 | ✅ 属实 | `internal/config/config.go` L159-177/L509-618/L1168-1291：YAML+env 双通道加载（小写归一）、`ResolveMailChannel` 显式胜出→块可用性门禁、空值推导（单生产块胜出 / 双全配 fail-closed / 均未配 mock）、`validateMail` 触碰即全键 + 解析门禁 |
| 迁移 0051 `mail_outbox` 双方言 | ✅ 属实 | `corepersistence/migration/migration.go` L52-74/L138-144（SQLite INTEGER / PG BIGINT + created_at 索引，v51 描述符） |
| mock 发布器有界保留（默认 500）同事务淘汰 | ✅ 属实 | `internal/mail/outbox.go`：`DefaultOutboxCap=500`；Send 内 INSERT+DELETE 同事务（`ORDER BY created_at DESC, id DESC LIMIT ?`）；List 无正文新→旧、Get 含正文、未知 id `ErrOutboxRecordNotFound` |
| Resend 适配器 fail-closed + 请求形状 + 探针 | ✅ 属实 | `internal/mail/resend.go`：构造器双键必填+裸地址校验；POST `{base}/emails` Bearer/from 盖章；非 2xx 报状态码不泄密钥；`Ping` GET /domains 5s 仅报状态码 |
| 独立管理取信面（不进 `/api/settings/*`） | ✅ 属实 | `handler/mail_outbox.go`：GET list `{items,total,page,pageSize}`（limit≤200）+ detail 404；Bearer + `settings.read` 门禁 |
| R7 密钥写后不可读回 + 主密钥加密落库 | ✅ 属实 | `internal/mail/secrets.go`（env SHA-256 或 data/ 下 0600 key 文件自动生成；AES-256-GCM；空串直通）；`runtime.go PublicView` 仅 `*Set` 布尔；迁移 0052 `_enc` 列 |
| 热切换：候选先验证后落库、失败保留旧 sender | ✅ 属实 | `runtime.go Update`：merge→`buildAdapter` 构造校验（SMTP 另加 5s ESMTP Ping）→任一失败在持久化前返回→成功才原子写行+刷缓存 |
| 试发走同一端口禁止旁路 | ✅ 属实 | `handler/mail_admin.go` `mailTestSend` 经 `svc.Send`（=Switcher）；PUT/test-send 记 `mail.channel-update`/`mail.test-send` operation log（常量见 operationlog repository.go L81-82） |
| 错误码入冻结清单 | ✅ 属实 | `errorcatalog.go` L188-190 三码双语表；`handler/error_contract_test.go` L74 pin 进合同清单 |
| 迁移 0053 + store 冻结账目同步 | ✅ 属实 | operationlog migration v53 `operation_log_mail_events`；`store/identity.go` `completeFingerprintCatalogHead = 53` 且注记 v51/v52 |
| R8 readyz 生产探针三态 | ✅ 属实 | `composition.go newMailRuntime` L697-751：boot=resend→`Resend.Ping`、boot=smtp→ESMTP Ping、mock/空→nil；三态测试 `TestMailRuntime*` 实测在库 |
| live 缝 env-gated skip | ✅ 属实 | `resend_live_test.go`：`MAIL_RESEND_TEST_API_KEY/FROM/TO` 缺一即 skip；离线 `TestResendPing` 断言 401 报状态码 |
| web 设置「邮件」tab（含 UX 精修三项） | ✅ 属实 | `modules/settings/schema/settings.json` tab-mail 单一 custom 节点、零 PATCH `/api/settings/default`；`src/components/mail-admin-tab.tsx` 渠道条件渲染、密钥 password 输入空=保持、outbox 表仅 mock、试发 subject/body；`registerCustomComponent` 注册并被 W25 守卫测试引用；双语 i18n 键抽查 11/11 存在 |
| 回归绿（api + web） | ✅ **本会话独立重跑** | `go build ./...` exit 0；`go test ./...` **45 包全部 ok exit 0**；web `npm run build`（tsc -b+vite）exit 0；vitest 全量 **78 文件 / 1099 用例全绿 exit 0** |
| 史不回退 / 边界守住 | ✅ 属实 | `capture.go`/`smtp.go` 及其测试原样保留且套件绿；composition 默认已切 Switcher-over-mock；settings.json 无 SMS/通知面；git 无 `*.key` 跟踪（c78bab6 已移除误提交 key 并 ignore） |

GOAL-009 `attachments/exit-denominator-evidence.md` 的判据 1～7 指针逐条落地核对，未发现死指针或名不副实。

## 对照成功标准（现行分母判据 1～7）

1～7 全部 satisfied（判据 3 的 live PASS 属操作员本地凭据证据，审计不可独立复跑；其代码缝与等价 harness 均已第一手核实，符合 VP「live 或等价 harness」口径——见 N-4）。

## Findings

| F-ID | 级别 | 严重度 | 意见 | 证据路径 |
|------|------|--------|------|----------|
| F-001 | **required** | med | **台账现势性复发（与 A-004 F-001/F-002 同类，漏改两个文件）**：① `GOAL-009-r8-evidence-readyz/00-meta.md` frontmatter 仍为 `progress: 0/3`（分母还与正文检查点表 C1～C4 的 4 不符），检查点表仍写「当前 **0/4**」且证据列为占位文本（"api 代码与测试"），与 goal-tree/workspace/Root「done · 4/4 · A-001 self pass」直接矛盾——status 改了、progress 与证据列没跟上；② `GOAL-001-outbound-mail/00-meta.md` frontmatter `progress: 7/8` 与正文路线图「当前 8/8」及 goal-tree 矛盾。不否定代码完成，但「完成」主张的权威台账自身不一致，后续读者以目标件为准时会得到相反结论 | `docs/workspaces/workspace-017-outbound-mail/GOAL-009-r8-evidence-readyz/00-meta.md` L9/L28-37；`GOAL-001-outbound-mail/00-meta.md` L9 vs L41 |
| F-002 | recommended | low | 操作员侧样例配置 `apps/api/configs/config.yaml` 的 `mail:` 节仍只有 SMTP 键位，未同步 R6 的 `channel:`/`resend:` 键位注释；嵌入默认 `internal/config/config.default.yaml`、`configs/.env.example`、apps/api/README.md 均已覆盖，故仅为样例一致性缺口 | `apps/api/configs/config.yaml` L180-192 |
| N-1 | note | — | web 测试计数台账漂移：E-002 记 1097、E-003 记 1100，本次实跑 1099（全绿）——计数随提交波动属正常，建议台账引用时注明「当日快照数」 | 本条目「成果」表回归行 |
| N-2 | note | — | `mailTestSend` 成功后经第二次 `PublicView()` 读渠道写入审计 detail——与并发热切换存在良性 TOCTOU，仅影响日志字段归因，不影响发送路径 | `handler/mail_admin.go` L154-158 |
| N-3 | note | — | `MAIL_SWITCH_REJECTED`(409) 同时承载「候选校验失败」与「DB 持久化失败」两类语义；前者是设计主语义，后者属实现细节搭车，可接受 | `handler/mail_admin.go` L114-118 |
| N-4 | note | — | live 投递 PASS（admin@eshowy.top → magicvr@hotmail.com）依赖本地凭据，独立审计不可复跑；缝+harness 已核实，与 A-004 N-3 同向，无需动作 | GOAL-009 E-003 / 证据包判据 3 |

required finding = **1（F-001）**。

## 必改项汇总

- **F-001**：同步两处 meta 现势性——GOAL-009 `00-meta.md`（frontmatter progress → 4/4；检查点表四行补「**完成**：<证据指针>」；修正 0/3 分母笔误）与 GOAL-001 `00-meta.md`（frontmatter progress → 8/8）。纯文档事务，一次提交可闭合。

## 与既有意见的异同

- 与 A-003（self）/A-004（independent）结论方向一致：现行分母代码级完成成立，未发现名不副实的交付主张。
- 差异：A-004 以判据抽查为主；本轮按用户指令把核对密度提高到**每个 E 条目主张逐条对源码**，并将两类回归（api 全量 / web 全量+build）在本会话独立重跑，因此能把「回归绿」从台账主张升级为第一手事实。
- 新增缺口：A-004 的 F-001/F-002 修的是 goal-tree 与索引，本轮发现同类现势性问题残留在 GOAL-009 / GOAL-001 两个 `00-meta.md`（F-001）——属该次整改的遗漏面，不是新缺陷类别。
- A-004 判据 4 称「tab-mail recordSource 不含密钥字段」，UX 精修后 tab 已收敛为单一 custom 节点、recordSource 移入组件内 GET `/api/mail/config`（密钥仅布尔）——语义不变，表述已过时但无碍。

## 结论 + 建议给编排器/用户的下一步

**verdict: conditional** —— 工作区 17 的**代码完成情况独立复核全面成立**（R1～R8 所有可证伪主张均与代码一致，双端全量回归独立重跑全绿）；唯一 required F-001 为纯文档台账现势性（med、低风险、易修），按尺度不可无条件视台账为权威，修复后即可无条件支持既有关门结论。

建议下一步：用 **`/govern`** 响应本意见——按 fixed 路径闭合 F-001（两处 meta 同步提交），F-002/N-1～N-4 可 note 留痕或顺手修；闭合后本意见合并效力为 pass，无需重开任何阶段。

## 声明

本意见不修改 status/progress；响应由 /govern 处理。
