---
id: r1-state-feedback-matrix
doc: evidence-attachment
status: active
parent: GOAL-002-r1-scope-semantics-freeze
created: 2026-09-16
updated: 2026-09-16
version: 0.1.0
---

# R1 · 状态、dirty-state 与反馈基线矩阵

本附件记录 2026-09-16 从当前代码得到的基线。`proposed` 行是待冻结方案，不是已发生事实。

## 1. 状态与 dirty-state 基线

| 场景 | 当前实现事实 | R1 待冻结语义 |
|------|--------------|---------------|
| 列表查询 | `SchemaCrudProvider.queries` 按 table id 保存 `q/filters/sort/order/page/pageSize`；`SchemaTable` 通过 `setTableQuery` 更新 | Saved View 是否保存 `q/filters/sort/order/pageSize`；默认不保存当前页、选择集、recordView/modal |
| 筛选/排序/分页 | SchemaTable 的筛选、排序和分页改变 query；search form 的 q 提交后写入 target table；当前不写入 App URL/history | 需明确这些变化是否视为可保存视图状态，而不是业务 dirty |
| 列配置 | DataTable 按 Schema 声明渲染全部 columns，当前没有可见性状态或用户列配置 | 是否在首波实现列可见性；字段序列化必须 allowlist 到当前 Schema columns |
| 默认表单 | `FormInner` 在本地保存 values，以 modal row、recordSource、route readOnly 值初始化；未注册 dirty 状态 | 非 search 的可编辑表单以初始化快照为 baseline；字段恢复到 baseline 时清 dirty |
| Search 表单 | `mode: search` 的 q/filters 绑定查询；当前 search form 也使用 FormInner 本地 values | 建议 search/query 改动不阻断离开；它们只进入 Saved View/查询状态，不视作未提交业务修改 |
| 内部导航 | `App.onNavigate` 直接 `history.pushState`，导航链接、面包屑、shell 菜单均调用它；没有离开确认 | dirty 时导航先进入确认状态；取消保持原页面，确认丢弃后再 pushState |
| 浏览器刷新/关闭 | 当前代码未发现 `beforeunload` 注册 | dirty 时只注册原生 `beforeunload`；没有 dirty 时不注册；文案由浏览器控制 |
| 返回/前进 | `popstate` 直接解析路径并切换 `path/routeQuery` | dirty 时需暂存目标并阻止状态提交，确认后完成历史移动；实现前需用浏览器自动化验证 |
| 提交成功 | `crud.submitForm` 成功后设置成功反馈、reloadList、关闭 modal、清 selected row；inline FormInner 保持组件值 | 成功提交清 dirty；失败/fieldErrors 保持 dirty；不能自动重发写请求 |
| reset/cancel | search reset 立即清 query；modal close 当前直接关闭；未发现 dirty-aware cancel | reset 到 baseline 清 dirty；有 dirty 的 modal cancel 需确认，取消留在 modal |

## 2. 反馈与错误合同基线

| 反馈面 | 当前实现事实 | R1 待冻结语义 |
|---------|--------------|---------------|
| 成功 Toast | `FeedbackRegion` 渲染 `role=status`，成功反馈 4 秒自动消失 | 保留 transient status toast；同一操作只产生一个成功事件 |
| 写入失败 | `SchemaCrudProvider` 将 ActionResult 设为 `kind:error`；FeedbackRegion 用 `role=alert`，不自动消失；表单另显示 fieldErrors/formError | 写入失败不自动 retry；保留当前表单值/dirty；只在用户重新点击时重试 |
| 列表失败 | DataTable 根据 `resolveAsyncDisplayState` 显示 `role=alert` 与 Retry；SchemaTable 的 retry 只增加本地 retry nonce | list GET、page schema GET 等幂等读取允许显式 retry；重试必须产生新请求而非重复写入 |
| API envelope | `readResourceApiError` 读取 HTTP status、`error/message`、`messageKey/params`、`fieldErrors`、`correlation_id` 或 `X-Request-ID` | 用户文案优先 catalog key；诊断 code/correlation 可辅助显示；不得把 token、密码或原始敏感载荷回显 |
| 401/重新认证 | AuthProvider 的 authFetch 在会话丢失时进入 `reauth-required`，由 HostFailureScreen 提供 reauth | 维持 Host 层终态，不由普通 Toast 伪装成业务成功/失败 |
| 403/权限 | SchemaCrudProvider 对声明的 permission target fail closed；动作未执行时设 `ACTION_NOT_EXECUTED` | Saved View 恢复遇到当前权限不能满足的字段/动作时必须 fail closed，不扩大权限 |
| 404/失效 | 当前 ResourceApiError 保留 HTTP 404 与 code；没有 Saved View 失效语义 | Saved View 读取/恢复遇到不存在或不再允许的引用时丢弃该视图并提示，不污染当前 query |
| 409/并发 | 当前错误 envelope 可保留冲突 code；写请求没有自动重试策略 | 冲突必须停在当前 dirty 草稿，给 reload/重新编辑入口，禁止盲目重复提交 |
| 429/维护/不可用 | 后端 operational handler 使用 `SERVICE_MAINTENANCE` / `SERVICE_UNAVAILABLE`；Host failure 已有 maintenance/unavailable kind，普通资源反馈尚未统一映射 | 维护/不可用显示可重试状态；若服务端给出 retry-after/after 秒数则遵守；不把维护当作 field error |
| 离线/超时 | HostFailureScreen 有 offline/timeout 映射；普通 resource fetch 当前主要显示 `REQUEST_FAILED` 文本 | 读取可 retry；写入只显示失败并保持 dirty，除非后续确认有幂等键合同 |
| 无障碍 | loading 使用 `role=status`，错误使用 `role=alert`；HostFailureScreen 首次终态聚焦标题并只按 failureId 宣告一次 | toast、列表错误、确认对话框均需可键盘到达；同一 failure 不重复宣告 |

## 3. Saved View 持久化取舍（待用户确认）

| 方案 | 优点 | 代价/风险 | 与当前 VP 边界 |
|------|------|-----------|----------------|
| A · 浏览器本地持久化（推荐用于首波） | 不新增 API/迁移；可按 `user.id + pageId + tableId` 隔离；失败可 fail closed；与现有 locale/timezone 偏好模式一致 | 只跨刷新/同一浏览器，不跨设备；用户清理站点数据会丢失；查询值存在浏览器存储 | 不新增业务域/第二服务；需明确“用户级”定义为本浏览器账号命名空间 |
| B · 现有数据库 + 专用 API | 跨浏览器/设备；服务端可做严格 owner 校验和版本化 | 新增表/迁移/handler/API/错误合同；扩展交付面；需要独立 data/privacy 审计 | 可纳入 VP，但会提高 R2 风险与独立审计成本 |
| C · 仅 URL/history | 可分享/可回退，零存储 | 不是 Saved View；可能把筛选值暴露到 URL/日志；无法提供命名视图 CRUD | 不满足首波 Saved View 退出判据，不建议采用 |

在 I-037-002 未经用户确认前，A/B 不得写成 accepted；C 仅作为排除项记录。
