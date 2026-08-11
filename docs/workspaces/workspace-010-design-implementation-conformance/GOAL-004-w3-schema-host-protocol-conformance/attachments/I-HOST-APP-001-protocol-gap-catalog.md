---
id: I-HOST-APP-001
title: Host/App 协议缺口与业务候选目录
status: active
created: 2026-08-12
updated: 2026-08-12
parent: GOAL-004-w3-schema-host-protocol-conformance
version: 0.1.0
---

# I-HOST-APP-001 · Host/App 协议缺口与业务候选目录

## 1. 用途与边界

本附件是上游 Schema-UI 协议增补的**候选输入**，不是协议权威，也不授权本仓提前实现。它同时覆盖：

- 当前 `apps/api` / `apps/web` 已使用、但 2.7.0 Host/App 契约未明确覆盖的能力；
- 现行协议已覆盖、实现却可能偏离的能力；
- 企业后台与业务应用中确实可能出现、当前产品尚未使用但应预留稳定协议位置的能力。

每个候选在 S2 必须得到一种处置：

| 处置 | 含义 |
|------|------|
| `adopt-now` | 本次协议定义完整数据形状、状态/错误、安全边界、能力声明和 fixtures |
| `reserve-extension` | 本次定义稳定 capability/extension point 与 fail-closed 行为，具体载荷后续版本补齐 |
| `explicitly-out` | 协议明确说明不负责，并指出所属层；不是静默遗漏 |

优先级只表示建议进入协议设计的顺序：`P0` 为当前整改或 Host 基线所需，`P1` 为常见业务能力，`P2` 为合理扩展候选。所有 P1/P2 仍须在 S2 得到处置，但不自动要求同一版本全部成为 mandatory capability。

## 2. 初步分类结论

### 2.1 已有协议覆盖，不应另造 UI 语义

| 能力 | 当前依据 | 初步结论 |
|------|----------|----------|
| App manifest 版本、required capabilities、pages、navigation | `docs/schemas/app-manifest.schema.json`；Web `validateAppManifest()` | 继续使用协议；未知 capability/version fail closed |
| 页面节点白名单与递归渲染 | `apps/web/src/renderer/render.ts` 的 `WHITELISTED_NODE_TYPES` / `RENDER_UNKNOWN_NODE_TYPE` | 未发现生产渲染路径绕过未知节点校验 |
| `form`、`table`、`recordView`、`actionButton` 等注册节点 | `docs/schemas/component-registry.json` 与 renderer parser | 注入 Form/Table renderer 是合法 Host adapter，不等于私有 DSL |
| Node loading/empty/error state | `docs/schemas/node.schema.json` | 页面局部状态已有协议；不要与全局启动/认证错误混为一谈 |
| action 的 navigate/reload/toast/modal/upload | `docs/schemas/action.schema.json` | 已覆盖部分命令；全局生命周期和 download/job 仍需增补 |
| route path/query mapping | page schema 与 component registry | 深链和 request mapping 可复用；Shell 级 404/redirect/恢复意图仍缺 |

### 2.2 实现整改候选

这些项目只有在 S3 固定的新协议确认语义后才进入 S4；当前不直接改代码。

| ID | 现象 | 证据 | 初步判断 | 新协议需先明确 |
|----|------|------|----------|----------------|
| IMP-001 | Settings PATCH 只解码匿名结构，未显式拒绝未知 JSON 字段 | `apps/api/internal/handler/settings.go:165-188` | 若 action/request object 是闭集，则为实现偏离 | request object 的 unknown-key、额外 body 字段和服务端 validator 规则 |
| IMP-002 | Users 导航文案同时来自 provider `Label: "Users"` 与 manifest `label`/`labelKey` | `apps/api/internal/modules/users/provider.go:85-93`；`manifest/fragment.json:21-30` | 单一语义源与 i18n precedence 疑似被破坏 | navigation projection 的 `label`/`labelKey` precedence 与 provider 可否覆写 |
| IMP-003 | Host 侧字段/状态若只靠 TypeScript interface 或 handler struct 约束 | auth、branding、session 等当前实现 | 在协议补齐后会成为 conformance 修复项 | schema validation 放在 API、Web bootstrap 和 fixtures 的最低要求 |
| IMP-004 | 争议中的 row selection → detail/drawer 交互可能依赖本地约定 | `recordView`/table Host integration 待逐项核验 | 未裁定；可能是实现偏离，也可能是协议缺口 | row context、selection binding、detail target、drawer/modal ownership |

