---
id: GOAL-008-w7-yaml-config
doc: audit-entry
record_id: A-003
source: independent
scope: S1–S4 全量记录 + S5 关门前独立审计（D-001/D-002/A-001/E-001～E-005/A-002 与 as-built）
verdict: pass
status: recorded
auditor: grok build（Grok 4.6）
parent: GOAL-008-w7-yaml-config
created: 2026-08-14
updated: 2026-08-14
version: 0.1.0
---

# A-003 · 独立审计 · S5 关门前交叉意见（2026-08-14）

- **source**：independent
- **auditor**：grok build（Grok 4.6）
- **类型** / **scope**：close-out · S1–S4 方案/实施/验证/go 判定 + 实现对照
- **verdict**：pass
- **工作区**：`workspace-010-design-implementation-conformance`（Root `GOAL-001-design-implementation-conformance`；`canonical_scope` 已校验；`shared_materials_catalog: none`；`primary_plan` = VP-010）

## 范围与区间

- **covered**：00-meta 路线图 S1–S4；D-001/D-002；E-001～E-005；A-001/A-002；I-001～I-004；`apps/api/internal/config/{config.go,config.default.yaml,config_test.go}`；`apps/api/configs/{config.yaml,.env.example}`；`apps/api/internal/handler/upload.go`；`apps/api/internal/composition/composition.go`；`apps/api/cmd/server/main.go`；`apps/api/Dockerfile`；根 `compose.yaml`。本轮复跑 `go test ./internal/config/ ./internal/handler/ ./internal/composition/`（exit 0）。另用一次性探针（跑完即删，不入库）核对空 YAML、多文档、省略键、注释插值、`${VAR:-}` 空默认、畸形 `.env`、生产空 JWT。
- **excluded**：未复跑 E-004 四路径真实进程启动；未改 status/progress/goal-tree；未读其他工作区正文（E-005 对外部判例的引用只记为该条目主张）。
- **P-005**：I-001～I-004 均 closed，无到期未闭合 required 信息项。

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| 分层加载存在且进程 env 覆盖 YAML | `config.go:99-107,195-213`；`config_test.go:144-221` |
| 显式 `CONFIG_FILE` / `CONFIG_ENV_FILE` 缺失 fail-closed | `config.go:138-142,246-249`；`config_test.go:223-229,338-344` |
| 默认路径缺失回退 embed | `config.go:145-151`；`config.default.yaml` `//go:embed`；镜像不拷贝 YAML（`apps/api/Dockerfile:31`） |
| 裸 `${VAR}` 未定义 fail-closed；`${VAR:-default}` 可用 | `config.go:301-321`；`config_test.go:231-263` |
| `${VAR:-}` 空默认合法；开发过 ValidateProd、生产被 32 字符门禁拦住 | 探针；`ValidateProd` `config.go:349-360`；`config.default.yaml:42,56` |
| 整行注释 / 行内 ` #` 中的 `${…}` 不计入引用 | `config.go:294-300`；探针（文档式示例与 ` # ${MUST_NOT_COUNT}` 均不报错） |
| `CONFIG_ENV_FILE` 不覆盖进程 env；畸形行 fail-closed | `config.go:269-270,261-265`；`config_test.go:316-336`；探针 |
| KnownFields 未知键 fail-closed | `config.go:162-166`；`config_test.go:282-292` |
| LoadError 经 ValidateProd 暴露；main 启动必调 | `config.go:339-342`；`config_test.go:346-355`；`main.go:30-33` |
| 生产 JWT 无绕过：未设 APP_ENV 拒绝；空串 APP_ENV 保持并拒绝 | `config.go:346-348`；`config.yaml:34`；探针 |
| upload 默认 1000 / 256MiB / 无类型限制；composition 注入 cfg | `upload.go:228-235`；`composition.go:208-211`；YAML `268435456` |
| compose 镜像走 embed + env；端口 :25080 一致 | `Dockerfile:31-36`；`compose.yaml:4-7,25-43`；YAML/字段默认 `:25080` |
| 包级单测本轮复跑绿 | `go test ./internal/config/ ./internal/handler/ ./internal/composition/ -count=1` exit 0 |
| I-001～I-004 关闭证据可指回 | D-002 §3 / §2 / §1；E-005 表 |

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| S1 方案冻结（17 项 + 优先级 + 插值 + 迁移） | 达成（表内默认值有漂移，见 F-001） | D-002；A-001 |
| S2 YAML 权威 + embed + env 覆盖 + upload 迁入 | 达成 | E-003；df5d440；上表实现行 |
| S3 单测 + 四路径实测 + 回归 | 单测本轮复跑达成；四路径采信 E-004（未独立复现进程） | E-004；本轮 `go test` |
| S4 go 不 held | 同意 E-005：未改 Profile/模块矩阵/Manifest/装配顺序 | E-005；`composition.go` 仅 RegisterUpload 变参注入 |
| 零迁移（旧 env-only） | 达成（进程 env 后覆盖；敏感项用 `${VAR:-}` 避免开发期插值拒绝） | `config.go:195-213`；`jwt_secret: ${AUTH_JWT_SECRET:-}` |
| I-001～I-004 门禁 | 达成 | 00-meta 信息表；对应 D/E |

