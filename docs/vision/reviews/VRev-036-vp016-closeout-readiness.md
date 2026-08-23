---
doc_type: vision-review
id: VRev-036
status: active
source: self
created: 2026-08-22
updated: 2026-08-22
version: 0.2.0
parent: null
---

# VRev-036 · VP-016 关门就绪度审视（2026-08-22）

| 字段 | 值 |
|------|-----|
| source | self |
| auditor | `/vision` · 会话编排（grok-4.6）；本轮对代码 / 测试 / live 双密钥 HTTP **独立复验**，不以 Goal 台账为关门充分条件 |
| scope | `VP-016-key-rotation-and-backup` 组合层关门就绪 · 退出判据 1–6 · 代码实现 · 本会话测试与 live 轮换路径 · Vision required · 有界 residual · 组合索引 |
| audit_type | vision-plan（关门就绪度 · 代码成果独立核验） |
| verdict | pass |
| 建议 class | editorial（组合层关门 + 索引同步 + residual 点名；不改 Charter 方向） |
| open required | 0 |

## 范围与结论

只读核对：`docs/architecture/principles.md` P-006、`docs/vision/alignment.md` §6/§7/§9、`charter.md` `@0.2.0`、[VP-016-key-rotation-and-backup](../plans/VP-016-key-rotation-and-backup.md)（审视时 `active` v0.2.0）、`roadmap.md`、`workspaces.md`、`revisions.md`（至 VR-038）、`reviews.md` 与 `reviews/VRev-001`～`VRev-035`、lead 工作区 `workspace-016-key-rotation-and-backup`（绑定与 Root 状态只作**指针**，不替代代码核验）。

**本轮独立核验（非转录治理记录）**：

1. **源码**：`apps/api/internal/auth/auth.go`（`previousSecret`、`NewWithRepositoryAndPrevious`、`issue()` 只用 `a.secret`、`verifyAccess` current→previous、`SignAccessToken` 无 `kid`、refresh 仍 opaque SHA-256）；`internal/config`（`AuthJWTSecretPrevious` / `AUTH_JWT_SECRET_PREVIOUS` / 生产同强度+同值守卫）；`internal/composition` `newAuthenticator` 接线；两份 YAML 注释键 + env 通道；`compose.yaml` 可选 `AUTH_JWT_SECRET_PREVIOUS: ${AUTH_JWT_SECRET_PREVIOUS:-}`。相对基线 `5195104` 的产品 diff = 上述 11 文件（含 README 一行与根 `compose.yaml`），**不含** `apps/web`、`internal/modules`、`docs/vision/charter.md`。
2. **测试（本会话复跑）**：`go test ./internal/auth -run TestDualKeyRotationOverlapWindow|TestParseAccessTokenExpiredAndWrongSecret -count=1` ok（3.314s）；`./internal/config -run TestJWTSecretPreviousConfig|TestValidateProd` ok（0.710s）；`./internal/composition -run TestNewAuthenticatorWiresPreviousSecret|TestSQLitePostRotationRecovery` ok（3.936s）；`TestPostgresPostRotationRecovery`（`R16_PG_DUMP_CONTAINER=r2-pg-probe`）ok（7.598s）；`./cmd/server` ok（11.484s）；`./internal/handler -run TestLogin|TestRefresh|TestAuth` ok（8.928s）。`go vet ./internal/auth ./internal/config ./internal/composition ./cmd/server` 0 finding。
3. **live（本会话，独立 SQLite + 独立二进制 + `127.0.0.1:25280`，不复用 GOAL-006 叙述）**：缺省单密钥（无 previous）→ `/healthz` `/readyz` **200**、login **200**。K1 签发后重启 `current=K2 previous=K1` → 旧 access `GET /api/accounts/me` **200**、新 login **200**、旧 refresh **200**。再重启仅 K2（退役 previous）→ 新 access **200**、旧 access **401**。`docker compose config`（不设 PREVIOUS）输出 `AUTH_JWT_SECRET_PREVIOUS: ""`。

未把 Goal `progress=5/5` 或 Root A-001/A-002 正文当作退出判据的充分证据。治理记录仅用于定位实现与对照开放 required。

**总判：pass（0 open required · 1 open recommended）。**

**关门的实质证据已齐备**（代码 + 本会话测试 + 本会话 live），可按 alignment §7 做**有界 closed**。Vision Review open required = 0；对齐链成立；激活后 Charter 仅 editorial（VR-037/VR-038），无 strategic 宽阻断。组合索引仍写 VP-016 `active`，且 `workspaces.md` 仍错误投影「Root active 0/5」——这是待用户书面确认后的投影同步，**不是**实现缺口。本轮用户意图为「审视是否满足关门条件；独立符合代码实现并检验无引入 bug；是的话关门」。