### 2.3 合法 Host 边界，但协议不足

- `LoginPage` 可以是 Host 认证壳，不要求伪装成普通业务 page；缺的是 auth discovery、challenge、session 和失败恢复契约。
- `/api/branding` 可以是公开 bootstrap 配置，不要求经 `/api/schema/{pageId}`；缺的是其 schema、缓存、安全和变更通知契约。
- `SchemaTable`、Form component、data/action adapters 可以由 Host 注入；缺的是 adapter capability、失败语义和扩展治理，不是要求消灭 adapter。
- `/api/accounts/me`、login/refresh/logout 可以是 Host session API；缺的是协议化状态机、token/cookie transport profile 和错误码。

## 3. Host/App 协议候选总表

### A. Auth、Identity 与 Session

| ID | 候选能力 | 真实业务场景 | 当前协议状态 | 最小协议内容 | 优先级 |
|----|----------|--------------|--------------|--------------|--------|
| AUTH-001 | 认证方式 discovery 与登录入口 | 密码、SSO、无密码登录并存 | absent | provider 列表、display key、启动方式、能力 ID、可用/禁用原因 | P0 |
| AUTH-002 | 会话 bootstrap / current principal | 启动时判断 anonymous/authenticated/locked/degraded | absent | principal、session state、expiry、roles/features provenance、缓存规则 | P0 |
| AUTH-003 | Access/refresh 生命周期 | token 轮换、过期、撤销、并发 401、auth-lost | absent | transport profile、刷新互斥、rotation/replay、clock skew、错误状态机 | P0 |
| AUTH-004 | 登录、登出与撤销 | 单设备登出、全设备登出、管理员撤销 | absent | command shape、幂等性、revoke scope、success/error、post-action | P0 |
| AUTH-005 | 认证后恢复意图 | 登录后回到原 deep link，并保留安全 query | absent | return target、allowlist、expiry、loop prevention | P0 |
| AUTH-006 | MFA / step-up challenge | 敏感操作二次认证、TOTP/WebAuthn/恢复码 | absent | challenge types、attempt state、resume action、cancel/expiry、a11y | P1 |
| AUTH-007 | OIDC/SAML/企业 SSO | 多 IdP、redirect/callback、账号绑定 | absent | redirect contract、state/nonce ownership、callback result、safe return URL | P1 |
| AUTH-008 | 密码与账号恢复 | 忘记密码、强制改密、锁定、邀请激活 | absent | flow states、one-time token handling、policy hints、privacy-safe errors | P1 |
| AUTH-009 | 会话/设备管理 | 查看活动会话、撤销指定设备 | absent | session list projection、device metadata、revoke action、审计引用 | P1 |
| AUTH-010 | 委托/模拟身份 | 客服或管理员临时切换身份 | absent | original/effective principal、明显横幅、时限、退出、审计要求 | P2 |
| AUTH-011 | Cookie/Bearer/CSRF transport profile | 浏览器 cookie 与 API token 部署差异 | absent | credential mode、CSRF、SameSite、storage 禁止项、CORS/security profile | P0 |

### B. Bootstrap、Discovery 与兼容协商

| ID | 候选能力 | 真实业务场景 | 当前协议状态 | 最小协议内容 | 优先级 |
|----|----------|--------------|--------------|--------------|--------|
| BOOT-001 | 确定性启动顺序 | branding、manifest、session、locale 相互依赖 | absent | 阶段 DAG、public/authenticated 边界、ready 判定、取消/重试 | P0 |
| BOOT-002 | Manifest discovery | 多部署 base path、CDN、版本化入口 | partial | well-known 发现、base URL、content type、redirect 与同源策略 | P0 |
| BOOT-003 | 版本与 capability 协商 | Host 与文档版本不一致 | partial | host capabilities、document requirements、partial support、fail-closed error | P0 |
| BOOT-004 | 工件 provenance / integrity | 固定 schema、registry、fixtures，避免漂移 | partial | source/version/commit/hash、bundle identity、验证失败行为 | P0 |
| BOOT-005 | 缓存与配置失效 | ETag、长会话中 branding/manifest 更新 | absent | ETag/version、TTL、stale policy、invalidation event、原子切换 | P0 |
| BOOT-006 | 环境与构建身份 | about/support、故障定位、灰度发布 | absent | app/build/protocol/environment/region 标识；明确不可含 secret | P1 |
| BOOT-007 | maintenance / upgrade-required / degraded | 后端维护、协议强制升级、部分依赖故障 | absent | boot terminal states、retryAt、read-only、support/correlation data | P0 |
| BOOT-008 | tenant/workspace preselection | 登录后先选组织或工作区再加载 manifest | absent | selection state、allowed contexts、default、switch invalidation | P1 |
| BOOT-009 | feature availability provenance | 灰度、许可、租户能力导致页面可用性不同 | needs clarification | flag source、stability、evaluation side、不得替代 permission | P1 |
| BOOT-010 | offline/reconnect policy | 临时断网、PWA/桌面壳恢复 | absent | offline state、cached-read policy、mutation prohibition/queue、reconnect | P2 |

