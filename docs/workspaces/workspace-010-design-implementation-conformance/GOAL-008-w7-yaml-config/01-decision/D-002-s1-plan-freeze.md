---
id: D-002
goal: GOAL-008-w7-yaml-config
status: accepted
date: 2026-08-14
scope: S1 方案冻结
parent: GOAL-008-w7-yaml-config
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# D-002 · S1 方案冻结：YAML 主配置体系

## 1. 背景与目标

现有 apps/api/config 为纯 env 读取（14 项），handler/upload.go 另读 3 项 upload env；敏感项（AUTH_JWT_SECRET、ADMIN_INITIAL_PASSWORD）靠 .env/进程 env。目标：

1. **YAML 为配置权威**（用户裁决：主流系统用 yaml；env 只放敏感信息，两者区分开）；
2. **现有部署零迁移**：已设置 env 继续优先；
3. 主配置集中、可读、可注释、可版本化。

## 2. 优先级与插值规则

加载优先级（高 → 低）：

1. **进程 env**（已设置时覆盖一切）；
2. **CONFIG_FILE 指定的 YAML**（默认 `configs/config.yaml`）；
3. **内置默认**（go:embed `configs/config.yaml` 作为 fallback 默认）——仅当 CONFIG_FILE 不存在时使用；
4. 字段级内置默认（代码内）。

插值：YAML 值支持 `${VAR}` 与 `${VAR:-default}`：
- `${VAR}`：VAR 未定义 → **fail-closed**（加载报错，进程拒绝启动）；
- `${VAR:-default}`：VAR 未定义 → 用 default。
- 与 compose.yaml 现有 :? / :- 语法同族（compose 先例，见 compose.yaml 的 AUTH_JWT_SECRET / APP_PROFILE）。

敏感字段（AUTH_JWT_SECRET、ADMIN_INITIAL_PASSWORD）在 YAML 中**只写引用**（`${AUTH_JWT_SECRET}`），真实值在 env（进程 env / .env / secret store）。

## 3. 配置字段清单（YAML schema）

全部 17 项迁入 YAML（app / http / auth / db / log / upload 小节）：

| YAML 路径 | 原 env | 敏感 | 默认 |
|-----------|--------|------|------|
| app.name | APP_NAME | 否 | schema-ui-core |
| app.env | APP_ENV | 否 | dev |
| app.profile | APP_PROFILE | 否 | mvp |
| app.modules_enabled | APP_MODULES_ENABLED | 否 | （空 = 全量） |
| http.addr | HTTP_ADDR | 否 | :8080 |
| http.read_timeout | HTTP_READ_TIMEOUT | 否 | 15s |
| http.write_timeout | HTTP_WRITE_TIMEOUT | 否 | 15s |
| http.idle_timeout | HTTP_IDLE_TIMEOUT | 否 | 60s |
| log.level | LOG_LEVEL | 否 | info |
| auth.jwt_secret | AUTH_JWT_SECRET | **是** | `${AUTH_JWT_SECRET}`（fail-closed） |
| auth.access_ttl | AUTH_ACCESS_TTL | 否 | 15m |
| auth.refresh_ttl | AUTH_REFRESH_TTL | 否 | 720h |
| auth.dev_session_enabled | AUTH_DEV_SESSION_ENABLED | 否 | false |
| db.path | DB_PATH | 否 | data/schema-ui.db |
| admin.initial_password | ADMIN_INITIAL_PASSWORD | **是** | `${ADMIN_INITIAL_PASSWORD}`（fail-closed） |
| upload.allowed_types | UPLOAD_ALLOWED_TYPES | 否 | （现有默认） |
| upload.max_files_per_user | UPLOAD_MAX_FILES_PER_USER | 否 | （现有默认） |
| upload.max_bytes_per_user | UPLOAD_MAX_BYTES_PER_USER | 否 | （现有默认） |

枚举（app.env: dev/test/prod；app.profile: mvp/admin/dev；log.level: debug/info/warn/error）沿用现状，不新增。

## 4. env 衔接（.env 与 CONFIG_ENV_FILE）

- YAML 是权威；.env **不是**配置源，只承载敏感值（AUTH_JWT_SECRET=...、ADMIN_INITIAL_PASSWORD=...），供 CONFIG_ENV_FILE 指向。
- CONFIG_ENV_FILE（可选 env，默认 configs/.env）：启动时读取（simple KEY=VALUE 解析），**不覆盖**已存在的进程 env（进程 env 优先）。
- compose.yaml 已有 .env + :? 插值，与本方案同族；compose 场景无需 CONFIG_ENV_FILE（env 由 compose 注入）。
- 真实 secret 仍走进程 env / secret store（生产）；.env 为开发便利且 gitignored（现状即 gitignored）。

## 5. 代码形态

- config.Load() 改为分层加载：显式路径参数（CONFIG_FILE env 或 flag）→ YAML 解析（gopkg.in/yaml.v3）→ 插值 → env 覆盖 → 默认兜底。
- Config 新增 UploadAllowedTypes / UploadMaxFilesPerUser / UploadMaxBytesPerUser（或等效小节结构）；handler/upload.go 的 3 个 os.Getenv 改为经 Config 注入。
- RegisterUpload 增加变参 UploadOption（函数式选项，默认回退现行为），避免破坏现有调用点。
- 保留 ValidateProd() / LogLevel() 等现有方法签名，调用点不迁移。

## 6. 未选方案与排除项

| 方案 | 未选原因 |
|------|----------|
| 纯 env 维持现状 | 用户裁决：yaml 主流；env 只放敏感 |
| 全部塞 env + 嵌套键名 | 与 yaml 决策冲突 |
| vault/secret manager 集成 | 本轮范围外（进程 env 已是标准 secret 通道） |
| 运行时热加载/重载 | 用户裁决仅重启生效（GOAL-013 同理） |
| Protocol / Profile / 模块矩阵（module-matrix）改造 | 明确排除，D-001 已定 |
| 装配语义（包注册顺序/Assembly 顺序） | 明确排除，D-001 已定 |

## 7. 影响面

- apps/api/config（Load + Config 结构 + 新字段）
- apps/api/handler/upload.go（3 env → Config）
- apps/api/composition（RegisterUpload 变参；无装配语义变化）
- apps/api/main.go（resolveJWTSecret 已 fail-closed，保持）
- apps/api/configs/config.yaml（新增，embed）+ configs/.env.example（新增，gitignored 的 .env 模板）
- deploy/compose.yaml 同步（env 注释/示例与 config.yaml 对应）
- **go（VP-008）判定**：配置载体变化，非装配语义/非门禁语义 → 不 held。