本意见原文**不**把组合索引改写成 `closed`。

### 核对事实

| 核对项 | 结论 | 证据（本会话独立核验，除非标明指针） |
|--------|------|------|
| 单愿景 / `vision_ref` | **pass** | 唯一 active Charter `schema-ui-core-admin-foundation@0.2.0`；VP-016 `vision_ref` 精确匹配 |
| 工作区绑定 | **pass** | `workspace-016` 唯一 lead / delivery；`plan_refs` / `primary_plan` / `vision_role: delivery` 合规；Charter `primary_workspace` 仍为 workspace-001 |
| 区证据指针（§7.1） | **pass（指针）** | goal-tree Root `done 5/5`、GOAL-002～006 `done`。**不**据此放行；放行依据是下表退出 1–5 的代码/测试/live |
| 实现层开放 required | **pass（指针）** | Root A-002 F-001 → `fixed`；开放 required = 0。本轮未改 Goal finding |
| 退出 1 · JWT current+previous；签发只用 current；重叠窗 previous 可验 | **pass** | `verifyAccess`；`issue()` 仅 `a.secret`。本会话双密钥 4 子用例绿。live：K2+prev=K1 时旧 access `/api/accounts/me` 200；退役 previous 后旧 access 401 |
| 退出 2 · 未配置 previous 时默认仍能开发与快测 | **pass** | YAML previous 缺省空；Compose 可选透传 `""`；生产缺 previous 不 fail-closed。live 缺省 healthz/readyz/login 200；`cmd/server` 本轮绿 |
| 退出 3 · 轮换后恢复：SQLite `VACUUM INTO` **与** PG `pg_dump`/`pg_restore` | **pass** | 本会话双方言循环均 PASS（非 skip）。T4 = `NewApp` Start；T5 = A1/A2/A3 |
| 退出 4 · 显式双密钥：一轮换路径 **与** 一轮换后恢复路径 | **pass** | 轮换路径 = auth 重叠窗 + composition 接线 + **本会话 live**；恢复路径 = 本会话 SQLite + PG 测试 |
| 退出 5 · 未进 A3/KMS/PITR/Admin/业务；未改 Charter；未假装热加载或第二套 dump | **pass** | `git diff --name-only 5195104`：`apps/web` / `internal/modules` / `charter.md` 为空。无 KMS/热加载 API。备份只消费既有 `VACUUM INTO` 与官方 `pg_dump`/`pg_restore` |
| 退出 6 · required = 0 | **pass** | Vision Review 与实现层开放 required 均为 0 |
| Vision required（§6 门禁 8） | **pass** | `reviews.md` open required = 0；VRev-035 为激活审视，本条为关门就绪首份 |
| Charter strategic 后 re-align | **pass** | 激活后仅 VR-037/VR-038（editorial）；无宽阻断 |
| 组合索引当前陈述 | **pass（待同步）** | Charter 关系节 / `roadmap.md` 第 16 行与 RT-K03「registered」/ `reviews.md` 摘要 / `workspaces.md`（仍写 Root 0/5）/ 区 `workspace.md` 仍写 VP-016 `active` |

## Findings

#### V-F069 · 组合层关门须同步索引，并显式映射 exit 1–6 ↔ 本轮独立证据、点名有界 residual

- level: `recommended`
- status: `open`
- severity: low
- impact: alignment §7.2 允许有界 closed，但 residual 必须点名到具体 workspace / goal id。若只让 Root `done` 而组合索引仍称 `active`、且 `workspaces.md` 仍写「Root 0/5」，后续读者会把 A5 读成未交付或与代码事实相反。
- finding: |
  1. 用户确认组合层关门时一次写清 exit 1–6 ↔ **本轮独立**证据（源码路径、本会话 `go test`/`go vet`、本会话 live 缺省 + 显式轮换/退役；治理 A 条目仅作指针）。
  2. residual 至少点名：`workspace-016` / `GOAL-001-key-rotation-and-backup` / **I-016-005**（立即失效未选，默认 previous 可验已交付，仍 collecting）。另点名旁路耦合：**`admin.mfa` TOTP 密文用当前 JWT secret 经 HKDF 派生 AES 钥封装**（GOAL-017 既有设计，本 VP 产品 diff **未改** `internal/modules/mfa`）。JWT previous 重叠窗**不**重包 MFA；`admin` profile 下轮换 current 后须重登记 MFA，或由未来 wrapping-key 合同承接。不进 A5 签名密钥分母（非 KMS、非第二套 dump）。缺省 `mvp` profile **不含** `admin.mfa`。
  3. 同步 `roadmap.md`（VP 行 + **RT-K03 → delivered**；RT-P05 轮换后恢复语义已由本 VP 补齐）/ `workspaces.md`（纠正 0/5）/ Charter 关系节 / `reviews.md` 摘要 / 区 `workspace.md`：VP-016 → `closed`；当前无 active **交付** VP（持续程序仍为 VP-009/VP-010）。将 VP 信息表 I-016-001～004 保持 verified。Root `done` 不能冒充 VP 层用户确认。