### C. Branding 与公共 App 配置

| ID | 候选能力 | 真实业务场景 | 当前协议状态 | 最小协议内容 | 优先级 |
|----|----------|--------------|--------------|--------------|--------|
| BRAND-001 | 产品标题、logo、favicon | 登录页、Shell、浏览器标签 | absent | light/dark assets、alt text、fallback、empty semantics | P0 |
| BRAND-002 | 资源安全与约束 | 外链 logo、CSP、错误 MIME、超大图片 | absent | same-origin/allowlist、MIME、尺寸/比例、CSP、失败回退 | P0 |
| BRAND-003 | 默认 locale/theme/timezone | 首次访问和未登录页面 | absent | 枚举、auto 语义、优先级、invalid fallback | P0 |
| BRAND-004 | tenant-specific branding | SaaS 白标与组织切换 | absent | scope、继承、缓存 key、切换失效、fallback chain | P1 |
| BRAND-005 | 法务、支持与页脚链接 | 隐私条款、支持中心、状态页 | absent | typed links、label key、external policy、locale variants | P1 |
| BRAND-006 | 配置变更通知 | 管理员改品牌后刷新 Shell | absent | namespace/version/event、affected scopes、refetch/atomic apply | P0 |
| BRAND-007 | 邮件/导出/打印品牌 profile | 发票、PDF、邮件模板 | absent | channel-specific asset refs 与 ownership；可独立 extension | P2 |

### D. Shell、Navigation 与 App Lifecycle

| ID | 候选能力 | 真实业务场景 | 当前协议状态 | 最小协议内容 | 优先级 |
|----|----------|--------------|--------------|--------------|--------|
| SHELL-001 | Shell region/slot 模型 | header/sidebar/main/user/footer | partial | 标准 slots、必需/可选、排序、缺失行为、responsive ownership | P0 |
| SHELL-002 | Navigation label/i18n precedence | provider projection 与 manifest labelKey 冲突 | partial | labelKey vs fallback label、projection precedence、duplicate detection | P0 |
| SHELL-003 | Navigation hierarchy、group、order、badge | 多模块菜单、计数徽标 | partial | tree/group/order/badge source、visibility、stable ID | P1 |
| SHELL-004 | Internal/external links | 文档、支持站、跨应用 | partial | target kind、allowlist、rel/referrer、安全确认 | P1 |
| SHELL-005 | 路由、deep link、redirect、not-found | 刷新深链、旧 URL 迁移 | partial | canonical route、redirect rules、404 ownership、loop prevention | P0 |
| SHELL-006 | 页面标题与 breadcrumbs | 扫描定位、浏览器历史 | absent | title/label key、dynamic params、breadcrumb ownership、document title precedence | P0 |
| SHELL-007 | 用户菜单 | 个人资料、偏好、帮助、登出 | partial | slot/items、commands、visibility、identity summary | P0 |
| SHELL-008 | 全局 context switcher | tenant/workspace/project 切换 | absent | context kinds、current/allowed、switch command、dirty/cache/session effects | P1 |
| SHELL-009 | 全局 search / command palette | 快速找页面、资源和命令 | absent | provider contract、result kinds、permission、keyboard/focus、telemetry privacy | P1 |
| SHELL-010 | Unsaved-change navigation guard | 编辑表单时离开页面 | absent | dirty signal、confirm action、router/back/refresh semantics | P0 |
| SHELL-011 | App update/reload 生命周期 | 新前端版本、协议工件变化 | absent | update-available、safe reload、dirty-state 协调、force-upgrade | P1 |
| SHELL-012 | Host extension namespace | 产品需受治理地扩展壳层 | absent | namespaced extension、capability negotiation、unknown extension fail-closed | P0 |