### 审计重点逐项

1. **优先级链 vs D-002**：高→低为进程 env → `CONFIG_FILE`（默认 `configs/config.yaml`）→ embed → 字段默认。相对 D-002 的细化成立且合理：显式 `CONFIG_FILE` 缺失 fail-closed；未设置且默认路径不存在才 embed。env-only 部署因后覆盖而零迁移成立。字段默认对**省略的字符串键**未兑现（F-002）。
2. **插值 fail-closed**：裸 `${VAR}` 未定义拒绝启动。`${VAR:-}` 空默认合法。整行 `#` 与行内 ` #` 排除生效，仓库内文档示例不会被误当引用。产品单测未覆盖注释排除（F-003）。引号内 ` #` 会被误切（F-003）。
3. **敏感项**：默认 `${AUTH_JWT_SECRET:-}` / `${ADMIN_INITIAL_PASSWORD:-}`。开发：`ValidateProd` 放行空 JWT，`resolveJWTSecret`/`resolveSeedHash` 走文档化 dev 回退。生产：`ValidateProd` 先拦 JWT 长度/字符类；空密码由 `main.resolveSeedHash` 拦（非 ValidateProd，但启动链连续，无绕过）。`APP_ENV==""` 拒绝猜测。`APP_ENV=dev`（非 `development`）走生产 JWT 规则，偏严不偏松。未见生产绕过路径。
4. **CONFIG_ENV_FILE**：`LookupEnv` 已存在则 skip；显式路径缺失 / 默认路径非 NotExist IO 错 / `KEY=VALUE` 解析失败均 fail-closed。默认 `configs/.env` 被 `apps/api/.gitignore` `.env` 忽略。
5. **KnownFields / 空文件 / 多文档**：未知键 fail-closed。空文件与仅注释文件是合法 YAML null，Decode 得 `EOF` → LoadError（偏严 fail-closed）。多文档只读第一份；第二份含未知键被忽略（F-005）。
6. **upload**：无 option 时 1000 / 256MiB / 空 allow-list；`composition` 传入 cfg；handler 不再 `os.Getenv("UPLOAD_*")`；quota 与危险类型硬拒绝路径仍在。旧调用点（如 filelibrary 测试）仍可零 option 运行。
7. **compose**：注释与实现一致。镜像仅拷贝二进制，WORKDIR `/app` 无 `configs/config.yaml` → embed；compose/Dockerfile 注入 `APP_ENV=production` 与 fail-closed secrets；`HTTP_ADDR` 不注入，embed `:25080` 对齐 `EXPOSE 25080` 与端口映射。
8. **边界**：HTTP 默认 `:25080` 在字段默认 / 两份 YAML / compose / Dockerfile 一致（D-002 表误写 `:8080`，见 F-001）。`DB_PATH` YAML 为 `./data/schema-ui.db`（相对 cwd）；容器用绝对 `/app/data/schema-ui.db`。`APP_ENV` 空串显式保持。

## Findings

### F-001 · D-002 §3 冻结默认值表与 as-built / 旧默认不一致

| 字段 | 值 |
|------|-----|
| level | recommended（non-blocking） |
| severity | med |
| status | open |
| evidence | D-002:47-62（`:8080`、15s/15s、`APP_ENV` 默认 `dev`、枚举 `dev/test/prod`、`APP_NAME=schema-ui-core`、敏感项裸 `${AUTH_JWT_SECRET}`）；as-built `config.go:110-126`、`config.default.yaml:19,26-28,42,56`（`:25080`、5s/10s、`env: ""`、`schema-ui-core-api`、`${VAR:-}`）；A-001 未记录该漂移 |

实现**正确**优先了「值语义不变 / 零迁移」（若按 D-002 表做裸 `${VAR}`，未设 `AUTH_JWT_SECRET` 的开发启动会在插值层被拒）。风险是后续按冻结表「回修」代码。关门前建议把 D-002 §3 修订为 as-built 合同，不必改代码。

