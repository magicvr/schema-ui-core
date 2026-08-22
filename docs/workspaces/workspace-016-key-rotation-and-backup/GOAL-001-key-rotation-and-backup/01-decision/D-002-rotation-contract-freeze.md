---
id: D-002
doc: decision-entry
goal: GOAL-001-key-rotation-and-backup
status: accepted
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
---

# D-002 · R1 轮换合同与配置面冻结（关闭 I-001 / I-002）

## 背景

R1 门禁要求在方案冻结前以证据关闭 I-001（current/previous 键名、生产 fail-closed、熵规则）与 I-002（本波密钥集合是否仅 JWT）。2026-08-22 对 `apps/api` 的代码核对结论如下（证据行号以当日 HEAD `5195104` 为准）：

| 证据点 | 位置 | 结论 |
|--------|------|------|
| YAML 键 `auth.jwt_secret` | `internal/config/config.go:202` | current 键沿用现名 |
| env 覆盖 `AUTH_JWT_SECRET` | `internal/config/config.go:496` | env 名沿用 |
| 生产熵规则 | `internal/config/config.go:953-964` + `minJWTSecretLen=32`（`:1210`）、`containsLettersAndDigits`（`:1214`） | 非开发环境：≥32 字符且同时含字母与数字 |
| 启动 fail-closed | `cmd/server/main.go:74-85`（`resolveJWTSecret`） | 非开发环境缺 current 直接启动失败；development 允许文档化不安全 dev key |
| 服务凭证密钥来源 | `internal/auth/auth.go:411-420`（`NewServiceCredentialToken` → `NewOpaqueToken`）+ `:57`（前缀 `sui_sc_`） | 256-bit CSPRNG 随机值，落库 SHA-256（`HashToken`），**与 JWT secret 无派生关系** |
| refresh 令牌形态 | `internal/auth/auth.go:400-409` + `issue()`（`:306-331`） | opaque 256-bit + SHA-256 落库，非 JWT |

## 决策

### 1. 配置面键名（I-001 · closed）

| 角色 | YAML | env | 说明 |
|------|------|-----|------|
| current | `auth.jwt_secret`（不变） | `AUTH_JWT_SECRET`（不变） | 唯一签发密钥 |
| previous | `auth.jwt_secret_previous`（新增） | `AUTH_JWT_SECRET_PREVIOUS`（新增） | 重叠窗校验密钥；缺省 = 单密钥模式 |

- **缺省单密钥**：previous 未配置（空）时，任何环境下行为与今日完全相同；轮换不是 mvp/dev/production 启动硬依赖。
- **生产 fail-closed 与熵规则**：非开发环境下，已配置的 previous 沿用与 current 完全相同的校验（≥32 字符、同时含字母与数字）。理由：previous 同样能验证签名，弱 previous 等于弱化验签面；统一规则避免"双标"。弱历史 key 若过不了该规则，操作者应先做一次合规轮换再进入双密钥态（fail-closed 取向与 GOAL-009 A-002 F-002-005 一致）。
- **同值守卫**：previous == current（非空相等）为配置错误，启动失败。防止"假轮换"造成重叠窗语义混乱。
- **secret 不入库、不进日志**：错误信息只点名键名，不携带值（沿既有 `ValidateProd` / S3 凭据错误先例）。
- **生效方式**：进程重启生效；无热加载（VP 冻结）。

### 2. 本波密钥集合 = 仅 JWT 应用签名密钥（I-002 · closed）

- 服务凭证 token 为 CSPRNG 随机值的 SHA-256 opaque hash，不与 JWT secret 共用、不受其轮换影响——**书面出局**，不进本轮换分母。
- refresh token 同为 opaque hash，不受签名密钥轮换影响（与 VP I-016-003 预期一致，R2 决策再正式引用）。
- `DB_PASSWORD`、S3 access/secret、`ADMIN_INITIAL_PASSWORD` 按 VP 边界出局。

### 3. R1 交付切片

R1 子目标（GOAL-002）只交付**配置面**：Config 字段、env/YAML 双通道解析、ValidateProd 规则、单元测试。Authenticator 双密钥验签消费属 R2（GOAL-003），本波不在 composition 接死代码。

### 4. 审计模式

按 D-001 §5：R1 合同冻结走 **self**；R2 生产路径实施走 **independent**（grok build · grok-4.6 · `/audit`）。

## 为什么

- 沿用现名可保证既有部署零迁移成本；新增 previous 键是纯增量。
- 统一熵规则使"哪把钥匙要强"没有歧义，审计面最小。
- 密钥集合书面出局服务凭证，回应 VP I-016-002 的条件项（若共用须纳入）；证据显示不共用，故出局成立。

## 未选方案

- previous 弱熵豁免（便于弱 key 逃逸）：拒绝——弱 previous 仍可伪造验签通过，风险大于迁移便利。
- 引入 JWT `kid` 头：留待 R2 决策（两把候选密钥的 HMAC 试验序即可满足本波语义），此处不预支。
- 时间窗式重叠（TTL 内可验）：拒绝——本波重叠窗 = previous 配置存续期，由操作者控制下线时机，不新增状态。
- 把服务凭证 hash 纳入轮换分母：证据显示不共用，纳入会违反 VP 退出分母冻结。