### E. Preferences、i18n 与时间语义

| ID | 候选能力 | 真实业务场景 | 当前协议状态 | 最小协议内容 | 优先级 |
|----|----------|--------------|--------------|--------------|--------|
| PREF-001 | locale/theme/timezone precedence | 系统默认、用户偏好、设备偏好冲突 | absent | system→tenant→user→session/device precedence 与 reset | P0 |
| PREF-002 | 偏好持久化与同步 | 多设备登录、匿名到登录迁移 | absent | storage owner、sync/version、conflict、logout handling | P1 |
| PREF-003 | locale catalog discovery/fallback | 部分模块缺翻译、catalog 更新 | partial | supported/fallback chain、catalog version、missing-key behavior | P0 |
| PREF-004 | 日期/数字/货币/周起始 | 跨区域后台报表 | absent | locale/timezone/currency/numbering system 的明确来源 | P1 |
| PREF-005 | RTL 与文字方向 | 阿拉伯语/希伯来语产品 | absent | document direction、component capability、asset mirroring | P2 |
| PREF-006 | 密度、对比度、减弱动画、字号 | 高密度后台与无障碍偏好 | absent | supported values、precedence、OS media query integration | P1 |

### F. 全局 State、Error 与 Recovery

| ID | 候选能力 | 真实业务场景 | 当前协议状态 | 最小协议内容 | 优先级 |
|----|----------|--------------|--------------|--------------|--------|
| ERROR-001 | 统一错误 envelope | API、bootstrap、actions 错误需一致处理 | partial | stable code、messageKey/params、field errors、correlationId、retryable | P0 |
| ERROR-002 | HTTP/transport 分类 | 401/403/404/409/422/429/5xx/timeout/offline/cancel | partial | code mapping、用户可见性、日志级别、默认恢复动作 | P0 |
| ERROR-003 | Boot/manifest/schema/capability 失败页 | 启动前无法渲染普通 page | absent | terminal/retryable states、诊断详情、安全摘要、retry/support actions | P0 |
| ERROR-004 | Auth lost / forbidden / reauth | token 失效或权限变更 | absent | clear session、return intent、reauth/forbidden distinction、loop prevention | P0 |
| ERROR-005 | Global not-found | 未知 route/pageRef/resource | absent | route vs resource 404、back/home/search actions、status semantics | P0 |
| ERROR-006 | 表单字段与跨字段验证映射 | 422 后聚焦首个错误 | partial | field path、form-level errors、message params、focus/announce behavior | P0 |
| ERROR-007 | Retry/backoff/rate limit | 瞬时故障、429、后台轮询 | absent | retryAfter、exponential policy、idempotency、cancel/manual retry | P1 |
| ERROR-008 | Crash boundary 与问题报告 | renderer/host 未捕获异常 | absent | safe fallback、correlation/build context、redaction、reload/report action | P1 |
| ERROR-009 | Read-only/degraded mode | 部分服务故障但允许查询 | absent | affected capabilities、mutation disable reason、banner、recovery event | P1 |

### G. 全局 UX Commands、Overlay 与通知

| ID | 候选能力 | 真实业务场景 | 当前协议状态 | 最小协议内容 | 优先级 |
|----|----------|--------------|--------------|--------------|--------|
| UX-001 | Toast 生命周期 | action 成功/失败反馈 | partial | level、dedupe key、duration、persistent/actionable、a11y live region | P0 |
| UX-002 | Modal/drawer ownership | 详情、确认、分步流程 | partial | route ownership、stack/back/escape、focus restore、deep link | P0 |
| UX-003 | Confirm/prompt | 删除、危险操作、输入确认词 | partial | severity、message key、required input、cancel semantics | P0 |
| UX-004 | Host notification center | 审批结果、系统事件、未读数；不等同于通知业务模块 | absent | item schema、read/archive actions、badge、pagination、permission | P1 |
| UX-005 | Clipboard/share/print/export | 复制 ID、分享链接、打印报表 | absent | command kind、permissions、success/error、sensitive-data policy | P1 |
| UX-006 | Background job progress | 批量导出、导入、长任务 | absent | job state/progress/cancel/result/expiry、poll/push transport | P1 |
| UX-007 | Global busy/operation coordination | 防重复提交、路由离开时请求未完成 | absent | operation identity、scope、cancelability、blocking policy | P1 |

