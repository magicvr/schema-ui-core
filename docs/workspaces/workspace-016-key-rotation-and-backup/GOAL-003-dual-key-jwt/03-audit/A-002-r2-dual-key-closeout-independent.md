---
id: A-002
doc: audit-entry
goal: GOAL-003-dual-key-jwt
status: recorded
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
auditor: grok-build (grok-4.6 · reasoning high)
---

> 编排器代贴留痕：本条由独立审计员（grok build headless `/audit` 流程）产出并按其要求代贴落盘；`source: independent` 保持不变，内容逐字未改（仅 Markdown 包裹与文件名）。项目级决策依据：`docs/architecture/independent-audit-execution.md`。

# A-002 · R2 JWT 双密钥实现 close-out 独立审计（2026-08-22）

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型 / scope**：close-out · GOAL-003 全部（workspace-016 纲领 R2「JWT 双密钥实现」：D-001 语义、`verifyAccess` 回退、composition 接线、单测/套件证据、越界核对）
- **verdict**：**pass**
- **工作区**：`workspace-016-key-rotation-and-backup`（Root `GOAL-001-key-rotation-and-backup`；`shared_materials_catalog: none`；本意见未把共享资料当证据）

本意见不修改 `status`/`progress`；响应由 `/govern` 处理。按编排器本轮书面要求，审计员**未落盘**；代贴至本文件并更新索引。

### 范围与区间

| 项 | 值 |
|----|----|
| 被审目标 | GOAL-003-dual-key-jwt（R2 JWT 双密钥实现） |
| 工作区 / Root | workspace-016-key-rotation-and-backup / GOAL-001-key-rotation-and-backup |
| 合同 | Root D-002（R1 配置面）；GOAL-003 D-001（I-003 裁定 + 实施方案） |
| 代码切片 | 工作树相对 HEAD `c96e963`（R1 checkpoint）：`apps/api/internal/auth/auth.go`、`auth_test.go`、`internal/composition/composition.go`。配置面 `AUTH_JWT_SECRET_PREVIOUS` / `ValidateProd` 为 R1 已冻结前置，本轮只消费、不重审 R1 关门。 |
| 不在本 scope | R3 轮换后恢复（I-004）；R4/R5 双路径证据；KMS / kid / refresh 模型变更 / 热加载 / Admin 密钥页 |
| P-005 | I-003 required、最晚阶段 = R2 接入前 → **verified**（D-001 + 实现与测试）。I-005 non-blocking、仍 collecting；D-001 按默认「previous 可验」实施，不阻断本关门。I-004 属 R3。 |
| 共享资料 | catalog = none；无固定引用被当作证据。 |

### 成果（有证据）

| 成果 | 证据 |
|------|------|
| I-003 关闭：重叠窗 = previous 配置存续期；不用 `kid`；refresh 为 opaque SHA-256 | D-001；Root `00-meta` I-003 行（工作树已改为 verified）；`SignAccessToken` 未设 `kid`（`auth.go:363-375`）；`Refresh` → `HashToken` / `RefreshTokenByHash`（`auth.go:256-293`）；`NewOpaqueToken`/`HashToken`（`auth.go:434-460`） |
| 签发只用 current | `issue()` 仍 `SignAccessToken(a.secret, …)`（`auth.go:338`）；`Login` / `Refresh` / `IssueTokensFor` 均走 `issue()`；测试子用例 3 |
| 校验 current → previous；两次都走 `ParseAccessToken`（强制 HMAC 方法 + `jwt.WithExpirationRequired()`） | `verifyAccess`（`auth.go:393-398`）；Middleware 改调 `verifyAccess`（`auth.go:571`）；`ParseAccessToken`（`auth.go:404-418`） |
| previous 空 = 单密钥（`len==0` 不回退）；既有 `New` / `NewWithRepository` 签名不变 | `verifyAccess` 短路径；`New`/`NewWithRepository` 仍不设 `previousSecret`；composition 在 previous 为空时传入 `[]byte("")`（len 0） |
| composition 接线；`NewApp` 对外签名不变 | `NewApp(cfg, secretValue, seedHash, logger)`（`composition.go:102`）未改；`newAuthenticator` 读 `cfg.AuthJWTSecretPrevious`（`composition.go:227-232`）；current 仍来自 `jwtSecret`（`resolveJWTSecret` 只解析 current，符合 D-001「previous 无需 main 层 fallback」） |
| 过期不延长寿命 | `ParseAccessToken` 两把钥匙都校验 `exp`；测试子用例 4 用 previous 签发 TTL=`-1s` 的 token，双密钥中间件仍 401 |
| 对外无验签状态 oracle | Middleware 在 `verifyAccess` 失败时固定写 `UNAUTHENTICATED` + `"invalid or expired access token"`（`auth.go:572-578`）；catalog 码 `UNAUTHENTICATED` 对非字段码丢弃 caller message（`errorcatalog.go:226-248`），jwt 内部错误不进信封 |
| 单测 4/4（本审计员独立复跑） | `go test ./internal/auth/ -run TestDualKeyRotationOverlapWindow -count=1 -v`：**PASS**（重叠窗通过 / 移除 previous 后 401 / 新签发只验 current / 过期 previous token 仍 401） |
| vet | 本审计员：`go vet ./...`（`apps/api`）**0 finding** |
| 越界未发生 | 工作树 diff 仅 auth 构造/验签 + composition 一处接线（另有 Root 台账 I-003 行与 R1 E-002 checkpoint 注记）；无 `kid`、无 refresh JWT、无密钥 setter/热加载、无 KMS；MFA 仍只用 current `jwtSecret`（`composition.go:355`） |

