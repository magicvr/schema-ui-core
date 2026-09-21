---
id: D-001-w34-availability-freeze
doc: decision
parent: GOAL-046-w34-batch-export-availability
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
status: accepted
---

# D-001 · W34 冻结：operator config 补齐 + 触发面可用性门禁

## 1. 诊断（`I-046-001`，已 verified）

用户报告（2026-09-19）：「无论是用户列表页还是角色列表页，选中列表项并点击『导出所选』的时候，都会报错宣示『未找到』，而并不会正常导出。」

**证据链（原样复现，非推断）**

| 步 | 动作 | 观察 |
|----|------|------|
| 1 | 以 `configs/config.yaml` 原样启动 API（临时 SQLite，避免触碰开发者库） | 启动成功 |
| 2 | `POST /api/jobs/batch-export`（admin Bearer，`Accept-Language: zh-CN`） | `404` · `{"error":"NOT_FOUND","message":"未找到","messageKey":"error.notFound"}` |
| 3 | `GET /api/jobs`（同一 token） | 同为 `404 NOT_FOUND` |
| 4 | `GET /api/export/users`（同步导出，对照） | 正常 —— 证明缺的只有 `admin.jobs` |
| 5 | 以 `profile: admin` 启动，同一请求 | `403 MUST_CHANGE_PASSWORD`（**路由存在**）；改密后 `202 queued → succeeded → result` 可下载 |
| 6 | 在 `configs/config.yaml` 补 `- admin.jobs` 后重启，同一请求 | users/roles 均 `202 → succeeded`，结果 2155 字节 CSV |

**根因**：`configs/config.yaml` 声明 `profile: custom` 并内联模块列表，其注释自称「the operator profile is a custom list = **the full admin preset** PLUS channel.telegram」。VP-038 把 `admin.jobs` 加入 admin preset（`kernel/profile.go`）时未同步这份列表 → 模块未装配 → 路由不存在 → 404。schema 节点 `users-batch-export` / `roles-batch-export` 属于 `admin.users` / `admin.roles`（所有 profile 都有），因此按钮照常渲染却必定失败。`messageKey: error.notFound` 正是 UI 显示的「未找到」文本来源（`errorcatalog` → `error.notFound` → zh「未找到」）。

**次要形态（同波次处理）**：`mvp` / `demo` preset 既无 `admin.jobs`（`jobs.write`）也无 `admin.data-transfer`（`data.export`），而节点同样发布 → 这两个 profile 下入口只有 404/403 一条路。实测 mvp：`POST /api/jobs/batch-export → 404`、`GET /api/export/users → 404`，而 `/api/schema/users` 仍含 `jobs-batch-export` 节点。

## 2. 授权与边界

- 载体 = `[workspace-010]` W34（VP-010 持续符合性程序波次）；**VP-038 保持 `closed`**，不重开、不改其 `status` 与 workspace-038 台账正文。
- 变更面：`apps/api/configs/config.yaml`（模块装配）、`apps/web/src/components/jobs-batch-export.tsx`（可用性渲染）、i18n 两个键、两处测试。
- **不做**：不改服务端两道门禁语义（`jobs.write` + `data.export` 仍是唯一授权方，本波只做客户端**镜像**以便 UI 不撒谎）；不改 Job 六态/批量导出契约；不把导出模块搬进 `mvp`/`demo`（是否进 preset 属 Profile 内容决策，须用户 P-004，本波不擅自决定）；不触碰 pinned 协议工件。

## 3. 不可用时的呈现：可见 + disabled + 说明（`I-046-002`，已 verified）

**决定**：当 `context.user.permissions` 不同时具备 `jobs.write` 与 `data.export` 时，触发面渲染为 **disabled 按钮 + 说明文案**（`data-jobs-batch-export-unavailable` / `-unavailable-note`），**不隐藏**。

**理由（两条独立、可核对）**

1. **W33 `D-001` §3 已冻结 fail-open 取向**：「本扩展是**布局**能力，不是权限或数据门禁。若采用『隐藏』，一处 schema 笔误就会让『导出所选』这个操作入口**静默消失**——这是比『位置不对』严重得多的失败模式；原地渲染则退化成 R3 的既有观感，用户仍能完成操作。」静默移除操作入口与本取向冲突。
2. **布局契约**：插槽宿主由表格在「有消费者注册」时渲染。若组件返回 `null`，宿主成为高度 0 的空 flex 子元素，破坏 `list-visual-surface` 的「page-actions 行控件同一高度」契约 —— 本波实施中已被该 e2e 实测捕获（`got 0, 32, 32, 32`），修正后才复绿。因此「隐藏」不是免费选项。

**同时**：上下文完全未提供 `permissions`（裸测试挂具、旧 host）时保持 **fail-open（视为可用）**，避免因读不到数据而禁用一个本来可用的入口。

**为何以这两个权限为判据**：服务端提交路由同时要求 `jobs.write`（授权变更异步运行时）与 `data.export`（授权数据外带），两者分别只由 `admin.jobs` 与 `admin.data-transfer` 贡献，因此「两权限齐备」与「路由存在且调用者被授权」在装配上等价；缺一即不可用（403/404）。这是**镜像**而非新增门禁语义。

## 4. 404 语义纠正

提交命中 `404` 时不再回显服务端 `message`（裸 `NOT_FOUND` / 「未找到」，会被读成「我的行不存在」——本次排查正是被它误导），改报 `schema.jobs.batchExport.unavailable`：「当前部署未启用批量导出（未装配 admin.jobs 模块）」。这是**诊断性兜底**：即使某部署仍发布节点而未装配模块，用户也能读到真实原因。

## 5. 守卫设计（`I-046-003`，已 verified）

原覆盖只到「admin profile 下可用」（`jobs-result-center.spec.ts` 为 admin-only），**没有任何断言覆盖「节点已发布但路由未挂载」**。本波补两层：

1. **配置层守卫**（`apps/api/internal/config/operator_config_coverage_test.go`）：operator config 内联列表必须是 **admin preset 的超集**（显式允许 `channel.telegram` / `biz.digital-offer` 等额外项，并显式要求 `admin.jobs`）。取「超集」而非「逐字相等」，是为了让守卫表达文件自述的不变量（镜像 admin preset + 显式额外项），而不是把某一次的快照冻死。
2. **组件层回归**（`jobs-batch-export.test.tsx`）：4 例可用性用例，覆盖门禁缺一、齐备、上下文未知、404 文案。

## 6. 授权范围（跨文件写入）

- 允许：`apps/api/configs/config.yaml`、`apps/api/internal/config/operator_config_coverage_test.go`、`apps/web/src/components/jobs-batch-export.tsx`、`apps/web/src/components/jobs-batch-export.test.tsx`、`apps/web/src/i18n/messages/{en-US,zh-CN}.json`。
- 愿景层登记：`docs/vision/roadmap.md`「未决项统一登记」+「最近更新」（只登记事实）。
- 禁止：VP-038 文件 `status`/判据、workspace-038 台账正文、`docs/schemas/**`、`apps/web/src/protocol/upstream/**`。