### H. Files、Download、Realtime 与并发

| ID | 候选能力 | 真实业务场景 | 当前协议状态 | 最小协议内容 | 优先级 |
|----|----------|--------------|--------------|--------------|--------|
| FILE-001 | Upload action | 附件、头像、导入文件 | covered/partial | 复核 metadata、scan/processing states 与现有 upload 的兼容扩展 | P1 |
| FILE-002 | Download action | 导出文件、受权附件下载 | absent | filename/content type/size、Content-Disposition、authz、range、错误 | P0 |
| FILE-003 | Preview/media | 图片、PDF、音视频预览 | absent | preview kind、safe URL、expiry、sandbox、fallback download | P1 |
| FILE-004 | Malware scan/quarantine | 企业上传安全 | absent | pending/clean/rejected/quarantined、retry/delete、用户提示 | P1 |
| FILE-005 | Chunked/resumable upload | 大文件、弱网络 | absent | session/chunks/resume/checksum/expiry/cancel | P2 |
| RT-001 | SSE/WebSocket capability | 实时通知、任务进度 | absent | transport discovery、auth、resume cursor、heartbeat、reconnect | P1 |
| RT-002 | Cache/data invalidation event | 配置或资源被其他用户修改 | absent | topic/resource/version、scope、refetch/patch semantics | P1 |
| RT-003 | Optimistic concurrency/conflict | 多人同时编辑 | partial | version/ETag、409 payload、compare/merge/reload actions | P1 |

### I. Tenancy、Entitlement 与上下文

| ID | 候选能力 | 真实业务场景 | 当前协议状态 | 最小协议内容 | 优先级 |
|----|----------|--------------|--------------|--------------|--------|
| TENANT-001 | tenant/org discovery 与选择 | 一个账号加入多个组织 | absent | allowed contexts、display metadata、default、selection command | P1 |
| TENANT-002 | context 切换失效语义 | 切组织后权限、branding、manifest 全变 | absent | session/manifest/cache/preferences invalidation 与原子重启 | P1 |
| TENANT-003 | entitlement/license gating | 套餐决定模块可用性 | absent | entitlement source/reason/expiry；不得伪装 permission | P1 |
| TENANT-004 | region/data-residency context | 多区域部署与合规提示 | absent | region display、switch/restriction、privacy notice；不暴露内部 secret | P2 |

### J. Observability、Security、Privacy 与 Accessibility

| ID | 候选能力 | 真实业务场景 | 当前协议状态 | 最小协议内容 | 优先级 |
|----|----------|--------------|--------------|--------------|--------|
| OBS-001 | Correlation/trace context | 用户报错与服务端日志关联 | absent | request/correlation IDs、响应暴露规则、错误页复制动作 | P0 |
| OBS-002 | Client telemetry event envelope | 启动失败、action 失败、性能分析 | absent | event name/version、context、sampling、PII redaction、consent | P1 |
| OBS-003 | Audit-link semantics | 敏感 action 对应审计记录 | absent | actor/effective actor、action/resource/session/correlation link | P1 |
| SEC-001 | Trusted URL/resource policy | branding、redirect、download、外链 | absent | same-origin/allowlist/scheme、CSP、open-redirect 防护 | P0 |
| SEC-002 | Sensitive data 与缓存/日志 | token、个人资料、错误详情 | absent | redaction、no-store、clipboard/export restrictions、diagnostic tiers | P0 |
| SEC-003 | Permission 与 feature/entitlement 分界 | menu 可见不等于 API 授权 | partial | authoritative enforcement、UI hint、denial behavior、provenance | P0 |
| A11Y-001 | Shell landmarks 与 focus restoration | route/modal/error 后键盘定位 | absent | landmarks、focus target/restore、skip link、tab order | P0 |
| A11Y-002 | Live region 与异步状态 | toast、loading、job、错误播报 | absent | politeness、dedupe、busy semantics、announcement content | P1 |
| A11Y-003 | Keyboard interaction contract | menu、table、command palette、drawer | partial | standard keys、escape/back、roving focus、shortcut conflict | P1 |
| A11Y-004 | Contrast/motion/direction capability | 高对比、减弱动画、RTL | absent | capability declaration、preference integration、unsupported behavior | P1 |

### K. 协议治理、扩展与弃用