### 对照成功标准

| 标准 | 判定 | 证据 |
|------|------|------|
| 新签发只用 current | 达成 | `issue()`；测试 3 |
| 重叠窗内 previous 可验未过期 access | 达成 | `verifyAccess`；测试 1。重叠窗长度 = previous 配置存续 + 重启退役（D-001），无时钟窗状态 |
| 未配置 previous = 今日单密钥 | 达成（结构 + 既有构造器回归） | `len(previousSecret)==0` 短路径；`New`/`NewWithRepository` 未改；auth/handler/composition 包测试通过。生产空 previous 走新构造器，语义等价见 F-001 |
| 回退不延长过期 token | 达成 | 两把钥匙都强制 `exp`；测试 4；current 已过期时第二把因签名失败同样拒绝（HMAC 不同钥） |
| 重启生效、无热加载 | 达成（结构） | 密钥仅构造器注入；无 `SetPrevious`/运行期换钥 API；退役 = 去掉配置并重启 |
| 不越界 kid / refresh 模型 / 热加载 / KMS | 达成 | diff 与签发/refresh 路径核对 |
| Root 方向级 2–4（R3/R4/R5） | 不在本 scope | 不得用本目标关门冒充 |

### Findings

#### F-001 · composition 双密钥接线缺少自动化钉死（recommended · low · open）

- **严重度**：low  
- **建议**：recommended  
- **状态**：open  
- **描述**：生产路径 `newAuthenticator` 在 previous 非空时把 `cfg.AuthJWTSecretPrevious` 传入 `NewWithRepositoryAndPrevious`，源码正确、`NewApp` 签名未改。但没有任何测试调用 `newAuthenticator`，也没有任何 `NewApp`/`cfg.AuthJWTSecretPrevious` 用例覆盖「composition 层双密钥真的验过 previous」。D-001 矩阵把双密钥断言放在 auth 包，把单密钥回归交给全套件；因此当前切片**行为正确**，但以后若把 `newAuthenticator` 退回 `NewWithRepository`，auth 包四子用例仍会绿、重叠窗会在真实进程里静默消失。  
- **证据**：`composition.go:227-232`；`grep newAuthenticator` 于 `*_test.go` = 0 命中；`NewApp(` 测试均未设 `AuthJWTSecretPrevious`（`composition_test.go` / `postgres_startup_test.go` 等）；双密钥四子用例全部直接 `NewWithRepositoryAndPrevious` / `NewWithRepository`（`auth_test.go:575-616`）。「移除 previous」子用例用的是旧构造器，不是生产用的「新构造器 + 空 previous」（二者 `len==0` 等价，但钉的不是 composition 接线）。  
- **建议动作**：补一条 composition 或 auth 中间件测试：`cfg.AuthJWTSecretPrevious` 非空时，经 `newAuthenticator`/`NewApp` 装配的中间件接受 old-key token。非本关门必改。

#### F-002 · 执行索引未登记 E-002（recommended · low · open）

- **严重度**：low  
- **建议**：recommended  
- **状态**：open  
- **描述**：`02-execution/E-002-full-suite-verification.md` 已存在，但 `02-execution.md` 索引只列 E-001。事实文件在，索引不完整，不影响实现正确性。  
- **证据**：`GOAL-003-dual-key-jwt/02-execution.md` 执行索引表 vs `02-execution/E-002-full-suite-verification.md`。  
- **建议动作**：编排器响应时把 E-002 补进索引。

#### F-003 · E-002「`go test ./...` exit 0」本环境未能整包复现（recommended · low · open）

