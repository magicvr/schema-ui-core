---
id: GOAL-012-w11-mfa-ux-review
doc: decision
status: active
parent: GOAL-001-design-implementation-conformance
created: 2026-08-15
updated: 2026-08-15
version: 0.1.0
---

# 决策记录 · GOAL-012

## 问题清单（2026-08-15 落盘）

来源：① 本会话以真实用户视角对已实现 api/web 功能页的审视（U 类）；② 用户实测补充 3 项 MFA 缺陷（M 类，P0）。

### M 类 · 个人中心 MFA 缺陷（用户报告，P0）

| ID | 现象 | 初步定位（代码阅读，待 I-001 复现验证） | 期望行为 |
|----|------|------------------------------------------|----------|
| M-01 | 绑定 MFA 时只有 secretBase32 / otpauthURL 文本输入框，**无二维码**供扫描 | apps/web/src/components/mfa-manager.tsx 的 data-mfa-enroll 区仅渲染两个只读 Input；API 已返回 otpauthURL（apps/api/internal/modules/mfa/service.go Enroll() → otpauthURL(...)），前端未渲染二维码 | 提供二维码（由 otpauthURL 直接生成）供扫码器扫描 |
| M-02 | **停用 MFA 不成功却直接退出登录**；重新登录仍要求 MFA 验证码，登录后 MFA 依然绑定，未解绑成功 | 根因①：无效动态码/恢复码被映射为 **HTTP 401**（apps/api/internal/handler/mfa.go writeMFAError：ErrMFAInvalid → StatusUnauthorized）；apps/web/src/account/auth-client.ts authFetch 将**任意 401** 视为会话过期 → refresh 失败或重试仍 401 → clearTokens() + onAuthLost 强制登出。根因②：即使解绑成功，服务端 Disable 后 BumpTokenVersionAndRevokeAll 吊销**全部**会话（含当前），前端随后 refresh()（GET /api/mfa/status）401 → 无提示登出 | 无效码应报错留页重填（不得登出）；解绑成功后给出明确提示并正常处理会话吊销（如「MFA 已解绑，请重新登录」） |
| M-03 | **绑定 MFA 输错动态码直接退出登录**（应报错让用户重填） | 同 M-02 根因①：Confirm → ErrMFAInvalid → 401 → authFetch 登出（mfa-manager.tsx confirm() → postJSON → authFetch） | 报错 + 保留 pending 状态允许重填（不登出） |

> 修复方向（待 I-001 裁决）：自服务端点（confirm/disable/rotate）的 ErrMFAInvalid 改映射为 **400**（而非 401），authFetch 不触发会话丢失；登录二步验证 /api/auth/mfa/verify 保持 401（由 mfaVerify 直接解析，不经过 authFetch，不受影响）。

### U 类 · UX 审视改进（按优先级分批）