| ID | 候选能力 | 真实业务场景 | 当前协议状态 | 最小协议内容 | 优先级 |
|----|----------|--------------|--------------|--------------|--------|
| GOV-001 | Host/App capability registry | 文档需要 Host 能力，Host 只实现子集 | partial | capability IDs、versions、dependencies、required/optional、error | P0 |
| GOV-002 | 扩展命名空间 | 产品特有能力又不能污染核心字段 | absent | owner/name/version、schema URI、negotiation、unknown handling | P0 |
| GOV-003 | 兼容与弃用生命周期 | 2.7.0 到 vNext 渐进迁移 | partial | since/deprecated/removed、兼容窗口、migration、fixture matrix | P0 |
| GOV-004 | Semantic validator profile | JSON Schema 通过但跨字段语义错误 | partial | structural/semantic/behavioral 层次、错误码、正反 fixtures | P0 |
| GOV-005 | 请求对象闭集规则 | API 是否可忽略未声明字段 | needs clarification | `additionalProperties` 默认、扩展字段规则、server enforcement | P0 |
| GOV-006 | i18n key/fallback precedence | label 与 labelKey 双源 | needs clarification | key/fallback、catalog owner、missing-key 与 duplicate behavior | P0 |
| GOV-007 | Conformance claim 与 evidence | Host 宣称支持某版本/能力 | absent | claim scope、test suite/version、evidence artifact、partial 禁止项 | P0 |

## 4. 建议的协议产物集合

命名由上游仓库最终决定；本附件只约定需要形成的语义块：

1. **Host/App 基础规范**：Host 与 page renderer 的边界、生命周期、capability negotiation、扩展治理、fail-closed 原则。
2. **Bootstrap schema + 状态机**：public/authenticated 阶段、manifest/provenance、缓存、maintenance/degraded/upgrade-required。
3. **Auth/session schema + 状态机**：discovery、principal/session、login/refresh/logout、reauth、transport/security profiles。
4. **Branding/config schema**：资产、locale/theme/timezone defaults、scope/precedence、缓存与 config-change event。
5. **Shell/app schema**：slots、navigation、routing、breadcrumbs/title、user menu、context switch、lifecycle。
6. **Global error envelope + recovery matrix**：HTTP/transport/domain errors、field mapping、correlation、retry/auth/not-found/crash。
7. **Preferences/i18n schema**：precedence、persistence、catalog/fallback、format 与 accessibility preferences。
8. **Action/async 扩展**：download、background job、notification、realtime、global overlays；复用现有 toast/modal/upload。
9. **Security/privacy/a11y normative requirements**：trusted URL、redaction/cache、permission boundary、focus/live-region/keyboard。
10. **Compatibility matrix 与 conformance fixtures**：2.7.0→vNext migration、每项 capability 的正反例、semantic validator 和 evidence format。

## 5. 明确非目标

- 不把 CSS、像素布局、视觉 token 实现细节写成跨产品协议；只规范必要的语义与无障碍约束。
- 不在协议中携带 OAuth client secret、数据库模型或供应商私有凭据。
- 不通过 Host 协议放行运行时远程插件市场；Charter 已明确远程插件市场不是当前目标。受控 extension namespace 不等于远程代码执行。
- 不用 feature flag、navigation visibility 或 entitlement 替代服务端 authorization。
- 不要求所有 Host 页面都由普通 page schema 渲染；但 Host UI 不得自行发明与协议冲突的业务语义。

## 6. S2/S3 出口门禁

> 以下均为未来出口条件，不是已完成事实。截至 2026-08-12，S2/S3 均未通过，I-003/I-005 仍为 `open`，且尚未发生正式 Goal Audit。

S2 方案冻结前必须满足：

- [ ] 本目录每一行均有 adopt/reserve/out 处置和理由；
- [ ] P0 项有 schema/状态机/能力/错误/安全/fixtures 的可核对提案；
- [ ] IMP-001～004 已被新协议裁定为实现偏离或合法 Host extension；
- [ ] 完成 `cross` 方案审视，required findings 已合法闭合。

S3 宣称“新协议到手”前必须满足：

- [ ] 上游合并版本或 commit 可核对；
- [ ] 本仓 provenance、schemas、registry、fixtures 与 capability matrix 已固定到同一来源；
- [ ] 2.7.0 的兼容、迁移、弃用与 fail-closed 行为已有正反证据；
- [ ] I-003/I-005 为 `verified`。只有此后才能进入 S4 修改 `apps/api` / `apps/web`。
