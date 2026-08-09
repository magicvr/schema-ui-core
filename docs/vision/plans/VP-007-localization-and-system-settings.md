---
doc_type: vision-plan
id: VP-007-localization-and-system-settings
title: 多语种与系统设置产品化
status: active
vision_ref: schema-ui-core-admin-foundation@0.2.0
lead_workspace: workspace-007-localization-and-system-settings
created: 2026-08-09
updated: 2026-08-09
version: 0.2.0
parent: null
---

# VP-007 · 多语种与系统设置产品化

## 意图

在 [VP-003](VP-003-modular-admin-architecture.md) 的单主线模块架构、[VP-005](VP-005-design-system-and-ui-experience.md) 的产品化界面和 [VP-006](VP-006-full-protocol-contract-v2-7-0.md) 的完整协议兼容面之上，建立覆盖前端与后端边界的多语种基础能力，并把既有 `admin.settings` 从仅含站点标题与 Logo URL 的设置面扩展为可直接用于 fork 项目的系统设置产品面。

本波次首发语言固定为 `zh-CN` 与 `en-US`；系统默认显示语种为 `auto`。有效语种的方向级优先级固定为：

```text
用户显式选择
  → 系统指定默认语种（非 auto）
  → 浏览器语言偏好（auto）
  → en-US 安全回退
```

用户显式语种选择必须位于普通用户可达的 Shell / 用户菜单，不得要求用户具备系统设置管理权限后才能切换。系统设置页只负责管理员可控的全局默认值与品牌配置。

### 用户确认（2026-08-09）

- 将“多语种 + 系统设置产品化”作为同一个新纲领波次落盘，不修订 Charter，不重开已关闭 VP 或历史工作区。
- 首发 `zh-CN` + `en-US`；全局默认语种为 `auto`，按浏览器语言偏好自动匹配，最终回退 `en-US`。
- 系统设置首版按 **General / Branding / Localization / Appearance** 四类组织。
- 多语种运行时覆盖 `mvp` 与 `admin` Profile；系统设置编辑页继续属于 `admin.settings`，仅在 `admin` Profile 且用户具备相应权限时出现。
- 后端反馈在有界可行时随前端当前语种显示；无论是否由后端生成最终文案，稳定错误码与前端本地化回退都必须保留。

### 交付范围

| 区域 | 方向级范围 |
|------|------------|
| 多语种基础 | BCP 47 语种标识、支持语种清单、语种解析与回退、翻译资源装载、缺失 key 可观察且安全回退 |
| 用户控制 | 登录前后均可切换语种；用户显式选择优先于系统/浏览器默认；具体持久化边界由信息门禁冻结 |
| 前端覆盖 | 登录页、Shell、导航、Renderer 状态与反馈、表单/表格、Settings，以及 S0 冻结的双 Profile Runtime Manifest 可达页面与 Schema 可见文本；同步文档 `lang` 与 locale-aware 日期/数字格式 |
| General | 站点标题；保持现有 Shell、登录页与浏览器标题联动 |
| Branding | Logo、favicon、浅色/深色 Logo，含预览、输入校验、清空/恢复默认；首版使用同源路径或 HTTPS URL |
| Localization | 默认语种 `auto | zh-CN | en-US`、默认时区；日期/数字格式默认随有效 locale，不在首版暴露任意格式模板 |
| Appearance | 全局默认主题 `auto | light | dark`；用户显式主题选择继续优先 |
| 权限与审计 | 沿用并验证 `settings.read` / `settings.write`；设置变更留下可核对操作日志与配置刷新信号 |
| 后端反馈 | 保持稳定错误码；兼容扩展 `messageKey` / 参数 / locale；有界可行时按请求语种返回已编目的用户可见消息并声明响应语种 |

### 最小可枚举证据面（F-V029）

退出判据 2 与 6 使用同一份、在 S0 落到 lead 工作区 Root `attachments/` 并由 Root 决策冻结的覆盖表。覆盖分母至少包含：