| ID | 优先级 | 页面/模块 | 问题 | 期望 |
|----|--------|-----------|------|------|
| U-01 | P0 | /users | 分配角色是 textarea 手输逗号分隔 role key（如 admin, editor），易输错/记错 | 多选下拉或复选框组，选项动态来自 /api/roles |
| U-02 | P0 | /roles | 权限项写死在 schema（仅 8 项 users/roles/settings/operations.*），新增 8+ 模块（scheduled-tasks、file-library、data-dictionary、recycle-bin、data-permission、mfa 等）的权限无法在前端勾选授权 | 权限元数据接口 + 模块分组权限矩阵（支持按模块全选） |
| U-03 | P1 | 全局 | 操作反馈为页面静态 Alert（FeedbackRegion）：不自动消失、无关闭按钮、挤压页面布局 | Toast/Message 浮层（自动消失、可关闭、不占布局） |
| U-04 | P1 | users/roles/data-dictionary/scheduled-tasks/file-library/recycle-bin 等 | 表格无搜索与筛选（后端已支持 q 搜索/sort，前端 schema 未配搜索表单与筛选器） | 统一搜索表单 + 状态/字段筛选 |
| U-05 | P1 | 所有表格页 | 行操作按钮过多（用户行 8 个：编辑/角色/密码/启用/禁用/解锁/重置MFA/删除），拥挤、小屏换行、高危操作易误触 | 保留 1～2 个高频操作，其余收进「更多（⋯）」下拉菜单 |
| U-06 | P1 | 表格分页 | 仅上一页/下一页+页码，无每页条数切换、无跳页 | pageSize 切换（10/20/50/100）+ 快速跳页 |
| U-07 | P1 | 表格 | 空状态仅文字（No rows.） | 图形化空状态 + 「新建」快捷引导 |
| U-08 | P2 | /dashboard | 仅 2 个 statCard（Users/Roles），3 列 grid 大片空白 | 卡片下钻跳转、快捷操作入口、chart 趋势 |
| U-09 | P2 | 登录页 | 验证码无「刷新」按钮（看不清只能故意输错/刷新页面）；MFA 两步 UI 同时平铺 TOTP 与恢复码输入框 + 双按钮（Sign In/Verify），回车语义混乱 | 验证码刷新按钮；两步验证过渡卡片（先 TOTP，可切换恢复码） |
| U-10 | P2 | 登录/改密 | 密码输入无显示/隐藏切换 | Eye/EyeOff 切换 |
| U-11 | P2 | /account | 会话列表缺设备/浏览器/UA/IP 指纹与「当前设备」标记；页面 4 块平铺过长 | 设备信息列 + 当前设备高亮；Tabs 分区（资料/安全/会话） |
| U-12 | P2 | 通知铃铛 / 通知页 | 未读计数仅在挂载时拉取（无轮询）；列表页标记已读后铃铛徽章不联动；下拉无「全部已读」；列表无未读/已读筛选 | 轮询/联动刷新；全部已读快捷按钮；已读筛选 |
| U-13 | P2 | /scheduled-tasks、/data-dictionary | Cron 表达式门槛高（纯手写 5 段式）；字典主子表分页跳转割裂 | 常用 Cron 预设 + 下次执行预览；一体化联动视图（左类型右条目） |
| U-14 | P2 | 全局 | 详情/长表单一律 Modal（操作日志详情等）；enabled/locked/mfaEnabled 等布尔无状态徽章 | 详情用 Drawer；状态 Badge（绿启用/灰禁用/红锁定） |

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | M-02/M-03 根因确认与修复方案（401 语义 vs authFetch 策略；解绑成功后的会话吊销提示） | 方案 | S2 | 复现验证 + 回归测试（E-002） | **closed** | — | D-001/D-002 |
| I-002 | required | UX P0 交互方案（角色多选数据源；权限动态化 API 形态） | 方案 | S3 | D-002 裁决（optionsSource + 目录端点；分组留 P2） | **closed** | — | D-002 |
| I-003 | non-blocking | 二维码渲染实现（新增依赖 vs 自绘） | 方案 | S2 | D-001 裁决（qrcode-generator） | **closed** | — | D-001 |
| I-004 | required | UX P1 范围确认（Toast 方案；搜索/筛选是否扩展协议能力） | 方案 | S4 | D-003 裁决（Toast 本地 UI；搜索复用既有 search-form 模式；select 筛选留 P2） | **closed** | — | D-003 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-15 | S1 范围与修复方向裁决（分批顺序；I-001 401→400；I-003 qrcode-generator） | accepted | `01-decision/D-001-s1-scope-and-fix-direction.md` |
| D-002 | 2026-08-15 | S2/S3 实施决策（MFA 分轨错误映射 + 解绑 UX；optionsSource 本地扩展 + RBAC 目录端点） | accepted | `01-decision/D-002-s2-s3-implementation.md` |
| D-003 | 2026-08-15 | I-004 闭合：Toast 本地 UI；搜索复用 search-form 模式；select 筛选留 P2 | accepted | `01-decision/D-003-s4-scope-confirmation.md` |

## 待决问题（P-004）

1. **UX 分批范围与顺序**：本波先做 M 类（P0）+ U-01/U-02（P0），U-03～U-07（P1）次批，U-08～U-14（P2）按资源选择——需用户确认分批与取舍。
2. **I-001 修复方向**：自服务端点无效码 401→400 的语义调整是否接受（影响 authFetch 行为边界）。