### F-002 · 省略的 YAML 字符串键会清零字段默认（D-002 第 4 层未兑现）

| 字段 | 值 |
|------|-----|
| level | recommended（non-blocking） |
| severity | med |
| status | open |
| evidence | `config.go:170-182` 无条件 `cfg.HTTPAddr = yf.HTTP.Addr` 等；本轮探针：仅含 `app.env` 的 YAML → `HTTPAddr=""``DBPath=""``AppName=""` |

已交付的 `configs/config.yaml` 与 embed 是全量键，compose/env-only 路径不受影响。操作员若把 `CONFIG_FILE` 写成「只覆盖几项」的 overlay，会丢掉 `:25080` / DB 路径。时长字段走 `orDuration`、upload 整数走 `> 0` 才保留默认，字符串层不一致。建议：省略键保留字段默认，或文档写明「必须从模板复制全量键」。

### F-003 · 注释排除行为成立，但 A-002 所称单测不存在；` #` 切割过宽

| 字段 | 值 |
|------|-----|
| level | recommended（non-blocking） |
| severity | low |
| status | open |
| evidence | A-002:24「有测试覆盖文档式示例」；`config_test.go` 无注释排除 / `${VAR:-}` 空默认用例。`config.go:298-299` `strings.Index(line, " #")`。探针：文档式 `# ${VAR}` 与行内 ` # ${MUST_NOT_COUNT}` 不 fail-closed；`name: "My App #1"` → `parse config YAML` LoadError |

行为对仓库内模板是安全的。建议补测并避免把引号内 ` #` 当注释。

### F-004 · 操作文档未同步「API 会读 configs/.env」

| 字段 | 值 |
|------|-----|
| level | recommended（non-blocking） |
| severity | low |
| status | open |
| evidence | `apps/api/README.md:29`「Go API 不会自动加载 `.env`」；`QUICKSTART.md:24` 仍只谈 `apps/api/.env.example`。对照 `config.go:236-245` 默认读取 `configs/.env`。`apps/api/.env.example:7` 已正确写明不加载 **apps/api/.env** |

D-001 含文档同步。不破坏零迁移（export / compose 仍有效），但会让人找不到新的 `configs/.env` 通道。建议改 README/QUICKSTART 指向 `CONFIG_ENV_FILE` / `configs/.env.example`。

### F-005 · 空文件与多文档流的 KnownFields 边界

| 字段 | 值 |
|------|-----|
| level | recommended（non-blocking） |
| severity | low |
| status | open |
| evidence | 探针：空文件 / 仅注释 → `parse config YAML: EOF`；`app.name: first-doc` + `---` + 含 `bogus_key` 的第二文档 → LoadError=nil 且 `AppName=first-doc`（第二文档未知键未触发 KnownFields） |

空文件被拒是偏严 fail-closed（应用「省略文件 → embed」，不要交空文件）。多文档非本项目模板形态。建议在 YAML 头注释写明「单文档、非空」；可选：Decode 第二文档若非 EOF 则 fail-closed。

## 必改项汇总

**无 required / 必改项。**

开放项均为 recommended（non-blocking）。按 P-003，不阻断将本目标标为 `done`。

## 与既有意见的异同

| 既有 | 本意见 |
|------|--------|
| A-001 self · S1 pass、无 finding | 同意方案可实施；补记 D-002 §3 表与旧默认不一致（F-001），A-001 未捕获 |
| A-002 self · S2/S3 pass、无 required；备注双份 YAML | 同意主路径与 fail-closed 三例；不同意「注释排除已有测试」（F-003）；双份 YAML 维护风险仍在，不另开 finding |
| A-002 / E-005「零调用点改动」 | 变参默认兼容成立，但 `composition.go:208-211` **应当且已经**传入 cfg（否则 YAML upload 不生效）。表述略满，不是缺陷 |
| E-004 四路径进程实测 | 未独立复跑；包级单测本轮绿。不因此降为 fail |

## 结论 + 建议给编排器/用户的下一步

**verdict: pass。无 required，可以关门。**

建议 `/govern` 响应本意见：F-001～F-005 可 `accepted-residual`（范围=文档/测试/overlay 卫生，不改启动语义；复审触发=下次改默认值或开放 overlay YAML）或顺手修订 D-002 表 + README/QUICKSTART + 补注释/空默认单测。不修代码也可合法闭合 recommended。闭合后勾 S5、`status: done`、同步 goal-tree（及 `workspace.md` 波次行仍写 0/5 的过期指针）。

go 判定维持 E-005：**不 held**。

## 声明

本意见不修改 status/progress/检查点/方案正文/goal-tree 状态列。响应与关门由 `/govern` 处理。