- evidence:
  - 本会话 live：缺省 healthz/readyz/login 200；K2+prev=K1 时旧 access 200 + 旧 refresh 200；退役 previous 后旧 access 401
  - 本会话测试：auth / config / composition（含双方言恢复）/ cmd/server / handler 指定套件全绿；vet 0
  - `git diff --name-only 5195104` 11 文件；`docs/vision/roadmap.md` RT-K03 仍写 VP-016 承接；`workspaces.md` 仍写 Root 0/5
  - `apps/api/internal/modules/mfa/service.go` `NewService` 注释：TOTP 钥从 JWT secret HKDF 派生；composition 只把 current `secret` 传入 MFA
  - alignment.md §7.2
- closure: |
  `/vision` 在用户书面确认组合层关门时按上列三项一并完成。本 finding 不阻断「就绪」结论，只约束关门落盘形状。
- 建议 class: `editorial`

### 不构成 fail / 不新开 required 的诚实边界

1. 本 `pass` **不是**「组合索引已 closed」：用户书面确认与索引原子同步仍待发生（本轮用户已给出「满足则关门」）。
2. `/vision` 本入口只写 `source: self`。用户要求「独立符合代码实现并检验无引入 bug」已在本报告用源码+测试+live 兑现；不是 `/vision-audit` 的 independent VRev。无独立 Vision Review 不是 alignment 强制项（强制时机仅为 Charter 初建与 strategic）。
3. MFA wrapping 与 JWT secret 共用是 **GOAL-017 既有耦合**，不是本 VP 引入的回归（产品 diff 无 `mfa/`）。它不否定 JWT 重叠窗合同。点名为有界 residual，不新开 required。
4. 本轮 live 使用独立 SQLite + 空 `CONFIG_ENV_FILE`，避免本机 `configs/.env` 的 postgres 方言干扰。这是取证隔离，不是「默认路径改成必须有 previous」。
5. 不把 progress=`5/5` 当作关门权威。
6. 架构 A3（多实例/Redis/队列）、KMS、PITR、热加载、Admin 密钥页本就不在退出分母。

### 声明

本意见不修改 Charter / VP / Goal status 或 progress；required/recommended finding 的响应由 `/vision` 追加在本报告中；实现层执行仍交 `/govern`。原 verdict 与 finding 原文不得改写。本入口不关闭 Goal finding。

### 门禁含义

- Vision Review **open required = 0**。
- **允许**：用户确认后，`/vision` 按 V-F069 执行 VP-016 有界组合层关门与索引同步。
- **禁止**：在无用户书面确认时把组合索引改成 VP-016 `closed`；把 Root `done` 或 Goal 审计原文冒充本轮代码核验；把 KMS / 热加载 / 第二套 dump / Admin 密钥页写成已交付。

### 响应（对 self 意见 · VRev-036 findings 闭合 · 2026-08-22）

| date | actor | summary |
|------|-------|---------|
| 2026-08-22 | `/vision` · 用户书面「审视 VP-016 是否满足关门条件（应该独立符合代码实现情况，并检验确保没有引入 bug）。是的话，关门 VP-016」 | **不回溯改写**原 verdict `pass` 与 finding 正文。**V-F069 → `fixed`**：VP-016 组合层确认 **有界 `closed`**（架构 A5）。关门记录含 exit 1–6 ↔ 本轮独立证据；residual 点名 `workspace-016` / `GOAL-001` / I-016-005 与 `admin.mfa` wrapping。`roadmap.md` / `workspaces.md` / Charter 关系节 / `reviews.md` / 区 `workspace.md` 原子同步（VR-039）。本 scope **0 open required、0 open recommended**。 |