1. **固定 UI 面**：匿名登录页；Shell header / top / sidebar / user navigation；语言切换器；通用 loading / empty / error / permission denied / success feedback；通用 table、search、form、modal 与 validation 文案；Admin Settings 四类设置面。
2. **Manifest / Schema 面**：S0 时 `mvp` 与 `admin` 两个真实 Runtime Manifest `pages[]` 的 **pageId 并集**、各 page 引用的 `schemaUrl`，以及本 VP 在关门前新增的任何可达 page/schema。覆盖表逐项记录 Profile、pageId、schemaUrl、翻译来源、fallback 与证据路径；不得用“代表性页面”替代该并集。`settings` 只在 `admin` Profile 可达，但仍属于分母。
3. **固定主流程**：
   - `M1`：匿名启动 → `auto`/显式 locale 生效 → 登录成功或失败反馈；
   - `M2`：登录 → Shell 导航 → `overview` / `users` / `roles` 的可达读路径，并以有权限账号完成至少一次用户或角色写表单及验证/权限失败反馈；
   - `M3`（Admin）：进入 `settings` → 修改一项 General / Branding / Localization / Appearance 配置 → Shell / 登录页或公开 bootstrap 的对应投影可观察生效；
   - `M4`：制造缺失翻译 key → 记录 locale、key 与 UI/schema 路径 → 使用安全文本 fallback，且主流程仍可完成。
4. **缺失翻译纪律**：缺失 key 不得渲染为空、不暴露未处理异常、不阻断操作；必须由测试或等价遥测观察到。任何未纳入上述分母的用户可见文案面，只能经用户书面 `accepted-residual`（明确范围、缓解、责任人和复审触发）排除，不能以“代表性”静默收缩。

## 方向级退出判据

在同时满足下列方向、且均有工作区 Q2 证据时，本 VP **可以**提议 `closed`：

1. **语种解析与用户控制闭合**
   `zh-CN`、`en-US` 与 `auto` 的解析、匹配、回退及切换行为可验证；用户显式选择、系统默认、浏览器偏好与 `en-US` 回退的优先级一致应用于登录前后；HTML `lang`、日期和数字格式跟随有效 locale。

2. **前端可见文本形成可维护翻译面**
   上述“最小可枚举证据面”的固定 UI、Runtime Manifest page/schema 分母和 M1～M4 主流程不存在必须依赖硬编码英文才能完成的路径；Manifest 的 `titleKey` / `labelKey` 必须同时具备真实解析、文本 fallback 和缺失 key 可观察证据。Schema 文本的本地化不得私自改写或冒充 `schema-ui-docs@v2.7.0` 语义。覆盖表未列明或无证据的页面不得以“代表性已通过”关闭本判据。

3. **系统设置产品面可用**
   授权管理员可在 General / Branding / Localization / Appearance 中读取和修改本 VP 固定的首版字段；站点标题、品牌资产、默认语种/时区和默认主题能在 Shell、登录页及公开启动配置中一致生效；预览、校验、恢复默认、权限失败和刷新行为可验证。

4. **Profile 与公开启动边界一致**
   `mvp` 与 `admin` 使用同一前端 build 和同一多语种运行时；`admin.settings` 编辑面只在 `admin` Profile 暴露，`mvp` 不因本 VP 静默扩张模块集合。登录前所需的非敏感品牌与 locale 启动信息有明确、兼容且可缓存验证的公开读取路径。

5. **错误与提示本地化保持兼容**
   现有稳定错误码继续可由客户端判断；前端能按错误码/key/参数以当前语种呈现用户可见反馈，这是不可降级的关门门禁。`I-L10N-004` 只能以二选一结果关闭：**(a)** 服务端 locale 路径纳入本波次，并对已编目的验证、认证与设置错误提供语言协商、响应 locale 和失败回退证据；或 **(b)** 用户书面 `accepted-residual`，明确本波次仅保证前端本地化 + 稳定错误码，记录服务端延期范围、理由、缓解、责任人和复审触发。路径 (b) 不解除本判据的前端半边。未编目/内部错误必须安全回退且不泄露诊断信息；不得以翻译字符串取代错误码契约，也不得以“成本无界”静默跳过服务端路径。