- **严重度**：low  
- **建议**：recommended  
- **状态**：open  
- **描述**：本审计员独立跑 `apps/api` 的 `go test ./... -count=1` **exit 1**。失败仅 `internal/store` 两条 Postgres 集成测试：`TestOpenPostgresProbeIntegration`、`TestPostgresMigrateRunnerIntegration`（`WasFresh() = false, want true`）。R2 diff **不包含** `internal/store`。JWT 相关包独立通过：`internal/auth`、`internal/config`、`internal/composition`、`internal/handler` 均为 `ok`；`go vet ./...` 0 finding。失败形态符合共享 probe DB 残留表导致 `WasFresh` 误判（测试只 DROP 了 `_r2_wasfresh_probe, r3_users, schema_migrations`）。因此：**不能把 E-002 的整包 exit 0 当作本会话已复现的事实**；也**不能**据此否定 R2 双密钥切片。  
- **证据**：本会话 `go test ./...` 输出；`store/postgres_test.go:501-502, 630-631`；`git diff --name-only` 相对 `c96e963` = `auth.go` / `auth_test.go` / `composition.go` + 两份 Root 台账。  
- **建议动作**：关门不阻断。若编排器要保留「全仓零回归」措辞，应在干净 PG probe 上重跑或把 E-002 收窄为「JWT 相关包 + vet」，并注明 store PG `WasFresh` 为环境残留、非本切片回归。

无 high/med required finding。无到期且影响本 scope 的 required 信息项。

### 必改项汇总（required）

空。

### 与既有意见的异同（A-001 self）

| 点 | A-001 (self) | 本意见 (independent) |
|----|----------------|----------------------|
| verdict | pass（待 independent） | **pass**（同意可关门，待编排器响应 recommended） |
| 必改 | 0 | 0 |
| 任意错误再试 previous | 书面接受（D-001 已选最简实现） | 同意；额外确认 jwt 错误不进 UNAUTHENTICATED 信封，不是用户状态 oracle |
| 双失败返回第二次错误 | 中间件信封不变 | 同意；`verifyAccess` 未导出，唯一生产调用点丢弃 `err` 正文 |
| 全仓零回归 | 引用 E-002 exit 0 | **不完全同意为已独立复现**：见 F-003；JWT 包与 vet 已复现 |
| composition | 「接线且 NewApp 签名不变」 | 同意源码；补充 F-001（缺 composition 级双密钥测试） |
| I-003 / 越界 | 达成 | 同意 |

无 verdict 冲突，无对同一必改项的一要一否。P-004 冲突裁决点不触发。

### 核对审计重点（摘要）

1. **与 D-001/D-002 合同一致性**：签发只用 current；校验 current 再 previous；重叠窗 = 配置存续期；previous 受 R1 `ValidateProd` 同强度/同值守卫（本轮只消费）。符合。  
2. **过期延长 / 状态 oracle**：两把钥匙都强制 `exp`；catalog 化 401 不泄漏哪把钥匙或 jwt 原因。previous 配置时无效 token 会多一次本地 HMAC（D-001 已选），不是账户状态 oracle。  
3. **未配置 previous**：`len==0` 即 `ParseAccessToken(current)`；既有构造器零迁移。  
4. **composition**：`NewApp` 签名未变；previous 从已注入的 `cfg` 读取；无 main 层遗漏。缺口仅测试钉死（F-001）。  
5. **测试可核对**：四子用例本会话 4/4 PASS；vet 0；整包 exit 0 未复现（F-003）。  
6. **越界**：无 kid、无 refresh JWT、无热加载、无 KMS。MFA 派生密钥仍绑 current，属既有 GOAL-017 耦合，不在 VP-016 分母；本切片正确地没有把 previous 传进 MFA。

### 结论 + 建议给编排器/用户的下一步

R2 实施切片与 D-001/D-002 一致，证据可核对，无未闭合 required finding，I-003 门禁已关闭。**independent verdict = pass。** GOAL-003 在编排器合并响应 A-001+A-002、处理 recommended（或不处理，因非必改）并同步 `goal-tree` 之后可以 `done`。

建议 `/govern` 下一句：

> 响应 GOAL-003 A-001（self pass）+ A-002（independent pass）：required=0，无冲突；F-001/F-002/F-003 均为 recommended。按 pass 关门 GOAL-003，更新 goal-tree 与 Root 路线图 R2=完成（progress 2/5）；Root `01-decision.md` I-003 镜像仍为 collecting、`workspace.md` 仍写 R2 未开始，一并改到与 `00-meta` 一致。F-001 可留待后续或本回合补 composition 钉死；F-002 补执行索引；F-003 不要写成「本切片引入 store 回归」。然后进入 R3（I-004）或先登记 R3 方案。

### 声明

本意见不修改 status/progress；响应由 `/govern` 处理。独立审计员未改代码、未改 goal-tree、未写入 `03-audit`（按本轮「不要写入任何文件」执行）。