6. **质量与关门证据完整**
   自动化或等价证据矩阵必须复用“最小可枚举证据面”的同一分母：行至少覆盖 `zh-CN` / `en-US` × `mvp` / `admin` × 匿名 / 已认证，列覆盖固定 UI、冻结 pageId/schema 集、M1～M4、权限正反例、缺失翻译、配置刷新和错误回退。Profile 不可达单元格须标 `N/A` 并写明模块边界，不能算作 pass；其余单元格须有证据或用户书面 residual。文档说明 Profile 可见性与配置优先级；lead 工作区 Root 完成约定范围，开放 required Goal/Vision findings = 0，并由用户确认关门。

## 信息门禁

| id | 问题 / 所需信息 | 级别 | 最晚阶段 | 验证 / 决策动作 | 初始状态 |
|----|-----------------|------|----------|-----------------|----------|
| `I-L10N-001` | Schema 驱动页面的字段标签、说明、动作和服务端文档如何本地化，且不创建宽于或冲突于 `schema-ui-docs@v2.7.0` 的私有协议语义？ | required | 多语种方案冻结前 | 盘点 v2.7.0 可用 key/text 字段与当前 Renderer；比较前端 key 解析、服务端 locale overlay、上游提案等路径并冻结兼容策略 | open |
| `I-L10N-002` | 用户显式语种选择首版持久化在浏览器本地还是账号资料；匿名到登录后的合并规则是什么？ | required | 用户控制实施前 | 结合现有 auth/session 与本地 theme 机制形成优先级和迁移矩阵；用户确认后冻结 | open |
| `I-L10N-003` | 公开品牌/locale 启动配置是兼容扩展 `/api/branding`，还是建立新的公开 bootstrap 契约？缓存与配置刷新语义如何保持一致？ | required | Settings API 方案冻结前 | 对当前 branding 消费端、Profile 路由和缓存头做差量盘点，冻结兼容路径 | open |
| `I-L10N-004` | 当前 `{error,message}` 错误 envelope、重复 `writeError` 与前端直显链路扩展到 locale 协商的真实成本和兼容边界是什么？ | required | 后端提示本地化实施前 | 盘点用户可见错误码及调用点；用认证/验证/设置错误验证 `Accept-Language`、`Content-Language`、key/params 与 fallback；关闭结论必须选择 exit 5 路径 (a) 实施证据或路径 (b) 用户书面 accepted-residual，禁止仅写“成本无界” | open |
| `I-L10N-005` | 默认时区的存储、展示和服务器时间语义如何定义，避免把显示时区与持久化时间混为一谈？ | required | Localization 设置实施前 | 固定 UTC 存储、显示转换、`auto`/指定时区和无效时区失败语义 | open |

这些未知不阻断 VP 落盘或激活，但在其“最晚阶段”前未关闭时必须阻断对应方案冻结或实施。任何范围收缩或 residual 必须按 P-004 由用户书面裁决。

## 建议实现阶段（供后续 `/govern` 建立 Root 纲领路线图时参考）

| 阶段 | 目的 | 建议检查点 |
|------|------|------------|
| S0 | 差距与契约冻结 | 关闭 I-L10N-001～005；冻结 locale、翻译资源、公开 bootstrap、错误兼容和时区语义 |
| S1 | 多语种核心 | locale resolver/provider、资源装载、缺失 key/fallback、用户切换、HTML lang 与格式化 |
| S2 | 前端与 Schema 覆盖 | 按冻结覆盖表完成固定 UI 与双 Profile Runtime Manifest page/schema 分母双语化；`titleKey` / `labelKey` 真解析、fallback 与缺失 key 可观察 |
| S3 | 系统设置产品化 | 四类设置、公开启动配置、品牌/locale/theme 生效、权限/审计/刷新闭环 |
| S4 | 后端用户可见反馈 | 稳定错误码 + 前端本地化保底；按 I-L10N-004 结论实施有界服务端 locale 协商 |
| S5 | 双 Profile 验证与关门 | 两语种 × 两 Profile × 登录前后/权限/失败矩阵、文档、审计与 close-out |

阶段和子目标由实现层根据 P-001/P-005 进一步冻结；本表不是 Goal progress 或已实施事实。

## Non-goals（非目标）

- 不改变 Charter 目的、成功边界或非目标；不建立第二 active Charter。
- 不重开 VP-002～006 或把本波次写入其已关闭工作区。
- 不把 `admin.settings` 静默加入默认 `mvp` Profile；不以更改默认 Profile 掩盖 Settings 可见性边界。
- 首版不提供 Logo/favicon 文件上传、对象存储或媒体管理；使用同源路径或 HTTPS URL。若要上传，应另行登记存储、安全与生命周期信息门禁。
- 不在本波次纳入 SMTP、SSO、密码策略、通知、遥测、数据保留、插件管理或任意自定义 HTML/CSS/脚本。
- 不翻译开发日志、内部诊断文本或以翻译文案替代可机读错误码。
- 不重新定义 `schema-ui-docs@v2.7.0` 协议语义；协议缺口回到兼容决策或上游提案。
- 不为 VP 建 Goal 五件套，不在 `docs/vision/` 记录 progress%。

## 与前后 VP 的关系

| VP | 关系 |
|----|------|
| VP-002 | 继承既有站点标题、Logo URL、Settings 权限/持久化/操作日志事实；不重开历史交付。 |
| VP-003 / VP-004 | 继承单主线、Profile、模块贡献和配置命名空间边界；新增 locale/settings 能力不得回到中央业务硬编码。 |
| VP-005 | 继承设计系统、主题与交互产品化基线；多语种后文本增长、状态和响应式布局不得造成回退。 |
| VP-006 | 继承 `schema-ui-docs@v2.7.0` 完整协议兼容面；翻译实现不得用私有扩展冒充上游契约。 |
| 后续业务 VP | 默认复用本 VP 的 locale、错误提示和系统设置基础，不在各业务模块重复建立翻译运行时。 |

## 工作区绑定

| workspace_id | root_goal | role | joined | notes |
|--------------|-----------|------|--------|-------|
| workspace-007-localization-and-system-settings | GOAL-001-localization-and-system-settings | lead | 2026-08-09 | 用户确认激活后唯一 lead / delivery；`/govern` 已 scaffold Root + S0–S5 纲领（2026-08-09）；`vision_role: delivery` |

用户已于 2026-08-09 确认将本 VP **激活**（`planned` → `active`），并指定唯一 lead / delivery 工作区 **`workspace-007-localization-and-system-settings`**（slug 按 VP-007 id 与既有 workspace-00N 惯例，用户本轮书面确认）。物理 scaffold（Root `GOAL-001-localization-and-system-settings`，`primary_plan` / `plan_refs` 均为本 VP）交 `/govern`。激活与建区 **不**构成任何多语种或 Settings 能力已交付。

## 关门记录

（仅 `closed` / `abandoned` 时填写。）

| date | outcome | summary | evidence_links | residuals |
|------|---------|---------|----------------|-----------|
| — | — | — | — | — |

## 规划修订短史

| date | version | change |
|------|---------|--------|
| 2026-08-09 | `0.1.0` | 用户确认新建本 VP：首发 `zh-CN` / `en-US`，系统默认 `auto`；系统设置按 General / Branding / Localization / Appearance 四类；多语种覆盖双 Profile，Settings 编辑面保持 Admin 权限边界；状态 `planned`，尚未激活或绑定工作区。 |
| 2026-08-09 | `0.1.1` | 响应 VRev-016：F-V029 → fixed，新增固定 UI + Runtime Manifest page/schema 并集 + M1～M4 的可枚举证据分母，并让 exit 6 复用同一矩阵；F-V030 → fixed，I-L10N-004 只允许服务端实施证据或用户书面 residual 两种关闭路径。保持 `planned`、0 区，未激活。 |
| 2026-08-09 | `0.2.0` | 用户确认激活：`planned → active`；`lead_workspace` = `workspace-007-localization-and-system-settings`（唯一 lead / delivery，角色 `delivery`）；物理 scaffold 交 `/govern`。激活不构成能力已交付。 |
